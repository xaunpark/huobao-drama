package services

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	models "github.com/drama-generator/backend/domain/models"
	"github.com/drama-generator/backend/application/prompts"
	"github.com/drama-generator/backend/pkg/ai"
	"github.com/drama-generator/backend/pkg/logger"
	"gorm.io/gorm"
)

// StyleDistillService handles batch distillation of general style/constraint prompts
// into per-shot concise styles after storyboard creation.
// See: plans/shot-style-distill.md (Phase 2-3)
type StyleDistillService struct {
	db         *gorm.DB
	aiService  *AIService
	promptI18n *PromptI18n
	log        *logger.Logger
}

func NewStyleDistillService(db *gorm.DB, aiService *AIService, promptI18n *PromptI18n, log *logger.Logger) *StyleDistillService {
	return &StyleDistillService{
		db:         db,
		aiService:  aiService,
		promptI18n: promptI18n,
		log:        log,
	}
}

// shotContext is the subset of storyboard fields sent to the LLM for distillation.
type shotContext struct {
	ShotNumber int    `json:"shot_number"`
	Action     string `json:"action,omitempty"`
	Result     string `json:"result,omitempty"`
	Location   string `json:"location,omitempty"`
	Atmosphere string `json:"atmosphere,omitempty"`
	ShotType   string `json:"shot_type,omitempty"`
	Angle      string `json:"angle,omitempty"`
	Movement   string `json:"movement,omitempty"`
	Characters string `json:"characters,omitempty"`
}

// distilledImageStyle is the expected JSON output from the image distill LLM call.
type distilledImageStyle struct {
	ShotNumber int    `json:"shot_number"`
	ImageStyle string `json:"image_style"`
}

// distilledVideoStyle is the expected JSON output from the video distill LLM call.
type distilledVideoStyle struct {
	ShotNumber int    `json:"shot_number"`
	VideoStyle string `json:"video_style"`
}

const maxShotsPerBatch = 20

// BatchDistillStyles runs 2 parallel LLM calls to distill style_prompt and video_constraint
// into per-shot image_style and video_style for all storyboards in an episode.
func (s *StyleDistillService) BatchDistillStyles(episodeID uint, dramaID uint) {
	s.log.Infow("Starting batch style distillation",
		"episode_id", episodeID,
		"drama_id", dramaID)

	// 1. Load all storyboards for this episode
	var storyboards []models.Storyboard
	if err := s.db.Where("episode_id = ?", episodeID).
		Order("storyboard_number ASC").
		Find(&storyboards).Error; err != nil {
		s.log.Errorw("Failed to load storyboards for distillation",
			"error", err, "episode_id", episodeID)
		return
	}

	if len(storyboards) == 0 {
		s.log.Infow("No storyboards found, skipping distillation", "episode_id", episodeID)
		return
	}

	// 2. Load drama to get style info
	var drama models.Drama
	if err := s.db.First(&drama, dramaID).Error; err != nil {
		s.log.Errorw("Failed to load drama for distillation",
			"error", err, "drama_id", dramaID)
		return
	}

	// 3. Resolve style_prompt and video_constraint from template
	stylePrompt := s.promptI18n.ResolveEffectiveStylePublic(dramaID, drama.Style, drama.CustomStyle)
	videoConstraint := s.promptI18n.WithDramaVideoConstraintPrompt(dramaID, "single")

	hasStylePrompt := stylePrompt != "" && len(stylePrompt) > 20
	hasVideoConstraint := videoConstraint != "" && len(videoConstraint) > 20

	if !hasStylePrompt && !hasVideoConstraint {
		s.log.Infow("No style_prompt or video_constraint to distill, skipping",
			"episode_id", episodeID, "drama_id", dramaID)
		return
	}

	// 4. Build shot context
	shotContexts := s.buildShotContexts(storyboards)

	// 5. Run distillation in parallel
	var wg sync.WaitGroup
	var imageStyles []distilledImageStyle
	var videoStyles []distilledVideoStyle
	var imageErr, videoErr error

	if hasStylePrompt {
		wg.Add(1)
		go func() {
			defer wg.Done()
			imageStyles, imageErr = s.distillImageStyles(stylePrompt, shotContexts)
		}()
	}

	if hasVideoConstraint {
		wg.Add(1)
		go func() {
			defer wg.Done()
			videoStyles, videoErr = s.distillVideoStyles(videoConstraint, shotContexts)
		}()
	}

	wg.Wait()

	// 6. Log errors but don't fail — fallback is to use original behavior
	if imageErr != nil {
		s.log.Errorw("Image style distillation failed, shots will use fallback style",
			"error", imageErr, "episode_id", episodeID)
	}
	if videoErr != nil {
		s.log.Errorw("Video style distillation failed, shots will use fallback behavior",
			"error", videoErr, "episode_id", episodeID)
	}

	// 7. Save distilled styles to DB
	s.saveDistilledStyles(storyboards, imageStyles, videoStyles)

	s.log.Infow("Batch style distillation completed",
		"episode_id", episodeID,
		"total_shots", len(storyboards),
		"image_styles_distilled", len(imageStyles),
		"video_styles_distilled", len(videoStyles))
}

// buildShotContexts creates the shot context array for the LLM prompt.
func (s *StyleDistillService) buildShotContexts(storyboards []models.Storyboard) []shotContext {
	contexts := make([]shotContext, 0, len(storyboards))

	for _, sb := range storyboards {
		ctx := shotContext{
			ShotNumber: sb.StoryboardNumber,
		}
		if sb.Action != nil {
			ctx.Action = *sb.Action
		}
		if sb.Result != nil {
			ctx.Result = *sb.Result
		}
		if sb.Location != nil {
			ctx.Location = *sb.Location
		}
		if sb.Atmosphere != nil {
			ctx.Atmosphere = *sb.Atmosphere
		}
		if sb.ShotType != nil {
			ctx.ShotType = *sb.ShotType
		}
		if sb.Angle != nil {
			ctx.Angle = *sb.Angle
		}
		if sb.Movement != nil {
			ctx.Movement = *sb.Movement
		}

		// Load character names for this storyboard
		var characters []models.Character
		if err := s.db.Model(&sb).Association("Characters").Find(&characters); err == nil && len(characters) > 0 {
			names := make([]string, len(characters))
			for i, c := range characters {
				names[i] = c.Name
			}
			ctx.Characters = strings.Join(names, ", ")
		}

		contexts = append(contexts, ctx)
	}

	return contexts
}

// distillImageStyles calls the LLM to distill style_prompt into per-shot image styles.
func (s *StyleDistillService) distillImageStyles(stylePrompt string, shots []shotContext) ([]distilledImageStyle, error) {
	template := prompts.Get("image_style_distill.txt")
	if template == "" {
		return nil, fmt.Errorf("image_style_distill.txt prompt not found")
	}

	var allResults []distilledImageStyle

	// Chunk shots if needed
	for i := 0; i < len(shots); i += maxShotsPerBatch {
		end := i + maxShotsPerBatch
		if end > len(shots) {
			end = len(shots)
		}
		batch := shots[i:end]

		shotsJSON, err := json.Marshal(batch)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal shot contexts: %w", err)
		}

		prompt := fmt.Sprintf(template, stylePrompt, string(shotsJSON))

		response, err := s.aiService.GenerateText(prompt, "", ai.WithMaxTokens(4000))
		if err != nil {
			return nil, fmt.Errorf("LLM call failed for image style distill: %w", err)
		}

		var batchResults []distilledImageStyle
		if err := parseJSONArray(response, &batchResults); err != nil {
			s.log.Warnw("Failed to parse image style distill response, trying robust extraction",
				"error", err, "response_length", len(response))
			if err := parseJSONArray(extractJSONArray(response), &batchResults); err != nil {
				return allResults, fmt.Errorf("failed to parse image style response: %w", err)
			}
		}

		allResults = append(allResults, batchResults...)
	}

	return allResults, nil
}

// distillVideoStyles calls the LLM to distill video_constraint into per-shot video styles.
func (s *StyleDistillService) distillVideoStyles(videoConstraint string, shots []shotContext) ([]distilledVideoStyle, error) {
	template := prompts.Get("video_style_distill.txt")
	if template == "" {
		return nil, fmt.Errorf("video_style_distill.txt prompt not found")
	}

	var allResults []distilledVideoStyle

	// Chunk shots if needed
	for i := 0; i < len(shots); i += maxShotsPerBatch {
		end := i + maxShotsPerBatch
		if end > len(shots) {
			end = len(shots)
		}
		batch := shots[i:end]

		shotsJSON, err := json.Marshal(batch)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal shot contexts: %w", err)
		}

		prompt := fmt.Sprintf(template, videoConstraint, string(shotsJSON))

		response, err := s.aiService.GenerateText(prompt, "", ai.WithMaxTokens(4000))
		if err != nil {
			return nil, fmt.Errorf("LLM call failed for video style distill: %w", err)
		}

		var batchResults []distilledVideoStyle
		if err := parseJSONArray(response, &batchResults); err != nil {
			s.log.Warnw("Failed to parse video style distill response, trying robust extraction",
				"error", err, "response_length", len(response))
			if err := parseJSONArray(extractJSONArray(response), &batchResults); err != nil {
				return allResults, fmt.Errorf("failed to parse video style response: %w", err)
			}
		}

		allResults = append(allResults, batchResults...)
	}

	return allResults, nil
}

// saveDistilledStyles updates each storyboard with its distilled image_style and video_style.
func (s *StyleDistillService) saveDistilledStyles(storyboards []models.Storyboard, imageStyles []distilledImageStyle, videoStyles []distilledVideoStyle) {
	// Build lookup maps by shot_number
	imageMap := make(map[int]string)
	for _, is := range imageStyles {
		imageMap[is.ShotNumber] = is.ImageStyle
	}
	videoMap := make(map[int]string)
	for _, vs := range videoStyles {
		videoMap[vs.ShotNumber] = vs.VideoStyle
	}

	for _, sb := range storyboards {
		updates := make(map[string]interface{})

		if style, ok := imageMap[sb.StoryboardNumber]; ok && style != "" {
			updates["image_style"] = style
		}
		if style, ok := videoMap[sb.StoryboardNumber]; ok && style != "" {
			updates["video_style"] = style
		}

		if len(updates) > 0 {
			if err := s.db.Model(&models.Storyboard{}).Where("id = ?", sb.ID).Updates(updates).Error; err != nil {
				s.log.Errorw("Failed to save distilled style for shot",
					"error", err,
					"storyboard_id", sb.ID,
					"shot_number", sb.StoryboardNumber)
			}
		}
	}
}

// parseJSONArray attempts to unmarshal a JSON string into the target slice.
func parseJSONArray(data string, target interface{}) error {
	return json.Unmarshal([]byte(data), target)
}

// extractJSONArray attempts to find a JSON array in a potentially messy LLM response.
// It looks for the first '[' and the last ']' and extracts the substring.
func extractJSONArray(response string) string {
	start := strings.Index(response, "[")
	end := strings.LastIndex(response, "]")
	if start == -1 || end == -1 || end <= start {
		return response
	}
	return response[start : end+1]
}
