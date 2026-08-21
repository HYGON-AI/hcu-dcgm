/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package dcgm

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type CompatibilityStatus string

const (
	CompatibilityPass CompatibilityStatus = "pass"
	CompatibilityWarn CompatibilityStatus = "warn"
)

type CompatibilityResult struct {
	Status            CompatibilityStatus `json:"status"`
	CardModel         string              `json:"cardModel"`
	CurrentDriver     string              `json:"currentDriver"`
	CurrentDTK        string              `json:"currentDTK"`
	RecommendedDriver string              `json:"recommendedDriver"`
	RecommendedDTK    string              `json:"recommendedDTK"`
}

type compatibilityTier struct {
	minimumDriver string
	supportedDTK  []string
}

type legacyCompatibilityRule struct {
	minimumDriver string
	minimumDTK    string
}

var officialCompatibilityMatrix = map[string][]compatibilityTier{
	"Z100": {
		{minimumDriver: "5.6.25", supportedDTK: []string{"21.04", "21.10", "22.04"}},
		{minimumDriver: "5.11.40", supportedDTK: []string{"22.04", "22.10", "23.04"}},
		{minimumDriver: "5.16.18", supportedDTK: []string{"22.10", "23.04"}},
		{minimumDriver: "5.16.29", supportedDTK: []string{"23.04", "23.10"}},
		{minimumDriver: "6.2.26", supportedDTK: []string{"24.04", "25.04"}},
		{minimumDriver: "6.3.8", supportedDTK: []string{"25.04"}},
		{minimumDriver: "6.3.30", supportedDTK: []string{"26.04"}},
	},
	"Z100L": {
		{minimumDriver: "5.6.25", supportedDTK: []string{"21.04", "21.10", "22.04"}},
		{minimumDriver: "5.11.40", supportedDTK: []string{"22.04", "22.10", "23.04"}},
		{minimumDriver: "5.16.18", supportedDTK: []string{"22.10", "23.04"}},
		{minimumDriver: "5.16.29", supportedDTK: []string{"23.04", "23.10"}},
		{minimumDriver: "6.2.26", supportedDTK: []string{"24.04", "25.04"}},
		{minimumDriver: "6.3.8", supportedDTK: []string{"25.04"}},
		{minimumDriver: "6.3.30", supportedDTK: []string{"26.04"}},
	},
	"K100": {
		{minimumDriver: "5.16.29", supportedDTK: []string{"23.04", "23.10"}},
		{minimumDriver: "6.2.26", supportedDTK: []string{"24.04", "25.04"}},
		{minimumDriver: "6.3.8", supportedDTK: []string{"25.04"}},
		{minimumDriver: "6.3.30", supportedDTK: []string{"26.04"}},
	},
	"K100_AI": {
		{minimumDriver: "6.2.26", supportedDTK: []string{"24.04", "25.04"}},
		{minimumDriver: "6.3.8", supportedDTK: []string{"25.04"}},
		{minimumDriver: "6.3.30", supportedDTK: []string{"26.04"}},
	},
	"BW10": {
		{minimumDriver: "6.3.30", supportedDTK: []string{"26.04"}},
	},
	"BW100": {
		{minimumDriver: "6.3.30", supportedDTK: []string{"26.04"}},
	},
	"BW150": {
		{minimumDriver: "6.3.30", supportedDTK: []string{"26.04"}},
	},
	"BW1000": {
		{minimumDriver: "6.3.8", supportedDTK: []string{"25.04"}},
		{minimumDriver: "6.3.30", supportedDTK: []string{"26.04"}},
	},
	"BW1100": {
		{minimumDriver: "6.3.30", supportedDTK: []string{"26.04"}},
	},
}

var legacyCompatibilityMatrix = map[string]legacyCompatibilityRule{
	"BW200": {minimumDriver: "6.3.0", minimumDTK: "25.00.00"},
}

func compatible(cardModel, driverVersion, dtkVersion string) (CompatibilityResult, error) {
	result := CompatibilityResult{
		CardModel:     cardModel,
		CurrentDriver: driverVersion,
		CurrentDTK:    dtkVersion,
	}
	normalizedDriver := normalizedVersion(driverVersion)
	normalizedDTK := normalizedVersion(dtkVersion)
	if normalizedDriver == "" || normalizedDTK == "" {
		return result, fmt.Errorf("无法解析卡型号 %s 的 Driver 或 DTK 版本", cardModel)
	}

	if tiers, exists := officialCompatibilityMatrix[cardModel]; exists {
		return validateOfficialCompatibility(result, normalizedDriver, normalizedDTK, tiers)
	}
	if rule, exists := legacyCompatibilityMatrix[cardModel]; exists {
		return validateLegacyCompatibility(result, normalizedDriver, normalizedDTK, rule)
	}
	return result, fmt.Errorf("不支持的卡型号: %s", cardModel)
}

func validateOfficialCompatibility(result CompatibilityResult, driverVersion, dtkVersion string, tiers []compatibilityTier) (CompatibilityResult, error) {
	matchedTier, matched := highestMatchingTier(driverVersion, tiers)
	if !matched {
		return compatibilityWarning(result, tiers[0].minimumDriver, tiers[0].supportedDTK), nil
	}
	if supportsDTKRelease(matchedTier.supportedDTK, dtkVersion) {
		result.Status = CompatibilityPass
		return result, nil
	}
	return compatibilityWarning(result, matchedTier.minimumDriver, matchedTier.supportedDTK), nil
}

func validateLegacyCompatibility(result CompatibilityResult, driverVersion, dtkVersion string, rule legacyCompatibilityRule) (CompatibilityResult, error) {
	if isVersionAtLeast(driverVersion, rule.minimumDriver) && isVersionAtLeast(dtkVersion, rule.minimumDTK) {
		result.Status = CompatibilityPass
		return result, nil
	}
	result.Status = CompatibilityWarn
	result.RecommendedDriver = ">= " + rule.minimumDriver
	result.RecommendedDTK = ">= " + rule.minimumDTK
	return result, nil
}

func compatibilityWarning(result CompatibilityResult, minimumDriver string, supportedDTK []string) CompatibilityResult {
	result.Status = CompatibilityWarn
	result.RecommendedDriver = ">= " + minimumDriver
	result.RecommendedDTK = formatSupportedDTKReleases(supportedDTK)
	return result
}

func formatSupportedDTKReleases(releases []string) string {
	formatted := make([]string, 0, len(releases))
	for _, release := range releases {
		formatted = append(formatted, normalizedVersion(release)+".*")
	}
	return strings.Join(formatted, "/")
}

func highestMatchingTier(driverVersion string, tiers []compatibilityTier) (compatibilityTier, bool) {
	var selected compatibilityTier
	matched := false
	for _, tier := range tiers {
		if isVersionAtLeast(driverVersion, tier.minimumDriver) {
			selected = tier
			matched = true
		}
	}
	return selected, matched
}

func normalizedVersion(version string) string {
	return regexp.MustCompile(`[0-9.]+`).FindString(strings.TrimSpace(version))
}

func isVersionAtLeast(version, minimum string) bool {
	versionParts := strings.Split(normalizedVersion(version), ".")
	minimumParts := strings.Split(normalizedVersion(minimum), ".")
	partCount := maximumVersionParts(len(versionParts), len(minimumParts))

	for index := range partCount {
		versionPart := versionPartAt(versionParts, index)
		minimumPart := versionPartAt(minimumParts, index)
		if versionPart != minimumPart {
			return versionPart > minimumPart
		}
	}
	return true
}

func maximumVersionParts(left, right int) int {
	if left > right {
		return left
	}
	return right
}

func versionPartAt(parts []string, index int) int {
	if index >= len(parts) || parts[index] == "" {
		return 0
	}
	part, _ := strconv.Atoi(parts[index])
	return part
}

func supportsDTKRelease(supportedReleases []string, version string) bool {
	version = normalizedVersion(version)
	for _, release := range supportedReleases {
		release = normalizedVersion(release)
		if version == release || strings.HasPrefix(version, release+".") {
			return true
		}
	}
	return false
}
