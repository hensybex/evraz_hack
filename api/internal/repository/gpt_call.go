// internal/repository/gpt_call.go

package repository

import (
	"evraz_api/internal/model"

	"gorm.io/gorm"
)

type GPTCallRepository interface {
	CreateOne(gptCall *model.GPTCall) (uint, error)
}

type GormGPTCallRepository struct {
	db *gorm.DB
}

func NewGormGPTCallRepository(db *gorm.DB) *GormGPTCallRepository {
	return &GormGPTCallRepository{db: db}
}

func (repo *GormGPTCallRepository) CreateOne(gptCall *model.GPTCall) (uint, error) {
	if err := repo.db.Create(gptCall).Error; err != nil {
		return 0, err
	}
	return gptCall.ID, nil
}
