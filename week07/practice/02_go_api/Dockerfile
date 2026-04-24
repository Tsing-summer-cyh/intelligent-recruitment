# ==========================================
# 第一阶段：编译构建阶段 (Builder)
# ==========================================
FROM golang:1.25-alpine AS builder

# 设置容器内工作目录
WORKDIR /app

# 先复制依赖配置文件并下载，利用 Docker 缓存机制加速后续构建
COPY go.mod go.sum ./
RUN go mod download

# 复制所有源代码到容器内
COPY . .

# 关键一步：静态编译 Go 程序
# CGO_ENABLED=0 确保编译出的二进制文件不依赖宿主机的 C 语言动态链接库，这样它才能在最简环境中运行
# GOOS=linux 确保编译出的是 Linux 系统下的可执行文件
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# ==========================================
# 第二阶段：极简运行环境 (Runner)
# ==========================================
# scratch 是 Docker 中的“空镜像”，里面什么都没有，体积为 0
FROM scratch

# 从第一阶段 (名字叫 builder) 中，把编译好的二进制文件 main 复制到当前空镜像的根目录下
COPY --from=builder /app/main /main

# 声明容器内部使用 8081 端口
EXPOSE 8081

# 容器启动时执行该二进制程序
CMD ["/main"]