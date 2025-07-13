# Smart Thermostat Example

这是一个智能温控器（thermostat）模拟应用示例。

## 功能
- 支持通过命令行参数指定设备ID、gRPC端口
- 定期（每5秒）模拟并上报温度、目标温度、模式等状态
- 通过HTTP API（兼容proxy）上报状态

## 命令行参数
- `-id`：设备ID（必需）
- `-grpc-port`：gRPC端口（默认50051，实际未用，仅为兼容）
- `-config`：配置文件路径（JSON，暂未用）
- `-help`：显示帮助

## 用法示例
```bash
./thermostat -id thermo-001 -grpc-port 50051
```

## 状态上报格式
```json
{
  "device_id": "thermo-001",
  "device_type": "thermostat",
  "timestamp": 1703123456,
  "data": {
    "room_temp": 22.3,
    "target_temp": 22.0,
    "mode": "auto"
  }
}
```

## Docker 构建与运行
```bash
docker build -t smart-thermostat .
docker run --rm smart-thermostat ./thermostat -id thermo-001 -grpc-port 50051
``` 