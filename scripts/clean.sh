#!/bin/bash

# 清理脚本：删除所有或指定容器，支持 --image 删除镜像

set -e

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

# 容器与镜像映射
get_container_image() {
    case $1 in
        cleaner)
            echo "brick-cleaner brick-smart-cleaner"
            ;;
        thermostat)
            echo "brick-thermostat brick-smart-thermostat"
            ;;
        lighting)
            echo "brick-lighting brick-smart-lighting"
            ;;
        *)
            ;;
    esac
}

# 删除容器
remove_container() {
    local container_name=$1
    if docker ps -a -q -f name=$container_name | grep -q .; then
        log_info "Removing container $container_name..."
        docker rm -f $container_name 2>/dev/null || true
        log_success "Container $container_name removed."
    else
        log_warning "Container $container_name does not exist."
    fi
}

# 删除镜像
remove_image() {
    local image=$1
    # 删除所有 tag 的镜像
    local image_ids=$(docker images --format '{{.Repository}}:{{.Tag}} {{.ID}}' | grep "^$image:" | awk '{print $2}')
    if [ -n "$image_ids" ]; then
        log_info "Removing all tags of image $image..."
        for imgid in $image_ids; do
            docker rmi -f $imgid 2>/dev/null || true
        done
        log_success "All tags of image $image removed."
    else
        log_warning "Image $image does not exist."
    fi
}

# 主清理逻辑
clean_target() {
    local target=$1
    local remove_image_flag=$2
    local container image
    read container image < <(get_container_image $target)
    if [ -n "$container" ]; then
        remove_container $container
        if [ "$remove_image_flag" = "true" ]; then
            remove_image $image
        fi
    fi
}

# 显示帮助
show_help() {
    echo "Usage: $0 [all|examples|cleaner|thermostat|lighting|help] [--image]"
    echo ""
    echo "Options:"
    echo "  all/examples   Clean all containers (default)"
    echo "  cleaner        Clean only cleaner container"
    echo "  thermostat     Clean only thermostat container"
    echo "  lighting       Clean only lighting container"
    echo "  --image        Also remove docker images"
    echo "  help, -h       Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                # Clean all containers"
    echo "  $0 all            # Clean all containers"
    echo "  $0 cleaner        # Clean only cleaner"
    echo "  $0 --image        # Clean all containers and images"
    echo "  $0 lighting --image # Clean lighting container and image"
}

# 解析参数
REMOVE_IMAGE=false
POSITIONAL=()
for arg in "$@"; do
    case $arg in
        --image)
            REMOVE_IMAGE=true
            shift
            ;;
        *)
            POSITIONAL+=("$arg")
            ;;
    esac
    shift $(( $# > 0 ? 1 : 0 ))
done
set -- "${POSITIONAL[@]}"

case "${1:-}" in
    help|-h|--help)
        show_help
        exit 0
        ;;
    all|examples|"")
        clean_target cleaner $REMOVE_IMAGE
        clean_target thermostat $REMOVE_IMAGE
        clean_target lighting $REMOVE_IMAGE
        ;;
    cleaner)
        clean_target cleaner $REMOVE_IMAGE
        ;;
    thermostat)
        clean_target thermostat $REMOVE_IMAGE
        ;;
    lighting)
        clean_target lighting $REMOVE_IMAGE
        ;;
    *)
        log_error "Unknown option: $1"
        show_help
        exit 1
        ;;
esac

log_success "Clean operation completed."