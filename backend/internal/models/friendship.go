package models

import "time"

type UserFriendship struct {
	ID            int        `json:"id" db:"id"`
	UserID        int        `json:"user_id" db:"user_id"`
	CharacterID   int        `json:"character_id" db:"character_id"`
	IsActive      bool       `json:"is_active" db:"is_active"`
	LastMessageAt *time.Time `json:"last_message_at" db:"last_message_at"`
	UnreadCount   int        `json:"unread_count" db:"unread_count"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

type FriendInfo struct {
	ID                   int        `json:"id"`
	Name                 string     `json:"name"`
	AvatarURL            string     `json:"avatar_url"`
	PersonalitySignature string     `json:"personality_signature"`
	LastMessage          string     `json:"last_message"`
	LastMessageAt        *time.Time `json:"last_message_at"`
	UnreadCount          int        `json:"unread_count"`
	IsOnline             bool       `json:"is_online"`
}

type AddFriendRequest struct {
	CharacterID int `json:"character_id" binding:"required"`
}

type SearchFriendsRequest struct {
	Keyword string `json:"keyword" binding:"required"`
}

type AvailableCharacter struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	Description          string `json:"description"`
	AvatarURL            string `json:"avatar_url"`
	PersonalitySignature string `json:"personality_signature"`
	IsAdded              bool   `json:"is_added"`
}
