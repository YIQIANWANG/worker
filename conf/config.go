package conf

// 文件路径
const (
	LogFilePath = "logs"
)

// 服务配置
const (
	PORT              = "7777"
	PrometheusPort    = "9200" // 暴露给Prometheus的端口
	ChunkMaxSize      = 100000 // Chunk最大为100KB
	HeartbeatInternal = 10     // 心跳间隔为10s
)

// MongoDB
const (
	PROTOCOL       = "mongodb"
	USERNAME       = "mongouser"
	PASSWORD       = "YqMTE*5873QpUJ"
	ADDRESS        = "9.134.32.73:27017,9.134.38.231:27017,9.134.47.32:27017"
	AUTHENTICATION = "somedb?authSource=admin"
	DATABASE       = "localhost" // 本地测试
	// DATABASE       = "stress" // 压力测试
	// DATABASE       = "cos"    // 生产环境
)