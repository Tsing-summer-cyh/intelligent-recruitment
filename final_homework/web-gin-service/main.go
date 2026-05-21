package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"web-gin-service/pb"
	"web-gin-service/utils"
)

var grpcClient pb.RecruitmentServiceClient

func InitGRPCClient() {
	conn, err := grpc.NewClient("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("连接 gRPC 服务失败: %v", err)
	}
	grpcClient = pb.NewRecruitmentServiceClient(conn)
	fmt.Println("🔗 成功连接到 Logic gRPC 核心服务！")
}

func main() {
	InitGRPCClient()
	r := gin.Default()
	utils.InitEino()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(corsConfig))

	api := r.Group("/api/v1")
	{
		// ==========================================
		// [公开接口] 游客免登录可访问
		// ==========================================
		publicGroup := api.Group("/public")
		{
			publicGroup.GET("/jobs", func(c *gin.Context) {
				req := &pb.ListHRJobsRequest{HrId: 0, Page: 1, Size: 100}
				resp, err := grpcClient.ListHRJobs(c.Request.Context(), req)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "查询失败"})
					return
				}
				c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": resp.Jobs})
			})
		}

		// ==========================================
		// [公开接口] 用户登录
		// ==========================================
		api.POST("/login", func(c *gin.Context) {
			var req pb.LoginRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "参数格式错误"})
				return
			}
			resp, err := grpcClient.Login(c.Request.Context(), &req)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"code": 1, "msg": err.Error()})
				return
			}
			token, err := utils.GenerateToken(resp.UserId, resp.Role)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "生成令牌失败"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "登录成功", "data": gin.H{
				"token": token, "role": resp.Role, "user_id": resp.UserId,
			}})
		})

		// ==========================================
		// [受保护接口] HR 专属路由组
		// ==========================================
		hrGroup := api.Group("/hr")
		hrGroup.Use(utils.JWTAuthMiddleware("hr"))
		{
			// 1. 发布岗位
			hrGroup.POST("/jobs", func(c *gin.Context) {
				var req pb.CreateJobRequest
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "参数格式错误"})
					return
				}
				hrID, _ := c.Get("userID")
				req.HrId = hrID.(int64)

				resp, err := grpcClient.CreateJob(c.Request.Context(), &req)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "发布失败"})
					return
				}
				c.JSON(http.StatusOK, gin.H{"code": 0, "msg": resp.Msg, "data": nil})
			})

			// 2. 获取个人发布的岗位列表
			hrGroup.GET("/jobs", func(c *gin.Context) {
				hrID, _ := c.Get("userID")
				req := &pb.ListHRJobsRequest{
					HrId: hrID.(int64),
					Page: 1,
					Size: 100,
				}

				resp, err := grpcClient.ListHRJobs(c.Request.Context(), req)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "查询失败"})
					return
				}
				c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": resp.Jobs})
			})

			// 3. 查看某个岗位的投递详情
			hrGroup.GET("/jobs/:id/applications", func(c *gin.Context) {
				jobID := c.Param("id")
				var id int64
				fmt.Sscanf(jobID, "%d", &id)

				req := &pb.GetApplicationsRequest{JobId: id}
				resp, err := grpcClient.GetJobApplications(c.Request.Context(), req)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "查询失败"})
					return
				}
				c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": resp})
			})

			// 4. 编辑/下架岗位
			hrGroup.PUT("/jobs/:id", func(c *gin.Context) {
				jobID := c.Param("id")
				var id int64
				fmt.Sscanf(jobID, "%d", &id)

				var req struct {
					Title       string `json:"title"`
					Description string `json:"description"`
					Status      int32  `json:"status"` // 1-上架, 0-下架
				}
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "参数格式错误"})
					return
				}

				hrID, _ := c.Get("userID")

				grpcReq := &pb.UpdateJobRequest{
					JobId:       id,
					HrId:        hrID.(int64),
					Title:       req.Title,
					Description: req.Description,
					Status:      req.Status,
				}

				resp, err := grpcClient.UpdateJob(c.Request.Context(), grpcReq)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "操作失败"})
					return
				}
				c.JSON(http.StatusOK, gin.H{"code": resp.Code, "msg": resp.Msg})
			})

			// 5. AI 智能问答接口（调用 gRPC，含历史持久化）
			hrGroup.POST("/ai/chat", func(c *gin.Context) {
				var req struct {
					Question string `json:"question"`
				}
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "参数格式错误"})
					return
				}

				hrID, _ := c.Get("userID")

				grpcReq := &pb.AIChatRequest{
					HrId:    hrID.(int64),
					Message: req.Question,
				}

				// 设置 60 秒超时，给 AI 足够时间响应
				ctx, cancel := context.WithTimeout(c.Request.Context(), 60*time.Second)
				defer cancel()

				resp, err := grpcClient.ChatWithAI(ctx, grpcReq)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "AI 思考失败: " + err.Error()})
					return
				}

				c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": resp.Reply})
			})

			// 6. 获取 AI 历史对话记录
			hrGroup.GET("/ai/history", func(c *gin.Context) {
				hrID, _ := c.Get("userID")

				req := &pb.GetChatHistoryRequest{HrId: hrID.(int64)}
				resp, err := grpcClient.GetChatHistory(c.Request.Context(), req)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "获取历史失败"})
					return
				}

				c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": resp.Messages})
			})
		}

		// ==========================================
		// [受保护接口] 候选人专属路由组
		// ==========================================
		candidateGroup := api.Group("/candidate")
		candidateGroup.Use(utils.JWTAuthMiddleware("candidate"))
		{
			// 简历格式强校验与真实 OSS 上传
			candidateGroup.POST("/upload/resume", func(c *gin.Context) {
				file, err := c.FormFile("file")
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "未读取到文件"})
					return
				}

				// 校验文件格式
				ext := strings.ToLower(filepath.Ext(file.Filename))
				if ext != ".pdf" && ext != ".doc" && ext != ".docx" {
					c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "非法文件格式！仅支持 PDF / DOC / DOCX"})
					return
				}

				// 校验文件大小 (最大 10MB)
				if file.Size > 10*1024*1024 {
					c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "文件过大！最大支持 10MB"})
					return
				}

				// 读取文件内容
				openedFile, err := file.Open()
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "读取文件失败"})
					return
				}
				defer openedFile.Close()

				content := make([]byte, file.Size)
				_, err = openedFile.Read(content)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "读取文件内容失败"})
					return
				}

				// 获取当前用户 ID
				userID, _ := c.Get("userID")

				// 调用 gRPC 上传到 OSS
				grpcReq := &pb.UploadResumeRequest{
					UserId:   userID.(int64),
					Filename: file.Filename,
					Content:  content,
				}

				resp, err := grpcClient.UploadResume(c.Request.Context(), grpcReq)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "上传失败: " + err.Error()})
					return
				}

				if resp.Code != 0 {
					c.JSON(http.StatusInternalServerError, gin.H{"code": resp.Code, "msg": resp.Msg})
					return
				}

				c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "简历上传成功", "data": resp.OssUrl})
			})

			// 获取职位大厅的所有岗位列表
			candidateGroup.GET("/jobs", func(c *gin.Context) {
				req := &pb.ListHRJobsRequest{
					HrId: 0,
					Page: 1,
					Size: 100,
				}

				resp, err := grpcClient.ListHRJobs(c.Request.Context(), req)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "获取职位大厅失败"})
					return
				}
				c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": resp.Jobs})
			})

			// 查询个人档案
			candidateGroup.GET("/profile", func(c *gin.Context) {
				userID, _ := c.Get("userID")

				req := &pb.GetProfileRequest{UserId: userID.(int64)}
				resp, err := grpcClient.GetProfile(c.Request.Context(), req)

				if err != nil {
					c.JSON(http.StatusOK, gin.H{"code": 0, "data": nil})
					return
				}

				c.JSON(http.StatusOK, gin.H{
					"code": 0,
					"msg":  "success",
					"data": gin.H{
						"real_name":          resp.RealName,
						"phone":              resp.Phone,
						"highest_education":  resp.HighestEducation,
						"university":         resp.University,
						"experience":         resp.Experience,
						"skills":             resp.Skills,
						"resume_url":         resp.ResumeOssUrl,
					},
				})
			})

			// 完善个人档案
			candidateGroup.PUT("/profile", func(c *gin.Context) {
				var tempForm struct {
					RealName          string `json:"real_name"`
					Phone             string `json:"phone"`
					HighestEducation  string `json:"highest_education"`
					University        string `json:"university"`
					Experience        string `json:"experience"`
					Skills            string `json:"skills"`
					ResumeUrl         string `json:"resume_url"`
				}

				if err := c.ShouldBindJSON(&tempForm); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "参数格式错误"})
					return
				}

				userID, _ := c.Get("userID")

				req := &pb.UpdateProfileRequest{
					UserId:            userID.(int64),
					RealName:          tempForm.RealName,
					Phone:             tempForm.Phone,
					HighestEducation:  tempForm.HighestEducation,
					University:        tempForm.University,
					Experience:        tempForm.Experience,
					Skills:            tempForm.Skills,
					ResumeOssUrl:      tempForm.ResumeUrl,
				}

				resp, err := grpcClient.UpdateProfile(c.Request.Context(), req)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": "保存失败: " + err.Error()})
					return
				}
				c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": resp})
			})

			// 投递岗位
			candidateGroup.POST("/jobs/:id/apply", func(c *gin.Context) {
				jobIDStr := c.Param("id")
				var jobID int64
				fmt.Sscanf(jobIDStr, "%d", &jobID)

				userID, _ := c.Get("userID")

				req := &pb.ApplyJobRequest{
					JobId:  jobID,
					UserId: userID.(int64),
				}

				resp, err := grpcClient.ApplyJob(c.Request.Context(), req)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"code": 1, "msg": err.Error()})
					return
				}
				c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": resp})
			})
		}
	}

	fmt.Println("🌐 Gin 网关服务已启动，正在监听 HTTP 端口 :8080 ...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("网关启动失败: %v", err)
	}
}