package model

// Storage 存储节点信息
type Storage struct {
	StorageAddress string `bson:"storageAddress" json:"storageAddress"`
	GroupID        string `bson:"groupID" json:"groupID"`
	Capacity       int    `bson:"capacity" json:"capacity"`
	AvailableCap   int    `bson:"availableCap" json:"availableCap"`
	UpdateTime     int64  `bson:"updateTime" json:"updateTime"`
}

// MappingInfo 映射信息
type MappingInfo struct {
	Start      uint32 `bson:"start" json:"start"`
	End        uint32 `bson:"end" json:"end"`
	GroupID    string `bson:"groupID" json:"groupID"`
	OldGroupID string `bson:"oldGroupID" json:"oldGroupID"`
}

// User 用户信息
type User struct {
	UserName       string   `bson:"userName" json:"userName"`
	Secret         string   `bson:"secret" json:"secret"`
	Buckets        []Bucket `bson:"buckets" json:"buckets"`
	FileCount      int      `bson:"fileCount" json:"fileCount"`
	OperationCount int      `bson:"operationCount" json:"operationCount"`
}

// Bucket 存储桶信息
type Bucket struct {
	BucketName     string `bson:"bucketName" json:"bucketName"`
	FileCount      int    `bson:"fileCount" json:"fileCount"`
	OperationCount int    `bson:"operationCount" json:"operationCount"`
}

// Chunk 文件块信息
type Chunk struct {
	ChunkID     string `bson:"chunkID" json:"chunkID"`
	UserName    string `bson:"userName" json:"userName"`
	BucketName  string `bson:"bucketName" json:"bucketName"`
	FileName    string `bson:"fileName" json:"fileName"`
	Size        int    `bson:"size" json:"size"`
	Sequence    int    `bson:"sequence" json:"sequence"`
	SplitNumber int    `bson:"splitNumber" json:"splitNumber"`
}
