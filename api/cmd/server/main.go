// cmd/server/main.go

package main

import (
	"evraz_api/internal/config"
	"evraz_api/internal/di"
	"evraz_api/internal/migration"
	"evraz_api/internal/model"
	"evraz_api/internal/router"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// AutoMigrate models
	if err := db.AutoMigrate(
		&model.GPTCall{},
		&model.ProgrammingLanguage{},
		&model.Project{},
		&model.ProjectFile{},
		&model.ProjectAnalysisResult{},
		&model.FileAnalysisResult{},
	); err != nil {
		log.Fatalf("Failed to automigrate: %v", err)
	}

	err = migration.ApplyCustomMigrations(db)
	if err != nil {
		log.Fatalf("Failed to apply custom migration : %v", err)
	}

	// Initialize DI container
	container := di.NewDIContainer(cfg, db)

	// Setup router
	r := router.SetupRouter(container)

	// Start server
	if err := r.Run("0.0.0.0:8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
