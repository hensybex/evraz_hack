// internal/service/mistral_service.go

package service

import (
	"bytes"
	"encoding/json"
	"evraz_api/internal/model"
	"evraz_api/internal/repository"
	"fmt"
	"strings"

	//"go_backend/internal/usecase"
	"io"
	"net/http"
	"os"
	"time"
)

// Define a new type for MistralModel with constants
type MistralModel string

const (
	Nemo      MistralModel = "open-mistral-nemo"
	Largest   MistralModel = "mistral-large-latest"
	Codestral MistralModel = "codestral-latest"
	Hack      MistralModel = "mistral-nemo-instruct-2407"
)

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens"`
	Temperature float64       `json:"temperature"`
}

type ChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int     `json:"prompt_tokens"`
		TotalTokens      int     `json:"total_tokens"`
		TokensPerSecond  float64 `json:"tokens_per_second"`
		CompletionTokens int     `json:"completion_tokens"`
	} `json:"usage"`
	// You can also include the request_id, response_id, model, etc. if needed
}

type MistralService struct {
	apiKey      string
	url         string
	model       string
	GPTCallRepo repository.GPTCallRepository
}

func NewMistralService(gptCallRepo repository.GPTCallRepository) *MistralService {
	apiKey := os.Getenv("MISTRAL_API_KEY")
	apiURL := os.Getenv("MISTRAL_API_URL")
	apiModel := os.Getenv("MISTRAL_API_MODEL")

	if apiKey == "" {
		fmt.Println("MISTRAL_API_KEY environment variable is not set")
	}
	if apiURL == "" {
		fmt.Println("MISTRAL_API_URL environment variable is not set, using default URL")
		apiURL = "https://api.mistral.ai/v1/chat/completions"
	}
	if apiModel == "" {
		fmt.Println("MISTRAL_API_MODEL environment variable is not set, using default model")
		apiModel = string(Nemo)
	}

	return &MistralService{
		apiKey:      apiKey,
		url:         apiURL,
		model:       apiModel,
		GPTCallRepo: gptCallRepo,
	}
}

/* func (ms *MistralService) CallMistral(prompt string, needJson bool, mistralModel MistralModel, entityType string, entityID uint) (string, uint, error) {
	fmt.Printf("Starting CallMistral")
	//fmt.Printf("Starting CallMistral with prompt: %s\n", prompt)

	messages := []ChatMessage{
		{
			Role:    "user",
			Content: prompt,
		},
	}

	requestBody := ChatRequest{
		Model:    string(mistralModel),
		Messages: messages,
	}

	if needJson {
		requestBody.ResponseFormat = map[string]string{"type": "json_object"}
	}

	jsonValue, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Printf("Failed to marshal request body: %v\n", err)
		return "", 0, err
	}

	//fmt.Printf("Request body JSON: %s\n", string(jsonValue))

	var resp *http.Response
	var body []byte
	maxAttempts := 5
	for attempts := 1; attempts <= maxAttempts; attempts++ {
		req, err := http.NewRequest("POST", ms.url, bytes.NewBuffer(jsonValue))
		if err != nil {
			fmt.Printf("Failed to create HTTP request: %v\n", err)
			return "", 0, err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+ms.apiKey)

		client := &http.Client{}
		resp, err = client.Do(req)
		if err != nil {
			fmt.Printf("HTTP request failed: %v\n", err)
			return "", 0, err
		}

		defer resp.Body.Close()
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Failed to read response body: %v\n", err)
			return "", 0, err
		}

		// Check if the response status is 429
		if resp.StatusCode == http.StatusTooManyRequests {
			fmt.Printf("Rate limit exceeded, attempt %d of %d\n", attempts, maxAttempts)
			if attempts < maxAttempts {
				// Exponential backoff
				sleepTime := time.Duration(attempts*attempts) * time.Second
				fmt.Printf("Sleeping for %v seconds before retrying...\n", sleepTime)
				time.Sleep(sleepTime)
				continue
			} else {
				return "", 0, fmt.Errorf("rate limit exceeded after %d attempts", maxAttempts)
			}
		}

		// Exit loop if status is not 429
		break
	}

	// Handle both JSON and non-JSON cases similarly
	var chatResponse ChatResponse
	err = json.Unmarshal(body, &chatResponse)
	if err != nil {
		fmt.Printf("Failed to unmarshal JSON response: %v\n", err)
		return "", 0, err
	}

	responseContent := ""
	if len(chatResponse.Choices) > 0 {
		responseContent = chatResponse.Choices[0].Message.Content
	} else {
		fmt.Printf("No choices found in JSON response\n")
		return "", 0, fmt.Errorf("no choices in JSON response")
	}

	// Log the GPT call using the provided GPTCallUsecase
	gptCall := model.GPTCall{
		FinalPrompt: prompt,
		Reply:       responseContent,
		EntityType:  entityType,
		EntityID:    entityID,
	}
	gptCallId, err := ms.GPTCallRepo.CreateOne(&gptCall)
	if err != nil {
		fmt.Printf("Failed to log GPT call: %v\n", err)
		return responseContent, 0, fmt.Errorf("failed to log GPT call: %w", err)
	}

	fmt.Printf("Successfully logged GPT call with ID: %d\n", gptCallId)

	return responseContent, gptCallId, nil
} */

/* func (ms *MistralService) CallMistral(prompt string, needJson bool, mistralModel MistralModel, entityType string, entityID uint) (string, uint, error) {
	fmt.Println("Starting CallMistral")

	// Prepare messages
	messages := []ChatMessage{
		{
			Role:    "user",
			Content: prompt,
		},
	}

	// If needJson is true, prepend a system message to request JSON response
	if needJson {
		messages = append([]ChatMessage{
			{
				Role:    "system",
				Content: "Please respond in JSON format.",
			},
		}, messages...)
	}

	// Construct the request body with model, messages, max_tokens, and temperature
	requestBody := ChatRequest{
		Model:       string(mistralModel),
		Messages:    messages,
		MaxTokens:   1000, // Adjust as needed
		Temperature: 0.3,  // Adjust as needed
	}

	jsonValue, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Printf("Failed to marshal request body: %v\n", err)
		return "", 0, err
	}

	// Initialize HTTP client
	client := &http.Client{}

	// Make the HTTP request with retry logic for rate limiting
	var resp *http.Response
	var body []byte
	maxAttempts := 5
	for attempts := 1; attempts <= maxAttempts; attempts++ {
		req, err := http.NewRequest("POST", ms.url, bytes.NewBuffer(jsonValue))
		if err != nil {
			fmt.Printf("Failed to create HTTP request: %v\n", err)
			return "", 0, err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", ms.apiKey) // Use the API key directly as per the API documentation

		resp, err = client.Do(req)
		if err != nil {
			fmt.Printf("HTTP request failed: %v\n", err)
			return "", 0, err
		}

		defer resp.Body.Close()
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Failed to read response body: %v\n", err)
			return "", 0, err
		}

		// Handle rate limiting
		if resp.StatusCode == http.StatusTooManyRequests {
			fmt.Printf("Rate limit exceeded, attempt %d of %d\n", attempts, maxAttempts)
			if attempts < maxAttempts {
				sleepTime := time.Duration(attempts*attempts) * time.Second
				fmt.Printf("Sleeping for %v seconds before retrying...\n", sleepTime)
				time.Sleep(sleepTime)
				continue
			} else {
				return "", 0, fmt.Errorf("rate limit exceeded after %d", attempts)
			}
		}

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Received non-200 response: %d\nResponse body: %s\n", resp.StatusCode, string(body))
			return "", 0, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
		}

		// Break out of loop if successful
		break
	}

	// Parse the response
	var chatResponse ChatResponse
	err = json.Unmarshal(body, &chatResponse)
	if err != nil {
		fmt.Printf("Failed to unmarshal JSON response: %v\nResponse body: %s\n", err, string(body))
		return "", 0, err
	}

	// Extract the response content
	responseContent := ""
	if len(chatResponse.Choices) > 0 {
		responseContent = chatResponse.Choices[0].Message.Content
	} else {
		fmt.Printf("No choices found in JSON response\n")
		return "", 0, fmt.Errorf("no choices in JSON response")
	}

	// Log the GPT call
	gptCall := model.GPTCall{
		FinalPrompt: prompt,
		Reply:       responseContent,
		EntityType:  entityType,
		EntityID:    entityID,
	}
	gptCallId, err := ms.GPTCallRepo.CreateOne(&gptCall)
	if err != nil {
		fmt.Printf("Failed to log GPT call: %v\n", err)
		return responseContent, 0, fmt.Errorf("failed to log GPT call: %w", err)
	}

	fmt.Printf("Successfully logged GPT call with ID: %d\n", gptCallId)

	return responseContent, gptCallId, nil
} */

func (ms *MistralService) CallMistral(prompt string, needJson bool, mistralModel MistralModel, entityType string, entityID uint) (string, uint, error) {
	fmt.Println("Starting CallMistral")

	// Prepare initial messages
	messages := []ChatMessage{
		{
			Role:    "user",
			Content: prompt,
		},
	}

	// If needJson is true, prepend a system message to request JSON response
	if needJson {
		messages = append([]ChatMessage{
			{
				Role:    "system",
				Content: "Please respond in JSON format.",
			},
		}, messages...)
	}

	// Initialize variables for response looping
	var fullResponse string
	var totalTokensUsed int
	var promptTokens int
	var completionTokens int
	maxAttempts := 25

	// Loop to handle token limits and fetch complete response
	for {
		// Construct the request body
		model := ms.model
		requestBody := ChatRequest{
			Model:       model,
			Messages:    messages,
			MaxTokens:   1024, // Max tokens per request
			Temperature: 0.3,  // Adjust as needed
		}

		jsonValue, err := json.Marshal(requestBody)
		if err != nil {
			fmt.Printf("Failed to marshal request body: %v\n", err)
			return "", 0, err
		}

		// Initialize HTTP client
		client := &http.Client{}

		// Make the HTTP request with retry logic
		var resp *http.Response
		var body []byte
		for attempts := 1; attempts <= maxAttempts; attempts++ {
			req, err := http.NewRequest("POST", ms.url, bytes.NewBuffer(jsonValue))
			if err != nil {
				fmt.Printf("Failed to create HTTP request: %v\n", err)
				return "", 0, err
			}

			req.Header.Set("Content-Type", "application/json")
			if ms.url == "https://api.mistral.ai/v1/chat/completions" {
				req.Header.Set("Authorization", "Bearer "+ms.apiKey)
			} else {
				req.Header.Set("Authorization", ms.apiKey)
			}

			resp, err = client.Do(req)
			if err != nil {
				fmt.Printf("HTTP request failed: %v\n", err)
				return "", 0, err
			}

			defer resp.Body.Close()
			body, err = io.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("Failed to read response body: %v\n", err)
				return "", 0, err
			}

			// Handle rate limiting and database connection errors
			if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode != http.StatusOK ||
				(resp.StatusCode == http.StatusBadRequest && strings.Contains(string(body), "too many clients already")) {

				condition := "rate limit exceeded"
				if resp.StatusCode == http.StatusBadRequest {
					condition = "too many clients"
				}
				fmt.Printf("%s, attempt %d of %d\n", condition, attempts, maxAttempts)

				if attempts < maxAttempts {
					sleepTime := time.Duration(attempts*attempts) * time.Second
					fmt.Printf("Sleeping for %v seconds before retrying...\n", sleepTime)
					time.Sleep(sleepTime)
					continue
				} else {
					return "", 0, fmt.Errorf("%s after %d attempts", condition, attempts)
				}
			}

			if resp.StatusCode != http.StatusOK {
				fmt.Printf("Received non-200 response: %d\nResponse body: %s\n", resp.StatusCode, string(body))
				continue
			}

			// Break out of retry loop if successful
			break
		}

		// Parse the response
		var chatResponse ChatResponse
		err = json.Unmarshal(body, &chatResponse)
		if err != nil {
			fmt.Printf("Failed to unmarshal JSON response: %v\nResponse body: %s\n", err, string(body))

			// Check if the error is due to incomplete JSON
			if err == io.ErrUnexpectedEOF || strings.Contains(err.Error(), "unexpected end of JSON input") {
				fmt.Println("Response is incomplete, requesting continuation")

				// Extract whatever content is available
				partialContent := extractPartialContent(string(body))

				// Append partial content to full response
				fullResponse += partialContent

				// Prepare messages to request continuation
				messages = append(messages, ChatMessage{
					Role:    "user",
					Content: partialContent,
				})
				continue
			} else {
				// Other parsing errors
				return "", 0, err
			}
		}

		// Extract the response content
		responseContent := ""
		if len(chatResponse.Choices) > 0 {
			responseContent = chatResponse.Choices[0].Message.Content
		} else {
			fmt.Printf("No choices found in JSON response\n")
			return "", 0, fmt.Errorf("no choices in JSON response")
		}

		// Append the response content to the full response
		fullResponse += responseContent
		totalTokensUsed += chatResponse.Usage.TotalTokens
		promptTokens += chatResponse.Usage.PromptTokens
		completionTokens += chatResponse.Usage.CompletionTokens

		// Check if the assistant indicates that the response is complete
		if isResponseComplete(responseContent) {
			break
		} else {
			// Response may still be incomplete, request continuation
			messages = append(messages, ChatMessage{
				Role:    "user",
				Content: responseContent,
			})
			continue
		}
	}

	// Log the GPT call with the full response
	gptCall := model.GPTCall{
		FinalPrompt:      prompt,
		Reply:            fullResponse,
		EntityType:       entityType,
		EntityID:         entityID,
		PromptTokens:     promptTokens,
		CompletionTokens: completionTokens,
		TotalTokens:      totalTokensUsed,
	}
	gptCallId, err := ms.GPTCallRepo.CreateOne(&gptCall)
	if err != nil {
		fmt.Printf("Failed to log GPT call: %v\n", err)
		return fullResponse, 0, fmt.Errorf("failed to log GPT call: %w", err)
	}

	fmt.Printf("Successfully logged GPT call with ID: %d\n", gptCallId)
	return fullResponse, gptCallId, nil
}

// Helper function to check if the assistant's response indicates completion
func isResponseComplete(response string) bool {
	// Implement logic to determine if the response is complete
	// This could be based on detecting a closing JSON brace, specific keywords, etc.
	// For example, if the response ends with a closing brace or square bracket
	trimmedResponse := strings.TrimSpace(response)
	if strings.HasSuffix(trimmedResponse, "}") || strings.HasSuffix(trimmedResponse, "]") {
		return true
	}
	return false
}

// Helper function to extract partial content from a possibly incomplete JSON string
func extractPartialContent(response string) string {
	// Attempt to find the 'content' field in the partial response
	// You may need to use a regex or string manipulation to extract the available content
	// For simplicity, we'll return the raw response
	return response
}
