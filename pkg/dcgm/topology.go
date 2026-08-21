/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package dcgm

/*
#cgo CFLAGS: -Wall -I./include
#cgo LDFLAGS: -L/opt/hyhal/lib -Wl,-rpath,/opt/hyhal/lib -lrocm_smi64 -Wl,--unresolved-symbols=ignore-in-object-files
#include <stdint.h>
#include <kfd_ioctl.h>
#include <rocm_smi64Config.h>
#include <rocm_smi.h>
*/
import "C"
import (
	"fmt"

	"github.com/golang/glog"
)

// rsmiTopoGetLinkWeight 获取2个gpu之间连接的权重
func rsmiTopoGetLinkWeight(dvIndSrc, dvIndDst int) (weight int64, err error) {
	var cweight C.uint64_t
	ret := C.rsmi_topo_get_link_weight(C.uint32_t(dvIndSrc), C.uint32_t(dvIndDst), &cweight)
	if err = errorString(ret); err != nil {
		return weight, fmt.Errorf("Error rsmiTopoGetLinkWeight:%S", err)
	}
	weight = int64(cweight)
	return
}

// rsmiTopoGetLinkType 获取2个gpu之间的hops和连接类型
func rsmiTopoGetLinkType(dvIndSrc, dvIndDst int) (hops int64, linkType RSMIIOLinkType, err error) {
	var chops C.uint64_t
	var clinkType C.RSMI_IO_LINK_TYPE
	ret := C.rsmi_topo_get_link_type(C.uint32_t(dvIndSrc), C.uint32_t(dvIndDst), &chops, &clinkType)
	if err = errorString(ret); err != nil {
		return hops, linkType, fmt.Errorf("Error rsmiTopoGetLinkType:%s", err)
	}
	hops = int64(chops)
	linkType = RSMIIOLinkType(clinkType)
	return
}

// rsmiTopoGetNumaBodeBumber 获取设备的numa cpu节点号
func rsmiTopoGetNumaBodeBumber(dvInd int) (numaNode int, err error) {
	var cnumaNode C.uint32_t
	ret := C.rsmi_topo_get_numa_node_number(C.uint32_t(dvInd), &cnumaNode)
	if err = errorString(ret); err != nil {
		return numaNode, fmt.Errorf("Error rsmiTopoGetNumaBodeBumber:%s", err)
	}
	numaNode = int(cnumaNode)
	return
}

func rsmiTopoIsHylink(srcDvInd, dstDvInd int) (bool, error) {
	var cIsHylink C.bool

	ret := C.rsmi_topo_is_hylink(
		C.uint32_t(srcDvInd),
		C.uint32_t(dstDvInd),
		&cIsHylink,
	)

	if err := errorString(ret); err != nil {
		return false, err
	}

	return bool(cIsHylink), nil
}

// rsmiDevXhclLinkNumber 获取指定 GPU 的 XHCL 链路数量
func rsmiDevXhclLinkNumber(dvInd int) (int, error) {
	var cLinkNum C.uint32_t

	ret := C.rsmi_dev_xhcl_link_number_get(
		C.uint32_t(dvInd),
		&cLinkNum,
	)

	glog.V(5).Infof(
		"rsmi_dev_xhcl_link_number_get dvInd: %v ret: %v linkNum: %v",
		dvInd,
		ret,
		cLinkNum,
	)

	if err := errorString(ret); err != nil {
		return 0, err
	}

	return int(cLinkNum), nil
}

// rsmiDevXhclLinkState 查询指定 GPU 的 XHCL 链路状态
func rsmiDevXhclLinkState(dvInd int, linkID int) (linkState uint32, err error) {
	var cLinkState C.uint32_t

	ret := C.rsmi_dev_xhcl_link_state_get(
		C.uint32_t(dvInd),
		C.uint32_t(linkID),
		&cLinkState,
	)

	glog.V(5).Infof(
		"rsmi_dev_xhcl_link_state_get dvInd: %v linkID: %v ret: %v linkState=%d",
		dvInd,
		linkID,
		ret,
		cLinkState,
	)

	if err = errorString(ret); err != nil {
		return 0, err
	}

	return uint32(cLinkState), nil
}

func DumpAllXhclLinkState(dvInd int) error {
	linkNum, err := rsmiDevXhclLinkNumber(dvInd)
	if err != nil {
		return err
	}

	for linkID := 0; linkID < linkNum; linkID++ {
		state, err := rsmiDevXhclLinkState(dvInd, linkID)
		if err != nil {
			fmt.Printf("GPU %d link %d error: %v\n", dvInd, linkID, err)
			continue
		}

		status := "DOWN"
		if state == 1 {
			status = "UP"
		}

		fmt.Printf(
			"GPU %d XHCL link %d/%d state: %s\n",
			dvInd, linkID, linkNum, status,
		)
	}
	return nil
}

// rsmiXhclLinkRemoteBdfidGet 获取指定 HCU 的某一条 XHCL 链路所连接的“远端设备”的 BDF ID。
//
// 该函数用于查询：
//   - HCU dvInd 上
//   - 第 linkID 条 XHCL 互联链路
//
// 实际连接到哪一个远端 PCI 设备。
//
// 返回的 bdfid 是远端设备的 PCI BDF（Bus / Device / Function）
// 编码值，可用于进一步映射到具体的 HCU index，
// 从而构建完整的物理互联拓扑关系（HCU ↔ HCU）。
//
// 参数说明：
//
//	dvInd   : 本端 HCU 的设备索引（0 ~ NumMonitorDevices()-1）
//	linkID  : XHCL 链路索引（通常来自 XhclLinkStates 返回的 LinkID）
//
// 返回值：
//
//	uint64  : 远端设备的 BDF ID（PCI Bus/Device/Function 编码）
//	error   : 当参数非法、链路不存在或底层接口失败时返回错误
//
// 可能的错误场景：
//   - RSMI_STATUS_INVALID_ARGS : dvInd / linkID 不合法
//   - RSMI_STATUS_NO_DATA      : 指定的 XHCL 链路不存在
//   - 其他 RSMI 错误码
//
// 使用场景示例：
//   - 解析 XHCL 链路的物理连接关系
//   - 将 XHCL link 映射为 HCU ↔ HCU 的邻接关系
//   - 拓扑完整性校验（是否存在断链 / 错链）
func rsmiXhclLinkRemoteBdfidGet(dvInd int, linkID int) (bdfid uint64, err error) {
	var cbdfid C.uint64_t // 必须用 C.uint64_t，避免 cgo 类型报错

	ret := C.rsmi_dev_xhcl_link_remote_bdfid_get(
		C.uint32_t(dvInd),
		C.uint32_t(linkID),
		&cbdfid, // 指针传给 C
	)

	if err = errorString(ret); err != nil {
		glog.Errorf(
			"rsmi_dev_xhcl_link_remote_bdfid_get failed: HCU: %v linkID: %v err: %v",
			dvInd, linkID, err,
		)
		return 0, err
	}

	// 转回 Go 的 uint64
	bdfid = uint64(cbdfid)

	// 打印返回值（重点）
	glog.V(5).Infof(
		"HCU %d XHCL link %d remote bdfid = 0x%x",
		dvInd,
		linkID,
		bdfid,
	)

	//// 拆解 bus:device.function
	//bus := (bdfid >> 8) & 0xff
	//device := (bdfid >> 3) & 0x1f
	//function := bdfid & 0x7
	//
	//glog.V(5).Infof(
	//	"HCU %d XHCL link %d remote BDF = %02x:%02x.%x",
	//	dvInd, linkID, bus, device, function,
	//)

	return bdfid, nil
}
