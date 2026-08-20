# HCU DCGM

## License

This project is licensed under the [Apache License 2.0](LICENSE).

Third-party driver interface headers distributed with the source tree are documented in [THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md). The notice file records their source boundary and licenses; the original copyright and license text in each header remains authoritative.


## 组件信息



HCU DCGM 为 HCU 管理提供 Golang 绑定接口，是管理和监控HCU的工具。包括健康状态监控、功率、时钟频率调控，以及资源使用情况统计等。



## 组件使用前置条件

前置条件：DCGM 编译和运行依赖 HCU 驱动提供的接口头文件及底层动态链接库。仓库内 `pkg/dcgm/include/` 保留了当前公开源码构建所需的驱动接口头文件副本；文件来源、版权和许可证见 [THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md)。其中 `rocm_smi.h` 和 `rocm_smi64Config.h` 标注为 AMD NCSA 许可证，公开发布前必须完成法务/合规审批，不得将其许可证改写为其他许可证。

运行时依赖 HCU 底层动态链接库 `libhydmi.so` 和 `librocm_smi64.so`（MIG 功能还需 `libhydmi_mig.so`），均随 HCU 驱动安装提供，安装方式如下。



#### 安装方式一：

1. 安装 HCU 驱动（上述动态链接库包含在驱动中，默认位于 `/opt/hyhal/lib`）

2. 确认驱动安装完成后，系统可正常加载上述库文件



#### 安装方式二：

适用于**没有 HCU 卡、未安装 HCU 驱动**的环境（如仅做编译验证、CI 或开发机无硬件场景）。需从已安装驱动的机器获取动态库，手动部署到本地目录（如 `/your/path/dcgm/lib`）：

1. 将 `librocm_smi64.so.2.8`、`libhydmi.so.1.5` 等版本文件拷贝至上述目录；若缺少 `librocm_smi64.so`、`libhydmi.so` 等软链接，可在该目录下创建：

   - `librocm_smi64.so.2` → `librocm_smi64.so.2.8`，`librocm_smi64.so` → `librocm_smi64.so.2`

   - `libhydmi.so.1` → `libhydmi.so.1.5`，`libhydmi.so` → `libhydmi.so.1`

   - （MIG 功能）`libhydmi_mig.so.1` → `libhydmi_mig.so.1.3`，`libhydmi_mig.so` → `libhydmi_mig.so.1`

   ```bash
   [root@worker-200 lib]# ls -lh
   total 7.6M
   lrwxrwxrwx. 1 root root   17 Jan 14 15:39 libhydmi_mig.so -> libhydmi_mig.so.1
   lrwxrwxrwx. 1 root root   19 Jan 14 15:39 libhydmi_mig.so.1 -> libhydmi_mig.so.1.3
   -rwxr-xr-x. 1 root root 3.0M Jan 14 15:39 libhydmi_mig.so.1.3
   lrwxrwxrwx. 1 root root   13 Jan 14 15:39 libhydmi.so -> libhydmi.so.1
   lrwxrwxrwx. 1 root root   15 Jan 14 15:39 libhydmi.so.1 -> libhydmi.so.1.5
   -rwxr-xr-x. 1 root root 2.0M Jan 14 15:39 libhydmi.so.1.5
   lrwxrwxrwx. 1 root root   18 Jan 14 15:39 librocm_smi64.so -> librocm_smi64.so.2
   lrwxrwxrwx. 1 root root   20 Jan 14 15:39 librocm_smi64.so.2 -> librocm_smi64.so.2.8
   -rwxr-xr-x. 1 root root 1.3M Jan 14 15:39 librocm_smi64.so.2.8
   ```

   

2. 将库目录加入动态库搜索路径：

   ```bash
   export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/your/path/dcgm/lib
   ```

   > 注意：无 HCU 硬件时仅能完成编译或部分接口调用，**无法**进行真实的设备监控与管理。

## 使用流程



> 运行环境要求：**Linux**、已安装 **HCU 驱动**、开启 **CGO**（`CGO_ENABLED=1`）。



### 获取源码



```bash

git clone <仓库地址>

cd dcgm-dcu

export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/opt/hyhal/lib

```



### 编译



```bash

# HTTP 服务（默认端口 16081）

go build -o dcgm-dcu ./pkg/service



# 命令行工具 dcgmi

go build -o dcgmi ./pkg/cmd

```



### 运行



```bash

# 启动 HTTP 服务

./dcgm-dcu -port 16081


# 测试是否启动成功，例如获取卡数量接口
curl -G http://localhost:16081/NumMonitorDevices


# 命令行示例

./dcgmi discovery

```



`pkg/samples/` 目录下提供了设备信息、拓扑、诊断等示例程序，可参考后自行编译运行。



### 作为 Go 库引用



在项目中引入依赖：



```bash

go get github.com/HYGON-AI/hcu-dcgm/v3/pkg/dcgm@latest

```



```go

import "github.com/HYGON-AI/hcu-dcgm/v3/pkg/dcgm"



func main() {

    if err := dcgm.Init(); err != nil {

        panic(err)

    }

    defer dcgm.ShutDown()

    // ...

}

```



本地开发时，可在调用方 `go.mod` 中使用 `replace` 指向本地克隆目录：



```

replace github.com/HYGON-AI/hcu-dcgm/v3 => /your/path/dcgm-dcu

```



核心 API 封装见 `pkg/dcgm/api.go`。



### Docker 部署



镜像不包含驱动动态库，运行时需挂载宿主机驱动目录。详见 `Dockerfile` 及 `deployment/` 下的部署脚本。



```bash

# 构建镜像（需先编译出 dcgm-dcu 二进制并置于项目根目录）

docker build -t dcgm-dcu:latest .



# 运行示例见 deployment/dcgm-dcu-docker.sh

```

<!-- quality-gate smoke test -->

