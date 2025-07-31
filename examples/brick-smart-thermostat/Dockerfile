# 基于proxy镜像
FROM brick-smart-template:latest AS proxy

# 构建阶段
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制thermostat源码
COPY . .

# 下载thermostat依赖
RUN go mod download

# 构建thermostat
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o thermostat ./main.go

# 生成 manifest.json
RUN echo '{"app_name": "thermostat", "health_check_interval": 3, "default_args": ["-id", "thermostat-001"]}' > /app/manifest.json

# 运行阶段
FROM alpine:latest

# 安装ca-certificates
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk --no-cache add ca-certificates

# 统一工作目录
WORKDIR /app

# 复制 proxy 二进制
COPY --from=proxy /app/app-proxy ./

# 复制 thermostat 二进制
COPY --from=builder /app/thermostat ./

# 创建日志目录
RUN mkdir -p /app/logs

# 运行 proxy
ENTRYPOINT ["./app-proxy"]