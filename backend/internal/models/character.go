package models

import "time"

type PresetCharacter struct {
	ID                   int       `json:"id" db:"id"`
	Name                 string    `json:"name" db:"name"`
	Description          string    `json:"description" db:"description"`
	AvatarURL            string    `json:"avatar_url" db:"avatar_url"`
	PersonalitySignature string    `json:"personality_signature" db:"personality_signature"`
	PersonalityTraits    string    `json:"personality_traits" db:"personality_traits"`
	BackgroundStory      string    `json:"background_story" db:"background_story"`
	VoiceSettings        string    `json:"voice_settings" db:"voice_settings"`
	SystemPrompt         string    `json:"system_prompt" db:"system_prompt"`
	SearchKeywords       string    `json:"search_keywords" db:"search_keywords"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
}

type CharacterResponse struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	Description          string `json:"description"`
	AvatarURL            string `json:"avatar_url"`
	PersonalitySignature string `json:"personality_signature"`
	PersonalityTraits    string `json:"personality_traits"`
	BackgroundStory      string `json:"background_story"`
	VoiceSettings        string `json:"voice_settings"`
	SystemPrompt         string `json:"system_prompt"`
	SearchKeywords       string `json:"search_keywords"`
	Skills               string `json:"skills"`
}
