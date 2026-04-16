package services

import (
	"fmt"
	"strings"

	"github.com/drama-generator/backend/application/prompts"
	"github.com/drama-generator/backend/domain/models"
	"github.com/drama-generator/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

// ============================================================================
// MV Maker Mode — Genre Profile System
// ============================================================================

// MVGenrePromptMap maps genre profile keys to default prompt files
var MVGenrePromptMap = map[string]string{
	"gaming_horror": "storyboard_mv_gaming_horror.txt",
	// Future genres:
	// "gaming_parody": "storyboard_mv_gaming_parody.txt",
	// "general":       "storyboard_mv_general.txt",
}

// processMVMakerGeneration handles the mv_maker split mode
// It reuses the nursery rhyme infrastructure (parser, struct, save) with genre-specific prompts
func (s *StoryboardService) processMVMakerGeneration(
	taskID, episodeID string, dramaID uint,
	model, genreProfile, scriptContent, characterList, sceneList, propList string,
) {
	// 1. Parse lyrics input — REUSE from nursery
	blocks, parseErr := parseLyricsInput(scriptContent)
	if parseErr != nil {
		s.log.Errorw("Failed to parse MV lyrics", "error", parseErr, "task_id", taskID)
		if updateErr := s.taskService.UpdateTaskError(taskID, fmt.Errorf("Lyrics parsing failed: %w", parseErr)); updateErr != nil {
			s.log.Errorw("Failed to update task error", "error", updateErr, "task_id", taskID)
		}
		return
	}

	s.log.Infow("Parsed MV lyrics",
		"task_id", taskID, "episode_id", episodeID,
		"block_count", len(blocks), "genre", genreProfile)

	if err := s.taskService.UpdateTaskStatus(taskID, "processing", 15, "Lyrics parsed, analyzing song structure..."); err != nil {
		s.log.Errorw("Failed to update task status", "error", err, "task_id", taskID)
	}

	// 2. Detect song structure — ENHANCED for MV genres (energy curve aware)
	structureType, structureReason := detectMVStructure(blocks)
	s.log.Infow("Detected MV song structure",
		"task_id", taskID, "structure_type", structureType, "reason", structureReason)

	// 3. Build analysis context — REUSE from nursery
	lyricsAnalysis := buildLyricsAnalysis(blocks, structureType, structureReason)

	// 4. Load genre-specific system prompt
	systemPrompt := s.loadMVGenrePrompt(dramaID, genreProfile)

	// 5. Build full prompt
	scriptLabel := s.promptI18n.FormatUserPrompt("script_content_label")
	charListLabel := s.promptI18n.FormatUserPrompt("character_list_label")
	charConstraint := s.promptI18n.FormatUserPrompt("character_constraint")
	sceneListLabel := s.promptI18n.FormatUserPrompt("scene_list_label")
	sceneConstraint := s.promptI18n.FormatUserPrompt("scene_constraint")
	propListLabel := s.promptI18n.FormatUserPrompt("prop_list_label")
	propConstraint := s.promptI18n.FormatUserPrompt("prop_constraint")
	formatInstructions := prompts.Get("storyboard_nursery_rhyme_format.txt") // REUSE format spec

	if err := s.taskService.UpdateTaskStatus(taskID, "processing", 20, "Calling AI for MV shot planning..."); err != nil {
		s.log.Errorw("Failed to update task status", "error", err, "task_id", taskID)
	}

	prompt := fmt.Sprintf(`%s

%s
%s

%s

%s
%s
%s

%s
%s
%s

%s
%s
%s

%s`,
		systemPrompt,
		scriptLabel, scriptContent,
		lyricsAnalysis,
		charListLabel, characterList, charConstraint,
		sceneListLabel, sceneList, sceneConstraint,
		propListLabel, propList, propConstraint,
		formatInstructions)

	// 6. Call AI
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
		s.log.Errorw("Failed to generate MV storyboard", "error", err, "task_id", taskID)
		if updateErr := s.taskService.UpdateTaskError(taskID, fmt.Errorf("MV generation failed: %w", err)); updateErr != nil {
			s.log.Errorw("Failed to update task error", "error", updateErr, "task_id", taskID)
		}
		return
	}

	if err := s.taskService.UpdateTaskStatus(taskID, "processing", 50, "AI response received, parsing shots..."); err != nil {
		s.log.Errorw("Failed to update task status", "error", err, "task_id", taskID)
	}

	// 7. Parse JSON → NurseryRhymeShot[] (reuse struct — fields are genre-neutral)
	var shots []NurseryRhymeShot
	if err := utils.SafeParseAIJSON(text, &shots); err != nil {
		s.log.Errorw("Failed to parse MV JSON", "error", err, "response", text[:min(500, len(text))], "task_id", taskID)
		if updateErr := s.taskService.UpdateTaskError(taskID, fmt.Errorf("Failed to parse MV result: %w", err)); updateErr != nil {
			s.log.Errorw("Failed to update task error", "error", updateErr, "task_id", taskID)
		}
		return
	}

	// Re-number shots sequentially
	for i := range shots {
		shots[i].ShotID = i + 1
	}

	// 8. Post-process: force-correct durations from timestamps
	for i := range shots {
		if shots[i].DurationSec <= 0 {
			startSec, err1 := parseTimestampToSeconds(shots[i].TimestampStart)
			endSec, err2 := parseTimestampToSeconds(shots[i].TimestampEnd)
			if err1 == nil && err2 == nil && endSec > startSec {
				shots[i].DurationSec = endSec - startSec
			} else {
				shots[i].DurationSec = 3 // Default fallback
			}
		}
	}

	if err := s.taskService.UpdateTaskStatus(taskID, "processing", 70, "Saving MV shots..."); err != nil {
		s.log.Errorw("Failed to update task status", "error", err, "task_id", taskID)
	}

	// Calculate total duration
	totalDuration := 0
	for _, shot := range shots {
		totalDuration += shot.DurationSec
	}

	s.log.Infow("MV storyboard generated",
		"task_id", taskID,
		"episode_id", episodeID,
		"count", len(shots),
		"total_duration_seconds", totalDuration,
		"genre_profile", genreProfile,
		"structure_type", structureType)

	// 9. Save to database — REUSE saveNurseryRhymeShots()
	if err := s.saveNurseryRhymeShots(episodeID, dramaID, shots); err != nil {
		s.log.Errorw("Failed to save MV shots", "error", err, "task_id", taskID)
		if updateErr := s.taskService.UpdateTaskError(taskID, fmt.Errorf("Failed to save shots: %w", err)); updateErr != nil {
			s.log.Errorw("Failed to update task error", "error", updateErr, "task_id", taskID)
		}
		return
	}

	if err := s.taskService.UpdateTaskStatus(taskID, "processing", 90, "Updating episode duration..."); err != nil {
		s.log.Errorw("Failed to update task status", "error", err, "task_id", taskID)
	}

	// Update episode duration
	durationMinutes := (totalDuration + 59) / 60
	if err := s.db.Model(&models.Episode{}).Where("id = ?", episodeID).Update("duration", durationMinutes).Error; err != nil {
		s.log.Errorw("Failed to update episode duration", "error", err, "task_id", taskID)
	}

	// Update task result
	resultData := gin.H{
		"storyboards":      shots,
		"total":            len(shots),
		"total_duration":   totalDuration,
		"duration_minutes": durationMinutes,
		"mode":             "mv_maker",
		"genre_profile":    genreProfile,
		"structure_type":   structureType,
		"structure_reason": structureReason,
	}

	if err := s.taskService.UpdateTaskResult(taskID, resultData); err != nil {
		s.log.Errorw("Failed to update task result", "error", err, "task_id", taskID)
		return
	}

	s.log.Infow("MV storyboard generation completed", "task_id", taskID, "episode_id", episodeID)
}

// ============================================================================
// MV Structure Detection
// ============================================================================

// detectMVStructure detects song structure type for MV genres
// Supports energy-curve aware classification:
//   - "dynamic_intensity": Has PRE-CHORUS or DROP → progressive energy build-release
//   - "narrative_arc": Has BRIDGE → tension-resolution storytelling
//   - "standard": Simple verse/chorus alternation
func detectMVStructure(blocks []LyricsBlock) (structureType string, reason string) {
	hasPreChorus := false
	hasDrop := false
	hasBridge := false
	hasBreakdown := false

	for _, b := range blocks {
		st := strings.ToLower(b.SectionType)
		switch st {
		case "pre_chorus", "pre-chorus":
			hasPreChorus = true
		case "drop":
			hasDrop = true
		case "breakdown":
			hasBreakdown = true
		case "bridge":
			hasBridge = true
		}
	}

	if hasDrop || hasBreakdown {
		return "dynamic_intensity",
			"Song has DROP/BREAKDOWN sections — use progressive energy curve: build → explode → valley → final peak"
	}
	if hasPreChorus {
		return "dynamic_intensity",
			"Song has PRE-CHORUS sections — use gradual tension build into chorus explosions"
	}
	if hasBridge {
		return "narrative_arc",
			"Song has BRIDGE section — use tension/resolution arc with dramatic shift at bridge"
	}
	return "standard",
		"Standard verse/chorus structure — use alternating energy levels between verse (low) and chorus (high)"
}

// ============================================================================
// Genre Prompt Loader
// ============================================================================

// loadMVGenrePrompt resolves the system prompt for a specific MV genre
func (s *StoryboardService) loadMVGenrePrompt(dramaID uint, genreProfile string) string {
	promptKey := "mv_maker_" + genreProfile
	return s.promptI18n.WithDramaMVMakerSystemPrompt(dramaID, promptKey)
}
