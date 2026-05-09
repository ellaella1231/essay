package services

import (
	"context"
	"testing"
)

func TestProcessEssayReturnsMockResult(t *testing.T) {
	service := NewAIService()

	result, err := service.ProcessEssay(context.Background(), "ignored.png", "My favorite sport", "rubric")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if result == nil {
		t.Fatal("expected result")
	}
	if result.Score <= 0 {
		t.Fatalf("expected positive score, got %d", result.Score)
	}
	if len(result.Errors) == 0 {
		t.Fatal("expected correction details")
	}
	if result.PerfectVersion == "" {
		t.Fatal("expected perfect essay")
	}
	if result.OriginalText == "" {
		t.Fatal("expected original text")
	}
}

func TestNormalizeCompatibleBaseURL(t *testing.T) {
	tests := map[string]string{
		"https://dashscope.aliyuncs.com/compatible-mode":                     "https://dashscope.aliyuncs.com/compatible-mode/v1",
		"https://dashscope.aliyuncs.com/compatible-mode/v1":                  "https://dashscope.aliyuncs.com/compatible-mode/v1",
		"https://dashscope.aliyuncs.com/compatible-mode/v1/":                 "https://dashscope.aliyuncs.com/compatible-mode/v1",
		"https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions": "https://dashscope.aliyuncs.com/compatible-mode/v1",
	}

	for input, want := range tests {
		if got := normalizeCompatibleBaseURL(input); got != want {
			t.Fatalf("normalizeCompatibleBaseURL(%q) = %q, want %q", input, got, want)
		}
	}
}
