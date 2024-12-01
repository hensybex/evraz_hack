// internal/usecase/project_analysis.go

package usecase

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"evraz_api/internal/dto/llm_responses"
	"evraz_api/internal/model"
	"evraz_api/internal/prompts"
	"evraz_api/internal/prompts/prompts_storage/file_prompts"
	helper_prompts "evraz_api/internal/prompts/prompts_storage/helpers"
	"evraz_api/internal/prompts/prompts_storage/project_prompts"
	"evraz_api/internal/prompts/types"
	"evraz_api/internal/repository"
	"evraz_api/internal/service"
	"evraz_api/internal/utils"
)

type ProjectAnalysisUsecase struct {
	ProjectRepo         repository.ProjectRepository
	ProjectFileRepo     repository.ProjectFileRepository
	ProjectAnalysisRepo repository.ProjectAnalysisRepository
	FileAnalysisRepo    repository.FileAnalysisRepository
	MistralService      service.MistralService
	Prompts             *prompts.Prompts
	PromptConstructor   *prompts.PromptConstructor
}

func NewProjectAnalysisUsecase(
	projectRepo repository.ProjectRepository,
	projectFileRepo repository.ProjectFileRepository,
	projectAnalysisRepo repository.ProjectAnalysisRepository,
	fileAnalysisRepo repository.FileAnalysisRepository,
	mistralService service.MistralService,
) *ProjectAnalysisUsecase {
	return &ProjectAnalysisUsecase{
		ProjectRepo:         projectRepo,
		ProjectFileRepo:     projectFileRepo,
		ProjectAnalysisRepo: projectAnalysisRepo,
		FileAnalysisRepo:    fileAnalysisRepo,
		MistralService:      mistralService,
		Prompts:             prompts.NewPrompts(),
		PromptConstructor:   prompts.NewPromptConstructor(),
	}
}

func (uc *ProjectAnalysisUsecase) AnalyzeProject(projectID uint) error {

	project, err := uc.ProjectRepo.GetOneByID(projectID)
	project.WasAnalyzed = true
	if err := uc.ProjectRepo.UpdateOneByID(project); err != nil {
		return fmt.Errorf("failed to update project GPTCallID: %w", err)
	}
	if err != nil {
		return fmt.Errorf("failed to retrieve project: %w", err)
	}

	// Map prompt names to actual prompts
	promptMap := map[string]types.Prompt{
		"ProjectStructure":        uc.Prompts.ProjectStructure,
		"KeyFiles":                uc.Prompts.KeyFiles,
		"ApplicationArchitecture": uc.Prompts.ApplicationArchitecture,
		"DependencyManagement":    uc.Prompts.DependencyManagement,
		"ProjectSettings":         uc.Prompts.ProjectSettings,
		"TestingStrategy":         uc.Prompts.TestingStrategy,
		"AdditionalTechnical":     uc.Prompts.AdditionalTechnical,
		"DateTimeHandling":        uc.Prompts.DateTimeHandling,
	}

	// Iterate over all prompts in promptMap
	for promptName, promptTemplate := range promptMap {
		var data types.PromptData
		var outputDBValue string
		var emptyValue bool

		switch promptName {
		case "ProjectStructure":
			data = project_prompts.ProjectStructureData{
				ProjectTree: project.Tree,
			}
		case "KeyFiles":
			emptyValue = false
			setupFilesContent := make(map[string]string)
			setupFiles := []string{"setup.py", "setup.cfg", "pyproject.toml", "README.md"}

			foundAny := false
			for _, filename := range setupFiles {
				content, err := uc.ProjectFileRepo.GetRootFileContentByName(project.ID, filename)
				if err == nil {
					foundAny = true
					setupFilesContent[filename] = content
				} else {
					content, err = uc.ProjectFileRepo.GetFileContentByName(project.ID, filename)
					if err == nil {
						foundAny = true
						setupFilesContent[filename] = content
					} else {
						setupFilesContent[filename] = filename + " is missing at the root of the project directory"
					}
				}
			}
			if foundAny == false {
				outputDBValue = "setup.py, setup.cfg, pyproject.toml и README.md отсутствуют в корне проекта"
				emptyValue = true
			}
			data = project_prompts.KeyFilesData{
				SetupFilesContent: setupFilesContent,
			}
		case "ApplicationArchitecture":
			data = project_prompts.ApplicationArchitectureData{
				ProjectStructure: project.Tree,
			}
		case "DependencyManagement":
			emptyValue = false
			dependencyFile := "requirements.txt"
			dependenciesContent, err := uc.ProjectFileRepo.GetRootFileContentByName(project.ID, dependencyFile)
			if err != nil {
				dependenciesContent, err = uc.ProjectFileRepo.GetFileContentByName(project.ID, dependencyFile)
				if err != nil {
					dependenciesContent = dependencyFile + " is missing at the root of the project directory"
					outputDBValue = dependenciesContent
					emptyValue = true
				}
			}
			data = project_prompts.DependencyManagementData{
				DependenciesContent: dependenciesContent,
			}
		case "ProjectSettings":
			emptyValue = false
			settingsFiles := []string{"settings.py", "config.yaml", "pyproject.toml"}
			settingsContent := make(map[string]string)

			foundAny := false
			for _, filename := range settingsFiles {
				content, err := uc.ProjectFileRepo.GetRootFileContentByName(project.ID, filename)
				if err == nil {
					foundAny = true
					settingsContent[filename] = content
				} else {
					content, err = uc.ProjectFileRepo.GetFileContentByName(project.ID, filename)
					if err == nil {
						foundAny = true
						settingsContent[filename] = content
					} else {
						settingsContent[filename] = filename + " is missing at the root of the project directory"
					}
				}
			}
			if foundAny == false {
				outputDBValue = "settings.py, config.yaml и pyproject.toml отсутствуют в корне проекта"
				emptyValue = true
			}
			data = project_prompts.ProjectSettingsData{
				SettingsFilesContent: settingsContent,
			}
		case "TestingStrategy":
			emptyValue = false
			// Construct the prompt to identify tests
			data = helper_prompts.ExtractTestsData{
				ProjectTree: project.Tree,
			}
			prompt, err := uc.PromptConstructor.GetPrompt(helper_prompts.ExtractTestsPrompt, data, "Russian (русский)", true)
			if err != nil {
				return fmt.Errorf("failed to construct prompt for %s: %w", promptName, err)
			}

			// Call the LLM
			analysisResult, _, err := uc.MistralService.CallMistral(prompt, true, service.Hack, "testsExtract", project.ID)
			if err != nil {
				return fmt.Errorf("failed to call Mistral service for %s: %w", promptName, err)
			}

			// Parse the response
			var testsExtractResponse struct {
				TestFilesRoutes []string `json:"test_files_routes"`
			}
			if err := utils.ExtractJSON(analysisResult, &testsExtractResponse); err != nil {
				emptyValue = true
				continue
			}

			// Initialize a variable to store concatenated paths and contents
			var concattedString string

			// Fetch file details and concatenate paths and contents
			for _, route := range testsExtractResponse.TestFilesRoutes {
				// Convert string ID to uint
				testsFileContent, err := uc.ProjectFileRepo.GetFileContentByPath(projectID, route)
				if err != nil {
					return fmt.Errorf("failed to fetch file at route %s: %w", route, err)
				}

				// Ensure both path and content are not empty before appending
				if route != "" && testsFileContent != "" {
					concattedString += fmt.Sprintf("Path: %s\nContent:\n%s\n\n", route, testsFileContent)
				}
			}

			// Pass concatenated string to the next prompt
			data = project_prompts.TestingStrategyData{
				ProjectTree:       project.Tree,
				TestsFilesContent: concattedString,
			}

		case "AdditionalTechnical":
			emptyValue = false
			filePath := "transaction_manager.py"
			foundAny := false
			transactionCode, err := uc.ProjectFileRepo.GetRootFileContentByName(project.ID, filePath)
			if err != nil {
				transactionCode, err = uc.ProjectFileRepo.GetFileContentByName(project.ID, filePath)
				if err != nil {
					transactionCode = filePath + " не найден"
					outputDBValue = transactionCode
					emptyValue = true
				} else {
					foundAny = true
				}
			} else {
				foundAny = true
			}
			filePath = "async_features.py"
			asyncCode, err := uc.ProjectFileRepo.GetRootFileContentByName(project.ID, filePath)
			if err != nil {
				asyncCode, err = uc.ProjectFileRepo.GetFileContentByName(project.ID, filePath)
				if err != nil {
					asyncCode = filePath + " не найден"
				} else {
					foundAny = true
				}
			} else {
				foundAny = true
			}
			if foundAny == false {
				outputDBValue = "transaction_manager.py и async_features.py отсутствуют в корне проекта"
				emptyValue = true
			}
			data = project_prompts.AdditionalTechnicalData{
				TransactionManagementCode: transactionCode,
				AsynchronousCodeUsage:     asyncCode,
			}
		case "DateTimeHandling":
			emptyValue = false
			filePath := "datetime_utils.py"
			dateTimeCode, err := uc.ProjectFileRepo.GetRootFileContentByName(project.ID, filePath)
			if err != nil {
				dateTimeCode, err = uc.ProjectFileRepo.GetFileContentByName(project.ID, filePath)
				if err != nil {
					dateTimeCode = filePath + " не найден"
					outputDBValue = dateTimeCode
					emptyValue = true
				}
			}
			data = project_prompts.DateTimeHandlingData{
				DateTimeCodeSamples: dateTimeCode,
			}
		default:
			continue
		}

		var analysisResult string
		var gptCallID uint
		var analysisDTO llm_responses.FileAnalysisResponse
		if emptyValue == false {

			// Construct the prompt
			prompt, err := uc.PromptConstructor.GetPrompt(promptTemplate, data, "Russian (русский)", true)
			if err != nil {
				return fmt.Errorf("failed to construct prompt for %s: %w", promptName, err)
			}

			// Call the LLM
			analysisResult, gptCallID, err = uc.MistralService.CallMistral(prompt, false, service.Hack, "project", project.ID)
			if err != nil {
				return fmt.Errorf("failed to call Mistral service for %s: %w", promptName, err)
			}
			if err := utils.ExtractJSON(analysisResult, &analysisDTO); err != nil {
				log.Println("Error while extracting Project Analysis Response")
				continue
			}
		} else {
			gptCallID = 0
			analysisResult = outputDBValue
			analysisDTO.Compliance = false
			analysisDTO.Issues = append(analysisDTO.Issues, outputDBValue)
			analysisDTO.Recommendations = append(analysisDTO.Recommendations, outputDBValue)
		}

		// Save the analysis result
		projectAnalysis := &model.ProjectAnalysisResult{
			ProjectID:       project.ID,
			PromptName:      promptName,
			Compliance:      fmt.Sprintf("%t", analysisDTO.Compliance),
			Issues:          strings.Join(analysisDTO.Issues, ", "),
			Recommendations: strings.Join(analysisDTO.Recommendations, ", "),
		}

		if err := uc.ProjectAnalysisRepo.CreateOne(projectAnalysis); err != nil {
			return fmt.Errorf("failed to save project analysis for %s: %w", promptName, err)
		}
		log.Println("Successfully created ProjectAnalysis in the database.")

		// Optionally, update the project with GPTCallID
		project.GPTCallID = &gptCallID
		project.WasAnalyzed = true
		if err := uc.ProjectRepo.UpdateOneByID(project); err != nil {
			return fmt.Errorf("failed to update project GPTCallID: %w", err)
		}
	}

	// Analyze each file
	files, err := uc.ProjectFileRepo.GetManyByProjectID(project.ID)
	if err != nil {
		return fmt.Errorf("failed to retrieve project files: %w", err)
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(files))
	semaphore := make(chan struct{}, 5)
	for _, file := range files {
		wg.Add(1)
		semaphore <- struct{}{} // Acquire a slot

		go func(fileID uint) {
			defer wg.Done()
			defer func() { <-semaphore }() // Release the slot

			if err := uc.AnalyzeFile(fileID); err != nil {
				errChan <- fmt.Errorf("failed to analyze file %d: %w", fileID, err)
			}
		}(file.ID)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	close(errChan)

	// Check for errors
	if len(errChan) > 0 {
		return <-errChan // Return the first error
	}

	return nil
}

func (uc *ProjectAnalysisUsecase) AnalyzeFile(fileID uint) error {

	file, err := uc.ProjectFileRepo.GetOneByID(fileID)
	if err != nil {
		return fmt.Errorf("failed to retrieve project file: %w", err)
	}

	project, err := uc.ProjectRepo.GetOneByID(file.ProjectID)
	if err != nil {
		return fmt.Errorf("failed to retrieve project for file: %w", err)
	}

	var targetExtension string
	if project.ProgrammingLanguageID == 1 {
		targetExtension = ".py"
	}
	if !strings.HasSuffix(file.Name, targetExtension) {
		return nil
	}

	// Construct the master prompt data
	masterData := file_prompts.FileMasterData{
		ProjectTree: project.Tree,
		FilePath:    file.Path,
		FileContent: file.Content,
	}

	// Construct the master prompt
	masterPrompt, err := uc.PromptConstructor.GetPrompt(uc.Prompts.FileMasterPrompt, masterData, "Russian (русский)", true)
	if err != nil {
		return fmt.Errorf("failed to construct master prompt: %w", err)
	}

	// Call the LLM to get applicable prompts
	masterResult, _, err := uc.MistralService.CallMistral(masterPrompt, false, service.Hack, "file", file.ID)
	if err != nil {
		return fmt.Errorf("failed to call Mistral service for master prompt: %w", err)
	}

	// Parse the LLM response to get a single int that indicates the file type
	var masterResponse struct {
		Value int `json:"value"`
	}
	if err := utils.ExtractJSON(masterResult, &masterResponse); err != nil {
		masterResponse.Value = 2
	}

	promptMap := map[string]types.Prompt{
		"CodingStandards":         uc.Prompts.CodingStandards,
		"ErrorHandlingAndLogging": uc.Prompts.ErrorHandlingAndLogging,
		"AdditionalTechnicalFile": uc.Prompts.AdditionalTechnicalFile,
		"DateTimeHandlingFile":    uc.Prompts.DateTimeHandlingFile,
	}

	// Prepare an array to hold the applicable prompts based on master response
	var applicablePrompts []string

	if masterResponse.Value == 0 {
		applicablePrompts = append(applicablePrompts, "ApplicationLayerCode")
	} else if masterResponse.Value == 1 {
		applicablePrompts = append(applicablePrompts, "AdaptersLayerCode")
	}

	// Always add universal prompts
	for k := range promptMap {
		applicablePrompts = append(applicablePrompts, k)
	}

	// Iterate over applicable prompts and perform analysis
	for _, promptName := range applicablePrompts {

		var promptTemplate types.Prompt
		switch promptName {
		case "ApplicationLayerCode":
			promptTemplate = uc.Prompts.ApplicationLayerCode
		case "AdaptersLayerCode":
			promptTemplate = uc.Prompts.AdaptersLayerCode
		case "CodingStandards":
			promptTemplate = uc.Prompts.CodingStandards
		case "ErrorHandlingAndLogging":
			promptTemplate = uc.Prompts.ErrorHandlingAndLogging
		case "AdditionalTechnicalFile":
			promptTemplate = uc.Prompts.AdditionalTechnicalFile
		case "DateTimeHandlingFile":
			promptTemplate = uc.Prompts.DateTimeHandlingFile
		default:
			continue
		}

		// Prepare prompt data based on prompt
		var data types.PromptData
		switch promptName {
		case "ApplicationLayerCode":
			data = file_prompts.ApplicationLayerCodeData{
				FilePath:    file.Path,
				FileContent: file.Content,
			}
		case "AdaptersLayerCode":
			data = file_prompts.AdaptersLayerCodeData{
				FilePath:    file.Path,
				FileContent: file.Content,
			}
		case "CodingStandards":
			data = file_prompts.CodingStandardsData{
				FilePath:    file.Path,
				FileContent: file.Content,
			}
		case "ErrorHandlingAndLogging":
			data = file_prompts.ErrorHandlingAndLoggingData{
				FilePath:    file.Path,
				FileContent: file.Content,
			}
		case "AdditionalTechnicalFile":
			data = file_prompts.AdditionalTechnicalFileData{
				FilePath:    file.Path,
				FileContent: file.Content,
			}
		case "DateTimeHandlingFile":
			data = file_prompts.DateTimeHandlingFileData{
				FilePath:    file.Path,
				FileContent: file.Content,
			}
		default:
			continue
		}

		// Construct the prompt
		prompt, err := uc.PromptConstructor.GetPrompt(promptTemplate, data, "Russian (русский)", true)
		if err != nil {
			return fmt.Errorf("failed to construct prompt: %w", err)
		}

		// Call the LLM
		analysisResult, gptCallID, err := uc.MistralService.CallMistral(prompt, false, service.Hack, "file", file.ID)
		if err != nil {
			return fmt.Errorf("failed to call Mistral service: %w", err)
		}

		// Parse and process the analysis result
		var analysisDTO llm_responses.FileAnalysisResponse
		if err := utils.ExtractJSON(analysisResult, &analysisDTO); err != nil {
			log.Println("Error while extracting Project Analysis Response")
			continue
		}

		// Save the analysis result
		fileAnalysis := &model.FileAnalysisResult{
			ProjectFileID:   file.ID,
			PromptName:      promptName,
			Compliance:      fmt.Sprintf("%t", analysisDTO.Compliance),
			Issues:          strings.Join(analysisDTO.Issues, ", "),
			Recommendations: strings.Join(analysisDTO.Recommendations, ", "),
		}
		if err := uc.FileAnalysisRepo.CreateOne(fileAnalysis); err != nil {
			return fmt.Errorf("failed to save file analysis: %w", err)
		}
		log.Println("Successfully created FileAnalysis in the database.")

		// Optionally, update the file with GPTCallID
		file.GPTCallID = &gptCallID
		file.WasAnalyzed = true
		if err := uc.ProjectFileRepo.UpdateOneByID(file); err != nil {
			return fmt.Errorf("failed to update file GPTCallID: %w", err)
		}
	}

	return nil
}
