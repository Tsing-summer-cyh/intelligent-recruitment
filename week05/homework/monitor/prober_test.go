package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// 测试 HTTP 探测成功的情况
func TestProbeHTTP_Success(t *testing.T) {
	// 启动一个本地 Mock 服务器
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello Go Native Cloud"))
	}))
	defer ts.Close()

	ctx := context.Background()

	// 场景 1：期待 200 OK
	err := probeHTTP(ctx, ts.URL, "200 OK")
	if err != nil {
		t.Errorf("期待成功，但得到错误: %v", err)
	}

	// 场景 2：期待包含特定文本
	err = probeHTTP(ctx, ts.URL, "Native")
	if err != nil {
		t.Errorf("期待文本匹配成功，但得到错误: %v", err)
	}
}

// 测试 Context 超时控制机制
func TestProbeHTTP_Timeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second) // 模拟慢速响应
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// 设置 1 秒超时
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := probeHTTP(ctx, ts.URL, "200 OK")
	if err == nil {
		t.Errorf("期待探测因超时失败，但返回了成功")
	}
}