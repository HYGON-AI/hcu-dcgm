/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package cli

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/HYGON-AI/hcu-dcgm/v3/pkg/dcgm"
)

var healthCheckCmd = &cobra.Command{
	Use:   "health [HCU IDs]",
	Short: "Run HCU health check",
	Long: `Run a health check on specified HCU devices.

If no device IDs are provided, the health check will run on all available HCUs.

Usage:
  dcgmi health [HCU IDs]

Examples:
  dcgmi health
  dcgmi health 0
  dcgmi health 0 1 2`,
	Run: func(cmd *cobra.Command, args []string) {
		var dvIdList []int
		var err error

		if len(args) == 0 {
			// 如果没有输入 ID，则检查全部 HCU
			numDevices, err := dcgm.NumMonitorDevices()
			if err != nil {
				fmt.Println("Error retrieving device count:", err)
				os.Exit(1)
			}
			for i := 0; i < numDevices; i++ {
				dvIdList = append(dvIdList, i)
			}
		} else {
			// 把命令行参数转成 int slice
			for _, arg := range args {
				id, convErr := strconv.Atoi(arg)
				if convErr != nil {
					fmt.Printf("Invalid device ID: %s\n", arg)
					os.Exit(1)
				}
				dvIdList = append(dvIdList, id)
			}
		}

		// 调用健康检查
		deviceHealths, err := dcgm.HealthCheckById(dvIdList, false)
		if err != nil {
			fmt.Println("Error performing health check:", err)
			os.Exit(1)
		}

		// 打印结果
		printHealthCheckResults(deviceHealths)
	},
}

func init() {
	rootCmd.AddCommand(healthCheckCmd)
}

// 打印健康检查结果的函数
func printHealthCheckResults(deviceHealths []dcgm.DeviceHealth) {
	// 初始化 tabwriter，用于生成表格样式的输出
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.AlignRight|tabwriter.Debug)

	// 表头
	fmt.Fprintln(w, "DEVICE ID\tOVERALL STATUS\tCHECK TYPE\tSTATUS\tDETAILS")

	// 遍历每个设备的健康检查结果
	for _, dh := range deviceHealths {
		for _, watch := range dh.Watches {
			// 打印健康检查信息
			var details string
			if watch.Error != "" {
				details = fmt.Sprintf("Error: %s", watch.Error)
			} else if watch.Result != nil {
				details = formatResult(watch.Result)
			}
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n", dh.HCU, dh.Status, watch.Type, watch.Status, details)
		}
		fmt.Fprintln(w, "----------------------------------------------------------------")
	}

	// 刷新 tabwriter 缓冲区，打印表格
	w.Flush()
}

// 格式化 Result 的辅助函数，将 JSON 数据转为简化的文本
func formatResult(result interface{}) string {
	switch v := result.(type) {
	case string:
		return v
	case []map[string]interface{}:
		// 如果是数组，则提取关键信息
		output := ""
		for _, item := range v {
			for k, val := range item {
				output += fmt.Sprintf("%s: %v ", k, val)
			}
			output += "| "
		}
		return output
	default:
		// 默认直接转成字符串
		return fmt.Sprintf("%v", result)
	}
}
