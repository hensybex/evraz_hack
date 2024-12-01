// internal/model/project_analysis_result.go

package model

import (
	"time"

	"gorm.io/gorm"
)

type ProjectAnalysisResult struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`

	ProjectID       uint   `gorm:"not null;index" json:"projectId"`
	PromptName      string `json:"promptName"` // New field to identify the prompt
	Compliance      string `gorm:"type:text" json:"compliance"`
	Issues          string `gorm:"type:text" json:"issues"`
	Recommendations string `gorm:"type:text" json:"recommendations"`

	Project Project `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE"`
}
