# HCU DCGM

## License

This project is licensed under the [Apache License 2.0](LICENSE).

Third-party driver interface headers distributed with the source tree are documented in [THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md). The notice file records their source boundary and licenses; the original copyright and license text in each header remains authoritative.

## 组件信息

HCU DCGM 是面向 HCU 的 Go 管理与监控组件，提供健康状态、功耗、时钟频率和资源使用情况等接口；同时提供 HTTP 服务、`dcgmi` 命令行工具和 Go 库。

核心 Go API 位于 `pkg/dcgm/api.go`，HTTP 路由位于 `pkg/service/router/`，示例程序位于 `samples/`。

## 前置条件

- Linux 系统。
- Go 1.22 或更高版本。
- 已启用 CGO（`CGO_ENABLED=1`）。
- 真实设备监控和管理需要安装 HCU 驱动。

构建和运行依赖 HCU 驱动提供的接口头文件及动态库。仓库内 `pkg/dcgm/include/` 保留了公开源码构建所需的驱动接口头文件副本；其来源、版权和许可证见 [THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md)。其中 `rocm_smi.h` 和 `rocm_smi64Config.h` 标注为 AMD NCSA 许可证，公开发布前必须完成法务与合规审批，不得改写其许可证。

运行时依赖 HCU 驱动提供的 `libhydmi.so` 和 `librocm_smi64.so`；使用 MIG 功能时还需要 `libhydmi_mig.so`。HCU 驱动默认将这些库安装在 `/opt/hyhal/lib`。

### 无硬件环境

没有 HCU 卡或未安装 HCU 驱动的开发、持续集成环境，可以从已安装驱动的机器复制所需动态库到本地目录，例如 `/your/path/hcu-dcgm/lib`，并补齐相应的软链接。

```bash
export LD_LIBRARY_PATH="${LD_LIBRARY_PATH}:/your/path/hcu-dcgm/lib"
```

无 HCU 硬件时可以进行编译和部分不依赖设备的测试，但不能进行真实设备监控或管理。

## 使用流程

### 获取源码

```bash
git clone <仓库地址> hcu-dcgm
cd hcu-dcgm

# 使用系统已安装的 HCU 驱动库。
export LD_LIBRARY_PATH="${LD_LIBRARY_PATH}:/opt/hyhal/lib"
```

### 编译

```bash
# HTTP 服务，默认监听 0.0.0.0:16081。
go build -o hcu-dcgm ./pkg/service

# 命令行工具。
go build -o dcgmi ./pkg/cmd
```

### 运行

```bash
# 启动 HTTP 服务。
./hcu-dcgm -port 16081

# 可选：指定监听地址；也可通过 HCU_DCGM_LISTEN 环境变量设置。
./hcu-dcgm -listen 127.0.0.1 -port 16081

# 查询已监视设备数量。
curl -G http://localhost:16081/NumMonitorDevices

# 命令行示例。
./dcgmi discovery
```

`dcgmi` 的完整命令帮助可通过 `./dcgmi --help` 查看；设备信息、拓扑和诊断示例见 `samples/`。

### 作为 Go 库引用

项目 Go 模块路径为 `github.com/HYGON-AI/hcu-dcgm/v3`。

```bash
go get github.com/HYGON-AI/hcu-dcgm/v3/pkg/dcgm@latest
```

```go
package main

import "github.com/HYGON-AI/hcu-dcgm/v3/pkg/dcgm"

func main() {
	if err := dcgm.Init(); err != nil {
		panic(err)
	}
	defer dcgm.ShutDown()

	// 调用 pkg/dcgm/api.go 提供的 API。
}
```

本地开发时，可在调用方的 `go.mod` 中使用 `replace` 指向本地源码：

```go
replace github.com/HYGON-AI/hcu-dcgm/v3 => /your/path/hcu-dcgm
```

## Docker 部署

`Dockerfile` 不包含 HCU 驱动动态库，容器运行时需要挂载宿主机的 `/opt/hyhal` 和设备节点。Dockerfile 期望项目根目录中存在已编译的 `hcu-dcgm` 二进制。

```bash
# 先在项目根目录构建服务二进制。
go build -o hcu-dcgm ./pkg/service

# 构建与部署脚本一致的镜像标签。
docker build -t hcu-dcgm:v2.0.0 .

# 使用 Docker 部署脚本启动容器。
bash deployment/hcu-dcgm-docker.sh
```

Kubernetes 部署配置见 `deployment/hcu-dcgm-k8s.yaml`。