package appmanager

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"

	"brick-smart-template/pkg/models"

	"github.com/sirupsen/logrus"
	"io/ioutil"
)

// Manager 应用管理器
type Manager struct {
	mu           sync.RWMutex
	appInfo      *models.AppInfo
	appState     *models.AppState
	cmd          *exec.Cmd
	ctx          context.Context
	cancel       context.CancelFunc
	healthTicker *time.Ticker
	logger       *logrus.Logger
	internalStatus map[string]interface{} // 存储app内部状态
	proxyID      string // 新增：proxy/app id
	lastProfile  string // 上次启动用的 profile
}

// NewManager 创建新的应用管理器
func NewManager(logger *logrus.Logger, proxyID string) *Manager {
	m := &Manager{
		appState: &models.AppState{
			Status: "ready",
		},
		logger: logger,
		proxyID: proxyID,
	}
	// 启动时尝试读取 /app/manifest.json
	manifestPath := "/app/manifest.json"
	if data, err := ioutil.ReadFile(manifestPath); err == nil {
		var manifest struct {
			AppName             string   `json:"app_name"`
			HealthCheckInterval int      `json:"health_check_interval"`
			DefaultArgs         []string `json:"default_args"`
		}
		if err := json.Unmarshal(data, &manifest); err == nil {
			m.appInfo = &models.AppInfo{
				Name:                manifest.AppName,
				Command:             "./" + manifest.AppName,
				Args:                manifest.DefaultArgs,
				Env:                 map[string]string{},
				AutoRestart:         false,
				MaxRestarts:         3,
				HealthCheckInterval: manifest.HealthCheckInterval,
			}
		}
	}
	return m
}

// ConfigureApp 配置应用
func (m *Manager) ConfigureApp(appInfo models.AppInfo) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.appInfo = &appInfo
	m.logger.Infof("Configured app: %s", appInfo.Name)
	return nil
}

// StartApp 启动应用（互斥、幂等、状态检查、自动补全 -id 参数）
func (m *Manager) StartApp(profile string) (*models.StartAppResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.appInfo == nil {
		return nil, fmt.Errorf("app not configured")
	}

	if m.appState.Status == models.AppStatusStarting || m.appState.Status == models.AppStatusRunning {
		// 幂等：已在运行直接返回
		return &models.StartAppResponse{
			Status:  "already_running",
			AppName: m.appInfo.Name,
			PID:     *m.appState.PID,
			Profile: profile,
		}, nil
	}

	// 解析profile
	var config map[string]interface{}
	if err := json.Unmarshal([]byte(profile), &config); err != nil {
		return nil, fmt.Errorf("invalid JSON profile: %v", err)
	}

	// 自动补全 -id 参数
	args := m.appInfo.Args
	idPresent := false
	for i, arg := range args {
		if arg == "-id" && i+1 < len(args) {
			args[i+1] = m.proxyID
			idPresent = true
		}
	}
	if !idPresent {
		args = append(args, "-id", m.proxyID)
	}

	// 设置环境变量
	env := os.Environ()
	for k, v := range m.appInfo.Env {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	env = append(env, fmt.Sprintf("APP_PROFILE=%s", profile))
	env = append(env, fmt.Sprintf("APP_NAME=%s", m.appInfo.Name))
	env = append(env, fmt.Sprintf("PROXY_GRPC_PORT=%s", os.Getenv("GRPC_PORT")))

	cmd := exec.CommandContext(context.Background(), m.appInfo.Command, args...)
	cmd.Env = env
	// 不再设置 cmd.Dir

	// 启动进程
	if err := cmd.Start(); err != nil {
		m.appState.Status = models.AppStatusError
		errorMsg := err.Error()
		m.appState.LastError = &errorMsg
		m.logger.Errorf("Failed to start app %s: %v", m.appInfo.Name, err)
		return nil, err
	}

	m.cmd = cmd
	pid := cmd.Process.Pid
	now := time.Now()

	// 更新状态
	m.appState.Status = models.AppStatusStarting
	m.appState.PID = &pid
	m.appState.StartTime = &now
	m.appState.Config = config
	m.appState.LastError = nil
	m.lastProfile = profile

	// 启动健康检查
	m.startHealthCheck()

	m.logger.Infof("Started app %s with PID %d", m.appInfo.Name, pid)

	return &models.StartAppResponse{
		Status:  "started",
		AppName: m.appInfo.Name,
		PID:     pid,
		Profile: profile,
	}, nil
}

// StopApp 停止应用（互斥、幂等、状态检查）
func (m *Manager) StopApp() (*models.StopAppResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.cmd == nil || m.appState.Status == models.AppStatusStopped {
		return &models.StopAppResponse{Status: "stopped"}, nil
	}

	m.appState.Status = models.AppStatusStopping

	// 发送SIGTERM
	if err := m.cmd.Process.Signal(os.Interrupt); err != nil {
		m.logger.Errorf("Failed to send SIGTERM: %v", err)
	}

	// 等待进程结束
	done := make(chan error, 1)
	go func() {
		done <- m.cmd.Wait()
	}()

	select {
	case <-done:
		// 进程正常结束
	case <-time.After(10 * time.Second):
		// 强制杀死
		if err := m.cmd.Process.Kill(); err != nil {
			m.logger.Errorf("Failed to kill process: %v", err)
		}
		<-done
	}

	now := time.Now()
	m.appState.Status = models.AppStatusStopped
	m.appState.StopTime = &now
	m.appState.PID = nil

	// 停止健康检查
	m.stopHealthCheck()

	// 停止后清空内部状态
	m.internalStatus = nil

	m.logger.Infof("Stopped app %s", m.appInfo.Name)

	return &models.StopAppResponse{Status: "stopped"}, nil
}

// RestartApp 重启应用（互斥、幂等、状态检查）
func (m *Manager) RestartApp() (*models.RestartAppResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.appInfo == nil {
		return nil, fmt.Errorf("app not configured")
	}

	profile := m.lastProfile
	if profile == "" {
		profile = "{}"
	}
	var config map[string]interface{}
	if err := json.Unmarshal([]byte(profile), &config); err != nil {
		return nil, fmt.Errorf("invalid JSON profile: %v", err)
	}

	// 如果应用正在运行，先停止
	if m.cmd != nil && m.appState.Status != models.AppStatusStopped {
		m.logger.Infof("Stopping app %s for restart", m.appInfo.Name)
		if err := m.cmd.Process.Signal(os.Interrupt); err != nil {
			m.logger.Errorf("Failed to send SIGTERM: %v", err)
		}
		done := make(chan error, 1)
		go func() {
			done <- m.cmd.Wait()
		}()
		select {
		case <-done:
			// 进程正常结束
		case <-time.After(10 * time.Second):
			if err := m.cmd.Process.Kill(); err != nil {
				m.logger.Errorf("Failed to kill process: %v", err)
			}
			<-done
		}
		m.stopHealthCheck()
	}

	// 自动补全 -id 参数
	args := m.appInfo.Args
	idPresent := false
	for i, arg := range args {
		if arg == "-id" && i+1 < len(args) {
			args[i+1] = m.proxyID
			idPresent = true
		}
	}
	if !idPresent {
		args = append(args, "-id", m.proxyID)
	}

	// 设置环境变量
	env := os.Environ()
	for k, v := range m.appInfo.Env {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	env = append(env, fmt.Sprintf("APP_PROFILE=%s", profile))
	env = append(env, fmt.Sprintf("APP_NAME=%s", m.appInfo.Name))
	env = append(env, fmt.Sprintf("PROXY_GRPC_PORT=%s", os.Getenv("GRPC_PORT")))

	cmd := exec.CommandContext(context.Background(), m.appInfo.Command, args...)
	cmd.Env = env
	// 不再设置 cmd.Dir

	// 启动进程
	if err := cmd.Start(); err != nil {
		m.appState.Status = models.AppStatusError
		errorMsg := err.Error()
		m.appState.LastError = &errorMsg
		m.logger.Errorf("Failed to restart app %s: %v", m.appInfo.Name, err)
		return nil, err
	}

	m.cmd = cmd
	pid := cmd.Process.Pid
	now := time.Now()

	// 更新状态
	m.appState.Status = models.AppStatusStarting
	m.appState.PID = &pid
	m.appState.StartTime = &now
	m.appState.Config = config
	m.appState.LastError = nil
	m.appState.RestartCount++

	// 启动健康检查
	m.startHealthCheck()

	m.logger.Infof("Restarted app %s with PID %d", m.appInfo.Name, pid)

	return &models.RestartAppResponse{
		Status:  "restarted",
		AppName: m.appInfo.Name,
		PID:     pid,
		Profile: profile,
	}, nil
}

// GetStatus 获取应用状态
func (m *Manager) GetStatus() *models.AppStatusResponse {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.appInfo == nil {
		// 优先用 manifest 信息
		manifestPath := "/app/manifest.json"
		if data, err := ioutil.ReadFile(manifestPath); err == nil {
			var manifest struct {
				AppName             string   `json:"app_name"`
				HealthCheckInterval int      `json:"health_check_interval"`
				DefaultArgs         []string `json:"default_args"`
			}
			if err := json.Unmarshal(data, &manifest); err == nil {
				return &models.AppStatusResponse{
					AppName: manifest.AppName,
					Status: "ready",
					RestartCount: 0,
				}
			}
		}
		// 其次用 APP_NAME 环境变量
		appName := os.Getenv("APP_NAME")
		if appName == "" {
			// 尝试用 hostname 作为容器名
			hostname, err := os.Hostname()
			if err == nil && hostname != "" {
				appName = hostname
			} else {
				appName = "cleaner"
			}
		}
		return &models.AppStatusResponse{
			AppName: appName,
			Status: "ready",
			RestartCount: 0,
		}
	}

	// 检查进程状态
	if m.cmd != nil && m.cmd.Process != nil {
		if m.appState.Status == models.AppStatusRunning {
			// 检查进程是否还在运行
			if m.cmd.ProcessState != nil && m.cmd.ProcessState.Exited() {
				m.appState.Status = models.AppStatusStopped
				now := time.Now()
				m.appState.StopTime = &now
			}
		}
	}

	return &models.AppStatusResponse{
		AppName:     m.appInfo.Name,
		Status:      string(m.appState.Status),
		PID:         m.appState.PID,
		StartTime:   m.appState.StartTime,
		StopTime:    m.appState.StopTime,
		RestartCount: m.appState.RestartCount,
		LastError:   m.appState.LastError,
		Config:      m.appState.Config,
	}
}

// UpdateInternalStatus 更新app内部状态
func (m *Manager) UpdateInternalStatus(status map[string]interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.internalStatus = status
}

// GetInternalStatus 获取app内部状态
func (m *Manager) GetInternalStatus() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.internalStatus
}

// startHealthCheck 启动健康检查
func (m *Manager) startHealthCheck() {
	if m.healthTicker != nil {
		m.healthTicker.Stop()
	}

	m.healthTicker = time.NewTicker(time.Duration(m.appInfo.HealthCheckInterval) * time.Second)
	go func() {
		for range m.healthTicker.C {
			m.checkHealth()
		}
	}()
}

// stopHealthCheck 停止健康检查
func (m *Manager) stopHealthCheck() {
	if m.healthTicker != nil {
		m.healthTicker.Stop()
		m.healthTicker = nil
	}
}

// checkHealth 健康检查
func (m *Manager) checkHealth() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.cmd == nil || m.cmd.Process == nil {
		return
	}

	// 检查进程是否还在运行
	if m.cmd.ProcessState != nil && m.cmd.ProcessState.Exited() {
		if m.appInfo.AutoRestart && m.appState.RestartCount < m.appInfo.MaxRestarts {
			m.logger.Infof("App %s crashed, restarting... (attempt %d/%d)", m.appInfo.Name, m.appState.RestartCount+1, m.appInfo.MaxRestarts)
			go m.restartApp()
		} else {
			m.appState.Status = models.AppStatusStopped
			now := time.Now()
			m.appState.StopTime = &now
			m.logger.Errorf("App %s crashed and max restarts reached (%d/%d)", m.appInfo.Name, m.appState.RestartCount, m.appInfo.MaxRestarts)
		}
	} else if m.appState.Status == models.AppStatusStarting {
		m.appState.Status = models.AppStatusRunning
		m.logger.Infof("App %s is now running", m.appInfo.Name)
	}
}

// restartApp 重启应用
func (m *Manager) restartApp() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.appState.RestartCount++
	m.appState.Status = models.AppStatusStarting

	// 启动新进程
	args := append([]string{m.appInfo.Command}, m.appInfo.Args...)
	cmd := exec.CommandContext(context.Background(), args[0], args[1:]...)
	cmd.Env = os.Environ()
	for k, v := range m.appInfo.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
	// 不再设置 cmd.Dir

	if err := cmd.Start(); err != nil {
		m.appState.Status = models.AppStatusError
		errorMsg := err.Error()
		m.appState.LastError = &errorMsg
		m.logger.Errorf("Failed to restart app: %v", err)
		return
	}

	m.cmd = cmd
	pid := cmd.Process.Pid
	now := time.Now()
	m.appState.PID = &pid
	m.appState.StartTime = &now
	m.appState.LastError = nil

	m.logger.Infof("Restarted app %s with PID %d", m.appInfo.Name, pid)
}

func (m *Manager) ProxyID() string {
	return m.proxyID
} 