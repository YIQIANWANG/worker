package operator

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"worker/model"
)

type StorageOperator struct {
}

func NewStorageOperator() *StorageOperator {
	return &StorageOperator{}
}

func (so *StorageOperator) PutChunk(storageAddress string, chunkID string, data []byte) error {
	reqUrl, _ := url.Parse("http://" + storageAddress + "/chunk")
	params := url.Values{}
	params.Set("chunkID", chunkID)
	reqUrl.RawQuery = params.Encode()

	resp, err := http.Post(reqUrl.String(), "text/plain", bytes.NewReader(data))
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	if resp.StatusCode == http.StatusBadRequest {
		body, _ := io.ReadAll(resp.Body)
		var storageResponse model.DefaultResponse
		_ = json.Unmarshal(body, &storageResponse)
		return errors.New(storageResponse.Message)
	}

	return nil
}

func (so *StorageOperator) GetChunk(storageAddress string, chunkID string) ([]byte, error) {
	reqUrl, _ := url.Parse("http://" + storageAddress + "/chunk")
	params := url.Values{}
	params.Set("chunkID", chunkID)
	reqUrl.RawQuery = params.Encode()

	resp, err := http.Get(reqUrl.String())
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	if resp.StatusCode == http.StatusBadRequest {
		body, _ := io.ReadAll(resp.Body)
		var storageResponse model.DefaultResponse
		_ = json.Unmarshal(body, &storageResponse)
		return nil, errors.New(storageResponse.Message)
	}

	body, _ := io.ReadAll(resp.Body)
	return body, nil
}

func (so *StorageOperator) DelChunk(storageAddress string, chunkID string) (int, error) {
	reqUrl, _ := url.Parse("http://" + storageAddress + "/chunk")
	params := url.Values{}
	params.Set("chunkID", chunkID)
	reqUrl.RawQuery = params.Encode()
	req, _ := http.NewRequest(http.MethodDelete, reqUrl.String(), nil)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return -1, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	body, _ := io.ReadAll(resp.Body)
	var storageResponse model.DelChunkResponse
	_ = json.Unmarshal(body, &storageResponse)
	if resp.StatusCode == http.StatusBadRequest {
		return -1, errors.New(storageResponse.Message)
	}

	return storageResponse.Size, nil
}

func (so *StorageOperator) GetChunkIDs(storageAddress string) ([]string, error) {
	reqUrl, _ := url.Parse("http://" + storageAddress + "/chunkIDs")

	resp, err := http.Get(reqUrl.String())
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	if resp.StatusCode == http.StatusBadRequest {
		body, _ := io.ReadAll(resp.Body)
		var storageResponse model.DefaultResponse
		_ = json.Unmarshal(body, &storageResponse)
		return nil, errors.New(storageResponse.Message)
	}

	body, _ := io.ReadAll(resp.Body)
	var chunkIDs []string
	_ = json.Unmarshal(body, &chunkIDs)
	return chunkIDs, nil
}

func (so *StorageOperator) PING(storageAddress string) {
	reqUrl, _ := url.Parse("http://" + storageAddress + "/PING")
	resp, err := http.Get(reqUrl.String())
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	return
}

func (so *StorageOperator) RESET(storageAddress string) {
	reqUrl, _ := url.Parse("http://" + storageAddress + "/RESET")
	resp, err := http.Get(reqUrl.String())
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	return
}
