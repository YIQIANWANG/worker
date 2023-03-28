package data

var Map MappingInfos

type MappingInfos []*MappingInfo

type MappingInfo struct {
	Start      uint32
	End        uint32
	GroupID    string
	OldGroupID string
}
