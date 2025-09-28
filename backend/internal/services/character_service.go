// Package services 提供角色相关的业务逻辑服务
package services

import (
	"database/sql"
	"fmt"
	"seven-ai-backend/internal/models"
	"strings"
)

// CharacterService 角色服务，处理预设角色的查询和管理
type CharacterService struct {
	db *sql.DB // 数据库连接
}

// NewCharacterService 创建角色服务实例
func NewCharacterService(db *sql.DB) *CharacterService {
	return &CharacterService{db: db}
}

// GetAllCharacters 获取所有预设角色列表
func (s *CharacterService) GetAllCharacters() ([]models.CharacterResponse, error) {
	rows, err := s.db.Query(`
		SELECT id, name, description, avatar_url, personality_signature, 
		       personality_traits, background_story, voice_settings, 
		       system_prompt, search_keywords, skills
		FROM preset_characters
		ORDER BY id
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query characters: %w", err)
	}
	defer rows.Close()

	var characters []models.CharacterResponse
	for rows.Next() {
		var char models.CharacterResponse
		var voiceSettings sql.NullString
		err := rows.Scan(
			&char.ID, &char.Name, &char.Description, &char.AvatarURL,
			&char.PersonalitySignature, &char.PersonalityTraits,
			&char.BackgroundStory, &voiceSettings,
			&char.SystemPrompt, &char.SearchKeywords, &char.Skills,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan character: %w", err)
		}

		// 处理NULL值
		if voiceSettings.Valid {
			char.VoiceSettings = voiceSettings.String
		} else {
			char.VoiceSettings = ""
		}

		characters = append(characters, char)
	}

	return characters, nil
}

// GetCharacterByID 根据ID获取单个角色信息
func (s *CharacterService) GetCharacterByID(characterID int) (*models.CharacterResponse, error) {
	var char models.CharacterResponse
	var voiceSettings sql.NullString
	err := s.db.QueryRow(`
		SELECT id, name, description, avatar_url, personality_signature,
		       personality_traits, background_story, voice_settings,
		       system_prompt, search_keywords, skills
		FROM preset_characters WHERE id = ?
	`, characterID).Scan(
		&char.ID, &char.Name, &char.Description, &char.AvatarURL,
		&char.PersonalitySignature, &char.PersonalityTraits,
		&char.BackgroundStory, &voiceSettings,
		&char.SystemPrompt, &char.SearchKeywords, &char.Skills,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("character not found")
		}
		return nil, fmt.Errorf("failed to query character: %w", err)
	}

	// 处理NULL值
	if voiceSettings.Valid {
		char.VoiceSettings = voiceSettings.String
	} else {
		char.VoiceSettings = ""
	}

	return &char, nil
}

// SearchCharacters 根据关键词搜索角色
func (s *CharacterService) SearchCharacters(keyword string) ([]models.CharacterResponse, error) {
	searchPattern := "%" + strings.ToLower(keyword) + "%"
	rows, err := s.db.Query(`
		SELECT id, name, description, avatar_url, personality_signature,
		       personality_traits, background_story, voice_settings,
		       system_prompt, search_keywords, skills
		FROM preset_characters 
		WHERE LOWER(name) LIKE ? OR LOWER(description) LIKE ? OR LOWER(search_keywords) LIKE ?
		ORDER BY id
	`, searchPattern, searchPattern, searchPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to search characters: %w", err)
	}
	defer rows.Close()

	var characters []models.CharacterResponse
	for rows.Next() {
		var char models.CharacterResponse
		var voiceSettings sql.NullString
		err := rows.Scan(
			&char.ID, &char.Name, &char.Description, &char.AvatarURL,
			&char.PersonalitySignature, &char.PersonalityTraits,
			&char.BackgroundStory, &voiceSettings,
			&char.SystemPrompt, &char.SearchKeywords, &char.Skills,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan character: %w", err)
		}

		// 处理NULL值
		if voiceSettings.Valid {
			char.VoiceSettings = voiceSettings.String
		} else {
			char.VoiceSettings = ""
		}

		characters = append(characters, char)
	}

	return characters, nil
}
