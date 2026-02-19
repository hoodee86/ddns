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

// LoadConfig 从config.json文件加载配置
func LoadConfig(filepath string) (*EmailConfig, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var config EmailConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 验证必要字段
	if config.SenderEmail == "" {
		return nil, fmt.Errorf("配置错误: sender_email 不能为空")
	}
	if config.SenderPassword == "" {
		return nil, fmt.Errorf("配置错误: sender_password 不能为空")
	}
	if config.RecipientEmail == "" {
		return nil, fmt.Errorf("配置错误: recipient_email 不能为空")
	}
	if config.SMTPServer == "" {
		return nil, fmt.Errorf("配置错误: smtp_server 不能为空")
	}
	if config.SMTPPort == 0 {
		return nil, fmt.Errorf("配置错误: smtp_port 不能为空")
	}

	return &config, nil
}
