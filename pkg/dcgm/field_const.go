/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package dcgm

const (
	HCU_TEMP int = iota + 101
	HCU_POWER_USAGE
	HCU_POWER_CAP
	HCU_UTILIZATION_RATE
	HCU_SCLK
	HCU_COMPUTE_UNIT_COUNT
	HCU_COMPUTE_UNIT_REMAINING_COUNT
	HCU_USED_MEMORY_BYTES
	HCU_MEMORY_CAP_BYTES
	HCU_MEMORY_REMAINING
	HCU_MEMORY_PERCENT
	HCU_PCIE_BW_MB
	HCU_PCIE_RECEIVE_MB
	HCU_PCIE_SENT_MB
	HCU_DF_BW_READ
	HCU_DF_BW_WRITE
	HCU_DF_BW_READ_WRITE
	HCU_VHCU_COUNT
)

const (
	VHCU_SCLK int = iota + 201
	VHCU_TEMP
	VHCU_UTILIZATION_RATE
	VHCU_USED_MEMORY_BYTES
)

var dcgmFields = map[string]int{
	"HCU_TEMP":                         HCU_TEMP,
	"HCU_POWER_USAGE":                  HCU_POWER_USAGE,
	"HCU_POWER_CAP":                    HCU_POWER_CAP,
	"HCU_UTILIZATION_RATE":             HCU_UTILIZATION_RATE,
	"HCU_SCLK":                         HCU_SCLK,
	"HCU_COMPUTE_UNIT_COUNT":           HCU_COMPUTE_UNIT_COUNT,
	"HCU_COMPUTE_UNIT_REMAINING_COUNT": HCU_COMPUTE_UNIT_REMAINING_COUNT,
	"HCU_USED_MEMORY_BYTES":            HCU_USED_MEMORY_BYTES,
	"HCU_MEMORY_CAP_BYTES":             HCU_MEMORY_CAP_BYTES,
	"HCU_MEMORY_REMAINING":             HCU_MEMORY_REMAINING,
	"HCU_MEMORY_PERCENT":               HCU_MEMORY_PERCENT,
	"HCU_PCIE_BW_MB":                   HCU_PCIE_BW_MB,
	"HCU_PCIE_RECEIVE_MB":              HCU_PCIE_RECEIVE_MB,
	"HCU_PCIE_SENT_MB":                 HCU_PCIE_SENT_MB,
	"HCU_DF_BW_READ":                   HCU_DF_BW_READ,
	"HCU_DF_BW_WRITE":                  HCU_DF_BW_WRITE,
	"HCU_DF_BW_READ_WRITE":             HCU_DF_BW_READ_WRITE,
	"HCU_VHCU_COUNT":                   HCU_VHCU_COUNT,
	"VHCU_SCLK":                        VHCU_SCLK,
	"VHCU_TEMP":                        VHCU_TEMP,
	"VHCU_UTILIZATION_RATE":            VHCU_UTILIZATION_RATE,
	"VHCU_USED_MEMORY_BYTES":           VHCU_USED_MEMORY_BYTES,
}

var FieldIdToName map[int]string
var profilingMetrics []MetricGroup

var unsupportedFieldsByName = map[string][]int{
	"K100AI": {},
	"K100":   {},
	"Z100":   {},
	"BW":     {},
}

func init() {
	FieldIdToName = make(map[int]string, len(dcgmFields))
	for name, id := range dcgmFields {
		FieldIdToName[id] = name
	}
	profilingMetrics = []MetricGroup{
		MetricGroup{1, 1, []int{HCU_UTILIZATION_RATE, HCU_SCLK}},
		{2, 1, []int{HCU_MEMORY_PERCENT}},
		{3, 1, []int{HCU_PCIE_BW_MB, HCU_PCIE_SENT_MB, HCU_PCIE_RECEIVE_MB}},
	}
}

func getFieldId(fieldName string) (int, bool) {
	id, ok := dcgmFields[fieldName]
	return id, ok
}
