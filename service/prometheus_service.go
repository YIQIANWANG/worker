package service

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
	"worker/conf"
	"worker/data"
	"worker/operator"
)

var (
	// Count 请求计数
	Count = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "request_count",
			Help: "The number of requests",
		},
		[]string{"method"},
	)
	// Duration 请求延迟
	Duration = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "request_duration",
			Help: "The duration of requests",
		},
		[]string{"method"},
	)
	// AvailableCapacity 每个Group的可用容量
	AvailableCapacity = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "available_capacity",
			Help: "The available capacity of groups",
		},
		[]string{"group"},
	)
	// ActiveNodes 每个Group的活跃结点
	ActiveNodes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "active_nodes",
			Help: "The number of active storages in group",
		},
		[]string{"group"},
	)
)

type PrometheusService struct {
	mongoOperator *operator.MongoOperator
}

func NewPrometheusService(mongoOperator *operator.MongoOperator) *PrometheusService {
	return &PrometheusService{mongoOperator: mongoOperator}
}

func (ps *PrometheusService) InitMetrics() {
	prometheus.MustRegister(Count)
	prometheus.MustRegister(Duration)
	prometheus.MustRegister(AvailableCapacity)
	prometheus.MustRegister(ActiveNodes)
	groups, _ := ps.mongoOperator.GetGroups()
	for _, group := range groups {
		AvailableCapacity.WithLabelValues(group.GroupID).Set(float64(group.AvailableCap))
	}
	for groupID, _ := range data.Groups.GroupInfos {
		ActiveNodes.WithLabelValues(groupID).Set(float64(len(data.Groups.GroupInfos[groupID])))
	}
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		err := http.ListenAndServe(":"+conf.PrometheusPort, nil)
		if err != nil {
			panic(err)
		}
	}()
}

func (ps *PrometheusService) StartReport() {
	go func() {
		for true {
			groups, _ := ps.mongoOperator.GetGroups()
			for _, group := range groups {
				AvailableCapacity.WithLabelValues(group.GroupID).Set(float64(group.AvailableCap))
			}
			for groupID, _ := range data.Groups.GroupInfos {
				ActiveNodes.WithLabelValues(groupID).Set(float64(len(data.Groups.GroupInfos[groupID])))
			}
			time.Sleep(conf.HeartbeatInternal * time.Second)
		}
	}()
}
