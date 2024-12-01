package model

import (
	"time"

	"gorm.io/gorm"
)

type FileAnalysisResult struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
	PromptName      string         `json:"promptName"`
	Compliance      string         `gorm:"type:text" json:"compliance"`
	Issues          string         `gorm:"type:text" json:"issues"`
	Recommendations string         `gorm:"type:text" json:"recommendations"`

	ProjectFileID uint        `gorm:"not null;index" json:"projectFileId"`
	ProjectFile   ProjectFile `gorm:"foreignKey:ProjectFileID;constraint:OnDelete:CASCADE"`
}
