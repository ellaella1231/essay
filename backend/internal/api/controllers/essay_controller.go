package controllers

import (
	"context"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"ai-essay-master-backend/internal/core/services"
)

type EssayController struct {
	aiService services.AIService
}

func NewEssayController(aiService services.AIService) *EssayController {
	return &EssayController{aiService: aiService}
}

func (c *EssayController) GradeEssay(ctx *gin.Context) {
	// 1. Get Prompt Text
	promptText := ctx.PostForm("prompt_text")
	grade := ctx.PostForm("grade")
	if grade == "" {
		grade = "Grade 7" // default
	}

	// 2. Read Rubric from hardcoded file in directory
	rubricPath := filepath.Join("rubrics", grade+".txt")
	rubricBytes, err := os.ReadFile(rubricPath)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load rubric file"})
		return
	}
	rubricText := string(rubricBytes)

	// 3. Handle Image Upload
	file, err := ctx.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Image file is required"})
		return
	}

	// Save file temporarily (in production, upload to S3/OSS)
	tempFilePath := filepath.Join(os.TempDir(), file.Filename)
	if err := ctx.SaveUploadedFile(file, tempFilePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}
	defer os.Remove(tempFilePath) // clean up

	// 4. Call AI Service (assuming AIService takes a local file path or public URL. 
	// For local files with OpenAI Vision, we need to pass base64 or have a public URL.
	// Let's assume aiService handles base64 if given a local path, or we pass the local path).
	result, err := c.aiService.ProcessEssay(context.Background(), tempFilePath, promptText, rubricText)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 5. Return JSON Result
	ctx.JSON(http.StatusOK, result)
}
