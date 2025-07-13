# 基于proxy镜像
FROM el/brick-smart-template:latest AS proxy

# 构建阶段
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制cleaner源码
COPY . .

# 下载cleaner依赖
RUN go mod download

# 构建cleaner
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o cleaner ./main.go

# 运行阶段
FROM alpine:latest

# 安装ca-certificates
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk --no-cache add ca-certificates

# 统一工作目录
WORKDIR /app

# 复制 proxy 二进制和版本文件
COPY --from=proxy /app/app-proxy ./
COPY --from=proxy /app/proxy.VERSION ./
COPY --from=proxy /app/proxy.build-info.json ./

# 复制 cleaner 二进制
COPY --from=builder /app/cleaner ./

# 定义构建参数
ARG VERSION=0.1.0-dev
ARG BUILD_TIME
ARG BUILD_DATE
ARG GIT_COMMIT
ARG GIT_BRANCH

# 生成app版本文件
RUN echo "$VERSION" > VERSION && \
    echo '{"version":"'$VERSION'","build_time":"'$BUILD_TIME'","build_date":"'$BUILD_DATE'","git_commit":"'$GIT_COMMIT'","git_branch":"'$GIT_BRANCH'"}' > build-info.json

# 生成 manifest.json
RUN echo '{"app_name": "cleaner", "health_check_interval": 3, "default_args": ["-id", "cleaner-001"]}' > /app/manifest.json

# 创建日志目录
RUN mkdir -p /app/logs

# 运行 proxy
ENTRYPOINT ["./app-proxy"]
