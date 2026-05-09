package domain

import "time"

type User struct {
	ID                 int        `json:"id"`
	Phone              string     `json:"phone"`
	FreeTrialCount     int        `json:"free_trial_count"`
	SubscriptionExpiry *time.Time `json:"subscription_expiry"`
	CreatedAt          time.Time  `json:"created_at"`
}

type GradeStandard struct {
	Grade      string `json:"grade"`
	RubricText string `json:"rubric_text"`
}

type Essay struct {
	ID             int       `json:"id"`
	UserID         int       `json:"user_id"`
	ImageURL       string    `json:"image_url"`
	PromptText     string    `json:"prompt_text"`
	StudentContent string    `json:"student_content"`
	PerfectVersion string    `json:"perfect_version"`
	Score          int       `json:"score"`
	CreatedAt      time.Time `json:"created_at"`
}

type Error struct {
	ID               int    `json:"id"`
	EssayID          int    `json:"essay_id"`
	UserID           int    `json:"user_id"`
	OriginalSegment  string `json:"original_segment"`
	SuggestedSegment string `json:"suggested_segment"`
	Explanation      string `json:"explanation"`
	IsLearned        bool   `json:"is_learned"`
}
