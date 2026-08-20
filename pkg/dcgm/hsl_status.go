/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package dcgm

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const hslManagerLogDir = "/var/log/hfm"
const hslManagerLogName = "hlinkmanager.log"

// resolveHSLLogPath 优先使用固定名称日志，若不存在则取目录中最新的 hfm.*.log。
func resolveHSLLogPath() (string, error) {
	fixed := filepath.Join(hslManagerLogDir, hslManagerLogName)
	if _, err := os.Stat(fixed); err == nil {
		return fixed, nil
	}

	entries, err := os.ReadDir(hslManagerLogDir)
	if err != nil {
		return "", fmt.Errorf("failed to list %s: %w", hslManagerLogDir, err)
	}

	var latest string
	for _, e := range entries {
		name := e.Name()
		if strings.HasPrefix(name, "hfm.") && strings.HasSuffix(name, ".log") {
			if name > latest {
				latest = name
			}
		}
	}
	if latest == "" {
		return "", fmt.Errorf("no HSL manager log found in %s", hslManagerLogDir)
	}
	return filepath.Join(hslManagerLogDir, latest), nil
}

var hslLogPrefixPattern = regexp.MustCompile(`^\[([^\]]+)\]\s*(.*)$`)

type HSLLogEntry struct {
	Header  string `json:"header,omitempty"`
	Level   string `json:"level"`
	Message string `json:"message"`
	Raw     string `json:"raw"`
}

type HSLLogStatus struct {
	DvInd     int           `json:"dvInd"`
	LogPath   string        `json:"logPath"`
	Status    string        `json:"status"`
	LastError string        `json:"lastError,omitempty"`
	Entries   []HSLLogEntry `json:"entries"`
}

// GetHSLLogStatus 返回日志中指定设备最近的链路状态记录。
// 日志格式由 hlinkmanager 决定，因此仅解析通用日志头，原始行始终保留。
func GetHSLLogStatus(dvInd, maxEntries int) (HSLLogStatus, error) {
	if dvInd < 0 {
		return HSLLogStatus{}, fmt.Errorf("invalid device index %d", dvInd)
	}
	if maxEntries <= 0 {
		maxEntries = 100
	}

	logPath, err := resolveHSLLogPath()
	if err != nil {
		return HSLLogStatus{}, err
	}

	file, err := os.Open(logPath)
	if err != nil {
		return HSLLogStatus{}, fmt.Errorf("failed to open HSL manager log %s: %w", logPath, err)
	}
	defer file.Close()

	status := HSLLogStatus{
		DvInd:   dvInd,
		LogPath: logPath,
		Status:  "UNKNOWN",
		Entries: make([]HSLLogEntry, 0, maxEntries),
	}
	devicePattern := regexp.MustCompile(`(?i)\b(?:device|dcu|hcu|gpu)\s*[:=#-]?\s*` + strconv.Itoa(dvInd) + `\b`)

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || !devicePattern.MatchString(line) {
			continue
		}
		entry := parseHSLLogEntry(line)
		status.Entries = append(status.Entries, entry)
		if len(status.Entries) > maxEntries {
			status.Entries = status.Entries[len(status.Entries)-maxEntries:]
		}
		if entry.Level == "ERROR" {
			status.Status = "ERROR"
			status.LastError = entry.Message
		} else if status.Status == "UNKNOWN" {
			status.Status = "OK"
		}
	}
	if err := scanner.Err(); err != nil {
		return HSLLogStatus{}, fmt.Errorf("failed to scan HSL manager log: %w", err)
	}
	return status, nil
}

func parseHSLLogEntry(line string) HSLLogEntry {
	entry := HSLLogEntry{Raw: line, Message: line, Level: "INFO"}
	if matches := hslLogPrefixPattern.FindStringSubmatch(line); len(matches) == 3 {
		entry.Header = strings.TrimSpace(matches[1])
		entry.Message = strings.TrimSpace(matches[2])
	}

	lower := strings.ToLower(entry.Message)
	switch {
	case strings.Contains(lower, "error"), strings.Contains(lower, "failed"), strings.Contains(lower, "link down"):
		entry.Level = "ERROR"
	case strings.Contains(lower, "warn"), strings.Contains(lower, "waiting"), strings.Contains(lower, "not ready"):
		entry.Level = "WARN"
	}
	return entry
}
