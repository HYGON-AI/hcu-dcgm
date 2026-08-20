/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package dcgm

import "fmt"

func getFieldValue(fieldEntityGroup Field_Entity_Group, entityId int, fieldId int) (value float64, err error) {
	switch fieldEntityGroup {
	case FE_HCU:
		return getDevFieldValue(entityId, fieldId)
	case FE_VHCU:
		return getVdcuFieldValue(entityId, fieldId)
	default:
		err = fmt.Errorf("currently not supported field entity group %s", fieldEntityGroup.String())
		return
	}
}

func getDevFieldValue(hcuIndex int, fieldId int) (value float64, err error) {
	switch fieldId {
	case HCU_TEMP:
		return Temperature(hcuIndex)
	case HCU_POWER_USAGE:
		power, err := Power(hcuIndex)
		return float64(power), err
	case HCU_POWER_CAP:
		powerCap, err := MaxPower(hcuIndex)
		return float64(powerCap), err
	case HCU_UTILIZATION_RATE:
		utilizationRate, err := HCUUse(hcuIndex)
		return float64(utilizationRate), err
	case HCU_USED_MEMORY_BYTES:
		return MemoryUsed(hcuIndex)
	case HCU_MEMORY_CAP_BYTES:
		return MemoryTotal(hcuIndex)
	case HCU_MEMORY_PERCENT:
		memoryPercent, err := MemoryPercent(hcuIndex)
		return float64(memoryPercent), err
	case HCU_PCIE_BW_MB, HCU_PCIE_RECEIVE_MB, HCU_PCIE_SENT_MB:
		pcieBandwidth, err := PcieBw(hcuIndex)
		if err != nil {
			return 0, err
		}
		switch fieldId {
		case HCU_PCIE_BW_MB:
			return pcieBandwidth.Sent + pcieBandwidth.Received, nil
		case HCU_PCIE_RECEIVE_MB:
			return pcieBandwidth.Received, nil
		case HCU_PCIE_SENT_MB:
			return pcieBandwidth.Sent, nil
		}
	case HCU_SCLK:
		return HCUClk(hcuIndex)
	case HCU_COMPUTE_UNIT_COUNT:
		deviceInfo, err := GetDeviceInfo(hcuIndex)
		return float64(deviceInfo.ComputeUnitCount), err
	case HCU_COMPUTE_UNIT_REMAINING_COUNT, HCU_MEMORY_REMAINING:
		cus, mem, err := DeviceRemainingInfo(hcuIndex)
		if err != nil {
			return 0, err
		}
		if fieldId == HCU_COMPUTE_UNIT_REMAINING_COUNT {
			return float64(cus), nil
		}
		return float64(mem), nil
	case HCU_DF_BW_READ, HCU_DF_BW_WRITE, HCU_DF_BW_READ_WRITE:
		bandwidth, err := DFBandwidth(hcuIndex, RSMI_DF_BW_TYPE_ALL)
		if err != nil {
			return 0, err
		}
		switch fieldId {
		case HCU_DF_BW_READ:
			return bandwidth.ReadBW, nil
		case HCU_DF_BW_WRITE:
			return bandwidth.WriteBW, nil
		case HCU_DF_BW_READ_WRITE:
			return bandwidth.ReadWriteBW, nil
		}
	case HCU_VHCU_COUNT:
		vDeviceCount, _, err := VDeviceByDvInd(hcuIndex)
		return float64(vDeviceCount), err
	default:
		err = fmt.Errorf("unknown field %v", fieldId)
		return 0, err
	}
	return 0, nil
}

func getVdcuFieldValue(vdcuIndex int, fieldId int) (value float64, err error) {
	vDcuInfo, err := dmiGetVDeviceInfo(vdcuIndex)
	if err != nil {
		return 0, err
	}
	switch fieldId {
	case VHCU_SCLK:
		return HCUClk(vDcuInfo.DvInd)
	case VHCU_TEMP:
		return Temperature(vDcuInfo.DvInd)
	case VHCU_UTILIZATION_RATE:
		return float64(vDcuInfo.VPercent), nil
	case VHCU_USED_MEMORY_BYTES:
		return float64(vDcuInfo.VMemoryUsed), nil
	default:
		err = fmt.Errorf("unknown field %v", fieldId)
		return 0, err
	}
}
