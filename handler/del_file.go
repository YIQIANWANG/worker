package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"worker/app"
	"worker/service"
)

func DelFile(c *gin.Context) {
	service.Count.WithLabelValues("DelFile").Inc()
	start := time.Now()

	userName := c.Query("userName")
	bucketName := c.Query("bucketName")
	fileName := c.Query("fileName")

	fileService := app.Default.GetFileService()
	err := fileService.DelFile(userName, bucketName, fileName)

	service.Duration.WithLabelValues("DelFile").Observe(float64(time.Since(start) / time.Millisecond))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Success."})
}
