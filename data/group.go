package data

var Groups *GlobalGroups

type GlobalGroups struct {
	GroupInfos    map[string][]Storage // 每个Group含有的Storage信息
	GroupsExist   map[string]bool      // 当前记录的GroupID
	StoragesExist map[string]bool      // 当前记录的StorageID
}

type Storage struct {
	StorageAddress string
	Latency        int64
}
