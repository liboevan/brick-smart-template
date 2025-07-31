#!/bin/bash

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

stop_container() {
    local container_name=$1
    if docker ps -q -f name=$container_name | grep -q .; then
        log_info "Stopping container $container_name..."
        docker stop $container_name 2>/dev/null || true
        docker rm $container_name 2>/dev/null || true
    fi
}

log_info "Starting Traefik service discovery..."
stop_container "traefik"
if docker run -d \
  --name traefik \
  -p 17110:17110 \
  -p 17111:17111 \
  -v /var/run/docker.sock:/var/run/docker.sock \
  traefik:v2.9 \
  --providers.docker=true \
  --providers.docker.exposedbydefault=false \
  --api.insecure=true \
  --api.dashboard=true \
  --entrypoints.traefik.address=:17110 \
  --entrypoints.web.address=:17111; then
    log_success "Traefik started successfully - Management UI: 17111, Web entrypoint: 17110"
    exit 0
else
    log_error "Failed to start Traefik"
    exit 1
fi