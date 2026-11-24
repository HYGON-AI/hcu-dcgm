# -*- coding: utf-8 -*-
"""最小集合通信用例：验证两个 vDCU 之间的 NCCL all-reduce。"""
import os

import torch
import torch.distributed as dist

rank = int(os.environ["RANK"])
world = int(os.environ["WORLD_SIZE"])

print("init rank=%d world=%d devices=%d" %
      (rank, world, torch.cuda.device_count()), flush=True)

dist.init_process_group("nccl", rank=rank, world_size=world)
t = torch.ones(4, device="cuda")
dist.all_reduce(t)
if rank == 0:
    print("ALLREDUCE_OK", t.tolist(), flush=True)
dist.destroy_process_group()
