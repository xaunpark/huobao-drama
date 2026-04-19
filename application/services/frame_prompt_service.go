package services

import (
	"fmt"
	"strings"

	"github.com/drama-generator/backend/application/prompts/fixed"
	"github.com/drama-generator/backend/domain/models"
	"github.com/drama-generator/backend/pkg/config"
	"github.com/drama-generator/backend/pkg/logger"
	"gorm.io/gorm"
)

// FramePromptService 处理帧提示词生成
type FramePromptService struct {
	db          *gorm.DB
	aiService   *AIService
	log         *logger.Logger
	config      *config.Config
	promptI18n  *PromptI18n
	taskService *TaskService
}

// NewFramePromptService 创建帧提示词服务
func NewFramePromptService(db *gorm.DB, cfg *config.Config, log *logger.Logger) *FramePromptService {
	return &FramePromptService{
		db:          db,
		aiService:   NewAIService(db, log),
		log:         log,
		config:      cfg,
		promptI18n:  NewPromptI18n(cfg),
		taskService: NewTaskService(db, log),
	}
}

// resolveStyleForShot returns the per-shot distilled image_style if available,
// otherwise falls back to the drama-level style key (used by resolveEffectiveStyle).
// This ensures that distilled, shot-specific styles take priority over the full
// project-level style_prompt template, reducing style bleeding across shots.
// See: plans/shot-style-distill.md (Phase 3, Task 10)
func resolveStyleForShot(sb models.Storyboard, dramaStyle string) string {
	if sb.ImageStyle != nil && *sb.ImageStyle != "" {
		return *sb.ImageStyle
	}
	return dramaStyle
}

// FrameType 帧类型
type FrameType string

const (
	FrameTypeFirst  FrameType = "first"  // 首帧
	FrameTypeKey    FrameType = "key"    // 关键帧
	FrameTypeLast   FrameType = "last"   // 尾帧
	FrameTypePanel  FrameType = "panel"  // 分镜板（3格组合）
	FrameTypeAction FrameType = "action" // 动作序列（5格）
	FrameTypeVideo  FrameType = "video"  // 视频 prompt (R2V)
)

// GenerateFramePromptRequest 生成帧提示词请求
type GenerateFramePromptRequest struct {
	StoryboardID string    `json:"storyboard_id"`
	FrameType    FrameType `json:"frame_type"`
	// 可选参数
	PanelCount int `json:"panel_count,omitempty"` // 分镜板格数，默认3
}

// FramePromptResponse 帧提示词响应
type FramePromptResponse struct {
	FrameType   FrameType          `json:"frame_type"`
	SingleFrame *SingleFramePrompt `json:"single_frame,omitempty"` // 单帧提示词
	MultiFrame  *MultiFramePrompt  `json:"multi_frame,omitempty"`  // 多帧提示词
}

// SingleFramePrompt 单帧提示词
type SingleFramePrompt struct {
	Prompt      string `json:"prompt"`
	Description string `json:"description"`
}

// MultiFramePrompt 多帧提示词
type MultiFramePrompt struct {
	Layout string              `json:"layout"` // horizontal_3, grid_2x2 等
	Frames []SingleFramePrompt `json:"frames"`
}

// GenerateFramePrompt 生成指定类型的帧提示词并保存到frame_prompts表
func (s *FramePromptService) GenerateFramePrompt(req GenerateFramePromptRequest, model string) (string, error) {
	// 查询分镜信息
	var storyboard models.Storyboard
	if err := s.db.Preload("Characters").First(&storyboard, req.StoryboardID).Error; err != nil {
		return "", fmt.Errorf("storyboard not found: %w", err)
	}

	// 创建任务
	task, err := s.taskService.CreateTask("frame_prompt_generation", req.StoryboardID)
	if err != nil {
		s.log.Errorw("Failed to create frame prompt generation task", "error", err, "storyboard_id", req.StoryboardID)
		return "", fmt.Errorf("创建任务失败: %w", err)
	}

	// 异步处理帧提示词生成
	go s.processFramePromptGeneration(task.ID, req, model)

	s.log.Infow("Frame prompt generation task created", "task_id", task.ID, "storyboard_id", req.StoryboardID, "frame_type", req.FrameType)
	return task.ID, nil
}

// processFramePromptGeneration 异步处理帧提示词生成
func (s *FramePromptService) processFramePromptGeneration(taskID string, req GenerateFramePromptRequest, model string) {
	// 更新任务状态为处理中
	s.taskService.UpdateTaskStatus(taskID, "processing", 0, "正在生成帧提示词...")

	// 查询分镜信息
	var storyboard models.Storyboard
	if err := s.db.Preload("Characters").First(&storyboard, req.StoryboardID).Error; err != nil {
		s.log.Errorw("Storyboard not found during frame prompt generation", "error", err, "storyboard_id", req.StoryboardID)
		s.taskService.UpdateTaskStatus(taskID, "failed", 0, "分镜信息不存在")
		return
	}

	// 获取场景信息
	var scene *models.Scene
	if storyboard.SceneID != nil {
		scene = &models.Scene{}
		if err := s.db.First(scene, *storyboard.SceneID).Error; err != nil {
			s.log.Warnw("Scene not found during frame prompt generation", "scene_id", *storyboard.SceneID, "task_id", taskID)
			scene = nil
		}
	}

	// 获取 drama 的 style 信息
	var episode models.Episode
	if err := s.db.Preload("Drama").First(&episode, storyboard.EpisodeID).Error; err != nil {
		s.log.Warnw("Failed to load episode and drama", "error", err, "episode_id", storyboard.EpisodeID)
	}
	dramaStyle := episode.Drama.Style

	response := &FramePromptResponse{
		FrameType: req.FrameType,
	}

	// 生成提示词
	var err error
	switch req.FrameType {
	case FrameTypeFirst:
		response.SingleFrame, err = s.generateFirstFrame(episode.Drama.ID, storyboard, scene, dramaStyle, model)
		if err != nil {
			s.taskService.UpdateTaskStatus(taskID, "failed", 0, "AI生成提示词失败: "+err.Error())
			return
		}
		// 保存单帧提示词
		s.saveFramePrompt(req.StoryboardID, string(req.FrameType), response.SingleFrame.Prompt, response.SingleFrame.Description, "")
	case FrameTypeKey:
		response.SingleFrame, err = s.generateKeyFrame(episode.Drama.ID, storyboard, scene, dramaStyle, model)
		if err != nil {
			s.taskService.UpdateTaskStatus(taskID, "failed", 0, "AI生成提示词失败: "+err.Error())
			return
		}
		s.saveFramePrompt(req.StoryboardID, string(req.FrameType), response.SingleFrame.Prompt, response.SingleFrame.Description, "")
	case FrameTypeLast:
		response.SingleFrame, err = s.generateLastFrame(episode.Drama.ID, storyboard, scene, dramaStyle, model)
		if err != nil {
			s.taskService.UpdateTaskStatus(taskID, "failed", 0, "AI生成提示词失败: "+err.Error())
			return
		}
		s.saveFramePrompt(req.StoryboardID, string(req.FrameType), response.SingleFrame.Prompt, response.SingleFrame.Description, "")
	case FrameTypePanel:
		count := req.PanelCount
		if count == 0 {
			count = 3
		}
		response.MultiFrame, err = s.generatePanelFrames(episode.Drama.ID, storyboard, scene, count, dramaStyle, model)
		if err != nil {
			s.taskService.UpdateTaskStatus(taskID, "failed", 0, "AI生成分镜板提示词失败: "+err.Error())
			return
		}
		// 保存多帧提示词（合并为一条记录）
		var prompts []string
		for _, frame := range response.MultiFrame.Frames {
			prompts = append(prompts, frame.Prompt)
		}
		combinedPrompt := strings.Join(prompts, "\n---\n")
		s.saveFramePrompt(req.StoryboardID, string(req.FrameType), combinedPrompt, "分镜板组合提示词", response.MultiFrame.Layout)
	case FrameTypeAction:
		response.MultiFrame, err = s.generateActionSequence(episode.Drama.ID, storyboard, scene, dramaStyle, model)
		if err != nil {
			s.taskService.UpdateTaskStatus(taskID, "failed", 0, "AI生成动作序列提示词失败: "+err.Error())
			return
		}
		var prompts []string
		for _, frame := range response.MultiFrame.Frames {
			prompts = append(prompts, frame.Prompt)
		}
		combinedPrompt := strings.Join(prompts, "\n---\n")
		s.saveFramePrompt(req.StoryboardID, string(req.FrameType), combinedPrompt, "动作序列组合提示词", response.MultiFrame.Layout)
	case FrameTypeVideo:
		response.SingleFrame, err = s.generateVideoPrompt(episode.Drama.ID, storyboard, scene, dramaStyle, episode.Drama.CustomStyle, model)
		if err != nil {
			s.taskService.UpdateTaskStatus(taskID, "failed", 0, "AI生成视频提示词失败: "+err.Error())
			return
		}
		// 保存视频提示词 (同步更新到 storyboard 表)
		s.saveFramePrompt(req.StoryboardID, string(req.FrameType), response.SingleFrame.Prompt, response.SingleFrame.Description, "")
		updates := map[string]interface{}{
			"video_prompt":        response.SingleFrame.Prompt,
			"video_prompt_source": "ai",
		}
		if err := s.db.Model(&models.Storyboard{}).Where("id = ?", req.StoryboardID).Updates(updates).Error; err != nil {
			s.log.Warnw("Failed to update storyboard video_prompt", "error", err, "storyboard_id", req.StoryboardID)
		}
	default:
		s.log.Errorw("Unsupported frame type during frame prompt generation", "frame_type", req.FrameType, "task_id", taskID)
		s.taskService.UpdateTaskStatus(taskID, "failed", 0, "不支持的帧类型")
		return
	}

	// 更新任务状态为完成
	s.taskService.UpdateTaskResult(taskID, map[string]interface{}{
		"response":      response,
		"storyboard_id": req.StoryboardID,
		"frame_type":    string(req.FrameType),
	})

	s.log.Infow("Frame prompt generation completed", "task_id", taskID, "storyboard_id", req.StoryboardID, "frame_type", req.FrameType)
}

// saveFramePrompt 保存帧提示词到数据库
func (s *FramePromptService) saveFramePrompt(storyboardID, frameType, prompt, description, layout string) {
	framePrompt := models.FramePrompt{
		StoryboardID: uint(mustParseUint(storyboardID)),
		FrameType:    frameType,
		Prompt:       prompt,
	}

	if description != "" {
		framePrompt.Description = &description
	}
	if layout != "" {
		framePrompt.Layout = &layout
	}

	// 先删除同类型的旧记录（保持最新）
	s.db.Where("storyboard_id = ? AND frame_type = ?", storyboardID, frameType).Delete(&models.FramePrompt{})

	// 插入新记录
	if err := s.db.Create(&framePrompt).Error; err != nil {
		s.log.Warnw("Failed to save frame prompt", "error", err, "storyboard_id", storyboardID, "frame_type", frameType)
	}
}

// mustParseUint 辅助函数
func mustParseUint(s string) uint64 {
	var result uint64
	fmt.Sscanf(s, "%d", &result)
	return result
}

// generateFirstFrame 生成首帧提示词
func (s *FramePromptService) generateFirstFrame(dramaID uint, sb models.Storyboard, scene *models.Scene, dramaStyle string, model string) (*SingleFramePrompt, error) {
	// 构建上下文信息
	contextInfo := s.buildStoryboardContext(sb, scene)

	// 使用国际化提示词 — 优先使用 shot-level distilled style
	shotStyle := resolveStyleForShot(sb, dramaStyle)
	dynamicPrompt := s.promptI18n.WithDramaFirstFramePrompt(dramaID, shotStyle)
	systemPrompt := dynamicPrompt + "\n\n" + fixed.Get("image_generation")
	userPrompt := s.promptI18n.FormatUserPrompt("frame_info", contextInfo)

	// 调用AI生成（如果指定了模型则使用指定的模型）
	var aiResponse string
	var err error
	if model != "" {
		client, getErr := s.aiService.GetAIClientForModel("text", model)
		if getErr != nil {
			s.log.Warnw("Failed to get client for specified model, using default", "model", model, "error", getErr)
			aiResponse, err = s.aiService.GenerateText(userPrompt, systemPrompt)
		} else {
			aiResponse, err = client.GenerateText(userPrompt, systemPrompt)
		}
	} else {
		aiResponse, err = s.aiService.GenerateText(userPrompt, systemPrompt)
	}
	if err != nil {
		s.log.Warnw("AI generation failed", "error", err)
		return nil, err
	}

	// 解析AI返回的JSON
	result := s.parseFramePromptJSON(aiResponse)
	if result == nil {
		s.log.Warnw("Failed to parse AI JSON response", "storyboard_id", sb.ID, "response", aiResponse)
		return nil, fmt.Errorf("解析AI结果失败")
	}

	return result, nil
}

// generateKeyFrame 生成关键帧提示词
func (s *FramePromptService) generateKeyFrame(dramaID uint, sb models.Storyboard, scene *models.Scene, dramaStyle string, model string) (*SingleFramePrompt, error) {
	// 构建上下文信息
	contextInfo := s.buildStoryboardContext(sb, scene)

	// 使用国际化提示词 — 优先使用 shot-level distilled style
	shotStyle := resolveStyleForShot(sb, dramaStyle)
	dynamicPrompt := s.promptI18n.WithDramaKeyFramePrompt(dramaID, shotStyle)
	systemPrompt := dynamicPrompt + "\n\n" + fixed.Get("image_generation")
	userPrompt := s.promptI18n.FormatUserPrompt("key_frame_info", contextInfo)

	// 调用AI生成
	var aiResponse string
	var err error
	if model != "" {
		client, getErr := s.aiService.GetAIClientForModel("text", model)
		if getErr != nil {
			s.log.Warnw("Failed to get client for specified model, using default", "model", model, "error", getErr)
			aiResponse, err = s.aiService.GenerateText(userPrompt, systemPrompt)
		} else {
			aiResponse, err = client.GenerateText(userPrompt, systemPrompt)
		}
	} else {
		aiResponse, err = s.aiService.GenerateText(userPrompt, systemPrompt)
	}
	if err != nil {
		s.log.Warnw("AI generation failed", "error", err)
		return nil, err
	}

	// 解析AI返回的JSON
	result := s.parseFramePromptJSON(aiResponse)
	if result == nil {
		s.log.Warnw("Failed to parse AI JSON response", "storyboard_id", sb.ID, "response", aiResponse)
		return nil, fmt.Errorf("解析AI结果失败")
	}

	return result, nil
}

// generateLastFrame 生成尾帧提示词
func (s *FramePromptService) generateLastFrame(dramaID uint, sb models.Storyboard, scene *models.Scene, dramaStyle string, model string) (*SingleFramePrompt, error) {
	// 构建上下文信息
	contextInfo := s.buildStoryboardContext(sb, scene)

	// 使用国际化提示词 — 优先使用 shot-level distilled style
	shotStyle := resolveStyleForShot(sb, dramaStyle)
	dynamicPrompt := s.promptI18n.WithDramaLastFramePrompt(dramaID, shotStyle)
	systemPrompt := dynamicPrompt + "\n\n" + fixed.Get("image_generation")
	userPrompt := s.promptI18n.FormatUserPrompt("last_frame_info", contextInfo)

	// 调用AI生成
	var aiResponse string
	var err error
	if model != "" {
		client, getErr := s.aiService.GetAIClientForModel("text", model)
		if getErr != nil {
			s.log.Warnw("Failed to get client for specified model, using default", "model", model, "error", getErr)
			aiResponse, err = s.aiService.GenerateText(userPrompt, systemPrompt)
		} else {
			aiResponse, err = client.GenerateText(userPrompt, systemPrompt)
		}
	} else {
		aiResponse, err = s.aiService.GenerateText(userPrompt, systemPrompt)
	}
	if err != nil {
		s.log.Warnw("AI generation failed", "error", err)
		return nil, err
	}

	// 解析AI返回的JSON
	result := s.parseFramePromptJSON(aiResponse)
	if result == nil {
		s.log.Warnw("Failed to parse AI JSON response", "storyboard_id", sb.ID, "response", aiResponse)
		return nil, fmt.Errorf("解析AI结果失败")
	}

	return result, nil
}

// generatePanelFrames 生成分镜板提示词（多格组合）
func (s *FramePromptService) generatePanelFrames(dramaID uint, sb models.Storyboard, scene *models.Scene, count int, dramaStyle string, model string) (*MultiFramePrompt, error) {
	layout := fmt.Sprintf("horizontal_%d", count)

	frames := make([]SingleFramePrompt, count)

	// 固定生成：首帧 -> 关键帧 -> 尾帧
	if count == 3 {
		f1, err1 := s.generateFirstFrame(dramaID, sb, scene, dramaStyle, model)
		if err1 != nil { return nil, err1 }
		frames[0] = *f1
		frames[0].Description = "第1格：初始状态"

		f2, err2 := s.generateKeyFrame(dramaID, sb, scene, dramaStyle, model)
		if err2 != nil { return nil, err2 }
		frames[1] = *f2
		frames[1].Description = "第2格：动作高潮"

		f3, err3 := s.generateLastFrame(dramaID, sb, scene, dramaStyle, model)
		if err3 != nil { return nil, err3 }
		frames[2] = *f3
		frames[2].Description = "第3格：最终状态"
	} else if count == 4 {
		f1, err1 := s.generateFirstFrame(dramaID, sb, scene, dramaStyle, model)
		if err1 != nil { return nil, err1 }
		frames[0] = *f1
		f2, err2 := s.generateKeyFrame(dramaID, sb, scene, dramaStyle, model)
		if err2 != nil { return nil, err2 }
		frames[1] = *f2
		f3, err3 := s.generateKeyFrame(dramaID, sb, scene, dramaStyle, model)
		if err3 != nil { return nil, err3 }
		frames[2] = *f3
		f4, err4 := s.generateLastFrame(dramaID, sb, scene, dramaStyle, model)
		if err4 != nil { return nil, err4 }
		frames[3] = *f4
	}

	return &MultiFramePrompt{
		Layout: layout,
		Frames: frames,
	}, nil
}

// generateActionSequence 生成动作序列提示词（3x3宫格）
func (s *FramePromptService) generateActionSequence(dramaID uint, sb models.Storyboard, scene *models.Scene, dramaStyle string, model string) (*MultiFramePrompt, error) {
	// 构建上下文信息
	contextInfo := s.buildStoryboardContext(sb, scene)

	// 使用国际化提示词 - 根据 pacing_mode 选择不同的提示词 — 优先使用 shot-level distilled style
	shotStyle := resolveStyleForShot(sb, dramaStyle)
	var dynamicPrompt string
	if sb.PacingMode != nil && *sb.PacingMode == "rapid_cut" {
		// Rapid cut mode: 3 panels = 3 distinct micro-shots
		dynamicPrompt = s.promptI18n.WithDramaRapidCutActionSequenceFramePrompt(dramaID, shotStyle)
		s.log.Infow("Using rapid cut action sequence prompt",
			"storyboard_id", sb.ID,
			"pacing_mode", *sb.PacingMode)
	} else {
		// Standard mode: 3 panels = Start → Peak → End of one action
		dynamicPrompt = s.promptI18n.WithDramaActionSequenceFramePrompt(dramaID, shotStyle)
	}
	systemPrompt := dynamicPrompt + "\n\n" + fixed.Get("image_generation")
	userPrompt := s.promptI18n.FormatUserPrompt("frame_info", contextInfo)

	// 调用AI生成
	var aiResponse string
	var err error
	if model != "" {
		client, getErr := s.aiService.GetAIClientForModel("text", model)
		if getErr != nil {
			s.log.Warnw("Failed to get client for specified model, using default", "model", model, "error", getErr)
			aiResponse, err = s.aiService.GenerateText(userPrompt, systemPrompt)
		} else {
			aiResponse, err = client.GenerateText(userPrompt, systemPrompt)
		}
	} else {
		aiResponse, err = s.aiService.GenerateText(userPrompt, systemPrompt)
	}

	if err != nil {
		s.log.Warnw("AI generation failed for action sequence", "error", err)
		return nil, err
	}

	// 解析AI返回的JSON
	result := s.parseFramePromptJSON(aiResponse)
	if result == nil {
		s.log.Warnw("Failed to parse AI JSON response for action sequence", "storyboard_id", sb.ID, "response", aiResponse)
		return nil, fmt.Errorf("解析AI结果失败")
	}

	// 动作序列是一个整体的1x3横向条带图片（Start → Peak → End），只返回一个prompt
	return &MultiFramePrompt{
		Layout: "horizontal_strip_3",
		Frames: []SingleFramePrompt{*result},
	}, nil
}

// buildStoryboardContext 构建镜头上下文信息
// Action and Result are placed first as they define the core of the shot.
// A hard constraints summary is appended to ensure AI maintains shot fidelity.
func (s *FramePromptService) buildStoryboardContext(sb models.Storyboard, scene *models.Scene) string {
	var parts []string

	// === PRIMARY: Action & Result define WHAT the shot IS ===
	// 动作（核心 - 必须第一位）
	if sb.Action != nil && *sb.Action != "" {
		parts = append(parts, s.promptI18n.FormatUserPrompt("action_label", *sb.Action))
	}

	// 结果（动作的最终状态 - Panel 3/End Frame 应该描绘这个结果）
	if sb.Result != nil && *sb.Result != "" {
		parts = append(parts, s.promptI18n.FormatUserPrompt("result_label", *sb.Result))
	}

	// === SECONDARY: Camera & composition define HOW the shot looks ===
	// 镜头参数
	if sb.ShotType != nil {
		parts = append(parts, s.promptI18n.FormatUserPrompt("shot_type_label", *sb.ShotType))
	}
	if sb.Angle != nil {
		parts = append(parts, s.promptI18n.FormatUserPrompt("angle_label", *sb.Angle))
	}
	if sb.Movement != nil {
		parts = append(parts, s.promptI18n.FormatUserPrompt("movement_label", *sb.Movement))
	}

	// === CONTEXT: Characters first (subject identity), then scene, atmosphere ===
	// 角色/主体 — 放在场景之前，避免场景名称语义覆盖主体身份
	if len(sb.Characters) > 0 {
		var charNames []string
		for _, char := range sb.Characters {
			name := char.Name
			if char.Appearance != nil && *char.Appearance != "" {
				name = fmt.Sprintf("%s (%s)", char.Name, *char.Appearance)
			}
			charNames = append(charNames, name)
		}
		parts = append(parts, s.promptI18n.FormatUserPrompt("characters_label", strings.Join(charNames, "; ")))
	}

	// 场景信息
	if scene != nil {
		parts = append(parts, s.promptI18n.FormatUserPrompt("scene_label", scene.Location, scene.Time))
	} else if sb.Location != nil && sb.Time != nil {
		parts = append(parts, s.promptI18n.FormatUserPrompt("scene_label", *sb.Location, *sb.Time))
	}

	// 氛围
	if sb.Atmosphere != nil && *sb.Atmosphere != "" {
		parts = append(parts, s.promptI18n.FormatUserPrompt("atmosphere_label", *sb.Atmosphere))
	}

	// 对白
	if sb.Dialogue != nil && *sb.Dialogue != "" {
		parts = append(parts, s.promptI18n.FormatUserPrompt("dialogue_label", *sb.Dialogue))
	}

	// 镜头描述（补充信息）
	if sb.Description != nil && *sb.Description != "" {
		parts = append(parts, s.promptI18n.FormatUserPrompt("shot_description_label", *sb.Description))
	}

	// === HARD CONSTRAINTS REMINDER ===
	// Append constraint summary to reinforce shot fidelity
	var constraints []string
	constraints = append(constraints, "\n**[HARD CONSTRAINTS - DO NOT VIOLATE]**")
	constraints = append(constraints, "- You MUST depict ONLY the action described above. Do NOT invent additional events, story beats, or dramatic escalation.")
	if sb.Action != nil && *sb.Action != "" {
		constraints = append(constraints, "- The visual focus/hero subject of this shot is defined by the Action field. Maintain this focus across ALL frames.")
	}
	if sb.ShotType != nil || sb.Angle != nil || sb.Movement != nil {
		constraints = append(constraints, "- Camera angle, shot type, and movement MUST remain consistent across all frames as specified above.")
	}
	constraints = append(constraints, "- If the Action describes FG/MG/BG layering, maintain that EXACT layering in every frame.")
	constraints = append(constraints, "- Frame 9 MUST depict the Result described above — this is the end state of the action.")
	if len(sb.Characters) > 0 {
		constraints = append(constraints, "- Do not change or substitute the subject/characters listed above.")
	}

	parts = append(parts, strings.Join(constraints, "\n"))

	return strings.Join(parts, "\n")
}

// buildFallbackPrompt 构建降级提示词（AI失败时使用）
func (s *FramePromptService) buildFallbackPrompt(sb models.Storyboard, scene *models.Scene, suffix string) string {
	var parts []string

	// 场景
	if scene != nil {
		parts = append(parts, fmt.Sprintf("%s, %s", scene.Location, scene.Time))
	}

	// 角色
	if len(sb.Characters) > 0 {
		for _, char := range sb.Characters {
			parts = append(parts, char.Name)
		}
	}

	// 氛围
	if sb.Atmosphere != nil {
		parts = append(parts, *sb.Atmosphere)
	}

	parts = append(parts, "anime style", suffix)
	return strings.Join(parts, ", ")
}

// generateVideoPrompt 生成视频提示词 (R2V)
func (s *FramePromptService) generateVideoPrompt(dramaID uint, sb models.Storyboard, scene *models.Scene, dramaStyle string, customStyle string, model string) (*SingleFramePrompt, error) {
	// 构建上下文信息
	contextInfo := s.buildStoryboardContext(sb, scene)

	// 使用国际化提示词
	dynamicPrompt := s.promptI18n.WithDramaVideoExtractionPrompt(dramaID, dramaStyle, customStyle)
	systemPrompt := dynamicPrompt + "\n\n" + fixed.Get("image_generation")
	userPrompt := s.promptI18n.FormatUserPrompt("frame_info", contextInfo)

	// 调用AI生成
	var aiResponse string
	var err error
	if model != "" {
		client, getErr := s.aiService.GetAIClientForModel("text", model)
		if getErr == nil {
			aiResponse, err = client.GenerateText(userPrompt, systemPrompt)
		} else {
			aiResponse, err = s.aiService.GenerateText(userPrompt, systemPrompt)
		}
	} else {
		aiResponse, err = s.aiService.GenerateText(userPrompt, systemPrompt)
	}

	if err != nil {
		s.log.Warnw("AI video prompt generation failed", "error", err)
		return nil, err
	}

	// 解析AI返回的JSON
	result := s.parseFramePromptJSON(aiResponse)
	if result == nil {
		s.log.Warnw("Failed to parse AI video prompt JSON response", "storyboard_id", sb.ID, "response", aiResponse)
		return nil, fmt.Errorf("解析AI结果失败")
	}

	return result, nil
}
