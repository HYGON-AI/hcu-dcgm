# -*- coding: utf-8 -*-
"""Qwen3.5-0.8B 小模型训练冒烟（transformers 框架）。

对应 v1.2 方案 5.4 节 TRN-201/202：小模型训练链路以 transformers 框架验证
（Megatron 版 0.8B 预训练受 transformer_engine 依赖限制，与 v1.2 实测一致）。
"""
import argparse
import json
import time

import torch
from transformers import AutoModelForCausalLM, AutoTokenizer


def main():
    p = argparse.ArgumentParser()
    p.add_argument("--model", default="/workspace/Qwen--Qwen3.5-0.8B")
    p.add_argument("--steps", type=int, default=30)
    p.add_argument("--batch", type=int, default=2)
    p.add_argument("--seq-len", type=int, default=128)
    args = p.parse_args()

    prop = torch.cuda.get_device_properties(0)
    torch.manual_seed(0)
    tokenizer = AutoTokenizer.from_pretrained(args.model, trust_remote_code=True)
    if tokenizer.pad_token_id is None:
        tokenizer.pad_token_id = tokenizer.eos_token_id
    model = AutoModelForCausalLM.from_pretrained(
        args.model, torch_dtype=torch.bfloat16).to("cuda")
    model.train()
    opt = torch.optim.AdamW(model.parameters(), lr=1e-5)

    # 优先使用 weibo_sentiment 真实语料（训练链路口径），缺失时合成
    samples = []
    try:
        import csv
        with open("/workspace/weibo_sentiment/weibo_senti_100k.csv",
                  encoding="utf-8") as f:
            for row in csv.reader(f):
                if len(row) > 1 and row[1].strip():
                    samples.append(row[1].strip())
                if len(samples) >= 800:
                    break
    except OSError:
        pass
    if len(samples) < 16:
        words = ["word%d" % i for i in range(2000)]
        samples = [" ".join(words[i % 2000:(i % 2000) + args.seq_len])
                   for i in range(args.batch * args.steps * 2)]

    n_tok = 0
    t0 = time.time()
    for step in range(args.steps):
        batch = samples[step * args.batch:(step + 1) * args.batch]
        enc = tokenizer(batch, padding=True, truncation=True,
                        max_length=args.seq_len, return_tensors="pt").to("cuda")
        out = model(**enc, labels=enc["input_ids"])
        loss = out.loss
        opt.zero_grad()
        loss.backward()
        torch.nn.utils.clip_grad_norm_(model.parameters(), 1.0)
        opt.step()
        n_tok += int(enc["input_ids"].numel())
        if step % 10 == 0 or step == args.steps - 1:
            print("step %d loss=%.4f" % (step, loss.item()), flush=True)
    elapsed = time.time() - t0

    print("TRAIN_RESULT " + json.dumps({
        "steps": args.steps,
        "loss": round(float(loss.item()), 4),
        "elapsed_s": round(elapsed, 2),
        "tokens": n_tok,
        "tok_per_s": round(n_tok / elapsed, 1),
        "mem_gb": round(torch.cuda.max_memory_allocated() / 1024 ** 3, 2),
        "visible_gb": round(prop.total_memory / 1024 ** 3, 2),
    }, ensure_ascii=False), flush=True)


if __name__ == "__main__":
    main()
