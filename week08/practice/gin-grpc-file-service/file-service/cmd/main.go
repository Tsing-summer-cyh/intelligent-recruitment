// file-service/cmd/main.go
package main

import (
	"log"
	"net"

	pb "file-service/gen"
	"file-service/internal/database"
	"file-service/internal/handler"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// 1. 初始化 SQLite 数据库
	database.InitDB()

	// 2. 监听本地 50051 端口
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("无法监听端口 50051: %v", err)
	}

	// 3. 创建 gRPC 服务器实例
	grpcServer := grpc.NewServer()

	// 4. 将我们刚才写的 handler 注册到 gRPC 服务器上
	pb.RegisterFileServiceServer(grpcServer, &handler.FileServiceServer{})

	// 注册反射服务（可选，方便使用 postman 或 grpcui 等工具调试）
	reflection.Register(grpcServer)

	log.Println("gRPC 文件服务已启动，正在监听端口 :50051 ...")
	
	// 5. 启动服务，阻塞等待
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("启动 gRPC 服务失败: %v", err)
	}
}