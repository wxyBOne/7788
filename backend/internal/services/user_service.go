package services

import (
	"database/sql"
	"fmt"
	"seven-ai-backend/internal/models"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	db        *sql.DB
	aiService *AIService
}

func NewUserService(db *sql.DB, aiService *AIService) *UserService {
	return &UserService{db: db, aiService: aiService}
}

func (s *UserService) CreateUser(req models.UserCreateRequest) (*models.UserResponse, error) {
	// 检查用户是否已存在
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", req.Email).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if count > 0 {
		return nil, fmt.Errorf("邮箱已存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 生成唯一的用户名
	username := s.generateUniqueUsername(req.Email)

	// 创建用户
	result, err := s.db.Exec(`
		INSERT INTO users (username, email, password_hash, created_at, updated_at) 
		VALUES (?, ?, ?, NOW(), NOW())
	`, username, req.Email, string(hashedPassword))
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get user ID: %w", err)
	}

	// 创建用户偏好设置
	_, err = s.db.Exec(`
		INSERT INTO user_preferences (user_id, voice_enabled, auto_save_memories, notification_enabled, language_preference)
		VALUES (?, TRUE, TRUE, TRUE, 'zh-CN')
	`, userID)
	if err != nil {
		// 记录错误但不影响用户创建
		fmt.Printf("Failed to create user preferences: %v\n", err)
	}

	// 自动添加赫敏为好友（character_id = 4）
	_, err = s.db.Exec(`
		INSERT INTO user_friendships (user_id, character_id, is_active, created_at, updated_at)
		VALUES (?, 4, TRUE, NOW(), NOW())
	`, userID)
	if err != nil {
		// 记录错误但不影响用户创建
		fmt.Printf("Failed to add Hermione as friend: %v\n", err)
	} else {
		// 生成赫敏的欢迎消息
		err = s.generateWelcomeMessage(int(userID), 4)
		if err != nil {
			fmt.Printf("Failed to generate welcome message: %v\n", err)
		}
	}

	// 返回用户信息
	user := &models.UserResponse{
		ID:        int(userID),
		Username:  username,
		Email:     req.Email,
		AvatarURL: "",
	}

	return user, nil
}

func (s *UserService) AuthenticateUser(req models.UserLoginRequest) (*models.UserResponse, error) {
	var user models.User
	err := s.db.QueryRow(`
		SELECT id, username, email, password_hash, avatar_url, created_at 
		FROM users WHERE email = ?
	`, req.Email).Scan(
		&user.ID, &user.Username, &user.Email,
		&user.PasswordHash, &user.AvatarURL, &user.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("密码错误")
	}

	// 处理AvatarURL的NULL值
	avatarURL := ""
	if user.AvatarURL.Valid {
		avatarURL = user.AvatarURL.String
	}

	return &models.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		AvatarURL: avatarURL,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *UserService) GetUserByID(userID int) (*models.UserResponse, error) {
	var user models.User
	err := s.db.QueryRow(`
		SELECT id, username, email, avatar_url, created_at 
		FROM users WHERE id = ?
	`, userID).Scan(
		&user.ID, &user.Username, &user.Email,
		&user.AvatarURL, &user.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	// 处理AvatarURL的NULL值
	avatarURL := ""
	if user.AvatarURL.Valid {
		avatarURL = user.AvatarURL.String
	}

	return &models.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		AvatarURL: avatarURL,
		CreatedAt: user.CreatedAt,
	}, nil
}

// SendResetCode 生成图片验证码
func (s *UserService) SendResetCode(email string) (*CaptchaData, error) {
	// 检查用户是否存在
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if count == 0 {
		return nil, fmt.Errorf("用户不存在")
	}

	// 生成图片验证码
	captchaService := NewCaptchaService()
	captchaData, err := captchaService.GenerateCaptcha()
	if err != nil {
		return nil, fmt.Errorf("failed to generate captcha: %w", err)
	}

	// TODO: 实际项目中应该：
	// 1. 将验证码存储到Redis或数据库，设置过期时间（5分钟）
	// 2. 限制发送频率

	return captchaData, nil
}

// ResetPassword 重置密码
func (s *UserService) ResetPassword(email, verificationCode, newPassword string) error {
	// 验证验证码（这里简化处理，实际应该从Redis或数据库验证）
	if verificationCode != "123456" { // 临时固定验证码
		return fmt.Errorf("验证码错误")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// 更新密码
	_, err = s.db.Exec(`
		UPDATE users 
		SET password_hash = ?, updated_at = NOW()
		WHERE email = ?
	`, string(hashedPassword), email)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// generateVerificationCode 生成6位数字验证码
func generateVerificationCode() string {
	// 简化实现，实际项目中应该使用更安全的随机数生成
	return "123456"
}

// generateWelcomeMessage 生成AI欢迎消息
func (s *UserService) generateWelcomeMessage(userID int, characterID int) error {
	// 获取角色信息
	var character models.CharacterResponse
	var voiceSettings sql.NullString
	err := s.db.QueryRow(`
		SELECT id, name, description, avatar_url, personality_signature,
		       personality_traits, background_story, voice_settings,
		       system_prompt, search_keywords
		FROM preset_characters WHERE id = ?
	`, characterID).Scan(
		&character.ID, &character.Name, &character.Description, &character.AvatarURL,
		&character.PersonalitySignature, &character.PersonalityTraits,
		&character.BackgroundStory, &voiceSettings,
		&character.SystemPrompt, &character.SearchKeywords,
	)
	if err != nil {
		return fmt.Errorf("failed to get character: %w", err)
	}

	// 处理NULL值
	if voiceSettings.Valid {
		character.VoiceSettings = voiceSettings.String
	} else {
		character.VoiceSettings = ""
	}

	// 构建欢迎消息提示
	welcomePrompt := fmt.Sprintf(`
%s

现在用户刚刚添加你为好友，请以%s的身份说一句欢迎的话。要求：
1. 符合角色的性格特点
2. 简短自然，不超过50字
3. 体现角色的个性特征
4. 不要过于热情或冷淡，保持角色设定

请直接回复欢迎消息，不要包含任何解释。
`, character.SystemPrompt, character.Name)

	// 调用AI生成欢迎消息
	messages := []Message{
		{Role: "system", Content: welcomePrompt},
	}

	response, err := s.aiService.ChatWithLLM(messages, "qwen3-max", 0.8, "text")
	if err != nil {
		fmt.Printf("AI call failed for user %d, character %d: %v\n", userID, characterID, err)
		return fmt.Errorf("failed to generate welcome message: %w", err)
	}

	fmt.Printf("AI response for user %d, character %d: %s\n", userID, characterID, response)

	// 后处理AI响应，移除角色名字前缀
	if strings.HasPrefix(response, character.Name+"。") {
		response = strings.TrimPrefix(response, character.Name+"。")
	} else if strings.HasPrefix(response, character.Name) {
		response = strings.TrimPrefix(response, character.Name)
	}
	response = strings.TrimSpace(response) // 移除可能存在的首尾空格

	// 保存欢迎消息到数据库
	sessionID := fmt.Sprintf("welcome_%d_%d_%d", userID, characterID, time.Now().Unix())
	_, err = s.db.Exec(`
		INSERT INTO conversations 
		(user_id, character_id, session_id, message_type, user_message, ai_response, is_ai_initiated, created_at)
		VALUES (?, ?, ?, 'text', '', ?, true, NOW())
	`, userID, characterID, sessionID, response)
	if err != nil {
		return fmt.Errorf("failed to save welcome message: %w", err)
	}

	return nil
}

// generateUniqueUsername 生成唯一的用户名
func (s *UserService) generateUniqueUsername(email string) string {
	// 使用邮箱前缀作为基础用户名
	baseUsername := strings.Split(email, "@")[0]
	username := baseUsername

	// 检查用户名是否已存在，如果存在则添加随机数
	for {
		var count int
		err := s.db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&count)
		if err != nil {
			break
		}
		if count == 0 {
			break
		}
		// 用户名已存在，添加随机数
		username = fmt.Sprintf("%s%d", baseUsername, time.Now().Unix()%10000)
	}

	return username
}
