/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package dcgm

/*
#cgo CFLAGS: -Wall -I./include
#cgo LDFLAGS: -L/opt/hyhal/lib -Wl,-rpath,/opt/hyhal/lib -lrocm_smi64 -Wl,--unresolved-symbols=ignore-in-object-files
#include <stdint.h>
#include <kfd_ioctl.h>
#include <rocm_smi64Config.h>
#include <rocm_smi.h>
*/
import "C"
import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"go.etcd.io/bbolt"
	"strconv"
)

var (
	POLICY_BUCKET = []byte("policy")
)

type PolicyAction int

const (
	ACTION_NONE PolicyAction = iota
	ACTION_HCU_RESET
	ACTION_COUNT
)

func (action PolicyAction) String() string {
	switch action {
	case ACTION_NONE:
		return "None"
	case ACTION_HCU_RESET:
		return "Reset HCU"
	}
	return "Unknown"
}

type PolicyValidation int

const (
	VALIDATION_NONE PolicyValidation = iota
	VALIDATION_SHORT
	VALIDATION_MEDIUM
	VALIDATION_LONG
	VALIDATION_COUNT
)

func (validation PolicyValidation) String() string {
	switch validation {
	case VALIDATION_NONE:
		return "None"
	case VALIDATION_SHORT:
		return "System Validation (Short)"
	case VALIDATION_MEDIUM:
		return "System Validation (Medium)"
	case VALIDATION_LONG:
		return "System Validation (Long)"
	}
	return "Unknown"
}

type PolicyConditions struct {
	MaxPagesEnable   bool    `json:"max_pages_enable"`
	MaxPages         int     `json:"max_pages"`
	MaxTempEnable    bool    `json:"max_temp_enable"`
	MaxTemp          float64 `json:"max_temp"`
	MaxPowerEnable   bool    `json:"max_power_enable"`
	MaxPower         int     `json:"max_power"`
	EccErrorsEnable  bool    `json:"ecc_errors_enable"`
	PcieErrorsEnable bool    `json:"pcie_errors_enable"`
}

type Policy struct {
	DvInd           int              `json:"dv_ind"`
	ActionIndex     PolicyAction     `json:"action_index"`
	ValidationIndex PolicyValidation `json:"validation_index"`
	Conditions      PolicyConditions `json:"conditions"`
}

func setPolicy(policyInfo Policy, hcuIndex int) error {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Update(func(tx *bbolt.Tx) error {
		policyBucket, err := tx.CreateBucketIfNotExists(POLICY_BUCKET)
		if err != nil {
			return err
		}
		info, err := json.Marshal(policyInfo)
		if err != nil {
			return err
		}
		err = policyBucket.Put(itob(hcuIndex), info)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func clearPolicy(hcuList []int) error {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	return db.Update(func(tx *bbolt.Tx) error {
		policyBucket := tx.Bucket(POLICY_BUCKET)
		if policyBucket == nil {
			return nil
		}
		for _, hcuIndex := range hcuList {
			if err = policyBucket.Delete(itob(hcuIndex)); err != nil {
				return err
			}
		}
		return nil
	})
}

func getPolicy(hcuList []int) (policyList []Policy, err error) {
	db, err := OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	policyList = make([]Policy, 0, len(hcuList))

	err = db.View(func(tx *bbolt.Tx) error {
		policyBucket := tx.Bucket(POLICY_BUCKET)

		for _, hcuIndex := range hcuList {
			policyInfo := Policy{
				DvInd:           hcuIndex,
				ActionIndex:     ACTION_NONE,
				ValidationIndex: VALIDATION_NONE,
				Conditions:      PolicyConditions{},
			}
			if policyBucket != nil {
				info := policyBucket.Get(itob(hcuIndex))
				if info != nil {
					if unmarshalErr := json.Unmarshal(info, &policyInfo); unmarshalErr != nil {
						return unmarshalErr
					}
				}
			}
			policyList = append(policyList, policyInfo)
		}
		return nil
	})
	return policyList, err
}

func judgePolicyConditions(hcuList []int) (errorDvInd int, err error) {
	db, err := OpenDB()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	err = db.View(func(tx *bbolt.Tx) error {
		policyBucket := tx.Bucket(POLICY_BUCKET)
		if policyBucket == nil {
			return nil
		}

		var policyInfo Policy
		for _, hcuIndex := range hcuList {
			info := policyBucket.Get(itob(hcuIndex))
			if info != nil {
				err = json.Unmarshal(info, &policyInfo)
				if err != nil {
					return err
				}
			} else {
				continue
			}
			conditions := policyInfo.Conditions
			if conditions.MaxPagesEnable {
				numPages, _, err := rsmiDevMemoryReservedPagesGet(hcuIndex)
				if err != nil {
					return err
				}
				if numPages >= conditions.MaxPages {
					errorDvInd = hcuIndex
					return fmt.Errorf("HCU %d: The maximum number of retired pages has violated policy manager values\nPage retirement count: %d", hcuIndex, numPages)
				}
			}
			if conditions.MaxTempEnable {
				devTemp, err := rsmiDevTempMetricGet(hcuIndex, 0, RSMI_TEMP_CURRENT)
				if err != nil {
					return err
				}
				temp, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(devTemp)/1000.0), 64)
				if temp >= conditions.MaxTemp {
					errorDvInd = hcuIndex
					return fmt.Errorf("HCU %d: The maximum thermal limit has violated policy manager values\nTemperature: %.0f", hcuIndex, temp)
				}
			}
			if conditions.MaxPowerEnable {
				powerAve, err := rsmiDevPowerAveGet(hcuIndex, 0)
				if err != nil {
					return err
				}
				power := powerAve / 1000000
				if int(power) >= conditions.MaxPower {
					errorDvInd = hcuIndex
					return fmt.Errorf("HCU %d: The maximum power limit has violated policy manager values.\nPower: %d", hcuIndex, power)
				}
			}
			if conditions.EccErrorsEnable {
				blockInfoList, err := EccBlocksInfo(hcuIndex)
				if err != nil {
					return err
				}
				for _, blockInfo := range blockInfoList {
					if blockInfo.CE > 0 {
						fmt.Printf("Warning: HCU %d ECC CE Errors: %d\n", hcuIndex, blockInfo.CE)
					}
					if blockInfo.UE > 0 {
						errorDvInd = hcuIndex
						return fmt.Errorf("HCU %d: ECC error has violated policy manager values.\nECC error count: %d", hcuIndex, blockInfo.UE)
					}
				}
			}
			if conditions.PcieErrorsEnable {
				count, err := rsmiDevPciReplayCounterGet(hcuIndex)
				if err != nil {
					return err
				}
				if count > 8 {
					errorDvInd = hcuIndex
					return fmt.Errorf("HCU %d: A PCIe replay event has violated policy manager values.\nPCIe replay count: %d", hcuIndex, count)
				}

			}
		}
		return nil
	})
	return errorDvInd, err
}

func takePolicyAction(hcuIndex int) error {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.View(func(tx *bbolt.Tx) error {
		poilcyBucket := tx.Bucket(POLICY_BUCKET)
		if poilcyBucket == nil {
			return fmt.Errorf("action not found")
		}
		info := poilcyBucket.Get(itob(hcuIndex))
		var policyInfo Policy
		if info == nil {
			return fmt.Errorf("action not found")
		} else {
			if err := json.Unmarshal(info, &policyInfo); err != nil {
				return err
			}
		}
		action := policyInfo.ActionIndex
		switch action {
		case ACTION_NONE:
			fmt.Println("No action required.")
		case ACTION_HCU_RESET:
			err := rsmiDevGpuReset(hcuIndex)
			if err != nil {
				return err
			}
			fmt.Println("HCU has been reset.")
		}
		return nil
	})
	return err
}

// rsmiDevPerfLevelSet 设置设备PowerPlay性能级别
func rsmiDevPerfLevelSet(dvInd int, devPerfLevel DevPerfLevel) (err error) {
	glog.V(5).Infof("dev_perf_level_set:", devPerfLevel)
	ret := C.rsmi_dev_perf_level_set(C.uint32_t(dvInd), C.rsmi_dev_perf_level_t(devPerfLevel))
	glog.V(5).Infof("dev_perf_level_set ret:%v,retstr:%v", ret, errorString(ret))
	if err = errorString(ret); err != nil {
		return fmt.Errorf("dev_perf_level_set:%s", err)
	}
	return
}

// rsmiDevClkRangeSet 设置设备时钟范围信息
func rsmiDevClkRangeSet(dvInd int, minClkValue, maxClkValue int64, clkType RSMIClkType) (err error) {
	ret := C.rsmi_dev_clk_range_set(C.uint32_t(dvInd), C.uint64_t(minClkValue), C.uint64_t(maxClkValue), C.rsmi_clk_type_t(clkType))
	glog.V(5).Infof("rsmi_dev_clk_range_set ret:%v, retstr:%v", ret, errorString(ret))
	if err = errorString(ret); err != nil {
		glog.Errorf("Error rsmi_dev_clk_range_set:%s", err)
		return fmt.Errorf("Error rsmi_dev_clk_range_set:%s", err)
	}
	return
}

// rsmiDevOdVoltInfoSet 设置设备电压曲线点
func rsmiDevOdVoltInfoSet(dvInd, vPoint, clkValue, voltValue int) (err error) {
	ret := C.rsmi_dev_od_volt_info_set(C.uint32_t(dvInd), C.uint32_t(vPoint), C.uint64_t(clkValue), C.uint64_t(voltValue))
	glog.V(5).Infof("rsmi_dev_od_volt_info_set ret:%v", ret)
	if err = errorString(ret); err != nil {
		return fmt.Errorf("Error rsmi_dev_od_volt_info_set:%s", err)
	}
	return
}

// rsmiDevOverdriveLevelSet 设置设备超速百分比
func rsmiDevOverdriveLevelSet(dvInd, od int) (err error) {
	ret := C.rsmi_dev_overdrive_level_set(C.uint32_t(dvInd), C.uint32_t(od))
	glog.V(5).Infof("rsmi_dev_overdrive_level_set ret:%v, retStr:%v", ret, errorString(ret))
	if err = errorString(ret); err != nil {
		return fmt.Errorf("Error rsmi_dev_overdrive_level_set:%s", err)
	}
	return
}

// rsmiDevGpuClkFreqSet 设置可用于指定时钟的频率集
func rsmiDevGpuClkFreqSet(dvInd int, clkType RSMIClkType, freqBitmask int64) (err error) {
	ret := C.rsmi_dev_gpu_clk_freq_set(C.uint32_t(dvInd), C.rsmi_clk_type_t(clkType), C.uint64_t(freqBitmask))
	glog.V(5).Infof("rsmi_dev_gpu_clk_freq_set: ret: %v, retStr:%v", ret, errorString(ret))
	if err = errorString(ret); err != nil {
		return fmt.Errorf("Error rsmi_dev_gpu_clk_freq_set:%s", err)
	}
	return nil
}

// rsmiDevCounterGroupSupported 判断设备是否支持特定事件组
func rsmiDevCounterGroupSupported(dvInd int, group RSMIEventGroup) (err error) {
	ret := C.rsmi_dev_counter_group_supported(C.uint32_t(dvInd), C.rsmi_event_group_t(group))
	if err = errorString(ret); err != nil {
		return fmt.Errorf("Error rsmi_dev_counter_group_supported:%s", err)
	}
	return
}

// rsmiDevCounterCreate 创建性能计数器对象
func rsmiDevCounterCreate(dvInd int, eventType RSMIEventType) (eventHandle EventHandle, err error) {
	var ceventHandle C.rsmi_event_handle_t
	ret := C.rsmi_dev_counter_create(C.uint32_t(dvInd), C.rsmi_event_type_t(eventType), &ceventHandle)
	if err = errorString(ret); err != nil {
		return eventHandle, fmt.Errorf("Error rsmi_dev_counter_create:%s", err)
	}
	eventHandle = EventHandle(ceventHandle)
	return
}

// rsmiDevCounterDestroy 释放性能计数器对象
func rsmiDevCounterDestroy(handle EventHandle) (err error) {
	var chandle C.rsmi_event_handle_t
	ret := C.rsmi_dev_counter_destroy(C.rsmi_event_handle_t(chandle))
	if err = errorString(ret); err != nil {
		return fmt.Errorf("Error rsmi_dev_counter_destroy:%s", err)
	}
	return
}

// rsmiCounterControl 发布性能计数器控制命令
func rsmiCounterControl(evtHandle EventHandle, cmd RSMICounterCommand) (err error) {
	ret := C.rsmi_counter_control(C.rsmi_event_handle_t(evtHandle), C.rsmi_counter_command_t(cmd), nil)

	if err := errorString(ret); err != nil {
		return fmt.Errorf("Error in rsmi_counter_control: %s", err)
	}
	return
}

// rsmiCounterRead 读取性能计数器的当前值
func rsmiCounterRead(handle EventHandle) (counterValue RSMICounterValue, err error) {
	var ccounterValue C.rsmi_counter_value_t
	ret := C.rsmi_counter_read(C.rsmi_event_handle_t(handle), &ccounterValue)
	if err = errorString(ret); err != nil {
		return counterValue, fmt.Errorf("Error rsmiCounterRead:%s", err)
	}
	counterValue = RSMICounterValue{
		Value:       uint64(ccounterValue.value),
		TimeEnabled: uint64(ccounterValue.time_enabled),
		TimeRunning: uint64(ccounterValue.time_running),
	}
	return
}

func rsmiCounterAvailableCountersGet(dvInd int, group RSMIEventGroup) (availAble int, err error) {
	var cavailAble C.uint32_t
	ret := C.rsmi_counter_available_counters_get(C.uint32_t(dvInd), C.rsmi_event_group_t(group), &cavailAble)
	if err = errorString(ret); err != nil {
		return availAble, fmt.Errorf("Error rsmiCounterAvailableCountersGet:%s", err)
	}
	availAble = int(cavailAble)
	return
}

// rsmiDevFanReset 将风扇复位为自动驱动控制
func rsmiDevFanReset(dvInd, sensorInd int) (err error) {
	ret := C.rsmi_dev_fan_reset(C.uint32_t(dvInd), C.uint32_t(sensorInd))
	glog.V(5).Infof("rsmi_dev_fan_reset_ret:", ret)
	if err = errorString(ret); err != nil {
		return fmt.Errorf("Error rsmi_dev_fan_reset: %s", err)
	}
	return nil
}

// rsmiDevPowerProfileSet 设置设备功率配置文件
func rsmiDevPowerProfileSet(dvInd int, reserved int, profile PowerProfilePresetMasks) (err error) {
	ret := C.rsmi_dev_power_profile_set(C.uint32_t(dvInd), C.uint32_t(reserved), C.rsmi_power_profile_preset_masks_t(profile))
	glog.V(5).Infof("rsmi_dev_power_profile_set ret:%v, retstr:%v", ret, errorString(ret))
	if err = errorString(ret); err != nil {
		glog.Errorf("Error rsmi_dev_power_profile_set:%v", err)
		return fmt.Errorf("Error rsmi_dev_power_profile_set:%s", err)
	}
	return
}

// rsmiDevHSLErrorReset 重置设备的HSL错误状态
func rsmiDevHSLErrorReset(dvInd int) (err error) {
	ret := C.rsmi_dev_xgmi_error_reset(C.uint32_t(dvInd))
	glog.V(5).Infof(" rsmi_dev_hsl_error_reset ret:%v,retStr:%v", ret, errorString(ret))
	if err = errorString(ret); err != nil {
		return fmt.Errorf("Error rsmiDevHSLErrorReset:%s", err)
	}
	return
}

// rsmiDevHSLErrorStatus 获取设备的HSL错误状态
func rsmiDevHSLErrorStatus(dvInd int) (status RSMIHSLStatus, err error) {
	var cStatus C.rsmi_xgmi_status_t
	ret := C.rsmi_dev_xgmi_error_status(C.uint32_t(dvInd), &cStatus)
	glog.V(5).Infof(" rsmi_dev_hsl_error_status ret:%v,retstr:%v", ret, errorString(ret))
	if err := errorString(ret); err != nil {
		return status, fmt.Errorf("Error RSMIDevHSLErrorStatus: %s", err)
	}
	status = RSMIHSLStatus(cStatus)
	glog.V(5).Infof("RSMIDevHSLErrorStatus:%v", status)
	return
}

// rsmiDevHSLHiveIdGet 获取设备的HSL hive id
func rsmiDevHSLHiveIdGet(dvInd int) (hiveId int64, err error) {
	var chiveId C.uint64_t
	ret := C.rsmi_dev_xgmi_hive_id_get(C.uint32_t(dvInd), &chiveId)
	glog.V(5).Infof("rsmi_dev_hsl_hive_id_get ret:%v", ret)
	if err = errorString(ret); err != nil {
		return hiveId, fmt.Errorf("Error rsmiDevHSLHiveIdGet:%s", err)
	}
	hiveId = int64(chiveId)
	glog.V(5).Infof("rsmi_dev_hsl_hive_id_get hiveId:%v", hiveId)
	return
}
