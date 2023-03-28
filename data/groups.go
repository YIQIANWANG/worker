package data

var Groups GroupInfos

type GroupInfos map[string]*GroupInfo

type GroupInfo struct {
	Storages     []*Storage
	Capacity     int
	AvailableCap int
}

type Storage struct {
	StorageAddress string
	Latency        int64
}
