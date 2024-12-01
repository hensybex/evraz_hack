// internal/migration/migration.go

package migration

import (
	"evraz_api/internal/model"
	"log"

	"gorm.io/gorm"
)

func ApplyCustomMigrations(db *gorm.DB) error {
	// Automigrate programming languages
	programmingLanguages := []model.ProgrammingLanguage{
		{Name: "Python"},
		{Name: "C#"},
		{Name: "Typescript"},
	}

	for _, language := range programmingLanguages {
		if err := db.FirstOrCreate(&language, model.ProgrammingLanguage{Name: language.Name}).Error; err != nil {
			return err
		}
	}

	log.Println("Custom migrations applied successfully.")
	return nil
}
