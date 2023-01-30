package model

// Storage 存储实例信息
type Storage struct {
	StorageAddress string `bson:"storageAddress" json:"storageAddress"`
	GroupID        string `bson:"groupID" json:"groupID"`
	UpdateTime     int64  `bson:"updateTime" json:"updateTime"`
}

// MappingInfo 映射信息
type MappingInfo struct {
	ShardIDStart uint32 `bson:"shardIDStart" json:"shardIDStart"`
	ShardIDEnd   uint32 `bson:"shardIDEnd" json:"shardIDEnd"`
	GroupID      string `bson:"groupID" json:"groupID"`
	OldGroupID   string `bson:"oldGroupID" json:"oldGroupID"`
}

// Group 存储实例组信息
type Group struct {
	GroupID      string `bson:"groupID" json:"groupID"`
	Capacity     int    `bson:"capacity" json:"capacity"`
	AvailableCap int    `bson:"availableCap" json:"availableCap"`
}

// User 用户信息
type User struct {
	UserName       string   `bson:"userName" json:"userName"`
	Secret         string   `bson:"secret" json:"secret"`
	Buckets        []Bucket `bson:"buckets" json:"buckets"`
	FileCount      int      `bson:"fileCount" json:"fileCount"`
	OperationCount int      `bson:"operationCount" json:"operationCount"`
}

// Bucket 桶信息
type Bucket struct {
	BucketName     string `bson:"bucketName" json:"bucketName"`
	FileCount      int    `bson:"fileCount" json:"fileCount"`
	OperationCount int    `bson:"operationCount" json:"operationCount"`
}

// Chunk 块信息
type Chunk struct {
	ChunkID     string `bson:"chunkID" json:"chunkID"`
	UserName    string `bson:"userName" json:"userName"`
	BucketName  string `bson:"bucketName" json:"bucketName"`
	FileName    string `bson:"fileName" json:"fileName"`
	Size        int    `bson:"size" json:"size"`
	Sequence    int    `bson:"sequence" json:"sequence"`
	SplitNumber int    `bson:"splitNumber" json:"splitNumber"`
}
