# 阶段一：编译
FROM golang:alpine AS builder
WORKDIR /app
# 设置 GOPROXY 加速下载
ENV GOPROXY=https://goproxy.cn,direct
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 阶段二：运行环境
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY .env .env
EXPOSE 8080
CMD ["./main"]