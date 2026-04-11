package services

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	models "github.com/drama-generator/backend/domain/models"
	"github.com/drama-generator/backend/application/prompts"
	"github.com/drama-generator/backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ============================================================================
// Nursery Rhyme Mode — Structs
// ============================================================================

// LyricsBlock represents a parsed segment from timestamped lyrics input
type LyricsBlock struct {
	BlockID        int    // Sequential ID
	StartTimeSec   int    // Start time in seconds
	EndTimeSec     int    // End time in seconds
	DurationSec    int    // Calculated duration
	LyricsText     string // The lyrics text
	SectionType    string // verse / chorus / bridge / intro / outro / instrumental
	SectionNumber  int    // Which verse/chorus (1, 2, 3...)
	VerseSubject   string // Subject from section header ("[VERSE 1: The Wheels]" → "The Wheels")
	IsInstrumental bool   // True if tagged [INSTRUMENTAL]
}

// NurseryRhymeShot is the AI output struct for nursery_rhyme mode
type NurseryRhymeShot struct {
	ShotID            int    `json:"shot_id"`
	LyricsBlockID     int    `json:"lyrics_block_id"`
	LyricsText        string `json:"lyrics_text"`
	TimestampStart    string `json:"timestamp_start"`
	TimestampEnd      string `json:"timestamp_end"`
	DurationSec       int    `json:"duration_sec"`
	SectionType       string `json:"section_type"`
	SectionNumber     int    `json:"section_number"`
	VerseSubject      string `json:"verse_subject"`
	ShotRole          string `json:"shot_role"`           // establishing / reveal / detail / group_payoff
	IsCallback        bool   `json:"is_callback"`
	CallbackToShotID  *int   `json:"callback_to_shot_id"`
	VisualDescription string `json:"visual_description"`
	Title             string `json:"title"`
	ShotType          string `json:"shot_type"`           // ELS / LS / MS / CU / ECU
	CameraMovement    string `json:"camera_movement"`
	AnimationHint     string `json:"animation_hint"`
	OverlayText       string `json:"overlay_text"`
	Location          string `json:"location"`
	Atmosphere        string `json:"atmosphere"`
	BgmPrompt         string `json:"bgm_prompt"`
	SoundEffect       string `json:"sound_effect"`
	TransitionIn      string `json:"transition_in"`
	Characters        []uint `json:"characters"`
	Props             []uint `json:"props"`
	SceneID           *uint  `json:"scene_id"`
}

// ============================================================================
// Lyrics Parser
// ============================================================================

// Regex patterns for lyrics parsing
var (
	// Section header: [VERSE 1: The Wheels] or [CHORUS] or [BRIDGE] etc.
	nurseryHeaderPattern = regexp.MustCompile(`(?i)^\[(?:(VERSE|CHORUS|BRIDGE|INTRO|OUTRO)\s*(\d*))\s*(?::\s*(.+?))?\]$`)
	// Timestamp: (0:05 – 0:11) or (0:05 - 0:11) or (00:05 – 00:11) or (1:05 – 1:11)
	nurseryTimestampPattern = regexp.MustCompile(`^\((\d{1,2}:\d{2})\s*[–\-]\s*(\d{1,2}:\d{2})\)\s*(.*)$`)
	// Instrumental tag: [INSTRUMENTAL]
	instrumentalTag = regexp.MustCompile(`(?i)\[INSTRUMENTAL\]`)
)

// parseTimestampToSeconds converts "M:SS" or "MM:SS" to total seconds
func parseTimestampToSeconds(ts string) (int, error) {
	parts := strings.Split(ts, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid timestamp format: %s", ts)
	}
	minutes, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid minutes in timestamp %s: %w", ts, err)
	}
	seconds, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("invalid seconds in timestamp %s: %w", ts, err)
	}
	return minutes*60 + seconds, nil
}

// parseLyricsInput parses timestamped lyrics format into LyricsBlock array
// Input format:
//
//	[VERSE 1: The Wheels]
//	(0:05 – 0:11) [INSTRUMENTAL] Bus driving establishing shot
//	(0:12 – 0:15) The wheels on the bus go round and round, round and round
func parseLyricsInput(script string) ([]LyricsBlock, error) {
	lines := strings.Split(strings.ReplaceAll(script, "\r\n", "\n"), "\n")

	var blocks []LyricsBlock
	blockID := 1

	// Track current section context
	currentSectionType := "verse"
	currentSectionNum := 1
	currentSubject := ""

	for _, rawLine := range lines {
		line := strings.TrimSpace(rawLine)
		if line == "" {
			continue
		}

		// Check for section header: [VERSE 1: The Wheels]
		if headerMatch := nurseryHeaderPattern.FindStringSubmatch(line); headerMatch != nil {
			sectionType := strings.ToLower(headerMatch[1])
			sectionNum := 1
			if headerMatch[2] != "" {
				if n, err := strconv.Atoi(headerMatch[2]); err == nil {
					sectionNum = n
				}
			}
			subject := strings.TrimSpace(headerMatch[3])

			currentSectionType = sectionType
			currentSectionNum = sectionNum
			currentSubject = subject
			continue
		}

		// Check for timestamp line: (0:05 – 0:11) lyrics text
		if tsMatch := nurseryTimestampPattern.FindStringSubmatch(line); tsMatch != nil {
			startTS := tsMatch[1]
			endTS := tsMatch[2]
			content := strings.TrimSpace(tsMatch[3])

			startSec, err := parseTimestampToSeconds(startTS)
			if err != nil {
				continue // Skip malformed timestamps
			}
			endSec, err := parseTimestampToSeconds(endTS)
			if err != nil {
				continue
			}

			// Detect [INSTRUMENTAL] tag
			isInstrumental := instrumentalTag.MatchString(content)
			if isInstrumental {
				content = strings.TrimSpace(instrumentalTag.ReplaceAllString(content, ""))
			}

			// Determine section type for this block
			sectionType := currentSectionType
			if isInstrumental && sectionType != "intro" && sectionType != "outro" {
				sectionType = "instrumental"
			}

			block := LyricsBlock{
				BlockID:        blockID,
				StartTimeSec:   startSec,
				EndTimeSec:     endSec,
				DurationSec:    endSec - startSec,
				LyricsText:     content,
				SectionType:    sectionType,
				SectionNumber:  currentSectionNum,
				VerseSubject:   currentSubject,
				IsInstrumental: isInstrumental,
			}

			blocks = append(blocks, block)
			blockID++
		}
	}

	if len(blocks) == 0 {
		return nil, fmt.Errorf("no valid timestamped lyrics blocks found in input")
	}

	return blocks, nil
}

// buildLyricsAnalysis creates a structured text analysis from parsed lyrics blocks
// This is injected into the AI prompt between the lyrics and format instructions
func buildLyricsAnalysis(blocks []LyricsBlock, structureType, structureReason string) string {
	var sb strings.Builder

	sb.WriteString("=== LYRICS STRUCTURE ANALYSIS ===\n")
	sb.WriteString(fmt.Sprintf("Total lyrics blocks: %d\n", len(blocks)))
	sb.WriteString(fmt.Sprintf("Structure type: %s\n", structureType))
	sb.WriteString(fmt.Sprintf("Reason: %s\n", structureReason))

	totalDuration := 0
	for _, b := range blocks {
		totalDuration += b.DurationSec
	}
	sb.WriteString(fmt.Sprintf("Total duration: %d seconds\n\n", totalDuration))

	sb.WriteString("=== BLOCK-BY-BLOCK BREAKDOWN ===\n")
	for _, b := range blocks {
		tag := b.SectionType
		if b.VerseSubject != "" {
			tag += ": " + b.VerseSubject
		}
		instrumental := ""
		if b.IsInstrumental {
			instrumental = " [INSTRUMENTAL]"
		}
		sb.WriteString(fmt.Sprintf("Block %d [%s]%s (%d:%02d – %d:%02d, %ds): %s\n",
			b.BlockID, tag, instrumental,
			b.StartTimeSec/60, b.StartTimeSec%60,
			b.EndTimeSec/60, b.EndTimeSec%60,
			b.DurationSec, b.LyricsText))
	}

	sb.WriteString("\n=== RULES ===\n")
	sb.WriteString("1. Create shots that FIT WITHIN each block's timestamp range\n")
	sb.WriteString("2. You MUST split one block into multiple shots if block duration is > 5s (max 5s per shot)\n")
	sb.WriteString("3. If you split a block, split its lyrics_text text proportionally among the shots\n")
	sb.WriteString("4. You MAY merge very short adjacent blocks into one shot (if total < 3s)\n")
	sb.WriteString("5. Every shot MUST have lyrics_block_id referencing the source block\n")
	sb.WriteString("6. Total shot durations MUST sum to match total lyrics duration\n")
	sb.WriteString("7. [INSTRUMENTAL] blocks MUST have visual description (establishing shots, transitions)\n")

	return sb.String()
}

// detectNurseryStructure detects Narrative vs Cumulative structure type
func detectNurseryStructure(blocks []LyricsBlock) (structureType string, reason string) {
	// Collect unique words per verse
	verseWords := make(map[int]map[string]bool)
	for _, b := range blocks {
		if b.SectionType != "verse" && b.SectionType != "chorus" {
			continue
		}
		if _, ok := verseWords[b.SectionNumber]; !ok {
			verseWords[b.SectionNumber] = make(map[string]bool)
		}
		words := strings.Fields(strings.ToLower(b.LyricsText))
		for _, w := range words {
			if len(w) > 3 { // Skip short words like "the", "and", "a"
				verseWords[b.SectionNumber][w] = true
			}
		}
	}

	if len(verseWords) < 2 {
		return "narrative", "Only one verse detected — defaulting to narrative structure"
	}

	// Check if later verses contain significant words from earlier verses
	// (cumulative songs repeat earlier subjects' sounds/words)
	overlapCount := 0
	totalComparisons := 0

	verseNums := make([]int, 0, len(verseWords))
	for n := range verseWords {
		verseNums = append(verseNums, n)
	}

	for i := 1; i < len(verseNums); i++ {
		laterVerse := verseWords[verseNums[i]]
		for j := 0; j < i; j++ {
			earlierVerse := verseWords[verseNums[j]]
			for word := range earlierVerse {
				totalComparisons++
				if laterVerse[word] {
					overlapCount++
				}
			}
		}
	}

	if totalComparisons > 0 {
		overlapRatio := float64(overlapCount) / float64(totalComparisons)
		if overlapRatio > 0.3 {
			return "cumulative",
				fmt.Sprintf("Later verses share %.0f%% words with earlier verses — indicating cumulative/additive structure (like Old MacDonald)", overlapRatio*100)
		}
	}

	return "narrative",
		"Each verse contains mostly unique content — indicating narrative structure (like Wheels on the Bus)"
}

// detectNurseryRhymeInput checks if script content matches nursery_rhyme format
// Used by auto-detect mode to suggest nursery_rhyme split mode
func detectNurseryRhymeInput(script string) bool {
	lines := strings.Split(script, "\n")
	timestampCount := 0
	headerCount := 0

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if nurseryTimestampPattern.MatchString(line) {
			timestampCount++
		}
		if nurseryHeaderPattern.MatchString(line) {
			headerCount++
		}
	}

	// Require >= 3 timestamp lines AND >= 1 section header
	return timestampCount >= 3 && headerCount >= 1
}

// ============================================================================
// Processing Function
// ============================================================================

// processNurseryRhymeGeneration handles the nursery_rhyme split mode
func (s *StoryboardService) processNurseryRhymeGeneration(
	taskID, episodeID string, dramaID uint,
	model, scriptContent, characterList, sceneList, propList string,
) {
	// 1. Parse lyrics input
	blocks, parseErr := parseLyricsInput(scriptContent)
	if parseErr != nil {
		s.log.Errorw("Failed to parse nursery rhyme lyrics", "error", parseErr, "task_id", taskID)
		if updateErr := s.taskService.UpdateTaskError(taskID, fmt.Errorf("Lyrics parsing failed: %w", parseErr)); updateErr != nil {
			s.log.Errorw("Failed to update task error", "error", updateErr, "task_id", taskID)
		}
		return
	}

	s.log.Infow("Parsed nursery rhyme lyrics",
		"task_id", taskID, "episode_id", episodeID,
		"block_count", len(blocks))

	if err := s.taskService.UpdateTaskStatus(taskID, "processing", 15, "Lyrics parsed, analyzing structure..."); err != nil {
		s.log.Errorw("Failed to update task status", "error", err, "task_id", taskID)
	}

	// 2. Detect structure type
	structureType, structureReason := detectNurseryStructure(blocks)
	s.log.Infow("Detected nursery rhyme structure",
		"task_id", taskID, "structure_type", structureType, "reason", structureReason)

	// 3. Build analysis context
	lyricsAnalysis := buildLyricsAnalysis(blocks, structureType, structureReason)

	// 4. Load system prompt
	systemPrompt := s.promptI18n.WithDramaNurseryRhymeSystemPrompt(dramaID)

	// 5. Build full prompt
	scriptLabel := s.promptI18n.FormatUserPrompt("script_content_label")
	charListLabel := s.promptI18n.FormatUserPrompt("character_list_label")
	charConstraint := s.promptI18n.FormatUserPrompt("character_constraint")
	sceneListLabel := s.promptI18n.FormatUserPrompt("scene_list_label")
	sceneConstraint := s.promptI18n.FormatUserPrompt("scene_constraint")
	propListLabel := s.promptI18n.FormatUserPrompt("prop_list_label")
	propConstraint := s.promptI18n.FormatUserPrompt("prop_constraint")
	formatInstructions := prompts.Get("storyboard_nursery_rhyme_format.txt")

	if err := s.taskService.UpdateTaskStatus(taskID, "processing", 20, "Calling AI for shot planning..."); err != nil {
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
		s.log.Errorw("Failed to generate nursery rhyme storyboard", "error", err, "task_id", taskID)
		if updateErr := s.taskService.UpdateTaskError(taskID, fmt.Errorf("Nursery rhyme generation failed: %w", err)); updateErr != nil {
			s.log.Errorw("Failed to update task error", "error", updateErr, "task_id", taskID)
		}
		return
	}

	if err := s.taskService.UpdateTaskStatus(taskID, "processing", 50, "AI response received, parsing shots..."); err != nil {
		s.log.Errorw("Failed to update task status", "error", err, "task_id", taskID)
	}

	// 7. Parse JSON → NurseryRhymeShot[]
	var shots []NurseryRhymeShot
	if err := utils.SafeParseAIJSON(text, &shots); err != nil {
		s.log.Errorw("Failed to parse nursery rhyme JSON", "error", err, "response", text[:min(500, len(text))], "task_id", taskID)
		if updateErr := s.taskService.UpdateTaskError(taskID, fmt.Errorf("Failed to parse nursery rhyme result: %w", err)); updateErr != nil {
			s.log.Errorw("Failed to update task error", "error", updateErr, "task_id", taskID)
		}
		return
	}

	// Re-number shots sequentially
	for i := range shots {
		shots[i].ShotID = i + 1
	}

	// 8. Post-process: force-correct durations from parsed timestamps
	for i := range shots {
		if shots[i].DurationSec <= 0 {
			// Try to calculate from timestamps
			startSec, err1 := parseTimestampToSeconds(shots[i].TimestampStart)
			endSec, err2 := parseTimestampToSeconds(shots[i].TimestampEnd)
			if err1 == nil && err2 == nil && endSec > startSec {
				shots[i].DurationSec = endSec - startSec
			} else {
				shots[i].DurationSec = 3 // Default fallback
			}
		}
	}

	if err := s.taskService.UpdateTaskStatus(taskID, "processing", 70, "Saving nursery rhyme shots..."); err != nil {
		s.log.Errorw("Failed to update task status", "error", err, "task_id", taskID)
	}

	// Calculate total duration
	totalDuration := 0
	for _, shot := range shots {
		totalDuration += shot.DurationSec
	}

	s.log.Infow("Nursery rhyme storyboard generated",
		"task_id", taskID,
		"episode_id", episodeID,
		"count", len(shots),
		"total_duration_seconds", totalDuration,
		"structure_type", structureType)

	// 9. Save to database
	if err := s.saveNurseryRhymeShots(episodeID, dramaID, shots); err != nil {
		s.log.Errorw("Failed to save nursery rhyme shots", "error", err, "task_id", taskID)
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
		"mode":             "nursery_rhyme",
		"structure_type":   structureType,
		"structure_reason": structureReason,
	}

	if err := s.taskService.UpdateTaskResult(taskID, resultData); err != nil {
		s.log.Errorw("Failed to update task result", "error", err, "task_id", taskID)
		return
	}

	s.log.Infow("Nursery rhyme storyboard generation completed", "task_id", taskID, "episode_id", episodeID)
}

// ============================================================================
// Save Function
// ============================================================================

// saveNurseryRhymeShots maps NurseryRhymeShot[] to models.Storyboard and saves to DB
// Follows the same pattern as saveVoiceoverShots()
func (s *StoryboardService) saveNurseryRhymeShots(episodeID string, dramaID uint, shots []NurseryRhymeShot) error {
	epID, err := strconv.ParseUint(episodeID, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid episode ID: %s", episodeID)
	}

	if len(shots) == 0 {
		return fmt.Errorf("AI returned 0 shots, refusing to save")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// Verify episode exists
		var episode models.Episode
		if err := tx.First(&episode, epID).Error; err != nil {
			return fmt.Errorf("episode not found: %s", episodeID)
		}

		// Clear existing image generation references
		var storyboardIDs []uint
		if err := tx.Model(&models.Storyboard{}).
			Where("episode_id = ?", uint(epID)).
			Pluck("id", &storyboardIDs).Error; err != nil {
			return err
		}

		if len(storyboardIDs) > 0 {
			if err := tx.Model(&models.ImageGeneration{}).
				Where("storyboard_id IN ?", storyboardIDs).
				Update("storyboard_id", nil).Error; err != nil {
				return err
			}
		}

		// Delete existing storyboards for this episode
		if result := tx.Where("episode_id = ?", uint(epID)).Delete(&models.Storyboard{}); result.Error != nil {
			return result.Error
		}

		// Batch-load all character descriptions for prompt generation
		allCharIDs := make(map[uint]bool)
		for _, shot := range shots {
			for _, cid := range shot.Characters {
				allCharIDs[cid] = true
			}
		}
		charDescMap := make(map[uint]string)
		if len(allCharIDs) > 0 {
			var ids []uint
			for id := range allCharIDs {
				ids = append(ids, id)
			}
			var chars []models.Character
			if err := tx.Where("id IN ?", ids).Find(&chars).Error; err == nil {
				for _, c := range chars {
					desc := c.Name
					if c.Appearance != nil && *c.Appearance != "" {
						desc += " (" + *c.Appearance + ")"
					} else if c.Description != nil && *c.Description != "" {
						desc += " (" + *c.Description + ")"
					}
					charDescMap[c.ID] = desc
				}
			}
		}

		// Helper functions
		strPtr := func(v string) *string {
			if v == "" {
				return nil
			}
			return &v
		}
		boolPtr := func(v bool) *bool {
			return &v
		}
		intPtr := func(v int) *int {
			return &v
		}

		// Save each nursery rhyme shot as a Storyboard
		for _, shot := range shots {
			// Build character descriptions for image prompt
			var charDescParts []string
			for _, cid := range shot.Characters {
				if d, ok := charDescMap[cid]; ok {
					charDescParts = append(charDescParts, d)
				}
			}
			charDescs := strings.Join(charDescParts, "; ")

			// Load props for prompt generation
			var propDescriptions string
			var loadedProps []models.Prop
			if len(shot.Props) > 0 {
				if err := tx.Where("id IN ?", shot.Props).Find(&loadedProps).Error; err == nil {
					var names []string
					for _, p := range loadedProps {
						desc := p.Name
						if p.Prompt != nil && *p.Prompt != "" {
							desc += " (" + *p.Prompt + ")"
						}
						names = append(names, desc)
					}
					propDescriptions = strings.Join(names, ", ")
				}
			}

			// Generate image prompt using existing helper
			sbForPrompt := Storyboard{
				ShotNumber:     shot.ShotID,
				Title:          shot.Title,
				ShotType:       shot.ShotType,
				Movement:       shot.CameraMovement,
				Location:       shot.Location,
				Atmosphere:     shot.Atmosphere,
				Action:         shot.VisualDescription,
				Duration:       shot.DurationSec,
				CharacterDescs: charDescs,
			}

			imagePrompt := s.generateImagePrompt(sbForPrompt, propDescriptions)
			videoPrompt := s.generateVideoPrompt(sbForPrompt)

			// Build the description field
			description := shot.VisualDescription

			// Set AudioMode to lyrics_sync for nursery rhyme mode
			audioMode := "lyrics_sync"

			// Map shot role — reuse existing ShotRole field
			shotRole := shot.ShotRole
			if shotRole == "" {
				shotRole = "establishing"
			}

			storyboard := models.Storyboard{
				EpisodeID:        uint(epID),
				SceneID:          shot.SceneID,
				StoryboardNumber: shot.ShotID,
				Title:            strPtr(shot.Title),
				Location:         strPtr(shot.Location),
				ShotType:         strPtr(shot.ShotType),
				Movement:         strPtr(shot.CameraMovement),
				Description:      &description,
				Action:           strPtr(shot.VisualDescription),
				Atmosphere:       strPtr(shot.Atmosphere),
				ImagePrompt:      &imagePrompt,
				VideoPrompt:      &videoPrompt,
				VideoPromptSource: "auto",
				BgmPrompt:        strPtr(shot.BgmPrompt),
				SoundEffect:      strPtr(shot.SoundEffect),
				Duration:         shot.DurationSec,
				// Reuse existing Voice-over fields for compatibility
				ScriptSegment:   strPtr(shot.LyricsText), // Lyrics appear in Narrator column
				AudioMode:       &audioMode,
				ShotRole:        &shotRole,
				VisualType:      strPtr("literal"),
				// Nursery Rhyme specific fields
				LyricsText:      strPtr(shot.LyricsText),
				SectionType:     strPtr(shot.SectionType),
				VerseSubject:    strPtr(shot.VerseSubject),
				OverlayText:     strPtr(shot.OverlayText),
				AnimationHint:   strPtr(shot.AnimationHint),
				IsCallback:      boolPtr(shot.IsCallback),
			}

			// Set callback shot number if applicable
			if shot.CallbackToShotID != nil {
				storyboard.CallbackShotNum = intPtr(*shot.CallbackToShotID)
			}

			if err := tx.Create(&storyboard).Error; err != nil {
				s.log.Errorw("Failed to create nursery rhyme shot", "error", err, "shot_id", shot.ShotID)
				return err
			}

			// Associate characters
			if len(shot.Characters) > 0 {
				var characters []models.Character
				if err := tx.Where("id IN ?", shot.Characters).Find(&characters).Error; err == nil && len(characters) > 0 {
					if err := tx.Model(&storyboard).Association("Characters").Append(characters); err != nil {
						s.log.Warnw("Failed to associate characters", "error", err, "shot_id", shot.ShotID)
					}
				}
			}

			// Associate props
			if len(loadedProps) > 0 {
				if err := tx.Model(&storyboard).Association("Props").Append(loadedProps); err != nil {
					s.log.Warnw("Failed to associate props", "error", err, "shot_id", shot.ShotID)
				}
			}
		}

		s.log.Infow("Nursery rhyme shots saved successfully", "episode_id", episodeID, "count", len(shots))
		return nil
	})
}
