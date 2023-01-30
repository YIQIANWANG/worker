package service

import (
	"log"
	"worker/model"
	"worker/operator"
)

type BucketService struct {
	mongoOperator *operator.MongoOperator
	fileService   *FileService
}

func NewBucketService(mongoOperator *operator.MongoOperator, fileService *FileService) *BucketService {
	return &BucketService{mongoOperator: mongoOperator, fileService: fileService}
}

func (bs *BucketService) CreateBucket(userName, bucketName string) error {
	bucket := &model.Bucket{
		BucketName: bucketName,
	}
	err := bs.mongoOperator.UpdateUserAppendBucket(userName, bucket)
	if err != nil {
		log.Println("Update User Bucket Failed: ", err)
		return err
	}
	return nil
}

func (bs *BucketService) DeleteBucket(userName, bucketName string) error {
	// 获取待删除文件
	fileNames, err := bs.mongoOperator.GetFileNamesByUserNameAndBucketName(userName, bucketName)
	if err != nil {
		log.Println("Get FileNames Failed: ", err)
		return err
	}

	// 文件依次删除
	for _, fileName := range fileNames {
		err = bs.fileService.DelFile(userName, bucketName, fileName)
		if err != nil {
			log.Println("Delete File Failed: ", err)
		}
	}

	// 更新桶数据
	err = bs.mongoOperator.UpdateUserDeleteBucket(userName, bucketName)
	if err != nil {
		log.Println("Update User Bucket Failed: ", err)
		return err
	}
	return nil
}
