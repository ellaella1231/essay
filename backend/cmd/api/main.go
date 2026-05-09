package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"ai-essay-master-backend/internal/core/services"
	"ai-essay-master-backend/internal/repository/memory"
)

func main() {
	loadLocalEnv(".env")

	userRepo := memory.NewUserRepository()
	log.Println("Using in-memory user repository")

	authService := services.NewAuthService(userRepo)
	essayService := services.NewAIService()
	logRuntimeConfig()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/ping", withCORS(func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"message": "pong"})
	}))
	mux.HandleFunc("/auth/verify", withCORS(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}

		var req struct {
			Phone string `json:"phone"`
			OTP   string `json:"otp"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
			return
		}

		token, user, err := authService.VerifyOTP(r.Context(), req.Phone, req.OTP)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"token": token,
			"user":  user,
		})
	}))
	mux.HandleFunc("/essays/grade", withCORS(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}
		if err := r.ParseMultipartForm(16 << 20); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid form data"})
			return
		}

		grade := r.FormValue("grade")
		if grade == "" {
			grade = "Grade 7"
		}
		promptText := r.FormValue("prompt_text")

		rubricPath := filepath.Join("rubrics", grade+".txt")
		rubricBytes, err := os.ReadFile(rubricPath)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to load rubric file"})
			return
		}

		tempFilePath, err := saveTempUpload(r.MultipartForm.File["image"])
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		defer os.Remove(tempFilePath)

		result, err := essayService.ProcessEssay(r.Context(), tempFilePath, promptText, string(rubricBytes))
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		writeJSON(w, http.StatusOK, result)
	}))

	log.Printf("API listening on http://127.0.0.1:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next(w, r)
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func saveTempUpload(files []*multipart.FileHeader) (string, error) {
	if len(files) == 0 {
		return "", http.ErrMissingFile
	}

	fileHeader := files[0]
	src, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	tempFile, err := os.CreateTemp("", "essay-*"+filepath.Ext(fileHeader.Filename))
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	if _, err := io.Copy(tempFile, src); err != nil {
		return "", err
	}

	return tempFile.Name(), nil
}

func loadLocalEnv(path string) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if key == "" {
			continue
		}
		_ = os.Setenv(key, value)
	}
}

func logRuntimeConfig() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("DASHSCOPE_API_KEY")
	}

	baseURL := os.Getenv("OPENAI_BASE_URL")
	if baseURL == "" {
		baseURL = os.Getenv("DASHSCOPE_BASE_URL")
	}
	if baseURL == "" {
		baseURL = "https://dashscope.aliyuncs.com/compatible-mode/v1"
	}

	model := os.Getenv("OPENAI_MODEL")
	if model == "" {
		model = os.Getenv("DASHSCOPE_MODEL")
	}
	if model == "" {
		model = "qwen-vl-max"
	}

	log.Printf("LLM config: model=%s base_url=%s api_key=%s", model, baseURL, maskSecret(apiKey))
}

func maskSecret(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "(missing)"
	}
	if len(value) <= 10 {
		return value
	}
	return value[:6] + "..." + value[len(value)-4:]
}
