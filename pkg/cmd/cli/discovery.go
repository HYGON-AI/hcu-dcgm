/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package cli

import (
	"fmt"
	"github.com/HYGON-AI/hcu-dcgm/v3/pkg/dcgm"
	"github.com/spf13/cobra"
	"slices"
	"strconv"
	"strings"
)

const NOT_APPLICABLE = "****"

var (
	infoFlags            string
	verboseFlag          bool
	computeHierarchyFlag bool
	hcuIndex             string
)

var discoveryCmd = &cobra.Command{
	Use:   "discovery",
	Short: "Used to discover and identify HCUs and their attributes.",
	Long:  `discovery -- Used to discover and identify HCUs and their attributes.`,
	Example: `  dcgmi discovery -l
  dcgmi discovery -i <flags> --hcuid <hcuId>
  dcgmi discovery -i <flags> -g <groupId> -v
  dcgmi discovery -c`,
	Run: func(cmd *cobra.Command, args []string) {
		// Main dispatcher logic
		switch {
		case listFlag:
			handleDiscoveryList()
		case computeHierarchyFlag:
			handleComputeHierarchy()
		case infoFlags != "":
			handleInfoFlagsOperations()
		default:
			cmd.Help()
		}
	},
}

func init() {
	rootCmd.AddCommand(discoveryCmd)

	discoveryCmd.Flags().BoolVarP(&listFlag, "list", "l", false, "List all the HCUs discovery on the host.")
	discoveryCmd.Flags().StringVarP(&infoFlags, "info", "i", "", "Specify which information to return.\n"+
		" a - device info\n p - power limits\n t - thermal limits\n c - clocks")
	discoveryCmd.Flags().StringVar(&hcuIndex, "hcuid", "", "The HCU ID to query.")
	discoveryCmd.Flags().StringVarP(&groupId, "group", "g", "", "The group ID to query.")
	discoveryCmd.Flags().BoolVarP(&computeHierarchyFlag, "compute-hierarchy", "c", false, "List all of the gpu instances and compute instances.")
	discoveryCmd.Flags().BoolVarP(&verboseFlag, "verbose", "v", false, "Display information per HCU.")
}

// Handlers for subcommands
func handleDiscoveryList() {
	if computeHierarchyFlag {
		fmt.Printf("PARSE ERROR: Argument: -l AND -c provided.\n             Only one is allowed.\n")
		return
	}

	numHcus, err := dcgm.NumMonitorDevices()
	if err != nil {
		fmt.Printf("Error getting device count: %v\n", err)
		return
	}
	suffix := ""
	if numHcus != 1 {
		suffix = "s"
	}
	fmt.Printf("%d HCU%s found.\n", numHcus, suffix)
	fmt.Println("+--------+----------------------------------------------------------------------+")
	fmt.Println("| HCU ID | Device Information                                                   |")
	fmt.Println("+--------+----------------------------------------------------------------------+")

	for i := 0; i < numHcus; i++ {
		hcuName, err := dcgm.DevName(i)
		if err != nil {
			hcuName = "N/A"
		}
		hcuTypeName, _, err := dcgm.DevTypeName(i)
		if err != nil {
			hcuTypeName = "N/A"
		}
		nameStr := hcuName + " - " + hcuTypeName
		hcuUniqueId, err := dcgm.GetDeviceUniqueId(i)
		if err != nil {
			hcuUniqueId = "N/A"
		}
		pciId, err := dcgm.PicBusInfo(i)
		if err != nil {
			pciId = "N/A"
		}
		fmt.Printf("| %-7d| Name: %-63s|\n", i, nameStr)
		fmt.Printf("|        | PCI Bus ID: %-57s|\n", pciId)
		fmt.Printf("|        | Device Unique ID: %-51s|\n", hcuUniqueId)
		fmt.Println("+--------+----------------------------------------------------------------------+")
	}
}

func handleComputeHierarchy() {
	if groupId != "" || hcuIndex != "" {
		fmt.Printf("For now, hierarchy must be used by itself.\n")
		return
	}

	if infoFlags != "" {
		fmt.Printf("PARSE ERROR: Argument: -i AND -c provided.\n             Only one is allowed.\n")
		return
	}

	fmt.Printf("+-------------------+--------------------------------------------------------------------+\n" +
		"| Instance Hierarchy                                                                     |\n" +
		"+===================+====================================================================+\n")
	migMode, _, err := dcgm.SystemMigMode()
	if err != nil || migMode != 1 {
		fmt.Println("+-------------------+--------------------------------------------------------------------+")
		return
	}
	migDevicesInfos, err := dcgm.MigInfos()
	if err != nil {
		fmt.Println("+-------------------+--------------------------------------------------------------------+")
		return
	}
	numHcus, err := dcgm.NumMonitorDevices()
	if err != nil {
		fmt.Println("+-------------------+--------------------------------------------------------------------+")
		return
	}
	migInfoList := make([][]dcgm.MigInfo, numHcus)
	for _, migInfo := range migDevicesInfos {
		migInfoList[migInfo.DvInd] = append(migInfoList[migInfo.DvInd], migInfo)
	}
	for i := 0; i < numHcus; i++ {
		if len(migInfoList[i]) > 0 {
			hcuUniqueId, err := dcgm.GetDeviceUniqueId(i)
			if err != nil {
				hcuUniqueId = "N/A"
			}
			fmt.Printf("| HCU %-14d| HCU %-63s|\n", i, hcuUniqueId)
			for _, migInfo := range migInfoList[i] {
				giIndexStr := fmt.Sprintf("%d/%d", i, migInfo.GpuInstanceId)
				ciIndexStr := fmt.Sprintf("%d/%d/%d", i, migInfo.GpuInstanceId, migInfo.ComputeInstanceId)
				ciStr := fmt.Sprintf("Compute Instance (%s) created by Profile %d", migInfo.Name, migInfo.CiProfileId)
				fmt.Printf("| -> I %-13s| GPU Instance created by Profile %-35d|\n", giIndexStr, migInfo.GiProfileId)
				fmt.Printf("|    -> CI %-9s| %-66s |\n", ciIndexStr, ciStr)
			}
			fmt.Println("+-------------------+--------------------------------------------------------------------+")
		}
	}
}

func queryHcuInfo(hcuIndex int, infoFlags string) error {
	fmt.Printf("+--------------------------+-------------------------------------------------+\n"+
		"| HCU ID: %-17d| Device Information                              |\n"+
		"+==========================+=================================================+\n", hcuIndex)
	if strings.Contains(infoFlags, "a") {
		nameStr, pciId, hcuUniqueId, hcuSerialNumber, vbiosVersion := getIdentifiers(hcuIndex)
		printIdentifiers(nameStr, pciId, hcuUniqueId, hcuSerialNumber, vbiosVersion)
	}
	if strings.Contains(infoFlags, "p") {
		powerAveStr, powerMaxStr, powerMinStr := getPowerLimits(hcuIndex)
		printPowerLimits(powerAveStr, powerMaxStr, powerMinStr)
	}
	if strings.Contains(infoFlags, "t") {
		tempCurrentStr, tempShutdownStr, tempCriticalStr, tempSlowdownStr := getThermals(hcuIndex)
		printThermals(tempCurrentStr, tempShutdownStr, tempCriticalStr, tempSlowdownStr)
	}
	if strings.Contains(infoFlags, "c") {
		mclkStrList, sclkStrList := getClocks(hcuIndex)
		printClocks(mclkStrList, sclkStrList)
	}
	return nil
}

func queryGroupInfo(groupId int, infoFlags string, verboseFlag bool) error {
	groupInfo, err := dcgm.GetGroupInfo(groupId)
	if err != nil {
		return err
	}
	entityList := groupInfo.EntityList
	if len(entityList) == 0 {
		return fmt.Errorf("no entity found for group %d", groupId)
	}
	var hcuInGroup []int
	for _, entity := range entityList {
		if entity.EntityGroupId == dcgm.FE_HCU {
			hcuInGroup = append(hcuInGroup, entity.EntityId)
		}
	}
	if verboseFlag {
		for _, entityId := range hcuInGroup {
			err = queryHcuInfo(entityId, infoFlags)
			if err != nil {
				return err
			}
		}
	} else {
		err = queryNonVerboseGroupInfo(infoFlags, hcuInGroup)
		if err != nil {
			return err
		}
	}
	return nil
}

func getIdentifiers(hcuIndex int) (nameStr, pciId, hcuUniqueId, hcuSerialNumber, vbiosVersion string) {
	hcuName, err := dcgm.DevName(hcuIndex)
	if err != nil {
		hcuName = "N/A"
	}
	hcuTypeName, _, err := dcgm.DevTypeName(hcuIndex)
	if err != nil {
		hcuTypeName = "N/A"
	}
	nameStr = hcuName + " - " + hcuTypeName
	pciId, err = dcgm.PicBusInfo(hcuIndex)
	if err != nil {
		pciId = "N/A"
	}
	hcuUniqueId, err = dcgm.GetDeviceUniqueId(hcuIndex)
	if err != nil {
		hcuUniqueId = "N/A"
	}
	hcuSerialNumber, err = dcgm.GetDeviceId(hcuIndex)
	if err != nil {
		hcuSerialNumber = "N/A"
	}
	vbiosVersion, err = dcgm.VbiosVersion(hcuIndex)
	if err != nil {
		vbiosVersion = "N/A"
	}
	return
}

func printIdentifiers(nameStr, pciId, hcuUniqueId, hcuSerialNumber, vbiosVersion string) {
	fmt.Printf("| Device Name              | %-48s|\n", nameStr)
	fmt.Printf("| PCI Bus ID               | %-48s|\n", pciId)
	fmt.Printf("| Unique ID                | %-48s|\n", hcuUniqueId)
	fmt.Printf("| Serial Number            | %-48s|\n", hcuSerialNumber)
	fmt.Printf("| VBIOS                    | %-48s|\n", vbiosVersion)
	fmt.Println("+--------------------------+-------------------------------------------------+")
}

func getPowerLimits(hcuIndex int) (powerAveStr, powerMaxStr, powerMinStr string) {
	powerAve, err := dcgm.Power(hcuIndex)
	if err != nil {
		powerAveStr = "N/A"
	} else {
		powerAveStr = strconv.Itoa(int(powerAve))
	}

	powerMax, powerMin, err := dcgm.DevPowerCapRange(hcuIndex)
	if err != nil {
		powerMinStr = "N/A"
		powerMaxStr = "N/A"
	} else {
		powerMax = (powerMax / 1000000)
		powerMin = (powerMin / 1000000)
		powerMaxStr = strconv.Itoa(int(powerMax))
		powerMinStr = strconv.Itoa(int(powerMin))
	}
	return
}

func printPowerLimits(powerAveStr, powerMaxStr, powerMinStr string) {
	fmt.Printf("| Power Ave Value (W)      | %-48s|\n", powerAveStr)
	fmt.Printf("| Power Max Value (W)      | %-48s|\n", powerMaxStr)
	fmt.Printf("| Power Min Value (W)      | %-48s|\n", powerMinStr)
	fmt.Println("+--------------------------+-------------------------------------------------+")
}

func getThermals(hcuIndex int) (tempCurrentStr, tempShutdownStr, tempCriticalStr, tempSlowdownStr string) {
	tempCurrent, err := dcgm.GetTempByMetric(hcuIndex, dcgm.RSMI_TEMP_CURRENT)
	if err != nil {
		tempCurrentStr = "N/A"
	} else {
		tempCurrentStr = strconv.FormatFloat(tempCurrent, 'f', -1, 64)
	}

	tempSlowdown, err := dcgm.GetTempByMetric(hcuIndex, dcgm.RSMI_TEMP_MAX)
	if err != nil {
		tempSlowdownStr = "N/A"
	} else {
		tempSlowdownStr = strconv.FormatFloat(tempSlowdown, 'f', -1, 64)
	}

	tempCritical, err := dcgm.GetTempByMetric(hcuIndex, dcgm.RSMI_TEMP_CRITICAL)
	if err != nil {
		tempCriticalStr = "N/A"
	} else {
		tempCriticalStr = strconv.FormatFloat(tempCritical, 'f', -1, 64)
	}

	tempShutdown, err := dcgm.GetTempByMetric(hcuIndex, dcgm.RSMI_TEMP_EMERGENCY)
	if err != nil {
		tempShutdownStr = "N/A"
	} else {
		tempShutdownStr = strconv.FormatFloat(tempShutdown, 'f', -1, 64)
	}
	return
}

func printThermals(tempCurrentStr, tempShutdownStr, tempCriticalStr, tempSlowdownStr string) {
	fmt.Printf("| Current Temperature (C)  | %-48s|\n", tempCurrentStr)
	fmt.Printf("| ShutDown Temperature (C) | %-48s|\n", tempShutdownStr)
	fmt.Printf("| Critical Temperature (C) | %-48s|\n", tempCriticalStr)
	fmt.Printf("| Slowdown Temperature (C) | %-48s|\n", tempSlowdownStr)
	fmt.Println("+--------------------------+-------------------------------------------------+")
}

func getClocks(hcuIndex int) (mclkStrList, sclkStrList []string) {
	mclkList, mclkCurrent, err := dcgm.GetClocksByType(hcuIndex, dcgm.RSMI_CLK_TYPE_MEM)

	if err != nil {
		mclkStrList = append(mclkStrList, "N/A")
	} else {
		for index, mclk := range mclkList {
			mclkStr := strconv.FormatUint(mclk, 10)
			if int(mclkCurrent) == index {
				mclkStr = mclkStr + " *"
			}
			mclkStrList = append(mclkStrList, mclkStr)
		}
	}
	sclkList, sclkCurrent, err := dcgm.GetClocksByType(hcuIndex, dcgm.RSMI_CLK_TYPE_SYS)

	if err != nil {
		sclkStrList = append(sclkStrList, "N/A")
	} else {
		for index, sclk := range sclkList {
			sclkStr := strconv.FormatUint(sclk, 10)
			if int(sclkCurrent) == index {
				sclkStr = sclkStr + " *"
			}
			sclkStrList = append(sclkStrList, sclkStr)
		}
	}
	return
}

func printClocks(mclkStrList, sclkStrList []string) {
	fmt.Println("| Supported Clocks (MHz)   | MCLK:                                           |")
	for _, mclk := range mclkStrList {
		fmt.Printf("|                          | %-48s|\n", mclk)
	}
	fmt.Println("|                          |                                                 |")
	fmt.Println("|                          | SCLK:                                           |")
	for _, sclk := range sclkStrList {
		fmt.Printf("|                          | %-48s|\n", sclk)
	}
	fmt.Println("+--------------------------+-------------------------------------------------+")
}

func queryNonVerboseGroupInfo(infoFlags string, hcuInGroup []int) error {
	fmt.Println("+--------------------------+-------------------------------------------------+")
	if len(hcuInGroup) == 1 {
		fmt.Println("| Group of 1 HCU           | Device Information                              |")

	} else {
		fmt.Printf("| Group of %d HCUs          | Device Information                              |\n", len(hcuInGroup))
	}
	fmt.Println("+==========================+=================================================+")
	if strings.Contains(infoFlags, "a") {
		tmpNameStr, tmpPciId, tmpHcuUniqueId, tmpHcuSerialNumber, tmpVbiosVersion := getIdentifiers(hcuInGroup[0])

		for i := 1; i < len(hcuInGroup); i++ {
			nameStr, pciId, hcuUniqueId, hcuSerialNumber, vbiosVersion := getIdentifiers(hcuInGroup[i])
			if tmpNameStr != NOT_APPLICABLE && tmpNameStr != nameStr {
				tmpNameStr = NOT_APPLICABLE
			}
			if tmpPciId != NOT_APPLICABLE && tmpPciId != pciId {
				tmpPciId = NOT_APPLICABLE
			}
			if tmpHcuUniqueId != NOT_APPLICABLE && tmpHcuUniqueId != hcuUniqueId {
				tmpHcuUniqueId = NOT_APPLICABLE
			}
			if tmpHcuSerialNumber != NOT_APPLICABLE && tmpHcuSerialNumber != hcuSerialNumber {
				tmpHcuSerialNumber = NOT_APPLICABLE
			}
			if tmpVbiosVersion != NOT_APPLICABLE && tmpVbiosVersion != vbiosVersion {
				tmpVbiosVersion = NOT_APPLICABLE
			}
		}
		printIdentifiers(tmpNameStr, tmpPciId, tmpHcuUniqueId, tmpHcuSerialNumber, tmpVbiosVersion)
	}
	if strings.Contains(infoFlags, "p") {
		tmpPowerAve, tmpPowerMax, tmpPowerMin := getPowerLimits(hcuInGroup[0])
		for i := 1; i < len(hcuInGroup); i++ {
			powerAve, powerMax, powerMin := getPowerLimits(hcuInGroup[i])
			if tmpPowerAve != NOT_APPLICABLE && tmpPowerAve != powerAve {
				tmpPowerAve = NOT_APPLICABLE
			}
			if tmpPowerMax != NOT_APPLICABLE && tmpPowerMax != powerMax {
				tmpPowerMax = NOT_APPLICABLE
			}
			if tmpPowerMin != NOT_APPLICABLE && tmpPowerMin != powerMin {
				tmpPowerMin = NOT_APPLICABLE
			}
		}
		printPowerLimits(tmpPowerAve, tmpPowerMax, tmpPowerMin)
	}
	if strings.Contains(infoFlags, "t") {
		tmpTempCurrent, tmpTempShutdown, tmpTempCritical, tmpTempSlowdown := getThermals(hcuInGroup[0])
		for i := 1; i < len(hcuInGroup); i++ {
			tempCurrent, tempShutdown, tempCritical, tempSlowdown := getThermals(hcuInGroup[i])
			if tmpTempCurrent != NOT_APPLICABLE && tmpTempCurrent != tempCurrent {
				tmpTempCurrent = NOT_APPLICABLE
			}
			if tmpTempShutdown != NOT_APPLICABLE && tmpTempShutdown != tempShutdown {
				tmpTempShutdown = NOT_APPLICABLE
			}
			if tmpTempCritical != NOT_APPLICABLE && tmpTempCritical != tempCritical {
				tmpTempCritical = NOT_APPLICABLE
			}
			if tmpTempSlowdown != NOT_APPLICABLE && tmpTempSlowdown != tempSlowdown {
				tmpTempSlowdown = NOT_APPLICABLE
			}
		}
		printThermals(tmpTempCurrent, tmpTempShutdown, tmpTempCritical, tmpTempSlowdown)
	}
	if strings.Contains(infoFlags, "c") {
		tmpMclkList, tmpSclkList := getClocks(hcuInGroup[0])
		for i := 1; i < len(hcuInGroup); i++ {
			mclkList, sclkList := getClocks(hcuInGroup[i])
			if tmpMclkList[0] != NOT_APPLICABLE && !slices.Equal(tmpMclkList, mclkList) {
				tmpMclkList = []string{NOT_APPLICABLE}
			}
			if tmpSclkList[0] != NOT_APPLICABLE && !slices.Equal(tmpSclkList, sclkList) {
				tmpSclkList = []string{NOT_APPLICABLE}
			}
		}
		printClocks(tmpMclkList, tmpSclkList)
	}
	fmt.Println("**** Non-homogenous settings across group. Use with –v flag to see details.")
	return nil
}

func queryAllHcusInfo(infoFlags string) error {
	numHcus, err := dcgm.NumMonitorDevices()
	if err != nil {
		return err
	}
	for i := 0; i < numHcus; i++ {
		err = queryHcuInfo(i, infoFlags)
		if err != nil {
			return err
		}
	}
	return nil
}

func handleInfoFlagsOperations() {
	for _, c := range infoFlags {
		switch c {
		case 'a', 't', 'p', 'c':
		default:
			fmt.Printf("Invalid input '%c'. Please include only valid tags.\n", c)
			return
		}
	}

	if groupId != "" && hcuIndex != "" {
		fmt.Printf("Both HCU and Group specified at command line.\n")
		return
	}

	if groupId != "" {
		groupId, err := strconv.Atoi(groupId)
		if err != nil {
			fmt.Printf("Invalid Group ID specified.\n")
			return
		}
		err = queryGroupInfo(groupId, infoFlags, verboseFlag)
		if err != nil {
			fmt.Printf("Failed to query group: %v\n", err)
			return
		}
	} else if hcuIndex != "" {
		hcuIndex, err := strconv.Atoi(hcuIndex)
		if err != nil {
			fmt.Printf("Invalid HCU index specified.\n")
			return
		}
		err = queryHcuInfo(hcuIndex, infoFlags)
		if err != nil {
			fmt.Printf("Failed to query HCU: %v\n", err)
			return
		}
	} else if verboseFlag {
		err := queryAllHcusInfo(infoFlags)
		if err != nil {
			fmt.Printf("Failed to query all HCUs: %v\n", err)
			return
		}
	}

}
