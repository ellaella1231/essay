package services

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
)

type AIService interface {
	ProcessEssay(ctx context.Context, imageURL, promptText, rubricText string) (*GradingResult, error)
}

type aiService struct {
	client *openai.Client
}

func NewAIService() AIService {
	apiKey := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(apiKey)
	return &aiService{client: client}
}

type GradingResult struct {
	Score          int           `json:"score"`
	PerfectVersion string        `json:"perfect_essay"`
	Errors         []ErrorDetail `json:"errors"`
}

type ErrorDetail struct {
	OriginalSegment  string `json:"original_segment"`
	SuggestedSegment string `json:"suggested_segment"`
	Explanation      string `json:"explanation"`
}

func (s *aiService) ProcessEssay(ctx context.Context, imagePath, promptText, rubricText string) (*GradingResult, error) {
	// Step 1: Vision - Transcribe handwriting
	
	// Read local image and encode to base64
	imgBytes, err := os.ReadFile(imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read image file: %w", err)
	}
	base64Img := base64.StdEncoding.EncodeToString(imgBytes)
	dataURI := fmt.Sprintf("data:image/jpeg;base64,%s", base64Img)

	ocrPrompt := "Transcribe the handwritten English text in this image exactly as written. Do not correct any mistakes."
	ocrReq := openai.ChatCompletionRequest{
		Model: openai.GPT4VisionPreview,
		Messages: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleUser,
				MultiContent: []openai.ChatMessagePart{
					{
						Type: openai.ChatMessagePartTypeText,
						Text: ocrPrompt,
					},
					{
						Type: openai.ChatMessagePartTypeImageURL,
						ImageURL: &openai.ChatMessageImageURL{
							URL: dataURI,
						},
					},
				},
			},
		},
	}

	ocrResp, err := s.client.CreateChatCompletion(ctx, ocrReq)
	if err != nil {
		return nil, fmt.Errorf("OCR failed: %w", err)
	}
	studentContent := ocrResp.Choices[0].Message.Content

	// Step 2: Analysis - Grade and Correct
	systemPrompt := fmt.Sprintf(`You are an expert English teacher. 
Grade the following student essay based on this rubric:
%s

Essay Prompt: %s

Student Essay:
%s

You must return a JSON object strictly matching this structure:
{
  "score": integer, // Score based on rubric
  "perfect_essay": string, // A "perfect" version maintaining original intent but elevating vocabulary/grammar
  "errors": [
    {
      "original_segment": string, // Exact text from student essay containing mistake
      "suggested_segment": string, // Corrected text
      "explanation": string // Why it was wrong
    }
  ]
}`, rubricText, promptText, studentContent)

	gradingReq := openai.ChatCompletionRequest{
		Model: openai.GPT4TurboPreview,
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONObject,
		},
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are an expert English teacher. You must return a JSON object strictly matching the required structure.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: systemPrompt,
			},
		},
	}

	gradingResp, err := s.client.CreateChatCompletion(ctx, gradingReq)
	if err != nil {
		return nil, fmt.Errorf("grading failed: %w", err)
	}

	var result GradingResult
	err = json.Unmarshal([]byte(gradingResp.Choices[0].Message.Content), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return &result, nil
}
