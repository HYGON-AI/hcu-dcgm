/*
 * SPDX-License-Identifier: Apache-2.0
 * Copyright (c) 2026 Hygon Information Technology Co., Ltd.
 */
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func main() {
	// 参数：run 文件名 可改为从 flag 传入
	var (
		runFile = flag.String("file", "dcgm-hcu-v1.0.1.run", "run installer filename (in current directory)")
		timeout = flag.Int("timeout", 0, "timeout in seconds (0 = no timeout)")
	)
	flag.Parse()

	// 确保在当前目录查找
	cwd, _ := os.Getwd()
	fullPath := fmt.Sprintf("%s/%s", cwd, *runFile)

	// 检查文件存在
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "error: file not found: %s\n", fullPath)
		os.Exit(2)
	}

	// 如果不是可执行，尝试 chmod +x
	info, err := os.Stat(fullPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: stat failed: %v\n", err)
		os.Exit(3)
	}
	mode := info.Mode()
	if mode&0111 == 0 {
		if err := os.Chmod(fullPath, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to chmod +x: %v\n", err)
			// 仍然尝试执行
		} else {
			fmt.Printf("made %s executable\n", fullPath)
		}
	}

	// 构造上下文（可带超时）
	var ctx context.Context
	var cancel context.CancelFunc
	if *timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(*timeout)*time.Second)
		defer cancel()
	} else {
		ctx, cancel = context.WithCancel(context.Background())
		defer cancel()
	}

	// 使用 bash -c "./file" 以支持脚本内部用到的 shell 特性
	cmd := exec.CommandContext(ctx, "bash", "-c", fmt.Sprintf("./%s", *runFile))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	fmt.Printf("Running %s ...\n\n", fullPath)
	start := time.Now()
	if err := cmd.Run(); err != nil {
		// 如果是因为 context 超时
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Fprintf(os.Stderr, "\nERROR: command timed out after %d seconds\n", *timeout)
			// 尝试返回 124 类似 timeout 工具行为
			os.Exit(124)
		}

		// 非零退出，尝试解析退出码或 signal
		if exitErr, ok := err.(*exec.ExitError); ok {
			if ws, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				if ws.Signaled() {
					sig := ws.Signal()
					fmt.Fprintf(os.Stderr, "\nProcess terminated by signal: %v (exit status = %d)\n", sig, 128+int(sig))
					os.Exit(128 + int(sig))
				}
				code := ws.ExitStatus()
				fmt.Fprintf(os.Stderr, "\nProcess exited with code: %d\n", code)
				os.Exit(code)
			}
		}

		// 其它错误
		fmt.Fprintf(os.Stderr, "\nFailed to run command: %v\n", err)
		os.Exit(1)
	}

	elapsed := time.Since(start)
	fmt.Printf("\nProcess finished successfully (elapsed %s)\n", elapsed)
	os.Exit(0)
}
