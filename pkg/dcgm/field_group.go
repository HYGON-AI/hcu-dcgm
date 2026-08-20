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
	"sort"
	"strings"
	"sync"
	"time"
)

var (
	FIELD_GROUP_INFO_BUCKET = []byte("FieldGroupInfo")
	FIELD_GROUP_ID_BUCKET   = []byte("FieldGroupId")
)

type FieldGroupInfo struct {
	FieldGroupId   int    `json:"field_group_id"`
	FieldGroupName string `json:"field_group_name"`
	FieldIds       []int  `json:"field_ids"`
}

func createFieldGroup(fieldGroupName string, fieldIds []int) (fieldGroupId int, err error) {
	db, err := OpenDB()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	err = db.Update(func(tx *bbolt.Tx) error {
		infoBucket, err := tx.CreateBucketIfNotExists(FIELD_GROUP_INFO_BUCKET)
		if err != nil {
			return err
		}
		idBucket, err := tx.CreateBucketIfNotExists(FIELD_GROUP_ID_BUCKET)
		if err != nil {
			return err
		}
		if infoBucket != nil {
			err = infoBucket.ForEach(func(k, v []byte) error {
				var info FieldGroupInfo
				if err = json.Unmarshal(v, &info); err != nil {
					return err
				}
				if info.FieldGroupName == fieldGroupName {
					return fmt.Errorf("field group name %s already exists", fieldGroupName)
				}
				return nil
			})
			if err != nil {
				return err
			}
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
				fieldGroupId = freeIds[0]
				remaining := freeIds[1:]
				if len(remaining) > 0 {
					newBytes, err := json.Marshal(&remaining)
					if err != nil {
						return err
					}
					idBucket.Put(FREE_IDS_KEY, newBytes)
				} else {
					idBucket.Delete(FREE_IDS_KEY)
				}
			} else {
				fieldGroupId = btoi(idBucket.Get(NEXT_ID_KEY))
				idBucket.Put(NEXT_ID_KEY, itob(fieldGroupId+1))
			}
		} else {
			fieldGroupId = btoi(idBucket.Get(NEXT_ID_KEY))
			idBucket.Put(NEXT_ID_KEY, itob(fieldGroupId+1))
		}
		fieldGroupInfo := FieldGroupInfo{
			FieldGroupId:   fieldGroupId,
			FieldGroupName: fieldGroupName,
			FieldIds:       fieldIds,
		}
		info, err := json.Marshal(fieldGroupInfo)
		if err != nil {
			return err
		}
		return infoBucket.Put(itob(fieldGroupId), info)
	})
	if err != nil {
		return 0, err
	}
	return fieldGroupId, nil
}

func destroyFieldGroup(fieldGroupId int) (err error) {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Update(func(tx *bbolt.Tx) error {
		infoBucket := tx.Bucket(FIELD_GROUP_INFO_BUCKET)
		if infoBucket == nil || infoBucket.Get(itob(fieldGroupId)) == nil {
			return fmt.Errorf("field group %d does not exist", fieldGroupId)
		}
		err = infoBucket.Delete(itob(fieldGroupId))
		if err != nil {
			return err
		}

		idBucket := tx.Bucket(FIELD_GROUP_ID_BUCKET)
		if idBucket == nil {
			idBucket, err = tx.CreateBucket(FIELD_GROUP_ID_BUCKET)
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
		freeIds = append(freeIds, fieldGroupId)
		updatedFreeIds, err := json.Marshal(freeIds)
		if err != nil {
			return err
		}
		return idBucket.Put(FREE_IDS_KEY, updatedFreeIds)
	})
	return err
}

func getFieldGroupInfo(fieldGroupId int) (fieldGroupInfo FieldGroupInfo, err error) {
	db, err := OpenDB()
	if err != nil {
		return fieldGroupInfo, err
	}
	defer db.Close()

	err = db.View(func(tx *bbolt.Tx) error {
		infoBucket := tx.Bucket(FIELD_GROUP_INFO_BUCKET)
		if infoBucket == nil {
			return fmt.Errorf("field group %d does not exist", fieldGroupId)
		}

		info := infoBucket.Get(itob(fieldGroupId))
		if info == nil {
			return fmt.Errorf("field group %d does not exist", fieldGroupId)
		}
		return json.Unmarshal(info, &fieldGroupInfo)
	})
	return fieldGroupInfo, err
}

func listAllFieldGroups() (fieldGroups []FieldGroupInfo, err error) {
	db, err := OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	err = db.View(func(tx *bbolt.Tx) error {
		infoBucket := tx.Bucket(FIELD_GROUP_INFO_BUCKET)
		if infoBucket == nil {
			fieldGroups = []FieldGroupInfo{}
			return nil
		}
		return infoBucket.ForEach(func(k, v []byte) error {
			var fieldGroupInfo FieldGroupInfo
			if err := json.Unmarshal(v, &fieldGroupInfo); err != nil {
				return err
			}
			fieldGroups = append(fieldGroups, fieldGroupInfo)
			return nil
		})
	})
	return fieldGroups, err
}

const (
	// DefaultUpdateFreq specifies the default update frequency
	DefaultUpdateFreq = 3 * time.Second
	// DefaultMaxKeepAge specifies the default maximum age to keep samples
	DefaultMaxKeepAge = 0
	// DefaultMaxKeepSamples specifies the default number of samples to keep
	DefaultMaxKeepSamples = 2
)

type FieldMeta struct {
	FieldId     int
	Name        string
	EntityLevel Field_Entity_Group
}

func getEntityLevel(name string) Field_Entity_Group {
	up := strings.ToUpper(name)
	switch {
	case strings.HasPrefix(up, "HCU_"):
		return FE_HCU
	case strings.HasPrefix(up, "VHCU_"):
		return FE_VHCU
	case strings.HasPrefix(up, "SWITCH_"):
		return FE_SWITCH
	case strings.HasPrefix(up, "GI_"):
		return FE_HCU_GI
	case strings.HasPrefix(up, "CI_"):
		return FE_HCU_CI
	case strings.HasPrefix(up, "LINK_"):
		return FE_LINK
	case strings.HasPrefix(up, "CPU_"):
		return FE_CPU
	case strings.HasPrefix(up, "CORE_"):
		return FE_CPU_CORE
	default:
		return FE_NONE
	}
}

func getFieldMetaById(fieldId int) FieldMeta {
	name, ok := FieldIdToName[fieldId]
	if !ok {
		return FieldMeta{
			FieldId:     fieldId,
			Name:        "unknown",
			EntityLevel: FE_NONE,
		}
	}

	entityLevel := getEntityLevel(name)
	return FieldMeta{
		FieldId:     fieldId,
		Name:        name,
		EntityLevel: entityLevel,
	}
}

func listFieldMeta() []FieldMeta {
	fieldMetas := make([]FieldMeta, 0, len(FieldIdToName))
	ids := make([]int, 0, len(FieldIdToName))
	for id := range FieldIdToName {
		ids = append(ids, id)
	}
	sort.Ints(ids)
	for _, fieldId := range ids {
		fieldName := FieldIdToName[fieldId]
		entityLevel := getEntityLevel(fieldName)
		fieldMetas = append(fieldMetas, FieldMeta{
			FieldId:     fieldId,
			Name:        fieldName,
			EntityLevel: entityLevel,
		})
	}
	return fieldMetas
}

type FieldValue_v1 struct {
	FieldId   int
	Timestamp int64
	Value     float64
	Err       error
	Tag       string
}

type FieldValue_v2 struct {
	EntityGroupId Field_Entity_Group
	EntityId      int
	FieldId       int
	Timestamp     int64
	Value         float64
	Err           error
	Tag           string
}

/* -----------------------
   Watcher + Request
------------------------*/

type FieldSample struct {
	ts    time.Time
	value float64
	err   error
}

type FieldWatchRequest struct {
	EntityGroupId  Field_Entity_Group
	EntityId       int
	Fields         []int
	UpdateFreq     time.Duration
	MaxKeepAge     time.Duration
	MaxKeepSamples int32
}

type FieldWatcher struct {
	req     FieldWatchRequest
	stopCh  chan struct{}
	wg      sync.WaitGroup
	mu      sync.Mutex
	samples map[int][]FieldSample
}

func newFieldWatcher(req FieldWatchRequest) *FieldWatcher {
	return &FieldWatcher{
		req:     req,
		stopCh:  make(chan struct{}),
		samples: make(map[int][]FieldSample),
	}
}

func (w *FieldWatcher) start() {
	w.collectOnce()
	ticker := time.NewTicker(w.req.UpdateFreq)
	w.wg.Add(1)

	go func() {
		defer w.wg.Done()
		for {
			select {
			case <-ticker.C:
				w.collectOnce()
			case <-w.stopCh:
				ticker.Stop()
				return
			}
		}
	}()
}

func (w *FieldWatcher) stop() {
	close(w.stopCh)
	w.wg.Wait()
}

func (w *FieldWatcher) collectOnce() {
	now := time.Now()

	w.mu.Lock()
	defer w.mu.Unlock()

	for _, f := range w.req.Fields {
		value, err := getFieldValue(w.req.EntityGroupId, w.req.EntityId, f)
		s := FieldSample{
			ts:    now,
			value: value,
			err:   err,
		}
		w.samples[f] = append(w.samples[f], s)

		if w.req.MaxKeepAge != DefaultMaxKeepAge {
			cutoff := now.Add(-w.req.MaxKeepAge)
			tmp := w.samples[f][:0]
			for _, v := range w.samples[f] {
				if v.ts.After(cutoff) {
					tmp = append(tmp, v)
				}
			}
			w.samples[f] = tmp
		}

		if int32(len(w.samples[f])) > w.req.MaxKeepSamples {
			drop := int32(len(w.samples[f])) - w.req.MaxKeepSamples
			w.samples[f] = w.samples[f][drop:]
		}
	}
}

func (w *FieldWatcher) latest(fieldId int) (FieldValue_v1, bool) {
	w.mu.Lock()
	defer w.mu.Unlock()

	arr := w.samples[fieldId]
	if len(arr) == 0 {
		return FieldValue_v1{}, false
	}

	last := arr[len(arr)-1]
	return FieldValue_v1{
		FieldId:   fieldId,
		Timestamp: last.ts.UnixNano() / 1e6,
		Value:     last.value,
		Err:       last.err,
	}, true
}

/* ---------------------------------
   Global Manager for all watchers
----------------------------------*/

type MetricsManager struct {
	mu       sync.Mutex
	watchers map[GroupEntityPair]*FieldWatcher
}

func newMetricsManager() *MetricsManager {
	return &MetricsManager{
		watchers: make(map[GroupEntityPair]*FieldWatcher),
	}
}

var manager = newMetricsManager()

func (m *MetricsManager) startWatcher(req FieldWatchRequest) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	watchKey := GroupEntityPair{
		req.EntityGroupId,
		req.EntityId,
	}

	if old, ok := m.watchers[watchKey]; ok {
		old.stop()
	}

	w := newFieldWatcher(req)
	m.watchers[watchKey] = w
	w.start()
	return nil
}

func (m *MetricsManager) getWatcher(fieldEntityGroup Field_Entity_Group, entityId int) (*FieldWatcher, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	watchKey := GroupEntityPair{
		fieldEntityGroup,
		entityId,
	}
	w, ok := m.watchers[watchKey]
	return w, ok
}

func (m *MetricsManager) stopWatcher(fieldEntityGroup Field_Entity_Group, entityId int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	watchKey := GroupEntityPair{
		fieldEntityGroup,
		entityId,
	}
	if w, ok := m.watchers[watchKey]; ok {
		w.stop()
		delete(m.watchers, watchKey)
	}
}

func watchFields(hcuIndex int, fieldIds []int) error {
	return watchFieldsWithEntity(FE_HCU, hcuIndex, fieldIds)
}

func watchFieldGroup(hcuIndex int, fieldGroupId int) error {
	return watchFieldGroupWithEntity(FE_HCU, hcuIndex, fieldGroupId)
}

func watchFieldsWithEntity(entityGroup Field_Entity_Group, entityId int, fieldIds []int) error {
	req := FieldWatchRequest{
		EntityGroupId:  entityGroup,
		EntityId:       entityId,
		Fields:         fieldIds,
		UpdateFreq:     DefaultUpdateFreq,
		MaxKeepAge:     DefaultMaxKeepAge,
		MaxKeepSamples: DefaultMaxKeepSamples,
	}
	return manager.startWatcher(req)
}

func watchFieldGroupWithEntity(entityGroup Field_Entity_Group, entityId int, fieldGroupId int) error {
	fieldGroupInfo, err := getFieldGroupInfo(fieldGroupId)
	if err != nil {
		return err
	}
	fieldIds := fieldGroupInfo.FieldIds

	return watchFieldsWithEntity(entityGroup, entityId, fieldIds)
}

func watchFieldsWithGroup(fieldGroupId int, groupId int) error {
	return watchFieldsWithGroupEx(fieldGroupId, groupId, DefaultUpdateFreq, DefaultMaxKeepAge, DefaultMaxKeepSamples)
}

func watchFieldsWithEntityGroup(
	fieldIds []int, groupId int, updateFreq time.Duration, maxKeepAge time.Duration, maxKeepSamples int32,
) error {
	groupInfo, err := getGroupInfo(groupId)
	if err != nil {
		return err
	}

	for _, e := range groupInfo.EntityList {
		req := FieldWatchRequest{
			EntityGroupId:  e.EntityGroupId,
			EntityId:       e.EntityId,
			Fields:         fieldIds,
			UpdateFreq:     updateFreq,
			MaxKeepAge:     maxKeepAge,
			MaxKeepSamples: maxKeepSamples,
		}
		err = manager.startWatcher(req)
		if err != nil {
			return err
		}
	}
	return nil
}

func watchFieldsWithGroupEx(
	fieldGroupId int, groupId int, updateFreq time.Duration, maxKeepAge time.Duration, maxKeepSamples int32,
) error {
	fieldGroupInfo, err := getFieldGroupInfo(fieldGroupId)
	if err != nil {
		return err
	}
	fields := fieldGroupInfo.FieldIds

	err = watchFieldsWithEntityGroup(fields, groupId, updateFreq, maxKeepAge, maxKeepSamples)
	if err != nil {
		return err
	}
	return nil
}

func unWatchFields(hcuIndex int) {
	unWatchFieldsWithEntity(FE_HCU, hcuIndex)
}

func unWatchFieldsWithEntity(entityGroup Field_Entity_Group, entityId int) {
	manager.stopWatcher(entityGroup, entityId)
}

func unWatchFieldsWithGroup(groupId int) error {
	groupInfo, err := getGroupInfo(groupId)
	if err != nil {
		return err
	}

	for _, entity := range groupInfo.EntityList {
		manager.stopWatcher(entity.EntityGroupId, entity.EntityId)
	}
	return nil
}

/* --------------------------
   Query APIs
--------------------------*/

func getLatestValuesForFields(hcuIndex int, fields []int) ([]FieldValue_v1, error) {
	return entityGetLatestValues(FE_HCU, hcuIndex, fields)
}

func entityGetLatestValues(entityGroup Field_Entity_Group, entityId int, fields []int) ([]FieldValue_v1, error) {
	watcher, ok := manager.getWatcher(entityGroup, entityId)
	if !ok {
		return nil, fmt.Errorf("no watcher for entity (%d, %d)", entityGroup, entityId)
	}
	var result []FieldValue_v1
	for _, f := range fields {
		if val, ok := watcher.latest(f); ok {
			result = append(result, val)
		}
	}
	return result, nil
}

func entitiesGetLatestValues(entities []GroupEntityPair, fields []int) ([]FieldValue_v2, error) {
	var result []FieldValue_v2

	for _, e := range entities {
		watcher, ok := manager.getWatcher(e.EntityGroupId, e.EntityId)
		if !ok {
			continue
		}
		for _, f := range fields {
			if val, ok := watcher.latest(f); ok {
				result = append(result, FieldValue_v2{
					EntityGroupId: e.EntityGroupId,
					EntityId:      e.EntityId,
					FieldId:       val.FieldId,
					Timestamp:     val.Timestamp,
					Value:         val.Value,
					Err:           val.Err,
					Tag:           val.Tag,
				})
			}
		}
	}
	return result, nil
}
