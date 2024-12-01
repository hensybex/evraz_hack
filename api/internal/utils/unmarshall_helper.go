package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// ExtractJSON attempts to extract and parse JSON data from a string with potential extra text.
// It returns an error if no valid JSON is found or if unmarshaling fails.
func ExtractJSON(input string, v interface{}) error {
	// Locate the first '{' and last '}' in the input string
	start := strings.Index(input, "{")
	end := strings.LastIndex(input, "}")

	// Ensure both '{' and '}' are found and properly positioned
	if start == -1 || end == -1 || start > end {
		return errors.New("no valid JSON object found")
	}

	// Extract the JSON substring
	jsonStr := input[start : end+1]

	// Try to unmarshal the cleaned JSON string
	if err := json.Unmarshal([]byte(jsonStr), v); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return nil
}
