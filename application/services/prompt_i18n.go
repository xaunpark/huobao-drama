package services

import (
	"fmt"
	"strings"

	"github.com/drama-generator/backend/application/prompts"
	"github.com/drama-generator/backend/domain/models"
	"github.com/drama-generator/backend/pkg/config"
)

// PromptI18n 提示词国际化工具
type PromptI18n struct {
	config          *config.Config
	templateService *PromptTemplateService // optional, nil = always use defaults
}

var globalTemplateService *PromptTemplateService

// NewPromptI18n 创建提示词国际化工具
func NewPromptI18n(cfg *config.Config) *PromptI18n {
	return &PromptI18n{
		config:          cfg,
		templateService: globalTemplateService,
	}
}

// SetTemplateService 设置模板服务（可选注入）
func (p *PromptI18n) SetTemplateService(svc *PromptTemplateService) {
	p.templateService = svc
	if globalTemplateService == nil {
		globalTemplateService = svc
	}
}

// resolvePrompt 使用模板服务解析 prompt，如果没有模板服务则返回默认
func (p *PromptI18n) resolvePrompt(dramaID uint, promptType string) string {
	if p.templateService != nil && dramaID > 0 {
		return p.templateService.ResolvePrompt(dramaID, promptType)
	}
	// Fallback to default
	defaultFile, ok := models.PromptTypeToDefaultFile[promptType]
	if !ok {
		return ""
	}
	return prompts.Get(defaultFile)
}

// formatPromptWithVars 安全地替换 %s 变量 - 用于自定义模板可能没有 %s 占位符
func formatPromptWithVars(template string, vars map[string]string) string {
	result := template
	for placeholder, value := range vars {
		result = strings.Replace(result, placeholder, value, 1)
	}
	return result
}

// resolveEffectiveStyle determines the actual style description to use.
// Priority: template's style_prompt > custom style text > dropdown style key.
// This ensures that when a template has its own style_prompt, the %s placeholders
// in extraction/image prompts receive the template's style description instead of
// the short dropdown value like "ghibli".
func (p *PromptI18n) resolveEffectiveStyle(dramaID uint, style string, customStyle string) string {
	// 1. Custom style always wins (user explicitly typed it)
	if style == "custom" && customStyle != "" {
		return customStyle
	}

	// 2. Check if drama has a template with style_prompt
	if p.templateService != nil && dramaID > 0 {
		templateStylePrompt := p.templateService.ResolvePromptIfCustom(dramaID, "style_prompt")
		if templateStylePrompt != "" {
			return templateStylePrompt
		}
	}

	// 3. Fallback to dropdown style key
	return p.GetStylePrompt(style, customStyle)
}

// ResolveEffectiveStylePublic is the exported version. Other services can call
// this to get the resolved style description (template override > custom > key).
func (p *PromptI18n) ResolveEffectiveStylePublic(dramaID uint, style string, customStyle string) string {
	return p.resolveEffectiveStyle(dramaID, style, customStyle)
}

// GetLanguage 获取当前语言设置
func (p *PromptI18n) GetLanguage() string {
	lang := p.config.App.Language
	if lang == "" {
		return "zh" // 默认中文
	}
	return lang
}

// IsEnglish 判断是否为英文模式（动态读取配置）
func (p *PromptI18n) IsEnglish() bool {
	return p.GetLanguage() == "en"
}

// GetStoryboardSystemPrompt 获取分镜生成系统提示词
func (p *PromptI18n) GetStoryboardSystemPrompt() string {
	if true || p.IsEnglish() {
		return prompts.Get("storyboard_story_breakdown.txt")
	}

	return prompts.Get("storyboard_story_breakdown.txt")
}

// GetSceneExtractionPrompt 获取场景提取提示词
func (p *PromptI18n) GetSceneExtractionPrompt(style string, customStyle string) string {
	imageRatio := "16:9"
	effectiveStyle := style
	if style == "custom" && customStyle != "" {
		effectiveStyle = customStyle
	}
	if true || p.IsEnglish() {
		return fmt.Sprintf(prompts.Get("scene_extraction.txt"), effectiveStyle, imageRatio)
	}

	return fmt.Sprintf(prompts.Get("scene_extraction.txt"), effectiveStyle, imageRatio)
}

// GetFirstFramePrompt 获取首帧提示词
func (p *PromptI18n) GetFirstFramePrompt(style string) string {
	imageRatio := "16:9"
	if true || p.IsEnglish() {
		return fmt.Sprintf(prompts.Get("image_first_frame.txt"), style, imageRatio)
	}

	return fmt.Sprintf(prompts.Get("image_first_frame.txt"), style, imageRatio)
}

// GetKeyFramePrompt 获取关键帧提示词
func (p *PromptI18n) GetKeyFramePrompt(style string) string {
	imageRatio := "16:9"
	if true || p.IsEnglish() {
		return fmt.Sprintf(prompts.Get("image_key_frame.txt"), style, imageRatio)
	}

	return fmt.Sprintf(prompts.Get("image_key_frame.txt"), style, imageRatio)
}

// GetActionSequenceFramePrompt 获取动作序列提示词
func (p *PromptI18n) GetActionSequenceFramePrompt(style string) string {
	imageRatio := "16:9"
	if true || p.IsEnglish() {
		return fmt.Sprintf(prompts.Get("image_action_sequence.txt"), style, imageRatio)
	}

	return fmt.Sprintf(prompts.Get("image_action_sequence.txt"), style, imageRatio)
}

// GetRapidCutActionSequenceFramePrompt 获取rapid cut多场景动作序列提示词
func (p *PromptI18n) GetRapidCutActionSequenceFramePrompt(style string) string {
	imageRatio := "16:9"
	return fmt.Sprintf(prompts.Get("image_action_sequence_rapid_cut.txt"), style, imageRatio)
}

// GetLastFramePrompt 获取尾帧提示词
func (p *PromptI18n) GetLastFramePrompt(style string) string {
	imageRatio := "16:9"
	if true || p.IsEnglish() {
		return fmt.Sprintf(prompts.Get("image_last_frame.txt"), style, imageRatio)
	}

	return fmt.Sprintf(prompts.Get("image_last_frame.txt"), style, imageRatio)
}

// GetVideoExtractionPrompt 获取视频提取提示词
func (p *PromptI18n) GetVideoExtractionPrompt(style string) string {
	if true || p.IsEnglish() {
		return fmt.Sprintf(prompts.Get("video_extraction.txt"), style)
	}

	return fmt.Sprintf(prompts.Get("video_extraction.txt"), style)
}

// GetOutlineGenerationPrompt 获取大纲生成提示词
func (p *PromptI18n) GetOutlineGenerationPrompt() string {
	if true || p.IsEnglish() {
		return prompts.Get("script_outline_generation.txt")
	}

	return prompts.Get("script_outline_generation.txt")
}

// GetCharacterExtractionPrompt 获取角色提取提示词
func (p *PromptI18n) GetCharacterExtractionPrompt(style string, customStyle string) string {
	imageRatio := "16:9"
	effectiveStyle := style
	if style == "custom" && customStyle != "" {
		effectiveStyle = customStyle
	}
	if true || p.IsEnglish() {
		return fmt.Sprintf(prompts.Get("character_extraction.txt"), effectiveStyle, imageRatio)
	}

	return fmt.Sprintf(prompts.Get("character_extraction.txt"), effectiveStyle, imageRatio)
}

// GetPropExtractionPrompt 获取道具提取提示词
func (p *PromptI18n) GetPropExtractionPrompt(style string, customStyle string) string {
	imageRatio := "1:1"
	effectiveStyle := style
	if style == "custom" && customStyle != "" {
		effectiveStyle = customStyle
	}

	if true || p.IsEnglish() {
		return fmt.Sprintf(prompts.Get("prop_extraction.txt"), effectiveStyle, imageRatio)
	}

	return fmt.Sprintf(prompts.Get("prop_extraction.txt"), effectiveStyle, imageRatio)
}

// GetEpisodeScriptPrompt 获取分集剧本生成提示词
func (p *PromptI18n) GetEpisodeScriptPrompt() string {
	if true || p.IsEnglish() {
		return prompts.Get("script_episode_generation.txt")
	}

	return prompts.Get("script_episode_generation.txt")
}

// FormatUserPrompt 格式化用户提示词的通用文本
func (p *PromptI18n) FormatUserPrompt(key string, args ...interface{}) string {
	templates := map[string]map[string]string{
		"en": {

			"outline_request":        "Please create a short drama outline for the following theme:\n\nTheme: %s",
			"genre_preference":       "\nGenre preference: %s",
			"style_requirement":      "\nStyle requirement: %s",
			"episode_count":          "\nNumber of episodes: %d episodes",
			"episode_importance":     "\n\n**Important: Must plan complete storylines for all %d episodes in the episodes array, each with clear story content!**",
			"character_request":      "Script content:\n%s\n\nPlease extract and organize detailed character profiles for all main characters from the script.",
			"episode_script_request": "Drama outline:\n%s\n%s\nPlease create detailed scripts for %d episodes based on the above outline and characters.\n\n**Important requirements:**\n- Must generate all %d episodes, from episode 1 to episode %d, cannot skip any\n- Each episode is about 3-5 minutes (150-300 seconds)\n- The duration field for each episode should be set reasonably based on script content length, not all the same value\n- The episodes array in the returned JSON must contain %d elements",
			"frame_info":             "Shot information:\n%s\n\nPlease directly generate the image prompt for the first frame without any explanation:",
			"key_frame_info":         "Shot information:\n%s\n\nPlease directly generate the image prompt for the key frame without any explanation:",
			"last_frame_info":        "Shot information:\n%s\n\nPlease directly generate the image prompt for the last frame without any explanation:",
			"script_content_label":   "【Script Content】",
			"storyboard_list_label":  "【Storyboard List】",
			"task_label":             "【Task】",
			"character_list_label":   "【Available Character List】",
			"scene_list_label":       "【Extracted Scene Backgrounds】",
			"prop_list_label":        "【Available Prop List】",
			"task_instruction":       "Break down the script into storyboard shots according to the rules and pacing defined in your role instructions.",
			"character_constraint":   "**Important**: In the characters field, only use character IDs (numbers) from the above character list. Do not create new characters or use other IDs.",
			"scene_constraint":       "**Important**: In the scene_id field, select the most matching background ID (number) from the above background list. If no suitable background exists, use null.",
			"prop_constraint":        "**Important**: In the props field, only use prop IDs (numbers) from the above prop list. Do not create new props or use other IDs. Use empty array [] if no props are featured.",
			"shot_description_label": "Shot description: %s",
			"scene_label":            "Scene: %s, %s",
			"characters_label":       "Characters: %s",
			"action_label":           "Action: %s",
			"result_label":           "Result: %s",
			"dialogue_label":         "Dialogue: %s",
			"atmosphere_label":       "Atmosphere: %s",
			"shot_type_label":        "Shot type: %s",
			"angle_label":            "Angle: %s",
			"movement_label":         "Movement: %s",
			"drama_info_template":    "Title: %s\nSummary: %s\nGenre: %s",
		},
		"zh": {

			"outline_request":        "Please create a short drama outline for the following theme:\n\nTheme: %s",
			"genre_preference":       "\nGenre preference: %s",
			"style_requirement":      "\nStyle requirement: %s",
			"episode_count":          "\nNumber of episodes: %d episodes",
			"episode_importance":     "\n\n**Important: Must plan complete storylines for all %d episodes in the episodes array, each with clear story content!**",
			"character_request":      "Script content:\n%s\n\nPlease extract and organize detailed character profiles for all main characters from the script.",
			"episode_script_request": "Drama outline:\n%s\n%s\nPlease create detailed scripts for %d episodes based on the above outline and characters.\n\n**Important requirements:**\n- Must generate all %d episodes, from episode 1 to episode %d, cannot skip any\n- Each episode is about 3-5 minutes (150-300 seconds)\n- The duration field for each episode should be set reasonably based on script content length, not all the same value\n- The episodes array in the returned JSON must contain %d elements",
			"frame_info":             "Shot information:\n%s\n\nPlease directly generate the image prompt for the first frame without any explanation:",
			"key_frame_info":         "Shot information:\n%s\n\nPlease directly generate the image prompt for the key frame without any explanation:",
			"last_frame_info":        "Shot information:\n%s\n\nPlease directly generate the image prompt for the last frame without any explanation:",
			"script_content_label":   "【Script Content】",
			"storyboard_list_label":  "【Storyboard List】",
			"task_label":             "【Task】",
			"character_list_label":   "【Available Character List】",
			"scene_list_label":       "【Extracted Scene Backgrounds】",
			"prop_list_label":        "【Available Prop List】",
			"task_instruction":       "Break down the script into storyboard shots according to the rules and pacing defined in your role instructions.",
			"character_constraint":   "**Important**: In the characters field, only use character IDs (numbers) from the above character list. Do not create new characters or use other IDs.",
			"scene_constraint":       "**Important**: In the scene_id field, select the most matching background ID (number) from the above background list. If no suitable background exists, use null.",
			"prop_constraint":        "**Important**: In the props field, only use prop IDs (numbers) from the above prop list. Do not create new props or use other IDs. Use empty array [] if no props are featured.",
			"shot_description_label": "Shot description: %s",
			"scene_label":            "Scene: %s, %s",
			"characters_label":       "Characters: %s",
			"action_label":           "Action: %s",
			"result_label":           "Result: %s",
			"dialogue_label":         "Dialogue: %s",
			"atmosphere_label":       "Atmosphere: %s",
			"shot_type_label":        "Shot type: %s",
			"angle_label":            "Angle: %s",
			"movement_label":         "Movement: %s",
			"drama_info_template":    "Title: %s\nSummary: %s\nGenre: %s",
		},
	}

	lang := "zh"
	if true || p.IsEnglish() {
		lang = "en"
	}

	template, ok := templates[lang][key]
	if !ok {
		return ""
	}

	if len(args) > 0 {
		return fmt.Sprintf(template, args...)
	}
	return template
}

// GetStylePrompt 获取风格提示词
func (p *PromptI18n) GetStylePrompt(style string, customStyle string) string {
	if style == "" {
		return ""
	}

	if style == "custom" {
		return fmt.Sprintf("You are an expert Art Director. The exact style you must consistently follow for all visual designs and prompts is: %s", customStyle)
	}

	stylePrompts := map[string]map[string]string{
		"zh": {
			"ghibli": prompts.Get("style_prompt.txt"),

			"guoman": prompts.Get("style_prompt.txt"),

			"wasteland": prompts.Get("style_prompt.txt"),

			"nostalgia": prompts.Get("style_prompt.txt"),

			"pixel": prompts.Get("style_prompt.txt"),

			"voxel": prompts.Get("style_prompt.txt"),

			"urban": prompts.Get("style_prompt.txt"),

			"guoman3d": prompts.Get("style_prompt.txt"),

			"chibi3d": prompts.Get("style_prompt.txt"),
		},
		"en": {
			"ghibli": prompts.Get("style_prompt.txt"),

			"guoman": prompts.Get("style_prompt.txt"),

			"wasteland": prompts.Get("style_prompt.txt"),

			"nostalgia": prompts.Get("style_prompt.txt"),

			"pixel": prompts.Get("style_prompt.txt"),

			"voxel": prompts.Get("style_prompt.txt"),

			"urban": prompts.Get("style_prompt.txt"),

			"guoman3d": prompts.Get("style_prompt.txt"),

			"chibi3d": prompts.Get("style_prompt.txt"),
		},
	}

	lang := "zh"
	if true || p.IsEnglish() {
		lang = "en"
	}

	if prompts, ok := stylePrompts[lang]; ok {
		if prompt, exists := prompts[style]; exists {
			return prompt
		}
	}

	return ""
}

// GetVideoConstraintPrompt 获取视频生成的约束提示词
// referenceMode: "single" (单图), "first_last" (首尾帧), "multiple" (多图), "action_sequence" (动作序列), "rapid_cut" (rapid cut多场景)
func (p *PromptI18n) GetVideoConstraintPrompt(referenceMode string) string {
	// Rapid cut mode—use rapid cut specific constraint
	if referenceMode == "rapid_cut" {
		return prompts.Get("video_constraint_rapid_cut.txt")
	}

	switch referenceMode {
	case "action_sequence":
		return "[System Note] The provided reference image is a storyboard grid showing a continuous sequence of actions from left to right. Do not generate a split-screen video. Instead, generate a single-frame, smooth video that transitions sequentially through these exact key poses in chronological order."
	case "first_frame_strict":
		return "[System Note] The provided reference image is the strict starting frame. The video must begin EXACTLY with this visual state, character pose, and environment, then smoothly start the action described above."
	case "last_frame_strict":
		return "[System Note] The provided reference image is the final target frame. The video must end EXACTLY matching this visual state. Plan the pacing so the action resolves into this specific pose at the end of the video."
	case "key_frame_style":
		return "[System Note] The provided reference image is for character, lighting, and aesthetic reference only. You do not need to strictly start or end with this exact pose, but maintain the exact same identity, clothing, and scene logic while performing the requested action."
	}

	return ""
}

// ==========================================
// WithDrama* Methods - Template-aware versions
// These methods resolve prompts via PromptTemplateService fallback.
// If dramaID is 0 or templateService is nil, they behave identically to originals.
// ==========================================

// WithDramaStoryboardSystemPrompt resolves storyboard system prompt for a specific drama
func (p *PromptI18n) WithDramaStoryboardSystemPrompt(dramaID uint) string {
	return p.resolvePrompt(dramaID, "storyboard_breakdown")
}

// WithDramaSceneExtractionPrompt resolves scene extraction prompt for a specific drama
func (p *PromptI18n) WithDramaSceneExtractionPrompt(dramaID uint, style string, customStyle string) string {
	imageRatio := "16:9"
	effectiveStyle := p.resolveEffectiveStyle(dramaID, style, customStyle)
	template := p.resolvePrompt(dramaID, "scene_extraction")
	// Use safe replacement for custom templates that may not have %s
	if strings.Contains(template, "%s") {
		return fmt.Sprintf(template, effectiveStyle, imageRatio)
	}
	return formatPromptWithVars(template, map[string]string{
		"{{STYLE}}": effectiveStyle,
		"{{RATIO}}": imageRatio,
	})
}

// WithDramaCharacterExtractionPrompt resolves character extraction prompt for a specific drama
func (p *PromptI18n) WithDramaCharacterExtractionPrompt(dramaID uint, style string, customStyle string) string {
	imageRatio := "16:9"
	effectiveStyle := p.resolveEffectiveStyle(dramaID, style, customStyle)
	template := p.resolvePrompt(dramaID, "character_extraction")
	if strings.Contains(template, "%s") {
		return fmt.Sprintf(template, effectiveStyle, imageRatio)
	}
	return formatPromptWithVars(template, map[string]string{
		"{{STYLE}}": effectiveStyle,
		"{{RATIO}}": imageRatio,
	})
}

// WithDramaPropExtractionPrompt resolves prop extraction prompt for a specific drama
func (p *PromptI18n) WithDramaPropExtractionPrompt(dramaID uint, style string, customStyle string) string {
	imageRatio := "1:1"
	effectiveStyle := p.resolveEffectiveStyle(dramaID, style, customStyle)
	template := p.resolvePrompt(dramaID, "prop_extraction")
	if strings.Contains(template, "%s") {
		return fmt.Sprintf(template, effectiveStyle, imageRatio)
	}
	return formatPromptWithVars(template, map[string]string{
		"{{STYLE}}": effectiveStyle,
		"{{RATIO}}": imageRatio,
	})
}

// WithDramaOutlineGenerationPrompt resolves outline generation prompt for a specific drama
func (p *PromptI18n) WithDramaOutlineGenerationPrompt(dramaID uint) string {
	return p.resolvePrompt(dramaID, "script_outline")
}

// WithDramaEpisodeScriptPrompt resolves episode script generation prompt for a specific drama
func (p *PromptI18n) WithDramaEpisodeScriptPrompt(dramaID uint) string {
	return p.resolvePrompt(dramaID, "script_episode")
}

// WithDramaFirstFramePrompt resolves first frame prompt for a specific drama
func (p *PromptI18n) WithDramaFirstFramePrompt(dramaID uint, style string) string {
	imageRatio := "16:9"
	effectiveStyle := p.resolveEffectiveStyle(dramaID, style, "")
	template := p.resolvePrompt(dramaID, "image_first_frame")
	if strings.Contains(template, "%s") {
		return fmt.Sprintf(template, effectiveStyle, imageRatio)
	}
	return formatPromptWithVars(template, map[string]string{
		"{{STYLE}}": effectiveStyle,
		"{{RATIO}}": imageRatio,
	})
}

// WithDramaKeyFramePrompt resolves key frame prompt for a specific drama
func (p *PromptI18n) WithDramaKeyFramePrompt(dramaID uint, style string) string {
	imageRatio := "16:9"
	effectiveStyle := p.resolveEffectiveStyle(dramaID, style, "")
	template := p.resolvePrompt(dramaID, "image_key_frame")
	if strings.Contains(template, "%s") {
		return fmt.Sprintf(template, effectiveStyle, imageRatio)
	}
	return formatPromptWithVars(template, map[string]string{
		"{{STYLE}}": effectiveStyle,
		"{{RATIO}}": imageRatio,
	})
}

// WithDramaLastFramePrompt resolves last frame prompt for a specific drama
func (p *PromptI18n) WithDramaLastFramePrompt(dramaID uint, style string) string {
	imageRatio := "16:9"
	effectiveStyle := p.resolveEffectiveStyle(dramaID, style, "")
	template := p.resolvePrompt(dramaID, "image_last_frame")
	if strings.Contains(template, "%s") {
		return fmt.Sprintf(template, effectiveStyle, imageRatio)
	}
	return formatPromptWithVars(template, map[string]string{
		"{{STYLE}}": effectiveStyle,
		"{{RATIO}}": imageRatio,
	})
}

// WithDramaActionSequenceFramePrompt resolves action sequence frame prompt for a specific drama
func (p *PromptI18n) WithDramaActionSequenceFramePrompt(dramaID uint, style string) string {
	imageRatio := "16:9"
	effectiveStyle := p.resolveEffectiveStyle(dramaID, style, "")
	template := p.resolvePrompt(dramaID, "image_action_sequence")
	if strings.Contains(template, "%s") {
		return fmt.Sprintf(template, effectiveStyle, imageRatio)
	}
	return formatPromptWithVars(template, map[string]string{
		"{{STYLE}}": effectiveStyle,
		"{{RATIO}}": imageRatio,
	})
}

// WithDramaVideoConstraintPrompt resolves video constraint prompt for a specific drama
func (p *PromptI18n) WithDramaVideoConstraintPrompt(dramaID uint, referenceMode string) string {
	// Rapid cut mode always uses the rapid cut constraint, bypassing templates
	if referenceMode == "rapid_cut" {
		return prompts.Get("video_constraint_rapid_cut.txt")
	}

	template := p.resolvePrompt(dramaID, "video_constraint")
	if template == "" {
		// Fallback to original logic
		return p.GetVideoConstraintPrompt(referenceMode)
	}
	return template
}

// WithDramaVideoExtractionPrompt resolves video extraction prompt for a specific drama
func (p *PromptI18n) WithDramaVideoExtractionPrompt(dramaID uint, style string, customStyle string) string {
	effectiveStyle := p.resolveEffectiveStyle(dramaID, style, customStyle)
	template := p.resolvePrompt(dramaID, "video_extraction")
	if strings.Contains(template, "%s") {
		return fmt.Sprintf(template, effectiveStyle)
	}
	return formatPromptWithVars(template, map[string]string{
		"{{STYLE}}": effectiveStyle,
	})
}

// WithDramaRapidCutActionSequenceFramePrompt resolves rapid cut action sequence prompt for a specific drama
func (p *PromptI18n) WithDramaRapidCutActionSequenceFramePrompt(dramaID uint, style string) string {
	imageRatio := "16:9"
	effectiveStyle := p.resolveEffectiveStyle(dramaID, style, "")
	// For now, rapid cut always uses the rapid cut prompt (no template override)
	template := prompts.Get("image_action_sequence_rapid_cut.txt")
	if strings.Contains(template, "%s") {
		return fmt.Sprintf(template, effectiveStyle, imageRatio)
	}
	return formatPromptWithVars(template, map[string]string{
		"{{STYLE}}": effectiveStyle,
		"{{RATIO}}": imageRatio,
	})
}

// WithDramaStylePrompt resolves style prompt for a specific drama
func (p *PromptI18n) WithDramaStylePrompt(dramaID uint, style string, customStyle string) string {
	if style == "custom" {
		return fmt.Sprintf("You are an expert Art Director. The exact style you must consistently follow for all visual designs and prompts is: %s", customStyle)
	}
	template := p.resolvePrompt(dramaID, "style_prompt")
	if template != "" {
		return template
	}
	// Fallback to original per-style logic
	return p.GetStylePrompt(style, customStyle)
}
