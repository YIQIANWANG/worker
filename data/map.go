package data

var Map *GlobalMap

type GlobalMap struct {
	MappingInfos []MappingInfo
}

type MappingInfo struct {
	Start      uint32
	End        uint32
	GroupID    string
	OldGroupID string
}
