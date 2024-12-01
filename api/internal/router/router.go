// internal/router/router.go

package router

import (
	"evraz_api/internal/di"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(container *di.DIContainer) *gin.Engine {
	r := gin.Default()

	config := cors.Config{
		AllowAllOrigins:  true, // Allow requests from any origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // Enable cookies or HTTP auth
		MaxAge:           12 * time.Hour,
	}

	r.Use(cors.New(config))

	// Auth routes
	apiGroup := r.Group("/api")
	projectsGroup := apiGroup.Group("/projects")
	{
		projectsGroup.POST("/upload_project", container.ProjectHandlers.UploadProject)
		projectsGroup.POST("/analyze", container.ProjectHandlers.AnalyzeProject)
		projectsGroup.POST("/files/:file_id/analyze", container.ProjectHandlers.AnalyzeFile)

		projectsGroup.GET("/all", container.ProjectHandlers.GetAllProjects)
		projectsGroup.GET("/:project_id/overview", container.ProjectHandlers.GetProjectOverview)
		projectsGroup.GET("/:project_id/generate_pdf", container.ProjectHandlers.GenerateProjectPDF)
	}
	filesGroup := apiGroup.Group("/files")
	{
		filesGroup.GET("/:file_id/analysis_results", container.ProjectHandlers.GetFileAnalysisResults)
	}

	return r
}
