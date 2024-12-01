// internal/utils/file_helpers.go

package utils

import (
	"path/filepath"
	"strings"
)

func GetFileType(filePath string) string {
	extension := filepath.Ext(filePath)
	switch extension {
	case ".py":
		return "python"
	case ".json", ".yaml", ".toml":
		return "config"
	default:
		return "other"
	}
}

func GetFileLayer(filePath string) string {
	if strings.Contains(filePath, "application") {
		return "application"
	} else if strings.Contains(filePath, "adapters") {
		return "adapters"
	} else {
		return "other"
	}
}

func SummarizeContent(content string) string {
	if strings.Contains(content, "datetime") || strings.Contains(content, "time") {
		return "Handles date and time operations"
	}
	if strings.Contains(content, "async") {
		return "Uses asynchronous code"
	}
	return "General code"
}

func GetSourceDirName(projectTree string) string {
	// Implement logic to extract source directory name from project tree
	// For example, parse the project tree and find the directory containing source code
	return "demo_project_backend"
}
