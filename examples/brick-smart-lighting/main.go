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

	"brick-smart-template/examples/brick-smart-lighting/pkg/lighting"
	"brick-smart-template/examples/brick-smart-lighting/pkg/httpclient"
)

func main() {
	var (
		id       = flag.String("id", "", "Lighting ID")
		httpPort = flag.String("http-port", "17100", "HTTP API server port")
		help     = flag.Bool("help", false, "Show help information")
	)
	flag.Parse()

	if *help {
		showHelp()
		os.Exit(0)
	}
	if *id == "" {
		log.Fatal("Error: -id parameter is required")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	httpClient := httpclient.NewClient(*httpPort)
	light := lighting.NewLighting(*id, httpClient)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		light.Start(ctx)
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Shutting down lighting...")
	cancel()
	wg.Wait()
	log.Println("Lighting stopped")
}

func showHelp() {
	fmt.Println("Brick Smart Lighting - A smart lighting simulation")
	fmt.Println("Usage: ./lighting [options]")
	fmt.Println("  -id ID              Lighting ID (required)")
	fmt.Println("  -http-port PORT     HTTP API server port (default: 17100)")
	fmt.Println("  -help               Show this help message")
} 