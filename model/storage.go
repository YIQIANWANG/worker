package model

type DefaultResponse struct {
	Message string `json:"message"`
}

type DelChunkResponse struct {
	Message string `json:"message"`
	Size    int    `json:"size"`
}
