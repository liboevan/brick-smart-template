.PHONY: build run clean gen-go-sum help

.DEFAULT_GOAL := help

build:
	@echo "Building all images..."
	./scripts/build.sh all

run:
	@echo "Running all containers..."
	./scripts/run.sh all

clean:
	@echo "Cleaning all containers and images..."
	./scripts/clean.sh all --image

gen-go-sum:
	@echo "Generating go.sum..."
	./scripts/gen-go-sum.sh

help:
	@echo "Available targets:"
	@echo "  build        - Build all docker images (proxy + examples)"
	@echo "  run          - Run all containers (proxy + examples)"
	@echo "  clean        - Clean all containers and images"
	@echo "  gen-go-sum   - Generate go.sum using Docker"
	@echo "  help         - Show this help message"
