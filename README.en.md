# Brick Smart Template

> 中文版请见: [README.md](README.md)

This is a project for concept and technical validation and practice, using the construction and management of smart device proxies and example applications as a case, and is not intended for any real-world production use.

## Project Overview

Brick Smart Template provides unified build, run, test, and clean scripts for a smart device proxy (proxy) and multiple example devices (cleaner, lighting, thermostat, etc.), suitable for IoT, smart home prototyping, and educational demos.

- **One-click build & run**: Quickly build, start, and clean all services via Makefile and scripts.
- **Multiple device examples**: Built-in simulation for cleaner, lighting, thermostat, and more.
- **Docker support**: All services can be containerized for consistent local/cloud environments.
- **Comprehensive docs**: Bilingual documentation, including device and test guides.

## Directory Structure

```
brick-smart-template/
├── Makefile                # Build entry
├── Dockerfile              # Proxy image build file
├── docker-compose.yml      # Multi-service orchestration
├── scripts/                # Build/run/test/clean scripts
├── examples/               # Example device code
├── docs/                   # Detailed docs (device, test, etc.)
└── README.md/README.en.md  # Project overview (CN/EN)
```

## Quick Start

```bash
make build      # Build all images
make run        # Start all containers
make test       # Run all API tests
make clean      # Clean all containers and images
```

## Key Documentation

- [Device Doc: cleaner](docs/cleaner.md)
- [Test Guide (EN)](docs/test.md)
- [中文说明](README.md)

## Use Cases

- IoT device development & testing
- Smart home system prototyping
- Education & demonstration
- Unified management & automation for multiple devices

For detailed test methods and API flows, see docs/test.md. 