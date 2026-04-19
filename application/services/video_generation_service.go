package services

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	models "github.com/drama-generator/backend/domain/models"
	"github.com/drama-generator/backend/infrastructure/external/ffmpeg"
	"github.com/drama-generator/backend/infrastructure/storage"
	"github.com/drama-generator/backend/pkg/logger"
	"github.com/drama-generator/backend/pkg/utils"
	"github.com/drama-generator/backend/pkg/video"
	"gorm.io/gorm"
)

type VideoGenerationService struct {
	db              *gorm.DB
	transferService *ResourceTransferService
	log             *logger.Logger
	localStorage    *storage.LocalStorage
	aiService       *AIService
	ffmpeg          *ffmpeg.FFmpeg
	promptI18n      *PromptI18n
}

func NewVideoGenerationService(db *gorm.DB, transferService *ResourceTransferService, localStorage *storage.LocalStorage, aiService *AIService, log *logger.Logger, promptI18n *PromptI18n) *VideoGenerationService {
	service := &VideoGenerationService{
		db:              db,
		localStorage:    localStorage,
		transferService: transferService,
		aiService:       aiService,
		log:             log,
		ffmpeg:          ffmpeg.NewFFmpeg(log),
		promptI18n:      promptI18n,
	}

	go service.RecoverPendingTasks()

	return service
}

type GenerateVideoRequest struct {
	StoryboardID *uint  `json:"storyboard_id"`
	DramaID      string `json:"drama_id" binding:"required"`
	ImageGenID   *uint  `json:"image_gen_id"`

	// 参考图模式：single, first_last, multiple, none
	ReferenceMode string `json:"reference_mode"`

	// 单图模式
	ImageURL       string  `json:"image_url"`
	ImageLocalPath *string `json:"image_local_path"` // 单图模式的本地路径

	// 首尾帧模式
	FirstFrameURL       *string `json:"first_frame_url"`
	FirstFrameLocalPath *string `json:"first_frame_local_path"` // 首帧本地路径
	LastFrameURL        *string `json:"last_frame_url"`
	LastFrameLocalPath  *string `json:"last_frame_local_path"` // 尾帧本地路径

	// 多图模式
	ReferenceImageURLs []string `json:"reference_image_urls"`

	Prompt       string  `json:"prompt" binding:"required,min=1"`
	Provider     string  `json:"provider"`
	Model        string  `json:"model"`
	Duration     *int    `json:"duration"`
	FPS          *int    `json:"fps"`
	AspectRatio  *string `json:"aspect_ratio"`
	Style        *string `json:"style"`
	MotionLevel  *int    `json:"motion_level"`
	CameraMotion *string `json:"camera_motion"`
	Seed         *int64  `json:"seed"`
	GenerationMode string  `json:"generation_mode"` // r2v, i2v, t2v
}

func (s *VideoGenerationService) GenerateVideo(request *GenerateVideoRequest) (*models.VideoGeneration, error) {
	if request.StoryboardID != nil {
		var storyboard models.Storyboard
		if err := s.db.Preload("Episode").Where("id = ?", *request.StoryboardID).First(&storyboard).Error; err != nil {
			return nil, fmt.Errorf("storyboard not found")
		}
		if fmt.Sprintf("%d", storyboard.Episode.DramaID) != request.DramaID {
			return nil, fmt.Errorf("storyboard does not belong to drama")
		}
	}

	if request.ImageGenID != nil {
		var imageGen models.ImageGeneration
		if err := s.db.Where("id = ?", *request.ImageGenID).First(&imageGen).Error; err != nil {
			return nil, fmt.Errorf("image generation not found")
		}
	}

	provider := request.Provider
	if provider == "" {
		provider = "doubao"
	}

	dramaID, _ := strconv.ParseUint(request.DramaID, 10, 32)

	videoGen := &models.VideoGeneration{
		StoryboardID: request.StoryboardID,
		DramaID:      uint(dramaID),
		ImageGenID:   request.ImageGenID,
		Provider:     provider,
		Prompt:       request.Prompt,
		Model:        request.Model,
		Duration:     request.Duration,
		FPS:          request.FPS,
		AspectRatio:  request.AspectRatio,
		Style:        request.Style,
		MotionLevel:  request.MotionLevel,
		CameraMotion: request.CameraMotion,
		Seed:         request.Seed,
		Status:       models.VideoStatusPending,
	}

	if request.GenerationMode != "" {
		videoGen.GenerationMode = &request.GenerationMode
	} else if request.ReferenceMode == "multiple" {
		r2vMode := "direct_r2v"
		videoGen.GenerationMode = &r2vMode
	} else {
		defaultMode := "shot_i2v"
		videoGen.GenerationMode = &defaultMode
	}

	// 根据参考图模式处理不同的参数
	if request.ReferenceMode != "" {
		videoGen.ReferenceMode = &request.ReferenceMode
	}

	switch request.ReferenceMode {
	case "single":
		// 单图模式 - 优先使用 local_path
		if request.ImageLocalPath != nil && *request.ImageLocalPath != "" {
			videoGen.ImageURL = request.ImageLocalPath
		} else if request.ImageURL != "" {
			videoGen.ImageURL = &request.ImageURL
		}
	case "first_last":
		// 首尾帧模式 - 优先使用 local_path
		if request.FirstFrameLocalPath != nil && *request.FirstFrameLocalPath != "" {
			videoGen.FirstFrameURL = request.FirstFrameLocalPath
		} else if request.FirstFrameURL != nil {
			videoGen.FirstFrameURL = request.FirstFrameURL
		}
		if request.LastFrameLocalPath != nil && *request.LastFrameLocalPath != "" {
			videoGen.LastFrameURL = request.LastFrameLocalPath
		} else if request.LastFrameURL != nil {
			videoGen.LastFrameURL = request.LastFrameURL
		}
	case "multiple":
		// 多图模式
		if len(request.ReferenceImageURLs) > 0 {
			referenceImagesJSON, err := json.Marshal(request.ReferenceImageURLs)
			if err == nil {
				referenceImagesStr := string(referenceImagesJSON)
				videoGen.ReferenceImageURLs = &referenceImagesStr
			}
		}
	case "none":
		// 无参考图，纯文本生成
	default:
		// 向后兼容：如果没有指定模式，根据提供的参数自动判断
		if request.ImageURL != "" {
			videoGen.ImageURL = &request.ImageURL
			mode := "single"
			videoGen.ReferenceMode = &mode
		} else if request.FirstFrameURL != nil || request.LastFrameURL != nil {
			videoGen.FirstFrameURL = request.FirstFrameURL
			videoGen.LastFrameURL = request.LastFrameURL
			mode := "first_last"
			videoGen.ReferenceMode = &mode
		} else if len(request.ReferenceImageURLs) > 0 {
			referenceImagesJSON, err := json.Marshal(request.ReferenceImageURLs)
			if err == nil {
				referenceImagesStr := string(referenceImagesJSON)
				videoGen.ReferenceImageURLs = &referenceImagesStr
				mode := "multiple"
				videoGen.ReferenceMode = &mode
			}
		}
	}

	if err := s.db.Create(videoGen).Error; err != nil {
		return nil, fmt.Errorf("failed to create record: %w", err)
	}

	// Start background goroutine to process video generation asynchronously
	// This allows the API to return immediately while video generation happens in background
	// CRITICAL: The goroutine will handle all video generation logic including API calls and polling
	go s.ProcessVideoGeneration(videoGen.ID)

	return videoGen, nil
}

func (s *VideoGenerationService) ProcessVideoGeneration(videoGenID uint) {
	var videoGen models.VideoGeneration
	if err := s.db.First(&videoGen, videoGenID).Error; err != nil {
		s.log.Errorw("Failed to load video generation", "error", err, "id", videoGenID)
		return
	}

	// 获取drama的style信息
	var drama models.Drama
	if err := s.db.First(&drama, videoGen.DramaID).Error; err != nil {
		s.log.Warnw("Failed to load drama for style", "error", err, "drama_id", videoGen.DramaID)
	}

	s.db.Model(&videoGen).Update("status", models.VideoStatusProcessing)

	client, err := s.getVideoClient(videoGen.Provider, videoGen.Model)
	if err != nil {
		s.log.Errorw("Failed to get video client", "error", err, "provider", videoGen.Provider, "model", videoGen.Model)
		s.updateVideoGenError(videoGenID, err.Error())
		return
	}

	s.log.Infow("Starting video generation", "id", videoGenID, "prompt", videoGen.Prompt, "provider", videoGen.Provider)

	var opts []video.VideoOption
	if videoGen.Model != "" {
		opts = append(opts, video.WithModel(videoGen.Model))
	}
	if videoGen.Duration != nil {
		opts = append(opts, video.WithDuration(*videoGen.Duration))
	}
	if videoGen.FPS != nil {
		opts = append(opts, video.WithFPS(*videoGen.FPS))
	}
	if videoGen.AspectRatio != nil {
		opts = append(opts, video.WithAspectRatio(*videoGen.AspectRatio))
	}
	if videoGen.Style != nil {
		opts = append(opts, video.WithStyle(*videoGen.Style))
	}
	if videoGen.MotionLevel != nil {
		opts = append(opts, video.WithMotionLevel(*videoGen.MotionLevel))
	}
	if videoGen.CameraMotion != nil {
		opts = append(opts, video.WithCameraMotion(*videoGen.CameraMotion))
	}
	if videoGen.Seed != nil {
		opts = append(opts, video.WithSeed(*videoGen.Seed))
	}

	if videoGen.GenerationMode != nil {
		opts = append(opts, video.WithGenerationMode(*videoGen.GenerationMode))
	} else if videoGen.ReferenceMode != nil {
		opts = append(opts, video.WithGenerationMode(*videoGen.ReferenceMode))
	}

	bypassBase64 := videoGen.Provider == "flow-tool"

	// 根据参考图模式添加相应的选项，并将本地图片转换为base64
	if videoGen.ReferenceMode != nil {
		switch *videoGen.ReferenceMode {
		case "first_last":
			// 首尾帧模式
			if videoGen.FirstFrameURL != nil {
				if bypassBase64 {
					opts = append(opts, video.WithFirstFrame(*videoGen.FirstFrameURL))
				} else {
					firstFrameBase64, err := s.convertImageToBase64(*videoGen.FirstFrameURL)
					if err != nil {
						s.log.Warnw("Failed to convert first frame to base64, using original URL", "error", err)
						opts = append(opts, video.WithFirstFrame(*videoGen.FirstFrameURL))
					} else {
						opts = append(opts, video.WithFirstFrame(firstFrameBase64))
					}
				}
			}
			if videoGen.LastFrameURL != nil {
				if bypassBase64 {
					opts = append(opts, video.WithLastFrame(*videoGen.LastFrameURL))
				} else {
					lastFrameBase64, err := s.convertImageToBase64(*videoGen.LastFrameURL)
					if err != nil {
						s.log.Warnw("Failed to convert last frame to base64, using original URL", "error", err)
						opts = append(opts, video.WithLastFrame(*videoGen.LastFrameURL))
					} else {
						opts = append(opts, video.WithLastFrame(lastFrameBase64))
					}
				}
			}
		case "multiple":
			// 多图模式
			if videoGen.ReferenceImageURLs != nil {
				var imageURLs []string
				if err := json.Unmarshal([]byte(*videoGen.ReferenceImageURLs), &imageURLs); err == nil {
					if bypassBase64 {
						opts = append(opts, video.WithReferenceImages(imageURLs))
					} else {
						var base64Images []string
						for _, imgURL := range imageURLs {
							base64Img, err := s.convertImageToBase64(imgURL)
							if err != nil {
								s.log.Warnw("Failed to convert reference image to base64, using original URL", "error", err, "url", imgURL)
								base64Images = append(base64Images, imgURL)
							} else {
								base64Images = append(base64Images, base64Img)
							}
						}
						opts = append(opts, video.WithReferenceImages(base64Images))
					}
				}
			}
		}
	}

	// 构造imageURL参数（单图模式使用，其他模式传空字符串）
	// 优先使用本地路径，避免外部URL过期导致 PUBLIC_ERROR_IP_INPUT_IMAGE
	imageURL := ""
	if videoGen.ImageURL != nil {
		resolvedImagePath := *videoGen.ImageURL

		// 如果存储的是外部URL（非本地路径），尝试从关联的ImageGeneration查找local_path
		if strings.HasPrefix(resolvedImagePath, "http://") || strings.HasPrefix(resolvedImagePath, "https://") {
			if videoGen.ImageGenID != nil {
				var imgGen models.ImageGeneration
				if err := s.db.Select("local_path, image_url").Where("id = ?", *videoGen.ImageGenID).First(&imgGen).Error; err == nil {
					if imgGen.LocalPath != nil && *imgGen.LocalPath != "" {
						s.log.Infow("Using local_path from ImageGeneration instead of external URL",
							"id", videoGenID,
							"image_gen_id", *videoGen.ImageGenID,
							"local_path", *imgGen.LocalPath,
							"original_url_prefix", resolvedImagePath[:min(len(resolvedImagePath), 80)])
						resolvedImagePath = *imgGen.LocalPath
					}
				}
			}
		}

		s.log.Infow("Resolving image for video generation",
			"id", videoGenID,
			"resolved_path", resolvedImagePath[:min(len(resolvedImagePath), 80)],
			"is_local", !strings.HasPrefix(resolvedImagePath, "http"))

		if bypassBase64 {
			imageURL = resolvedImagePath
		} else {
			base64Image, err := s.convertImageToBase64(resolvedImagePath)
			if err != nil {
				s.log.Errorw("Failed to convert image to base64",
					"error", err,
					"id", videoGenID,
					"resolved_path", resolvedImagePath[:min(len(resolvedImagePath), 80)])
				imageURL = resolvedImagePath
			} else {
				imageURL = base64Image
			}
		}
	}

	// 构建完整的提示词：风格提示词 + 约束提示词 + 用户提示词
	prompt := videoGen.Prompt

	// TODO(Phase 3): referenceMode detection below will be repurposed for
	// per-shot video_style prepending when StyleDistillService is implemented.
	// For now, the detection is disabled since video_constraint injection was removed.
	// See: plans/shot-style-distill.md (Phase 3: Integration, Task 11)
	// referenceMode := "single"
	// if videoGen.ReferenceMode != nil { referenceMode = *videoGen.ReferenceMode }
	// ... frame type detection for action/first/last/key ...

	// NOTE: video_constraint 不再在此处注入。video_constraint 是为 LLM 设计的系统指令，
	// 不应直接注入到视频生成 API（Kling/Doubao）。将在 Phase 3 中被 per-shot 的
	// distilled video_style 替代。
	// See: plans/shot-style-distill.md (Phase 1: Cleanup)

	// 打印完整的提示词信息
	s.log.Infow("Video generation prompts",
		"id", videoGenID,
		"user_prompt", videoGen.Prompt,
		"final_prompt", prompt)

	result, err := client.GenerateVideo(imageURL, prompt, opts...)
	if err != nil {
		s.log.Errorw("Video generation API call failed", "error", err, "id", videoGenID)
		s.updateVideoGenError(videoGenID, err.Error())
		return
	}

	// CRITICAL FIX: Validate TaskID before starting polling goroutine
	// Empty TaskID would cause polling to fail silently or cause issues
	if result.TaskID != "" {
		s.db.Model(&videoGen).Updates(map[string]interface{}{
			"task_id": result.TaskID,
			"status":  models.VideoStatusProcessing,
		})
		// Start background goroutine to poll task status
		// This allows the API to return immediately while video generation continues asynchronously
		// The goroutine will poll until completion, failure, or timeout (max 300 attempts * 10s = 50 minutes)
		go s.pollTaskStatus(videoGenID, result.TaskID, videoGen.Provider, videoGen.Model)
		return
	}

	if result.VideoURL != "" {
		s.completeVideoGeneration(videoGenID, result.VideoURL, &result.Duration, &result.Width, &result.Height, nil)
		return
	}

	s.updateVideoGenError(videoGenID, "no task ID or video URL returned")
}

func (s *VideoGenerationService) pollTaskStatus(videoGenID uint, taskID string, provider string, model string) {
	// CRITICAL FIX: Validate taskID parameter to prevent invalid API calls
	// Empty taskID would cause unnecessary API calls and potential errors
	if taskID == "" {
		s.log.Errorw("Invalid empty taskID for polling", "video_gen_id", videoGenID)
		s.updateVideoGenError(videoGenID, "invalid task ID for polling")
		return
	}

	client, err := s.getVideoClient(provider, model)
	if err != nil {
		s.log.Errorw("Failed to get video client for polling", "error", err)
		s.updateVideoGenError(videoGenID, "failed to get video client")
		return
	}

	// Polling configuration: max 200 attempts with 3 second intervals
	// Total maximum polling time: 200 * 3s = 10 minutes (600s)
	// This prevents infinite polling if the task never completes
	maxAttempts := 200
	interval := 3 * time.Second

	// Track consecutive polls where error is present but status isn't terminal.
	// Flow-Tool may return PENDING with error for 401 auto-requeue (transient).
	// We tolerate up to 3 consecutive error polls before treating as permanent failure.
	consecutiveErrors := 0
	const maxConsecutiveErrors = 3

	// Track consecutive polls where status is PENDING (job stuck in queue, never dispatched).
	// Flow-Tool may queue jobs but never dispatch them if all workers are busy/failed.
	// After ~45s of pure PENDING, we treat the job as stalled.
	consecutivePending := 0
	const maxConsecutivePending = 15 // 15 × 3s = ~45 seconds

	for attempt := 0; attempt < maxAttempts; attempt++ {
		var videoGen models.VideoGeneration
		if err := s.db.First(&videoGen, videoGenID).Error; err != nil {
			s.log.Errorw("Failed to load video generation", "error", err, "id", videoGenID)
			return
		}

		// Initial sleep for Upscaling (shorter than generation to detect failures faster)
		if attempt == 0 && videoGen.Status == models.VideoStatusUpscaling {
			time.Sleep(10 * time.Second)
		} else {
			// Sleep before each poll attempt to avoid overwhelming the API
			time.Sleep(interval)
		}

		// CRITICAL FIX: Check if status was manually changed (e.g., cancelled by user)
		// If status is no longer "processing" or "upscaling", stop polling to avoid unnecessary API calls
		// This prevents polling when the task has been cancelled or failed externally
		if videoGen.Status != models.VideoStatusProcessing && videoGen.Status != models.VideoStatusUpscaling {
			s.log.Infow("Video generation status changed, stopping poll", "id", videoGenID, "status", videoGen.Status)
			return
		}

		// Poll the video generation API for task status
		// Continue polling on transient errors (network issues, temporary API failures)
		// Only stop on permanent errors or task completion
		result, err := client.GetTaskStatus(taskID)
		if err != nil {
			s.log.Errorw("Failed to get task status", "error", err, "task_id", taskID, "attempt", attempt+1)
			
			// If the task is permanently gone (e.g. 404 Not Found) or server is permanently crashing (500), stop polling immediately!
			if strings.Contains(strings.ToLower(err.Error()), "status 404") || strings.Contains(strings.ToLower(err.Error()), "not found") {
				s.updateVideoGenError(videoGenID, "Task was lost or deleted on AI server (404 Not Found).")
				return
			}
			
			// If Flow-tool hits 500 Internal error while checking task, it usually means the task is corrupted.
			if strings.Contains(strings.ToLower(err.Error()), "status 500") || strings.Contains(strings.ToLower(err.Error()), "internal server error") {
				s.updateVideoGenError(videoGenID, "AI Server encountered an internal error while processing this task (500).")
				return
			}
			
			// Continue polling on other errors - might be transient network issue
			// Will eventually timeout after maxAttempts if error persists
			continue
		}

		// Check if task completed successfully
		// CRITICAL FIX: Validate that video URL exists when task is marked as completed
		// Some APIs may mark task as completed but fail to provide the video URL
		if result.Completed {
			if result.VideoURL != "" {
				// Successfully completed with video URL - download and update database
				s.completeVideoGeneration(videoGenID, result.VideoURL, &result.Duration, &result.Width, &result.Height, nil)
				return
			}
			// Task marked as completed but no video URL - this is an error condition
			s.updateVideoGenError(videoGenID, "task completed but no video URL")
			return
		}

		// Check if task failed with an error message
		if result.Error != "" {
			normalizedErrStatus := strings.ToUpper(strings.TrimSpace(result.Status))

			// If Flow-Tool explicitly marks job as FAILED/ERROR → terminal, fail immediately
			if normalizedErrStatus == "FAILED" || normalizedErrStatus == "ERROR" {
				s.log.Errorw("Job failed (terminal status from Flow-Tool)",
					"id", videoGenID,
					"flow_status", result.Status,
					"error", result.Error)
				s.updateVideoGenError(videoGenID, result.Error)
				return
			}

			// Status is NOT terminal (e.g. PENDING with error = 401 auto-requeue).
			// Allow a few retries before giving up.
			consecutiveErrors++
			consecutivePending = 0 // Reset pending counter — job is no longer just waiting
			if consecutiveErrors >= maxConsecutiveErrors {
				s.log.Errorw("Job error persisted across retries, marking as failed",
					"id", videoGenID,
					"flow_status", result.Status,
					"error", result.Error,
					"consecutive_errors", consecutiveErrors)
				s.updateVideoGenError(videoGenID, result.Error)
				return
			}
			s.log.Warnw("Job has error but status is not terminal, retrying",
				"id", videoGenID,
				"flow_status", result.Status,
				"error", result.Error,
				"consecutive_errors", consecutiveErrors,
				"attempt", attempt+1)
			continue
		}

		// No error — reset consecutive error counter
		consecutiveErrors = 0

		// Detect stale PENDING: job stuck in queue, never dispatched by Flow-Tool
		// Flow-Tool returns "PENDING" when job is queued but not yet assigned to a worker.
		// If all workers are busy or failed, the job will stay PENDING indefinitely.
		normalizedStatus := strings.ToUpper(strings.TrimSpace(result.Status))
		if normalizedStatus == "PENDING" {
			consecutivePending++
			if consecutivePending >= maxConsecutivePending {
				s.log.Errorw("Job stalled in PENDING queue, marking as failed",
					"id", videoGenID,
					"flow_status", result.Status,
					"consecutive_pending", consecutivePending,
					"stalled_seconds", consecutivePending*int(interval.Seconds()))
				s.updateVideoGenError(videoGenID, fmt.Sprintf("Job stalled in queue (PENDING for %d polls / ~%ds). Flow-Tool dispatcher may be stuck.", consecutivePending, consecutivePending*int(interval.Seconds())))
				return
			}
		} else {
			// Status is PROCESSING or something else active — job is being worked on, reset counter
			consecutivePending = 0
		}

		// Task still in progress - log with Flow-Tool status for debugging
		s.log.Infow("Video generation in progress",
			"id", videoGenID,
			"flow_status", result.Status,
			"attempt", attempt+1,
			"max_attempts", maxAttempts)
	}

	// CRITICAL FIX: Handle polling timeout gracefully
	// After maxAttempts (10 minutes), mark task as failed if still not completed
	// This prevents indefinite polling and resource waste
	s.updateVideoGenError(videoGenID, fmt.Sprintf("polling timeout after %d attempts (%.1f minutes)", maxAttempts, float64(maxAttempts*int(interval))/60.0))
}

func (s *VideoGenerationService) completeVideoGeneration(videoGenID uint, videoURL string, duration *int, width *int, height *int, firstFrameURL *string) {
	var localVideoPath *string

	// Check if it was upscaling to selectively bypass ffmpeg later
	isUpscale := false
	var oldVideoGen models.VideoGeneration
	if err := s.db.First(&oldVideoGen, videoGenID).Error; err == nil {
		if oldVideoGen.Status == models.VideoStatusUpscaling {
			isUpscale = true
		}
	}

	// 下载视频到本地存储并保存相对路径到数据库
	if s.localStorage != nil && videoURL != "" {
		downloadResult, err := s.localStorage.DownloadFromURLWithPath(videoURL, "videos")
		if err != nil {
			s.log.Warnw("Failed to download video to local storage",
				"error", err,
				"id", videoGenID,
				"original_url", videoURL)
		} else {
			localVideoPath = &downloadResult.RelativePath
			s.log.Infow("Video downloaded to local storage",
				"id", videoGenID,
				"original_url", videoURL,
				"local_path", downloadResult.RelativePath)
		}
	}

	if localVideoPath != nil {
		absPath := s.localStorage.GetAbsolutePath(*localVideoPath)
		fileInfo, err := os.Stat(absPath)

		// Set dynamic file size limit based on task type
		var minSize int64 = 51200 // 50 KB (Ngưỡng an toàn cho video gen ngắn)
		if isUpscale {
			minSize = 1048576 // 1 MB (Ngưỡng đặc trị cho video upscale 1080p)
		}

		if err == nil && fileInfo.Size() < minSize {
			s.log.Errorw("Downloaded video file is too small (likely corrupted)", "id", videoGenID, "size", fileInfo.Size(), "local_path", *localVideoPath, "is_upscale", isUpscale)
			
			errorMsg := "Tác vụ trả về file trống hoặc lỗi định dạng (<50KB)"
			if isUpscale {
				errorMsg = "Tác vụ trả về file trống hoặc lỗi định dạng (<1MB)"
			}
			
			s.updateVideoGenError(videoGenID, errorMsg)
			os.Remove(absPath)
			return
		}
	}

	// 如果视频已下载到本地，探测真实时长
	// 特别是当 AI 服务返回的 duration 为 0 或 nil 时，必须探测
	// Bỏ qua FFmpeg nếu tác vụ là Upscale
	shouldProbe := localVideoPath != nil && s.ffmpeg != nil && (duration == nil || *duration == 0) && !isUpscale
	if shouldProbe {
		absPath := s.localStorage.GetAbsolutePath(*localVideoPath)
		if probedDuration, err := s.ffmpeg.GetVideoDuration(absPath); err == nil {
			// 转换为整数秒（向上取整）
			durationInt := int(probedDuration + 0.5)
			duration = &durationInt
			s.log.Infow("Probed video duration (was 0 or nil)",
				"id", videoGenID,
				"duration_seconds", durationInt,
				"duration_float", probedDuration)
		} else {
			s.log.Errorw("Failed to probe video duration, duration will be 0",
				"error", err,
				"id", videoGenID,
				"local_path", *localVideoPath)
		}
	} else if localVideoPath != nil && s.ffmpeg != nil && duration != nil && *duration > 0 && !isUpscale {
		// 即使有 duration，也验证一下（可选）
		absPath := s.localStorage.GetAbsolutePath(*localVideoPath)
		if probedDuration, err := s.ffmpeg.GetVideoDuration(absPath); err == nil {
			durationInt := int(probedDuration + 0.5)
			if durationInt != *duration {
				s.log.Warnw("Probed duration differs from provided duration",
					"id", videoGenID,
					"provided", *duration,
					"probed", durationInt)
				// 使用探测到的时长（更准确）
				duration = &durationInt
			}
		}
	}

	// 下载首帧图片到本地存储（仅用于缓存，不更新数据库）
	if firstFrameURL != nil && *firstFrameURL != "" && s.localStorage != nil {
		_, err := s.localStorage.DownloadFromURL(*firstFrameURL, "video_frames")
		if err != nil {
			s.log.Warnw("Failed to download first frame to local storage",
				"error", err,
				"id", videoGenID,
				"original_url", *firstFrameURL)
		} else {
			s.log.Infow("First frame downloaded to local storage for caching",
				"id", videoGenID,
				"original_url", *firstFrameURL)
		}
	}

	// 数据库中保存原始URL和本地路径
	updates := map[string]interface{}{
		"video_url":  videoURL,
		"local_path": localVideoPath,
	}

	targetStatus := models.VideoStatusCompleted
	if isUpscale {
		targetStatus = models.VideoStatusUpscaled
		updates["is_upscaled"] = true // Keep for backward compatibility
	}
	updates["status"] = targetStatus

	// 只有当 duration 大于 0 时才保存，避免保存无效的 0 值
	if duration != nil && *duration > 0 {
		updates["duration"] = *duration
	}
	if width != nil {
		updates["width"] = *width
	}
	if height != nil {
		updates["height"] = *height
	}
	if firstFrameURL != nil {
		updates["first_frame_url"] = *firstFrameURL
	}

	if err := s.db.Model(&models.VideoGeneration{}).Where("id = ?", videoGenID).Updates(updates).Error; err != nil {
		s.log.Errorw("Failed to update video generation", "error", err, "id", videoGenID)
		return
	}

	var videoGen models.VideoGeneration
	if err := s.db.First(&videoGen, videoGenID).Error; err == nil {
		if videoGen.StoryboardID != nil {
			// 更新 Storyboard 的 video_url 和 duration
			storyboardUpdates := map[string]interface{}{
				"video_url": videoURL,
			}
			// 只有当 duration 大于 0 时才更新，避免用无效的 0 值覆盖
			if duration != nil && *duration > 0 {
				storyboardUpdates["duration"] = *duration
			}
			if err := s.db.Model(&models.Storyboard{}).Where("id = ?", *videoGen.StoryboardID).Updates(storyboardUpdates).Error; err != nil {
				s.log.Warnw("Failed to update storyboard", "storyboard_id", *videoGen.StoryboardID, "error", err)
			} else {
				s.log.Infow("Updated storyboard with video info", "storyboard_id", *videoGen.StoryboardID, "duration", duration)
			}
		}
	}

	s.log.Infow("Video generation completed", "id", videoGenID, "url", videoURL, "duration", duration)
}

func (s *VideoGenerationService) updateVideoGenError(videoGenID uint, errorMsg string) {
	var oldVideoGen models.VideoGeneration
	targetStatus := models.VideoStatusFailed
	if err := s.db.First(&oldVideoGen, videoGenID).Error; err == nil {
		if oldVideoGen.Status == models.VideoStatusUpscaling {
			targetStatus = models.VideoStatusUpscaleFailed
		}
	}

	if err := s.db.Model(&models.VideoGeneration{}).Where("id = ?", videoGenID).Updates(map[string]interface{}{
		"status":    targetStatus,
		"error_msg": errorMsg,
	}).Error; err != nil {
		s.log.Errorw("Failed to update video generation error", "error", err, "id", videoGenID)
	}
}

func (s *VideoGenerationService) getVideoClient(provider string, modelName string) (video.VideoClient, error) {
	// 根据模型名称获取AI配置
	var config *models.AIServiceConfig
	var err error

	if modelName != "" {
		config, err = s.aiService.GetConfigForModel("video", modelName)
		if err != nil {
			s.log.Warnw("Failed to get config for model, using default", "model", modelName, "error", err)
			config, err = s.aiService.GetDefaultConfig("video")
			if err != nil {
				return nil, fmt.Errorf("no video AI config found: %w", err)
			}
		}
	} else {
		config, err = s.aiService.GetDefaultConfig("video")
		if err != nil {
			return nil, fmt.Errorf("no video AI config found: %w", err)
		}
	}

	// 使用配置中的信息创建客户端
	baseURL := config.BaseURL
	apiKey := config.APIKey
	model := modelName
	if model == "" && len(config.Model) > 0 {
		model = config.Model[0]
	}

	// 根据配置中的 provider 创建对应的客户端
	var endpoint string
	var queryEndpoint string

	switch config.Provider {
	case "chatfire":
		endpoint = "/video/generations"
		queryEndpoint = "/video/task/{taskId}"
		return video.NewChatfireClient(baseURL, apiKey, model, endpoint, queryEndpoint), nil
	case "doubao", "volcengine", "volces":
		endpoint = ""
		queryEndpoint = ""
		return video.NewVolcesArkClient(baseURL, apiKey, model, endpoint, queryEndpoint), nil
	case "openai":
		// OpenAI Sora 使用 /v1/videos 端点
		return video.NewOpenAISoraClient(baseURL, apiKey, model), nil
	case "runway":
		return video.NewRunwayClient(baseURL, apiKey, model), nil
	case "pika":
		return video.NewPikaClient(baseURL, apiKey, model), nil
	case "minimax":
		return video.NewMinimaxClient(baseURL, apiKey, model), nil
	case "flowtool", "flow_tool":
		return video.NewFlowToolVideoClient(baseURL), nil
	default:
		return nil, fmt.Errorf("unsupported video provider: %s", provider)
	}
}

func (s *VideoGenerationService) RecoverPendingTasks() {
	var pendingVideos []models.VideoGeneration
	// Query for pending tasks with non-empty task_id
	// Note: Using IS NOT NULL and != '' to ensure we only get valid task IDs
	if err := s.db.Where("status IN ? AND task_id IS NOT NULL AND task_id != ''", []models.VideoStatus{models.VideoStatusProcessing, models.VideoStatusUpscaling}).Find(&pendingVideos).Error; err != nil {
		s.log.Errorw("Failed to load pending video tasks", "error", err)
		return
	}

	s.log.Infow("Recovering pending video generation tasks", "count", len(pendingVideos))

	for _, videoGen := range pendingVideos {
		if videoGen.TaskID == nil || *videoGen.TaskID == "" {
			s.log.Warnw("Skipping video generation with nil or empty TaskID", "id", videoGen.ID)
			continue
		}

		// Calculate how long this task has been stuck
		elapsed := time.Since(videoGen.UpdatedAt)
		
		if elapsed.Minutes() > 10 {
			// If it's been more than 10 minutes, it's a dead task. Mark as failed.
			s.log.Warnw("Task stuck for over 10 minutes. Force-failing it.", "id", videoGen.ID, "status", videoGen.Status, "elapsed_mins", elapsed.Minutes())
			
			errorMsg := "task timeout after backend restart (stuck for " + strconv.Itoa(int(elapsed.Minutes())) + " mins)"
			s.updateVideoGenError(videoGen.ID, errorMsg)
			continue
		}

		// Start goroutine to poll task status for each pending video
		// Each goroutine will poll independently until completion or timeout
		go s.pollTaskStatus(videoGen.ID, *videoGen.TaskID, videoGen.Provider, videoGen.Model)
	}
}

func (s *VideoGenerationService) GetVideoGeneration(id uint) (*models.VideoGeneration, error) {
	var videoGen models.VideoGeneration
	if err := s.db.First(&videoGen, id).Error; err != nil {
		return nil, err
	}
	return &videoGen, nil
}

func (s *VideoGenerationService) ListVideoGenerations(dramaID *uint, storyboardID *uint, status string, generationMode string, limit int, offset int) ([]*models.VideoGeneration, int64, error) {
	var videos []*models.VideoGeneration
	var total int64

	query := s.db.Model(&models.VideoGeneration{})

	if dramaID != nil {
		query = query.Where("drama_id = ?", *dramaID)
	}
	if storyboardID != nil {
		query = query.Where("storyboard_id = ?", *storyboardID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if generationMode != "" {
		query = query.Where("generation_mode = ?", generationMode)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&videos).Error; err != nil {
		return nil, 0, err
	}

	return videos, total, nil
}

func (s *VideoGenerationService) GenerateVideoFromImage(imageGenID uint) (*models.VideoGeneration, error) {
	var imageGen models.ImageGeneration
	if err := s.db.First(&imageGen, imageGenID).Error; err != nil {
		return nil, fmt.Errorf("image generation not found")
	}

	if imageGen.Status != models.ImageStatusCompleted || imageGen.ImageURL == nil {
		return nil, fmt.Errorf("image is not ready")
	}

	// 获取关联的Storyboard以获取时长
	var duration *int
	if imageGen.StoryboardID != nil {
		var storyboard models.Storyboard
		if err := s.db.Where("id = ?", *imageGen.StoryboardID).First(&storyboard).Error; err == nil {
			duration = &storyboard.Duration
			s.log.Infow("Using storyboard duration for video generation",
				"storyboard_id", *imageGen.StoryboardID,
				"duration", storyboard.Duration)
		}
	}

	req := &GenerateVideoRequest{
		DramaID:      fmt.Sprintf("%d", imageGen.DramaID),
		StoryboardID: imageGen.StoryboardID,
		ImageGenID:   &imageGenID,
		ImageURL:     *imageGen.ImageURL,
		Prompt:       imageGen.Prompt,
		Provider:     "doubao",
		Duration:     duration,
	}

	return s.GenerateVideo(req)
}

func (s *VideoGenerationService) BatchGenerateVideosForEpisode(episodeID string) ([]*models.VideoGeneration, error) {
	var episode models.Episode
	if err := s.db.Preload("Storyboards").Where("id = ?", episodeID).First(&episode).Error; err != nil {
		return nil, fmt.Errorf("episode not found")
	}

	var results []*models.VideoGeneration
	for _, storyboard := range episode.Storyboards {
		if storyboard.ImagePrompt == nil {
			continue
		}

		var imageGen models.ImageGeneration
		if err := s.db.Where("storyboard_id = ? AND status = ?", storyboard.ID, models.ImageStatusCompleted).
			Order("created_at DESC").First(&imageGen).Error; err != nil {
			s.log.Warnw("No completed image for storyboard", "storyboard_id", storyboard.ID)
			continue
		}

		videoGen, err := s.GenerateVideoFromImage(imageGen.ID)
		if err != nil {
			s.log.Errorw("Failed to generate video", "storyboard_id", storyboard.ID, "error", err)
			continue
		}

		results = append(results, videoGen)
	}

	return results, nil
}

func (s *VideoGenerationService) DeleteVideoGeneration(id uint) error {
	return s.db.Delete(&models.VideoGeneration{}, id).Error
}

func (s *VideoGenerationService) UpscaleVideo(videoGenID uint) error {
	var videoGen models.VideoGeneration
	if err := s.db.First(&videoGen, videoGenID).Error; err != nil {
		return err
	}

	if videoGen.TaskID == nil || *videoGen.TaskID == "" {
		return fmt.Errorf("视频没有TaskId，无法Upscale")
	}

	client, err := s.getVideoClient(videoGen.Provider, videoGen.Model)
	if err != nil {
		return err
	}

	type Upscaler interface {
		UpscaleVideo(taskID string) error
	}

	upscaler, ok := client.(Upscaler)
	if !ok {
		return fmt.Errorf("当前服务商 %s 不支持Upscale操作", videoGen.Provider)
	}

	err = upscaler.UpscaleVideo(*videoGen.TaskID)
	if err != nil {
		return fmt.Errorf("调用Upscale接口失败: %v", err)
	}

	// Update database status to upscaling
	if err := s.db.Model(&videoGen).Updates(map[string]interface{}{
		"status": models.VideoStatusUpscaling,
	}).Error; err != nil {
		s.log.Warnw("更新状态失败", "error", err, "id", videoGenID)
	}

	// poll task status asynchronously
	go s.pollTaskStatus(videoGenID, *videoGen.TaskID, videoGen.Provider, videoGen.Model)

	return nil
}

func (s *VideoGenerationService) ResetVideoStatus(videoGenID uint) error {
	var videoGen models.VideoGeneration
	if err := s.db.First(&videoGen, videoGenID).Error; err != nil {
		return err
	}

	targetStatus := videoGen.Status
	if videoGen.Status == models.VideoStatusUpscaling {
		targetStatus = models.VideoStatusCompleted
	} else if videoGen.Status == models.VideoStatusProcessing {
		targetStatus = models.VideoStatusPending
	} else {
		return nil // Already in a stable state
	}

	err := s.db.Model(&videoGen).Update("status", targetStatus).Error
	if err == nil {
		s.log.Infow("Video status reset by user", "id", videoGenID, "from", videoGen.Status, "to", targetStatus)
	}
	return err
}

// convertImageToBase64 将图片转换为base64格式
// 优先使用本地存储的图片，如果没有则使用URL
func (s *VideoGenerationService) convertImageToBase64(imageURL string) (string, error) {
	// 如果已经是base64格式，直接返回
	if strings.HasPrefix(imageURL, "data:") {
		return imageURL, nil
	}

	// 尝试从本地存储读取
	if s.localStorage != nil {
		var relativePath string

		// 1. 检查是否是本地URL（包含 /static/）
		if strings.Contains(imageURL, "/static/") {
			// 提取相对路径，例如从 "http://localhost:5678/static/images/xxx.jpg" 提取 "images/xxx.jpg"
			parts := strings.Split(imageURL, "/static/")
			if len(parts) == 2 {
				relativePath = parts[1]
			}
		} else if !strings.HasPrefix(imageURL, "http://") && !strings.HasPrefix(imageURL, "https://") {
			// 2. 如果不是 HTTP/HTTPS URL，视为相对路径（如 "images/xxx.jpg"）
			relativePath = imageURL
		}

		// 如果识别出相对路径，尝试读取本地文件
		if relativePath != "" {
			absPath := s.localStorage.GetAbsolutePath(relativePath)

			// 使用工具函数转换为base64
			base64Str, err := utils.ImageToBase64(absPath)
			if err == nil {
				s.log.Infow("Converted local image to base64", "path", relativePath)
				return base64Str, nil
			}
			s.log.Warnw("Failed to convert local image to base64, will try URL", "error", err, "path", absPath)
		}
	}

	// 如果本地读取失败或不是本地路径，尝试从URL下载并转换
	base64Str, err := utils.ImageToBase64(imageURL)
	if err != nil {
		return "", fmt.Errorf("failed to convert image to base64: %w", err)
	}

	urlLen := len(imageURL)
	if urlLen > 50 {
		urlLen = 50
	}
	s.log.Infow("Converted remote image to base64", "url", imageURL[:urlLen])
	return base64Str, nil
}
