package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"gopkg.in/gomail.v2"
)

// sendEmailNotification 发送 IP 变化通知邮件
func sendEmailNotification(config *EmailConfig, oldIP, newIP string) {
	m := gomail.NewMessage()
	m.SetHeader("From", config.SenderEmail)
	m.SetHeader("To", config.RecipientEmail)
	m.SetHeader("Subject", "IP 地址变化通知")

	body := fmt.Sprintf(
		"检测到 IP 地址变化\n\n旧 IP: %s\n新 IP: %s\n变化时间: %s",
		oldIP,
		newIP,
		time.Now().Format("2006-01-02 15:04:05"),
	)
	m.SetBody("text/plain; charset=UTF-8", body)

	d := gomail.NewDialer(config.SMTPServer, config.SMTPPort, config.SenderEmail, config.SenderPassword)

	if err := d.DialAndSend(m); err != nil {
		fmt.Printf("[%s] 邮件发送失败: %v\n", time.Now().Format("15:04:05"), err)
	} else {
		fmt.Printf("[%s] 邮件发送成功: IP 从 %s 变更为 %s\n", time.Now().Format("15:04:05"), oldIP, newIP)
	}
}

// updateAliyunDNS 更新阿里云 DNS 解析记录
func updateAliyunDNS(dnsClient *AliyunDNSClient, ip string) {
	if err := dnsClient.UpdateDNSRecord(ip); err != nil {
		fmt.Printf("[%s] 阿里云 DNS 更新失败: %v\n", time.Now().Format("15:04:05"), err)
	}
}

func main() {
	// 加载配置
	config, err := LoadConfig("config.json")
	if err != nil {
		fmt.Printf("配置加载失败: %v\n", err)
		return
	}

	// 初始化阿里云 DNS 客户端
	dnsClient, err := NewAliyunDNSClient(&config.Aliyun)
	if err != nil {
		fmt.Printf("阿里云 DNS 客户端初始化失败: %v\n", err)
		return
	}

	// Server 的地址（固定 IP 地址和端口）
	const serverURL = "http://43.110.38.163:8999"

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	fmt.Printf("Client 启动，每隔 30 秒向 %s 发送请求\n", serverURL)
	fmt.Printf("域名: %s.%s\n", config.Aliyun.RR, config.Aliyun.Domain)

	var lastIP string // 记录上一次的 IP 地址

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
				currentIP := strings.TrimSpace(string(body))
				fmt.Printf("[%s] 我的 IP 地址是: %s\n", time.Now().Format("15:04:05"), currentIP)

				// 检测 IP 是否变化（包括首次获取）
				if lastIP != currentIP {
					// 更新阿里云 DNS 解析
					go updateAliyunDNS(dnsClient, currentIP)

					// 非首次获取时发送邮件通知
					if lastIP != "" {
						fmt.Printf("[%s] 检测到 IP 变化: %s -> %s\n", time.Now().Format("15:04:05"), lastIP, currentIP)
						go sendEmailNotification(&config.Email, lastIP, currentIP)
					}

					lastIP = currentIP
				}
			}
		}

		time.Sleep(30 * time.Second) // 修改为 30 秒
	}
}
