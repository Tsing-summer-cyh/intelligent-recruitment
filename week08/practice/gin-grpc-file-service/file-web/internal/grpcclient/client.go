// file-web/internal/grpcclient/client.go
package grpcclient

import (
	"log"

	pb "file-web/gen" // 引入 file-web 自己生成的 proto 代码

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var FileClient pb.FileServiceClient

func InitClient() {
	// 连接本地的 50051 端口 (file-service)
	// 注意：新版 gRPC 推荐使用 grpc.NewClient 替代 grpc.Dial
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("无法连接到 gRPC 文件服务: %v", err)
	}
	
	// 初始化客户端实例
	FileClient = pb.NewFileServiceClient(conn)
	log.Println("成功连接到 gRPC 文件服务！")
}