/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/HYGON-AI/hcu-dcgm/v3/pkg/dcgm"
)

var vDeviceInfoCmd = &cobra.Command{
	Use:   "vdevice-info [device-index]",
	Short: "Get virtual device information",
	Long:  `Retrieve detailed information about a virtual device using its device index.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// 将设备索引字符串转换为整数
		dvInd, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Invalid device index:", err)
			os.Exit(1)
		}

		// 获取设备信息
		info, err := dcgm.VDeviceSingleInfo(dvInd)
		if err != nil {
			fmt.Println("Error fetching virtual device info:", err)
			os.Exit(1)
		}

		// 格式化输出
		fmt.Println("========== Virtual Device Info ==========")
		fmt.Printf("Name:               %s\n", info.Name)
		fmt.Printf("ComputeUnitCount:   %d\n", info.VComputeUnitCount)
		fmt.Printf("VMemoryTotal:      %d bytes (%.2f GB)\n", info.VMemoryTotal, float64(info.VMemoryTotal)/(1024*1024*1024))
		fmt.Printf("VMemoryUsed:       %d bytes (%.2f GB)\n", info.VMemoryUsed, float64(info.VMemoryUsed)/(1024*1024*1024))
		fmt.Printf("ContainerID:        %d\n", info.ContainerID)
		fmt.Printf("DvInd:           	%d\n", info.DvInd)
		fmt.Printf("VPercent:            %d%%\n", info.VPercent)
		fmt.Printf("VdvInd:       		%d\n", info.VdvInd)
		fmt.Printf("PciBusNumber:       %s\n", info.PciBusNumber)
		fmt.Println("=========================================")
	},
}

var destroyVDeviceCmd = &cobra.Command{
	Use:   "destroy-vdevice<dvInd>",
	Short: "Destroy a single virtual device",
	Long:  `This command destroys a single virtual device by its index.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		vDvInd, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Invalid virtual device index:", err)
			os.Exit(1)
		}

		err = dcgm.DestroySingleVDevice(vDvInd)
		if err != nil {
			fmt.Println("Error destroying virtual device:", err)
			os.Exit(1)
		}

		fmt.Printf("Virtual device %d destroyed successfully.\n", vDvInd)
	},
}
var allDeviceInfosCmd = &cobra.Command{
	Use:   "all-device-infos",
	Short: "Get information for all physical devices",
	Long:  `Retrieve detailed information about all physical devices.`,
	Run: func(cmd *cobra.Command, args []string) {
		infos, err := dcgm.AllDeviceInfos()
		if err != nil {
			fmt.Println("Error fetching all device infos:", err)
			os.Exit(1)
		}
		fmt.Println("==========allDevices==========")
		fmt.Printf(dataToJson(infos))

	},
}

func init() {
	rootCmd.AddCommand(vDeviceInfoCmd)
	rootCmd.AddCommand(destroyVDeviceCmd)
	rootCmd.AddCommand(allDeviceInfosCmd)
}
