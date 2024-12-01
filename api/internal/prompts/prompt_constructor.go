// internal/prompts/prompt_constructor.go

package prompts

import (
	"evraz_api/internal/prompts/types"
	"fmt"
)

// PromptConstructor is used to construct prompts with specific content
type PromptConstructor struct{}

// NewPromptConstructor initializes and returns a PromptConstructor
func NewPromptConstructor() *PromptConstructor {
	return &PromptConstructor{}
}

// PromptData is an interface that defines a method to transform data into PassedData
type PromptData interface {
	ToPassedData() []types.PassedData
}

// GetPrompt constructs a prompt based on the provided prompt project and data
func (pc *PromptConstructor) GetPrompt(prompt types.Prompt, data types.PromptData, language string, repeatLanguage bool) (string, error) {
	prompt.Language = language
	passedData := data.ToPassedData()

	// Construct the list of passed data
	passedDataList := "You will receive:\n"
	passedDataContentStr := ""
	for _, data := range passedData {
		passedDataList += fmt.Sprintf("%s - %s\n", data.Name, data.Description)
		passedDataContentStr += fmt.Sprintf("%s - %s\n\n", data.Name, data.Content)
	}

	// Construct the language instruction if applicable
	languageInstruction := ""
	if prompt.Language != "" {
		languageInstruction = fmt.Sprintf("Fully respond in %s.\n", prompt.Language)
	}

	// Construct the JSON instruction if applicable
	jsonInstruction := ""
	if len(prompt.JSONStruct) > 0 {
		jsonInstruction = "Your response should be a structured JSON with the following keys:\n"
		for _, js := range prompt.JSONStruct {
			jsonInstruction += fmt.Sprintf("%s: %s\n", js.Key, js.Description)
		}
	}

	// Construct the final prompt with optional language instruction
	finalPrompt := ""
	if languageInstruction != "" {
		finalPrompt = languageInstruction // Always include at the start
	}

	finalPrompt += prompt.BasePrompt + "\n\n" + passedDataList + prompt.BaseTaskDesc + "\n\n" + passedDataContentStr

	if languageInstruction != "" && repeatLanguage {
		finalPrompt += languageInstruction // Add again at the end if repeatLanguage is true
	}

	// Append JSON structure instruction if available
	finalPrompt += jsonInstruction

	return finalPrompt, nil
}
