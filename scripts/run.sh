#!/bin/bash

# 运行脚本：启动所有容器或指定容器

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
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

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# 检查镜像是否存在
check_image() {
    local image=$1
    if ! docker image inspect $image >/dev/null 2>&1; then
        log_error "Image $image not found. Please build it first:"
        echo "  ./scripts/build.sh examples"
        return 1
    fi
    return 0
}

# 停止并删除容器
stop_container() {
    local container_name=$1
    if docker ps -q -f name=$container_name | grep -q .; then
        log_info "Stopping container $container_name..."
        docker stop $container_name 2>/dev/null || true
        docker rm $container_name 2>/dev/null || true
    fi
}

# 运行容器
run_container() {
    local image=$1
    local container_name=$2
    local port=$3
    local enable_traefik=${4:-false}
    
    log_info "Starting $container_name on port $port..."
    
    # 停止已存在的容器
    stop_container $container_name
    
    # 运行新容器，参数用单横线
    if docker run -d \
      --name $container_name \
      -p $port:$port \
      ${enable_traefik:+
      --label "traefik.enable=true" \
      --label "traefik.http.routers.$container_name.rule=PathPrefix(\"/$container_name\")" \
      --label "traefik.http.services.$container_name.loadbalancer.server.port=$port" \
      --label "traefik.http.middlewares.strip-$container_name.stripprefix.prefixes=/$container_name" \
      --label "traefik.http.routers.$container_name.middlewares=strip-$container_name@docker" \
      } \
      $image -id $container_name -http-port $port; then
        log_success "$container_name started successfully on port $port"
        return 0
    else
        log_error "Failed to start $container_name"
        return 1
    fi
}

# 等待服务启动
wait_for_service() {
    local port=$1
    local service=$2
    local max_attempts=30
    local attempt=1
    
    log_info "Waiting for $service to start on port $port..."
    
    while [ $attempt -le $max_attempts ]; do
        if curl -f http://localhost:$port/health >/dev/null 2>&1; then
            log_success "$service is ready on port $port"
            return 0
        fi
        
        echo -n "."
        sleep 2
        attempt=$((attempt + 1))
    done
    
    log_error "$service failed to start on port $port"
    return 1
}

# 主函数
main() {
    echo "=========================================="
    echo "Brick Smart Template - 运行脚本"
    echo "=========================================="
    
    # 检查Docker
    if ! command -v docker &> /dev/null; then
        log_error "Docker is not installed or not in PATH"
        exit 1
    fi
    
    # 检查Docker是否运行
    if ! docker info &> /dev/null; then
        log_error "Docker is not running"
        exit 1
    fi
    
    log_info "Starting containers..."

    log_info "Starting application containers..."
    
    # 检查镜像
    check_image "brick-smart-cleaner" || exit 1
    #check_image "brick-smart-thermostat" || exit 1
    #check_image "brick-smart-lighting" || exit 1
    
    # 运行所有容器
    run_container "brick-smart-cleaner" "brick-cleaner" "17101" true
    # 暂时禁用其他容器
    # run_container "brick-smart-thermostat" "brick-thermostat" "17103" false
    # run_container "brick-smart-lighting" "brick-lighting" "17102" false
    
    # 等待服务启动
    log_info "Waiting for all services to start..."
    wait_for_service 17101 "cleaner"
    # wait_for_service 17102 "lighting"
    # wait_for_service 17103 "thermostat"
    
    echo ""
    echo "=========================================="
    log_success "All containers started successfully!"
    echo "=========================================="
    echo ""
    echo "Running containers:"
    echo "  - brick-cleaner     (port 17101)"
    # echo "  - brick-thermostat  (port 17103)"
    # echo "  - brick-lighting    (port 17102)"
    echo ""
    echo "Service URLs:"
    echo "  Cleaner:     http://localhost:17101"
    # echo "  Lighting:    http://localhost:17102"
    # echo "  Thermostat:  http://localhost:17103"
    echo ""
    echo "To check status:"
    echo "  curl http://localhost:17101/app/status"
    # echo "  curl http://localhost:17102/app/status"
    # echo "  curl http://localhost:17103/app/status"
    echo ""
    echo "To stop all containers:"
    echo "  ./scripts/stop-all.sh"
    echo ""
    echo "To view logs:"
    echo "  docker logs brick-cleaner"
    echo "  docker logs brick-thermostat"
    echo "  docker logs brick-lighting"
}

# 显示帮助
show_help() {
    echo "Usage: $0 [all|examples|cleaner|thermostat|lighting|help]"
    echo ""
    echo "Options:"
    echo "  all/examples   Run all containers (default)"
    echo "  cleaner        Run only cleaner container"
    echo "  thermostat     Run only thermostat container"
    echo "  lighting       Run only lighting container"
    echo "  help, -h       Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                # Run all containers"
    echo "  $0 all            # Run all containers"
    echo "  $0 examples       # Run all containers"
    echo "  $0 cleaner        # Run only cleaner"
    echo "  $0 thermostat     # Run only thermostat"
    echo "  $0 lighting       # Run only lighting"
}

# 解析参数
case "${1:-}" in
    help|-h|--help)
        show_help
        exit 0
        ;;
    all|examples|"")
        main
        ;;
    cleaner)
        check_image "brick-smart-cleaner" || exit 1
        run_container "brick-smart-cleaner" "brick-cleaner" "17101"
        wait_for_service 17101 "cleaner"
        echo ""
        log_success "Cleaner container started on port 17101"
        exit 0
        ;;
    thermostat)
        # check_image "brick-smart-thermostat" || exit 1
        # run_container "brick-smart-thermostat" "brick-thermostat" "17103" false
        # wait_for_service 17103 "thermostat"
        # echo ""
        # log_success "Thermostat container started on port 17103"
        # exit 0
        ;;
    lighting)
        # check_image "brick-smart-lighting" || exit 1
        # run_container "brick-smart-lighting" "brick-lighting" "17102" false
        # wait_for_service 17102 "lighting"
        # echo ""
        # log_success "Lighting container started on port 17102"
        # exit 0
        ;;
    *)
        log_error "Unknown option: $1"
        show_help
        exit 1
        ;;
esac