# -*- coding: utf-8 -*-
"""整卡集合通信长稳：100 次 64 MB all-reduce，统计 P50/P95/最大耗时。

经 torchrun --nproc_per_node=N 启动（对应 VDCU-COM-014）。
"""
import time

import torch
import torch.distributed as dist

dist.init_process_group("nccl")
rank = dist.get_rank()
t = torch.empty(64 * 1024 ** 2 // 4, dtype=torch.float32, device="cuda")
times = []
for i in range(100):
    t0 = time.time()
    dist.all_reduce(t)
    torch.cuda.synchronize()
    times.append((time.time() - t0) * 1000)
times.sort()
n = len(times)
if rank == 0:
    print("AR_LOOP_OK total=%d p50=%.2fms p95=%.2fms max=%.2fms" % (
        n, times[n // 2], times[int(n * 0.95)], times[-1]), flush=True)
dist.destroy_process_group()
