# vDCU 自动化测试工具（vdcu-test）使用手册

vdcu-test 是一套在海光 DCU 宿主机上自动执行 vDCU 虚拟化测试的脚本工具。**登录测试机，敲一条命令，13 个测试方向自动跑完，测试报告自动生成。**

本手册覆盖从部署、配置到手动/自动执行、查看报告、问题排查的完整过程，照做即可完成测试。

## 特点

- 所有环境信息（hy-smi 路径、测试卡号、镜像、仓库地址、容器挂载、用户名密码等）都只写在 `config.json` 里，代码中没有任何写死的内容；config.json 中的值均为**默认值**，现场不同改这一个文件即可（也可用 `VDCU_` 环境变量临时覆盖）；
- 支持两种用法：**零配置直接跑**（环境与默认值一致时）和**改配置后跑**；
- 13 个方向各有一个独立入口，也可一条命令串行跑完 13 个方向；每个方向独立生成报告，单方向失败不中断后续；
- 报告默认脱敏（IP/机器标识/镜像仓库地址替换为占位符），`reports/` 不入 git；
- 202 条用例全部可执行，组件未部署/环境受限的用例如实记 `SKIP` 或 `XFAIL`，不误报失败。

---

## 一、部署

把整个 `vdcu-test` 目录拷贝到测试宿主机（建议放 `/root/tz/` 下，本手册以 `/root/tz/test-tj/vdcu-test` 为例）：

```bash
scp -r vdcu-test root@<测试机IP>:/root/tz/test-tj/
```

如果是从 Windows 拷贝过去的，上机后执行一次（修正行尾与执行权限）：

```bash
cd /root/tz/test-tj/vdcu-test
find . -type f \( -name "*.sh" -o -name "*.py" -o -name "*.json" \) -exec sed -i 's/\r$//' {} +
chmod +x *.sh run.py
```

## 二、环境要求（上机先花 1 分钟确认）

| 检查项 | 怎么确认 | 不满足怎么办 |
| --- | --- | --- |
| 已登录测试机（root 或同等权限） | `whoami` | 自行登录，本工具不处理登录 |
| python3 | `python3 --version` | 安装 python3（3.10 已验证） |
| docker 可用 | `docker version` | 安装并启动 docker |
| hy-smi 可用 | `/opt/hyhal/bin/hy-smi virtual --show-device-info` | 路径不同 → 改 config.json 的 `hy_smi` |
| 测试镜像已导入 | `docker images \| grep bert` | 镜像名不同 → 改 config.json 的 `image`；构建方法见"附录：bert 镜像构建" |
| 数据盘有足够空间 | `df -h /data`（docker 数据目录所在盘） | 至少预留 50G；`run_all13.sh` 在磁盘不足时会自动跳过 K8s 镜像导入 |
| 要跑集群方向（K8S/VOL/HAMI）时 kubectl 可用 | `kubectl get nodes` | 无集群 → config.json 的 `cluster_tests` 改为 `0`，集群方向自动跳过；报 Unauthorized → 见"常见问题" |
| 测试卡空闲 | `hy-smi virtual --show-vdevice-info` | 默认测试 5/6/7 号卡，**不能有正在使用的业务 vDCU**，测试会创建/销毁 vDCU |

## 三、快速开始（30 秒上手）

环境与默认值一致时（hy-smi 在 `/opt/hyhal/bin/hy-smi`、通用镜像 `bert:v1`、测试卡 5/6/7），**什么都不用改**：

```bash
cd /root/tz/test-tj/vdcu-test

./run_all13.sh                 # 一键跑完 13 个测试方向，报告自动生成在 reports/
```

只想跑某一个方向：

```bash
./run_fun.sh                   # 例如功能方向；其余方向脚本见"13 个测试方向"
```

先看看有哪些用例：

```bash
python3 run.py --list          # 列出全部 202 条用例
```

环境与默认值不一致（卡号/镜像/路径不同）时，先按"五、配置"修改 `config.json`，再执行上面的命令。

## 四、手动测试流程（一个方向一个方向盯着看）

适合想逐个方向观察结果的场景。每个方向跑完实时打印结果、自动生成报告。

```bash
# 1. 登录并进入工具目录
ssh root@<测试机IP>
cd /root/tz/test-tj/vdcu-test

# 2. 跑之前先确认（第一次跑建议按顺序做）
cat config.json                          # 看一遍配置，卡号/镜像/模型不符就改
python3 run.py --list                    # 确认工具正常、用例可列出
hy-smi virtual --show-vdevice-info       # 确认测试卡上没有业务 vDCU

# 3. 按需逐个执行方向（每个方向跑完实时打印 PASS/FAIL/SKIP 并生成报告）
./run_fun.sh     # ① FUN  功能（vDCU 创建/销毁/配额）
./run_ctr.sh     # ② CTR  容器挂载与生命周期
./run_com.sh     # ③ COM  vDCU 间通信
./run_trn.sh     # ④ TRN  训练场景（Bert/ModernBERT/Qwen）
./run_inf.sh     # ⑤ INF  推理场景（vLLM 冒烟、bge）
./run_k8s.sh     # ⑥ K8S  集群编排（需 kubectl 可用）
./run_vol.sh     # ⑦ VOL  Volcano 调度（同上）
./run_iso.sh     # ⑧ ISO  隔离性与 QoS
./run_perf.sh    # ⑨ PERF 性能对比
./run_stb.sh     # ⑩ STB  稳定性
./run_mon.sh     # ⑪ MON  监控
./run_cmp.sh     # ⑫ CMP  兼容性
./run_sec.sh     # ⑬ SEC  安全

# 4. 只想跑一条用例 / 看清执行过程
python3 run.py --id VDCU-K8S-002          # 只跑一条
python3 run.py --id VDCU-FUN-002 --verbose  # 边跑边打印每条命令（排查问题最有用）
python3 run.py --category FUN,CTR         # 按类别跑多条
python3 run.py --id VDCU-FUN-013 --allow-destructive   # 全局销毁类用例需加此参数
# 全部参数：python3 run.py --help

# 5. 看报告
ls -lt reports/ | head                    # 最新的 report-<时间>.json / .md
cat reports/report-<时间>.md              # 汇总表：用例/状态/耗时/说明

# 6. 跑完清理
python3 run.py --cleanup                  # 清理测试卡上残留的 vDCU
hy-smi virtual --show-vdevice-info        # 复核已清空
```

每个脚本执行时：屏幕实时打印每条用例的编号、名称和结果；结束后打印汇总行 `PASS x | FAIL x | XFAIL x | XPASS x | SKIP x`；退出码 `0`=全部符合预期，`1`=有 FAIL/XPASS，`2`=环境预检失败。

## 五、自动测试流程（一次发起、无人值守）

适合一次跑完 13 个方向、最后只看汇总的场景。

```bash
cd /root/tz/test-tj/vdcu-test

# 一键启动（后台运行，关掉终端也不受影响）
nohup ./run_all13.sh > /tmp/all13.console.log 2>&1 &
echo $! > /tmp/all13.pid
```

自动化的行为设计：
- 按 `fun → ctr → com → trn → inf → k8s → vol → iso → perf → stb → mon → cmp → sec` 顺序串行执行；
- 每个方向跑完自动生成一份独立报告到 `reports/`；
- **单个方向失败不会中断后续方向**，rc 记入日志继续跑；
- 全部结束后自动销毁测试卡上的 vDCU、清理测试 Pod；
- 配置全部来自 config.json，环境变量可临时覆盖（如 `IMAGE=xxx HS=/xxx CLEANUP_CARDS="5 6" NS=xxx`）；
- K8s 镜像导入保护：`SKIP_IMPORT=1` 强制跳过；数据盘剩余不足 50G 时自动跳过（避免写爆磁盘）。

**过程监控**（不影响脚本执行）：

```bash
tail -f /tmp/all13.log          # 13 个方向实时进度（哪个方向开始/结束、rc 是多少）
tail -f /tmp/all13.console.log  # 用例级 PASS/FAIL 明细
ls -lt reports/ | head          # 已生成多少份报告
```

出现 `ALL_13_DONE` 说明 13 个方向全部执行完毕。

**跑完后的结果确认**：

```bash
cd /root/tz/test-tj/vdcu-test

# 每个方向的汇总（一份报告一个方向）
for f in reports/report-*.md; do echo "== $f =="; grep "用例总数" "$f"; done

# 失败用例清单
for f in reports/report-*.json; do
  python3 -c "import json,sys; d=json.load(open('$f')); \
  bad=[r for r in d['results'] if r['status'] in ('FAIL','XPASS')]; \
  [bad] and print('$f', [(b['id'],b['status']) for b in bad])"
done

# 确认测试卡已清理干净
hy-smi virtual --show-vdevice-info
```

**只要自动跑某一个方向**：`nohup ./run_fun.sh > /tmp/fun.console.log 2>&1 &`（换成任意方向即可）。

**连续多轮执行**（报告按时间戳归档互不覆盖）：

```bash
for i in 1 2 3; do
  python3 run.py --cleanup
  ./run_all13.sh >> /tmp/all13.multi.log 2>&1
done
```

**自动跑异常处理**：

| 现象 | 处理 |
| --- | --- |
| `ALL_13_DONE` 没出现 | `tail -50 /tmp/all13.log` 看卡在哪个方向；`ps aux \| grep run_` 确认进程 |
| 某方向 rc≠0 | 进 reports/ 看该方向最新 md 报告的 FAIL 行；`--verbose` 单跑定位 |
| 中途想停止 | `kill $(cat /tmp/all13.pid)`；随后手动 `python3 run.py --cleanup` |
| K8S 方向报 Unauthorized | 见"常见问题"，拿到新证书后单独重跑 `./run_k8s.sh` |

## 六、配置（config.json）

所有配置集中在 `config.json`，文件开头的 `_comments` 字段对每个配置项都有"默认值"标注，按现场修改后保存即可。也可用同名 `VDCU_` 环境变量临时覆盖：

```bash
VDCU_CARD=3 VDCU_CARD2=4 ./run_fun.sh      # 临时改用 3、4 号卡
```

**常用配置速查**：

| 配置项 | 默认值 | 什么时候要改 |
| --- | --- | --- |
| `card` / `card2` / `full_card` | `5` / `6` / `7` | 现场空闲的物理卡号不同 |
| `hy_smi` | `/opt/hyhal/bin/hy-smi` | 工具路径不同 |
| `work` | `/root/tz` | 模型/数据放在别的目录（不存在会自动创建） |
| `image` | `bert:v1` | 通用测试镜像名不同（构建方法见附录） |
| `k8s_image` | `bert:v1` | K8s Pod 用镜像不同（需已导入集群节点 containerd） |
| `inference_image` | 现场 vLLM 镜像 | vLLM 镜像仓库地址/版本不同（**必须已导入本机**，否则用例会启动超时） |
| `vllm_model` | `/workspace/openbmb--MiniCPM5-1B-SFT` | vLLM 模型路径不同（空 = 跳过该用例） |
| `emb_model` | `/workspace/AI-ModelScope--bge-large-zh-v1.5/...` | bge 模型路径不同（空 = 跳过 INF-206） |
| `cluster_tests` | `1` | 没有集群环境时改为 `0` |
| `cluster_mutating_tests` | `1` | 不想让测试改变集群状态（插件重启等）时改为 `0` |
| `cluster_namespace` / `k8s_node` | `vdcu-test` / 空 | 集群命名空间；测试镜像只导入本机时填本机节点名 |
| `cleanup_cards` | `5 6 7` | run_all13.sh 结束清理的卡号列表 |

其余配置项（`hyhal_dir`、`vdev_conf_dir`、`dri_dir`、`kfd_device`、`mkfd_device`、`dtk_hyhal_lib`、`assets_container_dir`、`workspace_container_dir` 等设备路径与容器挂载）是驱动/系统的标准路径，**一般不用动**，完整清单见 config.json 内的 `_comments`。

**用户名密码**：本工具代码里没有任何写死的用户名和密码。测试机的 root 登录由您自己完成；工具只在一种情况下用到账号——**私有镜像仓库需要 `docker login`**。在 config.json 里填：

```json
"registry":          "仓库地址:端口",   // 如 harbor.example.com:5443（留空则自动从 inference_image 解析）
"registry_user":     "您的仓库用户名",
"registry_password": "您的仓库密码"
```

运行时工具自动执行 `docker login <registry> -u <用户名> --password-stdin`，密码只走 stdin 管道传入，**不会出现在命令行、日志和测试报告里**。登录失败只告警不中断（镜像已提前导入本地时不需要登录）。

## 七、测试报告怎么看

- 报告自动生成在 `reports/`，每次运行一份 `report-<时间>.json` 和 `report-<时间>.md`；
- 状态含义：`PASS` 通过；`FAIL` 失败；`SKIP` 环境不具备（组件未部署/配置未开启），跳过；`XFAIL` 已知缺陷，按预期失败；`XPASS` 已知缺陷意外修复；
- 退出码：`0` 没有未预期失败；`1` 有 FAIL/XPASS（看报告定位）；`2` 环境预检失败（按提示逐条处理）；
- 报告默认**脱敏**：内网 IP、机器标识、镜像仓库地址替换为占位符，可放心对外；排查细节时把 config.json 的 `sanitize` 改为 `0` 再跑；
- 每条用例的命令与输出记录在 JSON 报告的 `steps` 字段（`$ 命令` + `[rc=返回码]` + 输出）；
- `reports/` 已在 `.gitignore` 中忽略，不会提交到 git。

## 八、常见问题排查

| 现象 | 原因与处理 |
| --- | --- |
| 一跑就报"环境预检失败"（退出码 2） | 按屏幕逐条提示处理（hy-smi/docker/kubectl/镜像/卡号），对照"二、环境要求"改 config.json 或补齐环境 |
| K8S/VOL/HAMI 方向全是 SKIP | 集群组件没部署或 `cluster_tests=0`，属设计行为（组件缺失记 SKIP 不误报）。要测集群方向：确保 kubectl 可用、`cluster_tests=1` |
| K8S 用例报 Unauthorized | kubeconfig 证书不被集群认可（常见于集群更换过客户端认证 CA）。找集群管理员要当前 CA 签发的 admin 证书，覆盖 `/root/.kube/config` 后重跑 `./run_k8s.sh` |
| 大量 SKIP | 组件未部署或配置未开启（如 Volcano/HAMi 未装、TE 内网源不可达、数据集未部署），报告里有具体 SKIP 原因，属正常行为 |
| SKIP 提示"需 --allow-destructive" | 全局销毁类用例（如 FUN-013 销毁全部 vDCU）需要显式授权：`python3 run.py --id VDCU-FUN-013 --allow-destructive`。该参数会清除测试卡上所有 vDCU，执行前确认没有业务在用 |
| 单条用例失败 | `python3 run.py --id <编号> --verbose` 重跑，看报告 steps 里的命令输出定位 |
| 卡上有脏 vDCU 导致用例报错 | 先 `python3 run.py --cleanup` 再重跑；仍异常时检查 `/etc/vdev/` 是否有残留配置（驱动 generic error 常见诱因） |
| 跑了一半 Ctrl-C | `python3 run.py --cleanup` 清理测试卡 |
| 数据盘满了 | docker 数据目录所在盘需预留空间；`run_all13.sh` 在不足 50G 时会自动跳过 K8s 镜像导入并告警；可清理悬空镜像 `docker image prune -f` |
| vLLM 用例报 startup timeout | 确认 config.json 的 `inference_image` 是本机已导入的镜像；`docker images \| grep vllm` 核对 |
| 报告含敏感信息 | 确认 config.json 的 `sanitize=1`（默认已开）；额外要隐藏的主机名/域名填 `sanitize_tokens` |

## 九、13 个测试方向一览

| 入口 | 方向 | 测什么 |
| --- | --- | --- |
| `run_fun.sh` | FUN 功能 | vDCU 创建、等分、配额、更新、销毁、权限、加密、异常处理 |
| `run_ctr.sh` | CTR 容器挂载 | 容器挂载、多容器、同卡/跨卡、CU 共享、设备可见、资源回收、vLLM 并发启动 |
| `run_com.sh` | COM 通信 | vDCU 间集合通信（当前驱动多卡通信受限，见 XFAIL 记录） |
| `run_trn.sh` | TRN 训练 | Bert/ModernBERT/Qwen 训练、多实例并发 |
| `run_inf.sh` | INF 推理 | 模型生成、vLLM 服务冒烟、多模态、Embedding |
| `run_k8s.sh` | K8S 编排 | 节点资源注册、Pod 申请/回收、资源不足 Pending、插件恢复、HAMi |
| `run_vol.sh` | VOL 调度 | Volcano 调度、PodGroup、Queue、优先级、混合负载 |
| `run_iso.sh` | ISO 隔离 | 显存/CU 配额隔离、噪声邻居、故障隔离 |
| `run_perf.sh` | PERF 性能 | 整卡基线、切分线性度、并发扩展、虚拟化开销 |
| `run_stb.sh` | STB 稳定性 | 生命周期循环、进程 kill 恢复、OOM 恢复 |
| `run_mon.sh` | MON 监控 | 指标可见范围、字段正确性、dcu-exporter 上报 |
| `run_cmp.sh` | CMP 兼容性 | 基础镜像、框架/运行时、vLLM 多版本兼容 |
| `run_sec.sh` | SEC 安全 | 最小设备挂载、越权防护 |

## 十、目录结构

```text
vdcu-test/
├── run.py                 # 主程序（python3 run.py --help 看全部参数）
├── run.sh                 # 通用入口（python3 run.py "$@"）
├── run_all13.sh           # 13 个方向串行总入口（含磁盘保护）
├── run_<方向>.sh × 13     # 各方向单独入口
├── config.json            # 全部环境配置（默认值已标注，改这一个文件即可）
├── vdcu/core.py           # 框架核心（用例注册、执行、报告、脱敏）
├── cases/                 # 各方向测试用例
├── assets/                # 容器内执行的负载脚本
└── reports/               # 测试报告（自动生成，不入 git）
```

## 附录：bert:v1 训练镜像构建

测试通用镜像按以下方法构建（Dockerfile 内容）：

```dockerfile
# 基础镜像来自内部镜像源 harbor
FROM harbor.sourcefind.cn:5443/hcu/admin/base/pytorch:2.12.0-ubuntu22.04-dtk26.04-py3.10-minimal

WORKDIR /workspace

# 训练依赖：transformers 需支持 ModernBERT（>=4.48）
RUN pip install transformers datasets scikit-learn pandas tqdm \
    -i https://pypi.tuna.tsinghua.edu.cn/simple
```

```bash
docker build --network=host -f Dockerfile -t bert:v1 <空目录上下文>   # 注意：构建上下文用空目录，避免把大数据集打进 tar
docker images bert:v1    # 验证：torch 2.12.0、transformers 5.20.2
```

K8s 方向还需把镜像导入集群节点的 containerd：`docker save bert:v1 | ctr -n k8s.io images import -`
