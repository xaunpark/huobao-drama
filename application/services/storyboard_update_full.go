package services

import (
	"fmt"
	"strings"

	"github.com/drama-generator/backend/domain/models"
)

// UpdateStoryboard 更新分镜的所有字段，并重新生成提示词
func (s *StoryboardService) UpdateStoryboard(storyboardID string, updates map[string]interface{}) error {
	// 查找分镜
	var storyboard models.Storyboard
	if err := s.db.Preload("Props").First(&storyboard, storyboardID).Error; err != nil {
		return fmt.Errorf("storyboard not found: %w", err)
	}



	// 构建用于重新生成提示词的Storyboard结构
	sb := Storyboard{
		ShotNumber: storyboard.StoryboardNumber,
	}

	// 从updates中提取字段并更新
	updateData := make(map[string]interface{})

	if val, ok := updates["title"].(string); ok && val != "" {
		updateData["title"] = val
		sb.Title = val
	}
	if val, ok := updates["shot_type"].(string); ok && val != "" {
		updateData["shot_type"] = val
		sb.ShotType = val
	}
	if val, ok := updates["angle"].(string); ok && val != "" {
		updateData["angle"] = val
		sb.Angle = val
	}
	if val, ok := updates["movement"].(string); ok && val != "" {
		updateData["movement"] = val
		sb.Movement = val
	}
	if val, ok := updates["location"].(string); ok && val != "" {
		updateData["location"] = val
		sb.Location = val
	}
	if val, ok := updates["time"].(string); ok && val != "" {
		updateData["time"] = val
		sb.Time = val
	}
	if val, ok := updates["action"].(string); ok && val != "" {
		updateData["action"] = val
		sb.Action = val
	}
	if val, ok := updates["dialogue"].(string); ok && val != "" {
		updateData["dialogue"] = val
		sb.Dialogue = val
	}
	if val, ok := updates["result"].(string); ok && val != "" {
		updateData["result"] = val
		sb.Result = val
	}
	if val, ok := updates["atmosphere"].(string); ok && val != "" {
		updateData["atmosphere"] = val
		sb.Atmosphere = val
	}
	if val, ok := updates["description"].(string); ok && val != "" {
		updateData["description"] = val
	}
	if val, ok := updates["image_prompt"].(string); ok && val != "" {
		updateData["image_prompt"] = val
	}
	if val, ok := updates["video_prompt"].(string); ok && val != "" {
		updateData["video_prompt"] = val
	}
	if val, ok := updates["bgm_prompt"].(string); ok && val != "" {
		updateData["bgm_prompt"] = val
		sb.BgmPrompt = val
	}
	if val, ok := updates["sound_effect"].(string); ok && val != "" {
		updateData["sound_effect"] = val
		sb.SoundEffect = val
	}
	if val, ok := updates["duration"].(float64); ok {
		updateData["duration"] = int(val)
		sb.Duration = int(val)
	}
	if val, ok := updates["scene_id"].(float64); ok {
		sceneID := uint(val)
		updateData["scene_id"] = sceneID
	}

	// 使用当前数据库值填充缺失字段（用于生成提示词）
	if sb.Title == "" && storyboard.Title != nil {
		sb.Title = *storyboard.Title
	}
	if sb.ShotType == "" && storyboard.ShotType != nil {
		sb.ShotType = *storyboard.ShotType
	}
	if sb.Angle == "" && storyboard.Angle != nil {
		sb.Angle = *storyboard.Angle
	}
	if sb.Movement == "" && storyboard.Movement != nil {
		sb.Movement = *storyboard.Movement
	}
	if sb.Location == "" && storyboard.Location != nil {
		sb.Location = *storyboard.Location
	}
	if sb.Time == "" && storyboard.Time != nil {
		sb.Time = *storyboard.Time
	}
	if sb.Action == "" && storyboard.Action != nil {
		sb.Action = *storyboard.Action
	}
	if sb.Dialogue == "" && storyboard.Dialogue != nil {
		sb.Dialogue = *storyboard.Dialogue
	}
	if sb.Result == "" && storyboard.Result != nil {
		sb.Result = *storyboard.Result
	}
	if sb.Atmosphere == "" && storyboard.Atmosphere != nil {
		sb.Atmosphere = *storyboard.Atmosphere
	}
	if sb.BgmPrompt == "" && storyboard.BgmPrompt != nil {
		sb.BgmPrompt = *storyboard.BgmPrompt
	}
	if sb.SoundEffect == "" && storyboard.SoundEffect != nil {
		sb.SoundEffect = *storyboard.SoundEffect
	}
	if sb.Duration == 0 {
		sb.Duration = storyboard.Duration
	}
	// Voice-over fields: always load from DB for video prompt generation
	if storyboard.NarratorScript != nil {
		sb.NarratorScript = *storyboard.NarratorScript
	}
	if storyboard.AudioMode != nil {
		sb.AudioMode = *storyboard.AudioMode
	}

	// Load character visual descriptions and voice styles for video prompt
	var charIDs []uint
	if err := s.db.Table("storyboard_characters").
		Where("storyboard_id = ?", storyboard.ID).
		Pluck("character_id", &charIDs).Error; err == nil && len(charIDs) > 0 {
		sb.CharacterDescs = s.buildCharacterDescs(charIDs)
		// Build voice styles
		var chars []models.Character
		if err := s.db.Where("id IN ?", charIDs).Find(&chars).Error; err == nil {
			var voiceParts []string
			for _, c := range chars {
				if c.VoiceStyle != nil && *c.VoiceStyle != "" {
					voiceParts = append(voiceParts, c.Name+": "+*c.VoiceStyle)
				}
			}
			sb.CharacterVoiceStyles = strings.Join(voiceParts, "; ")
		}
	}

	// Load narrator voice profile from Drama
	var episode models.Episode
	if err := s.db.Select("drama_id").First(&episode, storyboard.EpisodeID).Error; err == nil {
		sb.NarratorVoiceProfile = s.loadNarratorVoiceProfile(episode.DramaID)
	}

	// 只有在不是 AI 生成的情况下，才根据内容自动生成简单的 video_prompt
	// 或者如果 updates 中明确传了新的 video_prompt，则使用那个（见上文）
	if _, hasNewPrompt := updates["video_prompt"]; !hasNewPrompt {
		if storyboard.VideoPromptSource != "ai" {
			videoPrompt := s.generateVideoPrompt(sb)

			// Safeguard: ensure narrator text is in video prompt for narrator shots
			if sb.NarratorScript != "" && !strings.Contains(videoPrompt, "Narration") {
				s.log.Warnw("UpdateStoryboard: Narrator text missing from video prompt — injecting safeguard",
					"storyboard_id", storyboardID,
					"narrator_script_len", len(sb.NarratorScript),
					"audio_mode", sb.AudioMode,
				)
				videoPrompt = strings.Replace(videoPrompt,
					"The character's mouth is completely closed, silent scene. --no talking, speaking, moving lips",
					fmt.Sprintf("Narration (voice-over): %s. The character's mouth is strictly closed, silent expression, purely visual acting, no speaking, voiceover scene. --no talking, speaking, moving lips, open mouth, chatting", sb.NarratorScript),
					1)
			}

			updateData["video_prompt"] = videoPrompt
			updateData["video_prompt_source"] = "auto"
		}
	} else {
		// 如果手动传了，更新 source 为 manual
		updateData["video_prompt_source"] = "manual"
	}

	// 更新数据库
	if err := s.db.Model(&storyboard).Updates(updateData).Error; err != nil {
		return fmt.Errorf("failed to update storyboard: %w", err)
	}

	s.log.Infow("Storyboard updated successfully",
		"storyboard_id", storyboardID,
		"fields_updated", len(updateData))

	return nil
}
