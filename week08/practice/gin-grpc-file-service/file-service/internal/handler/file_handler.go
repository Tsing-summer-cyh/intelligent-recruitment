// file-service/internal/handler/file_handler.go
package handler

import (
	"context"
	"errors"

	pb "file-service/gen" // 引入我们刚刚生成的 proto 代码
	"file-service/internal/database"
	"file-service/internal/model"

	"gorm.io/gorm"
)

// FileServiceServer 必须嵌入 UnimplementedFileServiceServer 以保证向前兼容
type FileServiceServer struct {
	pb.UnimplementedFileServiceServer
}

// 1. 保存文件记录 [cite: 51]
func (s *FileServiceServer) SaveFileRecord(ctx context.Context, req *pb.SaveFileRequest) (*pb.SaveFileResponse, error) {
	// 将 gRPC 请求数据映射到数据库模型
	record := model.FileRecord{
		OriginalName: req.OriginalName,
		StoredName:   req.StoredName,
		Size:         req.Size,
		MimeType:     req.MimeType,
		Path:         req.Path,
	}

	// 插入 SQLite 数据库
	if err := database.DB.Create(&record).Error; err != nil {
		return nil, err
	}

	// 将结果转回 gRPC 响应格式并返回
	return &pb.SaveFileResponse{
		File: convertToPbRecord(&record),
	}, nil
}

// 2. 获取文件列表 [cite: 52]
func (s *FileServiceServer) GetFileList(ctx context.Context, req *pb.GetFileListRequest) (*pb.GetFileListResponse, error) {
	var records []model.FileRecord
	// 从数据库查询所有记录
	if err := database.DB.Find(&records).Error; err != nil {
		return nil, err
	}

	var pbFiles []*pb.FileRecord
	for _, r := range records {
		pbFiles = append(pbFiles, convertToPbRecord(&r))
	}

	return &pb.GetFileListResponse{
		Files: pbFiles,
	}, nil
}

// 3. 根据 ID 查询单条记录（供下载使用） [cite: 111-112]
func (s *FileServiceServer) GetFileRecordById(ctx context.Context, req *pb.GetFileByIdRequest) (*pb.GetFileByIdResponse, error) {
	var record model.FileRecord
	// 根据 ID 查询单条数据
	if err := database.DB.First(&record, req.Id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文件记录不存在")
		}
		return nil, err
	}

	return &pb.GetFileByIdResponse{
		File: convertToPbRecord(&record),
	}, nil
}

// 辅助函数：将数据库模型转换为 proto 消息结构
func convertToPbRecord(r *model.FileRecord) *pb.FileRecord {
	return &pb.FileRecord{
		Id:           r.ID,
		OriginalName: r.OriginalName,
		StoredName:   r.StoredName,
		Size:         r.Size,
		MimeType:     r.MimeType,
		Path:         r.Path,
	}
}