// internal/model/programming_languages.go

package model

import (
	"time"

	"gorm.io/gorm"
)

type ProgrammingLanguage struct {
	ID        uint           `gorm:"primaryKey" json:"id,omitempty"`
	CreatedAt time.Time      `json:"createdAt,omitempty"`
	UpdatedAt time.Time      `json:"updatedAt,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
	Name      string         `gorm:"not null" json:"name"`
}
