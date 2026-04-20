package handlers

import (
	"strconv"

	"github.com/drama-generator/backend/application/services"
	"github.com/drama-generator/backend/infrastructure/storage"
	"github.com/drama-generator/backend/pkg/logger"
	"github.com/drama-generator/backend/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type VideoReviewHandler struct {
	reviewService *services.VideoReviewService
	log           *logger.Logger
}

func NewVideoReviewHandler(db *gorm.DB, log *logger.Logger, aiService *services.AIService, taskService *services.TaskService, localStorage *storage.LocalStorage) *VideoReviewHandler {
	return &VideoReviewHandler{
		reviewService: services.NewVideoReviewService(db, log, aiService, taskService, localStorage),
		log:           log,
	}
}

// ReviewVideo triggers an async video review.
// POST /videos/:id/review
func (h *VideoReviewHandler) ReviewVideo(c *gin.Context) {
	videoGenID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid video ID")
		return
	}

	taskID, err := h.reviewService.ReviewVideoAsync(uint(videoGenID))
	if err != nil {
		h.log.Errorw("Failed to start video review", "videoGenID", videoGenID, "error", err)
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"task_id": taskID,
		"status":  "processing",
		"message": "Video review started",
	})
}

// GetVideoReview returns the latest review for a video generation.
// GET /videos/:id/review
func (h *VideoReviewHandler) GetVideoReview(c *gin.Context) {
	videoGenID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid video ID")
		return
	}

	review, err := h.reviewService.GetLatestReview(uint(videoGenID))
	if err != nil {
		h.log.Errorw("Failed to get video review", "videoGenID", videoGenID, "error", err)
		response.InternalError(c, err.Error())
		return
	}

	if review == nil {
		// No review exists yet — return empty success (not 404)
		response.Success(c, nil)
		return
	}

	response.Success(c, review)
}
