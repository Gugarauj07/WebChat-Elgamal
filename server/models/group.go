package models

import (
	"time"
)

type Group struct {
	ID         string    `gorm:"primaryKey"`
	Name       string    `gorm:"not null"`
	SenderKey  []byte    `gorm:"not null"`
	AdminID    string    `gorm:"not null;index"`
	CreatedAt  time.Time

	// Relacionamentos
	Admin        User           `gorm:"foreignKey:AdminID"`
	GroupMembers []GroupMember  `gorm:"foreignKey:GroupID"`
	Conversation []Conversation `gorm:"foreignKey:GroupID"`
} 