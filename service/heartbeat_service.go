package service

import (
	"log"
	"sort"
	"sync"
	"time"
	"worker/conf"
	"worker/data"
	"worker/operator"
)

type HeartbeatService struct {
	mongoOperator   *operator.MongoOperator
	storageOperator *operator.StorageOperator
}

func NewHeartbeatService(mongoOperator *operator.MongoOperator, storageOperator *operator.StorageOperator) *HeartbeatService {
	return &HeartbeatService{mongoOperator: mongoOperator, storageOperator: storageOperator}
}

func (hs *HeartbeatService) StartCheck() {
	go func() {
		for true {
			go func() {
				var err error
				data.Map, err = hs.GetGlobalMap()
				if err != nil {
					log.Println("Get GlobalMap Failed: ", err)
				}
			}()
			go func() {
				newGroups, err := hs.GetGlobalGroups()
				if err != nil {
					log.Println("Get GlobalGroups Failed: ", err)
				}
				hs.checkNew(newGroups)
				data.Groups = newGroups
			}()
			time.Sleep(conf.HeartbeatInternal * time.Second)
		}
	}()
}

func (hs *HeartbeatService) GetGlobalMap() (*data.GlobalMap, error) {
	rawInfos, err := hs.mongoOperator.GetMappingInfos()
	if err != nil {
		return nil, err
	}
	mappingInfos := make([]data.MappingInfo, 0)
	for _, rawInfo := range rawInfos {
		mappingInfos = append(mappingInfos, data.MappingInfo{
			ShardIDStart: rawInfo.ShardIDStart,
			ShardIDEnd:   rawInfo.ShardIDEnd,
			GroupID:      rawInfo.GroupID,
			OldGroupID:   rawInfo.OldGroupID,
		})
	}
	return &data.GlobalMap{MappingInfos: mappingInfos}, nil
}

func (hs *HeartbeatService) GetGlobalGroups() (*data.GlobalGroups, error) {
	rawInfos, err := hs.mongoOperator.GetStorages()
	if err != nil {
		return nil, err
	}
	groupInfos := make(map[string][]data.Storage)
	groupsExist := make(map[string]bool)
	storagesExist := make(map[string]bool)

	// 统计每个Group的可用Storage
	now := time.Now().Unix()
	for _, rawInfo := range rawInfos {
		if now-rawInfo.UpdateTime <= 15 {
			groupInfos[rawInfo.GroupID] = append(groupInfos[rawInfo.GroupID], data.Storage{StorageAddress: rawInfo.StorageAddress})
			groupsExist[rawInfo.GroupID] = true
			storagesExist[rawInfo.StorageAddress] = true
		} else {
			_ = hs.mongoOperator.DeleteStorage(rawInfo.StorageAddress)
		}
	}

	// 对每个组并发PING Storage，按时延排序
	for groupID, _ := range groupInfos {
		var lock sync.Mutex
		wg := sync.WaitGroup{}
		wg.Add(len(groupInfos[groupID]))
		for i := range groupInfos[groupID] {
			go func(i int) {
				start := time.Now()
				hs.storageOperator.PING(groupInfos[groupID][i].StorageAddress)
				lock.Lock()
				groupInfos[groupID][i].Latency = time.Since(start).Milliseconds()
				lock.Unlock()
				wg.Done()
			}(i)
		}
		wg.Wait()
		sort.Slice(groupInfos[groupID], func(i, j int) bool {
			return groupInfos[groupID][i].Latency < groupInfos[groupID][j].Latency
		})
	}

	return &data.GlobalGroups{
		GroupInfos:    groupInfos,
		GroupsExist:   groupsExist,
		StoragesExist: storagesExist,
	}, nil
}

// 检查是否有新的Group，或旧的Group里出现了新的Storage
func (hs *HeartbeatService) checkNew(newGroups *data.GlobalGroups) {
	for groupID, _ := range newGroups.GroupInfos {
		// 新的Group出现
		if !data.Groups.GroupsExist[groupID] {
			log.Println("Find New Group, GroupID: ", groupID)
			// 在映射信息中查找新加入Group对应的OldGroupID
			canFind := false
			var oldGroupID string
			for _, mappingInfo := range data.Map.MappingInfos {
				if mappingInfo.GroupID == groupID {
					canFind = true
					oldGroupID = mappingInfo.OldGroupID
					break
				}
			}
			// Map还没更新，下个心跳再查
			if !canFind {
				delete(newGroups.GroupsExist, groupID)
				continue
			}
			// 系统首次启动，不需要转移数据
			if canFind && oldGroupID == "" {
				continue
			}
			// 需要转移数据
			go func() {
				// 查到可用源storage
				srcAddress := data.Groups.GroupInfos[oldGroupID][0].StorageAddress
				// 筛选需要转移的ChunkID
				chunkIDs, err := hs.storageOperator.GetChunkIDs(srcAddress)
				if err != nil {
					log.Println("Get ChunkIDs failed: ", err)
					return
				}
				toTransChunkIDs := make([]string, 0)
				for _, chunkID := range chunkIDs {
					if getGroupIDsByChunkID(chunkID)[1] == oldGroupID {
						toTransChunkIDs = append(toTransChunkIDs, chunkID)
					}
				}
				// 并发转移
				n := len(newGroups.GroupInfos[groupID])
				wg := sync.WaitGroup{}
				wg.Add(n)
				for _, storage := range newGroups.GroupInfos[groupID] {
					go func(dstAddress string) {
						for i := range toTransChunkIDs {
							_ = hs.storageOperator.SyncChunk(dstAddress, srcAddress, toTransChunkIDs[i])
						}
						wg.Done()
					}(storage.StorageAddress)
				}
				wg.Wait()
				// 并发删除
				n = len(data.Groups.GroupInfos[oldGroupID])
				wg = sync.WaitGroup{}
				wg.Add(n)
				for _, storage := range data.Groups.GroupInfos[oldGroupID] {
					go func(oldStorageAddress string) {
						for i := range toTransChunkIDs {
							_ = hs.storageOperator.DelChunk(oldStorageAddress, toTransChunkIDs[i])
						}
						wg.Done()
					}(storage.StorageAddress)
				}
				wg.Wait()
				// 去掉配置中心对应表项的OldGroupID字段
				sharID := getShardIDByChunkID(toTransChunkIDs[0])
				for _, mappingInfo := range data.Map.MappingInfos {
					if mappingInfo.ShardIDStart <= sharID && sharID < mappingInfo.ShardIDEnd {
						_ = hs.mongoOperator.UpdateMappingInfoDeleteOldGroupID(mappingInfo.ShardIDStart, mappingInfo.ShardIDEnd)
						break
					}
				}
			}()
		} else {
			canSyncAddress := make([]string, 0)
			toSyncAddress := make([]string, 0)
			for _, storage := range newGroups.GroupInfos[groupID] {
				// 新的Storage出现
				if !data.Groups.StoragesExist[storage.StorageAddress] {
					log.Println("Find New Storage, StorageAddress: ", storage.StorageAddress)
					toSyncAddress = append(toSyncAddress, storage.StorageAddress)
				} else {
					canSyncAddress = append(canSyncAddress, storage.StorageAddress)
				}
			}
			for i := range toSyncAddress {
				go func(dstAddress string) {
					_ = hs.storageOperator.SyncAll(dstAddress, canSyncAddress[0])
				}(toSyncAddress[i])
			}
		}
	}
}
