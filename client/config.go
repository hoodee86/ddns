package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// EmailConfig 邮件配置结构体
type EmailConfig struct {
	SenderEmail    string `json:"sender_email"`
	SenderPassword string `json:"sender_password"`
	RecipientEmail string `json:"recipient_email"`
	SMTPServer     string `json:"smtp_server"`
	SMTPPort       int    `json:"smtp_port"`
}

// AliyunConfig 阿里云配置结构体
type AliyunConfig struct {
	AccessKeyID     string `json:"aliyun_access_key_id"`
	AccessKeySecret string `json:"aliyun_access_key_secret"`
	Domain          string `json:"aliyun_domain"`
	RR              string `json:"aliyun_rr"` // 主机记录，如 "@" 或 "www"
}

// AppConfig 应用总配置
type AppConfig struct {
	Email  EmailConfig  `json:"email"`
	Aliyun AliyunConfig `json:"aliyun"`
}

// LoadConfig 从config.json文件加载配置
func LoadConfig(filepath string) (*AppConfig, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var config AppConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 验证邮件配置字段
	if config.Email.SenderEmail == "" {
		return nil, fmt.Errorf("配置错误: email.sender_email 不能为空")
	}
	if config.Email.SenderPassword == "" {
		return nil, fmt.Errorf("配置错误: email.sender_password 不能为空")
	}
	if config.Email.RecipientEmail == "" {
		return nil, fmt.Errorf("配置错误: email.recipient_email 不能为空")
	}
	if config.Email.SMTPServer == "" {
		return nil, fmt.Errorf("配置错误: email.smtp_server 不能为空")
	}
	if config.Email.SMTPPort == 0 {
		return nil, fmt.Errorf("配置错误: email.smtp_port 不能为空")
	}

	// 验证阿里云配置字段
	if config.Aliyun.AccessKeyID == "" {
		return nil, fmt.Errorf("配置错误: aliyun.aliyun_access_key_id 不能为空")
	}
	if config.Aliyun.AccessKeySecret == "" {
		return nil, fmt.Errorf("配置错误: aliyun.aliyun_access_key_secret 不能为空")
	}
	if config.Aliyun.Domain == "" {
		return nil, fmt.Errorf("配置错误: aliyun.aliyun_domain 不能为空")
	}
	if config.Aliyun.RR == "" {
		config.Aliyun.RR = "@" // 默认为根域名
	}

	return &config, nil
}
