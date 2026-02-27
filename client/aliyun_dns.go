package main

import (
	"fmt"
	"time"

	alidns "github.com/alibabacloud-go/alidns-20150109/v4/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/alibabacloud-go/tea/tea"
)

// AliyunDNSClient 阿里云 DNS 客户端
type AliyunDNSClient struct {
	client *alidns.Client
	config *AliyunConfig
}

// NewAliyunDNSClient 创建阿里云 DNS 客户端
func NewAliyunDNSClient(config *AliyunConfig) (*AliyunDNSClient, error) {
	openApiConfig := &openapi.Config{
		AccessKeyId:     tea.String(config.AccessKeyID),
		AccessKeySecret: tea.String(config.AccessKeySecret),
		Endpoint:        tea.String("alidns.cn-hangzhou.aliyuncs.com"),
	}

	client, err := alidns.NewClient(openApiConfig)
	if err != nil {
		return nil, fmt.Errorf("创建阿里云 DNS 客户端失败: %w", err)
	}

	return &AliyunDNSClient{
		client: client,
		config: config,
	}, nil
}

// UpdateDNSRecord 更新 DNS 解析记录，如果记录不存在则新建
func (c *AliyunDNSClient) UpdateDNSRecord(ip string) error {
	domain := c.config.Domain
	rr := c.config.RR // 主机记录，如 "@" 表示根域名, "www" 表示 www 子域名

	// 1. 先查询已有的解析记录
	recordID, oldIP, err := c.getExistingRecord(domain, rr)
	if err != nil {
		return fmt.Errorf("查询 DNS 记录失败: %w", err)
	}

	// 2. 如果记录已存在且 IP 相同，无需更新
	if recordID != "" && oldIP == ip {
		fmt.Printf("[%s] DNS 记录已是最新，无需更新 (%s.%s -> %s)\n",
			time.Now().Format("15:04:05"), rr, domain, ip)
		return nil
	}

	// 3. 如果记录已存在，更新记录
	if recordID != "" {
		return c.updateRecord(recordID, rr, ip)
	}

	// 4. 如果记录不存在，新增记录
	return c.addRecord(domain, rr, ip)
}

// getExistingRecord 查询已有的 A 记录
func (c *AliyunDNSClient) getExistingRecord(domain, rr string) (recordID, ip string, err error) {
	request := &alidns.DescribeDomainRecordsRequest{
		DomainName: tea.String(domain),
		RRKeyWord:  tea.String(rr),
		Type:       tea.String("A"),
	}

	response, err := c.client.DescribeDomainRecords(request)
	if err != nil {
		return "", "", err
	}

	if response.Body != nil && response.Body.DomainRecords != nil {
		for _, record := range response.Body.DomainRecords.Record {
			if tea.StringValue(record.RR) == rr && tea.StringValue(record.Type) == "A" {
				return tea.StringValue(record.RecordId), tea.StringValue(record.Value), nil
			}
		}
	}

	return "", "", nil
}

// updateRecord 更新已有的 DNS 记录
func (c *AliyunDNSClient) updateRecord(recordID, rr, ip string) error {
	request := &alidns.UpdateDomainRecordRequest{
		RecordId: tea.String(recordID),
		RR:       tea.String(rr),
		Type:     tea.String("A"),
		Value:    tea.String(ip),
	}

	_, err := c.client.UpdateDomainRecord(request)
	if err != nil {
		return fmt.Errorf("更新 DNS 记录失败: %w", err)
	}

	fmt.Printf("[%s] DNS 记录更新成功: %s.%s -> %s\n",
		time.Now().Format("15:04:05"), rr, c.config.Domain, ip)
	return nil
}

// addRecord 新增 DNS 记录
func (c *AliyunDNSClient) addRecord(domain, rr, ip string) error {
	request := &alidns.AddDomainRecordRequest{
		DomainName: tea.String(domain),
		RR:         tea.String(rr),
		Type:       tea.String("A"),
		Value:      tea.String(ip),
	}

	_, err := c.client.AddDomainRecord(request)
	if err != nil {
		return fmt.Errorf("新增 DNS 记录失败: %w", err)
	}

	fmt.Printf("[%s] DNS 记录新增成功: %s.%s -> %s\n",
		time.Now().Format("15:04:05"), rr, domain, ip)
	return nil
}
