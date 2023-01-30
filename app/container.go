package app

import (
	"worker/dal"
	"worker/operator"
	"worker/service"
)

var Default *Container

func InitDefault() {
	var err error
	Default, err = InitContainer()
	if err != nil {
		panic(err)
	}
}

type Container struct {
	mongoClient       *dal.MongoClient
	mongoOperator     *operator.MongoOperator
	storageOperator   *operator.StorageOperator
	bucketService     *service.BucketService
	fileService       *service.FileService
	userService       *service.UserService
	heartbeatService  *service.HeartbeatService
	prometheusService *service.PrometheusService
}

func NewContainer(
	mongoClient *dal.MongoClient,
	mongoOperator *operator.MongoOperator,
	storageOperator *operator.StorageOperator,
	bucketService *service.BucketService,
	fileService *service.FileService,
	userService *service.UserService,
	heartbeatService *service.HeartbeatService,
	prometheusService *service.PrometheusService,
) *Container {
	return &Container{
		mongoClient:       mongoClient,
		mongoOperator:     mongoOperator,
		storageOperator:   storageOperator,
		bucketService:     bucketService,
		fileService:       fileService,
		userService:       userService,
		heartbeatService:  heartbeatService,
		prometheusService: prometheusService,
	}
}

func (c *Container) GetMongoClient() *dal.MongoClient {
	return c.mongoClient
}

func (c *Container) GetMongoOperator() *operator.MongoOperator {
	return c.mongoOperator
}

func (c *Container) GetStorageOperator() *operator.StorageOperator {
	return c.storageOperator
}

func (c *Container) GetBucketService() *service.BucketService {
	return c.bucketService
}

func (c *Container) GetFileService() *service.FileService {
	return c.fileService
}

func (c *Container) GetUserService() *service.UserService {
	return c.userService
}

func (c *Container) GetHeartbeatService() *service.HeartbeatService {
	return c.heartbeatService
}

func (c *Container) GetPrometheusService() *service.PrometheusService {
	return c.prometheusService
}
