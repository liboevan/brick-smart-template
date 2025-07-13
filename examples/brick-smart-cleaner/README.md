# Smart Cleaner - 智能扫地机示例

这是一个智能扫地机的示例应用，模拟真实的扫地机工作过程。

## 功能特性

- **循环清理**: 在6个房间中循环清理
- **进度跟踪**: 每个房间清理1分钟，显示百分比进度
- **状态上报**: 通过gRPC实时上报清理状态
- **随机事件**: 模拟遇到障碍物等真实情况

## 房间列表

1. 客厅
2. 卧室  
3. 厨房
4. 卫生间
5. 书房
6. 阳台

## 使用方法

### 命令行参数

```bash
./cleaner -id <设备ID> -grpc-port <gRPC端口> [-config <配置文件>]
```

### 参数说明

- `-id`: 扫地机设备ID（必需）
- `-grpc-port`: gRPC服务器端口（默认: 50051）
- `-config`: 配置文件路径（JSON格式，暂未使用）
- `-help`: 显示帮助信息

### 示例

```bash
# 基本使用
./cleaner -id cleaner-001 -grpc-port 50051

# 指定不同端口
./cleaner -id cleaner-002 -grpc-port 50052

# 显示帮助
./cleaner -help
```

## 状态信息

扫地机会上报以下状态信息：

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

### 字段说明

- `cycle_count`: 当前循环次数
- `room_id`: 正在清理的房间ID
- `room_name`: 房间名称
- `progress`: 清理进度（0-100%）
- `total_time`: 总运行时间（秒）
- `current_time`: 当前房间清理时间（秒）
- `status`: 状态（cleaning/completed）

## Docker 构建

```bash
# 构建镜像
docker build -t smart-cleaner .

# 运行容器
docker run --rm smart-cleaner ./cleaner -id cleaner-001 -grpc-port 50051
```

## 项目结构

```
smart-cleaner/
├── main.go                 # 主程序入口
├── go.mod                  # Go模块文件
├── Dockerfile              # Docker构建文件
├── README.md              # 项目文档
└── pkg/
    ├── cleaner/
    │   └── cleaner.go     # 扫地机核心逻辑
    └── grpcclient/
        └── client.go      # gRPC客户端
```

## 开发说明

### 扩展其他设备

这个示例展示了如何创建一个设备应用：

1. **定义设备类型**: 在状态消息中指定 `device_type`
2. **实现业务逻辑**: 在对应的包中实现设备特定逻辑
3. **状态上报**: 使用统一的gRPC客户端上报状态
4. **参数配置**: 通过命令行参数配置设备行为

### 消息格式

所有设备都使用统一的消息格式：

```json
{
  "device_id": "设备ID",
  "device_type": "设备类型",
  "timestamp": "时间戳",
  "data": {
    // 设备特定的数据
  }
}
```

这种设计使得代理服务可以统一处理不同类型的设备状态。 