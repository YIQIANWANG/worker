package handler

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"time"
	"worker/app"
	"worker/model"
	"worker/service"
)

func PutFile(c *gin.Context) {
	service.Count.WithLabelValues("PutFile").Inc()
	start := time.Now()

	// 获取表单
	var form model.PutFileForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	// userName := c.PostForm("userName")
	// bucketName := c.PostForm("bucketName")
	// fileName := c.PostForm("fileName")

	// 获取文件数据
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	defer func() {
		_ = file.Close()
	}()
	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	fileService := app.Default.GetFileService()
	err = fileService.PutFile(form.UserName, form.BucketName, form.FileName, fileData)

	service.Duration.WithLabelValues("PutFile").Observe(float64(time.Since(start) / time.Millisecond))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Success."})
}
