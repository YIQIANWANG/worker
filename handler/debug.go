package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"worker/data"
)

func PING(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "PONG!"})
}

func GROUPS(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Success.", "groups": data.Groups})
}
