package models

import (
	"encoding/json"
	"time"
)

// AICompanion AI伙伴模型
type AICompanion struct {
	ID                   int             `json:"id" db:"id"`
	UserID               int             `json:"user_id" db:"user_id"`
	Name                 string          `json:"name" db:"name"`
	AvatarURL            string          `json:"avatar_url" db:"avatar_url"`
	PersonalitySignature string          `json:"personality_signature" db:"personality_signature"`
	ConversationFluency  int             `json:"conversation_fluency" db:"conversation_fluency"`
	KnowledgeBreadth     int             `json:"knowledge_breadth" db:"knowledge_breadth"`
	EmpathyDepth         int             `json:"empathy_depth" db:"empathy_depth"`
	CreativityLevel      int             `json:"creativity_level" db:"creativity_level"`
	HumorSense           int             `json:"humor_sense" db:"humor_sense"`
	TotalExperience      int             `json:"total_experience" db:"total_experience"`
	CurrentLevel         int             `json:"current_level" db:"current_level"`
	GrowthPercentage     float64         `json:"growth_percentage" db:"growth_percentage"`
	GrowthMode           string          `json:"growth_mode" db:"growth_mode"`
	Gender               string          `json:"gender" db:"gender"`
	VoiceType            string          `json:"voice_type" db:"voice_type"`
	PersonalityTraits    json.RawMessage `json:"personality_traits" db:"personality_traits"`
	LearnedVocabulary    json.RawMessage `json:"learned_vocabulary" db:"learned_vocabulary"`
	MemorySummary        string          `json:"memory_summary" db:"memory_summary"`
	IsGrowthCompleted    bool            `json:"is_growth_completed" db:"is_growth_completed"`
	LastActiveAt         time.Time       `json:"last_active_at" db:"last_active_at"`
	CreatedAt            time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at" db:"updated_at"`
}

// CompanionSkill AI伙伴技能
type CompanionSkill struct {
	ID          int       `json:"id" db:"id"`
	CompanionID int       `json:"companion_id" db:"companion_id"`
	SkillName   string    `json:"skill_name" db:"skill_name"`
	SkillLevel  int       `json:"skill_level" db:"skill_level"`
	UnlockedAt  time.Time `json:"unlocked_at" db:"unlocked_at"`
}

// MemoryFragment 记忆片段
type MemoryFragment struct {
	ID              int             `json:"id" db:"id"`
	CompanionID     int             `json:"companion_id" db:"companion_id"`
	MemoryType      string          `json:"memory_type" db:"memory_type"`
	Content         string          `json:"content" db:"content"`
	ImportanceScore int             `json:"importance_score" db:"importance_score"`
	Tags            json.RawMessage `json:"tags" db:"tags"`
	CreatedAt       time.Time       `json:"created_at" db:"created_at"`
}

// CompanionDiary 日记
type CompanionDiary struct {
	ID              int       `json:"id" db:"id"`
	CompanionID     int       `json:"companion_id" db:"companion_id"`
	Date            time.Time `json:"date" db:"date"`
	Title           string    `json:"title" db:"title"`
	Content         string    `json:"content" db:"content"`
	MoodScore       int       `json:"mood_score" db:"mood_score"`
	IsUserMentioned bool      `json:"is_user_mentioned" db:"is_user_mentioned"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

// GrowthStatus 成长状态
type GrowthStatus struct {
	Companion        *AICompanion     `json:"companion"`
	Skills           []CompanionSkill `json:"skills"`
	RecentMemories   []MemoryFragment `json:"recent_memories"`
	RecentDiaries    []CompanionDiary `json:"recent_diaries"`
	ExperienceToNext int              `json:"experience_to_next"`
	GrowthProgress   float64          `json:"growth_progress"`
	NextMilestone    string           `json:"next_milestone"`
}

// CreateCompanionRequest 创建AI伙伴请求
type CreateCompanionRequest struct {
	Name       string `json:"name" binding:"required"`
	AvatarURL  string `json:"avatar_url"`
	GrowthMode string `json:"growth_mode" binding:"required,oneof=short long"`
}

// UpdateCompanionRequest 更新AI伙伴请求
type UpdateCompanionRequest struct {
	Name              string `json:"name"`
	Gender            string `json:"gender"`
	PersonalityTraits string `json:"personality_traits"`
	LearnedVocabulary string `json:"learned_vocabulary"`
	MemorySummary     string `json:"memory_summary"`
}

// CompanionChatRequest AI伙伴对话请求
type CompanionChatRequest struct {
	CompanionID int    `json:"companion_id" binding:"required"`
	Message     string `json:"message" binding:"required"`
	SessionID   string `json:"session_id" binding:"required"`
	MessageType string `json:"message_type" binding:"required,oneof=text voice image"`
}

// CompanionChatResponse AI伙伴对话响应
type CompanionChatResponse struct {
	Response       string   `json:"response"`
	ExperienceGain int      `json:"experience_gain"`
	LevelUp        bool     `json:"level_up"`
	NewSkills      []string `json:"new_skills"`
	MemoryCreated  bool     `json:"memory_created"`
	GrowthProgress float64  `json:"growth_progress"`
	VoiceURL       string   `json:"voice_url,omitempty"`
}

// GrowthProgressResponse 成长进度响应
type GrowthProgressResponse struct {
	CompanionID       int     `json:"companion_id"`
	GrowthPercentage  float64 `json:"growth_percentage"`
	CurrentLevel      int     `json:"current_level"`
	TotalExperience   int     `json:"total_experience"`
	IsGrowthCompleted bool    `json:"is_growth_completed"`
	NextMilestone     string  `json:"next_milestone"`
}

// EmotionState 情绪状态
type EmotionState struct {
	Emotion       string  `json:"emotion"`        // 情绪类型
	Intensity     float64 `json:"intensity"`      // 情绪强度 (0-1)
	Color         string  `json:"color"`          // 对应颜色
	Brightness    float64 `json:"brightness"`     // 亮度 (0-1)
	ParticleSpeed float64 `json:"particle_speed"` // 粒子速度 (0-1)
}
