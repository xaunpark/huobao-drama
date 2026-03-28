package handlers

import (
	"github.com/drama-generator/backend/application/services"
	"github.com/drama-generator/backend/pkg/config"
	"github.com/drama-generator/backend/pkg/logger"
	"github.com/drama-generator/backend/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RapidCutHandler struct {
	rapidCutService *services.RapidCutService
	taskService     *services.TaskService
	log             *logger.Logger
}

func NewRapidCutHandler(db *gorm.DB, cfg *config.Config, log *logger.Logger) *RapidCutHandler {
	return &RapidCutHandler{
		rapidCutService: services.NewRapidCutService(db, cfg, log),
		taskService:     services.NewTaskService(db, log),
		log:             log,
	}
}

// GenerateRapidCut generates rapid cut production shots for an episode
func (h *RapidCutHandler) GenerateRapidCut(c *gin.Context) {
	episodeID := c.Param("episode_id")

	var req struct {
		Model string `json:"model"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		req.Model = ""
	}

	taskID, err := h.rapidCutService.GenerateRapidCut(episodeID, req.Model)
	if err != nil {
		h.log.Errorw("Failed to generate rapid cut", "error", err, "episode_id", episodeID)
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"task_id": taskID,
		"status":  "pending",
		"message": "Rapid cut generation task created, processing in background...",
	})
}

// DeleteRapidCut removes all production shots for an episode (back to standard mode)
func (h *RapidCutHandler) DeleteRapidCut(c *gin.Context) {
	episodeID := c.Param("episode_id")

	if err := h.rapidCutService.DeleteRapidCut(episodeID); err != nil {
		h.log.Errorw("Failed to delete rapid cut", "error", err, "episode_id", episodeID)
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"message": "Rapid cut production shots deleted, back to standard mode",
	})
}

// GetRapidCutStatus checks if an episode has rapid cut production shots
func (h *RapidCutHandler) GetRapidCutStatus(c *gin.Context) {
	episodeID := c.Param("episode_id")

	hasRapidCut, err := h.rapidCutService.HasRapidCut(episodeID)
	if err != nil {
		h.log.Errorw("Failed to check rapid cut status", "error", err, "episode_id", episodeID)
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"has_rapid_cut": hasRapidCut,
	})
}
