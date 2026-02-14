package main

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

// getClientIP 从请求中获取客户端的 IP 地址
func getClientIP(r *http.Request) string {
	// 首先尝试获取 X-Forwarded-For 头（在代理服务器后）
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// X-Forwarded-For 可能包含多个 IP 地址，取第一个
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}

	// 尝试获取 X-Real-IP 头
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}

	// 从 RemoteAddr 中获取 IP 地址
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

// handler 处理客户端请求
func handler(w http.ResponseWriter, r *http.Request) {
	clientIP := getClientIP(r)
	fmt.Printf("来自 %s 的请求\n", clientIP)
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "%s", clientIP)
}

func main() {
	const port = ":8999"
	fmt.Printf("Server 启动在 %s\n", port)

	http.HandleFunc("/", handler)
	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Printf("Server 错误: %v\n", err)
	}
}
