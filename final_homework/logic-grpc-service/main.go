package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"

	"logic-grpc-service/db"
	"logic-grpc-service/pb"
	"logic-grpc-service/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
)

type server struct {
	pb.UnimplementedRecruitmentServiceServer
}

// 真实登录逻辑
func (s *server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	var user db.User

	// 去数据库里查这个用户名
	result := db.DB.Where("username = ?", req.Username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, fmt.Errorf("数据库查询错误")
	}

	// 校验密码
	if user.PasswordHash != req.Password {
		return nil, fmt.Errorf("密码错误")
	}

	fmt.Printf("✅ 数据库验证通过，用户登录成功: %s (ID: %d, 角色: %s)\n", user.Username, user.ID, user.Role)

	// 返回真实数据
	return &pb.LoginResponse{
		Token:  "jwt-token-to-be-generated-in-gin", // Gin层后续会负责生成真实的JWT
		Role:   user.Role,
		UserId: user.ID,
	}, nil
}

// 真实发布岗位逻辑
func (s *server) CreateJob(ctx context.Context, req *pb.CreateJobRequest) (*pb.CommonResponse, error) {
	// 将 proto 的请求体映射到我们的 GORM 模型
	newJob := db.Job{
		HrID:        req.HrId,
		Title:       req.Title,
		Description: req.Description,
		Status:      1, // 默认上架
	}

	// 写入数据库
	if err := db.DB.Create(&newJob).Error; err != nil {
		fmt.Printf("❌ 岗位发布失败: %v\n", err)
		return &pb.CommonResponse{Code: 1, Msg: "数据库写入失败"}, nil
	}

	fmt.Printf("📢 岗位发布成功: [%s]\n", newJob.Title)
	return &pb.CommonResponse{Code: 0, Msg: "发布成功"}, nil
}

// 查询候选人个人档案
func (s *server) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	var profile db.CandidateProfile
	result := db.DB.Where("user_id = ?", req.UserId).First(&profile)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// 档案不存在，返回空响应（不报错）
			return &pb.GetProfileResponse{}, nil
		}
		return nil, fmt.Errorf("数据库查询错误")
	}

	fmt.Printf("📖 查询候选人档案 (ID: %d): %s\n", req.UserId, profile.RealName)
	return &pb.GetProfileResponse{
		UserId:           profile.UserID,
		RealName:         profile.RealName,
		Phone:            profile.Phone,
		HighestEducation: profile.HighestEducation,
		University:       profile.University,
		Experience:       profile.Experience,
		Skills:           profile.Skills,
		ResumeOssUrl:     profile.ResumeOssUrl,
	}, nil
}

// 候选人更新个人档案
func (s *server) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.CommonResponse, error) {
	// 查找是否已存在档案
	var profile db.CandidateProfile
	result := db.DB.Where("user_id = ?", req.UserId).First(&profile)

	profile.UserID = req.UserId
	profile.RealName = req.RealName
	profile.Phone = req.Phone
	profile.HighestEducation = req.HighestEducation // 从前端传入的真实学历
	profile.University = req.University
	profile.Experience = req.Experience
	profile.Skills = req.Skills // 从前端传入的真实技能标签
	profile.ResumeOssUrl = req.ResumeOssUrl // 核心：只存 OSS 相对路径，不存文件实体！

	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// 不存在则创建
		db.DB.Create(&profile)
	} else {
		// 存在则更新
		db.DB.Save(&profile)
	}

	fmt.Printf("📝 候选人 (ID: %d) 档案已更新，简历路径: %s\n", req.UserId, req.ResumeOssUrl)
	return &pb.CommonResponse{Code: 0, Msg: "档案保存成功"}, nil
}

// 候选人投递岗位
func (s *server) ApplyJob(ctx context.Context, req *pb.ApplyJobRequest) (*pb.CommonResponse, error) {
	// 投递强制校验规则：必须有档案且有简历路径
	var profile db.CandidateProfile
	if err := db.DB.Where("user_id = ?", req.UserId).First(&profile).Error; err != nil {
		return &pb.CommonResponse{Code: 1, Msg: "请先完善个人档案"}, nil
	}
	if profile.ResumeOssUrl == "" {
		return &pb.CommonResponse{Code: 1, Msg: "请先上传合规简历至 OSS"}, nil
	}

	// 写入投递记录表
	application := db.JobApplication{
		JobID:       req.JobId,
		CandidateID: req.UserId,
	}
	if err := db.DB.Create(&application).Error; err != nil {
		return &pb.CommonResponse{Code: 1, Msg: "投递失败，请稍后重试"}, nil
	}

	fmt.Printf("✅ 候选人 (ID: %d) 成功投递岗位 (ID: %d)\n", req.UserId, req.JobId)
	return &pb.CommonResponse{Code: 0, Msg: "投递成功"}, nil
}

// 简历上传到 OSS
func (s *server) UploadResume(ctx context.Context, req *pb.UploadResumeRequest) (*pb.UploadResumeResponse, error) {
	// 生成 OSS 存储路径: resumes/user_{id}_{timestamp}_{filename}
	objectName := fmt.Sprintf("resumes/user_%d_%s", req.UserId, req.Filename)

	// 调用 OSS SDK 上传
	err := utils.UploadResume(objectName, req.Content)
	if err != nil {
		fmt.Printf("❌ 简历上传失败: %v\n", err)
		return &pb.UploadResumeResponse{Code: 1, Msg: "上传失败"}, nil
	}

	fmt.Printf("☁️ 候选人(ID:%d) 简历上传成功: %s\n", req.UserId, objectName)
	return &pb.UploadResumeResponse{Code: 0, Msg: "上传成功", OssUrl: objectName}, nil
}


// HR 查看岗位投递详情
func (s *server) GetJobApplications(ctx context.Context, req *pb.GetApplicationsRequest) (*pb.ApplicationsResponse, error) {
	var applications []db.JobApplication
	// 查询该岗位的所有投递记录
	db.DB.Where("job_id = ?", req.JobId).Find(&applications)

	var profiles []*pb.CandidateProfile
	for _, app := range applications {
		var p db.CandidateProfile
		db.DB.Where("user_id = ?", app.CandidateID).First(&p)

		// ⭐️ 核心安全机制：动态生成 10 分钟有效期的签名 URL，防盗链，保护隐私
		signedUrl, err := utils.GeneratePresignedURL(p.ResumeOssUrl)
		if err != nil {
			fmt.Printf("生成签名链接失败: %v\n", err)
			signedUrl = "生成链接失败"
		}

		// 组装返回给前端的数据
		profiles = append(profiles, &pb.CandidateProfile{
			RealName:         p.RealName,
			Phone:            p.Phone,
			University:       p.University,
			HighestEducation: p.HighestEducation,
			Skills:           p.Skills,
			ResumeUrl:        signedUrl, // 返回给前端的是可以直接点击下载的加密链接
		})
	}

	return &pb.ApplicationsResponse{Profiles: profiles}, nil
}

// 查询HR发布的岗位列表
func (s *server) ListHRJobs(ctx context.Context, req *pb.ListHRJobsRequest) (*pb.ListJobsResponse, error) {
	var jobs []db.Job

	if req.HrId == 0 {
		db.DB.Order("created_at desc").Find(&jobs)
	} else {
		db.DB.Where("hr_id = ?", req.HrId).Order("created_at desc").Find(&jobs)
	}

	var pbJobs []*pb.JobInfo
	for _, j := range jobs {
		pbJobs = append(pbJobs, &pb.JobInfo{
			Id:          j.ID,
			Title:       j.Title,
			Description: j.Description,
			HrId:        j.HrID,
			Status:      j.Status,
		})
	}
	return &pb.ListJobsResponse{Jobs: pbJobs}, nil
}

// 编辑/下架岗位
func (s *server) UpdateJob(ctx context.Context, req *pb.UpdateJobRequest) (*pb.CommonResponse, error) {
	var job db.Job

	// 查询岗位是否存在
	result := db.DB.First(&job, req.JobId)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return &pb.CommonResponse{Code: 1, Msg: "岗位不存在"}, nil
		}
		return &pb.CommonResponse{Code: 1, Msg: "数据库查询错误"}, nil
	}

	// 校验岗位归属权限：只有岗位创建者才能编辑/下架
	if job.HrID != req.HrId {
		return &pb.CommonResponse{Code: 1, Msg: "无权限操作此岗位"}, nil
	}

	// 更新字段
	if req.Title != "" {
		job.Title = req.Title
	}
	if req.Description != "" {
		job.Description = req.Description
	}
	job.Status = req.Status

	// 保存更新
	if err := db.DB.Save(&job).Error; err != nil {
		fmt.Printf("❌ 岗位更新失败: %v\n", err)
		return &pb.CommonResponse{Code: 1, Msg: "更新失败"}, nil
	}

	action := "更新"
	if req.Status == 0 {
		action = "下架"
	}
	fmt.Printf("📝 岗位(ID:%d) %s 成功\n", req.JobId, action)
	return &pb.CommonResponse{Code: 0, Msg: action + "成功"}, nil
}

// ==================== 数据库上下文构建函数 ====================

// BuildDatabaseContext 根据用户问题构建数据库上下文
func BuildDatabaseContext(hrID int64, question string) string {
	var contextBuilder strings.Builder
	intent := utils.DetectIntent(question)

	contextBuilder.WriteString("=== 智能招聘系统数据库实时数据 ===\n\n")

	// 1. 岗位数据（HR自己发布的）
	var jobs []db.Job
	db.DB.Where("hr_id = ?", hrID).Order("created_at desc").Find(&jobs)

	activeCount := 0
	pausedCount := 0
	for _, j := range jobs {
		if j.Status == 1 {
			activeCount++
		} else {
			pausedCount++
		}
	}

	contextBuilder.WriteString(fmt.Sprintf("【岗位统计】\n"))
	contextBuilder.WriteString(fmt.Sprintf("- 您发布的岗位总数: %d\n", len(jobs)))
	contextBuilder.WriteString(fmt.Sprintf("- 正在招聘的岗位: %d\n", activeCount))
	contextBuilder.WriteString(fmt.Sprintf("- 已下架的岗位: %d\n", pausedCount))

	if len(jobs) > 0 {
		contextBuilder.WriteString("\n【岗位详情列表】\n")
		for i, j := range jobs {
			if i >= 10 { // 最多列出10个
				contextBuilder.WriteString(fmt.Sprintf("... 还有 %d 个岗位\n", len(jobs)-10))
				break
			}
			status := "招聘中"
			if j.Status == 0 {
				status = "已下架"
			}
			contextBuilder.WriteString(fmt.Sprintf("- [%d] %s (%s) - 创建于 %s\n", j.ID, j.Title, status, j.CreatedAt.Format("2006-01-02")))
		}
	}

	// 2. 投递数据
	var allApplications []db.JobApplication

	// 获取该HR所有岗位的投递
	jobIDs := make([]int64, len(jobs))
	for i, j := range jobs {
		jobIDs[i] = j.ID
	}
	if len(jobIDs) > 0 {
		db.DB.Where("job_id IN ?", jobIDs).Order("created_at desc").Find(&allApplications)
	}

	contextBuilder.WriteString(fmt.Sprintf("\n【投递统计】\n"))
	contextBuilder.WriteString(fmt.Sprintf("- 总投递次数: %d\n", len(allApplications)))

	// 统计每个岗位的投递数
	if len(jobs) > 0 && len(allApplications) > 0 {
		contextBuilder.WriteString("\n【各岗位投递情况】\n")
		for _, j := range jobs {
			count := 0
			for _, a := range allApplications {
				if a.JobID == j.ID {
					count++
				}
			}
			if count > 0 {
				contextBuilder.WriteString(fmt.Sprintf("- %s: %d 次投递\n", j.Title, count))
			}
		}
	}

	// 3. 候选人数据（投递过该HR岗位的候选人）
	if len(allApplications) > 0 {
		candidateIDs := make(map[int64]bool)
		for _, a := range allApplications {
			candidateIDs[a.CandidateID] = true
		}

		contextBuilder.WriteString(fmt.Sprintf("\n【候选人统计】\n"))
		contextBuilder.WriteString(fmt.Sprintf("- 投递过您岗位的候选人数量: %d\n", len(candidateIDs)))

		// 列出候选人详情（最多5个）
		if len(candidateIDs) > 0 {
			contextBuilder.WriteString("\n【候选人信息示例】（最多显示5位）\n")
			count := 0
			for cid := range candidateIDs {
				if count >= 5 {
					break
				}
				var profile db.CandidateProfile
				if err := db.DB.Where("user_id = ?", cid).First(&profile).Error; err == nil {
					contextBuilder.WriteString(fmt.Sprintf("- 姓名: %s, 院校: %s, 学历: %s, 技能: %s\n",
						profile.RealName, profile.University, profile.HighestEducation, profile.Skills))
				}
				count++
			}
			if len(candidateIDs) > 5 {
				contextBuilder.WriteString(fmt.Sprintf("... 还有 %d 位候选人\n", len(candidateIDs)-5))
			}
		}
	}

	// 4. 根据意图补充特定数据
	if intent == "stats" || strings.Contains(question, "统计") || strings.Contains(question, "数据") {
		contextBuilder.WriteString("\n【综合数据分析】\n")
		if len(allApplications) > 0 && len(jobs) > 0 {
			avgApplications := float64(len(allApplications)) / float64(len(jobs))
			contextBuilder.WriteString(fmt.Sprintf("- 平均每个岗位投递数: %.1f\n", avgApplications))
		}

		// 今日投递（模拟）
		todayCount := 0
		for _, a := range allApplications {
			if a.CreatedAt.Format("2006-01-02") == "2024-01-02" { // 可以替换为当前日期
				todayCount++
			}
		}
		contextBuilder.WriteString(fmt.Sprintf("- 今日新增投递: %d\n", todayCount))
	}

	// 5. 历史对话上下文（如果用户提到之前的对话）
	if intent == "general" {
		var recentHistory []db.AiChatHistory
		db.DB.Where("hr_id = ?", hrID).Order("created_at desc").Limit(3).Find(&recentHistory)

		if len(recentHistory) > 0 {
			contextBuilder.WriteString("\n【最近的对话记录】\n")
			for _, h := range recentHistory {
				contextBuilder.WriteString(fmt.Sprintf("- 问: %s\n  答: %s\n", h.Question, h.Answer))
			}
		}
	}

	contextBuilder.WriteString("\n=== 以上是数据库真实数据，请据此回答用户问题 ===\n")

	return contextBuilder.String()
}

// ==================== AI 智能问答（核心增强） ====================

// AI 智能问答（含历史持久化 + 数据库上下文）
func (s *server) ChatWithAI(ctx context.Context, req *pb.AIChatRequest) (*pb.AIChatResponse, error) {
	// 1. 构建数据库上下文（核心：让 AI 了解真实数据）
	dbContext := BuildDatabaseContext(req.HrId, req.Message)

	fmt.Printf("📊 为 HR(ID:%d) 构建了数据库上下文，长度: %d 字符\n", req.HrId, len(dbContext))

	// 2. 获取历史对话上下文（最近5轮）
	var histories []db.AiChatHistory
	db.DB.Where("hr_id = ?", req.HrId).Order("created_at desc").Limit(5).Find(&histories)

	var historyContext strings.Builder
	if len(histories) > 0 {
		historyContext.WriteString("\n【历史对话记录】\n")
		for i := len(histories) - 1; i >= 0; i-- { // 按时间正序
			h := histories[i]
			historyContext.WriteString(fmt.Sprintf("HR: %s\nAI: %s\n", h.Question, h.Answer))
		}
	}

	// 3. 合成完整上下文
	fullContext := dbContext + historyContext.String()

	// 4. 调用 Eino 大模型生成回复
	answer, err := utils.ChatWithAI(ctx, req.Message, fullContext)
	if err != nil {
		fmt.Printf("❌ AI 生成失败: %v\n", err)
		return &pb.AIChatResponse{Reply: "抱歉，AI 暂时无法回答，请稍后再试。"}, nil
	}

	// 5. 将本次问答存入数据库（持久化）
	chatRecord := db.AiChatHistory{
		HrID:     req.HrId,
		Question: req.Message,
		Answer:   answer,
	}
	db.DB.Create(&chatRecord)

	fmt.Printf("🤖 HR(ID:%d) AI问答完成，问题: %s\n", req.HrId, req.Message)
	return &pb.AIChatResponse{Reply: answer}, nil
}

// 获取 HR 的历史对话记录
func (s *server) GetChatHistory(ctx context.Context, req *pb.GetChatHistoryRequest) (*pb.GetChatHistoryResponse, error) {
	var histories []db.AiChatHistory
	db.DB.Where("hr_id = ?", req.HrId).Order("created_at asc").Find(&histories)

	var messages []*pb.ChatMessage
	for _, h := range histories {
		// 用户提问
		messages = append(messages, &pb.ChatMessage{
			Role:      "user",
			Content:   h.Question,
			CreatedAt: h.CreatedAt.Unix(),
		})
		// AI 回复
		messages = append(messages, &pb.ChatMessage{
			Role:      "assistant",
			Content:   h.Answer,
			CreatedAt: h.CreatedAt.Unix(),
		})
	}

	fmt.Printf("📜 加载 HR(ID:%d) 历史对话 %d 条\n", req.HrId, len(histories))
	return &pb.GetChatHistoryResponse{Messages: messages}, nil
}


func main() {
	db.InitDB()
	utils.InitOSS()
	utils.InitEino() // 初始化 Eino AI 引擎

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("端口监听失败: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterRecruitmentServiceServer(s, &server{})
	reflection.Register(s)

	fmt.Println("🚀 Logic gRPC 服务已启动，正在监听端口 :50051 ...")
	fmt.Println("🧠 AI 助手已就绪，可查询真实数据库数据")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("启动 gRPC 服务失败: %v", err)
	}
}