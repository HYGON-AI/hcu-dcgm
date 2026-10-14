# 3 HCU-DCGM二次开发

## 3.1 开发环境配置

### 3.1.1 开发环境配置

1.HCU-DCGM的Golang版本SDK目前没有开源，若有项目需求请联系售前同事申请。克隆或拷贝的HCU-DCGM项目需与二次开发项目放置于同级目录中。

2.修改二次开发项目中go.mod文件。

replace g.sugon.com/das/hcu-dcgm => /your/path/hcu-dcgm

3. 在二次开发项目中执行 go mod tidy 引入HCU-DCGM包。

4. 在go文件中import相关依赖包后使用调用，其中api.go为封装HCU-DCGM的API调用，提供与DCGM库交互的各种API接口，处理具体的功能调用。一个简单的调用实例如下：

```
package main

import (
	"flag"
	"log"
	"github.com/golang/glog"
	"g.sugon.com/das/hcu-dcgm/pkg/dcgm"
)

func main() {
	//初始化dcgm服务
	dcgm.Init()
	log.Println("go-dcgm Init ...")
	numDevices, _ := dcgm.NumMonitorDevices()
	log.Println("HCU number of devices: ", numDevices)
	defer dcgm.ShutDown()
	log.Println("go-dcgm ShutDown ...")
}
```



### 3.1.2 部署环境配置

基于HCU-DCGM二次开发的应用，运行依赖驱动链接库libhydmi、librocm_smi64和libhydmi_mig，6.3.*版本以上的驱动这些链接库默认在驱动安装目录/opt/hyhal/lib中。

```
进入/your/path/dcgm/lib目录创建对应缺失链接库的链接：

cd /your/path/dcgm/lib
#创建 libhydmi.so 相关软链<可选>
ln -s libhydmi.so.1.4 libhydmi.so.1
ln -s libhydmi.so.1 libhydmi.so

#创建 libroom_smi64.so 相关软链<可选>
ln -s libroom_smi64.so.2.8 libroom_smi64.so.2
ln -s libroom_smi64.so.2 libroom_smi64.so

#创建 libhydmi_mig.so 相关软链<可选>
ln -s libhydmi_mig.so.1.4 libhydmi_mig.so.1
ln -s libhydmi_mig.so.1 libhydmi_mig.so

补齐缺失链接库后，用户执行 export LD_LIBRARY_PATH=/your/path/dcgm/lib:$LD_LIBRARY_PATH 即可完成环境变量配置。
```



## 3.2 接口说明

### 3.2.1 Device

#### 3.2.1.1 获取所有物理设备信息（监控指标）

函数原型：AllDeviceInfos() ([]PhysicalDeviceInfo, error)

功能说明：所有物理设备的详细信息

参数说明：

**表 3.1 获取所有物理设备信息（监控指标）接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| device | Object | Device 物理设备的详细信息 |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |

#### 3.2.1.2 删除健康检查项的配置

函数原型：DeleteHealthCheckConfig() (err error)

功能说明：删除健康检查项的配置

#### 3.2.1.3 获取设备品牌名称

函数原型：DevBrand(dvInd int) (brand string, err error)

功能说明：根据设备索引获取品牌名称

参数说明：

**表 3.2 获取设备品牌名称接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | integer | Y | 设备索引 |

**表 3.3 获取设备品牌名称接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| brand | String | 设备品牌名称 |

#### 3.2.1.4 获取 GPU 度量信息

函数原型：DevGpuMetricsInfo(dvInd int) (gpuMetrics RSMIGPUMetrics, err error)

功能说明：获取 GPU 度量信息

参数说明：

**表 3.4 获取 GPU 度量信息接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | integer | Y | 设备索引 |

**表 3.5 获取 GPU 度量信息接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| averageDclk0Frequency | integer | AverageDclk0Frequency 平均显示时钟0频率 |
| averageDclk1Frequency | integer | AverageDclk1Frequency 平均显示时钟1频率 |
| averageGfxActivity | integer | AverageGfxActivity 平均图形活动 |
| averageGfxclkFrequency | integer | AverageGfxclkFrequency 平均图形时钟频率 |
| averageMmActivity | integer | AverageMmActivity 平均多媒体活动 |
| averageSocclkFrequency | integer | AverageSocclkFrequency 平均SoC时钟频率 |
| averageSocketPower | integer | AverageSocketPower 平均插座功率 |
| averageUclkFrequency | integer | AverageUclkFrequency 平均内存时钟频率 |
| averageUmcActivity | integer | AverageUmcActivity 平均内存控制器活动 |
| averageVclk0Frequency | integer | AverageVclk0Frequency 平均视频时钟0频率 |
| averageVclk1Frequency | integer | AverageVclk1Frequency 平均视频时钟1频率 |
| commonHeader | object | CommonHeader 公共表头 |
| currentDclk0 | integer | CurrentDclk0 当前显示时钟0 |
| currentDclk1 | integer | CurrentDclk1 当前显示时钟1 |
| currentFanSpeed | integer | CurrentFanSpeed 当前风扇速度 |
| currentGfxclk | integer | CurrentGfxclk 当前图形时钟 |
| currentSocclk | integer | CurrentSocclk 当前SoC时钟 |
| currentUclk | integer | CurrentUclk 当前内存时钟 |
| currentVclk0 | integer | CurrentVclk0 当前视频时钟0 |
| currentVclk1 | integer | CurrentVclk1 当前视频时钟1 |
| energyAccumulator | integer | EnergyAccumulator 能量累加器 |
| gfxActivityAcc | integer | GfxActivityAcc 图形活动累加器 |
| memActivityAcc | integer | MemActivityAcc 内存活动累加器 |
| padding | integer | Padding 填充 |
| pcieLinkSpeed | integer | PcieLinkSpeed PCIe链路速度（0.1 GT/s） |
| pcieLinkWidth | integer | PcieLinkWidth PCIe链路宽度 |
| systemClockCounter | integer | SystemClockCounter 系统时钟计数器 |
| temperatureEdge | integer | TemperatureEdge 边缘温度 |
| temperatureHotspot | integer | TemperatureHotspot 热点温度 |
| temperatureMem | integer | TemperatureMem 内存温度 |
| temperatureVRGfx | integer | TemperatureVRGfx VR图形温度 |
| temperatureVRMem | integer | TemperatureVRMem VR内存温度 |
| temperatureVRSoc | integer | TemperatureVRSoc VRSoC温度 |
| tempetureHBM | array | TempetureHBM 高带宽内存温度 |
| throttleStatus | integer | ThrottleStatus 节流状态 |

#### 3.2.1.5 获取设备ID的十六进制值

函数原型：DevTypeID (dvInd int) (id int, err error)

功能说明：获取设备ID的十六进制值

参数说明：

**表 3.6 获取设备ID的十六进制值接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | integer | Y | 设备索引 |

**表 3.7 获取设备ID的十六进制值接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| id | integer | ID的十六进制值 |

#### 3.2.1.6 获取设备SKU

函数原型：DevSku(dvInd int) (sku int, err error)

功能说明：获取设备SKU

参数说明：

**表 3.8 获取设备SKU接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | integer | Y | 设备索引 |

**表 3.9 获取设备SKU接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| sku | integer | 设备SKU |

#### 3.2.1.7 获取设备供应商名称

函数原型：DevVendorName(dvInd int) (vendorName string, err error)

功能说明：获取设备供应商名称

参数说明：

**表 3.10 获取设备供应商名称接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | integer | Y | 设备索引 |

**表 3.11 获取设备供应商名称接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| vendorName | string | 设备供应商名称 |

#### 3.2.1.8 获取设备显存供应商名称

函数原型：DevVramVendor(dvInd int) (vramVendor string, err error)

功能说明：获取设备显存供应商名称

参数说明：

**表 3.12 获取设备显存供应商名称接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | integer | Y | 设备索引 |

**表 3.13 获取设备显存供应商名称接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| vramVendor | string | 显存供应商名称 |

#### 3.2.1.9 获取设备信息列表

函数原型：DeviceInfos() (deviceInfos []DeviceInfo, err error)

功能说明：获取设备信息列表

参数说明：

**表 3.14 获取设备信息列表接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | integer | Y | 设备索引 |

**表 3.15 获取设备信息列表接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| DeviceId | string | 设备唯一标识 |
| ComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| DevType | string | 设备类型id |
| DevTypeName | string | 设备类型名称 |
| DvInd | integer | DvInd物理设备的设备ID |
| MemoryTotal | float64 | VMemoryTotal虚拟设备的全局内存大小 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| PowerTotal | float64 | 设备总功率 |
| VDeviceCount | integer | 虚拟卡数量 |

#### 3.2.1.10 获取所有设备的型号信息

函数原型：GetDeviceModelInfos() (devices []DeviceModelInfo)

功能说明：获取所有设备的型号信息

参数说明：

**表 3.16 获取所有设备的型号信息接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| cucount | number | CUCount CU数量 |
| memorySize | number | MemorySize 内存大小 |
| model | string | Model HCU类型 |

#### 3.2.1.11 获取健康检查项的配置

函数原型：GetHealthCheckConfig() (healthConfig HealthCheckConfig, err error)

功能说明：获取健康检查项的配置

参数说明：

**表 3.17 获取健康检查项的配置接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| enabled | boolean | Enabled 是否开启健康检查 |
| options | array | Options 健康检查项 |

#### 3.2.1.12 获取指定组中设备的健康检查结果

函数原型：HealthCheckByGroupId(groupId int, checkHealthConfig bool) (deviceHealths []DeviceHealth, err error)

功能说明：获取指定组中设备的健康检查结果

参数说明：

**表 3.18 获取指定组中设备的健康检查结果接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| groupId | int | Y | 组的id |
| checkHealthConfig | bool | Y | 是否使用配置项 |

**表 3.19 获取指定组中设备的健康检查结果接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| hcu | integer | HCU编号 |
| status | string | Status 状态 |
| watches | array:SystemWatch | Watches 检查详情 |
| error | string | Error 错误的详细信息 |
| result | object | Result 结果 |
| status | string | Status 状态 |
| type | string | Type 检查项的类型 |

#### 3.2.1.13 获取指定设备的健康检查结果

函数原型：HealthCheckById(dvIdList []int, checkHealthConfig bool) (deviceHealths []DeviceHealth, err error)

功能说明：获取指定设备的健康检查结果

参数说明：

**表 3.20 获取指定设备的健康检查结果接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvIdList | array | Y | 设备ID列表 |

**表 3.21 获取指定设备的健康检查结果接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| hcu | integer | HCU编号 |
| status | string | Status 状态 |
| watches | array:router.SystemWatch | Watches 检查详情 |
| error | string | Error 错误的详细信息 |
| result | object | Result 结果 |
| status | string | Status 状态 |
| type | string | Type 检查项的类型 |

#### 3.2.1.14 获取 HCU 数量

函数原型：NumMonitorDevices() (int, error)

功能说明：获取 HCU 数量

参数说明：

**表 3.22 获取 HCU 数量接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| number | int | HCU数量 |

#### 3.2.1.15 运行设备诊断

函数原型：RunDiag(level int) (diagResults DiagResults, err error)

功能说明：运行设备诊断

参数说明：

**表 3.23 运行设备诊断接口参数说明**

| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| --- | --- | --- | --- | --- |
| level | integer | query | Y | 诊断级别 |

**表 3.24 运行设备诊断接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| deviceNumber | string | DeviceNumber 设备数量信息 |
| perHCU | array:router.HCUResult | PerHCU HCU 的诊断结果 |
| hcu | integer | HCU HCU 的编号 |
| diagResults | array:router.DiagResult | DiagResults HCU 上的所有单独诊断测试的结果 |
| errorCode | integer | ErrorCode 错误码 |
| errorMessage | string | ErrorMessage 错误信息 |
| status | string | Status 诊断的结果状态 |
| testName | string | TestName 测试的名称 |
| testOutput | string | TestOutput 测试执行后生成的详细结果 |
| rc | integer | RC 返回代码 |
| software | array:router.DiagResult | Software 软件层面的诊断结果 |
| errorCode | integer | ErrorCode 错误码 |
| errorMessage | string | ErrorMessage 错误信息 |
| status | string | Status 诊断的结果状态 |
| testName | string | TestName 测试的名称 |
| testOutput | string | TestOutput 测试执行后生成的详细结果 |

#### 3.2.1.16 停止设备诊断

函数原型：StopDiag()

功能说明：停止设备诊断

#### 3.2.1.17 Memory Bandwidth 压测

函数原型：BandwidthTestResult(dvIdList []int) (map[int]float64, error)

功能说明：内存带宽测试

参数说明：

**表 3.25 Memory Bandwidth 压测接口参数说明**

| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| --- | --- | --- | --- | --- |
| dvIdList | []int | query | Y | 设备ID列表 |

**表 3.26 Memory Bandwidth 压测接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| result | map[int]float64 | 带宽测试结果，key 为设备 ID，value 为带宽值 |
| error | error | 测试执行状态或结果有效性说明 |

#### 3.2.1.18 PCIe 带宽压测

函数原型：PcieBandwidthTestResult() (PcieResult, error)

功能说明：PCIe 带宽测试（Sys<->Fb）

参数说明：

**表 3.27 PCIe 带宽压测接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| Matrix | [][]float64 | PCIe 带宽矩阵（单位：GB/s），`Matrix[i][j]` 表示设备 i → j 的单向带宽 |
| DeviceCount | int | 检测到的 HCU 设备数量，用于校验 `Matrix` 的有效维度 |
| CPUCount | int | 参与 PCIe 拓扑识别的 CPU root 数量 |
| HCUs | []HCUBW | 每个 HCU 的 Sys→Fb / Fb→Sys 带宽信息，仅在存在 HCU 时有效 |
| DvInd | int | HCU 设备索引（device index） |
| SysToFb | float64 | Sys → FB 的 PCIe 带宽（GB/s） |
| FbToSys | float64 | FB → Sys 的 PCIe 带宽（GB/s） |
| LogFile | string | 原始 PCIe 带宽测试日志文件路径 |
| error | error | 非 nil 表示测试执行失败、解析失败或结果无效；此时其他字段可能仍包含部分有效信息 |

#### 3.2.1.19 XHCL 带宽测试

函数原型：XHCLTestResult() ([]XHCLResult, error)

功能说明：XHCL（跨 HCU）互联带宽测试

参数说明：

**表 3.28 XHCL 带宽测试接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| SrcHCUId | int | 源 HCU（起点）设备索引 |
| DstHCUId | int | 目标 HCU（终点）设备索引 |
| BandwidthGBs | float64 | 源 HCU → 目标 HCU 的 XHCL 带宽，单位：GB/s |
| error | error | 非 nil 表示测试执行失败、解析失败或结果无效 |

#### 3.2.1.20 TargetStress（GEMM）压测

函数原型：TargetStressTestResult() (TargetStressResult, error)

功能说明：多种 GEMM 模式的 TargetStress 压测

参数说明：

**表 3.29 TargetStress（GEMM）压测接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| Results | []GemmTestResult | 目标压测（Target Stress）结果列表，每条表示某个 HCU 上某个 GEMM 的一次测试结果 |
| HCUId | int | HCU设备索引 |
| GemmName | string | GEMM 测试名称（如 hgemm、sgemm、dgemm 等） |
| Mean | float64 | 测得的平均性能值（单位以底层测试工具或日志定义为准，如 GFLOPS / MB/s） |
| Failed | bool | 测试是否失败：true`表示该 GEMM 测试失败或出错，false表示成功 |
| LogDir | string | 保存本次 Target Stress 测试日志的目录路径 |
| error | error | 非 nil 表示测试执行失败、部分 GEMM 测试异常或结果无效 |

#### 3.2.1.21 MemtestCL 压测

函数原型：MemtestCLTestResult(dvIdList []int) (MemtestCLAllResult, error)

功能说明：MemtestCL 对所有 HCU 进行测试

参数说明：

**表 3.30 MemtestCL 压测接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvIdList | []int | Y | 参与 memtestCL 测试的 HCU 设备索引列表 |

**表 3.31 MemtestCL 压测接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| Results | []MemtestCLResult | memtestCL 测试结果列表，每条对应一个 HCU |
| HCUId |  |  |

#### 3.2.1.22 初始化 (DCGM) 库

函数原型：func Init() (err error)

功能说明：初始化 (DCGM) 库。

#### 3.2.1.23 关闭 Data Center GPU Manager (DCGM) 库

函数原型：func ShutDown() error

功能说明：关闭 Data Center GPU Manager (DCGM) 库。

#### 3.2.1.24 根据设备索引获取利用率计数器

函数原型：func UtilizationCount(dvInd int, utilizationCounters []UtilizationCounter, count int) (timestamp int64, err error)

功能说明：根据设备索引获取利用率计数器

参数说明：

**表 3.189 根据设备索引获取利用率计数器接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |
| utilizationCounters | []UtilizationCounter | 是 | 根据函数签名推断的参数 |
| count | int | 是 | 根据函数签名推断的参数 |

**表 3.190 根据设备索引获取利用率计数器接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| timestamp | int64 | 函数返回值 |

#### 3.2.1.25 获取指定设备的 GPU 与显存粗粒度利用率计数器快照

函数原型：func DeviceUtilizationInfoGet(dvInd int) (info DeviceUtilizationInfo, err error)

功能说明：获取指定设备的 GPU 与显存粗粒度利用率计数器快照。

参数说明：

**表 3.191 获取指定设备的 GPU 与显存粗粒度利用率计数器快照接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.192 获取指定设备的 GPU 与显存粗粒度利用率计数器快照接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| info | DeviceUtilizationInfo | 函数返回值 |

#### 3.2.1.26 根据设备 ID 获取设备名称

函数原型：func DevName(dvInd int) (name string, err error)

功能说明：根据设备 ID 获取设备名称。

参数说明：

**表 3.193 根据设备 ID 获取设备名称接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.194 根据设备 ID 获取设备名称接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| name | string | 函数返回值 |

#### 3.2.1.27 根据设备 ID 获取设备的唯一标识符

函数原型：func GetDeviceId(dvInd int) (deviceId string, err error)

功能说明：根据设备 ID 获取设备的唯一标识符。

参数说明：

**表 3.195 根据设备 ID 获取设备的唯一标识符接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.196 根据设备 ID 获取设备的唯一标识符接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| deviceId | string | 函数返回值 |

#### 3.2.1.28 根据设备索引获取设备的唯一标识符

函数原型：func GetDeviceUniqueId(dvInd int) (string, error)

功能说明：根据设备索引获取设备的唯一标识符。

参数说明：

**表 3.197 GetDeviceUniqueId接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.198 GetDeviceUniqueId接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| 返回值1 | string | 函数返回值 |

#### 3.2.1.29 由 GetTempBySensor 统一转换并保留两位小数

函数原型：func GetDeviceTemperatureInfo(dvInd int) (temperatures DeviceTemperatureInfo, err error)

功能说明：由 GetTempBySensor 统一转换并保留两位小数。

参数说明：

**表 3.199 由 GetTempBySensor 统一转换并保留两位小数接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.200 由 GetTempBySensor 统一转换并保留两位小数接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| temperatures | DeviceTemperatureInfo | 函数返回值 |

#### 3.2.1.30 解析指定设备当前的降频状态

函数原型：func DevThrottleFlags(dvInd int) (flags ThrottleFlags, err error)

功能说明：解析指定设备当前的降频状态

参数说明：

**表 3.201 解析指定设备当前的降频状态接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.202 解析指定设备当前的降频状态接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| flags | ThrottleFlags | 函数返回值 |

#### 3.2.1.31 根据设备索引和传感器ID获取设备的最大和最小功率值范围

函数原型：func DevPowerCapRange(dvInd int) (powerMax, powerMin int64, err error)

功能说明：根据设备索引和传感器ID获取设备的最大和最小功率值范围

参数说明：

**表 3.203 根据设备索引和传感器ID获取设备的最大和最小功率值范围接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.204 根据设备索引和传感器ID获取设备的最大和最小功率值范围接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| 返回值1 | powerMax | 函数返回值 |
| powerMin | int64 | 函数返回值 |

#### 3.2.1.32 设置指定设备的功率上限，单位微瓦。有效范围可通过 /DevPowerCapRange 查询

函数原型：func DevPowerCapSet(dvInd int, cap uint64) error

功能说明：设置指定设备的功率上限，单位微瓦。有效范围可通过 /DevPowerCapRange 查询。

参数说明：

**表 3.205 设置指定设备的功率上限，单位微瓦。有效范围可通过 /DevPowerCapRange 查询接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |
| cap | uint64 | 是 | 根据函数签名推断的参数 |

#### 3.2.1.33 读取默认功率上限后重新写入，恢复出厂默认值

函数原型：func DevPowerCapReset(dvInd int) error

功能说明：读取默认功率上限后重新写入，恢复出厂默认值。

参数说明：

**表 3.206 读取默认功率上限后重新写入，恢复出厂默认值接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |

#### 3.2.1.34 设置设备性能等级，推荐使用此新版接口

函数原型：func DevPerfLevelSetV1(dvInd int, level DevPerfLevel) error

功能说明：设置设备性能等级，推荐使用此新版接口。

参数说明：

**表 3.207 设置设备性能等级，推荐使用此新版接口接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |
| level | DevPerfLevel | 是 | 根据函数签名推断的参数 |

#### 3.2.1.35 收集所有设备的监控指标信息

函数原型：func CollectDeviceMetrics() (monitorInfos []MonitorInfo, err error)

功能说明：收集所有设备的监控指标信息。

**表 3.208 收集所有设备的监控指标信息接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| monitorInfos | []MonitorInfo | 函数返回值 |

#### 3.2.1.36 根据设备的 dvInd 获取物理设备信息

函数原型：func GetDeviceByDvInd(dvInd int) (physicalDeviceInfo PhysicalDeviceInfo, err error)

功能说明：根据设备的 dvInd 获取物理设备信息

参数说明：

**表 3.209 根据设备的 dvInd 获取物理设备信息接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.210 根据设备的 dvInd 获取物理设备信息接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| physicalDeviceInfo | PhysicalDeviceInfo | 函数返回值 |

#### 3.2.1.37 根据物理设备索引获取该设备下的虚拟设备数量及虚拟设备索引列表

函数原型：func VDeviceByDvInd(dvInd int) (vDeviceCount int, vDevInds []int, err error)

功能说明：根据物理设备索引获取该设备下的虚拟设备数量及虚拟设备索引列表。

参数说明：

**表 3.211 VDeviceByDvInd接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.212 VDeviceByDvInd接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| vDeviceCount | int | 函数返回值 |
| vDevInds | []int | 函数返回值 |

#### 3.2.1.38 根据设备索引返回设备类型名称

函数原型：func DevTypeName(dvInd int) (devTypeName string, computeUnit float64, err error)

功能说明：根据设备索引返回设备类型名称

参数说明：

**表 3.213 根据设备索引返回设备类型名称接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.214 根据设备索引返回设备类型名称接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| devTypeName | string | 函数返回值 |
| computeUnit | float64 | 函数返回值 |

#### 3.2.1.39 根据设备索引获取设备子系统 ID

函数原型：func DevSubsystemId(dvInd int) (subsystemId string, err error)

功能说明：根据设备索引获取设备子系统 ID。

参数说明：

**表 3.215 DevSubsystemId接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.216 DevSubsystemId接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| subsystemId | string | 函数返回值 |

#### 3.2.1.40 根据设备索引获取设备子系统名称

函数原型：func DevSubsystemName(dvInd int) (name string, err error)

功能说明：根据设备索引获取设备子系统名称。

参数说明：

**表 3.217 DevSubsystemName接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.218 DevSubsystemName接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| name | string | 函数返回值 |

#### 3.2.1.41 返回所有设备设备状态信息

函数原型：func DeviceStatus() (deviceStatusInfos []DeviceStatusInfo, err error)

功能说明：返回所有设备设备状态信息

**表 3.219 返回所有设备设备状态信息接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| deviceStatusInfos | []DeviceStatusInfo | 函数返回值 |

#### 3.2.1.42 根据设备索引和 GPU 块获取 ECC 错误计数

函数原型：func EccCount(dvInd int, block RSMIGpuBlock) (errorCount RSMIErrorCount, err error)

功能说明：根据设备索引和 GPU 块获取 ECC 错误计数。

参数说明：

**表 3.220 EccCount接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |
| block | RSMIGpuBlock | 是 | 根据函数签名推断的参数 |

**表 3.221 EccCount接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| errorCount | RSMIErrorCount | 函数返回值 |

#### 3.2.1.43 根据设备索引和温度指标获取设备温度，并转换为摄氏度后保留两位小数

函数原型：func GetTempByMetric(dvInd int, metric RSMITemperatureMetric) (temp float64, err error)

功能说明：根据设备索引和温度指标获取设备温度，并转换为摄氏度后保留两位小数。

参数说明：

**表 3.222 GetTempByMetric接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |
| metric | RSMITemperatureMetric | 是 | 根据函数签名推断的参数 |

**表 3.223 GetTempByMetric接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| temp | float64 | 函数返回值 |

#### 3.2.1.44 设置设备超速百分比

函数原型：func DevOverdriveLevelSet(dvInd, od int) (err error)

功能说明：设置设备超速百分比

参数说明：

**表 3.224 设置设备超速百分比接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | — | 是 | 根据函数签名推断的参数 |
| od | int | 是 | 根据函数签名推断的参数 |

#### 3.2.1.45 获取设备的超速百分比

函数原型：func DevOverdriveLevelGet(dvInd int) (od int, err error)

功能说明：获取设备的超速百分比

参数说明：

**表 3.225 获取设备的超速百分比接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.226 获取设备的超速百分比接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| od | int | 函数返回值 |

#### 3.2.1.46 设置电压曲线

函数原型：func DevOdVoltInfoSet(dvInd, vPoint, clkValue, voltValue int) (err error)

功能说明：设置电压曲线

参数说明：

**表 3.227 设置电压曲线接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | — | 是 | 根据函数签名推断的参数 |
| vPoint | — | 是 | 根据函数签名推断的参数 |
| clkValue | — | 是 | 根据函数签名推断的参数 |
| voltValue | int | 是 | 根据函数签名推断的参数 |

#### 3.2.1.47 获取指定设备和传感器的功耗配置文件预设状态

函数原型：func DevPowerProfilePresetsGet(dvInd, sensorInd int) (powerProfileStatus PowerProfileStatus, err error)

功能说明：获取指定设备和传感器的功耗配置文件预设状态。

参数说明：

**表 3.228 DevPowerProfilePresetsGet接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | — | 是 | 根据函数签名推断的参数 |
| sensorInd | int | 是 | 根据函数签名推断的参数 |

**表 3.229 DevPowerProfilePresetsGet接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| powerProfileStatus | PowerProfileStatus | 函数返回值 |

#### 3.2.1.48 显示指定设备列表的简要硬件信息

函数原型：func ShowAllConciseHw(dvIdList []int)

功能说明：显示指定设备列表的简要硬件信息

参数说明：

**表 3.230 显示指定设备列表的简要硬件信息接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvIdList | []int | 是 | 根据函数签名推断的参数 |

#### 3.2.1.49 获取指定设备上的进程信息

函数原型：func ProcessInfoByDevice(pid int, dvInd int) (proc ProcessInfos, err error)

功能说明：获取指定设备上的进程信息

参数说明：

**表 3.231 获取指定设备上的进程信息接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| pid | int | 是 | 根据函数签名推断的参数 |
| dvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.232 获取指定设备上的进程信息接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| proc | ProcessInfos | 函数返回值 |

#### 3.2.1.50 获取设备电压/频率曲线信息(K100 AI不支持)

函数原型：func DevOdVoltInfoGet(deInd int) (odv OdVoltFreqData, err error)

功能说明：获取设备电压/频率曲线信息(K100 AI不支持)

参数说明：

**表 3.233 获取设备电压/频率曲线信息(K100 AI不支持)接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| deInd | int | 是 | 根据函数签名推断的参数 |

**表 3.234 获取设备电压/频率曲线信息(K100 AI不支持)接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| odv | OdVoltFreqData | 函数返回值 |

#### 3.2.1.51 显示指定设备的产品名称、供应商、系列、型号和SKU信息

函数原型：func ShowProductName(dvIdList []int) (deviceProductInfos []DeviceproductInfo, err error)

功能说明：获取并显示指定设备的产品名称、供应商、系列、型号和SKU信息。

参数说明：

**表 3.235 显示指定设备的产品名称、供应商、系列、型号和SKU信息接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvIdList | []int | 是 | 根据函数签名推断的参数 |

**表 3.236 显示指定设备的产品名称、供应商、系列、型号和SKU信息接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| deviceProductInfos | []DeviceproductInfo | 函数返回值 |

#### 3.2.1.52 显示指定设备的序列号信息

函数原型：func ShowSerialNumber(dvIdList []int) (deviceSerialInfos []DeviceSerialInfo, err error)

功能说明：获取并显示指定设备的序列号信息。

参数说明：

**表 3.237 显示指定设备的序列号信息接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvIdList | []int | 是 | 根据函数签名推断的参数 |

**表 3.238 显示指定设备的序列号信息接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| deviceSerialInfos | []DeviceSerialInfo | 函数返回值 |

#### 3.2.1.53 显示指定设备的唯一ID信息

函数原型：func ShowUId(dvIdList []int) (deviceUIdInfos []DeviceUIdInfo, err error)

功能说明：获取并显示指定设备的唯一ID信息。

参数说明：

**表 3.239 显示指定设备的唯一ID信息接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvIdList | []int | 是 | 根据函数签名推断的参数 |

**表 3.240 显示指定设备的唯一ID信息接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| deviceUIdInfos | []DeviceUIdInfo | 函数返回值 |

#### 3.2.1.54 获取当前系统中的设备数量

函数原型：func DeviceCount() (count int, err error)

功能说明：获取当前系统中的设备数量。

**表 3.241 获取当前系统中的设备数量接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| count | int | 函数返回值 |

#### 3.2.1.55 返回物理设备时间忙碌百分比

函数原型：func DevBusyPercent(dvInd int) (utilizationRate float64, err error)

功能说明：返回物理设备时间忙碌百分比

参数说明：

**表 3.242 返回物理设备时间忙碌百分比接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.243 返回物理设备时间忙碌百分比接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| utilizationRate | float64 | 函数返回值 |

#### 3.2.1.56 根据虚拟设备索引返回当前虚拟设备的使用百分比

函数原型：func VDevBusyPercent(vDvInd int) (percent int, err error)

功能说明：根据虚拟设备索引返回当前虚拟设备的使用百分比

参数说明：

**表 3.244 根据虚拟设备索引返回当前虚拟设备的使用百分比接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| vDvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.245 根据虚拟设备索引返回当前虚拟设备的使用百分比接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| percent | int | 函数返回值 |

#### 3.2.1.57 根据设备索引返回设备的详细信息

函数原型：func GetDeviceInfo(dvInd int) (deviceInfo DMIDeviceInfo, err error)

功能说明：根据设备索引返回设备的详细信息

参数说明：

**表 3.246 根据设备索引返回设备的详细信息接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.247 根据设备索引返回设备的详细信息接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| deviceInfo | DMIDeviceInfo | 函数返回值 |

#### 3.2.1.58 根据参数设置是否开启健康检查，并选择需要的检查项

函数原型：func SetHealthCheckConfig(enabled bool, options []int) (err error)

功能说明：根据参数设置是否开启健康检查，并选择需要的检查项

参数说明：

**表 3.248 根据参数设置是否开启健康检查，并选择需要的检查项接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| enabled | bool | 是 | 根据函数签名推断的参数 |
| options | []int | 是 | 根据函数签名推断的参数 |

#### 3.2.1.59 返回所有可用的健康检查项及其配置信息

函数原型：func HealthCheckConfigList() (map[string]interface{}, error)

功能说明：返回所有可用的健康检查项及其配置信息

**表 3.249 返回所有可用的健康检查项及其配置信息接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| 返回值1 | map[string]interface{} | 函数返回值 |

#### 3.2.1.60 获取当前系统中的设备数量

函数原型：func DeviceGetCount() (deviceCount int, err error)

功能说明：获取当前系统中的设备数量。

**表 3.250 DeviceGetCount接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| deviceCount | int | 函数返回值 |

#### 3.2.1.61 根据设备索引获取设备句柄

函数原型：func DeviceGetHandleByIndex(dvInd int) (device MIGDevice, err error)

功能说明：根据设备索引获取设备句柄。

参数说明：

**表 3.251 DeviceGetHandleByIndex接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.252 DeviceGetHandleByIndex接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| device | MIGDevice | 函数返回值 |

#### 3.2.1.62 获取指定设备支持的最大 MIG 设备数量

函数原型：func DeviceGetMaxMigDeviceCountByIndex(dvInd int) (count int, err error)

功能说明：获取指定设备支持的最大 MIG 设备数量。

参数说明：

**表 3.253 DeviceGetMaxMigDeviceCountByIndex接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.254 DeviceGetMaxMigDeviceCountByIndex接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| count | int | 函数返回值 |

#### 3.2.1.63 根据设备索引和 MIG 设备索引获取设备属性

函数原型：func DeviceGetAttributesByIndex(dvInd int, migId int) (attr NvmlDeviceAttributes, err error)

功能说明：根据设备索引和 MIG 设备索引获取设备属性。

参数说明：

**表 3.255 DeviceGetAttributesByIndex接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |
| migId | int | 是 | 根据函数签名推断的参数 |

**表 3.256 DeviceGetAttributesByIndex接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| attr | NvmlDeviceAttributes | 函数返回值 |

#### 3.2.1.64 获取指定设备当前和待生效的 MIG 模式

函数原型：func DeviceGetMigModeByIndex(dvInd int) (currentMode uint32, pendingMode uint32, err error)

功能说明：获取指定设备当前和待生效的 MIG 模式。

参数说明：

**表 3.257 DeviceGetMigModeByIndex接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.258 DeviceGetMigModeByIndex接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| currentMode | uint32 | 函数返回值 |
| pendingMode | uint32 | 函数返回值 |

#### 3.2.1.65 获取指定设备的 GPU 实例配置文件信息

函数原型：func DeviceGetGpuInstanceProfileInfo(dvInd int, profileIdx uint32) (info NvmlGpuInstanceProfileInfo, err error)

功能说明：获取指定设备的 GPU 实例配置文件信息。

参数说明：

**表 3.259 DeviceGetGpuInstanceProfileInfo接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |
| profileIdx | uint32 | 是 | 根据函数签名推断的参数 |

**表 3.260 DeviceGetGpuInstanceProfileInfo接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| info | NvmlGpuInstanceProfileInfo | 函数返回值 |

#### 3.2.1.66 获取指定设备指定 GPU 实例配置文件的剩余容量

函数原型：func DeviceGetGpuInstanceRemainingCapacity(dvInd int, profileId uint32) (count uint32, err error)

功能说明：获取指定设备指定 GPU 实例配置文件的剩余容量。

参数说明：

**表 3.261 DeviceGetGpuInstanceRemainingCapacity接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |
| profileId | uint32 | 是 | 根据函数签名推断的参数 |

**表 3.262 DeviceGetGpuInstanceRemainingCapacity接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| count | uint32 | 函数返回值 |

#### 3.2.1.67 获取指定设备指定配置文件的 GPU 实例信息列表

函数原型：func DeviceGetGpuInstancesInfo(dvInd int, profileId uint32) ([]GpuInstanceInfo, error)

功能说明：获取指定设备指定配置文件的 GPU 实例信息列表。

参数说明：

**表 3.263 DeviceGetGpuInstancesInfo接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |
| profileId | uint32 | 是 | 根据函数签名推断的参数 |

**表 3.264 DeviceGetGpuInstancesInfo接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| 返回值1 | []GpuInstanceInfo | 函数返回值 |

#### 3.2.1.68 获取指定设备上 MIG 设备对应的 GPU 实例 ID

函数原型：func DeviceGetGpuInstanceId(dvInd int, migId int) (gpuInstanceId uint32, err error)

功能说明：获取指定设备上 MIG 设备对应的 GPU 实例 ID。

参数说明：

**表 3.265 DeviceGetGpuInstanceId接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |
| migId | int | 是 | 根据函数签名推断的参数 |

**表 3.266 DeviceGetGpuInstanceId接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| gpuInstanceId | uint32 | 函数返回值 |

#### 3.2.1.69 获取指定设备上 MIG 设备对应的计算实例 ID

函数原型：func DeviceGetComputeInstanceId(dvInd int, migId int) (computeInstanceId uint32, err error)

功能说明：获取指定设备上 MIG 设备对应的计算实例 ID。

参数说明：

**表 3.267 DeviceGetComputeInstanceId接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |
| migId | int | 是 | 根据函数签名推断的参数 |

**表 3.268 DeviceGetComputeInstanceId接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| computeInstanceId | uint32 | 函数返回值 |

#### 3.2.1.70 获取指定设备当前可用的 MIG 设备数量及索引列表

函数原型：func AvailableMigDeviceIds(dvInd int) (int, []int, error)

功能说明：获取指定设备当前可用的 MIG 设备数量及索引列表。

参数说明：

**表 3.269 AvailableMigDeviceIds接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.270 AvailableMigDeviceIds接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| 返回值1 | int | 函数返回值 |
| 返回值2 | []int | 函数返回值 |

#### 3.2.1.71 根据设备索引和 MIG 设备索引获取 MIG 设备句柄

函数原型：func DeviceGetMigDeviceHandleByIndex(dvInd int, migId int) (migDevice MIGDevice, err error)

功能说明：根据设备索引和 MIG 设备索引获取 MIG 设备句柄。

参数说明：

**表 3.271 DeviceGetMigDeviceHandleByIndex接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |
| migId | int | 是 | 根据函数签名推断的参数 |

**表 3.272 DeviceGetMigDeviceHandleByIndex接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| migDevice | MIGDevice | 函数返回值 |

#### 3.2.1.72 获取指定 MIG 设备的 PCI 信息

函数原型：func DeviceGetPciInfo(device MIGDevice) (pciInfo NvmlPciInfo, err error)

功能说明：获取指定 MIG 设备的 PCI 信息。

参数说明：

**表 3.273 DeviceGetPciInfo接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| device | MIGDevice | 是 | 根据函数签名推断的参数 |

**表 3.274 DeviceGetPciInfo接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| pciInfo | NvmlPciInfo | 函数返回值 |

#### 3.2.1.73 获取系统的着色引擎（SE）数量

函数原型：func MIGSECount() (int, error)

功能说明：获取系统的着色引擎（SE）数量。

**表 3.275 MIGSECount接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| 返回值1 | int | 函数返回值 |

#### 3.2.1.74 重置指定 HCU 设备

函数原型：func DevHCUReset(dvInd int) (err error)

功能说明：重置指定 HCU 设备。

参数说明：

**表 3.276 DevHCUReset接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |

#### 3.2.1.75 返回设备的 CU（Compute Unit）数量

函数原型：func DevCuNum(dvInd int) (int, error)

功能说明：返回设备的 CU（Compute Unit）数量

参数说明：

**表 3.277 返回设备的 CU（Compute Unit）数量接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.278 返回设备的 CU（Compute Unit）数量接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| 返回值1 | int | 函数返回值 |

#### 3.2.1.76 返回设备的 KFD node ID，用于拓扑关联

函数原型：func DevNodeId(dvInd int) (uint32, error)

功能说明：返回设备的 KFD node ID，用于拓扑关联

参数说明：

**表 3.279 返回设备的 KFD node ID，用于拓扑关联接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.280 返回设备的 KFD node ID，用于拓扑关联接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| 返回值1 | uint32 | 函数返回值 |

#### 3.2.1.77 返回显存 overdrive 百分比（0-20）

函数原型：func DevMemOverdriveLevelGet(dvInd int) (int, error)

功能说明：返回显存 overdrive 百分比（0-20）

参数说明：

**表 3.281 返回显存 overdrive 百分比（0-20）接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.282 返回显存 overdrive 百分比（0-20）接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| 返回值1 | int | 函数返回值 |

#### 3.2.1.78 返回 GFX（shader）时钟的最大额定频率（MHz）

函数原型：func DevGfxClockMax(dvInd int) (uint64, error)

功能说明：返回 GFX（shader）时钟的最大额定频率（MHz）

参数说明：

**表 3.283 返回 GFX（shader）时钟的最大额定频率（MHz）接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.284 返回 GFX（shader）时钟的最大额定频率（MHz）接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| 返回值1 | uint64 | 函数返回值 |

#### 3.2.1.79 返回显存时钟的最大额定频率（MHz）

函数原型：func DevMemClockMax(dvInd int) (uint64, error)

功能说明：返回显存时钟的最大额定频率（MHz）

参数说明：

**表 3.285 返回显存时钟的最大额定频率（MHz）接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.286 返回显存时钟的最大额定频率（MHz）接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| 返回值1 | uint64 | 函数返回值 |

#### 3.2.1.80 查询设备 CU 在采样窗口内的 wave 占用周期占比（全 CU 平均）

函数原型：func DevCuUtil(dvInd int, duration int) (float32, error)

功能说明：查询设备 CU 在采样窗口内的 wave 占用周期占比（全 CU 平均）。

参数说明：

**表 3.287 查询设备 CU 在采样窗口内的 wave 占用周期占比（全 CU 平均）接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |
| duration | int | 是 | 根据函数签名推断的参数 |

**表 3.288 查询设备 CU 在采样窗口内的 wave 占用周期占比（全 CU 平均）接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| 返回值1 | float32 | 函数返回值 |

#### 3.2.1.81 查询设备 CU 在采样窗口内的 wave 驻留数量占比（全 CU 平均）

函数原型：func DevWaveUtil(dvInd int, duration int) (float32, error)

功能说明：查询设备 CU 在采样窗口内的 wave 驻留数量占比（全 CU 平均）。

参数说明：

**表 3.289 查询设备 CU 在采样窗口内的 wave 驻留数量占比（全 CU 平均）接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |
| duration | int | 是 | 根据函数签名推断的参数 |

**表 3.290 查询设备 CU 在采样窗口内的 wave 驻留数量占比（全 CU 平均）接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| 返回值1 | float32 | 函数返回值 |

#### 3.2.1.82 所有 HCU 卡的健康状态

函数原型：func HCUHealthCheck() ([]HCUHealthStatus, error)

功能说明：检查所有 HCU 卡的健康状态

**表 3.291 所有 HCU 卡的健康状态接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| 返回值1 | []HCUHealthStatus | 函数返回值 |

### 3.2.2 Error

#### 3.2.2.1 获取 ECC block 信息

函数原型：EccBlocksInfo(dvInd int) (blocksInfos []BlocksInfo, err error)

功能说明：获取 ECC block 信息

参数说明：

**表 3.32 获取 ECC block 信息接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | integer | Y | 设备索引 |

**表 3.33 获取 ECC block 信息接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| block | string | Block block |
| ce | integer | CE CE |
| state | string | State 状态 |
| ue | integer | UE UE |

#### 3.2.2.2 获取GPU块的ECC状态

函数原型：EccStatus(dvInd int, block RSMIGpuBlock) (state string, err error)

功能说明：获取GPU块的ECC状态

参数说明：

**表 3.34 获取GPU块的ECC状态接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | integer | Y | 设备索引 |
| block | string | Y | GPU块 |

#### 3.2.2.3 重置指定设备的XGMI错误状态

函数原型：ResetXGMIErr(dvIdList []int) (failedMessage []FailedMessage)

功能说明：重置指定设备的XGMI错误状态

参数说明：

**表 3.35 重置指定设备的XGMI错误状态接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvIdList | array | Y | 设备ID列表 |

**表 3.36 重置指定设备的XGMI错误状态接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| errorMsg | string | ErrorMsg 错误信息 |
| id | integer | ID 设备ID |

#### 3.2.2.4 获取XGMI错误状态

函数原型：XGMIErrorStatus(dvInd int) (status RSMIXGMIStatus, err error)

功能说明：获取XGMI错误状态

参数说明：

**表 3.37 获取XGMI错误状态接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | integer | Y | 设备索引 |

**表 3.38 获取XGMI错误状态接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| status | RSMIXGMIStatus | XGMI错误状态码 |

#### 3.2.2.5 获取指定设备的XGMI hive id

函数原型：XGMIHiveIdGet(dvInd int) (hiveId int64, err error)

功能说明：获取指定设备的XGMI hive id

参数说明：

**表 3.39 获取指定设备的XGMI hive id接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | integer | Y | 设备索引 |

**表 3.40 获取指定设备的XGMI hive id接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| hiveId | integer | 设备的XGMI hive id |

### 3.2.3 Event

#### 3.2.3.1 打印设备的事件列表

函数原型：PrintEventList(dvInd int, delay int, eventList []string)

功能说明：打印设备的事件列表

参数说明：

**表 3.41 打印设备的事件列表接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | integer | Y | 设备索引 |
| delay | integer | Y | 延迟时间（秒） |
| eventList | array | Y | 事件列表 |

#### 3.2.3.2 显示设备的事件

函数原型：ShowEvents(dvIdList []int, eventTypes []string)

功能说明：显示设备的事件

参数说明：

**表 3.42 显示设备的事件接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvIdList | array | Y | 设备ID列表 |
| eventTypes | array | Y | 事件类型列表 |

#### 3.2.2.6 获取指定设备启用的 ECC block 位掩码

函数原型：func EccEnabled(dvInd int) (enabledBlocks int64, err error)

功能说明：获取指定设备启用的 ECC block 位掩码。

参数说明：

**表 3.292 EccEnabled接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.293 EccEnabled接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| enabledBlocks | int64 | 函数返回值 |

#### 3.2.2.7 获取指定物理设备的HSL（高速互连链路）错误状态

函数原型：func HSLErrorStatus(dvInd int) (status RSMIHSLStatus, err error)

功能说明：获取指定物理设备的HSL（高速互连链路）错误状态。

参数说明：

**表 3.294 获取指定物理设备的HSL（高速互连链路）错误状态接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.295 获取指定物理设备的HSL（高速互连链路）错误状态接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| status | RSMIHSLStatus | 函数返回值 |

#### 3.2.2.8 读取显存 bank 级 ECC 统计信息

函数原型：func MemBankEccInfo(dvInd int) (map[string]string, error)

功能说明：读取显存 bank 级 ECC 统计信息。

参数说明：

**表 3.296 读取显存 bank 级 ECC 统计信息接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.297 读取显存 bank 级 ECC 统计信息接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| 返回值1 | map[string]string | 函数返回值 |

#### 3.2.2.9 启用或禁用指定 RAS block

函数原型：func SetRasFeature(dvInd int, block string, enabled bool) error

功能说明：启用或禁用指定 RAS block。

参数说明：

**表 3.298 启用或禁用指定 RAS block接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |
| block | string | 是 | 根据函数签名推断的参数 |
| enabled | bool | 是 | 根据函数签名推断的参数 |

#### 3.2.2.10 向指定设备写入驱动定义的 RAS 注入参数

函数原型：func InjectRasError(dvInd int, parameters []string) error

功能说明：向指定设备写入驱动定义的 RAS 注入参数。

参数说明：

**表 3.299 向指定设备写入驱动定义的 RAS 注入参数接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |
| parameters | []string | 是 | 根据函数签名推断的参数 |

### 3.2.4 Memory

#### 3.2.4.1 获取设备的指定内存使用情况

函数原型：MemInfo(dvInd int, memType string) (memUsed int64, memTotal int64, err error)

功能说明：获取设备的指定内存使用情况

参数说明：

**表 3.43 获取设备的指定内存使用情况接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | integer | Y | 设备索引 |
| memType | string | Y | 内存类型（可选值: vram, vis_vram, gtt） |

**表 3.44 获取设备的指定内存使用情况接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| memUsed | int64 | 内存使用量 |
| memTotal | memTotal | 内存总量 |

#### 3.2.4.2 获取内存使用百分比

函数原型：MemoryPercent(dvInd int) (busyPercent int, err error)

功能说明：获取内存使用百分比

参数说明：

**表 3.45 获取内存使用百分比接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | integer | Y | 设备索引 |

**表 3.46 获取内存使用百分比接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| busyPercent | integer | 内存使用百分比 |

#### 3.2.4.3 显示设备的退役页信息

函数原型：ShowRetiredPages(dvIdList []int, retiredType string)

功能说明：显示设备的退役页信息

参数说明：

**表 3.47 显示设备的退役页信息接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvIdList | array | Y | 设备ID列表 |
| retiredType | object | N | 退役类型 (默认为'all') |

#### 3.2.4.4 获取设备内存总量

函数原型：MemoryTotal(dvInd int) (memoryTotal float64, err error)

功能说明：获取设备内存总量

参数说明：

**表 3.48 获取设备内存总量接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | integer | Y | 设备索引 |

**表 3.49 获取设备内存总量接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| memoryTotal | float64 | 内存总量 |

#### 3.2.4.5 获取设备内存使用量

函数原型：MemoryUsed(dvInd int) (memoryUsed float64, err error)

功能说明：获取设备内存使用量

参数说明：

**表 3.50 获取设备内存使用量接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | integer | Y | 设备索引 |

**表 3.51 获取设备内存使用量接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| memoryUsed | float64 | 内存使用量 |

#### 3.2.4.6 检测内存保留页

函数原型：func MemoryReservedPages(dvInd int) (records []RetiredPageRecord, err error)

功能说明：检测内存保留页

参数说明：

**表 3.300 检测内存保留页接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.301 检测内存保留页接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| records | []RetiredPageRecord | 函数返回值 |

#### 3.2.4.7 返回驱动当前可向应用分配的近似剩余显存量（字节）

函数原型：func MemoryAvailable(dvInd int) (available int64, err error)

功能说明：返回驱动当前可向应用分配的近似剩余显存量（字节）

参数说明：

**表 3.302 返回驱动当前可向应用分配的近似剩余显存量（字节）接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.303 返回驱动当前可向应用分配的近似剩余显存量（字节）接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| available | int64 | 函数返回值 |

#### 3.2.4.8 返回显存颗粒当前温度（摄氏度）

函数原型：func MemoryTemperature(dvInd int) (float64, error)

功能说明：返回显存颗粒当前温度（摄氏度）

参数说明：

**表 3.304 返回显存颗粒当前温度（摄氏度）接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.305 返回显存颗粒当前温度（摄氏度）接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| 返回值1 | float64 | 函数返回值 |

#### 3.2.4.9 获取并展示指定设备的内存使用情况，包括不同类型的内存

函数原型：func ShowMemInfo(dvIdList []int, memTypes []string)

功能说明：获取并展示指定设备的内存使用情况，包括不同类型的内存。

参数说明：

**表 3.306 获取并展示指定设备的内存使用情况，包括不同类型的内存接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvIdList | []int | 是 | 根据函数签名推断的参数 |
| memTypes | []string | 是 | 根据函数签名推断的参数 |

#### 3.2.4.10 获取并展示指定设备的当前内存使用百分比和其他相关的利用率数据

函数原型：func ShowMemUse(dvIdList []int)

功能说明：获取并展示指定设备的当前内存使用百分比和其他相关的利用率数据。

参数说明：

**表 3.307 获取并展示指定设备的当前内存使用百分比和其他相关的利用率数据接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvIdList | []int | 是 | 根据函数签名推断的参数 |

#### 3.2.4.11 获取并展示指定设备的内存供应商信息

函数原型：func ShowMemVendor(dvIdList []int) (deviceMemVendorInfos []DeviceMemVendorInfo, err error)

功能说明：获取并展示指定设备的内存供应商信息。

参数说明：

**表 3.308 获取并展示指定设备的内存供应商信息接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvIdList | []int | 是 | 根据函数签名推断的参数 |

**表 3.309 获取并展示指定设备的内存供应商信息接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| deviceMemVendorInfos | []DeviceMemVendorInfo | 函数返回值 |

#### 3.2.4.12 Memtestcl压力测试，测试显存稳定性

函数原型：func MemtestCL(dvIdList []int) error

功能说明：Memtestcl压力测试，测试显存稳定性

参数说明：

**表 3.310 Memtestcl压力测试，测试显存稳定性接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvIdList | []int | 是 | 根据函数签名推断的参数 |

### 3.2.5 PCIe

#### 3.2.5.1 获取可用的 PCIe 带宽列表

函数原型：DevPciBandwidth(dvInd int) (rsmiPcieBandwidth RSMIPcieBandwidth, err error)

功能说明：获取可用的 PCIe 带宽列表

参数说明：

**表 3.52 获取可用的 PCIe 带宽列表接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | integer | Y | 设备索引 |

**表 3.53 获取可用的 PCIe 带宽列表接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| lanes | array | lanes 表示 PCIe 通道的配置 |
| transferRate | object | TransferRate 表示传输速率的频率信息 |

#### 3.2.5.2 设置设备允许的 PCIe 带宽

函数原型：DevPciBandwidthSet(dvInd int, bwBitmask int64) (err error)

功能说明：设置设备允许的 PCIe 带宽

参数说明：

**表 3.54 设置设备允许的 PCIe 带宽接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | integer | Y | 设备索引 |
| bwBitmask | integer | Y | 带宽掩码 |

#### 3.2.5.3 展示设备的PCIe带宽使用情况

函数原型：ShowPcieBw(dvIdList []int) (pcieBandwidthInfos []PcieBandwidthInfo, err error)

功能说明：设备的PCIe带宽使用情况

参数说明：

**表 3.55 展示设备的PCIe带宽使用情况接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvIdList | array | Y | 设备ID列表 |

**表 3.56 展示设备的PCIe带宽使用情况接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| bw | number | Bw bw |
| deviceID | integer | DeviceID 设备索引号 |
| maxPktSz | integer | MaxPktSz 最大PktSz |
| received | integer | Received 接收 |
| sent | integer | Sent 发送 |

#### 3.2.5.4 展示设备的PCIe重放计数

函数原型：ShowPcieReplayCount(dvIdList []int) (pcieReplayCountInfos []PcieReplayCountInfo, err error)

功能说明：设备的PCIe重放计数

参数说明：

**表 3.57 展示设备的PCIe重放计数接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvIdList | array | Y | 设备ID列表 |

**表 3.58 展示设备的PCIe重放计数接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| count | integer | 重放总数 |
| deviceID | integer | DeviceID 设备索引号 |

#### 3.2.5.5 获取总线信息

函数原型：PicBusInfo(dvInd int) (picID string, err error)

功能说明：获取总线信息

参数说明：

**表 3.59 获取总线信息接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | integer | Y | 设备索引 |

**表 3.60 获取总线信息接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| picID | string | 设备总线信息 |

#### 3.2.5.6 设备的PCIe带宽

函数原型：ShowPcieBw (dvInd int) (pcieBandwidthInfo PcieBandwidthInfo, err error)

功能说明：设备的PCIe带宽

参数说明：

**表 3.61 设备的PCIe带宽接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | integer | Y | 设备索引 |

**表 3.62 设备的PCIe带宽接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| bw | float64 | Bw bw |
| DvInd | integer | DeviceID 设备索引号 |
| received | integer | Received 接收 |
| sent | integer | Sent 发送 |

#### 3.2.5.7 获取指定设备的 PCIe 实际流量和带宽信息

函数原型：func PcieBw(dvInd int) (pcieBandwidthInfo PcieBandwidthInfo, err error)

功能说明：获取指定设备的 PCIe 实际流量和带宽信息。

参数说明：

**表 3.311 PcieBw接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.312 PcieBw接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| pcieBandwidthInfo | PcieBandwidthInfo | 函数返回值 |

#### 3.2.5.8 获取指定设备的 PCIe 带宽测试信息

函数原型：func PcieBwTest(dvInd int) (pcieBandwidthInfo PcieBandwidthInfo, err error)

功能说明：获取指定设备的 PCIe 带宽测试信息。

参数说明：

**表 3.313 PcieBwTest接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.314 PcieBwTest接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| pcieBandwidthInfo | PcieBandwidthInfo | 函数返回值 |

#### 3.2.5.9 显存带宽测试

函数原型：func BandwidthTest(dvIdList []int) bool

功能说明：显存带宽测试

参数说明：

**表 3.315 显存带宽测试接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvIdList | []int | 是 | 根据函数签名推断的参数 |

**表 3.316 显存带宽测试接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| 返回值1 | bool | 函数返回值 |

#### 3.2.5.10 PCIe带宽测试

函数原型：func PcieBandwidthTest() bool

功能说明：PCIe带宽测试

**表 3.317 PCIe带宽测试接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| 返回值1 | bool | 函数返回值 |

#### 3.2.5.11 xHCL带宽测试

函数原型：func XHCLTest() bool

功能说明：xHCL带宽测试

**表 3.318 xHCL带宽测试接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| 返回值1 | bool | 函数返回值 |

#### 3.2.5.12 获取DF带宽信息

函数原型：func DFBandwidth(dvInd int, bandwidthType int) (dfBandwidthInfo DFBandwidthInfo, err error)

功能说明：获取DF带宽信息

参数说明：

**表 3.319 获取DF带宽信息接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |
| bandwidthType | int | 是 | 根据函数签名推断的参数 |

**表 3.320 获取DF带宽信息接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| dfBandwidthInfo | DFBandwidthInfo | 函数返回值 |

#### 3.2.5.13 获取指定设备指定通道的 UMC 带宽信息

函数原型：func UMCBandwidth(dvInd int, chanId int, delay int) (info UMCBandwidthInfo, err error)

功能说明：获取指定设备指定通道的 UMC 带宽信息。

参数说明：

**表 3.321 UMCBandwidth接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |
| chanId | int | 是 | 根据函数签名推断的参数 |
| delay | int | 是 | 根据函数签名推断的参数 |

**表 3.322 UMCBandwidth接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| info | UMCBandwidthInfo | 函数返回值 |

#### 3.2.5.14 获取指定物理设备的 XHCL 带宽信息

函数原型：func XHCLBandwidth(dvInd int, linkId int, direction int, delay int) (info XhclBandwidthInfo, err error)

功能说明：获取指定物理设备的 XHCL 带宽信息

参数说明：

**表 3.323 获取指定物理设备的 XHCL 带宽信息接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |
| linkId | int | 是 | 根据函数签名推断的参数 |
| direction | int | 是 | 根据函数签名推断的参数 |
| delay | int | 是 | 根据函数签名推断的参数 |

**表 3.324 获取指定物理设备的 XHCL 带宽信息接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| info | XhclBandwidthInfo | 函数返回值 |

#### 3.2.5.15 查询所有监控设备的 UMC 带宽

函数原型：func GetHyUMCStatus() ([]DeviceUmcSum, error)

功能说明：查询所有监控设备的 UMC 带宽

**表 3.325 查询所有监控设备的 UMC 带宽接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| 返回值1 | []DeviceUmcSum | 函数返回值 |

### 3.2.6 Performance

#### 3.2.6.1 设置 GPU 时钟频率

函数原型：DevGpuClkFreqSet(dvInd int, clkType RSMIClkType, freqBitmask int64) (err error)

功能说明：设置 GPU 时钟频率

参数说明：

**表 3.63 设置 GPU 时钟频率接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | integer | Y | 设备索引 |
| freqBitmask | integer | Y | 频率掩码 |

#### 3.2.6.2 获取设备的当前性能水平

函数原型：PerfLevel(dvInd int) (perf string, err error)

功能说明：获取设备的当前性能水平

参数说明：

**表 3.64 获取设备的当前性能水平接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | integer | Y | 设备索引 |

**表 3.65 获取设备的当前性能水平接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| perf | string | 性能水平 |

#### 3.2.6.3 重置设备时钟

函数原型：ResetClocks(dvIdList []int) (failedMessage []FailedMessage)

功能说明：重置设备时钟

参数说明：

**表 3.66 重置设备时钟接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvIdList | array | Y | 设备ID列表 |

**表 3.67 重置设备时钟接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| errorMsg | string | ErrorMsg 错误信息 |
| id | integer | ID 设备ID |

#### 3.2.6.4 重置Performance Determinism

函数原型：ResetPerfDeterminism(dvIdList []int) (failedMessage []FailedMessage)

功能说明：重置Performance Determinism

参数说明：

**表 3.68 重置Performance Determinism接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvIdList | array | Y | 设备ID列表 |

**表 3.69 重置Performance Determinism接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| errorMsg | string | ErrorMsg 错误信息 |
| id | integer | ID 设备ID |

#### 3.2.6.5 重置指定设备的电源配置文件和性能级别

函数原型：ResetProfile(dvIdList []int) (failedMessage []FailedMessage)

功能说明：重置指定设备的电源配置文件和性能级别

参数说明：

**表 3.70 重置指定设备的电源配置文件和性能级别接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvIdList | array | Y | 设备ID列表 |

**表 3.71 重置指定设备的电源配置文件和性能级别接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| errorMsg | string | ErrorMsg 错误信息 |
| id | integer | ID 设备ID |

#### 3.2.6.6 设置设备的时钟OverDrive

函数原型：SetClockOverDrive(dvIdList []int, clktype string, value string, autoRespond bool) (failedMessage []FailedMessage)

功能说明：设置设备的时钟OverDrive

参数说明：

**表 3.72 设置设备的时钟OverDrive接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvIdList | array | Y | 设备ID列表 |
| clktype | string | Y | 时钟类型（sclk 或 mclk） |
| value | string | Y | OverDrive值，表示为百分比（0-20%） |
| autoRespond | boolean | N | 自动响应超出规格的警告 |

**表 3.73 设置设备的时钟OverDrive接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| errorMsg | string | ErrorMsg 错误信息 |
| id | integer | ID 设备ID |

#### 3.2.6.7 设置设备的时钟频率范围

函数原型：SetClockRange(dvIdList []int, clkType string, minvalue string, maxvalue string, autoRespond bool) (failedMessage []FailedMessage)

功能说明：设置设备的时钟频率范围

参数说明：

**表 3.74 设置设备的时钟频率范围接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvIdList | array | Y | 设备ID列表 |
| clkType | string | Y | 时钟类型（sclk 或 mclk） |
| minvalue | string | Y | 最小值（MHz） |
| maxvalue | string | Y | 最大值（MHz） |
| autoRespond | boolean | N | 自动响应超出规格的警告 |

**表 3.75 设置设备的时钟频率范围接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| errorMsg | string | ErrorMsg 错误信息 |
| id | integer | ID 设备ID |

#### 3.2.6.8 设置设备的性能确定性

函数原型：SetPerfDeterminism(dvIdList []int, clkvalue string) (failedMessage []FailedMessage, err error)

功能说明：设置设备的性能确定性

参数说明：

**表 3.76 设置设备的性能确定性接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvIdList | array | Y | 设备ID列表 |

**表 3.77 设置设备的性能确定性接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| errorMsg | string | ErrorMsg 错误信息 |
| id | integer | ID 设备ID |

#### 3.2.6.9 设置设备性能等级

函数原型：SetPerformanceLevel(dvIdList []int, level string) (failedMessages []FailedMessage)

功能说明：设置设备性能等级

参数说明：

**表 3.78 设置设备性能等级接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvIdList | array | Y | 设备ID列表 |

**表 3.79 设置设备性能等级接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| errorMsg | string | ErrorMsg 错误信息 |
| id | integer | ID 设备ID |

#### 3.2.6.10 设置设备的PowerPlay表级别

函数原型：SetPowerPlayTableLevel(dvIdList []int, clkType string, point string, clk string, volt string, autoRespond bool) (failedMessage []FailedMessage)

功能说明：设置设备的PowerPlay表级别

参数说明：

**表 3.80 设置设备的PowerPlay表级别接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvIdList | array | Y | 设备ID列表 |
| clkType | string | Y | 时钟类型（sclk 或 mclk） |
| point | string | Y | 电压点 |
| clk | string | Y | 时钟值（MHz） |
| volt | string | Y | 电压值（mV） |
| autoRespond | boolean | N | 自动响应超出规格的警告 |

**表 3.81 设置设备的PowerPlay表级别接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| errorMsg | string | ErrorMsg 错误信息 |
| id | integer | ID 设备ID |

#### 3.2.6.11 显示时钟信息

函数原型：ShowClocks(dvIdList []int)

功能说明：显示时钟信息

参数说明：

**表 3.82 显示时钟信息接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvIdList | array | Y | 设备ID列表 |

#### 3.2.6.12 获取设备粗粒度利用率

函数原型：GetCoarseGrainUtil(dvInd int, typeName string) (utilizationCounters []RSMIUtilizationCounter, err error)

功能说明：获取设备粗粒度利用率

参数说明：

**表 3.83 获取设备粗粒度利用率接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | integer | Y | 设备索引 |
| typeName | string | N | 利用率计数器类型名称 |

**表 3.84 获取设备粗粒度利用率接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| type | integer | Type 表示利用率计数器的类型 |
| value | integer | Value 表示计数器的值 |

#### 3.2.6.13 根据设备索引返回当前HCU的使用百分比

函数原型：func HCUUse(dvInd int) (percent int, err error)

功能说明：根据设备索引返回当前HCU的使用百分比

参数说明：

**表 3.326 根据设备索引返回当前HCU的使用百分比接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.327 根据设备索引返回当前HCU的使用百分比接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| percent | int | 函数返回值 |

#### 3.2.6.14 获取指定设备当前的 HCU（GPU）时钟频率

函数原型：func HCUClk(dvInd int) (clk float64, err error)

功能说明：获取指定设备当前的 HCU（GPU）时钟频率。

参数说明：

**表 3.328 HCUClk接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.329 HCUClk接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| clk | float64 | 函数返回值 |

#### 3.2.6.15 获取指定设备当前的显存时钟频率

函数原型：func HCUMclk(dvInd int) (clk float64, err error)

功能说明：获取指定设备当前的显存时钟频率。

参数说明：

**表 3.330 HCUMclk接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.331 HCUMclk接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| clk | float64 | 函数返回值 |

#### 3.2.6.16 设置设备的性能确定性模式(K100 AI不支持)

函数原型：func PerfDeterminismMode(dvInd int, clkValue int64) (err error)

功能说明：设置设备的性能确定性模式(K100 AI不支持)

参数说明：

**表 3.332 设置设备的性能确定性模式(K100 AI不支持)接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |
| clkValue | int64 | 是 | 根据函数签名推断的参数 |

#### 3.2.6.17 根据设备索引和时钟类型获取支持的频率列表及当前频率

函数原型：func GetClocksByType(dvInd int, clkType RSMIClkType) (frequencyList []uint64, currentFrequency uint32, err error)

功能说明：根据设备索引和时钟类型获取支持的频率列表及当前频率。

参数说明：

**表 3.333 GetClocksByType接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |
| clkType | RSMIClkType | 是 | 根据函数签名推断的参数 |

**表 3.334 GetClocksByType接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| frequencyList | []uint64 | 函数返回值 |
| currentFrequency | uint32 | 函数返回值 |

#### 3.2.6.18 HCU使用率

函数原型：func ShowHCUUse(dvIdList []int) (deviceUseInfos []DeviceUseInfo, err error)

功能说明：HCU使用率

参数说明：

**表 3.335 HCU使用率接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvIdList | []int | 是 | 根据函数签名推断的参数 |

**表 3.336 HCU使用率接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| deviceUseInfos | []DeviceUseInfo | 函数返回值 |

#### 3.2.6.19 显示指定设备的有效电流或电压范围

函数原型：func ShowRange(dvIdList []int, rangeType string)

功能说明：获取并显示指定设备的有效电流或电压范围。

参数说明：

**表 3.337 显示指定设备的有效电流或电压范围接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvIdList | []int | 是 | 根据函数签名推断的参数 |
| rangeType | string | 是 | 根据函数签名推断的参数 |

#### 3.2.6.20 查询指定设备当前的 HCU 占用率（瞬时值）

函数原型：func HCUCuUsage(dvInd int) (utilizationRate float64, err error)

功能说明：查询指定设备当前的 HCU 占用率（瞬时值）。只要某个 CU 内存在至少一个活跃 wave，即认为该 CU 活跃；占用率 = 活跃 CU 数量 / CU 总数量 × 100%。对应 hy-smi -u。

参数说明：

**表 3.338 查询指定设备当前的 HCU 占用率（瞬时值）接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.339 查询指定设备当前的 HCU 占用率（瞬时值）接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| utilizationRate | float64 | 函数返回值 |

#### 3.2.6.21 在 sampleDurationMs 时间窗口内周期性采样 HCU 活跃状态，返回活跃采样次数占比

函数原型：func HCUSampledUsage(dvInd int, sampleDurationMs int) (utilizationRate float64, err error)

功能说明：在 sampleDurationMs 时间窗口内周期性采样 HCU 活跃状态，返回活跃采样次数占比。对应 hy-smi --showhcuutil ；hy-smi 默认采样窗口为 1s（1000ms），采样间隔约 10ms。

参数说明：

**表 3.340 在 sampleDurationMs 时间窗口内周期性采样 HCU 活跃状态，返回活跃采样次数占比接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |
| sampleDurationMs | int | 是 | 根据函数签名推断的参数 |

**表 3.341 在 sampleDurationMs 时间窗口内周期性采样 HCU 活跃状态，返回活跃采样次数占比接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| utilizationRate | float64 | 函数返回值 |

#### 3.2.6.22 在 sampleDurationMs 窗口内统计各 CU 活跃采样占比并取全 CU 平均值

函数原型：func HCUCUSampledUsage(dvInd int, sampleDurationMs int) (utilizationRate float64, err error)

功能说明：在 sampleDurationMs 窗口内统计各 CU 活跃采样占比并取全 CU 平均值。只要 CU 内存在至少一个活跃 wave 即视为活跃。对应 hy-smi --showcuutil 。

参数说明：

**表 3.342 在 sampleDurationMs 窗口内统计各 CU 活跃采样占比并取全 CU 平均值接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |
| sampleDurationMs | int | 是 | 根据函数签名推断的参数 |

**表 3.343 在 sampleDurationMs 窗口内统计各 CU 活跃采样占比并取全 CU 平均值接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| utilizationRate | float64 | 函数返回值 |

#### 3.2.6.23 在 sampleDurationMs 窗口内统计各 CU 上活跃 wave 占 wave 总数的比例

函数原型：func HCUWaveSampledUsage(dvInd int, sampleDurationMs int) (utilizationRate float64, err error)

功能说明：在 sampleDurationMs 窗口内统计各 CU 上活跃 wave 占 wave 总数的比例，并对所有 CU 取平均。对应 hy-smi --showwaveutil 。

参数说明：

**表 3.344 在 sampleDurationMs 窗口内统计各 CU 上活跃 wave 占 wave 总数的比例接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |
| sampleDurationMs | int | 是 | 根据函数签名推断的参数 |

**表 3.345 在 sampleDurationMs 窗口内统计各 CU 上活跃 wave 占 wave 总数的比例接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| utilizationRate | float64 | 函数返回值 |

#### 3.2.6.24 按 SE 维度返回活跃 CU 占比（瞬时值）

函数原型：func HCUSEUsage(dvInd int) (shaderEngineUsage SEUsageInfo, err error)

功能说明：按 SE 维度返回活跃 CU 占比（瞬时值）。只要 CU 内存在至少一个活跃 wave 即视为活跃；每个 SE 占用率 = SE 内活跃 CU 数 / 该 SE 内 CU 总数 × 100%。对应 hy-smi --showseuse 。

参数说明：

**表 3.346 按 SE 维度返回活跃 CU 占比（瞬时值接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.347 按 SE 维度返回活跃 CU 占比（瞬时值）接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| shaderEngineUsage | SEUsageInfo | 函数返回值 |

#### 3.2.6.25 target stress-Gemm压力测试

函数原型：func TargetStressTest()

功能说明：target stress-Gemm压力测试

#### 3.2.6.26 测试板卡稳定性

函数原型：func EDPpTest()

功能说明：测试板卡稳定性

#### 3.2.6.27 对外接口：运行 EDPp 测试并以结构化形式返回结果

函数原型：func EDPpTestResult(dvIdList ...int) (EDPPResult, error)

功能说明：对外接口：运行 EDPp 测试并以结构化形式返回结果。

参数说明：

**表 3.348 对外接口：运行 EDPp 测试并以结构化形式返回结果接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvIdList | ...int | 是 | 根据函数签名推断的参数 |

**表 3.349 对外接口：运行 EDPp 测试并以结构化形式返回结果接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| 返回值1 | EDPPResult | 函数返回值 |

#### 3.2.6.28 获取指定 GPU 实例的计算实例配置文件信息

函数原型：func GpuInstanceGetComputeInstanceProfileInfo(dvInd int, giId uint32, profileId uint32, engProfileId uint32) (info NvmlComputeInstanceProfileInfo, err error)

功能说明：获取指定 GPU 实例的计算实例配置文件信息。

参数说明：

**表 3.350 GpuInstanceGetComputeInstanceProfileInfo接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |
| giId | uint32 | 是 | 根据函数签名推断的参数 |
| profileId | uint32 | 是 | 根据函数签名推断的参数 |
| engProfileId | uint32 | 是 | 根据函数签名推断的参数 |

**表 3.351 GpuInstanceGetComputeInstanceProfileInfo接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| info | NvmlComputeInstanceProfileInfo | 函数返回值 |

#### 3.2.6.29 获取指定 GPU 实例配置、计算实例配置和引擎配置对应的计算实例配置文件列表

函数原型：func ComputeInstanceProfileInfoList(dvInd int, profileId, ciProfileId, engProfileId uint32) ([]NvmlComputeInstanceProfileInfo, error)

功能说明：获取指定 GPU 实例配置、计算实例配置和引擎配置对应的计算实例配置文件列表。

参数说明：

**表 3.352 ComputeInstanceProfileInfoList接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |
| profileId | — | 是 | 根据函数签名推断的参数 |
| ciProfileId | — | 是 | 根据函数签名推断的参数 |
| engProfileId | uint32 | 是 | 根据函数签名推断的参数 |

**表 3.353 ComputeInstanceProfileInfoList接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| 返回值1 | []NvmlComputeInstanceProfileInfo | 函数返回值 |

#### 3.2.6.30 读取 GFX boost 等级信息

函数原型：func GfxBoostInfo(dvInd int) (string, error)

功能说明：读取 GFX boost 等级信息。

参数说明：

**表 3.354 读取 GFX boost 等级信息接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.355 读取 GFX boost 等级信息接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| 返回值1 | string | 函数返回值 |

#### 3.2.6.31 读取 GFX 频率限制历史（frequency dynamic range history）

函数原型：func FdrHistory(dvInd int) (string, error)

功能说明：读取 GFX 频率限制历史（frequency dynamic range history）。

参数说明：

**表 3.356 读取 GFX 频率限制历史（frequency dynamic range history）接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.357 读取 GFX 频率限制历史（frequency dynamic range history）接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| 返回值1 | string | 函数返回值 |

#### 3.2.6.32 设置设备 GFX boost 等级，具体取值范围由驱动校验

函数原型：func SetGfxBoostLevel(dvInd, level int) error

功能说明：设置设备 GFX boost 等级，具体取值范围由驱动校验。

参数说明：

**表 3.358 设置设备 GFX boost 等级，具体取值范围由驱动校验接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | — | 是 | 根据函数签名推断的参数 |
| level | int | 是 | 根据函数签名推断的参数 |

### 3.2.7 Physical State

#### 3.2.7.1 获取设备的风扇速度

函数原型：DevFanRpms(dvInd int) (speed int64, err error)

功能说明：获取设备的风扇速度

参数说明：

**表 3.85 获取设备的风扇速度接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | integer | Y | 设备索引 |

**表 3.86 获取设备的风扇速度接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| speed | integer | 风扇转速 |

#### 3.2.7.2 复位风扇控制

函数原型：ResetFans(dvIdList []int) (err error)

功能说明：复位风扇控制

参数说明：

**表 3.87 复位风扇控制接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvIdList | array | Y | 设备ID列表 |

#### 3.2.7.3 设置风扇转速

函数原型：SetFanSpeed(dvIdList []int, fan string)

功能说明：设置风扇转速

参数说明：

**表 3.88 设置风扇转速接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvIdList | array | Y | 设备ID列表 |
| fan | string | Y | 风扇速度值（0-255,单位:RPM） |

#### 3.2.7.4 获取设备温度

函数原型：Temperature(dvInd int) (temp float64, err error)

功能说明：获取设备温度

参数说明：

**表 3.89 获取设备温度接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | integer | Y | 设备索引 |

**表 3.90 获取设备温度接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| temp | float64 | 设备温度（摄氏度） |

#### 3.2.7.5 获取风扇转速

函数原型：FanSpeedInfo(dvInd int) (fanLevel int64, fanPercentage float64, err error)

功能说明：获取风扇转速

参数说明：

**表 3.91 获取风扇转速接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | integer | Y | 设备索引 |

**表 3.92 获取风扇转速接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| fanLevel | integer | 当前风扇转速 |
| fanPercentage | integer | 风扇转速百分比 |

#### 3.2.7.6 显示设备温度传感器数据

函数原型：ShowCurrentTemps(dvIdList []int) (temperatureInfos []TemperatureInfo, err error)

功能说明：显示设备温度传感器数据

参数说明：

**表 3.93 显示设备温度传感器数据接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvIdList | array | Y | 设备ID列表 |

**表 3.94 显示设备温度传感器数据接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| deviceID | integer | DeviceID 设备索引号 |
| sensorTemps | object | SensorTemps 传感器名称到温度的映射 |

#### 3.2.7.7 本函数转换为摄氏度并保留两位小数

函数原型：func GetTempBySensor(dvInd int, sensorType int, metric RSMITemperatureMetric) (temp float64, err error)

功能说明：本函数转换为摄氏度并保留两位小数。

参数说明：

**表 3.359 本函数转换为摄氏度并保留两位小数接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |
| sensorType | int | 是 | 根据函数签名推断的参数 |
| metric | RSMITemperatureMetric | 是 | 根据函数签名推断的参数 |

**表 3.360 本函数转换为摄氏度并保留两位小数接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| temp | float64 | 函数返回值 |

#### 3.2.7.8 返回板卡边缘温度（摄氏度）

函数原型：func BoardTemperature(dvInd int) (float64, error)

功能说明：返回板卡边缘温度（摄氏度）

参数说明：

**表 3.361 返回板卡边缘温度（摄氏度）接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.362 返回板卡边缘温度（摄氏度）接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| 返回值1 | float64 | 函数返回值 |

#### 3.2.7.9 显示指定设备的当前风扇转速和风扇级别

函数原型：func ShowCurrentFans(dvIdList []int, printJSON bool)

功能说明：显示指定设备的当前风扇转速和风扇级别

参数说明：

**表 3.363 显示指定设备的当前风扇转速和风扇级别接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvIdList | []int | 是 | 根据函数签名推断的参数 |
| printJSON | bool | 是 | 根据函数签名推断的参数 |

### 3.2.8 Power

#### 3.2.8.1 设置设备 PowerPlay 性能级别

函数原型：DevPerfLevelSet(dvInd int, level RSMIDevPerfLevel) error

功能说明：设置设备 PowerPlay 性能级别

参数说明：

**表 3.95 设置设备 PowerPlay 性能级别接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | integer | Y | 设备索引 |
| level | string | Y | 要设置的性能级别 |

#### 3.2.8.2 设置设备功率配置文件

函数原型：DevPowerProfileSet(dvInd int, reserved int, profile RSNIPowerProfilePresetMasks) (err error)

功能说明：设置设备功率配置文件

参数说明：

**表 3.96 设置设备功率配置文件接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | integer | Y | 设备索引 |
| reserved | integer | Y | 保留参数，通常为0 |
| profile | integer | Y | 功率配置文件的枚举值 |

#### 3.2.8.3 获取设备的最大功率

函数原型：MaxPower(dvInd int) (power int64, err error)

功能说明：获取设备的最大功率

参数说明：

**表 3.97 获取设备的最大功率接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | integer | Y | 设备索引 |

**表 3.98 获取设备的最大功率接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| PowerMax | integer | 最大功率 |

#### 3.2.8.4 获取设备的平均功耗

函数原型：Power(dvInd int) (power int64, err error)

功能说明：获取设备的平均功耗

参数说明：

**表 3.99 获取设备的平均功耗接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | integer | Y | 设备索引 |

**表 3.100 获取设备的平均功耗接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| power | integer | 平均功率 |

#### 3.2.8.5 设置功率配置

函数原型：SetProfile(dvIdList []int, profile string) (failedMessages []FailedMessage)

功能说明：设置功率配置

参数说明：

**表 3.101 设置功率配置接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvIdList | array | Y | 设备ID列表 |
| profile | string | Y | 功率配置文件名称:CUSTOM、VIDEO、POWER SAVING、COMPUTE、VR、3D FULL SCREEN、BOOTUP DEFAULT |

**表 3.102 设置功率配置接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| errorMsg | string | ErrorMsg 错误信息 |
| id | integer | ID 设备ID |

#### 3.2.8.6 展示设备的平均功率消耗

函数原型：ShowPower(dvIdList []int) (devicePowerInfos []DevicePowerInfo, err error)

功能说明：展示设备的平均功率消耗

参数说明：

**表 3.103 展示设备的平均功率消耗接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvIdList | array | Y | 设备ID列表 |

**表 3.104 展示设备的平均功率消耗接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| deviceID | integer | DeviceID 设备索引号 |
| power | integer | Power 设备功率 |

#### 3.2.8.7 展示设备的HCU内存时钟频率和电压

函数原型：ShowPowerPlayTable(dvIdList []int) (devicePowerPlayInfos []DevicePowerPlayInfo, err error)

功能说明：展示设备的HCU内存时钟频率和电压

参数说明：

**表 3.105 展示设备的HCU内存时钟频率和电压接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvIdList | array | Y | 设备ID列表 |

**表 3.106 展示设备的HCU内存时钟频率和电压接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| deviceID | integer | DeviceID 设备索引号 |
| od_MCLK | string | OD_MCLK MCLK |
| od_RANGE | array | OD_RANGE RANGE |
| od_SCLK | array | OD_SCLK SCLK |
| od_VDDC_CURVE | array | OD_VDDC_CURVE DDC_CURVE |

#### 3.2.8.8 显示设备的电源配置文件

函数原型：ShowProfile(dvIdList []int) (deviceProfiles []DeviceProfile, err error)

功能说明：显示设备的电源配置文件

参数说明：

**表 3.107 显示设备的电源配置文件接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvIdList | array | Y | 设备ID列表 |

**表 3.108 显示设备的电源配置文件接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| deviceID | integer | DeviceID 设备索引号 |
| profiles | array | Profiles 文件信息 |

#### 3.2.8.9 展示设备的能量消耗

函数原型：ShowEnergy(dvIdList []int)

功能说明：展示设备的能量消耗

参数说明：

**表 3.109 展示设备的能量消耗接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvIdList | array | Y | 设备ID列表 |

#### 3.2.8.10 显示设备的电压信息

函数原型：ShowVoltage(dvIdList []int) (deviceVoltageInfos []DeviceVoltageInfo, err error)

功能说明：显示设备的电压信息

参数说明：

**表 3.110 显示设备的电压信息接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvIdList | array | Y | 设备ID列表 |

**表 3.111 显示设备的电压信息接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| deviceID | integer | DeviceID 设备索引号 |
| voltage | integer | Voltage 电压信息 |

#### 3.2.8.11 显示设备的电压曲线点

函数原型：ShowVoltageCurve(dvIdList []int)

功能说明：显示设备的电压曲线点

参数说明：

**表 3.112 显示设备的电压曲线点接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvIdList | array | Y | 设备ID列表 |

#### 3.2.8.12 设置设备运行时电源管理模式。

函数原型：func SetLowPowerMode(dvInd int, enabled bool) error

功能说明：设置设备运行时电源管理模式。enabled=true 写入 auto，false 写入 on。

参数说明：

**表 3.364 设置设备运行时电源管理模式。enabled=true 写入 auto，false 写入 on接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |
| enabled | bool | 是 | 根据函数签名推断的参数 |

#### 3.2.8.13 设置设备运行时自动挂起延时，单位为毫秒

函数原型：func SetLowPowerDelay(dvInd, delayMs int) error

功能说明：设置设备运行时自动挂起延时，单位为毫秒。

参数说明：

**表 3.365 设置设备运行时自动挂起延时，单位为毫秒接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | — | 是 | 根据函数签名推断的参数 |
| delayMs | int | 是 | 根据函数签名推断的参数 |

### 3.2.9 System

#### 3.2.9.1 获取设备总线信息

函数原型：GetBus(dvInd int) (picId string, err error)

功能说明：获取设备总线信息

参数说明：

**表 3.113 获取设备总线信息接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | integer | Y | 设备索引 |

**表 3.114 获取设备总线信息接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| picId | string | 设备总线ID |

#### 3.2.9.2 获取指定PID的进程名

函数原型：ProcessName(pid int) string

功能说明：获取指定PID的进程名

参数说明：

**表 3.115 获取指定PID的进程名接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| pid | integer | Y | 进程ID |

**表 3.116 获取指定PID的进程名接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| name | string | 进程名称 |

#### 3.2.9.3 展示系统中正在运行的KFD进程信息

函数原型：ShowPids() (err error)

功能说明：展示系统中正在运行的KFD进程信息

#### 3.2.9.4 获取计算进程列表

函数原型：PidList() (pidList []string, err error)

功能说明：获取计算进程列表

#### 3.2.9.5 getPid 获取特定应用程序的进程 ID

函数原型：func PidByName(name string) (pid string, err error)

功能说明：getPid 获取特定应用程序的进程 ID

参数说明：

**表 3.366 getPid 获取特定应用程序的进程 ID接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| name | string | 是 | 根据函数签名推断的参数 |

**表 3.367 getPid 获取特定应用程序的进程 ID接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| pid | string | 函数返回值 |

#### 3.2.9.6 获取指定进程使用的 HCU 设备索引列表

函数原型：func ProcessHCU(pid int) (dvIndices []int, err error)

功能说明：获取指定进程使用的 HCU 设备索引列表。

参数说明：

**表 3.368 ProcessHCU接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| pid | int | 是 | 根据函数签名推断的参数 |

**表 3.369 ProcessHCU接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| dvIndices | []int | 函数返回值 |

#### 3.2.9.7 获取并进程信息和使用的HCU设备信息

函数原型：func ProcessHCUInfo() ([]Process, error)

功能说明：获取并进程信息和使用的HCU设备信息。

**表 3.370 获取并进程信息和使用的HCU设备信息接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| 返回值1 | []Process | 函数返回值 |

#### 3.2.9.8 根据进程 ID 获取进程名称

函数原型：func GetProcessName(pid int) string

功能说明：根据进程 ID 获取进程名称。

参数说明：

**表 3.371 GetProcessName接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| pid | int | 是 | 根据函数签名推断的参数 |

**表 3.372 GetProcessName接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| 返回值1 | string | 函数返回值 |

#### 3.2.9.9 获取系统中的计算进程信息及进程数量

函数原型：func ComputeProcessInfoGet() (processInfo []ProcessInfos, numItems int, err error)

功能说明：获取系统中的计算进程信息及进程数量。

**表 3.373 ComputeProcessInfoGet接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| processInfo | []ProcessInfos | 函数返回值 |
| numItems | int | 函数返回值 |

#### 3.2.9.10 根据PID获取进程信息

函数原型：func ProcessInfoByPid(pid uint32) (RsmiProcessInfoV2, error)

功能说明：根据PID获取进程信息

参数说明：

**表 3.374 根据PID获取进程信息接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| pid | uint32 | 是 | 根据函数签名推断的参数 |

**表 3.375 根据PID获取进程信息接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| 返回值1 | RsmiProcessInfoV2 | 函数返回值 |

#### 3.2.9.11 获取指定进程的进程信息

函数原型：func ProcessInfo(pid int) (proc ProcessInfos, err error)

功能说明：获取指定进程的进程信息

参数说明：

**表 3.376 获取指定进程的进程信息接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| pid | int | 是 | 根据函数签名推断的参数 |

**表 3.377 获取指定进程的进程信息接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| proc | ProcessInfos | 函数返回值 |

#### 3.2.9.12 获取系统当前和待生效的 MIG 模式

函数原型：func SystemMigMode() (currentMode, pendingMode int, err error)

功能说明：获取系统当前和待生效的 MIG 模式。

**表 3.378 SystemMigMode接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| 返回值1 | currentMode | 函数返回值 |
| pendingMode | int | 函数返回值 |

### 3.2.10 Topo

#### 3.2.10.1 显示 hcu 拓扑跳数

函数原型：ShowHopsTopology(dvIdList []int, printJSON bool)

功能说明：显示 HCU 拓扑跳数

参数说明：

**表 3.117 显示 hcu 拓扑跳数接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvIdList | array | Y | 设备ID列表 |

#### 3.2.10.2 显示完整的硬件拓扑信息

函数原型：ShowHwTopology(dvIdList []int)

功能说明：显示完整的硬件拓扑信息

参数说明：

**表 3.118 显示完整的硬件拓扑信息接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvIdList | array | Y | 设备ID列表 |

#### 3.2.10.3 显示 NUMA 节点信息

函数原型：ShowNumaTopology(dvIdList []int) (numaInfos []NumaInfo, err error)

功能说明：显示 NUMA 节点信息

参数说明：

**表 3.119 显示 NUMA 节点信息接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvIdList | array | Y | 设备ID列表 |

**表 3.120 显示 NUMA 节点信息接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| DeviceID | integer | 设备索引号 |
| NumaNode | integer | numaNode值 |
| NumaAffinity | integer | 关联信息 |

#### 3.2.10.4 显示 HCU 拓扑链接类型

函数原型：ShowTypeTopology(dvIdList []int, printJSON bool)

功能说明：显示 HCU 拓扑链接类型

参数说明：

**表 3.121 显示 HCU 拓扑链接类型接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvIdList | array | Y | 设备ID列表 |

#### 3.2.10.5 显示 GPU 拓扑权重

函数原型：ShowWeightTopology(dvIdList []int, printJSON bool)

功能说明：显示 HCU 拓扑权重

参数说明：

**表 3.122 显示 GPU 拓扑权重接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvIdList | array | Y | 设备ID列表 |

#### 3.2.10.6 重置设备的HSL错误状态

函数原型：func ResetHSLErr(dvIdList []int) (failedMessage []FailedMessage)

功能说明：重置设备的HSL错误状态

参数说明：

**表 3.379 重置设备的HSL错误状态接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvIdList | []int | 是 | 根据函数签名推断的参数 |

**表 3.380 重置设备的HSL错误状态接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| failedMessage | []FailedMessage | 函数返回值 |

#### 3.2.10.7 获取设备的HSL hive id

函数原型：func HSLHiveIdGet(dvInd int) (hiveId int64, err error)

功能说明：获取设备的HSL hive id

参数说明：

**表 3.381 获取设备的HSL hive id接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.382 获取设备的HSL hive id接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| hiveId | int64 | 函数返回值 |

#### 3.2.10.8 显示一组 GPU 设备的 HSL 错误状态

函数原型：func ShowHSLErr(dvIdList []int, printJSON bool)

功能说明：显示一组 GPU 设备的 HSL 错误状态。

参数说明：

**表 3.383 显示一组 GPU 设备的 HSL 错误状态接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvIdList | []int | 是 | 根据函数签名推断的参数 |
| printJSON | bool | 是 | 根据函数签名推断的参数 |

#### 3.2.10.9 查询源设备到多个目标设备之间的连接类型（如 PCIe、HSL）

函数原型：func GetTopoLinkType(srcDvInd int, dstDvIndList []int) (linkTypeList []string, err error)

功能说明：查询源设备到多个目标设备之间的连接类型（如 PCIe、HSL）。

参数说明：

**表 3.384 查询源设备到多个目标设备之间的连接类型（如 PCIe、HSL）接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| srcDvInd | int | 是 | 根据函数签名推断的参数 |
| dstDvIndList | []int | 是 | 根据函数签名推断的参数 |

**表 3.385 查询源设备到多个目标设备之间的连接类型（如 PCIe、HSL）接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| linkTypeList | []string | 函数返回值 |

#### 3.2.10.10 判断源设备与目标设备之间是否存在 HSL 直连

函数原型：func TopoIsHylink(srcDvInd, dstDvInd int) (bool, error)

功能说明：判断源设备与目标设备之间是否存在 HSL 直连。

参数说明：

**表 3.386 判断源设备与目标设备之间是否存在 HSL 直连接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| srcDvInd | — | 是 | 根据函数签名推断的参数 |
| dstDvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.387 判断源设备与目标设备之间是否存在 HSL 直连接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| 返回值1 | bool | 函数返回值 |

#### 3.2.10.11 查询指定 HCU 的所有 XHCL 链路当前状态

函数原型：func XhclLinkStates(dvInd int) ([]XhclLinkState, error)

功能说明：查询指定 HCU 的所有 XHCL 链路当前状态。

参数说明：

**表 3.388 查询指定 HCU 的所有 XHCL 链路当前状态接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.389 查询指定 HCU 的所有 XHCL 链路当前状态接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| 返回值1 | []XhclLinkState | 函数返回值 |

#### 3.2.10.12 枚举指定 GPU 上所有处于 UP 状态的 XHCL 链路，并返回每条链路所连接的远端设备 BDF ID

函数原型：func DumpXhclRemoteBdfids(dvInd int) ([]XhclRemoteBdf, error)

功能说明：枚举指定 GPU 上所有处于 UP 状态的 XHCL 链路，并返回每条链路所连接的远端设备 BDF ID。

参数说明：

**表 3.390 枚举指定 GPU 上所有处于 UP 状态的 XHCL 链路，并返回每条链路所连接的远端设备 BDF ID接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.391 枚举指定 GPU 上所有处于 UP 状态的 XHCL 链路，并返回每条链路所连接的远端设备 BDF ID接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| 返回值1 | []XhclRemoteBdf | 函数返回值 |

#### 3.2.10.13 枚举整机 HCU 互联关系，包括链路类型（PCIe / HSL / HYSWITCH / NONE）及对应权重

函数原型：func DiscoverInterconnectTopology() (matrix HCUInterconnectMatrix, err error)

功能说明：枚举整机 HCU 互联关系，包括链路类型（PCIe / HSL / HYSWITCH / NONE）及对应权重。

**表 3.392 枚举整机 HCU 互联关系，包括链路类型（PCIe / HSL / HYSWITCH / NONE）及对应权重接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| matrix | HCUInterconnectMatrix | 函数返回值 |

#### 3.2.10.14 获取指定 GPU 的互联拓扑信息

函数原型：func GetGpuInterconnectInfo(srcDvInd int, dstDvIndList []int) (info GpuInterconnectInfo, err error)

功能说明：获取指定 GPU 的互联拓扑信息

参数说明：

**表 3.393 获取指定 GPU 的互联拓扑信息接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| srcDvInd | int | 是 | 根据函数签名推断的参数 |
| dstDvIndList | []int | 是 | 根据函数签名推断的参数 |

**表 3.394 获取指定 GPU 的互联拓扑信息接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| info | GpuInterconnectInfo | 函数返回值 |

#### 3.2.10.15 GetHyLinkLinkStatus 查询所有监控设备的 XHCL link 带宽，按 device 顺序返回切片

函数原型：func GetHyLinkStatus() ([]DeviceLinkSum, error)

功能说明：GetHyLinkLinkStatus 查询所有监控设备的 XHCL link 带宽，按 device 顺序返回切片。

**表 3.395 GetHyLinkLinkStatus 查询所有监控设备的 XHCL link 带宽，按 device 顺序返回切片接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| 返回值1 | []DeviceLinkSum | 函数返回值 |

#### 3.2.10.16 查询指定设备（索引 dvInd）上每个 link 的 recv/send 明细及汇总

函数原型：func HyLinkStatusByHcuId(dvInd int) (DeviceLinkSum, error)

功能说明：查询指定设备（索引 dvInd）上每个 link 的 recv/send 明细及汇总。

参数说明：

**表 3.396 查询指定设备（索引 dvInd）上每个 link 的 recv/send 明细及汇总接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.397 查询指定设备（索引 dvInd）上每个 link 的 recv/send 明细及汇总接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| 返回值1 | DeviceLinkSum | 函数返回值 |

#### 3.2.10.17 查询两设备之间是否支持 P2P 访问

函数原型：func IsP2PAccessible(srcDvInd, dstDvInd int) (bool, error)

功能说明：查询两设备之间是否支持 P2P 访问

参数说明：

**表 3.398 查询两设备之间是否支持 P2P 访问接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| srcDvInd | — | 是 | 根据函数签名推断的参数 |
| dstDvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.399 查询两设备之间是否支持 P2P 访问接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| 返回值1 | bool | 函数返回值 |

### 3.2.11 VDevice

#### 3.2.11.1 创建虚拟设备

函数原型：CreateVDevices(dvInd int, vDevCount int, vDevCUs []int, vDevMemSize []int) (vdevIDs []int, err error)

功能说明：创建虚拟设备

参数说明：

**表 3.123 创建虚拟设备接口参数说明**

| 物理参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | integer | Y | 设备索引 |
| vDevCount | integer | Y | 要创建的虚拟设备数量 |
| vDevCUs | array | Y | 每个虚拟设备的计算单元数量，多个值使用多个 vDevCUs 参数传递，例如：vDevCUs=10&vDevCUs=20 |
| vDevMemSize | array | Y | 每个虚拟设备的内存大小，多个值使用多个 vDevMemSize 参数传递，例如：vDevMemSize=1024&vDevMemSize=2048 |

#### 3.2.11.2 销毁单个虚拟设备

函数原型：DestroySingleVDevice(vDvInd int) (err error)

功能说明：销毁单个虚拟设备

参数说明：

**表 3.124 销毁单个虚拟设备接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| vDvInd | integer | Y | 虚拟设备索引 |

#### 3.2.11.3 销毁所有虚拟设备

函数原型：DestroyVDevice(dvInd int) (err error)

功能说明：销毁所有虚拟设备

参数说明：

**表 3.125 销毁所有虚拟设备接口参数说明**

| 物理参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | integer | Y | 设备索引 |

#### 3.2.11.4 获取当前虚拟机的加密状态

函数原型：EncryptionVMStatus() (status bool, err error)

功能说明：获取当前虚拟机的加密状态

#### 3.2.11.5 设置虚拟机加密状态

函数原型：SetEncryptionVMStatus(status bool)

功能说明：设置虚拟机加密状态

参数说明：

**表 3.126 设置虚拟机加密状态接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| status | boolean | Y | 加密状态 |

#### 3.2.11.6 启动指定的虚拟设备

函数原型：StartVDevice(vDvInd int) (err error)

功能说明：启动指定的虚拟设备

参数说明：

**表 3.127 启动指定的虚拟设备接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| vDvInd | integer | Y | 虚拟设备索引 |

#### 3.2.11.7 停止指定的虚拟设备

函数原型：StopVDevice(vDvInd int) (err error)

功能说明：停止指定的虚拟设备

参数说明：

**表 3.128 停止指定的虚拟设备接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| vDvInd | integer | Y | 虚拟设备索引 |

#### 3.2.11.8 更新虚拟设备资源

函数原型：UpdateSingleVDevice(vDvInd int, vDevCUs int, vDevMemSize int) (err error)

功能说明：更新虚拟设备资源

参数说明：

**表 3.129 更新虚拟设备资源接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| vDvInd | integer | Y | 虚拟设备索引 |
| vDevCUs | integer | Y | 更新后的计算单元数量 |
| vDevMemSize | integer | Y | 更新后的内存大小 |

#### 3.2.11.9 获取单个虚拟设备的信息

函数原型：VDeviceSingleInfo(vDvInd int) (vDeviceInfo DMIVDeviceInfo, err error)

功能说明：获取单个虚拟设备的信息

参数说明：

**表 3.130 获取单个虚拟设备的信息接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| vDvInd | integer | Y | 虚拟设备索引 |

**表 3.131 获取单个虚拟设备的信息接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| DvInd | integer | DvInd物理设备的设备ID |
| VMemoryTotal | integer | VMemoryTotal虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed虚拟设备的已使用内存大小 |
| VdvInd | integer | VdvInd虚拟设备的索引号 |

#### 3.2.11.10 获取设备剩余信息

函数原型：DeviceRemainingInfo(dvInd int) (cus, memories uint64, err error)

功能说明：获取设备剩余信息

参数说明：

**表 3.132 获取设备剩余信息接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | integer | Y | 物理设备索引 |

**表 3.133 获取设备剩余信息接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| cus | uint64 | CU剩余数量 |
| memories | uint64 | 内存剩余量 |

#### 3.2.11.11 获取虚拟设备数量

函数原型：VDeviceCount() (count int, err error)

功能说明：获取虚拟设备数量

参数说明：

**表 3.134 获取虚拟设备数量接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| count | integer | 虚拟卡总数量 |

#### 3.2.11.12 获取所有虚拟卡信息

函数原型：VDeviceInfos() (vDeviceInfos []VDeviceInfo, err error)

功能说明：获取所有虚拟卡信息

参数说明：

**表 3.135 获取所有虚拟卡信息接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| vDvInd | integer | Y | 虚拟设备索引 |

**表 3.136 获取所有虚拟卡信息接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| DvInd | integer | DvInd物理设备的设备ID |
| VMemoryTotal | integer | VMemoryTotal虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed虚拟设备的已使用内存大小 |
| VdvInd | integer | VdvInd虚拟设备的索引号 |

#### 3.2.11.13 获取指定 GPU 实例的计算实例剩余容量

函数原型：func GpuInstanceGetComputeInstanceRemainingCapacity(dvInd int, giId uint32, profileId uint32) (uint32, error)

功能说明：获取指定 GPU 实例的计算实例剩余容量。

参数说明：

**表 3.400 GpuInstanceGetComputeInstanceRemainingCapacity接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |
| giId | uint32 | 是 | 根据函数签名推断的参数 |
| profileId | uint32 | 是 | 根据函数签名推断的参数 |

**表 3.401 GpuInstanceGetComputeInstanceRemainingCapacity接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| 返回值1 | uint32 | 函数返回值 |

#### 3.2.11.14 获取指定设备 GPU 实例配置文件对应的所有计算实例剩余容量

函数原型：func GpuInstancesComputeInstanceRemainingCapacityList(dvInd int, profileId uint32) ([]ComputeInstanceRemainInfo, error)

功能说明：获取指定设备 GPU 实例配置文件对应的所有计算实例剩余容量。

参数说明：

**表 3.402 GpuInstancesComputeInstanceRemainingCapacityList接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |
| profileId | uint32 | 是 | 根据函数签名推断的参数 |

**表 3.403 GpuInstancesComputeInstanceRemainingCapacityList接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| 返回值1 | []ComputeInstanceRemainInfo | 函数返回值 |

#### 3.2.11.15 获取当前主机上所有GPU的所有MIG分区信息

函数原型：func MigInfos() ([]MigInfo, error)

功能说明：获取当前主机上所有GPU的所有MIG分区信息。

**表 3.404 获取当前主机上所有GPU的所有MIG分区信息接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| 返回值1 | []MigInfo | 函数返回值 |

#### 3.2.11.16 遍历所有可用的GPU设备，依次收集每个GPU上所有MIG设备的属性、GPU实例ID和计算实例ID

函数原型：func MigInfoByDvInd(dvInd int) ([]MigInfo, error)

功能说明：遍历所有可用的GPU设备，依次收集每个GPU上所有MIG设备的属性、GPU实例ID和计算实例ID。

参数说明：

**表 3.405 遍历所有可用的GPU设备，依次收集每个GPU上所有MIG设备的属性、GPU实例ID和计算实例ID接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.406 遍历所有可用的GPU设备，依次收集每个GPU上所有MIG设备的属性、GPU实例ID和计算实例ID接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| 返回值1 | []MigInfo | 函数返回值 |

#### 3.2.11.17 根据设备UUID获取该GPU的MIG分区信息

函数原型：func MigInfoByUUID(uuid string) (migInfo MigInfo, err error)

功能说明：根据设备UUID获取该GPU的MIG分区信息。

参数说明：

**表 3.407 根据设备UUID获取该GPU的MIG分区信息接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| uuid | string | 是 | 根据函数签名推断的参数 |

**表 3.408 根据设备UUID获取该GPU的MIG分区信息接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| migInfo | MigInfo | 函数返回值 |

#### 3.2.11.18 获取系统支持的 MIG 配置列表

函数原型：func MIGConfigs() ([]MigConfig, error)

功能说明：获取系统支持的 MIG 配置列表。

**表 3.409 MIGConfigs接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| 返回值1 | []MigConfig | 函数返回值 |

#### 3.2.11.19 根据 SE 数量、GPU 实例切片数、计算实例切片数和显存大小生成 MIG 名称

函数原型：func MIGName(seCount int, giSliceCount int, ciSliceCount int, memorySizeMB uint64) string

功能说明：根据 SE 数量、GPU 实例切片数、计算实例切片数和显存大小生成 MIG 名称。

参数说明：

**表 3.410 MIGName接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| seCount | int | 是 | 根据函数签名推断的参数 |
| giSliceCount | int | 是 | 根据函数签名推断的参数 |
| ciSliceCount | int | 是 | 根据函数签名推断的参数 |
| memorySizeMB | uint64 | 是 | 根据函数签名推断的参数 |

**表 3.411 MIGName接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| 返回值1 | string | 函数返回值 |

### 3.2.12 Version

#### 3.2.12.1 获取RSMI版本

函数原型：HCUVersion () (version DevVersion, err error)

功能说明：获取RSMI版本

参数说明：

**表 3.137 获取RSMI版本接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| Major | uint32 | 主版本号 |
| Minor | uint32 | 次版本号 |
| Patch | uint32 | 补丁版本号 |
| Build | string | 构建字符串，包含构建的额外信息 |

#### 3.2.12.2 获取设备的VBIOS版本

函数原型：VbiosVersion(dvInd int) (vbios string, err error)

功能说明：获取设备的VBIOS版本

参数说明：

**表 3.138 获取设备的VBIOS版本接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | integer | Y | 设备索引 |

**表 3.139 获取设备的VBIOS版本接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| vbios | string | VBIOS版本 |

#### 3.2.12.3 获取当前系统的驱动程序版本

函数原型：Version(component RSMISwComponent) (varStr string, err error)

功能说明：获取当前系统的驱动程序版本

参数说明：

**表 3.140 获取当前系统的驱动程序版本接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| component | string | Y | 驱动组件:FIRST、DRIVER、LAST |

**表 3.141 获取当前系统的驱动程序版本接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| varStr | string | 驱动版本号 |

#### 3.2.12.4 获取设备固件版本信息

函数原型：ShowFwInfo(dvIdList []int, fwType []string) (fwInfos []FirmwareInfo, err error)

功能说明：显示设备固件版本信息

参数说明：

**表 3.142 获取设备固件版本信息接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvIdList | array | Y | 设备ID列表 |
| fwType | array | Y | 固件类型列表 |

**表 3.143 获取设备固件版本信息接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| deviceID | integer | DeviceID 设备索引号 |
| firmwareVer | map | FirmwareVer 固件块名称到版本信息的映射 |

#### 3.2.12.5 获取指定设备的VBIOS版本

函数原型：ShowVbiosVersion(dvIdList []int) (deviceVBIOSInfos []DeviceVBIOSInfo, err error)

功能说明：获取指定设备固件版本信息

参数说明：

**表 3.144 获取指定设备的VBIOS版本接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvIdList | array | Y | 设备ID列表 |

**表 3.145 获取指定设备的VBIOS版本接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| deviceID | integer | DeviceID 设备索引号 |
| vbios | string | VBIOS 版本信息 |

#### 3.2.12.6 获取DTK版本

函数原型：DTKVersion() (dtkVersion string, err error)

功能说明：获取DTK版本

参数说明：

**表 3.146 获取DTK版本接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dtkVersion | string | DTK版本号 |
| Minor | uint32 | 次版本号 |
| Patch | uint32 | 补丁版本号 |
| Build | string | 构建字符串，包含构建的额外信息 |

#### 3.2.12.7 返回当前设备的RSMI版本信息

函数原型：func HCUVersion() (version DevVersion, err error)

功能说明：返回当前设备的RSMI版本信息

**表 3.412 返回当前设备的RSMI版本信息接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| version | DevVersion | 函数返回值 |

#### 3.2.12.8 根据卡型号、驱动版本和 DTK 版本，返回通过或告警结果及推荐版本

函数原型：func Compatible(cardModel, driverVersion, dtkVersion string) (CompatibilityResult, error)

功能说明：根据卡型号、驱动版本和 DTK 版本，返回通过或告警结果及推荐版本。

参数说明：

**表 3.413 根据卡型号、驱动版本和 DTK 版本，返回通过或告警结果及推荐版本接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| cardModel | — | 是 | 根据函数签名推断的参数 |
| driverVersion | — | 是 | 根据函数签名推断的参数 |
| dtkVersion | string | 是 | 根据函数签名推断的参数 |

**表 3.414 根据卡型号、驱动版本和 DTK 版本，返回通过或告警结果及推荐版本接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| 返回值1 | CompatibilityResult | 函数返回值 |

#### 3.2.12.9 校验卡型号、驱动版本和 DTK 版本的兼容性

函数原型：func CompatibleTest(cardModel, driverVersion, dtkVersion string) error

功能说明：校验卡型号、驱动版本和 DTK 版本的兼容性。

参数说明：

**表 3.415 CompatibleTest接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| cardModel | — | 是 | 根据函数签名推断的参数 |
| driverVersion | — | 是 | 根据函数签名推断的参数 |
| dtkVersion | string | 是 | 根据函数签名推断的参数 |

#### 3.2.12.10 返回海光扩展版本（hy_major/hy_minor）

函数原型：func HyVersion() (HyVersionInfo, error)

功能说明：返回海光扩展版本（hy_major/hy_minor）

**表 3.416 返回海光扩展版本（hy_major/hy_minor）接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| 返回值1 | HyVersionInfo | 函数返回值 |

#### 3.2.12.11 读取 HSW固件版本号

函数原型：func HswFwVersion(dvInd int) (string, error)

功能说明：读取 HSW固件版本号

参数说明：

**表 3.417 读取 HSW固件版本号**接口参数

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| dvInd | int | 是 | 根据函数签名推断的参数 |

**表 3.418 读取 HSW固件版本号返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| dvInd | integer | DvInd 设备的索引号 |
| pciBusNumber | string | PciBusNumber 设备的总线编号 |
| deviceId | string | DeviceId 设备的唯一标识符 |
| devTypeName | string | DevTypeName 设备的子系统名称 |
| temperature | float64 | Temperature 设备当前的温度 |
| powerUsed | float64 | PowerUsed设备当前的功耗 |
| powerTotal | float64 | PowerTotal 设备的功耗上限 |
| memoryTotal | float64 | MemoryTotal 设备的内存容量 |
| memoryUsed | float64 | MemoryUsed 设备已使用的内存 |
| utilizationRate | float64 | UtilizationRate 设备的利用率 |
| pcieBwMb | float64 | PcieBwMb 设备的PCIe带宽 |
| pcieSent | float64 | PcieSent 设备的PCIe发送数据量 |
| pcieReceived | float64 | PcieReceived 设备的PCIe接收数据量 |
| clk | float64 | Clk 设备的当前时钟频率 |
| computeUnitCount | float64 | ComputeUnitCount 设备的计算单元总数 |
| computeUnitRemainingCount | uint64 | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |
| memoryRemaining | uint64 | MemoryRemaining 设备剩余可用的内存量 |
| percent | integer | Percent 物理设备使用百分比 |
| maxVDeviceCount | integer | MaxVDeviceCount 物理设备上支持的最大虚拟设备数量 |
| VDeviceCount | integer | VDeviceCount 虚拟设备数量 |
| blocksInfos | array | BlocksInfo 设备的block信息 |
| block | string | Block 类型 |
| state | string | tate 状态 |
| CE | integer | CE错误数 |
| UE | integer | UE错误数 |
| DFBandwidthInfo | array | DF内存带宽信息 |
| readBW | float64 | ReadBW 读速度 |
| writeBW | float64 | WriteBW 写速度 |
| readWriteBW | float64 | ReadWriteBW 读写速度 |
| virtualDevices | array | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |
| VComputeUnitCount | integer | VComputeUnitCount虚拟设备的计算单元数量 |
| containerID | integer | ContainerID 虚拟设备的容器ID |
| VdvInd | integer | VdvInd 虚拟设备的设备ID |
| DvInd | integer | DvInd 物理设备的ID |
| VMemoryTotal | integer | VMemoryTotal 虚拟设备的全局内存大小 |
| name | string | Name 虚拟设备的名称 |
| pciBusNumber | string | PciBusNumber 虚拟设备的总线编号 |
| VPercent | integer | VPercent 虚拟设备的使用百分比 |
| VMemoryUsed | integer | VMemoryUsed 虚拟设备的已使用内存大小 |
| 返回值1 | string | 函数返回值 |

### 3.2.13 Group

#### 3.2.13.1 获取所有设备组信息

函数原型：ListAllGroups() (groups []GroupInfo, err error)

功能说明：获取所有设备组信息

参数说明：

**表 3.147 获取所有设备组信息接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| GroupId | integer | 设备组ID |
| GroupName | string | 设备组名称 |
| EntityGroupId | integer | 设备类型ID，1代表HCU设备，2代表vHCU设备 |
| EntityId | integer | 设备ID |

#### 3.2.13.2 获取指定设备组信息

函数原型：GetGroupInfo(groupId int) (groupInfo GroupInfo, err error)

功能说明：获取指定设备组信息

参数说明：

**表 3.148 获取指定设备组信息接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| groupId | integer | Y | 设备组ID |

**表 3.149 获取指定设备组信息接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| GroupId | integer | 设备组ID |
| GroupName | string | 设备组名称 |
| EntityGroupId | integer | 设备类型ID，1代表HCU设备，2代表vHCU设备 |
| EntityId | integer | 设备ID |

#### 3.2.13.3 创建设备组

函数原型：CreateGroup(groupName string) (groupId int, err error)

功能说明：创建设备组

参数说明：

**表 3.150 创建设备组接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| groupName | string | Y | 设备组名称 |

**表 3.151 创建设备组接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| GroupId | integer | 设备组ID |

#### 3.2.13.4 向指定设备组中添加HCU设备

函数原型：AddToGroup(groupId int, hcuIndex int) (err error)

功能说明：向指定设备组中添加HCU设备

参数说明：

**表 3.152 向指定设备组中添加HCU设备接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| groupId | integer | Y | 设备组ID |
| hcuIndex | integer | Y | HCU设备ID |

#### 3.2.13.5 向指定设备组中添加设备

函数原型：AddEntityToGroup(groupId int, entityList []GroupEntityPair) error

功能说明：向指定设备组中添加设备

参数说明：

**表 3.153 向指定设备组中添加设备接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| groupId | integer | Y | 设备组ID |
| EntityGroupId | integer | Y | 设备类型ID，1代表HCU设备，2代表vHCU设备 |
| EntityId | integer | Y | 设备ID |

#### 3.2.13.6 从设备组中移除设备

函数原型：RemoveEntityFromGroup(groupId int, entityList []GroupEntityPair) error

功能说明：从设备组中移除设备

参数说明：

**表 3.154 从设备组中移除设备接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| groupId | integer | Y | 设备组ID |
| EntityGroupId | integer | Y | 设备类型ID，1代表HCU设备，2代表vHCU设备 |
| EntityId | integer | Y | 设备ID |

#### 3.2.13.7 删除指定设备组

函数原型：DestroyGroup(groupId int) error

功能说明：删除指定设备组

参数说明：

**表 3.155 删除指定设备组接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| groupId | integer | Y | 设备组ID |

#### 3.2.13.8 创建默认设备组

函数原型：func CreateDefaultGroup(groupName string) (groupId int, err error)

功能说明：创建默认设备组。

参数说明：

**表 3.419 CreateDefaultGroup接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| groupName | string | 是 | 根据函数签名推断的参数 |

**表 3.420 CreateDefaultGroup接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| groupId | int | 函数返回的设备组 ID |

#### 3.2.13.9 获取设备组中的 HCU 列表及设备组名称

函数原型：func GetHCUListFromGroup(groupId int) (hcuIds []int, groupName string, err error)

功能说明：获取设备组中的 HCU 列表及设备组名称。

参数说明：

**表 3.421 GetHCUListFromGroup接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| groupId | int | 是 | 根据函数签名推断的参数 |

**表 3.422 GetHCUListFromGroup接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| hcuIds | []int | 设备组中的 HCU 设备索引列表 |
| groupName | string | 设备组名称 |

### 3.2.14 Field Group

#### 3.2.14.1 获取所有字段组信息

函数原型：ListAllFieldGroups() (fieldGroups []FieldGroupInfo, err error)

功能说明：获取所有字段组信息

参数说明：

**表 3.156 获取所有字段组信息接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| FieldGroupId | integer | 字段组ID |
| FieldGroupName | string | 字段组名称 |
| FieldIds | []int | 字段ID列表 |

#### 3.2.14.2 获取指定字段组信息

函数原型：GetFieldGroupInfo(fieldGroupId int) (fieldGroupInfo FieldGroupInfo, err error)

功能说明：获取指定字段组信息

参数说明：

**表 3.157 获取指定字段组信息接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| fieldGroupId | integer | Y | ID |

**表 3.158 获取指定字段组信息接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| FieldGroupId | integer | 字段组ID |
| FieldGroupName | string | 字段组名称 |
| FieldIds | []int | 字段ID列表 |

#### 3.2.14.3 创建字段组

函数原型：CreateFieldGroup(fieldGroupName string, fieldIds []int) (fieldGroupId int, err error)

功能说明：创建字段组

参数说明：

**表 3.159 创建字段组接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| fieldGroupName | string | Y | 字段组名称 |
| fieldIds | []int | Y | 字段ID列表 |

**表 3.160 创建字段组接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| fieldGroupId | integer | 字段组ID |

#### 3.2.14.4 删除指定字段组

函数原型：DestroyFieldGroup(fieldGroupId int) error

功能说明：删除指定字段组

参数说明：

**表 3.161 删除指定字段组接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| fieldGroupId | integer | Y | 字段组ID |

### 3.2.15 Watch Field

#### 3.2.15.1 获取所有字段信息

函数原型：ListFieldMeta() []FieldMeta

功能说明：获取所有字段信息

参数说明：

**表 3.162 获取所有字段信息接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| FieldId | integer | 字段ID |
| Name | string | 字段名称 |
| EntityLevel | integer | 字段类型ID |

#### 3.2.15.2 获取指定字段信息

函数原型：GetFieldMetaById(fieldId int) FieldMeta

功能说明：获取指定字段ID的字段信息

参数说明：

**表 3.163 获取指定字段信息接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| fieldId | integer | 字段ID |

#### 3.2.15.3 监测字段值

函数原型：WatchFields(hcuIndex int, fieldIds []int) error

功能说明：监测字段值

参数说明：

**表 3.164 监测字段值接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| hcuIndex | integer | Y | 设备序列 |
| fieldIds | []int | Y | 字段ID列表 |

#### 3.2.15.4 指定设备序列和字段列表监测字段值

函数原型：WatchFields(hcuIndex int, fieldIds []int) error

功能说明：指定设备序列和字段列表监测字段值

参数说明：

**表 3.165 指定设备序列和字段列表监测字段值接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| hcuIndex | integer | Y | 设备序列 |
| fieldIds | []int | Y | 字段ID列表 |

#### 3.2.15.5 指定设备序列和字段组监测字段值

函数原型：WatchFieldGroup(hcuIndex int, fieldGroupId int) error

功能说明：指定设备序列和字段组监测字段值

参数说明：

**表 3.166 指定设备序列和字段组监测字段值接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| hcuIndex | integer | Y | 设备序列 |
| fieldGroupId | integer | Y | 字段组ID |

#### 3.2.15.6 指定设备类型和字段列表监测字段值

函数原型：WatchFieldsWithEntity(entityGroup Field_Entity_Group, entityId int, fieldIds []int) error

功能说明：指定设备类型和字段列表监测字段值

参数说明：

**表 3.167 指定设备类型和字段列表监测字段值接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| entityGroup | integer | Y | 设备类型ID |
| entityId | integer | Y | 设备ID |
| fieldIds | []int | Y | 字段ID列表 |

#### 3.2.15.7 指定设备类型和字段组监测字段值

函数原型：WatchFieldGroupWithEntity(entityGroup Field_Entity_Group, entityId int, fieldGroupId int) error

功能说明：指定设备类型和字段组监测字段值

参数说明：

**表 3.168 指定设备类型和字段组监测字段值接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| entityGroup | integer | Y | 设备类型ID |
| entityId | integer | Y | 设备ID |
| fieldGroupId | integer | Y | 字段组列表 |

#### 3.2.15.8 指定设备组和字段组监测字段值

函数原型：WatchFieldsWithGroup(fieldGroupId int, groupId int) error

功能说明：指定设备组和字段组监测字段值

参数说明：

**表 3.169 指定设备组和字段组监测字段值接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| fieldGroupId | integer | Y | 字段组ID |
| groupId | integer | Y | 设备组ID |

#### 3.2.15.9 自定义频次监测字段值

函数原型：WatchFieldsWithGroupEx(fieldGroupId int, groupId int, updateFreq time.Duration, maxKeepAge time.Duration, maxKeepSamples int32) error

功能说明：自定义频次监测字段值

参数说明：

**表 3.170 自定义频次监测字段值接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| fieldGroupId | integer | Y | 字段组ID |
| groupId | integer | Y | 设备组ID |
| updateFreq | time.Duration | Y | 更新频次时间间隔 |
| maxKeepAge | time.Duration | Y | 字段值保留最长时间 |
| maxKeepSamples | time.Duration | Y | 字段值保留最大个数 |

#### 3.2.15.10 取消监测字段值

函数原型：UnWatchFields(hcuIndex int)

功能说明：取消监测字段值

参数说明：

**表 3.171 取消监测字段值接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| hcuIndex | integer | Y | 设备序列 |

#### 3.2.15.11 指定设备类型取消监测字段值

函数原型：UnWatchFieldsWithEntity(entityGroup Field_Entity_Group, entityId int)

功能说明：指定设备类型取消监测字段值

参数说明：

**表 3.172 指定设备类型取消监测字段值接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| entityGroup | integer | Y | 设备类型ID |
| entityId | integer | Y | 设备ID |

#### 3.2.15.12 指定设备组取消监测字段值

函数原型：UnWatchFieldsWithGroup(groupId int) error

功能说明：指定设备组取消监测字段值

参数说明：

**表 3.173 指定设备组取消监测字段值接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| groupId | integer | Y | 设备组ID |

#### 3.2.15.13 获取最新字段值

函数原型：GetLatestValuesForFields(hcuIndex int, fields []int) ([]FieldValue_v1, error)

功能说明：获取最新字段值

参数说明：

**表 3.174 获取最新字段值接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| hcuIndex | integer | Y | 设备序列 |
| fields | []int | Y | 字段ID列表 |

**表 3.175 获取最新字段值接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| FieldId | integer | 字段ID |
| Timestamp | int64 | 时间戳 |
| Value | float | 字段值 |
| Err | error | 报错信息 |
| Tag | string | 额外注解信息 |

#### 3.2.15.14 指定设备类型获取最新字段值

函数原型：EntityGetLatestValues(entityGroup Field_Entity_Group, entityId int, fields []int) ([]FieldValue_v1, error)

功能说明：指定设备类型获取最新字段值

参数说明：

**表 3.176 指定设备类型获取最新字段值接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| entityGroup | integer | Y | 设备类型 |
| entityId | integer | Y | 设备ID |
| fields | []int | Y | 字段ID列表 |

**表 3.177 指定设备类型获取最新字段值接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| FieldId | integer | 字段ID |
| Timestamp | int64 | 时间戳 |
| Value | float | 字段值 |
| Err | error | 报错信息 |
| Tag | string | 额外注解信息 |

#### 3.2.15.15 指定多个设备获取最新字段值

函数原型：EntitiesGetLatestValues(entities []GroupEntityPair, fields []int) ([]FieldValue_v2, error)

功能说明：指定多个设备获取最新字段值

参数说明：

**表 3.178 指定多个设备获取最新字段值接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| entityGroup | integer | Y | 设备类型 |
| entityId | integer | Y | 设备ID |
| fields | []int | Y | 字段ID列表 |

**表 3.179 指定多个设备获取最新字段值接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| EntityGroupId | integer | 设备类型 |
| EntityId | integer | 设备ID |
| FieldId | integer | 字段ID |
| Timestamp | int64 | 时间戳 |
| Value | float | 字段值 |
| Err | error | 报错信息 |
| Tag | string | 额外注解信息 |

#### 3.2.15.16 按指定采样参数监测设备组的字段值

函数原型：func WatchFieldsWithEntityGroup(fieldIds []int, groupId int, updateFreq time.Duration, maxKeepAge time.Duration, maxKeepSamples int32) error

功能说明：按指定采样参数监测设备组的字段值。

参数说明：

**表 3.423 WatchFieldsWithEntityGroup接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| fieldIds | []int | 是 | 需要监测的字段 ID 列表 |
| groupId | int | 是 | 设备组 ID |
| updateFreq | time.Duration | 是 | 字段值更新频率 |
| maxKeepAge | time.Duration | 是 | 字段值最大保留时间 |
| maxKeepSamples | int32 | 是 | 字段值最大保留样本数 |

### 3.2.16 Policy

#### 3.2.16.1 获取策略

函数原型：GetPolicy(hcuList []int) (policyList []Policy, err error)

功能说明：获取指定HCU设备列表的已设定策略

参数说明：

**表 3.423 WatchFieldsWithEntityGroup接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| fieldIds | []int | 是 | 需要监测的字段 ID 列表 |
| groupId | int | 是 | 设备组 ID |
| updateFreq | time.Duration | 是 | 字段值更新频率 |
| maxKeepAge | time.Duration | 是 | 字段值最大保留时间 |
| maxKeepSamples | int32 | 是 | 字段值最大保留样本数 |

**表 3.180 获取策略接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| hcuList | []int | Y | 设备ID列表 |

**表 3.181 获取策略接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| HcuIndex | integer | 设备序列 |
| ActionIndex | integer | 执行动作序列 |
| ValidationIndex | integer | 执行动作后验证序列 |
| Conditions | PolicyConditions | 策略执行条件集合 |

#### 3.2.16.2 设置策略

函数原型： SetPolicy(policyInfo Policy, hcuIndex int)

功能说明：为指定HCU设备设定策略

参数说明：

**表 3.182 设置策略接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| HcuIndex | integer | Y | 设备序列 |
| ActionIndex | integer | Y | 执行动作序列 |
| ValidationIndex | integer | N | 执行动作后验证序列 |
| Conditions | PolicyConditions | Y | 策略执行条件集合 |

#### 3.2.16.3 清除策略

函数原型：ClearPolicy(hcuList []int) error

功能说明：清除指定HCU设备列表的已设定策略

参数说明：

**表 3.183 清除策略接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| hcuList | []int | Y | 设备ID列表 |

#### 3.2.16.4 判定策略条件

函数原型：JudgePolicyConditions(hcuList []int) (hcuIndex int, err error)

功能说明：判断指定HCU设备列表的策略条件是否已经满足

参数说明：

**表 3.184 判定策略条件接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| hcuList | []int | Y | 设备ID列表 |

**表 3.185 判定策略条件接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| hcuIndex | integer | 到达条件的设备序列 |

#### 3.2.16.5 执行动作

函数原型：TakePolicyAction(hcuIndex int) error

功能说明：策略条件满足后调用该接口执行动作

参数说明：

**表 3.186 执行动作接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| hcuIndex | integer | Y | 设备序列 |

### 3.2.17 Profile

#### 3.2.17.1 获取性能指标

函数原型：GetSupportedMetricGroups(hcuIndex int) ([]MetricGroup, error)

功能说明：获取性能指标

参数说明：

**表 3.187 获取性能指标接口参数说明**

| 参数名 | 数据类型 | 是否必填 | 说明 |
| --- | --- | --- | --- |
| hcuIndex | integer | Y | 设备序列 |

**表 3.188 获取性能指标接口返回值说明**

| 返回属性名 | 类型 | 说明 |
| --- | --- | --- |
| Major | integer | 性能指标大类 |
| Minor | integer | 性能指标小类 |
| FieldIds | []int | 字段ID列表 |

