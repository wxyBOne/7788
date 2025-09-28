// Package services 提供AI伙伴相关的业务逻辑服务
package services

import (
	"database/sql"
	"fmt"
	"seven-ai-backend/internal/models"
)

// CompanionService AI伙伴服务，处理AI伙伴的创建、管理和成长
type CompanionService struct {
	db        *sql.DB    // 数据库连接
	aiService *AIService // AI服务
}

// NewCompanionService 创建AI伙伴服务实例
func NewCompanionService(db *sql.DB, aiService *AIService) *CompanionService {
	return &CompanionService{
		db:        db,
		aiService: aiService,
	}
}

// CreateCompanion 创建AI伙伴
func (s *CompanionService) CreateCompanion(userID int, req models.CreateCompanionRequest) (*models.AICompanion, error) {
	// 检查用户是否已有AI伙伴
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM ai_companions WHERE user_id = ?", userID).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("检查AI伙伴数量失败: %v", err)
	}
	if count > 0 {
		return nil, fmt.Errorf("用户已有AI伙伴，无法重复创建")
	}

	// 创建AI伙伴
	query := `
		INSERT INTO ai_companions (
			user_id, name, avatar_url, personality_signature,
			conversation_fluency, knowledge_breadth, empathy_depth, 
			creativity_level, humor_sense, growth_mode, gender,
			created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	result, err := s.db.Exec(query,
		userID,
		req.Name,
		"/src/img/BlankAI.png", // 默认粒子小球头像
		"我...我是谁？你...你是谁？",     // 初始个性签名
		1, 1, 1, 1, 1,          // 初始能力值都是1
		req.GrowthMode,
		"unknown", // 初始性别未知
	)

	if err != nil {
		return nil, fmt.Errorf("创建AI伙伴失败: %v", err)
	}

	companionID, _ := result.LastInsertId()

	// 获取创建的AI伙伴信息
	companion, err := s.GetCompanion(int(companionID))
	if err != nil {
		return nil, fmt.Errorf("获取AI伙伴信息失败: %v", err)
	}

	return companion, nil
}

// GetUserCompanions 获取用户的AI伙伴列表
func (s *CompanionService) GetUserCompanions(userID int) ([]*models.AICompanion, error) {
	query := `
		SELECT id, user_id, name, avatar_url, personality_signature,
			   conversation_fluency, knowledge_breadth, empathy_depth,
			   creativity_level, humor_sense, total_experience, current_level,
			   growth_percentage, growth_mode, gender, voice_type,
			   personality_traits, learned_vocabulary, memory_summary,
			   is_growth_completed, last_active_at, created_at, updated_at
		FROM ai_companions 
		WHERE user_id = ?
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("查询AI伙伴列表失败: %v", err)
	}
	defer rows.Close()

	var companions []*models.AICompanion
	for rows.Next() {
		var companion models.AICompanion
		err := rows.Scan(
			&companion.ID, &companion.UserID, &companion.Name, &companion.AvatarURL,
			&companion.PersonalitySignature, &companion.ConversationFluency,
			&companion.KnowledgeBreadth, &companion.EmpathyDepth,
			&companion.CreativityLevel, &companion.HumorSense,
			&companion.TotalExperience, &companion.CurrentLevel,
			&companion.GrowthPercentage, &companion.GrowthMode,
			&companion.Gender, &companion.VoiceType,
			&companion.PersonalityTraits, &companion.LearnedVocabulary,
			&companion.MemorySummary, &companion.IsGrowthCompleted,
			&companion.LastActiveAt, &companion.CreatedAt, &companion.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("扫描AI伙伴数据失败: %v", err)
		}
		companions = append(companions, &companion)
	}

	return companions, nil
}

// GetCompanion 获取单个AI伙伴信息
func (s *CompanionService) GetCompanion(companionID int) (*models.AICompanion, error) {
	query := `
		SELECT id, user_id, name, avatar_url, personality_signature,
			   conversation_fluency, knowledge_breadth, empathy_depth,
			   creativity_level, humor_sense, total_experience, current_level,
			   growth_percentage, growth_mode, gender, voice_type,
			   personality_traits, learned_vocabulary, memory_summary,
			   is_growth_completed, last_active_at, created_at, updated_at
		FROM ai_companions 
		WHERE id = ?
	`

	var companion models.AICompanion
	err := s.db.QueryRow(query, companionID).Scan(
		&companion.ID, &companion.UserID, &companion.Name, &companion.AvatarURL,
		&companion.PersonalitySignature, &companion.ConversationFluency,
		&companion.KnowledgeBreadth, &companion.EmpathyDepth,
		&companion.CreativityLevel, &companion.HumorSense,
		&companion.TotalExperience, &companion.CurrentLevel,
		&companion.GrowthPercentage, &companion.GrowthMode,
		&companion.Gender, &companion.VoiceType,
		&companion.PersonalityTraits, &companion.LearnedVocabulary,
		&companion.MemorySummary, &companion.IsGrowthCompleted,
		&companion.LastActiveAt, &companion.CreatedAt, &companion.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("AI伙伴不存在")
		}
		return nil, fmt.Errorf("获取AI伙伴信息失败: %v", err)
	}

	return &companion, nil
}

// UpdateCompanion 更新AI伙伴信息
func (s *CompanionService) UpdateCompanion(companionID int, req models.UpdateCompanionRequest) error {
	query := `
		UPDATE ai_companions 
		SET name = ?, gender = ?, personality_traits = ?, 
			learned_vocabulary = ?, memory_summary = ?, updated_at = NOW()
		WHERE id = ?
	`

	_, err := s.db.Exec(query,
		req.Name, req.Gender, req.PersonalityTraits,
		req.LearnedVocabulary, req.MemorySummary, companionID,
	)

	if err != nil {
		return fmt.Errorf("更新AI伙伴信息失败: %v", err)
	}

	return nil
}


// GetGrowthStatus 获取AI伙伴成长状态
func (s *CompanionService) GetGrowthStatus(companionID int) (*models.GrowthProgressResponse, error) {
	companion, err := s.GetCompanion(companionID)
	if err != nil {
		return nil, err
	}

	// 计算下一个里程碑
	var nextMilestone string
	if companion.GrowthPercentage < 30 {
		nextMilestone = "学会基本对话"
	} else if companion.GrowthPercentage < 70 {
		nextMilestone = "形成稳定人格"
	} else if companion.GrowthPercentage < 100 {
		nextMilestone = "完全成熟"
	} else {
		nextMilestone = "成长完成"
	}

	return &models.GrowthProgressResponse{
		CompanionID:       companion.ID,
		GrowthPercentage:  companion.GrowthPercentage,
		CurrentLevel:      companion.CurrentLevel,
		TotalExperience:   companion.TotalExperience,
		IsGrowthCompleted: companion.IsGrowthCompleted,
		NextMilestone:     nextMilestone,
	}, nil
}

// GetDiary 获取AI伙伴日记
func (s *CompanionService) GetDiary(companionID int, limit int) ([]*models.CompanionDiary, error) {
	query := `
		SELECT id, companion_id, date, title, content, mood_score, 
			   is_user_mentioned, created_at
		FROM companion_diaries 
		WHERE companion_id = ?
		ORDER BY date DESC
		LIMIT ?
	`

	rows, err := s.db.Query(query, companionID, limit)
	if err != nil {
		return nil, fmt.Errorf("查询日记失败: %v", err)
	}
	defer rows.Close()

	var diaries []*models.CompanionDiary
	for rows.Next() {
		var diary models.CompanionDiary
		err := rows.Scan(
			&diary.ID, &diary.CompanionID, &diary.Date,
			&diary.Title, &diary.Content, &diary.MoodScore,
			&diary.IsUserMentioned, &diary.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("扫描日记数据失败: %v", err)
		}
		diaries = append(diaries, &diary)
	}

	return diaries, nil
}

// AddExperience 为AI伙伴添加经验值
func (s *CompanionService) AddExperience(companionID int, experience int) error {
	// 更新经验值
	query := `
		UPDATE ai_companions 
		SET total_experience = total_experience + ?, 
			current_level = FLOOR(total_experience / 100) + 1,
			growth_percentage = LEAST(100.0, (total_experience / 1000.0) * 100),
			updated_at = NOW()
		WHERE id = ?
	`

	_, err := s.db.Exec(query, experience, companionID)
	if err != nil {
		return fmt.Errorf("更新经验值失败: %v", err)
	}

	// 检查是否完成成长
	var growthPercentage float64
	err = s.db.QueryRow("SELECT growth_percentage FROM ai_companions WHERE id = ?", companionID).Scan(&growthPercentage)
	if err != nil {
		return fmt.Errorf("查询成长进度失败: %v", err)
	}

	if growthPercentage >= 100.0 {
		_, err = s.db.Exec("UPDATE ai_companions SET is_growth_completed = TRUE WHERE id = ?", companionID)
		if err != nil {
			return fmt.Errorf("更新成长完成状态失败: %v", err)
		}
	}

	return nil
}

// GetEmotionState 获取AI伙伴当前情绪状态（用于粒子小球外观）
func (s *CompanionService) GetEmotionState(companionID int) (*models.EmotionState, error) {
	// 这里可以根据AI伙伴的成长状态、最近对话等计算情绪
	// 暂时返回默认的平静状态
	return &models.EmotionState{
		Emotion:       "平静",
		Intensity:     0.5,
		Color:         "#52b4b4", // 柔绿色
		Brightness:    0.7,
		ParticleSpeed: 0.5,
	}, nil
}
