/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package dcgm

import "strings"

type MetricGroup struct {
	Major    int
	Minor    int
	FieldIds []int
}

func listMetricGroups() []MetricGroup {
	return profilingMetrics
}

func toSet(items []int) map[int]struct{} {
	set := make(map[int]struct{}, len(items))
	for _, v := range items {
		set[v] = struct{}{}
	}
	return set
}

func filterMetricsByName(hcuName string, groups []MetricGroup) []MetricGroup {
	unsupported, ok := unsupportedFieldsByName[hcuName]
	if !ok || len(unsupported) == 0 {
		return groups
	}

	unsupportedSet := toSet(unsupported)
	result := make([]MetricGroup, 0, len(groups))

	for _, g := range groups {
		filtered := make([]int, 0, len(g.FieldIds))
		for _, m := range g.FieldIds {
			if _, found := unsupportedSet[m]; !found {
				filtered = append(filtered, m)
			}
		}

		if len(filtered) > 0 {
			result = append(result, MetricGroup{
				Major:    g.Major,
				Minor:    g.Minor,
				FieldIds: filtered,
			})
		}
	}
	return result
}

func getSupportedMetricGroups(hcuIndex int) ([]MetricGroup, error) {
	typeName, _, err := DevTypeName(hcuIndex)
	if err != nil {
		return nil, err
	}
	switch {
	case strings.HasPrefix(typeName, "K100_AI"):
		return filterMetricsByName("K100_AI", profilingMetrics), nil
	case strings.HasPrefix(typeName, "K100"):
		return filterMetricsByName("K100", profilingMetrics), nil
	case strings.HasPrefix(typeName, "Z100"):
		return filterMetricsByName("Z100", profilingMetrics), nil
	case strings.HasPrefix(typeName, "BW"):
		return filterMetricsByName("BW", profilingMetrics), nil
	default:
		return profilingMetrics, nil
	}
}
