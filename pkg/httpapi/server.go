package httpapi

import (
	"net/http"

	"brick-smart-template/pkg/appmanager"
	"brick-smart-template/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Server HTTP API服务器
type Server struct {
	router  *gin.Engine
	manager *appmanager.Manager
	logger  *logrus.Logger
}

// NewServer 创建新的HTTP服务器
func NewServer(manager *appmanager.Manager, logger *logrus.Logger) *Server {
	server := &Server{
		router:  gin.Default(),
		manager: manager,
		logger:  logger,
	}

	server.setupRoutes()
	return server
}

// setupRoutes 设置路由
func (server *Server) setupRoutes() {
	// 健康检查
	server.router.GET("/health", server.healthCheck)

	// 应用管理API
	appGroup := server.router.Group("/app")
	{
		appGroup.POST("/configure", server.configureApp)
		appGroup.POST("/start", server.startApp)
		appGroup.POST("/restart", server.restartApp)
		appGroup.POST("/stop", server.stopApp)
		appGroup.GET("/status", server.getAppStatus)
		appGroup.GET("/data", server.getInternalStatus)
		appGroup.GET("/process", server.getProcessStatus)
	}

	// 状态报告API (用于gRPC的替代)
	server.router.POST("/app/status/report", server.reportStatus)
}

// healthCheck 健康检查
func (server *Server) healthCheck(c *gin.Context) {
	status := server.manager.GetStatus()
	
	response := models.HealthCheckResponse{
		Status:      "healthy",
		ProxyStatus: "running",
		AppStatus:   status.Status,
	}

	c.JSON(http.StatusOK, response)
}

// configureApp 配置应用
func (server *Server) configureApp(c *gin.Context) {
	var request models.ConfigureAppRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		server.logger.Errorf("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := server.manager.ConfigureApp(request.AppInfo); err != nil {
		server.logger.Errorf("Failed to configure app: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := models.ConfigureAppResponse{
		Status:  "configured",
		AppName: request.AppInfo.Name,
	}

	c.JSON(http.StatusOK, response)
}

// startApp 启动应用
func (server *Server) startApp(c *gin.Context) {
	var request models.StartAppRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		server.logger.Errorf("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 校验 app_name 和 id
	proxyID := server.manager.ProxyID()
	proxyAppName := ""
	if server.manager.GetStatus() != nil {
		proxyAppName = server.manager.GetStatus().AppName
	}
	if request.ID != "" && request.ID != proxyID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id mismatch: expected " + proxyID})
		return
	}
	if request.AppName != "" && proxyAppName != "" && request.AppName != proxyAppName {
		c.JSON(http.StatusBadRequest, gin.H{"error": "app_name mismatch: expected " + proxyAppName})
		return
	}

	response, err := server.manager.StartApp(request.Profile)
	if err != nil {
		server.logger.Errorf("Failed to start app: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// restartApp 重启应用
func (server *Server) restartApp(c *gin.Context) {
	response, err := server.manager.RestartApp()
	if err != nil {
		server.logger.Errorf("Failed to restart app: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

// stopApp 停止应用
func (server *Server) stopApp(c *gin.Context) {
	response, err := server.manager.StopApp()
	if err != nil {
		server.logger.Errorf("Failed to stop app: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// getAppStatus 获取应用状态
func (server *Server) getAppStatus(c *gin.Context) {
	status := server.manager.GetStatus()
	processID := server.manager.ProxyID()
	c.JSON(http.StatusOK, gin.H{
		"process_id": processID,
		"status": status,
	})
}

// getInternalStatus 获取应用内部状态
func (server *Server) getInternalStatus(c *gin.Context) {
	internalStatus := server.manager.GetInternalStatus()
	c.JSON(http.StatusOK, internalStatus)
}

// getProcessStatus 获取合并后的状态和数据
func (server *Server) getProcessStatus(c *gin.Context) {
    status := server.manager.GetStatus()
    data := server.manager.GetInternalStatus()
    processID := server.manager.ProxyID()
    c.JSON(http.StatusOK, gin.H{
        "app_name": status.AppName,
        "process_id": processID,
        "process_status": status,
        "process_data": data,
    })
}

// reportStatus 报告状态 (gRPC的HTTP替代)
func (server *Server) reportStatus(c *gin.Context) {
	var request struct {
		Status string                 `json:"status"`
		Data   map[string]interface{} `json:"data"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		server.logger.Errorf("Invalid status report: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 只更新内部状态
	// 同时更新内部状态
	server.manager.UpdateInternalStatus(request.Data)
	
	server.logger.Infof("Received status report: %s", request.Status)

	c.JSON(http.StatusOK, gin.H{"status": "received"})
}

// Run 启动HTTP服务器
func (server *Server) Run(addr string) error {
	server.logger.Infof("Starting HTTP server on %s", addr)
	return server.router.Run(addr)
} 