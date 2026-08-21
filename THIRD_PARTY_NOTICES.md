# Third-Party Notices

本文件记录 HCU DCGM 源码树中随项目分发的外部驱动接口头文件。以下文件不是 HCU DCGM 项目自行编写的业务源码，而是从 HCU Driver 安装路径取得、并为保证公开源码可直接构建而保留在仓库中的未修改接口文件副本。

## Hygon DMI headers

- **Files:**
  - `pkg/dcgm/include/dmi.h`
  - `pkg/dcgm/include/dmi_error.h`
  - `pkg/dcgm/include/dmi_mig.h`
  - `pkg/dcgm/include/dmi_virtual.h`
- **Original upstream:** Hygon/HCU DMI driver interface.
- **HCU Driver:** `HCU Driver` `6.3.*` series.
- **Installation path:** `/opt/hyhal/include/dmi/`.
- **Copyright:** SW Group, Chengdu Haiguang IC Design Co., Ltd.
- **License:** Apache License 2.0 (`Apache-2.0`), as identified in the files' headers.
- **Modification status:** Copied from the HCU Driver installation without local modification. Repository copies are maintained for compatibility with the HCU Driver `6.3.*` series; different patch versions can contain version or generated-metadata differences.

## AMD KFD interface header

- **File:** `pkg/dcgm/include/kfd_ioctl.h`
- **Original upstream:** AMD KFD/Linux KFD interface.
- **HCU Driver:** `HCU Driver` `6.3.*` series.
- **Installation path:** `/opt/hyhal/include/rocm_smi/kfd_ioctl.h`.
- **Copyright:** Advanced Micro Devices, Inc.
- **License:** MIT, as identified in the file's header.
- **Modification status:** Copied from the HCU Driver installation without local modification. Repository copies are maintained for compatibility with the HCU Driver `6.3.*` series; different patch versions can contain version or generated-metadata differences.

## AMD ROCm SMI headers

- **Files:**
  - `pkg/dcgm/include/rocm_smi.h`
  - `pkg/dcgm/include/rocm_smi64Config.h`
- **Original upstream:** AMD ROCm SMI.
- **HCU Driver:** `HCU Driver` `6.3.*` series.
- **Installation path:** `/opt/hyhal/include/rocm_smi/`.
- **Copyright:** Advanced Micro Devices, Inc.
- **License:** University of Illinois/NCSA Open Source License (`NCSA`), as identified in each file's header.
- **Modification status:** Copied from the HCU Driver installation without local modification. `rocm_smi64Config.h` is a driver-build-generated configuration header; its version metadata can vary between compatible HCU Driver `6.3.*` patch versions.
- **Compliance status:** The original AMD copyright and NCSA license text remain intact. Public-distribution handling follows the legal/compliance direction for AMD/ROCm copyright and license-bearing files.

## Distribution and dependency boundary

- These headers are distributed as source files because the public release tree must contain the interfaces required by the cgo compilation units.
- The corresponding runtime libraries are not included in this repository. They are supplied by the installed HCU Driver, including `libhydmi.so`, `librocm_smi64.so`, and, for MIG functionality, `libhydmi_mig.so`.
- The original copyright and license text in each external header must remain intact in source distributions.
