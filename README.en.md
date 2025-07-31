# Brick Smart Template

> Chinese version: [README.md](README.md)

This is a template project for building and managing smart device proxies and example applications, focusing on rapid prototyping and educational demonstrations for IoT and smart home scenarios. This project does not include production environment functionality and is intended for proof-of-concept and technical practice only.

## Project Overview

Brick Smart Template provides unified build, run, test, and clean scripts for a smart device proxy and multiple example devices (cleaner, lighting, thermostat, etc.).

- **One-click build & run**: Quickly build, start, and clean all services via Makefile and scripts
- **Multiple device examples**: Built-in simulation for cleaner, lighting, thermostat, and other smart devices
- **Docker support**: All services can be containerized to ensure environment consistency

## System Architecture

### Core Components
- **brick-proxy**: Smart device proxy that manages device connections and communication
- **Device examples**: Simulated implementations of various smart devices (cleaner, lighting, etc.)
- **Traefik reverse proxy**: Centralized service access management

## Directory Structure

```
brick-smart-template/
├── Makefile                # Build entry
├── Dockerfile              # Proxy image build file
├── docker-compose.yml      # Multi-service orchestration
├── scripts/                # Build/run/test/clean scripts
├── examples/               # Example device code
│   ├── brick-smart-cleaner/
│   ├── brick-smart-lighting/
│   └── brick-smart-thermostat/
├── docs/                   # Detailed documentation
├── README.md               # Project overview (Chinese)
└── README.en.md            # Project overview (English)
```

## Quick Start

```bash
make build      # Build all images
make run        # Start all containers
make test       # Run all API tests
make clean      # Clean all containers and images
```

## Key Documentation

Detailed documentation for device examples can be found in the README of each example project:
- [Cleaner Example](examples/brick-smart-cleaner/README.md)
- [Lighting Example](examples/brick-smart-lighting/README.md)
- [Thermostat Example](examples/brick-smart-thermostat/README.md)

## Use Cases

- IoT device development & testing
- Smart home system prototyping
- Education & demonstration
- Unified management & automation for multiple devices

## Traefik Reverse Proxy Configuration

This project uses Traefik as a reverse proxy to centrally manage access to all services.

### Key Features
- Automatic Docker container service discovery
- Service differentiation through path prefixes
- Automatic path prefix stripping without modifying service code

### Service Access Method
All services are accessed through Traefik's port 17111 using the following URL format:
```
http://localhost:17111/<service-name>/<api-endpoint>
```

### Example
Accessing the cleaner service's health API:
```
http://localhost:17111/brick-cleaner/health
```

### Configuration Notes
- Traefik configuration file: `scripts/start_traefik.sh`
- Service registration is implemented through Docker labels, see the `run_container` function in `scripts/run.sh`