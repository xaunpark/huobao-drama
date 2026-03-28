package services

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	models "github.com/drama-generator/backend/domain/models"
	"github.com/drama-generator/backend/application/prompts"
	"github.com/drama-generator/backend/pkg/config"
	"github.com/drama-generator/backend/pkg/logger"
	"github.com/drama-generator/backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RapidCutService struct {
	db          *gorm.DB
	aiService   *AIService
	taskService *TaskService
	log         *logger.Logger
	config      *config.Config
	promptI18n  *PromptI18n
}

func NewRapidCutService(db *gorm.DB, cfg *config.Config, log *logger.Logger) *RapidCutService {
	return &RapidCutService{
		db:          db,
		aiService:   NewAIService(db, log),
		taskService: NewTaskService(db, log),
		log:         log,
		config:      cfg,
		promptI18n:  NewPromptI18n(cfg),
	}
}

// RapidCutShot represents the AI output for a merged production shot
type RapidCutShot struct {
	SourceShotNumbers []int  `json:"source_shot_numbers"`
	Title             string `json:"title"`
	Action            string `json:"action"`
	Result            string `json:"result"`
	ShotType          string `json:"shot_type"`
	Angle             string `json:"angle"`
	Movement          string `json:"movement"`
	Location          string `json:"location"`
	Time              string `json:"time"`
	Atmosphere        string `json:"atmosphere"`
	Emotion           string `json:"emotion"`
	Duration          int    `json:"duration"`
	Dialogue          string `json:"dialogue"`
	BgmPrompt         string `json:"bgm_prompt"`
	SoundEffect       string        `json:"sound_effect"`
	Characters        []interface{} `json:"characters"`
	SceneID           interface{}   `json:"scene_id"`
	IsPrimary         bool          `json:"is_primary"`
}

// GenerateRapidCut creates production shots by merging editorial shots
func (s *RapidCutService) GenerateRapidCut(episodeID string, model string) (string, error) {
	epID, err := strconv.ParseUint(episodeID, 10, 32)
	if err != nil {
		return "", fmt.Errorf("invalid episode ID: %s", episodeID)
	}

	// Verify episode exists
	var episode struct {
		ID      string
		DramaID string
	}
	if err := s.db.Table("episodes").
		Select("id, drama_id").
		Where("id = ?", episodeID).
		First(&episode).Error; err != nil {
		return "", fmt.Errorf("episode not found")
	}

	// Get editorial shots (non-production, ordered)
	var editorialShots []models.Storyboard
	if err := s.db.Where("episode_id = ? AND is_production = ?", uint(epID), false).
		Order("storyboard_number ASC").
		Preload("Characters").
		Preload("Props").
		Find(&editorialShots).Error; err != nil {
		return "", fmt.Errorf("failed to get editorial shots: %w", err)
	}

	if len(editorialShots) == 0 {
		return "", fmt.Errorf("no editorial shots found for this episode. Generate storyboard first")
	}

	// Create async task
	task, err := s.taskService.CreateTask("rapid_cut_generation", episodeID)
	if err != nil {
		return "", fmt.Errorf("failed to create task: %w", err)
	}

	s.log.Infow("Generating rapid cut production shots",
		"task_id", task.ID,
		"episode_id", episodeID,
		"editorial_shot_count", len(editorialShots))

	dramaIDUint, _ := strconv.ParseUint(episode.DramaID, 10, 32)
	go s.processRapidCutGeneration(task.ID, episodeID, uint(dramaIDUint), model, editorialShots)

	return task.ID, nil
}

// processRapidCutGeneration handles the async rapid cut generation
func (s *RapidCutService) processRapidCutGeneration(taskID, episodeID string, dramaID uint, model string, editorialShots []models.Storyboard) {
	if err := s.taskService.UpdateTaskStatus(taskID, "processing", 10, "Preparing rapid cut merge..."); err != nil {
		s.log.Errorw("Failed to update task status", "error", err, "task_id", taskID)
		return
	}

	// Build the shots list for the AI prompt
	shotsJSON := s.buildShotsDescription(editorialShots)

	// Build the merge prompt
	mergePrompt := prompts.Get("rapid_cut_merge.txt")

	prompt := fmt.Sprintf(`%s

[Input: Editorial Storyboard Shots]
The following is the complete list of editorial shots to merge into rapid cut production units:

%s

[Task]
Analyze the above shots and merge adjacent short shots into production units following the rules. 
Keep complex/long shots solo. Output a JSON array of production units.`, mergePrompt, shotsJSON)

	s.log.Infow("Rapid cut prompt built",
		"task_id", taskID,
		"prompt_length", len(prompt),
		"shot_count", len(editorialShots))

	// Get AI client
	client, getErr := s.aiService.GetAIClientForModel("text", model)
	if model != "" && getErr != nil {
		s.log.Warnw("Failed to get client for specified model, using default", "model", model, "error", getErr, "task_id", taskID)
	}

	var text string
	var err error
	if model != "" && getErr == nil {
		text, err = client.GenerateText(prompt, "")
	} else {
		text, err = s.aiService.GenerateText(prompt, "")
	}

	if err != nil {
		s.log.Errorw("Failed to generate rapid cut", "error", err, "task_id", taskID)
		if updateErr := s.taskService.UpdateTaskError(taskID, fmt.Errorf("rapid cut generation failed: %w", err)); updateErr != nil {
			s.log.Errorw("Failed to update task error", "error", updateErr, "task_id", taskID)
		}
		return
	}

	// Parse AI response
	var rapidCutShots []RapidCutShot
	if err := utils.SafeParseAIJSON(text, &rapidCutShots); err != nil {
		s.log.Errorw("Failed to parse rapid cut JSON", "error", err, "response", text[:min(500, len(text))], "task_id", taskID)
		if updateErr := s.taskService.UpdateTaskError(taskID, fmt.Errorf("failed to parse rapid cut result: %w", err)); updateErr != nil {
			s.log.Errorw("Failed to update task error", "error", updateErr, "task_id", taskID)
		}
		return
	}

	if err := s.taskService.UpdateTaskStatus(taskID, "processing", 50, "Rapid cut merge complete, saving production shots..."); err != nil {
		s.log.Errorw("Failed to update task status", "error", err, "task_id", taskID)
	}

	s.log.Infow("Rapid cut merge result",
		"task_id", taskID,
		"production_shot_count", len(rapidCutShots),
		"original_shot_count", len(editorialShots))

	// Save production shots
	epID, _ := strconv.ParseUint(episodeID, 10, 32)
	if err := s.saveProductionShots(uint(epID), dramaID, editorialShots, rapidCutShots); err != nil {
		s.log.Errorw("Failed to save production shots", "error", err, "task_id", taskID)
		if updateErr := s.taskService.UpdateTaskError(taskID, fmt.Errorf("failed to save production shots: %w", err)); updateErr != nil {
			s.log.Errorw("Failed to update task error", "error", updateErr, "task_id", taskID)
		}
		return
	}

	// Update task result
	resultData := gin.H{
		"production_shots":    len(rapidCutShots),
		"editorial_shots":     len(editorialShots),
		"compression_ratio":   fmt.Sprintf("%.1f%%", float64(len(rapidCutShots))/float64(len(editorialShots))*100),
	}
	if err := s.taskService.UpdateTaskResult(taskID, resultData); err != nil {
		s.log.Errorw("Failed to update task result", "error", err, "task_id", taskID)
	}

	s.log.Infow("Rapid cut generation completed",
		"task_id", taskID,
		"episode_id", episodeID,
		"production_shots", len(rapidCutShots))
}

// buildShotsDescription creates a text description of shots for the AI
func (s *RapidCutService) buildShotsDescription(shots []models.Storyboard) string {
	var parts []string
	for _, shot := range shots {
		charIDs := make([]uint, 0)
		for _, c := range shot.Characters {
			charIDs = append(charIDs, c.ID)
		}

		desc := fmt.Sprintf(`Shot #%d:
  Title: %s
  Action: %s
  Result: %s
  Shot Type: %s | Angle: %s | Movement: %s
  Location: %s | Time: %s
  Atmosphere: %s
  Dialogue: %s
  Duration: %ds
  Characters: %v
  Scene ID: %v`,
			shot.StoryboardNumber,
			derefStr(shot.Title),
			derefStr(shot.Action),
			derefStr(shot.Result),
			derefStr(shot.ShotType), derefStr(shot.Angle), derefStr(shot.Movement),
			derefStr(shot.Location), derefStr(shot.Time),
			derefStr(shot.Atmosphere),
			derefStr(shot.Dialogue),
			shot.Duration,
			charIDs,
			shot.SceneID)
		parts = append(parts, desc)
	}
	return strings.Join(parts, "\n\n")
}

// derefStr safely dereferences a string pointer
func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// saveProductionShots saves the merged production shots to the database
func (s *RapidCutService) saveProductionShots(episodeID, dramaID uint, editorialShots []models.Storyboard, rapidCutShots []RapidCutShot) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// First, delete any existing production shots for this episode
		if err := tx.Where("episode_id = ? AND is_production = ?", episodeID, true).
			Delete(&models.Storyboard{}).Error; err != nil {
			return fmt.Errorf("failed to delete existing production shots: %w", err)
		}

		// Build a map from shot_number → Storyboard for quick lookup
		shotMap := make(map[int]*models.Storyboard)
		for i := range editorialShots {
			shotMap[editorialShots[i].StoryboardNumber] = &editorialShots[i]
		}

		pacingMode := "rapid_cut"

		for i, rcShot := range rapidCutShots {
			// Serialize source shot numbers as IDs (lookup actual DB IDs)
			sourceIDs := make([]uint, 0)
			for _, shotNum := range rcShot.SourceShotNumbers {
				if shot, ok := shotMap[shotNum]; ok {
					sourceIDs = append(sourceIDs, shot.ID)
				}
			}
			sourceIDsJSON, _ := json.Marshal(sourceIDs)

			// Build production shot with same struct as editorial
			var titlePtr, locationPtr, timePtr, shotTypePtr, anglePtr, movementPtr *string
			var actionPtr, resultPtr, atmospherePtr, dialoguePtr *string
			var bgmPromptPtr, soundEffectPtr *string

			if rcShot.Title != "" {
				titlePtr = &rcShot.Title
			}
			if rcShot.Location != "" {
				locationPtr = &rcShot.Location
			}
			if rcShot.Time != "" {
				timePtr = &rcShot.Time
			}
			if rcShot.ShotType != "" {
				shotTypePtr = &rcShot.ShotType
			}
			if rcShot.Angle != "" {
				anglePtr = &rcShot.Angle
			}
			if rcShot.Movement != "" {
				movementPtr = &rcShot.Movement
			}
			if rcShot.Action != "" {
				actionPtr = &rcShot.Action
			}
			if rcShot.Result != "" {
				resultPtr = &rcShot.Result
			}
			if rcShot.Atmosphere != "" {
				atmospherePtr = &rcShot.Atmosphere
			}
			if rcShot.Dialogue != "" {
				dialoguePtr = &rcShot.Dialogue
			}
			if rcShot.BgmPrompt != "" {
				bgmPromptPtr = &rcShot.BgmPrompt
			}
			if rcShot.SoundEffect != "" {
				soundEffectPtr = &rcShot.SoundEffect
			}

			// Generate description
			description := fmt.Sprintf("【Rapid Cut Production Shot】\n【Merged from shots】%v\n【Action】%s\n【Result】%s\n【Emotion】%s",
				rcShot.SourceShotNumbers, rcShot.Action, rcShot.Result, rcShot.Emotion)

			// Generate video prompt (same format as storyboard_service)
			// Safely extract SceneID
			var finalSceneID *uint
			if rcShot.SceneID != nil {
				switch v := rcShot.SceneID.(type) {
				case float64:
					if v > 0 {
						val := uint(v)
						finalSceneID = &val
					}
				case int:
					if v > 0 {
						val := uint(v)
						finalSceneID = &val
					}
				case string:
					if parsed, err := strconv.ParseUint(v, 10, 32); err == nil && parsed > 0 {
						val := uint(parsed)
						finalSceneID = &val
					}
				}
			}

			// Generate video prompt (same format as storyboard_service)
			videoPrompt := s.generateRapidCutVideoPrompt(rcShot)

			productionShot := models.Storyboard{
				EpisodeID:        episodeID,
				SceneID:          finalSceneID,
				StoryboardNumber: i + 1,
				Title:            titlePtr,
				Location:         locationPtr,
				Time:             timePtr,
				ShotType:         shotTypePtr,
				Angle:            anglePtr,
				Movement:         movementPtr,
				Action:           actionPtr,
				Result:           resultPtr,
				Atmosphere:       atmospherePtr,
				Dialogue:         dialoguePtr,
				Description:      &description,
				VideoPrompt:      &videoPrompt,
				BgmPrompt:        bgmPromptPtr,
				SoundEffect:      soundEffectPtr,
				Duration:         rcShot.Duration,
				IsProduction:     true,
				PacingMode:       &pacingMode,
				SourceShotIDs:    sourceIDsJSON,
			}

			if err := tx.Create(&productionShot).Error; err != nil {
				s.log.Errorw("Failed to create production shot", "error", err, "shot_number", i+1)
				return err
			}

			// Associate characters (union from all source shots)
			var finalCharIDs []uint
			if len(rcShot.Characters) > 0 {
				for _, c := range rcShot.Characters {
					switch v := c.(type) {
					case float64:
						if v > 0 {
							finalCharIDs = append(finalCharIDs, uint(v))
						}
					case int:
						if v > 0 {
							finalCharIDs = append(finalCharIDs, uint(v))
						}
					case string:
						if parsed, err := strconv.ParseUint(v, 10, 32); err == nil && parsed > 0 {
							finalCharIDs = append(finalCharIDs, uint(parsed))
						}
					}
				}
			}

			if len(finalCharIDs) > 0 {
				var characters []models.Character
				if err := tx.Where("id IN ?", finalCharIDs).Find(&characters).Error; err == nil && len(characters) > 0 {
					if err := tx.Model(&productionShot).Association("Characters").Append(characters); err != nil {
						s.log.Warnw("Failed to associate characters to production shot", "error", err, "shot_number", i+1)
					}
				}
			}

			// Associate props (union from all source shots)
			propIDs := s.collectPropsFromSources(shotMap, rcShot.SourceShotNumbers)
			if len(propIDs) > 0 {
				var props []models.Prop
				if err := tx.Where("id IN ?", propIDs).Find(&props).Error; err == nil && len(props) > 0 {
					if err := tx.Model(&productionShot).Association("Props").Append(props); err != nil {
						s.log.Warnw("Failed to associate props to production shot", "error", err, "shot_number", i+1)
					}
				}
			}

			s.log.Infow("Production shot created",
				"shot_number", i+1,
				"source_shots", rcShot.SourceShotNumbers,
				"character_count", len(rcShot.Characters),
				"duration", rcShot.Duration)
		}

		s.log.Infow("All production shots saved",
			"episode_id", episodeID,
			"count", len(rapidCutShots))
		return nil
	})
}

// collectPropsFromSources collects unique prop IDs from source editorial shots
func (s *RapidCutService) collectPropsFromSources(shotMap map[int]*models.Storyboard, sourceNumbers []int) []uint {
	propIDSet := make(map[uint]bool)
	for _, num := range sourceNumbers {
		if shot, ok := shotMap[num]; ok {
			for _, prop := range shot.Props {
				propIDSet[prop.ID] = true
			}
		}
	}
	result := make([]uint, 0, len(propIDSet))
	for id := range propIDSet {
		result = append(result, id)
	}
	return result
}

// generateRapidCutVideoPrompt generates video prompt for a rapid cut production shot
func (s *RapidCutService) generateRapidCutVideoPrompt(rcShot RapidCutShot) string {
	var parts []string
	videoRatio := "16:9"

	// 1. Multi-beat action (core)
	if rcShot.Action != "" {
		parts = append(parts, fmt.Sprintf("Action sequence: %s", rcShot.Action))
	}

	// 2. Result
	if rcShot.Result != "" {
		parts = append(parts, fmt.Sprintf("Final result: %s", rcShot.Result))
	}

	// 3. Camera transitions
	if rcShot.Movement != "" {
		parts = append(parts, fmt.Sprintf("Camera movement transitions: %s", rcShot.Movement))
	}

	// 4. Shot type transitions
	if rcShot.ShotType != "" {
		parts = append(parts, fmt.Sprintf("Shot type transitions: %s", rcShot.ShotType))
	}
	if rcShot.Angle != "" {
		parts = append(parts, fmt.Sprintf("Camera angle transitions: %s", rcShot.Angle))
	}

	// 5. Scene/environment
	if rcShot.Location != "" {
		locDesc := rcShot.Location
		if rcShot.Time != "" {
			locDesc += ", " + rcShot.Time
		}
		parts = append(parts, fmt.Sprintf("Scene transitions: %s", locDesc))
	}

	// 6. Atmosphere
	if rcShot.Atmosphere != "" {
		parts = append(parts, fmt.Sprintf("Atmosphere evolution: %s", rcShot.Atmosphere))
	}

	// 7. Sound effects
	if rcShot.SoundEffect != "" {
		parts = append(parts, fmt.Sprintf("Sound effects: %s", rcShot.SoundEffect))
	}

	// 8. Video ratio
	parts = append(parts, fmt.Sprintf("=VideoRatio: %s", videoRatio))

	if len(parts) > 0 {
		return strings.Join(parts, ". ")
	}
	return "Cinematic rapid cut video sequence"
}

// DeleteRapidCut removes all production shots for an episode
func (s *RapidCutService) DeleteRapidCut(episodeID string) error {
	epID, err := strconv.ParseUint(episodeID, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid episode ID: %s", episodeID)
	}

	// Get production shot IDs first (for cleaning up associations)
	var productionShotIDs []uint
	if err := s.db.Model(&models.Storyboard{}).
		Where("episode_id = ? AND is_production = ?", uint(epID), true).
		Pluck("id", &productionShotIDs).Error; err != nil {
		return fmt.Errorf("failed to get production shots: %w", err)
	}

	if len(productionShotIDs) == 0 {
		return nil // Nothing to delete
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// Clean up image_generations references
		if err := tx.Model(&models.ImageGeneration{}).
			Where("storyboard_id IN ?", productionShotIDs).
			Update("storyboard_id", nil).Error; err != nil {
			s.log.Warnw("Failed to clean image generation references", "error", err)
		}

		// Delete production shots
		result := tx.Where("episode_id = ? AND is_production = ?", uint(epID), true).
			Delete(&models.Storyboard{})
		if result.Error != nil {
			return fmt.Errorf("failed to delete production shots: %w", result.Error)
		}

		s.log.Infow("Production shots deleted",
			"episode_id", episodeID,
			"deleted_count", result.RowsAffected)
		return nil
	})
}

// HasRapidCut checks if an episode has production shots
func (s *RapidCutService) HasRapidCut(episodeID string) (bool, error) {
	epID, err := strconv.ParseUint(episodeID, 10, 32)
	if err != nil {
		return false, fmt.Errorf("invalid episode ID")
	}

	var count int64
	if err := s.db.Model(&models.Storyboard{}).
		Where("episode_id = ? AND is_production = ?", uint(epID), true).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
