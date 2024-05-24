# 构建最小运行时镜像
FROM golang:1.22-alpine AS builder
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 go build main.go

WORKDIR /app
# configs为你的配置文件目录
RUN cp /build/main /app

#FROM alpine:3.18
#FROM lampnick/runtime:chromium-alpine

FROM chromedp_debian:latest
COPY  --from=builder /app /app
WORKDIR /app

# 项目启动命令
ENTRYPOINT ["./main"]