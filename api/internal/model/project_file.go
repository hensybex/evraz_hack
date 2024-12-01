package model

import (
	"time"

	"gorm.io/gorm"
)

type ProjectFile struct {
	ID          uint           `gorm:"primaryKey" json:"id,omitempty"`
	CreatedAt   time.Time      `json:"createdAt,omitempty"`
	UpdatedAt   time.Time      `json:"updatedAt,omitempty"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
	ProjectID   uint           `gorm:"not null;index" json:"project_id"`
	Path        string         `json:"path"`
	Name        string         `json:"name"`
	Content     string         `json:"content"`
	WasAnalyzed bool           `json:"was_analyzed"`
	GPTCallID   *uint          `json:"gpt_call_id,omitempty"`

	Project             Project              `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	FileAnalysisResults []FileAnalysisResult `gorm:"foreignKey:ProjectFileID"`
}
