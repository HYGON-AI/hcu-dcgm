# -*- coding: utf-8 -*-
"""ModernBERT 短跑冒烟：用于 vDCU 训练链路验证（对应 VDCU-TRN-206）。"""
import argparse
import json
import time

import pandas as pd
import torch
from torch.utils.data import DataLoader, Dataset
from transformers import AutoModelForSequenceClassification, AutoTokenizer


class TextDS(Dataset):
    def __init__(self, texts, labels, tokenizer, max_len):
        self.texts = texts
        self.labels = labels
        self.tokenizer = tokenizer
        self.max_len = max_len

    def __len__(self):
        return len(self.texts)

    def __getitem__(self, idx):
        enc = self.tokenizer(
            str(self.texts[idx]),
            truncation=True,
            padding="max_length",
            max_length=self.max_len,
            return_tensors="pt",
        )
        item = {k: v.squeeze(0) for k, v in enc.items()}
        item["labels"] = torch.tensor(int(self.labels[idx]), dtype=torch.long)
        return item


def main():
    p = argparse.ArgumentParser()
    p.add_argument("--model", default="/workspace/answerdotai--ModernBERT-base")
    p.add_argument("--data", default="/workspace/weibo_sentiment/weibo_senti_100k.csv")
    p.add_argument("--samples", type=int, default=800)
    p.add_argument("--epochs", type=int, default=1)
    p.add_argument("--batch", type=int, default=8)
    p.add_argument("--max-len", type=int, default=64)
    args = p.parse_args()

    torch.manual_seed(42)
    device = torch.device("cuda:0")
    torch.cuda.set_device(device)
    prop = torch.cuda.get_device_properties(0)
    print("DEVICE %s MEM_GB %.2f" % (prop.name, prop.total_memory / 1024 ** 3), flush=True)

    tokenizer = AutoTokenizer.from_pretrained(args.model)
    df = pd.read_csv(args.data, header=0, names=["label", "text"], nrows=args.samples)
    df = df.dropna()
    texts = df["text"].astype(str).tolist()
    labels = df["label"].astype(int).tolist()

    ds = TextDS(texts, labels, tokenizer, args.max_len)
    dl = DataLoader(ds, batch_size=args.batch, shuffle=True, num_workers=0)

    model = AutoModelForSequenceClassification.from_pretrained(
        args.model, num_labels=2).to(device)
    opt = torch.optim.AdamW(model.parameters(), lr=2e-5)

    t0 = time.time()
    model.train()
    total, correct, loss_sum = 0, 0, 0.0
    for epoch in range(args.epochs):
        for batch in dl:
            batch = {k: v.to(device) for k, v in batch.items()}
            opt.zero_grad()
            out = model(**batch)
            loss = out.loss
            loss.backward()
            opt.step()
            loss_sum += loss.item()
            pred = out.logits.argmax(dim=-1)
            correct += (pred == batch["labels"]).sum().item()
            total += batch["labels"].size(0)
    elapsed = time.time() - t0

    res = {
        "samples": total,
        "epochs": args.epochs,
        "elapsed_s": round(elapsed, 2),
        "throughput": round(total / elapsed, 2) if elapsed > 0 else 0,
        "loss": round(loss_sum / max(1, len(dl) * args.epochs), 4),
        "acc": round(correct / max(1, total), 4),
        "mem_gb": round(prop.total_memory / 1024 ** 3, 2),
    }
    print("SMOKE_RESULT " + json.dumps(res), flush=True)


if __name__ == "__main__":
    main()
