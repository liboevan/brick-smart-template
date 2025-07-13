package models

import (
	"time"
)

// AppStatus 应用状态枚举
type AppStatus string

const (
	AppStatusIdle     AppStatus = "idle"
	AppStatusStarting AppStatus = "starting"
	AppStatusRunning  AppStatus = "running"
	AppStatusStopping AppStatus = "stopping"
	AppStatusStopped  AppStatus = "stopped"
	AppStatusError    AppStatus = "error"
)

// AppInfo 应用配置信息
type AppInfo struct {
	Name                string            `json:"name"`
	Command             string            `json:"command"`
	Args                []string          `json:"args"`
	Env                 map[string]string `json:"env"`
	AutoRestart         bool              `json:"auto_restart"`
	MaxRestarts         int               `json:"max_restarts"`
	HealthCheckInterval int               `json:"health_check_interval"`
}

// AppState 应用运行时状态
type AppState struct {
	Status       AppStatus              `json:"status"`
	PID          *int                   `json:"pid,omitempty"`
	StartTime    *time.Time             `json:"start_time,omitempty"`
	StopTime     *time.Time             `json:"stop_time,omitempty"`
	RestartCount int                    `json:"restart_count"`
	LastError    *string                `json:"last_error,omitempty"`
	Config       map[string]interface{} `json:"config,omitempty"`
}

// StatusReport gRPC状态报告
type StatusReport struct {
	Timestamp time.Time              `json:"timestamp"`
	Status    string                 `json:"status"`
	Data      map[string]interface{} `json:"data"`
}

// HTTP请求/响应结构
type ConfigureAppRequest struct {
	AppInfo AppInfo `json:"app_info"`
}

type ConfigureAppResponse struct {
	Status  string `json:"status"`
	AppName string `json:"app_name"`
}

type StartAppRequest struct {
	Profile string `json:"profile"`
	ID      string `json:"id"`
	AppName string `json:"app_name"`
}

type StartAppResponse struct {
	Status   string `json:"status"`
	AppName  string `json:"app_name"`
	PID      int    `json:"pid"`
	Profile  string `json:"profile"`
}

type StopAppResponse struct {
	Status string `json:"status"`
}

type RestartAppRequest struct {
	Profile string `json:"profile"`
}

type RestartAppResponse struct {
	Status   string `json:"status"`
	AppName  string `json:"app_name"`
	PID      int    `json:"pid"`
	Profile  string `json:"profile"`
}

type AppStatusResponse struct {
	AppName     string                 `json:"app_name"`
	Status      string                 `json:"status"`
	PID         *int                   `json:"pid,omitempty"`
	StartTime   *time.Time             `json:"start_time,omitempty"`
	StopTime    *time.Time             `json:"stop_time,omitempty"`
	RestartCount int                   `json:"restart_count"`
	LastError   *string                `json:"last_error,omitempty"`
	Config      map[string]interface{} `json:"config,omitempty"`
}

type HealthCheckResponse struct {
	Status      string `json:"status"`
	ProxyStatus string `json:"proxy_status"`
	AppStatus   string `json:"app_status"`
} 