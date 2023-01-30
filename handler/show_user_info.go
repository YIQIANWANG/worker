package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"worker/app"
)

func ShowUserInfo(c *gin.Context) {
	userName := c.Query("userName")
	bucketName := c.Query("bucketName")

	userService := app.Default.GetUserService()
	if bucketName == "" {
		user, err := userService.GetUserStatistics(userName)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Success.", "userStatistics": user})
	} else {
		fileNames, err := userService.GetUserFiles(userName, bucketName)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Success.", "files": fileNames})
	}
}
