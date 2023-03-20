package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"worker/app"
)

func ShowUser(c *gin.Context) {
	userName := c.Query("userName")

	userService := app.Default.GetUserService()
	user, err := userService.ShowUser(userName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Success.", "userStatistics": user})
}
