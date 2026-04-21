package services

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/drama-generator/backend/application/prompts"
	models "github.com/drama-generator/backend/domain/models"
	"github.com/drama-generator/backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ============================================================================
// Narrative MV Mode — Story Bible Structs
// ============================================================================

// StoryBible is the parsed representation of a [STORY_BIBLE] structured input
type StoryBible struct {
	WorldDescription string
	EmotionalCore    string
	Characters       []NarrativeCharacter
	PrologueDesc     string
	PrologueDuration int // seconds
	MusicSegments    []NarrativeMusicSegment
	LyricsBlocks     []LyricsBlock // reused from storyboard_nursery_service.go
	EpilogueDesc     string
	EpilogueDuration int // seconds
	TotalMusicSec    int // computed: sum of all MusicSegment durations
}

// NarrativeCharacter represents a character entry from [CHARACTERS] section
type NarrativeCharacter struct {
	Name        string
	Description string
	Role        string
}

// NarrativeMusicSegment represents a music segment from [MUSIC_SEGMENTS] section
type NarrativeMusicSegment struct {
	StartSec   int
	EndSec     int
	Name       string // "INTRO", "VERSE 1", "CHORUS"...
	Emotion    string // text after "emotion:"
	SyncPoints []NarrativeSyncPoint
}

// NarrativeSyncPoint is an indented [SYNC_POINT] line under a MUSIC_SEGMENT
type NarrativeSyncPoint struct {
	Description string // text before " — "
	SyncType    string // "convergent" | "irony" | "parallel"
}

// NarrativePlan is the output of Phase 1 AI call (story planning)
type NarrativePlan struct {
	NarrativeThread []struct {
		Motif     string   `json:"motif"`
		AppearsIn []string `json:"appears_in"`
		Meaning   string   `json:"meaning"`
	} `json:"narrative_thread"`
	LightingArc    string            `json:"lighting_arc"`
	CharacterArcs  map[string]string `json:"character_arcs"`
	SyncPoints     []struct {
		MusicTimestamp string `json:"music_timestamp"`
		Segment        string `json:"segment"`
		SyncType       string `json:"sync_type"`
		Description    string `json:"description"`
	} `json:"sync_points"`
	TotalMusicDurationSec int `json:"total_music_duration_sec"`
}

// NarrativeShot is the AI output struct for narrative_mv mode (Phase 2)
type NarrativeShot struct {
	ShotID            int    `json:"shot_id"`
	NarrativePart     string `json:"narrative_part"`      // "prologue" | "music_film" | "epilogue"
	HasMusic          bool   `json:"has_music"`
	TimestampStart    string `json:"timestamp_start"`     // assigned by backend post-processing, not AI
	TimestampEnd      string `json:"timestamp_end"`       // assigned by backend post-processing
	DurationSec       int    `json:"duration_sec"`        // AI output — the only time field AI generates
	MusicSegment      string `json:"music_segment"`       // "VERSE 1", "CHORUS", "" for no-music parts
	MusicSyncType     string `json:"music_sync_type"`     // "parallel" | "convergent" | "irony" | ""
	LyricsAnchor      string `json:"lyrics_anchor"`       // lyric line this shot is anchored to
	VisualDescription string `json:"visual_description"`
	ActingNote        string `json:"acting_note"`
	ShotType          string `json:"shot_type"`
	CameraAngle       string `json:"camera_angle"`
	CameraMovement    string `json:"camera_movement"`
	NarrativeFunction string `json:"narrative_function"` // setup_character | plot_reveal | emotional_peak | resolution
	Subtext           string `json:"subtext"`            // what character is hiding (for visual irony)
	Location          string `json:"location"`
	Atmosphere        string `json:"atmosphere"`
	Characters        []uint `json:"characters"`
	SceneID           *uint  `json:"scene_id"`
	Title             string `json:"title"`
}

// ============================================================================
// Regex Patterns for Story Bible Parsing
// ============================================================================

var (
	// narrativeSectionPattern matches top-level section markers: [STORY_BIBLE], [CHARACTERS], etc.
	narrativeSectionPattern = regexp.MustCompile(`(?i)^\[(STORY_BIBLE|CHARACTERS|PROLOGUE|MUSIC_SEGMENTS|LYRICS|EPILOGUE)\](.*)$`)

	// narrativeMusicSegmentPattern matches: (0:00 - 0:45) INTRO — emotion: solitude, false calm
	narrativeMusicSegmentPattern = regexp.MustCompile(`^\((\d{1,2}:\d{2})\s*[-–]\s*(\d{1,2}:\d{2})\)\s+([^—]+?)(?:\s+—\s+emotion:\s+(.+))?$`)

	// narrativeSyncPointPattern matches indented [SYNC_POINT] lines
	// Must be indented ≥2 spaces or 1 tab
	narrativeSyncPointPattern = regexp.MustCompile(`^[\t ][\t ]+\[SYNC_POINT\]\s+(.+?)\s+—\s+(convergent|irony|parallel)$`)

	// narrativeDurationPattern matches "duration: 75s" on section header lines
	narrativeDurationPattern = regexp.MustCompile(`(?i)duration:\s*(\d+)s`)
)

// ============================================================================
// Story Bible Parser
// ============================================================================

// parseStoryBible parses the structured [MARKER]-based Story Bible input
func parseStoryBible(script string) (*StoryBible, error) {
	lines := strings.Split(strings.ReplaceAll(script, "\r\n", "\n"), "\n")

	bible := &StoryBible{}

	// Current section state
	currentSection := ""
	var currentLines []string

	flushSection := func() {
		content := strings.TrimSpace(strings.Join(currentLines, "\n"))
		switch currentSection {
		case "STORY_BIBLE":
			bible.WorldDescription = content
			// Extract "Emotional core:" if present
			for _, line := range currentLines {
				if strings.HasPrefix(strings.ToLower(strings.TrimSpace(line)), "emotional core:") {
					idx := strings.Index(strings.ToLower(line), "emotional core:")
					bible.EmotionalCore = strings.TrimSpace(line[idx+len("emotional core:"):])
				}
			}
		case "CHARACTERS":
			bible.Characters = parseNarrativeCharacters(currentLines)
		case "PROLOGUE":
			bible.PrologueDesc = content
		case "EPILOGUE":
			bible.EpilogueDesc = content
		case "MUSIC_SEGMENTS":
			bible.MusicSegments = parseNarrativeMusicSegments(currentLines)
			for _, seg := range bible.MusicSegments {
				bible.TotalMusicSec += seg.EndSec - seg.StartSec
			}
		case "LYRICS":
			lyricsText := strings.Join(currentLines, "\n")
			blocks, err := parseLyricsInput(lyricsText)
			if err == nil {
				bible.LyricsBlocks = blocks
			}
			// Non-fatal if no lyrics blocks — [LYRICS] section is optional
		}
		currentLines = nil
	}

	for _, rawLine := range lines {
		// Check for section header
		if match := narrativeSectionPattern.FindStringSubmatch(strings.TrimSpace(rawLine)); match != nil {
			// Flush previous section
			flushSection()
			currentSection = strings.ToUpper(match[1])
			headerExtra := strings.TrimSpace(match[2])

			// Parse inline attributes on section header (e.g. "duration: 75s")
			switch currentSection {
			case "PROLOGUE":
				if dm := narrativeDurationPattern.FindStringSubmatch(headerExtra); dm != nil {
					bible.PrologueDuration, _ = strconv.Atoi(dm[1])
				}
			case "EPILOGUE":
				if dm := narrativeDurationPattern.FindStringSubmatch(headerExtra); dm != nil {
					bible.EpilogueDuration, _ = strconv.Atoi(dm[1])
				}
			}
			continue
		}

		if currentSection != "" {
			currentLines = append(currentLines, rawLine)
		}
	}
	flushSection()

	if bible.WorldDescription == "" && bible.PrologueDesc == "" {
		return nil, fmt.Errorf("invalid Story Bible: no [STORY_BIBLE] or [PROLOGUE] section found")
	}

	return bible, nil
}

// parseNarrativeCharacters parses pipe-delimited character lines: "Name | Description | Role"
func parseNarrativeCharacters(lines []string) []NarrativeCharacter {
	var chars []NarrativeCharacter
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "|", 3)
		char := NarrativeCharacter{}
		if len(parts) >= 1 {
			char.Name = strings.TrimSpace(parts[0])
		}
		if len(parts) >= 2 {
			char.Description = strings.TrimSpace(parts[1])
		}
		if len(parts) >= 3 {
			// Strip "Role:" prefix if present
			role := strings.TrimSpace(parts[2])
			if idx := strings.Index(strings.ToLower(role), "role:"); idx >= 0 {
				role = strings.TrimSpace(role[idx+5:])
			}
			char.Role = role
		}
		if char.Name != "" {
			chars = append(chars, char)
		}
	}
	return chars
}

// parseNarrativeMusicSegments parses timestamped music segment lines with optional [SYNC_POINT] children
func parseNarrativeMusicSegments(lines []string) []NarrativeMusicSegment {
	var segments []NarrativeMusicSegment
	var current *NarrativeMusicSegment

	for _, rawLine := range lines {
		// Check for indented [SYNC_POINT] first (before trimming)
		if current != nil {
			if spMatch := narrativeSyncPointPattern.FindStringSubmatch(rawLine); spMatch != nil {
				syncType := spMatch[2]
				if syncType == "" {
					syncType = "parallel" // default
				}
				current.SyncPoints = append(current.SyncPoints, NarrativeSyncPoint{
					Description: strings.TrimSpace(spMatch[1]),
					SyncType:    syncType,
				})
				continue
			}
		}

		line := strings.TrimSpace(rawLine)
		if line == "" {
			continue
		}

		if match := narrativeMusicSegmentPattern.FindStringSubmatch(line); match != nil {
			// Flush previous segment
			if current != nil {
				segments = append(segments, *current)
			}
			startSec, _ := parseTimestampToSeconds(match[1])
			endSec, _ := parseTimestampToSeconds(match[2])
			name := strings.TrimSpace(match[3])
			emotion := strings.TrimSpace(match[4])

			current = &NarrativeMusicSegment{
				StartSec: startSec,
				EndSec:   endSec,
				Name:     name,
				Emotion:  emotion,
			}
		}
	}
	// Flush last segment
	if current != nil {
		segments = append(segments, *current)
	}
	return segments
}

// ============================================================================
// Narrative MV — Main Pipeline
// ============================================================================

// processNarrativeMVGeneration handles the narrative_mv split mode
// Uses a 2-phase AI pipeline: story planning → shot generation
func (s *StoryboardService) processNarrativeMVGeneration(
	taskID, episodeID string, dramaID uint,
	model, scriptContent, characterList, sceneList, propList string,
) {
	// 1. Parse Story Bible input
	bible, err := parseStoryBible(scriptContent)
	if err != nil {
		s.log.Errorw("Failed to parse Story Bible", "error", err, "task_id", taskID)
		if updateErr := s.taskService.UpdateTaskError(taskID, fmt.Errorf("Story Bible parsing failed: %w", err)); updateErr != nil {
			s.log.Errorw("Failed to update task error", "error", updateErr, "task_id", taskID)
		}
		return
	}

	s.log.Infow("Parsed Story Bible",
		"task_id", taskID, "episode_id", episodeID,
		"characters", len(bible.Characters),
		"music_segments", len(bible.MusicSegments),
		"lyrics_blocks", len(bible.LyricsBlocks),
		"total_music_sec", bible.TotalMusicSec)

	if err := s.taskService.UpdateTaskStatus(taskID, "processing", 10, "Story Bible parsed, running story planning AI..."); err != nil {
		s.log.Errorw("Failed to update task status", "error", err, "task_id", taskID)
	}

	// 2. Phase 1: Story Planning AI call
	plannerPrompt := s.promptI18n.WithDramaNarrativeMVPlannerPrompt(dramaID)
	phase1Prompt := buildNarrativePhase1Prompt(plannerPrompt, scriptContent, bible)

	client, getErr := s.aiService.GetAIClientForModel("text", model)
	if model != "" && getErr != nil {
		s.log.Warnw("Failed to get client for specified model, using default", "model", model, "error", getErr, "task_id", taskID)
	}

	callAI := func(p string) (string, error) {
		if model != "" && getErr == nil {
			return client.GenerateText(p, "")
		}
		return s.aiService.GenerateText(p, "")
	}

	phase1Text, err := callAI(phase1Prompt)
	if err != nil {
		s.log.Errorw("Failed Phase 1 (story planning)", "error", err, "task_id", taskID)
		if updateErr := s.taskService.UpdateTaskError(taskID, fmt.Errorf("story planning AI failed: %w", err)); updateErr != nil {
			s.log.Errorw("Failed to update task error", "error", updateErr, "task_id", taskID)
		}
		return
	}

	var narrativePlan NarrativePlan
	if err := utils.SafeParseAIJSON(phase1Text, &narrativePlan); err != nil {
		s.log.Warnw("Failed to parse NarrativePlan JSON, will proceed without plan context",
			"error", err, "task_id", taskID, "response", phase1Text[:min(300, len(phase1Text))])
		// Non-fatal: proceed with Phase 2 using an empty plan
	} else {
		// Override TotalMusicSec from plan if AI computed a different value (trust bible's computed value)
		narrativePlan.TotalMusicDurationSec = bible.TotalMusicSec
	}

	if err := s.taskService.UpdateTaskStatus(taskID, "processing", 30, "Story plan complete, generating shot list..."); err != nil {
		s.log.Errorw("Failed to update task status", "error", err, "task_id", taskID)
	}

	// 3. Phase 2: Shot Generation AI call
	directorPrompt := s.promptI18n.WithDramaNarrativeMVDirectorPrompt(dramaID)
	visualStylePrompt := s.promptI18n.WithDramaNarrativeMVVisualStylePrompt(dramaID)
	formatInstructions := prompts.Get("storyboard_narrative_format.txt")
	phase2Prompt := buildNarrativePhase2Prompt(directorPrompt, visualStylePrompt, scriptContent, bible, narrativePlan, characterList, sceneList, propList, formatInstructions)

	phase2Text, err := callAI(phase2Prompt)
	if err != nil {
		s.log.Errorw("Failed Phase 2 (shot generation)", "error", err, "task_id", taskID)
		if updateErr := s.taskService.UpdateTaskError(taskID, fmt.Errorf("shot generation AI failed: %w", err)); updateErr != nil {
			s.log.Errorw("Failed to update task error", "error", updateErr, "task_id", taskID)
		}
		return
	}

	if err := s.taskService.UpdateTaskStatus(taskID, "processing", 55, "AI response received, parsing shots..."); err != nil {
		s.log.Errorw("Failed to update task status", "error", err, "task_id", taskID)
	}

	// 4. Parse NarrativeShot[]
	var shots []NarrativeShot
	if err := utils.SafeParseAIJSON(phase2Text, &shots); err != nil {
		s.log.Errorw("Failed to parse NarrativeShot JSON",
			"error", err, "response", phase2Text[:min(500, len(phase2Text))], "task_id", taskID)
		if updateErr := s.taskService.UpdateTaskError(taskID, fmt.Errorf("failed to parse shot list: %w", err)); updateErr != nil {
			s.log.Errorw("Failed to update task error", "error", updateErr, "task_id", taskID)
		}
		return
	}

	// 5. Post-process: re-number + assign timestamps + fix durations
	shots = postProcessNarrativeShots(shots)

	// 6. Validate Part 2 duration coverage
	validateMusicFilmDuration(shots, bible.TotalMusicSec, taskID, s.log)

	if err := s.taskService.UpdateTaskStatus(taskID, "processing", 70, "Saving narrative shots..."); err != nil {
		s.log.Errorw("Failed to update task status", "error", err, "task_id", taskID)
	}

	// 7. Save to database
	if err := s.saveNarrativeShots(episodeID, dramaID, shots); err != nil {
		s.log.Errorw("Failed to save narrative shots", "error", err, "task_id", taskID)
		if updateErr := s.taskService.UpdateTaskError(taskID, fmt.Errorf("failed to save shots: %w", err)); updateErr != nil {
			s.log.Errorw("Failed to update task error", "error", updateErr, "task_id", taskID)
		}
		return
	}

	// Trigger style distillation asynchronously (non-blocking)
	epIDUint, _ := strconv.ParseUint(episodeID, 10, 32)
	go s.styleDistillService.BatchDistillStyles(uint(epIDUint), dramaID)

	if err := s.taskService.UpdateTaskStatus(taskID, "processing", 90, "Updating episode duration..."); err != nil {
		s.log.Errorw("Failed to update task status", "error", err, "task_id", taskID)
	}

	// Calculate total duration
	totalDuration := 0
	for _, shot := range shots {
		totalDuration += shot.DurationSec
	}

	// Count shots by part for logging
	var prologueCount, musicCount, epilogueCount int
	for _, shot := range shots {
		switch shot.NarrativePart {
		case "prologue":
			prologueCount++
		case "music_film":
			musicCount++
		case "epilogue":
			epilogueCount++
		}
	}

	s.log.Infow("Narrative MV storyboard generated",
		"task_id", taskID, "episode_id", episodeID,
		"total_shots", len(shots),
		"prologue_shots", prologueCount,
		"music_film_shots", musicCount,
		"epilogue_shots", epilogueCount,
		"total_duration_seconds", totalDuration)

	durationMinutes := (totalDuration + 59) / 60
	if err := s.db.Model(&models.Episode{}).Where("id = ?", episodeID).Update("duration", durationMinutes).Error; err != nil {
		s.log.Errorw("Failed to update episode duration", "error", err, "task_id", taskID)
	}

	resultData := gin.H{
		"storyboards":       shots,
		"total":             len(shots),
		"total_duration":    totalDuration,
		"duration_minutes":  durationMinutes,
		"mode":              "narrative_mv",
		"prologue_shots":    prologueCount,
		"music_film_shots":  musicCount,
		"epilogue_shots":    epilogueCount,
	}

	if err := s.taskService.UpdateTaskResult(taskID, resultData); err != nil {
		s.log.Errorw("Failed to update task result", "error", err, "task_id", taskID)
		return
	}

	s.log.Infow("Narrative MV generation completed", "task_id", taskID, "episode_id", episodeID)
}

// ============================================================================
// Prompt Builders
// ============================================================================

// buildNarrativePhase1Prompt constructs the story planning prompt for Phase 1
func buildNarrativePhase1Prompt(plannerSystemPrompt, rawScript string, bible *StoryBible) string {
	var sb strings.Builder

	sb.WriteString(plannerSystemPrompt)
	sb.WriteString("\n\n=== STORY BIBLE INPUT ===\n")
	sb.WriteString(rawScript)
	sb.WriteString("\n\n=== COMPUTED ANALYSIS ===\n")
	sb.WriteString(fmt.Sprintf("Total music duration: %d seconds\n", bible.TotalMusicSec))
	sb.WriteString(fmt.Sprintf("Characters defined: %d\n", len(bible.Characters)))
	sb.WriteString(fmt.Sprintf("Music segments: %d\n", len(bible.MusicSegments)))

	if len(bible.MusicSegments) > 0 {
		sb.WriteString("\nMusic segment breakdown:\n")
		for _, seg := range bible.MusicSegments {
			sb.WriteString(fmt.Sprintf("  %s (%d:%02d – %d:%02d, %ds)",
				seg.Name,
				seg.StartSec/60, seg.StartSec%60,
				seg.EndSec/60, seg.EndSec%60,
				seg.EndSec-seg.StartSec))
			if seg.Emotion != "" {
				sb.WriteString(fmt.Sprintf(" — emotion: %s", seg.Emotion))
			}
			sb.WriteString("\n")
			for _, sp := range seg.SyncPoints {
				sb.WriteString(fmt.Sprintf("    [SYNC POINT: %s — %s]\n", sp.Description, sp.SyncType))
			}
		}
	}

	return sb.String()
}

// buildNarrativePhase2Prompt constructs the shot director prompt for Phase 2
// visualStylePrompt is optional — when set (e.g. CG5 3-act template), it injects
// specific visual language rules BEFORE the generic director instructions.
func buildNarrativePhase2Prompt(
	directorSystemPrompt, visualStylePrompt, rawScript string,
	bible *StoryBible,
	plan NarrativePlan,
	characterList, sceneList, propList string,
	formatInstructions string,
) string {
	var sb strings.Builder

	// Inject visual style template first (e.g. CG5 3-act cinematography rules)
	if visualStylePrompt != "" {
		sb.WriteString(visualStylePrompt)
		sb.WriteString("\n\n=== DIRECTOR INSTRUCTIONS (applied on top of visual style above) ===\n")
	}
	sb.WriteString(directorSystemPrompt)
	sb.WriteString("\n\n=== STORY BIBLE ===\n")
	sb.WriteString(rawScript)

	sb.WriteString("\n\n=== NARRATIVE PLAN (from Phase 1 analysis) ===\n")
	if plan.LightingArc != "" {
		sb.WriteString(fmt.Sprintf("Lighting arc: %s\n", plan.LightingArc))
	}
	if len(plan.NarrativeThread) > 0 {
		sb.WriteString("Narrative motifs to maintain:\n")
		for _, thread := range plan.NarrativeThread {
			sb.WriteString(fmt.Sprintf("  - %s: %s (appears in: %s)\n",
				thread.Motif, thread.Meaning, strings.Join(thread.AppearsIn, ", ")))
		}
	}
	if len(plan.CharacterArcs) > 0 {
		sb.WriteString("Character arcs:\n")
		for name, arc := range plan.CharacterArcs {
			sb.WriteString(fmt.Sprintf("  - %s: %s\n", name, arc))
		}
	}
	if len(plan.SyncPoints) > 0 {
		sb.WriteString("Confirmed story-music sync points:\n")
		for _, sp := range plan.SyncPoints {
			sb.WriteString(fmt.Sprintf("  - At %s (%s): [%s] %s\n",
				sp.MusicTimestamp, sp.Segment, sp.SyncType, sp.Description))
		}
	}

	sb.WriteString(fmt.Sprintf("\n=== DURATION REQUIREMENTS ===\n"))
	sb.WriteString(fmt.Sprintf("Prologue target duration: %d seconds\n", bible.PrologueDuration))
	sb.WriteString(fmt.Sprintf("Music Film duration (MUST cover): %d seconds\n", bible.TotalMusicSec))
	sb.WriteString(fmt.Sprintf("Epilogue target duration: %d seconds\n", bible.EpilogueDuration))

	sb.WriteString(fmt.Sprintf("\n\nCharacter List: %s\n", characterList))
	sb.WriteString(fmt.Sprintf("Scene List: %s\n", sceneList))
	sb.WriteString(fmt.Sprintf("Prop List: %s\n", propList))

	sb.WriteString("\n\n")
	sb.WriteString(formatInstructions)

	return sb.String()
}

// ============================================================================
// Post-Processing (A3.5 — Timestamp Assignment)
// ============================================================================

// formatSecondsToTimestamp converts integer seconds to "M:SS" format
func formatSecondsToTimestamp(totalSec int) string {
	minutes := totalSec / 60
	seconds := totalSec % 60
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

// defaultDurationForPart returns a fallback duration for a shot based on its narrative part
func defaultDurationForPart(part string) int {
	switch part {
	case "prologue":
		return 5
	case "music_film":
		return 4
	case "epilogue":
		return 6
	default:
		return 4
	}
}

// postProcessNarrativeShots re-numbers shots and assigns cumulative timestamps
// AI only outputs duration_sec; timestamps are computed here (Step A3.5 in plan)
func postProcessNarrativeShots(shots []NarrativeShot) []NarrativeShot {
	cumulativeSec := 0
	for i := range shots {
		shots[i].ShotID = i + 1

		// Apply fallback for missing or invalid durations
		if shots[i].DurationSec <= 0 {
			shots[i].DurationSec = defaultDurationForPart(shots[i].NarrativePart)
		}

		// Assign absolute video timestamps based on cumulative duration
		shots[i].TimestampStart = formatSecondsToTimestamp(cumulativeSec)
		cumulativeSec += shots[i].DurationSec
		shots[i].TimestampEnd = formatSecondsToTimestamp(cumulativeSec)
	}
	return shots
}

// validateMusicFilmDuration logs a warning if Part 2 duration deviates more than ±10%
func validateMusicFilmDuration(shots []NarrativeShot, expectedSec int, taskID string, log interface{ Warnw(string, ...interface{}) }) {
	if expectedSec <= 0 {
		return
	}
	musicFilmTotal := 0
	for _, shot := range shots {
		if shot.NarrativePart == "music_film" {
			musicFilmTotal += shot.DurationSec
		}
	}
	if musicFilmTotal == 0 {
		log.Warnw("Music Film part has zero shots — possible AI output error", "task_id", taskID)
		return
	}
	tolerance := float64(expectedSec) * 0.10
	diff := abs(musicFilmTotal - expectedSec)
	if float64(diff) > tolerance {
		log.Warnw("Music Film duration mismatch",
			"task_id", taskID,
			"expected_sec", expectedSec,
			"actual_sec", musicFilmTotal,
			"diff_sec", diff)
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// ============================================================================
// Save Function
// ============================================================================

// saveNarrativeShots maps NarrativeShot[] to models.Storyboard and saves to DB
// Mirrors saveNurseryRhymeShots() pattern
func (s *StoryboardService) saveNarrativeShots(episodeID string, dramaID uint, shots []NarrativeShot) error {
	epID, err := strconv.ParseUint(episodeID, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid episode ID: %s", episodeID)
	}

	if len(shots) == 0 {
		return fmt.Errorf("AI returned 0 shots, refusing to save")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
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

		// Batch-load character descriptions
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

		strPtr := func(v string) *string {
			if v == "" {
				return nil
			}
			return &v
		}
		boolPtr := func(v bool) *bool { return &v }

		for _, shot := range shots {
			// Build character descriptions
			var charDescParts []string
			for _, cid := range shot.Characters {
				if d, ok := charDescMap[cid]; ok {
					charDescParts = append(charDescParts, d)
				}
			}
			charDescs := strings.Join(charDescParts, "; ")

			// Generate image and video prompts using existing helper
			sbForPrompt := Storyboard{
				ShotNumber:     shot.ShotID,
				Title:          shot.Title,
				ShotType:       shot.ShotType,
				Angle:          shot.CameraAngle,
				Movement:       shot.CameraMovement,
				Location:       shot.Location,
				Atmosphere:     shot.Atmosphere,
				Action:         shot.VisualDescription,
				Duration:       shot.DurationSec,
				CharacterDescs: charDescs,
			}
			imagePrompt := s.generateImagePrompt(sbForPrompt, "")
			videoPrompt := s.generateVideoPrompt(sbForPrompt)

			audioMode := "narrator_only" // default for narrative film

			hasMusicVal := shot.HasMusic

			storyboard := models.Storyboard{
				EpisodeID:        uint(epID),
				SceneID:          shot.SceneID,
				StoryboardNumber: shot.ShotID,
				Title:            strPtr(shot.Title),
				Location:         strPtr(shot.Location),
				ShotType:         strPtr(shot.ShotType),
				Angle:            strPtr(shot.CameraAngle),
				Movement:         strPtr(shot.CameraMovement),
				Description:      strPtr(shot.VisualDescription),
				Action:           strPtr(shot.VisualDescription),
				Atmosphere:       strPtr(shot.Atmosphere),
				ImagePrompt:      &imagePrompt,
				VideoPrompt:      &videoPrompt,
				VideoPromptSource: "auto",
				Duration:         shot.DurationSec,
				// Use LyricsText field to store lyric anchor for post-prod reference
				LyricsText:      strPtr(shot.LyricsAnchor),
				SectionType:     strPtr(shot.MusicSegment),
				AudioMode:       &audioMode,
				ShotRole:        strPtr(shot.NarrativeFunction),
				VisualType:      strPtr("literal"),
				// Narrative MV specific fields
				NarrativePart:   strPtr(shot.NarrativePart),
				HasMusic:        boolPtr(hasMusicVal),
				MusicSegment:    strPtr(shot.MusicSegment),
				MusicSyncType:   strPtr(shot.MusicSyncType),
				ActingNote:      strPtr(shot.ActingNote),
				LyricsAnchor:    strPtr(shot.LyricsAnchor),
			}

			if err := tx.Create(&storyboard).Error; err != nil {
				s.log.Errorw("Failed to create narrative shot", "error", err, "shot_id", shot.ShotID)
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
		}

		s.log.Infow("Narrative shots saved successfully", "episode_id", episodeID, "count", len(shots))
		return nil
	})
}
