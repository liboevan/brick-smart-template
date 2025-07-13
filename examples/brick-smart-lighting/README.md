# Smart Lighting Example

这是一个智能照明（lighting）模拟应用示例。

## 功能
- 支持通过命令行参数指定设备ID、gRPC端口
- 定期（每5秒）模拟并上报开关、亮度、模式等状态
- 通过HTTP API（兼容proxy）上报状态

## 命令行参数
- `-id`：设备ID（必需）
- `-grpc-port`：gRPC端口（默认50051，实际未用，仅为兼容）
- `-config`：配置文件路径（JSON，暂未用）
- `-help`：显示帮助

## 用法示例
```bash
./lighting -id light-001 -grpc-port 50051
```

## 状态上报格式
```json
{
  "device_id": "light-001",
  "device_type": "lighting",
  "timestamp": 1703123456,
  "data": {
    "is_on": true,
    "brightness": 80,
    "mode": "normal"
  }
}
```

## Docker 构建与运行
```bash
docker build -t smart-lighting .
docker run --rm smart-lighting ./lighting -id light-001 -grpc-port 50051
``` 