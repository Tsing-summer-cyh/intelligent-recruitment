// file-web/internal/handler/file_handler.go
package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	pb "file-web/gen"
	"file-web/internal/grpcclient"

	"github.com/gin-gonic/gin"
)

// 1. 批量上传文件接口 (POST /api/files/uploads) [cite: 74-75]
func UploadFiles(c *gin.Context) {
	// 从请求中获取 multipart form
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "解析表单失败"})
		return
	}

	// 假设前端上传文件的字段名是 "files"
	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "没有找到上传的文件"})
		return
	}

	var results []interface{}

	for _, fileHeader := range files {
		// 打开上传的文件流
		file, err := fileHeader.Open()
		if err != nil {
			continue
		}
		defer file.Close()

		// 计算文件的 SHA256 Hash [cite: 79]
		hash := sha256.New()
		if _, err := io.Copy(hash, file); err != nil {
			continue
		}
		hashString := hex.EncodeToString(hash.Sum(nil))

		// 基于 hash 生成新文件名 (取前16位作为文件名) [cite: 80]
		ext := filepath.Ext(fileHeader.Filename)
		storedName := hashString[:16] + ext
		savePath := filepath.Join("uploads", storedName)

		// 保存文件到本地 uploads 目录 [cite: 81]
		if err := c.SaveUploadedFile(fileHeader, savePath); err != nil {
			continue
		}

		// 调用 gRPC 服务保存文件记录 [cite: 82]
		req := &pb.SaveFileRequest{
			OriginalName: fileHeader.Filename,
			StoredName:   storedName,
			Size:         fileHeader.Size,
			MimeType:     fileHeader.Header.Get("Content-Type"),
			Path:         savePath,
		}

		resp, err := grpcclient.FileClient.SaveFileRecord(c.Request.Context(), req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "保存记录失败: " + err.Error()})
			return
		}

		results = append(results, resp.File)
	}

	// 返回结果 [cite: 83]
	c.JSON(http.StatusOK, gin.H{"uploaded": results})
}

// 2. 获取文件列表接口 (GET /api/files) [cite: 94-95]
func GetFiles(c *gin.Context) {
	resp, err := grpcclient.FileClient.GetFileList(c.Request.Context(), &pb.GetFileListRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取列表失败"})
		return
	}
	c.JSON(http.StatusOK, resp.Files) // 返回示例中的 JSON 数组格式 [cite: 96-106]
}

// 3. 下载文件接口 (GET /api/files/download/:id) [cite: 107-108]
func DownloadFile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	// 调用 gRPC 查询文件记录 [cite: 111]
	resp, err := grpcclient.FileClient.GetFileRecordById(c.Request.Context(), &pb.GetFileByIdRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	// 检查本地文件是否存在
	if _, err := os.Stat(resp.File.Path); os.IsNotExist(err) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "物理文件已丢失"})
		return
	}

	// 返回文件内容给客户端，并设置原始文件名 [cite: 113-116]
	c.FileAttachment(resp.File.Path, resp.File.OriginalName)
}