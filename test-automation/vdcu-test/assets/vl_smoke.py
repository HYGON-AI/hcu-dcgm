# -*- coding: utf-8 -*-
"""Qwen3-VL-2B-Instruct 多模态推理冒烟（transformers 框架）。

对应 v1.2 方案 6.2 节 INF-201~203：整卡与 vDCU 形态的图文推理链路验证。
"""
import argparse
import json
import time

import torch
from PIL import Image


def main():
    p = argparse.ArgumentParser()
    p.add_argument("--model", default="/workspace/Qwen--Qwen3-VL-2B-Instruct")
    p.add_argument("--max-new-tokens", type=int, default=16)
    args = p.parse_args()

    prop = torch.cuda.get_device_properties(0)
    torch.manual_seed(0)

    try:
        from transformers import (AutoProcessor,
                                  Qwen3VLForConditionalGeneration)
        processor = AutoProcessor.from_pretrained(
            args.model, trust_remote_code=True)
        model = Qwen3VLForConditionalGeneration.from_pretrained(
            args.model, torch_dtype=torch.float16,
            attn_implementation="eager").to("cuda").eval()
    except ImportError:
        from transformers import (AutoModelForImageTextToText,
                                  AutoProcessor)
        processor = AutoProcessor.from_pretrained(
            args.model, trust_remote_code=True)
        model = AutoModelForImageTextToText.from_pretrained(
            args.model, torch_dtype=torch.float16,
            attn_implementation="eager").to("cuda").eval()

    img = Image.new("RGB", (224, 224), color=(120, 60, 30))
    # Qwen3-VL 的图像占位 token 由 chat template 注入，必须先走模板再编码
    messages = [{"role": "user", "content": [
        {"type": "image", "image": img},
        {"type": "text", "text": "用一句话描述这张图片。"}]}]
    try:
        text = processor.apply_chat_template(
            messages, tokenize=False, add_generation_prompt=True,
            images=[img])
    except TypeError:
        text = processor.apply_chat_template(
            messages, tokenize=False, add_generation_prompt=True)
    inputs = processor(text=[text], images=[img],
                       return_tensors="pt").to("cuda")
    t0 = time.time()
    with torch.no_grad():
        out = model.generate(**inputs, max_new_tokens=args.max_new_tokens,
                             do_sample=False)
    elapsed = time.time() - t0
    new_tokens = int(out.shape[1] - inputs["input_ids"].shape[1])
    text = ""
    try:
        text = processor.batch_decode(
            out[:, inputs["input_ids"].shape[1]:],
            skip_special_tokens=True)[0][:80]
    except Exception:
        pass

    print("VL_RESULT " + json.dumps({
        "new_tokens": new_tokens,
        "elapsed_s": round(elapsed, 2),
        "tok_per_s": round(new_tokens / elapsed, 2) if elapsed > 0 else 0,
        "visible_gb": round(prop.total_memory / 1024 ** 3, 2),
        "text": text,
    }, ensure_ascii=False), flush=True)


if __name__ == "__main__":
    main()
