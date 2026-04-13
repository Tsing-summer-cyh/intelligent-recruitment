package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// RunProbes 启动高并发探测
func RunProbes(targets []Target, timeout time.Duration, verbose bool) []ProbeResult {
	var wg sync.WaitGroup
	// 使用带缓冲的 channel 安全收集结果，避免并发写入切片导致的数据竞态
	resultsChan := make(chan ProbeResult, len(targets))

	for _, t := range targets {
		wg.Add(1)
		go func(target Target) {
			defer wg.Done()

			// 按文档要求：人为加入延迟展示并发性能
			time.Sleep(1 * time.Second)

			// 注入超时控制上下文
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			// 执行探测（包含重试机制）
			result := probeWithRetry(ctx, target, verbose)
			resultsChan <- result
		}(t)
	}

	// 独立协程等待所有任务结束，然后关闭 channel
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// 主协程遍历 channel 聚合结果
	var finalResults []ProbeResult
	for r := range resultsChan {
		finalResults = append(finalResults, r)
	}
	return finalResults
}

// probeWithRetry 实现重试控制
func probeWithRetry(ctx context.Context, target Target, verbose bool) ProbeResult {
	var lastResult ProbeResult
	attempts := target.RetryCount + 1

	for i := 0; i < attempts; i++ {
		start := time.Now()
		var err error

		// 动态识别协议
		if strings.HasPrefix(target.URL, "http://") || strings.HasPrefix(target.URL, "https://") {
			err = probeHTTP(ctx, target.URL, target.Expect)
		} else {
			err = probeTCP(ctx, target.URL)
		}

		lastResult = ProbeResult{
			Target:       target,
			Status:       err == nil,
			ResponseTime: time.Since(start),
		}

		if err != nil {
			lastResult.ErrorMsg = err.Error()
			if verbose {
				fmt.Printf("[FAIL] %s (尝试 %d/%d) - 失败原因: %v\n", target.Name, i+1, attempts, err)
			}
			continue // 失败则继续下一轮重试
		}

		// 成功则跳出循环
		if verbose {
			fmt.Printf("[OK] %s (尝试 %d/%d) - 耗时: %v\n", target.Name, i+1, attempts, lastResult.ResponseTime)
		}
		break
	}
	return lastResult
}

// probeHTTP 执行 HTTP 探测
func probeHTTP(ctx context.Context, urlStr, expect string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 校验 200 OK
	if expect == "200 OK" {
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("非预期状态码: %s", resp.Status)
		}
		return nil
	}

	// 校验期待失败的情况（如无效服务）
	if expect == "Fail" {
		return fmt.Errorf("期待探测失败，但请求成功了")
	}

	// 校验文本包含
	if expect != "" && expect != "Connected" {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		if !strings.Contains(string(body), expect) {
			return fmt.Errorf("响应体中未包含预期内容: %s", expect)
		}
	}

	return nil
}

// probeTCP 执行 TCP 端口探测
func probeTCP(ctx context.Context, addr string) error {
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", addr)
	if err != nil {
		return err
	}
	conn.Close()
	return nil
}