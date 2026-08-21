/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package dcgm

import "C"
import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang/glog"
)

// @Summary 初始化 DCGM
// @Description 初始化 (DCGM) 库。
// @Produce json
// @Success 200 {object} string "成功初始化"
// @Failure 500 {object} error "初始化失败"
// @Router /Init [post]
func Init() (err error) {
	initConfig()

	// 获取设备计数，即使有错误也继续执行
	devCount, listErr := listFilesInDevDri()
	if listErr != nil {
		// 记录错误但不立即返回，继续执行
		glog.V(5).Infof("listFilesInDevDri 返回错误但继续执行: %v", listErr)
	}
	glog.V(5).Infof("devCount:%v", devCount)

	maxRetries := 12                   // 最大重试次数
	retryCount := 0                    // 记录连续返回相同设备数量的次数
	lastNumDevices := -1               // 记录上一次获取的设备数量
	restartTimeout := 10 * time.Second // 每次重试等待10秒
	initFailCount := 0                 // rsmiInit 连续失败的计数
	maxInitFails := 6                  // 连续失败最大次数

	live := checkDriverLive()
	if !live {
		err = fmt.Errorf("HCU driver is not live, please ensure the driver is installed correctly and running.")
		// 如果有 listErr，将其与当前错误合并
		if listErr != nil {
			err = fmt.Errorf("%v; 同时出现: %v", err, listErr)
		}
		return err
	}

	for {
		initErr := rsmiInit() // 初始化rsmi，使用新变量避免覆盖 err
		if initErr == nil {
			ShutDown()
			for retryCount < maxRetries {
				rsmiInit()
				numDevices, numErr := NumMonitorDevices() // 获取HCU设备数量
				deviceCount, countErr := DeviceCount()

				// 如果有错误，记录下来但不中断流程
				if numErr != nil {
					glog.V(5).Infof("NumMonitorDevices 错误: %v", numErr)
				}
				if countErr != nil {
					glog.V(5).Infof("DeviceCount 错误: %v", countErr)
				}

				if devCount == numDevices && devCount == deviceCount {
					glog.V(5).Infof("HCU initialization is complete:%v", numDevices)
					probeDriverCaps()
					// 返回 listErr（如果有），否则返回 nil
					return listErr
				} else {
					if numDevices == lastNumDevices {
						retryCount++ // 记录连续返回相同设备数量的次数
					} else {
						retryCount = 0 // 数量变化时重置计数
					}

					glog.V(5).Infof("retryCount:%v", retryCount)
					if retryCount >= maxRetries {
						glog.Warningf(
							"设备数量连续 %d 次相同但仍不一致 (devCount=%d, numDevices=%d, deviceCount=%d)，继续启动",
							maxRetries, devCount, numDevices, deviceCount,
						)
						if numErr != nil {
							glog.Warningf("NumMonitorDevices 错误: %v", numErr)
						}
						if countErr != nil {
							glog.Warningf("DeviceCount 错误: %v", countErr)
						}
						if listErr != nil {
							glog.Warningf("listFilesInDevDri 错误: %v", listErr)
						}
						probeDriverCaps()
						return listErr
					}
					lastNumDevices = numDevices // 更新记录的设备数量
					ShutDown()                  // 数量不相等，执行关闭操作
				}
				time.Sleep(restartTimeout) // 等待10秒
			}
		} else {
			initFailCount++ // 初始化失败，计数加一
			glog.V(5).Infof("初始化失败: %v. 10秒后重试...\n", initErr)

			if initFailCount >= maxInitFails {
				glog.Errorf("rsmiInit 连续 %d 次失败，终止初始化: %v", maxInitFails, initErr)
				// 返回初始化错误，同时包含 listErr（如果有）
				if listErr != nil {
					return fmt.Errorf("rsmiInit 连续失败: %v; 同时出现: %v", initErr, listErr)
				}
				return initErr
			}
		}
		time.Sleep(restartTimeout) // 等待10秒后再次重试
	}
	glog.V(5).Infof("================================================================")
	// 返回 listErr（如果有），否则返回 nil
	return listErr
}

// @Summary 关闭 DCGM
// @Description 关闭 Data Center GPU Manager (DCGM) 库。
// @Produce json
// @Success 200 {object} string "成功关闭"
// @Failure 500 {object} error "关闭失败"
// @Router /ShutDown [post]
func ShutDown() error {
	return rsmiShutdown()
}

// @Summary 获取 GPU 数量
// @Description 获取监视的 GPU 数量。
// @Produce json
// @Success 200 {int} int "GPU 数量"
// @Failure 500 {object} error "获取 GPU 数量失败"
// @Router /NumMonitorDevices [get]
func NumMonitorDevices() (int, error) {
	return rsmiNumMonitorDevices()
}

// 获取设备利用率计数器
// @Summary 获取设备利用率计数器
// @Description 根据设备索引获取利用率计数器
// @Param dvInd query int true "设备索引"
// @Param utilizationCounters body []UtilizationCounter true "利用率计数器对象列表"
// @Param count query int true "计数器的数量"
// @Success 200 {object} int64 "返回的时间戳"
// @Failure 400 {object} error "请求失败"
// @Router /utilizationcount [post]
func UtilizationCount(dvInd int, utilizationCounters []UtilizationCounter, count int) (timestamp int64, err error) {
	return rsmiUtilizationCountGet(dvInd, utilizationCounters, count)
}

// DeviceUtilizationInfoGet 获取指定设备的 GPU 与显存粗粒度利用率计数器快照。
//
// @Summary 获取设备 GPU 与显存利用率计数器
// @Description 内部固定查询 RSMI_COARSE_GRAIN_GFX_ACTIVITY（GPU）和 RSMI_COARSE_GRAIN_MEM_ACTIVITY（显存）
// @Description 两个粗粒度计数器，一次调用返回两者及采样时间戳。
// @Description 若硬件不支持（计数器值为 0xFFFFFFFF），Supported 字段为 false，计数器字段为 0。
// @Param dvInd 设备索引，从 0 开始。
// @Return info DeviceUtilizationInfo，包含 GfxActivity、MemActivity、Timestamp 和 Supported。
// @Return err 底层 RSMI 调用失败时返回错误。
func DeviceUtilizationInfoGet(dvInd int) (info DeviceUtilizationInfo, err error) {
	counters := []UtilizationCounter{
		{Type: RSMI_COARSE_GRAIN_GFX_ACTIVITY},
		{Type: RSMI_COARSE_GRAIN_MEM_ACTIVITY},
	}
	ts, err := rsmiUtilizationCountGet(dvInd, counters, len(counters))
	if err != nil {
		return info, fmt.Errorf("get device utilization info: %w", err)
	}
	info.Timestamp = ts
	// 0xFFFFFFFF 表示硬件/驱动不支持该计数器
	const unsupported = uint64(0xFFFFFFFF)
	if counters[0].Value == unsupported || counters[1].Value == unsupported {
		info.Supported = false
		return info, nil
	}
	info.GfxActivity = counters[0].Value
	info.MemActivity = counters[1].Value
	info.Supported = true
	return info, nil
}

// @Description 根据设备 ID 获取设备名称。
// @Produce json
// @Param dvInd path int true "设备 ID"
// @Success 200 {string} name "设备名称"
// @Failure 400 {object} error "请求错误"
// @Failure 404 {object} error "设备未找到"
// @Router /DevName [get]
func DevName(dvInd int) (name string, err error) {
	return rsmiDevNameGet(dvInd)
}

// @Summary 设备的唯一标识符
// @Description 根据设备 ID 获取设备的唯一标识符。
// @Produce json
// @Param dvInd path int true "设备 ID"
// @Success 200 {string} deviceId "设备的唯一标识符"
// @Failure 400 {object} error "请求错误"
// @Failure 404 {object} error "设备未找到"
// @Router /DevName [get]
func GetDeviceId(dvInd int) (deviceId string, err error) {
	deviceId, err = rsmiDevSerialNumberGet(dvInd)
	return
}

func GetDeviceUniqueId(dvInd int) (string, error) {
	uniqueIdInt, err := rsmiDevUniqueIdGet(dvInd)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("0x%016x", uint64(uniqueIdInt)), nil
}

// 获取设备SKU
// @Summary 获取设备SKU
// @Description 根据设备索引获取SKU
// @Param dvInd query int true "设备索引"
// @Success 200 {int} sku "返回设备SKU"
// @Failure 400 {object} error "请求失败"
// @Router /DevSku [get]
func DevSku(dvInd int) (sku int, err error) {
	return rsmiDevSkuGet(dvInd)
}

// 获取设备品牌名称
// @Summary 获取设备品牌名称
// @Description 根据设备索引获取品牌名称
// @Param dvInd query int true "设备索引"
// @Success 200 {string} brand "设备品牌名称"
// @Failure 400 {object} error "请求失败"
// @Router /DevBrand [get]
func DevBrand(dvInd int) (brand string, err error) {
	return rsmiDevBrandGet(dvInd)
}

// 获取设备供应商名称
// @Summary 获取设备供应商名称
// @Description 根据设备索引获取供应商名称
// @Param dvInd query int true "设备索引"
// @Success 200 {string} vendorName "返回设备供应商名称"
// @Failure 400 {object} error "请求失败"
// @Router /DevVendorName [get]
func DevVendorName(dvInd int) (vendorName string, err error) {
	return rsmiDevVendorNameGet(dvInd)
}

// 获取设备显存供应商名称
// @Summary 获取设备显存供应商名称
// @Description 根据设备索引获取显存供应商名称
// @Param dvInd query int true "设备索引"
// @Success 200 {string} vramVendor "返回显存供应商名称"
// @Failure 400 {object} error "请求失败"
// @Router /DevVramVendor [get]
func DevVramVendor(dvInd int) (vramVendor string, err error) {
	return rsmiDevVramVendorGet(dvInd)
}

// @Summary 获取可用的 PCIe 带宽列表
// @Description 根据设备 ID 获取设备的可用 PCIe 带宽列表。
// @Produce json
// @Param dvInd path int true "设备 ID"
// @Success 200 {object} PcieBandwidth "PCIe 带宽列表"
// @Failure 400 {object} error "请求错误"
// @Failure 404 {object} error "设备未找到"
// @Router /DevPciBandwidth [get]
func DevPciBandwidth(dvInd int) (PcieBandwidth PcieBandwidth, err error) {
	return rsmiDevPciBandwidthGet(dvInd)
}

// 设置设备的PCIe带宽
// @Summary 设置设备的PCIe带宽
// @Description 根据设备索引和带宽位掩码设置可用的PCIe带宽
// @Param dvInd query int true "设备索引"
// @Param bwBitmask query int64 true "带宽位掩码，指示要启用(1)或禁用(0)的带宽索引"
// @Success 200 {string} string "设置成功"
// @Failure 400 {object} error "请求失败"
// @Router /DevPciBandwidthSet [post]
func DevPciBandwidthSet(dvInd int, bwBitmask int64) (err error) {
	return rsmiDevPciBandwidthSet(dvInd, bwBitmask)
}

// @Summary 获取内存使用百分比
// @Description 根据设备 ID 获取设备内存的CollectDeviceMetrics使用百分比。
// @Produce json
// @Param dvInd path int true "设备 ID"
// @Success 200 {int} busyPercent "内存使用百分比"
// @Failure 400 {object} error "请求错误"
// @Failure 404 {object} error "设备未找到"
// @Router /MemoryPercent [get]
func MemoryPercent(dvInd int) (busyPercent int, err error) {
	return rsmiDevMemoryBusyPercentGet(dvInd)
}

func MemoryTotal(dvInd int) (memoryTotal float64, err error) {
	//获取设备内存总量
	memory, _ := rsmiDevMemoryTotalGet(dvInd, RSMI_MEM_TYPE_FIRST)
	memoryTotal, _ = strconv.ParseFloat(fmt.Sprintf("%f", float64(memory)/1.0), 64)
	glog.V(5).Infof("HCU[%v] 内存总量: %.0f", dvInd, memoryTotal)
	return
}

func MemoryUsed(dvInd int) (memoryUsed float64, err error) {
	memory, _ := rsmiDevMemoryUsageGet(dvInd, RSMI_MEM_TYPE_FIRST)
	memoryUsed, _ = strconv.ParseFloat(fmt.Sprintf("%f", float64(memory)/1.0), 64)
	glog.V(5).Infof("HCU[%v] 内存使用量: %.0f", dvInd, memoryUsed)
	return
}

// 检测内存保留页
func MemoryReservedPages(dvInd int) (records []RetiredPageRecord, err error) {
	_, records, err = rsmiDevMemoryReservedPagesGet(dvInd)
	return
}

// MemoryAvailable 返回驱动当前可向应用分配的近似剩余显存量（字节）
// 计算方式：total_vram - used_vram，对应 DCGM_FI_DEV_FB_RESERVED 的 HCU 近似语义
func MemoryAvailable(dvInd int) (available int64, err error) {
	total, err := rsmiDevMemoryTotalGet(dvInd, RSMI_MEM_TYPE_FIRST)
	if err != nil {
		return 0, err
	}
	used, err := rsmiDevMemoryUsageGet(dvInd, RSMI_MEM_TYPE_FIRST)
	if err != nil {
		return 0, err
	}
	available = total - used
	return
}

// GetTempBySensor 获取指定设备、温度传感器和度量维度的温度值。
//
// @Summary 获取指定温度传感器的度量值
// @Description 调用 RSMI `rsmi_dev_temp_metric_get` 获取温度。底层返回单位为毫摄氏度，
// @Description 本函数转换为摄氏度并保留两位小数。
// @Param dvInd 设备索引；取值为 Init 后枚举到的物理设备序号，从 0 开始。
// @Param sensorType 温度传感器类型：0=EDGE（GPU 边缘）、1=JUNCTION（热点）、2=MEMORY（显存）、
// @Param sensorType 3=VR_GFX、4=VR_SOC、5=VR_MEM0、6=VR_MEM1、7-10=HBM_0-HBM_3、11=CORE。
// @Param metric 温度度量类型：RSMI_TEMP_CURRENT（当前值）、RSMI_TEMP_MAX/MIN（上下限）、
// @Param metric RSMI_TEMP_CRITICAL（临界上限）、RSMI_TEMP_EMERGENCY（紧急上限）、
// @Param metric RSMI_TEMP_LOWEST/HIGHEST（历史最低/最高）及相应的 *_HYST、CRIT_MIN、OFFSET 枚举值。
// @Return temp 温度值，单位为摄氏度。
// @Return err 底层 RSMI 调用失败时返回错误；设备或传感器/度量组合不可用时可能为
// @Return err RSMI_STATUS_NOT_SUPPORTED，参数不合法或组合不满足驱动前置条件时可能为 RSMI_STATUS_INVALID_ARGS。
func GetTempBySensor(dvInd int, sensorType int, metric RSMITemperatureMetric) (temp float64, err error) {
	raw, err := rsmiDevTempMetricGet(dvInd, sensorType, metric)
	if err != nil {
		return 0, err
	}
	temp, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", float64(raw)/1000.0), 64)
	return
}

// GetDeviceTemperatureInfo 获取指定设备四个常用传感器的当前温度。
//
// @Summary 获取设备四路当前温度
// @Description 依次读取 EDGE、JUNCTION、MEMORY、CORE 传感器的 RSMI_TEMP_CURRENT 指标，
// @Description 返回值的四个字段均为摄氏度。底层 RSMI 返回单位为毫摄氏度，
// @Description 由 GetTempBySensor 统一转换并保留两位小数。
// @Param dvInd 设备索引；取值为 Init 后枚举到的物理设备序号，从 0 开始。
// @Return temperatures 指定设备的 EDGE(0)、JUNCTION(1)、MEMORY(2)、CORE(11) 当前温度。
// @Return err 任一传感器读取失败时，仍返回其他成功温度；失败字段为零值。
// @Return err 错误使用 errors.Join 聚合，并保留失败传感器名称及原始 RSMI 错误。
func GetDeviceTemperatureInfo(dvInd int) (temperatures DeviceTemperatureInfo, err error) {
	var errs []error

	temperatures.Edge, err = GetTempBySensor(dvInd, 0, RSMI_TEMP_CURRENT)
	if err != nil {
		errs = append(errs, fmt.Errorf("get edge temperature: %w", err))
	}

	temperatures.Junction, err = GetTempBySensor(dvInd, 1, RSMI_TEMP_CURRENT)
	if err != nil {
		errs = append(errs, fmt.Errorf("get junction temperature: %w", err))
	}

	temperatures.Memory, err = GetTempBySensor(dvInd, 2, RSMI_TEMP_CURRENT)
	if err != nil {
		errs = append(errs, fmt.Errorf("get memory temperature: %w", err))
	}

	temperatures.Core, err = GetTempBySensor(dvInd, 11, RSMI_TEMP_CURRENT)
	if err != nil {
		errs = append(errs, fmt.Errorf("get core temperature: %w", err))
	}

	return temperatures, errors.Join(errs...)
}

// MemoryTemperature 返回显存颗粒当前温度（摄氏度）
// 对应 DCGM_FI_DEV_MEMORY_TEMP；底层传感器 RSMI_TEMP_TYPE_MEMORY(2)
func MemoryTemperature(dvInd int) (float64, error) {
	return GetTempBySensor(dvInd, 2, RSMI_TEMP_CURRENT)
}

// BoardTemperature 返回板卡边缘温度（摄氏度）
// 对应 DCGM_FI_DEV_BOARD_TEMP；底层传感器 RSMI_TEMP_TYPE_EDGE(0)
// Temperature() 语义相同，BoardTemperature 命名更明确
func BoardTemperature(dvInd int) (float64, error) {
	return GetTempBySensor(dvInd, 0, RSMI_TEMP_CURRENT)
}

// DevThrottleFlags 解析指定设备当前的降频状态
// 内部调用 DevGpuMetricsInfo 读取 ThrottleStatus 位掩码，再按 bit 拆分
// 各 bit 偏移见 ThrottleFlags 结构体注释；如硬件厂商变更定义请同步修正
func DevThrottleFlags(dvInd int) (flags ThrottleFlags, err error) {
	m, err := DevGpuMetricsInfo(dvInd)
	if err != nil {
		return flags, err
	}
	s := m.ThrottleStatus
	flags = ThrottleFlags{
		Thermal:    s&(1<<0) != 0,
		Power:      s&(1<<1) != 0,
		Slowdown:   s&(1<<2) != 0,
		BoardLimit: s&(1<<3) != 0,
	}
	return
}

// 获取设备温度值
//func DevTemp(dvInd int) int64 {
//	return go_rsmi_dev_temp_metric_get(dvInd)
//}

// @Summary 设置设备 PowerPlay 性能级别
// @Description 根据设备 ID 设置 PowerPlay 性能级别。
// @Produce json
// @Param dvInd path int true "设备 ID"
// @Param level query string true "要设置的性能级别"
// @Success 200 {string} string "操作成功"
// @Failure 400 {object} error "请求错误"
// @Failure 404 {object} error "设备未找到"
// @Router /DevPerfLevelSet [post]
func DevPerfLevelSet(dvInd int, level DevPerfLevel) error {
	return rsmiDevPerfLevelSet(dvInd, level)
}

// DevGpuMetricsInfo 获取 GPU 度量信息
// @Summary 获取 GPU 度量信息
// @Description 根据设备 ID 获取 GPU 的度量信息。
// @Produce json
// @Param dvInd query int true "设备 ID"
// @Success 200 {object} RSMIGPUMetrics "GPU 度量信息"
// @Failure 400 {object} error "请求错误"
// @Failure 404 {object} error "设备未找到"
// @Router /DevGpuMetricsInfo [get]
func DevGpuMetricsInfo(dvInd int) (gpuMetrics RSMIGPUMetrics, err error) {
	return rsmiDevGpuMetricsInfoGet(dvInd)
}

// 获取设备功率有效值范围
// @Summary 获取设备功率有效值范围
// @Description 根据设备索引和传感器ID获取设备的最大和最小功率值范围
// @Param dvInd query int true "设备索引"
// @Success 200 {int64} powerMax 返回最大值"
// @Success 200 {int64} powerMin 返回最小值"
// @Failure 400 {object} error "请求失败"
// @Router /DevPowerCapRange [get]
func DevPowerCapRange(dvInd int) (powerMax, powerMin int64, err error) {

	return rsmiDevPowerCapRangeGet(dvInd, 0)
}

// DevPowerCapSet 设置设备功率上限
// @Summary 设置设备功率上限
// @Description 设置指定设备的功率上限，单位微瓦。有效范围可通过 /DevPowerCapRange 查询。
// @Produce json
// @Param dvInd path int true "设备索引"
// @Param cap body int true "功率上限（微瓦）"
// @Success 200 {string} string "操作成功"
// @Failure 400 {object} error "请求失败"
// @Router /powercap/{dvInd} [post]
func DevPowerCapSet(dvInd int, cap uint64) error {
	return rsmiDevPowerCapSet(dvInd, 0, cap)
}

// DevPowerCapReset 复位设备功率上限至默认值
// @Summary 复位设备功率上限
// @Description 读取默认功率上限后重新写入，恢复出厂默认值。
// @Produce json
// @Param dvInd path int true "设备索引"
// @Success 200 {string} string "操作成功"
// @Failure 400 {object} error "请求失败"
// @Router /powercap/{dvInd}/reset [post]
func DevPowerCapReset(dvInd int) error {
	defaultCap, err := rsmiDevPowerCapDefaultGet(dvInd)
	if err != nil {
		return err
	}
	return rsmiDevPowerCapSet(dvInd, 0, defaultCap)
}

// DevPerfLevelSetV1 使用新版接口设置设备性能等级
// @Summary 设置设备性能等级（v1）
// @Description 使用 rsmi_dev_perf_level_set_v1 设置设备性能等级，推荐使用此新版接口。
// @Produce json
// @Param dvInd path int true "设备索引"
// @Param level query string true "性能等级 (auto, low, high, normal, manual, stable_std, stable_peak, stable_min_mclk, stable_min_sclk)"
// @Success 200 {string} string "操作成功"
// @Failure 400 {object} error "请求失败"
// @Router /perflevel/v1/{dvInd} [post]
func DevPerfLevelSetV1(dvInd int, level DevPerfLevel) error {
	return rsmiDevPerfLevelSetV1(dvInd, level)
}

// @Summary 获取设备监控中的指标
// @Description 收集所有设备的监控指标信息。
// @Produce json
// @Success 200 {array} MonitorInfo "设备监控指标信息列表"
// @Failure 400 {object} error "请求错误"
// @Failure 404 {object} error "设备未找到"
// @Router /CollectDeviceMetrics [get]
func CollectDeviceMetrics() (monitorInfos []MonitorInfo, err error) {
	numMonitorDevices, err := rsmiNumMonitorDevices()
	if err != nil {
		return nil, err
	}
	var wg sync.WaitGroup
	monitorInfos = make([]MonitorInfo, numMonitorDevices)
	deviceResults := make(chan MonitorInfo, numMonitorDevices) // Create a channel to collect results

	for i := 0; i < numMonitorDevices; i++ {
		wg.Add(1)
		go func(deviceIndex int) {
			defer wg.Done()

			var wgDevice sync.WaitGroup
			var muDevice sync.Mutex
			monitorInfo := MonitorInfo{MinorNumber: deviceIndex}

			// Collect PCI ID
			wgDevice.Add(1)
			go func() {
				defer wgDevice.Done()
				bdfid, err := rsmiDevPciIdGet(deviceIndex)
				if err != nil {
					glog.Errorf("Failed to get PCI ID for device %d: %v", deviceIndex, err)
					return
				}
				domain := (bdfid >> 32) & 0xffffffff
				bus := (bdfid >> 8) & 0xff
				dev := (bdfid >> 3) & 0x1f
				function := bdfid & 0x7
				pciBusNumber := fmt.Sprintf("%04x:%02x:%02x.%x", domain, bus, dev, function)
				muDevice.Lock()
				monitorInfo.PciBusNumber = pciBusNumber
				muDevice.Unlock()
			}()

			// Collect Device Serial Number
			wgDevice.Add(1)
			go func() {
				defer wgDevice.Done()
				deviceId, _ := rsmiDevSerialNumberGet(deviceIndex)
				muDevice.Lock()
				monitorInfo.DeviceId = deviceId
				muDevice.Unlock()
			}()

			// Collect Device Type ID
			wgDevice.Add(1)
			go func() {
				defer wgDevice.Done()
				//devTypeId, _ := rsmiDevIdGet(deviceIndex)
				info, _ := GetDeviceInfo(deviceIndex)
				devTypeName := NormalizeDevTypeName(info.Name)
				muDevice.Lock()
				monitorInfo.SubSystemName = devTypeName
				muDevice.Unlock()
			}()

			// Collect Temperature
			wgDevice.Add(1)
			go func() {
				defer wgDevice.Done()
				temperature, _ := rsmiDevTempMetricGet(deviceIndex, 0, RSMI_TEMP_CURRENT)
				t, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(temperature)/1000.0), 64)
				muDevice.Lock()
				monitorInfo.Temperature = t
				muDevice.Unlock()
			}()

			// Collect Power Usage
			wgDevice.Add(1)
			go func() {
				defer wgDevice.Done()
				powerUsage, _ := rsmiDevPowerAveGet(deviceIndex, 0)
				pu, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(powerUsage)/1000000.0), 64)
				muDevice.Lock()
				monitorInfo.PowerUsage = pu
				muDevice.Unlock()
			}()

			// Collect Power Cap
			wgDevice.Add(1)
			go func() {
				defer wgDevice.Done()
				powerCap, _ := rsmiDevPowerCapGet(deviceIndex, 0)
				pc, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(powerCap)/1000000.0), 64)
				muDevice.Lock()
				monitorInfo.PowerCap = pc
				muDevice.Unlock()
			}()

			// Collect Memory Capacity
			wgDevice.Add(1)
			go func() {
				defer wgDevice.Done()
				memoryCap, _ := rsmiDevMemoryTotalGet(deviceIndex, RSMI_MEM_TYPE_FIRST)
				mc, _ := strconv.ParseFloat(fmt.Sprintf("%f", float64(memoryCap)/1.0), 64)
				muDevice.Lock()
				monitorInfo.MemoryCap = mc
				muDevice.Unlock()
			}()

			// Collect Memory Usage
			wgDevice.Add(1)
			go func() {
				defer wgDevice.Done()
				memoryUsed, _ := rsmiDevMemoryUsageGet(deviceIndex, RSMI_MEM_TYPE_FIRST)
				mu, _ := strconv.ParseFloat(fmt.Sprintf("%f", float64(memoryUsed)/1.0), 64)
				muDevice.Lock()
				monitorInfo.MemoryUsed = mu
				muDevice.Unlock()
			}()

			// Collect Utilization Rate
			wgDevice.Add(1)
			go func() {
				defer wgDevice.Done()
				utilizationRate, _ := rsmiDevBusyPercentGet(deviceIndex)
				ur, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(utilizationRate)/1.0), 64)
				muDevice.Lock()
				monitorInfo.UtilizationRate = ur
				muDevice.Unlock()
			}()

			// Collect PCIe Throughput
			wgDevice.Add(1)
			go func() {
				defer wgDevice.Done()
				sent, received, maxPktSz, _ := rsmiDevPciThroughputGet(deviceIndex)
				_, _, pcieBwMb := pcieThroughputMB(sent, received, maxPktSz)
				muDevice.Lock()
				monitorInfo.PcieBwMb = pcieBwMb
				muDevice.Unlock()
			}()

			// Collect GPU Clock Frequencies
			wgDevice.Add(1)
			go func() {
				defer wgDevice.Done()
				clk, _ := rsmiDevGpuClkFreqGet(deviceIndex, RSMI_CLK_TYPE_SYS)
				sclk, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(clk.Frequency[clk.Current])/1000000.0), 64)
				supported := clk.NumSupported
				var sclkFrequency []string
				for i := 0; i < int(supported); i++ {
					freq := fmt.Sprintf("%d", int(clk.Frequency[i]/1000000))
					sclkFrequency = append(sclkFrequency, freq)
				}
				muDevice.Lock()
				monitorInfo.Clk = sclk
				monitorInfo.SclkFrequency = sclkFrequency
				muDevice.Unlock()
			}()

			wgDevice.Add(1)
			go func() {
				defer wgDevice.Done()
				soc, _ := rsmiDevGpuClkFreqGet(deviceIndex, RSMI_CLK_TYPE_SOC)
				socclk, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(soc.Frequency[soc.Current])/1000000.0), 64)
				supported := soc.NumSupported
				var socclkFrequency []string
				for i := 0; i < int(supported); i++ {
					freq := fmt.Sprintf("%d", int(soc.Frequency[i]/1000000))
					socclkFrequency = append(socclkFrequency, freq)
				}

				muDevice.Lock()
				monitorInfo.Socclk = socclk
				monitorInfo.SocclkFrequency = socclkFrequency
				muDevice.Unlock()
			}()

			// Collect Performance Level
			wgDevice.Add(1)
			go func() {
				defer wgDevice.Done()
				perf, err := PerfLevel(deviceIndex)
				if err != nil {
					glog.Errorf("Failed to get performance level for device %d: %v", deviceIndex, err)
					return
				}
				muDevice.Lock()
				monitorInfo.PerfLevel = perf
				muDevice.Unlock()
			}()

			wgDevice.Wait()

			deviceResults <- monitorInfo // Send result to channel
		}(i)
	}

	// Close the channel once all Goroutines are done
	go func() {
		wg.Wait()
		close(deviceResults)
	}()

	// Collect results from channel
	for monitorInfo := range deviceResults {
		monitorInfos[monitorInfo.MinorNumber] = monitorInfo
	}

	glog.V(5).Infof("monitorInfos: %", monitorInfos)
	return
}

/*func CollectVDeviceMetrics() (devices []PhysicalDeviceInfo, err error) {

}*/

func DevGpuClkFreqSet(dvInd int, clkType RSMIClkType, freqBitmask int64) (err error) {
	return rsmiDevGpuClkFreqSet(dvInd, clkType, freqBitmask)
}

// GetDeviceByDvInd 根据设备的 dvInd 获取物理设备信息
// @Summary 获取物理设备信息
// @Description 根据设备的 dvInd 获取物理设备信息
// @Tags Device
// @Param dvInd path int true "设备的 MinorNumber"
// @Success 200 {object} PhysicalDeviceInfo "返回物理设备信息"
// @Failure 404 {string} string "设备未找到"
// @Failure 500 {string} string "内部服务器错误"
// @Router /GetDeviceByDvInd [get]
func GetDeviceByDvInd(dvInd int) (physicalDeviceInfo PhysicalDeviceInfo, err error) {
	// 检查设备索引是否有效
	numDevices, err := rsmiNumMonitorDevices()
	if err != nil {
		return physicalDeviceInfo, err
	}

	if dvInd < 0 || dvInd >= numDevices {
		return physicalDeviceInfo, fmt.Errorf("无效的设备索引: %d, 有效范围: 0-%d", dvInd, numDevices-1)
	}

	//物理设备支持最大虚拟化设备数量
	maxVDeviceCount, _ := dmiGetMaxVDeviceCount()
	//物理设备使用百分比
	devPercent, _ := dmiGetDevBusyPercent(dvInd)

	bdfid, err := rsmiDevPciIdGet(dvInd)
	if err != nil {
		return physicalDeviceInfo, err
	}
	// 解析BDFID
	domain := (bdfid >> 32) & 0xffffffff
	bus := (bdfid >> 8) & 0xff
	dev := (bdfid >> 3) & 0x1f
	function := bdfid & 0x7
	// 格式化PCI ID
	pciBusNumber := fmt.Sprintf("%04x:%02x:%02x.%x", domain, bus, dev, function)
	//设备序列号
	deviceId, _ := rsmiDevSerialNumberGet(dvInd)
	info, _ := GetDeviceInfo(dvInd)
	devTypeId, _ := DevTypeID(dvInd)
	//型号名称
	devTypeName := NormalizeDevTypeName(info.Name)
	//获取设备子系统名称
	subsystemTypeId, _ := DevSubsystemId(dvInd)
	subsystemTypeName := NormalizeDevTypeName(info.Name)
	//设备温度
	temperature, _ := rsmiDevTempMetricGet(dvInd, 0, RSMI_TEMP_CURRENT)
	t, err := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(temperature)/1000.0), 64)
	if err != nil {
		return physicalDeviceInfo, err
	}
	//设备平均功耗
	powerUsage, _ := rsmiDevPowerAveGet(dvInd, 0)
	pu, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(powerUsage)/1000000.0), 64)
	//获取设备功率上限
	powerCap, _ := rsmiDevPowerCapGet(dvInd, 0)
	pc, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(powerCap)/1000000.0), 64)
	//获取设备内存总量
	memoryCap, _ := rsmiDevMemoryTotalGet(dvInd, RSMI_MEM_TYPE_FIRST)
	mc, _ := strconv.ParseFloat(fmt.Sprintf("%f", float64(memoryCap)/1.0), 64)
	//获取设备内存使用量
	memoryUsed, _ := rsmiDevMemoryUsageGet(dvInd, RSMI_MEM_TYPE_FIRST)
	mu, _ := strconv.ParseFloat(fmt.Sprintf("%f", float64(memoryUsed)/1.0), 64)
	//获取设备设备忙碌时间百分比
	utilizationRate, _ := rsmiDevBusyPercentGet(dvInd)
	ur, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(utilizationRate)/1.0), 64)
	//获取pcie流量信息
	sent, received, maxPktSz, _ := rsmiDevPciThroughputGet(dvInd)
	pcieSent, pcieReceived, pcieBwMb := pcieThroughputMB(sent, received, maxPktSz)
	//获取设备系统时钟速度列表
	clk, _ := rsmiDevGpuClkFreqGet(dvInd, RSMI_CLK_TYPE_SYS)
	sclk, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(clk.Frequency[clk.Current])/1000000.0), 64)

	computeUnit := float64(info.ComputeUnitCount)
	blockInfos, _ := EccBlocksInfo(dvInd)
	cus, memories, _ := DeviceRemainingInfo(dvInd)
	dfBandwidthInfo, _ := DFBandwidth(dvInd, RSMI_DF_BW_TYPE_ALL)

	device := Device{
		DvInd:                     dvInd,
		PciBusNumber:              pciBusNumber,
		DeviceId:                  deviceId,
		DevTypeId:                 devTypeId,
		DevTypeName:               devTypeName,
		SubsystemTypeId:           subsystemTypeId,
		SubsystemTypeName:         subsystemTypeName,
		Temperature:               t,
		PowerUsed:                 pu,
		PowerTotal:                pc,
		MemoryTotal:               mc,
		MemoryUsed:                mu,
		UtilizationRate:           ur,
		PcieSent:                  pcieSent,
		PcieReceived:              pcieReceived,
		PcieBwMb:                  pcieBwMb,
		Clk:                       sclk,
		ComputeUnitCount:          computeUnit,
		MaxVDeviceCount:           maxVDeviceCount,
		Percent:                   devPercent,
		VDeviceCount:              0,
		ComputeUnitRemainingCount: cus,
		MemoryRemaining:           memories,
		BlocksInfos:               blockInfos,
		DFBandwidthInfo:           dfBandwidthInfo,
	}

	// 创建PhysicalDeviceInfo
	physicalDeviceInfo = PhysicalDeviceInfo{
		Device:         device,
		VirtualDevices: []VDeviceInfo{},
	}

	//获取虚拟设备数量
	vDeviceCount := numDevices * 4

	// 获取所有虚拟设备信息并关联到对应的物理设备
	for j := 0; j < vDeviceCount; j++ {
		vDeviceInfo, err := dmiGetVDeviceInfo(j)
		if err == nil && vDeviceInfo.DvInd == dvInd {
			vDevPercent, _ := dmiGetVDevBusyPercent(j)
			vDeviceInfo.VPercent = vDevPercent
			vDeviceInfo.VdvInd = j
			// 更新虚拟设备的 PciBusNumber，使用物理设备的 pciBusNumber
			vDeviceInfo.PciBusNumber = pciBusNumber
			//设备的类型名称
			vDeviceInfo.Name = devTypeName
			// 将虚拟设备添加到物理设备的 VirtualDevices 列表中
			physicalDeviceInfo.VirtualDevices = append(physicalDeviceInfo.VirtualDevices, vDeviceInfo)
		}
	}

	// 更新物理设备的 VDeviceCount，等于当前虚拟设备的数量
	physicalDeviceInfo.Device.VDeviceCount = len(physicalDeviceInfo.VirtualDevices)

	glog.V(5).Infof("physicalDevice:%v", physicalDeviceInfo)
	return physicalDeviceInfo, nil
}

func VDeviceByDvInd(dvInd int) (vDeviceCount int, vDevInds []int, err error) {
	// 检查物理设备索引是否有效
	numDevices, err := rsmiNumMonitorDevices()
	if err != nil {
		return 0, nil, err
	}
	if dvInd < 0 || dvInd >= numDevices {
		return 0, nil, fmt.Errorf("无效的物理设备索引: %d, 有效范围: 0-%d", dvInd, numDevices-1)
	}
	// 假设每个物理设备最多有4个虚拟设备，这与DeviceInfos函数中的假设一致
	maxVDeviceCount := numDevices * 4

	// 遍历所有可能的虚拟设备，查找属于指定物理设备的虚拟设备
	for j := 0; j < maxVDeviceCount; j++ {
		vDeviceInfo, err := dmiGetVDeviceInfo(j)
		if err == nil && vDeviceInfo.DvInd == dvInd {
			// 找到一个属于指定物理设备的虚拟设备
			vDevInds = append(vDevInds, j)
			glog.V(5).Infof("找到虚拟设备 %d，对应物理设备 %d", j, dvInd)
		}
	}
	// 返回虚拟设备数量和索引列表
	vDeviceCount = len(vDevInds)
	glog.V(5).Infof("物理设备 %d 上的虚拟设备数量: %d, 索引: %v", dvInd, vDeviceCount, vDevInds)
	return vDeviceCount, vDevInds, nil
}

func AllDeviceInfos() ([]PhysicalDeviceInfo, error) {
	var allDevices []PhysicalDeviceInfo
	// 获取物理设备数量
	deviceCount, err := rsmiNumMonitorDevices()
	if err != nil {
		return nil, err
	}

	// 用于保存所有物理设备的信息
	deviceMap := make(map[int]*PhysicalDeviceInfo)

	// 获取所有物理设备信息
	for i := 0; i < deviceCount; i++ {
		//物理设备支持最大虚拟化设备数量
		maxVDeviceCount, _ := dmiGetMaxVDeviceCount()
		//物理设备使用百分比
		devPercent, _ := dmiGetDevBusyPercent(i)

		bdfid, err := rsmiDevPciIdGet(i)
		if err != nil {
			return nil, err
		}
		// 解析BDFID
		domain := (bdfid >> 32) & 0xffffffff
		bus := (bdfid >> 8) & 0xff
		dev := (bdfid >> 3) & 0x1f
		function := bdfid & 0x7
		// 格式化PCI ID
		pciBusNumber := fmt.Sprintf("%04x:%02x:%02x.%x", domain, bus, dev, function)
		//设备序列号
		deviceId, _ := rsmiDevSerialNumberGet(i)
		info, _ := GetDeviceInfo(i)
		//获取设备类型标识id
		devTypeId, _ := DevTypeID(i)
		//型号名称
		devTypeName := NormalizeDevTypeName(info.Name)
		//获取设备子系统名称
		subsystemTypeId, _ := DevSubsystemId(i)
		subsystemTypeName := NormalizeDevTypeName(info.Name)
		//设备温度
		temperature, _ := rsmiDevTempMetricGet(i, 0, RSMI_TEMP_CURRENT)
		t, err := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(temperature)/1000.0), 64)
		if err != nil {
			return nil, err
		}
		//glog.V(5).Infof("HCU[%v] temperature cap:%v ",i,t)
		//设备平均功耗
		powerUsage, _ := rsmiDevPowerAveGet(i, 0)
		pu, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(powerUsage)/1000000.0), 64)
		//glog.V(5).Infof(" HCU[%v] power usage : %.0f", i, pu)
		//获取设备功率上限
		powerCap, _ := rsmiDevPowerCapGet(i, 0)
		pc, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(powerCap)/1000000.0), 64)
		//glog.V(5).Infof("🔋 HCU[%v] power cap : %.0f", i, pc)
		//获取设备内存总量
		memoryCap, _ := rsmiDevMemoryTotalGet(i, RSMI_MEM_TYPE_FIRST)
		mc, _ := strconv.ParseFloat(fmt.Sprintf("%f", float64(memoryCap)/1.0), 64)
		//glog.V(5).Infof("HCU[%v] memory total: %.0f", i, mc)
		//获取设备内存使用量
		memoryUsed, _ := rsmiDevMemoryUsageGet(i, RSMI_MEM_TYPE_FIRST)
		mu, _ := strconv.ParseFloat(fmt.Sprintf("%f", float64(memoryUsed)/1.0), 64)
		//glog.V(5).Infof(" HCU[%v] memory used : %.0f ", i, mu)
		//获取设备设备忙碌时间百分比
		utilizationRate, _ := rsmiDevBusyPercentGet(i)
		ur, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(utilizationRate)/1.0), 64)
		//glog.V(5).Infof(" HCU[%v] utilization rate : %.0f", i, ur)
		//获取pcie流量信息
		sent, received, maxPktSz, _ := rsmiDevPciThroughputGet(i)
		pcieSent, pcieReceived, pcieBwMb := pcieThroughputMB(sent, received, maxPktSz)
		//glog.V(5).Infof(" HCU[%v] PCIE  bandwidth : %.0f", i, pcieBwMb)
		//获取设备系统时钟速度列表
		clk, _ := rsmiDevGpuClkFreqGet(i, RSMI_CLK_TYPE_SYS)
		sclk, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(clk.Frequency[clk.Current])/1000000.0), 64)
		//glog.V(5).Infof(" HCU[%v] SCLK : %.0f", i, sclk)
		computeUnit := float64(info.ComputeUnitCount)
		blockInfos, err := EccBlocksInfo(i)
		cus, memories, _ := DeviceRemainingInfo(i)
		//dfBandwidthInfo, err := DFBandwidth(i, RSMI_DF_BW_TYPE_ALL)
		device := Device{
			DvInd:                     i,
			PciBusNumber:              pciBusNumber,
			DeviceId:                  deviceId,
			DevTypeId:                 devTypeId,
			DevTypeName:               devTypeName,
			SubsystemTypeId:           subsystemTypeId,
			SubsystemTypeName:         subsystemTypeName,
			Temperature:               t,
			PowerUsed:                 pu,
			PowerTotal:                pc,
			MemoryTotal:               mc,
			MemoryUsed:                mu,
			UtilizationRate:           ur,
			PcieSent:                  pcieSent,
			PcieReceived:              pcieReceived,
			PcieBwMb:                  pcieBwMb,
			Clk:                       sclk,
			ComputeUnitCount:          computeUnit,
			MaxVDeviceCount:           maxVDeviceCount,
			Percent:                   devPercent,
			VDeviceCount:              0,
			ComputeUnitRemainingCount: cus,
			MemoryRemaining:           memories,
			BlocksInfos:               blockInfos,
			//DFBandwidthInfo:           dfBandwidthInfo,
		} // 创建PhysicalDeviceInfo并存入map
		pdi := PhysicalDeviceInfo{
			Device:         device,
			VirtualDevices: []VDeviceInfo{},
		}
		deviceMap[device.DvInd] = &pdi
	}

	// 获取虚拟设备数量
	//vDeviceCount, err := dmiGetVDeviceCount()
	vDeviceCount := deviceCount * 4
	if err != nil {
		return nil, err
	}
	// 获取所有虚拟设备信息并关联到对应的物理设备
	for j := 0; j < vDeviceCount; j++ {
		vDeviceInfo, err := dmiGetVDeviceInfo(j)
		glog.V(5).Infof("vDeviceInfo warning: %v", err)
		if err == nil {
			vDevPercent, _ := dmiGetVDevBusyPercent(j)
			vDeviceInfo.VPercent = vDevPercent
			vDeviceInfo.VdvInd = j
			// 找到对应的物理设备并将虚拟设备添加到其VirtualDevices中
			if pdi, exists := deviceMap[vDeviceInfo.DvInd]; exists {
				// 更新虚拟设备的 PciBusNumber，使用物理设备的 pciBusNumber
				vDeviceInfo.PciBusNumber = pdi.Device.PciBusNumber
				//设备类型的名称
				vDeviceInfo.Name = pdi.Device.DevTypeName
				// 将虚拟设备添加到物理设备的 VirtualDevices 列表中
				pdi.VirtualDevices = append(pdi.VirtualDevices, vDeviceInfo)
				// 更新物理设备的 VDeviceCount，等于当前虚拟设备的数量
				pdi.Device.VDeviceCount = len(pdi.VirtualDevices)
			}
		}
		if err != nil {
			glog.Errorf("Error getting virtual device info for virtual device %d: %s", j, err)
		}
	}

	// 将map中的所有PhysicalDeviceInfo转为slice
	for _, pdi := range deviceMap {
		allDevices = append(allDevices, *pdi)
	}
	glog.V(5).Infof("allDevices:%v", (allDevices))
	return allDevices, nil
}

// PicBusInfo 获取设备的总线信息
// @Summary 获取设备的总线信息
// @Description 根据设备索引返回对应的总线信息（BDF格式）
// @Produce json
// @Param dvInd query int true "设备索引"
// @Success 200 {string} string "返回设备的总线信息"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /PicBusInfo [get]
func PicBusInfo(dvInd int) (picID string, err error) {
	bdfid, err := rsmiDevPciIdGet(dvInd)
	if err != nil {
		return "", err
	}
	// Parse BDFID
	domain := (bdfid >> 32) & 0xffffffff
	bus := (bdfid >> 8) & 0xff
	devID := (bdfid >> 3) & 0x1f
	function := bdfid & 0x7
	// Format and return the bus identifier
	picID = fmt.Sprintf("%04x:%02x:%02x.%x", domain, bus, devID, function)
	return
}

// FanSpeedInfo 获取风扇转速信息
// @Summary 获取风扇转速信息
// @Description 根据设备索引返回当前风扇转速及其占最大转速的百分比
// @Produce json
// @Param dvInd query int true "设备索引"
// @Success 200 {int64} fanLevel "返回当前风扇转速"
// @Success 200 {float64} fanPercentage "返回风扇转速百分比"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /FanSpeedInfo [get]
func FanSpeedInfo(dvInd int) (fanLevel int64, fanPercentage float64, err error) {
	// 当前转速
	fanLevel, err = rsmiDevFanSpeedGet(dvInd, 0)
	if err != nil {
		return 0, 0, err
	}
	// 最大转速
	fanMax, err := rsmiDevFanSpeedMaxGet(dvInd, 0)
	if err != nil {
		return 0, 0, err
	}
	// Calculate fan speed percentage
	fanPercentage = (float64(fanLevel) / float64(fanMax)) * 100
	return
}

// HCUUse 当前HCU使用的百分比
// @Summary 获取当前HCU使用的百分比
// @Description 根据设备索引返回当前HCU的使用百分比
// @Produce json
// @Param dvInd query int true "设备索引"
// @Success 200 {int} percent "返回HCU使用的百分比"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /HCUUse [get]
func HCUUse(dvInd int) (percent int, err error) {
	percent, err = dmiGetDevBusyPercent(dvInd)
	if err != nil {
		return 0, err
	}
	return
}

// DevID 设备ID的十六进制值
// @Summary 获取设备ID的十六进制值
// @Description 根据设备索引返回设备ID的十六进制值
// @Produce json
// @Param dvInd query int true "设备索引"
// @Success 200 {string} id "返回设备ID的十六进制值"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /DevID [get]
func DevTypeID(dvInd int) (id string, err error) {
	devId, err := rsmiDevIdGet(dvInd)
	id = fmt.Sprintf("%x", devId)
	glog.V(5).Infof("DevID:%v", id)
	return
}

// DevName 设备类型名称
// @Summary 获取设备类型名称
// @Description 根据设备索引返回设备类型名称
// @Produce json
// @Param dvInd query int true "设备索引"
// @Success 200 {string} id "返回设备类型名称"
// @Success 200 {float64} id "返回设备CU数量"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /DevTypeName [get]
func DevTypeName(dvInd int) (devTypeName string, computeUnit float64, err error) {
	//型号名称
	info, err := GetDeviceInfo(dvInd)
	if err != nil {
		return
	}
	devTypeName = NormalizeDevTypeName(info.Name)
	computeUnit = float64(info.ComputeUnitCount)
	return
}

func DevSubsystemId(dvInd int) (subsystemId string, err error) {
	id, err := rsmiDevSubsystemIdGet(dvInd)
	subsystemId = fmt.Sprintf("%x", id)
	glog.V(5).Infof("DevSubsystemId:%v", subsystemId)
	return
}

func DevSubsystemName(dvInd int) (name string, err error) {
	name, err = rsmiDevSubsystemNameGet(dvInd)
	return
}

// MaxPower 设备的最大功率
// @Summary 获取设备的最大功率
// @Description 根据设备索引返回设备的最大功率（以瓦特为单位）
// @Produce json
// @Param dvInd query int true "设备索引"
// @Success 200 {int64} powerMax "返回设备的最大功率"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /MaxPower [get]
func MaxPower(dvInd int) (powerMax int64, err error) {
	powerMax, err = rsmiDevPowerCapGet(dvInd, 0)
	if err != nil {
		return 0, err
	}
	glog.V(5).Infof("Max power: %v", (powerMax / 1000000))
	return (powerMax / 1000000), nil
}

// MemInfo 获取设备的指定内存使用情况
// @Summary 获取设备的指定内存使用情况
// @Description 根据设备索引和内存类型返回内存的使用量和总量
// @Produce json
// @Param dvInd query int true "设备索引"
// @Param memType query string true "内存类型（可选值: vram, vis_vram, gtt）"
// @Success 200 {int64} memUsed "返回指定内存类型的使用量"
// @Success 200 {int64} memTotal "返回指定内存类型的总量"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /MemInfo [get]
func MemInfo(dvInd int, memType string) (memUsed int64, memTotal int64, err error) {
	memType = strings.ToUpper(memType)
	if !contains(memoryTypeL, memType) {
		//fmt.Println(dvInd, fmt.Sprintf("Invalid memory type %s", memType))
		return 0, 0, fmt.Errorf("invalid memory type")
	}
	memTypeIndex := RSMIMemoryType(indexOf(memoryTypeL, memType))
	memUsed, err = rsmiDevMemoryUsageGet(dvInd, memTypeIndex)
	if err != nil {
		return memUsed, memTotal, err
	}
	//fmt.Println(dvInd, fmt.Sprintf("memUsed: %d", memUsed))
	memTotal, err = rsmiDevMemoryTotalGet(dvInd, memTypeIndex)
	if err != nil {
		return memUsed, memTotal, err
	}
	//fmt.Println(dvInd, fmt.Sprintf("memTotal: %d", memTotal))
	return
}

func HCUClk(dvInd int) (clk float64, err error) {
	hcuClk, _ := rsmiDevGpuClkFreqGet(dvInd, RSMI_CLK_TYPE_SYS)
	clk, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", float64(hcuClk.Frequency[hcuClk.Current])/1000000.0), 64)
	return
}

func HCUMclk(dvInd int) (clk float64, err error) {
	mclk, _ := rsmiDevGpuClkFreqGet(dvInd, RSMI_CLK_TYPE_MEM)
	clk, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", float64(mclk.Frequency[mclk.Current])/1000000.0), 64)
	return
}

// DeviceInfos 获取设备信息列表
// @Summary 获取设备信息列表
// @Description 返回所有设备的详细信息列表
// @Produce json
// @Success 200 {array} DeviceInfo "返回设备信息列表"
// @Failure 500 {object} error "服务器内部错误"
// @Router /DeviceInfos [get]
func DeviceInfos() (deviceInfos []DeviceInfo, err error) {
	numDevices, err := rsmiNumMonitorDevices()
	if err != nil {
		return nil, err
	}

	// 初始化设备信息数组
	for i := 0; i < numDevices; i++ {
		bdfid, err := rsmiDevPciIdGet(i)
		if err != nil {
			return nil, err
		}
		// 解析BDFID
		domain := (bdfid >> 32) & 0xffffffff
		bus := (bdfid >> 8) & 0xff
		dev := (bdfid >> 3) & 0x1f
		function := bdfid & 0x7
		// 格式化PCI ID
		pciBusNumber := fmt.Sprintf("%04x:%02x:%02x.%x", domain, bus, dev, function)
		//设备序列号
		deviceId, _ := rsmiDevSerialNumberGet(i)
		//获取设备类型标识id
		//devTypeId, _ := rsmiDevIdGet(i)
		//devType := fmt.Sprintf("%x", devTypeId)
		devType, _ := DevTypeID(i)
		//型号名称
		//devTypeName := type2name[devType]
		//devTypeName := NormalizeCardSeriesName(type2name[devType])
		info, _ := GetDeviceInfo(i)
		devTypeName := NormalizeDevTypeName(info.Name)
		//获取设备子系统名称
		subsystemTypeId, _ := DevSubsystemId(i)
		//subsystemTypeName := type2name[subsystemTypeId]
		//subsystemTypeName := NormalizeCardSeriesName(type2name[subsystemTypeId])
		subsystemTypeName := NormalizeDevTypeName(info.Name)
		//获取设备内存总量
		memoryTotal, _ := rsmiDevMemoryTotalGet(i, RSMI_MEM_TYPE_FIRST)
		mt, _ := strconv.ParseFloat(fmt.Sprintf("%f", float64(memoryTotal)/1.0), 64)
		glog.V(5).Infof("HCU[%v] 内存总量: %.0f", i, mt)
		//获取设备功率上限
		powerTotal, _ := rsmiDevPowerCapGet(i, 0)
		pt, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(powerTotal)/1000000.0), 64)
		computeUnit := float64(info.ComputeUnitCount)
		glog.V(5).Infof("HCU[%v] 计算单元: %.0f", i, computeUnit)

		deviceInfo := DeviceInfo{
			DvInd:             i,
			DeviceId:          deviceId,
			DevType:           devType,
			DevTypeName:       devTypeName,
			SubsystemTypeId:   subsystemTypeId,
			SubsystemTypeName: subsystemTypeName,
			PciBusNumber:      pciBusNumber,
			MemoryTotal:       mt,
			PowerTotal:        pt,
			ComputeUnit:       computeUnit,
			VDeviceCount:      0, // 初始化虚拟设备计数为0
		}
		deviceInfos = append(deviceInfos, deviceInfo)
	}

	// 获取虚拟设备数量并关联到物理设备
	vDeviceCount := numDevices * 4 // 假设每个物理设备最多有4个虚拟设备
	glog.V(5).Infof("开始遍历虚拟设备，总数: %d", vDeviceCount)

	// 遍历所有可能的虚拟设备，直接更新对应物理设备的VDeviceCount
	for j := 0; j < vDeviceCount; j++ {
		vDeviceInfo, err := dmiGetVDeviceInfo(j)
		if err == nil {
			// 获取虚拟设备对应的物理设备ID
			physicalDeviceID := vDeviceInfo.DvInd
			glog.V(5).Infof("虚拟设备 %d 信息: DeviceID=%d", j, physicalDeviceID)
			vDeviceInfo.Name = deviceInfos[physicalDeviceID].DevTypeName
			// 检查物理设备ID是否在有效范围内
			if physicalDeviceID >= 0 && physicalDeviceID < len(deviceInfos) {
				// 直接增加对应物理设备的虚拟设备计数
				deviceInfos[physicalDeviceID].VDeviceCount++
				glog.V(5).Infof("找到虚拟设备 %d，对应物理设备 %d，更新虚拟设备计数为 %d",
					j, physicalDeviceID, deviceInfos[physicalDeviceID].VDeviceCount)
			} else {
				glog.Warningf("虚拟设备 %d 的物理设备ID %d 超出范围", j, physicalDeviceID)
			}
		} else {
			// 如果获取虚拟设备信息失败，可能是该索引没有对应的虚拟设备
			glog.V(3).Infof("索引 %d 没有对应的虚拟设备: %v", j, err)
		}
	}

	glog.V(5).Infof("deviceInfos: %s", (deviceInfos))
	return deviceInfos, nil
}

// DeviceStatus 获取设备状态信息
// @Summary 获取设备状态信息
// @Description 返回所有设备设备状态信息
// @Produce json
// @Success 200 {array} DeviceStatusInfo "返回设备设备状态信息"
// @Failure 500 {object} error "服务器内部错误"
// @Router /DeviceStatus [get]
func DeviceStatus() (deviceStatusInfos []DeviceStatusInfo, err error) {
	numDevices, err := rsmiNumMonitorDevices()
	if err != nil {
		return nil, err
	}
	// 初始化设备信息数组
	for i := 0; i < numDevices; i++ {
		//设备温度
		temperature, _ := rsmiDevTempMetricGet(i, 0, RSMI_TEMP_CURRENT)
		t, err := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(temperature)/1000.0), 64)
		if err != nil {
			return nil, err
		}
		//设备平均功耗
		powerUsed, _ := rsmiDevPowerAveGet(i, 0)
		pu, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(powerUsed)/1000000.0), 64)
		//获取设备内存使用量
		memoryUsed, _ := rsmiDevMemoryUsageGet(i, RSMI_MEM_TYPE_FIRST)
		mu, _ := strconv.ParseFloat(fmt.Sprintf("%f", float64(memoryUsed)/1.0), 64)
		//获取设备设备忙碌时间百分比
		utilizationRate, _ := rsmiDevBusyPercentGet(i)
		ur, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(utilizationRate)/1.0), 64)
		//获取pcie流量信息
		sent, received, maxPktSz, _ := rsmiDevPciThroughputGet(i)
		pcieSent, pcieReceived, pcieBwMb := pcieThroughputMB(sent, received, maxPktSz)
		clk, err := rsmiDevGpuClkFreqGet(i, RSMI_CLK_TYPE_SYS)
		var sclk = 0.0
		if err == nil {
			// 类型转换：将 uint32 转换为 int
			currentIndex := int(clk.Current)
			frequencyLength := len(clk.Frequency)
			// 添加安全索引检查（兼容类型）
			if currentIndex >= 0 && currentIndex < frequencyLength {
				freqMHz := float64(clk.Frequency[currentIndex]) / 1_000_000.0
				sclk = math.Round(freqMHz*100) / 100
			} else {
				log.Printf("Invalid clock index: %d (max %d)", currentIndex, frequencyLength-1)
			}
		} else {
			log.Printf("Failed to get clock: %v", err)
		}

		//物理设备使用百分比
		devPercent, _ := dmiGetDevBusyPercent(i)
		blockInfos, err := EccBlocksInfo(i)
		cus, memories, _ := DeviceRemainingInfo(i)
		deviceStatusInfo := DeviceStatusInfo{
			DvInd:                     i,
			Temperature:               t,
			PowerUsed:                 pu,
			MemoryUsed:                mu,
			UtilizationRate:           ur,
			PcieBwMb:                  pcieBwMb,
			PcieSent:                  pcieSent,
			PcieReceived:              pcieReceived,
			Clk:                       sclk,
			Percent:                   devPercent,
			ComputeUnitRemainingCount: cus,
			MemoryRemaining:           memories,
			BlocksInfos:               blockInfos,
		}
		deviceStatusInfos = append(deviceStatusInfos, deviceStatusInfo)
	}
	glog.V(5).Infof("deviceStatusInfos:%v", (deviceStatusInfos))
	return
}

// ProcessName 获取指定PID的进程名
// @Summary 获取指定PID的进程名
// @Description 根据进程ID（PID）返回对应的进程名称
// @Produce json
// @Param pid query int true "进程ID"
// @Success 200 {string} string "返回进程名称"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /ProcessName [get]
func ProcessName(pid int) string {
	if pid < 1 {
		glog.V(5).Infof("PID must be greater than 0")
		return "UNKNOWN"
	}
	cmd := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "comm=")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		glog.V(5).Infof("Error executing command:", err)
		return "UNKNOWN"
	}
	pName := out.String()
	if pName == "" {
		return "UNKNOWN"
	}
	// Remove the substrings surrounding from process name (b' and \n')
	pName = strings.TrimPrefix(pName, "b'")
	pName = strings.TrimSuffix(pName, "\\n'")
	glog.V(5).Infof("Process name: %s\n", pName)
	return strings.TrimSpace(pName)
}

// PerfLevel 获取设备的当前性能水平
// @Summary 获取设备的当前性能水平
// @Description 返回指定设备的当前性能等级
// @Produce json
// @Param dvInd query int true "设备索引"
// @Success 200 {string} string "返回当前性能水平"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /PerfLevel [get]
func PerfLevel(dvInd int) (perf string, err error) {
	level, err := rsmiDevPerfLevelGet(dvInd)
	if err != nil {
		return perf, err
	}
	perf = perfLevelString(int(level))
	glog.V(5).Infof("Perf level: %v", perf)
	return
}

// getPid 获取特定应用程序的进程 ID
func PidByName(name string) (pid string, err error) {
	glog.V(5).Infof("pidName: %s\n", name)
	cmd := exec.Command("pidof", name)
	output, err := cmd.Output()
	glog.V(5).Infof("output:", output)
	if err != nil {
		glog.V(5).Infof("Error: %v\nOutput: %s", err, string(output))
	} else {
		glog.V(5).Infof("Output: %s", string(output))
	}
	// 移除末尾的换行符并返回 PID
	pid = strings.TrimSpace(string(output))
	glog.V(5).Infof("pid: %s\n", pid)
	return
}

// Power 获取设备的平均功耗
// @Summary 获取设备的平均功耗
// @Description 返回指定设备的平均功耗
// @Produce json
// @Param dvInd query int true "设备索引"
// @Success 200 {int64} int64 "返回平均功耗（瓦特）"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /Power [get]
func Power(dvInd int) (power int64, err error) {
	powerAve, err := rsmiDevPowerAveGet(dvInd, 0)
	power = powerAve / 1000000
	glog.V(5).Infof("Power: %v", power)
	if err != nil {
		return power, err
	}
	return
}

// EccStatus 获取GPU块的ECC状态
// @Summary 获取GPU块的ECC状态
// @Description 返回指定GPU块的ECC状态
// @Produce json
// @Param dvInd query int true "设备索引"
// @Param block query string true "GPU块"
// @Success 200 {string} string "返回ECC状态"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /EccStatus [get]
func EccStatus(dvInd int, block RSMIGpuBlock) (state string, err error) {
	eccStatus, err := rsmiDevEccStatusGet(dvInd, block)
	state = rasErrStaleMachine[eccStatus]
	return
}

func EccCount(dvInd int, block RSMIGpuBlock) (errorCount RSMIErrorCount, err error) {
	errorCount, err = rsmiDevEccCountGet(dvInd, block)
	return
}

// EccBlocksInfo 获取ECC块信息
// @Summary 获取ECC块信息
// @Description 根据设备索引返回ECC块的详细信息，包括块类型、状态、CE错误数和UE错误数
// @Produce json
// @Param dvInd query int true "设备索引"
// @Success 200 {array} BlocksInfo "返回包含每个ECC块信息的数组"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /EccBlocksInfo [get]
func EccBlocksInfo(dvInd int) (blocksInfos []BlocksInfo, err error) {
	// 定义所有的RSMIGpuBlock值
	blocks := []RSMIGpuBlock{
		RSMIGpuBlockATHUB,
		RSMIGpuBlockDF,
		RSMIGpuBlockFuse,
		RSMIGpuBlockGFX,
		RSMIGpuBlockHDP,
		RSMIGpuBlockMMHUB,
		RSMIGpuBlockMP0,
		RSMIGpuBlockMP1,
		RSMIGpuBlockPCIEBIF,
		RSMIGpuBlockSDMA,
		RSMIGpuBlockSEM,
		RSMIGpuBlockSMN,
		RSMIGpuBlockUMC,
		RSMIGpuBlockHSLWAFL,
	}

	// 遍历所有的block，分别调用EccStatus和EccCount
	for _, block := range blocks {
		state, err := EccStatus(dvInd, block)
		if err != nil {
			glog.Errorf("EccStatus 调用错误: block: %v, 错误: %v\n", block, err)
			continue
		}
		//glog.V(5).Infof("EccStatus - block: %v, state: %v\n", block, state)

		// 当状态是“ENABLED”时，调用EccCount接口获取错误计数
		if state == "ENABLED" {
			errorCount, err := EccCount(dvInd, block)
			if err != nil {
				glog.Errorf("EccCount 调用错误: block: %v, 错误: %v\n", block, err)
				continue
			}
			//glog.V(5).Infof("EccCount - block: %v, CorrectableErr: %v, UncorrectableErr: %v\n", block, errorCount.CorrectableErr, errorCount.UncorrectableErr)
			// 将block信息添加到结果集中
			blocksInfos = append(blocksInfos, BlocksInfo{
				Block: ConvertFromRSMIGpuBlock(block),
				State: state,
				CE:    int64(errorCount.CorrectableErr),
				UE:    int64(errorCount.UncorrectableErr),
			})
		} else {
			// 状态不是ENABLED时，只添加状态信息，不获取错误计数
			blocksInfos = append(blocksInfos, BlocksInfo{
				Block: ConvertFromRSMIGpuBlock(block),
				State: state,
				CE:    0,
				UE:    0,
			})
		}
	}
	//glog.V(5).Infof("blocksInfos:%v", (blocksInfos))
	return
}

func EccEnabled(dvInd int) (enabledBlocks int64, err error) {
	return rsmiDevEccEnabledGet(dvInd)
}

// 设置设备的性能确定性模式(K100 AI不支持)
func PerfDeterminismMode(dvInd int, clkValue int64) (err error) {
	return rsmiPerfDeterminismModeSet(dvInd, clkValue)
}

// Temperature 获取设备温度
// @Summary 获取设备温度
// @Description 返回指定设备的当前温度
// @Produce json
// @Param dvInd query int true "设备索引"

// @Success 200 {float64} float64 "返回温度（摄氏度）"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /Temperature [get]
func Temperature(dvInd int) (temp float64, err error) {
	deviceTemp, err := rsmiDevTempMetricGet(dvInd, 0, RSMI_TEMP_CURRENT)
	temp, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", float64(deviceTemp)/1000.0), 64)
	glog.V(5).Infof("device Temperature:%v", temp)
	return
}

func GetTempByMetric(dvInd int, metric RSMITemperatureMetric) (temp float64, err error) {
	deviceTemp, err := rsmiDevTempMetricGet(dvInd, 0, metric)
	if err != nil {
		return 0, err
	}
	temp, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", float64(deviceTemp)/1000.0), 64)
	glog.V(5).Infof("get device Temperature %v by metric %v", temp, metric)
	return
}

// DevVersion 获取当前运行的RSMI版本
// @Summary 获取当前运行的RSMI版本
// @Description 返回当前设备的RSMI版本信息
// @Produce json
// @Success 200 {object} DevVersion "返回包含RSMI版本信息的对象"
// @Failure 500 {object} error "服务器内部错误"
// @Router /HCUVersion [get]
func HCUVersion() (version DevVersion, err error) {
	return rsmiVersionGet()
}

func DTKVersion() (dtkVersion string, err error) {
	return getDTKVersionByReadFile()
}

// VbiosVersion 获取设备的VBIOS版本
// @Summary 获取设备的VBIOS版本
// @Description 返回指定设备的VBIOS版本
// @Produce json
// @Param dvInd query int true "设备索引"
// @Success 200 {string} string "返回VBIOS版本"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /VbiosVersion [get]
func VbiosVersion(dvInd int) (vbios string, err error) {
	vbios, err = rsmiDevVbiosVersionGet(dvInd, 256)
	glog.V(5).Infof("VbiosVersion:%v", vbios)
	return
}

// Version 获取当前系统的驱动程序版本
// @Summary 获取当前系统的驱动程序版本
// @Description 返回指定组件的驱动程序版本
// @Produce json
// @Param component query string true "驱动组件"
// @Success 200 {string} string "返回驱动程序版本"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /Version [get]
func Version(component SwComponent) (varStr string, err error) {
	varStr, err = rsmiVersionStrGet(component, 256)
	glog.V(5).Infof("component; Version:%v,%v", component, varStr)
	return
}

// 设置设备超速百分比
func DevOverdriveLevelSet(dvInd, od int) (err error) {
	return rsmiDevOverdriveLevelSet(dvInd, od)
}

// 获取设备的超速百分比
func DevOverdriveLevelGet(dvInd int) (od int, err error) {
	return rsmiDevOverdriveLevelGet(dvInd)
}

// ResetClocks 将设备的时钟重置为默认值
// @Summary 重置设备时钟
// @Description 重置指定设备的时钟和性能等级为默认值
// @Produce json
// @Param dvIdList body []int true "设备ID列表"
// @Success 200 {array} FailedMessage "返回失败消息列表"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /ResetClocks [post]
func ResetClocks(dvIdList []int) (failedMessage []FailedMessage) {
	errorMap := make(map[int][]string)
	glog.V(5).Infof(" Reset Clocks ")
	for _, device := range dvIdList {
		// Reset OverDrive
		err := rsmiDevOverdriveLevelSet(device, 0)
		if err != nil {
			errorMap[device] = append(errorMap[device], "Unable to reset OverDrive")
			glog.Errorf("Unable to reset OverDrive, device: %v, error: %v", device, err)
		}
		// Reset PerfLevel
		err = rsmiDevPerfLevelSet(device, RSMI_DEV_PERF_LEVEL_AUTO)
		if err != nil {
			errorMap[device] = append(errorMap[device], "Unable to reset clocks")
			glog.Errorf("Unable to reset clocks, device: %v, error: %v", device, err)
		}

		// Set performance level to auto
		err = rsmiDevPerfLevelSet(device, RSMI_DEV_PERF_LEVEL_AUTO)
		if err != nil {
			errorMap[device] = append(errorMap[device], "Unable to set performance level to auto")
			glog.Errorf("Unable to set performance level to auto, device: %v, error: %v", device, err)
		}
	}
	for id, msg := range errorMap {
		failedMessage = append(failedMessage, FailedMessage{ID: id, ErrorMsg: strings.Join(msg, "; ")})
	}
	return
}

// ResetFans 复位风扇驱动控制
// @Summary 复位风扇控制
// @Description 重置指定设备的风扇控制为默认值
// @Produce json
// @Param dvIdList body []int true "设备ID列表"
// @Success 200 {string} string "复位成功"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /ResetFans [post]
func ResetFans(dvIdList []int) (err error) {
	for _, id := range dvIdList {
		err := rsmiDevFanReset(id, 0)
		glog.V(5).Infof("Resetting fan :%v", id)
		if err != nil {
			glog.Errorf("Unable reset Fan dvId:%v ,err:%v", id, err)
		}
	}
	return
}

// ResetProfile 重置设备的配置文件
// @Summary 重置指定设备的电源配置文件和性能级别
// @Produce json
// @Param dvIdList body []int true "设备ID列表"
// @Success 200 {array} FailedMessage "返回失败的设备及其错误信息"
// @Router /ResetProfile [post]
func ResetProfile(dvIdList []int) (failedMessage []FailedMessage) {
	errorMap := make(map[int][]string)
	for _, id := range dvIdList {
		err := rsmiDevPowerProfileSet(id, 0, RSMI_PWR_PROF_PRST_BOOTUP_DEFAULT)
		if err != nil {
			errorMap[id] = append(errorMap[id], "Unable to reset OverDrive")
			glog.Errorf("Unable to reset OverDrive, device: %v, error: %v", id, err)
		}
		// Reset PerfLevel
		err = rsmiDevPerfLevelSet(id, RSMI_DEV_PERF_LEVEL_AUTO)
		if err != nil {
			errorMap[id] = append(errorMap[id], "Unable to reset PerfLevel")
			glog.Errorf("Unable to reset PerfLevel, device: %v, error: %v", id, err)
		}
	}
	for id, msg := range errorMap {
		failedMessage = append(failedMessage, FailedMessage{ID: id, ErrorMsg: strings.Join(msg, "; ")})
	}
	return
}

// ResetHSLErr 重置设备的HSL错误状态
// @Summary 重置指定设备的HSL错误状态
// @Produce json
// @Param dvIdList body []int true "设备ID列表"
// @Success 200 {array} FailedMessage "返回失败的设备及其错误信息"
// @Router /ResetHSLErr [post]
func ResetHSLErr(dvIdList []int) (failedMessage []FailedMessage) {
	errorMap := make(map[int][]string)
	for _, id := range dvIdList {
		err := rsmiDevHSLErrorReset(id)
		if err != nil {
			errorMap[id] = append(errorMap[id], "Unable to reset HSL error")
			glog.Errorf("Unable to reset HSL error, device: %v, error: %v", id, err)
		}
	}
	for id, msg := range errorMap {
		failedMessage = append(failedMessage, FailedMessage{ID: id, ErrorMsg: strings.Join(msg, "; ")})
	}
	return
}

// HSLErrorStatus 获取HSL错误状态
// @Summary 获取HSL错误状态
// @Description 获取指定物理设备的HSL（高速互连链路）错误状态。
// @Tags HSL状态
// @Param dvInd query int true "物理设备的索引"
// @Success 200 {integer} int "返回HSL错误状态码"
// @Failure 400 {string} string "获取HSL错误状态失败"
// @Router /HSLErrorStatus [get]
func HSLErrorStatus(dvInd int) (status RSMIHSLStatus, err error) {
	return rsmiDevHSLErrorStatus(dvInd)
}

// HSLHiveIdGet 获取设备的HSL hive id
// @Summary 获取指定设备的HSL hive id
// @Produce json
// @Param dvInd query int true "设备索引"
// @Success 200 {int64} int64 "返回设备的HSL hive id"
// @Router /HSLHiveIdGet [get]
func HSLHiveIdGet(dvInd int) (hiveId int64, err error) {
	return rsmiDevHSLHiveIdGet(dvInd)
}

// ResetPerfDeterminism 重置Performance Determinism
// @Summary 重置指定设备的性能决定性设置
// @Produce json
// @Param dvIdList body []int true "设备ID列表"
// @Success 200 {array} FailedMessage "返回失败的设备及其错误信息"
// @Router /ResetPerfDeterminism [post]
func ResetPerfDeterminism(dvIdList []int) (failedMessage []FailedMessage) {
	errorMap := make(map[int][]string)
	for _, device := range dvIdList {
		// Set performance level to auto
		err := rsmiDevPerfLevelSet(device, RSMI_DEV_PERF_LEVEL_AUTO)
		if err != nil {
			errorMap[device] = append(errorMap[device], "Unable to diable performance determinism")
			glog.Errorf("Unable to diable performance determinism, device: %v, error: %v", device, err)
		}
	}
	for id, msg := range errorMap {
		failedMessage = append(failedMessage, FailedMessage{ID: id, ErrorMsg: strings.Join(msg, "; ")})
	}
	return
}

// 为设备选定的时钟类型设定相应的频率范围
func SetClockRange(dvIdList []int, clkType string, minvalue string, maxvalue string, autoRespond bool) (failedMessage []FailedMessage) {
	errorMap := make(map[int][]string)
	if clkType != "sclk" && clkType != "mclk" {
		glog.V(5).Infof("device :%v,Invalid range identifier %v", dvIdList, clkType)
		glog.V(5).Infof("Unsupported range type %s", clkType)
		return
	}
	minVal, errMin := strconv.ParseInt(minvalue, 10, 64)
	maxVal, errMax := strconv.ParseInt(maxvalue, 10, 64)
	if errMin != nil || errMax != nil {
		glog.Errorf("Unable to set %s range", clkType)
		glog.V(5).Infof("%s or %s is not an integer", minvalue, maxvalue)
		return
	}
	confirmOutOfSpecWarning(autoRespond)
	for _, device := range dvIdList {
		err := rsmiDevClkRangeSet(device, minVal, maxVal, rsmiClkNamesDict[clkType])
		if err == nil {
			glog.Errorf("device:%v Successfully set %v from %v(MHz) to %v(MHz)", clkType, minVal, maxVal)
		} else {
			glog.Errorf("device:%v Unable to set %v from %v(MHz) to %v(MHz)", device, clkType, minVal, maxVal)
			errorMap[device] = append(errorMap[device], err.Error())
			glog.Errorf("Unable to diable performance determinism, device: %v, error: %v", device, err)

		}
	}
	for id, msg := range errorMap {
		failedMessage = append(failedMessage, FailedMessage{ID: id, ErrorMsg: strings.Join(msg, "; ")})
	}
	glog.V(5).Infof("SetClockRange failedMessage:%v", failedMessage)
	return
}

// 设置电压曲线
func DevOdVoltInfoSet(dvInd, vPoint, clkValue, voltValue int) (err error) {
	return rsmiDevOdVoltInfoSet(dvInd, vPoint, clkValue, voltValue)
}

// SetPowerPlayTableLevel 设置 PowerPlay 级别
// @Summary 设置设备的 PowerPlay 表级别
// @Description 该函数为设备列表设置 PowerPlay 表级别。它会检查输入值的有效性并相应地调整电压设置。
// @Tags 设备
// @Param dvIdList body []int true "设备 ID 列表"
// @Param clkType query string true "时钟类型（sclk 或 mclk）"
// @Param point query string true "电压点"
// @Param clk query string true "时钟值（以 MHz 为单位）"
// @Param volt query string true "电压值（以 mV 为单位）"
// @Param autoRespond query bool false "自动响应超出规格的警告"
// @Success 200 {string} string "成功设置 PowerPlay 表级别"
// @Failure 400 {string} string "输入无效或无法设置 PowerPlay 表级别"
// @Router /SetPowerPlayTableLevel [post]
func SetPowerPlayTableLevel(dvIdList []int, clkType string, point string, clk string, volt string, autoRespond bool) (failedMessage []FailedMessage) {
	value := fmt.Sprintf("%s %s %s", point, clk, volt)
	_, errPoint := strconv.Atoi(point)
	_, errClk := strconv.Atoi(clk)
	_, errVolt := strconv.Atoi(volt)

	// 创建一个 errorMap 用来记录错误信息
	errorMap := make(map[int][]string)

	if errPoint != nil || errClk != nil || errVolt != nil {
		glog.V(5).Infof("Unable to set PowerPlay table level")
		glog.V(5).Infof("Non-integer characters are present in %s", value)
		// 这里可以返回错误信息
		failedMessage = append(failedMessage, FailedMessage{ID: -1, ErrorMsg: "Invalid non-integer characters in parameters"})
		return
	}

	confirmOutOfSpecWarning(autoRespond)

	for _, device := range dvIdList {
		pointVal, _ := strconv.Atoi(point)
		clkVal, _ := strconv.Atoi(clk)
		voltVal, _ := strconv.Atoi(volt)

		if clkType == "sclk" || clkType == "mclk" {
			err := rsmiDevOdVoltInfoSet(device, pointVal, clkVal, voltVal)
			if err == nil {
				glog.V(5).Infof("device:%v Successfully set voltage point %v to %v(MHz) %v(mV)", device, point, clk, volt)
			} else {
				errorMsg := fmt.Sprintf("Unable to set voltage point %v to %v(MHz) %v(mV)", point, clk, volt)
				glog.Errorf("device:%v %s", device, errorMsg)
				errorMap[device] = append(errorMap[device], errorMsg)
			}
		} else {
			errorMsg := fmt.Sprintf("Unsupported range type %s", clkType)
			glog.Errorf("device:%v Unable to set %s range", device, clkType)
			glog.V(5).Infof("Unsupported range type %s", clkType)
			errorMap[device] = append(errorMap[device], errorMsg)
		}
	}

	// 将 errorMap 转换为 failedMessage 列表
	for id, msg := range errorMap {
		failedMessage = append(failedMessage, FailedMessage{ID: id, ErrorMsg: strings.Join(msg, "; ")})
	}

	return
}

// SetClockOverDrive 设置时钟速度为 OverDrive
// @Summary 为设备设置时钟 OverDrive
// @Description 该函数为设备列表设置时钟 OverDrive 级别。它会调整时钟速度，并在需要时确保性能级别设置为手动模式。
// @Tags 设备
// @Param dvIdList body []int true "设备 ID 列表"
// @Param clktype query string true "时钟类型（sclk 或 mclk）"
// @Param value query string true "OverDrive 值，表示为百分比（0-20%）"
// @Param autoRespond query bool false "自动响应超出规格的警告"
// @Success 200 {string} string "成功设置时钟 OverDrive"
// @Failure 400 {string} string "输入无效或无法设置时钟 OverDrive"
// @Router /SetClockOverDrive [post]
func SetClockOverDrive(dvIdList []int, clktype string, value string, autoRespond bool) (failedMessage []FailedMessage) {
	glog.V(5).Infof("Set Clock OverDrive Range: 0 to 20%")
	intValue, err := strconv.Atoi(value)
	if err != nil {
		glog.V(5).Infof("Unable to set OverDrive level")
		glog.Errorf("%s it is not an integer", value)
		failedMessage = append(failedMessage, FailedMessage{ID: -1, ErrorMsg: "Invalid non-integer value for OverDrive"})
		return
	}

	confirmOutOfSpecWarning(autoRespond)

	for _, device := range dvIdList {
		if intValue < 0 {
			glog.Errorf("Unable to set OverDrive for device: %v", device)
			glog.V(5).Infof("Overdrive cannot be less than 0%")
			failedMessage = append(failedMessage, FailedMessage{ID: device, ErrorMsg: "OverDrive cannot be less than 0%"})
			continue
		}
		if intValue > 20 {
			glog.V(5).Infof("device:%v, Setting OverDrive to 20%%", device)
			glog.V(5).Infof("OverDrive cannot be set to a value greater than 20%")
			intValue = 20
		}
		perf, _ := PerfLevel(device)
		if perf != "MANUAL" {
			err := rsmiDevPerfLevelSet(device, RSMI_DEV_PERF_LEVEL_MANUAL)
			if err == nil {
				glog.V(5).Infof("device:%v Performance level set to manual", device)
			} else {
				glog.Errorf("device:%v Unable to set performance level to manual")
				failedMessage = append(failedMessage, FailedMessage{ID: device, ErrorMsg: err.Error()})
				continue
			}
		}
		if clktype == "mclk" {
			fsFile := fmt.Sprintf("/sys/class/drm/card%d/device/pp_mclk_od", device)
			if _, err := os.Stat(fsFile); os.IsNotExist(err) {
				glog.V(5).Infof("Unable to write to sysfs file")
				glog.Warning("File does not exist: ", fsFile)
				failedMessage = append(failedMessage, FailedMessage{ID: device, ErrorMsg: "Sysfs file does not exist for mclk OverDrive"})
				continue
			}
			f, err := os.OpenFile(fsFile, os.O_WRONLY, 0644)
			if err != nil {
				glog.V(5).Infof("Unable to open sysfs file %v", fsFile)
				glog.Warning("IO or OS error")
				failedMessage = append(failedMessage, FailedMessage{ID: device, ErrorMsg: "Unable to open sysfs file for mclk OverDrive"})
				continue
			}
			defer f.Close()
			_, err = f.WriteString(fmt.Sprintf("%v", intValue))
			if err != nil {
				glog.V(5).Infof("Unable to write to sysfs file %v", fsFile)
				glog.Warning("IO or OS error")
				failedMessage = append(failedMessage, FailedMessage{ID: device, ErrorMsg: "Unable to write to sysfs file for mclk OverDrive"})
				continue
			}
			glog.V(5).Infof("device%v Successfully set %s OverDrive to %d%%", device, clktype, intValue)
		} else if clktype == "sclk" {
			err := rsmiDevOverdriveLevelSet(device, intValue)
			if err == nil {
				glog.V(5).Infof("device:%v Successfully set %s OverDrive to %d%%", device, clktype, intValue)
			} else {
				glog.Errorf("device:%v Unable to set %s OverDrive to %d%%", device, clktype, intValue)
				failedMessage = append(failedMessage, FailedMessage{ID: device, ErrorMsg: err.Error()})
			}
		} else {
			glog.Errorf("device:%v Unable to set OverDrive", device)
			glog.Errorf("Unsupported clock type %v", clktype)
			failedMessage = append(failedMessage, FailedMessage{ID: device, ErrorMsg: "Unsupported clock type"})
		}
	}
	return
}

// SetPerfDeterminism 设置时钟频率级别以启用性能确定性
// @Summary 设置时钟频率级别以启用性能确定性
// @Description 根据设备ID列表和给定的时钟频率值，设置设备的性能确定性模式
// @Tags Device
// @Accept  json
// @Produce  json
// @Param dvIdList body []int true "设备 ID 列表"
// @Param clkvalue query string true "时钟频率值"
// @Success 200 {array} FailedMessage
// @Failure 400 {object} FailedMessage
// @Router /SetPerfDeterminism [post]
func SetPerfDeterminism(dvIdList []int, clkvalue string) (failedMessage []FailedMessage, err error) {
	// 验证 clkvalue 是否为有效的整数
	intValue, err := strconv.ParseInt(clkvalue, 10, 64)
	if err != nil {
		glog.Errorf("Unable to set Performance Determinism")
		glog.Errorf("clkvalue:%v is not an integer", clkvalue)
		return failedMessage, fmt.Errorf("clkvalue:%v is not an integer", clkvalue)
	}

	errorMap := make(map[int][]string)
	// 遍历每个设备并设置性能确定性模式
	for _, device := range dvIdList {
		err := rsmiPerfDeterminismModeSet(device, intValue)
		if err != nil {
			errorMap[device] = append(errorMap[device], err.Error())
			glog.Errorf("Unable to set performance determinism and clock frequency to %v for device %v", clkvalue, device)
		}
	}
	for id, msg := range errorMap {
		failedMessage = append(failedMessage, FailedMessage{ID: id, ErrorMsg: strings.Join(msg, "; ")})
	}
	return
}

// SetFanSpeed 设置风扇转速 [0-255]
// @Summary 设置风扇转速
// @Description 根据设备ID列表和给定的风扇速度，设置设备的风扇速度
// @Tags Device
// @Accept  json
// @Produce  json
// @Param dvIdList body []int true "设备 ID 列表"
// @Param fan query string true "风扇速度值或百分比（如 50%）"
// @Success 200 {string} string "成功信息"
// @Failure 400 {string} string "失败信息"
// @Router /SetFanSpeed [post]
func SetFanSpeed(dvIdList []int, fan string) {
	for _, device := range dvIdList {
		var fanLevel int64
		var err error
		lastChar := fan[len(fan)-1:]
		if lastChar == "%" {
			percentValue, err := strconv.Atoi(fan[:len(fan)-1])
			if err != nil {
				glog.Errorf("Invalid fan speed percentage: %s", fan)
				continue
			}
			fanLevel = int64(percentValue * 255 / 100)
		} else {
			fanLevel, err = strconv.ParseInt(fan, 10, 64)
			if err != nil {
				glog.Errorf("Invalid fan speed value: %s", fan)
				continue
			}
		}
		glog.V(5).Infof("Setting fan speed fanLevel value to %v", fanLevel)
		err = rsmiDevFanSpeedSet(device, 0, fanLevel)
		if err != nil {
			log.Printf("Failed to set fan speed for device %d", device)
		}
	}
}

// DevFanRpms 获取设备的风扇速度
// @Summary 获取设备的风扇速度
// @Description 获取指定设备的风扇速度（RPM）
// @Tags Device
// @Accept  json
// @Produce  json
// @Param dvInd path int true "设备索引"
// @Success 200 {integer} int64 "风扇速度 (RPM)"
// @Failure 400 {string} string "失败信息"
// @Router /DevFanRpms/{dvInd} [get]
func DevFanRpms(dvInd int) (speed int64, err error) {
	return rsmiDevFanRpmsGet(dvInd, 0)
}

// SetPerformanceLevel 设置设备性能等级
// @Summary 设置设备性能等级
// @Description 根据设备ID列表和给定的性能等级，设置设备的性能等级
// @Tags Device
// @Accept  json
// @Produce  json
// @Param dvIdList body []int true "设备 ID 列表"
// @Param level query string true "性能等级 (auto, low, high, normal)"
// @Success 200 {array} FailedMessage
// @Failure 400 {object} FailedMessage
// @Router /SetPerformanceLevel [post]
func SetPerformanceLevel(dvIdList []int, level string) (failedMessages []FailedMessage) {
	for _, device := range dvIdList {
		devPerfLevel, valid := validLevels[level]
		if !valid {
			glog.Errorf("device :%v Unable to set Performance Level, Invalid Performance level: %v", device, level)
		} else {
			err := rsmiDevPerfLevelSet(device, devPerfLevel)
			if err != nil {
				glog.Errorf("device:%v Failed to set performance level to %v", device, level)
				failedMessages = append(failedMessages, FailedMessage{
					ID:       device,
					ErrorMsg: fmt.Sprintf("Failed to set performance level to %v", level),
				})
			}
		}
	}
	return
}

// SetProfile 设置功率配置
// @Summary 设置功率配置
// @Description 根据设备ID列表和给定的功率配置文件，设置设备的功率配置
// @Tags Power
// @Accept  json
// @Produce  json
// @Param dvIdList body []int true "设备 ID 列表"
// @Param profile query string true "功率配置文件名称"
// @Success 200 {array} FailedMessage "设置成功的消息列表"
// @Failure 400 {object} FailedMessage "失败的消息列表"
// @Router /SetProfile [post]
func SetProfile(dvIdList []int, profile string) (failedMessages []FailedMessage) {

	for _, device := range dvIdList {
		// 获取先前的配置文件
		status, err := rsmiDevPowerProfilePresetsGet(device, 0)
		glog.V(5).Infof("status.Current: %v, int:%v", status.Current, int(status.Current))
		if err == nil {
			previousProfile := profileString(int(status.Current))

			// 确定期望的配置文件
			glog.V(5).Infof("previousProfile value: %v", previousProfile)
			glog.V(5).Infof("desiredProfile value: %v", profile)
			glog.V(5).Infof("previousProfile and desiredProfile:%v", profile == previousProfile)
			if profile == "UNKNOWN" {
				glog.Errorf("device:%v Unable to set profile to: %v (UNKNOWN profile)", device, profile)
				failedMessages = append(failedMessages, FailedMessage{ID: device, ErrorMsg: fmt.Sprintf("Unable to set profile to: %s (UNKNOWN profile)", profile)})
				continue
			}

			// 设置配置文件
			if previousProfile == profile {
				glog.V(5).Infof("device:%v Profile was already set to%v", device, previousProfile)
			} else {
				err := rsmiDevPowerProfileSet(device, 0, profileEnum(profile))
				if err == nil {
					// 获取当前配置文件
					profileStatus, err := rsmiDevPowerProfilePresetsGet(device, 0)
					if err == nil {
						currentProfile := profileString(int(profileStatus.Current))
						if currentProfile == profile {
							glog.V(5).Infof("device:%v Successfully set profile to:%v", device, profile)
						} else {
							glog.Errorf("device:%v Failed to set profile to: %v", device, profile)
							failedMessages = append(failedMessages, FailedMessage{ID: device, ErrorMsg: fmt.Sprintf("Failed to set profile to: %s", profile)})
						}
					}
				} else {
					glog.Errorf("device:%v Failed to set profile to: %v", device, err.Error())
					failedMessages = append(failedMessages, FailedMessage{ID: device, ErrorMsg: fmt.Sprintf("Failed to set profile to: %s", profile)})
				}
			}
		}
	}

	return
}

// DevPowerProfileSet 设置设备功率配置文件
// @Summary 设置设备功率配置文件
// @Description 设置指定设备的功率配置文件
// @Tags Power
// @Accept  json
// @Produce  json
// @Param dvInd path int true "设备索引"
// @Param reserved query int true "保留参数，通常为0"
// @Param profile query int true "功率配置文件的枚举值"
// @Success 200 {string} string "成功信息"
// @Failure 400 {string} string "失败信息"
// @Router /DevPowerProfileSet [post]
func DevPowerProfileSet(dvInd int, reserved int, profile PowerProfilePresetMasks) (err error) {
	return rsmiDevPowerProfileSet(dvInd, reserved, profile)
}

func DevPowerProfilePresetsGet(dvInd, sensorInd int) (powerProfileStatus PowerProfileStatus, err error) {
	return rsmiDevPowerProfilePresetsGet(dvInd, sensorInd)
}

// GetBus 获取设备总线信息
// @Summary 获取设备总线信息
// @Description 获取指定设备的总线信息
// @Tags Device
// @Accept  json
// @Produce  json
// @Param dvInd path int true "设备索引"
// @Success 200 {string} string "设备总线ID"
// @Failure 400 {string} string "失败信息"
// @Router /GetBus/{device} [get]
func GetBus(dvInd int) (picId string, err error) {

	bdfid, err := rsmiDevPciIdGet(dvInd)
	if err != nil {
		return picId, err
	}
	domain := (bdfid >> 32) & 0xffffffff
	bus := (bdfid >> 8) & 0xff
	dev := (bdfid >> 3) & 0x1f
	function := bdfid & 0x7
	picId = fmt.Sprintf("%04x:%02x:%02x.%x", domain, bus, dev, function)
	return
}

// ShowAllConciseHw 显示设备硬件信息
// @Summary 显示设备硬件信息
// @Description 显示指定设备列表的简要硬件信息
// @Tags Hardware
// @Accept  json
// @Produce  json
// @Param dvIdList body []int true "设备 ID 列表"
// @Success 200 {string} string "成功信息"
// @Failure 400 {string} string "失败信息"
// @Router /ShowAllConciseHw [post]
func ShowAllConciseHw(dvIdList []int) {
	header := []string{"HCU", "DID", "GFX RAS", "SDMA RAS", "UMC RAS", "VBIOS", "BUS"}
	headWidths := make([]int, len(header))
	for i, head := range header {
		headWidths[i] = len(head) + 2
	}

	values := make(map[string][]string)
	for _, device := range dvIdList {
		gpuid, _ := rsmiDevIdGet(device)
		gfxRas, _ := EccStatus(device, RSMIGpuBlockGFX)
		sdmaRas, _ := EccStatus(device, RSMIGpuBlockSDMA)
		umcRas, _ := EccStatus(device, RSMIGpuBlockUMC)
		vbios, _ := VbiosVersion(device)
		bus, _ := GetBus(device)
		values[fmt.Sprintf("card%d", device)] = []string{
			fmt.Sprintf("GPU%d", device), strconv.Itoa(gpuid), gfxRas, sdmaRas, umcRas, vbios, bus,
		}
	}

	valWidths := make(map[int][]int)
	for _, device := range dvIdList {
		valWidths[device] = make([]int, len(values[fmt.Sprintf("card%d", device)]))
		for i, val := range values[fmt.Sprintf("card%d", device)] {
			valWidths[device][i] = len(val) + 2
		}
	}
	maxWidths := headWidths
	for _, device := range dvIdList {
		for col := range valWidths[device] {
			if valWidths[device][col] > maxWidths[col] {
				maxWidths[col] = valWidths[device][col]
			}
		}
	}

	for i, head := range header {
		fmt.Printf("%-*s", maxWidths[i], head)
	}
	fmt.Println()

	for _, device := range dvIdList {
		for i, val := range values[fmt.Sprintf("card%d", device)] {
			fmt.Printf("%-*s", maxWidths[i], val)
		}
		fmt.Println()
	}

}

// ShowClocks 显示时钟信息
// @Summary 显示时钟信息
// @Description 显示指定设备的时钟信息
// @Tags Clock
// @Accept  json
// @Produce  json
// @Param dvIdList body []int true "设备 ID 列表"
// @Success 200 {string} string "成功信息"
// @Failure 400 {string} string "失败信息"
// @Router /ShowClocks [post]
func ShowClocks(dvIdList []int) {

	for _, device := range dvIdList {
		for clkType, clkID := range rsmiClkNamesDict {

			freq, err := rsmiDevGpuClkFreqGet(device, clkID)
			if err == nil {
				glog.V(5).Infof("device:%v Supported %v frequencies on GPU%v", device, clkType, device)
				for x := 0; x < int(freq.NumSupported); x++ {
					fr := fmt.Sprintf("%dMhz", freq.Frequency[x]/1000000)
					if uint32(x) == freq.Current {
						glog.V(5).Infof("Device %d: %d %s *", device, x, fr)
					} else {
						glog.V(5).Infof("Device %d: %d %s", device, x, fr)
					}
				}
			} else {
				glog.Errorf("device:%v clkType:%v frequency is unsupported", device, clkType)

			}
		}
		bw, err := rsmiDevPciBandwidthGet(device)
		if err == nil {
			glog.V(5).Infof("Supported PCIe frequencies on GPU%d", device)
			for x := 0; x < int(bw.TransferRate.NumSupported); x++ {
				fr := fmt.Sprintf("%.1fGT/s x%d", float64(bw.TransferRate.Frequency[x])/1000000000, bw.Lanes[x])
				if uint32(x) == bw.TransferRate.Current {
					glog.V(5).Infof("Device %d: %d %s *", device, x, fr)
				} else {
					glog.V(5).Infof("Device %d: %d %s", device, x, fr)
				}
			}
		}
	}
}

func GetClocksByType(dvInd int, clkType RSMIClkType) (frequencyList []uint64, currentFrequency uint32, err error) {
	frequencyStruct, err := rsmiDevGpuClkFreqGet(dvInd, clkType)
	if err != nil {
		return nil, 0, err
	}
	frequencyList = make([]uint64, frequencyStruct.NumSupported)
	for i := 0; i < int(frequencyStruct.NumSupported); i++ {
		fr := frequencyStruct.Frequency[i] / 1000000
		frequencyList[i] = fr
	}
	currentFrequency = frequencyStruct.Current
	return
}

// ShowCurrentFans 展示风扇转速和风扇级别
// @Summary 展示风扇转速和风扇级别
// @Description 显示指定设备的当前风扇转速和风扇级别
// @Tags Fan
// @Accept  json
// @Produce  json
// @Param dvIdList body []int true "设备 ID 列表"
// @Param printJSON query bool true "是否以 JSON 格式打印输出"
// @Success 200 {string} string "成功信息"
// @Failure 400 {string} string "失败信息"
// @Router /ShowCurrentFans [post]
func ShowCurrentFans(dvIdList []int, printJSON bool) {
	glog.V(5).Infof("--------- Current Fan Metric ---------")
	var sensorInd uint32 = 0

	for _, device := range dvIdList {
		fanLevel, fanSpeed, err := FanSpeedInfo(device)
		if err != nil {
			glog.Errorf("Unable to detect fan speed for GPU %v: %v", device, err)
			continue
		}

		fanSpeed = float64(int64(fanSpeed + 0.5)) // 四舍五入

		if fanLevel == 0 || fanSpeed == 0 {
			glog.V(5).Infof("Device %v: Unable to detect fan speed", device)
			glog.V(5).Infof("Current fan speed is: %v", fanSpeed)
			glog.V(5).Infof("Current fan level is: %v", fanLevel)
			glog.V(5).Infof("GPU might be cooled with a non-PWM fan")
			continue
		}
		if printJSON {
			glog.V(5).Infof("Device %v: Fan speed (level): %v", device, fanLevel)
			glog.V(5).Infof("Device %v: Fan speed (%%): %.0f", device, fanSpeed)
		} else {
			glog.V(5).Infof("Device %v: Fan Level: %d (%.0f%%)", device, fanLevel, fanSpeed)
		}

		rpmSpeed, err := rsmiDevFanRpmsGet(device, int(sensorInd))
		if err == nil {
			glog.V(5).Infof("Device %v: Fan RPM: %v", device, rpmSpeed)
		} else {
			glog.Errorf("Device %v: Error getting fan RPM: %v", device, err)
		}
	}
	glog.V(5).Infof("--------------------------------------")
}

// ShowCurrentTemps 显示所有设备的所有可用温度传感器的温度
// @Summary 显示设备温度传感器数据
// @Tags Temperature
// @Param dvIdList query []int true "设备 ID 列表"
// @Success 200 {object} TemperatureInfo "温度信息列表"
// @Failure 400 {object} error "错误信息"
// @Router /ShowCurrentTemps [get]
func ShowCurrentTemps(dvIdList []int) (temperatureInfos []TemperatureInfo, err error) {
	glog.V(5).Infof("--------- Temperature ---------")
	for _, device := range dvIdList {
		sensorTemps := make(map[string]float64)
		for _, sensor := range tempTypeList {
			temp, err := Temperature(device)
			if err != nil {
				glog.Errorf("Error getting temperature for device %d sensor %d: %v", device, sensor.Type, err)
			} else {
				glog.V(5).Infof("Device %d Temperature (Sensor %s): %.2f°C", device, sensor.Name, temp)
				sensorTemps[sensor.Name] = temp
			}
			deviceTempInfo := TemperatureInfo{
				DeviceID:    device,
				SensorTemps: sensorTemps,
			}
			temperatureInfos = append(temperatureInfos, deviceTempInfo)
		}
	}
	glog.V(5).Infof("--------------------------------")
	glog.V(5).Infof("temperatureInfos:%v", (temperatureInfos))
	return
}

// ShowFwInfo 显示给定设备列表中指定固件类型的固件版本信息
// @Summary 显示设备固件版本信息
// @Tags Firmware
// @Param dvIdList query []int true "设备 ID 列表"
// @Param fwType query []string true "固件类型列表"
// @Success 200 {object} []FirmwareInfo "固件版本信息列表"
// @Failure 400 {object} error "错误信息"
// @Router /ShowFwInfo [get]
func ShowFwInfo(dvIdList []int, fwType []string) (fwInfos []FirmwareInfo, err error) {
	var firmwareBlocks []string
	if len(fwType) == 0 || contains(fwType, "all") {
		firmwareBlocks = fwBlockNames
	} else {
		for _, name := range fwType {
			if contains(fwBlockNames, strings.ToUpper(name)) {
				firmwareBlocks = append(firmwareBlocks, strings.ToUpper(name))
			}
		}
	}
	for _, device := range dvIdList {
		fwVerMap := make(map[string]string)
		for _, fwName := range firmwareBlocks {
			fwNameUpper := strings.ToUpper(fwName)
			fwVersion, err := rsmiDevFirmwareVersionGet(device, RSMIFwBlock(indexOf(fwBlockNames, fwNameUpper)))
			if err != nil {
				glog.Errorf("Error getting firmware version for device %v firmware block %v: %v", device, fwNameUpper, err)
				continue
			}

			var formattedFwVersion string
			if fwNameUpper == "VCN" || fwNameUpper == "VCE" || fwNameUpper == "UVD" || fwNameUpper == "SOS" {
				formattedFwVersion = fmt.Sprintf("0x%s", strings.ToUpper(fmt.Sprintf("%08x", fwVersion)))
			} else if fwNameUpper == "TA XGMI" || fwNameUpper == "TA RAS" || fwNameUpper == "SMC" {
				formattedFwVersion = fmt.Sprintf("%02d.%02d.%02d.%02d",
					(fwVersion>>24)&0xFF, (fwVersion>>16)&0xFF, (fwVersion>>8)&0xFF, fwVersion&0xFF)
			} else if fwNameUpper == "ME" || fwNameUpper == "MC" || fwNameUpper == "CE" {
				formattedFwVersion = fmt.Sprintf("\t\t%d", fwVersion)
			} else {
				formattedFwVersion = fmt.Sprintf("\t%d", fwVersion)
			}

			fwVerMap[fwNameUpper] = formattedFwVersion
			glog.V(5).Infof("Device %v %v firmware version: %v", device, fwNameUpper, formattedFwVersion)
		}
		fwInfos = append(fwInfos, FirmwareInfo{
			DeviceID:    device,
			FirmwareVer: fwVerMap,
		})
	}
	glog.V(5).Infof("fwInfos:%v", (fwInfos))
	return
}

// PidList 获取进程列表
// @Summary 获取计算进程列表
// @Tags Process
// @Success 200 {array} string "进程 ID 列表"
// @Failure 400 {object} error "错误信息"
// @Router /PidList [get]
func PidList() (pidList []string, err error) {
	processInfo, numItems, err := rsmiComputeProcessInfoGet()
	if err != nil {
		return nil, err
	}
	if numItems == 0 {
		return
	}
	for i := 0; i < numItems; i++ {
		pidList = append(pidList, fmt.Sprintf("%d", processInfo[i].ProcessID))
	}
	glog.V(5).Infof("pidList:%v", pidList)
	return
}

// GetCoarseGrainUtil 获取设备的粗粒度利用率
// @Summary 获取设备粗粒度利用率
// @Tags Utilization
// @Param dvInd query int true "设备 ID"
// @Param typeName query string false "利用率计数器类型名称"
// @Success 200 {array} UtilizationCounter "利用率计数器列表"
// @Failure 400 {object} error "错误信息"
// @Router /GetCoarseGrainUtil [get]
func GetCoarseGrainUtil(dvInd int, typeName string) (utilizationCounters []UtilizationCounter, err error) {
	var length int

	if typeName != "" {
		// 获取特定类型的利用率计数器
		var i UtilizationCounterType
		var found bool
		for index, name := range utilizationCounterName {
			if name == typeName {
				i = UtilizationCounterType(index)
				found = true
				break
			}
		}
		if !found {
			glog.V(5).Infof("No such coarse grain counter type: %v", typeName)
			return nil, fmt.Errorf("no such coarse grain counter type")
		}
		length = 1
		utilizationCounters = make([]UtilizationCounter, length)
		utilizationCounters[0].Type = i
	} else {
		// 获取所有类型的利用率计数器
		length = int(RSMI_UTILIZATION_COUNTER_LAST) + 1
		utilizationCounters = make([]UtilizationCounter, length)
		for i := 0; i < length; i++ {
			utilizationCounters[i].Type = UtilizationCounterType(i)
		}
	}
	_, err = rsmiUtilizationCountGet(dvInd, utilizationCounters, length)
	if err != nil {
		return nil, err
	}
	return
}

// ShowHCUUse HCU使用率
// @Summary 显示设备的 HCU 使用率
// @Tags HCU
// @Param dvIdList query []int true "设备 ID 列表"
// @Success 200 {object} []DeviceUseInfo "设备使用信息列表"
// @Failure 400 {object} error "错误信息"
// @Router /ShowHCUUse [get]
func ShowHCUUse(dvIdList []int) (deviceUseInfos []DeviceUseInfo, err error) {
	fmt.Printf(" time GPU is busy\n ")

	for _, device := range dvIdList {
		deviceUseInfo := DeviceUseInfo{
			DeviceID:    device,
			Utilization: make(map[string]uint64),
		}

		// 获取 GPU 使用百分比
		percent, err := HCUUse(device)
		if err != nil {
			fmt.Printf("Device %d: GPU use Unsupported\n", device)
			deviceUseInfo.GPUUsePercent = -1

		} else {
			fmt.Printf("Device %d: GPU use (%%) %d\n", device, percent)
			deviceUseInfo.GPUUsePercent = percent
		}

		// 获取粗粒度利用率
		typeName := "GFX Activity"
		utilCounters, err := GetCoarseGrainUtil(device, typeName)
		if err != nil {
			fmt.Printf("Device %d: Error getting coarse grain utilization: %v\n", device, err)
		} else {
			for _, counter := range utilCounters {
				fmt.Printf("Device %d: %s %d\n", device, utilizationCounterName[counter.Type], counter.Value)
				if int(counter.Type) < len(utilizationCounterName) {
					deviceUseInfo.Utilization[utilizationCounterName[counter.Type]] = counter.Value
				}
			}
		}
		deviceUseInfos = append(deviceUseInfos, deviceUseInfo)
	}
	glog.V(5).Infof("deviceUseInfos: %v", (deviceUseInfos))
	return
}

// ShowEnergy 展示设备消耗的能量
// @Summary 展示设备的能量消耗
// @Description 获取并展示指定设备的能量消耗情况。
// @Tags 设备
// @Param dvIdList query []int true "设备ID列表"
// @Success 200 {string} string "成功返回设备的能量消耗信息"
// @Failure 400 {string} string "请求参数错误"
// @Router /showEnergy [get]
func ShowEnergy(dvIdList []int) {
	for _, device := range dvIdList {
		power, counterResolution, _, err := rsmiDevEnergyCountGet(device)
		if err != nil {
			glog.Errorf("Error getting energy count for device %d: %v\n", device, err)
			continue
		}
		fmt.Printf("Device %d Energy counter: %d\n", device, power)
		fmt.Printf("Device %d Accumulated Energy (uJ): %.2f\n", device, float64(power)*float64(counterResolution))
	}
}

// ShowMemInfo 展示设备的内存信息
// @Summary 展示设备内存信息
// @Description 获取并展示指定设备的内存使用情况，包括不同类型的内存。
// @Tags 设备
// @Param dvIdList query []int true "设备ID列表"
// @Param memTypes query []string true "内存类型列表，如 'all' 或指定类型"
// @Success 200 {string} string "成功返回设备的内存信息"
// @Failure 400 {string} string "请求参数错误"
// @Router /showMemInfo [get]
func ShowMemInfo(dvIdList []int, memTypes []string) {
	var returnTypes []string

	if len(memTypes) == 1 && memTypes[0] == "all" {
		returnTypes = memoryTypeL
	} else {
		for _, memType := range memTypes {
			if contains(memoryTypeL, memType) {
				returnTypes = append(returnTypes, memType)
			} else {
				log.Printf("Invalid memory type: %s", memType)
				return
			}
		}
	}

	fmt.Println(" Memory Usage (Bytes) ")
	for _, device := range dvIdList {
		for _, mem := range returnTypes {
			memInfoUsed, memInfoTotal, err := MemInfo(device, mem)
			if err != nil {
				log.Printf("Error getting memory info for device %d: %v", device, err)
				continue
			}
			fmt.Println("device ", device, fmt.Sprintf("%s Total Memory (B)", mem), memInfoTotal)
			fmt.Println("device ", device, fmt.Sprintf("%s Total Used Memory (B)", mem), memInfoUsed)
		}
	}
	fmt.Println("End of Memory Usage")
}

// ShowMemUse 展示设备的内存使用情况
// @Summary 展示设备内存使用情况
// @Description 获取并展示指定设备的当前内存使用百分比和其他相关的利用率数据。
// @Tags 设备
// @Param dvIdList query []int true "设备ID列表"
// @Success 200 {string} string "成功返回设备的内存使用信息"
// @Failure 400 {string} string "请求参数错误"
// @Router /showMemUse [get]
func ShowMemUse(dvIdList []int) {
	fmt.Println("Current Memory Use")
	for _, device := range dvIdList {
		busyPercent, err := rsmiDevMemoryBusyPercentGet(device)
		if err == nil {
			fmt.Println("device: ", device, "GPU memory use (%)", busyPercent)
		}
		typeName := "Memory Activity"
		utilCounters, err := GetCoarseGrainUtil(device, typeName)
		if err == nil {
			for _, utCounter := range utilCounters {
				fmt.Println("device: ", device, utilizationCounterName[utCounter.Type], utCounter.Value)
			}
		} else {
			glog.Errorf("Device %d: Failed to get coarse grain util counters: %v", device, err)
		}
	}
}

// ShowMemVendor 展示设备供应商信息
// @Summary 展示设备的内存供应商信息
// @Description 获取并展示指定设备的内存供应商信息。
// @Tags 设备
// @Param dvIdList query []int true "设备ID列表"
// @Success 200 {object} []DeviceMemVendorInfo "成功返回设备的内存供应商信息"
// @Failure 400 {string} string "请求参数错误"
// @Router /showMemVendor [get]
func ShowMemVendor(dvIdList []int) (deviceMemVendorInfos []DeviceMemVendorInfo, err error) {
	for _, device := range dvIdList {
		vendor, err := rsmiDevVramVendorGet(device)
		if err == nil {
			glog.V(5).Infof("device:%v  GPU memory vendor:%v", device, vendor)
			deviceMemVendorInfos = append(deviceMemVendorInfos, DeviceMemVendorInfo{DeviceID: device, Vendor: vendor})
		} else {
			glog.Warning("GPU memory vendor missing or not supported")
			deviceMemVendorInfos = append(deviceMemVendorInfos, DeviceMemVendorInfo{DeviceID: device, Vendor: ""})
		}
	}
	glog.V(5).Infof("GPU memory vendor: %v", (deviceMemVendorInfos))
	return
}

// ShowPcieBw 展示设备的PCIe带宽使用情况
// @Summary 展示设备的PCIe带宽使用情况
// @Description 获取并展示指定设备的PCIe带宽使用情况，包括发送和接收的带宽。
// @Tags 设备
// @Param dvIdList query []int true "设备ID列表"
// @Success 200 {object} []PcieBandwidthInfo "成功返回设备的PCIe带宽使用信息"
// @Failure 400 {string} string "请求参数错误"
// @Router /showPcieBw [get]
func ShowPcieBw(dvIdList []int) (pcieBandwidthInfos []PcieBandwidthInfo, err error) {
	for _, device := range dvIdList {
		sent, received, maxPktSz, err := rsmiDevPciThroughputGet(device)
		if err == nil {
			sentMB, receivedMB, bw := pcieThroughputMB(sent, received, maxPktSz)
			bwstr := fmt.Sprintf("%.3f", bw)
			glog.V(5).Infof("device:%v Estimated maximum PCIe bandwidth over the last second (MB/s):%v", device, bwstr)
			pcieBandwidthInfos = append(pcieBandwidthInfos, PcieBandwidthInfo{DvInd: device, Sent: sentMB, Received: receivedMB, Bw: bw})
		} else {
			glog.Warning("GPU PCIe bandwidth usage not supported")
			pcieBandwidthInfos = append(pcieBandwidthInfos, PcieBandwidthInfo{DvInd: device, Sent: 0, Received: 0, Bw: 0})
		}
	}
	glog.V(5).Infof("pcieBandwidthInfos:%v", (pcieBandwidthInfos))
	return
}

func PcieBw(dvInd int) (pcieBandwidthInfo PcieBandwidthInfo, err error) {
	sent, received, maxPktSz, err := rsmiDevPciThroughputGet(dvInd)
	if err == nil {
		sentMB, receivedMB, bw := pcieThroughputMB(sent, received, maxPktSz)
		pcieBandwidthInfo = PcieBandwidthInfo{
			DvInd:    dvInd,
			Sent:     sentMB,
			Received: receivedMB,
			Bw:       bw,
		}
		glog.V(5).Infof("device:%v Estimated maximum PCIe bandwidth over the last second (MB/s):%v", pcieBandwidthInfo, bw)
	}
	glog.V(5).Infof("pcieBandwidthInfo:%v", (pcieBandwidthInfo))
	return
}

func PcieBwTest(dvInd int) (pcieBandwidthInfo PcieBandwidthInfo, err error) {
	sent, received, maxPktSz, err := rsmiDevPciThroughputGet(dvInd)
	if err == nil {
		sentMB, receivedMB, bw := pcieThroughputMB(sent, received, maxPktSz)
		pcieBandwidthInfo = PcieBandwidthInfo{
			DvInd:    dvInd,
			Sent:     sentMB,
			Received: receivedMB,
			Bw:       bw,
		}
		glog.V(5).Infof("device:%v Estimated maximum PCIe bandwidth over the last second (MB/s):%v", pcieBandwidthInfo, bw)
	}
	glog.V(5).Infof("pcieBandwidthInfo:%v", (pcieBandwidthInfo))
	return
}

// ShowPcieReplayCount 展示设备的PCIe重放计数
// @Summary 展示设备的PCIe重放计数
// @Description 获取并展示指定设备的PCIe重放计数。
// @Tags 设备
// @Param dvIdList query []int true "设备ID列表"
// @Success 200 {object} []PcieReplayCountInfo "设备的PCIe重放计数信息"
// @Failure 400 {string} string "请求参数错误"
// @Router /showPcieReplayCount [get]
func ShowPcieReplayCount(dvIdList []int) (pcieReplayCountInfos []PcieReplayCountInfo, err error) {
	for _, device := range dvIdList {
		count, err := rsmiDevPciReplayCounterGet(device)
		if err == nil {
			glog.V(5).Infof("device:%v PCIe Replay Count:%v", device, count)
			pcieReplayCountInfos = append(pcieReplayCountInfos, PcieReplayCountInfo{DeviceID: device, Count: count})
		} else {
			glog.Warning("GPU PCIe replay count not supported")
			pcieReplayCountInfos = append(pcieReplayCountInfos, PcieReplayCountInfo{DeviceID: device, Count: 0})
		}
	}
	glog.V(5).Infof("pcieReplayCountInfos:%v", (pcieReplayCountInfos))
	return
}

// 获取指定进程的进程信息
func ProcessInfo(pid int) (proc ProcessInfos, err error) {
	return rsmiComputeProcessInfoByPidGet(pid)
}

// 获取指定设备上的进程信息
func ProcessInfoByDevice(pid int, dvInd int) (proc ProcessInfos, err error) {
	return rsmiProcessInfoByDevice(pid, dvInd)
}

func ProcessHCU(pid int) (dvIndices []int, err error) {
	return rsmiComputeProcessGpusGet(pid)
}

// @Tags System
// ProcessHCUInfo 进程列表信息
// @Summary 进程列表信息
// @Description 获取并进程信息和使用的HCU设备信息。
// @Success 200 {string} string "成功返回进程信息"
// @Failure 400 {string} string "请求错误"
// @Router /showPids [get]
func ProcessHCUInfo() ([]Process, error) {
	var (
		processes []Process // 最终返回的进程信息
		errs      []error   // 错误信息列表
		pidList   []int     //进程id
	)

	// 调用获取进程信息的函数
	processInfo, numItems, err := rsmiComputeProcessInfoGet()
	if err != nil {
		return nil, err
	}
	// 如果没有获取到进程信息
	if numItems == 0 {
		return nil, nil
	}
	// 遍历进程信息，提取 ProcessID 并转换为 int
	for i := 0; i < numItems; i++ {
		pidList = append(pidList, int(processInfo[i].ProcessID))
	}
	for _, pid := range pidList {
		// 获取进程基本信息
		procInfo, err := ProcessInfo(pid)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to get process info for PID %d: %w", pid, err))
			continue
		}
		// 获取 GPU 索引信息
		dvIndices, err := ProcessHCU(pid)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to get HCU info for PID %d: %w", pid, err))
			dvIndices = nil // 确保错误时 GPU 索引为空
		}
		processName := ProcessName(pid)
		// 构造 Process 并添加到结果列表
		processes = append(processes, Process{
			ProcessID:    procInfo.ProcessID,
			Pasid:        procInfo.Pasid,
			VramUsage:    procInfo.VramUsage,
			SdmaUsage:    procInfo.SdmaUsage,
			CuOccupancy:  procInfo.CuOccupancy,
			ProcessName:  processName,
			MinorNumbers: dvIndices,
		})
	}
	glog.V(5).Infof("GPU process:%v", (processes))
	// 如果有错误，返回错误信息
	if len(errs) > 0 {
		return processes, fmt.Errorf("encountered errors during processing: %v", errs)
	}
	return processes, nil
}

// ShowPids 展示进程信息
// @Summary 展示系统中正在运行的KFD进程信息
// @Description 获取并展示当前系统中运行的KFD进程的详细信息。
// @Tags 系统
// @Success 200 {string} string "成功返回进程信息"
// @Failure 400 {string} string "请求错误"
// @Router /showPids [get]
func ShowPids() error {
	pidList, err := PidList()
	if err != nil {
		fmt.Printf("Error getting PID list: %v\n", err)
		return err
	}

	title := "KFD PROCESSES"
	headers := []string{"PID", "PROCESS NAME", "HCU", "VRAM USED", "SDMA USED", "CU OCCUPANCY"}

	var rows [][]string
	if len(pidList) == 0 {
		rows = append(rows, []string{"None", "", "", "", "", ""})
	} else {
		for _, pidStr := range pidList {
			pid, _ := strconv.Atoi(pidStr)

			gpuNumber := "UNKNOWN"
			vramUsage := "UNKNOWN"
			sdmaUsage := "UNKNOWN"
			cuOccupancy := "UNKNOWN"

			if dvIndices, err := rsmiComputeProcessGpusGet(pid); err == nil {
				gpuNumber = fmt.Sprintf("%d", len(dvIndices))
			}
			if proc, err := rsmiComputeProcessInfoByPidGet(pid); err == nil {
				// VRAM 单位转换为 MB
				vramInMB := proc.VramUsage / (1024 * 1024)
				vramUsage = fmt.Sprintf("%dMB", vramInMB)
				sdmaUsage = fmt.Sprintf("%d", proc.SdmaUsage)
				cuOccupancy = fmt.Sprintf("%d", proc.CuOccupancy)
			}

			rows = append(rows, []string{
				pidStr,
				GetProcessName(pid),
				gpuNumber,
				vramUsage,
				sdmaUsage,
				cuOccupancy,
			})
		}
	}

	printAsciiTable(title, headers, rows)
	return nil
}

func GetProcessName(pid int) string {
	if pid < 1 {
		log.Println("PID must be greater than 0")
		return "UNKNOWN"
	}

	cmd := exec.Command("sh", "-c", fmt.Sprintf("ps -p %d -o comm=", pid))
	output, err := cmd.Output()
	if err != nil {
		log.Println("Error executing command:", err)
		return "UNKNOWN"
	}

	pName := strings.TrimSpace(string(output))
	if pName == "" {
		pName = "UNKNOWN"
	}

	return pName
}

// ShowPower 展示设备的平均功率
// @Summary 展示设备的平均功率消耗
// @Description 获取并展示指定设备的平均图形功率消耗。
// @Tags 设备
// @Param dvIdList query []int true "设备ID列表"
// @Success 200 {object} []DevicePowerInfo "设备的功率信息"
// @Failure 400 {string} string "请求参数错误"
// @Router /showPower [get]
func ShowPower(dvIdList []int) (devicePowerInfos []DevicePowerInfo, err error) {
	fmt.Println("========== Power Consumption ==========")

	for _, device := range dvIdList {
		power, err := Power(device)
		if err != nil {
			glog.Errorf("device:%v Unable to get Average Graphics Package Power Consumption", device)
			devicePowerInfos = append(devicePowerInfos, DevicePowerInfo{DeviceID: device, Power: -1})
			continue
		}
		if power != 0 {
			fmt.Println("device:", device, "Average Graphics Package Power (W)", fmt.Sprintf("%d", power))
			devicePowerInfos = append(devicePowerInfos, DevicePowerInfo{DeviceID: device, Power: power})
		} else {
			glog.Errorf("device:%v Unable to get Average Graphics Package Power Consumption", device)
			devicePowerInfos = append(devicePowerInfos, DevicePowerInfo{DeviceID: device, Power: -1})
		}
	}
	glog.V(5).Infof("devicePowerInfos:%v", (devicePowerInfos))
	return
}

// 获取设备电压/频率曲线信息(K100 AI不支持)
func DevOdVoltInfoGet(deInd int) (odv OdVoltFreqData, err error) {
	odv, err = rsmiDevOdVoltInfoGet(deInd)
	return
}

// ShowPowerPlayTable 展示设备的GPU内存时钟频率和电压
// @Summary 展示设备的GPU内存时钟频率和电压
// @Description 获取并展示指定设备的GPU内存时钟频率和电压表。
// @Tags 设备
// @Param dvIdList query []int true "设备ID列表"
// @Success 200 {object} []DevicePowerPlayInfo "设备的GPU时钟频率和电压信息"
// @Failure 400 {string} string "请求参数错误"
// @Router /showPowerPlayTable [get]
func ShowPowerPlayTable(dvIdList []int) (devicePowerPlayInfos []DevicePowerPlayInfo, err error) {
	fmt.Println("========== GPU Memory clock frequencies and voltages ==========")
	for _, device := range dvIdList {
		odv, err := rsmiDevOdVoltInfoGet(device)
		if err != nil {
			log.Printf("Error retrieving voltage info for device %d: %v\n", device, err)
			continue
		}

		od_sclk := []string{
			fmt.Sprintf("0: %dMhz", odv.CurrSclkRange.LowerBound/1000000),
			fmt.Sprintf("1: %dMhz", odv.CurrSclkRange.UpperBound/1000000),
		}

		od_mclk := fmt.Sprintf("1: %dMhz", odv.CurrMclkRange.UpperBound/1000000)

		od_vddc_curve := make([]string, 3)
		for position := 0; position < 3; position++ {
			od_vddc_curve[position] = fmt.Sprintf("%d: %dMhz %dmV", position,
				odv.Curve.VcPoints[position].Frequency/1000000,
				odv.Curve.VcPoints[position].Voltage)
		}

		od_range := []string{
			fmt.Sprintf("SCLK: %dMhz %dMhz", odv.SclkFreqLimits.LowerBound/1000000, odv.SclkFreqLimits.UpperBound/1000000),
			fmt.Sprintf("MCLK: %dMhz %dMhz", odv.MclkFreqLimits.LowerBound/1000000, odv.MclkFreqLimits.UpperBound/1000000),
		}

		for position := 0; position < 3; position++ {
			od_range = append(od_range, fmt.Sprintf("VDDC_CURVE_SCLK[%d]: %dMhz", position, odv.Curve.VcPoints[position].Frequency/1000000))
			od_range = append(od_range, fmt.Sprintf("VDDC_CURVE_VOLT[%d]: %dmV", position, odv.Curve.VcPoints[position].Voltage))
		}

		powerPlayInfo := DevicePowerPlayInfo{
			DeviceID:  device,
			SCLK:      od_sclk,
			MCLK:      od_mclk,
			DDC_CURVE: od_vddc_curve,
			RANGE:     od_range,
		}

		devicePowerPlayInfos = append(devicePowerPlayInfos, powerPlayInfo)
		glog.V(5).Infof("DevicePowerPlayInfo:%v", (devicePowerPlayInfos))
	}

	fmt.Println("===============================================================")
	return
}

// ShowProductName 显示设备列表中所请求的产品名称
// @Summary 显示设备的产品名称
// @Description 获取并显示指定设备的产品名称、供应商、系列、型号和SKU信息。
// @Tags 设备
// @Param dvIdList query []int true "设备ID列表"
// @Success 200 {array} DeviceproductInfo "设备的产品信息列表"
// @Failure 400 {string} string "请求参数错误"
// @Router /showProductName [get]
func ShowProductName(dvIdList []int) (deviceProductInfos []DeviceproductInfo, err error) {
	fmt.Println("========== Product Info ==========")
	for _, device := range dvIdList {
		deviceProductInfo := DeviceproductInfo{DeviceID: device}

		// Retrieve card vendor
		vendor, err := rsmiDevVendorNameGet(device)
		if err != nil {
			log.Printf("Incompatible device. GPU[%d]: Expected vendor name: Advanced Micro Devices, Inc. [AMD/ATI]\nGPU[%d]: Actual vendor name: %s\n", device, device, vendor)
			continue
		}
		deviceProductInfo.CardVendor = vendor

		// Retrieve the device series
		series, err := rsmiDevNameGet(device)
		if err == nil {
			deviceProductInfo.CardSeries = series
			fmt.Printf("GPU[%d] Card series: %s\n", device, series)
		}

		// Retrieve the device model
		model, err := rsmiDevSubsystemNameGet(device)
		if err == nil {
			deviceProductInfo.CardModel = model
			fmt.Printf("GPU[%d] Card model: %s\n", device, model)
		}

		fmt.Printf("GPU[%d] Card vendor: %s\n", device, vendor)

		// Retrieve the device SKU
		vbios, err := rsmiDevVbiosVersionGet(device, 256)
		if err == nil {
			deviceProductInfo.CardSKU = vbios
			fmt.Printf("GPU[%d] Card SKU: %s\n", device, deviceProductInfo.CardSKU)
		}

		deviceProductInfos = append(deviceProductInfos, deviceProductInfo)

	}

	fmt.Println("==================================")
	glog.V(5).Infof("deviceProductInfos:%v", (deviceProductInfos))
	return
}

// ShowProfile 可用电源配置文件
// @Summary 显示设备的电源配置文件
// @Description 获取并显示指定设备的电源配置文件，包括可用的电源配置选项。
// @Tags 设备
// @Param dvIdList query []int true "设备ID列表"
// @Success 200 {array} DeviceProfile "设备的电源配置文件信息列表"
// @Failure 400 {string} string "请求参数错误"
// @Router /showProfile [get]
func ShowProfile(dvIdList []int) (deviceProfiles []DeviceProfile, err error) {
	fmt.Println(" Show Power Profiles ")

	for _, device := range dvIdList {
		status, err := rsmiDevPowerProfilePresetsGet(device, 0)
		if err != nil {
			log.Printf("Error getting power profile presets: %v", err)
			continue
		}

		binaryMaskString := fmt.Sprintf("%07b", status.AvailableProfiles)
		bitMaskPosition := 0
		profileNumber := 0
		var profiles []string

		for bitMaskPosition < 7 {
			if binaryMaskString[6-bitMaskPosition] == '1' {
				profileNumber++
				var profileInfo string
				if 1<<bitMaskPosition == int(status.Current) {
					profileInfo = fmt.Sprintf("%d. Available power profile (#%d of 7): %s*", profileNumber, bitMaskPosition+1, profileString(1<<bitMaskPosition))
				} else {
					profileInfo = fmt.Sprintf("%d. Available power profile (#%d of 7): %s", profileNumber, bitMaskPosition+1, profileString(1<<bitMaskPosition))
				}
				profiles = append(profiles, profileInfo)
			}
			bitMaskPosition++
		}

		deviceProfiles = append(deviceProfiles, DeviceProfile{DeviceID: device, Profiles: profiles})
	}
	glog.V(5).Infof("deviceProfiles: %v", (deviceProfiles))
	return
}

// ShowRange 电流或电压范围
// @Summary 显示设备的电流或电压范围
// @Description 获取并显示指定设备的有效电流或电压范围。
// @Tags 设备
// @Param dvIdList query []int true "设备ID列表"
// @Param rangeType query string true "范围类型 (sclk, mclk, voltage)"
// @Success 200 {string} string "设备的电流或电压范围信息"
// @Failure 400 {string} string "请求参数错误"
// @Router /showRange [get]
func ShowRange(dvIdList []int, rangeType string) {
	if rangeType != "sclk" && rangeType != "mclk" && rangeType != "voltage" {
		fmt.Println(0, fmt.Sprintf("Invalid range identifier %s", rangeType))
		return
	}

	fmt.Println(fmt.Sprintf(" Show Valid %s Range ", rangeType))

	for _, device := range dvIdList {
		odvf, err := rsmiDevOdVoltInfoGet(device)
		if err != nil {
			log.Printf("Error getting OD volt info: %v", err)
			fmt.Println(device, fmt.Sprintf("Unable to display %s range", rangeType))
			continue
		}
		switch rangeType {
		case "sclk":
			fmt.Println(device, fmt.Sprintf("Valid sclk range: %dMhz - %dMhz",
				odvf.CurrSclkRange.LowerBound/1000000, odvf.CurrSclkRange.UpperBound/1000000))
		case "mclk":
			fmt.Println(device, fmt.Sprintf("Valid mclk range: %dMhz - %dMhz",
				odvf.CurrMclkRange.LowerBound/1000000, odvf.CurrMclkRange.UpperBound/1000000))
		case "voltage":
			numRegions, regions, err := rsmiDevOdVoltCurveRegionsGet(device)
			if err != nil {
				log.Printf("Error getting OD volt curve regions: %v", err)
				fmt.Println(device, fmt.Sprintf("Unable to display %s range", rangeType))
				continue
			}
			for i := 0; i < numRegions; i++ {
				fmt.Println(device, fmt.Sprintf("Region %d: Valid voltage range: %dmV - %dmV",
					i, regions[i], regions[i].VoltRange.UpperBound))
			}
		}
	}

	fmt.Println(" End of Range Display ")
}

// ShowRetiredPages 显示设备列表中指定类型的退役页
// @Summary 显示设备的退役页信息
// @Description 获取并显示指定设备的退役内存页信息。
// @Tags 设备
// @Param dvIdList query []int true "设备ID列表"
// @Param retiredType query string false "退役类型 (默认为'all')"
// @Success 200 {string} string "设备的退役页信息"
// @Failure 400 {string} string "请求参数错误"
// @Router /showRetiredPages [get]
func ShowRetiredPages(dvIdList []int, retiredType string) {
	fmt.Println(" Pages Info ")
	if retiredType == "" {
		retiredType = "all"
	}

	for _, device := range dvIdList {
		_, records, err := rsmiDevMemoryReservedPagesGet(device)
		if err != nil {
			log.Printf("Unable to retrieve reserved page info for device %d: %v", device, err)
			continue
		}

		var data [][]string
		for _, rec := range records {
			status := MemoryPageStatusStr[rec.Status]
			if status == retiredType || retiredType == "all" {
				data = append(data, []string{
					fmt.Sprintf("0x%X", rec.PageAddress),
					fmt.Sprintf("0x%X", rec.PageSize),
					status,
				})
			}
		}

		if len(data) > 0 {
			printTableLog([]string{"Page address", "Page size", "Status"}, data, device, retiredType+" PAGES INFO")
		}
	}
	fmt.Println(" Pages Info ")
}

// ShowSerialNumber 设备序列号
// @Summary 显示设备的序列号
// @Description 获取并显示指定设备的序列号信息。
// @Tags 设备
// @Param dvIdList query []int true "设备ID列表"
// @Success 200 {array} DeviceSerialInfo "设备的序列号信息列表"
// @Failure 400 {string} string "请求参数错误"
// @Router /showSerialNumber [get]
func ShowSerialNumber(dvIdList []int) (deviceSerialInfos []DeviceSerialInfo, err error) {
	fmt.Println("----- Serial Number -----")
	for _, device := range dvIdList {
		serialNumber, err := rsmiDevSerialNumberGet(device)
		deviceSerialInfo := DeviceSerialInfo{
			DeviceID: device,
		}
		if err == nil && serialNumber != "" {
			deviceSerialInfo.SerialNumber = serialNumber
		} else {
			deviceSerialInfo.SerialNumber = "N/A"
		}
		deviceSerialInfos = append(deviceSerialInfos, deviceSerialInfo)
		fmt.Printf("Device %d - Serial Number: %s\n", device, deviceSerialInfo.SerialNumber)
	}
	fmt.Println("------------------------")
	glog.V(5).Infof("deviceSerialInfos:%v", (deviceSerialInfos))
	return
}

// ShowUId 唯一设备ID
// @Summary 显示设备的唯一ID
// @Description 获取并显示指定设备的唯一ID信息。
// @Tags 设备
// @Param dvIdList query []int true "设备ID列表"
// @Success 200 {array} DeviceUIdInfo "设备的唯一ID信息列表"
// @Failure 400 {string} string "请求参数错误"
// @Router /showUId [get]
func ShowUId(dvIdList []int) (deviceUIdInfos []DeviceUIdInfo, err error) {
	fmt.Println("----- Unique ID -----")
	for _, device := range dvIdList {
		uniqueId, err := rsmiDevUniqueIdGet(device)
		deviceUIdInfo := DeviceUIdInfo{
			DeviceID: device,
		}
		if err == nil && uniqueId != 0 {
			deviceUIdInfo.UId = fmt.Sprintf("0x%x", uniqueId)
		} else {
			deviceUIdInfo.UId = "N/A"
		}
		deviceUIdInfos = append(deviceUIdInfos, deviceUIdInfo)
		fmt.Printf("Device %d - Unique ID: %s\n", device, deviceUIdInfo.UId)
	}
	fmt.Println("---------------------")
	return
}

// ShowVbiosVersion 打印并返回设备的VBIOS版本信息
// @Summary 显示设备的VBIOS版本
// @Description 获取并显示指定设备的VBIOS版本信息。
// @Tags 设备
// @Param dvIdList query []int true "设备ID列表"
// @Success 200 {array} DeviceVBIOSInfo "设备的VBIOS版本信息列表"
// @Failure 400 {string} string "请求参数错误"
// @Router /showVbiosVersion [get]
func ShowVbiosVersion(dvIdList []int) (deviceVBIOSInfos []DeviceVBIOSInfo, err error) {
	fmt.Println("----- VBIOS -----")
	for _, device := range dvIdList {
		vbios, err := VbiosVersion(device)
		if err != nil {
			fmt.Printf("Error fetching VBIOS version for device %d: %v\n", device, err)
			deviceVBIOSInfos = append(deviceVBIOSInfos, DeviceVBIOSInfo{
				DeviceID: device,
				VBIOS:    "Error",
			})
		} else {
			fmt.Printf("Device %d VBIOS version: %s\n", device, vbios)
			deviceVBIOSInfos = append(deviceVBIOSInfos, DeviceVBIOSInfo{
				DeviceID: device,
				VBIOS:    vbios,
			})
		}
	}
	fmt.Println("---------------")
	glog.V(5).Infof("deviceVBIOSInfos:%v", (deviceVBIOSInfos))
	return
}

// ShowEvents 显示设备的事件
// @Summary 显示设备的事件
// @Description 获取并显示指定设备的事件信息。
// @Tags 设备
// @Param dvIdList query []int true "设备ID列表"
// @Param eventTypes query []string true "事件类型列表"
// @Success 200 {string} string "成功返回设备的事件信息"
// @Failure 400 {string} string "请求参数错误"
// @Router /showEvents [get]
func ShowEvents(dvIdList []int, eventTypes []string) {
	fmt.Println("----- Show Events -----")
	fmt.Println("Press 'q' or 'ctrl + c' to quit")

	var eventTypeList []string
	for _, event := range eventTypes { // 清理列表中的错误值
		cleanEvent := strings.ReplaceAll(event, ",", "")
		if contains(notificationTypeNames, strings.ToUpper(cleanEvent)) {
			eventTypeList = append(eventTypeList, strings.ToUpper(cleanEvent))
		} else {
			fmt.Printf("Ignoring unrecognized event type %s\n", cleanEvent)
		}
	}

	if len(eventTypeList) == 0 {
		eventTypeList = notificationTypeNames
	}

	var wg sync.WaitGroup
	for _, device := range dvIdList {
		wg.Add(1)
		go func(device int) {
			defer wg.Done()
			printEventList(device, 1000, eventTypeList)
		}(device)
		time.Sleep(250 * time.Millisecond)
	}

	go func() {
		var input string
		for {
			fmt.Scanln(&input)
			if input == "q" {
				for _, device := range dvIdList {
					if err := rsmiEventNotificationStop(device); err != nil {
						fmt.Printf("GPU[%d]: Unable to end event notifications: %v\n", device, err)
					}
				}
				break
			}
		}
	}()

	wg.Wait()
}

// ShowVoltage 当前电压信息
// @Summary 显示设备的电压信息
// @Description 获取并显示指定设备的当前电压信息。
// @Tags 设备
// @Param dvIdList query []int true "设备ID列表"
// @Success 200 {array} DeviceVoltageInfo "设备的电压信息列表"
// @Failure 400 {string} string "请求参数错误"
// @Router /showVoltage [get]
func ShowVoltage(dvIdList []int) (deviceVoltageInfos []DeviceVoltageInfo, err error) {
	for _, device := range dvIdList {
		// 默认电压类型和度量标准
		vtype := RSMI_VOLT_TYPE_FIRST
		met := RSMI_VOLT_CURRENT //
		voltage := rsmiDevVoltMetricGet(device, vtype, met)
		if voltage != 0 {
			fmt.Printf("Device %d: Voltage (mV) = %d\n", device, voltage)
			deviceVoltageInfos = append(deviceVoltageInfos, DeviceVoltageInfo{
				DeviceID: device,
				Voltage:  voltage,
			})
		} else {
			log.Printf("GPU %d voltage not supported\n", device)
		}
	}
	glog.V(5).Infof("deviceVoltageInfos:%v", (deviceVoltageInfos))
	return
}

// ShowVoltageCurve 电压曲线点
// @Summary 显示设备的电压曲线点
// @Description 获取并显示指定设备的电压曲线点信息。
// @Tags 设备
// @Param dvIdList query []int true "设备ID列表"
// @Success 200 {string} string "设备的电压曲线点信息"
// @Failure 400 {string} string "请求参数错误"
// @Router /showVoltageCurve [get]
func ShowVoltageCurve(dvIdList []int) {
	fmt.Println("------------ Voltage Curve Points ------------")
	for _, device := range dvIdList {
		odv, err := rsmiDevOdVoltInfoGet(device)
		if err != nil {
			log.Printf("GPU %d: Voltage Curve is not supported: %v\n", device, err)
			continue
		}

		for position, point := range odv.Curve.VcPoints {
			fmt.Printf("Device %d: Voltage point %d: %d MHz %d mV\n", device, position, point.Frequency/1000000, point.Voltage)
		}
	}
	fmt.Println("----------------------------------------------")
}

// ShowHSLErr 显示指定设备的 HSL 错误状态。
//
// @Summary 显示 HSL 错误状态
// @Description 显示一组 GPU 设备的 HSL 错误状态。
// @Tags Topology
// @Param dvIdList query []int true "设备 ID 列表"
// @Param printJSON query bool false "是否以 JSON 格式输出"
// @Success 200 {string} string "HSL 错误状态信息"
// @Router /showHSLErr [get]
func ShowHSLErr(dvIdList []int, printJSON bool) {
	fmt.Println("------------ HSL Error Status ------------")
	for _, device := range dvIdList {
		status, err := rsmiDevHSLErrorStatus(device)
		if err != nil {
			log.Printf("Error retrieving HSL status for device %d: %v\n", device, err)
			continue
		}

		var desc string
		switch status {
		case RSMIHSLStatusNoErrors:
			desc = "No errors detected since last read"
		case RSMIHSLStatusError:
			desc = "Single error detected since last read"
		case RSMIHSLStatusMultipleErrors:
			desc = "Multiple errors detected since last read"
		default:
			log.Printf("Invalid return value from hsl_error for device %d\n", device)
			continue
		}

		if printJSON {
			fmt.Printf("Device %d: HSL Error count: %d\n", device, status)
		} else {
			fmt.Printf("Device %d: HSL Error count: %d (%s)\n", device, status, desc)
		}
	}
	fmt.Println("-------------------------------------------")
}

// ShowWeightTopology 显示 GPU 拓扑中两台设备之间的权重。
// @Summary 显示 GPU 拓扑权重
// @Description 显示 GPU 设备之间的权重信息。
// @Tags Topology
// @Param dvIdList query []int true "设备 ID 列表"
// @Param printJSON query bool false "是否以 JSON 格式输出"
// @Success 200 {string} string "GPU 拓扑权重信息"
// @Router /showWeightTopology [get]
func ShowWeightTopology(dvIdList []int, printJSON bool) {
	// 初始化矩阵存储设备间的权重
	gpuLinksWeight := make([][]int64, len(dvIdList))
	for i := range gpuLinksWeight {
		gpuLinksWeight[i] = make([]int64, len(dvIdList))
	}

	fmt.Println("------------ Weight between two GPUs ------------")
	for _, srcDevice := range dvIdList {
		for _, destDevice := range dvIdList {
			if srcDevice == destDevice {
				gpuLinksWeight[srcDevice][destDevice] = 0
			} else {
				weight, err := rsmiTopoGetLinkWeight(srcDevice, destDevice)
				if err != nil {
					log.Printf("Cannot read Link Weight between device %d and %d: %v\n", srcDevice, destDevice, err)
					continue
				}
				gpuLinksWeight[srcDevice][destDevice] = weight
			}
		}
	}

	if printJSON {
		formatMatrixToJSON(dvIdList, gpuLinksWeight, "(Topology) Weight between DRM devices %d and %d")
		return
	}

	// 打印矩阵表格
	printTableRow("", "      ")
	for _, row := range dvIdList {
		printTableRow("%-12s", fmt.Sprintf("GPU%d", row))
	}
	fmt.Println()
	for _, gpu1 := range dvIdList {
		printTableRow("%-6s", fmt.Sprintf("GPU%d", gpu1))
		for _, gpu2 := range dvIdList {
			if gpu1 == gpu2 {
				printTableRow("%-12s", "0")
			} else {
				printTableRow("%-12d", gpuLinksWeight[gpu1][gpu2])
			}
		}
		fmt.Println()
	}
	fmt.Println("-------------------------------------------------")
}

// ShowHopsTopology 显示 GPU 拓扑中两台设备之间的跳数。
// @Summary 显示 GPU 拓扑跳数
// @Description 显示 GPU 设备之间的跳数信息。
// @Tags Topology
// @Param dvIdList query []int true "设备 ID 列表"
// @Param printJSON query bool false "是否以 JSON 格式输出"
// @Success 200 {string} string "GPU 拓扑跳数信息"
// @Router /showHopsTopology [get]

func ShowHopsTopology(dvIdList []int, printJSON bool) {
	// 初始化矩阵存储设备间的跳数
	gpuLinksHops := make([][]int64, len(dvIdList))
	for i := range gpuLinksHops {
		gpuLinksHops[i] = make([]int64, len(dvIdList))
	}

	fmt.Println("------------ Hops between two GPUs ------------")
	for _, srcDevice := range dvIdList {
		for _, destDevice := range dvIdList {
			if srcDevice == destDevice {
				gpuLinksHops[srcDevice][destDevice] = 0
			} else {
				hops, _, err := rsmiTopoGetLinkType(srcDevice, destDevice)
				if err != nil {
					log.Printf("Cannot read Link Hops between device %d and %d: %v\n", srcDevice, destDevice, err)
					continue
				}
				gpuLinksHops[srcDevice][destDevice] = hops
			}
		}
	}

	if printJSON {
		formatMatrixToJSON(dvIdList, gpuLinksHops, "(Topology) Hops between DRM devices %d and %d")
		return
	}

	// 打印矩阵表格
	printTableRow("", "      ")
	for _, row := range dvIdList {
		printTableRow("%-12s", fmt.Sprintf("GPU%d", row))
	}
	fmt.Println()
	for _, gpu1 := range dvIdList {
		printTableRow("%-6s", fmt.Sprintf("GPU%d", gpu1))
		for _, gpu2 := range dvIdList {
			if gpu1 == gpu2 {
				printTableRow("%-12s", "0")
			} else {
				printTableRow("%-12d", gpuLinksHops[gpu1][gpu2])
			}
		}
		fmt.Println()
	}
	fmt.Println("-------------------------------------------------")
}

// ShowTypeTopology 显示 GPU 拓扑中两台设备之间的链接类型。
// @Summary 显示 GPU 拓扑链接类型
// @Description 显示 GPU 设备之间的链接类型信息。
// @Tags Topology
// @Param dvIdList query []int true "设备 ID 列表"
// @Param printJSON query bool false "是否以 JSON 格式输出"
// @Success 200 {string} string "GPU 拓扑链接类型信息"
// @Router /showTypeTopology [get]
func ShowTypeTopology(dvIdList []int, printJSON bool) {
	// 初始化矩阵存储设备间的链接类型
	gpuLinksType := make([][]string, len(dvIdList))
	for i := range gpuLinksType {
		gpuLinksType[i] = make([]string, len(dvIdList))
	}

	fmt.Println("------------ Link Type between two GPUs ------------")
	for _, srcDevice := range dvIdList {
		for _, destDevice := range dvIdList {
			if srcDevice == destDevice {
				gpuLinksType[srcDevice][destDevice] = "0"
			} else {
				_, linkType, err := rsmiTopoGetLinkType(srcDevice, destDevice)
				if err != nil {
					log.Printf("Cannot read Link Type between device %d and %d: %v\n", srcDevice, destDevice, err)
					continue
				}
				switch linkType {
				case 1:
					gpuLinksType[srcDevice][destDevice] = LinkTypePCIE
				case 2:
					gpuLinksType[srcDevice][destDevice] = LinkTypeHSL
				default:
					gpuLinksType[srcDevice][destDevice] = LinkTypeUnknown
				}
			}
		}
	}

	if printJSON {
		formatMatrixToStrJSON(dvIdList, gpuLinksType, "(Topology) Link type between DRM devices %d and %d")
		return
	}

	// 打印矩阵表格
	printTableRow("", "      ")
	for _, row := range dvIdList {
		printTableRow("%-12s", fmt.Sprintf("GPU%d", row))
	}
	fmt.Println()
	for _, gpu1 := range dvIdList {
		printTableRow("%-6s", fmt.Sprintf("GPU%d", gpu1))
		for _, gpu2 := range dvIdList {
			if gpu1 == gpu2 {
				printTableRow("%-12s", "0")
			} else {
				printTableRow("%-12s", gpuLinksType[gpu1][gpu2])
			}
		}
		fmt.Println()
	}
	fmt.Println("----------------------------------------------------")
}

// ShowNumaTopology 显示指定设备的 NUMA 节点信息。
// @Summary 显示 NUMA 节点信息
// @Description 显示一组 HCU 设备的 NUMA 节点和关联信息。
// @Tags Topology
// @Param dvIdList query []int true "设备 ID 列表"
// @Success 200 {string} string "NUMA 节点信息"
// @Router /showNumaTopology [get]
func ShowNumaTopology(dvIdList []int) (numaInfos []NumaInfo, err error) {
	for _, device := range dvIdList {
		// 获取 NUMA 节点编号
		numaNode, err := rsmiTopoGetNumaBodeBumber(device)
		if err != nil {
			return nil, err
		}

		// 获取 NUMA 关联信息
		numaAffinity, err := rsmiTopoNumaAffinityGet(device)
		if err != nil {
			return nil, err
		}
		// 将设备和 NUMA 信息存储在结构体中并添加到切片中
		numaInfo := NumaInfo{
			DeviceID:     device,
			NumaNode:     numaNode,
			NumaAffinity: numaAffinity,
		}
		numaInfos = append(numaInfos, numaInfo)
		glog.V(5).Infof("Device %d: Numa Node: %d, Numa Affinity: %d\n", device, numaNode, numaAffinity)
		glog.V(5).Infof("numaInfos: %v", (numaInfos))
	}
	return
}

// ShowHwTopology 显示指定设备的完整硬件拓扑信息。
// @Summary 显示完整的硬件拓扑信息
// @Description 显示一组 GPU 设备的权重、跳数、链接类型和 NUMA 节点信息。
// @Tags Topology
// @Param dvIdList query []int true "设备 ID 列表"
// @Success 200 {string} string "完整的硬件拓扑信息"
// @Router /showHwTopology [get]
func ShowHwTopology(dvIdList []int) {
	ShowWeightTopology(dvIdList, true)

	ShowHopsTopology(dvIdList, true)

	ShowTypeTopology(dvIdList, true)

	ShowNumaTopology(dvIdList)
}

// GetTopoLinkType 获取指定源设备与一组目标设备之间的互联链路类型。
// @Summary 查询设备之间的互联链路类型
// @Description 查询源设备到多个目标设备之间的连接类型（如 PCIe、HSL）。
// @Tags Topology
// @Param srcDvInd query int true "源设备 ID（HCU Index）"
// @Param dstDvIndList query []int true "目标设备 ID 列表（GPU Index）"
// @Success 200 {array} string "链路类型列表（PCIe / HSL / Unknown / 0）"
// @Failure 500 {string} string "查询链路类型失败"
// @Router /topology/linkType [get]
func GetTopoLinkType(srcDvInd int, dstDvIndList []int) (linkTypeList []string, err error) {
	linkTypeList = make([]string, len(dstDvIndList))

	for i, destDevice := range dstDvIndList {

		// 1. src == dst，直接标记并跳过
		if srcDvInd == destDevice {
			linkTypeList[i] = LinkTypeUnknown
			continue
		}

		// 2. 调用 RSMI 接口
		_, linkType, err := rsmiTopoGetLinkType(srcDvInd, destDevice)
		if err != nil {
			glog.V(5).Infof(
				"rsmiTopoGetLinkType failed: src=%v dst=%v err=%v",
				srcDvInd, destDevice, err,
			)
			return nil, err
		}

		switch linkType {
		case 1:
			linkTypeList[i] = LinkTypePCIE
		case 2:
			linkTypeList[i] = LinkTypeHSL
		default:
			linkTypeList[i] = LinkTypeUnknown
		}
	}

	return linkTypeList, nil
}

// TopoIsHylink 判断两个设备之间是否通过 Hylink（HSL）直连。
// @Summary 判断是否为 HSL（Hylink）直连
// @Description 判断源设备与目标设备之间是否存在 HSL 直连。
// @Tags Topology
// @Param srcDvInd query int true "源设备 ID（HCU Index）"
// @Param dstDvInd query int true "目标设备 ID（HCU Index）"
// @Success 200 {boolean} bool "是否为 Hylink（HSL）直连"
// @Failure 500 {string} string "查询 Hylink 状态失败"
// @Router /topology/isHylink [get]
func TopoIsHylink(srcDvInd, dstDvInd int) (bool, error) {
	return rsmiTopoIsHylink(srcDvInd, dstDvInd)
}

// XhclLinkStates 获取指定 HCU 的所有 XHCL（HSL）链路状态。
// @Summary 获取 GPU 的 XHCL 链路状态
// @Description 查询指定 HCU 的所有 XHCL 链路当前状态。
// @Tags Topology
// @Param dvInd query int true "设备 ID（GPU Index）"
// @Success 200 {array} XhclLinkState "XHCL 链路状态列表"
// @Failure 500 {string} string "查询 XHCL 链路状态失败"
func XhclLinkStates(dvInd int) ([]XhclLinkState, error) {
	linkNum, err := rsmiDevXhclLinkNumber(dvInd)
	if err != nil {
		return nil, err
	}

	states := make([]XhclLinkState, 0, linkNum)

	glog.V(5).Infof(
		"GetXhclLinkStates dvInd: %v linkNum: %v",
		dvInd,
		linkNum,
	)

	for linkID := 0; linkID < linkNum; linkID++ {
		state, err := rsmiDevXhclLinkState(dvInd, linkID)
		if err != nil {
			glog.V(5).Infof(
				"GetXhclLinkStates dvInd: %v linkID: %v error: %v",
				dvInd,
				linkID,
				err,
			)
			continue
		}

		linkState := XhclLinkState{
			LinkID: linkID,
			State:  state,
			Up:     state == XhclLinkUp,
		}

		glog.V(5).Infof(
			"GPU %d XHCL link %d/%d state=%d up=%v",
			dvInd,
			linkID,
			linkNum,
			state,
			linkState.Up,
		)

		states = append(states, linkState)
	}

	return states, nil
}

// DumpXhclRemoteBdfids 获取指定 GPU 上所有 UP 状态的 XHCL 链路对应的远端设备 BDF ID。
//
// @Summary 获取 GPU 的 XHCL 远端设备 BDF 信息
// @Description 枚举指定 GPU 上所有处于 UP 状态的 XHCL 链路，并返回每条链路所连接的远端设备 BDF ID。
// @Tags Topology
// @Param dvInd query int true "设备 ID（GPU Index）"
// @Success 200 {array} XhclRemoteBdf "XHCL 链路远端设备 BDF 列表"
// @Failure 500 {string} string "查询 XHCL 远端 BDF 信息失败"
func DumpXhclRemoteBdfids(dvInd int) ([]XhclRemoteBdf, error) {
	states, err := XhclLinkStates(dvInd)
	if err != nil {
		return nil, err
	}

	results := make([]XhclRemoteBdf, 0)

	for _, s := range states {
		// 只处理 UP 状态的链路
		if !s.Up {
			continue
		}

		bdfid, err := rsmiXhclLinkRemoteBdfidGet(dvInd, s.LinkID)
		if err != nil {
			glog.Warningf(
				"HCU %d link %d get remote bdfid failed: %v",
				dvInd, s.LinkID, err,
			)
			continue
		}

		glog.V(5).Infof(
			"HCU %d XHCL link %d -> remote BDFID: 0x%x",
			dvInd,
			s.LinkID,
			bdfid,
		)

		results = append(results, XhclRemoteBdf{
			LinkID: s.LinkID,
			BdfID:  bdfid,
		})
	}

	return results, nil
}

// DiscoverInterconnectTopology 枚举整机 HCU 的互联关系，返回 HCU × HCU 的互联矩阵。
//
// @Summary 获取整机 HCU 的互联矩阵信息
// @Description 枚举整机 HCU 互联关系，包括链路类型（PCIe / HSL / HYSWITCH / NONE）及对应权重。
//   - 自动获取 HCU 数量
//   - 构建 HCU × HCU 的互联矩阵
//   - 判断每一对 HCU 之间的链路类型
//   - 计算对应的链路权重（PCIe / NUMA / HSL / HYSWITCH）
//
// @Tags Topology
// @Success 200 {object} HCUInterconnectMatrix "HCU 互联矩阵信息"
// @Failure 500 {string} string "查询 HCU 互联信息失败"
func DiscoverInterconnectTopology() (matrix HCUInterconnectMatrix, err error) {
	deviceCount, err := NumMonitorDevices()
	if err != nil {
		return matrix, err
	}

	matrix.DeviceCount = deviceCount
	glog.V(5).Infof("DiscoverInterconnectTopology start, deviceCount=%d", deviceCount)

	matrix.Matrix = make([][]HCULinkInfo, deviceCount)
	for i := 0; i < deviceCount; i++ {
		matrix.Matrix[i] = make([]HCULinkInfo, deviceCount)
	}

	dvIndToBdf := make(map[int]uint64)
	for i := 0; i < deviceCount; i++ {
		bdfid, err := rsmiDevPciIdGet(i)
		if err != nil {
			return matrix, err
		}
		dvIndToBdf[i] = uint64(bdfid)

		glog.V(5).Infof(
			"HCU %d BDFID=0x%x",
			i, dvIndToBdf[i],
		)
	}

	for src := 0; src < deviceCount; src++ {
		for dst := 0; dst < deviceCount; dst++ {
			linkInfo := HCULinkInfo{
				SrcDvInd: src,
				DstDvInd: dst,
				BdfID:    dvIndToBdf[dst],
				LinkType: LinkTypeNONE,
				Weight:   0,
				Hops:     0,
			}

			if src == dst {
				linkInfo.Weight = -1
				matrix.Matrix[src][dst] = linkInfo
				continue
			}

			linkTypes, err := GetTopoLinkType(src, []int{dst})
			if err != nil {
				return matrix, err
			}

			if len(linkTypes) == 0 {
				matrix.Matrix[src][dst] = linkInfo
				continue
			}

			switch linkTypes[0] {
			case LinkTypePCIE:
				linkInfo.LinkType = LinkTypePCIE

				numaInfos, err := ShowNumaTopology([]int{src, dst})
				if err != nil {
					return matrix, err
				}

				if len(numaInfos) == 2 &&
					numaInfos[0].NumaNode == numaInfos[1].NumaNode {
					linkInfo.Weight = 1
				} else {
					linkInfo.Weight = 0
				}

				glog.V(5).Infof(
					"HCU %d -> %d PCIE link, weight=%d",
					src, dst, linkInfo.Weight,
				)

			case LinkTypeHSL:
				isHyswitch, err := TopoIsHylink(src, dst)
				if err != nil {
					return matrix, err
				}

				if isHyswitch {
					linkInfo.LinkType = LinkTypeHSLHyswitch
					linkInfo.Weight = deviceCount - 1

					glog.V(5).Infof(
						"HCU %d -> %d HYSWITCH link, weight=%d",
						src, dst, linkInfo.Weight,
					)
				} else {
					linkInfo.LinkType = LinkTypeHSL

					xhclLinks, err := DumpXhclRemoteBdfids(src)
					if err != nil {
						return matrix, err
					}

					for _, l := range xhclLinks {
						if l.BdfID == dvIndToBdf[dst] {
							linkInfo.Weight++
						}
					}

					glog.V(5).Infof(
						"HCU %d -> %d HSL link, weight=%d",
						src, dst, linkInfo.Weight,
					)
				}

			default:
				linkInfo.LinkType = LinkTypeUnknown
			}

			matrix.Matrix[src][dst] = linkInfo
		}
	}

	for i := range matrix.Matrix {
		for j := range matrix.Matrix[i] {
			matrix.Matrix[i][j].PciID = formatBDFID(matrix.Matrix[i][j].BdfID)
		}
	}

	return matrix, nil
}

// GetGpuInterconnectInfo 获取指定 GPU 的互联拓扑信息
//
// 该函数从单卡视角出发，综合 PCIe / HSL / Hyswitch 信息，
// 判断 GPU 的互联类型、互联规模（几卡互联），以及与各直连 GPU
// 之间的链路权重（XHCL link 数量）。
//
// @param srcDvInd     当前 GPU 的设备索引
// @param dstDvIndList 系统中其他 GPU 的设备索引列表
//
// @return info  GPU 的互联拓扑信息
// @return err   查询或分析过程中发生的错误
func GetGpuInterconnectInfo(srcDvInd int, dstDvIndList []int) (info GpuInterconnectInfo, err error) {
	glog.V(5).Infof("GetGpuInterconnectInfo start, srcDvInd=%v, dstDvIndList=%v", srcDvInd, dstDvIndList)

	// ---------- 0. 建立 BDFID -> dvInd 映射 ----------
	bdfToDvInd := make(map[uint64]int)
	for _, dv := range dstDvIndList {
		bdfid, err := rsmiDevPciIdGet(dv)
		if err != nil {
			return info, err
		}
		bdfToDvInd[uint64(bdfid)] = dv

	}

	// ---------- 1. 判断每个 dst GPU 的链路类型 ----------
	pcieLinks := make([]int, 0) // 保存 PCIe 直连的 GPU 索引
	hslLinks := make([]int, 0)  // 保存 HSL 直连的 GPU 索引
	isHyswitch := false         // 是否有 Hyswitch

	for _, dst := range dstDvIndList {
		if dst == srcDvInd {
			continue
		}

		linkTypes, err := GetTopoLinkType(srcDvInd, []int{dst})
		if err != nil {
			return info, err
		}
		if len(linkTypes) == 0 {
			continue
		}

		switch linkTypes[0] {
		case LinkTypePCIE:
			pcieLinks = append(pcieLinks, dst)
			glog.V(5).Infof("GPU %v direct PCIe link detected with GPU %v", srcDvInd, dst)
		case LinkTypeHSL:
			hslLinks = append(hslLinks, dst)
			hylink, err := TopoIsHylink(srcDvInd, dst)
			if err != nil {
				return info, err
			}
			if hylink {
				isHyswitch = true
			}
		default:
			glog.V(5).Infof("GPU %v unknown link type with GPU %v", srcDvInd, dst)
		}
	}

	// ---------- 2. 处理 HSL / Hyswitch 链路 ----------
	bdfWeightMap := make(map[uint64]*DirectLinkInfo) // BDFID -> DirectLinkInfo
	if len(hslLinks) > 0 {
		xhclLinks, err := DumpXhclRemoteBdfids(srcDvInd)
		if err != nil {
			return info, err
		}

		for _, link := range xhclLinks {
			if infoLink, ok := bdfWeightMap[link.BdfID]; ok {
				infoLink.Weight++
			} else {
				dvInd := -1
				if v, ok := bdfToDvInd[link.BdfID]; ok {
					dvInd = v
				}
				bdfWeightMap[link.BdfID] = &DirectLinkInfo{
					RemoteDvInd: dvInd,
					RemoteBdfID: link.BdfID,
					LinkType:    LinkTypeHSL,
					Weight:      1,
				}
			}
		}
	}

	// ---------- 3. 处理 PCIe 链路 ----------
	for _, dst := range pcieLinks {
		bdfid, err := rsmiDevPciIdGet(dst)
		if err != nil {
			glog.V(5).Infof("rsmiDevPciIdGet failed for dvInd=%v, err=%v", dst, err)
			continue
		}

		bdfWeightMap[uint64(bdfid)] = &DirectLinkInfo{
			RemoteDvInd: dst,
			RemoteBdfID: uint64(bdfid),
			LinkType:    LinkTypePCIE,
			Weight:      1,
		}
	}

	// ---------- 4. 构造返回 Links ----------
	for _, link := range bdfWeightMap {
		info.Links = append(info.Links, *link)
	}

	// ---------- 5. 填充总体互联类型 ----------
	switch {
	case len(pcieLinks) > 0 && len(hslLinks) > 0:
		info.LinkType = LinkTypeHybrid
	case len(hslLinks) > 0 && isHyswitch:
		info.LinkType = LinkTypeHSLHyswitch
	case len(hslLinks) > 0:
		info.LinkType = LinkTypeHSL
	case len(pcieLinks) > 0:
		info.LinkType = LinkTypePCIE
	default:
		info.LinkType = LinkTypeUnknown
	}

	// ---------- 6. 填充互联规模 ----------
	info.CardCount = len(info.Links) + 1 // +1 表示自身

	glog.V(5).Infof("GetGpuInterconnectInfo done, srcDvInd=%v, linkType=%v, cardCount=%v",
		srcDvInd, info.LinkType, info.CardCount)

	return info, nil
}

/*************************************VHCU******************************************/
// DeviceCount 返回设备的数量。
// @Summary 获取设备数量
// @Description 获取当前系统中的设备数量。
// @Tags Device
// @Success 200 {int} int "设备数量"
// @Failure 500 {object} string "内部服务器错误"
// @Router /deviceCount [get]
func DeviceCount() (count int, err error) {
	return dmiGetDeviceCount()
}

// VDeviceSingleInfo
// @Summary 获取单个虚拟设备的信息
// @Description 根据设备索引获取对应的虚拟设备信息
// @Tags VirtualDevice
// @Param vDvInd query int true "设备索引"
// @Success 200 {object} VDeviceInfo "虚拟设备信息"
// @Failure 400 {string} string "请求参数错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /VDeviceSingleInfo [get]
func VDeviceSingleInfo(vDvInd int) (vDeviceInfo VDeviceInfo, err error) {
	glog.V(5).Infof("VDeviceSingleInfo vDvInd:%v", vDvInd)
	vDeviceInfo, err = dmiGetVDeviceInfo(vDvInd)
	devTypeName, _, _ := DevTypeName(vDeviceInfo.DvInd)
	vDeviceInfo.Name = NormalizeDevTypeName(devTypeName)
	return
}

// VDeviceInfos 获取所有虚拟设备信息
func VDeviceInfos() (vDeviceInfos []VDeviceInfo, err error) {
	numDevices, err := rsmiNumMonitorDevices()
	if err != nil {
		return nil, err
	}

	// 假设每个物理设备最多 4 个虚拟设备
	vDeviceCount := numDevices * 4

	for i := 0; i < vDeviceCount; i++ {
		vDeviceInfo, err := dmiGetVDeviceInfo(i)
		if err == nil {
			vDeviceInfo.VdvInd = i

			bdfid, err := rsmiDevPciIdGet(vDeviceInfo.DvInd)
			if err != nil {
				return nil, err
			}

			devTypeName, _, _ := DevTypeName(vDeviceInfo.DvInd)
			vDeviceInfo.Name = NormalizeDevTypeName(devTypeName)

			vPercent, _ := dmiGetVDevBusyPercent(i)
			vDeviceInfo.VPercent = vPercent

			// 解析 PCI BDF
			domain := (bdfid >> 32) & 0xffffffff
			bus := (bdfid >> 8) & 0xff
			dev := (bdfid >> 3) & 0x1f
			function := bdfid & 0x7

			vDeviceInfo.PciBusNumber = fmt.Sprintf(
				"%04x:%02x:%02x.%x",
				domain, bus, dev, function,
			)

			vDeviceInfos = append(vDeviceInfos, vDeviceInfo)
		}
	}

	glog.V(5).Infof("vDeviceInfos: %+v", vDeviceInfos)
	return
}

// VDeviceCount 返回虚拟设备的数量。
// @Summary 获取虚拟设备数量
// @Description 获取当前系统中的虚拟设备数量。
// @Tags Device
// @Success 200 {int} int "虚拟设备数量"
// @Failure 500 {object} string "内部服务器错误"
// @Router /vDeviceCount [get]
func VDeviceCount() (count int, err error) { return dmiGetVDeviceCount() }

// DeviceRemainingInfo 返回指定物理设备的剩余计算单元（CU）和内存信息。
// @Summary 获取设备剩余信息
// @Description 获取指定设备的剩余计算单元和内存信息。
// @Tags Device
// @Param dvInd path int true "物理设备索引"
// @Success 200 {string} uint64 "剩余的CU信息"
// @Success 200 {string} uint64 "剩余的内存信息"
// @Failure 400 {object} string "无效的设备索引"
// @Failure 500 {object} string "内部服务器错误"
// @Router /deviceRemainingInfo/{dvInd} [get]
func DeviceRemainingInfo(dvInd int) (cus, memories uint64, err error) {
	return dmiGetDeviceRemainingInfo(dvInd)
}

// CreateVDevices 创建指定数量的虚拟设备
// @Summary 创建虚拟设备
// @Description 在指定的物理设备上创建指定数量的虚拟设备，返回创建的虚拟设备ID集合。
// @Tags 虚拟设备
// @Param dvInd query int true "物理设备的索引"
// @Param vDevCount query int true "要创建的虚拟设备数量"
// @Param vDevCUs query []int true "每个虚拟设备的计算单元数量"
// @Param vDevMemSize query []int true "每个虚拟设备的内存大小"
// @Success 200 {array} int "虚拟设备创建成功，返回虚拟设备ID集合"
// @Failure 400 {string} string "创建虚拟设备失败"
// @Router /CreateVDevices [post]
func CreateVDevices(dvInd int, vDevCount int, vDevCUs []int, vDevMemSize []int) (vdevIDs []int, err error) {
	return dmiCreateVDevices(dvInd, vDevCount, vDevCUs, vDevMemSize)
}

// DestroyVDevice 销毁指定物理设备上的所有虚拟设备
// @Summary 销毁所有虚拟设备
// @Description 销毁指定物理设备上的所有虚拟设备。
// @Tags 虚拟设备
// @Param dvInd query int true "物理设备的索引"
// @Success 200 {string} string "虚拟设备销毁成功"
// @Failure 400 {string} string "虚拟设备销毁失败"
// @Router /DestroyVDevice [delete]
func DestroyVDevice(dvInd int) (err error) {
	return dmiDestroyVDevices(dvInd)
}

// DestroySingleVDevice 销毁指定虚拟设备
// @Summary 销毁单个虚拟设备
// @Description 销毁指定索引的虚拟设备。
// @Tags 虚拟设备
// @Param vDvInd query int true "虚拟设备的索引"
// @Success 200 {string} string "虚拟设备销毁成功"
// @Failure 400 {string} string "虚拟设备销毁失败"
// @Router /DestroySingleVDevice [delete]
func DestroySingleVDevice(vDvInd int) (err error) {
	return dmiDestroySingleVDevice(vDvInd)
}

// UpdateSingleVDevice 更新指定设备资源大小
// @Summary 更新虚拟设备资源
// @Description 更新指定虚拟设备的计算单元和内存大小。如果 vDevCUs 或 vDevMemSize 为 -1，则对应的资源不更改。
// @Tags 虚拟设备
// @Param vDvInd query int true "虚拟设备的索引"
// @Param vDevCUs query int true "更新后的计算单元数量"
// @Param vDevMemSize query int true "更新后的内存大小"
// @Success 200 {string} string "虚拟设备更新成功"
// @Failure 400 {string} string "虚拟设备更新失败"
// @Router /UpdateSingleVDevice [put]
func UpdateSingleVDevice(vDvInd int, vDevCUs int, vDevMemSize int) (err error) {
	return dmiUpdateSingleVDevice(vDvInd, vDevCUs, vDevMemSize)
}

// StartVDevice 启动虚拟设备
// @Summary 启动指定的虚拟设备
// @Description 启动虚拟设备，指定设备索引
// @Tags VirtualDevice
// @Param vDvInd path int true "虚拟设备索引"
// @Success 200 {string} string "操作成功"
// @Failure 400 {string} string "操作失败"
// @Router /StartVDevice/{vDvInd} [get]
func StartVDevice(vDvInd int) (err error) {
	return dmiStartVDevice(vDvInd)
}

// DevBusyPercent 返回物理设备时间忙碌百分比
// @Summary 返回物理设备时间忙碌百分比
// @Description 返回物理设备时间忙碌百分比
// @Produce json
// @Param dvInd query int true "设备索引"
// @Success 200 {int} float64 "返回设备的当前时间忙碌百分比"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /DevBusyPercent [get]
func DevBusyPercent(dvInd int) (utilizationRate float64, err error) {
	utilization, _ := rsmiDevBusyPercentGet(dvInd)
	utilizationRate, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", float64(utilization)/1.0), 64)
	return
}

// HCUCuUsage 获取指定设备的 HCU 瞬时占用率。
//
// @Summary 获取 HCU 瞬时占用率
// @Description 查询指定设备当前的 HCU 占用率（瞬时值）。只要某个 CU 内存在至少一个活跃 wave，即认为该 CU 活跃；占用率 = 活跃 CU 数量 / CU 总数量 × 100%。对应 hy-smi -u / rsmi_dev_cu_usage_get。
// @Tags Utilization
// @Produce json
// @Param dvInd path int true "设备索引，从 0 开始，对应 hy-smi -d"
// @Success 200 {number} float64 "HCU 瞬时占用率"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /HCUCuUsage/{dvInd} [get]
//
// 参数说明：
//   - dvInd：物理设备索引，与 NumMonitorDevices 枚举顺序一致
//
// 返回值说明：
//   - utilizationRate：HCU 瞬时占用率，浮点数；GPU 空闲时通常为 0
//   - err：调用失败时返回错误信息（如 RSMI_STATUS_NOT_SUPPORTED）
func HCUCuUsage(dvInd int) (utilizationRate float64, err error) {
	utilization, err := rsmiDevCuUsageGet(dvInd)
	if err != nil {
		return 0, err
	}
	utilizationRate, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", float64(utilization)), 64)
	return
}

// HCUSampledUsage 获取指定设备在过去采样窗口内的 HCU 占用情况。
//
// @Summary 获取 HCU 采样占用情况
// @Description 在 sampleDurationMs 时间窗口内周期性采样 HCU 活跃状态，返回活跃采样次数占比。对应 hy-smi --showhcuutil / rsmi_dev_hcu_util_get；hy-smi 默认采样窗口为 1s（1000ms），采样间隔约 10ms。
// @Tags Utilization
// @Produce json
// @Param dvInd path int true "设备索引，从 0 开始"
// @Param sampleDurationMs query int false "采样时间窗口（毫秒），默认 1000"
// @Success 200 {number} float64 "HCU 采样占用率"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /HCUSampledUsage/{dvInd} [get]
//
// 参数说明：
//   - dvInd：物理设备索引
//   - sampleDurationMs：采样时间窗口（毫秒）；hy-smi 命令行等效为 1s，建议传 1000
//
// 返回值说明：
//   - utilizationRate：HCU 活跃采样占比，即 HCU 为活跃状态的次数 / 采样总次数
//   - err：调用失败时返回错误信息
func HCUSampledUsage(dvInd int, sampleDurationMs int) (utilizationRate float64, err error) {
	utilization, err := rsmiDevHcuUtilGet(dvInd, sampleDurationMs)
	if err != nil {
		return 0, err
	}
	utilizationRate, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", float64(utilization)), 64)
	return
}

// HCUCUSampledUsage 获取指定设备 CU 在过去采样窗口内的占用情况。
//
// @Summary 获取 CU 采样占用情况
// @Description 在 sampleDurationMs 窗口内统计各 CU 活跃采样占比并取全 CU 平均值。只要 CU 内存在至少一个活跃 wave 即视为活跃。对应 hy-smi --showcuutil / rsmi_dev_cu_util_get。
// @Tags Utilization
// @Produce json
// @Param dvInd path int true "设备索引，从 0 开始"
// @Param sampleDurationMs query int false "采样时间窗口（毫秒），默认 1000"
// @Success 200 {number} float64 "CU 平均采样占用率"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /HCUCUSampledUsage/{dvInd} [get]
//
// 参数说明：
//   - dvInd：物理设备索引
//   - sampleDurationMs：采样时间窗口（毫秒），建议 1000（与 hy-smi 默认 1s 一致）
//
// 返回值说明：
//   - utilizationRate：所有 CU 活跃采样占比的平均值
//   - err：调用失败时返回错误信息
func HCUCUSampledUsage(dvInd int, sampleDurationMs int) (utilizationRate float64, err error) {
	utilization, err := rsmiDevCuUtilGet(dvInd, sampleDurationMs)
	if err != nil {
		return 0, err
	}
	utilizationRate, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", float64(utilization)), 64)
	return
}

// HCUWaveSampledUsage 获取指定设备 Wave 在过去采样窗口内的占用情况。
//
// @Summary 获取 Wave 采样占用情况
// @Description 在 sampleDurationMs 窗口内统计各 CU 上活跃 wave 占 wave 总数的比例，并对所有 CU 取平均。对应 hy-smi --showwaveutil / rsmi_dev_wave_util_get。
// @Tags Utilization
// @Produce json
// @Param dvInd path int true "设备索引，从 0 开始"
// @Param sampleDurationMs query int false "采样时间窗口（毫秒），默认 1000"
// @Success 200 {number} float64 "Wave 平均采样占用率"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /HCUWaveSampledUsage/{dvInd} [get]
//
// 参数说明：
//   - dvInd：物理设备索引
//   - sampleDurationMs：采样时间窗口（毫秒），建议 1000
//
// 返回值说明：
//   - utilizationRate：活跃 wave 采样占比在所有 CU 上的平均值
//   - err：调用失败时返回错误信息
func HCUWaveSampledUsage(dvInd int, sampleDurationMs int) (utilizationRate float64, err error) {
	utilization, err := rsmiDevWaveUtilGet(dvInd, sampleDurationMs)
	if err != nil {
		return 0, err
	}
	utilizationRate, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", float64(utilization)), 64)
	return
}

// HCUSEUsage 获取指定设备各 Shader Engine（SE）的瞬时占用率。
//
// @Summary 获取 SE 瞬时占用率
// @Description 按 SE 维度返回活跃 CU 占比（瞬时值）。只要 CU 内存在至少一个活跃 wave 即视为活跃；每个 SE 占用率 = SE 内活跃 CU 数 / 该 SE 内 CU 总数 × 100%。对应 hy-smi --showseuse / rsmi_dev_se_util_get。
// @Tags Utilization
// @Produce json
// @Param dvInd path int true "设备索引，从 0 开始"
// @Success 200 {object} SEUsageInfo "各 SE 占用率"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /HCUSEUsage/{dvInd} [get]
//
// 参数说明：
//   - dvInd：物理设备索引
//
// 返回值说明：
//   - shaderEngineUsage：SE 占用信息，Percent[i] 为第 i 个 SE 的占用率；数组长度 MAX_SE_CNT（8）
//   - err：调用失败时返回错误信息
func HCUSEUsage(dvInd int) (shaderEngineUsage SEUsageInfo, err error) {
	return rsmiDevSeUtilGet(dvInd)
}

// VDevBusyPercent 返回虚拟设备使用百分比
// @Summary 获取虚拟设备使用百分比
// @Description 根据虚拟设备索引返回当前虚拟设备的使用百分比
// @Produce json
// @Param vDvInd query int true "虚拟设备索引"
// @Success 200 {int} int "返回虚拟设备的当前使用百分比"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /VDevBusyPercent [get]
func VDevBusyPercent(vDvInd int) (percent int, err error) {
	return dmiGetVDevBusyPercent(vDvInd)
}

// StopVDevice 停止虚拟设备
// @Summary 停止指定的虚拟设备
// @Description 停止虚拟设备，指定设备索引
// @Tags VirtualDevice
// @Param vDvInd path int true "虚拟设备索引"
// @Success 200 {string} string "操作成功"
// @Failure 400 {string} string "操作失败"
// @Router /StopVDevice/{vDvInd} [get]
func StopVDevice(vDvInd int) (err error) {
	return dmiStopVDevice(vDvInd)
}

// SetEncryptionVMStatus 设置虚拟机加密状态
// @Summary 设置虚拟机加密状态
// @Description 根据提供的状态开启或关闭虚拟机加密
// @Tags VirtualDevice
// @Param status query bool true "加密状态"
// @Success 200 {string} string "操作成功"
// @Failure 400 {string} string "操作失败"
// @Router /SetEncryptionVMStatus [post]
func SetEncryptionVMStatus(status bool) (err error) {
	return dmiSetEncryptionVMStatus(status)
}

// EncryptionVMStatus 获取加密虚拟机状态
// @Summary 获取当前虚拟机的加密状态
// @Description 返回虚拟机是否处于加密状态
// @Tags VirtualDevice
// @Success 200 {boolean} boolean "加密状态"
// @Failure 400 {string} string "操作失败"
// @Router /EncryptionVMStatus [get]
func EncryptionVMStatus() (status bool, err error) {
	return dmiGetEncryptionVMStatus()
}

// PrintEventList 打印事件列表
// @Summary 打印设备的事件列表
// @Description 打印指定设备的事件列表，并设置延迟
// @Tags Event
// @Param dvInd path int true "设备索引"
// @Param delay query int true "延迟时间（秒）"
// @Param eventList query []string true "事件列表"
// @Success 200 {string} string "操作成功"
// @Failure 400 {string} string "操作失败"
// @Router /PrintEventList/{device} [get]
func PrintEventList(dvInd int, delay int, eventList []string) {
	printEventList(dvInd, delay, eventList)
}

// GetDeviceInfo 获取设备信息
// @Summary 获取设备信息
// @Description 根据设备索引返回设备的详细信息
// @Produce json
// @Param dvInd query int true "设备索引"
// @Success 200 {object} DMIDeviceInfo "返回包含设备详细信息的对象"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /GetDeviceInfo [get]
func GetDeviceInfo(dvInd int) (deviceInfo DMIDeviceInfo, err error) {
	return dmiGetDeviceInfo(dvInd)
}

// RunDiag 运行设备诊断
// @Summary 运行设备诊断
// @Description 执行设备的诊断测试并返回诊断结果
// @Produce json
// @Success 200 {object} DiagResults "返回诊断测试结果"
// @Failure 500 {object} error "服务器内部错误"
// @Router /RunDiag [get]
func RunDiag(level int) (diagResults DiagResults, err error) {
	return runDiag(level)
}

// StopDiag 请求停止诊断
func StopDiag() {
	atomic.StoreInt32(&stopDiagFlag, 1)
}

// BandwidthTest 显存带宽测试（CLI 用：会打印 banner）
func BandwidthTest(dvIdList []int) bool {
	// 创建一个可被取消的上下文
	ctx := context.Background()

	// verbose=true -> runBandwidthTest 会打印 "===== 带宽测试结果 =====" banner
	_, err := runBandwidthTest(ctx, dvIdList, true)
	if err != nil {
		glog.Errorf("带宽测试被取消或出错: %v", err)
		return false
	}

	return true
}

// BandwidthTest
// @Summary      设备带宽测试
// @Description  对指定设备执行带宽测试，返回每个设备的带宽结果。
// @Description  当 error 不为空时，返回结果可能仍包含设备带宽数据。
// @Tags         Bandwidth
// @Accept       json
// @Produce      json
// @Success      200 {object} BandwidthTestResp
// @Failure      200 {object} Response
// @Router       /bandwidth/test [get]

func BandwidthTestResult(dvIdList []int) (map[int]float64, error) {
	ctx := context.Background()

	// verbose=false -> runBandwidthTest 不会打印到 stdout（适合 diag 调用）
	bwMap, err := runBandwidthTest(ctx, dvIdList, false)
	if err != nil {
		glog.Errorf("带宽测试被取消或出错: %v", err)
		return nil, err
	}

	return bwMap, nil
}

// PcieBandwidthTest PCIe带宽测试
func PcieBandwidthTest() bool {
	return runPcieBandwidthTest()
}

// 对外 API：返回结构化结果。
func PcieBandwidthTestResult() (PcieResult, error) {
	return runPcieBandwidthTestWithResult()
}

// XHCLTest xHCL带宽测试
func XHCLTest() bool {
	return hcuXHCLTest()
}

// XHCLTestResult 对外提供结构化结果接口
func XHCLTestResult() ([]XHCLResult, error) {
	return runHcuXHCLTestWithResult()
}

// TargetStressTest target stress-Gemm压力测试
func TargetStressTest() {
	targetStressTest()
}

// TargetStressTestResult 对外接口：运行 TargetStress 测试并返回结构化结果
func TargetStressTestResult() (TargetStressResult, error) {
	return runTargetStressTestWithResult()
}

// MemtestCL  Memtestcl压力测试，测试显存稳定性
func MemtestCL(dvIdList []int) error {
	return memtestCL(dvIdList)
}

// MemtestCLTestResult 对外接口：运行 memtestCL 并返回结构化结果（保留原 MemtestCL 不变）
func MemtestCLTestResult(dvIdList []int) (MemtestCLAllResult, error) {
	return runMemtestCLWithResult(dvIdList)
}

// EDPpTest 测试板卡稳定性
func EDPpTest() {
	edpppTest()
}

// EDPpTestResult 对外接口：运行 EDPp 测试并以结构化形式返回结果。
// 不传 dvIdList 时保持生产默认行为：测试全部 HCU。
func EDPpTestResult(dvIdList ...int) (EDPPResult, error) {
	return runEdppTestWithResult(dvIdList...)
}

/**************health ****************/
// @Summary 设置健康检查配置
// @Description 根据参数设置是否开启健康检查，并选择需要的检查项
// @Produce json
// @Param enabled query bool true "是否开启健康检查"
// @Param options query []int true "检查项选项，包含如下检查项：1-NumaTopology Health, 2-PcieBandwidth Health, 3-Power Health, 4-Memory Health, 5-Temperature Health, 6-Performance Health, 7-EccBlocks Health, 8-HCUUsage Health"
// @Success 200 {string} string "健康检查配置已更新"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /SetHealthCheckConfig [post]
func SetHealthCheckConfig(enabled bool, options []int) (err error) {
	var validOptions []string
	var invalidOptions []int
	// 校验并转换 options 数字
	for _, opt := range options {
		if checkType, exists := HealthType[opt]; exists {
			validOptions = append(validOptions, checkType)
		} else {
			invalidOptions = append(invalidOptions, opt)
		}
	}
	// 如果有不符合要求的数字，返回错误
	if len(invalidOptions) > 0 {
		return fmt.Errorf("无效的选项: %v", invalidOptions)
	}

	// 所有选项都通过校验，调用 setHealthCheckConfig
	return setHealthCheckConfig(enabled, validOptions)
}

// GetHealthCheckConfig 获取健康检查项信息
// @Summary 获取健康检查项配置信息
// @Description 返回当前健康检查的启用状态和已选择的检查项
// @Produce json
// @Success 200 {object} HealthCheckConfig "返回当前健康检查配置，包括启用状态和检查项列表"
// @Failure 500 {object} error "服务器内部错误"
// @Router /GetHealthCheckConfig [get]
func GetHealthCheckConfig() (healthConfig HealthCheckConfig, err error) {
	return getHealthCheckConfig()
}

// DeleteHealthCheckConfig 删除健康检查项信息
// @Summary 删除健康检查配置
// @Description 清除当前健康检查的配置信息
// @Produce json
// @Success 200 {string} string "健康检查配置已删除"
// @Failure 500 {object} error "服务器内部错误"
// @Router /DeleteHealthCheckConfig [delete]
func DeleteHealthCheckConfig() (err error) {
	return deleteHealthCheckConfig()
}

// HealthCheckConfigList 列出所有健康检查项
// @Summary 列出健康检查项
// @Description 返回所有可用的健康检查项及其配置信息
// @Produce json
// @Success 200 {object} map[string]interface{} "返回健康检查项列表"
// @Failure 500 {object} error "服务器内部错误"
// @Router /HealthCheckConfigList [get]
func HealthCheckConfigList() (map[string]interface{}, error) {
	return healthCheckConfigList()
}

// HealthCheckById 设备健康检查
// @Summary 获取指定设备的健康检查结果
// @Description 根据设备 ID 列表返回对应设备的健康检查状态信息
// @Produce json
// @Param dvIdList query []int true "设备 ID 列表"
// @Success 200 {array} DeviceHealth "返回指定设备的健康检查结果"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /HealthCheckById [get]
func HealthCheckById(dvIdList []int, checkHealthConfig bool) (deviceHealths []DeviceHealth, err error) {
	return healthCheckById(dvIdList, checkHealthConfig)
}

// HealthCheckByGroupId 设备组健康检查
// @Summary 获取指定设备组的健康检查结果
// @Description 根据设备组 ID 返回该组内所有设备的健康检查状态信息
// @Produce json
// @Param groupId query string true "设备组 ID"
// @Param checkHealthConfig query bool false "是否使用健康检查配置" default(true)
// @Success 200 {array} DeviceHealth "返回设备组内所有设备的健康检查结果"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /HealthCheckByGroupId [get]
func HealthCheckByGroupId(groupId int, checkHealthConfig bool) (deviceHealths []DeviceHealth, err error) {
	// 根据组ID获取该组内的HCU设备列表
	dvIdList, groupName, err := getHCUListFromGroup(groupId)
	if err != nil {
		return nil, fmt.Errorf("获取组内设备列表失败: %v", err)
	}
	if len(dvIdList) == 0 {
		return nil, fmt.Errorf("设备组 '%s' (ID: %d) 中没有找到HCU设备", groupName, groupId)
	}
	glog.V(5).Infof("开始检查设备组 '%s' (ID: %d) 的健康状态，包含 %d 个设备",
		groupName, groupId, len(dvIdList))
	return healthCheckById(dvIdList, checkHealthConfig)
}

// GetDeviceModelInfos 获取设备型号信息
// @Summary 获取所有设备的型号信息
// @Description 返回系统中所有设备的型号信息列表
// @Produce json
// @Success 200 {array} DeviceModelInfo "返回设备型号信息列表"
// @Failure 500 {object} error "服务器内部错误"
// @Router /GetDeviceModelInfos [get]
func GetDeviceModelInfos() (devices []DeviceModelInfo) {
	for model, name := range modelName {
		cuCount := computeUnit[model]
		memorySize := memorySize[model]
		devices = append(devices, DeviceModelInfo{
			Model:      name,
			CUCount:    cuCount,
			MemorySize: memorySize,
		})
	}
	glog.V(5).Infof("GetDeviceModelInfos: %v", (devices))
	return devices
}

// Compatible 校验卡型号、驱动版本和 DTK 版本的兼容性。
// 版本不匹配返回 Status 为 warn 的结果；仅未知卡型号或无法解析版本时返回 error。
// @Summary 校验兼容性
// @Description 根据卡型号、驱动版本和 DTK 版本，返回通过或告警结果及推荐版本。
// @Produce json
// @Param cardModel query string true "卡型号，例如 K100, K100_AI"
// @Param driverVersion query string true "驱动版本，例如 5.16.18, 6.2.26"
// @Param dtkVersion query string true "DTK 版本，例如 23.10, 24.04"
// @Success 200 {object} CompatibilityResult "兼容性检查结果"
// @Failure 400 {object} error "请求参数错误"
// @Failure 500 {object} error "服务器内部错误"
// @Router /compatible [get]
func Compatible(cardModel, driverVersion, dtkVersion string) (CompatibilityResult, error) {
	return compatible(cardModel, driverVersion, dtkVersion)
}

// 获取DF带宽信息
func DFBandwidth(dvInd int, bandwidthType int) (dfBandwidthInfo DFBandwidthInfo, err error) {
	return dfBandwidth(dvInd, bandwidthType)
}

func UMCBandwidth(dvInd int, chanId int, delay int) (info UMCBandwidthInfo, err error) {
	return umcBandwidth(dvInd, chanId, delay)
}

// XHCLBandwidth 获取指定物理设备的 XHCL 带宽信息
// dvInd: 物理设备索引
// linkId: XHCL 链路 ID
// direction: 带宽方向
// delay: 采样延迟
func XHCLBandwidth(dvInd int, linkId int, direction int, delay int) (info XhclBandwidthInfo, err error) {
	return xhclBandwidth(dvInd, linkId, direction, delay)
}

// GetHyLinkLinkStatus 查询所有监控设备的 XHCL link 带宽，按 device 顺序返回切片。
// 返回：
//   - []DeviceLinkSum: 按 device index 升序（0..N-1）的结果切片
//   - error: 若任一设备发生查询错误，则返回一个汇总 error（同时切片仍包含每台设备的数据/错误信息）
func GetHyLinkStatus() ([]DeviceLinkSum, error) {
	deviceCount, err := rsmiNumMonitorDevices()
	if err != nil {
		return nil, fmt.Errorf("failed to get monitor device count: %w", err)
	}
	if deviceCount <= 0 {
		return nil, fmt.Errorf("no monitor devices found (count=%d)", deviceCount)
	}

	deviceList := make([]DeviceLinkSum, 0, deviceCount)
	var allDeviceErrors []string

	for i := 0; i < deviceCount; i++ {
		item, err := HyLinkStatusByHcuId(i)
		if err != nil {
			if item.Err != "" {
				allDeviceErrors = append(allDeviceErrors, item.Err)
			} else {
				allDeviceErrors = append(allDeviceErrors, fmt.Sprintf("device %d: %v", i, err))
			}
		}
		deviceList = append(deviceList, item)
	}

	if len(allDeviceErrors) > 0 {
		return deviceList, fmt.Errorf("some devices failed: %s", strings.Join(allDeviceErrors, " | "))
	}
	return deviceList, nil
}

// HyLinkStatusByHcuId 查询指定设备（索引 dvInd）上每个 link 的 recv/send 明细及汇总。
// 查询失败的方向保持零值，并在结果中聚合错误。
func HyLinkStatusByHcuId(dvInd int) (DeviceLinkSum, error) {
	var result DeviceLinkSum

	deviceCount, err := rsmiNumMonitorDevices()
	if err != nil {
		return result, fmt.Errorf("failed to get monitor device count: %w", err)
	}
	if deviceCount <= 0 {
		return result, fmt.Errorf("no monitor devices found (count=%d)", deviceCount)
	}
	if dvInd < 0 || dvInd >= deviceCount {
		return result, fmt.Errorf("device index out of range: %d (count=%d)", dvInd, deviceCount)
	}

	glog.V(5).Infof("HyLinkStatusByHcuId: query device index=%d (deviceCount=%d)", dvInd, deviceCount)

	result = DeviceLinkSum{
		DvInd: dvInd,
		Links: make([]LinkBandwidth, 0, MAX_XHCL_LINK_NUM),
	}
	if MAX_XHCL_LINK_NUM <= 0 {
		return result, fmt.Errorf("invalid MAX_XHCL_LINK_NUM: %d", MAX_XHCL_LINK_NUM)
	}

	roundToTwoDecimal := func(value float64) float64 {
		return math.Round(value*100) / 100
	}

	var perDeviceErrors []string
	var recvInfo *XhclBandwidthInfo
	var sendInfo *XhclBandwidthInfo

	if info, queryErr := xhclBandwidth(dvInd, linkAllID, 0, delay); queryErr != nil {
		perDeviceErrors = append(perDeviceErrors, fmt.Sprintf("direction=0: %v", queryErr))
		glog.V(4).Infof("device %d direction=0 xhclBandwidth failed: %v", dvInd, queryErr)
	} else {
		recvInfo = &info
	}

	if info, queryErr := xhclBandwidth(dvInd, linkAllID, 1, delay); queryErr != nil {
		perDeviceErrors = append(perDeviceErrors, fmt.Sprintf("direction=1: %v", queryErr))
		glog.V(4).Infof("device %d direction=1 xhclBandwidth failed: %v", dvInd, queryErr)
	} else {
		sendInfo = &info
	}

	for linkID := 0; linkID < MAX_XHCL_LINK_NUM; linkID++ {
		var recv, send float64
		if recvInfo != nil {
			recv = recvInfo.Bw[linkID]
			result.Recv += recv
		}
		if sendInfo != nil {
			send = sendInfo.Bw[linkID]
			result.Send += send
		}
		result.Links = append(result.Links, LinkBandwidth{
			LinkId: linkID,
			Recv:   roundToTwoDecimal(recv),
			Send:   roundToTwoDecimal(send),
		})
	}

	result.Recv = roundToTwoDecimal(result.Recv)
	result.Send = roundToTwoDecimal(result.Send)
	if len(perDeviceErrors) == 0 {
		return result, nil
	}

	result.Err = fmt.Sprintf("device %d: %s", dvInd, strings.Join(perDeviceErrors, "; "))
	return result, fmt.Errorf("some queries failed: %s", strings.Join(perDeviceErrors, " | "))
}

// GetHyUMCStatus 查询所有监控设备的 UMC 带宽
// 返回:
//   - []DeviceUmcSum: 按 device index 升序（0..N-1）的结果切片
//   - error: 若任一设备发生查询错误，则返回一个汇总 error（同时切片仍包含每台设备的数据/错误信息）
func GetHyUMCStatus() ([]DeviceUmcSum, error) {
	// 获取并校验设备数量
	deviceCount, err := rsmiNumMonitorDevices()
	if err != nil {
		return nil, fmt.Errorf("failed to get monitor device count: %w", err)
	}
	if deviceCount <= 0 {
		return nil, fmt.Errorf("no monitor devices found (count=%d)", deviceCount)
	}

	results := make([]DeviceUmcSum, 0, deviceCount)
	var deviceErrors []string
	roundToTwoDecimal := func(v float64) float64 { return math.Round(v*100) / 100 }

	// 遍历每个设备（按索引 0..deviceCount-1）
	for dvInd := 0; dvInd < deviceCount; dvInd++ {
		item := DeviceUmcSum{
			DvInd:     dvInd,
			Read:      0,
			Write:     0,
			ReadWrite: 0,
			Err:       "",
		}

		var perDeviceErrors []string

		// 获取该设备全部 channel 的 UMC 带宽信息
		info, e := umcBandwidth(dvInd, chanAllID, delay)
		if e != nil {
			perDeviceErrors = append(perDeviceErrors, fmt.Sprintf("umcBandwidth error: %v", e))
		} else {
			var readSum, writeSum, readWriteSum float64
			for i := 0; i < MAX_UMC_CHAN_NUM; i++ {
				readSum += info.ReadBW[i]
				writeSum += info.WriteBW[i]
				readWriteSum += info.ReadWriteBW[i]
			}
			item.Read = roundToTwoDecimal(readSum)
			item.Write = roundToTwoDecimal(writeSum)
			item.ReadWrite = roundToTwoDecimal(readWriteSum)
		}

		if len(perDeviceErrors) > 0 {
			msg := fmt.Sprintf("device %d: %s", dvInd, strings.Join(perDeviceErrors, "; "))
			item.Err = msg
			deviceErrors = append(deviceErrors, msg)
		}

		results = append(results, item)
	}

	if len(deviceErrors) > 0 {
		return results, fmt.Errorf("some devices failed: %s", strings.Join(deviceErrors, " | "))
	}
	return results, nil
}

func CompatibleTest(cardModel, driverVersion, dtkVersion string) error {
	_, err := compatible(cardModel, driverVersion, dtkVersion)
	return err
}

func ComputeProcessInfoGet() (processInfo []ProcessInfos, numItems int, err error) {
	return rsmiComputeProcessInfoGet()
}

func ProcessInfoByPid(pid uint32) (RsmiProcessInfoV2, error) {
	return getProcessInfoByPID(pid)
}

func DeviceGetCount() (deviceCount int, err error) {
	return nvmlDeviceGetCount()
}

func DeviceGetHandleByIndex(dvInd int) (device MIGDevice, err error) {
	return nvmlDeviceGetHandleByIndex(dvInd)
}

func DeviceGetMaxMigDeviceCountByIndex(dvInd int) (count int, err error) {
	return nvmlDeviceGetMaxMigDeviceCountByIndex(dvInd)
}

func DeviceGetAttributesByIndex(dvInd int, migId int) (attr NvmlDeviceAttributes, err error) {
	return nvmlDeviceGetAttributesByIndex(dvInd, migId)
}

func DeviceGetMigModeByIndex(dvInd int) (currentMode uint32, pendingMode uint32, err error) {
	return nvmlDeviceGetMigModeByIndex(dvInd)
}

func DeviceGetGpuInstanceProfileInfo(dvInd int, profileIdx uint32) (info NvmlGpuInstanceProfileInfo, err error) {
	return nvmlDeviceGetGpuInstanceProfileInfo(dvInd, profileIdx)
}

func DeviceGetGpuInstanceRemainingCapacity(dvInd int, profileId uint32) (count uint32, err error) {
	return nvmlDeviceGetGpuInstanceRemainingCapacity(dvInd, profileId)
}

func GpuInstanceGetComputeInstanceProfileInfo(dvInd int, giId uint32, profileId uint32, engProfileId uint32) (info NvmlComputeInstanceProfileInfo, err error) {
	return nvmlGpuInstanceGetComputeInstanceProfileInfo(dvInd, giId, profileId, engProfileId)
}

func DeviceGetGpuInstancesInfo(dvInd int, profileId uint32) ([]GpuInstanceInfo, error) {
	return nvmlDeviceGetGpuInstances(dvInd, profileId)
}

func ComputeInstanceProfileInfoList(dvInd int, profileId, ciProfileId, engProfileId uint32) ([]NvmlComputeInstanceProfileInfo, error) {
	return allComputeInstanceProfileInfo(dvInd, profileId, ciProfileId, engProfileId)
}

func GpuInstanceGetComputeInstanceRemainingCapacity(dvInd int, giId uint32, profileId uint32) (uint32, error) {
	return nvmlGpuInstanceGetComputeInstanceRemainingCapacity(dvInd, giId, profileId)
}

func GpuInstancesComputeInstanceRemainingCapacityList(dvInd int, profileId uint32) ([]ComputeInstanceRemainInfo, error) {
	return nvmlAllGpuInstancesComputeInstanceRemainingCapacity(dvInd, profileId)
}

func DeviceGetGpuInstanceId(dvInd int, migId int) (gpuInstanceId uint32, err error) {
	return nvmlDeviceGetGpuInstanceId(dvInd, migId)
}

func DeviceGetComputeInstanceId(dvInd int, migId int) (computeInstanceId uint32, err error) {
	return nvmlDeviceGetComputeInstanceId(dvInd, migId)
}

func AvailableMigDeviceIds(dvInd int) (int, []int, error) {
	return availableMigDeviceIds(dvInd)
}

func DeviceGetMigDeviceHandleByIndex(dvInd int, migId int) (migDevice MIGDevice, err error) {
	return nvmlDeviceGetMigDeviceHandleByIndex(dvInd, migId)
}

func DeviceGetPciInfo(device MIGDevice) (pciInfo NvmlPciInfo, err error) {
	return nvmlDeviceGetPciInfo(device)
}

// MigInfos 获取当前主机上所有GPU的所有MIG分区信息。
// MigInfos 组装所有 MigInfo 并返回
func MigInfos() ([]MigInfo, error) {
	// 1. 拿到所有的底层配置
	configs, err := migConfigs()
	if err != nil {
		return nil, fmt.Errorf("migConfigs failed: %w", err)
	}

	// 2. 计算 seCount
	seCount, err := getSECount()
	if err != nil {
		return nil, fmt.Errorf("getSECount failed: %w", err)
	}

	// 3. 遍历 configs，映射到 MigInfo
	infos := make([]MigInfo, 0, len(configs))
	for _, cfg := range configs {
		name := formatMIGName(
			seCount,
			cfg.Gi.GpuSliceCount,
			cfg.Ci.GpuSliceCount,
			uint64(cfg.Gi.MemorySizeMB),
		)
		infos = append(infos, MigInfo{
			DvInd:             cfg.DvInd,
			Name:              name,
			UUID:              "MIG-" + cfg.Ci.MigUUID,
			ComputeUnit:       uint32(cfg.Ci.CuCount),
			MemoryTotal:       uint64(cfg.Gi.MemorySizeMB),
			GpuInstanceId:     uint32(cfg.Ci.GiId),
			ComputeInstanceId: uint32(cfg.Ci.Id),
			PciBusNumber:      cfg.Ci.Pci,
			GiProfileId:       cfg.Gi.ProfileId,
			CiProfileId:       cfg.Ci.ProfileId,
		})
	}

	return infos, nil
}

// 遍历所有可用的GPU设备，依次收集每个GPU上所有MIG设备的属性、GPU实例ID和计算实例ID。
// 返回值:
//   - []MigInfo：每个GPU的MIG分区信息切片
//   - error：获取过程中遇到的错误
//func MigInfos() ([]MigInfo, error) {
//	var results []MigInfo
//	gpuCount, err := DeviceGetCount()
//	if err != nil {
//		return nil, fmt.Errorf("获取GPU数量失败: %v", err)
//	}
//
//	globalMigId := 0 // 全局递增MIG id
//
//	for gpuIdx := 0; gpuIdx < gpuCount; gpuIdx++ {
//		device, _ := nvmlDeviceGetHandleByIndex(gpuIdx)
//		pciInfo, _ := nvmlDeviceGetPciInfo(device)
//		glog.V(5).Infof("HCU %v ,nvmlDeviceGetPciInfo: %v", gpuIdx, (pciInfo))
//		foundAny := false
//		maxMigPerGpu, _ := DeviceGetMaxMigDeviceCountByIndex(gpuIdx)
//		glog.V(5).Infof("MaxMigDeviceCount for GPU%d: %v", gpuIdx, maxMigPerGpu)
//		for local := 0; local < maxMigPerGpu; local++ {
//			migDev, err := DeviceGetMigDeviceHandleByIndex(gpuIdx, globalMigId)
//			glog.V(5).Infof("DeviceGetMigDeviceHandleByIndex:%v", migDev)
//			if err != nil || migDev == nil {
//				break // 当前GPU的MIG设备查完
//			}
//
//			attr, err := DeviceGetAttributesByIndex(gpuIdx, globalMigId)
//			if err != nil {
//				glog.Warningf("获取MIG属性失败: gpu=%d mig=%d err=%v", gpuIdx, globalMigId, err)
//				globalMigId++
//				continue
//			}
//			gi, err := DeviceGetGpuInstanceId(gpuIdx, globalMigId)
//			if err != nil {
//				glog.Warningf("获取GI失败: gpu=%d mig=%d err=%v", gpuIdx, globalMigId, err)
//				globalMigId++
//				continue
//			}
//			ci, err := DeviceGetComputeInstanceId(gpuIdx, globalMigId)
//			if err != nil {
//				glog.Warningf("获取CI失败: gpu=%d mig=%d err=%v", gpuIdx, globalMigId, err)
//				globalMigId++
//				continue
//			}
//
//			results = append(results, MigInfo{
//				DvInd:             gpuIdx,
//				MigId:             globalMigId,
//				Name:              attr.Name,
//				UUID:              attr.UUID,
//				ComputeUnit:       attr.CUCount,
//				MemoryTotal:       attr.MemorySizeMB,
//				GpuInstanceId:     gi,
//				ComputeInstanceId: ci,
//				PciBusNumber:      pciInfo.BusID,
//			})
//
//			globalMigId++
//			foundAny = true
//		}
//		if !foundAny {
//			glog.Infof("物理GPU=%d 没有有效MIG设备", gpuIdx)
//		}
//	}
//
//	return results, nil
//}

// MigInfoByDvInd 根据设备索引（dvInd）获取该GPU的所有MIG分区信息。
// 参数:
//   - dvInd int：目标GPU的设备索引号
//
// 返回值:
//   - []MigInfo：该GPU的所有MIG分区信息切片
//   - error：未找到或获取失败时返回错误
func MigInfoByDvInd(dvInd int) ([]MigInfo, error) {
	// 调用MigInfos获取所有GPU的MIG信息
	allInfos, err := MigInfos()
	if err != nil {
		return nil, err
	}

	// 过滤出目标GPU的MIG分区信息
	var filtered []MigInfo
	for _, info := range allInfos {
		if info.DvInd == dvInd {
			filtered = append(filtered, info)
		}
	}
	if len(filtered) == 0 {
		return nil, fmt.Errorf("未找到索引为 %d 的GPU的MIG分区信息", dvInd)
	}
	return filtered, nil
}

// MigInfoByUUID 根据设备UUID获取该GPU的MIG分区信息。
// 参数:
//   - uuid string：目标GPU的UUID
//
// 返回值:
//   - MigInfo：目标GPU的MIG分区信息（如未找到则返回error）
func MigInfoByUUID(uuid string) (migInfo MigInfo, err error) {
	// 调用MigInfos获取所有GPU的MIG信息
	allInfos, err := MigInfos()
	if err != nil {
		return
	}

	// 查找目标UUID的MIG分区信息
	for _, info := range allInfos {
		if info.UUID == uuid {
			return info, nil
		}
	}

	return migInfo, fmt.Errorf("未找到UUID为 %s 的GPU的MIG分区信息", uuid)
}

func SystemMigMode() (currentMode, pendingMode int, err error) {
	return nvmlGetSystemMigMode()
}

func MIGConfigs() ([]MigConfig, error) {
	return migConfigs()
}

func MIGSECount() (int, error) {
	return getSECount()
}

func MIGName(seCount int, giSliceCount int, ciSliceCount int, memorySizeMB uint64) string {
	return formatMIGName(seCount, giSliceCount, ciSliceCount, memorySizeMB)
}

func DevHCUReset(dvInd int) (err error) {
	return rsmiDevGpuReset(dvInd)
}

func CreateGroup(groupName string) (groupId int, err error) {
	return createGroup(groupName)
}

func CreateDefaultGroup(groupName string) (groupId int, err error) {
	return createDefaultGroup(groupName)
}

func AddToGroup(groupId int, dvInd int) (err error) {
	return addToGroup(groupId, dvInd)
}

func AddEntityToGroup(groupId int, entityList []GroupEntityPair) error {
	return addEntityToGroup(groupId, entityList)
}

func RemoveEntityFromGroup(groupId int, entityList []GroupEntityPair) error {
	return removeEntityFromGroup(groupId, entityList)
}

func DestroyGroup(groupId int) error {
	return destroyGroup(groupId)
}

func GetGroupInfo(groupId int) (groupInfo GroupInfo, err error) {
	return getGroupInfo(groupId)
}

func GetHCUListFromGroup(groupId int) (hcuIds []int, groupName string, err error) {
	return getHCUListFromGroup(groupId)
}

func ListAllGroups() (groups []GroupInfo, err error) {
	return listAllGroups()
}

func CreateFieldGroup(fieldGroupName string, fieldIds []int) (fieldGroupId int, err error) {
	return createFieldGroup(fieldGroupName, fieldIds)
}

func DestroyFieldGroup(fieldGroupId int) error {
	return destroyFieldGroup(fieldGroupId)
}

func GetFieldGroupInfo(fieldGroupId int) (fieldGroupInfo FieldGroupInfo, err error) {
	return getFieldGroupInfo(fieldGroupId)
}

func ListAllFieldGroups() (fieldGroups []FieldGroupInfo, err error) {
	return listAllFieldGroups()
}

func SetPolicy(policyInfo Policy, dvInd int) error {
	return setPolicy(policyInfo, dvInd)
}

func ClearPolicy(dvIdList []int) error {
	return clearPolicy(dvIdList)
}

func GetPolicy(dvIdList []int) (policyList []Policy, err error) {
	return getPolicy(dvIdList)
}

func JudgePolicyConditions(dvIdList []int) (dvInd int, err error) {
	return judgePolicyConditions(dvIdList)
}

func TakePolicyAction(dvInd int) error {
	return takePolicyAction(dvInd)
}

func GetFieldMetaById(fieldId int) FieldMeta {
	return getFieldMetaById(fieldId)
}

func ListFieldMeta() []FieldMeta {
	return listFieldMeta()
}

func WatchFields(dvInd int, fieldIds []int) error {
	return watchFields(dvInd, fieldIds)
}

func WatchFieldGroup(dvInd int, fieldGroupId int) error {
	return watchFieldGroup(dvInd, fieldGroupId)
}

func WatchFieldsWithEntity(entityGroup Field_Entity_Group, entityId int, fieldIds []int) error {
	return watchFieldsWithEntity(entityGroup, entityId, fieldIds)
}

func WatchFieldGroupWithEntity(entityGroup Field_Entity_Group, entityId int, fieldGroupId int) error {
	return watchFieldGroupWithEntity(entityGroup, entityId, fieldGroupId)
}

func WatchFieldsWithGroup(fieldGroupId int, groupId int) error {
	return watchFieldsWithGroup(fieldGroupId, groupId)
}

func WatchFieldsWithEntityGroup(
	fieldIds []int, groupId int, updateFreq time.Duration, maxKeepAge time.Duration, maxKeepSamples int32,
) error {
	return watchFieldsWithEntityGroup(fieldIds, groupId, updateFreq, maxKeepAge, maxKeepSamples)
}

func WatchFieldsWithGroupEx(fieldGroupId int, groupId int, updateFreq time.Duration, maxKeepAge time.Duration, maxKeepSamples int32) error {
	return watchFieldsWithGroupEx(fieldGroupId, groupId, updateFreq, maxKeepAge, maxKeepSamples)
}

func UnWatchFields(dvInd int) {
	unWatchFields(dvInd)
}

func UnWatchFieldsWithEntity(entityGroup Field_Entity_Group, entityId int) {
	unWatchFieldsWithEntity(entityGroup, entityId)
}

func UnWatchFieldsWithGroup(groupId int) error {
	return unWatchFieldsWithGroup(groupId)
}

func GetLatestValuesForFields(dvInd int, fields []int) ([]FieldValue_v1, error) {
	return getLatestValuesForFields(dvInd, fields)
}

func EntityGetLatestValues(entityGroup Field_Entity_Group, entityId int, fields []int) ([]FieldValue_v1, error) {
	return entityGetLatestValues(entityGroup, entityId, fields)
}

func EntitiesGetLatestValues(entities []GroupEntityPair, fields []int) ([]FieldValue_v2, error) {
	return entitiesGetLatestValues(entities, fields)
}

func GetSupportedMetricGroups(dvInd int) ([]MetricGroup, error) {
	return getSupportedMetricGroups(dvInd)
}

// DevCuNum 返回设备的 CU（Compute Unit）数量
func DevCuNum(dvInd int) (int, error) {
	return rsmiDevCuNumGet(dvInd)
}

// DevNodeId 返回设备的 KFD node ID，用于拓扑关联
func DevNodeId(dvInd int) (uint32, error) {
	return rsmiDevNodeIdGet(dvInd)
}

// DevMemOverdriveLevelGet 返回显存 overdrive 百分比（0-20）
// 对应 hy-smi --showmemoverdrive，监控显存超频状态
func DevMemOverdriveLevelGet(dvInd int) (int, error) {
	return rsmiDevMemOverdriveLevelGet(dvInd)
}

// IsP2PAccessible 查询两设备之间是否支持 P2P 访问
// 对应 hy-smi --showtopoaccess；srcDvInd → dstDvInd 方向
func IsP2PAccessible(srcDvInd, dstDvInd int) (bool, error) {
	return rsmiIsP2PAccessible(srcDvInd, dstDvInd)
}

// HyVersion 返回海光扩展版本（hy_major/hy_minor）
// 对应 hy-smi -libversion 中的 hy 段；与 Version() 的 ROCm 段互补
func HyVersion() (HyVersionInfo, error) {
	return rsmiHyVersionGet()
}

// DevGfxClockMax 返回 GFX（shader）时钟的最大额定频率（MHz）
// 对应 DCGM_FI_DEV_SM_CLOCK_MAX
func DevGfxClockMax(dvInd int) (uint64, error) {
	freqs, _, err := GetClocksByType(dvInd, RSMI_CLK_TYPE_SYS)
	if err != nil {
		return 0, err
	}
	var max uint64
	for _, f := range freqs {
		if f > max {
			max = f
		}
	}
	return max, nil
}

// DevMemClockMax 返回显存时钟的最大额定频率（MHz）
// 对应 DCGM_FI_DEV_MEM_CLOCK_MAX
func DevMemClockMax(dvInd int) (uint64, error) {
	freqs, _, err := GetClocksByType(dvInd, RSMI_CLK_TYPE_MEM)
	if err != nil {
		return 0, err
	}
	var max uint64
	for _, f := range freqs {
		if f > max {
			max = f
		}
	}
	return max, nil
}

// DevCuUtil 查询设备 CU 在采样窗口内的 wave 占用周期占比（全 CU 平均）。
//
// 语义：duration 窗口内，每个 CU 至少分配了 1 个 wave 的时钟周期占比，再对所有 CU 取平均。
// duration 单位由驱动定义（实测传 1000 即约 1s 窗口）；传 0 使用驱动默认窗口。
// 返回值 percent 范围通常为 0~1。
func DevCuUtil(dvInd int, duration int) (float32, error) {
	return rsmiDevCuUtilGet(dvInd, duration)
}

// DevWaveUtil 查询设备 CU 在采样窗口内的 wave 驻留数量占比（全 CU 平均）。
//
// 语义：duration 窗口内各 CU 上驻留 wave 数量占最大可驻留 wave 数量的比例，再对所有 CU 取平均。
// 与 DevCuUtil 互补：DevCuUtil 看"有没有 wave"，DevWaveUtil 看"wave 填了多满"。
// 返回值 percent 范围通常为 0~1。
func DevWaveUtil(dvInd int, duration int) (float32, error) {
	return rsmiDevWaveUtilGet(dvInd, duration)
}

// ── sysfs 只读功能（第2批）──────────────────────────────────────────────────────

// HswFwVersion 读取 HSW（海光交换机？）固件版本号。
//
// sysfs 路径：/sys/class/hsd/hsd{dvInd}/build_info
// 对应 hy-smi --showfwinfo 中的 HSW 条目。
func HswFwVersion(dvInd int) (string, error) {
	path := fmt.Sprintf("/sys/class/hsd/hsd%d/build_info", dvInd)
	return readSysfsFile(path)
}

// GfxBoostInfo 读取 GFX boost 等级信息。
//
// sysfs 路径：/sys/class/drm/card{dvInd}/device/pp_gfx_boost
// 对应 hy-smi --showboost。内容为驱动定义的键值对文本。
func GfxBoostInfo(dvInd int) (string, error) {
	cardIdx := dvIndexToCardIndex(dvInd)
	path := fmt.Sprintf("/sys/class/drm/card%d/device/pp_gfx_boost", cardIdx)
	return readSysfsFile(path)
}

// FdrHistory 读取 GFX 频率限制历史（frequency dynamic range history）。
//
// sysfs 路径：/sys/class/drm/card{dvInd}/device/pp_gfx_limit_info
// 对应 hy-smi --showfdrhistory。
func FdrHistory(dvInd int) (string, error) {
	cardIdx := dvIndexToCardIndex(dvInd)
	path := fmt.Sprintf("/sys/class/drm/card%d/device/pp_gfx_limit_info", cardIdx)
	return readSysfsFile(path)
}

// MemBankEccInfo 读取显存 bank 级 ECC 统计信息。
//
// sysfs 路径：/sys/class/drm/card{dvInd}/device/mem_banks/*/properties
// 对应 hy-smi --showmemeccinfo。每个 bank 一个子目录，properties 为键值对文本。
// 返回：map[bankName]propertiesText
func MemBankEccInfo(dvInd int) (map[string]string, error) {
	cardIdx := dvIndexToCardIndex(dvInd)
	baseDir := fmt.Sprintf("/sys/class/drm/card%d/device/mem_banks", cardIdx)

	entries, err := listSysfsDir(baseDir)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string, len(entries))
	for _, entry := range entries {
		propPath := fmt.Sprintf("%s/properties", entry)
		content, err := readSysfsFile(propPath)
		if err != nil {
			// 跳过不可读的 bank（可能是权限受限或空目录）
			glog.V(5).Infof("MemBankEccInfo: skip %s: %v", propPath, err)
			continue
		}
		result[fmt.Sprintf("%s", entry[len(baseDir)+1:])] = content
	}
	return result, nil
}

// SetGfxBoostLevel 设置设备 GFX boost 等级，具体取值范围由驱动校验。
func SetGfxBoostLevel(dvInd, level int) error {
	cardIdx := dvIndexToCardIndex(dvInd)
	path := fmt.Sprintf("/sys/class/drm/card%d/device/pp_gfx_boost", cardIdx)
	return writeSysfsFile(path, strconv.Itoa(level))
}

// SetLowPowerMode 设置设备运行时电源管理模式。enabled=true 写入 auto，false 写入 on。
func SetLowPowerMode(dvInd int, enabled bool) error {
	cardIdx := dvIndexToCardIndex(dvInd)
	path := fmt.Sprintf("/sys/class/drm/card%d/device/power/control", cardIdx)
	mode := "on"
	if enabled {
		mode = "auto"
	}
	return writeSysfsFile(path, mode)
}

// SetLowPowerDelay 设置设备运行时自动挂起延时，单位为毫秒。
func SetLowPowerDelay(dvInd, delayMs int) error {
	cardIdx := dvIndexToCardIndex(dvInd)
	path := fmt.Sprintf("/sys/class/drm/card%d/device/power/autosuspend_delay_ms", cardIdx)
	return writeSysfsFile(path, strconv.Itoa(delayMs))
}

// SetRasFeature 启用或禁用指定 RAS block。
func SetRasFeature(dvInd int, block string, enabled bool) error {
	if err := validateDebugfsToken(block); err != nil {
		return err
	}
	action := "disable"
	if enabled {
		action = "enable"
	}
	cardIdx := dvIndexToCardIndex(dvInd)
	path := fmt.Sprintf("/sys/kernel/debug/dri/%d/ras_ctrl", cardIdx)
	return writeDebugfs(path, action+" "+block)
}

// InjectRasError 向指定设备写入驱动定义的 RAS 注入参数。
// parameters 顺序与 ras_ctrl 的 inject 命令一致，参数本身不得包含空白符。
func InjectRasError(dvInd int, parameters []string) error {
	if len(parameters) == 0 {
		return fmt.Errorf("RAS inject parameters must not be empty")
	}
	for _, parameter := range parameters {
		if err := validateDebugfsToken(parameter); err != nil {
			return err
		}
	}
	cardIdx := dvIndexToCardIndex(dvInd)
	path := fmt.Sprintf("/sys/kernel/debug/dri/%d/ras_ctrl", cardIdx)
	return writeDebugfs(path, "inject "+strings.Join(parameters, " "))
}

func validateDebugfsToken(value string) error {
	if value == "" || strings.ContainsAny(value, " \t\r\n\x00") {
		return fmt.Errorf("invalid debugfs command token %q", value)
	}
	return nil
}

// HCUHealthCheck 执行 hy-smi --healthcheck，返回所有 HCU 卡的健康状态
func HCUHealthCheck() ([]HCUHealthStatus, error) {
	return hcuHealthCheck()
}
