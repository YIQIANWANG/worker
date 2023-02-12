package service

import (
	"log"
	"math"
	"sync"
	"worker/conf"
	"worker/data"
	"worker/model"
	"worker/operator"
)

type FileService struct {
	mongoOperator   *operator.MongoOperator
	storageOperator *operator.StorageOperator
}

func NewFileService(mongoOperator *operator.MongoOperator, storageOperator *operator.StorageOperator) *FileService {
	return &FileService{mongoOperator: mongoOperator, storageOperator: storageOperator}
}

func (fs *FileService) PutFile(userName, bucketName, fileName string, fileData []byte) error {
	// 文件拆分
	length := len(fileData)
	splitNumber := int(math.Ceil(float64(length) / conf.ChunkMaxSize))
	chunksData := make([][]byte, splitNumber)
	index := 0
	for sequence := 0; sequence < splitNumber; sequence++ {
		if index+conf.ChunkMaxSize < length {
			chunksData[sequence] = fileData[index : index+conf.ChunkMaxSize]
			index = index + conf.ChunkMaxSize
		} else {
			chunksData[sequence] = fileData[index:]
		}
	}

	// 并发处理Chunk
	wg := sync.WaitGroup{}
	wg.Add(splitNumber)
	for sequence := 0; sequence < splitNumber; sequence++ {
		go func(sequence int) {
			_ = fs.putChunk(userName, bucketName, fileName, chunksData[sequence], sequence, splitNumber)
			wg.Done()
		}(sequence)
	}
	wg.Wait()

	// 更新业务数据
	err := fs.mongoOperator.UpdateUserCount(userName, bucketName, 1)
	if err != nil {
		log.Println("Put File Failed: ", err)
		return err
	}
	return nil
}

func (fs *FileService) GetFile(userName, bucketName, fileName string) ([]byte, error) {
	// 读其中第一块Chunk
	firstChunkID := getChunkID(userName, bucketName, fileName, 0)
	firstChunk, err := fs.mongoOperator.GetChunk(firstChunkID)
	if err != nil {
		log.Println("Get File Failed: ", err)
		return nil, err
	}

	// 并发处理Chunk
	splitNumber := firstChunk.SplitNumber
	chunksData := make([][]byte, splitNumber)
	wg := sync.WaitGroup{}
	wg.Add(splitNumber)
	for sequence := 0; sequence < splitNumber; sequence++ {
		go func(sequence int) {
			chunkID := getChunkID(userName, bucketName, fileName, sequence)
			chunksData[sequence], _ = fs.getChunk(chunkID)
			wg.Done()
		}(sequence)
	}
	wg.Wait()
	fileData := make([]byte, 0)
	for _, chunkData := range chunksData {
		fileData = append(fileData, chunkData...)
	}

	// 更新业务数据
	err = fs.mongoOperator.UpdateUserCount(firstChunk.UserName, firstChunk.BucketName, 0)
	if err != nil {
		log.Println("Get File Failed: ", err)
		return nil, err
	}
	return fileData, nil
}

func (fs *FileService) DelFile(userName, bucketName, fileName string) error {
	// 读其中第一块Chunk
	firstChunkID := getChunkID(userName, bucketName, fileName, 0)
	firstChunk, err := fs.mongoOperator.GetChunk(firstChunkID)
	if err != nil {
		log.Println("Del File Failed: ", err)
		return err
	}

	// 并发处理Chunk
	splitNumber := firstChunk.SplitNumber
	wg := sync.WaitGroup{}
	wg.Add(splitNumber)
	for sequence := 0; sequence < splitNumber; sequence++ {
		go func(sequence int) {
			chunkID := getChunkID(userName, bucketName, fileName, sequence)
			_ = fs.delChunk(chunkID)
			wg.Done()
		}(sequence)
	}
	wg.Wait()

	// 更新业务数据
	err = fs.mongoOperator.UpdateUserCount(firstChunk.UserName, firstChunk.BucketName, -1)
	if err != nil {
		log.Println("Del File Failed: ", err)
		return err
	}
	return nil
}

func (fs *FileService) putChunk(userName, bucketName, fileName string, chunkData []byte, sequence, splitNumber int) error {
	// 查询GroupID
	chunkID := getChunkID(userName, bucketName, fileName, sequence)
	groupID := getGroupIDsByChunkID(chunkID)[0]

	// 选择Storages
	var storagesAddress []string
	for _, storage := range data.Groups.GroupInfos[groupID] {
		storagesAddress = append(storagesAddress, storage.StorageAddress)
	}

	// 请求Storages
	wg := sync.WaitGroup{}
	wg.Add(len(storagesAddress))
	for _, address := range storagesAddress {
		go func(address string) {
			_ = fs.storageOperator.PutChunk(address, chunkID, chunkData)
			wg.Done()
		}(address)
	}
	wg.Wait()

	// 更新业务数据
	size := len(chunkData)
	wg = sync.WaitGroup{}
	wg.Add(2)
	go func() {
		err := fs.mongoOperator.UpdateGroupAvailableCap(groupID, -size)
		if err != nil {
			log.Println("Put Chunk Failed: ", err)
		}
		wg.Done()
	}()
	go func() {
		chunk := &model.Chunk{
			ChunkID:     chunkID,
			UserName:    userName,
			BucketName:  bucketName,
			FileName:    fileName,
			Size:        size,
			Sequence:    sequence,
			SplitNumber: splitNumber,
		}
		err := fs.mongoOperator.InsertChunk(chunk)
		if err != nil {
			log.Println("Put Chunk Failed: ", err)
		}
		wg.Done()
	}()
	wg.Wait()
	return nil
}

func (fs *FileService) getChunk(chunkID string) ([]byte, error) {
	// 查询GroupID
	groupIDs := getGroupIDsByChunkID(chunkID)

	// 选择Storages
	storagesAddress := make([][]string, 2)
	for _, storage := range data.Groups.GroupInfos[groupIDs[0]] {
		storagesAddress[0] = append(storagesAddress[0], storage.StorageAddress)
	}
	for _, storage := range data.Groups.GroupInfos[groupIDs[1]] {
		storagesAddress[1] = append(storagesAddress[1], storage.StorageAddress)
	}

	// 请求Storage
	failedAddress := make([]string, 0)
	successAddress := ""
	var chunkData []byte
	var err error
	for _, address := range storagesAddress[0] {
		chunkData, err = fs.storageOperator.GetChunk(address, chunkID)
		if err != nil {
			failedAddress = append(failedAddress, address)
		} else {
			successAddress = address
			break
		}
	}
	if successAddress == "" {
		for _, address := range storagesAddress[1] {
			chunkData, err = fs.storageOperator.GetChunk(address, chunkID)
			if err == nil {
				successAddress = address
				break
			}
		}
	}
	// 同步读取失败的Storage需要执行写操作
	if n := len(failedAddress); n > 0 {
		wg := sync.WaitGroup{}
		wg.Add(n)
		for _, address := range failedAddress {
			go func(address string) {
				_ = fs.storageOperator.PutChunk(address, chunkID, chunkData)
				wg.Done()
			}(address)
		}
		wg.Wait()
	}
	return chunkData, nil
}

func (fs *FileService) delChunk(chunkID string) error {
	// 查询GroupID
	groupIDs := getGroupIDsByChunkID(chunkID)

	// 选择Storages
	var storagesAddress [2][]string
	for _, storage := range data.Groups.GroupInfos[groupIDs[0]] {
		storagesAddress[0] = append(storagesAddress[0], storage.StorageAddress)
	}
	for _, storage := range data.Groups.GroupInfos[groupIDs[1]] {
		storagesAddress[1] = append(storagesAddress[1], storage.StorageAddress)
	}

	// 请求Storages
	var size [2]int
	wg := [2]sync.WaitGroup{}
	wg[0].Add(len(storagesAddress[0]))
	wg[1].Add(len(storagesAddress[1]))
	for _, address := range storagesAddress[0] {
		go func(address string) {
			size[0], _ = fs.storageOperator.DelChunk(address, chunkID)
			wg[0].Done()
		}(address)
	}
	for _, address := range storagesAddress[1] {
		go func(address string) {
			size[1], _ = fs.storageOperator.DelChunk(address, chunkID)
			wg[1].Done()
		}(address)
	}
	wg[0].Done()
	wg[1].Done()

	// 更新业务数据
	wg2 := sync.WaitGroup{}
	wg2.Add(1)
	go func() {
		err := fs.mongoOperator.DeleteChunk(chunkID)
		if err != nil {
			log.Println("Del Chunk Failed: ", err)
		}
		wg2.Done()
	}()
	if size[0] > 0 {
		wg2.Add(1)
		go func() {
			err := fs.mongoOperator.UpdateGroupAvailableCap(groupIDs[0], size[0])
			if err != nil {
				log.Println("Del Chunk Failed: ", err)
			}
			wg2.Done()
		}()
	}
	if size[1] > 0 {
		wg2.Add(1)
		go func() {
			err := fs.mongoOperator.UpdateGroupAvailableCap(groupIDs[1], size[1])
			if err != nil {
				log.Println("Del Chunk Failed: ", err)
			}
			wg2.Done()
		}()
	}
	wg2.Wait()

	return nil
}
