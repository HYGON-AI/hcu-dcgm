# -*- coding: utf-8 -*-
"""小模型生成冒烟：验证 vDCU 上的推理链路（对应 VDCU-INF-001 轻量版）。"""
import argparse
import json
import time

import torch
from transformers import AutoModelForCausalLM, AutoTokenizer


def main():
    p = argparse.ArgumentParser()
    p.add_argument("--model", default="/workspace/Qwen--Qwen3.5-0.8B")
    p.add_argument("--prompt", default="用一句话介绍 GPU 虚拟化。")
    p.add_argument("--max-new-tokens", type=int, default=32)
    args = p.parse_args()

    torch.manual_seed(0)
    prop = torch.cuda.get_device_properties(0)
    tokenizer = AutoTokenizer.from_pretrained(args.model)
    model = AutoModelForCausalLM.from_pretrained(
        args.model, torch_dtype=torch.float16).to("cuda")
    model.eval()

    inputs = tokenizer(args.prompt, return_tensors="pt").to("cuda")
    t0 = time.time()
    with torch.no_grad():
        out = model.generate(**inputs, max_new_tokens=args.max_new_tokens,
                             do_sample=False, pad_token_id=tokenizer.eos_token_id)
    elapsed = time.time() - t0
    new_tokens = int(out.shape[1] - inputs["input_ids"].shape[1])

    print("GEN_RESULT " + json.dumps({
        "new_tokens": new_tokens,
        "elapsed_s": round(elapsed, 2),
        "tok_per_s": round(new_tokens / elapsed, 2) if elapsed > 0 else 0,
        "mem_gb": round(prop.total_memory / 1024 ** 3, 2),
        "text": tokenizer.decode(out[0][inputs["input_ids"].shape[1]:],
                                 skip_special_tokens=True)[:80],
    }, ensure_ascii=False), flush=True)


if __name__ == "__main__":
    main()
