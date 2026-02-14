package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	// Server 的地址（固定 IP 地址和端口）
	const serverURL = "http://43.110.38.163:8999"

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	fmt.Printf("Client 启动，每隔 1 秒向 %s 发送请求\n", serverURL)

	for {
		resp, err := client.Get(serverURL)
		if err != nil {
			fmt.Printf("[%s] 请求失败: %v\n", time.Now().Format("15:04:05"), err)
		} else {
			body, err := io.ReadAll(resp.Body)
			resp.Body.Close()

			if err != nil {
				fmt.Printf("[%s] 读取响应失败: %v\n", time.Now().Format("15:04:05"), err)
			} else {
				fmt.Printf("[%s] 我的 IP 地址是: %s\n", time.Now().Format("15:04:05"), string(body))
			}
		}

		time.Sleep(1 * time.Second)
	}
}
