package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
	"worker/app"
	"worker/conf"
	"worker/service"
)

func init() {
	InitLog()
	app.InitDefault()
	InitHeartBeat()
	InitPrometheus()
}

func InitLog() {
	if service.PathExists(conf.LogFilePath) {
		_ = os.RemoveAll(conf.LogFilePath)
	}
	err := os.MkdirAll(conf.LogFilePath, os.ModePerm)
	if err != nil {
		panic(err)
	}

	// Log
	logFileName := fmt.Sprintf("%s/%s-%s-%s", conf.LogFilePath, time.Now().Format("2006"), time.Now().Format("01"), time.Now().Format("02"))
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm) // 创建、追加、读写，777（所有权限）
	if err != nil {
		panic(err)
	}
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func InitHeartBeat() {
	heartbeatService := app.Default.GetHeartbeatService()
	heartbeatService.InitData()
	heartbeatService.StartCheck()
}

func InitPrometheus() {
	prometheusService := app.Default.GetPrometheusService()
	prometheusService.InitMetrics()
	prometheusService.StartReport()
}
