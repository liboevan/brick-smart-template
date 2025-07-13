package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Client HTTP客户端
type Client struct {
	httpPort string
}

// NewClient 创建新的HTTP客户端
func NewClient(httpPort string) *Client {
	if httpPort == "" {
		httpPort = "17100"
	}
	return &Client{
		httpPort: httpPort,
	}
}

// ReportStatus 上报设备状态
func (c *Client) ReportStatus(ctx context.Context, status map[string]interface{}) error {
	// 构建HTTP请求数据
	requestData := map[string]interface{}{
		"status": "running",
		"data":   status["data"],
	}

	requestJSON, err := json.Marshal(requestData)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}

	// 通过HTTP API发送状态报告
	url := fmt.Sprintf("http://localhost:%s/app/status/report", c.httpPort)
	
	// 发送HTTP POST请求
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestJSON))
	if err != nil {
		return fmt.Errorf("failed to send status report: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status report failed with status: %d", resp.StatusCode)
	}

	log.Printf("Status reported successfully for device %s", status["device_id"])
	return nil
} 