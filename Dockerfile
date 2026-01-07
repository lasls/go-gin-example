# 构建阶段
FROM golang:alpine AS builder

# 设置环境变量
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOPROXY=https://goproxy.cn,direct

# 工作目录
WORKDIR /app

# 复制 go mod 文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN go build -a -installsuffix cgo -o go-gin-example .

# 运行阶段
FROM alpine:latest

# 安装必要的包，包括 CA 证书和 DNS 解析
RUN apk --no-cache add ca-certificates tzdata

# 创建非 root 用户
RUN addgroup -g 65532 nonroot &&\
    adduser -D -u 65532 -G nonroot nonroot

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/go-gin-example .

# 复制配置文件
COPY conf/app.ini ./conf/app.ini

# 更改文件权限
RUN chmod +x ./go-gin-example
RUN chown -R nonroot:nonroot /root/

# 切换到非 root 用户
USER nonroot

# 暴露端口
EXPOSE 8000

# 启动应用
ENTRYPOINT ["./go-gin-example"]