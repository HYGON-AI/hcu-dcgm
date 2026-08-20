/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package dcgm

/*
#cgo CFLAGS: -Wall -I./include
#cgo LDFLAGS: -ldl
#include <stdlib.h>
#include "driver_caps.h"
*/
import "C"
import (
	"fmt"
	"sync"
	"unsafe"

	"github.com/golang/glog"
)

// driverCaps 记录当前已加载驱动 .so 中可选 API 是否可用。
type driverCaps struct {
	HasXhclBandwidth bool
	HasUmcBandwidth  bool
	HasDevCuUsage    bool
	HasDevHcuUtil    bool
	HasDevCuUtil     bool
	HasDevWaveUtil   bool
	HasDevSeUtil     bool
}

var (
	driverCapability driverCaps
	capsOnce         sync.Once
)

var (
	errXhclBandwidthUnsupported = fmt.Errorf("driver does not support xhcl bandwidth (rsmi_dev_xhcl_bandwidth_get)")
	errUmcBandwidthUnsupported  = fmt.Errorf("driver does not support umc bandwidth (rsmi_dev_umc_bandwidth_get)")
	errDevCuUsageUnsupported    = fmt.Errorf("driver does not support cu usage (rsmi_dev_cu_usage_get)")
	errDevHcuUtilUnsupported    = fmt.Errorf("driver does not support hcu util (rsmi_dev_hcu_util_get)")
	errDevCuUtilUnsupported     = fmt.Errorf("driver does not support cu util (rsmi_dev_cu_util_get)")
	errDevWaveUtilUnsupported   = fmt.Errorf("driver does not support wave util (rsmi_dev_wave_util_get)")
	errDevSeUtilUnsupported     = fmt.Errorf("driver does not support se util (rsmi_dev_se_util_get)")
)

// probeDriverCaps 在 Init 成功后探测可选 RSMI 符号，仅执行一次。
func probeDriverCaps() {
	capsOnce.Do(func() {
		driverCapability.HasXhclBandwidth = rsmiSymbolExists("rsmi_dev_xhcl_bandwidth_get")
		driverCapability.HasUmcBandwidth = rsmiSymbolExists("rsmi_dev_umc_bandwidth_get")
		driverCapability.HasDevCuUsage = rsmiSymbolExists("rsmi_dev_cu_usage_get")
		driverCapability.HasDevHcuUtil = rsmiSymbolExists("rsmi_dev_hcu_util_get")
		driverCapability.HasDevCuUtil = rsmiSymbolExists("rsmi_dev_cu_util_get")
		driverCapability.HasDevWaveUtil = rsmiSymbolExists("rsmi_dev_wave_util_get")
		driverCapability.HasDevSeUtil = rsmiSymbolExists("rsmi_dev_se_util_get")
		glog.Infof("driver caps: xhcl_bandwidth=%v umc_bandwidth=%v cu_usage=%v hcu_util=%v cu_util=%v wave_util=%v se_util=%v",
			driverCapability.HasXhclBandwidth, driverCapability.HasUmcBandwidth,
			driverCapability.HasDevCuUsage, driverCapability.HasDevHcuUtil,
			driverCapability.HasDevCuUtil, driverCapability.HasDevWaveUtil, driverCapability.HasDevSeUtil)
	})
}

func rsmiSymbolExists(name string) bool {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return C.dcgm_lookup_symbol(cname) != nil
}

func ensureXhclBandwidth() error {
	probeDriverCaps()
	if !driverCapability.HasXhclBandwidth {
		return errXhclBandwidthUnsupported
	}
	return nil
}

func ensureUmcBandwidth() error {
	probeDriverCaps()
	if !driverCapability.HasUmcBandwidth {
		return errUmcBandwidthUnsupported
	}
	return nil
}

func ensureDevCuUsage() error {
	probeDriverCaps()
	if !driverCapability.HasDevCuUsage {
		return errDevCuUsageUnsupported
	}
	return nil
}

func ensureDevHcuUtil() error {
	probeDriverCaps()
	if !driverCapability.HasDevHcuUtil {
		return errDevHcuUtilUnsupported
	}
	return nil
}

func ensureDevCuUtil() error {
	probeDriverCaps()
	if !driverCapability.HasDevCuUtil {
		return errDevCuUtilUnsupported
	}
	return nil
}

func ensureDevWaveUtil() error {
	probeDriverCaps()
	if !driverCapability.HasDevWaveUtil {
		return errDevWaveUtilUnsupported
	}
	return nil
}

func ensureDevSeUtil() error {
	probeDriverCaps()
	if !driverCapability.HasDevSeUtil {
		return errDevSeUtilUnsupported
	}
	return nil
}
