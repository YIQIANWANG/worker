package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
	"worker/app"
	"worker/conf"
	"worker/data"
	"worker/service"
)

func init() {
	LogInit()
	app.InitDefault()
	DataInit()
	PrometheusInit()
}

func LogInit() {
	if !service.PathExists(conf.LogFilePath) {
		err := os.MkdirAll(conf.LogFilePath, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	logFile := fmt.Sprintf("%s/%s-%s-%s", conf.LogFilePath, time.Now().Format("2006"), time.Now().Format("01"), time.Now().Format("02"))
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm) // 创建、追加、读写，777（所有权限）
	if err != nil {
		panic(err)
	}
	multiWriter := io.MultiWriter(os.Stdout, file)
	log.SetOutput(multiWriter)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func DataInit() {
	heartbeatService := app.Default.GetHeartbeatService()
	var err error
	data.Map, err = heartbeatService.GetGlobalMap()
	if err != nil {
		panic(err)
	}
	data.Groups, err = heartbeatService.GetGlobalGroups()
	if err != nil {
		panic(err)
	}
	heartbeatService.StartCheck()
}

func PrometheusInit() {
	prometheusService := app.Default.GetPrometheusService()
	prometheusService.InitMetrics()
	prometheusService.StartReport()
}
