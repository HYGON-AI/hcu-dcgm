#!/usr/bin/env bash
#
# HCU DCGM 接口兼容性验证工具
#
# 默认验证：
#   1. HTTP GET 接口：检查路由、参数和运行时返回值
#   2. dcgmi 只读命令：检查 CLI 初始化、命令解析和 API 调用
#
# 默认不会执行诊断压测。需要时使用 --include-diag，或显式指定
# --diag-level 1..4；level 3/4 可能执行压力测试并改变设备负载。
#
# 兼容旧用法：
#   bash test_get_apis.sh [dvInd]
#
# 常用用法：
#   ./test_get_apis.sh
#   ./test_get_apis.sh --device 0 --url http://localhost:16081
#   ./test_get_apis.sh --http-only
#   ./test_get_apis.sh --dcgmi-only --dcgmi /usr/local/bin/dcgmi
#   ./test_get_apis.sh --include-diag --diag-level 1
#   ./test_get_apis.sh --fail-on-warn
#
# 环境变量：
#   DCGM_URL              HTTP 服务地址，默认 http://localhost:16081
#   DCGMI_BIN             dcgmi 可执行文件路径；不设置时从 PATH 查找
#   DCGM_TEST_OUTPUT_DIR  报告根目录；默认使用 ${TMPDIR:-/tmp}/dcgm-compatibility

set -uo pipefail

SCRIPT_NAME="$(basename "$0")"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
BASE="${DCGM_URL:-http://localhost:16081}"
DCGMI_BIN="${DCGMI_BIN:-}"
DV="0"
TIMEOUT="8"
CONNECT_TIMEOUT="3"
OUTPUT_ROOT="${DCGM_TEST_OUTPUT_DIR:-${TMPDIR:-/tmp}/dcgm-compatibility}"
RUN_HTTP=1
RUN_DCGMI=1
RUN_COMPAT=0
INCLUDE_DIAG=0
DIAG_LEVEL="1"
FAIL_ON_WARN=0
COMPAT_CARD_MODEL="${DCGM_COMPAT_CARD_MODEL:-}"
COMPAT_DRIVER_VERSION="${DCGM_COMPAT_DRIVER_VERSION:-}"
COMPAT_DTK_VERSION="${DCGM_COMPAT_DTK_VERSION:-}"
POSITIONAL_DEVICE_SET=0

usage() {
    cat <<EOF
用法：
  ${SCRIPT_NAME} [dvInd]
  ${SCRIPT_NAME} [选项]

选项：
  -u, --url URL             DCGM HTTP 地址（默认：${BASE}）
  -d, --device INDEX        设备索引 dvInd（默认：0）
      --card-model MODEL      可选：验证 GET /compatible 的卡型号
      --driver-version VER    可选：验证 GET /compatible 的驱动版本
      --dtk-version VER       可选：验证 GET /compatible 的 DTK 版本
      --compatibility ...     仅执行 /compatible 版本兼容性检查（需同时提供上面三个参数）
      --http-only            只验证 HTTP 接口
      --dcgmi-only           只验证 dcgmi 命令
      --include-diag         额外执行 dcgmi diag r 1
      --diag-level LEVEL     执行 dcgmi diag r LEVEL（1..4）
  -t, --timeout SECONDS     单个 HTTP/CLI 测试超时（默认：${TIMEOUT}）
      --connect-timeout SEC HTTP 建连超时（默认：${CONNECT_TIMEOUT}）
  -o, --output-dir DIR      报告根目录
      --fail-on-warn        将 WARN 结果也作为失败退出
  -h, --help                显示帮助

退出码：
  0  没有 FAIL（默认允许 WARN）
  1  存在 FAIL，或 --fail-on-warn 时存在 WARN
  2  参数错误或必需工具不可用

说明：
  - PASS 表示调用成功并且没有检测到明确错误。
  - WARN 通常表示接口存在，但当前机器缺少设备、进程、虚拟设备或 sysfs 条件。
  - FAIL 表示连接失败、非预期 HTTP 状态、dcgmi 非零退出或输出包含明确错误。
EOF
}

fatal() {
    printf 'ERROR: %s\n' "$*" >&2
    exit 2
}

is_integer() {
    [[ "$1" =~ ^[0-9]+$ ]]
}

while (($# > 0)); do
    case "$1" in
        -u|--url)
            (($# >= 2)) || fatal "$1 需要参数"
            BASE="$2"
            shift 2
            ;;
        -d|--device)
            (($# >= 2)) || fatal "$1 需要参数"
            DV="$2"
            POSITIONAL_DEVICE_SET=1
            shift 2
            ;;
        --dcgmi)
            (($# >= 2)) || fatal "$1 需要参数"
            DCGMI_BIN="$2"
            shift 2
            ;;
        --card-model)
            (($# >= 2)) || fatal "$1 需要参数"
            COMPAT_CARD_MODEL="$2"
            shift 2
            ;;
        --driver-version)
            (($# >= 2)) || fatal "$1 需要参数"
            COMPAT_DRIVER_VERSION="$2"
            shift 2
            ;;
        --dtk-version)
            (($# >= 2)) || fatal "$1 需要参数"
            COMPAT_DTK_VERSION="$2"
            shift 2
            ;;
        --compatibility)
            RUN_HTTP=0
            RUN_DCGMI=0
            RUN_COMPAT=1
            shift
            ;;
        --http-only)
            RUN_DCGMI=0
            shift
            ;;
        --dcgmi-only)
            RUN_HTTP=0
            shift
            ;;
        --include-diag)
            INCLUDE_DIAG=1
            shift
            ;;
        --diag-level)
            (($# >= 2)) || fatal "$1 需要参数"
            DIAG_LEVEL="$2"
            INCLUDE_DIAG=1
            shift 2
            ;;
        -t|--timeout)
            (($# >= 2)) || fatal "$1 需要参数"
            TIMEOUT="$2"
            shift 2
            ;;
        --connect-timeout)
            (($# >= 2)) || fatal "$1 需要参数"
            CONNECT_TIMEOUT="$2"
            shift 2
            ;;
        -o|--output-dir)
            (($# >= 2)) || fatal "$1 需要参数"
            OUTPUT_ROOT="$2"
            shift 2
            ;;
        --fail-on-warn)
            FAIL_ON_WARN=1
            shift
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        --)
            shift
            (($# <= 1)) || fatal "只支持一个位置参数 dvInd"
            if (($# == 1)); then
                [[ "$POSITIONAL_DEVICE_SET" == 0 ]] || fatal "--device 与位置参数不能同时使用"
                DV="$1"
            fi
            shift || true
            ;;
        -*)
            fatal "未知选项：$1"
            ;;
        *)
            [[ "$POSITIONAL_DEVICE_SET" == 0 ]] || fatal "设备索引重复指定"
            [[ "$DV" == "0" ]] || fatal "只支持一个位置参数 dvInd"
            DV="$1"
            POSITIONAL_DEVICE_SET=1
            shift
            ;;
    esac
done

is_integer "$DV" || fatal "设备索引必须是非负整数：$DV"
is_integer "$TIMEOUT" || fatal "--timeout 必须是非负整数：$TIMEOUT"
is_integer "$CONNECT_TIMEOUT" || fatal "--connect-timeout 必须是非负整数：$CONNECT_TIMEOUT"
is_integer "$DIAG_LEVEL" || fatal "--diag-level 必须是 1..4：$DIAG_LEVEL"
((DIAG_LEVEL >= 1 && DIAG_LEVEL <= 4)) || fatal "--diag-level 必须是 1..4：$DIAG_LEVEL"
[[ "$RUN_HTTP$RUN_DCGMI$RUN_COMPAT" != "000" ]] || fatal "至少选择 HTTP、dcgmi 或 /compatible 测试之一"
if ((RUN_COMPAT)); then
    [[ -n "$COMPAT_CARD_MODEL" && -n "$COMPAT_DRIVER_VERSION" && -n "$COMPAT_DTK_VERSION" ]] || fatal "兼容性检查需要同时提供 --card-model、--driver-version、--dtk-version"
elif ((RUN_HTTP)) && [[ -n "$COMPAT_CARD_MODEL$COMPAT_DRIVER_VERSION$COMPAT_DTK_VERSION" ]]; then
    [[ -n "$COMPAT_CARD_MODEL" && -n "$COMPAT_DRIVER_VERSION" && -n "$COMPAT_DTK_VERSION" ]] || fatal "兼容性检查需要同时提供 --card-model、--driver-version、--dtk-version"
    RUN_COMPAT=1
fi

if ! mkdir -p "$OUTPUT_ROOT"; then
    fatal "无法创建报告目录：$OUTPUT_ROOT"
fi
RUN_DIR="${OUTPUT_ROOT}/$(date +%Y%m%d_%H%M%S)_pid_$$"
mkdir -p "$RUN_DIR" || fatal "无法创建本次报告目录：$RUN_DIR"
LOG_FILE="${RUN_DIR}/compatibility.log"
RESULTS_FILE="${RUN_DIR}/results.tsv"
HTTP_LAST_BODY=""

pass=0
warn=0
fail=0
skip=0
total=0
sequence=0

if [[ -t 1 ]]; then
    GREEN=$'\033[0;32m'
    RED=$'\033[0;31m'
    YELLOW=$'\033[0;33m'
    CYAN=$'\033[0;36m'
    DIM=$'\033[2m'
    RESET=$'\033[0m'
else
    GREEN=''
    RED=''
    YELLOW=''
    CYAN=''
    DIM=''
    RESET=''
fi

log() {
    printf '%s\n' "$*" >> "$LOG_FILE"
}

section() {
    log ""
    log "================================================================"
    log " $1"
    log "================================================================"
    printf '\n%s== %s ==%s\n' "$CYAN" "$1" "$RESET"
}

compact_file() {
    local file="$1"
    if [[ -s "$file" ]]; then
        tr '\r\n\t' '   ' < "$file" | sed 's/  */ /g' | cut -c1-220
    else
        printf '%s' ""
    fi
}

first_nonempty_line() {
    local file="$1"
    if [[ -s "$file" ]]; then
        sed '/^[[:space:]]*$/d' "$file" | head -1 | tr '\r\n\t' '   ' | cut -c1-220
    fi
}

http_reason() {
    local body_file="$1"
    local error_file="$2"
    local reason

    reason="$(sed -n 's/.*"error"[[:space:]]*:[[:space:]]*"\([^"\\]*\).*/\1/p' "$body_file" 2>/dev/null | head -1)"
    if [[ -z "$reason" ]]; then
        reason="$(sed -n 's/.*"message"[[:space:]]*:[[:space:]]*"\([^"\\]*\).*/\1/p' "$body_file" 2>/dev/null | head -1)"
    fi
    if [[ -z "$reason" ]]; then
        reason="$(compact_file "$body_file")"
    fi
    if [[ -z "$reason" ]]; then
        reason="$(compact_file "$error_file")"
    fi
    [[ -n "$reason" ]] || reason="无响应内容"
    printf '%s' "$reason"
}

sanitize_tsv() {
    tr '\r\n\t' '   ' <<< "$1" | cut -c1-500
}

record_result() {
    local status="$1" category="$2" label="$3" target="$4" reason="$5"
    local status_color="$RESET"
    case "$status" in
        PASS) status_color="$GREEN"; ((pass++)) ;;
        WARN) status_color="$YELLOW"; ((warn++)) ;;
        FAIL) status_color="$RED"; ((fail++)) ;;
        SKIP) status_color="$YELLOW"; ((skip++)) ;;
    esac
    ((total++))
    printf '%b%-5s%b %-5s %-38s %s\n' "$status_color" "$status" "$RESET" "$category" "$label" "$reason"
    printf '%s\t%s\t%s\t%s\t%s\n' \
        "$status" "$category" "$(sanitize_tsv "$label")" \
        "$(sanitize_tsv "$target")" "$(sanitize_tsv "$reason")" >> "$RESULTS_FILE"
}

skip_test() {
    local category="$1" label="$2" target="$3" reason="$4"
    log ""
    log "[$category] $label"
    log "  Target: $target"
    log "  Status: SKIP"
    log "  Reason: $reason"
    record_result "SKIP" "$category" "$label" "$target" "$reason"
}

# HTTP GET 测试。第五个参数为 environment 时，非 2xx 结果记为 WARN，
# 用于区分“接口存在但当前机器不具备条件”和真正的接口回归。
http_get() {
    local label="$1" path="$2" note="${3:-}" severity="${4:-fail}"
    local body_file error_file http_code curl_rc reason status

    ((sequence++))
    body_file="${RUN_DIR}/http_${sequence}.body"
    error_file="${RUN_DIR}/http_${sequence}.stderr"
    HTTP_LAST_BODY="$body_file"
    : > "$body_file"
    : > "$error_file"

    http_code="$(curl -sS --connect-timeout "$CONNECT_TIMEOUT" --max-time "$TIMEOUT" \
        -o "$body_file" -w '%{http_code}' "${BASE}${path}" 2>"$error_file")"
    curl_rc=$?
    [[ "$http_code" =~ ^[0-9]{3}$ ]] || http_code="000"

    reason="$(http_reason "$body_file" "$error_file")"
    if ((curl_rc != 0)) || [[ "$http_code" == "000" ]]; then
        status="FAIL"
        reason="连接或请求失败（curl rc=${curl_rc}）：${reason}"
    elif [[ "$http_code" =~ ^2[0-9][0-9]$ ]]; then
        status="PASS"
        reason="HTTP ${http_code}"
    elif [[ "$severity" == "environment" ]]; then
        status="WARN"
        reason="HTTP ${http_code}；当前环境可能不支持：${reason}"
    else
        status="FAIL"
        reason="HTTP ${http_code}：${reason}"
    fi

    log ""
    log "[HTTP] $label"
    log "  GET: ${BASE}${path}"
    [[ -n "$note" ]] && log "  Note: $note"
    log "  curl_rc: $curl_rc"
    log "  http_code: $http_code"
    log "  status: $status"
    log "  reason: $reason"
    if [[ -s "$error_file" ]]; then
        log "  stderr:"
        cat "$error_file" >> "$LOG_FILE"
    fi
    log "  response:"
    cat "$body_file" >> "$LOG_FILE"

    record_result "$status" "HTTP" "$label" "$path" "$reason"
}

cli_output_has_error() {
    local output_file="$1"
    grep -Eiq '(^|[^[:alpha:]])(error|failed|failure|panic|segmentation fault|invalid|parse error|initialization failed)([^[:alpha:]]|$)' "$output_file"
}

cli_output_has_warning() {
    local output_file="$1"
    grep -Eiq '(warning|warn|unsupported|not supported|not found|no [^[:space:]]+ found|unavailable|cannot access)' "$output_file"
}

cli_warning_reason() {
    local output_file="$1"
    local reason
    reason="$(grep -Ei '(warning|warn|unsupported|not supported|not found|no [^[:space:]]+ found|unavailable|cannot access)' "$output_file" 2>/dev/null | head -1 | tr '\r\n\t' '   ' | cut -c1-220)"
    [[ -n "$reason" ]] || reason="$(first_nonempty_line "$output_file")"
    printf '%s' "$reason"
}

cli_environment_stderr() {
    local output_file="$1"
    grep -Eiq '(virtual device.*(not exist|does not exist)|unsupported|not supported|not found|unavailable|cannot access)' "$output_file"
}

# dcgmi 命令测试。命令参数以数组形式传入，避免设备索引或路径被再次拆词。
run_dcgmi() {
    local label="$1"
    shift
    local stdout_file stderr_file combined_file rc reason status command_text arg

    ((sequence++))
    stdout_file="${RUN_DIR}/dcgmi_${sequence}.stdout"
    stderr_file="${RUN_DIR}/dcgmi_${sequence}.stderr"
    combined_file="${RUN_DIR}/dcgmi_${sequence}.output"

    command_text="$(printf '%q ' "$DCGMI_BIN" "$@")"
    timeout "$TIMEOUT" "$DCGMI_BIN" "$@" >"$stdout_file" 2>"$stderr_file"
    rc=$?
    if ((rc == 124)); then
        status="FAIL"
        reason="命令执行超时（${TIMEOUT}s）"
    fi
    cat "$stdout_file" "$stderr_file" > "$combined_file"

    if ((rc != 0)) && [[ "${status:-}" != "FAIL" ]]; then
        status="FAIL"
        reason="退出码 ${rc}：$(first_nonempty_line "$stderr_file")"
        [[ "$reason" == "退出码 ${rc}：" ]] && reason="退出码 ${rc}：$(first_nonempty_line "$stdout_file")"
    elif cli_environment_stderr "$stderr_file"; then
        status="WARN"
        reason="命令成功，但标准错误提示当前环境不完整：$(first_nonempty_line "$stderr_file")"
    elif cli_output_has_error "$stderr_file"; then
        status="FAIL"
        reason="命令成功但标准错误包含错误：$(first_nonempty_line "$stderr_file")"
    elif cli_output_has_warning "$combined_file"; then
        status="WARN"
        reason="命令成功，但输出提示当前环境不完整：$(cli_warning_reason "$combined_file")"
    else
        status="PASS"
        reason="退出码 0；$(first_nonempty_line "$stdout_file")"
        [[ "$reason" == "退出码 0；" ]] && reason="退出码 0"
    fi

    log ""
    log "[DCGMI] $label"
    log "  Command: $command_text"
    log "  exit_code: $rc"
    log "  status: $status"
    log "  reason: $reason"
    log "  stdout:"
    cat "$stdout_file" >> "$LOG_FILE"
    if [[ -s "$stderr_file" ]]; then
        log "  stderr:"
        cat "$stderr_file" >> "$LOG_FILE"
    fi

    record_result "$status" "DCGMI" "$label" "$command_text" "$reason"
}

preflight() {
    section "环境预检查"
    log "开始时间: $(date '+%Y-%m-%d %H:%M:%S %z')"
    log "服务地址: $BASE"
    log "设备索引: $DV"
    log "报告目录: $RUN_DIR"

    if ((RUN_HTTP || RUN_COMPAT)); then
        if ! command -v curl >/dev/null 2>&1; then
            record_result "FAIL" "CHECK" "curl" "PATH" "未找到 curl；HTTP 测试无法执行"
        else
            record_result "PASS" "CHECK" "curl" "$(command -v curl)" "工具可用"
        fi
    fi

    if ((RUN_DCGMI)); then
        if ! command -v timeout >/dev/null 2>&1; then
            record_result "FAIL" "CHECK" "timeout" "PATH" "未找到 timeout；无法安全限制 dcgmi 单命令执行时间"
        fi
        if [[ -z "$DCGMI_BIN" ]]; then
            DCGMI_BIN="$(command -v dcgmi 2>/dev/null || true)"
            [[ -z "$DCGMI_BIN" && -x "${SCRIPT_DIR}/dcgmi" ]] && DCGMI_BIN="${SCRIPT_DIR}/dcgmi"
        fi
        if [[ -z "$DCGMI_BIN" || ! -x "$DCGMI_BIN" ]]; then
            record_result "FAIL" "CHECK" "dcgmi" "${DCGMI_BIN:-PATH}" "未找到可执行的 dcgmi；CLI 测试无法执行"
        else
            record_result "PASS" "CHECK" "dcgmi" "$DCGMI_BIN" "工具可用"
        fi
    fi

    if ((INCLUDE_DIAG)); then
        printf '%s提示：将执行 dcgmi diag r %s；该命令可能耗时并改变设备负载。%s\n' "$YELLOW" "$DIAG_LEVEL" "$RESET"
        log "WARNING: diagnostic test enabled: dcgmi diag r $DIAG_LEVEL"
    fi
    if ((RUN_COMPAT)); then
        log "兼容性检查: ${COMPAT_CARD_MODEL}, driver=${COMPAT_DRIVER_VERSION}, dtk=${COMPAT_DTK_VERSION}"
    fi
}

run_http_tests() {
    section "HTTP GET 接口"

    http_get "NumMonitorDevices"            "/NumMonitorDevices"
    local device_count="$(grep -oE '[0-9]+' "$HTTP_LAST_BODY" 2>/dev/null | head -1 || true)"
    if ! is_integer "${device_count:-}"; then
        device_count=""
    fi

    http_get "DevTypeID"                    "/DevTypeID?dvInd=${DV}"
    http_get "DevTypeName"                  "/DevTypeName?dvInd=${DV}"
    http_get "DevSubsystemId"               "/DevSubsystemId?dvInd=${DV}"
    http_get "DevSubsystemName"             "/DevSubsystemName?dvInd=${DV}"
    http_get "MaxPower"                     "/MaxPower?dvInd=${DV}"
    http_get "MemInfo"                      "/MemInfo?dvInd=${DV}&memType=vram"
    http_get "HSLErrorStatus"               "/HSLErrorStatus?dvInd=${DV}"
    http_get "HSLHiveId"                    "/HSLHiveId?dvInd=${DV}"
    http_get "Version"                      "/Version?component=DRIVER"
    if ((RUN_COMPAT)); then
        http_get "compatible" "/compatible?cardModel=${COMPAT_CARD_MODEL}&driverVersion=${COMPAT_DRIVER_VERSION}&dtkVersion=${COMPAT_DTK_VERSION}" \
            "卡型号=${COMPAT_CARD_MODEL}，Driver=${COMPAT_DRIVER_VERSION}，DTK=${COMPAT_DTK_VERSION}"
    fi

    http_get "AllDeviceInfos"               "/AllDeviceInfos"
    http_get "DeviceInfos"                  "/DeviceInfos"
    http_get "DeviceStatus"                 "/DeviceStatus"
    http_get "CollectDeviceMetrics"         "/CollectDeviceMetrics"
    http_get "firmware/info"                "/firmware/info"
    http_get "process/list"                 "/process/list"
    local process_list_body="$HTTP_LAST_BODY"
    http_get "device/showPids"              "/device/showPids"
    http_get "discoverInterconnectTopology" "/discoverInterconnectTopology"
    http_get "deviceCount"                  "/deviceCount"
    http_get "vDeviceCount"                 "/vDeviceCount"
    local vdevice_count="$(grep -oE '[0-9]+' "$HTTP_LAST_BODY" 2>/dev/null | head -1 || true)"
    if ! is_integer "${vdevice_count:-}"; then
        vdevice_count="0"
    fi
    http_get "VDeviceInfos"                 "/VDeviceInfos"
    http_get "EncryptionVMStatus"           "/EncryptionVMStatus"
    http_get "EccBlocksInfo"                "/EccBlocksInfo"
    http_get "RSMIVersion"                  "/RSMIVersion"
    http_get "health/checkConfig"           "/health/checkConfig"
    http_get "deviceModelInfos"             "/deviceModelInfos"
    http_get "version/hy"                   "/version/hy"
    http_get "linktopo"                     "/linktopo"
    http_get "fields"                       "/fields"
    http_get "policies"                     "/policies"
    http_get "groups"                       "/groups"
    http_get "fieldgroups"                  "/fieldgroups"
    http_get "migInfos"                     "/migInfos" "MIG 依赖设备和驱动环境" environment
    http_get "GetHyLinkLinkStatus"          "/GetHyLinkLinkStatus" "依赖 HSL 驱动环境" environment
    http_get "GetHyUMCStatus"               "/GetHyUMCStatus" "依赖 UMC 驱动环境" environment

    if ((RUN_COMPAT)); then
        :
    else
        skip_test "HTTP" "compatible" "/compatible" "需要 cardModel、driverVersion、dtkVersion；使用 --compatibility 并提供三个参数后执行"
    fi
    skip_test "HTTP" "policies/result" "/policies/result" "需要 POST 请求体"
    skip_test "HTTP" "ProcessName(无pid)" "/ProcessName" "需要有效 PID；无 PID 请求没有兼容性意义"
    skip_test "HTTP" "processInfo(无pid)" "/processInfo" "需要有效 PID；无 PID 请求没有兼容性意义"
    http_get "processHCUInfo"               "/processHCUInfo"

    section "HTTP 按设备索引接口（dvInd=${DV}）"
    http_get "devicename"                   "/devicename/${DV}"
    http_get "DevSku"                       "/DevSku/${DV}"
    http_get "DevBrand"                     "/DevBrand/${DV}"
    http_get "DevVendorName"                "/DevVendorName/${DV}"
    http_get "DevVramVendor"                "/DevVramVendor/${DV}"
    http_get "DevPciBandwidth"              "/DevPciBandwidth/${DV}"
    http_get "MemoryPercent"                "/MemoryPercent/${DV}"
    http_get "PerfLevel"                    "/PerfLevel/${DV}"
    http_get "DevGpuMetricsInfo"            "/DevGpuMetricsInfo/${DV}"
    http_get "DeviceInfo"                   "/DeviceInfo/${DV}"
    http_get "PicbusInfo"                   "/PicbusInfo/${DV}"
    http_get "DFBandwidth"                  "/DFBandwidth/${DV}"
    http_get "HCUCuUsage"                   "/HCUCuUsage/${DV}"
    http_get "HCUSampledUsage"              "/HCUSampledUsage/${DV}"
    http_get "HCUCUSampledUsage"            "/HCUCUSampledUsage/${DV}"
    http_get "HCUWaveSampledUsage"          "/HCUWaveSampledUsage/${DV}"
    http_get "HCUSEUsage"                   "/HCUSEUsage/${DV}"
    http_get "FanSpeedInfo"                 "/FanSpeedInfo/${DV}"
    http_get "HCUUse"                       "/HCUUse/${DV}"
    http_get "Power"                        "/Power/${DV}"
    http_get "EccStatus(GFX)"               "/EccStatus/${DV}?block=GFX"
    http_get "Temperature"                  "/Temperature/${DV}"
    http_get "VbiosVersion"                 "/VbiosVersion/${DV}"
    http_get "DevFanRpms"                   "/DevFanRpms/${DV}"
    http_get "GetBus"                       "/GetBus/${DV}"
    http_get "hsl/linkstate"                "/hsl/linkstate/${DV}"
    http_get "throttle/flags"               "/throttle/flags/${DV}"
    http_get "memory/temperature"           "/memory/temperature/${DV}"
    http_get "board/temperature"            "/board/temperature/${DV}"
    http_get "memory/available"             "/memory/available/${DV}"
    http_get "clock/mem"                    "/clock/mem/${DV}"
    http_get "clock/gfx/max"               "/clock/gfx/max/${DV}"
    http_get "clock/mem/max"               "/clock/mem/max/${DV}"
    http_get "device/cunum"                 "/device/cunum/${DV}"
    http_get "device/nodeid"                "/device/nodeid/${DV}"
    http_get "device/memoverdrive"          "/device/memoverdrive/${DV}"
    http_get "device/info"                  "/device/info/${DV}"
    http_get "deviceRemainingInfo"          "/deviceRemainingInfo/${DV}"
    http_get "device/cuutil"                "/device/cuutil/${DV}"
    http_get "device/waveutil"              "/device/waveutil/${DV}"
    http_get "metricgroups/hcus"            "/metricgroups/hcus/${DV}"
    http_get "policies/hcus"                "/policies/hcus/${DV}"
    http_get "HyLinkStatusByDcuId"          "/HyLinkStatusByDcuId/${DV}" "依赖 HSL 设备状态" environment
    http_get "hswfw[sysfs]"                 "/hswfw/${DV}" "依赖 sysfs/debugfs" environment
    http_get "boost/GET[sysfs]"             "/boost/${DV}" "依赖 sysfs/debugfs" environment
    http_get "fdrhistory[sysfs]"            "/fdrhistory/${DV}" "依赖 sysfs/debugfs" environment
    http_get "memeccinfo[sysfs]"            "/memeccinfo/${DV}" "依赖 sysfs/debugfs" environment
    http_get "hslstatus[log]"               "/hslstatus/${DV}" "依赖 HSL 日志" environment
    skip_test "HTTP" "fields/hcus/values" "/fields/hcus/${DV}/values" "需要 POST 请求体 fieldIds"
    skip_test "HTTP" "PrintEventList" "/PrintEventList/${DV}" "接口可能长时间保持连接；不纳入默认冒烟测试"

    section "HTTP 条件接口"
    if ((vdevice_count > 0)); then
        http_get "VDeviceSingleInfo vDvInd=0"   "/VDeviceSingleInfo/0"
        http_get "vdevice/busypercent vDvInd=0" "/vdevice/busypercent/0"
    else
        skip_test "HTTP" "VDeviceSingleInfo" "/VDeviceSingleInfo/0" "当前没有虚拟设备"
        skip_test "HTTP" "vdevice/busypercent" "/vdevice/busypercent/0" "当前没有虚拟设备"
    fi

    if [[ -n "$device_count" && "$device_count" =~ ^[0-9]+$ && "$device_count" -gt 1 ]]; then
        http_get "topology/p2p 0->1" "/topology/p2p/${DV}/1"
    else
        skip_test "HTTP" "topology/p2p" "/topology/p2p/${DV}/1" "需要至少两个物理设备"
    fi

    http_get "fields/101 fieldId=101"    "/fields/101"
    skip_test "HTTP" "groups/0"             "/groups/0" "需要现有 group fixture；不创建或修改资源"
    skip_test "HTTP" "fieldgroups/1"         "/fieldgroups/1" "需要现有 field group fixture；不创建或修改资源"
    skip_test "HTTP" "diag/status/1"         "/diag/status/1" "需要现有诊断任务 ID；不使用虚构任务验证"

    local pids
    pids="$(grep -oE '"'[0-9]{1,}'"' "$process_list_body" 2>/dev/null | head -1 | tr -d '"' || true)"
    if ! is_integer "${pids:-}"; then
        pids="$(grep -h -oE '"'[0-9]{1,}'"' "${RUN_DIR}"/http_*.body 2>/dev/null | head -1 | tr -d '"' || true)"
    fi
    if [[ -n "${pids:-}" ]]; then
        http_get "process/gpuusage pid=${pids}" "/process/gpuusage/${pids}"
        http_get "process/name pid=${pids}"     "/process/name/${pids}"
    else
        skip_test "HTTP" "process/gpuusage" "/process/gpuusage/<pid>" "无法从 process/list 获取有效 PID"
        skip_test "HTTP" "process/name" "/process/name/<pid>" "无法从 process/list 获取有效 PID"
    fi
    skip_test "HTTP" "fields/entities/1/0/values" "/fields/entities/1/0/values" "需要 POST 请求体 fieldIds"
    skip_test "HTTP" "fields/groups/1/values"      "/fields/groups/1/values" "需要 POST 请求体 fieldIds"
}

run_dcgmi_tests() {
    section "dcgmi 只读命令"
    run_dcgmi "--help"                         --help
    run_dcgmi "discovery -l"                   discovery -l
    run_dcgmi "discovery -i a --hcuid ${DV}"  discovery -i a --hcuid "$DV"
    run_dcgmi "all-device-infos"               all-device-infos
    run_dcgmi "topo --hcuid ${DV}"             topo --hcuid "$DV"
    run_dcgmi "profile -l"                     profile -l
    run_dcgmi "process list"                   process list
    run_dcgmi "process info"                   process info
    run_dcgmi "health ${DV}"                   health "$DV"
    run_dcgmi "policy --get"                   policy --get
    run_dcgmi "group -l -j"                    group -l -j
    run_dcgmi "fieldgroup -l -j"               fieldgroup -l -j
    run_dcgmi "dmon -l"                        dmon -l

    if ((INCLUDE_DIAG)); then
        run_dcgmi "diag r ${DIAG_LEVEL}"       diag r "$DIAG_LEVEL"
    else
        skip_test "DCGMI" "diag r ${DIAG_LEVEL}" "dcgmi diag r ${DIAG_LEVEL}" "默认跳过诊断/压测；使用 --include-diag 显式启用"
    fi
}

write_summary() {
    local tested=$((total - skip))
    local pass_rate=0
    ((tested > 0)) && pass_rate=$((pass * 100 / tested))

    {
        printf '\n'
        printf '================================================================\n'
        printf ' HCU DCGM 接口兼容性验证完成\n'
        printf ' 时间       : %s\n' "$(date '+%Y-%m-%d %H:%M:%S %z')"
        printf ' 服务       : %s\n' "$BASE"
        printf ' 设备索引   : %s\n' "$DV"
        printf ' 总项目     : %d\n' "$total"
        printf ' 已执行     : %d\n' "$tested"
        printf ' PASS       : %d\n' "$pass"
        printf ' WARN       : %d\n' "$warn"
        printf ' FAIL       : %d\n' "$fail"
        printf ' SKIP       : %d\n' "$skip"
        printf ' 通过率     : %d%%（仅已执行项目）\n' "$pass_rate"
        printf ' 详细日志   : %s\n' "$LOG_FILE"
        printf ' 结果明细   : %s\n' "$RESULTS_FILE"
        printf '================================================================\n'
    } | tee -a "$LOG_FILE"
}

: > "$LOG_FILE"
printf 'status\tcategory\tlabel\ttarget\treason\n' > "$RESULTS_FILE"
{
    printf 'HCU DCGM compatibility test\n'
    printf 'script: %s\n' "$SCRIPT_NAME"
    printf 'started: %s\n' "$(date '+%Y-%m-%d %H:%M:%S %z')"
} >> "$LOG_FILE"

preflight

if ((RUN_HTTP || RUN_COMPAT)); then
    if command -v curl >/dev/null 2>&1; then
        if ((RUN_HTTP)); then
            run_http_tests
        else
            section "HTTP /compatible"
            http_get "compatible" "/compatible?cardModel=${COMPAT_CARD_MODEL}&driverVersion=${COMPAT_DRIVER_VERSION}&dtkVersion=${COMPAT_DTK_VERSION}" \
                "卡型号=${COMPAT_CARD_MODEL}，Driver=${COMPAT_DRIVER_VERSION}，DTK=${COMPAT_DTK_VERSION}"
        fi
    else
        skip_test "HTTP" "HTTP test suite" "$BASE" "curl 不可用"
    fi
fi

if ((RUN_DCGMI)); then
    if [[ -n "$DCGMI_BIN" && -x "$DCGMI_BIN" ]]; then
        run_dcgmi_tests
    else
        skip_test "DCGMI" "dcgmi test suite" "dcgmi" "dcgmi 不可用"
    fi
fi

write_summary

if ((fail > 0)); then
    exit 1
fi
if ((FAIL_ON_WARN && warn > 0)); then
    exit 1
fi
exit 0
