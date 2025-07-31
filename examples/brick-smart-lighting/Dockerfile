# 基于proxy镜像
FROM brick-smart-template:latest AS proxy

# 构建阶段
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制lighting源码
COPY . .

# 下载lighting依赖
RUN go mod download

# 构建lighting
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o lighting ./main.go

# 生成 manifest.json
RUN echo '{"app_name": "lighting", "health_check_interval": 3, "default_args": ["-id", "lighting-001"]}' > /app/manifest.json

# 运行阶段
FROM alpine:latest

# 安装ca-certificates
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk --no-cache add ca-certificates

# 统一工作目录
WORKDIR /app

# 复制 proxy 二进制
COPY --from=proxy /app/app-proxy ./

# 复制 lighting 二进制
COPY --from=builder /app/lighting ./

# 创建日志目录
RUN mkdir -p /app/logs

# 运行 proxy
ENTRYPOINT ["./app-proxy"]