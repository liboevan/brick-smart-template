# Smart Cleaner 使用指南

## 目标场景

1. 启动一个smart-cleaner container，启动时运行proxy，不运行cleaner app
2. 调用proxy API可以启动和停止cleaner app
3. cleaner app在运行过程中会通过HTTP API向proxy上报状态
4. proxy提供API可查询状态，只存储最新状态

## 快速开始

### 1. 启动Proxy容器

```bash
cd brick-smart-template
docker-compose up -d
```

这会在端口17100启动proxy服务。

### 2. 配置Smart Cleaner应用

```bash
curl -X POST http://localhost:17100/app/configure \
  -H "Content-Type: application/json" \
  -d '{
    "app_info": {
      "name": "smart-cleaner",
      "command": "./cleaner",
      "args": ["-id", "cleaner-001", "-grpc-port", "50051"],
      "working_dir": "/app/myapp",
      "env": {
        "PROXY_HTTP_PORT": "17100"
      }
    }
  }'
```

### 3. 启动Smart Cleaner

```bash
curl -X POST http://localhost:17100/app/start \
  -H "Content-Type: application/json" \
  -d '{
    "profile": "{\"cleaner_id\": \"cleaner-001\", \"rooms\": 6}"
  }'
```

### 4. 查询应用状态

```bash
curl http://localhost:17100/app/status
```

### 5. 停止Smart Cleaner

```bash
curl -X POST http://localhost:17100/app/stop
```

## API接口

### 健康检查
```
GET /health
```

### 配置应用
```
POST /app/configure
Content-Type: application/json

{
  "app_info": {
    "name": "smart-cleaner",
    "command": "./cleaner",
    "args": ["-id", "cleaner-001", "-grpc-port", "50051"],
    "working_dir": "/app/myapp",
    "env": {
      "PROXY_HTTP_PORT": "17100"
    }
  }
}
```

### 启动应用
```
POST /app/start
Content-Type: application/json

{
  "profile": "{\"cleaner_id\": \"cleaner-001\", \"rooms\": 6}"
}
```

### 停止应用
```
POST /app/stop
```

### 查询状态
```
GET /app/status
```

### 状态上报（内部使用）
```
POST /app/status/report
Content-Type: application/json

{
  "status": "running",
  "data": {
    "cycle_count": 1,
    "room_id": 1,
    "room_name": "客厅",
    "progress": 45,
    "total_time": 120.5,
    "current_time": 27.3,
    "status": "cleaning"
  }
}
```

## 状态信息

Smart Cleaner会定期上报以下状态信息：

```json
{
  "device_id": "cleaner-001",
  "device_type": "cleaner",
  "timestamp": 1703123456,
  "data": {
    "cycle_count": 1,
    "room_id": 1,
    "room_name": "客厅",
    "progress": 45,
    "total_time": 120.5,
    "current_time": 27.3,
    "status": "cleaning"
  }
}
```

## 测试脚本

运行完整的集成测试：

```bash
chmod +x test-smart-cleaner.sh
./test-smart-cleaner.sh
```

## 项目结构

```
brick-smart-template/
├── docker-compose.yml          # 容器编排
├── Dockerfile                  # Proxy镜像
├── cmd/proxy/main.go          # Proxy主程序
├── pkg/
│   ├── appmanager/            # 应用管理器
│   ├── httpapi/               # HTTP API服务
│   └── grpcservice/           # gRPC服务
└── examples/
    └── smart-cleaner/         # Smart Cleaner示例
        ├── main.go            # 主程序
        ├── Dockerfile         # Cleaner镜像
        └── pkg/
            ├── cleaner/       # 清理逻辑
            └── grpcclient/    # 状态上报客户端
```

## 部署选项

### 选项1：提供Proxy二进制文件
- 开发者下载proxy二进制文件
- 自行配置和运行

### 选项2：提供Proxy Docker镜像
- 开发者拉取proxy镜像
- 使用docker-compose或docker run启动

## 扩展性

这个架构支持：
- 多种设备类型（不仅限于cleaner）
- 统一的状态管理
- 灵活的配置方式
- 容器化部署 