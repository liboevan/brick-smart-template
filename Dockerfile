# 构建阶段
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制go mod文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app-proxy ./cmd/proxy

# 定义构建参数
ARG VERSION=0.1.0-dev
ARG BUILD_TIME
ARG BUILD_DATE
ARG GIT_COMMIT
ARG GIT_BRANCH

# 生成proxy版本文件（builder阶段）
RUN echo "$VERSION" > /app/proxy.VERSION && \
    echo '{"version":"'$VERSION'","build_time":"'$BUILD_TIME'","build_date":"'$BUILD_DATE'","git_commit":"'$GIT_COMMIT'","git_branch":"'$GIT_BRANCH'"}' > /app/proxy.build-info.json

# 运行阶段
FROM alpine:latest

# 安装ca-certificates
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk --no-cache add ca-certificates

# 统一工作目录
WORKDIR /app

# 复制二进制和版本文件
COPY --from=builder /app/app-proxy ./

# 复制版本信息
COPY --from=builder /app/proxy.VERSION ./
COPY --from=builder /app/proxy.build-info.json ./

# 生成 manifest.json
RUN echo '{"app_name": "proxy", "health_check_interval": 3}' > /app/manifest.json

# 运行应用
CMD ["./app-proxy"]
