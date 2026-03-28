package handlers

import (
	"github.com/drama-generator/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// GetCharacterFullPrompt returns the full composed prompt that would be sent to the AI for character image generation
func (h *CharacterLibraryHandler) GetCharacterFullPrompt(c *gin.Context) {
	characterID := c.Param("id")

	fullPrompt, err := h.libraryService.BuildCharacterFullPrompt(characterID)
	if err != nil {
		if err.Error() == "character not found" {
			response.NotFound(c, "角色不存在")
			return
		}
		if err.Error() == "unauthorized" {
			response.Forbidden(c, "无权限")
			return
		}
		h.log.Errorw("Failed to build character full prompt", "error", err)
		response.InternalError(c, "获取完整提示词失败")
		return
	}

	response.Success(c, gin.H{
		"prompt": fullPrompt,
	})
}
