// Copyright (c) 2026 Hygon Information Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0
package dcgm

import (
	"errors"
	"testing"
)

func TestParseDevicePerformanceLevel(t *testing.T) {
	level, err := parseDevicePerformanceLevel("stable_peak")
	if err != nil {
		t.Fatal(err)
	}
	if level != RSMI_DEV_PERF_LEVEL_STABLE_PEAK {
		t.Fatalf("level = %v, want stable peak", level)
	}

	if _, err := parseDevicePerformanceLevel("turbo"); err == nil {
		t.Fatal("unsupported performance level must fail")
	}
}

func TestDeviceConfigPathValidation(t *testing.T) {
	path, err := deviceConfigPath("production_01")
	if err != nil || path != "/opt/dcgm/device-configs/production_01.json" {
		t.Fatalf("path = %q, err = %v", path, err)
	}

	_, err = deviceConfigPath("../escape")
	if !errors.Is(err, ErrInvalidDeviceConfigName) {
		t.Fatalf("err = %v, want ErrInvalidDeviceConfigName", err)
	}
}
