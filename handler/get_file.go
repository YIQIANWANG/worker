package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"worker/app"
	"worker/service"
)

func GetFile(c *gin.Context) {
	service.Count.WithLabelValues("GetFile").Inc()
	start := time.Now()

	userName := c.Query("userName")
	bucketName := c.Query("bucketName")
	fileName := c.Query("fileName")

	fileService := app.Default.GetFileService()
	fileData, err := fileService.GetFile(userName, bucketName, fileName)

	service.Duration.WithLabelValues("GetFile").Observe(float64(time.Since(start) / time.Millisecond))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	c.Data(http.StatusOK, "text/plain", fileData)
}
