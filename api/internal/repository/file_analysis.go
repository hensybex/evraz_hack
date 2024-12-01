// internal/repository/file_analysis.go

package repository

import (
	"evraz_api/internal/model"
	"log"

	"gorm.io/gorm"
)

type FileAnalysisRepository interface {
	CreateOne(analysis *model.FileAnalysisResult) error
	GetManyByFileID(projectFileID uint) ([]model.FileAnalysisResult, error)
	// Add more methods as needed
}

type GormFileAnalysisRepository struct {
	db *gorm.DB
}

func NewGormFileAnalysisRepository(db *gorm.DB) *GormFileAnalysisRepository {
	return &GormFileAnalysisRepository{db: db}
}

func (repo *GormFileAnalysisRepository) CreateOne(analysis *model.FileAnalysisResult) error {
	log.Printf("Attempting to create FileAnalysis: %+v", analysis)
	err := repo.db.Create(analysis).Error
	if err != nil {
		log.Printf("Error during FileAnalysis creation: %v", err)
	} else {
		log.Println("FileAnalysis creation succeeded.")
	}
	return err
}

func (repo *GormFileAnalysisRepository) GetManyByFileID(projectFileID uint) ([]model.FileAnalysisResult, error) {
	var results []model.FileAnalysisResult
	if err := repo.db.Where("project_file_id = ?", projectFileID).Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}
