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

	"brick-smart-template/examples/brick-smart-thermostat/pkg/thermostat"
	"brick-smart-template/examples/brick-smart-thermostat/pkg/httpclient"
)

func main() {
	var (
		id       = flag.String("id", "", "Thermostat ID")
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
	thermo := thermostat.NewThermostat(*id, httpClient)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		thermo.Start(ctx)
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Shutting down thermostat...")
	cancel()
	wg.Wait()
	log.Println("Thermostat stopped")
}

func showHelp() {
	fmt.Println("Brick Smart Thermostat - A smart thermostat simulation")
	fmt.Println("Usage: ./thermostat [options]")
	fmt.Println("  -id ID              Thermostat ID (required)")
	fmt.Println("  -http-port PORT     HTTP API server port (default: 17100)")
	fmt.Println("  -help               Show this help message")
} 