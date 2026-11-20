# vDCU 虚拟化自动化测试

本目录提供 vDCU 功能、容器业务以及 Kubernetes、Volcano、HAMi 集群场景的可执行测试用例。框架只使用 Python 3 标准库，脚本在被测宿主机上运行，结果写入 `reports/`。

## 1. 运行环境

| 项目 | 配置或说明 |
| --- | --- |
| 宿主机 | Ubuntu 22.04，使用 root 或具有 Docker、kubectl 权限的账号 |
| 加速卡 | 海光 BW1000，实际卡号通过 `config.json` 配置 |
| 管理工具 | `/opt/hyhal/bin/hy-smi` |
| 容器运行时 | Docker；Kubernetes 节点使用 containerd |
| 通用测试镜像 | `bert:v1` |
| 模型、数据和脚本目录 | `/root/tz` |
| Kubernetes | 使用现场集群版本，不在脚本中写死版本号 |

脚本会根据节点实际注册的资源名选择 Kubernetes 资源。目前已验证的普通设备插件资源名为 `hygon.com/dcu`；部署了 vHCU/HAMi 后，也可以通过配置使用对应的资源名。Volcano 或 HAMi 未部署时，相应用例报告为 `SKIP`，不会误报为失败。

## 2. 目录结构

```text
vdcu-test-automation/
├── run.py                 # 参数解析、预检、执行和报告
├── run.sh                 # Linux 快捷入口
├── config.json            # 宿主机、镜像和集群配置
├── vdcu/core.py           # 用例注册、命令执行、断言和报告
├── cases/
│   ├── fun.py             # vDCU 创建、查询、更新、销毁和资源约束
│   ├── ctr.py             # 容器挂载、复用、隔离和生命周期
│   ├── workload.py        # Bert、模型生成、vLLM、通信和隔离
│   ├── system.py          # 监控、兼容性和性能比例检查
│   └── cluster.py         # Kubernetes、Volcano 和 HAMi
├── assets/                # 容器内执行的最小业务脚本
└── reports/               # JSON 和 Markdown 结果
```

## 3. 快速开始

```bash
cd /root/tz/vdcu-test-automation

# 检查环境并列出用例
python3 run.py --list

# 先执行单机 vDCU 和容器用例
python3 run.py --category FUN,CTR,TRN,INF,COM,ISO,MON,CMP

# 执行 Kubernetes、Volcano、HAMi 用例
export VDCU_CLUSTER_TESTS=1
export VDCU_CLUSTER_MUTATING_TESTS=0
python3 run.py --category K8S,VOL,HAMI --no-report

# 查看单条用例的命令和输出
python3 run.py --id VDCU-K8S-002 --verbose

# 全部测试结束后清理脚本创建的测试资源
python3 run.py --cleanup
```

`VDCU_CLUSTER_MUTATING_TESTS=1` 才会执行可能改变集群状态的资源不足、设备插件重启等检查。执行前应确认集群允许进行这些操作。

退出码为 `0` 表示没有未预期失败，`1` 表示存在 `FAIL` 或 `XPASS`，`2` 表示预检失败。

## 4. 配置

`config.json` 可直接修改，也可使用同名的 `VDCU_` 环境变量覆盖。

| 配置项 | 默认值 | 说明 |
| --- | --- | --- |
| `hy_smi` | `/opt/hyhal/bin/hy-smi` | vDCU 管理工具绝对路径 |
| `card`、`card2` | `5`、`6` | 单机测试使用的两张物理卡 |
| `full_card` | `7` | 整卡基线用物理卡 |
| `physical_cards` | `8` | 物理卡数量，用于状态检查 |
| `image` | `bert:v1` | Bert 和通用容器测试镜像 |
| `inference_image` | 现场可用的 vLLM 镜像 | vLLM 冒烟镜像 |
| `vllm_model` | 空 | 模型路径为空时跳过 vLLM 用例 |
| `work` | `/root/tz` | 模型、数据和日志目录 |
| `cluster_namespace` | `default` | 集群测试命名空间 |
| `k8s_image` | `bert:v1` | Kubernetes Pod 测试镜像 |
| `cluster_tests` | `0` | 设为 `1` 才执行集群用例 |
| `cluster_mutating_tests` | `0` | 设为 `1` 才执行变更型集群用例 |

请先确认 `card`、`card2` 和 `full_card` 没有承载业务，再执行会创建或销毁 vDCU 的用例。

## 5. 测试内容

脚本按以下类别组织用例：

- `FUN`：设备查询、创建、等分、配额、更新、销毁、权限、加密开关和异常处理。
- `CTR`：单容器、多容器、同卡、跨卡、CU 共享、设备可见性和资源回收。
- `TRN`、`INF`：Bert 训练、模型生成和 vLLM OpenAI 接口冒烟。
- `COM`、`ISO`：通信入口以及显存、CU 隔离检查；已知缺陷会以 `XFAIL` 和 ISSUE 号记录。
- `K8S`：节点资源注册、Pod 资源申请、设备可见性、Pending、删除回收和设备插件恢复。
- `HAMI`：HAMi/vHCU 资源注册、Pod 申请、设备可见性和资源回收。
- `VOL`：Volcano 组件、deviceshare、调度器、PodGroup、Queue、优先级、调度策略以及训练/推理混合负载入口。
- `MON`、`CMP`、`PERF`：管理工具字段、运行时兼容性和与整卡基线的比例检查。

### Kubernetes 资源申请

`VDCU-K8S-001` 会读取节点 `capacity` 和 `allocatable`，自动识别可用的 DCU 资源名。`VDCU-K8S-002` 创建临时 Pod，申请一份资源，等待 `Running`，在容器内探测设备，再删除 Pod 并检查资源回收。镜像没有 `hy-smi` 时，用例会保留 Pod 运行结果并说明设备探测限制。

### Volcano 与 HAMi

Volcano 和 HAMi 测试只在组件及其 CRD、资源名实际存在时执行。组件未部署、资源名不匹配或配置不完整时，报告会给出具体原因并标记 `SKIP`。因此在不同现场环境中可以直接复用脚本，不需要修改用例逻辑来掩盖缺失组件。

## 6. 报告

默认报告文件位于 `reports/`，包含每条用例的状态、执行耗时、关键命令输出和失败原因。报告中的 `PASS`、`FAIL`、`SKIP`、`XFAIL`、`XPASS` 分别表示实际通过、实际失败、环境未提供、已知缺陷预期失败和已知缺陷意外通过。

执行集群用例前后建议保存以下现场信息，便于结果复核：

```bash
kubectl version
kubectl get nodes -o wide
kubectl get pods -A
kubectl get crd | grep -Ei 'volcano|hami|device'
kubectl get nodes -o json
```
