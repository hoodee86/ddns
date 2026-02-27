package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClientConnection(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("127.0.0.1"))
	}))
	defer server.Close()

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(server.URL)
	if err != nil {
		t.Fatalf("connection failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestLoadConfig(t *testing.T) {
	// 测试加载当前配置文件
	config, err := LoadConfig("config.json")
	if err != nil {
		t.Fatalf("加载配置失败: %v", err)
	}

	if config.Email.SenderEmail == "" {
		t.Error("email.sender_email 不应为空")
	}
	if config.Aliyun.Domain == "" {
		t.Error("aliyun.aliyun_domain 不应为空")
	}
	if config.Aliyun.Domain != "yccloud.org.cn" {
		t.Errorf("域名期望 yccloud.org.cn，实际 %s", config.Aliyun.Domain)
	}
}
