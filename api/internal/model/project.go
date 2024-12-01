// internal/model/project.go

package model

import (
	"time"

	"gorm.io/gorm"
)

type Project struct {
	ID                    uint           `gorm:"primaryKey" json:"id,omitempty"`
	CreatedAt             time.Time      `json:"createdAt,omitempty"`
	UpdatedAt             time.Time      `json:"updatedAt,omitempty"`
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
	ProgrammingLanguageID uint           `gorm:"not null;index" json:"programming_language_id"`
	GPTCallID             *uint          `json:"gpt_call_id,omitempty"`
	Name                  string         `json:"name"`
	Path                  string         `json:"path"`
	Tree                  string         `json:"tree"`
	WasAnalyzed           bool           `json:"was_analyzed"`

	ProgrammingLanguage ProgrammingLanguage `gorm:"foreignKey:ProgrammingLanguageID"`
}
