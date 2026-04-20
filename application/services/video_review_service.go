package services

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/drama-generator/backend/domain/models"
	"github.com/drama-generator/backend/infrastructure/external/ffmpeg"
	"github.com/drama-generator/backend/infrastructure/storage"
	"github.com/drama-generator/backend/pkg/ai"
	"github.com/drama-generator/backend/pkg/logger"
	"gorm.io/gorm"
)

const (
	reviewFPS  = 8.0
	reviewCols = 8
	// Gemini model name for Locally-AI
	reviewModel = "gemini/auto"
)

// Dimension weights for overall score calculation (must sum to 1.0)
var dimensionWeights = map[string]float64{
	"character_consistency": 0.35,
	"prompt_adherence":      0.25,
	"text_accuracy":         0.15,
	"entity_stability":      0.25,
}

type VideoReviewService struct {
	db           *gorm.DB
	log          *logger.Logger
	aiService    *AIService
	taskService  *TaskService
	ffmpeg       *ffmpeg.FFmpeg
	localStorage *storage.LocalStorage
}

func NewVideoReviewService(db *gorm.DB, log *logger.Logger, aiService *AIService, taskService *TaskService, localStorage *storage.LocalStorage) *VideoReviewService {
	return &VideoReviewService{
		db:           db,
		log:          log,
		aiService:    aiService,
		taskService:  taskService,
		ffmpeg:       ffmpeg.NewFFmpeg(log),
		localStorage: localStorage,
	}
}

// ReviewVideoAsync triggers an async video review task.
// Returns the task ID for frontend polling.
func (s *VideoReviewService) ReviewVideoAsync(videoGenID uint) (string, error) {
	// Verify video exists and has a local path or URL
	var video models.VideoGeneration
	if err := s.db.First(&video, videoGenID).Error; err != nil {
		return "", fmt.Errorf("video generation not found: %w", err)
	}

	if video.LocalPath == nil && video.VideoURL == nil {
		return "", fmt.Errorf("video %d has no file path or URL", videoGenID)
	}

	// Create async task
	task, err := s.taskService.CreateTask("video_review", fmt.Sprintf("%d", videoGenID))
	if err != nil {
		return "", fmt.Errorf("failed to create review task: %w", err)
	}

	// Run review in background goroutine
	go s.executeReview(videoGenID, task.ID)

	return task.ID, nil
}

// GetLatestReview returns the most recent review for a video generation.
func (s *VideoReviewService) GetLatestReview(videoGenID uint) (*models.VideoReview, error) {
	var review models.VideoReview
	err := s.db.Where("video_gen_id = ?", videoGenID).
		Order("created_at DESC").
		First(&review).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No review yet — not an error
		}
		return nil, err
	}
	return &review, nil
}

// executeReview is the core review logic, run in a goroutine.
func (s *VideoReviewService) executeReview(videoGenID uint, taskID string) {
	defer func() {
		if r := recover(); r != nil {
			s.log.Errorw("Video review panicked", "videoGenID", videoGenID, "panic", r)
			s.taskService.UpdateTaskError(taskID, fmt.Errorf("review panicked: %v", r))
		}
	}()

	s.log.Infow("Starting video review", "videoGenID", videoGenID, "taskID", taskID)
	s.taskService.UpdateTaskStatus(taskID, "processing", 10, "Preparing video for review...")

	// 1. Fetch video info
	var video models.VideoGeneration
	if err := s.db.First(&video, videoGenID).Error; err != nil {
		s.taskService.UpdateTaskError(taskID, fmt.Errorf("video not found: %w", err))
		return
	}

	// 2. Resolve video file path (local_path is stored as RELATIVE path in DB)
	videoPath := ""
	if video.LocalPath != nil && *video.LocalPath != "" {
		// Convert relative path to absolute using localStorage
		if s.localStorage != nil {
			videoPath = s.localStorage.GetAbsolutePath(*video.LocalPath)
		} else {
			// Fallback: try as-is (might be absolute already)
			videoPath = *video.LocalPath
		}
	} else if video.VideoURL != nil && *video.VideoURL != "" {
		// For URL-only videos, we'd need to download first
		// For now, only local path is supported
		s.taskService.UpdateTaskError(taskID, fmt.Errorf("video review requires local_path, but video %d only has URL", videoGenID))
		return
	}

	if videoPath == "" {
		s.taskService.UpdateTaskError(taskID, fmt.Errorf("video %d has no local_path set", videoGenID))
		return
	}

	if _, err := os.Stat(videoPath); os.IsNotExist(err) {
		s.taskService.UpdateTaskError(taskID, fmt.Errorf("video file not found at: %s (local_path in DB: %s)", videoPath, *video.LocalPath))
		return
	}

	// 3. Get storyboard context for prompt and characters
	var storyboard models.Storyboard
	imagePrompt := ""
	videoPrompt := ""
	action := ""
	var additionalImages []string

	if video.StoryboardID != nil {
		if err := s.db.Preload("Characters").First(&storyboard, *video.StoryboardID).Error; err == nil {
			if storyboard.ImagePrompt != nil {
				imagePrompt = *storyboard.ImagePrompt
			}
			if storyboard.VideoPrompt != nil {
				videoPrompt = *storyboard.VideoPrompt
			}
			if storyboard.Action != nil {
				action = *storyboard.Action
			}

			// Load character images
			for _, char := range storyboard.Characters {
				if char.LocalPath != nil && *char.LocalPath != "" {
					absoluteCharPath := ""
					if s.localStorage != nil {
						absoluteCharPath = s.localStorage.GetAbsolutePath(*char.LocalPath)
					} else {
						absoluteCharPath = *char.LocalPath
					}

					if _, err := os.Stat(absoluteCharPath); err == nil {
						additionalImages = append(additionalImages, absoluteCharPath)
					}
				}
			}
		}
	}
	// Fallback: use video's own prompt
	if imagePrompt == "" {
		imagePrompt = video.Prompt
	}

	s.taskService.UpdateTaskStatus(taskID, "processing", 20, "Generating contact sheet...")

	// 4. Create contact sheet
	tmpDir, err := os.MkdirTemp("", "video-review-*")
	if err != nil {
		s.taskService.UpdateTaskError(taskID, fmt.Errorf("failed to create temp dir: %w", err))
		return
	}
	defer os.RemoveAll(tmpDir)

	contactSheetPath, totalFrames, err := s.ffmpeg.CreateContactSheet(videoPath, reviewFPS, tmpDir, reviewCols)
	if err != nil {
		s.taskService.UpdateTaskError(taskID, fmt.Errorf("contact sheet creation failed: %w", err))
		return
	}

	s.taskService.UpdateTaskStatus(taskID, "processing", 40, "Sending to Gemini Vision for analysis...")

	// 5. Build review prompt
	prompt := s.buildReviewPrompt(totalFrames, reviewFPS, imagePrompt, videoPrompt, action)

	// 6. Get Locally-AI config
	config, err := s.aiService.GetDefaultConfig("text")
	if err != nil {
		s.taskService.UpdateTaskError(taskID, fmt.Errorf("no AI text config found: %w", err))
		return
	}

	// 7. Send multimodal request to Locally-AI (Gemini)
	imagePaths := []string{contactSheetPath}
	imagePaths = append(imagePaths, additionalImages...)
	visionReq := ai.BuildVisionRequest(reviewModel, prompt, imagePaths)
	responseText, err := ai.SendMultimodal(config.BaseURL, config.APIKey, visionReq)
	if err != nil {
		s.taskService.UpdateTaskError(taskID, fmt.Errorf("Gemini Vision API failed: %w", err))
		return
	}

	s.taskService.UpdateTaskStatus(taskID, "processing", 70, "Parsing review results...")

	// 8. Parse response JSON
	reviewData, err := parseReviewResponse(responseText)
	if err != nil {
		s.taskService.UpdateTaskError(taskID, fmt.Errorf("failed to parse review response: %w (raw: %s)", err, truncate(responseText, 500)))
		return
	}

	// 9. Calculate overall score
	dims := reviewData.Dimensions
	overallScore := computeOverallScore(dims)

	// Check for critical errors and enforce scoring rules
	hasCritical := false
	for _, e := range reviewData.Errors {
		if strings.ToUpper(e.Severity) == "CRITICAL" {
			hasCritical = true
			break
		}
	}

	// Enforce: critical errors cap character_consistency at 3.0
	if hasCritical {
		if val, ok := dims["character_consistency"]; ok && val != nil && *val > 3.0 {
			capVal := 3.0
			dims["character_consistency"] = &capVal
			overallScore = computeOverallScore(dims)
		}
	}

	// Force score cap for critical errors
	if hasCritical && overallScore > 5.9 {
		overallScore = 5.9
	}

	verdict := computeVerdict(overallScore)

	// 10. Marshal JSON fields
	dimsJSON, _ := json.Marshal(dims)
	errorsJSON, _ := json.Marshal(reviewData.Errors)

	fixGuide := ""
	if reviewData.FixGuide != "" {
		fixGuide = reviewData.FixGuide
	}

	// 11. Save to DB
	review := &models.VideoReview{
		VideoGenID:     videoGenID,
		StoryboardID:   video.StoryboardID,
		OverallScore:   overallScore,
		Verdict:        verdict,
		Dimensions:     string(dimsJSON),
		Errors:         string(errorsJSON),
		FixGuide:       fixGuide,
		FramesAnalyzed: totalFrames,
		FPSUsed:        reviewFPS,
		HasCritical:    hasCritical,
	}

	if err := s.db.Create(review).Error; err != nil {
		s.taskService.UpdateTaskError(taskID, fmt.Errorf("failed to save review: %w", err))
		return
	}

	s.log.Infow("Video review completed",
		"videoGenID", videoGenID,
		"score", overallScore,
		"verdict", verdict,
		"hasCritical", hasCritical,
	)

	// 12. Mark task complete
	s.taskService.UpdateTaskResult(taskID, review)
}

// buildReviewPrompt loads the rubric template and fills in placeholders
func (s *VideoReviewService) buildReviewPrompt(nFrames int, fps float64, imagePrompt, videoPrompt, action string) string {
	// Read template
	templateBytes, err := os.ReadFile("application/prompts/video_review_rubric.txt")
	if err != nil {
		s.log.Warnw("Failed to read review rubric template, using inline", "error", err)
		return fmt.Sprintf("Review this contact sheet of %d frames from an AI-generated video. Score quality 0-10. Return JSON only.", nFrames)
	}

	template := string(templateBytes)
	template = strings.ReplaceAll(template, "{n_frames}", fmt.Sprintf("%d", nFrames))
	template = strings.ReplaceAll(template, "{fps}", fmt.Sprintf("%.0f", fps))
	template = strings.ReplaceAll(template, "{cols}", fmt.Sprintf("%d", reviewCols))
	template = strings.ReplaceAll(template, "{image_prompt}", imagePrompt)
	template = strings.ReplaceAll(template, "{video_prompt}", videoPrompt)
	template = strings.ReplaceAll(template, "{action}", action)

	return template
}

// --- Response parsing ---

type reviewResponseData struct {
	Dimensions map[string]*float64 `json:"dimensions"`
	Errors     []reviewError       `json:"errors"`
	Segments   []reviewSegment     `json:"usable_segments"`
	FixGuide   string              `json:"fix_guide"`
}

type reviewError struct {
	Severity             string `json:"severity"`
	TimeRange            string `json:"time_range"`
	Description          string `json:"description"`
	Category             string `json:"category"`
	NarrativelyJustified bool   `json:"narratively_justified"`
}

type reviewSegment struct {
	TimeRange string  `json:"time_range"`
	Score     float64 `json:"score"`
}

func parseReviewResponse(raw string) (*reviewResponseData, error) {
	raw = strings.TrimSpace(raw)

	// Strip markdown fences if present
	if strings.HasPrefix(raw, "```") {
		parts := strings.SplitN(raw, "\n", 2)
		if len(parts) > 1 {
			raw = parts[1]
		}
		if idx := strings.LastIndex(raw, "```"); idx >= 0 {
			raw = raw[:idx]
		}
		raw = strings.TrimSpace(raw)
	}

	// Find JSON object in response
	if !strings.HasPrefix(raw, "{") {
		start := strings.Index(raw, "{")
		if start < 0 {
			return nil, fmt.Errorf("no JSON object found in response")
		}
		raw = raw[start:]
	}

	// Find matching closing brace
	depth := 0
	end := -1
	for i, ch := range raw {
		if ch == '{' {
			depth++
		} else if ch == '}' {
			depth--
			if depth == 0 {
				end = i + 1
				break
			}
		}
	}
	if end > 0 {
		raw = raw[:end]
	}

	var data reviewResponseData
	if err := json.Unmarshal([]byte(raw), &data); err != nil {
		return nil, fmt.Errorf("JSON parse error: %w", err)
	}

	// Validate dimensions exist
	if data.Dimensions == nil {
		return nil, fmt.Errorf("response missing 'dimensions' field")
	}

	// Fill default values for missing dimensions
	for key := range dimensionWeights {
		if _, ok := data.Dimensions[key]; !ok {
			if key != "text_accuracy" {
				defaultVal := 5.0
				data.Dimensions[key] = &defaultVal
			}
		}
	}

	return &data, nil
}

func computeOverallScore(dims map[string]*float64) float64 {
	totalScore := 0.0
	totalWeight := 0.0

	for key, weight := range dimensionWeights {
		scorePtr, ok := dims[key]
		if ok && scorePtr != nil {
			totalScore += *scorePtr * weight
			totalWeight += weight
		}
	}

	if totalWeight == 0 {
		return 5.0
	}

	finalScore := totalScore / totalWeight
	return math.Round(finalScore*100) / 100 // round to 2 decimal places
}

func computeVerdict(score float64) string {
	switch {
	case score >= 9.0:
		return "excellent"
	case score >= 7.5:
		return "good"
	case score >= 6.0:
		return "acceptable"
	case score >= 4.0:
		return "poor"
	default:
		return "unusable"
	}
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
