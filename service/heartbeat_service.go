package service

import (
	"log"
	"math/rand"
	"sync"
	"time"
	"worker/conf"
	"worker/data"
	"worker/operator"
)

type HeartbeatService struct {
	mongoOperator   *operator.MongoOperator
	storageOperator *operator.StorageOperator
}

func NewHeartbeatService(mongoOperator *operator.MongoOperator, storageOperator *operator.StorageOperator) *HeartbeatService {
	return &HeartbeatService{mongoOperator: mongoOperator, storageOperator: storageOperator}
}

func (hs *HeartbeatService) InitData() {
	var err error
	data.Map, err = hs.getGlobalMap()
	if err != nil {
		panic(err)
	}
	data.Groups, err = hs.getGlobalGroups()
	if err != nil {
		panic(err)
	}
}

func (hs *HeartbeatService) StartCheck() {
	go func() {
		for true {
			go func() {
				newMap, err := hs.getGlobalMap()
				if err != nil {
					log.Println("GetGlobalMap Failed: ", err)
					return
				}
				data.Map = newMap
			}()
			go func() {
				newGroups, err := hs.getGlobalGroups()
				if err != nil {
					log.Println("GetGlobalGroups Failed: ", err)
					return
				}
				hs.check(newGroups)
			}()
			time.Sleep(conf.HeartbeatInternal * time.Second)
		}
	}()
}

func (hs *HeartbeatService) getGlobalMap() (*data.GlobalMap, error) {
	rawInfos, err := hs.mongoOperator.GetMappingInfos()
	if err != nil {
		return nil, err
	}
	mappingInfos := make([]data.MappingInfo, 0)
	for _, rawInfo := range rawInfos {
		mappingInfos = append(mappingInfos, data.MappingInfo{
			Start:      rawInfo.Start,
			End:        rawInfo.End,
			GroupID:    rawInfo.GroupID,
			OldGroupID: rawInfo.OldGroupID,
		})
	}
	return &data.GlobalMap{MappingInfos: mappingInfos}, nil
}

func (hs *HeartbeatService) getGlobalGroups() (*data.GlobalGroups, error) {
	rawInfos, err := hs.mongoOperator.GetStorages()
	if err != nil {
		return nil, err
	}
	groupInfos := make(map[string]*data.GroupInfo)
	groupSet := make(map[string]bool)
	storageSet := make(map[string]bool)

	// 统计每个Group的可用Storage
	now := time.Now().Unix()
	for _, rawInfo := range rawInfos {
		if now-rawInfo.UpdateTime <= 1.5*conf.HeartbeatInternal {
			if groupInfos[rawInfo.GroupID] == nil {
				groupInfos[rawInfo.GroupID] = &data.GroupInfo{}
			}
			groupInfos[rawInfo.GroupID].Storages = append(groupInfos[rawInfo.GroupID].Storages, data.Storage{StorageAddress: rawInfo.StorageAddress})
			groupInfos[rawInfo.GroupID].Capacity = rawInfo.Capacity
			if groupInfos[rawInfo.GroupID].AvailableCap > rawInfo.AvailableCap {
				groupInfos[rawInfo.GroupID].AvailableCap = rawInfo.AvailableCap
			}
			groupSet[rawInfo.GroupID] = true
			storageSet[rawInfo.StorageAddress] = true
		} else {
			_ = hs.mongoOperator.DeleteStorage(rawInfo.StorageAddress)
		}
	}

	/*
		// 对每个Group并发PING Storage，按时延排序
		for groupID, _ := range groupInfos {
			var lock sync.Mutex
			wg := sync.WaitGroup{}
			wg.Add(len(groupInfos[groupID].Storages))
			for i := range groupInfos[groupID].Storages {
				go func(i int) {
					start := time.Now()
					hs.storageOperator.PING(groupInfos[groupID].Storages[i].StorageAddress)
					lock.Lock()
					groupInfos[groupID].Storages[i].Latency = time.Since(start).Milliseconds()
					lock.Unlock()
					wg.Done()
				}(i)
			}
			wg.Wait()
			sort.Slice(groupInfos[groupID], func(i, j int) bool {
				return groupInfos[groupID].Storages[i].Latency < groupInfos[groupID].Storages[j].Latency
			})
		}
	*/

	return &data.GlobalGroups{
		GroupInfos: groupInfos,
		GroupSet:   groupSet,
		StorageSet: storageSet,
	}, nil
}

// 新Group进行数据迁移，对旧Group新Storage进行数据同步
func (hs *HeartbeatService) check(newGroups *data.GlobalGroups) {
	for groupID := range newGroups.GroupInfos {
		// 新Group出现
		if !data.Groups.GroupSet[groupID] {
			log.Println("Find New Group, GroupID: ", groupID)
			// 在映射信息中查找新加入Group对应的OldGroupID
			found := false
			var oldGroupID string
			for _, mappingInfo := range data.Map.MappingInfos {
				if mappingInfo.GroupID == groupID {
					found = true
					oldGroupID = mappingInfo.OldGroupID
					break
				}
			}
			// 若Map还未更新，Groups需要在下个心跳更新
			if !found {
				delete(newGroups.GroupSet, groupID)
				return
			}
			// 系统首次启动，不需要转移数据
			if found && oldGroupID == "" {
				continue
			}
			dst := make([]string, 0)
			src := make([]string, 0)
			for _, storage := range newGroups.GroupInfos[groupID].Storages {
				dst = append(dst, storage.StorageAddress)
			}
			for _, storage := range newGroups.GroupInfos[oldGroupID].Storages {
				src = append(src, storage.StorageAddress)
			}
			// 获取需要转移的ChunkID
			chunkIDs, err := hs.storageOperator.GetChunkIDs(src[rand.Intn(len(src))])
			if err != nil {
				log.Println("TransChunks Failed: ", err)
				return
			}
			toTransChunkIDs := make([]string, 0)
			for _, chunkID := range chunkIDs {
				if getGroupIDsByChunkID(chunkID)[1] == oldGroupID {
					toTransChunkIDs = append(toTransChunkIDs, chunkID)
				}
			}
			if len(toTransChunkIDs) > 0 {
				// 数据迁移
				hs.transChunks(dst, src, toTransChunkIDs)
			}
			// MappingInfos更新
			for _, mappingInfo := range data.Map.MappingInfos {
				if mappingInfo.GroupID == groupID {
					err := hs.mongoOperator.UpdateMappingInfoDeleteOldGroupID(mappingInfo.Start, mappingInfo.End)
					if err != nil {
						log.Println("Trans Chunks failed: ", err)
					}
					break
				}
			}
		} else {
			dst := make([]string, 0)
			src := make([]string, 0)
			for _, storage := range newGroups.GroupInfos[groupID].Storages {
				// 旧Group新Storage出现
				if !data.Groups.StorageSet[storage.StorageAddress] {
					log.Printf("Find New Storage in %s, StorageAddress: %s\n", groupID, storage.StorageAddress)
					dst = append(dst, storage.StorageAddress)
				} else {
					src = append(src, storage.StorageAddress)
				}
			}
			if len(dst) > 0 {
				// 数据同步
				hs.syncChunks(dst, src)
			}
		}
	}
	data.Groups = newGroups
}

func (hs *HeartbeatService) transChunks(dst, src []string, toTransChunkIDs []string) {
	failed := make(chan string, len(toTransChunkIDs))
	wg := sync.WaitGroup{}
	wg.Add(len(toTransChunkIDs))
	for _, toTransChunkID := range toTransChunkIDs {
		go func(toTransChunkID string) {
			defer wg.Done()
			done := false
			for i := 0; i < 3 && !done; i++ {
				// 获取数据
				chunkData, err := hs.storageOperator.GetChunk(src[rand.Intn(len(src))], toTransChunkID)
				if err != nil {
					continue
				}
				// 并发写入数据
				written := make(chan string, len(dst))
				wgWrite := sync.WaitGroup{}
				wgWrite.Add(len(dst))
				for index := 0; index < len(dst); index++ {
					go func(index int) {
						defer wgWrite.Done()
						err := hs.storageOperator.PutChunk(dst[index], toTransChunkID, chunkData)
						if err != nil {
							return
						}
						written <- toTransChunkID
					}(index)
				}
				wgWrite.Wait()
				if len(written) < len(dst) {
					continue
				}
				// 并发删除数据
				deleted := make(chan string, len(src))
				wgDelete := sync.WaitGroup{}
				wgDelete.Add(len(src))
				for index := 0; index < len(src); index++ {
					go func(index int) {
						defer wgDelete.Done()
						err := hs.storageOperator.DelChunk(src[index], toTransChunkID)
						if err != nil {
							return
						}
						deleted <- toTransChunkID
					}(index)
				}
				wgDelete.Wait()
				if len(deleted) < len(src) {
					continue
				}
				done = true
			}
			if !done {
				failed <- toTransChunkID
			}
		}(toTransChunkID)
	}
	wg.Wait()
	if len(failed) > 0 {
		str := ""
		for len(failed) > 0 {
			str += <-failed + "  "
		}
		log.Println("TransChunks Failed, ChunkID: ", str)
	}
}

func (hs *HeartbeatService) syncChunks(dst, src []string) {
	// 获取需要同步的ChunkID
	chunkIDs, err := hs.storageOperator.GetChunkIDs(src[rand.Intn(len(src))])
	if err != nil {
		log.Println("SyncChunks Failed: ", err)
		return
	}
	// 并发同步
	failed := make(chan string, len(dst))
	wg := sync.WaitGroup{}
	wg.Add(len(dst))
	for index := 0; index < len(dst); index++ {
		go func(index int) {
			defer wg.Done()
			done := make(map[string]bool)
			for i := 0; i < 3 && len(done) < len(chunkIDs); i++ {
				for _, chunkID := range chunkIDs {
					if !done[chunkID] {
						chunkData, err := hs.storageOperator.GetChunk(src[rand.Intn(len(src))], chunkID)
						if err != nil {
							continue
						}
						err = hs.storageOperator.PutChunk(dst[index], chunkID, chunkData)
						if err != nil {
							continue
						}
						done[chunkID] = true
					}
				}
			}
			if len(done) < len(chunkIDs) {
				failed <- dst[index]
			}
		}(index)
	}
	wg.Wait()
	if len(failed) > 0 {
		str := ""
		for len(failed) > 0 {
			str += <-failed + "  "
		}
		log.Println("SyncChunks Failed, StorageAddress: ", str)
	}
}
