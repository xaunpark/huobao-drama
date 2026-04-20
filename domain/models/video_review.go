package models

import (
	"time"

	"gorm.io/gorm"
)

// VideoReview stores AI quality review results for a specific video generation.
// Score is bound to VideoGenID — when the video changes (new gen), the old review
// naturally becomes irrelevant (frontend queries by current video_gen_id).
type VideoReview struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	VideoGenID   uint  `gorm:"not null;index" json:"video_gen_id"`
	StoryboardID *uint `gorm:"index" json:"storyboard_id,omitempty"`

	OverallScore float64 `gorm:"not null" json:"overall_score"`       // 0.0 - 10.0
	Verdict      string  `gorm:"type:varchar(20);not null" json:"verdict"` // excellent/good/acceptable/poor/unusable
	Dimensions   string  `gorm:"type:text" json:"dimensions"`         // JSON: {character_consistency, prompt_adherence, ...}
	Errors       string  `gorm:"type:text" json:"errors"`             // JSON array: [{severity, time_range, description}]
	FixGuide     string  `gorm:"type:text" json:"fix_guide"`

	FramesAnalyzed int     `json:"frames_analyzed"`
	FPSUsed        float64 `json:"fps_used"`
	HasCritical    bool    `gorm:"default:false" json:"has_critical_errors"`
}

func (VideoReview) TableName() string {
	return "video_reviews"
}
