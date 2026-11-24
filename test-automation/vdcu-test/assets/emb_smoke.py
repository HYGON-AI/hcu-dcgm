# -*- coding: utf-8 -*-
"""Qwen3-Embedding-0.6B 向量化推理冒烟（transformers 框架）。

对应 v1.2 方案 6.2 节 INF-204/205：整卡与 vDCU 形态的 Embedding 推理链路验证。
"""
import argparse
import json
import time

import torch
import torch.nn.functional as F
from transformers import AutoModel, AutoTokenizer


def main():
    p = argparse.ArgumentParser()
    p.add_argument("--model", default="/workspace/Qwen--Qwen3-Embedding-0.6B")
    p.add_argument("--slow-tokenizer", action="store_true",
                   help="使用 BertTokenizer（bge 等 BERT 系模型，"
                        "避免 fast tokenizer 后端不兼容）")
    args = p.parse_args()

    prop = torch.cuda.get_device_properties(0)
    torch.manual_seed(0)
    if args.slow_tokenizer:
        from transformers import BertTokenizer
        tokenizer = BertTokenizer.from_pretrained(args.model)
    else:
        tokenizer = AutoTokenizer.from_pretrained(args.model, trust_remote_code=True)
    if tokenizer.pad_token_id is None:
        tokenizer.pad_token_id = tokenizer.eos_token_id
    model = AutoModel.from_pretrained(
        args.model, trust_remote_code=True,
        torch_dtype=torch.float16).to("cuda").eval()

    texts = ["GPU 虚拟化是什么？", "什么是 vDCU？", "DCU 显存切分", "Kubernetes 调度",
             "容器化部署", "分布式训练", "推理服务部署", "性能隔离"] * 2
    enc = tokenizer(texts, padding=True, truncation=True, max_length=128,
                    return_tensors="pt").to("cuda")
    t0 = time.time()
    with torch.no_grad():
        out = model(**enc)
    elapsed = time.time() - t0
    hidden = out.last_hidden_state
    mask = enc["attention_mask"].unsqueeze(-1).float()
    pooled = (hidden * mask).sum(1) / mask.sum(1).clamp(min=1)
    sim = float(F.cosine_similarity(pooled[:1], pooled[1:], dim=-1)[0])

    print("EMB_RESULT " + json.dumps({
        "sentences": len(texts),
        "dim": int(pooled.shape[1]),
        "elapsed_ms": round((elapsed * 1000), 1),
        "visible_gb": round(prop.total_memory / 1024 ** 3, 2),
        "sim_0_1": round(sim, 4),
    }, ensure_ascii=False), flush=True)


if __name__ == "__main__":
    main()
