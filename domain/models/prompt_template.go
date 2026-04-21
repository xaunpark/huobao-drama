package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// PromptTemplate 自定义提示词模板
// 用户可以创建自定义模板并在不同项目(Drama)之间复用
type PromptTemplate struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string         `gorm:"type:varchar(200);not null" json:"name"`
	Description *string        `gorm:"type:text" json:"description"`
	Prompts     datatypes.JSON `gorm:"type:json" json:"prompts"` // key-value: prompt_type -> dynamic content
	CreatedAt   time.Time      `gorm:"not null;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"not null;autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (p *PromptTemplate) TableName() string {
	return "prompt_templates"
}

// PromptTemplatePrompts 用于解析 Prompts JSON 字段的结构
// 每个 key 对应一个 prompt 类型，value 为用户自定义的 Dynamic 内容
// 如果某个 key 不存在或为空字符串，系统将使用默认的 embed 文件
type PromptTemplatePrompts struct {
	StoryboardBreakdown string `json:"storyboard_breakdown,omitempty"`
	CharacterExtraction string `json:"character_extraction,omitempty"`
	SceneExtraction     string `json:"scene_extraction,omitempty"`
	PropExtraction      string `json:"prop_extraction,omitempty"`
	ScriptOutline       string `json:"script_outline,omitempty"`
	ScriptEpisode       string `json:"script_episode,omitempty"`
	ImageFirstFrame     string `json:"image_first_frame,omitempty"`
	ImageKeyFrame       string `json:"image_key_frame,omitempty"`
	ImageLastFrame      string `json:"image_last_frame,omitempty"`
	ImageActionSequence string `json:"image_action_sequence,omitempty"`
	VideoConstraint     string `json:"video_constraint,omitempty"`
	StylePrompt         string `json:"style_prompt,omitempty"`
	VideoExtraction     string `json:"video_extraction,omitempty"`
	VisualUnitBreakdown    string `json:"visual_unit_breakdown,omitempty"`    // Voice-over AI Director shot planning
	NurseryRhymeBreakdown  string `json:"nursery_rhyme_breakdown,omitempty"`  // Nursery rhyme lyrics-synced shot planning
	MVMakerGamingHorror    string `json:"mv_maker_gaming_horror,omitempty"`   // MV Maker: Gaming horror genre prompt
	MVMakerCinematicMovie  string `json:"mv_maker_cinematic_movie,omitempty"` // MV Maker: Cinematic movie genre prompt
	NarrativeMVPlanner     string `json:"narrative_mv_planner,omitempty"`    // Narrative MV: story planning prompt (Phase 1)
	NarrativeMVDirector    string `json:"narrative_mv_director,omitempty"`   // Narrative MV: shot director prompt (Phase 2)
	NarrativeMVCG5         string `json:"narrative_mv_cg5,omitempty"`        // Narrative MV: CG5 3-act visual style template
}

// PromptTypeToDefaultFile 将 prompt type key 映射到默认的 embed 文件名
var PromptTypeToDefaultFile = map[string]string{
	"storyboard_breakdown": "storyboard_story_breakdown.txt",
	"character_extraction": "character_extraction.txt",
	"scene_extraction":     "scene_extraction.txt",
	"prop_extraction":      "prop_extraction.txt",
	"script_outline":       "script_outline_generation.txt",
	"script_episode":       "script_episode_generation.txt",
	"image_first_frame":    "image_first_frame.txt",
	"image_key_frame":      "image_key_frame.txt",
	"image_last_frame":     "image_last_frame.txt",
	"image_action_sequence": "image_action_sequence.txt",
	"video_constraint":     "video_constraint_prefixes.txt",
	"style_prompt":            "style_prompt.txt",
	"video_extraction":        "video_extraction.txt",
	"visual_unit_breakdown":    "storyboard_visual_unit.txt",
	"nursery_rhyme_breakdown":  "storyboard_nursery_rhyme.txt",
	"mv_maker_gaming_horror":   "storyboard_mv_gaming_horror.txt",
	"mv_maker_cinematic_movie": "storyboard_mv_cinematic_movie.txt",
	"narrative_mv_planner":     "storyboard_narrative_planner.txt",
	"narrative_mv_director":    "storyboard_narrative_director.txt",
	"narrative_mv_cg5":         "storyboard_narrative_cg5.txt",
}
