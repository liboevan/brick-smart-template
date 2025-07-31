# Brick Smart Template

> English version: [README.en.md](README.en.md)

这是一个用于构建和管理智能设备代理及示例应用的模板项目，专注于物联网、智能家居场景的快速原型开发和教学演示。本项目不包含生产环境功能，仅用于概念验证和技术实践。

## 项目简介

Brick Smart Template 提供智能设备代理（proxy）和多种智能设备（cleaner、lighting、thermostat等）示例的统一构建、运行、测试和清理脚本。

- **一键构建与运行**：通过Makefile和脚本快速完成所有服务的构建、启动与清理
- **多设备示例**：内置扫地机、灯光、温控等多种智能设备模拟
- **Docker支持**：所有服务均可容器化部署，确保环境一致性

## 系统架构

### 核心组件
- **brick-proxy**：智能设备代理，管理设备连接和通信
- **设备示例**：多种智能设备模拟实现（cleaner、lighting等）
- **Traefik反向代理**：统一管理服务访问入口

## 目录结构

```
brick-smart-template/
├── Makefile                # 构建入口
├── Dockerfile              # 代理镜像构建文件
├── docker-compose.yml      # 多服务编排
├── scripts/                # 构建/运行/测试/清理等脚本
├── examples/               # 设备示例代码
│   ├── brick-smart-cleaner/
│   ├── brick-smart-lighting/
│   └── brick-smart-thermostat/
├── docs/                   # 详细文档
├── README.md               # 项目说明（中文）
└── README.en.md            # 项目说明（英文）
```

## 快速开始

```bash
make build      # 构建所有镜像
make run        # 启动所有容器
make test       # 执行所有API测试
make clean      # 清理所有容器和镜像
```

## 主要文档

设备示例的详细文档请查看各示例项目中的README：
- [扫地机示例](examples/brick-smart-cleaner/README.md)
- [灯光示例](examples/brick-smart-lighting/README.md)
- [温控器示例](examples/brick-smart-thermostat/README.md)



## 适用场景

- 物联网设备开发与测试
- 智能家居系统原型
- 教学与演示
- 多设备统一管理与自动化

## Traefik 反向代理配置

本项目使用Traefik作为反向代理，统一管理各服务的访问入口。

### 主要特性
- 自动发现Docker容器服务
- 通过路径前缀区分不同服务
- 自动剥离路径前缀，无需修改服务代码

### 服务访问方式
所有服务通过Traefik的17111端口访问，URL格式为：
```
http://localhost:17111/<service-name>/<api-endpoint>
```

### 示例
访问cleaner服务的health API：
```
http://localhost:17111/brick-cleaner/health
```

### 配置说明
- Traefik配置文件：`scripts/start_traefik.sh`
- 服务注册通过Docker标签实现，详见`scripts/run.sh`中的`run_container`函数
