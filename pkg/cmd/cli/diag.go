/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package cli

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/HYGON-AI/hcu-dcgm/v3/pkg/dcgm"
)

// ==================== diag 命令 ====================
var diagCmd = &cobra.Command{
	Use:   "diag",
	Short: "Run diagnostics commands",
	Example: `  dcgmi diag -g <groupId> -i <flags>
  dcgmi diag r <diagLevel>
  dcgmi diag bandwidth <hcuId>
  dcgmi diag memtestCL <hcuId>
  dcgmi diag pcie 
  dcgmi diag xhcl
  dcgmi diag edpp
  dcgmi diag gemm`,
	Run: func(cmd *cobra.Command, args []string) {
		switch {
		case groupId != "":
			handleDiagGroup()
		case infoFlags != "":
			fmt.Println("Error: No group has been specified.")
		default:
			cmd.Help()
		}
	},
}

func formatDiagFailure(testName string, err error, logPaths ...string) string {
	lines := []string{fmt.Sprintf("%s failed: %v", testName, err)}
	for _, logPath := range logPaths {
		if logPath != "" {
			lines = append(lines, fmt.Sprintf("       Logs: %s", logPath))
		}
	}
	return strings.Join(lines, "\n") + "\n"
}

func formatPCIESummary(result dcgm.PcieResult) string {
	if len(result.HCUs) == 0 {
		return ""
	}

	hcus := append([]dcgm.HCUBW(nil), result.HCUs...)
	sort.Slice(hcus, func(i, j int) bool {
		return hcus[i].DvInd < hcus[j].DvInd
	})

	var output strings.Builder
	output.WriteString("=== PCIe Bandwidth Summary ===\n")
	for _, hcu := range hcus {
		fmt.Fprintf(&output, "HCU %d  System memory->HCU: %.2f GB/s  HCU->System memory: %.2f GB/s\n", hcu.DvInd, hcu.SysToFb, hcu.FbToSys)
	}
	return output.String()
}

func formatGEMMSummary(result dcgm.TargetStressResult) string {
	if len(result.Results) == 0 {
		return ""
	}

	resultsByHCU := make(map[int][]dcgm.GemmTestResult)
	for _, gemm := range result.Results {
		resultsByHCU[gemm.HCUId] = append(resultsByHCU[gemm.HCUId], gemm)
	}

	hcuIDs := make([]int, 0, len(resultsByHCU))
	for hcuID := range resultsByHCU {
		hcuIDs = append(hcuIDs, hcuID)
	}
	sort.Ints(hcuIDs)

	var output strings.Builder
	output.WriteString("=== GEMM Performance Summary ===\n")
	for _, hcuID := range hcuIDs {
		gemmResults := resultsByHCU[hcuID]
		passed := 0
		metrics := make([]string, 0, len(gemmResults))
		for _, gemm := range gemmResults {
			if gemm.Failed {
				metrics = append(metrics, fmt.Sprintf("%s=FAIL", gemm.GemmName))
				continue
			}
			passed++
			metrics = append(metrics, fmt.Sprintf("%s=%.2f", gemm.GemmName, gemm.Mean))
		}

		status := "PASS"
		if passed != len(gemmResults) {
			status = "WARN"
		}
		fmt.Fprintf(&output, "HCU %d  %s %d/%d  mean: %s\n", hcuID, status, passed, len(gemmResults), strings.Join(metrics, ", "))
	}
	return output.String()
}

// ==================== diag run ====================
var runDiagCmd = &cobra.Command{
	Use:   "r [level]",
	Short: "Run general diagnostics. Usage: dcgmi diag r <level> (1~4)",
	Long: `Run a comprehensive diagnostic check on all devices with specified level.

Level 1: Basic checks
  - Memory Check
  - picBus Check
  - Power Check
  - DTK Version Check
  - Driver Version Check
  - RSMI Version Check
  - VBIOS Version Check
  - card Compatibility Check

Level 2: Level 1 + 
  - PCIe Bandwidth Check
  - Memory Bandwidth Check
  - Memory Reserved Pages Check

Level 3: Level 2 +
  - PCIe Bandwidth Stress Test
  - xHCL Bandwidth Test
  - Target Stress Test

Level 4: Level 3 +
  - MemtestCL
  - EDPpTest`,
	Args: cobra.ExactArgs(1), // 必须传一个参数
	Run: func(cmd *cobra.Command, args []string) {
		level, err := strconv.Atoi(args[0])
		if err != nil || level < 1 || level > 4 {
			fmt.Println("Invalid level. Please provide a value between 1 and 4.")
			os.Exit(1)
		}

		results, err := dcgm.RunDiag(level)
		if err != nil {
			fmt.Print(formatDiagFailure("Diagnostics", err, "logs"))
			os.Exit(1)
		}
		printDiagnosticResults(results)
	},
}

// ==================== diag bandwidth ====================
var bandwidthCmd = &cobra.Command{
	Use:   "bandwidth [HCU IDs]",
	Short: "Run HCU bandwidth test",
	Long: `Run a HCU memory bandwidth stress test on specified devices.

If no device IDs are provided, the test will run on all available HCUs.

Usage:
  dcgmi diag bandwidth [HCU IDs]

Examples:
  dcgmi diag bandwidth
  dcgmi diag bandwidth 0
  dcgmi diag bandwidth 0 1 2

[HCU IDs] are the numeric IDs of the HCU devices you want to test.`,
	Run: func(cmd *cobra.Command, args []string) {
		var dvIdList []int

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
				id, err := strconv.Atoi(arg)
				if err != nil {
					fmt.Printf("Invalid HCU ID '%s'\n", arg)
					os.Exit(1)
				}
				dvIdList = append(dvIdList, id)
			}
		}

		fmt.Printf("Running memory bandwidth test for HCU: %v\n", dvIdList)
		bwResults, err := dcgm.BandwidthTestResult(dvIdList)
		if err != nil {
			fmt.Print(formatDiagFailure("Bandwidth test", err, dcgm.DiagLogDirBandwidth))
			os.Exit(1)
		}
		for hcu, gbps := range bwResults {
			fmt.Printf("HCU %d  Bandwidth: %.2f GB/s\n", hcu, gbps)
		}
		fmt.Println("memory bandwidth test: Done.")
	},
}

// ==================== diag pcieBandwidth ====================
var pcieCmd = &cobra.Command{
	Use:   "pcie",
	Short: "Run PCIe memory bandwidth test",
	Long: `Run a PCIe memory bandwidth stress test on all hcu.

It evaluates bandwidth for Sys->Fb, Fb->Sys, and XHCL P2P transfers.

Usage:
  dcgmi diag pcie`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("=== Running PCIe Bandwidth Test ===")
		pcieResult, err := dcgm.PcieBandwidthTestResult()
		if err != nil {
			fmt.Print(formatDiagFailure("PCIe bandwidth test", err, pcieResult.LogFile))
			os.Exit(1)
		}
		fmt.Print(formatPCIESummary(pcieResult))
		fmt.Println("PCIe bandwidth test: Done.")
		if pcieResult.LogFile != "" {
			fmt.Printf("       Log: %s\n", pcieResult.LogFile)
		}
	},
}

// ==================== diag xhcl ====================
var xhclCmd = &cobra.Command{
	Use:   "xhcl",
	Short: "Run XHCL stress test",
	Long: `Run an XHCL P2P bandwidth stress test across all hcu.

Usage:
  dcgmi diag xhcl`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running XHCL stress test on all hcu...")
		xhclResults, err := dcgm.XHCLTestResult()
		if err != nil {
			fmt.Print(formatDiagFailure("XHCL stress test", err, dcgm.XHCLLogDir))
			os.Exit(1)
		}
		for _, r := range xhclResults {
			fmt.Printf("HCU %d->%d  Bandwidth: %.2f GB/s\n", r.SrcHCUId, r.DstHCUId, r.BandwidthGBs)
		}
		fmt.Println("XHCL stress test: Done.")
	},
}

var edppDeviceIDs []int

// ==================== diag edpp ====================
var edppCmd = &cobra.Command{
	Use:   "edpp",
	Short: "Run HCU EDPp stability test",
	Long: `Run the EDPp stress test on HCU devices.

By default, this test runs on all available HCUs. Use -d/--device to run it on
one or more specified HCUs for development validation.

Usage:
  dcgmi diag edpp
  dcgmi diag edpp -d 0
  dcgmi diag edpp -d 0 -d 1

Examples:
  dcgmi diag edpp
  dcgmi diag edpp -d 0`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(edppDeviceIDs) == 0 {
			fmt.Println("Running EDPp test for all HCUs...")
		} else {
			fmt.Printf("Running EDPp test for HCU(s): %v...\n", edppDeviceIDs)
		}
		result, err := dcgm.EDPpTestResult(edppDeviceIDs...)
		fmt.Print(formatEDPPResult(result))
		if err != nil {
			return fmt.Errorf("EDPp test failed: %w; logs: %s", err, result.LogDir)
		}
		fmt.Println("EDPp test: Done.")
		return nil
	},
}

func formatEDPPResult(result dcgm.EDPPResult) string {
	var output strings.Builder
	for _, hcuResult := range result.HCUEdppResults {
		for _, pattern := range hcuResult.PatternResults {
			fmt.Fprintf(&output, "HCU %d  Pattern=%s, ECC=%d, Mem=%d, Compute=%d\n", hcuResult.HCUId, pattern.PatternName, pattern.ECCCount, pattern.MemoryErrorCount, pattern.ComputeErrorCount)
		}
		for _, gemmPower := range hcuResult.GemmPowerResults {
			fmt.Fprintf(&output, "HCU %d  %s\n", hcuResult.HCUId, gemmPower.Summary())
			if gemmPower.Error != "" {
				fmt.Fprintf(&output, "       Error: %s\n", gemmPower.Error)
			}
		}
	}
	return output.String()
}

// ==================== diag gemm ====================
var gemmCmd = &cobra.Command{
	Use:   "gemm",
	Short: "Run GEMM target stress test",
	Long: `Run a GEMM target stress test on all hcu.

This test stresses compute performance using GEMM workloads.

Usage:
  dcgmi diag gemm`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running GEMM target stress test on all hcu...")
		gemmResult, err := dcgm.TargetStressTestResult()
		if err != nil {
			fmt.Print(formatDiagFailure("GEMM stress test", err, gemmResult.LogDir))
			os.Exit(1)
		}
		fmt.Print(formatGEMMSummary(gemmResult))
		fmt.Println("GEMM stress test: Done.")
		if gemmResult.LogDir != "" {
			fmt.Printf("       Logs: %s\n", gemmResult.LogDir)
		}
	},
}

// ==================== diag memtestCL ====================
var memtestCLCmd = &cobra.Command{
	Use:   "memtestCL [HCU IDs]",
	Short: "Run HCU memory stress test",
	Long: `Run a memory stress and integrity test on specified HCU devices.

If no device IDs are provided, the test will run on all available HCUs.
It detects memory errors under heavy load and produces detailed logs.

Usage:
  dcgmi diag memtestCL [HCU IDs]

Examples:
  dcgmi diag memtestCL
  dcgmi diag memtestCL 0
  dcgmi diag memtestCL 0 1 2`,
	Run: func(cmd *cobra.Command, args []string) {
		var dvIdList []int

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
				id, err := strconv.Atoi(arg)
				if err != nil {
					fmt.Printf("Invalid HCU ID '%s'\n", arg)
					os.Exit(1)
				}
				dvIdList = append(dvIdList, id)
			}
		}

		fmt.Printf("Running memtestCL stress test for HCU(s): %v\n", dvIdList)
		if err := dcgm.MemtestCL(dvIdList); err != nil {
			fmt.Print(formatDiagFailure("MemtestCL test", err, dcgm.MemtestCLLogDir))
			os.Exit(1)
		}
		//fmt.Println("Successfully completed memtestCL stress test.")
		fmt.Println("memtestCL stress test: Done.")
	},
}

// -------------------- 辅助打印函数 --------------------
func printDiagnosticResults(results dcgm.DiagResults) {
	fmt.Println("Successfully ran diagnostic for HCU.")
	fmt.Println("+---------------------------+------------------------------------------------+")
	fmt.Println("| Diagnostic                | Result                                         |")
	fmt.Println("+===========================+================================================+")

	// 输出 Software 部分
	if len(results.Software) > 0 {
		fmt.Println("| Software Diagnostics      |                                                |")
		fmt.Println("+---------------------------+------------------------------------------------+")
		for _, result := range results.Software {
			printTestResult(result)
		}
	}

	// 输出 PerHCU 部分
	if len(results.PerHCU) > 0 {
		fmt.Println("| Hardware Diagnostics      |                                                |")
		fmt.Println("+---------------------------+------------------------------------------------+")
		for _, hcu := range results.PerHCU {
			fmt.Printf("| HCU: %d                     |------------------------------------------------|\n", hcu.HCU)
			for _, result := range hcu.DiagResults {
				printTestResult(result)
			}
		}
	}

	fmt.Println("+---------------------------+------------------------------------------------+")
}

func printTestResult(result dcgm.DiagResult) {
	fmt.Printf("| %-25s | %-46s |\n", result.TestName, diagnosticStatusLabel(result.Status))
	if result.ErrorMessage != "" {
		fmt.Printf("| %-25s | %-46s |\n", "Error Message", result.ErrorMessage)
	}
	if result.TestOutput != "" {
		fmt.Printf("| %-25s | %-46s |\n", "Test Output", result.TestOutput)
	}
	if result.ErrorCode != 0 {
		fmt.Printf("| %-25s | %-46d |\n", "Error Code", result.ErrorCode)
	}
	fmt.Println("+------------------------------------------------+")
}

func diagnosticStatusLabel(status string) string {
	if status == dcgm.DiagResultWarn {
		return "Warning"
	}
	return status
}

// -------------------- 初始化 --------------------
func init() {
	diagCmd.AddCommand(runDiagCmd)
	diagCmd.AddCommand(bandwidthCmd)
	diagCmd.AddCommand(pcieCmd)
	rootCmd.AddCommand(diagCmd)
	diagCmd.AddCommand(xhclCmd)
	diagCmd.AddCommand(gemmCmd)
	diagCmd.AddCommand(memtestCLCmd)
	diagCmd.AddCommand(edppCmd)
	edppCmd.Flags().IntSliceVarP(&edppDeviceIDs, "device", "d", nil, "HCU ID to test. Repeat or comma-separate for multiple HCUs; omit to test all HCUs.")

	diagCmd.Flags().StringVarP(&groupId, "group", "g", "", "The group ID to query.")
	diagCmd.Flags().StringVarP(&infoFlags, "info", "i", "", "Specify which information to return.\n"+
		" b - memory bandwidth\n m - memtestCL stress")
}

func handleDiagGroup() {
	if infoFlags == "" {
		fmt.Println("Error: No info flag has been specified.")
		return
	}
	for _, c := range infoFlags {
		switch c {
		case 'b', 'm':
		default:
			fmt.Printf("Invalid input '%c'. Please include only valid tags.\n", c)
			return
		}
	}
	groupIdInt, err := strconv.Atoi(groupId)
	if err != nil {
		fmt.Println("Error: Invalid Group ID given.")
		return
	}
	hcuInGroup, _, err := dcgm.GetHCUListFromGroup(groupIdInt)
	if err != nil {
		fmt.Printf("Error getting group info: %v\n", err)
		return
	}
	if len(hcuInGroup) == 0 {
		fmt.Printf("Failed to query group: no entity found for group %v\n", groupId)
		return
	}
	if strings.Contains(infoFlags, "b") {
		fmt.Printf("Running memory bandwidth test for HCU(s): %v\n", hcuInGroup)
		bwResults, err := dcgm.BandwidthTestResult(hcuInGroup)
		if err != nil {
			fmt.Print(formatDiagFailure("Bandwidth test", err, dcgm.DiagLogDirBandwidth))
			return
		}
		fmt.Println("===== 带宽测试结果 =====")
		for _, hcu := range hcuInGroup {
			fmt.Printf("HCU%d: %.2f GB/s\n", hcu, bwResults[hcu])
		}
		fmt.Println("Successfully completed memory bandwidth test.")
	}
	if strings.Contains(infoFlags, "m") {
		fmt.Printf("Running memtestCL stress test for HCU(s): %v\n", hcuInGroup)
		if err := dcgm.MemtestCL(hcuInGroup); err != nil {
			fmt.Print(formatDiagFailure("MemtestCL test", err, dcgm.MemtestCLLogDir))
			os.Exit(1)
		}
		fmt.Println("Successfully completed memtestCL stress test ✅")
	}
}
