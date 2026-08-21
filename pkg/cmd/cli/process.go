/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/HYGON-AI/hcu-dcgm/v3/pkg/dcgm"
)

// -------------------- process list --------------------
// 显示系统中受管理的进程 ID 列表（ASCII 表格风格）
var processListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all managed process IDs",
	Long:  `Retrieve a list of process IDs (PIDs) managed by the system in a table format.`,
	Run: func(cmd *cobra.Command, args []string) {
		pidList, err := dcgm.PidList()
		if err != nil {
			fmt.Println("Error fetching PID list:", err)
			os.Exit(1)
		}
		title := "MANAGED PID LIST"
		headers := []string{"PID"}
		var rows [][]string
		for _, pid := range pidList {
			rows = append(rows, []string{pid})
		}
		// 如果没有 PID，则显示 None
		if len(rows) == 0 {
			rows = append(rows, []string{"None"})
		}

		printAsciiTable(title, headers, rows)
	},
}

// -------------------- process kfd --------------------
// 显示正在运行的 KFD 进程的详细信息
var processKfdCmd = &cobra.Command{
	Use:   "kfd",
	Short: "Show running KFD process information",
	Long: `Retrieve and display detailed information about 
KFD processes currently running on the system.

Example:
  dcgmi process kfd`,
	Run: func(cmd *cobra.Command, args []string) {
		err := dcgm.ShowPids()
		if err != nil {
			fmt.Println("Error displaying KFD process information:", err)
			os.Exit(1)
		}
	},
}

// -------------------- process info --------------------
// 显示与 HCU 相关的进程信息
var processInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show HCU process information",
	Long: `Retrieve detailed information about HCU processes,
including memory usage, CU occupancy, and associated HCUs.

Example:
  dcgmi process info`,
	Run: func(cmd *cobra.Command, args []string) {
		processes, err := dcgm.ProcessHCUInfo()
		if err != nil {
			fmt.Printf("Error fetching HCU process info: %v\n", err)
			os.Exit(1)
		}

		title := "HCU PROCESSES"
		headers := []string{"PID", "PROCESS_NAME", "PASID", "VRAM_USAGE", "SDMA_USAGE", "CU_OCC", "HCU"}

		var rows [][]string
		if len(processes) == 0 {
			rows = append(rows, []string{"None", "", "", "", "", "", ""})
		} else {
			for _, process := range processes {
				vramInMB := process.VramUsage / (1024 * 1024)
				vramUsageStr := fmt.Sprintf("%dMB", vramInMB)

				minorNumbers := fmt.Sprintf("[%s]",
					strings.Trim(strings.Replace(fmt.Sprint(process.MinorNumbers), " ", " ", -1), "[]"))

				rows = append(rows, []string{
					fmt.Sprintf("%d", process.ProcessID),
					process.ProcessName,
					fmt.Sprintf("%d", process.Pasid),
					vramUsageStr,
					fmt.Sprintf("%d", process.SdmaUsage),
					fmt.Sprintf("%d", process.CuOccupancy),
					minorNumbers,
				})
			}
		}

		printAsciiTable(title, headers, rows)
	},
}

// -------------------- process root --------------------
// process 根命令，用来统一管理进程相关子命令
var processCmd = &cobra.Command{
	Use:   "process",
	Short: "Manage and inspect HCU processes",
	Long: `Provides commands to list and inspect processes
associated with HCU devices.`,
}

func init() {
	// 注册 process 子命令和子命令集合
	processCmd.AddCommand(processListCmd)
	processCmd.AddCommand(processKfdCmd)
	processCmd.AddCommand(processInfoCmd)

	rootCmd.AddCommand(processCmd)
}

// -------------------- 辅助函数：打印 ASCII 表格 --------------------
func printAsciiTable(title string, headers []string, rows [][]string) {
	colWidths := make([]int, len(headers))
	for i, h := range headers {
		colWidths[i] = len(h)
	}
	for _, row := range rows {
		for i, col := range row {
			if len(col) > colWidths[i] {
				colWidths[i] = len(col)
			}
		}
	}

	printLine := func() {
		for _, w := range colWidths {
			fmt.Print("+", strings.Repeat("-", w+2))
		}
		fmt.Println("+")
	}

	fmt.Println()
	fmt.Printf("+%s+\n", strings.Repeat("=", sum(colWidths)+3*len(colWidths)-1))
	fmt.Printf("| %-*s |\n", sum(colWidths)+3*len(colWidths)-1, title)
	fmt.Printf("+%s+\n", strings.Repeat("=", sum(colWidths)+3*len(colWidths)-1))

	printLine()
	// 打印表头
	fmt.Print("|")
	for i, h := range headers {
		fmt.Printf(" %-*s |", colWidths[i], h)
	}
	fmt.Println()
	printLine()
	// 打印数据行
	for _, row := range rows {
		fmt.Print("|")
		for i, col := range row {
			fmt.Printf(" %-*s |", colWidths[i], col)
		}
		fmt.Println()
	}
	printLine()
	fmt.Println()
}

// 计算列总宽度
func sum(arr []int) int {
	total := 0
	for _, v := range arr {
		total += v
	}
	return total
}
