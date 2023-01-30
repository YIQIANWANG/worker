package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"worker/app"
	"worker/model"
)

func CreateBucket(c *gin.Context) {
	var form model.CreateBucketForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	// userName := c.PostForm("userName")
	// bucketName := c.PostForm("bucketName")

	bucketService := app.Default.GetBucketService()
	err := bucketService.CreateBucket(form.UserName, form.BucketName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Success."})
}
