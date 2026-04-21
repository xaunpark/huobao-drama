package services

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/drama-generator/backend/application/prompts"
	"github.com/drama-generator/backend/application/prompts/fixed"
	"github.com/drama-generator/backend/domain/models"
	"github.com/drama-generator/backend/pkg/logger"
	"gorm.io/gorm"
)

// PromptTemplateService handles CRUD for prompt templates
type PromptTemplateService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewPromptTemplateService(db *gorm.DB, log *logger.Logger) *PromptTemplateService {
	return &PromptTemplateService{db: db, log: log}
}

// --- CRUD ---

type CreatePromptTemplateRequest struct {
	Name        string                        `json:"name" binding:"required,min=1,max=200"`
	Description string                        `json:"description"`
	Prompts     models.PromptTemplatePrompts   `json:"prompts"`
}

type UpdatePromptTemplateRequest struct {
	Name        string                        `json:"name" binding:"omitempty,min=1,max=200"`
	Description string                        `json:"description"`
	Prompts     *models.PromptTemplatePrompts  `json:"prompts"`
}

func (s *PromptTemplateService) List() ([]models.PromptTemplate, error) {
	var templates []models.PromptTemplate
	err := s.db.Order("updated_at DESC").Find(&templates).Error
	return templates, err
}

func (s *PromptTemplateService) GetByID(id uint) (*models.PromptTemplate, error) {
	var tpl models.PromptTemplate
	if err := s.db.First(&tpl, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("prompt template not found")
		}
		return nil, err
	}
	return &tpl, nil
}

func (s *PromptTemplateService) Create(req *CreatePromptTemplateRequest) (*models.PromptTemplate, error) {
	promptsJSON, err := json.Marshal(req.Prompts)
	if err != nil {
		return nil, err
	}

	tpl := &models.PromptTemplate{
		Name:    req.Name,
		Prompts: promptsJSON,
	}
	if req.Description != "" {
		tpl.Description = &req.Description
	}

	if err := s.db.Create(tpl).Error; err != nil {
		s.log.Errorw("Failed to create prompt template", "error", err)
		return nil, err
	}

	s.log.Infow("Prompt template created", "id", tpl.ID, "name", tpl.Name)
	return tpl, nil
}

func (s *PromptTemplateService) Update(id uint, req *UpdatePromptTemplateRequest) (*models.PromptTemplate, error) {
	var tpl models.PromptTemplate
	if err := s.db.First(&tpl, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("prompt template not found")
		}
		return nil, err
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Prompts != nil {
		promptsJSON, err := json.Marshal(req.Prompts)
		if err != nil {
			return nil, err
		}
		updates["prompts"] = promptsJSON
	}
	updates["updated_at"] = time.Now()

	if err := s.db.Model(&tpl).Updates(updates).Error; err != nil {
		s.log.Errorw("Failed to update prompt template", "error", err)
		return nil, err
	}

	s.log.Infow("Prompt template updated", "id", id)
	return &tpl, nil
}

func (s *PromptTemplateService) Delete(id uint) error {
	// Check if any drama is using this template
	var count int64
	s.db.Model(&models.Drama{}).Where("prompt_template_id = ?", id).Count(&count)
	if count > 0 {
		return errors.New("template is in use by one or more projects, cannot delete")
	}

	result := s.db.Delete(&models.PromptTemplate{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("prompt template not found")
	}
	s.log.Infow("Prompt template deleted", "id", id)
	return nil
}

func (s *PromptTemplateService) Duplicate(id uint) (*models.PromptTemplate, error) {
	src, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	newName := src.Name + " (Copy)"
	newTpl := &models.PromptTemplate{
		Name:        newName,
		Description: src.Description,
		Prompts:     src.Prompts,
	}

	if err := s.db.Create(newTpl).Error; err != nil {
		return nil, err
	}

	s.log.Infow("Prompt template duplicated", "source_id", id, "new_id", newTpl.ID)
	return newTpl, nil
}

// --- Default Prompts (trả về Dynamic mặc định từ embed files) ---

// GetDefaultPrompts returns the default dynamic content from embedded .txt files
// so that frontend can display them as placeholders or let user "Load Default"
func (s *PromptTemplateService) GetDefaultPrompts() models.PromptTemplatePrompts {
	return models.PromptTemplatePrompts{
		StoryboardBreakdown: prompts.Get("storyboard_story_breakdown.txt"),
		CharacterExtraction: prompts.Get("character_extraction.txt"),
		SceneExtraction:     prompts.Get("scene_extraction.txt"),
		PropExtraction:      prompts.Get("prop_extraction.txt"),
		ScriptOutline:       prompts.Get("script_outline_generation.txt"),
		ScriptEpisode:       prompts.Get("script_episode_generation.txt"),
		ImageFirstFrame:     prompts.Get("image_first_frame.txt"),
		ImageKeyFrame:       prompts.Get("image_key_frame.txt"),
		ImageLastFrame:      prompts.Get("image_last_frame.txt"),
		ImageActionSequence: prompts.Get("image_action_sequence.txt"),
		VideoConstraint:     prompts.Get("video_constraint_prefixes.txt"),
		StylePrompt:         prompts.Get("style_prompt.txt"),
		VideoExtraction:     prompts.Get("video_extraction.txt"),
		VisualUnitBreakdown: prompts.Get("storyboard_visual_unit.txt"),
	}
}

// --- PromptResolver: Fallback Logic ---

// ResolvePrompt returns the final Dynamic prompt for a given drama and prompt type.
// Fallback chain: Template override → Default embed file
// The Fixed format instructions are always appended separately by the caller.
func (s *PromptTemplateService) ResolvePrompt(dramaID uint, promptType string) string {
	// 1. Find drama's template
	var drama models.Drama
	if err := s.db.Select("prompt_template_id").First(&drama, dramaID).Error; err != nil {
		// Drama not found or error → use default
		return s.getDefaultPrompt(promptType)
	}

	// 2. No template assigned → use default
	if drama.PromptTemplateID == nil {
		return s.getDefaultPrompt(promptType)
	}

	// 3. Load template
	var tpl models.PromptTemplate
	if err := s.db.First(&tpl, *drama.PromptTemplateID).Error; err != nil {
		return s.getDefaultPrompt(promptType)
	}

	// 4. Parse template prompts JSON
	var tplPrompts models.PromptTemplatePrompts
	if err := json.Unmarshal(tpl.Prompts, &tplPrompts); err != nil {
		s.log.Warnw("Failed to parse template prompts, using default", "template_id", tpl.ID, "error", err)
		return s.getDefaultPrompt(promptType)
	}

	// 5. Check if prompt type has a custom value
	customValue := s.getPromptFromStruct(tplPrompts, promptType)
	if customValue != "" {
		return customValue
	}

	// 6. Fallback to default
	return s.getDefaultPrompt(promptType)
}

// ResolvePromptIfCustom returns the template's custom prompt ONLY if the drama has
// a template and the specific prompt field is non-empty. Unlike ResolvePrompt, this
// never falls back to default embed files. Returns "" if no custom override exists.
func (s *PromptTemplateService) ResolvePromptIfCustom(dramaID uint, promptType string) string {
	var drama models.Drama
	if err := s.db.Select("prompt_template_id").First(&drama, dramaID).Error; err != nil {
		return ""
	}
	if drama.PromptTemplateID == nil {
		return ""
	}
	var tpl models.PromptTemplate
	if err := s.db.First(&tpl, *drama.PromptTemplateID).Error; err != nil {
		return ""
	}
	var tplPrompts models.PromptTemplatePrompts
	if err := json.Unmarshal(tpl.Prompts, &tplPrompts); err != nil {
		return ""
	}
	return s.getPromptFromStruct(tplPrompts, promptType)
}

// GetFixedPrompt returns the fixed format instructions for a prompt type
func (s *PromptTemplateService) GetFixedPrompt(promptType string) string {
	return fixed.Get(promptType)
}

// getDefaultPrompt returns default embed file content by prompt type key
func (s *PromptTemplateService) getDefaultPrompt(promptType string) string {
	filename, ok := models.PromptTypeToDefaultFile[promptType]
	if !ok {
		return ""
	}
	return prompts.Get(filename)
}

// getPromptFromStruct extracts prompt value from struct by type key
func (s *PromptTemplateService) getPromptFromStruct(p models.PromptTemplatePrompts, promptType string) string {
	switch promptType {
	case "storyboard_breakdown":
		return p.StoryboardBreakdown
	case "character_extraction":
		return p.CharacterExtraction
	case "scene_extraction":
		return p.SceneExtraction
	case "prop_extraction":
		return p.PropExtraction
	case "script_outline":
		return p.ScriptOutline
	case "script_episode":
		return p.ScriptEpisode
	case "image_first_frame":
		return p.ImageFirstFrame
	case "image_key_frame":
		return p.ImageKeyFrame
	case "image_last_frame":
		return p.ImageLastFrame
	case "image_action_sequence":
		return p.ImageActionSequence
	case "video_constraint":
		return p.VideoConstraint
	case "style_prompt":
		return p.StylePrompt
	case "video_extraction":
		return p.VideoExtraction
	case "visual_unit_breakdown":
		return p.VisualUnitBreakdown
	case "narrative_mv_planner":
		return p.NarrativeMVPlanner
	case "narrative_mv_director":
		return p.NarrativeMVDirector
	default:
		return ""
	}
}
