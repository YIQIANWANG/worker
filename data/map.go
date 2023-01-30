package data

var Map *GlobalMap

type GlobalMap struct {
	MappingInfos []MappingInfo
}

type MappingInfo struct {
	ShardIDStart uint32
	ShardIDEnd   uint32
	GroupID      string
	OldGroupID   string
}
