package services

import (
	"database/sql"
	"fmt"
	"seven-ai-backend/internal/models"
	"strings"
)

type CharacterService struct {
	db *sql.DB
}

func NewCharacterService(db *sql.DB) *CharacterService {
	return &CharacterService{db: db}
}

func (s *CharacterService) GetAllCharacters() ([]models.CharacterResponse, error) {
	rows, err := s.db.Query(`
		SELECT id, name, description, avatar_url, personality_signature, 
		       personality_traits, background_story, voice_settings, 
		       system_prompt, search_keywords
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
		err := rows.Scan(
			&char.ID, &char.Name, &char.Description, &char.AvatarURL,
			&char.PersonalitySignature, &char.PersonalityTraits,
			&char.BackgroundStory, &char.VoiceSettings,
			&char.SystemPrompt, &char.SearchKeywords,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan character: %w", err)
		}
		characters = append(characters, char)
	}

	return characters, nil
}

func (s *CharacterService) GetCharacterByID(characterID int) (*models.CharacterResponse, error) {
	var char models.CharacterResponse
	err := s.db.QueryRow(`
		SELECT id, name, description, avatar_url, personality_signature,
		       personality_traits, background_story, voice_settings,
		       system_prompt, search_keywords
		FROM preset_characters WHERE id = ?
	`, characterID).Scan(
		&char.ID, &char.Name, &char.Description, &char.AvatarURL,
		&char.PersonalitySignature, &char.PersonalityTraits,
		&char.BackgroundStory, &char.VoiceSettings,
		&char.SystemPrompt, &char.SearchKeywords,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("character not found")
		}
		return nil, fmt.Errorf("failed to query character: %w", err)
	}

	return &char, nil
}

func (s *CharacterService) SearchCharacters(keyword string) ([]models.CharacterResponse, error) {
	searchPattern := "%" + strings.ToLower(keyword) + "%"
	rows, err := s.db.Query(`
		SELECT id, name, description, avatar_url, personality_signature,
		       personality_traits, background_story, voice_settings,
		       system_prompt, search_keywords
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
		err := rows.Scan(
			&char.ID, &char.Name, &char.Description, &char.AvatarURL,
			&char.PersonalitySignature, &char.PersonalityTraits,
			&char.BackgroundStory, &char.VoiceSettings,
			&char.SystemPrompt, &char.SearchKeywords,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan character: %w", err)
		}
		characters = append(characters, char)
	}

	return characters, nil
}
