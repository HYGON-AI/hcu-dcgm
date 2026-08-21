/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package cli

import (
	"fmt"
	"github.com/HYGON-AI/hcu-dcgm/v3/pkg/dcgm"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var (
	entityId string
	fieldId  string
	delay    int
	count    int
)

var dmonCmd = &cobra.Command{
	Use:   "dmon",
	Short: "Used to monitor HCUs and their stats.",
	Long:  `dmon -- Used to monitor HCUs and their stats.`,
	Example: `  dcgmi dmon -l 
  dcgmi dmon -e <fieldId>
  dcgmi dmon -f <fieldGroupId>
  dcgmi dmon -e <fieldId> -i <entityId> 
  dcgmi dmon -e <fieldId> -g <groupId> 
  dcgmi dmon -f <fieldGroupId> -i <entityId>
  dcgmi dmon -f <fieldGroupId> -g <groupId>
  dcgmi dmon -e <fieldId> -i <entityId> -d <delaySeconds> -c <count>`,
	Run: func(cmd *cobra.Command, args []string) {
		if listFlag {
			if groupId != "" || entityId != "" || fieldId != "" || delay != 3 || count != -1 {
				fmt.Println("Error: Invalid parameters with list arg. Usage : dmon -l.")
				return
			}
			handleFieldList()
			return
		}
		if groupId != "" && entityId != "" {
			fmt.Println("Error: Only one of --entityid and --groupid can be provided.")
			return
		}
		if count <= 0 && count != -1 {
			fmt.Println("Error: Positive value expected, negative value found for arg count.")
			return
		}
		if delay < 1 {
			fmt.Println("Error: Invalid value for arg delay.")
			return
		}
		if fieldId != "" {
			if fieldGroupId != "" {
				fmt.Println("PARSE ERROR: Argument: -f AND -e provided. Only one is allowed.")
				return
			}
			if entityId != "" {
				handleFieldWithEntity()
			} else if groupId != "" {
				handleFieldWithGroup()
			} else {
				handleFieldWithAllHcu()
			}
			return
		}
		if fieldGroupId != "" {
			if entityId != "" {
				handleFieldGroupWithEntity()
			} else if groupId != "" {
				handleFieldGroupWithGroup()
			} else {
				handleFieldGroupWithAllHcu()
			}
			return
		}
		fmt.Println("PARSE ERROR: Required argument missing: {fieldgroupid | fieldid | list}")
	},
}

func init() {
	rootCmd.AddCommand(dmonCmd)
	dmonCmd.Flags().BoolVarP(&listFlag, "list", "l", false, "List to look up the field names and field ids.")
	dmonCmd.Flags().StringVarP(&fieldGroupId, "fieldgroupid", "f", "", "The field group to query on the host.")
	dmonCmd.Flags().StringVarP(&fieldId, "fieldid", "e", "", "Field identifier to view/inject.")
	dmonCmd.Flags().StringVarP(&entityId, "entityid", "i", "", "Comma-separated list of entities to run the dmon on.\n(csv hcuIds like 0,1 or entityIds like hcu:0,hcu:1).\nDefault is -1 which runs for all supported HCU.")
	dmonCmd.Flags().StringVarP(&groupId, "groupid", "g", "", "The group to query on the host.")
	dmonCmd.Flags().IntVarP(&delay, "delay", "d", 3, "Time in seconds between each run.")
	dmonCmd.Flags().IntVarP(&count, "count", "c", -1, "Integer representing How many times to loop \nbefore exiting. [default- runs forever.]")
}

func handleFieldList() {
	fieldMetaList := dcgm.ListFieldMeta()

	fmt.Println("+-------------------+----------------------------------------------------------+")
	fmt.Println("| Field ID          | Field Name                                               |")
	fmt.Println("+-------------------+----------------------------------------------------------+")
	for _, fieldMeta := range fieldMetaList {
		fmt.Printf("| %-17d | %-56s |\n", fieldMeta.FieldId, fieldMeta.Name)
	}
	fmt.Println("+-------------------+----------------------------------------------------------+")
}

func parseFieldList(fieldIdStr string) (fieldIdList []int, fieldNameList []string, err error) {
	items := strings.Split(fieldIdStr, ",")
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		fieldId, err := strconv.Atoi(item)
		if err != nil {
			return fieldIdList, fieldNameList, fmt.Errorf("invalid field ID in '%s': %v\n", item, err)
		}
		fieldName := dcgm.FieldIdToName[fieldId]
		fieldIdList = append(fieldIdList, fieldId)
		fieldNameList = append(fieldNameList, fieldName)
	}
	return fieldIdList, fieldNameList, nil
}

func parseFieldGroup(fieldGroupIdStr string) (fieldIdList []int, fieldNameList []string, err error) {
	fieldGroupId, err := strconv.Atoi(fieldGroupIdStr)
	if err != nil {
		return fieldIdList, fieldNameList, fmt.Errorf("invalid field group ID in '%s': %v\n", fieldGroupIdStr, err)
	}
	fieldGroupInfo, err := dcgm.GetFieldGroupInfo(fieldGroupId)
	if err != nil {
		return fieldIdList, fieldNameList, fmt.Errorf("invalid field group ID in '%s': %v\n", fieldGroupIdStr, err)
	}
	fieldIdList = fieldGroupInfo.FieldIds
	for _, fieldId := range fieldIdList {
		fieldName := dcgm.FieldIdToName[fieldId]
		fieldNameList = append(fieldNameList, fieldName)
	}
	return fieldIdList, fieldNameList, nil
}

func printLatestValues(entityList []dcgm.GroupEntityPair, fieldIdList []int, fieldNameList []string) {
	for _, entity := range entityList {
		err := dcgm.WatchFieldsWithEntity(entity.EntityGroupId, entity.EntityId, fieldIdList)
		if err != nil {
			fmt.Printf("Failed to get field values: %v\n", err)
			for _, entity = range entityList {
				dcgm.UnWatchFieldsWithEntity(entity.EntityGroupId, entity.EntityId)
			}
			return
		}
		fieldValueList, err := dcgm.EntityGetLatestValues(entity.EntityGroupId, entity.EntityId, fieldIdList)
		if err != nil {
			fmt.Printf("Failed to get field values: %v\n", err)
			for _, entity = range entityList {
				dcgm.UnWatchFieldsWithEntity(entity.EntityGroupId, entity.EntityId)
			}
			return
		}
		entityStr := fmt.Sprintf("%s %d", entity.EntityGroupId, entity.EntityId)
		fmt.Printf("%-23s", entityStr)
		for index, fieldValue := range fieldValueList {
			width := len(fieldNameList[index])
			format := fmt.Sprintf("%%-%d.3f", width+2)
			fmt.Printf(format, fieldValue.Value)
		}
		fmt.Println("")
	}
}

func loopPrintFieldValues(entityList []dcgm.GroupEntityPair, fieldIdList []int, fieldNameList []string) {
	for _, entity := range entityList {
		for _, fieldId := range fieldIdList {
			fieldMeta := dcgm.GetFieldMetaById(fieldId)
			if entity.EntityGroupId != fieldMeta.EntityLevel {
				fmt.Printf("Entity %s %d has no field %d(%s)\n", entity.EntityGroupId, entity.EntityId, fieldId, fieldMeta.Name)
				return
			}
		}
	}
	fmt.Print("#Entity ID             ")
	for _, fieldName := range fieldNameList {
		width := len(fieldName)
		format := fmt.Sprintf("%%-%ds", width+2)
		fmt.Printf(format, fieldName)
	}
	fmt.Println("")

	executed := 0
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	ticker := time.NewTicker(time.Duration(delay) * time.Second)
	defer ticker.Stop()

loop:
	for {
		select {
		case <-stop:
			fmt.Println("\ndmon was stopped do to receiving a signal.")
			for _, entity := range entityList {
				dcgm.UnWatchFieldsWithEntity(entity.EntityGroupId, entity.EntityId)
			}
			break loop
		case <-ticker.C:
			printLatestValues(entityList, fieldIdList, fieldNameList)
			executed++
			if count != -1 && executed >= count {
				for _, entity := range entityList {
					dcgm.UnWatchFieldsWithEntity(entity.EntityGroupId, entity.EntityId)
				}
				break loop
			}
		}
	}
}

func handleFieldWithEntity() {
	fieldIdList, fieldNameList, err := parseFieldList(fieldId)
	if err != nil {
		fmt.Printf("Invalid --fieldid input: %v\n", err)
		return
	}
	entityList, err := parseEntityList(entityId)
	if err != nil {
		fmt.Printf("Invalid --entityid input: %v\n", err)
		return
	}
	loopPrintFieldValues(entityList, fieldIdList, fieldNameList)
}

func handleFieldWithGroup() {
	fieldIdList, fieldNameList, err := parseFieldList(fieldId)
	if err != nil {
		fmt.Printf("Invalid --fieldid input: %v\n", err)
		return
	}
	groupInt, err := strconv.Atoi(groupId)
	if err != nil {
		fmt.Printf("Error getting entity list: %v\n", err)
		return
	}
	groupInfo, err := dcgm.GetGroupInfo(groupInt)
	if err != nil {
		fmt.Printf("Error getting entity list: %v\n", err)
		return
	}
	entityList := groupInfo.EntityList
	loopPrintFieldValues(entityList, fieldIdList, fieldNameList)
}

func handleFieldWithAllHcu() {
	fieldIdList, fieldNameList, err := parseFieldList(fieldId)
	if err != nil {
		fmt.Printf("Invalid --fieldid input: %v\n", err)
		return
	}
	numDevices, err := dcgm.NumMonitorDevices()
	if err != nil {
		fmt.Printf("Error getting entity list: %v\n", err)
		return
	}
	entityList := make([]dcgm.GroupEntityPair, numDevices)
	for i := 0; i < numDevices; i++ {
		entityList[i] = dcgm.GroupEntityPair{
			EntityGroupId: dcgm.FE_HCU,
			EntityId:      i,
		}
	}
	loopPrintFieldValues(entityList, fieldIdList, fieldNameList)
}

func handleFieldGroupWithEntity() {
	entityList, err := parseEntityList(entityId)
	if err != nil {
		fmt.Printf("Invalid --entityid input: %v\n", err)
		return
	}
	fieldIdList, fieldNameList, err := parseFieldGroup(fieldGroupId)
	if err != nil {
		fmt.Printf("Invalid --fieldgroup input: %v\n", err)
		return
	}
	loopPrintFieldValues(entityList, fieldIdList, fieldNameList)
}

func handleFieldGroupWithGroup() {
	fieldIdList, fieldNameList, err := parseFieldGroup(fieldGroupId)
	if err != nil {
		fmt.Printf("Invalid --fieldgroup input: %v\n", err)
		return
	}
	groupInt, err := strconv.Atoi(groupId)
	if err != nil {
		fmt.Printf("Error getting entity list: %v\n", err)
		return
	}
	groupInfo, err := dcgm.GetGroupInfo(groupInt)
	if err != nil {
		fmt.Printf("Error getting entity list: %v\n", err)
		return
	}
	entityList := groupInfo.EntityList
	loopPrintFieldValues(entityList, fieldIdList, fieldNameList)
}

func handleFieldGroupWithAllHcu() {
	fieldIdList, fieldNameList, err := parseFieldGroup(fieldGroupId)
	if err != nil {
		fmt.Printf("Invalid --fieldgroup input: %v\n", err)
		return
	}
	numDevices, err := dcgm.NumMonitorDevices()
	if err != nil {
		fmt.Printf("Error getting entity list: %v\n", err)
		return
	}
	entityList := make([]dcgm.GroupEntityPair, numDevices)
	for i := 0; i < numDevices; i++ {
		entityList[i] = dcgm.GroupEntityPair{
			EntityGroupId: dcgm.FE_HCU,
			EntityId:      i,
		}
	}
	loopPrintFieldValues(entityList, fieldIdList, fieldNameList)
}
