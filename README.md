# Brick Smart Template

> English version: [README.en.md](README.en.md)

这是一个以构建和管理智能设备代理及示例应用为例、用于概念和技术验证与实践、无实际生产用途的项目。

## 项目简介

Brick Smart Template 提供了智能设备代理（proxy）和多种智能设备（如 cleaner、lighting、thermostat）示例的统一构建、运行、测试和清理脚本，适合物联网、智能家居等场景的快速原型开发和教学演示。

- **一键构建与运行**：通过 Makefile 和脚本，快速完成所有服务的构建、启动与清理。
- **多设备示例**：内置扫地机、灯光、温控等多种智能设备模拟。
- **Docker 支持**：所有服务均可容器化部署，便于本地和云端环境一致性。
- **文档完善**：中英文双语文档，详细介绍各设备和测试方法。

## 目录结构

```
brick-smart-template/
├── Makefile                # 构建入口
├── Dockerfile              # 代理镜像构建文件
├── docker-compose.yml      # 多服务编排
├── scripts/                # 构建/运行/测试/清理等脚本
├── examples/               # 设备示例代码
├── docs/                   # 详细文档（设备说明、测试说明等）
└── README.md/README.en.md  # 项目说明（中/英文）
```

## 快速开始

```bash
make build      # 构建所有镜像
make run        # 启动所有容器
make test       # 执行所有API测试
make clean      # 清理所有容器和镜像
```

## 主要文档

- [设备文档 cleaner](docs/cleaner.md)
- [测试说明 Test Guide (EN)](docs/test.md)
- [English README](README.en.md)

## 适用场景

- 物联网设备开发与测试
- 智能家居系统原型
- 教学与演示
- 多设备统一管理与自动化

如需详细测试方法、API调用流程等，请参见 docs/test.md。
