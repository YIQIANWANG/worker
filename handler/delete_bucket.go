package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"worker/app"
)

func DeleteBucket(c *gin.Context) {
	userName := c.Query("userName")
	bucketName := c.Query("bucketName")

	bucketService := app.Default.GetBucketService()
	err := bucketService.DeleteBucket(userName, bucketName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Success."})
}
