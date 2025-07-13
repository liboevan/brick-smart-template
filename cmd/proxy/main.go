package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"brick-smart-template/pkg/appmanager"
	"brick-smart-template/pkg/httpapi"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	// 解析命令行参数
	var (
		httpPort = flag.String("http-port", "", "HTTP API server port (e.g., 8000)")
		grpcPort = flag.String("grpc-port", "", "gRPC server port (e.g., 50051)")
		configFile = flag.String("config", "", "Configuration file path")
		id = flag.String("id", "", "Proxy/App ID (用于 app 启动和校验)")
		help = flag.Bool("help", false, "Show help information")
	)
	flag.Parse()

	// 显示帮助信息
	if *help {
		showHelp()
		os.Exit(0)
	}

	// 初始化日志
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	// 加载配置
	loadConfig()

	// 应用命令行参数
	if *httpPort != "" {
		viper.Set("http.addr", ":"+*httpPort)
	}
	if *grpcPort != "" {
		viper.Set("grpc.addr", ":"+*grpcPort)
	}
	if *configFile != "" {
		viper.SetConfigFile(*configFile)
	}

	proxyID := *id
	if proxyID == "" {
		proxyID = os.Getenv("PROXY_ID")
	}
	if proxyID == "" {
		proxyID = "default-proxy"
	}

	// 创建应用管理器
	manager := appmanager.NewManager(logger, proxyID)

	// 创建HTTP服务器
	httpServer := httpapi.NewServer(manager, logger)

	// 启动HTTP服务器
	go func() {
		httpAddr := viper.GetString("http.addr")
		if httpAddr == "" {
			httpAddr = ":8000"
		}
		
		logger.Infof("Starting HTTP server on %s", httpAddr)
		if err := httpServer.Run(httpAddr); err != nil {
			logger.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// 等待中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.Info("Shutting down...")

	logger.Info("Proxy stopped")
}

// showHelp 显示帮助信息
func showHelp() {
	fmt.Println("Brick Smart Template App Proxy")
	fmt.Println("A proxy service for managing application containers")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  ./app-proxy [options]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -http-port PORT     HTTP API server port (e.g., 8000)")
	fmt.Println("  -grpc-port PORT     gRPC server port (e.g., 50051)")
	fmt.Println("  -config FILE        Configuration file path")
	fmt.Println("  -help               Show this help message")
	fmt.Println()
	fmt.Println("Environment Variables:")
	fmt.Println("  PROXY_HTTP_ADDR     HTTP server address (default: :8000)")
	fmt.Println("  PROXY_GRPC_ADDR     gRPC server address (default: :50051)")
	fmt.Println("  PROXY_LOG_LEVEL     Log level (default: info)")
	fmt.Println("  PROXY_SHUTDOWN_TIMEOUT  Shutdown timeout (default: 30s)")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  ./app-proxy -http-port 8080")
	fmt.Println("  ./app-proxy -http-port 9000 -grpc-port 50052")
	fmt.Println("  ./app-proxy -config /path/to/config.yaml")
	fmt.Println("  PROXY_HTTP_ADDR=:8080 ./app-proxy")
}

func loadConfig() {
	// 设置默认值
	viper.SetDefault("http.addr", ":8000")
	viper.SetDefault("grpc.addr", ":50051")
	viper.SetDefault("shutdown.timeout", "30s")
	viper.SetDefault("log.level", "info")

	// 从环境变量读取
	viper.SetEnvPrefix("PROXY")
	viper.AutomaticEnv()

	// 从配置文件读取
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		// 配置文件不存在，使用默认值
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(err)
		}
	}
} 