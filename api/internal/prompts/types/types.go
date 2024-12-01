// internal/prompts/types/types.go

package types

// PassedData represents the structure of data passed to the prompt
type PassedData struct {
	Name        string
	Description string
	Content     string
}

// JSONStruct represents the structure of the JSON response expected
type JSONStruct struct {
	Key         string
	Description string
	Example     string
}

// Prompt represents a single prompt with a base prompt and task description
type Prompt struct {
	BasePrompt   string
	BaseTaskDesc string
	PassedData   []PassedData
	JSONStruct   []JSONStruct
	Language     string
}

// PromptData is an interface that all prompt data types implement
type PromptData interface {
	ToPassedData() []PassedData
}
