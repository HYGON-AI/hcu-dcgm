/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package dcgm

import (
	"fmt"
	"strings"

	"go.etcd.io/bbolt"
)

// 默认值常量
const (
	healthCheckBucket = "HealthCheckBucket"
	healthCheckKey    = "health_check"
)

// 初始化配置
//func initConfig() {
//	configPath, err := getConfigPath("config.json")
//	if err != nil {
//		glog.Warningf("Failed to get config path: %v", err)
//		return
//	}
//
//	file, err := os.Open(configPath)
//	if err != nil {
//		glog.Warningf("Config file not found, using default values: %v", err)
//		return
//	}
//	defer file.Close()
//
//	var config Config
//	if err := json.NewDecoder(file).Decode(&config); err != nil {
//		glog.Warningf("Error decoding config file, using default values: %v", err)
//		return
//	}
//
//	// 设置配置变量
//	minPowerThreshold = cfg.Health.MinPowerThreshold
//	maxPowerThreshold = cfg.Health.MaxPowerThreshold
//	maxMemoryUsageThreshold = cfg.Health.MaxMemoryUsageThreshold
//	maxTemperatureThreshold = cfg.Health.MaxTemperatureThreshold
//	glog.V(5).Infof("Loaded configuration: %+v", config)
//}

// setHealthCheck 设置健康检查信息
func setHealthCheckConfig(enabled bool, options []string) (err error) {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	healthConfig := HealthCheckConfig{
		Enabled: enabled,
		Options: options,
	}

	return Create(db, healthCheckBucket, healthCheckKey, healthConfig)
}

// getHealthCheckConfig 获取健康检查配置
func getHealthCheckConfig() (healthConfig HealthCheckConfig, err error) {
	db, err := OpenDB()
	if err != nil {
		return healthConfig, err
	}
	defer db.Close()

	err = db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(healthCheckBucket))
		if bucket == nil {
			// bucket 不存在，没配置返回空
			return nil
		}

		val := bucket.Get([]byte(healthCheckKey))
		if val == nil {
			// key 不存在，没配置返回空
			return nil
		}

		return Read(db, healthCheckBucket, healthCheckKey, &healthConfig)
	})

	return healthConfig, err
}

// deleteHealthCheck 删除健康检查信息
func deleteHealthCheckConfig() (err error) {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	return Delete(db, healthCheckBucket, healthCheckKey)
}

// HealthCheckList 列出所有健康检查项（仅示例用途）
func healthCheckConfigList() (map[string]interface{}, error) {
	db, err := OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var healthChecks map[string]interface{}
	err = ListItems(db, healthCheckBucket, &healthChecks)
	if err != nil {
		return nil, err
	}

	return healthChecks, nil
}

// ---------- 通用检查函数 ----------
func performCheck(
	device int,
	checkType string,
	enabled bool,
	options []string,
	threshold interface{},
	getValue func(int) (interface{}, error),
) SystemWatch {

	// enabled=true 且 options 非空 且不包含该项 → skip
	if enabled && len(options) > 0 && !containsOption(options, checkType) {
		return SystemWatch{
			Type:   checkType,
			Status: HealthStatusSkipped,
		}
	}

	value, err := getValue(device)
	if err != nil {
		return SystemWatch{
			Type:   checkType,
			Status: HealthStatusFailure,
			Error:  err.Error(),
		}
	}

	switch checkType {

	case PowerHealth:
		powerVal := value.(int64)
		rng := threshold.([2]int64)
		min, max := rng[0], rng[1]

		if powerVal < min || powerVal > max {
			return SystemWatch{
				Type:   checkType,
				Status: HealthStatusWarning,
				Result: fmt.Sprintf(
					"Current Power: %d (out of safe range %d - %d)",
					powerVal, min, max,
				),
				Error: "Power out of safe range",
			}
		}

		return SystemWatch{
			Type:   checkType,
			Status: HealthStatusHealthy,
			Result: fmt.Sprintf(
				"Current Power: %d, Safe Range: [%d, %d]",
				powerVal, min, max,
			),
		}

	case MemoryHealth:
		percent := value.(float64)
		max := threshold.(float64)

		if percent > max {
			return SystemWatch{
				Type:   checkType,
				Status: HealthStatusWarning,
				Result: fmt.Sprintf(
					"Memory usage %.2f%% exceeds threshold %.2f%%",
					percent, max,
				),
				Error: "Memory usage exceeds safe threshold",
			}
		}

		return SystemWatch{
			Type:   checkType,
			Status: HealthStatusHealthy,
			Result: fmt.Sprintf(
				"Memory usage %.2f%% within safe threshold %.2f%%",
				percent, max,
			),
		}

	case TemperatureHealth:
		temp := value.(float64)
		max := threshold.(float64)

		if temp > max {
			return SystemWatch{
				Type:   checkType,
				Status: HealthStatusWarning,
				Result: fmt.Sprintf(
					"Temperature %.2f°C exceeds threshold %.2f°C",
					temp, max,
				),
				Error: "Temperature exceeds safe threshold",
			}
		}

		return SystemWatch{
			Type:   checkType,
			Status: HealthStatusHealthy,
			Result: fmt.Sprintf(
				"Temperature %.2f°C within safe threshold %.2f°C",
				temp, max,
			),
		}

	case PerformanceHealth:
		return SystemWatch{
			Type:   checkType,
			Status: HealthStatusHealthy,
			Result: fmt.Sprintf("PerfLevel: %v", value),
		}
	}

	return SystemWatch{Type: checkType, Status: HealthStatusSkipped}
}

// ---------- healthCheckById ----------
func healthCheckById(
	dvIdList []int,
	checkHealthConfig bool,
) (deviceHealths []DeviceHealth, err error) {

	var options []string
	if checkHealthConfig {
		cfg, err := getHealthCheckConfig()
		if err != nil {
			return nil, err
		}
		// 注意：只有 options 非空才生效
		options = cfg.Options
	}

	for _, device := range dvIdList {
		dh := DeviceHealth{
			HCU:    uint(device),
			Status: HealthStatusHealthy,
		}

		// ---------- 1. NUMA ----------
		if !checkHealthConfig || len(options) == 0 || containsOption(options, NumaTopologyHealth) {
			infos, err := ShowNumaTopology([]int{device})
			if err != nil {
				dh.Watches = append(dh.Watches, SystemWatch{
					Type:   NumaTopologyHealth,
					Status: HealthStatusFailure,
					Error:  err.Error(),
				})
				dh.Status = HealthStatusFailure
			} else {
				var desc string
				for _, i := range infos {
					desc += fmt.Sprintf(
						"DeviceID:%d NumaNode:%d Affinity:%d ",
						i.DeviceID, i.NumaNode, i.NumaAffinity,
					)
				}
				dh.Watches = append(dh.Watches, SystemWatch{
					Type:   NumaTopologyHealth,
					Status: HealthStatusHealthy,
					Result: desc,
				})
			}
		}

		// ---------- 2. PCIe Bandwidth ----------
		if !checkHealthConfig || len(options) == 0 || containsOption(options, PcieBandwidthHealth) {
			infos, err := ShowPcieBw([]int{device})
			if err != nil {
				dh.Watches = append(dh.Watches, SystemWatch{
					Type:   PcieBandwidthHealth,
					Status: HealthStatusFailure,
					Error:  err.Error(),
				})
				dh.Status = HealthStatusFailure
			} else {
				var desc string
				for _, i := range infos {
					desc += fmt.Sprintf(
						"Sent:%.0f Recv:%.0f BW:%.2fMB/s ",
						i.Sent, i.Received, i.Bw,
					)
				}
				dh.Watches = append(dh.Watches, SystemWatch{
					Type:   PcieBandwidthHealth,
					Status: HealthStatusHealthy,
					Result: desc,
				})
			}
		}

		// ---------- 3. Power ----------
		pw := performCheck(
			device,
			PowerHealth,
			checkHealthConfig,
			options,
			[2]int64{
				int64(minPowerThreshold),
				int64(maxPowerThreshold),
			},
			func(d int) (interface{}, error) {
				return Power(d)
			},
		)
		dh.Watches = append(dh.Watches, pw)

		// ---------- 4. Memory ----------
		mem := performCheck(
			device,
			MemoryHealth,
			checkHealthConfig,
			options,
			float64(maxMemoryUsageThreshold),
			func(d int) (interface{}, error) {
				used, total, err := MemInfo(d, "vram")
				if err != nil {
					return 0.0, err
				}
				return float64(used) / float64(total) * 100, nil
			},
		)
		dh.Watches = append(dh.Watches, mem)

		// ---------- 5. Temperature ----------
		tmp := performCheck(
			device,
			TemperatureHealth,
			checkHealthConfig,
			options,
			float64(maxTemperatureThreshold),
			func(d int) (interface{}, error) {
				return Temperature(d)
			},
		)
		dh.Watches = append(dh.Watches, tmp)

		// ---------- 6. Performance ----------
		perf := performCheck(
			device,
			PerformanceHealth,
			checkHealthConfig,
			options,
			nil,
			func(d int) (interface{}, error) {
				return PerfLevel(d)
			},
		)
		dh.Watches = append(dh.Watches, perf)

		// ---------- 7. ECC Blocks ----------
		if !checkHealthConfig || len(options) == 0 || containsOption(options, EccBlocksHealth) {
			infos, err := EccBlocksInfo(device)
			if err != nil {
				dh.Watches = append(dh.Watches, SystemWatch{
					Type:   EccBlocksHealth,
					Status: HealthStatusFailure,
					Error:  err.Error(),
				})
				dh.Status = HealthStatusFailure
			} else {
				var warns []string
				for _, b := range infos {
					if b.CE > 0 || b.UE > 0 {
						warns = append(warns,
							fmt.Sprintf("%s CE:%d UE:%d", b.Block, b.CE, b.UE))
					}
				}
				if len(warns) > 0 {
					dh.Watches = append(dh.Watches, SystemWatch{
						Type:   EccBlocksHealth,
						Status: HealthStatusWarning,
						Result: strings.Join(warns, "; "),
					})
					dh.Status = HealthStatusWarning
				} else {
					dh.Watches = append(dh.Watches, SystemWatch{
						Type:   EccBlocksHealth,
						Status: HealthStatusHealthy,
						Result: "All ECC blocks are healthy",
					})
				}
			}
		}

		// ---------- 8. HCU Usage ----------
		if !checkHealthConfig || len(options) == 0 || containsOption(options, HCUUsageHealth) {
			use, err := HCUUse(device)
			if err != nil {
				dh.Watches = append(dh.Watches, SystemWatch{
					Type:   HCUUsageHealth,
					Status: HealthStatusFailure,
					Error:  err.Error(),
				})
				dh.Status = HealthStatusFailure
			} else {
				dh.Watches = append(dh.Watches, SystemWatch{
					Type:   HCUUsageHealth,
					Status: HealthStatusHealthy,
					Result: fmt.Sprintf("HCU Usage: %.2f%%", float64(use)),
				})
			}
		}

		deviceHealths = append(deviceHealths, dh)
	}

	return
}

func containsOption(options []string, option string) bool {
	for _, opt := range options {
		if opt == option {
			return true
		}
	}
	return false
}
