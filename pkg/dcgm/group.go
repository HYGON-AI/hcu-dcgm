/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package dcgm

import "C"
import (
	"encoding/json"
	"fmt"

	"go.etcd.io/bbolt"
)

// DCGM_GROUP_MAX_ENTITIES represents the maximum number of entities allowed in a group
const (
	DCGM_GROUP_MAX_ENTITIES = 1024
)

var (
	GROUP_INFO_BUCKET = []byte("GroupInfo")
	GROUP_ID_BUCKET   = []byte("GroupId")
	NEXT_ID_KEY       = []byte("next_id")
	FREE_IDS_KEY      = []byte("free_ids")
)

// Field_Entity_Group represents the type of DCGM entity
type Field_Entity_Group int

const (
	// FE_NONE represents no entity type
	FE_NONE Field_Entity_Group = iota
	// FE_HCU represents a HCU device entity
	FE_HCU
	// FE_VHCU represents a virtual HCU entity
	FE_VHCU
	// FE_SWITCH represents an HYSwitch entity
	FE_SWITCH
	// FE_HCU_GI represents a HCU GPU instance entity
	FE_HCU_GI
	// FE_HCU_CI represents a HCU compute instance entity
	FE_HCU_CI
	// FE_LINK represents an HYLink entity
	FE_LINK
	// FE_CPU represents a CPU entity
	FE_CPU
	// FE_CPU_CORE represents a CPU core entity
	FE_CPU_CORE
	// FE_COUNT represents the total number of entity types
	FE_COUNT
)

// String returns a string representation of the Field_Entity_Group
func (e Field_Entity_Group) String() string {
	switch e {
	case FE_HCU:
		return "HCU"
	case FE_VHCU:
		return "vHCU"
	case FE_SWITCH:
		return "HySwitch"
	case FE_HCU_GI:
		return "HCU GPU Instance"
	case FE_HCU_CI:
		return "HCU Compute Instance"
	case FE_LINK:
		return "HyLink"
	case FE_CPU:
		return "CPU"
	case FE_CPU_CORE:
		return "CPU Core"
	}
	return "unknown"
}

// GroupEntityPair represents a DCGM entity and its group identifier
type GroupEntityPair struct {
	// EntityGroupId specifies the type of the entity
	EntityGroupId Field_Entity_Group `json:"entity_group_id"`
	// EntityId is the unique identifier for this entity
	EntityId int `json:"entity_id"`
}

// GroupInfo contains information about a DCGM group
type GroupInfo struct {
	GroupId    int               `json:"group_id"`
	GroupName  string            `json:"group_name"`
	EntityList []GroupEntityPair `json:"entity_list"`
}

// createGroup creates a new empty HCU group with the specified name
func createGroup(groupName string) (groupId int, err error) {
	db, err := OpenDB()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	err = db.Update(func(tx *bbolt.Tx) error {
		infoBucket, err := tx.CreateBucketIfNotExists(GROUP_INFO_BUCKET)
		if err != nil {
			return err
		}
		idBucket, err := tx.CreateBucketIfNotExists(GROUP_ID_BUCKET)
		if err != nil {
			return err
		}
		if idBucket.Get(NEXT_ID_KEY) == nil {
			idBucket.Put(NEXT_ID_KEY, itob(0))
		}
		ids := idBucket.Get(FREE_IDS_KEY)
		var freeIds []int
		if ids != nil {
			if err = json.Unmarshal(ids, &freeIds); err != nil {
				return err
			}
			if len(freeIds) > 0 {
				groupId = freeIds[0]
				remaining := freeIds[1:]
				if len(remaining) > 0 {
					newBytes, err := json.Marshal(remaining)
					if err != nil {
						return err
					}
					idBucket.Put(FREE_IDS_KEY, newBytes)
				} else {
					idBucket.Delete(FREE_IDS_KEY)
				}
			} else {
				groupId = btoi(idBucket.Get(NEXT_ID_KEY))
				idBucket.Put(NEXT_ID_KEY, itob(groupId+1))
			}
		} else {
			groupId = btoi(idBucket.Get(NEXT_ID_KEY))
			idBucket.Put(NEXT_ID_KEY, itob(groupId+1))
		}
		groupInfo := GroupInfo{
			GroupId:    groupId,
			GroupName:  groupName,
			EntityList: []GroupEntityPair{},
		}
		info, err := json.Marshal(groupInfo)
		if err != nil {
			return err
		}
		return infoBucket.Put(itob(groupId), info)
	})
	if err != nil {
		return 0, err
	}
	return groupId, nil
}

// createDefaultGroup creates a new group with default HCUs and the specified name
func createDefaultGroup(groupName string) (groupId int, err error) {
	numDevices, err := rsmiNumMonitorDevices()
	if err != nil {
		return 0, err
	}

	entityList := make([]GroupEntityPair, numDevices)
	for i := 0; i < numDevices; i++ {
		entityList[i] = GroupEntityPair{
			EntityGroupId: FE_HCU,
			EntityId:      i,
		}
	}

	groupId, err = createGroup(groupName)
	if err != nil {
		return 0, err
	}

	err = addEntityToGroup(groupId, entityList)
	return groupId, err
}

// addToGroup adds a HCU to an existing group
func addToGroup(groupId int, hcuIndex int) (err error) {
	_, err = GetDeviceId(hcuIndex)
	if err != nil {
		return fmt.Errorf("HCU %d does not exist", hcuIndex)
	}

	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Update(func(tx *bbolt.Tx) error {
		infoBucket := tx.Bucket(GROUP_INFO_BUCKET)
		if infoBucket == nil {
			return fmt.Errorf("group %d does not exist", groupId)
		}
		info := infoBucket.Get(itob(groupId))
		if info == nil {
			return fmt.Errorf("group %d does not exist", groupId)
		}

		var groupInfo GroupInfo
		if err = json.Unmarshal(info, &groupInfo); err != nil {
			return err
		}

		existing := make(map[string]bool)
		for _, entity := range groupInfo.EntityList {
			k := fmt.Sprintf("%d-%d", entity.EntityGroupId, entity.EntityId)
			existing[k] = true
		}
		k := fmt.Sprintf("1-%d", hcuIndex)
		if !existing[k] {
			entity := GroupEntityPair{
				EntityGroupId: FE_HCU,
				EntityId:      hcuIndex,
			}
			groupInfo.EntityList = append(groupInfo.EntityList, entity)
		}

		updatedInfo, err := json.Marshal(groupInfo)
		if err != nil {
			return err
		}
		return infoBucket.Put(itob(groupId), updatedInfo)
	})
	return err
}

// addEntityToGroup adds an entityList to an existing group
func addEntityToGroup(groupId int, entityList []GroupEntityPair) (err error) {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Update(func(tx *bbolt.Tx) error {
		infoBucket := tx.Bucket(GROUP_INFO_BUCKET)
		if infoBucket == nil {
			return fmt.Errorf("group %d does not exist", groupId)
		}
		info := infoBucket.Get(itob(groupId))
		if info == nil {
			return fmt.Errorf("group %d does not exist", groupId)
		}

		var groupInfo GroupInfo
		if err = json.Unmarshal(info, &groupInfo); err != nil {
			return err
		}

		existing := make(map[string]bool)
		for _, entity := range groupInfo.EntityList {
			k := fmt.Sprintf("%d-%d", entity.EntityGroupId, entity.EntityId)
			existing[k] = true
		}
		for _, entity := range entityList {
			k := fmt.Sprintf("%d-%d", entity.EntityGroupId, entity.EntityId)
			if !existing[k] {
				groupInfo.EntityList = append(groupInfo.EntityList, entity)
				existing[k] = true
			}
		}
		updatedInfo, err := json.Marshal(groupInfo)
		if err != nil {
			return err
		}
		return infoBucket.Put(itob(groupId), updatedInfo)
	})
	return err
}

func removeEntityFromGroup(groupId int, entityList []GroupEntityPair) (err error) {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	delSet := make(map[GroupEntityPair]struct{}, len(entityList))
	for _, e := range entityList {
		delSet[e] = struct{}{}
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		infoBucket := tx.Bucket(GROUP_INFO_BUCKET)
		if infoBucket == nil {
			return fmt.Errorf("group %d does not exist", groupId)
		}

		info := infoBucket.Get(itob(groupId))
		if info == nil {
			return fmt.Errorf("group %d does not exist", groupId)
		}

		var groupInfo GroupInfo
		if err = json.Unmarshal(info, &groupInfo); err != nil {
			return err
		}

		filtered := groupInfo.EntityList[:0]
		for _, entity := range groupInfo.EntityList {
			if _, found := delSet[entity]; !found {
				filtered = append(filtered, entity)
			}
		}
		groupInfo.EntityList = filtered
		updatedInfo, err := json.Marshal(groupInfo)
		if err != nil {
			return err
		}
		return infoBucket.Put(itob(groupId), updatedInfo)
	})
	return err
}

// destroyGroup destroys an existing group
func destroyGroup(groupId int) (err error) {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Update(func(tx *bbolt.Tx) error {
		infoBucket := tx.Bucket(GROUP_INFO_BUCKET)
		if infoBucket == nil || infoBucket.Get(itob(groupId)) == nil {
			return fmt.Errorf("group %d does not exist", groupId)
		}
		err = infoBucket.Delete(itob(groupId))
		if err != nil {
			return err
		}

		idBucket := tx.Bucket(GROUP_ID_BUCKET)
		if idBucket == nil {
			idBucket, err = tx.CreateBucket(GROUP_ID_BUCKET)
		}
		ids := idBucket.Get(FREE_IDS_KEY)
		var freeIds []int
		if ids == nil {
			freeIds = []int{}
		} else {
			if err = json.Unmarshal(ids, &freeIds); err != nil {
				return err
			}
		}
		freeIds = append(freeIds, groupId)
		updatedFreeIds, err := json.Marshal(freeIds)
		if err != nil {
			return err
		}
		return idBucket.Put(FREE_IDS_KEY, updatedFreeIds)
	})
	return err
}

// getGroupInfo retrieves information about a DCGM group
func getGroupInfo(groupId int) (groupInfo GroupInfo, err error) {
	db, err := OpenDB()
	if err != nil {
		return groupInfo, err
	}
	defer db.Close()

	err = db.View(func(tx *bbolt.Tx) error {
		infoBucket := tx.Bucket(GROUP_INFO_BUCKET)
		if infoBucket == nil {
			return fmt.Errorf("group %d does not exist", groupId)
		}

		info := infoBucket.Get(itob(groupId))
		if info == nil {
			return fmt.Errorf("group %d does not exist", groupId)
		}
		return json.Unmarshal(info, &groupInfo)
	})
	return groupInfo, err
}

func listAllGroups() (groups []GroupInfo, err error) {
	db, err := OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	err = db.View(func(tx *bbolt.Tx) error {
		infoBucket := tx.Bucket(GROUP_INFO_BUCKET)
		if infoBucket == nil {
			groups = []GroupInfo{}
			return nil
		}
		return infoBucket.ForEach(func(k, v []byte) error {
			var groupinfo GroupInfo
			if err := json.Unmarshal(v, &groupinfo); err != nil {
				return err
			}
			groups = append(groups, groupinfo)
			return nil
		})
	})
	return groups, err
}

func getHCUListFromGroup(groupId int) (hcuInGroup []int, groupName string, err error) {
	groupInfo, err := getGroupInfo(groupId)
	if err != nil {
		return hcuInGroup, groupName, err
	}
	groupName = groupInfo.GroupName
	entityList := groupInfo.EntityList
	for _, entity := range entityList {
		if entity.EntityGroupId == FE_HCU {
			hcuInGroup = append(hcuInGroup, entity.EntityId)
		}
	}
	return hcuInGroup, groupName, nil
}
