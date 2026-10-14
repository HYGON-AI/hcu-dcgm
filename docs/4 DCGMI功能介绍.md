# 4 DCGMI功能介绍

## 4.1 Diag

Diag提供HCU诊断与压力测试功能，支持单卡或多卡运行。可执行基础诊断（分级 1~4）、内存、PCIe、xHCL压力测试、稳定性压力测试（EDPp、GEMM）、显存完整性测试（memtestCL），也可按组或单独指定设备运行，评估硬件健康状态与性能稳定性。Diag提供HCU诊断与压力测试功能，支持单卡或多卡运行。可执行基础诊断（分级 1~4）、内存、PCIe、xHCL压力测试、稳定性压力测试（EDPp、GEMM）、显存完整性测试（memtestCL），也可按组或单独指定设备运行，评估硬件健康状态与性能稳定性。 当前命令默认测试全部HCU，也可通过-d/--device重复指定一个或多个HCU。

```text
[root@worker-200 hcu-dcgm]# --collect-timeout 200
bash: --collect-timeout: command not found...
[root@worker-200 hcu-dcgm]# dcgmi diag -h
Run diagnostics commands
```

```text
Usage:
dcgmi diag [flags]
dcgmi diag [command]
```

```text
Examples:
dcgmi diag -g <groupId> -i <flags>
dcgmi diag r <diagLevel>
dcgmi diag bandwidth <hcuId>
dcgmi diag memtestCL <hcuId>
dcgmi diag pcie
dcgmi diag xhcl
dcgmi diag edpp
dcgmi diag gemm
```

```text
Available Commands:
bandwidth   Run HCU bandwidth test
edpp        Run HCU EDPp stability test
gemm        Run GEMM target stress test
memtestCL   Run HCU memory stress test
pcie        Run PCIe memory bandwidth test
r           Run general diagnostics. Usage: dcgmi diag r <level> (1~4)
xhcl        Run XHCL stress test
```

```text
Flags:
-g, --group string   The group ID to query.
-h, --help           help for diag
-i, --info string    Specify which information to return.
b - memory bandwidth
m - memtestCL stress
```

```text
Use "dcgmi diag [command] --help" for more information about a command.
```

### 4.1.1 运行诊断测试

**执行命令**

```bash
dcgmi diag r 2
```

```text
[root@worker-200 hcu-dcgm]# dcgmi diag r 2
Successfully ran diagnostic for HCU.
+---------------------------+------------------------------------------------+
| Diagnostic                | Result                                         |
+===========================+================================================+
| Software Diagnostics      |                                                |
+---------------------------+------------------------------------------------+
| DTK Version Check         | pass                                           |
| Test Output               | 26.04                                          |
+------------------------------------------------+
| Driver Version Check      | pass                                           |
| Test Output               | 6.3.31-V1.5.1                                  |
+------------------------------------------------+
| RSMI Version Check        | pass                                           |
| Test Output               | RSMI Version: Major 7, Minor 6, Patch 0, Build 0 |
+------------------------------------------------+
| Device 0 VBIOS Version Check | pass                                           |
| Test Output               | VBIOS version for device 0: 5.714.001200p.604260 |
+------------------------------------------------+
| card K100_AI Compatibility Check | pass                                           |
| Test Output               | card K100_AI is compatible                     |
+------------------------------------------------+
| Device 1 VBIOS Version Check | pass                                           |
| Test Output               | VBIOS version for device 1: 5.714.001200p.604260 |
+------------------------------------------------+
| card K100_AI Compatibility Check | pass                                           |
| Test Output               | card K100_AI is compatible                     |
+------------------------------------------------+
| Memory Bandwidth Summary  | pass                                           |
| Test Output               | Parsed 2 entries                               |
+------------------------------------------------+
| Hardware Diagnostics      |                                                |
+---------------------------+------------------------------------------------+
| HCU: 0                     |------------------------------------------------|
| Memory Check              | Warning                                        |
| Error Message             | Issues detected. Memory usage: 0.00%. Block: UMC, State: ENABLED, CE: 1, UE: 0;  |
| Test Output               | Total: 63.98G, Used: 0.00G, Usage: 0.00% | ECC Status: Block: UMC, State: ENABLED, CE: 1, UE: 0;  |
+------------------------------------------------+
| picBus Check              | pass                                           |
| Test Output               | 0000:00:07.0                                   |
+------------------------------------------------+
| Power Check               | pass                                           |
| Test Output               | Max Power: 400, Current Power: 112             |
+------------------------------------------------+
| PCIe Bandwidth Check      | pass                                           |
| Test Output               | PCIe Bandwidth Test Results | Supported Transfer Rates: 2, Current Transfer Rate: 1 | Frequencies: [2500000000 16000000000 0 0 0] | PCIe Lanes: [16 16 0 0 0] |
+------------------------------------------------+
| Memory Bandwidth          | pass                                           |
| Test Output               | Bandwidth: 742.998 GB/s                        |
+------------------------------------------------+
| HCU: 1                     |------------------------------------------------|
| Memory Check              | pass                                           |
| Test Output               | Total: 63.98G, Used: 0.00G, Usage: 0.00% | ECC Status: All ECC blocks are healthy. |
+------------------------------------------------+
| picBus Check              | pass                                           |
| Test Output               | 0000:00:08.0                                   |
+------------------------------------------------+
| Power Check               | pass                                           |
| Test Output               | Max Power: 400, Current Power: 107             |
+------------------------------------------------+
| PCIe Bandwidth Check      | pass                                           |
| Test Output               | PCIe Bandwidth Test Results | Supported Transfer Rates: 2, Current Transfer Rate: 1 | Frequencies: [2500000000 16000000000 0 0 0] | PCIe Lanes: [16 16 0 0 0] |
+------------------------------------------------+
| Memory Bandwidth          | pass                                           |
| Test Output               | Bandwidth: 739.725 GB/s                        |
+------------------------------------------------+
+---------------------------+------------------------------------------------+
```

该命令中，r 2表示level2级别诊断，诊断测试一共分为4个水平。

### 4.1.2 测试指定设备组内HCU的显存带宽

**执行命令**

```bash
dcgmi diag -g 0 -i b
```

```text
[root@worker-200 hcu-dcgm]# dcgmi diag -g 0 -i b
Running memory bandwidth test for HCU(s): [0 1]
===== 带宽测试结果 =====
HCU0: 729.89 GB/s
HCU1: 744.19 GB/s
Successfully completed memory bandwidth test.
```

### 4.1.3 测试指定设备组内HCU的显存稳定性

```text
[root@localhost ~]# dcgmi diag -g 0 -i m
Running memtestCL stress test for HCU(s):
[HCU 0] ✅ 测试通过
------ Test Summary ------
Memtest86 walking 8-bit             : 0 failed iterations
True walking ones (32-bit)          : 0 failed iterations
Final                               : Final error count: 0 errors
True walking ones (8-bit)           : 0 failed iterations
True walking zeros (32-bit)         : 0 failed iterations
Memtest86 Modulo-20                 : 0 failed iterations
Integer logic                       : 0 failed iterations
Moving inversions (ones and zeros)  : 0 failed iterations
Random blocks                       : 0 failed iterations
True walking zeros (8-bit)          : 0 failed iterations
Moving inversions (random)          : 0 failed iterations
Integer logic (4 loops)             : 0 failed iterations
Integer logic (local memory)        : 0 failed iterations
Integer logic (4 loops, local memory) : 0 failed iterations
-------------Test end-------------
[HCU 1] ✅ 测试通过
------ Test Summary ------
True walking ones (8-bit)           : 0 failed iterations
True walking zeros (32-bit)         : 0 failed iterations
Integer logic                       : 0 failed iterations
Integer logic (4 loops)             : 0 failed iterations
True walking zeros (8-bit)          : 0 failed iterations
Memtest86 Modulo-20                 : 0 failed iterations
Final                               : Final error count: 0 errors
Moving inversions (ones and zeros)  : 0 failed iterations
Memtest86 walking 8-bit             : 0 failed iterations
Random blocks                       : 0 failed iterations
Integer logic (local memory)        : 0 failed iterations
Integer logic (4 loops, local memory) : 0 failed iterations
Moving inversions (random)          : 0 failed iterations
True walking ones (32-bit)          : 0 failed iterations
-----------Test end---------------
Successfully completed memtestCL stress test
```

### 4.1.4 测试指定HCU的显存带宽

**执行命令**

```bash
dcgmi diag bandwidth 0
```

```text
[root@localhost ~]# dcgmi diag bandwidth 0
Running memory bandwidth test for HCU: [0]
===== 带宽测试结果 =====
HCU0: 753.01 GB/s
memory bandwidth test: Done.
```

### 4.1.5 测试HCU的稳定性

**执行命令**

```bash
dcgmi diag edpp
或：dcgmi diag edpp -d 0
```

```text
[root@localhost ~]# dcgmi diag edpp
Running EDPp test for all HCUs...
===== Edpp Test Pattern: 10KHz start =====
----- Hardware Error Summary for pattern 10KHz -----
HCU 0 -> ECC: 0, Mem: 0, Compute: 0
HCU 1 -> ECC: 0, Mem: 0, Compute: 0

===== Edpp Test Pattern: 10KHz end =====

===== Edpp Test Pattern: 5KHz start =====
----- Hardware Error Summary for pattern 5KHz -----
HCU 0 -> ECC: 0, Mem: 0, Compute: 0
HCU 1 -> ECC: 0, Mem: 0, Compute: 0

===== Edpp Test Pattern: 5KHz end =====

===== Edpp Test Pattern: 4KHz start =====
----- Hardware Error Summary for pattern 4KHz -----
HCU 0 -> ECC: 0, Mem: 0, Compute: 0
HCU 1 -> ECC: 0, Mem: 0, Compute: 0
===== Edpp Test Pattern: 4KHz end =====

===== Edpp Test Pattern: 3KHz start =====
----- Hardware Error Summary for pattern 3KHz -----
HCU 0 -> ECC: 0, Mem: 0, Compute: 0
HCU 1 -> ECC: 0, Mem: 0, Compute: 0
===== Edpp Test Pattern: 3KHz end =====

===== Edpp Test Pattern: 2KHz start =====
----- Hardware Error Summary for pattern 2KHz -----
HCU 0 -> ECC: 0, Mem: 0, Compute: 0
HCU 1 -> ECC: 0, Mem: 0, Compute: 0

===== Edpp Test Pattern: 2KHz end =====

===== Edpp Test Pattern: 1.5KHz start =====
----- Hardware Error Summary for pattern 1.5KHz -----
HCU 0 -> ECC: 0, Mem: 0, Compute: 0
HCU 1 -> ECC: 0, Mem: 0, Compute: 0
===== Edpp Test Pattern: 1.5KHz end =====

===== Edpp Test Pattern: 1KHz start =====
----- Hardware Error Summary for pattern 1KHz -----
HCU 0 -> ECC: 0, Mem: 0, Compute: 0
HCU 1 -> ECC: 0, Mem: 0, Compute: 0
===== Edpp Test Pattern: 1KHz end =====

===== Edpp Test Pattern: 800Hz start =====
----- Hardware Error Summary for pattern 800Hz -----
HCU 0 -> ECC: 0, Mem: 0, Compute: 0
HCU 1 -> ECC: 0, Mem: 0, Compute: 0
===== Edpp Test Pattern: 800Hz end =====

===== Edpp Test Pattern: 500Hz start =====
----- Hardware Error Summary for pattern 500Hz -----
HCU 0 -> ECC: 0, Mem: 0, Compute: 0
HCU 1 -> ECC: 0, Mem: 0, Compute: 0
===== Edpp Test Pattern: 500Hz end =====

===== Edpp Test Pattern: 200Hz start =====
----- Hardware Error Summary for pattern 200Hz -----
HCU 0 -> ECC: 0, Mem: 0, Compute: 0
HCU 1 -> ECC: 0, Mem: 0, Compute: 0
===== Edpp Test Pattern: 200Hz end =====

===== Edpp Test Pattern: 100Hz start =====
----- Hardware Error Summary for pattern 100Hz -----
HCU 0 -> ECC: 0, Mem: 0, Compute: 0
HCU 1 -> ECC: 0, Mem: 0, Compute: 0
===== Edpp Test Pattern: 100Hz end =====

===== Edpp Test Pattern: 50Hz start =====
----- Hardware Error Summary for pattern 50Hz -----
HCU 0 -> ECC: 0, Mem: 0, Compute: 0
HCU 1 -> ECC: 0, Mem: 0, Compute: 0
===== Edpp Test Pattern: 50Hz end =====

===== Edpp Test Pattern: 100ms start =====
----- Hardware Error Summary for pattern 100ms -----
HCU 0 -> ECC: 0, Mem: 0, Compute: 0
HCU 1 -> ECC: 0, Mem: 0, Compute: 0
===== Edpp Test Pattern: 100ms end =====

===== EDPp Test Final Summary =====
[HCU 0]
Test Item                           Error Count
4KHz_ECC                            0
1.5KHz_ECC                          0
500Hz_ECC                           0
500Hz_Mem                           0
200Hz_Compute                       0
100ms_ECC                           0
100ms_Mem                           0
10KHz_ECC                           0
5KHz_ECC                            0
1.5KHz_Mem                          0
200Hz_ECC                           0
200Hz_Mem                           0
100Hz_ECC                           0
10KHz_Mem                           0
3KHz_ECC                            0
1KHz_Compute                        0
800Hz_Mem                           0
800Hz_Compute                       0
500Hz_Compute                       0
2KHz_ECC                            0
1KHz_ECC                            0
100Hz_Mem                           0
50Hz_Mem                            0
10KHz_Compute                       0
4KHz_Compute                        0
2KHz_Mem                            0
1KHz_Mem                            0
5KHz_Compute                        0
4KHz_Mem                            0
3KHz_Mem                            0
3KHz_Compute                        0
2KHz_Compute                        0
1.5KHz_Compute                      0
50Hz_ECC                            0
100Hz_Compute                       0
5KHz_Mem                            0
800Hz_ECC
```

### 4.1.6 测试指定HCU的显存稳定性

**执行命令**

```bash
dcgmi diag memtestCL 0
```

```text
[root@localhost ~]# dcgmi diag memtestCL 0
Running memtestCL stress test for HCU(s): [0]
[HCU 0] ✅ 测试通过
------ Test Summary ------
Integer logic                       : 0 failed iterations
Integer logic (4 loops, local memory) : 0 failed iterations
Memtest86 Modulo-20                 : 0 failed iterations
Integer logic (local memory)        : 0 failed iterations
Final                               : Final error count: 0 errors
True walking zeros (8-bit)          : 0 failed iterations
Moving inversions (random)          : 0 failed iterations
True walking ones (32-bit)          : 0 failed iterations
Random blocks                       : 0 failed iterations
Moving inversions (ones and zeros)  : 0 failed iterations
Memtest86 walking 8-bit             : 0 failed iterations
True walking ones (8-bit)           : 0 failed iterations
True walking zeros (32-bit)         : 0 failed iterations
Integer logic (4 loops)             : 0 failed iterations
--------------------------
memtestCL stress test: Done.
```

### 4.1.7 测试HCU的PCIe带宽

**执行命令**

```bash
dcgmi diag pcie
```

```text
[root@localhost ~]# dcgmi diag pcie
=== Running PCIe Bandwidth Test ===

  Unidirectional copy peak bandwidth GB/s (PCIe)

           device0     device1     device2
device0          N/A      29.031      29.032
device1       28.213         N/A      28.218
device2       28.212      28.214         N/A

HCU0 PCIe Sys->Fb: 29.031 GB/s, Fb->Sys: 28.213 GB/s
HCU1 PCIe Sys->Fb: 29.032 GB/s, Fb->Sys: 28.212 GB/s
PCIe bandwidth test: Done.
```

### 4.1.8 测试HCU的XHCL带宽

**执行命令**

```bash
dcgmi diag xhcl
```

```text
[root@localhost ~]# dcgmi diag xhcl
Running XHCL stress test on all HCU...

  Bidirectional copy peak bandwidth GB/s (XHCL)

           device0     device1     device2
device0          N/A      35.965      35.869
device1       35.965         N/A      38.379
device2       35.869      38.379         N/A

HCU0 <-> HCU1 XHCL: 38.379 GB/s
HCU1 <-> HCU0 XHCL: 38.379 GB/s
XHCL stress test: Done.
```

### 4.1.9 HCU的GEMM性能测试

**执行命令**

```bash
dcgmi diag gemm
```

```text
[root@localhost ~]# dcgmi diag gemm
Running GEMM target stress test on all HCU...
TargetStress: HCU0, hgemm mean 98.31
TargetStress: HCU0, bf16gemm mean 97.64
TargetStress: HCU0, sgemm mean 15.54
TargetStress: HCU0, tf32gemm mean 13.66
TargetStress: HCU0, i8gemm mean 148.39
TargetStress: HCU0, dgemm mean 0.21
TargetStress: HCU1, hgemm mean 97.98
TargetStress: HCU1, bf16gemm mean 98.32
TargetStress: HCU1, sgemm mean 15.55
TargetStress: HCU1, tf32gemm mean 13.70
TargetStress: HCU1, i8gemm mean 151.05
TargetStress: HCU1, dgemm mean 0.21
TargetStress completed!GEMM stress test: Done.
```

## 4.2 Discovery

Discovery用于查看当前节点可用的HCU设备，以及输出HCU的一些基本设备信息。

```text
[root@localhost ~]# dcgmi discovery -h discovery -- Used to discover and identify HCUs and their attributes.  Usage:Usage:   dcgmi discovery [flags]  Examples:Examples:   dcgmi discovery -l    dcgmi discovery -i <flags> --hcuid <HCUId>    dcgmi discovery -i <flags> -g <groupId> -v   dcgmi discovery -c  Flags:Flags:   -c, --compute-hierarchy   List all of the gpu instances and compute instances.       --hcuid string        The HCU ID to query.   -g, --group string        The group ID to query.   -h, --help                help for discovery   -i, --info string         Specify which information to return.                  a - device info                  p - power limits                  t - thermal limits                  c - clocks   -l, --list                List all the HCUs discovery on the host.   -v, --verbose             Display information per HCU.
```

### 4.2.1 查看节点全部HCU基本信息

**执行命令**

```bash
dcgmi discovery -l
```

```text
[root@localhost ~]# dcgmi discovery -l 2 HCUs found. +--------+----------------------+ | HCU ID | Device Information               | +--------+----------------------+ | 0      | Name: KONGMING - K100_AI                     | |        | PCI Bus ID: 0000:00:08.0                     | |        | Device Unique ID: 0x73873c7bc4ee4901                     | +--------+----------------------+ | 1      | Name: KONGMING - K100_AI                     | |        | PCI Bus ID: 0000:00:09.0                     | |        | Device Unique ID: 0x73873c7bc4d210a1                     | +--------+----------------------+
```

该命令中，可通过-l返回当前节点可用HCU的数量，并依次返回每个HCU设备的名称、PCI Bus ID、Unique ID等。

### 4.2.2 查看指定HCU的指定信息

**执行命令**

```bash
dcgmi discovery -i aptc --hcuid 0
```

```text
[root@localhost ~]# dcgmi discovery -i aptc --hcuid 0 +-------------+ | HCU ID: 0                | Device Information                  | +=============+ | Device Name              | KONGMING - K100_AI                  | | PCI Bus ID               | 0000:00:08.0            | | Unique ID                | 0x73873c7bc4ee4901                  | | Serial Number            | TPXS320016070601                    | | VBIOS                    | 5.714.001200p.604260                | +-------------+ | Power Ave Value (W)      | 333                     | | Power Max Value (W)      | 400                     | | Power Min Value (W)      | 0                       | +-------------+ | Current Temperature (C)  | 68                      | | ShutDown Temperature (C) | 110                     | | Critical Temperature (C) | 100                     | | Slowdown Temperature (C) | N/A                     | +-------------+ | Supported Clocks (MHz)   | MCLK:                   | |              | 875 *                   | |              | SCLK:                   | |              | 600 *                   | |              | 700                     | |              | 750                     | |              | 800                     | |              | 900                     | |              | 1000                    | |              | 1106                    | |              | 1200                    | |              | 1270                    | |              | 1319                    | |              | 1400                    | +-------------+
```

该命令中，可通过-i指定需要显示的信息，其中“a”表示设备基本信息；“p”表示设备功率信息；“t”表示设备温度信息；“c”表示设备时钟信息。可通过--hcuid指定HCU设备。在返回结果中，若当前节点无法正确获取某项信息，则将该项信息值以“N/A”显示；时钟信息中返回了MCLK和SCLK的支持列表，“*”表示当前值。

### 4.2.3 查看设备组内HCU的指定信息

**执行命令**

```bash
dcgmi discovery -i at -g 0
```

```text
[root@localhost ~]# dcgmi discovery -i at -g 0
+--------------------------+-------------------------------------------------+
| Group of 2 HCUs          | Device Information                              |
+==========================+=================================================+
| Device Name              | KONGMING - K100_AI                              |
| PCI Bus ID               | ****                                            |
| Unique ID                | ****                                            |
| Serial Number            | ****                                            |
| VBIOS                    | 5.714.001200p.604260                            |
+--------------------------+-------------------------------------------------+
| Current Temperature (C)  | ****                                            |
| ShutDown Temperature (C) | 110                                             |
| Critical Temperature (C) | 100                                             |
| Slowdown Temperature (C) | N/A                                             |
+--------------------------+-------------------------------------------------+
**** Non-homogenous settings across group. Use with –v flag to see details.
```

该命令中，可通过-g指定设备组。返回结果显示了指定设备组中HCU的信息，若HCU之间的某项信息值有差异，则将该项值以“****”显示。

**执行命令**

```bash
dcgmi discovery -i a -g 0 -v
```

```text
[root@localhost ~]# dcgmi discovery -i a -g 0 -v +-------------+ | HCU ID: 0                | Device Information                  | +=============+ | Device Name              | KONGMING - K100_AI                  | | PCI Bus ID               | 0000:00:08.0            | | Unique ID                | 0x73873c7bc4ee4901                  | | Serial Number            | TPXS320016070601                    | | VBIOS                    | 5.714.001200p.604260                | +-------------+ | HCU ID: 1                | Device Information                  | +=============+ | Device Name              | KONGMING - K100_AI                  | | PCI Bus ID               | 0000:00:09.0            | | Unique ID                | 0x73873c7bc4d210a1                  | | Serial Number            | TPXS320010040501                    | | VBIOS                    | 5.714.001200p.604260                | +-------------+
```

该命令中，可通过-v显示设备组中每个HCU的信息。

### 4.2.4 查看节点全部的MIG设备资源

**执行命令**

```bash
dcgmi discovery -c
```

```text
[root@localhost ~]# dcgmi discovery -c
+-------------------+--------------------------------------------------------------------+
| Instance Hierarchy                                                                     |
+===================+====================================================================+
| HCU 0             | HCU 0x73873c7bc4ee4901                                             |
| -> I 0/1          | GPU Instance created by Profile 1                                  |
|    -> CI 0/1/2    | Compute Instance (MIG 4g.31gb) created by Profile 1                |
+-------------------+--------------------------------------------------------------------+
| HCU 1             | HCU 0x73873c7bc4d210a1                                             |
| -> I 1/5          | GPU Instance created by Profile 3                                  |
|    -> CI 1/5/0    | Compute Instance (MIG 2g.15gb) created by Profile 0                |
+-------------------+--------------------------------------------------------------------+
```

该命令中，可通过-c查看当前节点中所有的GI（GPU Instance）和CI（Compute Instance）信息。除表头外，第一行表示GI、CI所归属的HCU，并显示了HCU的Unique ID；第二行表示GI信息，其中“0/1”表示其Instance ID为1，所属物理HCU的ID为0；第三行表示CI信息，其中“0/1/2”表示其Instance ID为2，所属GI ID为1，所属物理HCU的ID为0。

若当前节点MIG模式未开启，则返回空列表。

```text
[root@localhost ~]# hy-smi mig -L
Error: mig mode is not enabled.
[root@localhost ~]# dcgmi discovery -c
+-------------------+--------------------------------------------------------------------+
| Instance Hierarchy                                                                     |
+===================+====================================================================+
+-------------------+--------------------------------------------------------------------+
```

## 4.3 FileldGroup

FieldGroup用于创建和维护字段组。字段组是一些用于监控和统计的指标的集合。

```text
[root@localhost ~]# dcgmi fieldgroup -h
fieldgroup -- Used to create and maintain groups of field IDs. Groups of field
 IDs can then be uniformly controlled through other DCGMI subsystems.

Usage:
  dcgmi fieldgroup [flags]

Examples:
  dcgmi fieldgroup -l -j
  dcgmi fieldgroup -c <fieldGroupName> -f <fieldIds>
  dcgmi fieldgroup -i -g <fieldGroupId> -j
  dcgmi fieldgroup -d -g <fieldGroupId>

Flags:
  -c, --create string       Create a field group on the host.
  -d, --delete              Delete a field group on the host.
  -g, --fieldgroup string   Specify the field group ID to operate on.
  -f, --fieldids string     Comma-separated list of the field ids to add to a field group when creating a new one.
  -h, --help                help for fieldgroup
  -i, --info                Display the information for the specified field group ID.
  -j, --json                Print the output in a json format.
  -l, --list                List the field groups that currently exist for the host.
```

### 4.3.1 查看全部字段组信息

**执行命令**

```bash
dcgmi fieldgroup -l
```

```text
[root@localhost ~]# dcgmi fieldgroup -l
+-------------------+----------------------------------------------------------+
| FIELD GROUPS                                                                 |
| 2  field groups found.                                                       |
+===================+==========================================================+
| Field Groups      |                                                          |
| -> 1              |                                                          |
|    -> ID          | 1                                                        |
|    -> Name        | fg0                                                      |
|    -> Field IDs   | 1003, 1004                                               |
| -> 2              |                                                          |
|    -> ID          | 2                                                        |
|    -> Name        | fg1                                                      |
|    -> Field IDs   | 1001, 1002                                               |
+-------------------+----------------------------------------------------------+
```

该命令中，可通过-l查看当前节点所有的字段组列表，列表信息包括组ID、名称、包含的字段ID等。

### 4.3.2 创建字段组

**执行命令**

```bash
dcgmi fieldgroup -c fg0 -f 1003,1004
```

```text
[root@localhost ~]# dcgmi fieldgroup -c fg0 -f 1003,1004
Successfully created field group "fg0" with a field group ID of 1
[root@localhost ~]# dcgmi fieldgroup -l
+-------------------+----------------------------------------------------------+
| FIELD GROUPS                                                                 |
| 1 field group found.                                                         |
+===================+==========================================================+
| Field Groups      |                                                          |
| -> 1              |                                                          |
|    -> ID          | 1                                                        |
|    -> Name        | fg0                                                      |
|    -> Field IDs   | 1003, 1004                                               |
+-------------------+----------------------------------------------------------+
```

该命令中，可通过-c指定待创建的字段组的名称；通过-f指定待添加到字段组中的字段ID。若添加多个字段，需用逗号分隔。

### 4.3.3 查看指定字段组信息

**执行命令**

```bash
dcgmi fieldgroup -i -g 1
```

```text
[root@localhost ~]# dcgmi fieldgroup -i -g 1
+-------------------+----------------------------------------------------------+
| FIELD GROUP INFO                                                             |
+===================+==========================================================+
| -> 1              |                                                          |
|    -> ID          | 1                                                        |
|    -> Name        | fg0                                                      |
|    -> Field IDs   | 1003, 1004                                               |
+-------------------+----------------------------------------------------------+
```

该命令中，通过-g指定字段组，并通过-i获取该字段组信息。

### 4.3.4 删除指定字段组

**执行命令**

```bash
dcgmi fieldgroup -d -g 1
```

```text
[root@localhost ~]# dcgmi fieldgroup -d -g 1
Successfully removed field group 1
[root@localhost ~]# dcgmi fieldgroup -i -g 1
Failed to get field group info: field group 1 does not exist
```

该命令中，通过-g指定字段组，并通过-d删除该字段组。

## 4.4 Group

Group用于创建和管理HCU设备组，从而便于对HCU设备进行其他操作。

```text
[root@localhost ~]# dcgmi group -h
group -- Used to create and maintain groups of HCUs. Groups of HCUs can then be
uniformly controlled through other DCGMI subsystems.

Usage:
  dcgmi group [flags]

Examples:
  dcgmi group -l -j
  dcgmi group -c <groupName> --default
  dcgmi group -c <groupName> -a <entityId>
  dcgmi group -d <groupId>
  dcgmi group -g <groupId> -i -j
  dcgmi group -g <groupId> -a <entityId>
  dcgmi group -g <groupId> -r <entityId>

Flags:
  -a, --add string      Add device(s) to group. (csv HCUIds like 0,1 or entityIds like hcu:0,hcu:1)
  -c, --create string   Create a group on the host.
      --default         Adds all available HCUs to the group being created.
  -d, --delete string   Delete a group on the host.
  -g, --group string    Specify the group ID to operate on.
  -h, --help            help for group
  -i, --info            Display the information for the specified group ID.
  -j, --json            Print the output in a json format.
  -l, --list            List the groups that currently exist for the host.
  -r, --remove string   Remove device(s) from group. (csv hcuIds like 0,1 or entityIds like hcu:0,hcu:1)
```

### 4.4.1 查看所有设备组信息

**执行命令**

```bash
dcgmi group -l
```

```text
[root@localhost ~]# dcgmi group -l
+-------------------+----------------------------------------------------------+
| GROUPS                                                                       |
| 2  groups found.                                                             |
+===================+==========================================================+
| Groups            |                                                          |
| -> 0              |                                                          |
|    -> Group ID    | 0                                                        |
|    -> Group Name  | group0                                                   |
|    -> Entities    | HCU 0, HCU 1                                             |
| -> 1              |                                                          |
|    -> Group ID    | 1                                                        |
|    -> Group Name  | group1                                                   |
|    -> Entities    | vHCU 0, vHCU 1                                           |
+-------------------+----------------------------------------------------------+
```

该命令中，可通过-l查看当前节点所有的设备组列表，列表信息包括组ID、名称、包含的设备等。

### 4.4.2 创建设备组

**执行命令**

```bash
dcgmi group -c dg0 -a 0,1
```

```text
[root@localhost ~]# dcgmi group -c dg0 -a 0,1
Successfully created group "dg0" with a group ID of 2
Successfully added entities to group 2
[root@localhost ~]# dcgmi group -g 2 -i
+-------------------+----------------------------------------------------------+
| GROUP INFO                                                                   |
+===================+==========================================================+
| -> 2              |                                                          |
|    -> Group ID    | 2                                                        |
|    -> Group Name  | dg0                                                      |
|    -> Entities    | HCU 0, HCU 1                                             |
+-------------------+----------------------------------------------------------+
```

该命令中，可通过-c指定待创建的设备组的名称；通过-a指定待添加到设备组中的设备ID。若添加多个设备，需用逗号分隔。若使用--default，则会将节点所有可用HCU添加到组中。

### 4.4.3 查看指定设备组信息

**执行命令**

```bash
dcgmi group -g 0 -i
```

```text
[root@localhost ~]# dcgmi group -g 0 -i
+-------------------+----------------------------------------------------------+
| GROUP INFO                                                                   |
+===================+==========================================================+
| -> 0              |                                                          |
|    -> Group ID    | 0                                                        |
|    -> Group Name  | group0                                                   |
|    -> Entities    | HCU 0, HCU 1                                             |
+-------------------+----------------------------------------------------------+
```

该命令中，通过-g指定设备组，并通过-i获取该设备组信息。

### 4.4.4 向设备组中添加设备

**执行命令**

```bash
dcgmi group -g 2 -a 0
```

```text
[root@localhost ~]# dcgmi group -c dg1
Successfully created group "dg1" with a group ID of 3
[root@localhost ~]# dcgmi group -g 3 -a 0
Successfully added entities to group 3
[root@localhost ~]# dcgmi group -g 3 -i
+-------------------+----------------------------------------------------------+
| GROUP INFO                                                                   |
+===================+==========================================================+
| -> 3              |                                                          |
|    -> Group ID    | 3                                                        |
|    -> Group Name  | dg1                                                      |
|    -> Entities    | HCU 0                                                    |
+-------------------+----------------------------------------------------------+
```

该命令中，通过-g指定设备组，并通过-a指定待添加到设备组中的设备ID。若添加多个设备，需用逗号分隔。

### 4.4.5 从设备组中移除设备

**执行命令**

```bash
dcgmi group -g 2 -r 0
```

```text
[root@localhost ~]# dcgmi group -g 2 -r 0
Successfully removed entities from group 2
[root@localhost ~]# dcgmi group -g 2 -i
+-------------------+----------------------------------------------------------+
| GROUP INFO                                                                   |
+===================+==========================================================+
| -> 2              |                                                          |
|    -> Group ID    | 2                                                        |
|    -> Group Name  | dg0                                                      |
|    -> Entities    | HCU 1                                                    |
+-------------------+----------------------------------------------------------+
```

该命令中，通过-g指定设备组，并通过-r从指定设备组中待移除的设备ID。若移除多个设备，需用逗号分隔。

### 4.4.6 删除指定设备组

**执行命令**

```bash
dcgmi group -d 2
```

```text
[root@localhost ~]# dcgmi group -d 2
Successfully removed group 2
[root@localhost ~]# dcgmi group -g 2 -i
Failed to get group info: group 2 does not exist
```

该命令中，通过-d指定待删除的设备组。

## 4.5 Health

用于监控检查HCU基本健康状态。

### 4.5.1 查看全部HCU健康信息

**执行命令**

```bash
dcgmi health
```

```text
[root@localhost ~]# dcgmi health
  DEVICE ID|  OVERALL STATUS|            CHECK TYPE|   STATUS|DETAILS
          0|         Warning|   NumaTopology Health|  Healthy|DeviceID:0 NumaNode:0 Affinity:-1
          0|         Warning|  PcieBandwidth Health|  Healthy|Sent:0 Recv:0 BW:0.00MB/s
          0|         Warning|          Power Health|  Healthy|Current Power: 84, Safe Range: [10, 400]
          0|         Warning|         Memory Health|  Healthy|Memory usage 0.00% within safe threshold 90.00%
          0|         Warning|    Temperature Health|  Healthy|Temperature 52.00°C within safe threshold 90.00°C
          0|         Warning|    Performance Health|  Healthy|PerfLevel: AUTO
          0|         Warning|      EccBlocks Health|  Warning|UMC CE:1 UE:0
          0|         Warning|       HCUUsage Health|  Healthy|HCU Usage: 0.00%
+----------------------------------------------------------------+
  1|  Healthy|   NumaTopology Health|  Healthy|DeviceID:1 NumaNode:0 Affinity:-1
  1|  Healthy|  PcieBandwidth Health|  Healthy|Sent:0 Recv:0 BW:0.00MB/s
  1|  Healthy|          Power Health|  Healthy|Current Power: 80, Safe Range: [10, 400]
  1|  Healthy|         Memory Health|  Healthy|Memory usage 0.00% within safe threshold 90.00%
  1|  Healthy|    Temperature Health|  Healthy|Temperature 56.00°C within safe threshold 90.00°C
  1|  Healthy|    Performance Health|  Healthy|PerfLevel: AUTO
  1|  Healthy|      EccBlocks Health|  Healthy|All ECC blocks are healthy
  1|  Healthy|       HCUUsage Health|  Healthy|HCU Usage: 0.00%
+----------------------------------------------------------------+
```

### 4.5.2 查看指定HCU健康信息

**执行命令**

```bash
dcgmi health 0
```

```text
[root@localhost ~]# dcgmi health 0
  DEVICE ID|  OVERALL STATUS|            CHECK TYPE|   STATUS|DETAILS
          0|         Warning|   NumaTopology Health|  Healthy|DeviceID:0 NumaNode:0 Affinity:-1
          0|         Warning|  PcieBandwidth Health|  Healthy|Sent:0 Recv:0 BW:0.00MB/s
          0|         Warning|          Power Health|  Healthy|Current Power: 84, Safe Range: [10, 400]
          0|         Warning|         Memory Health|  Healthy|Memory usage 0.00% within safe threshold 90.00%
          0|         Warning|    Temperature Health|  Healthy|Temperature 52.00°C within safe threshold 90.00°C
          0|         Warning|    Performance Health|  Healthy|PerfLevel: AUTO
          0|         Warning|      EccBlocks Health|  Warning|UMC CE:1 UE:0
          0|         Warning|       HCUUsage Health|  Healthy|HCU Usage: 0.00%
+----------------------------------------------------------------+
```

## 4.6 Policy

Policy用于定义在特定事件发生时触发的操作，可以实现异常的自动处理。

```text
[root@localhost ~]# dcgmi policy -h
policy -- Used to control policies for groups of HCUs. Policies control actions
 which are triggered by specific events.

Usage:
  dcgmi policy [flags]

Examples:
  dcgmi policy --get -j
  dcgmi policy -g <groupId> --reg
  dcgmi policy -g <groupId> --set <actn,val> -M <maxRetiredPages> -T
        <maxTemp> -P <maxPower> -e -p
  dcgmi policy -g <groupId> --clear

Flags:
      --clear             Clear the current violation policy.
  -e, --eccerrors         Add ECC errors to the policy
                          conditions.
      --get               Get the current violation policy.
  -g, --group string      The group ID to query.
  -h, --help              help for policy
  -j, --json              Print the output in a json format.
  -M, --maxpages string   Specify the maximum number of retired pages that
                          will trigger a violation.
  -P, --maxpower string   Specify the maximum power a group's HCUs can reach
                          before triggering a violation.
  -T, --maxtemp string    Specify the maximum temperature a group's HCUs can
                          reach before triggering a violation.
  -p, --pcierrors         Add PCIe replay errors to the policy conditions.
      --reg               Register this process for policy updates.  This
                          process will sit in an infinite loop waiting for
                           updates from the policy manager.
      --set string        Set the current violation policy. Use csv action
                          ,validation (ie. 1,2)
                          -----
                          Action to take when any of the violations
                           specified occur.
                          0 - None
                          1 - HCU Reset
                          -----
                          Validation to take after the violation action has
                          been performed.
                          0 - None
                          1 - System Validation (short)
                          2 - System Validation (medium)
                          3 - System Validation (long)
  -v, --verbose           Display policy information per HCU.
```

### 4.6.1 查看已设定的策略

**执行命令**

```bash
dcgmi policy -g 0 --get
```

```text
[root@localhost ~]# dcgmi policy -g 0 --get
Policy information
+-----------------------------+------------------------------------------------+
| Policy Information                                                           |
| group0                                                                       |
+=============================+================================================+
| Violation conditions        | ECC errors                                     |
|                             | PCI errors and replays                         |
|                             | Max retired pages threshold - 80               |
|                             | Max temperature threshold - 100                |
|                             | Max power threshold - 150                      |
| Action on violation         | None                                           |
| Validation after action     | System Validation (Short)                      |
+-----------------------------+------------------------------------------------+
**** Non-homogenous settings across group. Use with –v flag to see details.
```

该命令中，可通过-g指定设备组，并通过--get获取指定设备组中的HCU的策略。若没有通过-g指定设备组，则--get获取节点所有HCU上的策略。此外，若HCU之间的某项设定值有差异，则将该项值以“****”显示。

**执行命令**

```bash
dcgmi policy -g 0 --get -v
```

```text
[root@localhost ~]# dcgmi policy -g 0 --get -v Policy information +------------+ | Policy Information                       | | HCU ID: 0                    | +============+ | Violation conditions        | ECC errors             | |                 | PCI errors and replays             | |                 | Max retired pages threshold - 80               | |                 | Max temperature threshold - 100                | |                 | Max power threshold - 150                      | | Action on violation         | None                   | | Validation after action     | System Validation (Short)                      | +------------+ | Policy Information                       | | HCU ID: 1                    | +============+ | Violation conditions        | ECC errors             | |                 | PCI errors and replays             | |                 | Max retired pages threshold - 80               | |                 | Max temperature threshold - 100                | |                 | Max power threshold - 150                      | | Action on violation         | None                   | | Validation after action     | System Validation (Short)                      | +------------+
```

该命令中，可通过-v显示设备组中每个HCU的策略。

### 4.6.2 设定新策略

**执行命令**

```bash
dcgmi policy -g 0 --set 0,1 -T 100 -P 150 -M 80 -p -e
```

```text
[root@localhost ~]# dcgmi policy -g 0 --set 0,1 -T 100 -P 150 -M 80 -p -e
Policy successfully set.
[root@localhost ~]# dcgmi policy -g 0 --get
Policy information
+-----------------------------+------------------------------------------------+
| Policy Information                                                           |
| group0                                                                       |
+=============================+================================================+
| Violation conditions        | ECC errors                                     |
|                             | PCI errors and replays                         |
|                             | Max retired pages threshold - 80               |
|                             | Max temperature threshold - 100                |
|                             | Max power threshold - 150                      |
| Action on violation         | None                                           |
| Validation after action     | System Validation (Short)                      |
+-----------------------------+------------------------------------------------+
**** Non-homogenous settings across group. Use with –v flag to see details.
```

该命令中，可通过-g指定设备组，并通过--set设定指定设备组中的HCU的策略。--set值为逗号分隔的两个数字，分别代表监控的条件发生时所触发的操作和操作执行后的验证（验证在当前版本不可用）。若没有通过-g指定设备组，则--set为节点所有HCU设定策略。可通过-T设定最大温度值条件；可通过-P设定最大功率值条件；可通过-M设定最大显存无效页面条件；可通过-p启用检测PCIe errors条件；可通过-e启用检测ECC errors条件。

### 4.6.3 监控策略条件

**执行命令**

```bash
dcgmi policy -g 0 --reg
```

```text
[root@localhost ~]# dcgmi policy -g 0 --reg
Listening for violations.
Timestamp: Wed Dec 31 14:35:04 2025
HCU 0: The maximum power limit has violated policy manager values.
Power: 400
Action will be performed on HCU 0 in 10 seconds.
No action required.
- Ctrl+C pressed in Terminal
```

该命令中，可通过-g指定设备组，并通过--reg监控指定设备组中的HCU的策略条件。若没有通过-g指定设备组，则--reg监控所有HCU的策略条件。可通过Ctrl+C停止监控。

### 4.6.4 清除已设定的策略

**执行命令**

```bash
dcgmi policy -g 0 --clear
```

```text
[root@localhost ~]# dcgmi policy -g 0 --clear
Policy successfully set.
[root@localhost ~]# dcgmi policy -g 0 --get
Policy information
+-----------------------------+------------------------------------------------+
| Policy Information                                                           |
| group0                                                                       |
+=============================+================================================+
| Violation conditions        | None                                           |
| Action on violation         | None                                           |
| Validation after action     | None                                           |
+-----------------------------+------------------------------------------------+
**** Non-homogenous settings across group. Use with –v flag to see details.
```

该命令中，可通过-g指定设备组，并通过--clear清除指定设备组中的HCU已设定的策略。若没有通过-g指定设备组，则--clear清除所有HCU已设定的策略。

## 4.8 Topo

Topo用于查看HCU设备拓扑信息。

```text
[root@localhost ~]# dcgmi topo -h  policy -- Used to find the topology of HCUs on the system.  Usage:Usage:   dcgmi topo [flags]  Examples:Examples:   dcgmi topo -g <groupId>   dcgmi topo --hcuid <hcuId>  Flags:Flags:       --hcuid string   The HCU ID to query.   -g, --group string   The group ID to query.   -h, --help           help for topo
```

### 4.8.1 查看全部HCU拓扑信息

**执行命令**

```bash
dcgmi topo
```

```text
[root@localhost ~]# dcgmi topo
+-------------------+------------------------+
| Topology Information | DCGM_ALL_SUPPORTED_HCUS |
+-------------------+------------------------+
| Numa Affinity     | -1                     |
| CPU Core Affinity | 0-15                   |
| NUMA Optimal      | True                   |
| Worst Path        | Connected via PCIE     |
+-------------------+------------------------+
```

该命令会返回节点所有HCU对应的Numa Affinity、所在Numa节点上的CPU Core、是否位于同一Numa节点、最差通信路径等信息。

### 4.8.2 查看指定设备组内HCU拓扑信息

**执行命令**

```bash
dcgmi topo -g 0
```

```text
[root@localhost ~]# dcgmi topo -g 0
+-------------------+--------------------+
| Topology Information                  |
| group0                                |
+===================+====================+
| Numa Affinity     | -1                 |
| CPU Core Affinity | 0-15               |
| NUMA Optimal      | True               |
| Worst Path        | Connected via PCIE |
+-------------------+--------------------+
```

该命令中，可通过-g指定设备组，从而获取指定设备组中的HCU对应的Numa Affinity、所在Numa节点上的CPU Core、是否位于统一Numa节点、最差通信路径等信息。

### 4.8.3 查看指定HCU拓扑信息

**执行命令**

```bash
dcgmi topo --hcuid 0
```

```text
[root@localhost ~]# dcgmi topo --hcuid 0
+-------------------+--------------------+
| Topology Information                   |
| HCU ID: 0                              |
+===================+====================+
| Numa Node         | 0                  |
| Numa Affinity     | -1                 |
| CPU Core Affinity | 0-15               |
| To HCU 1          | Connected via PCIE |
+-------------------+--------------------+
```

该命令中，可通过--hcuid指定HCU，从而获取该HCU对应的Numa节点、Numa Affinity、所在Numa节点上的CPU Core、同节点其他HCU的通信方式等信息。

## 4.9 Dmon

Dmon用于实时查看统计监控信息。

```text
[root@localhost ~]# dcgmi dmon -h
dmon -- Used to monitor HCUs and their stats.

Usage:
  dcgmi dmon [flags]

Examples:
  dcgmi dmon -l
  dcgmi dmon -e <fieldId>
  dcgmi dmon -f <fieldGroupId>
  dcgmi dmon -e <fieldId> -i <entityId>
  dcgmi dmon -e <fieldId> -g <groupId>
  dcgmi dmon -f <fieldGroupId> -i <entityId>
  dcgmi dmon -f <fieldGroupId> -g <groupId>
  dcgmi dmon -e <fieldId> -i <entityId> -d <delaySeconds> -c <count>

Flags:
  -c, --count int             Integer representing How many times to loop
                              before exiting. [default- runs forever.] (default -1)
  -d, --delay int             Time in seconds between each run. (default 3)
  -i, --entityid string       Comma-separated list of entities to run the dmon on.
                              (csv hcuIds like 0,1 or entityIds like hcu:0,hcu:1).
                              Default is -1 which runs for all supported HCU.
  -f, --fieldgroupid string   The field group to query on the host.
  -e, --fieldid string        Field identifier to view/inject.
  -g, --groupid string        The group to query on the host.
  -h, --help                  help for dmon
  -l, --list                  List to look up the field names and field ids.
```

### 4.9.1 查看全部监控字段信息

**执行命令**

```bash
dcgmi dmon -l
```

```text
[root@localhost ~]# dcgmi dmon -l
+-------------------+----------------------------------------------------------+
| Field ID          | Field Name                                               |
+-------------------+----------------------------------------------------------+
| 101               | HCU_TEMP                                                 |
| 102               | HCU_POWER_USAGE                                          |
| 103               | HCU_POWER_CAP                                            |
| 104               | HCU_UTILIZATION_RATE                                     |
| 105               | HCU_SCLK                                                 |
| 106               | HCU_COMPUTE_UNIT_COUNT                                   |
| 107               | HCU_COMPUTE_UNIT_REMAINING_COUNT                         |
| 108               | HCU_USED_MEMORY_BYTES                                    |
| 109               | HCU_MEMORY_CAP_BYTES                                     |
| 110               | HCU_MEMORY_REMAINING                                     |
| 111               | HCU_MEMORY_PERCENT                                       |
| 112               | HCU_PCIE_BW_MB                                           |
| 113               | HCU_PCIE_RECEIVE_MB                                      |
| 114               | HCU_PCIE_SENT_MB                                         |
| 115               | HCU_DF_BW_READ                                           |
| 116               | HCU_DF_BW_WRITE                                          |
| 117               | HCU_DF_BW_READ_WRITE                                     |
| 118               | HCU_VHCU_COUNT                                           |
| 201               | VHCU_SCLK                                                |
| 202               | VHCU_TEMP                                                |
| 203               | VHCU_UTILIZATION_RATE                                    |
| 204               | VHCU_USED_MEMORY_BYTES                                   |
+-------------------+----------------------------------------------------------+
```

该命令会返回节点HCU可用的监控字段名称和ID信息。

### 4.9.2 查看指定字段的统计信息

**执行命令**

```bash
dcgmi dmon -e 101,102
```

```text
[root@localhost ~]# dcgmi dmon -e 101,102
#Entity ID             HCU_TEMP  HCU_POWER_USAGE
HCU 0                  51.000    86.000
HCU 1                  55.000    82.000
HCU 0                  50.000    86.000
HCU 1                  55.000    82.000
^C
dmon was stopped do to receiving a signal.
```

该命令中，可通过-e指定字段ID。若需指定多个字段ID，需用逗号分隔。未指定设备或设备组的情况下，会返回节点上所有HCU的指定字段值。Ctrl + C可停止信息的输出。

### 4.9.3 查看指定字段组的统计信息

**执行命令**

```bash
dcgmi dmon -f 0
```

```text
[root@localhost ~]# dcgmi dmon -f 0
#Entity ID             HCU_TEMP  HCU_PCIE_BW_MB  HCU_POWER_USAGE
HCU 0                  52.000    3.324           114.000
HCU 1                  56.000    0.001           79.000
HCU 0                  51.000    0.000           84.000
HCU 1                  55.000    0.000           83.000
^C
dmon was stopped do to receiving a signal.
```

该命令中，可通过-f指定字段组ID，字段组ID及其包含字段可通过dcgmi fieldgroup -l查看。未指定设备或设备组的情况下，会返回节点上所有HCU的指定字段值。Ctrl + C可停止信息的输出。

### 4.9.4 查看指定字段及设备的统计信息

**执行命令**

```bash
dcgmi dmon -e 101,102 -i 0
```

```text
[root@localhost ~]# dcgmi dmon -e 101,102 -i 0
#Entity ID             HCU_TEMP  HCU_POWER_USAGE
HCU 0                  50.000    86.000
HCU 0                  51.000    114.000
^C
dmon was stopped do to receiving a signal.
```

该命令中，可通过-e指定字段ID。若需指定多个字段ID，需用逗号分隔。通过-i指定设备，默认为HCU，也可通过<设备类型:设备序列>的方式指定，多个设备之间通过逗号分隔（例如vhcu:0,vhcu:1）。Ctrl + C可停止信息的输出。

### 4.9.5 查看指定字段及设备组的统计信息

**执行命令**

```bash
dcgmi dmon -e 101,102 -g 0
```

```text
[root@localhost ~]# dcgmi dmon -e 101,102 -g 0
#Entity ID             HCU_TEMP  HCU_POWER_USAGE
HCU 0                  65.000    244.000
HCU 1                  55.000    84.000
HCU 0                  66.000    239.000
HCU 1                  55.000    84.000
^C
dmon was stopped do to receiving a signal.
```

该命令中，可通过-e指定字段ID；通过-g指定设备组。若需指定多个字段ID，需用逗号分隔。设备组及其包含设备可通过dcgmi group -l查看。Ctrl + C可停止信息的输出。

### 4.9.6 查看指定字段组及设备的统计信息

**执行命令**

```bash
dcgmi dmon -f 0 -i 0,1
```

```text
[root@localhost ~]# dcgmi dmon -f 0 -i 0,1
#Entity ID             HCU_TEMP  HCU_PCIE_BW_MB  HCU_POWER_USAGE
HCU 0                  69.000    0.000           238.000
HCU 1                  54.000    0.000           84.000
HCU 0                  70.000    0.040           242.000
HCU 1                  55.000    0.000           83.000
^C
dmon was stopped do to receiving a signal.
```

该命令中，可通过-f指定字段组ID，字段组ID及其包含字段可通过dcgmi fieldgroup -l查看。通过-i指定设备，默认为HCU，也可通过<设备类型:设备序列>的方式指定，多个设备之间通过逗号分隔（例如vhcu:0,vhcu:1）。Ctrl + C可停止信息的输出。

### 4.9.7 查看指定字段组及设备组的统计信息

**执行命令**

```bash
dcgmi dmon -f 0 -g 0
```

```text
[root@localhost ~]# dcgmi dmon -f 0 -g 0
#Entity ID             HCU_TEMP  HCU_PCIE_BW_MB  HCU_POWER_USAGE
HCU 0                  71.000    0.034           245.000
HCU 1                  54.000    0.002           84.000
HCU 0                  71.000    0.033           242.000
HCU 1                  54.000    0.000           83.000
^C
dmon was stopped do to receiving a signal.
```

该命令中，可通过-f指定字段组ID，字段组ID及其包含字段可通过dcgmi fieldgroup -l查看。通过-g指定设备组，设备组及其包含设备可通过dcgmi group -l查看。Ctrl + C可停止信息的输出。

### 4.9.8 查看自定义频次统计信息

**执行命令**

```bash
dcgmi dmon -e 101,102 -i 0,1 -d 5 -c 3
```

```text
[root@localhost ~]# dcgmi dmon -e 101,102 -i 0,1 -d 5 -c 3
#Entity ID             HCU_TEMP  HCU_POWER_USAGE
HCU 0                  70.000    241.000
HCU 1                  53.000    82.000
HCU 0                  71.000    243.000
HCU 1                  53.000    83.000
HCU 0                  71.000    245.000
HCU 1                  52.000    83.000
```

该命令中，可通过-d指定每次输出的时间间隔（单位为秒）。通过-c指定输出次数。

## 4.10 Profile

Profile用于性能指标字段的查看，可配合Dmon获取字段值。

```text
[root@localhost ~]# dcgmi profile -h
profile -- View available profiling metrics for HCUs

Usage:
  dcgmi profile [flags]

Examples:
  dcgmi profile -l
  dcgmi profile -l -i <entityId>
  dcgmi profile -l -g <groupId>

Flags:
  -i, --entityId string   Comma-seperated list of entity IDs to query.
                          Default is supported HCUs on the system. Run
                          dcgmi discovery -l to check list of HCUs available
  -g, --group string      The group of HCUs to query on the host.
  -h, --help              help for profile
  -l, --list              List available profiling metric
```

### 4.10.1 查看性能指标信息

**执行命令**

```bash
dcgmi profile -l
```

```text
[root@localhost ~]# dcgmi profile -l
+----------------+----------+------------------------------------------------------+
| Group.Subgroup | Field ID | Field Tag                                            |
+----------------+----------+------------------------------------------------------+
| A.1            | 104      | utilization_rate                                     |
| A.1            | 105      | sclk                                                 |
| B.1            | 111      | memory_percent                                       |
| C.1            | 112      | pcie_bw_mb                                           |
| C.1            | 114      | pcie_sent_mb                                         |
| C.1            | 113      | pcie_receive_mb                                      |
+----------------+----------+------------------------------------------------------+
```

该命令中，可通过-l查看性能指标ID、归属组、名称等信息。

### 4.10.2 查看指定设备可用指标

**执行命令**

```bash
dcgmi profile -l -i 0
```

```text
[root@localhost ~]# dcgmi profile -l -i 0
+----------------+----------+------------------------------------------------------+
| Group.Subgroup | Field ID | Field Tag                                            |
+----------------+----------+------------------------------------------------------+
| A.1            | 104      | utilization_rate                                     |
| A.1            | 105      | sclk                                                 |
| B.1            | 111      | memory_percent                                       |
| C.1            | 112      | pcie_bw_mb                                           |
| C.1            | 114      | pcie_sent_mb                                         |
| C.1            | 113      | pcie_receive_mb                                      |
+----------------+----------+------------------------------------------------------+
```

该命令中，可通过-i指定设备，并通过-l查看性能指标ID、归属组、名称等信息。

### 4.10.3 查看指定设备组可用指标

**执行命令**

```bash
dcgmi profile -l -g 0
```

```text
[root@localhost ~]# dcgmi profile -l -g 0
+----------------+----------+------------------------------------------------------+
| Group.Subgroup | Field ID | Field Tag                                            |
+----------------+----------+------------------------------------------------------+
| A.1            | 104      | utilization_rate                                     |
| A.1            | 105      | sclk                                                 |
| B.1            | 111      | memory_percent                                       |
| C.1            | 112      | pcie_bw_mb                                           |
| C.1            | 114      | pcie_sent_mb                                         |
| C.1            | 113      | pcie_receive_mb                                      |
+----------------+----------+------------------------------------------------------+
```

该命令中，可通过-g指定设备组，并通过-l查看性能指标ID、归属组、名称等信息。

## 4.11 Idmon

Idmon用于按固定时间间隔采样并输出HCU的时钟、风扇、功耗、温度和利用率指标，支持文本或CSV格式输出.

### 4.11.1 采样HCU监控指标

**执行命令**

```bash
dcgmi idmon -i 0 -s cptu -c 10 -d 1
或：dcgmi idmon -i 0,1 -s cfptu -t csv -f metrics.csv -o
```

```text
命令用法：dcgmi idmon -i 0 -s cptu -c 10 -d 1
```

该命令中，可通过-i指定HCU索引，-s选择采样指标（c为时钟、f为风扇、p为功耗、t为温度、u为利用率），-c指定采样次数，-d指定采样间隔；可通过-t选择text或csv格式，并通过-f指定输出文件、-o输出时间列。

## 4.12 Process

Process用于查看系统中与HCU相关的进程，包括受管理进程列表、KFD进程信息以及进程的显存和计算资源使用情况.

### 4.12.1 查看受管理进程列表

**执行命令**

```bash
dcgmi process list
```

```text
命令用法：dcgmi process list
```

该命令用于列出系统当前受DCGM管理的进程ID。

### 4.12.2 查看KFD进程信息

**执行命令**

```bash
dcgmi process kfd
```

```text
命令用法：dcgmi process kfd
```

该命令用于显示系统中正在运行的KFD进程详细信息。

### 4.12.3 查看HCU进程信息

**执行命令**

```bash
dcgmi process info
```

```text
命令用法：dcgmi process info
```

该命令用于显示HCU相关进程信息，包括进程显存使用量、CU占用情况、SDMA使用量以及关联HCU。

## 4.13 vDevice

vDevice命令用于查看虚拟设备信息、销毁单个虚拟设备以及查看全部物理设备信息.

### 4.13.1 查看虚拟设备信息

**执行命令**

```bash
dcgmi vdevice-info <device-index>
```

```text
命令用法：dcgmi vdevice-info <device-index>
```

该命令根据虚拟设备索引显示虚拟设备名称、计算单元数量、虚拟显存、容器ID、设备索引和PCI总线号等信息。

### 4.13.2 销毁单个虚拟设备

**执行命令**

```bash
dcgmi destroy-vdevice<dvInd>
```

```text
命令用法：dcgmi destroy-vdevice<dvInd>
```

该命令根据虚拟设备索引销毁单个虚拟设备。

### 4.13.3 查看全部物理设备信息

**执行命令**

```bash
dcgmi all-device-infos
```

```text
命令用法：dcgmi all-device-infos
```

该命令以JSON格式输出当前节点全部物理设备的详细信息。

