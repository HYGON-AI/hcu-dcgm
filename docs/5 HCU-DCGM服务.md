# 5 HCU-DCGM服务

## 5.1 接口使用说明

### 5.1.1 Device

#### 5.1.1.1 获取所有物理设备信息

**表 5.1 获取所有物理设备信息接口说明**

| 接口名称 | 获取所有物理设备信息 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 该接口返回所有物理设备的详细信息。 |  |  |  |
| URL | /AllDeviceInfos |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 所有设备的详细信息 |  |  |  |
| 400 | 无效的请求参数 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| blocksInfos | array:router.BlocksInfo |  | 设备Block信息列表 |  |
| block | string |  | Block Block 名称 |  |
| ce | integer |  | CE Correctable Error 数量 |  |
| state | string |  | State Block 状态 |  |
| ue | integer |  | UE Uncorrectable Error 数量 |  |
| clk | number |  | 当前时钟频率 |  |
| computeUnitCount | number |  | 计算单元总数 |  |
| computeUnitRemainingCount | integer |  | 剩余可用计算单元数量 |  |
| deviceId | string |  | 设备的唯一标识符 |  |
| maxVDeviceCount | integer |  | 最大虚拟设备数量 |  |
| memoryCap | number |  | 显存容量 |  |
| memoryRemaining | integer |  | 剩余可用显存 |  |
| memoryUsed | number |  | 已使用显存 |  |
| minorNumber | integer |  | 设备的索引号 |  |
| pciBusNumber | string |  | 设备的总线编号 |  |
| pcieBwMb | number |  | PCIe带宽 |  |
| powerCap | number |  | 功耗上限 |  |
| powerUsage | number |  | 当前功耗 |  |
| subSystemName | string |  | 设备的子系统名称 |  |
| temperature | number |  | 当前温度（摄氏度） |  |
| utilizationRate | number |  | 设备利用率 |  |
| vDeviceCount | integer |  | 当前虚拟设备数量 |  |
| virtualDevices | array:router.DMIVDeviceInfo |  | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |  |
| computeUnitCount | integer |  | 虚拟设备计算单元数量 |  |
| containerID | integer |  | 所属容器ID |  |
| deviceID | integer |  | 设备ID |  |
| globalMemSize | integer |  | 全局内存大小 |  |
| name | string |  | 虚拟设备名称 |  |
| pciBusNumber | string |  | 虚拟设备总线编号 |  |
| percent | integer |  | 使用百分比 |  |
| usageMemSize | integer |  | 已使用内存大小 |  |
| vminorNumber | integer |  | 虚拟设备索引号 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {    "powerCap": 0.0,    "clk": 0.0,    "computeUnitCount": 0.0,    "memoryRemaining": 0,    "subSystemName": "string",    "utilizationRate": 0.0,    "powerUsage": 0.0,    "blocksInfos": [        {            "block": "string",            "ce": 0,            "state": "string",            "ue": 0        }    ],    "computeUnitRemainingCount": 0,    "vDeviceCount": 0,    "memoryUsed": 0.0,    "deviceId": "string",    "pcieBwMb": 0.0,    "pciBusNumber": "string",    "maxVDeviceCount": 0,    "minorNumber": 0,    "virtualDevices": [        {            "computeUnitCount": 0,            "containerID": 0,            "deviceID": 0,            "globalMemSize": 0,            "name": "string",            "pciBusNumber": "string",            "percent": 0,            "usageMemSize": 0,            "vminorNumber": 0        }    ],    "memoryCap": 0.0,    "temperature": 0.0} |  |  |  |

#### 5.1.1.2 获取设备监控中的指标

**表 5.2 获取设备监控中的指标接口说明**

| 接口名称 | 获取设备监控中的指标 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 收集所有设备的监控指标信息。 |  |  |  |
| URL | /CollectDeviceMetrics |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 |  |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 设备监控指标信息列表 |  |  |  |
| 400 | 请求错误 |  |  |  |
| 404 | 设备未找到 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| clk | number(float64) |  | Clk 备系统时钟速度列表 |  |
| deviceId | string |  | DeviceId 设备序列号 |  |
| memoryCap | number(float64) |  | MemoryCap 设备内存总量 |  |
| memoryUsed | number(float64) |  | MemoryUsed 设备内存使用量 |  |
| minorNumber | integer |  | MinorNumber 设备索引号 |  |
| pciBusNumber | string |  | PciBusNumber PCI ID |  |
| pcieBwMb | number(float64) |  | PcieBwMb pcie流量信息 |  |
| powerCap | number(float64) |  | PowerCap 设备功率上限 |  |
| powerUsage | number(float64) |  | PowerUsage 设备平均功耗 |  |
| subSystemName | string |  | SubSystemName 型号名称 |  |
| temperature | number(float64) |  | Temperature 设备温度 |  |
| utilizationRate | number(float64) |  | UtilizationRate 设备忙碌时间百分比 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {    "pcieBwMb": null,    "powerCap": null,    "pciBusNumber": "string",    "clk": null,    "minorNumber": 0,    "subSystemName": "string",    "utilizationRate": null,    "powerUsage": null,    "memoryCap": null,    "temperature": null,    "memoryUsed": null,    "deviceId": "string"} |  |  |  |

#### 5.1.1.3 删除健康检查项的配置

**表 5.3 删除健康检查项的配置接口说明**

| 接口名称 | 删除健康检查项的配置 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 删除当前的健康检查项配置 |  |  |  |
| URL | /DeleteHealthCheckConfig |  |  |  |
| 请求方式 | delete |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 删除成功 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 |  |  |  |  |

#### 5.1.1.4 获取设备品牌名称

**表 5.4 获取设备品牌名称接口说明**

| 接口名称 | 获取设备品牌名称 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 根据设备索引获取品牌名称 |  |  |  |
| URL | /DevBrand/{dvInd} |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 设备索引 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回设备品牌名称 |  |  |  |
| 400 | 请求失败 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| brand | string |  | Brand 设备品牌名称 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"brand":"string"} |  |  |  |

#### 5.1.1.5 获取 GPU 度量信息

**表 5.5 获取 GPU 度量信息接口说明**

| 接口名称 | 获取 GPU 度量信息 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 根据设备 ID 获取 GPU 的度量信息。 |  |  |  |
| URL | /DevGpuMetricsInfo/{dvInd} |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 |  |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 设备 ID |
| 状态码 | 描述 |  | 说明 |  |
| 200 | GPU 度量信息 |  |  |  |
| 400 | 请求错误 |  |  |  |
| 404 | 设备未找到 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| averageDclk0Frequency | integer(int32) |  | AverageDclk0Frequency 平均显示时钟0频率 |  |
| averageDclk1Frequency | integer(int32) |  | AverageDclk1Frequency 平均显示时钟1频率 |  |
| averageGfxActivity | integer(int32) |  | AverageGfxActivity 平均图形活动 |  |
| averageGfxclkFrequency | integer(int32) |  | AverageGfxclkFrequency 平均图形时钟频率 |  |
| averageMmActivity | integer(int32) |  | AverageMmActivity 平均多媒体活动 |  |
| averageSocclkFrequency | integer(int32) |  | AverageSocclkFrequency 平均SoC时钟频率 |  |
| averageSocketPower | integer(int32) |  | AverageSocketPower 平均插座功率 |  |
| averageUclkFrequency | integer(int32) |  | AverageUclkFrequency 平均内存时钟频率 |  |
| averageUmcActivity | integer(int32) |  | AverageUmcActivity 平均内存控制器活动 |  |
| averageVclk0Frequency | integer(int32) |  | AverageVclk0Frequency 平均视频时钟0频率 |  |
| averageVclk1Frequency | integer(int32) |  | AverageVclk1Frequency 平均视频时钟1频率 |  |
| commonHeader | object |  | CommonHeader 公共表头 |  |
| currentDclk0 | integer(int32) |  | CurrentDclk0 当前显示时钟0 |  |
| currentDclk1 | integer(int32) |  | CurrentDclk1 当前显示时钟1 |  |
| currentFanSpeed | integer(int32) |  | CurrentFanSpeed 当前风扇速度 |  |
| currentGfxclk | integer(int32) |  | CurrentGfxclk 当前图形时钟 |  |
| currentSocclk | integer(int32) |  | CurrentSocclk 当前SoC时钟 |  |
| currentUclk | integer(int32) |  | CurrentUclk 当前内存时钟 |  |
| currentVclk0 | integer(int32) |  | CurrentVclk0 当前视频时钟0 |  |
| currentVclk1 | integer(int32) |  | CurrentVclk1 当前视频时钟1 |  |
| energyAccumulator | integer(int64) |  | EnergyAccumulator 能量累加器 |  |
| gfxActivityAcc | integer(int32) |  | GfxActivityAcc 图形活动累加器 |  |
| memActivityAcc | integer(int32) |  | MemActivityAcc 内存活动累加器 |  |
| padding | integer(int32) |  | Padding 填充 |  |
| pcieLinkSpeed | integer(int32) |  | PcieLinkSpeed PCIe链路速度（0.1 GT/s） |  |
| pcieLinkWidth | integer(int32) |  | PcieLinkWidth PCIe链路宽度 |  |
| systemClockCounter | integer(int64) |  | SystemClockCounter 系统时钟计数器 |  |
| temperatureEdge | integer(int32) |  | TemperatureEdge 边缘温度 |  |
| temperatureHotspot | integer(int32) |  | TemperatureHotspot 热点温度 |  |
| temperatureMem | integer(int32) |  | TemperatureMem 内存温度 |  |
| temperatureVRGfx | integer(int32) |  | TemperatureVRGfx VR图形温度 |  |
| temperatureVRMem | integer(int32) |  | TemperatureVRMem VR内存温度 |  |
| temperatureVRSoc | integer(int32) |  | TemperatureVRSoc VRSoC温度 |  |
| tempetureHBM | array |  | TempetureHBM 高带宽内存温度 |  |
| throttleStatus | integer(int32) |  | ThrottleStatus 节流状态 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {    "averageDclk0Frequency": 0,    "energyAccumulator": 0,    "currentVclk1": 0,    "gfxActivityAcc": 0,    "currentFanSpeed": 0,    "currentVclk0": 0,    "averageSocketPower": 0,    "temperatureEdge": 0,    "temperatureMem": 0,    "systemClockCounter": 0,    "currentGfxclk": 0,    "tempetureHBM": [        {}    ],    "averageGfxActivity": 0,    "temperatureHotspot": 0,    "memActivityAcc": 0,    "pcieLinkSpeed": 0,    "temperatureVRMem": 0,    "averageUmcActivity": 0,    "padding": 0,    "averageMmActivity": 0,    "pcieLinkWidth": 0,    "averageDclk1Frequency": 0,    "averageGfxclkFrequency": 0,    "temperatureVRSoc": 0,    "averageVclk0Frequency": 0,    "currentDclk0": 0,    "currentUclk": 0,    "averageUclkFrequency": 0,    "averageVclk1Frequency": 0,    "currentDclk1": 0,    "temperatureVRGfx": 0,    "throttleStatus": 0,    "averageSocclkFrequency": 0,    "commonHeader": {},    "currentSocclk": 0} |  |  |  |

#### 5.1.1.6 获取设备ID的十六进制值

**表 5.6 获取设备ID的十六进制值接口说明**

| 接口名称 | 获取设备ID的十六进制值 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 根据设备索引返回设备ID的十六进制值 |  |  |  |
| URL | /DevID |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | query | Y | 设备索引 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 设备ID十六进制值 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| id | string |  | ID 设备ID的十六进制值 |  |
| 示例 |  |  |  |  |
| 请求参数 | dvInd=0 |  |  |  |
| 返回值 | {"id":"string"} |  |  |  |

#### 5.1.17 获取设备SKU

**表 5.7 获取设备SKU接口说明**

| 接口名称 | 获取设备SKU |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 根据设备索引获取SKU |  |  |  |
| URL | /DevSku/{dvInd} |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 设备索引 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回设备SKU |  |  |  |
| 400 | 获取设备SKU失败 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| sku | integer |  | Sku 设备 SKU |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"sku":0} |  |  |  |

#### 5.1.1.8 获取设备供应商名称

**表 5.8 获取设备供应商名称接口说明**

| 接口名称 | 获取设备供应商名称 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 根据设备索引获取供应商名称 |  |  |  |
| URL | /DevVendorName/{dvInd} |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 设备索引 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回设备供应商名称 |  |  |  |
| 400 | 请求失败 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| bname | string |  | BName 设备供应商名称 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"bname":"string"} |  |  |  |

#### 5.1.1.9 获取设备显存供应商名称

**表 5.9 获取设备显存供应商名称接口说明**

| 接口名称 | 获取设备显存供应商名称 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 根据设备索引获取显存供应商名称 |  |  |  |
| URL | /DevVramVendor/{dvInd} |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 设备索引 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回显存供应商名称 |  |  |  |
| 400 | 请求失败 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| vendorName | string |  | VendorName 显存供应商名称 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"vendorName":"string"} |  |  |  |

#### 5.1.1.10 获取设备信息列表

**表 5.10 获取设备信息列表接口说明**

| 接口名称 | 获取设备信息列表 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 返回所有设备的详细信息列表 |  |  |  |
| URL | /DeviceInfos |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回设备信息列表 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| computeUnit | number(float64) |  | ComputeUnit 设备的计算单元数量 |  |
| devType | string |  | DevType 设备类型 |  |
| devTypeName | string |  | DevTypeName 设备类型名称 |  |
| deviceId | string |  | DeviceId 设备ID |  |
| dvInd | integer |  | DvInd 设备索引 |  |
| memoryTotal | number(float64) |  | MemoryTotal 设备的内存总量 |  |
| memoryUsed | number(float64) |  | MemoryUsed 设备的已使用内存量 |  |
| pciBusNumber | string |  | PciBusNumber 设备的总线号 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {    "devType": "string",    "pciBusNumber": "string",    "memoryTotal": null,    "devTypeName": "string",    "computeUnit": null,    "dvInd": 0,    "memoryUsed": null,    "deviceId": "string"} |  |  |  |

#### 5.1.1.11 获取所有设备的型号信息

**表 5.11 获取所有设备的型号信息接口说明**

| 接口名称 | 获取所有设备的型号信息 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 返回系统中所有设备的型号信息列表 |  |  |  |
| URL | /GetDeviceModelInfos |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 设备型号信息列表 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| devices | array:router.DeviceModelInfo |  | Devices 设备型号信息列表 |  |
| cuCount | number |  | CUCount CU 数量 |  |
| memorySize | number |  | MemorySize 内存大小（单位：MB 或 GB，视实现而定） |  |
| model | string |  | Model HCU 类型 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {    "devices": [        {            "cuCount": 0.0,            "memorySize": 0.0,            "model": "string"        }    ]} |  |  |  |

#### 5.1.1.12 获取健康检查项的配置

**表 5.12 获取健康检查项的配置接口说明**

| 接口名称 | 获取健康检查项的配置 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 获取当前健康检查项的配置信息 |  |  |  |
| URL | /GetHealthCheckConfig |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回健康检查项信息 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| enabled | boolean |  | Enabled 是否开启健康检查 |  |
| options | array |  | Options 健康检查项列表 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"options":[{}],"enabled":true} |  |  |  |

#### 5.1.1.13 获取指定设备的健康检查结果

**表 5.13 获取指定设备的健康检查结果接口说明**

| 接口名称 | 获取指定设备的健康检查结果 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 根据设备 ID 列表返回对应设备的健康检查状态信息 |  |  |  |
| URL | /HealthCheckById |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvIdList | array | body | Y | 设备 ID 列表 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回指定设备的健康检查结果 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 内部服务器错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| deviceHealths | array:router.DeviceHealth |  | DeviceHealths 设备健康结果列表 |  |
| hcu | integer |  | HCU HCU编号 |  |
| status | string |  | Status 状态 |  |
| watches | array:router.SystemWatch |  | Watches 检查详情 |  |
| error | string |  | Error 错误的详细信息 |  |
| result | object |  | Result 结果 |  |
| status | string |  | Status 状态 |  |
| type | string |  | Type 检查项的类型 |  |
| 示例 |  |  |  |  |
| 请求参数 | -d '[{}]' |  |  |  |
| 返回值 | {    "deviceHealths": [        {            "hcu": 0,            "status": "string",            "watches": [                {                    "error": "string",                    "result": {},                    "status": "string",                    "type": "string"                }            ]        }    ]} |  |  |  |

#### 5.1.1.14 获取 GPU 数量

**表 5.14 获取 GPU 数量接口说明**

| 接口名称 | 获取 GPU 数量 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 获取监视的 GPU 数量 |  |  |  |
| URL | /NumMonitorDevices |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | GPU 数量 |  |  |  |
| 500 | 获取 GPU 数量失败 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| gpuCount | integer |  | GpuCount GPU 数量 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"gpuCount":0} |  |  |  |

#### 5.1.1.15 运行设备诊断

**表 5.15 运行设备诊断接口说明**

| 接口名称 | 运行设备诊断 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 运行设备的诊断工具，返回诊断结果 |  |  |  |
| URL | /diag/run/{level} |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| level | integer | query | Y | 诊断级别 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回诊断结果 |  |  |  |
| 400 | 参数无效 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| deviceNumber | integer |  | DeviceNumber 设备数量信息 |  |
| perHCU | array:router.HCUResult |  | PerHCU 每个 HCU 的诊断结果 |  |
| hcu | integer |  | HCU HCU 编号 |  |
| diagResults | array:router.DiagResult |  | DiagResults HCU 下的诊断结果列表 |  |
| errorCode | integer |  | ErrorCode 错误码 |  |
| errorMessage | string |  | ErrorMessage 错误信息 |  |
| status | string |  | Status 诊断结果状态 |  |
| testName | string |  | TestName 测试名称 |  |
| testOutput | string |  | TestOutput 测试输出信息 |  |
| rc | integer |  | RC 返回码 |  |
| software | array:router.DiagResult |  | Software 软件层诊断结果 |  |
| errorCode | integer |  | ErrorCode 错误码 |  |
| errorMessage | string |  | ErrorMessage 错误信息 |  |
| status | string |  | Status 诊断结果状态 |  |
| testName | string |  | TestName 测试名称 |  |
| testOutput | string |  | TestOutput 测试输出信息 |  |
| 示例 |  |  |  |  |
| 请求参数 | level=0 |  |  |  |
| 返回值 | {    "perHCU": [        {            "hcu": 0,            "diagResults": [                {                    "errorCode": 0,                    "errorMessage": "string",                    "status": "string",                    "testName": "string",                    "testOutput": "string"                }            ],            "rc": 0        }    ],    "software": [        {            "errorCode": 0,            "errorMessage": "string",            "status": "string",            "testName": "string",            "testOutput": "string"        }    ],    "deviceNumber": 0} |  |  |  |
| 16）设置是否开启健康检查 |  |  |  |  |

#### 5.1.1.16 设置是否开启健康检查

**表 5.16 设置是否开启健康检查接口说明**

| 接口名称 | 设置是否开启健康检查 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 配置是否启用健康检查，并选择检查项 |  |  |  |
| URL | /SetHealthCheckConfig |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| enabled | object | body | Y | 是否开启 |
| options | array | body | Y | 检查项 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 设置成功 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| 示例 |  |  |  |  |
| 请求参数 | -d '{"enabled":{},"options":[{}]}' |  |  |  |
| 返回值 |  |  |  |  |

#### 5.1.1.17 显示设备硬件信息

**表 5.17 显示设备硬件信息接口说明**

| 接口名称 | 显示设备硬件信息 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 显示指定设备列表的简要硬件信息 |  |  |  |
| URL | /ShowAllConciseHw |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvIdList | array | body | Y | 设备 ID 列表 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 成功信息 |  |  |  |
| 400 | 失败信息 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 |  |  |  |  |

#### 5.1.1.18 控制设备的性能级别、时钟频率和风扇重置

**表 5.18 控制设备的性能级别、时钟频率和风扇重置接口说明**

| 接口名称 | 控制设备的性能级别、时钟频率和风扇重置 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 根据传入的设备控制信息，设置设备的性能级别、时钟频率，并可选择性重置风扇 |  |  |  |
| URL | /device/control |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| deviceControl | object:router.DeviceControlInfo | body | Y | 设备控制信息 |
| dvInd | integer |  | N | DvInd 设备索引号 |
| perfLevel | string |  | N | PerfLevel 性能水平 |
| resetFan | boolean |  | N | ResetFan 是否重置风扇控制 |
| sclkClock | string |  | N | SclkClock sclk时钟频率 600、700、750、800、900、1000、1106、1200、1270、1319、1400、1500、1600 |
| socclkClock | string |  | N | SocclkClock soclk时钟频率 309、523、566、618、680、755、850、971 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 操作成功 |  |  |  |
| 400 | 无效的请求参数或操作失败 |  |  |  |
| 500 | 内部服务器错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| errors | array |  | Errors 执行失败的错误信息 |  |
| success | boolean |  | Success 是否执行成功 |  |
| 示例 |  |  |  |  |
| 请求参数 | -d '{"dvInd":0,"perfLevel":"string","resetFan":true,"sclkClock":"string","socclkClock":"string"}' |  |  |  |
| 返回值 | {"errors":[{}],"success":true} |  |  |  |

#### 5.1.1.19 获取设备信息

**表 5.19 获取设备信息（/device/info）接口说明**

| 接口名称 | 获取设备信息 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 根据设备索引获取对应的设备信息 |  |  |  |
| URL | /device/info/{dvInd} |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 |  |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 设备索引 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 设备信息 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 内部服务器错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| deviceInfo | object |  | DeviceInfo 设备信息 |  |
| name | string |  | 设备名称 |  |
| computeUnitCount | int |  | 计算单元总数 |  |
| ComputeUnitRemainingCount | integer(uint64) |  | 剩余计算单元数 |  |
| MemoryRemaining | integer(uint64) |  | 剩余内存 |  |
| GlobalMemSize | integer(uint64) |  | 全局内存大小 |  |
| UsageMemSize | integer(uint64) |  | 已使用内存大小 |  |
| deviceID | int |  | 设备索引号 |  |
| percent | int |  | 使用率（百分比） |  |
| maxVDeviceCount | int |  | 最大虚拟设备数量 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {    "deviceInfo": {        "name": "K100_AI",        "computeUnitCount": 120,        "ComputeUnitRemainingCount": 0,        "MemoryRemaining": 0,        "GlobalMemSize": 0,        "UsageMemSize": 0,        "deviceID": 1,        "percent": 0,        "maxVDeviceCount": 0    }} |  |  |  |

#### 5.1.1.20 设备的产品信息列表

**表 5.20 显示设备的产品名称接口说明**

| 接口名称 | 显示设备的产品名称 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 获取并显示指定设备的产品名称、供应商、系列、型号和 SKU 信息 |  |  |  |
| URL | /device/product |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 |  |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvIdList | array | body | Y | 设备 ID 列表 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 设备的产品信息列表 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| cardModel | string |  | CardModel 设备型号 example: "HHCU-1000" |  |
| cardSeries | string |  | CardSeries 设备系列名称 example: "Hygon HCU" |  |
| cardSku | string |  | CardSKU SKU 信息 example: "SKU-001" |  |
| cardVendor | string |  | CardVendor 设备供应商 example: "Hygon" |  |
| deviceId | integer |  | DeviceID 设备索引号 example: 0 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {    "cardSku": "string",    "cardModel": "string",    "cardVendor": "string",    "deviceId": 0,    "cardSeries": "string"} |  |  |  |

#### 5.1.1.21 显示设备的序列号

**表 5.21 显示设备的序列号接口说明**

| 接口名称 | 显示设备的序列号 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 获取并显示指定设备的序列号信息 |  |  |  |
| URL | /device/serialnumber |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 |  |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvIdList | array | body | Y | 设备ID列表 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 设备的序列号信息列表 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| deviceId | integer |  | DeviceID 设备索引号 example: 0 |  |
| serialNumber | string |  | SerialNumber 设备序列号 example: "SN123456789" |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {    "serialNumber": "string",    "deviceId": 0} |  |  |  |

#### 5.1.1.22 获取设备子系统 ID

**表 5.22 获取设备子系统 ID接口说明**

| 接口名称 | 获取设备子系统 ID |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 根据物理设备索引 dvInd，返回设备子系统 ID |  |  |  |
| URL | /device/subsystem/id |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | query | Y | 物理设备索引（Device Index） |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 设备子系统 ID |  |  |  |
| 400 | Invalid device ID |  |  |  |
| 500 | Internal Server Error |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| subsystemId | string |  | SubsystemId 子系统 ID |  |
| 示例 |  |  |  |  |
| 请求参数 | dvInd=0 |  |  |  |
| 返回值 | {"subsystemId":"string"} |  |  |  |

#### 5.1.1.23 获取设备子系统名称

**表 5.23 获取设备子系统名称接口说明**

| 接口名称 | 获取设备子系统名称 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 根据物理设备索引 dvInd，返回设备子系统名称 |  |  |  |
| URL | /device/subsystem/name |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | query | Y | 物理设备索引（Device Index） |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 设备子系统名称 |  |  |  |
| 400 | Invalid device ID |  |  |  |
| 500 | Internal Server Error |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| subsystemName | string |  | SubsystemName 子系统名称 |  |
| 示例 |  |  |  |  |
| 请求参数 | dvInd=0 |  |  |  |
| 返回值 | {"subsystemName":"string"} |  |  |  |

#### 5.1.1.24 获取设备类型名称

**表 5.24 获取设备类型名称接口说明**

| 接口名称 | 获取设备类型名称 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 根据物理设备索引 dvInd，返回设备类型名称及单位 |  |  |  |
| URL | /device/type/name |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | query | Y | 物理设备索引（Device Index） |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 设备类型名称信息 |  |  |  |
| 400 | Invalid device ID |  |  |  |
| 500 | Internal Server Error |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| devTypeName | string |  | DevTypeName 设备类型名称 |  |
| unit | number |  | Unit 单位 |  |
| 示例 |  |  |  |  |
| 请求参数 | dvInd=0 |  |  |  |
| 返回值 | {"devTypeName":"string","unit":0.0} |  |  |  |

#### 5.1.1.25 获取物理设备的虚拟设备索引

**表 5.25 获取物理设备的虚拟设备索引接口说明**

| 接口名称 | 获取物理设备的虚拟设备索引 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 根据物理设备索引 dvInd，返回该物理设备下的虚拟设备数量及虚拟设备索引列表 |  |  |  |
| URL | /device/{dvInd}/vdevices |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 物理设备索引（Device Index） |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 虚拟设备索引信息 |  |  |  |
| 400 | Invalid device ID |  |  |  |
| 500 | Internal Server Error |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| vDevInds | array |  | VDevInds 虚拟设备索引列表 |  |
| vDeviceCount | integer |  | VDeviceCount 虚拟设备数量 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {    "vDeviceCount": 0,    "vDevInds": [        {}    ]} |  |  |  |

#### 5.1.1.26 获取设备数量

**表 5.26 获取设备数量接口说明**

| 接口名称 | 获取设备数量 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 获取当前系统中的设备数量 |  |  |  |
| URL | /deviceCount |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 |  |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 设备数量信息 |  |  |  |
| 500 | 内部服务器错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| count | integer |  | Count 设备数量 example: 4 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"count":0} |  |  |  |

#### 5.1.1.27 获取设备信息

**表 5.27 获取设备信息（/deviceinfo）接口说明**

| 接口名称 | 获取设备信息 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 根据设备 ID 获取物理设备的详细信息 |  |  |  |
| URL | /deviceinfo/{dvInd} |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | Device ID |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 设备信息 |  |  |  |
| 400 | Invalid device ID |  |  |  |
| 500 | Internal Server Error |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| blocksInfos | array:router.BlocksInfo |  | 设备Block信息列表 |  |
| block | string |  | Block Block 名称 |  |
| ce | integer |  | CE Correctable Error 数量 |  |
| state | string |  | State Block 状态 |  |
| ue | integer |  | UE Uncorrectable Error 数量 |  |
| clk | number |  | 当前时钟频率（MHz） |  |
| computeUnitCount | number |  | 计算单元总数 |  |
| computeUnitRemainingCount | integer |  | 剩余可用计算单元数量 |  |
| deviceId | string |  | 设备的唯一标识符 |  |
| maxVDeviceCount | integer |  | 最大虚拟设备数量 |  |
| memoryCap | number |  | 显存容量 |  |
| memoryRemaining | integer |  | 剩余可用显存 |  |
| memoryUsed | number |  | 已使用显存 |  |
| minorNumber | integer |  | 设备的索引号 |  |
| pciBusNumber | string |  | 设备的总线编号 |  |
| pcieBwMb | number |  | PCIe带宽 |  |
| powerCap | number |  | 功耗上限 |  |
| powerUsage | number |  | 当前功耗 |  |
| subSystemName | string |  | 设备的子系统名称 |  |
| temperature | number |  | 当前温度 |  |
| utilizationRate | number |  | 设备利用率（百分比） |  |
| vDeviceCount | integer |  | 当前虚拟设备数量 |  |
| virtualDevices | array:router.DMIVDeviceInfo |  | VirtualDevices 该物理设备上关联的虚拟设备信息列表 |  |
| computeUnitCount | integer |  | 虚拟设备计算单元数量 |  |
| containerID | integer |  | 所属容器ID |  |
| deviceID | integer |  | 设备ID |  |
| globalMemSize | integer |  | 全局内存大小 |  |
| name | string |  | 虚拟设备名称 |  |
| pciBusNumber | string |  | 虚拟设备总线编号 |  |
| percent | integer |  | 使用百分比 |  |
| usageMemSize | integer |  | 已使用内存大小 |  |
| vminorNumber | integer |  | 虚拟设备索引号 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {    "powerCap": 0.0,    "clk": 0.0,    "computeUnitCount": 0.0,    "memoryRemaining": 0,    "subSystemName": "string",    "utilizationRate": 0.0,    "powerUsage": 0.0,    "blocksInfos": [        {            "block": "string",            "ce": 0,            "state": "string",            "ue": 0        }    ],    "computeUnitRemainingCount": 0,    "vDeviceCount": 0,    "memoryUsed": 0.0,    "deviceId": "string",    "pcieBwMb": 0.0,    "pciBusNumber": "string",    "maxVDeviceCount": 0,    "minorNumber": 0,    "virtualDevices": [        {            "computeUnitCount": 0,            "containerID": 0,            "deviceID": 0,            "globalMemSize": 0,            "name": "string",            "pciBusNumber": "string",            "percent": 0,            "usageMemSize": 0,            "vminorNumber": 0        }    ],    "memoryCap": 0.0,    "temperature": 0.0} |  |  |  |

#### 5.1.1.28 获取设备名称

**表 5.28 获取设备名称接口说明**

| 接口名称 | 获取设备名称 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 根据设备索引 ID 获取设备名称 |  |  |  |
| URL | /devices/{dvInd}/name |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 |  |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 设备索引ID |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 成功返回设备名称 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| deviceName | string |  |  |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"deviceName":"string"} |  |  |  |

#### 5.1.1.29 获取 DF 带宽信息

**表 5.29 获取 DF 带宽信息接口说明**

| 接口名称 | 获取 DF 带宽信息 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 根据物理设备索引 dvInd 获取 DF（Data Fabric）带宽信息 |  |  |  |
| URL | /dfbandwidth/{dvInd} |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 物理设备索引（Device Index） |
| 状态码 | 描述 |  | 说明 |  |
| 200 | DF 带宽信息 |  |  |  |
| 400 | Invalid device ID |  |  |  |
| 500 | Internal Server Error |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| dfBandwidthInfo | object |  | DFBandwidthInfo DF 带宽信息 |  |
| readBW | number（float64） |  | DF 读取带宽 |  |
| writeBW | number（float64） |  | DF 写入带宽 |  |
| readWriteBW | number（float64） |  | DF 读写综合带宽 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {    "dfBandwidthInfo": {        "readBW": 0,        "writeBW": 0,        "readWriteBW": 0    }} |  |  |  |

#### 5.1.1.30 显示设备的 HCU 使用率

**表 5.30 显示设备的 HCU 使用率接口说明**

| 接口名称 | 显示设备的 HCU 使用率 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 无 |  |  |  |
| URL | /gpu/use |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvIdList | array | body | Y | 设备 ID 列表 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 设备使用信息列表 |  |  |  |
| 400 | 错误信息 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| deviceUseInfos | array:router.DeviceUseInfo |  | DeviceUseInfos 设备使用信息列表 |  |
| deviceId | integer |  | DeviceID 设备索引号 |  |
| gpuUsePercent | integer |  | GPUUsePercent GPU 使用率（百分比） |  |
| utilization | object |  | Utilization 各类利用率统计 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {    "deviceUseInfos": [        {            "deviceId": 0,            "gpuUsePercent": 0,            "utilization": {}        }    ]} |  |  |  |

#### 5.1.1.31 获取 MIG 分区信息

**表 5.31 获取 MIG 分区信息接口说明**

| 接口名称 | 获取 MIG 分区信息 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 获取系统中所有 MIG 分区的详细信息 |  |  |  |
| URL | /mig/infos |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | MIG 分区信息列表 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| migInfos | array:router.MigInfo |  | MigInfos MIG 分区信息列表 |  |
| ciProfileId | integer |  | CiProfileId CI Profile ID |  |
| computeInstanceId | integer |  | ComputeInstanceId Compute Instance ID |  |
| computeUnit | integer |  | ComputeUnit 计算单元数量 |  |
| dvInd | integer |  | DvInd 物理设备索引 |  |
| giProfileId | integer |  | GiProfileId GI Profile ID |  |
| gpuInstanceId | integer |  | GpuInstanceId GPU Instance ID |  |
| memoryTotal | integer |  | MemoryTotal 分区总显存（Bytes） |  |
| migId | integer |  | MigId MIG 分区 ID |  |
| name | string |  | Name MIG 分区名称 |  |
| pciBusNumber | string |  | PciBusNumber PCI Bus 号 |  |
| uuid | string |  | UUID MIG 分区 UUID |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {    "migInfos": [        {            "ciProfileId": 0,            "computeInstanceId": 0,            "computeUnit": 0,            "dvInd": 0,            "giProfileId": 0,            "gpuInstanceId": 0,            "memoryTotal": 0,            "migId": 0,            "name": "string",            "pciBusNumber": "string",            "uuid": "string"        }    ]} |  |  |  |

#### 5.1.1.32 显示设备的唯一ID

**表 5.32 显示设备的唯一ID接口说明**

| 接口名称 | 显示设备的唯一ID |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 获取并显示指定设备的唯一ID信息。 |  |  |  |
| URL | /showUId |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 |  |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvIdList | array | body | Y | 设备ID列表 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 设备的唯一ID信息列表 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| deviceId | integer |  | DeviceID 设备索引号 example: 0 |  |
| uid | string |  | UId 设备唯一ID example: "UID123456789" |  |
| 示例 |  |  |  |  |
| 请求参数 | -d '[{}]' |  |  |  |
| 返回值 | {    "uid": "string",    "deviceId": 0} |  |  |  |

#### 5.1.1.33 获取 UMC 带宽信息

**表 5.33 获取 UMC 带宽信息接口说明**

| 接口名称 | 获取 UMC 带宽信息 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 根据物理设备索引和 UMC 通道 ID 获取 UMC 带宽信息 |  |  |  |
| URL | /umc/bandwidth |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| body | object:router.UMCBandwidthReq | body | Y | UMC 带宽查询参数 |
| chanId | integer |  | N | ChanId UMC 通道 ID（0 ~ 31） |
| delay | integer |  | N | Delay 采样延迟（单位：秒） |
| dvInd | integer |  | N | DvInd 物理设备索引 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | UMC 带宽信息 |  |  |  |
| 400 | Invalid parameters |  |  |  |
| 500 | Internal Server Error |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| umcBandwidth | object |  | UMCBandwidth UMC 带宽信息 |  |
| readBW | array |  | 各通道读带宽 |  |
| writeBW | array |  | 各通道写带宽 |  |
| readWriteBW | array |  | ReadWriteBW |  |
| 示例 |  |  |  |  |
| 请求参数 | -d '{"chanId":0,"delay":0,"dvInd":0}' |  |  |  |
| 返回值 | {  "umcBandwidth": {    "readBW": [123.45, 130.12, 128.90],    "writeBW": [98.30, 100.56, 97.80],    "readWriteBW": [221.75, 230.68, 226.70]  }} |  |  |  |

#### 5.1.1.34 获取设备状态信息

**表 5.34 获取设备状态信息接口说明**

| 接口名称 | 获取设备状态信息 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 返回所有设备设备状态信息 |  |  |  |
| URL | /DeviceStatus |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回设备设备状态信息 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| blocksInfos | array:router.BlocksInfo |  | BlocksInfo 设备的block信息 |  |
| block | string |  | Block Block 名称 |  |
| ce | integer |  | CE Correctable Error 数量 |  |
| state | string |  | State Block 状态 |  |
| ue | integer |  | UE Uncorrectable Error 数量 |  |
| clk | number(float64) |  | Clk 设备的当前时钟频率 |  |
| computeUnitRemainingCount | integer(int64) |  | ComputeUnitRemainingCount 设备剩余可用的计算单元数量 |  |
| dvInd | integer |  | DvInd 设备索引 |  |
| memoryRemaining | integer(int64) |  | MemoryRemaining 设备剩余可用的内存量 |  |
| memoryUsed | number(float64) |  | MemoryUsed 设备已使用的内存 |  |
| pcieBwMb | number(float64) |  | PcieBwMb 设备的PCIe带宽 (MB/s) |  |
| pcieReceived | number(float64) |  | PcieReceived 设备的PCIe接收数据量 |  |
| pcieSent | number(float64) |  | PcieSent 设备的PCIe发送数据量 |  |
| percent | integer |  | Percent 物理设备使用百分比 |  |
| powerUsed | number(float64) |  | PowerUsed 设备当前的功耗 |  |
| temperature | number(float64) |  | Temperature 设备当前的温度 |  |
| utilizationRate | number(float64) |  | UtilizationRate 设备时间忙碌百分比 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {    "powerUsed": null,    "clk": null,    "memoryRemaining": 0,    "utilizationRate": null,    "blocksInfos": [        {            "block": "string",            "ce": 0,            "state": "string",            "ue": 0        }    ],    "computeUnitRemainingCount": 0,    "pcieSent": null,    "memoryUsed": null,    "percent": 0,    "pcieBwMb": null,    "temperature": null,    "dvInd": 0,    "pcieReceived": null} |  |  |  |

#### 5.1.1.35 校验兼容性

**表 5.35 校验兼容性接口说明**

| 接口名称 | 校验兼容性 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 根据卡型号、驱动版本和 DTK 版本，检查是否符合兼容性要求 |  |  |  |
| URL | /compatible |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| cardModel | string | query | Y | 卡型号，例如 K100, K100_AI |
| driverVersion | string | query | Y | 驱动版本，例如 5.16.18, 6.2.26 |
| dtkVersion | string | query | Y | DTK 版本，例如 23.10, 24.04 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 兼容性校验通过 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| 示例 |  |  |  |  |
| 请求参数 | cardModel=string&driverVersion=string&dtkVersion=string |  |  |  |
| 返回值 |  |  |  |  |

#### 5.1.1.36 获取 HCU 瞬时占用率

**表 5.159 获取 HCU 瞬时占用率接口说明**

| 接口名称 | 获取 HCU 瞬时占用率 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 获取 HCU 瞬时占用率 |  |  |  |
| URL | /HCUCuUsage/:dvInd |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 物理设备索引 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 请求成功 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| data | number |  | HCU 瞬时占用率 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"data":0} |  |  |  |

#### 5.1.1.37 获取 HCU 采样占用情况

**表 5.160 获取 HCU 采样占用情况接口说明**

| 接口名称 | 获取 HCU 采样占用情况 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 获取 HCU 采样占用情况 |  |  |  |
| URL | /HCUSampledUsage/:dvInd |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 物理设备索引 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 请求成功 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| data | number |  | HCU 采样占用率 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"data":0} |  |  |  |

#### 5.1.1.38 获取 CU 采样占用情况

**表 5.161 获取 CU 采样占用情况接口说明**

| 接口名称 | 获取 CU 采样占用情况 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 获取 CU 采样占用情况 |  |  |  |
| URL | /HCUCUSampledUsage/:dvInd |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 物理设备索引 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 请求成功 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| data | number |  | CU 采样占用率 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"data":0} |  |  |  |

#### 5.1.1.39 获取 Wave 采样占用情况

**表 5.162 获取 Wave 采样占用情况接口说明**

| 接口名称 | 获取 Wave 采样占用情况 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 获取 Wave 采样占用情况 |  |  |  |
| URL | /HCUWaveSampledUsage/:dvInd |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 物理设备索引 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 请求成功 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| data | number |  | Wave 采样占用率 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"data":0} |  |  |  |

#### 5.1.1.40 获取 SE 瞬时占用率

**表 5.163 获取 SE 瞬时占用率接口说明**

| 接口名称 | 获取 SE 瞬时占用率 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 获取 SE 瞬时占用率 |  |  |  |
| URL | /HCUSEUsage/:dvInd |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 物理设备索引 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 请求成功 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| data | number |  | SE 瞬时占用率 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"data":0} |  |  |  |

#### 5.1.1.41 查询 CU 数量

**表 5.174 查询 CU 数量接口说明**

| 接口名称 | 查询 CU 数量 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 查询 CU 数量 |  |  |  |
| URL | /device/cunum/:dvInd |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 物理设备索引 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 请求成功 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| cuCount | integer |  | CU 数量 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"cuCount":0} |  |  |  |

#### 5.1.1.42 查询 KFD node ID

**表 5.175 查询 KFD node ID接口说明**

| 接口名称 | 查询 KFD node ID |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 查询 KFD node ID |  |  |  |
| URL | /device/nodeid/:dvInd |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 物理设备索引 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 请求成功 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| nodeId | integer |  | KFD node ID |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"nodeId":0} |  |  |  |

#### 5.1.1.43 保存设备配置

**表 5.195 保存设备配置接口说明**

| 接口名称 | 保存设备配置 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 保存设备配置 |  |  |  |
| URL | /config/save |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| body | object | body | Y | 请求体参数，具体字段以接口实现为准 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 请求成功 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| message | string |  | 操作结果 |  |
| 示例 |  |  |  |  |
| 请求参数 | {} |  |  |  |
| 返回值 | {"message":"成功"} |  |  |  |

#### 5.1.1.44 加载设备配置

**表 5.196 加载设备配置接口说明**

| 接口名称 | 加载设备配置 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 加载设备配置 |  |  |  |
| URL | /config/load |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| body | object | body | Y | 请求体参数，具体字段以接口实现为准 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 请求成功 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| data | object |  | 设备配置 |  |
| 示例 |  |  |  |  |
| 请求参数 | {} |  |  |  |
| 返回值 | {"data":{}} |  |  |  |

### 5.1.2 Bandwidth

#### 5.1.2.1 查询所有设备的 HyLink Link 带宽状态

**表 5.36 查询所有设备的 HyLink Link 带宽状态接口说明**

| 接口名称 | 查询所有设备的 HyLink Link 带宽状态 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 无需请求参数，自动遍历所有监控设备， 返回每张卡的 Link 带宽聚合结果（Recv / Send）以及每个 Link 的明细。 |  |  |  |
| URL | /bandwidth/hylink/status |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 查询成功 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| devices | array |  | 所有设备的 HyLink 带宽汇总信息 |  |
| dvInd | integer |  | 设备索引号 |  |
| recv | number (float64) |  | 该设备所有 Link 的接收方向带宽总和 |  |
| send | number (float64) |  | 该设备所有 Link 的发送方向带宽总和 |  |
| err | string |  | 查询该设备时的错误信息，成功时为空字符串 |  |
| links | array |  | 该设备下每条 Link 的带宽明细 |  |
| linkId | integer |  | Link 索引号 |  |
| recv | number (float64) |  | 该 Link 的接收方向带宽 |  |
| send | number (float64) |  | 该 Link 的发送方向带宽 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {  "devices": [    {      "dvInd": 0,      "recv": 123.45,      "send": 98.76,      "err": "",      "links": [        {          "linkId": 0,          "recv": 60.12,          "send": 48.30        }      ]    }  ]} |  |  |  |

#### 5.1.2.2 查询 XHCL 带宽信息

**表 5.37 查询 XHCL 带宽信息接口说明**

| 接口名称 | 查询 XHCL 带宽信息 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 根据物理设备索引和链路参数获取 XHCL 带宽信息 |  |  |  |
| URL | /bandwidth/xhcl |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| req | object:router.XHCLBandwidthReq | body | Y | XHCL 带宽查询参数 |
| delay | integer |  | N | Delay 采样延迟（秒） |
| direction | integer |  | N | Direction 带宽方向 |
| dvInd | integer |  | N | DvInd 物理设备索引 |
| linkId | integer |  | N | LinkId XHCL 链路 ID |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 查询成功 |  |  |  |
| 400 | 参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| bw | array (float64) |  | XHCL 各链路的带宽值数组 | 数组最大长度为7 |
| bw[0] | number (float64) |  | LinkId = 0 的带宽 |  |
| bw[0] | bw[0]number (float64) |  | LinkId = n 的带宽，n < MAX_XHCL_LINK_NUM |  |
| 示例 |  |  |  |  |
| 请求参数 | '{"delay":0,"direction":0,"dvInd":0,"linkId":0}' |  |  |  |
| 返回值 | {  "bw": [    120.5,    118.3,    0,    0,    0,    0,    0,    0  ]} |  |  |  |

### 5.1.3 Diag

#### 5.1.3.1 启动 Memory Bandwidth 测试

**表 5.38 启动 Memory Bandwidth 压测接口说明**

| 接口名称 | 启动 Memory Bandwidth 压测 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 异步启动内存带宽测试，立即返回 job id。 测试结果会被封装为 DiagResults 并持久化到 logs/.json。 |  |  |  |
| URL | /bandwidth/start |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | OK |  |  |  |
| 409 | 已有作业正在运行 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| jobId | string |  | JobID 诊断任务 ID |  |
| status | string |  | Status 当前任务状态（pending / running） |  |
| statusUrl | string |  | StatusURL 查询任务状态的接口地址 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {  "jobId": "string",  "statusUrl": "string",  "status": "string"} |  |  |  |

#### 5.1.3.2 启动异步诊断任务

**表 5.39 启动异步诊断任务接口说明**

| 接口名称 | 启动异步诊断任务 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 启动指定级别（1~4）的异步诊断任务，立即返回 jobId；同一时间只允许一个诊断任务运行 |  |  |  |
| URL | /diag/start/{level} |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| level | integer | path | Y | 诊断级别（1~4） |
| 状态码 | 描述 |  | 说明 |  |
| 202 | Accepted |  |  |  |
| 400 | Bad Request |  |  |  |
| 409 | 已有诊断任务正在运行 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| jobId | string |  | JobID 诊断任务 ID |  |
| status | string |  | Status 当前任务状态（pending / running） |  |
| statusUrl | string |  | StatusURL 查询任务状态的接口地址 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | **{  "jobId": "string",  "statusUrl": "string",  "status": "string"} |  |  |  |

#### 5.1.3.3 查询诊断作业状态

**表 5.40 查询诊断作业状态接口说明**

| 接口名称 | 查询诊断作业状态 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 根据 job id 查询诊断作业当前状态及结果。 优先从内存中查询运行中或刚完成的作业； 若内存中不存在，则从磁盘 logs/.json 中读取历史结果。 |  |  |  |
| URL | /diag/status/{id} |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| id | string | path | Y | 作业 ID |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回作业状态及结果（若已完成） |  |  |  |
| 404 | 作业不存在 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| endedAt | string |  | EndedAt 作业结束（成功/失败/取消）的时间戳（零值表示尚未结束） example: 2025-12-29T10:30:00Z |  |
| errorMessage | string |  | ErrorMessage 如果作业执行失败或中止，这里保存错误信息或说明 example: GPU not found |  |
| id | string |  | ID 唯一作业标识符（由服务生成，例如 job--） example: job-1698493923-01 |  |
| level | integer |  | Level 要执行的诊断等级（1..4） example: 2 |  |
| result | object |  | Result 作业完成后保存的结构化诊断结果 |  |
| startedAt | string |  | StartedAt 作业实际开始执行的时间戳（零值表示尚未开始） example: 2025-12-29T10:00:00Z |  |
| status | string |  | Status 作业当前状态 example: pending |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {  "result": {},  "level": 0,  "endedAt": "string",  "errorMessage": "string",  "startedAt": "string",  "id": "string",  "status": "string"} |  |  |  |

#### 5.1.3.4 停止当前诊断作业

**表 5.41 停止当前诊断作业接口说明**

| 接口名称 | 停止当前诊断作业 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 请求停止当前正在运行的诊断任务（非强制杀死）。 后台任务会在安全检查点检测到停止标志后终止后续步骤。 |  |  |  |
| URL | /diag/stop |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 停止请求已接收 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 |  |  |  |  |

#### 5.1.3.5 启动 EDPp 测试

**表 5.42 启动 EDPp 测试接口说明**

| 接口名称 | 启动 EDPp 测试 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 异步启动 EDPp 稳定性与错误注入测试。 按 HCU 和 Pattern 聚合结果并持久化。 |  |  |  |
| URL | /edpp/start |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | OK |  |  |  |
| 409 | 已有作业正在运行 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| jobId | string |  | JobID 诊断任务 ID |  |
| status | string |  | Status 当前任务状态（pending / running） |  |
| statusUrl | string |  | StatusURL 查询任务状态的接口地址 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {  "jobId": "string",  "statusUrl": "string",  "status": "string"} |  |  |  |

#### 5.1.3.6 启动 PCIe 带宽测试

**表 5.43 启动 PCIe 带宽压测接口说明**

| 接口名称 | 启动 PCIe 带宽压测 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 异步启动 PCIe 带宽测试（Sys<->Fb）。 返回 job id，可通过 /diag/status/{id} 查询结果。 |  |  |  |
| URL | /pcie/start |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | OK |  |  |  |
| 409 | 已有作业正在运行 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| jobId | string |  | JobID 诊断任务 ID |  |
| status | string |  | Status 当前任务状态（pending / running） |  |
| statusUrl | string |  | StatusURL 查询任务状态的接口地址 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {  "jobId": "string",  "statusUrl": "string",  "status": "string"} |  |  |  |

#### 5.1.3.7 启动 TargetStress（GEMM）测试

**表 5.44 启动 TargetStress（GEMM）压测接口说明**

| 接口名称 | 启动 TargetStress（GEMM）压测 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 异步执行多种 GEMM 模式的 TargetStress 压测。 返回 job id，结果通过状态接口获取。 |  |  |  |
| URL | /targetstress/start |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | OK |  |  |  |
| 409 | 已有作业正在运行 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| jobId | string |  | JobID 诊断任务 ID |  |
| status | string |  | Status 当前任务状态（pending / running） |  |
| statusUrl | string |  | StatusURL 查询任务状态的接口地址 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {  "jobId": "string",  "statusUrl": "string",  "status": "string"} |  |  |  |

#### 5.1.3.8 启动 XHCL 带宽测试

**表 5.45 启动 XHCL 带宽测试接口说明**

| 接口名称 | 启动 XHCL 带宽测试 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 异步启动 XHCL（跨 HCU）互联带宽测试。 测试结果按 HCU 聚合并持久化。 |  |  |  |
| URL | /xhcl/start |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | OK |  |  |  |
| 409 | 已有作业正在运行 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| jobId | string |  | JobID 诊断任务 ID |  |
| status | string |  | Status 当前任务状态（pending / running） |  |
| statusUrl | string |  | StatusURL 查询任务状态的接口地址 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {  "jobId": "string",  "statusUrl": "string",  "status": "string"} |  |  |  |

#### 5.1.3.9 启动 MemtestCL测试

**表 5.46 异步启动 MemtestCL接口说明**

| 接口名称 | 异步启动 MemtestCL |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 异步启动 MemtestCL 对所有 HCU 进行测试，测试结果会落盘，可通过 /diag/status/ 查询作业状态和结果 |  |  |  |
| URL | /memtestcl/start |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回作业 ID 和状态查询 URL |  |  |  |
| 409 | 已有正在运行的作业 |  |  |  |
| 500 | 启动作业失败 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| jobId | string |  | JobID 诊断任务 ID |  |
| status | string |  | Status 当前任务状态（pending / running） |  |
| statusUrl | string |  | StatusURL 查询任务状态的接口地址 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {  "jobId": "string",  "statusUrl": "string",  "status": "string"} |  |  |  |

### 5.1.4 Error

#### 5.1.4.1 获取 ECC block 信息

**表 5.47 获取 ECC block 信息接口说明**

| 接口名称 | 获取 ECC block 信息 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 根据设备索引获取 ECC block 信息 |  |  |  |
| URL | /EccBlocksInfo |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | query | Y | 设备索引 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | ECC block 信息 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 内部服务器错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| block | string |  | Block Block 名称 |  |
| ce | integer |  | CE Correctable Error 数量 |  |
| state | string |  | State Block 状态 |  |
| ue | integer |  | UE Uncorrectable Error 数量 |  |
| 示例 |  |  |  |  |
| 请求参数 | dvInd=0 |  |  |  |
| 返回值 | {  "ue": 0,  "block": "string",  "ce": 0,  "state": "string"} |  |  |  |

#### 5.1.4.2 获取HCU块的ECC状态

**表 5.48 获取GPU块的ECC状态接口说明**

| 接口名称 | 获取GPU块的ECC状态 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 返回指定HCU块的ECC状态 |  |  |  |
| URL | /EccStatus/{dvInd} |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 设备索引 |
| block | string | query | Y | GPU块 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回ECC状态 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| eccStatus | string |  | EccStatus ECC状态，字符串类型 |  |
| 示例 |  |  |  |  |
| 请求参数 | block=string |  |  |  |
| 返回值 | {"eccStatus":"string"} |  |  |  |

#### 5.1.4.3 重置指定设备的XGMI错误状态

**表 5.49 重置指定设备的XGMI错误状态接口说明**

| 接口名称 | 重置指定设备的XGMI错误状态 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 无 |  |  |  |
| URL | /ResetXGMIErr |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvIdList | array | body | Y | 设备ID列表 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回失败的设备及其错误信息 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| failedMessages | array:router.FailedMessage |  | FailedMessages 包含重置clock操作失败时的设备ID和错误信息 |  |
| errorMsg | string |  | ErrorMsg 错误信息 |  |
| id | integer |  | ID 设备ID |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {  "failedMessages": [    {      "errorMsg": "string",      "id": 0    }  ]} |  |  |  |

#### 5.1.4.4 获取XGMI错误状态

**表 5.50 获取XGMI错误状态接口说明**

| 接口名称 | 获取XGMI错误状态 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 获取指定物理设备的XGMI（高速互连链路）错误状态。 |  |  |  |
| URL | /XGMIErrorStatus |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 |  |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | query | Y | 物理设备的索引 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回XGMI错误状态码 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| status | integer |  | Status XGMI状态码 |  |
| 示例 |  |  |  |  |
| 请求参数 | dvInd=0 |  |  |  |
| 返回值 | {"status":0} |  |  |  |

#### 5.1.4.5 获取指定设备的XGMI hive id

**表 5.51 获取指定设备的XGMI hive id接口说明**

| 接口名称 | 获取指定设备的XGMI hive id |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 无 |  |  |  |
| URL | /XGMIHiveIdGet |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | query | Y | 设备索引 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回设备的XGMI hive id |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| hiveId | integer |  | HiveId 设备的XGMI hive id |  |
| 示例 |  |  |  |  |
| 请求参数 | dvInd=0 |  |  |  |
| 返回值 | {"hiveId":0} |  |  |  |

#### 5.1.4.6 显示 XGMI 错误状态

**表 5.52 显示 XGMI 错误状态接口说明**

| 接口名称 | 显示 XGMI 错误状态 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 显示一组 GPU 设备的 XGMI 错误状态。 |  |  |  |
| URL | /showXgmiErr |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 |  |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvIdList | array | body | Y | 设备 ID 列表 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | XGMI 错误状态信息 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 |  |  |  |  |

#### 5.1.4.7 获取显存 ECC 信息（bank 级）

**表 5.189 获取显存 ECC 信息（bank 级）接口说明**

| 接口名称 | 获取显存 ECC 信息（bank 级） |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 获取显存 ECC 信息（bank 级） |  |  |  |
| URL | /memeccinfo/:dvInd |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 物理设备索引 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 请求成功 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| data | object |  | 显存 bank 级 ECC 信息 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"data":{}} |  |  |  |

#### 5.1.4.8 设置 RAS 控制

**表 5.192 设置 RAS 控制接口说明**

| 接口名称 | 设置 RAS 控制 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 设置 RAS 控制 |  |  |  |
| URL | /ras/:dvInd |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 物理设备索引 |
| body | object | body | Y | 请求体参数，具体字段以接口实现为准 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 请求成功 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| message | string |  | 操作结果 |  |
| 示例 |  |  |  |  |
| 请求参数 | {} |  |  |  |
| 返回值 | {"message":"成功"} |  |  |  |

#### 5.1.4.9 获取 HSL 链路日志状态

**表 5.193 获取 HSL 链路日志状态接口说明**

| 接口名称 | 获取 HSL 链路日志状态 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 获取 HSL 链路日志状态 |  |  |  |
| URL | /hslstatus/:dvInd |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 物理设备索引 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 请求成功 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| data | object |  | HSL 链路日志状态 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"data":{}} |  |  |  |

### 5.1.5 Event

#### 5.1.5.1 打印设备的事件列表

**表 5.53 打印设备的事件列表接口说明**

| 接口名称 | 打印设备的事件列表 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 打印指定设备的事件列表，并设置延迟 |  |  |  |
| URL | /PrintEventList/{device} |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 |  |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| device | integer | path | Y | 设备索引 |
| delay | integer | query | Y | 延迟时间（秒） |
| eventList | array | query | Y | 事件列表 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 操作成功 |  |  |  |
| 400 | 操作失败 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| 示例 |  |  |  |  |
| 请求参数 | delay=0&eventList=[{}] |  |  |  |
| 返回值 |  |  |  |  |

#### 5.1.5.2 显示设备的事件

**表 5.54 显示设备的事件接口说明**

| 接口名称 | 显示设备的事件 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 获取并显示指定设备的事件信息。 |  |  |  |
| URL | /showEvents |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 |  |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvIdList | array | body | Y | 设备ID列表 |
| eventTypes | array | body | Y | 事件类型列表 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 成功返回设备的事件信息 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| 示例 |  |  |  |  |
| 请求参数 | -d '{"dvIdList":[{}],"eventTypes":[{}]}' |  |  |  |
| 返回值 |  |  |  |  |

### 5.1.6 Health

#### 5.1.6.1 获取设备健康状态

**表 5.55 获取设备健康状态接口说明**

| 接口名称 | 获取设备健康状态 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 根据 groupId 获取设备健康状态，可选择是否检查 HealthConfig |  |  |  |
| URL | /health/group |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json,application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| request | object:router.HealthCheckByGroupReq | body | Y | 请求参数 |
| checkHealthConfig | boolean |  | N | CheckHealthConfig 是否检查 HealthConfig |
| groupId | integer |  | Y | GroupId 设备组ID |
| 状态码 | 描述 |  | 说明 |  |
| 200 | OK |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| deviceHealths | array:router.DeviceHealth |  | DeviceHealths 设备健康结果列表 |  |
| hcu | integer |  | HCU HCU编号 |  |
| status | string |  | Status 状态 |  |
| watches | array:router.SystemWatch |  | Watches 检查详情 |  |
| error | string |  | Error 错误的详细信息 |  |
| result | object |  | Result 结果 |  |
| status | string |  | Status 状态 |  |
| type | string |  | Type 检查项的类型 |  |
| 示例 |  |  |  |  |
| 请求参数 | -d '{"checkHealthConfig":true,"groupId":0}' |  |  |  |
| 返回值 | {  "deviceHealths": [    {      "hcu": 0,      "status": "string",      "watches": [        {          "error": "string",          "result": {},          "status": "string",          "type": "string"        }      ]    }  ]} |  |  |  |

### 5.1.7 HyLink

#### 5.1.7.1 查询指定 HCU 的 HyLink 带宽状态

**表 5.56 查询指定 HCU 的 HyLink 带宽状态接口说明**

| 接口名称 | 查询指定 HCU 的 HyLink 带宽状态 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 根据物理设备索引（DvInd）查询该 HCU 的 HyLink 带宽信息。 返回结果包含： - 该设备所有 Link 的接收 / 发送带宽汇总 - 每条 Link 的接收 / 发送带宽明细 |  |  |  |
| URL | /hylink/status/{dvInd} |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | HCU 设备索引（从 0 开始） |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 查询成功 |  |  |  |
| 400 | 参数错误（非法的 dvInd） |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| deviceLinkBandwidth | object:router.DeviceLinkSum |  | 设备的 HyLink 带宽汇总结果 |  |
| dvInd | integer (int) |  | 设备索引（从 0 开始） |  |
| recv | number (float64) |  | 接收方向（direction=0）所有 Link 的带宽总和 |  |
| send | number (float64) |  | 发送方向（direction=1）所有 Link 的带宽总和 |  |
| err | string |  | 查询过程中可能的错误信息（成功时为空字符串） |  |
| links | array:router.LinkBandwidth |  | 每条 Link 的接收/发送带宽明细 |  |
| linkId | integer (int) |  | Link 索引 |  |
| recv | number (float64) |  | 接收方向（direction=0）的带宽 |  |
| send | number (float64) |  | 发送方向（direction=1）的带宽 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {  "deviceLinkBandwidth": {    "dvInd": 0,    "recv": 1024.5,    "send": 2048.7,    "err": "",    "links": [      {        "linkId": 0,        "recv": 512.2,        "send": 1024.3      },      {        "linkId": 1,        "recv": 512.3,        "send": 1024.4      }    ]  }} |  |  |  |

#### 5.1.7.2 获取 Hy UMC 带宽状态

**表 5.57 获取 Hy UMC 带宽状态接口说明**

| 接口名称 | 获取 Hy UMC 带宽状态 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 查询所有监控设备的 UMC 汇总带宽信息（Read / Write / ReadWrite） |  |  |  |
| URL | /umc/status |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 成功返回设备 UMC 带宽信息 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| deviceUmcBandwidth | array:router.DeviceUmcSum |  | DeviceUmcBandwidth 所有设备的 UMC 带宽汇总信息列表 |  |
| dvInd | integer |  | DvInd 设备索引（从 0 开始） |  |
| err | string |  | Err 错误信息；成功时为空字符串 |  |
| read | number |  | Read 所有 UMC channel 的读带宽之和（失败时为 0） |  |
| readWrite | number |  | ReadWrite 所有 UMC channel 的读写带宽之和（失败时为 0） |  |
| write | number |  | Write 所有 UMC channel 的写带宽之和（失败时为 0） |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {  "deviceUmcBandwidth": [    {      "dvInd": 0,      "err": "string",      "read": 0.0,      "readWrite": 0.0,      "write": 0.0    }  ]} |  |  |  |

### 5.1.8 Memory

#### 5.1.8.1 获取设备的指定内存使用情况

**表 5.58 获取设备的指定内存使用情况接口说明**

| 接口名称 | 获取设备的指定内存使用情况 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 根据设备索引和内存类型返回内存的使用量和总量 |  |  |  |
| URL | /MemInfo |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | query | Y | 设备索引 |
| memType | string | query | Y | 内存类型（可选值: vram, vis_vram, gtt） |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回指定内存类型的使用量和总量 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| memTotal | integer |  | MemTotal 内存总量 |  |
| memUsed | integer |  | MemUsed 已使用的内存量 |  |
| 示例 |  |  |  |  |
| 请求参数 | dvInd=0&memType=string |  |  |  |
| 返回值 | {"memUsed":0,"memTotal":0} |  |  |  |

#### 5.1.8.2 获取内存使用百分比

**表 5.59 获取内存使用百分比接口说明**

| 接口名称 | 获取内存使用百分比 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 根据设备 ID 获取设备内存的使用百分比。 |  |  |  |
| URL | /MemoryPercent/{dvInd} |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 |  |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 设备 ID |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 内存使用百分比 |  |  |  |
| 400 | 请求错误 |  |  |  |
| 404 | 设备未找到 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| busyPercent | integer |  | BusyPercent 内存使用百分比 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"busyPercent":0} |  |  |  |

#### 5.1.8.3 显示设备的退役页信息

**表 5.60 显示设备的退役页信息接口说明**

| 接口名称 | 显示设备的退役页信息 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 获取并显示指定设备的退役内存页信息 |  |  |  |
| URL | /device/retiredpages |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 |  |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvIdList | array | body | Y | 设备ID列表 |
| retiredType | object | body | N | 退役类型 (默认为'all') |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 设备的退役页信息 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| 示例 |  |  |  |  |
| 请求参数 | -d '{"dvIdList":[{}],"retiredType":{}}' |  |  |  |
| 返回值 |  |  |  |  |

#### 5.1.8.4 展示设备内存信息

**表 5.61 展示设备内存信息接口说明**

| 接口名称 | 展示设备内存信息 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 获取并展示指定设备的内存使用情况，包括不同类型的内存。 |  |  |  |
| URL | /memory/info |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 |  |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvIdList | array | body | Y | 设备 ID 列表 |
| memTypes | array | body | Y | 内存类型列表，如 'all' 或指定类型 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 成功返回设备的内存信息 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| 示例 |  |  |  |  |
| 请求参数 | -d '{"dvIdList":[{}],"memTypes":[{}]}' |  |  |  |
| 返回值 |  |  |  |  |

#### 5.1.8.5 展示设备内存使用情况

**表 5.62 展示设备内存使用情况接口说明**

| 接口名称 | 展示设备内存使用情况 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 获取并展示指定设备的当前内存使用百分比和其他相关的利用率数据。 |  |  |  |
| URL | /memory/use |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 |  |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvIdList | array | body | Y | 设备 ID 列表 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 成功返回设备的内存使用信息 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 |  |  |  |  |

#### 5.1.8.6 展示设备的内存供应商信息

**表 5.63 展示设备的内存供应商信息接口说明**

| 接口名称 | 展示设备的内存供应商信息 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 获取并展示指定设备的内存供应商信息 |  |  |  |
| URL | /memory/vendor |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvIdList | array | body | Y | 设备 ID 列表 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 设备内存供应商信息 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| deviceMemVendorInfos | array:router.DeviceMemVendorInfo |  | DeviceMemVendorInfos 设备内存供应商信息列表 |  |
| deviceId | integer |  | DeviceID 设备索引号 |  |
| vendor | string |  | Vendor 内存供应商信息 |  |
| 示例 |  |  |  |  |
| 请求参数 | -d '[{}]' |  |  |  |
| 返回值 | {  "deviceMemVendorInfos": [    {      "deviceId": 0,      "vendor": "string"    }  ]} |  |  |  |

#### 5.1.8.7 查询显存温度

**表 5.168 查询显存温度接口说明**

| 接口名称 | 查询显存温度 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 查询显存温度 |  |  |  |
| URL | /memory/temperature/:dvInd |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 物理设备索引 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 请求成功 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| temperature | number |  | 显存温度 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"temperature":0} |  |  |  |

#### 5.1.8.8 查询可用显存量

**表 5.170 查询可用显存量接口说明**

| 接口名称 | 查询可用显存量 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 查询可用显存量 |  |  |  |
| URL | /memory/available/:dvInd |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 物理设备索引 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 请求成功 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| memoryAvailable | integer |  | 可用显存量（字节） |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"memoryAvailable":0} |  |  |  |

### 5.1.9 PCIe

#### 5.1.9.1 获取可用的 PCIe 带宽列表

**表 5.64 获取可用的 PCIe 带宽列表接口说明**

| 接口名称 | 获取可用的 PCIe 带宽列表 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 根据设备 ID 获取设备的可用 PCIe 带宽列表 |  |  |  |
| URL | /DevPciBandwidth/{dvInd} |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 设备 ID |
| 状态码 | 描述 |  | 说明 |  |
| 200 | PCIe 带宽列表 |  |  |  |
| 400 | 请求错误 |  |  |  |
| 404 | 设备未找到 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| pcieBandwidth | object |  | PCIe 带宽信息 |  |
| transferRate | object |  | PCIe 传输速率信息 |  |
| numSupported | integer (uint32) |  | 支持的 PCIe 速率数量 |  |
| current | integer (uint32) |  | 当前使用的 PCIe 速率索引 |  |
| frequency | array:number (uint64) |  | 支持的 PCIe 速率列表 |  |
| lanes | array:integer (uint32[33]) |  | PCIe 各速率对应的 Lane 数量 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {  "pcieBandwidth": {    "transferRate": {      "numSupported": 4,      "current": 2,      "frequency": [2, 4, 8, 16]    },    "lanes": [16, 16, 8, 8]  }} |  |  |  |

#### 5.1.9.2 设置设备允许的 PCIe 带宽

**表 5.65 设置设备允许的 PCIe 带宽接口说明**

| 接口名称 | 设置设备允许的 PCIe 带宽 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 根据设备索引和带宽掩码限制设备允许的 PCIe 带宽 |  |  |  |
| URL | /DevPciBandwidthSet |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | query | Y | 设备索引 |
| bwBitmask | integer(int64) | query | Y | 带宽掩码 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 操作成功 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| 示例 |  |  |  |  |
| 请求参数 | dvInd=0&bwBitmask=0 |  |  |  |
| 返回值 |  |  |  |  |

#### 5.1.9.3 展示设备的 PCIe 带宽使用情况

**表 5.66 展示设备的 PCIe 带宽使用情况接口说明**

| 接口名称 | 展示设备的 PCIe 带宽使用情况 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 获取并展示指定设备的 PCIe 带宽使用情况，包括发送和接收带宽 |  |  |  |
| URL | /pcie/bandwidth |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvIdList | array | body | Y | 设备 ID 列表 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 设备 PCIe 带宽信息 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| pcieBandwidthInfos | array:router.PcieBandwidthInfo |  | PcieBandwidthInfos PCIe 带宽信息列表 |  |
| bw | number |  | Bw 总带宽 |  |
| dvInd | integer |  | DeviceID 设备索引号 |  |
| received | number |  | Received PCIe 接收带宽 |  |
| sent | number |  | Sent PCIe 发送带宽 |  |
| 示例 |  |  |  |  |
| 请求参数 | -d '[{}]' |  |  |  |
| 返回值 | {  "pcieBandwidthInfos": [    {      "bw": 0.0,      "dvInd": 0,      "received": 0.0,      "sent": 0.0    }  ]} |  |  |  |
|  |  |  |  |  |

#### 5.1.9.4 展示设备的 PCIe 重放计数

**表 5.67 展示设备的 PCIe 重放计数接口说明**

| 接口名称 | 展示设备的 PCIe 重放计数 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 获取并展示指定设备的 PCIe 重放计数 |  |  |  |
| URL | /pcie/replaycount |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 |  |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvIdList | array | body | Y | 设备 ID 列表 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 设备的 PCIe 重放计数信息 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| pcieReplayCountInfos | array:router.PcieReplayCountInfo |  | PcieReplayCountInfos 各设备 PCIe 重放计数信息 |  |
| count | integer |  | Count PCIe 重放总数 |  |
| deviceId | integer |  | DeviceID 设备索引号 |  |
| 示例 |  |  |  |  |
| 请求参数 | -d '[{}]' |  |  |  |
| 返回值 | {  "pcieReplayCountInfos": [    {      "count": 0,      "deviceId": 0    }  ]} |  |  |  |

#### 5.1.9.5 获取总线信息

**表 5.68 获取总线信息接口说明**

| 接口名称 | 获取总线信息 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 获取设备的总线信息 (BDF格式) |  |  |  |
| URL | /picbusinfo/{dvInd} |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | Device ID |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 总线信息 |  |  |  |
| 400 | Invalid device ID |  |  |  |
| 500 | Internal Server Error |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| busInfo | string |  | BusInfo BDF 格式的总线信息 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"busInfo":"string"} |  |  |  |

### 5.1.10 Performance

#### 5.1.10.1 设置 GPU 时钟频率

**表 5.69 设置 GPU 时钟频率接口说明**

| 接口名称 | 设置 GPU 时钟频率 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 设置 GPU 上指定时钟的允许频率。clkType 设置为默认值，无需传递。 |  |  |  |
| URL | /DevGpuClkFreqSet |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | query | Y | 设备索引 |
| freqBitmask | integer(int64) | query | Y | 频率掩码 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 操作结果 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| message | string |  | Message 提示信息 |  |
| 示例 |  |  |  |  |
| 请求参数 | dvInd=0&freqBitmask=0 |  |  |  |
| 返回值 | {"message":"string"} |  |  |  |

#### 5.1.10.2 获取设备的当前性能水平

**表 5.70 获取设备的当前性能水平接口说明**

| 接口名称 | 获取设备的当前性能水平 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 返回指定设备的当前性能等级 |  |  |  |
| URL | /PerfLevel/{dvInd} |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 设备索引 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回当前性能水平 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| perfLevel | string |  | PerfLevel 当前性能等级 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"perfLevel":"string"} |  |  |  |

#### 5.1.10.3 重置设备时钟

**表 5.71 重置设备时钟接口说明**

| 接口名称 | 重置设备时钟 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 重置指定设备的时钟和性能等级为默认值 |  |  |  |
| URL | /ResetClocks |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvIdList | array | body | Y | 设备ID列表 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回失败消息列表 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| failedMessages | array:router.FailedMessage |  | FailedMessages 包含重置clock操作失败时的设备ID和错误信息 |  |
| errorMsg | string |  | ErrorMsg 错误信息 |  |
| id | integer |  | ID 设备ID |  |
| 示例 |  |  |  |  |
| 请求参数 | -d '[{}]' |  |  |  |
| 返回值 | {  "failedMessages": [    {      "errorMsg": "string",      "id": 0    }  ]} |  |  |  |

#### 5.1.10.4 重置Performance Determinism

**表 5.72 重置Performance Determinism接口说明**

| 接口名称 | 重置Performance Determinism |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 该接口用于重置指定设备的性能决定性设置。请求体中需要包含设备ID列表。 |  |  |  |
| URL | /ResetPerfDeterminism |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvIdList | array | body | Y | 设备ID列表 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 成功或失败信息 |  |  |  |
| 400 | 无效的请求体或部分设备失败 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| failedMessages | array:router.FailedMessage |  | FailedMessages 包含操作失败的设备及其错误信息 |  |
| errorMsg | string |  | ErrorMsg 错误信息 |  |
| id | integer |  | ID 设备ID |  |
| 示例 |  |  |  |  |
| 请求参数 | -d '[{}]' |  |  |  |
| 返回值 | {  "failedMessages": [    {      "errorMsg": "string",      "id": 0    }  ] |  |  |  |

#### 5.1.10.5 重置指定设备的电源配置文件和性能级别

**表 5.73 重置指定设备的电源配置文件和性能级别接口说明**

| 接口名称 | 重置指定设备的电源配置文件和性能级别 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 无 |  |  |  |
| URL | /ResetProfile |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvIdList | array | body | Y | 设备ID列表 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回失败的设备及其错误信息 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| failedMessages | array:router.FailedMessage |  | FailedMessages 包含重置clock操作失败时的设备ID和错误信息 |  |
| errorMsg | string |  | ErrorMsg 错误信息 |  |
| id | integer |  | ID 设备ID |  |
| 示例 |  |  |  |  |
| 请求参数 | -d '[{}]' |  |  |  |
| 返回值 | {  "failedMessages": [    {      "errorMsg": "string",      "id": 0    }  ]} |  |  |  |

#### 5.1.10.6 设置设备的时钟OverDrive

**表 5.74 设置设备的时钟OverDrive接口说明**

| 接口名称 | 设置设备的时钟OverDrive |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 设置设备的时钟OverDrive |  |  |  |
| URL | /SetClockOverDrive |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvIdList | array | body | Y | 设备ID列表 |
| clktype | string | query | Y | 时钟类型（sclk 或 mclk） |
| value | string | query | Y | OverDrive值，表示为百分比（0-20%） |
| autoRespond | boolean | query | N | 自动响应超出规格的警告 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 成功设置时钟OverDrive |  |  |  |
| 400 | 返回失败消息列表 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| failedMessages | array:router.FailedMessage |  | FailedMessages 包含操作失败的设备及其错误信息 |  |
| errorMsg | string |  | ErrorMsg 错误信息 |  |
| id | integer |  | ID 设备ID |  |
| 示例 |  |  |  |  |
| 请求参数 | clktype=string&value=string&autoRespond=true -d '[{}]' |  |  |  |
| 返回值 | {  "failedMessages": [    {      "errorMsg": "string",      "id": 0    }  ]} |  |  |  |

#### 5.1.10.7 设置设备的时钟频率范围

**表 5.75 设置设备的时钟频率范围接口说明**

| 接口名称 | 设置设备的时钟频率范围 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 设置设备的时钟频率范围（sclk 或 mclk） |  |  |  |
| URL | /SetClockRange |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvIdList | array | body | Y | 设备ID列表 |
| clkType | string | query | Y | 时钟类型（sclk 或 mclk） |
| minvalue | string | query | Y | 最小值（MHz） |
| maxvalue | string | query | Y | 最大值（MHz） |
| autoRespond | boolean | query | N | 自动响应超出规格的警告 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回失败信息或成功 |  |  |  |
| 400 | 无效的请求参数或无法设置时钟范围 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| failedMessages | array:router.FailedMessage |  | FailedMessages 包含操作失败的设备及其错误信息 |  |
| errorMsg | string |  | ErrorMsg 错误信息 |  |
| id | integer |  | ID 设备ID |  |
| 示例 |  |  |  |  |
| 请求参数 | clkType=string&minvalue=string&maxvalue=string&autoRespond=true -d '[{}]' |  |  |  |
| 返回值 | {  "failedMessages": [    {      "errorMsg": "string",      "id": 0    }  ]} |  |  |  |

#### 5.1.10.8 设置设备的性能确定性

**表 5.76 设置设备的性能确定性接口说明**

| 接口名称 | 设置设备的性能确定性 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 设置设备的性能确定性 |  |  |  |
| URL | /SetPerfDeterminism |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvIdList | array | body | Y | 设备ID列表 |
| clkvalue | string | query | Y | 时钟频率值 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回失败的设备及其错误信息 |  |  |  |
| 400 | 无效的请求体或无法设置性能确定性 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| errorMsg | string |  | ErrorMsg 错误信息 |  |
| id | integer |  | ID 设备ID |  |
| 示例 |  |  |  |  |
| 请求参数 | clkvalue=string -d '[{}]' |  |  |  |
| 返回值 | {  "id": 0,  "errorMsg": "string"} |  |  |  |

#### 5.1.10.9 设置设备性能等级

**表 5.77 设置设备性能等级接口说明**

| 接口名称 | 设置设备性能等级 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 根据设备ID列表和给定的性能等级，设置设备的性能等级 |  |  |  |
| URL | /SetPerformanceLevel |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| deviceList | array | body | Y | 设备 ID 列表 |
| level | string | query | Y | 性能等级 (auto, low, high, normal) |
| 状态码 | 描述 |  | 说明 |  |
| 200 | OK |  |  |  |
| 400 | Bad Request |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| errorMsg | string |  | ErrorMsg 错误信息 |  |
| id | integer |  | ID 设备ID |  |
| 示例 |  |  |  |  |
| 请求参数 | level=string -d '[{}]' |  |  |  |
| 返回值 | {  "id": 0,  "errorMsg": "string"} |  |  |  |

#### 5.1.10.10 设置设备的PowerPlay表级别

**表 5.78 设置设备的PowerPlay表级别接口说明**

| 接口名称 | 设置设备的PowerPlay表级别 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 设置设备的PowerPlay表级别 |  |  |  |
| URL | /SetPowerPlayTableLevel |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvIdList | array | body | Y | 设备ID列表 |
| clkType | string | query | Y | 时钟类型（sclk 或 mclk） |
| point | string | query | Y | 电压点 |
| clk | string | query | Y | 时钟值（MHz） |
| volt | string | query | Y | 电压值（mV） |
| autoRespond | boolean | query | N | 自动响应超出规格的警告 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 成功设置PowerPlay表级别 |  |  |  |
| 400 | 返回失败消息列表 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| failedMessages | array:router.FailedMessage |  | FailedMessages 包含操作失败的设备及其错误信息 |  |
| errorMsg | string |  | ErrorMsg 错误信息 |  |
| id | integer |  | ID 设备ID |  |
| 示例 |  |  |  |  |
| 请求参数 | clkType=string&point=string&clk=string&volt=string&autoRespond=true -d '[{}]' |  |  |  |
| 返回值 | {  "failedMessages": [    {      "errorMsg": "string",      "id": 0    }  ]} |  |  |  |

#### 5.1.10.11 显示时钟信息

**表 5.79 显示时钟信息接口说明**

| 接口名称 | 显示时钟信息 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 显示指定设备的时钟信息 |  |  |  |
| URL | /ShowClocks |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvIdList | array | body | Y | 设备 ID 列表 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 成功信息 |  |  |  |
| 400 | 失败信息 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| 示例 |  |  |  |  |
| 请求参数 | -d '[{}]' |  |  |  |
| 返回值 |  |  |  |  |

#### 5.1.10.12 获取设备粗粒度利用率

**表 5.80 获取设备粗粒度利用率接口说明**

| 接口名称 | 获取设备粗粒度利用率 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 无 |  |  |  |
| URL | /utilization/coarse |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| body | object:router.GetCoarseGrainUtilReq | body | Y | 请求参数 |
| device | integer |  | N | Device 设备 ID |
| typeName | string |  | N | TypeName 利用率计数器类型名称（可选） |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 利用率计数器列表 |  |  |  |
| 400 | 错误信息 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| utilizationCounters | array:router.UtilizationCounter |  | UtilizationCounters 利用率计数器列表 |  |
| type | integer |  | Type 利用率计数器类型 |  |
| value | integer |  | Value 利用率值 |  |
| 示例 |  |  |  |  |
| 请求参数 | -d '{"device":0,"typeName":"string"}' |  |  |  |
| 返回值 | {  "utilizationCounters": [    {      "type": 0,      "value": 0    }  ]} |  |  |  |

#### 5.1.10.13 查询设备降频状态

**表 5.167 查询设备降频状态接口说明**

| 接口名称 | 查询设备降频状态 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 查询设备降频状态 |  |  |  |
| URL | /throttle/flags/:dvInd |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 物理设备索引 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 请求成功 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| data | object |  | 设备降频状态 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"data":{}} |  |  |  |

#### 5.1.10.14 查询显存当前频率

**表 5.171 查询显存当前频率接口说明**

| 接口名称 | 查询显存当前频率 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 查询显存当前频率 |  |  |  |
| URL | /clock/mem/:dvInd |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 物理设备索引 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 请求成功 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| clock | number |  | 显存当前频率 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"clock":0} |  |  |  |

#### 5.1.10.15 查询 GFX 最大额定频率

**表 5.172 查询 GFX 最大额定频率接口说明**

| 接口名称 | 查询 GFX 最大额定频率 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 查询 GFX 最大额定频率 |  |  |  |
| URL | /clock/gfx/max/:dvInd |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 物理设备索引 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 请求成功 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| clock | number |  | GFX 最大额定频率 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"clock":0} |  |  |  |

#### 5.1.10.16 查询显存最大额定频率

**表 5.173 查询显存最大额定频率接口说明**

| 接口名称 | 查询显存最大额定频率 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 查询显存最大额定频率 |  |  |  |
| URL | /clock/mem/max/:dvInd |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 物理设备索引 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 请求成功 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| clock | number |  | 显存最大额定频率 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"clock":0} |  |  |  |

#### 5.1.10.17 查询显存 overdrive 百分比

**表 5.176 查询显存 overdrive 百分比接口说明**

| 接口名称 | 查询显存 overdrive 百分比 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 查询显存 overdrive 百分比 |  |  |  |
| URL | /device/memoverdrive/:dvInd |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 物理设备索引 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 请求成功 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| overdrive | number |  | 显存 overdrive 百分比 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"overdrive":0} |  |  |  |

#### 5.1.10.18 查询设备 CU wave 占用周期占比

**表 5.179 查询设备 CU wave 占用周期占比接口说明**

| 接口名称 | 查询设备 CU wave 占用周期占比 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 查询设备 CU wave 占用周期占比 |  |  |  |
| URL | /device/cuutil/:dvInd |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 物理设备索引 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 请求成功 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| utilization | number |  | CU wave 占用周期占比 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"utilization":0} |  |  |  |

#### 5.1.10.19 查询设备 wave 驻留数量占比

**表 5.180 查询设备 wave 驻留数量占比接口说明**

| 接口名称 | 查询设备 wave 驻留数量占比 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 查询设备 wave 驻留数量占比 |  |  |  |
| URL | /device/waveutil/:dvInd |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 物理设备索引 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 请求成功 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| utilization | number |  | Wave 驻留数量占比 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"utilization":0} |  |  |  |

#### 5.1.10.20 设置设备性能等级（v1）

**表 5.185 设置设备性能等级（v1）接口说明**

| 接口名称 | 设置设备性能等级（v1） |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 设置设备性能等级（v1） |  |  |  |
| URL | /perflevel/v1/:dvInd |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 物理设备索引 |
| body | object | body | Y | 请求体参数，具体字段以接口实现为准 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 请求成功 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| message | string |  | 操作结果 |  |
| 示例 |  |  |  |  |
| 请求参数 | {} |  |  |  |
| 返回值 | {"message":"成功"} |  |  |  |

#### 5.1.10.21 获取 GFX boost 状态

**表 5.187 获取 GFX boost 状态接口说明**

| 接口名称 | 获取 GFX boost 状态 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 获取 GFX boost 状态 |  |  |  |
| URL | /boost/:dvInd |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 物理设备索引 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 请求成功 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| data | object |  | GFX boost 状态 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"data":{}} |  |  |  |

#### 5.1.10.22 获取频率限制历史

**表 5.188 获取频率限制历史接口说明**

| 接口名称 | 获取频率限制历史 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 获取频率限制历史 |  |  |  |
| URL | /fdrhistory/:dvInd |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 物理设备索引 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 请求成功 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| data | array |  | 频率限制历史 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"data":[]} |  |  |  |

#### 5.1.10.23 设置 GFX boost 等级

**表 5.190 设置 GFX boost 等级接口说明**

| 接口名称 | 设置 GFX boost 等级 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 设置 GFX boost 等级 |  |  |  |
| URL | /boost/:dvInd |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 物理设备索引 |
| body | object | body | Y | 请求体参数，具体字段以接口实现为准 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 请求成功 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| message | string |  | 操作结果 |  |
| 示例 |  |  |  |  |
| 请求参数 | {} |  |  |  |
| 返回值 | {"message":"成功"} |  |  |  |

### 5.1.11 Physical State

#### 5.1.11.1 获取设备的风扇速度

**表 5.81 获取设备的风扇速度接口说明**

| 接口名称 | 获取设备的风扇速度 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 获取指定设备的风扇速度（RPM） |  |  |  |
| URL | /DevFanRpms/{dvInd} |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 设备索引 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回风扇速度 |  |  |  |
| 400 | 请求参数错误或获取失败 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| speed | integer |  | Speed 风扇速度，单位 RPM |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"speed":0} |  |  |  |

#### 5.1.11.2 复位风扇控制

**表 5.82 复位风扇控制接口说明**

| 接口名称 | 复位风扇控制 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 重置指定设备的风扇控制为默认值 |  |  |  |
| URL | /ResetFans |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvIdList | array | body | Y | 设备ID列表 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 复位成功 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| 示例 |  |  |  |  |
| 请求参数 | -d '[{}]' |  |  |  |
| 返回值 |  |  |  |  |

#### 5.1.11.3 设置风扇转速

**表 5.83 设置风扇转速接口说明**

| 接口名称 | 设置风扇转速 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 根据设备ID列表和给定的风扇速度，设置设备的风扇速度 |  |  |  |
| URL | /SetFanSpeed |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvIdList | array | body | Y | 设备 ID 列表 |
| fan | string | query | Y | 风扇速度值（0-255,单位:RPM） |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 成功信息 |  |  |  |
| 400 | 失败信息 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| 示例 |  |  |  |  |
| 请求参数 | fan=string -d '[{}]' |  |  |  |
| 返回值 |  |  |  |  |

#### 5.1.11.4 获取设备温度

**表 5.84 获取设备温度接口说明**

| 接口名称 | 获取设备温度 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 返回指定设备的当前温度 |  |  |  |
| URL | /Temperature/{dvInd} |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 设备索引 |
| sensorType | integer | query | Y | 传感器类型: 0: Edge GPU temperature; 1: Junction/hotspot temperature; 2: VRAM temperature |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回设备温度 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| temp | number |  | Temp 当前温度（摄氏度） |  |
| 示例 |  |  |  |  |
| 请求参数 | sensorType=0 |  |  |  |
| 返回值 | {"temp":0.0} |  |  |  |

#### 5.1.11.5 展示风扇转速和风扇级别

**表 5.85 展示风扇转速和风扇级别接口说明**

| 接口名称 | 展示风扇转速和风扇级别 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 显示指定设备的当前风扇转速和风扇级别 |  |  |  |
| URL | /fans/current |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvIdList | array | body | Y | 设备 ID 列表 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 成功信息 |  |  |  |
| 400 | 失败信息 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| 示例 |  |  |  |  |
| 请求参数 | -d '[{}]' |  |  |  |
| 返回值 |  |  |  |  |

#### 5.1.11.6 获取风扇转速

**表 5.86 获取风扇转速接口说明**

| 接口名称 | 获取风扇转速 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 获取指定设备的风扇转速及其占最大转速的百分比 |  |  |  |
| URL | /fanspeedinfo/{dvInd} |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | Device ID |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 风扇转速信息 |  |  |  |
| 400 | Invalid device ID |  |  |  |
| 500 | Internal Server Error |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| fanLevel | integer |  | FanLevel 风扇转速 |  |
| fanPercentage | number |  | FanPercentage 风扇转速占最大转速的百分比 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {  "fanLevel": 0,  "fanPercentage": 0.0} |  |  |  |

#### 5.1.11.7 获取HCU使用率

**表 5.87 获取HCU使用率接口说明**

| 接口名称 | 获取HCU使用率 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 获取指定设备的HCU当前使用百分比 |  |  |  |
| URL | /gpuuse/{dvInd} |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | Device ID |
| 状态码 | 描述 |  | 说明 |  |
| 200 | HCU 使用率信息 |  |  |  |
| 400 | Invalid device ID |  |  |  |
| 500 | Internal Server Error |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| gpuUsage | integer |  | GPUUsage HCU 当前使用百分比 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"gpuUsage":0} |  |  |  |

#### 5.1.11.8 显示设备温度传感器数据

**表 5.88 显示设备温度传感器数据接口说明**

| 接口名称 | 显示设备温度传感器数据 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 获取指定设备列表的温度传感器数据 |  |  |  |
| URL | /temps/current |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvIdList | array | body | Y | 设备 ID 列表 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回设备温度信息列表 |  |  |  |
| 400 | 请求参数错误或获取温度失败 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| temperatureInfos | array:router.TemperatureInfo |  | TemperatureInfos 温度信息列表 |  |
| deviceID | integer |  | DeviceID 设备索引号 |  |
| sensorTemps | object |  | SensorTemps 传感器名称到温度的映射 key 为传感器名称（如 "Edge", "Junction", "VRAM"） value 为对应温度（摄氏度） |  |
| 示例 |  |  |  |  |
| 请求参数 | -d '[{}]' |  |  |  |
| 返回值 | {  "temperatureInfos": [    {      "deviceID": 0,      "sensorTemps": {}    }  ]} |  |  |  |

#### 5.1.11.9 查询板卡边缘温度

**表 5.169 查询板卡边缘温度接口说明**

| 接口名称 | 查询板卡边缘温度 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 查询板卡边缘温度 |  |  |  |
| URL | /board/temperature/:dvInd |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 物理设备索引 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 请求成功 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| temperature | number |  | 板卡边缘温度 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"temperature":0} |  |  |  |

### 5.1.12 Power

#### 5.1.12.1 设置设备 PowerPlay 性能级别

**表 5.89 设置设备 PowerPlay 性能级别接口说明**

| 接口名称 | 设置设备 PowerPlay 性能级别 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 根据设备 ID 设置 PowerPlay 性能级别。 |  |  |  |
| URL | /DevPerfLevelSet/{dvInd} |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 |  |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 设备 ID |
| level | string | query | Y | 要设置的性能级别 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 操作成功 |  |  |  |
| 400 | 请求错误 |  |  |  |
| 404 | 设备未找到 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| 示例 |  |  |  |  |
| 请求参数 | level=string |  |  |  |
| 返回值 |  |  |  |  |

#### 5.1.12.2 设置设备功率配置文件

**表 5.90 设置设备功率配置文件接口说明**

| 接口名称 | 设置设备功率配置文件 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 设置指定设备的功率配置文件 |  |  |  |
| URL | /DevPowerProfileSet |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 设备索引 |
| reserved | integer | query | Y | 保留参数，通常为0 |
| profile | integer | query | Y | 功率配置文件的枚举值 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 成功信息 |  |  |  |
| 400 | 失败信息 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| 示例 |  |  |  |  |
| 请求参数 | reserved=0&profile=0 |  |  |  |
| 返回值 |  |  |  |  |

#### 5.1.12.3 获取设备的最大功率

**表 5.91 获取设备的最大功率接口说明**

| 接口名称 | 获取设备的最大功率 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 根据设备索引返回设备的最大功率（以瓦特为单位） |  |  |  |
| URL | /MaxPower |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | query | Y | 设备索引 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 设备最大功率 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| power | integer |  | Power 设备的最大功率（瓦特） |  |
| 示例 |  |  |  |  |
| 请求参数 | dvInd=0 |  |  |  |
| 返回值 | {"power":0} |  |  |  |

#### 5.1.12.4 获取设备的平均功耗

**表 5.92 获取设备的平均功耗接口说明**

| 接口名称 | 获取设备的平均功耗 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 返回指定设备的平均功耗（瓦特） |  |  |  |
| URL | /Power/{dvInd} |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 设备索引 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回平均功耗 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| power | integer |  | Power 平均功耗，单位瓦特 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"power":0} |  |  |  |

#### 5.1.12.5 设置功率配置

**表 5.93 设置功率配置接口说明**

| 接口名称 | 设置功率配置 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 根据设备ID列表和给定的功率配置文件，设置设备的功率配置 |  |  |  |
| URL | /SetProfile |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvIdList | array | body | Y | 设备 ID 列表 |
| profile | string | query | Y | 功率配置文件名称:CUSTOM、VIDEO、POWER SAVING、COMPUTE、VR、3D FULL SCREEN、BOOTUP DEFAULT |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 设置成功的消息列表 |  |  |  |
| 400 | 失败的消息列表 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| errorMsg | string |  | ErrorMsg 错误信息 |  |
| id | integer |  | ID 设备ID |  |
| 示例 |  |  |  |  |
| 请求参数 | profile=string -d '[{}]' |  |  |  |
| 返回值 | {  "id": 0,  "errorMsg": "string"} |  |  |  |

#### 5.1.12.6 展示设备的平均功率消耗

**表 5.94 展示设备的平均功率消耗接口说明**

| 接口名称 | 展示设备的平均功率消耗 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 获取并展示指定设备的平均图形功率消耗 |  |  |  |
| URL | /device/power |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 |  |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvIdList | array | body | Y | 设备 ID 列表 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 设备的功率信息 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| deviceId | integer |  | DeviceID 设备索引号 |  |
| power | integer |  | Power 设备功率 |  |
| 示例 |  |  |  |  |
| 请求参数 | -d '[{}]' |  |  |  |
| 返回值 | {  "power": 0,  "deviceId": 0} |  |  |  |

#### 5.1.12.7 展示设备的HCU内存时钟频率和电压

**表 5.95 展示设备的HCU内存时钟频率和电压接口说明**

| 接口名称 | 展示设备的HCU内存时钟频率和电压 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 获取并展示指定设备的HCU内存时钟频率和电压表 |  |  |  |
| URL | /device/powerplay |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 |  |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvIdList | array | body | Y | 设备 ID 列表 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 设备的HCU时钟频率和电压信息 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| ddcCurve | array |  | DDC_CURVE 电压曲线 |  |
| deviceId | integer |  | DeviceID 设备索引号 |  |
| mclk | string |  | MCLK MCLK 时钟值 |  |
| range | array |  | RANGE 时钟范围 |  |
| sclk | array |  | SCLK SCLK 时钟列表 |  |
| 示例 |  |  |  |  |
| 请求参数 | -d '[{}]' |  |  |  |
| 返回值 | {  "range": [    {}  ],  "sclk": [    {}  ],  "deviceId": 0,  "mclk": "string",  "ddcCurve": [    {}  ]} |  |  |  |

#### 5.1.12.8 显示设备的电源配置文件

**表 5.96 显示设备的电源配置文件接口说明**

| 接口名称 | 显示设备的电源配置文件 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 获取并显示指定设备的电源配置文件，包括可用的电源配置选项 |  |  |  |
| URL | /device/profile |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 |  |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvIdList | array | body | Y | 设备 ID 列表 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 设备的电源配置文件信息列表 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| deviceId | integer |  | DeviceID 设备索引号 example: 0 |  |
| profiles | array |  | Profiles 文件信息列表 example: ["Profile1", "Profile2"] |  |
| 示例 |  |  |  |  |
| 请求参数 | -d '[{}]' |  |  |  |
| 返回值 | {  "profiles": [    {}  ],  "deviceId": 0} |  |  |  |

#### 5.1.12.9 显示设备的电流或电压范围（K100_AI卡不支持该操作）

**表 5.97 显示设备的电流或电压范围（K100_AI卡不支持该操作）接口说明**

| 接口名称 | 显示设备的电流或电压范围（K100_AI卡不支持该操作） |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 获取并显示指定设备的有效电流或电压范围 |  |  |  |
| URL | /device/range |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 |  |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvIdList | array | body | Y | 设备ID列表 |
| rangeType | object | body | Y | 范围类型 (sclk, mclk, voltage) |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 设备的电流或电压范围信息 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| 示例 |  |  |  |  |
| 请求参数 | -d '{"dvIdList":[{}],"rangeType":{}}' |  |  |  |
| 返回值 |  |  |  |  |

#### 5.1.12.10 展示设备的能量消耗

**表 5.98 展示设备的能量消耗接口说明**

| 接口名称 | 展示设备的能量消耗 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 获取并展示指定设备的能量消耗情况。 |  |  |  |
| URL | /energy |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 |  |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvIdList | array | body | Y | 设备 ID 列表 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 成功返回设备的能量消耗信息 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| 示例 |  |  |  |  |
| 请求参数 | -d '[{}]' |  |  |  |
| 返回值 |  |  |  |  |

#### 5.1.12.11 显示设备的电压信息

**表 5.99 显示设备的电压信息接口说明**

| 接口名称 | 显示设备的电压信息 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 获取并显示指定设备的当前电压信息。 |  |  |  |
| URL | /showVoltage |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 |  |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvIdList | array | body | Y | 设备ID列表 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 设备的电压信息列表 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| deviceId | integer |  | DeviceID 设备索引号 example: 0 |  |
| voltage | integer |  | Voltage 电压信息 example: 120 |  |
| 示例 |  |  |  |  |
| 请求参数 | -d '[{}]' |  |  |  |
| 返回值 | {  "deviceId": 0,  "voltage": 0} |  |  |  |

#### 5.1.12.12 显示设备的电压曲线点

**表 5.100 显示设备的电压曲线点接口说明**

| 接口名称 | 显示设备的电压曲线点 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 获取并显示指定设备的电压曲线点信息。 |  |  |  |
| URL | /showVoltageCurve |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 |  |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvIdList | array | body | Y | 设备ID列表 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 设备的电压曲线点信息 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| 示例 |  |  |  |  |
| 请求参数 | -d '[{}]' |  |  |  |
| 返回值 |  |  |  |  |

#### 5.1.12.13 设置设备功率上限

**表 5.183 设置设备功率上限接口说明**

| 接口名称 | 设置设备功率上限 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 设置设备功率上限 |  |  |  |
| URL | /powercap/:dvInd |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 物理设备索引 |
| body | object | body | Y | 请求体参数，具体字段以接口实现为准 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 请求成功 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| message | string |  | 操作结果 |  |
| 示例 |  |  |  |  |
| 请求参数 | {} |  |  |  |
| 返回值 | {"message":"成功"} |  |  |  |

#### 5.1.12.14 复位设备功率上限

**表 5.184 复位设备功率上限接口说明**

| 接口名称 | 复位设备功率上限 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 复位设备功率上限 |  |  |  |
| URL | /powercap/:dvInd/reset |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 物理设备索引 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 请求成功 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| message | string |  | 操作结果 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"message":"成功"} |  |  |  |

#### 5.1.12.15 设置低功耗控制

**表 5.191 设置低功耗控制接口说明**

| 接口名称 | 设置低功耗控制 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 设置低功耗控制 |  |  |  |
| URL | /lowpower/:dvInd |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 物理设备索引 |
| body | object | body | Y | 请求体参数，具体字段以接口实现为准 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 请求成功 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| message | string |  | 操作结果 |  |
| 示例 |  |  |  |  |
| 请求参数 | {} |  |  |  |
| 返回值 | {"message":"成功"} |  |  |  |

### 5.1.13 Process

#### 5.1.13.1 获取指定进程的计算信息

**表 5.101 获取指定进程的计算信息接口说明**

| 接口名称 | 获取指定进程的计算信息 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 根据进程 PID 查询该进程的计算资源使用情况，包括显存使用量、SDMA 使用量和 CU 占用率 |  |  |  |
| URL | /process/info |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| pid | integer | query | Y | 进程 PID |
| 状态码 | 描述 |  | 说明 |  |
| 200 | OK |  |  |  |
| 400 | Bad Request |  |  |  |
| 500 | Internal Server Error |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| data | object |  | 接口返回数据 |  |
| processInfo | object |  | 指定进程的计算资源使用信息 |  |
| processIDmessage | integer (uint32)string |  | 进程 ID |  |
| pasid | integer (uint32) |  | 进程地址空间 ID |  |
| vramUsage | integer (uint64) |  | 显存使用量（字节） |  |
| sdmaUsage | integer (uint64) |  | SDMA 使用量 |  |
| cuOccupancy | integer (uint32) |  | CU 占用率 |  |
| error | object |  | 错误信息（失败时返回） |  |
| message | string |  | 接口返回描述信息 |  |
| 示例 |  |  |  |  |
| 请求参数 | pid=0 |  |  |  |
| 返回值 | {  "data": {    "processInfo": {      "processID": 12345,      "pasid": 678,      "vramUsage": 2147483648,      "sdmaUsage": 102400,      "cuOccupancy": 75    }  },  "error": {},  "message": "success"} |  |  |  |

#### 5.1.13.2 查询进程各 GPU 占用率

**表 5.182 查询进程各 GPU 占用率接口说明**

| 接口名称 | 查询进程各 GPU 占用率 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 查询进程各 GPU 占用率（近似 DCGM_FI_PROC_ENC_UTIL） |  |  |  |
| URL | /process/gpuusage/:pid |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| pid | integer | path | Y | 进程 ID |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 请求成功 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| data | array |  | 进程各 GPU 占用率 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"data":[]} |  |  |  |

### 5.1.14 System

#### 5.1.14.1 获取设备总线信息

**表 5.102 获取设备总线信息接口说明**

| 接口名称 | 获取设备总线信息 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 获取指定设备的总线信息 |  |  |  |
| URL | /GetBus/{dvInd} |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 设备索引 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回设备总线信息 |  |  |  |
| 400 | 请求参数错误或获取失败 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| picId | string |  | PicID 设备总线ID |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"picId":"string"} |  |  |  |

#### 5.1.14.2 获取指定PID的进程名

**表 5.103 获取指定PID的进程名接口说明**

| 接口名称 | 获取指定PID的进程名 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 根据进程ID（PID）返回对应的进程名称 |  |  |  |
| URL | /ProcessName |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| pid | integer | query | Y | 进程ID |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回进程名称 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| processName | string |  | ProcessName 进程名称 |  |
| 示例 |  |  |  |  |
| 请求参数 | pid=0 |  |  |  |
| 返回值 | {"processName":"string"} |  |  |  |

#### 5.1.14.3 展示系统中正在运行的KFD进程信息

**表 5.104 展示系统中正在运行的KFD进程信息接口说明**

| 接口名称 | 展示系统中正在运行的KFD进程信息 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 获取并展示当前系统中运行的KFD进程的详细信息。 |  |  |  |
| URL | /pids |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 |  |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 成功返回进程信息 |  |  |  |
| 400 | 请求错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 |  |  |  |  |

#### 5.1.14.4 获取进程列表信息

**表 5.105 获取进程列表信息接口说明**

| 接口名称 | 获取进程列表信息 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 获取并返回进程信息和使用的 HCU 设备信息 |  |  |  |
| URL | /process-info |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 进程信息列表 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| processes | array:router.Process |  | Processes 进程信息列表 |  |
| cuOccupancy | integer |  | CuOccupancy CU 占用率 |  |
| minorNumbers | array |  | MinorNumbers 使用的设备索引号 |  |
| pasid | integer |  | Pasid 进程地址空间 ID |  |
| processId | integer |  | ProcessID 进程 ID |  |
| processName | string |  | ProcessName 进程名称 |  |
| sdmaUsage | integer |  | SdmaUsage SDMA 使用量 |  |
| vramUsage | integer |  | VramUsage 显存使用量（单位：Bytes） |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {  "processes": [    {      "cuOccupancy": 0,      "minorNumbers": [        {}      ],      "pasid": 0,      "processId": 0,      "processName": "string",      "sdmaUsage": 0,      "vramUsage": 0    }  ]} |  |  |  |

#### 5.1.14.5 获取计算进程列表

**表 5.106 获取计算进程列表接口说明**

| 接口名称 | 获取计算进程列表 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 返回当前系统计算进程的 PID 列表 |  |  |  |
| URL | /process/list |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 进程 ID 列表 |  |  |  |
| 400 | 获取进程列表失败 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| pidList | array |  | PidList 计算进程 ID 列表 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"pidList":[{}]} |  |  |  |

#### 5.1.14.6 进程列表信息

**表 5.107 进程列表信息接口说明**

| 接口名称 | 进程列表信息 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 获取系统中计算进程的信息及其使用的 HCU 设备索引 |  |  |  |
| URL | /processHCUInfo |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 成功返回进程信息 |  |  |  |
| 400 | 请求错误或获取进程信息失败 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| processInfo | array:router.Process |  | ProcessInfo 进程信息列表 |  |
| cuOccupancy | integer |  | CuOccupancy CU 占用率 |  |
| minorNumbers | array |  | MinorNumbers 使用的设备索引号 |  |
| pasid | integer |  | Pasid 进程地址空间 ID |  |
| processId | integer |  | ProcessID 进程 ID |  |
| processName | string |  | ProcessName 进程名称 |  |
| sdmaUsage | integer |  | SdmaUsage SDMA 使用量 |  |
| vramUsage | integer |  | VramUsage 显存使用量（单位：Bytes） |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {  "processInfo": [    {      "cuOccupancy": 0,      "minorNumbers": [        {}      ],      "pasid": 0,      "processId": 0,      "processName": "string",      "sdmaUsage": 0,      "vramUsage": 0    }  ]} |  |  |  |

#### 5.1.14.7 获取指定 PID 的进程名

**表 5.164 获取指定 PID 的进程名接口说明**

| 接口名称 | 获取指定 PID 的进程名 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 获取指定 PID 的进程名 |  |  |  |
| URL | /process/name/:pid |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| pid | integer | path | Y | 进程 ID |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 请求成功 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| processName | string |  | 进程名称 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"processName":"string"} |  |  |  |

### 5.1.15 Topo

#### 5.1.15.1 显示 GPU 拓扑跳数

**表 5.108 显示 GPU 拓扑跳数接口说明**

| 接口名称 | 显示 GPU 拓扑跳数 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 显示 GPU 设备之间的跳数信息。 |  |  |  |
| URL | /showHopsTopology |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 |  |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvIdList | array | body | Y | 设备 ID 列表 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | GPU 拓扑跳数信息 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| 示例 |  |  |  |  |
| 请求参数 | -d '[{}]' |  |  |  |
| 返回值 |  |  |  |  |

#### 5.1.15.2 显示完整的硬件拓扑信息

**表 5.109 显示完整的硬件拓扑信息接口说明**

| 接口名称 | 显示完整的硬件拓扑信息 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 显示一组 GPU 设备的权重、跳数、链接类型和 NUMA 节点信息。 |  |  |  |
| URL | /showHwTopology |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 |  |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvIdList | array | body | Y | 设备 ID 列表 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 完整的硬件拓扑信息 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| 示例 |  |  |  |  |
| 请求参数 | -d '[{}]' |  |  |  |
| 返回值 |  |  |  |  |

#### 5.1.15.3 显示 NUMA 节点信息

**表 5.110 显示 NUMA 节点信息接口说明**

| 接口名称 | 显示 NUMA 节点信息 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 显示一组 GPU 设备的 NUMA 节点和关联信息。 |  |  |  |
| URL | /showNumaTopology |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 |  |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvIdList | array | body | Y | 设备 ID 列表 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | NUMA 节点信息列表 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| deviceId | integer |  | DeviceID 设备索引号 example: 0 |  |
| numaAffinity | integer |  | NumaAffinity NUMA 关联信息 example: 255 |  |
| numaNode | integer |  | NumaNode NUMA 节点 example: 1 |  |
| 示例 |  |  |  |  |
| 请求参数 | -d '[{}]' |  |  |  |
| 返回值 | {  "numaNode": 0,  "numaAffinity": 0,  "deviceId": 0} |  |  |  |

#### 5.1.15.4 显示 GPU 拓扑链接类型

**表 5.111 显示 GPU 拓扑链接类型接口说明**

| 接口名称 | 显示 GPU 拓扑链接类型 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 显示 GPU 设备之间的链接类型信息。 |  |  |  |
| URL | /showTypeTopology |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 |  |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvIdList | array | body | Y | 设备 ID 列表 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | GPU 拓扑链接类型信息 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| 示例 |  |  |  |  |
| 请求参数 | -d '[{}]' |  |  |  |
| 返回值 |  |  |  |  |

#### 5.1.15.5 显示 GPU 拓扑权重

**表 5.112 显示 GPU 拓扑权重接口说明**

| 接口名称 | 显示 GPU 拓扑权重 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 显示 GPU 设备之间的权重信息。 |  |  |  |
| URL | /showWeightTopology |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 |  |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvIdList | array | body | Y | 设备 ID 列表 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | GPU 拓扑权重信息 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| 示例 |  |  |  |  |
| 请求参数 | -d '[{}]' |  |  |  |
| 返回值 |  |  |  |  |

#### 5.1.15.6 显示整机 HCU 互联矩阵信息

**表 5.165 显示整机 HCU 互联矩阵信息接口说明**

| 接口名称 | 显示整机 HCU 互联矩阵信息 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 显示整机 HCU 互联矩阵信息 |  |  |  |
| URL | /discoverInterconnectTopology |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 请求成功 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| data | array |  | HCU 互联拓扑信息 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"data":[]} |  |  |  |

#### 5.1.15.7 查询指定设备的 HSL/XHCL 链路通断状态

**表 5.166 查询指定设备的 HSL/XHCL 链路通断状态接口说明**

| 接口名称 | 查询指定设备的 HSL/XHCL 链路通断状态 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 查询指定设备的 HSL/XHCL 链路通断状态 |  |  |  |
| URL | /hsl/linkstate/:dvInd |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 物理设备索引 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 请求成功 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| data | object |  | HSL/XHCL 链路状态 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"data":{}} |  |  |  |

#### 5.1.15.8 查询两设备间 P2P 可达性

**表 5.177 查询两设备间 P2P 可达性接口说明**

| 接口名称 | 查询两设备间 P2P 可达性 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 查询两设备间 P2P 可达性 |  |  |  |
| URL | /topology/p2p/:srcDvInd/:dstDvInd |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| srcDvInd | integer | path | Y | 源设备索引 |
| dstDvInd | integer | path | Y | 目标设备索引 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 请求成功 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| accessible | boolean |  | 两设备间是否支持 P2P 访问 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"accessible":false} |  |  |  |

#### 5.1.15.9 获取物理链接拓扑

**表 5.194 获取物理链接拓扑接口说明**

| 接口名称 | 获取物理链接拓扑 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 获取物理链接拓扑 |  |  |  |
| URL | /linktopo |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 请求成功 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| data | array |  | 物理链接拓扑信息 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"data":[]} |  |  |  |

### 5.1.16 vDevice

#### 5.1.16.1 创建虚拟设备

**表 5.113 创建虚拟设备接口说明**

| 接口名称 | 创建虚拟设备 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 在指定的物理设备上创建指定数量的虚拟设备，返回创建的虚拟设备ID集合 |  |  |  |
| URL | /CreateVDevices |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 |  |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | query | Y | 物理设备的索引 |
| vDevCount | integer | query | Y | 要创建的虚拟设备数量 |
| vDevCUs | array | query | Y | 每个虚拟设备的计算单元数量，多个值使用多个 vDevCUs 参数传递，例如：vDevCUs=10&vDevCUs=20 |
| vDevMemSize | array | query | Y | 每个虚拟设备的内存大小，多个值使用多个 vDevMemSize 参数传递，例如：vDevMemSize=1024&vDevMemSize=2048 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 虚拟设备创建成功，返回虚拟设备ID集合 |  |  |  |
| 400 | 请求参数无效 |  |  |  |
| 500 | 创建虚拟设备失败 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| vdevIDs | array |  | VDevIDs 创建成功的虚拟设备ID集合 |  |
| 示例 |  |  |  |  |
| 请求参数 | dvInd=0&vDevCount=0&vDevCUs=[{}]&vDevMemSize=[{}] |  |  |  |
| 返回值 | {"vdevIDs":[{}]} |  |  |  |

#### 5.1.16.2 销毁单个虚拟设备

**表 5.114 销毁单个虚拟设备接口说明**

| 接口名称 | 销毁单个虚拟设备 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 销毁指定索引的虚拟设备 |  |  |  |
| URL | /DestroySingleVDevice |  |  |  |
| 请求方式 | delete |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 |  |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| vDvInd | integer | query | Y | 虚拟设备的索引 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 虚拟设备销毁成功 |  |  |  |
| 400 | 虚拟设备销毁失败 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| 示例 |  |  |  |  |
| 请求参数 | vDvInd=0 |  |  |  |
| 返回值 |  |  |  |  |

#### 5.1.16.3 销毁所有虚拟设备

**表 5.115 销毁所有虚拟设备接口说明**

| 接口名称 | 销毁所有虚拟设备 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 销毁指定物理设备上的所有虚拟设备 |  |  |  |
| URL | /DestroyVDevice |  |  |  |
| 请求方式 | delete |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 |  |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | query | Y | 物理设备的索引 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 虚拟设备销毁成功 |  |  |  |
| 400 | 虚拟设备销毁失败 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| 示例 |  |  |  |  |
| 请求参数 | dvInd=0 |  |  |  |
| 返回值 |  |  |  |  |

#### 5.1.16.4 获取当前虚拟机的加密状态

**表 5.116 获取当前虚拟机的加密状态接口说明**

| 接口名称 | 获取当前虚拟机的加密状态 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 返回虚拟机是否处于加密状态 |  |  |  |
| URL | /EncryptionVMStatus |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 |  |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 加密状态 |  |  |  |
| 400 | 操作失败 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| status | boolean |  | Status 是否处于加密状态 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"status":true} |  |  |  |

#### 5.1.16.5 设置虚拟机加密状态

**表 5.117 设置虚拟机加密状态接口说明**

| 接口名称 | 设置虚拟机加密状态 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 根据提供的状态开启或关闭虚拟机加密 |  |  |  |
| URL | /SetEncryptionVMStatus |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 |  |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| status | boolean | query | Y | 加密状态 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 操作成功 |  |  |  |
| 400 | 操作失败 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| 示例 |  |  |  |  |
| 请求参数 | status=true |  |  |  |
| 返回值 |  |  |  |  |

#### 5.1.16.6 启动指定的虚拟设备

**表 5.118 启动指定的虚拟设备接口说明**

| 接口名称 | 启动指定的虚拟设备 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 启动虚拟设备，指定设备索引 |  |  |  |
| URL | /StartVDevice/{vDvInd} |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 |  |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| vDvInd | integer | path | Y | 虚拟设备索引 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 操作成功 |  |  |  |
| 400 | 操作失败 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 |  |  |  |  |

#### 5.1.16.7 停止指定的虚拟设备

**表 5.119 停止指定的虚拟设备接口说明**

| 接口名称 | 停止指定的虚拟设备 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 停止虚拟设备，指定设备索引 |  |  |  |
| URL | /StopVDevice/{vDvInd} |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 |  |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| vDvInd | integer | path | Y | 虚拟设备索引 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 操作成功 |  |  |  |
| 400 | 操作失败 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 |  |  |  |  |

#### 5.1.16.8 更新虚拟设备资源

**表 5.120 更新虚拟设备资源接口说明**

| 接口名称 | 更新虚拟设备资源 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 更新指定虚拟设备的计算单元和内存大小 |  |  |  |
| URL | /UpdateSingleVDevice |  |  |  |
| 请求方式 | put |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 |  |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| vDvInd | integer | query | Y | 虚拟设备的索引 |
| vDevCUs | integer | query | Y | 更新后的计算单元数量 |
| vDevMemSize | integer | query | Y | 更新后的内存大小 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 虚拟设备更新成功 |  |  |  |
| 400 | 虚拟设备更新失败 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| 示例 |  |  |  |  |
| 请求参数 | vDvInd=0&vDevCUs=0&vDevMemSize=0 |  |  |  |
| 返回值 |  |  |  |  |

#### 5.1.16.9 获取单个虚拟设备的信息

**表 5.121 获取单个虚拟设备的信息接口说明**

| 接口名称 | 获取单个虚拟设备的信息 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 根据设备索引获取对应的虚拟设备信息 |  |  |  |
| URL | /VDeviceSingleInfo |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 |  |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| vDvInd | integer | query | Y | 设备索引 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 虚拟设备信息 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 内部服务器错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| containerID | integer |  | ContainerID 虚拟设备的容器ID |  |
| dvInd | integer |  | DvInd 物理设备的设备ID |  |
| name | string |  | Name 虚拟设备的名称 |  |
| pciBusNumber | string |  | PciBusNumber 虚拟设备的总线编号 |  |
| subsystemTypeName | string |  | SubsystemTypeName 设备子系统名称 |  |
| vComputeUnitCount | integer |  | VComputeUnitCount 虚拟设备的计算单元数量 |  |
| vMemoryTotal | integer |  | VMemoryTotal 虚拟设备的全局内存大小 |  |
| vMemoryUsed | integer |  | VMemoryUsed 虚拟设备的已使用内存大小 |  |
| vPercent | integer |  | VPercent 虚拟设备的使用百分比 |  |
| vdvInd | integer |  | VdvInd 虚拟设备的索引号 |  |
| 示例 |  |  |  |  |
| 请求参数 | vDvInd=0 |  |  |  |
| 返回值 | {  "pciBusNumber": "string",  "vComputeUnitCount": 0,  "vPercent": 0,  "vdvInd": 0,  "name": "string",  "dvInd": 0,  "vMemoryUsed": 0,  "containerID": 0,  "vMemoryTotal": 0,  "subsystemTypeName": "string"} |  |  |  |

#### 5.1.16.10 获取设备剩余信息

**表 5.122 获取设备剩余信息接口说明**

| 接口名称 | 获取设备剩余信息 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 获取指定物理设备的剩余计算单元和内存信息 |  |  |  |
| URL | /deviceRemainingInfo/{dvInd} |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 |  |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 物理设备索引 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 设备剩余信息 |  |  |  |
| 400 | 无效的设备索引 |  |  |  |
| 500 | 内部服务器错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| cus | integer |  | CUs 剩余计算单元数量 |  |
| memories | integer |  | Memories 剩余内存大小（单位：字节） |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {  "cus": 0,  "memories": 0} |  |  |  |

#### 5.1.16.11 获取虚拟设备数量

**表 5.123 获取虚拟设备数量接口说明**

| 接口名称 | 获取虚拟设备数量 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 获取当前系统中的虚拟设备数量 |  |  |  |
| URL | /vDeviceCount |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 |  |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 虚拟设备数量 |  |  |  |
| 500 | 内部服务器错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| vDeviceCount | integer |  | VDeviceCount 虚拟设备数量 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"vDeviceCount":0} |  |  |  |

#### 5.1.16.12 获取所有虚拟设备信息

**表 5.124 获取所有虚拟设备信息接口说明**

| 接口名称 | 获取所有虚拟设备信息 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 查询当前系统中所有虚拟设备（VDevice）的详细信息，包括所属物理设备、使用率及 PCI 信息 |  |  |  |
| URL | /vdevice/infos |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | OK |  |  |  |
| 500 | Internal Server Error |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| data | object |  | 接口返回数据 |  |
| vDeviceInfoserror | arrayobject |  | 虚拟设备信息列表 |  |
| namemessage | stringstring |  | 虚拟设备名称 |  |
| subsystemTypeName | string |  | 设备子系统名称 |  |
| vComputeUnitCount | integer (int) |  | 虚拟设备的计算单元数量 |  |
| vMemoryTotal | integer (uint64) |  | 虚拟设备的全局内存大小（字节） |  |
| vMemoryUsed | integer (uint64) |  | 虚拟设备已使用的内存大小（字节） |  |
| containerID | integer (uint64) |  | 虚拟设备所属的容器 ID |  |
| dvInd | integer (int) |  | 所属物理设备 ID |  |
| vPercent | integer (int) |  | 虚拟设备使用率百分比 |  |
| vdvInd | integer (int) |  | 虚拟设备索引号 |  |
| pciBusNumber | string |  | 虚拟设备 PCI 总线编号 |  |
| error | object |  | 错误信息（失败时返回） |  |
| message | string |  | 接口返回描述信息 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {  "data": {    "vDeviceInfos": [      {        "name": "vdevice-0",        "subsystemTypeName": "compute",        "vComputeUnitCount": 16,        "vMemoryTotal": 8589934592,        "vMemoryUsed": 2147483648,        "containerID": 123456,        "dvInd": 0,        "vPercent": 50,        "vdvInd": 0,        "pciBusNumber": "0000:65:00.0"      },      {        "name": "vdevice-1",        "subsystemTypeName": "compute",        "vComputeUnitCount": 8,        "vMemoryTotal": 4294967296,        "vMemoryUsed": 1073741824,        "containerID": 123457,        "dvInd": 1,        "vPercent": 25,        "vdvInd": 1,        "pciBusNumber": "0000:66:00.0"      }    ]  },  "error": {},  "message": "success"} |  |  |  |

#### 5.1.16.13 查询虚拟设备 compute 忙碌占比

**表 5.181 查询虚拟设备 compute 忙碌占比接口说明**

| 接口名称 | 查询虚拟设备 compute 忙碌占比 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 查询虚拟设备 compute 忙碌占比（近似 MIG_ENC_UTIL/DEC_UTIL） |  |  |  |
| URL | /vdevice/busypercent/:vDvInd |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| vDvInd | integer | path | Y | 虚拟设备索引 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 请求成功 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| busyPercent | number |  | 虚拟设备 compute 忙碌占比 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"busyPercent":0} |  |  |  |

### 5.1.17 Version

#### 5.1.17.1 获取RSMI版本

**表 5.125 获取RSMI版本接口说明**

| 接口名称 | 获取RSMI版本 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 返回当前系统的RSMI版本 |  |  |  |
| URL | /RsmiVersion |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回RSMI版本信息 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| rsmiVersion | object |  | RsmiVersion 当前RSMI版本 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"rsmiVersion":{}} |  |  |  |

#### 5.1.17.2 获取设备的VBIOS版本

**表 5.126 获取设备的VBIOS版本接口说明**

| 接口名称 | 获取设备的VBIOS版本 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 返回指定设备的VBIOS版本 |  |  |  |
| URL | /VbiosVersion/{dvInd} |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 设备索引 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回设备VBIOS版本信息 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| vbios | string |  | Vbios VBIOS版本号 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"vbios":"string"} |  |  |  |

#### 5.1.17.3 获取当前系统的驱动程序版本

**表 5.127 获取当前系统的驱动程序版本接口说明**

| 接口名称 | 获取当前系统的驱动程序版本 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 返回指定组件的驱动程序版本 |  |  |  |
| URL | /Version |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| component | string | query | Y | 驱动组件:FIRST、DRIVER、LAST |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回驱动程序版本信息 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| version | string |  | Version 驱动程序版本 |  |
| 示例 |  |  |  |  |
| 请求参数 | component=string |  |  |  |
| 返回值 | {"version":"string"} |  |  |  |

#### 5.1.17.4 显示设备固件版本信息

**表 5.128 显示设备固件版本信息接口说明**

| 接口名称 | 显示设备固件版本信息 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 获取指定设备列表和固件类型列表的固件版本信息 |  |  |  |
| URL | /firmware/info |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvIdList | array | query | Y | 设备 ID 列表 |
| fwType | array | query | Y | 固件类型列表 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 固件版本信息列表 |  |  |  |
| 400 | 请求参数错误或获取固件信息失败 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| fwInfos | array:router.FirmwareInfoResp |  | FwInfos 设备固件信息列表 |  |
| deviceID | integer |  | DeviceID 设备索引号 |  |
| firmwareVer | array:router.FirmwareBlock |  | FirmwareVer 固件块列表 |  |
| blockName | string |  | BlockName 固件块名称 |  |
| version | string |  | Version 固件版本 |  |
| 示例 |  |  |  |  |
| 请求参数 | dvIdList=[{}]&fwType=[{}] |  |  |  |
| 返回值 | {  "fwInfos": [    {      "deviceID": 0,      "firmwareVer": [        {          "blockName": "string",          "version": "string"        }      ]    }  ]} |  |  |  |

#### 5.1.17.5 显示设备的VBIOS版本

**表 5.129 显示设备的VBIOS版本接口说明**

| 接口名称 | 显示设备的VBIOS版本 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 获取并显示指定设备的VBIOS版本信息。 |  |  |  |
| URL | /showVbiosVersion |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 |  |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvIdList | array | body | Y | 设备ID列表 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 设备的VBIOS版本信息列表 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| deviceId | integer |  | DeviceID 设备索引号 example: 0 |  |
| vbios | string |  | VBIOS 版本信息 example: "113-C56801-102" |  |
| 示例 |  |  |  |  |
| 请求参数 | -d '[{}]' |  |  |  |
| 返回值 | {  "deviceId": 0,  "vbios": "string"} |  |  |  |

#### 5.1.17.6 查询海光扩展版本

**表 5.178 查询海光扩展版本接口说明**

| 接口名称 | 查询海光扩展版本 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 查询海光扩展版本 |  |  |  |
| URL | /version/hy |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 请求成功 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| hyMajor | integer |  | 海光扩展主版本 |  |
| hyMinor | integer |  | 海光扩展次版本 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"hyMajor":0,"hyMinor":0} |  |  |  |

#### 5.1.17.7 获取 HSW 固件版本

**表 5.186 获取 HSW 固件版本接口说明**

| 接口名称 | 获取 HSW 固件版本 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 获取 HSW 固件版本 |  |  |  |
| URL | /hswfw/:dvInd |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | 物理设备索引 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 请求成功 |  |  |  |
| 400 | 请求参数错误 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| version | string |  | HSW 固件版本 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"version":"string"} |  |  |  |

### 5.1.18 Group

#### 5.1.18.1 获取所有设备组信息

**表 5.130 获取所有设备组信息接口说明**

| 接口名称 | 获取所有设备组信息 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 返回节点所有设备组信息 |  |  |  |
| URL | /groups |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回设备组信息 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| group_id | integer |  | 设备组ID |  |
| group_name | string |  | 设备组名称 |  |
| entity_group_id | integer |  | 设备类型ID |  |
| entity_id | integer |  | 设备ID |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"message":"成功","data":{"groups":[{"group_id":0,"group_name":"group0","entity_list":[{"entity_group_id":1,"entity_id":0},{"entity_group_id":1,"entity_id":1}]},{"group_id":1,"group_name":"group1","entity_list":[{"entity_group_id":1,"entity_id":1}]}]}} |  |  |  |

#### 5.1.18.2 获取指定设备组信息

**表 5.131 获取指定设备组信息接口说明**

| 接口名称 | 获取指定设备组信息 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 返回指定设备组信息 |  |  |  |
| URL | /groups/{groupId} |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| groupId | integer | path | Y | 设备组ID |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回设备组信息 |  |  |  |
| 400 | 请求失败 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| group_id | integer |  | 设备组ID |  |
| group_name | string |  | 设备组名称 |  |
| entity_group_id | integer |  | 设备类型ID |  |
| entity_id | integer |  | 设备ID |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"message":"成功","data":{"groupInfo":{"group_id":0,"group_name":"group0","entity_list":[{"entity_group_id":1,"entity_id":0},{"entity_group_id":1,"entity_id":1}]}}} |  |  |  |

#### 5.1.18.3 创建设备组

**表 5.132 创建设备组接口说明**

| 接口名称 | 创建设备组 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 给定名称，创建新的设备组 |  |  |  |
| URL | /groups |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| groupName | string | body | Y | 设备组名称 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回创建成功的设备组ID |  |  |  |
| 400 | 请求失败 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| groupId | integer |  | 设备组ID |  |
| 示例 |  |  |  |  |
| 请求参数 | -d '{"groupName": "group0"}' |  |  |  |
| 返回值 | {"message":"成功","data":{"groupId":0}} |  |  |  |

#### 5.1.18.4 向指定设备组中添加HCU设备

**表 5.133 向指定设备组中添加HCU设备接口说明**

| 接口名称 | 向指定设备组中添加HCU设备 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 给定设备组ID，向其中添加HCU设备 |  |  |  |
| URL | /groups/{groupId}/hcus |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| groupId | integer | path | Y | 设备组ID |
| dvInd | integer | body | Y | HCU设备ID |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回添加成功信息 |  |  |  |
| 400 | 请求失败 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| 示例 |  |  |  |  |
| 请求参数 | -d '{"dvInd": 0}' |  |  |  |
| 返回值 | {"message":"成功","data":null} |  |  |  |

#### 5.1.18.5 向指定设备组中添加设备

**表 5.134 向指定设备组中添加设备接口说明**

| 接口名称 | 向指定设备组中添加设备 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 给定设备组ID，向其中添加设备 |  |  |  |
| URL | /groups/{groupId}/entities |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| groupId | integer | path | Y | 设备组ID |
| entity_group_id | integer | body | Y | 设备类型ID |
| entity_id | integer | body | Y | 设备ID |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回添加成功信息 |  |  |  |
| 400 | 请求失败 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| 示例 |  |  |  |  |
| 请求参数 | -d '{"entityList": [{"entity_group_id": 1, "entity_id": 0}, {"entity_group_id": 1, "entity_id": 1}]}' |  |  |  |
| 返回值 | {"message":"成功","data":null} |  |  |  |

#### 5.1.18.6 从设备组中移除设备

**表 5.135 从设备组中移除设备接口说明**

| 接口名称 | 从设备组中移除设备 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 给定设备组ID，从中移除设备 |  |  |  |
| URL | /groups/{groupId}/entities |  |  |  |
| 请求方式 | delete |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| groupId | integer | path | Y | 设备组ID |
| entity_group_id | integer | body | Y | 设备类型ID |
| entity_id | integer | body | Y | 设备ID |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回移除成功信息 |  |  |  |
| 400 | 请求失败 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| 示例 |  |  |  |  |
| 请求参数 | -d '{"entityList": [{"entity_group_id": 1, "entity_id": 0}, {"entity_group_id": 1, "entity_id": 1}]}' |  |  |  |
| 返回值 | {"message":"成功","data":null} |  |  |  |

#### 5.1.18.7 删除指定设备组

**表 5.136 删除指定设备组接口说明**

| 接口名称 | 删除指定设备组 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 删除给定设备组ID对应的设备组 |  |  |  |
| URL | /groups/{groupId} |  |  |  |
| 请求方式 | delete |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| groupId | integer | path | Y | 设备组ID |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回删除成功信息 |  |  |  |
| 400 | 请求失败 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"message":"成功","data":null} |  |  |  |

### 5.1.19 Field Group

#### 5.1.19.1 获取所有字段组信息

**表 5.137 获取所有字段组信息接口说明**

| 接口名称 | 获取所有字段组信息 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 返回节点所有字段组信息 |  |  |  |
| URL | /fieldgroups |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回字段组信息 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| field_group_id | integer |  | 字段组ID |  |
| fieldgroup_name | string |  | 字段组名称 |  |
| field_ids | []int |  | 字段ID列表 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"message":"成功","data":{"fieldGroups":[{"field_group_id":0,"field_group_name":"fg0","field_ids":[101,112,102]},{"field_group_id":1,"field_group_name":"fg1","field_ids":[103,104]}]}} |  |  |  |

#### 5.1.19.2 获取指定字段组信息

**表 5.138 获取指定字段组信息接口说明**

| 接口名称 | 获取指定字段组信息 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 返回指定字段组信息 |  |  |  |
| URL | /fieldgroups/{fieldGroupId} |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回字段组信息 |  |  |  |
| 400 | 请求失败 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| field_group_id | integer |  | 字段组ID |  |
| fieldgroup_name | string |  | 字段组名称 |  |
| field_ids | []int |  | 字段ID列表 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"message":"成功","data":{"fieldGroupInfo":{"field_group_id":0,"field_group_name":"fg0","field_ids":[101,112,102]}}} |  |  |  |

#### 5.1.19.3 创建字段组

**表 5.139 创建字段组接口说明**

| 接口名称 | 创建字段组 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 给定字段组名称以及字段列表，创建新的字段组 |  |  |  |
| URL | /fieldgroups |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| fieldGroupName | string | body | Y | 字段组名称 |
| fieldIds | []int | body | Y | 字段列表 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回创建成功的字段组ID |  |  |  |
| 400 | 请求失败 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| fieldGroupId | integer |  | 字段组ID |  |
| 示例 |  |  |  |  |
| 请求参数 | -d '{"fieldGroupName": "fg0", "fieldIds": [101,102]}' |  |  |  |
| 返回值 | {"message":"成功","data":{"fieldGroupId":0}} |  |  |  |

#### 5.1.19.4 删除指定字段组

**表 5.140 删除指定字段组接口说明**

| 接口名称 | 删除指定字段组 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 删除给定字段组ID对应的字段组 |  |  |  |
| URL | /fieldgroups/{fieldgroupId} |  |  |  |
| 请求方式 | delete |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| fieldGroupId | integer | path | Y | 字段组ID |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回删除成功信息 |  |  |  |
| 400 | 请求失败 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"message":"成功","data":null} |  |  |  |

### 5.1.20 Policy

#### 5.1.20.1 获取所有策略

**表 5.141 获取所有策略接口说明**

| 接口名称 | 获取所有策略 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 返回节点所有策略 |  |  |  |
| URL | /policies |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回策略信息 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| hcu_index | integer |  | HCU设备序列 |  |
| action_index | integer |  | 执行动作序列 |  |
| validation_index | integer |  | 执行动作后验证序列 |  |
| conditions | PolicyConditions |  | 策略执行条件集合 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"message":"成功","data":{"policyList":[{"hcu_index":0,"action_index":0,"validation_index":0,"conditions":{"max_pages_enable":true,"max_pages":100,"max_temp_enable":false,"max_temp":0,"max_power_enable":false,"max_power":0,"ecc_errors_enable":true,"pcie_errors_enable":false}},{"hcu_index":1,"action_index":0,"validation_index":0,"conditions":{"max_pages_enable":false,"max_pages":0,"max_temp_enable":false,"max_temp":0,"max_power_enable":false,"max_power":0,"ecc_errors_enable":false,"pcie_errors_enable":false}}]}} |  |  |  |

#### 5.1.20.2 获取指定设备的策略

**表 5.142 获取指定设备的策略接口说明**

| 接口名称 | 获取指定设备的策略 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 给定HCU设备序列，返回其已设定策略 |  |  |  |
| URL | /policies/hcus/{dvInd} |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | HCU设备序列 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回策略信息 |  |  |  |
| 400 | 请求失败 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| hcu_index | integer |  | HCU设备序列 |  |
| action_index | integer |  | 执行动作序列 |  |
| validation_index | integer |  | 执行动作后验证序列 |  |
| conditions | PolicyConditions |  | 策略执行条件集合 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"message":"成功","data":{"policyList":[{"hcu_index":0,"action_index":0,"validation_index":0,"conditions":{"max_pages_enable":true,"max_pages":100,"max_temp_enable":false,"max_temp":0,"max_power_enable":false,"max_power":0,"ecc_errors_enable":true,"pcie_errors_enable":false}}]}} |  |  |  |

#### 5.1.20.3 设置策略

**表 5.143 设置策略接口说明**

| 接口名称 | 设置策略 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 给定HCU设备序列，设置策略 |  |  |  |
| URL | /policies/hcus/{dvInd} |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | HCU设备序列 |
| hcu_index | integer | body | Y | HCU设备序列 |
| action_index | integer | body | Y | 执行动作序列 |
| validation_index | integer | body | N | 执行动作后验证序列 |
| conditions | PolicyConditions | body | Y | 策略执行条件集合 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回设置成功信息 |  |  |  |
| 400 | 请求失败 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| hcu_index | integer |  | HCU设备序列 |  |
| action_index | integer |  | 执行动作序列 |  |
| validation_index | integer |  | 执行动作后验证序列 |  |
| conditions | PolicyConditions |  | 策略执行条件集合 |  |
| 示例 |  |  |  |  |
| 请求参数 | -d '{"hcu_index": 0, "action_index": 0, "validation_index": 0, "conditions": {"max_pages_enable": true, "max_pages": 100, "ecc_errors_enable": true}}' |  |  |  |
| 返回值 | {"message":"成功","data":{"policyList":[{"hcu_index":0,"action_index":0,"validation_index":0,"conditions":{"max_pages_enable":true,"max_pages":100,"max_temp_enable":false,"max_temp":0,"max_power_enable":false,"max_power":0,"ecc_errors_enable":true,"pcie_errors_enable":false}}]}} |  |  |  |

#### 5.1.20.4 清除策略

**表 5.144 清除策略接口说明**

| 接口名称 | 清除策略 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 给定HCU设备序列，清除其上设定的策略 |  |  |  |
| URL | /policies/hcus/{dvInd} |  |  |  |
| 请求方式 | delete |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | HCU设备序列 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回清除成功信息 |  |  |  |
| 400 | 请求失败 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"message":"成功","data":null} |  |  |  |

#### 5.1.20.5 判定策略条件

**表 5.145 判定策略条件接口说明**

| 接口名称 | 判定策略条件 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 给定HCU设备序列列表，获取设定策略的监测结果 |  |  |  |
| URL | /policies/result |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| hcuList | []int | body | Y | HCU设备序列列表 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回策略信息 |  |  |  |
| 400 | 请求失败 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| error | string |  | 错误信息 |  |
| errorHcuIndex | integer |  | 到达条件的设备序列 |  |
| 示例 |  |  |  |  |
| 请求参数 | -d '{"hcuList": [0,1]}' |  |  |  |
| 返回值 | {"message":"成功","data":{"error":null,"errorhcuIndex":null}} |  |  |  |

#### 5.1.20.6 执行动作

**表 5.146 执行动作接口说明**

| 接口名称 | 执行动作 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 给定HCU设备序列，执行动作 |  |  |  |
| URL | /policies/action/{dvInd} |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | HCU设备序列 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回设置成功信息 |  |  |  |
| 400 | 请求失败 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"message":"成功","data":null} |  |  |  |

### 5.1.21 Field

#### 5.1.21.1 获取所有字段信息

**表 5.147 获取所有字段信息接口说明**

| 接口名称 | 获取所有字段信息 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 返回所有字段信息 |  |  |  |
| URL | /fields |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回字段信息 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| FieldId | integer |  | 字段ID |  |
| Name | string |  | 字段名称 |  |
| EntityLevel | integer |  | 字段对应设备类型ID |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"message":"成功","data":{"fieldMetaList":[{"FieldId":101,"Name":"HCU_TEMP","EntityLevel":1},{"FieldId":102,"Name":"HCU_POWER_USAGE","EntityLevel":1},{"FieldId":103,"Name":"HCU_POWER_CAP","EntityLevel":1},{"FieldId":104,"Name":"HCU_UTILIZATION_RATE","EntityLevel":1},{"FieldId":105,"Name":"HCU_SCLK","EntityLevel":1},{"FieldId":106,"Name":"HCU_COMPUTE_UNIT_COUNT","EntityLevel":1},{"FieldId":107,"Name":"HCU_COMPUTE_UNIT_REMAINING_COUNT","EntityLevel":1},{"FieldId":108,"Name":"HCU_USED_MEMORY_BYTES","EntityLevel":1},{"FieldId":109,"Name":"HCU_MEMORY_CAP_BYTES","EntityLevel":1},{"FieldId":110,"Name":"HCU_MEMORY_REMAINING","EntityLevel":1},{"FieldId":111,"Name":"HCU_MEMORY_PERCENT","EntityLevel":1},{"FieldId":112,"Name":"HCU_PCIE_BW_MB","EntityLevel":1},{"FieldId":113,"Name":"HCU_PCIE_RECEIVE_MB","EntityLevel":1},{"FieldId":114,"Name":"HCU_PCIE_SENT_MB","EntityLevel":1},{"FieldId":115,"Name":"HCU_DF_BW_READ","EntityLevel":1},{"FieldId":116,"Name":"HCU_DF_BW_WRITE","EntityLevel":1},{"FieldId":117,"Name":"HCU_DF_BW_READ_WRITE","EntityLevel":1},{"FieldId":118,"Name":"HCU_VHCU_COUNT","EntityLevel":1},{"FieldId":201,"Name":"VHCU_SCLK","EntityLevel":2},{"FieldId":202,"Name":"VHCU_TEMP","EntityLevel":2},{"FieldId":203,"Name":"VHCU_UTILIZATION_RATE","EntityLevel":2},{"FieldId":204,"Name":"VHCU_USED_MEMORY_BYTES","EntityLevel":2}]}} |  |  |  |

#### 5.1.21.2 获取指定字段信息

**表 5.148 获取指定字段信息接口说明**

| 接口名称 | 获取指定字段信息 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 返回指定字段信息 |  |  |  |
| URL | /fields/{fieldId} |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回字段信息 |  |  |  |
| 400 | 请求失败 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| FieldId | integer |  | 字段ID |  |
| Name | string |  | 字段名称 |  |
| EntityLevel | integer |  | 字段对应设备类型ID |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"message":"成功","data":{"fieldMeta":{"FieldId":101,"Name":"HCU_TEMP","EntityLevel":1}}} |  |  |  |

#### 5.1.21.3 监测指定HCU设备的字段值

**表 5.149 监测指定HCU设备的字段值接口说明**

| 接口名称 | 监测指定HCU设备的字段值 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 给定HCU设备序列，监测其字段值 |  |  |  |
| URL | /fields/hcus/{dvInd} |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | HCU设备序列 |
| fieldIdList | []int | body | Y | 字段ID列表 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回监测成功信息 |  |  |  |
| 400 | 请求失败 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| 示例 |  |  |  |  |
| 请求参数 | -d '{"fieldIdList": [101,102]}' |  |  |  |
| 返回值 | {"message":"成功","data":null} |  |  |  |

#### 5.1.21.4 获取指定HCU设备监测项的最新值

**表 5.150 获取指定HCU设备监测项的最新值接口说明**

| 接口名称 | 获取指定HCU设备监测项的最新值 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 指定HCU设备序列，返回其监测项的最新值 |  |  |  |
| URL | /fields/hcus/{dvInd}/values |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | HCU设备序列 |
| fieldIdList | []int | body | Y | 字段ID列表 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回字段值信息 |  |  |  |
| 400 | 请求失败 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| FieldId | integer |  | 字段ID |  |
| Timestamp | int64 |  | 时间戳 |  |
| Value | float |  | 字段值 |  |
| Err | error |  | 报错信息 |  |
| Tag | string |  | 额外注解信息 |  |
| 示例 |  |  |  |  |
| 请求参数 | -d '{"fieldIdList": [101,102]}' |  |  |  |
| 返回值 | {"message":"成功","data":{"fieldValueList":[{"FieldId":101,"Timestamp":1767143904286,"Value":54,"Err":null,"Tag":""},{"FieldId":102,"Timestamp":1767143904286,"Value":84,"Err":null,"Tag":""}]}} |  |  |  |

#### 5.1.21.5 取消监测指定HCU设备

**表 5.151 取消监测指定HCU设备接口说明**

| 接口名称 | 取消监测指定HCU设备 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 给定HCU设备序列，取消监测该设备 |  |  |  |
| URL | /fields/hcus/{dvInd} |  |  |  |
| 请求方式 | delete |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | HCU设备序列 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回取消监测成功信息 |  |  |  |
| 400 | 请求失败 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"message":"成功","data":null} |  |  |  |

#### 5.1.21.6 监测指定设备组的字段值

**表 5.152 监测指定设备组的字段值接口说明**

| 接口名称 | 监测指定设备组的字段值 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 给定设备组ID，监测其内设备字段值 |  |  |  |
| URL | /fields/groups/{groupId} |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| groupId | integer | path | Y | 设备组ID |
| fieldIdList | []int | body | Y | 字段ID列表 |
| updateFreq | float | body | N | 接口内部每次获取字段值的时间间隔（单位为秒） |
| maxKeepAge | float | body | N | 接口内部保留字段值的时长（单位为秒） |
| maxKeepSamples | integer | body | N | 接口内部保留字段值的样本数 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回监测成功信息 |  |  |  |
| 400 | 请求失败 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| 示例 |  |  |  |  |
| 请求参数 | -d '{"fieldIdList": [101,102], "updateFreq": 3, "maxKeepAge": 0, "maxKeepSamples": 3}' |  |  |  |
| 返回值 | {"message":"成功","data":null} |  |  |  |

#### 5.1.21.7 获取指定设备组监测项的最新值

**表 5.153 获取指定设备组监测项的最新值接口说明**

| 接口名称 | 获取指定设备组监测项的最新值 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 指定设备组ID，返回其设备监测项的最新值 |  |  |  |
| URL | /fields/groups/{groupId}/values |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| groupId | integer | path | Y | 设备组ID |
| fieldIdList | []int | body | Y | 字段ID列表 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回字段值信息 |  |  |  |
| 400 | 请求失败 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| entityGroupId | integer |  | 设备类型ID |  |
| entityId | integer |  | 设备ID |  |
| FieldId | integer |  | 字段ID |  |
| Timestamp | int64 |  | 时间戳 |  |
| Value | float |  | 字段值 |  |
| Err | error |  | 报错信息 |  |
| Tag | string |  | 额外注解信息 |  |
| 示例 |  |  |  |  |
| 请求参数 | -d '{"fieldIdList": [101,102]}' |  |  |  |
| 返回值 | { "message": "成功", "data": {"entityList": [{ "entityGroupId": "1", "entityId": "0", "fieldValueList": [{"FieldId":101,"Timestamp":1767143904286,"Value":54,"Err":null,"Tag":""},{"FieldId":102,"Timestamp":1767143904286,"Value":84,"Err":null,"Tag":""}]},{"entityGroupId": "1", "entityId": "1", "fieldValueList": [ {"FieldId":101,"Timestamp":1767143904286,"Value":55,"Err":null,"Tag":""}, {"FieldId":102,"Timestamp":1767143904286,"Value":88,"Err":null,"Tag":""} ]}]}} |  |  |  |

#### 5.1.21.8 取消监测指定设备组

**表 5.154 取消监测指定设备组接口说明**

| 接口名称 | 取消监测指定设备组 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 给定设备组，取消监测该组内设备 |  |  |  |
| URL | /fields/groups/{groupId} |  |  |  |
| 请求方式 | delete |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| groupId | integer | path | Y | 设备组ID |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回取消监测成功信息 |  |  |  |
| 400 | 请求失败 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"message":"成功","data":null} |  |  |  |

#### 5.1.21.9 监测指定设备的字段值

**表 5.155 监测指定设备的字段值接口说明**

| 接口名称 | 监测指定设备的字段值 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 给定设备类型ID和设备ID，监测其字段值 |  |  |  |
| URL | /fields/entities/{entityGroupId}/{entityId} |  |  |  |
| 请求方式 | post |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| entityGroupId | integer | path | Y | 设备类型ID |
| entityId | integer | path | Y | 设备ID |
| fieldIdList | []int | body | Y | 字段ID列表 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回监测成功信息 |  |  |  |
| 400 | 请求失败 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| 示例 |  |  |  |  |
| 请求参数 | -d '{"fieldIdList": [201,202]}' |  |  |  |
| 返回值 | {"message":"成功","data":null} |  |  |  |

#### 5.1.21.10 获取指定设备监测项的最新值

**表 5.156 获取指定设备监测项的最新值接口说明**

| 接口名称 | 获取指定设备监测项的最新值 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 指定设备类型ID和设备ID，返回其监测项的最新值 |  |  |  |
| URL | /fields/entities/{entityGroupId}/{entityId}/values |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 | application/json |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| entityGroupId | integer | path | Y | 设备类型ID |
| entityId | integer | path | Y | 设备ID |
| fieldIdList | []int | body | Y | 字段ID列表 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回字段值信息 |  |  |  |
| 400 | 请求失败 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| FieldId | integer |  | 字段ID |  |
| Timestamp | int64 |  | 时间戳 |  |
| Value | float |  | 字段值 |  |
| Err | error |  | 报错信息 |  |
| Tag | string |  | 额外注解信息 |  |
| 示例 |  |  |  |  |
| 请求参数 | -d '{"fieldIdList": [201,202]}' |  |  |  |
| 返回值 | {"message":"成功","data":{"fieldValueList":[{"FieldId":201,"Timestamp":1768634751369,"Value":600,"Err":null,"Tag":""},{"FieldId":202,"Timestamp":1768634751369,"Value":37,"Err":null,"Tag":""}]}} |  |  |  |

#### 5.1.21.11 取消监测指定设备

**表 5.157 取消监测指定设备接口说明**

| 接口名称 | 取消监测指定设备 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 给定设备类型ID和设备ID，取消监测该设备 |  |  |  |
| URL | /fields/entities/{entityGroupId}/{entityId} |  |  |  |
| 请求方式 | delete |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| entityGroupId | integer | path | Y | 设备类型ID |
| entityId | integer | path | Y | 设备ID |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回取消监测成功信息 |  |  |  |
| 400 | 请求失败 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"message":"成功","data":null} |  |  |  |

### 5.1.22 Metric Group

#### 5.1.22.1 获取性能指标项

**表 5.158 获取性能指标项接口说明**

| 接口名称 | 获取性能指标项 |  |  |  |
| --- | --- | --- | --- | --- |
| 接口描述 | 指定HCU设备序列，返回性能指标项 |  |  |  |
| URL | /metricgroups/hcus/{dvInd} |  |  |  |
| 请求方式 | get |  |  |  |
| 请求类型 |  |  |  |  |
| 返回类型 | application/json |  |  |  |
| 参数名 | 数据类型 | 参数类型 | 是否必填 | 说明 |
| dvInd | integer | path | Y | HCU设备序列 |
| 状态码 | 描述 |  | 说明 |  |
| 200 | 返回指标项信息 |  |  |  |
| 400 | 请求失败 |  |  |  |
| 500 | 服务器内部错误 |  |  |  |
| 返回属性名 | 类型 |  | 说明 |  |
| Major | integer |  | 性能指标大类序列 |  |
| Minor | integer |  | 性能指标小类序列 |  |
| FieldIds | []int |  | 字段ID列表 |  |
| 示例 |  |  |  |  |
| 请求参数 |  |  |  |  |
| 返回值 | {"message":"成功","data":{"metricGroupList":[{"Major":1,"Minor":1,"FieldIds":[104,105]},{"Major":2,"Minor":1,"FieldIds":[111]},{"Major":3,"Minor":1,"FieldIds":[112,114,113]}]}} |  |  |  |
