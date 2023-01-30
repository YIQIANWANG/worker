package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"sync"
	"worker/app"
	"worker/conf"
	"worker/data"
)

func PING(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "PONG!"})
}

func RESET(c *gin.Context) {
	// 删除更早的日志
	_ = os.RemoveAll(conf.LogFilePath)
	// 重置所有Storage
	storageOperator := app.Default.GetStorageOperator()
	for groupID, _ := range data.Groups.GroupInfos {
		wg := sync.WaitGroup{}
		wg.Add(len(data.Groups.GroupInfos[groupID]))
		for _, storage := range data.Groups.GroupInfos[groupID] {
			go func(storageAddress string) {
				storageOperator.RESET(storageAddress)
				wg.Done()
			}(storage.StorageAddress)
		}
		wg.Wait()
	}

	log.Println("RESET")
	c.JSON(http.StatusOK, gin.H{"message": "Success."})
}

func GROUPS(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Success.", "groups": data.Groups})
}
