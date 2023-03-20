package main

import (
	"github.com/gin-gonic/gin"
	"worker/conf"
	"worker/handler"
)

func main() {
	router := gin.Default()
	router.POST("/bucket", handler.CreateBucket)
	router.GET("/bucket", handler.ShowBucket)
	router.DELETE("/bucket", handler.DeleteBucket)
	router.POST("/file", handler.PutFile)
	router.GET("/file", handler.GetFile)
	router.DELETE("/file", handler.DelFile)
	router.GET("/user", handler.ShowUser)
	router.GET("/PING", handler.PING)
	router.GET("/GROUPS", handler.GROUPS)

	err := router.Run(":" + conf.PORT)
	if err != nil {
		panic(err)
	}
}
