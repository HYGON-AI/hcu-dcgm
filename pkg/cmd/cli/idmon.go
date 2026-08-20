/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package cli

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/HYGON-AI/hcu-dcgm/v3/pkg/dcgm"
)

var (
	idmonDevices  string
	idmonSelect   string
	idmonCount    int
	idmonDelay    int
	idmonShowTime bool
	idmonFile     string
	idmonFormat   string
)

type idmonMetric struct {
	Key    rune
	Header string
	Read   func(int) (string, error)
}

var idmonMetrics = []idmonMetric{
	{Key: 'c', Header: "GFX_CLOCK_MHZ", Read: func(dvInd int) (string, error) {
		value, err := dcgm.HCUClk(dvInd)
		return strconv.FormatFloat(value, 'f', 2, 64), err
	}},
	{Key: 'f', Header: "FAN_PERCENT", Read: func(dvInd int) (string, error) {
		_, value, err := dcgm.FanSpeedInfo(dvInd)
		return strconv.FormatFloat(value, 'f', 2, 64), err
	}},
	{Key: 'p', Header: "POWER_W", Read: func(dvInd int) (string, error) {
		value, err := dcgm.Power(dvInd)
		return strconv.FormatInt(value, 10), err
	}},
	{Key: 't', Header: "TEMP_C", Read: func(dvInd int) (string, error) {
		value, err := dcgm.Temperature(dvInd)
		return strconv.FormatFloat(value, 'f', 2, 64), err
	}},
	{Key: 'u', Header: "UTIL_PERCENT", Read: func(dvInd int) (string, error) {
		value, err := dcgm.DevBusyPercent(dvInd)
		return strconv.FormatFloat(value, 'f', 2, 64), err
	}},
}

var idmonCmd = &cobra.Command{
	Use:   "idmon",
	Short: "Periodically sample selected HCU metrics.",
	Example: `  dcgmi idmon -i 0 -s cptu -c 10 -d 1
  dcgmi idmon -i 0,1 -s cfptu -t csv -f metrics.csv -o`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if idmonCount == 0 || idmonCount < -1 {
			return fmt.Errorf("count must be -1 or a positive integer")
		}
		if idmonDelay < 1 {
			return fmt.Errorf("delay must be at least one second")
		}
		if idmonFormat != "text" && idmonFormat != "csv" {
			return fmt.Errorf("type must be text or csv")
		}

		devices, err := parseIDMonDevices(idmonDevices)
		if err != nil {
			return err
		}
		metrics, err := selectIDMonMetrics(idmonSelect)
		if err != nil {
			return err
		}
		return runIDMon(devices, metrics)
	},
}

func init() {
	rootCmd.AddCommand(idmonCmd)
	idmonCmd.Flags().StringVarP(&idmonDevices, "index", "i", "", "Comma-separated HCU indexes; empty selects all devices.")
	idmonCmd.Flags().StringVarP(&idmonSelect, "select", "s", "cfptu", "Metrics: c=clock, f=fan, p=power, t=temperature, u=utilization.")
	idmonCmd.Flags().IntVarP(&idmonCount, "count", "c", -1, "Number of samples; -1 runs until interrupted.")
	idmonCmd.Flags().IntVarP(&idmonDelay, "delay", "d", 1, "Seconds between samples.")
	idmonCmd.Flags().BoolVarP(&idmonShowTime, "showtime", "o", false, "Include an RFC3339 timestamp column.")
	idmonCmd.Flags().StringVarP(&idmonFile, "file", "f", "", "Write output to a file instead of stdout.")
	idmonCmd.Flags().StringVarP(&idmonFormat, "type", "t", "text", "Output format: text or csv.")
}

func parseIDMonDevices(value string) ([]int, error) {
	if strings.TrimSpace(value) == "" {
		count, err := dcgm.NumMonitorDevices()
		if err != nil {
			return nil, err
		}
		devices := make([]int, count)
		for i := range devices {
			devices[i] = i
		}
		return devices, nil
	}

	seen := make(map[int]struct{})
	devices := make([]int, 0)
	for _, item := range strings.Split(value, ",") {
		device, err := strconv.Atoi(strings.TrimSpace(item))
		if err != nil || device < 0 {
			return nil, fmt.Errorf("invalid HCU index %q", item)
		}
		if _, exists := seen[device]; exists {
			continue
		}
		seen[device] = struct{}{}
		devices = append(devices, device)
	}
	if len(devices) == 0 {
		return nil, fmt.Errorf("at least one HCU index is required")
	}
	return devices, nil
}

func selectIDMonMetrics(value string) ([]idmonMetric, error) {
	requested := make(map[rune]struct{})
	for _, key := range strings.ToLower(strings.TrimSpace(value)) {
		requested[key] = struct{}{}
	}
	if len(requested) == 0 {
		return nil, fmt.Errorf("at least one metric must be selected")
	}

	metrics := make([]idmonMetric, 0, len(requested))
	for _, metric := range idmonMetrics {
		if _, ok := requested[metric.Key]; ok {
			metrics = append(metrics, metric)
			delete(requested, metric.Key)
		}
	}
	if len(requested) != 0 {
		unknown := make([]string, 0, len(requested))
		for key := range requested {
			unknown = append(unknown, string(key))
		}
		return nil, fmt.Errorf("unsupported metric selection %q", strings.Join(unknown, ""))
	}
	return metrics, nil
}

func runIDMon(devices []int, metrics []idmonMetric) error {
	writer, closeWriter, err := idmonOutputWriter(idmonFile)
	if err != nil {
		return err
	}
	defer closeWriter()

	header := []string{"HCU"}
	if idmonShowTime {
		header = append([]string{"TIMESTAMP"}, header...)
	}
	for _, metric := range metrics {
		header = append(header, metric.Header)
	}

	writeRows, flush, err := newIDMonRowWriter(writer, header, idmonFormat)
	if err != nil {
		return err
	}
	defer flush()

	writeSample := func() error {
		timestamp := time.Now().Format(time.RFC3339)
		rows := make([][]string, 0, len(devices))
		for _, device := range devices {
			row := []string{strconv.Itoa(device)}
			if idmonShowTime {
				row = append([]string{timestamp}, row...)
			}
			for _, metric := range metrics {
				value, err := metric.Read(device)
				if err != nil {
					value = "N/A"
					fmt.Fprintf(os.Stderr, "idmon: HCU %d %s: %v\n", device, metric.Header, err)
				}
				row = append(row, value)
			}
			rows = append(rows, row)
		}
		return writeRows(rows)
	}

	executed := 0
	if err := writeSample(); err != nil {
		return err
	}
	executed++
	if idmonCount != -1 && executed >= idmonCount {
		return nil
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(stop)
	ticker := time.NewTicker(time.Duration(idmonDelay) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-stop:
			return nil
		case <-ticker.C:
			if err := writeSample(); err != nil {
				return err
			}
			executed++
			if idmonCount != -1 && executed >= idmonCount {
				return nil
			}
		}
	}
}

func idmonOutputWriter(path string) (io.Writer, func(), error) {
	if path == "" {
		return os.Stdout, func() {}, nil
	}
	file, err := os.Create(path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create output file: %w", err)
	}
	return file, func() { _ = file.Close() }, nil
}

func newIDMonRowWriter(writer io.Writer, header []string, format string) (func([][]string) error, func(), error) {
	if format == "csv" {
		csvWriter := csv.NewWriter(writer)
		if err := csvWriter.Write(header); err != nil {
			return nil, nil, err
		}
		csvWriter.Flush()
		write := func(rows [][]string) error {
			csvWriter.WriteAll(rows)
			csvWriter.Flush()
			return csvWriter.Error()
		}
		return write, csvWriter.Flush, nil
	}

	tabWriter := tabwriter.NewWriter(writer, 0, 4, 2, ' ', 0)
	if _, err := fmt.Fprintln(tabWriter, strings.Join(header, "\t")); err != nil {
		return nil, nil, err
	}
	if err := tabWriter.Flush(); err != nil {
		return nil, nil, err
	}
	write := func(rows [][]string) error {
		for _, row := range rows {
			if _, err := fmt.Fprintln(tabWriter, strings.Join(row, "\t")); err != nil {
				return err
			}
		}
		return tabWriter.Flush()
	}
	return write, func() { _ = tabWriter.Flush() }, nil
}
