# 2 HCU-DCGM部署

## 2.1 基础软硬件环境

### 2.1.1 硬件环境

需要宿主机安装配置了HCU设备。

### 2.1.2 软件环境

需要在宿主机安装HCU驱动程序，可以在光合开发者社区下载： 光合开发者社区

### 2.1.3 组件列表

**表 2.1 HCU-DCGM 组件列表**

| hcu-dcgm-*.run |
| --- |
| hcu-dcgm-*.tar |

### 2.1.4 使用约束

#### 2.1.4.1 HCU-DCGM服务形式约束

驱动需要6.2.26版本及以上，6.2.26版本以下驱动部分功能不正常，具体如下：

•DF宽带信息无法获取；

•HCU物理卡上切分vHCU虚拟卡的剩余资源信息获取异常；

•虚拟化模块功能不支持多线程，vHCU设备信息无法实时获取。

#### 2.1.4.2 HCU-DCGM命令形式（dcgmi）约束

•驱动需要6.3.8版本及以上；

•DTK版本需要25.0.4.1版本及以上；

•支持系统发行版：rocky8.6、rocky9.2、centos8和ubuntu22.04。

## 2.2 HCU-DCGM服务部署

### 2.2.1 物理机部署

```bash
[root@worker-200 dcgm-dcu]# ./hcu-dcgm-v3.0.0.run
Extracting payload (tail -n +121) to /tmp/hcu-dcgm.tar.gz ...
Checking for running hcu-dcgm processes...
Removing old installation at /opt/hcu-dcgm ...
Extracting payload to /opt ...
HCU DCGM has been installed to /opt/hcu-dcgm.
dcgmi (if present) installed to /usr/local/bin/dcgmi.
To use dcgm libraries, you may source the wrapper:
  source /opt/hcu-dcgm/dcgm
Starting hcu-dcgm with environment sources...
hcu-dcgm started, logging to /opt/hcu-dcgm/hcu-dcgm.log.
```

 2.1 物理机安装

```bash
[root@worker-200 dcgm-dcu]# dcgmi -h
Command-line interface for managing and interacting with DCGM. Use dcgm-cli [command] --help for more information on a command.

Usage:
  dcgmi [command]

Available Commands:
  all-device-infos       Get information for all physical devices
  completion             Generate the autocompletion script for the specified shell
  destroy-vdevice<dvInd> Destroy a single virtual device
  diag                   Run diagnostics commands
  discovery              Used to discover and identify HCUs and their attributes.
  dmon                   Used to monitor HCUs and their stats.
  fieldgroup             Used to create and maintain groups of field IDs.
  group                  Used to create and maintain groups of HCUs.
  health                 Run HCU health check
  help                   Help about any command
  idmon                  Periodically sample selected HCU metrics.
  policy                 Used to control policies for groups of HCUs.
  process                Manage and inspect HCU processes
  profile                View available profiling metrics for HCUs
  topo                   HCU Topology
  vdevice-info           Get virtual device information

Flags:
  -h, --help   help for dcgmi

Use "dcgmi [command] --help" for more information about a command.
```

2.2 帮助命令

使用curl命令来查看指标：

```bash
curl localhost:16081/CollectDeviceMetrics
```

接口调用成功可看到宿主机上HCU卡的部分信息：

```json
[
  {
    "MinorNumber": 0,
    "PciBusNumber": "0000:00:08.0",
    "DeviceId": "TPXS320016070601",
    "SubSystemName": "K100_AI",
    "Temperature": 51,
    "PowerUsage": 84,
    "PowerCap": 400,
    "MemoryCap": 68702699520,
    "MemoryUsed": 2150400,
    "UtilizationRate": 0,
    "PcieBwMb": 0,
    "Clk": 600,
    "SclkFrequency": [
      "600",
      "700",
      "750",
      "800",
      "900",
      "1000",
      "1106",
      "1200",
      "1270",
      "1319",
      "1400"
    ],
    "Socclk": 309,
    "SocclkFrequency": [
      "309",
      "523",
      "566",
      "618",
      "680",
      "755",
      "850",
      "971"
    ],
    "PerfLevel": "AUTO"
  },
  {
    "MinorNumber": 1,
    "PciBusNumber": "0000:00:09.0",
    "DeviceId": "TPXS320010040501",
    "SubSystemName": "K100_AI",
    "Temperature": 55,
    "PowerUsage": 80,
    "PowerCap": 400,
    "MemoryCap": 68702699520,
    "MemoryUsed": 2150400,
    "UtilizationRate": 0,
    "PcieBwMb": 0,
    "Clk": 600,
    "SclkFrequency": [
      "600",
      "700",
      "750",
      "800",
      "900",
      "1000",
      "1106",
      "1200",
      "1270",
      "1319",
      "1400"
    ],
    "Socclk": 309,
    "SocclkFrequency": [
      "309",
      "523",
      "566",
      "618",
      "680",
      "755",
      "850",
      "971"
    ],
    "PerfLevel": "AUTO"
  }
]
```

注：

```bash
nohup /opt/hcu-dcgm/bin/hcu-dcgm \
  -port=16081 \
  -logtostderr \
  -v=2 \
  > /opt/hcu-dcgm/bin/hcu-dcgm.log 2>&1 &
```

1）二进制可执行文件同级目录下新建config文件夹。例如可执行文件hcu-dcgm存放路径为/home/dcgm，则在/home/dcgm路径下新建文件夹名字为config。

2）在/home/dcgm/config路径下创建文件config.json，文件内容如下：

```json
{
  "health": {
    "minPowerThreshold": 10,
    "maxPowerThreshold": 400,
    "maxMemoryUsageThreshold": 95,
    "maxTemperatureThreshold": 95.0
  },
  "gemm": {
    "targetStressTimeSeconds": 1,
    "dgemm": 20.6,
    "sgemm": 40.6,
    "hgemm": 200.0,
    "int8gemm": 330.2,
    "bf16gemm": 200.12,
    "tf32gemm": 110.1
  },
  "bandwidth": {
    "precision": 1,
    "arraysize": 72001536,
    "numtimes": 300
  }
}
```

参数说明：minPowerThreshold和maxPowerThreshold分别是最小功率和最大功率阈值（范围在0-400）；maxMemoryUsageThreshold为最大内存使用量阈值（范围在0-100）；MaxTemperatureThreshold为最大温度阈值（范围在0-100）。

### 2.2.2 docker容器部署

HCU-DCGM服务支持docker容器形式部署，启动方式分为两种形式：

```bash
docker run -d \
  --name hcu-dcgm \
  --privileged \
  --device=/dev/kfd \
  --device=/dev/mkfd \
  --device=/dev/dri \
  -v /opt/hyhal:/opt/hyhal \
  -p 16081:16081 \
  image.sourcefind.cn:5000/hcu/admin/base/hcu-dcgm:v2.1.0 \
  /bin/bash -c "/usr/local/bin/start-dcgm.sh"
```

1）容器启动后自动执行 /usr/local/bin/start-dcgm.sh；

2）DCGM 服务随容器启动。

1）先启动容器

```bash
docker run -d \
  --name hcu-dcgm \
  --privileged \
  --device=/dev/kfd \
  --device=/dev/mkfd \
  --device=/dev/dri \
  -v /opt/hyhal:/opt/hyhal \
  -p 16081:16081 \
  harbor.sourcefind.cn:5443/hcu/admin/base/hcu-dcgm:v3.0.0 \
  tail -f /dev/null
```

2）进入容器

```bash
docker exec -it hcu-dcgm /bin/bash
```

3）在容器启动 DCGM

```bash
/usr/local/bin/start-dcgm.sh
```

```bash
docker exec hcu-dcgm tail -f /home/dcgm/dcgm.log
```

```bash
[root@worker-200 dcgm-dcu]# docker logs -f hcu-dcgm
[INFO] First startup: extracting DTK ...
[INFO] DTK extracted to /opt/dtk-26.04
[INFO] Removed /opt/dtk-26.04.tar.gz
[INFO] sourcing /opt/hyhal/env.sh
[INFO] sourcing /opt/dtk/env.sh
[INFO] Starting hcu-dcgm ...
```

2.3 容器内 DCGM 启动日志

```bash
[GIN-debug] GET    /fields/groups/:groupId/values --> github.com/HYGON-AI/hcu-dcgm/v3/pkg/service/router.GetLatestValuesForFieldsWithGroup (3 handlers)
[GIN-debug] GET    /metricgroups/hcus/:dvInd --> github.com/HYGON-AI/hcu-dcgm/v3/pkg/service/router.GetSupportedMetricGroups (3 handlers)
[GIN-debug] GET    /device/cuutil/:dvInd     --> github.com/HYGON-AI/hcu-dcgm/v3/pkg/service/router.GetDevCuUtil (3 handlers)
[GIN-debug] GET    /device/waveutil/:dvInd   --> github.com/HYGON-AI/hcu-dcgm/v3/pkg/service/router.GetDevWaveUtil (3 handlers)
[GIN-debug] GET    /vdevice/busypercent/:vDvInd --> github.com/HYGON-AI/hcu-dcgm/v3/pkg/service/router.GetVDevBusyPercent (3 handlers)
[GIN-debug] POST   /powercap/:dvInd          --> github.com/HYGON-AI/hcu-dcgm/v3/pkg/service/router.SetPowerCap (3 handlers)
[GIN-debug] POST   /powercap/:dvInd/reset    --> github.com/HYGON-AI/hcu-dcgm/v3/pkg/service/router.ResetPowerCap (3 handlers)
[GIN-debug] POST   /perflevel/v1/:dvInd      --> github.com/HYGON-AI/hcu-dcgm/v3/pkg/service/router.SetPerfLevelV1 (3 handlers)
[GIN-debug] GET    /hswfw/:dvInd             --> github.com/HYGON-AI/hcu-dcgm/v3/pkg/service/router.GetHswFwVersion (3 handlers)
[GIN-debug] GET    /boost/:dvInd             --> github.com/HYGON-AI/hcu-dcgm/v3/pkg/service/router.GetGfxBoostInfo (3 handlers)
[GIN-debug] GET    /fdrhistory/:dvInd        --> github.com/HYGON-AI/hcu-dcgm/v3/pkg/service/router.GetFdrHistory (3 handlers)
[GIN-debug] GET    /memeccinfo/:dvInd        --> github.com/HYGON-AI/hcu-dcgm/v3/pkg/service/router.GetMemBankEccInfo (3 handlers)
[GIN-debug] POST   /boost/:dvInd             --> github.com/HYGON-AI/hcu-dcgm/v3/pkg/service/router.SetGfxBoostLevel (3 handlers)
[GIN-debug] POST   /lowpower/:dvInd          --> github.com/HYGON-AI/hcu-dcgm/v3/pkg/service/router.SetLowPower (3 handlers)
[GIN-debug] POST   /ras/:dvInd               --> github.com/HYGON-AI/hcu-dcgm/v3/pkg/service/router.SetRasControl (3 handlers)
[GIN-debug] GET    /hslstatus/:dvInd         --> github.com/HYGON-AI/hcu-dcgm/v3/pkg/service/router.GetHSLLogStatus (3 handlers)
[GIN-debug] GET    /linktopo                 --> github.com/HYGON-AI/hcu-dcgm/v3/pkg/service/router.GetPhysicalLinkTopology (3 handlers)
[GIN-debug] POST   /config/save              --> github.com/HYGON-AI/hcu-dcgm/v3/pkg/service/router.SaveDeviceConfig (3 handlers)
[GIN-debug] POST   /config/load              --> github.com/HYGON-AI/hcu-dcgm/v3/pkg/service/router.LoadDeviceConfig (3 handlers)
I1014 06:20:27.229841      44 main.go:61] DCGM 服务监听地址: 0.0.0.0:16081
[GIN-debug] [WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.
Please check https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies for details.
[GIN-debug] Listening and serving HTTP on 0.0.0.0:16081
```

2.4 DCGM 服务接口注册与监听信息

```json
[
  {
    "MinorNumber": 0,
    "PciBusNumber": "0000:00:08.0",
    "DeviceId": "TPXS320016070601",
    "SubSystemName": "K100_AI",
    "Temperature": 51,
    "PowerUsage": 84,
    "PowerCap": 400,
    "MemoryCap": 68702699520,
    "MemoryUsed": 2150400,
    "UtilizationRate": 0,
    "PcieBwMb": 0,
    "Clk": 600,
    "SclkFrequency": [
      "600",
      "700",
      "750",
      "800",
      "900",
      "1000",
      "1106",
      "1200",
      "1270",
      "1319",
      "1400"
    ],
    "Socclk": 309,
    "SocclkFrequency": [
      "309",
      "523",
      "566",
      "618",
      "680",
      "755",
      "850",
      "971"
    ],
    "PerfLevel": "AUTO"
  },
  {
    "MinorNumber": 1,
    "PciBusNumber": "0000:00:09.0",
    "DeviceId": "TPXS320010040501",
    "SubSystemName": "K100_AI",
    "Temperature": 55,
    "PowerUsage": 80,
    "PowerCap": 400,
    "MemoryCap": 68702699520,
    "MemoryUsed": 2150400,
    "UtilizationRate": 0,
    "PcieBwMb": 0,
    "Clk": 600,
    "SclkFrequency": [
      "600",
      "700",
      "750",
      "800",
      "900",
      "1000",
      "1106",
      "1200",
      "1270",
      "1319",
      "1400"
    ],
    "Socclk": 309,
    "SocclkFrequency": [
      "309",
      "523",
      "566",
      "618",
      "680",
      "755",
      "850",
      "971"
    ],
    "PerfLevel": "AUTO"
  }
]
```

