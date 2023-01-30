package model

type PutFileForm struct {
	UserName   string `json:"userName" form:"userName"`
	BucketName string `json:"bucketName" form:"bucketName"`
	FileName   string `json:"fileName" form:"fileName"`
}

type CreateBucketForm struct {
	UserName   string `json:"userName" form:"userName"`
	BucketName string `json:"bucketName" form:"bucketName"`
}
