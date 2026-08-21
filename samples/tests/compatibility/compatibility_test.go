// Copyright (c) 2026 Hygon Information Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0

package compatibility_test

import (
	"strings"
	"testing"

	"github.com/HYGON-AI/hcu-dcgm/v3/pkg/dcgm"
)

func TestOfficialCompatibilityMatrix(t *testing.T) {
	testCases := []struct {
		name          string
		cardModel     string
		driverVersion string
		dtkVersion    string
		wantStatus    dcgm.CompatibilityStatus
		wantDriver    string
		wantDTK       string
		wantError     string
	}{
		{
			name:          "BW1100 accepts deployed driver and DTK",
			cardModel:     "BW1100",
			driverVersion: "6.3.31-V1.5.1",
			dtkVersion:    "26.04",
			wantStatus:    dcgm.CompatibilityPass,
		},
		{
			name:          "BW1000 accepts driver suffix and DTK patch release",
			cardModel:     "BW1000",
			driverVersion: "rock-6.3.27-V1.2.5",
			dtkVersion:    "DTK-25.04.4",
			wantStatus:    dcgm.CompatibilityPass,
		},
		{
			name:          "BW1000 warns about a different DTK minor release",
			cardModel:     "BW1000",
			driverVersion: "rock-6.3.27-V1.2.5",
			dtkVersion:    "25.05.0",
			wantStatus:    dcgm.CompatibilityWarn,
			wantDriver:    ">= 6.3.8",
			wantDTK:       "25.04.*",
		},
		{
			name:          "Z100 warns when driver is below first tier",
			cardModel:     "Z100",
			driverVersion: "5.6.24",
			dtkVersion:    "21.04",
			wantStatus:    dcgm.CompatibilityWarn,
			wantDriver:    ">= 5.6.25",
			wantDTK:       "21.04.*/21.10.*/22.04.*",
		},
		{
			name:          "unknown model is rejected",
			cardModel:     "UNKNOWN",
			driverVersion: "6.3.30",
			dtkVersion:    "26.04",
			wantError:     "不支持的卡型号: UNKNOWN",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := dcgm.Compatible(testCase.cardModel, testCase.driverVersion, testCase.dtkVersion)
			if testCase.wantError != "" {
				if err == nil || !strings.Contains(err.Error(), testCase.wantError) {
					t.Fatalf("Compatible() error = %v, want containing %q", err, testCase.wantError)
				}
				return
			}
			if err != nil {
				t.Fatalf("Compatible() returned unexpected error: %v", err)
			}
			if result.Status != testCase.wantStatus {
				t.Fatalf("Compatible() status = %q, want %q", result.Status, testCase.wantStatus)
			}
			if result.RecommendedDriver != testCase.wantDriver {
				t.Fatalf("Compatible() recommended driver = %q, want %q", result.RecommendedDriver, testCase.wantDriver)
			}
			if result.RecommendedDTK != testCase.wantDTK {
				t.Fatalf("Compatible() recommended DTK = %q, want %q", result.RecommendedDTK, testCase.wantDTK)
			}
		})
	}
}

func TestOfficialCompatibilityMatrixCoversSupportedModels(t *testing.T) {
	models := []struct {
		cardModel     string
		driverVersion string
		dtkVersion    string
	}{
		{"Z100", "5.6.25", "21.04"},
		{"Z100L", "5.6.25", "21.04"},
		{"K100", "5.16.29", "23.04"},
		{"K100_AI", "6.2.26", "24.04"},
		{"BW10", "6.3.30", "26.04"},
		{"BW100", "6.3.30", "26.04"},
		{"BW150", "6.3.30", "26.04"},
		{"BW1000", "6.3.8", "25.04"},
		{"BW1100", "6.3.30", "26.04"},
	}

	for _, model := range models {
		t.Run(model.cardModel, func(t *testing.T) {
			result, err := dcgm.Compatible(model.cardModel, model.driverVersion, model.dtkVersion)
			if err != nil {
				t.Fatalf("Compatible() returned unexpected error: %v", err)
			}
			if result.Status != dcgm.CompatibilityPass {
				t.Fatalf("Compatible() status = %q, want pass", result.Status)
			}
		})
	}
}

func TestCompatibleRetainsBW200LegacyMinimums(t *testing.T) {
	result, err := dcgm.Compatible("BW200", "6.3.0", "25.00.00")
	if err != nil || result.Status != dcgm.CompatibilityPass {
		t.Fatalf("Compatible() = (%+v, %v), want pass", result, err)
	}

	result, err = dcgm.Compatible("BW200", "6.2.99", "24.04")
	if err != nil || result.Status != dcgm.CompatibilityWarn {
		t.Fatalf("Compatible() = (%+v, %v), want warning", result, err)
	}
	if result.RecommendedDriver != ">= 6.3.0" || result.RecommendedDTK != ">= 25.00.00" {
		t.Fatalf("Compatibility warning recommendations = (%q, %q)", result.RecommendedDriver, result.RecommendedDTK)
	}
}
