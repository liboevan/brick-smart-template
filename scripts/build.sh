#!/bin/bash

set -e

VERSION=${VERSION:-"0.1.0-dev"}
BUILD_TIME=${BUILD_TIME:-$(date -u +"%Y-%m-%dT%H:%M:%SZ")}
BUILD_DATE=${BUILD_DATE:-$(date -u +"%Y-%m-%d")}

echo "Building brick-smart-template..."
echo "Version: $VERSION"
echo "Build time: $BUILD_TIME"

# 检查是否跳过proxy构建
if [ "$1" = "skip-proxy" ] || [ "$2" = "skip-proxy" ]; then
    echo "Skipping proxy build as requested..."
else
echo "Building proxy image..."
docker build --build-arg VERSION=$VERSION --build-arg BUILD_TIME=$BUILD_TIME --build-arg BUILD_DATE=$BUILD_DATE -t el/brick-smart-template:$VERSION -t el/brick-smart-template:latest .
echo "Proxy image built successfully!"
fi

# 检查是否有参数来构建 examples
if [ "$1" = "examples" ] || [ "$1" = "all" ] || [ "$1" = "skip-proxy" ]; then
    echo ""
    echo "Building examples..."
    
    echo "Building brick-smart-cleaner..."
    cd examples/brick-smart-cleaner
    docker build --build-arg VERSION=$VERSION --build-arg BUILD_TIME=$BUILD_TIME --build-arg BUILD_DATE=$BUILD_DATE -t el/brick-smart-cleaner:$VERSION -t el/brick-smart-cleaner:latest .
    cd ../..

    echo "Building brick-smart-thermostat..."
    cd examples/brick-smart-thermostat
    docker build --build-arg VERSION=$VERSION --build-arg BUILD_TIME=$BUILD_TIME --build-arg BUILD_DATE=$BUILD_DATE -t el/brick-smart-thermostat:$VERSION -t el/brick-smart-thermostat:latest .
    cd ../..

    echo "Building brick-smart-lighting..."
    cd examples/brick-smart-lighting
    docker build --build-arg VERSION=$VERSION --build-arg BUILD_TIME=$BUILD_TIME --build-arg BUILD_DATE=$BUILD_DATE -t el/brick-smart-lighting:$VERSION -t el/brick-smart-lighting:latest .
    cd ../..

    echo ""
    echo "All images built successfully!"
    echo "Available images:"
    echo "- el/brick-smart-template:$VERSION"
    echo "- el/brick-smart-template:latest"
    echo "- el/brick-smart-cleaner:$VERSION"
    echo "- el/brick-smart-cleaner:latest"
    echo "- el/brick-smart-thermostat:$VERSION"
    echo "- el/brick-smart-thermostat:latest"
    echo "- el/brick-smart-lighting:$VERSION"
    echo "- el/brick-smart-lighting:latest"
else
    echo ""
    echo "Usage:"
    echo "  ./scripts/build.sh                    - Build proxy only"
    echo "  ./scripts/build.sh examples           - Build proxy + all examples"
    echo "  ./scripts/build.sh all                - Build proxy + all examples"
    echo "  ./scripts/build.sh skip-proxy         - Build examples only (skip proxy)"
    echo "  ./scripts/build.sh examples skip-proxy - Build examples only (skip proxy)"
    echo ""
    echo "Each image will have two tags:"
    echo "  - el/image:$VERSION (version tag)"
    echo "  - el/image:latest (latest tag)"
fi