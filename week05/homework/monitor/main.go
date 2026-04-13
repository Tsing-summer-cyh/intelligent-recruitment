package main

import (
	"flag"
	"fmt"
	"log"
	"time"
)

func main() {
	// 1. 解析命令行参数
	configFile := flag.String("config", "config.json", "指定配置文件的路径")
	timeoutSec := flag.Int("timeout", 3, "单次探测超时时间（单位：秒）")
	verbose := flag.Bool("v", false, "开启详细模式，实时打印每一个目标的探测状态")
	flag.Parse()

	// 2. 加载配置
	cfg, err := LoadConfig(*configFile)
	if err != nil {
		log.Fatalf("加载配置文件失败: %v\n", err)
	}

	if *verbose {
		fmt.Println("======== 开始并发健康探测 ========")
	}

	// 3. 启动并发探测引擎
	timeout := time.Duration(*timeoutSec) * time.Second
	results := RunProbes(cfg.Targets, timeout, *verbose)

	// 4. 生成分析报表
	GenerateReport(results)
}