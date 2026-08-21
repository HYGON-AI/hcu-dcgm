/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package dcgm

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const deviceConfigDirectory = "/opt/dcgm/device-configs"

var (
	ErrInvalidDeviceConfigName = errors.New("invalid device config name")
	deviceConfigNamePattern    = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)
)

type DeviceConfigSnapshot struct {
	DvInd              int     `json:"dvInd"`
	PowerCapUW         *uint64 `json:"powerCapUW,omitempty"`
	PerformanceLevel   *string `json:"performanceLevel,omitempty"`
	OverdriveLevel     *int    `json:"overdriveLevel,omitempty"`
	LowPowerMode       *string `json:"lowPowerMode,omitempty"`
	AutosuspendDelayMS *int    `json:"autosuspendDelayMs,omitempty"`
}

type DeviceConfigIssue struct {
	DvInd int    `json:"dvInd"`
	Field string `json:"field"`
	Error string `json:"error"`
}

type DeviceConfigDocument struct {
	Version int                    `json:"version"`
	Name    string                 `json:"name"`
	SavedAt time.Time              `json:"savedAt"`
	Devices []DeviceConfigSnapshot `json:"devices"`
	Issues  []DeviceConfigIssue    `json:"issues,omitempty"`
}

type DeviceConfigApplyResult struct {
	Name    string              `json:"name"`
	Devices int                 `json:"devices"`
	Issues  []DeviceConfigIssue `json:"issues,omitempty"`
}

// SaveDeviceConfig 捕获设备可对称恢复的配置，并原子写入服务固定目录。
func SaveDeviceConfig(name string, deviceIDs []int) (DeviceConfigDocument, error) {
	path, err := deviceConfigPath(name)
	if err != nil {
		return DeviceConfigDocument{}, err
	}
	if len(deviceIDs) == 0 {
		count, err := NumMonitorDevices()
		if err != nil {
			return DeviceConfigDocument{}, err
		}
		deviceIDs = make([]int, count)
		for i := range deviceIDs {
			deviceIDs[i] = i
		}
	}

	document := DeviceConfigDocument{
		Version: 1,
		Name:    name,
		SavedAt: time.Now().UTC(),
		Devices: make([]DeviceConfigSnapshot, 0, len(deviceIDs)),
		Issues:  make([]DeviceConfigIssue, 0),
	}
	for _, deviceID := range deviceIDs {
		snapshot, issues := captureDeviceConfig(deviceID)
		document.Devices = append(document.Devices, snapshot)
		document.Issues = append(document.Issues, issues...)
	}

	data, err := json.MarshalIndent(document, "", "  ")
	if err != nil {
		return DeviceConfigDocument{}, fmt.Errorf("failed to encode device config: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
		return DeviceConfigDocument{}, fmt.Errorf("failed to create device config directory: %w", err)
	}
	temporaryPath := path + ".tmp"
	if err := os.WriteFile(temporaryPath, data, 0600); err != nil {
		return DeviceConfigDocument{}, fmt.Errorf("failed to write device config: %w", err)
	}
	if err := os.Rename(temporaryPath, path); err != nil {
		_ = os.Remove(temporaryPath)
		return DeviceConfigDocument{}, fmt.Errorf("failed to replace device config: %w", err)
	}
	return document, nil
}

// LoadDeviceConfig 读取命名快照并逐设备恢复；单字段失败不会阻断其它字段。
func LoadDeviceConfig(name string) (DeviceConfigApplyResult, error) {
	path, err := deviceConfigPath(name)
	if err != nil {
		return DeviceConfigApplyResult{}, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return DeviceConfigApplyResult{}, fmt.Errorf("failed to read device config %s: %w", name, err)
	}
	var document DeviceConfigDocument
	if err := json.Unmarshal(data, &document); err != nil {
		return DeviceConfigApplyResult{}, fmt.Errorf("failed to decode device config %s: %w", name, err)
	}
	if document.Version != 1 {
		return DeviceConfigApplyResult{}, fmt.Errorf("unsupported device config version %d", document.Version)
	}

	result := DeviceConfigApplyResult{Name: name, Devices: len(document.Devices)}
	for _, snapshot := range document.Devices {
		result.Issues = append(result.Issues, applyDeviceConfig(snapshot)...)
	}
	return result, nil
}

func captureDeviceConfig(dvInd int) (DeviceConfigSnapshot, []DeviceConfigIssue) {
	snapshot := DeviceConfigSnapshot{DvInd: dvInd}
	issues := make([]DeviceConfigIssue, 0)

	if capUW, err := rsmiDevPowerCapGet(dvInd, 0); err == nil {
		value := uint64(capUW)
		snapshot.PowerCapUW = &value
	} else {
		issues = appendConfigIssue(issues, dvInd, "powerCapUW", err)
	}
	if level, err := PerfLevel(dvInd); err == nil {
		snapshot.PerformanceLevel = &level
	} else {
		issues = appendConfigIssue(issues, dvInd, "performanceLevel", err)
	}
	if level, err := DevOverdriveLevelGet(dvInd); err == nil {
		snapshot.OverdriveLevel = &level
	} else {
		issues = appendConfigIssue(issues, dvInd, "overdriveLevel", err)
	}

	cardIdx := dvIndexToCardIndex(dvInd)
	if mode, err := readSysfsFile(fmt.Sprintf("/sys/class/drm/card%d/device/power/control", cardIdx)); err == nil {
		snapshot.LowPowerMode = &mode
	} else {
		issues = appendConfigIssue(issues, dvInd, "lowPowerMode", err)
	}
	if delay, err := readSysfsFile(fmt.Sprintf("/sys/class/drm/card%d/device/power/autosuspend_delay_ms", cardIdx)); err == nil {
		if value, parseErr := strconv.Atoi(strings.TrimSpace(delay)); parseErr == nil {
			snapshot.AutosuspendDelayMS = &value
		} else {
			issues = appendConfigIssue(issues, dvInd, "autosuspendDelayMs", parseErr)
		}
	} else {
		issues = appendConfigIssue(issues, dvInd, "autosuspendDelayMs", err)
	}
	return snapshot, issues
}

func applyDeviceConfig(snapshot DeviceConfigSnapshot) []DeviceConfigIssue {
	issues := make([]DeviceConfigIssue, 0)
	apply := func(field string, err error) {
		if err != nil {
			issues = appendConfigIssue(issues, snapshot.DvInd, field, err)
		}
	}
	if snapshot.PowerCapUW != nil {
		apply("powerCapUW", DevPowerCapSet(snapshot.DvInd, *snapshot.PowerCapUW))
	}
	if snapshot.OverdriveLevel != nil {
		apply("overdriveLevel", DevOverdriveLevelSet(snapshot.DvInd, *snapshot.OverdriveLevel))
	}
	if snapshot.AutosuspendDelayMS != nil {
		apply("autosuspendDelayMs", SetLowPowerDelay(snapshot.DvInd, *snapshot.AutosuspendDelayMS))
	}
	if snapshot.LowPowerMode != nil {
		switch strings.ToLower(*snapshot.LowPowerMode) {
		case "auto":
			apply("lowPowerMode", SetLowPowerMode(snapshot.DvInd, true))
		case "on":
			apply("lowPowerMode", SetLowPowerMode(snapshot.DvInd, false))
		default:
			apply("lowPowerMode", fmt.Errorf("unsupported power control mode %q", *snapshot.LowPowerMode))
		}
	}
	if snapshot.PerformanceLevel != nil {
		level, err := parseDevicePerformanceLevel(*snapshot.PerformanceLevel)
		if err == nil {
			err = DevPerfLevelSetV1(snapshot.DvInd, level)
		}
		apply("performanceLevel", err)
	}
	return issues
}

func parseDevicePerformanceLevel(value string) (DevPerfLevel, error) {
	levels := map[string]DevPerfLevel{
		"AUTO":            RSMI_DEV_PERF_LEVEL_AUTO,
		"LOW":             RSMI_DEV_PERF_LEVEL_LOW,
		"HIGH":            RSMI_DEV_PERF_LEVEL_HIGH,
		"MANUAL":          RSMI_DEV_PERF_LEVEL_MANUAL,
		"STABLE_STD":      RSMI_DEV_PERF_LEVEL_STABLE_STD,
		"STABLE_PEAK":     RSMI_DEV_PERF_LEVEL_STABLE_PEAK,
		"STABLE_MIN_MCLK": RSMI_DEV_PERF_LEVEL_STABLE_MIN_MCLK,
		"STABLE_MIN_SCLK": RSMI_DEV_PERF_LEVEL_STABLE_MIN_SCLK,
		"DETERMINISM":     RSMI_DEV_PERF_LEVEL_DETERMINISM,
	}
	level, ok := levels[strings.ToUpper(strings.TrimSpace(value))]
	if !ok {
		return RSMI_DEV_PERF_LEVEL_UNKNOWN, fmt.Errorf("unsupported performance level %q", value)
	}
	return level, nil
}

func appendConfigIssue(issues []DeviceConfigIssue, dvInd int, field string, err error) []DeviceConfigIssue {
	return append(issues, DeviceConfigIssue{DvInd: dvInd, Field: field, Error: err.Error()})
}

func deviceConfigPath(name string) (string, error) {
	if !deviceConfigNamePattern.MatchString(name) {
		return "", fmt.Errorf("%w: use letters, numbers, underscore or hyphen", ErrInvalidDeviceConfigName)
	}
	return filepath.Join(deviceConfigDirectory, name+".json"), nil
}
