package services

import (
	"fmt"

	"github.com/drama-generator/backend/pkg/config"
	"github.com/drama-generator/backend/application/prompts"
)

// PromptI18n 提示词国际化工具
type PromptI18n struct {
	config *config.Config
}

// NewPromptI18n 创建提示词国际化工具
func NewPromptI18n(cfg *config.Config) *PromptI18n {
	return &PromptI18n{config: cfg}
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
func (p *PromptI18n) GetSceneExtractionPrompt(style string) string {
	// 默认图片比例
	imageRatio := "16:9"

	if true || p.IsEnglish() {
		return fmt.Sprintf(prompts.Get("scene_extraction.txt"), style, imageRatio)
	}

	return fmt.Sprintf(prompts.Get("scene_extraction.txt"), style, imageRatio)
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

// GetLastFramePrompt 获取尾帧提示词
func (p *PromptI18n) GetLastFramePrompt(style string) string {
	imageRatio := "16:9"
	if true || p.IsEnglish() {
		return fmt.Sprintf(prompts.Get("image_last_frame.txt"), style, imageRatio)
	}

	return fmt.Sprintf(prompts.Get("image_last_frame.txt"), style, imageRatio)
}

// GetOutlineGenerationPrompt 获取大纲生成提示词
func (p *PromptI18n) GetOutlineGenerationPrompt() string {
	if true || p.IsEnglish() {
		return prompts.Get("script_outline_generation.txt")
	}

	return prompts.Get("script_outline_generation.txt")
}

// GetCharacterExtractionPrompt 获取角色提取提示词
func (p *PromptI18n) GetCharacterExtractionPrompt(style string) string {
	imageRatio := "16:9"
	if true || p.IsEnglish() {
		return fmt.Sprintf(prompts.Get("character_extraction.txt"), style, imageRatio)
	}

	return fmt.Sprintf(prompts.Get("character_extraction.txt"), style, imageRatio)
}

// GetPropExtractionPrompt 获取道具提取提示词
func (p *PromptI18n) GetPropExtractionPrompt(style string) string {
	imageRatio := "1:1"

	if true || p.IsEnglish() {
		return fmt.Sprintf(prompts.Get("prop_extraction.txt"), style, imageRatio)
	}

	return fmt.Sprintf(prompts.Get("prop_extraction.txt"), style, imageRatio)
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
			"character_request":      "Script content:\n%s\n\nPlease extract and organize detailed character profiles for up to %d main characters from the script.",
			"episode_script_request": "Drama outline:\n%s\n%s\nPlease create detailed scripts for %d episodes based on the above outline and characters.\n\n**Important requirements:**\n- Must generate all %d episodes, from episode 1 to episode %d, cannot skip any\n- Each episode is about 3-5 minutes (150-300 seconds)\n- The duration field for each episode should be set reasonably based on script content length, not all the same value\n- The episodes array in the returned JSON must contain %d elements",
			"frame_info":             "Shot information:\n%s\n\nPlease directly generate the image prompt for the first frame without any explanation:",
			"key_frame_info":         "Shot information:\n%s\n\nPlease directly generate the image prompt for the key frame without any explanation:",
			"last_frame_info":        "Shot information:\n%s\n\nPlease directly generate the image prompt for the last frame without any explanation:",
			"script_content_label":   "【Script Content】",
			"storyboard_list_label":  "【Storyboard List】",
			"task_label":             "【Task】",
			"character_list_label":   "【Available Character List】",
			"scene_list_label":       "【Extracted Scene Backgrounds】",
			"task_instruction":       "Break down the novel script into storyboard shots based on **independent action units**.",
			"character_constraint":   "**Important**: In the characters field, only use character IDs (numbers) from the above character list. Do not create new characters or use other IDs.",
			"scene_constraint":       "**Important**: In the scene_id field, select the most matching background ID (number) from the above background list. If no suitable background exists, use null.",
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
			"character_request":      "Script content:\n%s\n\nPlease extract and organize detailed character profiles for up to %d main characters from the script.",
			"episode_script_request": "Drama outline:\n%s\n%s\nPlease create detailed scripts for %d episodes based on the above outline and characters.\n\n**Important requirements:**\n- Must generate all %d episodes, from episode 1 to episode %d, cannot skip any\n- Each episode is about 3-5 minutes (150-300 seconds)\n- The duration field for each episode should be set reasonably based on script content length, not all the same value\n- The episodes array in the returned JSON must contain %d elements",
			"frame_info":             "Shot information:\n%s\n\nPlease directly generate the image prompt for the first frame without any explanation:",
			"key_frame_info":         "Shot information:\n%s\n\nPlease directly generate the image prompt for the key frame without any explanation:",
			"last_frame_info":        "Shot information:\n%s\n\nPlease directly generate the image prompt for the last frame without any explanation:",
			"script_content_label":   "【Script Content】",
			"storyboard_list_label":  "【Storyboard List】",
			"task_label":             "【Task】",
			"character_list_label":   "【Available Character List】",
			"scene_list_label":       "【Extracted Scene Backgrounds】",
			"task_instruction":       "Break down the novel script into storyboard shots based on **independent action units**.",
			"character_constraint":   "**Important**: In the characters field, only use character IDs (numbers) from the above character list. Do not create new characters or use other IDs.",
			"scene_constraint":       "**Important**: In the scene_id field, select the most matching background ID (number) from the above background list. If no suitable background exists, use null.",
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
func (p *PromptI18n) GetStylePrompt(style string) string {
	if style == "" {
		return ""
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
// referenceMode: "single" (单图), "first_last" (首尾帧), "multiple" (多图), "action_sequence" (动作序列)
func (p *PromptI18n) GetVideoConstraintPrompt(referenceMode string) string {
	// 动作序列图（九宫格）的约束提示词
	actionSequencePrompts := map[string]string{
		"zh": prompts.Get("video_constraint_prefixes.txt"),

		"en": prompts.Get("video_constraint_prefixes.txt"),
	}

	// 通用约束提示词（单图、首尾帧、多图）
	generalPrompts := map[string]string{
		"zh": prompts.Get("video_constraint_prefixes.txt"),

		"en": prompts.Get("video_constraint_prefixes.txt"),
	}

	lang := "zh"
	if true || p.IsEnglish() {
		lang = "en"
	}

	// 如果是动作序列模式，返回九宫格约束提示词
	if referenceMode == "action_sequence" {
		if prompt, ok := actionSequencePrompts[lang]; ok {
			return prompt
		}
	}

	// 其他模式返回通用约束提示词
	if prompt, ok := generalPrompts[lang]; ok {
		return prompt
	}

	return ""
}
