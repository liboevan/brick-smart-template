package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"brick-smart-template/examples/brick-smart-cleaner/pkg/cleaner"
	"brick-smart-template/examples/brick-smart-cleaner/pkg/httpclient"
)

func main() {
	// 解析命令行参数
	var (
		id       = flag.String("id", "", "Cleaner ID")
		httpPort = flag.String("http-port", "17100", "HTTP API server port")
		help     = flag.Bool("help", false, "Show help information")
	)
	flag.Parse()

	// 显示帮助信息
	if *help {
		showHelp()
		os.Exit(0)
	}

	// 验证必需参数
	if *id == "" {
		log.Fatal("Error: -id parameter is required")
	}

	// 创建上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 创建HTTP客户端
	httpClient := httpclient.NewClient(*httpPort)

	// 创建扫地机实例
	cleanerInstance := cleaner.NewCleaner(*id, httpClient)

	// 启动清理任务
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		cleanerInstance.StartCleaning(ctx)
	}()

	// 等待中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down cleaner...")
	cancel()
	wg.Wait()
	log.Println("Cleaner stopped")
}

func showHelp() {
	fmt.Println("Brick Smart Cleaner - A robotic vacuum cleaner simulation")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  ./cleaner [options]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -id ID              Cleaner ID (required)")
	fmt.Println("  -http-port PORT     HTTP API server port (default: 17100)")
	fmt.Println("  -help               Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  ./cleaner -id cleaner-001 -http-port 17100")
	fmt.Println("  ./cleaner -id cleaner-002 -http-port 17200")
	fmt.Println()
	fmt.Println("Behavior:")
	fmt.Println("  - Continuously cleans 6 rooms in a cycle")
	fmt.Println("  - Each room takes 1 minute to clean")
	fmt.Println("  - Reports progress via HTTP API")
	fmt.Println("  - Shows cleaning progress as percentage")
} 