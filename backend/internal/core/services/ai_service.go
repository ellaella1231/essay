package services

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type AIService interface {
	ProcessEssay(ctx context.Context, imageURL, promptText, rubricText string) (*GradingResult, error)
}

type aiService struct {
	apiKey     string
	baseURL    string
	model      string
	httpClient *http.Client
}

func NewAIService() AIService {
	baseURL := strings.TrimSpace(os.Getenv("OPENAI_BASE_URL"))
	if baseURL == "" {
		baseURL = strings.TrimSpace(os.Getenv("DASHSCOPE_BASE_URL"))
	}
	if baseURL == "" {
		baseURL = "https://dashscope.aliyuncs.com/compatible-mode/v1"
	}

	model := strings.TrimSpace(os.Getenv("OPENAI_MODEL"))
	if model == "" {
		model = strings.TrimSpace(os.Getenv("DASHSCOPE_MODEL"))
	}
	if model == "" {
		model = "qwen-vl-max"
	}

	apiKey := strings.TrimSpace(os.Getenv("OPENAI_API_KEY"))
	if apiKey == "" {
		apiKey = strings.TrimSpace(os.Getenv("DASHSCOPE_API_KEY"))
	}

	return &aiService{
		apiKey:  apiKey,
		baseURL: normalizeCompatibleBaseURL(baseURL),
		model:   model,
		httpClient: &http.Client{
			Timeout: 90 * time.Second,
		},
	}
}

type GradingResult struct {
	Score          int           `json:"score"`
	PerfectVersion string        `json:"perfect_essay"`
	OriginalText   string        `json:"original_text,omitempty"`
	Errors         []ErrorDetail `json:"errors"`
}

type ErrorDetail struct {
	OriginalSegment  string `json:"original_segment"`
	SuggestedSegment string `json:"suggested_segment"`
	Explanation      string `json:"explanation"`
}

func (s *aiService) ProcessEssay(ctx context.Context, imagePath, promptText, rubricText string) (*GradingResult, error) {
	if s.apiKey == "" {
		return mockGradingResult(promptText), nil
	}

	imageDataURI, err := buildImageDataURI(imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare image: %w", err)
	}

	studentContent, err := s.transcribeEssay(ctx, imageDataURI)
	if err != nil {
		return nil, err
	}

	result, err := s.gradeEssay(ctx, studentContent, promptText, rubricText)
	if err != nil {
		return nil, err
	}
	result.OriginalText = studentContent
	return result, nil
}

func (s *aiService) transcribeEssay(ctx context.Context, imageDataURI string) (string, error) {
	reqBody := map[string]any{
		"model": s.model,
		"messages": []map[string]any{
			{
				"role": "user",
				"content": []map[string]any{
					{
						"type": "text",
						"text": "Transcribe the handwritten English text in this image exactly as written. Do not correct any mistakes.",
					},
					{
						"type": "image_url",
						"image_url": map[string]string{
							"url": imageDataURI,
						},
					},
				},
			},
		},
	}

	resp, err := s.createChatCompletion(ctx, reqBody)
	if err != nil {
		return "", fmt.Errorf("OCR failed: %w", err)
	}

	content := strings.TrimSpace(resp.FirstMessageContent())
	if content == "" {
		return "", fmt.Errorf("OCR failed: empty transcription")
	}
	return content, nil
}

func (s *aiService) gradeEssay(ctx context.Context, studentContent, promptText, rubricText string) (*GradingResult, error) {
	systemPrompt := fmt.Sprintf(`You are an expert English teacher.
Grade the following student essay based on this rubric:
%s

Essay Prompt: %s

Student Essay:
%s

You must return a JSON object strictly matching this structure:
{
  "score": integer,
  "perfect_essay": string,
  "errors": [
    {
      "original_segment": string,
      "suggested_segment": string,
      "explanation": string
    }
  ]
}`, rubricText, promptText, studentContent)

	reqBody := map[string]any{
		"model": s.model,
		"response_format": map[string]string{
			"type": "json_object",
		},
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You are an expert English teacher. You must return a JSON object strictly matching the required structure.",
			},
			{
				"role":    "user",
				"content": systemPrompt,
			},
		},
	}

	resp, err := s.createChatCompletion(ctx, reqBody)
	if err != nil {
		return nil, fmt.Errorf("grading failed: %w", err)
	}

	content := strings.TrimSpace(resp.FirstMessageContent())
	if content == "" {
		return nil, fmt.Errorf("grading failed: empty model response")
	}

	var result GradingResult
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, fmt.Errorf("grading failed: invalid JSON response: %w", err)
	}
	return &result, nil
}

func (s *aiService) createChatCompletion(ctx context.Context, payload map[string]any) (*chatCompletionResponse, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.baseURL+"/chat/completions", strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("provider returned %s: %s", resp.Status, strings.TrimSpace(string(raw)))
	}

	var parsed chatCompletionResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, err
	}
	if len(parsed.Choices) == 0 {
		return nil, fmt.Errorf("provider returned no choices")
	}
	return &parsed, nil
}

func buildImageDataURI(imagePath string) (string, error) {
	imageBytes, err := os.ReadFile(imagePath)
	if err != nil {
		return "", err
	}

	mimeType := mime.TypeByExtension(strings.ToLower(filepath.Ext(imagePath)))
	if mimeType == "" {
		mimeType = http.DetectContentType(imageBytes)
	}
	if mimeType == "" {
		mimeType = "image/jpeg"
	}

	return fmt.Sprintf("data:%s;base64,%s", mimeType, base64.StdEncoding.EncodeToString(imageBytes)), nil
}

func normalizeCompatibleBaseURL(baseURL string) string {
	trimmed := strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if strings.HasSuffix(trimmed, "/chat/completions") {
		return strings.TrimSuffix(trimmed, "/chat/completions")
	}
	if strings.HasSuffix(trimmed, "/v1") {
		return trimmed
	}
	return trimmed + "/v1"
}

type chatCompletionResponse struct {
	Choices []struct {
		Message struct {
			Content any `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (r *chatCompletionResponse) FirstMessageContent() string {
	if len(r.Choices) == 0 {
		return ""
	}

	switch content := r.Choices[0].Message.Content.(type) {
	case string:
		return content
	case []any:
		var parts []string
		for _, item := range content {
			itemMap, ok := item.(map[string]any)
			if !ok {
				continue
			}
			if text, ok := itemMap["text"].(string); ok && text != "" {
				parts = append(parts, text)
			}
		}
		return strings.Join(parts, "\n")
	default:
		return ""
	}
}

func mockGradingResult(promptText string) *GradingResult {
	perfectEssay := "I really enjoy playing football. I play with my friends every day, and it makes me very happy."
	if strings.TrimSpace(promptText) != "" {
		perfectEssay = fmt.Sprintf("Prompt: %s\n\n%s", strings.TrimSpace(promptText), perfectEssay)
	}

	return &GradingResult{
		Score:          85,
		PerfectVersion: perfectEssay,
		OriginalText:   "I very like play football. Every day I playing with my friends. It make me happy.",
		Errors: []ErrorDetail{
			{
				OriginalSegment:  "I very like play",
				SuggestedSegment: "I really enjoy playing",
				Explanation:      "Use a natural adverb-verb phrase, and follow enjoy with an -ing verb.",
			},
			{
				OriginalSegment:  "Every day I playing",
				SuggestedSegment: "I play every day",
				Explanation:      "Use the present simple tense for repeated daily actions.",
			},
			{
				OriginalSegment:  "It make me happy",
				SuggestedSegment: "It makes me happy",
				Explanation:      "It is third-person singular, so the verb takes an s.",
			},
		},
	}
}
