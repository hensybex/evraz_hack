// internal/model/gpt_call.go

package model

import (
	"gorm.io/gorm"
	"time"
)

// GPTCall represents a record of a call to the GPT-based service.
type GPTCall struct {
	ID               uint           `json:"id,omitempty"`
	CreatedAt        time.Time      `json:"createdAt,omitempty"`
	UpdatedAt        time.Time      `json:"updatedAt,omitempty"`
	DeletedAt        gorm.DeletedAt `json:"deletedAt,omitempty" gorm:"index"`
	FinalPrompt      string         `json:"finalPrompt,omitempty"`
	Reply            string         `json:"reply,omitempty"`
	EntityType       string         `json:"entity_type,omitempty"`
	EntityID         uint           `json:"entity_id,omitempty"`
	PromptTokens     int            `json:"promptTokens,omitempty"`
	CompletionTokens int            `json:"completionTokens,omitempty"`
	TotalTokens      int            `json:"totalTokens,omitempty"`
}
