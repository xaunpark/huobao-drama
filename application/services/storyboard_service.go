package services

import (
	"encoding/json"
	"regexp"
	"strconv"

	"fmt"
	"strings"

	models "github.com/drama-generator/backend/domain/models"
	"github.com/drama-generator/backend/pkg/config"
	"github.com/drama-generator/backend/application/prompts"
	"github.com/drama-generator/backend/pkg/logger"
	"github.com/drama-generator/backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type StoryboardService struct {
	db          *gorm.DB
	aiService   *AIService
	taskService *TaskService
	log         *logger.Logger
	config      *config.Config
	promptI18n  *PromptI18n
}

func NewStoryboardService(db *gorm.DB, cfg *config.Config, log *logger.Logger) *StoryboardService {
	return &StoryboardService{
		db:          db,
		aiService:   NewAIService(db, log),
		taskService: NewTaskService(db, log),
		log:         log,
		config:      cfg,
		promptI18n:  NewPromptI18n(cfg),
	}
}

type Storyboard struct {
	ShotNumber  int    `json:"shot_number"`
	Title       string `json:"title"`        // 镜头标题
	ShotType    string `json:"shot_type"`    // 景别
	Angle       string `json:"angle"`        // 镜头角度
	Time        string `json:"time"`         // 时间
	Location    string `json:"location"`     // 地点
	SceneID     *uint  `json:"scene_id"`     // 背景ID（AI直接返回，可为null）
	Movement    string `json:"movement"`     // 运镜
	Action      string `json:"action"`       // 动作
	Dialogue    string `json:"dialogue"`     // 对话/独白
	Result      string `json:"result"`       // 画面结果
	Atmosphere  string `json:"atmosphere"`   // 环境氛围
	Emotion     string `json:"emotion"`      // 情绪
	Duration    int    `json:"duration"`     // 时长（秒）
	BgmPrompt   string `json:"bgm_prompt"`   // 配乐提示词
	SoundEffect string `json:"sound_effect"` // 音效描述
	Characters  []uint `json:"characters"`   // 涉及的角色ID列表
	Props       []uint `json:"props"`        // 涉及的道具ID列表
	IsPrimary   bool   `json:"is_primary"`   // 是否主镜
}

type GenerateStoryboardResult struct {
	Storyboards []Storyboard `json:"storyboards"`
	Total       int          `json:"total"`
}

// VoiceoverShot — AI output struct for visual_unit mode (voice-over director)
type VoiceoverShot struct {
	ShotID            int      `json:"shot_id"`
	ScriptSegment     string   `json:"script_segment"`
	ScriptStartChar   int      `json:"script_start_char"`
	ScriptEndChar     int      `json:"script_end_char"`
	EstimatedDuration int      `json:"estimated_duration_sec"`
	VisualType        string   `json:"visual_type"`
	ShotRole          string   `json:"shot_role"`
	VisualDescription string   `json:"visual_description"`
	ReasonForShot     string   `json:"reason_for_shot"`
	TriggeredRules    []string `json:"triggered_rules"`
	Title             string   `json:"title"`
	ShotType          string   `json:"shot_type"`
	Angle             string   `json:"angle"`
	Movement          string   `json:"movement"`
	Location          string   `json:"location"`
	Time              string   `json:"time"`
	Atmosphere        string   `json:"atmosphere"`
	// Audio strategy
	AudioMode       string `json:"audio_mode"`
	NarratorEnabled bool   `json:"narrator_enabled"`
	NarratorDucking bool   `json:"narrator_ducking"`
	DialogueType    string `json:"dialogue_type"`
	DialogueText    string `json:"dialogue_text"`
	AmbienceType    string `json:"ambience_type"`
	AmbienceLevel   string `json:"ambience_level"`
	MusicMood       string `json:"music_mood"`
	MusicLevel      string `json:"music_level"`
	SoundEffect     string `json:"sound_effect"`
	BgmPrompt       string `json:"bgm_prompt"`
	// References
	Characters []uint `json:"characters"`
	Props      []uint `json:"props"`
	SceneID    *uint  `json:"scene_id"`
}

func (s *StoryboardService) GenerateStoryboard(episodeID string, model string, splitMode string) (string, error) {
	// 从数据库获取剧集信息
	var episode struct {
		ID            string
		ScriptContent *string
		Description   *string
		DramaID       string
	}

	err := s.db.Table("episodes").
		Select("episodes.id, episodes.script_content, episodes.description, episodes.drama_id").
		Joins("INNER JOIN dramas ON dramas.id = episodes.drama_id").
		Where("episodes.id = ?", episodeID).
		First(&episode).Error

	if err != nil {
		return "", fmt.Errorf("剧集不存在或无权限访问")
	}

	// 获取剧本内容
	var scriptContent string
	if episode.ScriptContent != nil && *episode.ScriptContent != "" {
		scriptContent = *episode.ScriptContent
	} else if episode.Description != nil && *episode.Description != "" {
		scriptContent = *episode.Description
	} else {
		return "", fmt.Errorf("剧本内容为空，请先生成剧集内容")
	}

	// 获取该剧本的所有角色
	var characters []models.Character
	if err := s.db.Where("drama_id = ?", episode.DramaID).Order("name ASC").Find(&characters).Error; err != nil {
		return "", fmt.Errorf("获取角色列表失败: %w", err)
	}

	// 构建角色列表字符串（包含ID和名称）
	characterList := "无角色"
	if len(characters) > 0 {
		var charInfoList []string
		for _, char := range characters {
			charInfoList = append(charInfoList, fmt.Sprintf(`{"id": %d, "name": "%s"}`, char.ID, char.Name))
		}
		characterList = fmt.Sprintf("[%s]", strings.Join(charInfoList, ", "))
	}

	// 获取该项目已提取的场景列表（项目级）
	var scenes []models.Scene
	if err := s.db.Where("drama_id = ?", episode.DramaID).Order("location ASC, time ASC").Find(&scenes).Error; err != nil {
		s.log.Warnw("Failed to get scenes", "error", err)
	}

	// 构建场景列表字符串（包含ID、地点、时间）
	sceneList := "无场景"
	if len(scenes) > 0 {
		var sceneInfoList []string
		for _, bg := range scenes {
			sceneInfoList = append(sceneInfoList, fmt.Sprintf(`{"id": %d, "location": "%s", "time": "%s"}`, bg.ID, bg.Location, bg.Time))
		}
		sceneList = fmt.Sprintf("[%s]", strings.Join(sceneInfoList, ", "))
	}

	// 获取该项目已提取的道具列表（项目级）
	var props []models.Prop
	if err := s.db.Where("drama_id = ?", episode.DramaID).Find(&props).Error; err != nil {
		s.log.Warnw("Failed to get props", "error", err)
	}

	// 构建道具列表字符串（包含ID、名称）
	propList := "无道具"
	if len(props) > 0 {
		var propInfoList []string
		for _, p := range props {
			propInfoList = append(propInfoList, fmt.Sprintf(`{"id": %d, "name": "%s"}`, p.ID, p.Name))
		}
		propList = fmt.Sprintf("[%s]", strings.Join(propInfoList, ", "))
	}


	// Auto-detect split mode: if script has timestamp patterns, use "preserve" mode
	effectiveSplitMode := splitMode
	if effectiveSplitMode == "" || effectiveSplitMode == "auto" {
		if detectTimestampPattern(scriptContent) {
			effectiveSplitMode = "preserve"
			s.log.Infow("Auto-detected timestamp pattern in script, using preserve mode",
				"episode_id", episodeID)
		} else {
			effectiveSplitMode = "breakdown"
		}
	}

	// 创建异步任务
	task, err := s.taskService.CreateTask("storyboard_generation", episodeID)
	if err != nil {
		s.log.Errorw("Failed to create task", "error", err)
		return "", fmt.Errorf("创建任务失败: %w", err)
	}

	s.log.Infow("Generating storyboard asynchronously",
		"task_id", task.ID,
		"episode_id", episodeID,
		"drama_id", episode.DramaID,
		"script_length", len(scriptContent),
		"character_count", len(characters),
		"characters", characterList,
		"scene_count", len(scenes),
		"scenes", sceneList,
		"prop_count", len(props),
		"props", propList,
		"split_mode", effectiveSplitMode)

	// 启动后台goroutine处理AI调用和后续逻辑
	dramaIDUint, _ := strconv.ParseUint(episode.DramaID, 10, 32)
	go s.processStoryboardGeneration(task.ID, episodeID, uint(dramaIDUint), model, effectiveSplitMode, scriptContent, characterList, sceneList, propList)

	// 立即返回任务ID
	return task.ID, nil
}

// detectTimestampPattern checks if the script contains timestamp patterns
// like (0:00 – 0:06), (0:06 – 0:12), indicating pre-defined shot boundaries.
func detectTimestampPattern(script string) bool {
	// Match patterns like (0:00 – 0:06), (00:00 - 00:06), (0:00–0:06), etc.
	// Also match patterns like (0:00 ~ 0:06) or [0:00 - 0:06]
	patterns := []string{
		`\(?\d{1,2}:\d{2}\s*[–\-~]\s*\d{1,2}:\d{2}\)?`,  // Timestamp ranges
		`(?i)^\s*shot\s+\d+`,                                // "Shot 1", "Shot 2"
		`(?i)^\s*scene\s+\d+`,                               // "Scene 1", "Scene 2"
	}

	timestampCount := 0
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllString(script, -1)
		timestampCount += len(matches)
	}

	// If we detect 5+ timestamp/shot markers, it's a pre-structured script
	return timestampCount >= 5
}

// processStoryboardGeneration 后台处理故事板生成
func (s *StoryboardService) processStoryboardGeneration(taskID, episodeID string, dramaID uint, model, splitMode, scriptContent, characterList, sceneList, propList string) {
	// 更新任务状态为处理中
	if err := s.taskService.UpdateTaskStatus(taskID, "processing", 10, "准备生成分镜头..."); err != nil {
		s.log.Errorw("Failed to update task status", "error", err, "task_id", taskID)
		return
	}

	s.log.Infow("Processing storyboard generation", "task_id", taskID, "episode_id", episodeID, "split_mode", splitMode)

	// Route to visual_unit mode if selected
	if splitMode == "visual_unit" {
		s.log.Infow("Using VISUAL_UNIT mode — AI Director voice-over shot planning", "task_id", taskID)
		s.processVisualUnitGeneration(taskID, episodeID, dramaID, model, scriptContent, characterList, sceneList, propList)
		return
	}

	// Choose system prompt based on split mode
	var systemPrompt string
	if splitMode == "preserve" {
		systemPrompt = prompts.Get("storyboard_preserve_shots.txt")
		s.log.Infow("Using PRESERVE mode — keeping script shot structure", "task_id", taskID)
	} else {
		systemPrompt = s.promptI18n.WithDramaStoryboardSystemPrompt(dramaID)
		s.log.Infow("Using BREAKDOWN mode — AI action unit analysis", "task_id", taskID)
	}

	scriptLabel := s.promptI18n.FormatUserPrompt("script_content_label")
	taskLabel := s.promptI18n.FormatUserPrompt("task_label")
	var taskInstruction string
	if splitMode == "preserve" {
		taskInstruction = "Preserve each shot/block from the script as a separate storyboard entry. Do NOT merge or skip any shots. Enrich each with cinematography metadata."
	} else {
		taskInstruction = s.promptI18n.FormatUserPrompt("task_instruction")
	}
	charListLabel := s.promptI18n.FormatUserPrompt("character_list_label")
	charConstraint := s.promptI18n.FormatUserPrompt("character_constraint")
	sceneListLabel := s.promptI18n.FormatUserPrompt("scene_list_label")
	sceneConstraint := s.promptI18n.FormatUserPrompt("scene_constraint")
	propListLabel := s.promptI18n.FormatUserPrompt("prop_list_label")
	propConstraint := s.promptI18n.FormatUserPrompt("prop_constraint")
	formatInstructions := prompts.Get("storyboard_format_instructions.txt")

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
%s

%s`, systemPrompt, scriptLabel, scriptContent, taskLabel, taskInstruction, charListLabel, characterList, charConstraint, sceneListLabel, sceneList, sceneConstraint, propListLabel, propList, propConstraint, formatInstructions)

	client, getErr := s.aiService.GetAIClientForModel("text", model)
	if model != "" && getErr != nil {
		s.log.Warnw("Failed to get client for specified model, using default", "model", model, "error", getErr, "task_id", taskID)
	}

	var text string
	var err error
	// Bỏ MaxTokens theo yêu cầu user (không giới hạn)
	if model != "" && getErr == nil {
		text, err = client.GenerateText(prompt, "")
	} else {
		text, err = s.aiService.GenerateText(prompt, "")
	}

	if err != nil {
		s.log.Errorw("Failed to generate storyboard", "error", err, "task_id", taskID)
		if updateErr := s.taskService.UpdateTaskError(taskID, fmt.Errorf("生成分镜头失败: %w", err)); updateErr != nil {
			s.log.Errorw("Failed to update task error", "error", updateErr, "task_id", taskID)
		}
		return
	}

	var result GenerateStoryboardResult
	var allStoryboards []Storyboard
	if err := utils.SafeParseAIJSON(text, &allStoryboards); err == nil {
		result.Storyboards = allStoryboards
		s.log.Infow("Parsed storyboard as array format", "count", len(allStoryboards), "task_id", taskID)
	} else {
		if err := utils.SafeParseAIJSON(text, &result); err != nil {
			s.log.Errorw("Failed to parse storyboard JSON in both formats", "error", err, "response", text[:min(500, len(text))], "task_id", taskID)
			if updateErr := s.taskService.UpdateTaskError(taskID, fmt.Errorf("解析分镜头结果失败: %w", err)); updateErr != nil {
				s.log.Errorw("Failed to update task error", "error", updateErr, "task_id", taskID)
			}
			return
		}
		allStoryboards = result.Storyboards
		s.log.Infow("Parsed storyboard as object format", "count", len(allStoryboards), "task_id", taskID)
	}

	// 更新任务进度
	if err := s.taskService.UpdateTaskStatus(taskID, "processing", 50, "分镜头全部生成完成，正在解析数据..."); err != nil {
		s.log.Errorw("Failed to update task status", "error", err, "task_id", taskID)
		return
	}

	// 重新编号所有的分镜，以保证连续性
	for i := range allStoryboards {
		allStoryboards[i].ShotNumber = i + 1
	}

	// 计算总时长（所有分镜时长之和）
	totalDuration := 0
	for _, sb := range allStoryboards {
		totalDuration += sb.Duration
	}

	s.log.Infow("Storyboard generated",
		"task_id", taskID,
		"episode_id", episodeID,
		"count", len(allStoryboards),
		"total_duration_seconds", totalDuration)

	// 更新任务进度
	if err := s.taskService.UpdateTaskStatus(taskID, "processing", 70, "正在保存分镜头..."); err != nil {
		s.log.Errorw("Failed to update task status", "error", err, "task_id", taskID)
		return
	}

	// 保存分镜头到数据库
	if err := s.saveStoryboards(episodeID, allStoryboards); err != nil {
		s.log.Errorw("Failed to save storyboards", "error", err, "task_id", taskID)
		if updateErr := s.taskService.UpdateTaskError(taskID, fmt.Errorf("保存分镜头失败: %w", err)); updateErr != nil {
			s.log.Errorw("Failed to update task error", "error", updateErr, "task_id", taskID)
		}
		return
	}

	// 更新任务进度
	if err := s.taskService.UpdateTaskStatus(taskID, "processing", 90, "正在更新剧集时长..."); err != nil {
		s.log.Errorw("Failed to update task status", "error", err, "task_id", taskID)
		return
	}

	// 更新剧集时长（秒转分钟，向上取整）
	durationMinutes := (totalDuration + 59) / 60
	if err := s.db.Model(&models.Episode{}).Where("id = ?", episodeID).Update("duration", durationMinutes).Error; err != nil {
		s.log.Errorw("Failed to update episode duration", "error", err, "task_id", taskID)
		// 不中断流程，只记录错误
	} else {
		s.log.Infow("Episode duration updated",
			"task_id", taskID,
			"episode_id", episodeID,
			"duration_seconds", totalDuration,
			"duration_minutes", durationMinutes)
	}

	// 更新任务结果
	resultData := gin.H{
		"storyboards":      allStoryboards,
		"total":            len(allStoryboards),
		"total_duration":   totalDuration,
		"duration_minutes": durationMinutes,
	}

	if err := s.taskService.UpdateTaskResult(taskID, resultData); err != nil {
		s.log.Errorw("Failed to update task result", "error", err, "task_id", taskID)
		return
	}

	s.log.Infow("Storyboard generation completed", "task_id", taskID, "episode_id", episodeID)
}

// ScriptSegmentInfo represents a parsed segment from marked script input
type ScriptSegmentInfo struct {
	Type      string // "narrator", "dialogue", "crowd", "sfx"
	Character string // Character name (for dialogue/crowd)
	Text      string // The actual text content
	LineNum   int    // Original line number
}

// ShotBlock represents a user-defined shot with pre-grouped content (structured input)
type ShotBlock struct {
	ShotNumber int
	Duration   int    // 0 if not specified by user
	ShotType   string // "" if not specified
	AudioMode  string // "" if not specified
	Lines      []ScriptSegmentInfo
	RawContent string // Original text between markers (for script_segment)
}

// shotHeaderPattern matches "// SHOT 01", "// SHOT 2 | 6s", "// SHOT 03 | 5s | CU | narrator_only"
var shotHeaderPattern = regexp.MustCompile(
	`^//\s*SHOT\s+(\d+)` +
		`(?:\s*\|\s*(\d+)s)?` +
		`(?:\s*\|\s*([\w-]+))?` +
		`(?:\s*\|\s*(narrator_only|dialogue_dominant))?`,
)

// detectStructuredShots checks if script contains >= 3 "// SHOT" markers
func detectStructuredShots(script string) bool {
	re := regexp.MustCompile(`(?m)^//\s*SHOT\s+\d+`)
	matches := re.FindAllString(script, -1)
	return len(matches) >= 3
}

// parseStructuredShots splits script into ShotBlocks by "// SHOT" markers
func parseStructuredShots(script string) []ShotBlock {
	lines := strings.Split(script, "\n")
	tagPattern := regexp.MustCompile(`^\s*\[([^\]]+)\]\s*(.*)$`)

	var blocks []ShotBlock
	var current *ShotBlock
	var contentLines []string

	flushBlock := func() {
		if current == nil {
			return
		}
		// Build RawContent from non-empty content lines
		var rawParts []string
		for _, cl := range contentLines {
			if strings.TrimSpace(cl) != "" {
				rawParts = append(rawParts, strings.TrimSpace(cl))
			}
		}
		current.RawContent = strings.Join(rawParts, "\n")

		// Parse segments within the block using tag detection
		for i, cl := range contentLines {
			trimmed := strings.TrimSpace(cl)
			if trimmed == "" {
				continue
			}
			// Skip markdown metadata (headers, tables, etc.) but NOT // SHOT lines
			if isMarkdownMetadata(trimmed) {
				continue
			}

			matches := tagPattern.FindStringSubmatch(trimmed)
			if matches != nil {
				tag := strings.TrimSpace(matches[1])
				text := strings.TrimSpace(matches[2])
				tagUpper := strings.ToUpper(tag)

				seg := ScriptSegmentInfo{LineNum: i + 1, Text: text}
				switch tagUpper {
				case "CROWD":
					seg.Type = "crowd"
					seg.Character = "CROWD"
				case "SFX":
					seg.Type = "sfx"
				case "BGM":
					seg.Type = "bgm"
				case "CAM", "CAMERA":
					seg.Type = "camera"
				case "VFX":
					seg.Type = "vfx"
				case "NOTE", "DIR":
					seg.Type = "note"
				case "NARRATOR":
					seg.Type = "narrator"
				default:
					seg.Type = "dialogue"
					seg.Character = tag
				}
				current.Lines = append(current.Lines, seg)
			} else {
				current.Lines = append(current.Lines, ScriptSegmentInfo{
					Type:    "narrator",
					Text:    trimmed,
					LineNum: i + 1,
				})
			}
		}

		blocks = append(blocks, *current)
	}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if matches := shotHeaderPattern.FindStringSubmatch(trimmed); matches != nil {
			// Flush previous block
			flushBlock()

			shotNum, _ := strconv.Atoi(matches[1])
			dur := 0
			if matches[2] != "" {
				dur, _ = strconv.Atoi(matches[2])
			}

			current = &ShotBlock{
				ShotNumber: shotNum,
				Duration:   dur,
				ShotType:   matches[3],
				AudioMode:  matches[4],
			}
			contentLines = nil
			continue
		}

		// Accumulate content lines for current block
		if current != nil {
			contentLines = append(contentLines, line)
		}
	}

	// Flush last block
	flushBlock()

	return blocks
}

// buildStructuredAnalysis creates AI context text from parsed shot blocks
func buildStructuredAnalysis(blocks []ShotBlock) string {
	var sb strings.Builder
	sb.WriteString("[Pre-Structured Script — ENRICH-ONLY Mode]\n")
	sb.WriteString(fmt.Sprintf("This script contains %d pre-defined shots marked with // SHOT markers.\n", len(blocks)))
	sb.WriteString("You MUST output EXACTLY this many shots. Do NOT split or merge any shot.\n\n")

	sb.WriteString("Shot Summary:\n")
	for _, b := range blocks {
		durStr := "auto"
		if b.Duration > 0 {
			durStr = fmt.Sprintf("%ds", b.Duration)
		}
		typeStr := "auto"
		if b.ShotType != "" {
			typeStr = b.ShotType
		}
		modeStr := "auto"
		if b.AudioMode != "" {
			modeStr = b.AudioMode
		}

		dialogueCount := 0
		sfxCount := 0
		bgmCount := 0
		camCount := 0
		vfxCount := 0
		noteCount := 0
		for _, seg := range b.Lines {
			switch seg.Type {
			case "dialogue", "crowd":
				dialogueCount++
			case "sfx":
				sfxCount++
			case "bgm":
				bgmCount++
			case "camera":
				camCount++
			case "vfx":
				vfxCount++
			case "note":
				noteCount++
			}
		}

		// Build compact tag stats
		var tagParts []string
		if dialogueCount > 0 {
			tagParts = append(tagParts, fmt.Sprintf("dlg=%d", dialogueCount))
		}
		if sfxCount > 0 {
			tagParts = append(tagParts, fmt.Sprintf("sfx=%d", sfxCount))
		}
		if bgmCount > 0 {
			tagParts = append(tagParts, fmt.Sprintf("bgm=%d", bgmCount))
		}
		if camCount > 0 {
			tagParts = append(tagParts, fmt.Sprintf("cam=%d", camCount))
		}
		if vfxCount > 0 {
			tagParts = append(tagParts, fmt.Sprintf("vfx=%d", vfxCount))
		}
		if noteCount > 0 {
			tagParts = append(tagParts, fmt.Sprintf("note=%d", noteCount))
		}
		tagStr := ""
		if len(tagParts) > 0 {
			tagStr = ", tags: " + strings.Join(tagParts, " ")
		}

		sb.WriteString(fmt.Sprintf("- SHOT %d: duration=%s, type=%s, mode=%s%s\n",
			b.ShotNumber, durStr, typeStr, modeStr, tagStr))
	}

	sb.WriteString("\nRULES FOR STRUCTURED INPUT:\n")
	sb.WriteString("1. Output EXACTLY the shots listed above — same count, same order\n")
	sb.WriteString("2. Each shot's script_segment = the FULL content between its // SHOT markers\n")
	sb.WriteString("3. If duration was pre-specified (not 'auto'), use that exact value for estimated_duration_sec\n")
	sb.WriteString("4. If shot_type was pre-specified (not 'auto'), use that exact value\n")
	sb.WriteString("5. If audio_mode was pre-specified (not 'auto'), use that exact value\n")
	sb.WriteString("6. For fields marked 'auto', infer the best value from shot content and [Tags]\n")
	sb.WriteString("\nTAG MAPPING RULES:\n")
	sb.WriteString("- Lines WITHOUT [tags] = NARRATOR → audio_mode='narrator_only'\n")
	sb.WriteString("- [Character Name] = DIALOGUE → audio_mode='dialogue_dominant', use text as dialogue_text\n")
	sb.WriteString("- [CROWD] = CROWD → audio_mode='dialogue_dominant', dialogue_type='crowd'\n")
	sb.WriteString("- [SFX] text → put in sound_effect field. Does NOT change audio_mode\n")
	sb.WriteString("- [BGM] text → put in bgm_prompt field. Does NOT change audio_mode\n")
	sb.WriteString("- [CAM] text → use as camera movement/angle instruction for the shot\n")
	sb.WriteString("- [VFX] text → append to visual_description as visual effects instruction\n")
	sb.WriteString("- [NOTE] text → use as director's note, incorporate into visual_description or reason_for_shot\n")
	sb.WriteString("- Do NOT invent dialogue that isn't in the script\n")

	return sb.String()
}

// parseMarkedScript detects [Character], [CROWD], [SFX] tags in the script
// and returns a structured analysis string to inject into the AI prompt.
// If no tags are found, returns empty string (pure narrator mode — backward compatible).
// Automatically strips markdown metadata (headers, tables, shot annotations, etc.)
func parseMarkedScript(script string) (segments []ScriptSegmentInfo, analysisText string) {
	lines := strings.Split(script, "\n")
	tagPattern := regexp.MustCompile(`^\s*\[([^\]]+)\]\s*(.*)$`)

	// First pass: check if any REAL script tags exist (skip metadata lines)
	hasAnyTag := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || isMarkdownMetadata(trimmed) {
			continue
		}
		if tagPattern.MatchString(trimmed) {
			hasAnyTag = true
			break
		}
	}

	// No tags found — pure narrator script, backward compatible
	if !hasAnyTag {
		return nil, ""
	}

	// Second pass: parse segments, skipping metadata
	inCodeBlock := false
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Track code block state (``` ... ```)
		if strings.HasPrefix(trimmed, "```") {
			inCodeBlock = !inCodeBlock
			continue
		}
		if inCodeBlock {
			continue
		}

		// Skip empty lines and markdown metadata
		if trimmed == "" || isMarkdownMetadata(trimmed) {
			continue
		}

		matches := tagPattern.FindStringSubmatch(trimmed)
		if matches != nil {
			tag := strings.TrimSpace(matches[1])
			text := strings.TrimSpace(matches[2])
			tagUpper := strings.ToUpper(tag)

			seg := ScriptSegmentInfo{LineNum: i + 1, Text: text}
			switch tagUpper {
			case "CROWD":
				seg.Type = "crowd"
				seg.Character = "CROWD"
			case "SFX":
				seg.Type = "sfx"
				seg.Character = ""
			case "BGM":
				seg.Type = "bgm"
				seg.Character = ""
			case "CAM", "CAMERA":
				seg.Type = "camera"
				seg.Character = ""
			case "VFX":
				seg.Type = "vfx"
				seg.Character = ""
			case "NOTE", "DIR":
				seg.Type = "note"
				seg.Character = ""
			case "NARRATOR":
				seg.Type = "narrator"
				seg.Character = ""
			default:
				seg.Type = "dialogue"
				seg.Character = tag
			}
			segments = append(segments, seg)
		} else {
			// No tag = narrator line
			segments = append(segments, ScriptSegmentInfo{
				Type:    "narrator",
				Text:    trimmed,
				LineNum: i + 1,
			})
		}
	}

	// Build analysis text
	var sb strings.Builder
	sb.WriteString("[Pre-Marked Script Analysis]\n")
	sb.WriteString("This script contains EXPLICIT dialogue markers. You MUST respect them:\n\n")

	narratorCount := 0
	dialogueCount := 0
	crowdCount := 0
	sfxCount := 0
	charNames := make(map[string]bool)

	for _, seg := range segments {
		switch seg.Type {
		case "narrator":
			narratorCount++
		case "dialogue":
			dialogueCount++
			charNames[seg.Character] = true
		case "crowd":
			crowdCount++
		case "sfx":
			sfxCount++
		}
	}

	sb.WriteString(fmt.Sprintf("- Total narrator segments: %d\n", narratorCount))
	sb.WriteString(fmt.Sprintf("- Total dialogue segments: %d\n", dialogueCount))
	sb.WriteString(fmt.Sprintf("- Total crowd segments: %d\n", crowdCount))
	sb.WriteString(fmt.Sprintf("- Total SFX cues: %d\n", sfxCount))

	if len(charNames) > 0 {
		names := make([]string, 0, len(charNames))
		for name := range charNames {
			names = append(names, name)
		}
		sb.WriteString(fmt.Sprintf("- Speaking characters: %s\n", strings.Join(names, ", ")))
	}

	sb.WriteString("\nRULES FOR MARKED SCRIPTS:\n")
	sb.WriteString("1. Lines WITHOUT [tags] are NARRATOR — set audio_mode to 'narrator_only'\n")
	sb.WriteString("2. Lines with [Character Name] are DIALOGUE — set audio_mode to 'dialogue_dominant', dialogue_text = the text, narrator_enabled = false\n")
	sb.WriteString("3. Lines with [CROWD] are CROWD — set audio_mode to 'dialogue_dominant', dialogue_type = 'crowd', narrator_ducking = true\n")
	sb.WriteString("4. Lines with [SFX] are sound effect cues — include in sound_effect field\n")
	sb.WriteString("5. DO NOT invent or add dialogue that isn't in the script\n")
	sb.WriteString("6. DO NOT move dialogue to a different position in the story\n")
	sb.WriteString("7. A dialogue line CAN be its own shot, or combined with adjacent narrator if very short\n")
	sb.WriteString("8. The script_segment field must include the ORIGINAL text WITH the [tag] markers\n")

	return segments, sb.String()
}

// isMarkdownMetadata returns true for lines that are markdown formatting/metadata
// and should be SKIPPED when parsing a marked script.
// This allows users to paste full production documents with structure maps,
// shot annotations, timing tables, etc. without manual cleanup.
func isMarkdownMetadata(line string) bool {
	// Markdown headers: # Title, ## Section, ### Subsection
	if strings.HasPrefix(line, "#") {
		return true
	}

	// Horizontal rules: ---, ***, ___
	stripped := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(line, " ", ""), "\t", ""), "\r", "")
	if len(stripped) >= 3 {
		allDash := true
		allStar := true
		allUnder := true
		for _, ch := range stripped {
			if ch != '-' {
				allDash = false
			}
			if ch != '*' {
				allStar = false
			}
			if ch != '_' {
				allUnder = false
			}
		}
		if allDash || allStar || allUnder {
			return true
		}
	}

	// Markdown table rows: | Column | Column |
	if strings.HasPrefix(line, "|") && strings.Count(line, "|") >= 2 {
		return true
	}

	// Bold shot annotations: **// SHOT 01 | 00:00–00:06 | 6s | BRAND IDENT**
	if strings.HasPrefix(line, "**//") || strings.HasPrefix(line, "**/ /") {
		return true
	}

	// Non-bold shot comment lines: // SHOT 01 ...
	// EXCEPT: preserve "// SHOT XX" markers used for structured input
	if strings.HasPrefix(line, "//") {
		if shotHeaderPattern.MatchString(line) {
			return false // This is a structural shot marker, NOT metadata
		}
		return true
	}

	// Lines starting with ** that contain metadata keywords (not dialogue)
	if strings.HasPrefix(line, "**") && strings.HasSuffix(line, "**") {
		lower := strings.ToLower(line)
		metaKeywords := []string{"kênh:", "format:", "thời lượng", "cú pháp:", "channel:", "duration:", "syntax:"}
		for _, kw := range metaKeywords {
			if strings.Contains(lower, kw) {
				return true
			}
		}
	}

	// Lines that are pure bold formatting with metadata patterns: **Kênh:** ..., **Format:** ...
	if strings.HasPrefix(line, "**") && strings.Contains(line, ":**") {
		return true
	}

	return false
}

// processVisualUnitGeneration handles the visual_unit split mode (AI Director for voice-over videos)
func (s *StoryboardService) processVisualUnitGeneration(taskID, episodeID string, dramaID uint, model, scriptContent, characterList, sceneList, propList string) {
	// Detect structured shot markers (// SHOT XX)
	isStructured := detectStructuredShots(scriptContent)
	var structuredBlocks []ShotBlock

	var systemPrompt string
	var scriptAnalysis string

	if isStructured {
		// === STRUCTURED MODE: User pre-defined shot boundaries ===
		structuredBlocks = parseStructuredShots(scriptContent)
		s.log.Infow("Detected structured shot markers (// SHOT), using ENRICH-ONLY mode",
			"task_id", taskID, "episode_id", episodeID, "shot_count", len(structuredBlocks))

		// Try custom template first, fallback to structured prompt
		customPrompt := s.promptI18n.WithDramaVisualUnitSystemPrompt(dramaID)
		defaultPrompt := prompts.Get("storyboard_visual_unit.txt")
		if customPrompt != defaultPrompt {
			// User has a custom template — use it (it should handle structured mode)
			systemPrompt = customPrompt
		} else {
			// No custom template — use dedicated structured prompt
			systemPrompt = prompts.Get("storyboard_visual_unit_structured.txt")
		}

		scriptAnalysis = buildStructuredAnalysis(structuredBlocks)
	} else {
		// === FREE-FORM MODE (unchanged behavior) ===
		systemPrompt = s.promptI18n.WithDramaVisualUnitSystemPrompt(dramaID)
		_, scriptAnalysis = parseMarkedScript(scriptContent)
	}

	scriptLabel := s.promptI18n.FormatUserPrompt("script_content_label")
	charListLabel := s.promptI18n.FormatUserPrompt("character_list_label")
	charConstraint := s.promptI18n.FormatUserPrompt("character_constraint")
	sceneListLabel := s.promptI18n.FormatUserPrompt("scene_list_label")
	sceneConstraint := s.promptI18n.FormatUserPrompt("scene_constraint")
	propListLabel := s.promptI18n.FormatUserPrompt("prop_list_label")
	propConstraint := s.promptI18n.FormatUserPrompt("prop_constraint")
	formatInstructions := prompts.Get("storyboard_visual_unit_format.txt")

	var prompt string
	if scriptAnalysis != "" {
		// Has analysis (structured or marked script) — inject between script and format
		logMsg := "Detected marked script with dialogue tags"
		if isStructured {
			logMsg = "Using structured shot analysis for ENRICH-ONLY mode"
		}
		s.log.Infow(logMsg, "task_id", taskID, "episode_id", episodeID)

		prompt = fmt.Sprintf(`%s

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
			scriptAnalysis,
			charListLabel, characterList, charConstraint,
			sceneListLabel, sceneList, sceneConstraint,
			propListLabel, propList, propConstraint,
			formatInstructions)
	} else {
		// Pure narrator script — original flow
		prompt = fmt.Sprintf(`%s

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
			charListLabel, characterList, charConstraint,
			sceneListLabel, sceneList, sceneConstraint,
			propListLabel, propList, propConstraint,
			formatInstructions)
	}

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
		s.log.Errorw("Failed to generate visual unit storyboard", "error", err, "task_id", taskID)
		if updateErr := s.taskService.UpdateTaskError(taskID, fmt.Errorf("Visual unit generation failed: %w", err)); updateErr != nil {
			s.log.Errorw("Failed to update task error", "error", updateErr, "task_id", taskID)
		}
		return
	}

	// Parse AI output as VoiceoverShot array
	var shots []VoiceoverShot
	if err := utils.SafeParseAIJSON(text, &shots); err != nil {
		s.log.Errorw("Failed to parse visual unit JSON", "error", err, "response", text[:min(500, len(text))], "task_id", taskID)
		if updateErr := s.taskService.UpdateTaskError(taskID, fmt.Errorf("Failed to parse visual unit result: %w", err)); updateErr != nil {
			s.log.Errorw("Failed to update task error", "error", updateErr, "task_id", taskID)
		}
		return
	}

	// Re-number shots sequentially
	for i := range shots {
		shots[i].ShotID = i + 1
	}

	// Post-validation for structured mode: warn if shot count mismatches
	if isStructured && len(structuredBlocks) > 0 {
		expected := len(structuredBlocks)
		actual := len(shots)
		if actual != expected {
			s.log.Warnw("Structured mode: AI output count mismatch — expected from // SHOT markers",
				"expected", expected, "actual", actual, "task_id", taskID)
		}

		// Wire user-specified metadata from ShotBlocks into AI output
		for i := range shots {
			if i >= len(structuredBlocks) {
				break
			}
			block := structuredBlocks[i]

			// Override with user-specified values if provided
			if block.Duration > 0 {
				shots[i].EstimatedDuration = block.Duration
			}
			if block.ShotType != "" {
				shots[i].ShotType = block.ShotType
			}
			if block.AudioMode != "" {
				shots[i].AudioMode = block.AudioMode
			}

			// Use block RawContent as script_segment if AI didn't capture it correctly
			if block.RawContent != "" && (shots[i].ScriptSegment == "" || shots[i].ScriptSegment == "null") {
				shots[i].ScriptSegment = block.RawContent
			}

			// Wire expanded tags from parsed lines into shot fields
			for _, seg := range block.Lines {
				switch seg.Type {
				case "sfx":
					if seg.Text != "" {
						if shots[i].SoundEffect == "" {
							shots[i].SoundEffect = seg.Text
						} else {
							shots[i].SoundEffect += "; " + seg.Text
						}
					}
				case "bgm":
					if seg.Text != "" {
						shots[i].BgmPrompt = seg.Text
					}
				case "camera":
					if seg.Text != "" {
						// Camera tag overrides movement field
						shots[i].Movement = seg.Text
					}
				case "vfx":
					if seg.Text != "" {
						// Append VFX instructions to visual description
						if shots[i].VisualDescription != "" {
							shots[i].VisualDescription += " [VFX: " + seg.Text + "]"
						} else {
							shots[i].VisualDescription = "[VFX: " + seg.Text + "]"
						}
					}
				case "note":
					if seg.Text != "" {
						// Append director's note to reason_for_shot
						if shots[i].ReasonForShot != "" {
							shots[i].ReasonForShot += " | Director note: " + seg.Text
						} else {
							shots[i].ReasonForShot = "Director note: " + seg.Text
						}
					}
				case "dialogue":
					// Override dialogue_text with user-specified dialogue
					if seg.Text != "" {
						if shots[i].DialogueText == "" {
							shots[i].DialogueText = seg.Text
						} else {
							shots[i].DialogueText += "\n" + seg.Text
						}
					}
				case "crowd":
					if seg.Text != "" {
						shots[i].DialogueText = seg.Text
						shots[i].DialogueType = "crowd"
					}
				}
			}
		}
	}

	if err := s.taskService.UpdateTaskStatus(taskID, "processing", 50, "Visual unit shots generated, parsing data..."); err != nil {
		s.log.Errorw("Failed to update task status", "error", err, "task_id", taskID)
		return
	}

	// Calculate total duration
	totalDuration := 0
	for _, shot := range shots {
		totalDuration += shot.EstimatedDuration
	}

	s.log.Infow("Visual unit storyboard generated",
		"task_id", taskID,
		"episode_id", episodeID,
		"count", len(shots),
		"total_duration_seconds", totalDuration)

	if err := s.taskService.UpdateTaskStatus(taskID, "processing", 70, "Saving visual unit shots..."); err != nil {
		s.log.Errorw("Failed to update task status", "error", err, "task_id", taskID)
		return
	}

	// Save to database
	if err := s.saveVoiceoverShots(episodeID, shots); err != nil {
		s.log.Errorw("Failed to save voiceover shots", "error", err, "task_id", taskID)
		if updateErr := s.taskService.UpdateTaskError(taskID, fmt.Errorf("Failed to save shots: %w", err)); updateErr != nil {
			s.log.Errorw("Failed to update task error", "error", updateErr, "task_id", taskID)
		}
		return
	}

	if err := s.taskService.UpdateTaskStatus(taskID, "processing", 90, "Updating episode duration..."); err != nil {
		s.log.Errorw("Failed to update task status", "error", err, "task_id", taskID)
		return
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
		"mode":             "visual_unit",
	}

	if err := s.taskService.UpdateTaskResult(taskID, resultData); err != nil {
		s.log.Errorw("Failed to update task result", "error", err, "task_id", taskID)
		return
	}

	s.log.Infow("Visual unit storyboard generation completed", "task_id", taskID, "episode_id", episodeID)
}

// saveVoiceoverShots maps VoiceoverShot[] to models.Storyboard and saves to DB
func (s *StoryboardService) saveVoiceoverShots(episodeID string, shots []VoiceoverShot) error {
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

		// Save each voiceover shot as a Storyboard
		for _, shot := range shots {
			// Build visual description as the Description field
			description := shot.VisualDescription

			// Generate image/video prompts from VoiceoverShot fields using existing helpers
			sbForPrompt := Storyboard{
				ShotNumber:  shot.ShotID,
				Title:       shot.Title,
				ShotType:    shot.ShotType,
				Angle:       shot.Angle,
				Movement:    shot.Movement,
				Location:    shot.Location,
				Time:        shot.Time,
				Atmosphere:  shot.Atmosphere,
				Action:      shot.VisualDescription, // Use visual description as action for prompts
				Dialogue:    shot.DialogueText,
				SoundEffect: shot.SoundEffect,
				Duration:    shot.EstimatedDuration,
			}

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

			imagePrompt := s.generateImagePrompt(sbForPrompt, propDescriptions)
			videoPrompt := s.generateVideoPrompt(sbForPrompt)

			// Convert string fields to pointers (nil if empty)
			strPtr := func(v string) *string {
				if v == "" {
					return nil
				}
				return &v
			}
			intPtr := func(v int) *int {
				return &v
			}
			boolPtr := func(v bool) *bool {
				return &v
			}

			// Convert triggered_rules to JSON
			var splitRulesJSON datatypes.JSON
			if len(shot.TriggeredRules) > 0 {
				if data, err := json.Marshal(shot.TriggeredRules); err == nil {
					splitRulesJSON = datatypes.JSON(data)
				}
			}

			storyboard := models.Storyboard{
				EpisodeID:        uint(epID),
				SceneID:          shot.SceneID,
				StoryboardNumber: shot.ShotID,
				Title:            strPtr(shot.Title),
				Location:         strPtr(shot.Location),
				Time:             strPtr(shot.Time),
				ShotType:         strPtr(shot.ShotType),
				Angle:            strPtr(shot.Angle),
				Movement:         strPtr(shot.Movement),
				Description:      &description,
				Action:           strPtr(shot.VisualDescription),
				Atmosphere:       strPtr(shot.Atmosphere),
				Dialogue:         strPtr(shot.DialogueText),
				ImagePrompt:      &imagePrompt,
				VideoPrompt:      &videoPrompt,
				VideoPromptSource: "auto",
				BgmPrompt:        strPtr(shot.BgmPrompt),
				SoundEffect:      strPtr(shot.SoundEffect),
				Duration:         shot.EstimatedDuration,
				// Voice-over Director fields
				ScriptSegment:   strPtr(shot.ScriptSegment),
				ScriptStartChar: intPtr(shot.ScriptStartChar),
				ScriptEndChar:   intPtr(shot.ScriptEndChar),
				ShotReason:      strPtr(shot.ReasonForShot),
				SplitRules:      splitRulesJSON,
				VisualType:      strPtr(shot.VisualType),
				ShotRole:        strPtr(shot.ShotRole),
				// Audio Strategy fields
				AudioMode:       strPtr(shot.AudioMode),
				NarratorEnabled: boolPtr(shot.NarratorEnabled),
				NarratorDucking: boolPtr(shot.NarratorDucking),
				DialogueType:    strPtr(shot.DialogueType),
				AmbienceType:    strPtr(shot.AmbienceType),
				AmbienceLevel:   strPtr(shot.AmbienceLevel),
				MusicMood:       strPtr(shot.MusicMood),
				MusicLevel:      strPtr(shot.MusicLevel),
			}

			if err := tx.Create(&storyboard).Error; err != nil {
				s.log.Errorw("Failed to create voiceover shot", "error", err, "shot_id", shot.ShotID)
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

		s.log.Infow("Voiceover shots saved successfully", "episode_id", episodeID, "count", len(shots))
		return nil
	})
}

// generateImagePrompt 生成专门用于图片生成的提示词（首帧静态画面）
func (s *StoryboardService) generateImagePrompt(sb Storyboard, propsDesc string) string {
	var parts []string

	// 0. 道具描述
	if propsDesc != "" {
		parts = append(parts, fmt.Sprintf("Props in scene: %s", propsDesc))
	}

	// 1. 完整的场景背景描述
	if sb.Location != "" {
		locationDesc := sb.Location
		if sb.Time != "" {
			locationDesc += ", " + sb.Time
		}
		parts = append(parts, locationDesc)
	}

	// 2. 角色初始静态姿态（去除动作过程，只保留起始状态）
	if sb.Action != "" {
		initialPose := extractInitialPose(sb.Action)
		if initialPose != "" {
			parts = append(parts, initialPose)
		}
	}

	// 3. 情绪氛围
	if sb.Emotion != "" {
		parts = append(parts, sb.Emotion)
	}

	// 4. 动漫风格
	parts = append(parts, "anime style, first frame")

	if len(parts) > 0 {
		return strings.Join(parts, ", ")
	}
	return "anime scene"
}

// extractInitialPose 提取初始静态姿态（去除动作过程）
func extractInitialPose(action string) string {
	// 去除动作过程关键词，保留初始状态描述
	processWords := []string{
		"然后", "接着", "接下来", "随后", "紧接着",
		"向下", "向上", "向前", "向后", "向左", "向右",
		"开始", "继续", "逐渐", "慢慢", "快速", "突然", "猛然",
	}

	result := action
	for _, word := range processWords {
		if idx := strings.Index(result, word); idx > 0 {
			// 在动作过程词之前截断
			result = result[:idx]
			break
		}
	}

	// 清理末尾标点
	result = strings.TrimRight(result, "，。,. ")
	return strings.TrimSpace(result)
}

// extractSimpleLocation 提取简化的场景地点（去除详细描述）
func extractSimpleLocation(location string) string {
	// 在"·"符号处截断，只保留主场景名称
	if idx := strings.Index(location, "·"); idx > 0 {
		return strings.TrimSpace(location[:idx])
	}

	// 如果有逗号，只保留第一部分
	if idx := strings.Index(location, "，"); idx > 0 {
		return strings.TrimSpace(location[:idx])
	}
	if idx := strings.Index(location, ","); idx > 0 {
		return strings.TrimSpace(location[:idx])
	}

	// 限制长度不超过15个字符
	maxLen := 15
	if len(location) > maxLen {
		return strings.TrimSpace(location[:maxLen])
	}

	return strings.TrimSpace(location)
}

// extractSimplePose 提取简单的核心姿态关键词（不超过10个字）
func extractSimplePose(action string) string {
	// 只提取前面最多10个字符作为核心姿态
	runes := []rune(action)
	maxLen := 10
	if len(runes) > maxLen {
		// 在标点符号处截断
		truncated := runes[:maxLen]
		for i := maxLen - 1; i >= 0; i-- {
			if truncated[i] == '，' || truncated[i] == '。' || truncated[i] == ',' || truncated[i] == '.' {
				truncated = runes[:i]
				break
			}
		}
		return strings.TrimSpace(string(truncated))
	}
	return strings.TrimSpace(action)
}

// extractFirstFramePose 从动作描述中提取首帧静态姿态
func extractFirstFramePose(action string) string {
	// 去除表示动作过程的关键词，保留初始状态
	processWords := []string{
		"然后", "接着", "向下", "向前", "走向", "冲向", "转身",
		"开始", "继续", "逐渐", "慢慢", "快速", "突然",
	}

	pose := action
	for _, word := range processWords {
		// 简单处理：在这些词之前截断
		if idx := strings.Index(pose, word); idx > 0 {
			pose = pose[:idx]
			break
		}
	}

	// 清理末尾标点
	pose = strings.TrimRight(pose, "，。,.")
	return strings.TrimSpace(pose)
}

// extractCompositionType 从镜头类型中提取构图类型（去除运镜）
func extractCompositionType(shotType string) string {
	// 去除运镜相关描述
	cameraMovements := []string{
		"晃动", "摇晃", "推进", "拉远", "跟随", "环绕",
		"运镜", "摄影", "移动", "旋转",
	}

	comp := shotType
	for _, movement := range cameraMovements {
		comp = strings.ReplaceAll(comp, movement, "")
	}

	// 清理多余的标点和空格
	comp = strings.ReplaceAll(comp, "··", "·")
	comp = strings.ReplaceAll(comp, "·", " ")
	comp = strings.TrimSpace(comp)

	return comp
}

// generateVideoPrompt 生成专门用于视频生成的提示词（包含运镜和动态元素）
// NOTE: BGM is intentionally excluded — video generation AI produces visuals, not audio.
// Including BGM can conflict with style DNA (e.g., styles that require ZERO music).
// Sound effects are included as environmental context to help AI understand the scene's physics.
func (s *StoryboardService) generateVideoPrompt(sb Storyboard) string {
	var parts []string
	videoRatio := "16:9"
	// 1. 人物动作（核心 - 定义shot内容）
	if sb.Action != "" {
		parts = append(parts, fmt.Sprintf("Action: %s", sb.Action))
	}



	// 2. 结果（动作的最终视觉状态 - 紧跟Action以保持叙事连贯）
	if sb.Result != "" {
		parts = append(parts, fmt.Sprintf("Result: %s", sb.Result))
	}

	// 3. 镜头运动（视频特有）
	if sb.Movement != "" {
		parts = append(parts, fmt.Sprintf("Camera movement: %s", sb.Movement))
	}

	// 4. 镜头类型和角度
	if sb.ShotType != "" {
		parts = append(parts, fmt.Sprintf("Shot type: %s", sb.ShotType))
	}
	if sb.Angle != "" {
		parts = append(parts, fmt.Sprintf("Camera angle: %s", sb.Angle))
	}

	// 5. 场景环境
	if sb.Location != "" {
		locationDesc := sb.Location
		if sb.Time != "" {
			locationDesc += ", " + sb.Time
		}
		parts = append(parts, fmt.Sprintf("Scene: %s", locationDesc))
	}

	// 6. 环境氛围
	if sb.Atmosphere != "" {
		parts = append(parts, fmt.Sprintf("Atmosphere: %s", sb.Atmosphere))
	}

	// 7. 对话与口型约束
	if sb.Dialogue != "" {
		parts = append(parts, fmt.Sprintf("Dialogue: %s", sb.Dialogue))
		
		// 自动判断如果是旁白/独白，则禁止嘴部动作；如果是对话，则要求说话
		dialogueLower := strings.ToLower(strings.TrimSpace(sb.Dialogue))
		isVoiceover := strings.HasPrefix(dialogueLower, "(vo)") || 
		   strings.HasPrefix(dialogueLower, "(monologue)") || 
		   strings.Contains(dialogueLower, "voiceover") || 
		   strings.HasPrefix(dialogueLower, "【旁白") || 
		   strings.HasPrefix(dialogueLower, "[旁白") || 
		   strings.HasPrefix(dialogueLower, "(narrator") ||
		   strings.Contains(dialogueLower, "（旁白）")
		
		if isVoiceover {
			parts = append(parts, "The character's mouth is strictly closed, silent expression, purely visual acting, no speaking, voiceover scene. --no talking, speaking, moving lips, open mouth, chatting")
		} else {
			parts = append(parts, "The character is actively speaking, lip-syncing naturally to the dialog, mouth moving")
		}
	} else {
		// 如果完全没有对话，要求保持嘴部闭合
		parts = append(parts, "The character's mouth is completely closed, silent scene. --no talking, speaking, moving lips")
	}

	// 8. 音效（作为环境物理上下文，帮助AI理解场景的物理特性）
	// BGM intentionally omitted - could conflict with style DNA (some styles require ZERO music)
	if sb.SoundEffect != "" {
		parts = append(parts, fmt.Sprintf("Sound effects: %s", sb.SoundEffect))
	}

	// 9. 视频比例
	parts = append(parts, fmt.Sprintf("=VideoRatio: %s", videoRatio))
	if len(parts) > 0 {
		return strings.Join(parts, ". ")
	}
	return "Cinematic video scene"
}

func (s *StoryboardService) saveStoryboards(episodeID string, storyboards []Storyboard) error {
	// 验证 episodeID
	epID, err := strconv.ParseUint(episodeID, 10, 32)
	if err != nil {
		s.log.Errorw("Invalid episode ID", "episode_id", episodeID, "error", err)
		return fmt.Errorf("无效的章节ID: %s", episodeID)
	}

	// 防御性检查：如果AI返回的分镜数量为0，不应该删除旧分镜
	if len(storyboards) == 0 {
		s.log.Errorw("AI返回的分镜数量为0，拒绝保存以避免删除现有分镜", "episode_id", episodeID)
		return fmt.Errorf("AI生成分镜失败：返回的分镜数量为0")
	}

	s.log.Infow("开始保存分镜头",
		"episode_id", episodeID,
		"episode_id_uint", uint(epID),
		"storyboard_count", len(storyboards))

	// 开启事务
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 验证该章节是否存在
		var episode models.Episode
		if err := tx.First(&episode, epID).Error; err != nil {
			s.log.Errorw("Episode not found", "episode_id", episodeID, "error", err)
			return fmt.Errorf("章节不存在: %s", episodeID)
		}

		s.log.Infow("找到章节信息",
			"episode_id", episode.ID,
			"episode_number", episode.EpisodeNum,
			"drama_id", episode.DramaID,
			"title", episode.Title)

		// 获取该剧集所有的分镜ID（使用 uint 类型）
		var storyboardIDs []uint
		if err := tx.Model(&models.Storyboard{}).
			Where("episode_id = ?", uint(epID)).
			Pluck("id", &storyboardIDs).Error; err != nil {
			return err
		}

		s.log.Infow("查询到现有分镜",
			"episode_id_string", episodeID,
			"episode_id_uint", uint(epID),
			"existing_storyboard_count", len(storyboardIDs),
			"storyboard_ids", storyboardIDs)

		// 如果有分镜，先清理关联的image_generations的storyboard_id
		if len(storyboardIDs) > 0 {
			if err := tx.Model(&models.ImageGeneration{}).
				Where("storyboard_id IN ?", storyboardIDs).
				Update("storyboard_id", nil).Error; err != nil {
				return err
			}
			s.log.Infow("已清理关联的图片生成记录", "count", len(storyboardIDs))
		}

		// 删除该剧集已有的分镜头（使用 uint 类型确保类型匹配）
		s.log.Warnw("准备删除分镜数据",
			"episode_id_string", episodeID,
			"episode_id_uint", uint(epID),
			"episode_id_from_db", episode.ID,
			"will_delete_count", len(storyboardIDs))

		result := tx.Where("episode_id = ?", uint(epID)).Delete(&models.Storyboard{})
		if result.Error != nil {
			s.log.Errorw("删除旧分镜失败", "episode_id", uint(epID), "error", result.Error)
			return result.Error
		}

		s.log.Infow("已删除旧分镜头",
			"episode_id", uint(epID),
			"deleted_count", result.RowsAffected)

		// 注意：不删除背景，因为背景是在分镜拆解前就提取好的
		// AI会直接返回scene_id，不需要在这里做字符串匹配

		// 保存新的分镜头
		for _, sb := range storyboards {
			// 构建描述信息，包含对话
			description := fmt.Sprintf("【镜头类型】%s\n【运镜】%s\n【动作】%s\n【对话】%s\n【结果】%s\n【情绪】%s",
				sb.ShotType, sb.Movement, sb.Action, sb.Dialogue, sb.Result, sb.Emotion)

			// 取出道具名称和描述用于 Prompt 生成
			var propDescriptions string
			var loadedProps []models.Prop
			if len(sb.Props) > 0 {
				if err := tx.Where("id IN ?", sb.Props).Find(&loadedProps).Error; err == nil {
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

			// 生成两种专用提示词
			imagePrompt := s.generateImagePrompt(sb, propDescriptions) // 专用于图片生成
			videoPrompt := s.generateVideoPrompt(sb) // 专用于视频生成

			// 处理 dialogue 字段
			var dialoguePtr *string
			if sb.Dialogue != "" {
				dialoguePtr = &sb.Dialogue
			}

			// 使用AI直接返回的SceneID
			if sb.SceneID != nil {
				s.log.Infow("Background ID from AI",
					"shot_number", sb.ShotNumber,
					"scene_id", *sb.SceneID)
			}

			// 处理 title 字段
			var titlePtr *string
			if sb.Title != "" {
				titlePtr = &sb.Title
			}

			// 处理shot_type、angle、movement字段
			var shotTypePtr, anglePtr, movementPtr *string
			if sb.ShotType != "" {
				shotTypePtr = &sb.ShotType
			}
			if sb.Angle != "" {
				anglePtr = &sb.Angle
			}
			if sb.Movement != "" {
				movementPtr = &sb.Movement
			}

			// 处理bgm_prompt、sound_effect字段
			var bgmPromptPtr, soundEffectPtr *string
			if sb.BgmPrompt != "" {
				bgmPromptPtr = &sb.BgmPrompt
			}
			if sb.SoundEffect != "" {
				soundEffectPtr = &sb.SoundEffect
			}

			// 处理result、atmosphere字段
			var resultPtr, atmospherePtr *string
			if sb.Result != "" {
				resultPtr = &sb.Result
			}
			if sb.Atmosphere != "" {
				atmospherePtr = &sb.Atmosphere
			}

			scene := models.Storyboard{
				EpisodeID:        uint(epID),
				SceneID:          sb.SceneID,
				StoryboardNumber: sb.ShotNumber,
				Title:            titlePtr,
				Location:         &sb.Location,
				Time:             &sb.Time,
				ShotType:         shotTypePtr,
				Angle:            anglePtr,
				Movement:         movementPtr,
				Description:      &description,
				Action:           &sb.Action,
				Result:           resultPtr,
				Atmosphere:       atmospherePtr,
				Dialogue:         dialoguePtr,
				ImagePrompt:      &imagePrompt,
				VideoPrompt:      &videoPrompt,
				VideoPromptSource: "auto",
				BgmPrompt:        bgmPromptPtr,
				SoundEffect:      soundEffectPtr,
				Duration:         sb.Duration,
			}

			if err := tx.Create(&scene).Error; err != nil {
				s.log.Errorw("Failed to create scene", "error", err, "shot_number", sb.ShotNumber)
				return err
			}

			// 关联角色
			if len(sb.Characters) > 0 {
				var characters []models.Character
				if err := tx.Where("id IN ?", sb.Characters).Find(&characters).Error; err != nil {
					s.log.Warnw("Failed to load characters for association", "error", err, "character_ids", sb.Characters)
				} else if len(characters) > 0 {
					if err := tx.Model(&scene).Association("Characters").Append(characters); err != nil {
						s.log.Warnw("Failed to associate characters", "error", err, "shot_number", sb.ShotNumber)
					} else {
						s.log.Infow("Characters associated successfully",
							"shot_number", sb.ShotNumber,
							"character_ids", sb.Characters,
							"count", len(characters))
					}
				}
			}

			// 关联道具
			if len(loadedProps) > 0 {
				if err := tx.Model(&scene).Association("Props").Append(loadedProps); err != nil {
					s.log.Warnw("Failed to associate props", "error", err, "shot_number", sb.ShotNumber)
				} else {
					s.log.Infow("Props associated successfully",
						"shot_number", sb.ShotNumber,
						"prop_ids", sb.Props,
						"count", len(loadedProps))
				}
			}
		}

		s.log.Infow("Storyboards saved successfully", "episode_id", episodeID, "count", len(storyboards))
		return nil
	})
}

// CreateStoryboardRequest 创建分镜请求
type CreateStoryboardRequest struct {
	EpisodeID        uint    `json:"episode_id"`
	SceneID          *uint   `json:"scene_id"`
	StoryboardNumber int     `json:"storyboard_number"`
	Title            *string `json:"title"`
	Location         *string `json:"location"`
	Time             *string `json:"time"`
	ShotType         *string `json:"shot_type"`
	Angle            *string `json:"angle"`
	Movement         *string `json:"movement"`
	Description      *string `json:"description"`
	Action           *string `json:"action"`
	Result           *string `json:"result"`
	Atmosphere       *string `json:"atmosphere"`
	Dialogue         *string `json:"dialogue"`
	BgmPrompt        *string `json:"bgm_prompt"`
	SoundEffect      *string `json:"sound_effect"`
	Duration         int     `json:"duration"`
	Characters       []uint  `json:"characters"`
}

// CreateStoryboard 创建单个分镜
func (s *StoryboardService) CreateStoryboard(req *CreateStoryboardRequest) (*models.Storyboard, error) {
	// 构建Storyboard对象
	sb := Storyboard{
		ShotNumber:  req.StoryboardNumber,
		ShotType:    getString(req.ShotType),
		Angle:       getString(req.Angle),
		Time:        getString(req.Time),
		Location:    getString(req.Location),
		SceneID:     req.SceneID,
		Movement:    getString(req.Movement),
		Action:      getString(req.Action),
		Dialogue:    getString(req.Dialogue),
		Result:      getString(req.Result),
		Atmosphere:  getString(req.Atmosphere),
		Emotion:     "", // 可以后续添加
		Duration:    req.Duration,
		BgmPrompt:   getString(req.BgmPrompt),
		SoundEffect: getString(req.SoundEffect),
		Characters:  req.Characters,
	}
	if req.Title != nil {
		sb.Title = *req.Title
	}

	// 生成提示词
	imagePrompt := s.generateImagePrompt(sb, "")
	videoPrompt := s.generateVideoPrompt(sb)

	// 构建 description
	desc := ""
	if req.Description != nil {
		desc = *req.Description
	}

	modelSB := &models.Storyboard{
		EpisodeID:        req.EpisodeID,
		SceneID:          req.SceneID,
		StoryboardNumber: req.StoryboardNumber,
		Title:            req.Title,
		Location:         req.Location,
		Time:             req.Time,
		ShotType:         req.ShotType,
		Angle:            req.Angle,
		Movement:         req.Movement,
		Description:      &desc,
		Action:           req.Action,
		Result:           req.Result,
		Atmosphere:       req.Atmosphere,
		Dialogue:         req.Dialogue,
		ImagePrompt:      &imagePrompt,
		VideoPrompt:      &videoPrompt,
		VideoPromptSource: "auto",
		BgmPrompt:        req.BgmPrompt,
		SoundEffect:      req.SoundEffect,
		Duration:         req.Duration,
	}

	if err := s.db.Create(modelSB).Error; err != nil {
		return nil, fmt.Errorf("failed to create storyboard: %w", err)
	}

	// 关联角色
	if len(req.Characters) > 0 {
		var characters []models.Character
		if err := s.db.Where("id IN ?", req.Characters).Find(&characters).Error; err != nil {
			s.log.Warnw("Failed to find characters for new storyboard", "error", err)
		} else if len(characters) > 0 {
			s.db.Model(modelSB).Association("Characters").Append(characters)
		}
	}

	s.log.Infow("Storyboard created", "id", modelSB.ID, "episode_id", req.EpisodeID)
	return modelSB, nil
}

// DeleteStoryboard 删除分镜
func (s *StoryboardService) DeleteStoryboard(storyboardID uint) error {
	result := s.db.Where("id = ? ", storyboardID).Delete(&models.Storyboard{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("storyboard not found")
	}
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func getString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
