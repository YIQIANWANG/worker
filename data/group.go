package data

var Groups *GlobalGroups

type GlobalGroups struct {
	GroupInfos map[string]*GroupInfo
	GroupSet   map[string]bool // 当前活跃的Group集合
	StorageSet map[string]bool // 当前活跃的Storage集合
}

type GroupInfo struct {
	Storages     []Storage
	Capacity     int
	AvailableCap int
}

type Storage struct {
	StorageAddress string
	Latency        int64
}
