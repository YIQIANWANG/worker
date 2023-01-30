package service

import (
	"encoding/base64"
	"hash/adler32"
	"math"
	"os"
	"strconv"
	"worker/data"
)

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil || !os.IsNotExist(err) {
		return true
	}
	return false
}

func getChunkID(userName, bucketName, fileName string, sequence int) string {
	str := userName + "/" + bucketName + "/" + fileName + "/" + strconv.Itoa(sequence)
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func getShardIDByChunkID(chunkID string) uint32 {
	shardID := adler32.Checksum([]byte(chunkID))
	if shardID == math.MaxUint32 {
		shardID = math.MaxUint32 - 1
	}
	return shardID
}

func getGroupIDsByChunkID(chunkID string) []string {
	shardID := getShardIDByChunkID(chunkID)
	for _, mappingInfo := range data.Map.MappingInfos {
		if mappingInfo.ShardIDStart <= shardID && shardID < mappingInfo.ShardIDEnd {
			return []string{mappingInfo.GroupID, mappingInfo.OldGroupID}
		}
	}
	return nil
}
