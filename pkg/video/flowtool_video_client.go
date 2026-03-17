package video

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/drama-generator/backend/pkg/utils"
)

// FlowToolVideoClient implements VideoClient for the Flow-Tool localhost API.
// Flow-Tool supports T2V, I2V_S, I2V_SE, and R2V generation modes.
type FlowToolVideoClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// flowToolUploadRequest represents the upload request payload.
type flowToolUploadRequest struct {
	ImageData string `json:"image_data"`
	MimeType  string `json:"mime_type"`
}

// flowToolUploadResponse represents the upload response.
type flowToolUploadResponse struct {
	Success bool    `json:"success"`
	MediaID string  `json:"media_id"`
	Error   *string `json:"error"`
}

// flowToolJobRequest represents the job creation request payload.
type flowToolJobRequest struct {
	Prompt            string   `json:"prompt"`
	Mode              string   `json:"mode"`
	Quality           string   `json:"quality,omitempty"`
	Ratio             string   `json:"ratio,omitempty"`
	ReferenceImageIDs []string `json:"reference_image_ids,omitempty"`
	StartImageID      *string  `json:"start_image_id,omitempty"`
	EndImageID        *string  `json:"end_image_id,omitempty"`
	WaitForResult     bool     `json:"wait_for_result"`
}

// flowToolJobResponse represents the job creation / status response.
type flowToolJobResponse struct {
	JobID      string   `json:"job_id"`
	Status     string   `json:"status"`
	Progress   int      `json:"progress"`
	Message    string   `json:"message"`
	ResultURLs []string `json:"result_urls"`
	Error      *string  `json:"error"`
}

// NewFlowToolVideoClient creates a new FlowToolVideoClient.
// baseURL is typically "http://localhost:8000".
func NewFlowToolVideoClient(baseURL string) *FlowToolVideoClient {
	return &FlowToolVideoClient{
		BaseURL: strings.TrimRight(baseURL, "/"),
		HTTPClient: &http.Client{
			Timeout: 300 * time.Second,
		},
	}
}

// uploadImage uploads a base64 image (or data URI) to Flow-Tool and returns the media_id.
// Supports: data URI (data:image/...), raw base64 string, and HTTP/HTTPS URLs.
func (c *FlowToolVideoClient) uploadImage(base64Data string) (string, error) {
	imageData := base64Data

	// If the input is an HTTP/HTTPS URL, download and convert to base64 first
	if strings.HasPrefix(imageData, "http://") || strings.HasPrefix(imageData, "https://") {
		converted, err := utils.ImageToBase64(imageData)
		if err != nil {
			return "", fmt.Errorf("failed to download image from URL for upload: %w", err)
		}
		imageData = converted
	}

	// Ensure the data has a data URI prefix; if not, add a default one
	if !strings.HasPrefix(imageData, "data:") {
		imageData = "data:image/png;base64," + imageData
	}

	// Detect mime type from the data URI
	mimeType := "image/png"
	if strings.Contains(imageData, "data:image/jpeg") || strings.Contains(imageData, "data:image/jpg") {
		mimeType = "image/jpeg"
	} else if strings.Contains(imageData, "data:image/webp") {
		mimeType = "image/webp"
	} else if strings.Contains(imageData, "data:image/gif") {
		mimeType = "image/gif"
	}

	uploadReq := flowToolUploadRequest{
		ImageData: imageData,
		MimeType:  mimeType,
	}

	jsonData, err := json.Marshal(uploadReq)
	if err != nil {
		return "", fmt.Errorf("marshal upload request: %w", err)
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/v1/upload", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("create upload request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("send upload request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read upload response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("upload API error (status %d): %s", resp.StatusCode, string(body))
	}

	var uploadResp flowToolUploadResponse
	if err := json.Unmarshal(body, &uploadResp); err != nil {
		return "", fmt.Errorf("parse upload response: %w, body: %s", err, string(body))
	}

	if !uploadResp.Success {
		errMsg := "unknown error"
		if uploadResp.Error != nil {
			errMsg = *uploadResp.Error
		}
		return "", fmt.Errorf("upload failed: %s", errMsg)
	}

	return uploadResp.MediaID, nil
}

// GenerateVideo implements VideoClient.GenerateVideo for Flow-Tool.
// Determines the correct Flow-Tool mode based on the reference images and options provided:
//   - No images → T2V (text to video)
//   - Single image (imageURL) → I2V_S (image to video with start frame)
//   - FirstFrame + LastFrame → I2V_SE (image to video with start + end frames)
//   - Multiple reference images → R2V (reference images to video)
func (c *FlowToolVideoClient) GenerateVideo(imageURL, prompt string, opts ...VideoOption) (*VideoResult, error) {
	options := &VideoOptions{
		Duration:    5,
		AspectRatio: "16:9",
	}

	for _, opt := range opts {
		opt(options)
	}

	// Determine mode based on provided images
	mode := "T2V"
	var startImageID *string
	var endImageID *string
	var referenceImageIDs []string

	// Priority: first_last > multiple > single > none
	if options.FirstFrameURL != "" && options.LastFrameURL != "" {
		// I2V_SE mode: start + end frames
		mode = "I2V_SE"

		firstMediaID, err := c.uploadImage(options.FirstFrameURL)
		if err != nil {
			return nil, fmt.Errorf("failed to upload first frame: %w", err)
		}
		startImageID = &firstMediaID

		lastMediaID, err := c.uploadImage(options.LastFrameURL)
		if err != nil {
			return nil, fmt.Errorf("failed to upload last frame: %w", err)
		}
		endImageID = &lastMediaID

	} else if len(options.ReferenceImageURLs) > 0 {
		// R2V mode: multiple reference images
		mode = "R2V"
		for _, refURL := range options.ReferenceImageURLs {
			mediaID, err := c.uploadImage(refURL)
			if err != nil {
				return nil, fmt.Errorf("failed to upload reference image: %w", err)
			}
			referenceImageIDs = append(referenceImageIDs, mediaID)
		}

	} else if options.FirstFrameURL != "" {
		// I2V_S mode: only start frame
		mode = "I2V_S"
		firstMediaID, err := c.uploadImage(options.FirstFrameURL)
		if err != nil {
			return nil, fmt.Errorf("failed to upload first frame: %w", err)
		}
		startImageID = &firstMediaID

	} else if imageURL != "" {
		// I2V_S mode: single image as start frame
		mode = "I2V_S"
		mediaID, err := c.uploadImage(imageURL)
		if err != nil {
			return nil, fmt.Errorf("failed to upload start image: %w", err)
		}
		startImageID = &mediaID
	}
	// else: T2V mode (no images)

	// Map aspect ratio to Flow-Tool ratio format
	ratio := "landscape" // default
	if options.AspectRatio == "9:16" || options.AspectRatio == "portrait" {
		ratio = "portrait"
	}

	// Map quality
	quality := "fast_lp"

	jobReq := flowToolJobRequest{
		Prompt:            prompt,
		Mode:              mode,
		Quality:           quality,
		Ratio:             ratio,
		ReferenceImageIDs: referenceImageIDs,
		StartImageID:      startImageID,
		EndImageID:        endImageID,
		WaitForResult:     false,
	}

	jsonData, err := json.Marshal(jobReq)
	if err != nil {
		return nil, fmt.Errorf("marshal job request: %w", err)
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/v1/jobs", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create job request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send job request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read job response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("job API error (status %d): %s", resp.StatusCode, string(body))
	}

	var jobResp flowToolJobResponse
	if err := json.Unmarshal(body, &jobResp); err != nil {
		return nil, fmt.Errorf("parse job response: %w, body: %s", err, string(body))
	}

	if jobResp.Error != nil && *jobResp.Error != "" {
		return nil, fmt.Errorf("flowtool job error: %s", *jobResp.Error)
	}

	return &VideoResult{
		TaskID:    jobResp.JobID,
		Status:    jobResp.Status,
		Completed: false,
	}, nil
}

// GetTaskStatus implements VideoClient.GetTaskStatus for Flow-Tool.
func (c *FlowToolVideoClient) GetTaskStatus(taskID string) (*VideoResult, error) {
	endpoint := fmt.Sprintf("%s/v1/jobs/%s", c.BaseURL, taskID)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("create status request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send status request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read status response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status API error (status %d): %s", resp.StatusCode, string(body))
	}

	var jobResp flowToolJobResponse
	if err := json.Unmarshal(body, &jobResp); err != nil {
		return nil, fmt.Errorf("parse status response: %w, body: %s", err, string(body))
	}

	result := &VideoResult{
		TaskID: jobResp.JobID,
		Status: jobResp.Status,
	}

	// Map Flow-Tool status to VideoResult
	switch strings.ToLower(jobResp.Status) {
	case "success", "completed", "done":
		result.Completed = true
		if len(jobResp.ResultURLs) > 0 {
			result.VideoURL = jobResp.ResultURLs[0]
		}
	case "failed", "error":
		errMsg := "video generation failed"
		if jobResp.Error != nil && *jobResp.Error != "" {
			errMsg = *jobResp.Error
		} else if jobResp.Message != "" {
			errMsg = jobResp.Message
		}
		result.Error = errMsg
	default:
		// queued, pending, processing → still in progress
		result.Completed = false
	}

	return result, nil
}
