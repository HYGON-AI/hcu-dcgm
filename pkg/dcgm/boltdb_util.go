/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package dcgm

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"go.etcd.io/bbolt"
)

// 数据库文件名
const dbName = "boltdb.db"

// 获取数据库的固定路径，放在 /opt/dcgm 下
func getDBPath() string {
	basePath := "/opt/dcgm"
	// 检查 /opt/dcgm 文件夹是否存在，如果不存在则创建
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		if err := os.MkdirAll(basePath, 0755); err != nil {
			panic("无法创建目录 /opt/dcgm: " + err.Error())
		}
	}
	// 返回固定的数据库路径
	return filepath.Join(basePath, dbName)
}

// OpenDB 打开数据库，如果目录不存在则自动创建
func OpenDB() (*bbolt.DB, error) {
	dbPath := getDBPath()

	// 确保目录存在
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %v", err)
	}

	return bbolt.Open(dbPath, 0600, nil)
}

// Create 插入或更新数据到指定桶
func Create(db *bbolt.DB, bucket, key string, value interface{}) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return fmt.Errorf("failed to create bucket: %v", err)
		}

		data, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal data: %v", err)
		}

		return b.Put([]byte(key), data)
	})
}

// Read 从指定桶查询数据
func Read(db *bbolt.DB, bucket, key string, result interface{}) error {
	return db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("bucket %s not found", bucket)
		}

		data := b.Get([]byte(key))
		if data == nil {
			return fmt.Errorf("key %s not found", key)
		}

		return json.Unmarshal(data, result)
	})
}

// Delete 从指定桶删除数据
func Delete(db *bbolt.DB, bucket, key string) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("bucket %s not found", bucket)
		}

		return b.Delete([]byte(key))
	})
}

// ListKeys 获取桶中的所有键
func ListKeys(db *bbolt.DB, bucket string) ([]string, error) {
	var keys []string
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("bucket %s not found", bucket)
		}

		return b.ForEach(func(k, v []byte) error {
			keys = append(keys, string(k))
			return nil
		})
	})
	return keys, err
}

// ListItems 获取桶中的所有数据项
func ListItems(db *bbolt.DB, bucket string, result interface{}) error {
	items := make(map[string]interface{})
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("bucket %s not found", bucket)
		}

		return b.ForEach(func(k, v []byte) error {
			var item interface{}
			if err := json.Unmarshal(v, &item); err != nil {
				return fmt.Errorf("failed to unmarshal item: %v", err)
			}
			items[string(k)] = item
			return nil
		})
	})

	if err == nil {
		// 将 items 赋值给 result
		b, _ := json.Marshal(items)
		_ = json.Unmarshal(b, result)
	}

	return err
}
