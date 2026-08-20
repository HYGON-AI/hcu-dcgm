/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package dcgm

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/golang/glog"
)

// ---------------- 配置结构体 ----------------

// Health 模块配置
type HealthConfig struct {
	MinPowerThreshold       int64   `json:"minPowerThreshold"`
	MaxPowerThreshold       int64   `json:"maxPowerThreshold"`
	MaxMemoryUsageThreshold int     `json:"maxMemoryUsageThreshold"`
	MaxTemperatureThreshold float64 `json:"maxTemperatureThreshold"`
}

// GEMM 模块配置
type GemmConfig struct {
	TargetStressTimeSeconds int     `json:"targetStressTimeSeconds"`
	Dgemm                   float64 `json:"dgemm"`
	Sgemm                   float64 `json:"sgemm"`
	Hgemm                   float64 `json:"hgemm"`
	Int8gemm                float64 `json:"int8gemm"`
	Bf16gemm                float64 `json:"bf16gemm"`
	Tf32gemm                float64 `json:"tf32gemm"`
}

// Bandwidth 模块配置
type BandwidthConfig struct {
	Precision int `json:"precision"`
	Arraysize int `json:"arraysize"`
	Numtimes  int `json:"numtimes"`
}

// 整体配置
type Config struct {
	Health    HealthConfig    `json:"health"`
	Gemmperf  GemmConfig      `json:"gemm"`
	Bandwidth BandwidthConfig `json:"bandwidth"`
}

// ---------------- 默认值 ----------------
var cfg = Config{
	Health: HealthConfig{
		MinPowerThreshold:       10,
		MaxPowerThreshold:       400,
		MaxMemoryUsageThreshold: 90,
		MaxTemperatureThreshold: 90.0,
	},
	Gemmperf: GemmConfig{
		TargetStressTimeSeconds: 1,
		Dgemm:                   20.6,
		Sgemm:                   40.6,
		Hgemm:                   200.0,
		Int8gemm:                330.1,
		Bf16gemm:                200.12,
		Tf32gemm:                110.1,
	},
	Bandwidth: BandwidthConfig{
		Precision: 1,
		Arraysize: 72001536,
		Numtimes:  200,
	},
}

// ---------------- 初始化配置 ----------------

// initConfig 初始化配置，JSON 覆盖默认值
func initConfig() {
	configPath, err := getConfigPath("config.json")
	if err != nil {
		glog.Warningf("Failed to get config path: %v, using defaults", err)
	} else {
		file, err := os.Open(configPath)
		if err == nil {
			defer file.Close()
			if err := json.NewDecoder(file).Decode(&cfg); err != nil {
				glog.Warningf("Error decoding config file, using defaults: %v", err)
			}
		}
	}

	// 无论是否读取到配置文件，都初始化模块变量
	InitHealthConfig(cfg)
	InitGemmConfig(cfg)
	InitBandwidthConfig(cfg)

	glog.V(5).Infof("Loaded configuration: %+v", cfg)
}

// getConfigPath 获取配置文件路径
func getConfigPath(filename string) (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}
	configDir := filepath.Join(filepath.Dir(execPath), "config")
	return filepath.Join(configDir, filename), nil
}

// ---------------- Health 配置初始化 ----------------
var (
	minPowerThreshold       int64
	maxPowerThreshold       int64
	maxMemoryUsageThreshold int
	maxTemperatureThreshold float64
)

// ---------------- GEMM 配置初始化 ----------------
var (
	TargetStressTimeSeconds int
	gemmList                []Gemm
)

type Gemm struct {
	Name     string
	Idx      int
	Expected float64
}

// ---------------- bandwidth 配置初始化 ----------------
var (
	bandwidthPrecision int
	bandwidthArraysize int
	bandwidthNumtimes  int
)

func InitBandwidthConfig(c Config) {
	bandwidthPrecision = c.Bandwidth.Precision
	bandwidthArraysize = c.Bandwidth.Arraysize
	bandwidthNumtimes = c.Bandwidth.Numtimes
}

func InitGemmConfig(c Config) {
	TargetStressTimeSeconds = c.Gemmperf.TargetStressTimeSeconds
	gemmList = []Gemm{
		{"hgemm", 0, float64(c.Gemmperf.Hgemm)},
		{"bf16gemm", 1, float64(c.Gemmperf.Bf16gemm)},
		{"sgemm", 2, float64(c.Gemmperf.Sgemm)},
		{"tf32gemm", 3, float64(c.Gemmperf.Tf32gemm)},
		{"i8gemm", 4, float64(c.Gemmperf.Int8gemm)},
		{"dgemm", 6, float64(c.Gemmperf.Dgemm)},
	}
}

func InitHealthConfig(c Config) {
	minPowerThreshold = c.Health.MinPowerThreshold
	maxPowerThreshold = c.Health.MaxPowerThreshold
	maxMemoryUsageThreshold = c.Health.MaxMemoryUsageThreshold
	maxTemperatureThreshold = c.Health.MaxTemperatureThreshold
}
