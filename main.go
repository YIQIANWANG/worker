package main

import (
	"github.com/gin-gonic/gin"
	"worker/conf"
	"worker/handler"
)

func main() {
	router := gin.Default()
	router.POST("/bucket", handler.CreateBucket)
	router.DELETE("/bucket", handler.DeleteBucket)
	router.GET("/user", handler.ShowUserInfo)
	router.POST("/file", handler.PutFile)
	router.GET("/file", handler.GetFile)
	router.DELETE("/file", handler.DelFile)
	router.GET("/PING", handler.PING)
	router.GET("/RESET", handler.RESET)
	router.GET("/GROUPS", handler.GROUPS)

	err := router.Run(":" + conf.PORT)
	if err != nil {
		panic(err)
	}
}
