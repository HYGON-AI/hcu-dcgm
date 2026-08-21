/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package cli

import (
	"fmt"
	"github.com/HYGON-AI/hcu-dcgm/v3/pkg/dcgm"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "View available profiling metrics for HCUs",
	Long:  `profile -- View available profiling metrics for HCUs`,
	Example: `  dcgmi profile -l 
  dcgmi profile -l -i <entityId>
  dcgmi profile -l -g <groupId>`,
	Run: func(cmd *cobra.Command, args []string) {
		if !listFlag {
			fmt.Println("PARSE ERROR: Required argument missing: list")
			return
		} else if entityId != "" && groupId != "" {
			fmt.Println("Error: Both entity and group IDs specified. Please use only one at a time.")
			return
		} else if entityId != "" {
			handleProfileWithEntity()
			return
		} else if groupId != "" {
			handleProfileWithGroup()
			return
		} else {
			handleProfileList()
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(profileCmd)
	profileCmd.Flags().BoolVarP(&listFlag, "list", "l", false, "List available profiling metrics")
	profileCmd.Flags().StringVarP(&entityId, "entityId", "i", "", "Comma-seperated list of entity IDs to query.\n"+
		"Default is supported HCUs on the system. Run\ndcgmi discovery -l to check list of HCUs available")
	profileCmd.Flags().StringVarP(&groupId, "group", "g", "", "The group of HCUs to query on the host.")
}

func setHcuIdFromEntityList(entityId string) (hcuIndex int, err error) {
	entityList, err := parseEntityList(entityId)
	if err != nil {
		return 0, fmt.Errorf("Invalid --entityid input: %v\n", err)
	}
	for _, entity := range entityList {
		if entity.EntityGroupId == dcgm.FE_HCU {
			mHcuId := entity.EntityId
			return mHcuId, nil
		}
	}
	return 0, fmt.Errorf("Error: No HCUs found in the provided entity list.")
}

func printProfilingMetrics(hcuIndex int) {
	metricGroups, err := dcgm.GetSupportedMetricGroups(hcuIndex)
	if err != nil {
		fmt.Printf("Error getting supported metric groups: %v\n", err)
		return
	}
	fmt.Printf("+----------------+----------+------------------------------------------------------+\n" +
		"| Group.Subgroup | Field ID | Field Tag                                            |\n" +
		"+----------------+----------+------------------------------------------------------+\n")
	for _, metricGroup := range metricGroups {
		groupStr := fmt.Sprintf("%s.%d", string(rune('A'-1+metricGroup.Major)), metricGroup.Minor)
		for _, fieldId := range metricGroup.FieldIds {
			fieldName := dcgm.FieldIdToName[fieldId]
			fmt.Printf("| %-15s| %-9d| %-53s|\n", groupStr, fieldId, strings.ToLower(strings.TrimPrefix(fieldName, "HCU_")))
		}
	}
	fmt.Println("+----------------+----------+------------------------------------------------------+")
}

func handleProfileWithEntity() {
	hcuIndex, err := setHcuIdFromEntityList(entityId)
	if err != nil {
		fmt.Println(err)
		return
	}
	printProfilingMetrics(hcuIndex)
}

func handleProfileWithGroup() {
	groupId, err := strconv.Atoi(groupId)
	if err != nil {
		fmt.Printf("Invalid Group ID specified.\n")
		return
	}
	hcuInGroup, _, err := dcgm.GetHCUListFromGroup(groupId)
	if err != nil {
		fmt.Printf("Error getting group info: %v\n", err)
		return
	}
	if len(hcuInGroup) == 0 {
		fmt.Printf("Failed to query group: no HCU found for group %v\n", groupId)
		return
	}
	printProfilingMetrics(hcuInGroup[0])
}

func handleProfileList() {
	numDevices, err := dcgm.NumMonitorDevices()
	if err != nil {
		fmt.Printf("Failed to get HCUs: %v\n", err)
		return
	}
	if numDevices < 1 {
		fmt.Println("Error: found 0 HCUs")
		return
	}
	printProfilingMetrics(0)
}
