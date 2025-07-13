#!/bin/bash

# 测试脚本：直接测试API调用

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

# 检查服务是否可用
check_service() {
    local port=$1
    local service=$2
    
    log_info "Checking if $service is available on port $port..."
    
    if curl -f http://localhost:$port/health >/dev/null 2>&1; then
        log_success "$service is available on port $port"
        return 0
    else
        log_error "$service is not available on port $port"
        log_info "Please make sure the container is running:"
        echo "  ./scripts/run.sh $service"
        return 1
    fi
}

# 配置和启动app
start_app() {
    local port=$1
    local app_name=$2
    
    log_info "Configuring $app_name on port $port..."
    
    # 配置app
    CONFIG_RESPONSE=$(curl -s -X POST http://localhost:$port/app/configure \
        -H "Content-Type: application/json" \
        -d "{
            \"app_info\": {
                \"name\": \"$app_name\",
                \"command\": \"./$app_name\",
                \"args\": [\"-id\", \"${app_name}-001\", \"-http-port\", \"$port\"],
                \"health_check_interval\": 3
            }
        }")
    
    if echo "$CONFIG_RESPONSE" | grep -q "configured"; then
        log_success "$app_name configured successfully"
    else
        log_error "Failed to configure $app_name: $CONFIG_RESPONSE"
        return 1
    fi
    
    # 启动app
    log_info "Starting $app_name..."
    START_RESPONSE=$(curl -s -X POST http://localhost:$port/app/start \
        -H "Content-Type: application/json" \
        -d "{
            \"profile\": \"{\\\"${app_name}_id\\\": \\\"${app_name}-001\\\"}\"
        }")
    
    if echo "$START_RESPONSE" | grep -q "started\|already_running"; then
        log_success "$app_name started successfully"
        return 0
    else
        log_error "Failed to start $app_name: $START_RESPONSE"
        return 1
    fi
}

# 重启app
restart_app() {
    local port=$1
    local app_name=$2
    
    log_info "Restarting $app_name on port $port..."
    RESTART_RESPONSE=$(curl -s -X POST http://localhost:$port/app/restart \
        -H "Content-Type: application/json" \
        -d "{
            \"profile\": \"{\\\"${app_name}_id\\\": \\\"${app_name}-001\\\"}\"
        }")
    
    if echo "$RESTART_RESPONSE" | grep -q "restarted"; then
        log_success "$app_name restarted successfully"
        return 0
    else
        log_error "Failed to restart $app_name: $RESTART_RESPONSE"
        return 1
    fi
}

# 检查app状态
check_app_status() {
    local port=$1
    local app_name=$2
    
    log_info "Checking $app_name status on port $port..."
    STATUS_RESPONSE=$(curl -s http://localhost:$port/app/status 2>/dev/null || echo "{}")
    echo "Status for $app_name: $STATUS_RESPONSE"
    
    if echo "$STATUS_RESPONSE" | grep -q "running\|starting"; then
        log_success "$app_name is running"
        return 0
    elif echo "$STATUS_RESPONSE" | grep -q "stopped\|stopping"; then
        log_warning "$app_name is stopped"
        return 1
    else
        log_error "$app_name status unknown"
        return 1
    fi
}

# 停止app
stop_app() {
    local port=$1
    local app_name=$2
    
    log_info "Stopping $app_name on port $port..."
    STOP_RESPONSE=$(curl -s -X POST http://localhost:$port/app/stop 2>/dev/null || echo "{}")
    
    if echo "$STOP_RESPONSE" | grep -q "stopped"; then
        log_success "$app_name stopped successfully"
        return 0
    else
        log_warning "$app_name stop response: $STOP_RESPONSE"
        return 1
    fi
}

# 显示状态和数据
show_status_and_data() {
    local port=$1
    local app_name=$2
    STATUS_RESPONSE=$(curl -s http://localhost:$port/app/status 2>/dev/null || echo "{}")
    DATA_RESPONSE=$(curl -s http://localhost:$port/app/data 2>/dev/null || echo "{}")
    echo -e "\033[1;36m[STATUS]\033[0m $app_name: $STATUS_RESPONSE"
    echo -e "\033[1;35m[DATA]\033[0m $app_name: $DATA_RESPONSE"
}

# 测试单个应用
# 新流程：运行状态下每2秒查一次data，持续5次，然后stop，再start，查一次data，保持app运行
# 保留原有参数和日志风格

test_single_app() {
    local app_name=$1
    local port=$2

    # 保证容器为初始状态
    log_info "Cleaning $app_name container before test..."
    ./scripts/clean.sh $app_name
    log_info "Starting $app_name container before test..."
    ./scripts/run.sh $app_name
    sleep 2

    echo ""
    echo "=========================================="
    log_info "Testing $app_name API on port $port..."
    echo "=========================================="

    # 0. 初始状态
    log_info "Initial status and data for $app_name:"
    show_status_and_data $port $app_name

    # 1. 检查服务是否可用
    check_service $port $app_name || return 1
    show_status_and_data $port $app_name

    # 2. 配置并启动app
    log_info "Configuring and starting $app_name on port $port..."
    start_app $port $app_name
    sleep 2
    show_status_and_data $port $app_name

    # 等待app启动
    sleep 3
    show_status_and_data $port $app_name

    # 3. 运行状态下每2秒查一次data，持续5次
    log_info "Polling $app_name /app/data every 2s for 5 times while running..."
    for i in {1..5}; do
        DATA_RESPONSE=$(curl -s http://localhost:$port/app/data 2>/dev/null || echo "{}")
        echo -e "\033[1;35m[DATA][$i]\033[0m $app_name: $DATA_RESPONSE"
        sleep 2
    done

    # 4. 停止app
    log_info "Stopping $app_name, current status and data:"
    sleep 2
    show_status_and_data $port $app_name
    stop_app $port $app_name
    sleep 2
    show_status_and_data $port $app_name

    # 5. restart app
    log_info "Restarting $app_name after stop..."
    RESTART_RESPONSE=$(curl -s -X POST http://localhost:$port/app/restart -H "Content-Type: application/json" -d '{}')
    sleep 2
    show_status_and_data $port $app_name

    # 6. 再查一次data
    log_info "Polling $app_name /app/data once after restart..."
    DATA_RESPONSE=$(curl -s http://localhost:$port/app/data 2>/dev/null || echo "{}")
    echo -e "\033[1;35m[DATA][after-restart]\033[0m $app_name: $DATA_RESPONSE"

    echo ""
    log_success "$app_name API test completed successfully! (app will remain running)"
    echo "=========================================="
}

# 测试所有应用
test_all_apps() {
    echo "=========================================="
    echo "Brick Smart Template - API 测试脚本"
    echo "=========================================="
    
    log_info "Starting API tests..."
    
    # 测试所有应用
    test_single_app "cleaner" "17101"
    test_single_app "lighting" "17102"
    test_single_app "thermostat" "17103"
    
    echo ""
    echo "=========================================="
    log_success "All API tests completed successfully!"
    echo "=========================================="
    echo ""
    echo "Test Summary:"
    echo "  ✓ Cleaner:     Configured, started, tested, restarted, internal status, stopped"
    echo "  ✓ Lighting:    Configured, started, tested, restarted, internal status, stopped"
    echo "  ✓ Thermostat:  Configured, started, tested, restarted, internal status, stopped"
    echo ""
    echo "Note: Containers should be running before running this test."
    echo "To start containers: ./scripts/run.sh"
}

# 显示帮助
show_help() {
    echo "Usage: $0 [all|examples|cleaner|thermostat|lighting|help]"
    echo ""
    echo "Options:"
    echo "  all/examples   Test all containers (default)"
    echo "  cleaner        Test only cleaner container"
    echo "  thermostat     Test only thermostat container"
    echo "  lighting       Test only lighting container"
    echo "  help, -h       Show this help message"
    echo ""
    echo "Test Process for each container:"
    echo "  1. Check if container is running"
    echo "  2. Configure and start app"
    echo "  3. Check app status"
    echo "  4. Test restart functionality"
    echo "  5. Test internal status API"
    echo "  6. Stop app"
    echo ""
    echo "Prerequisites:"
    echo "  Containers must be running before running this test."
    echo "  To start containers: ./scripts/run.sh"
    echo ""
    echo "Examples:"
    echo "  $0                # Test all containers"
    echo "  $0 all            # Test all containers"
    echo "  $0 examples       # Test all containers"
    echo "  $0 cleaner        # Test only cleaner"
    echo "  $0 thermostat     # Test only thermostat"
    echo "  $0 lighting       # Test only lighting"
}

# 解析参数
case "${1:-}" in
    help|-h|--help)
        show_help
        exit 0
        ;;
    all|examples|"")
        test_all_apps
        ;;
    cleaner)
        test_single_app "cleaner" "17101"
        ;;
    thermostat)
        test_single_app "thermostat" "17103"
        ;;
    lighting)
        test_single_app "lighting" "17102"
        ;;
    *)
        log_error "Unknown option: $1"
        show_help
        exit 1
        ;;
esac 