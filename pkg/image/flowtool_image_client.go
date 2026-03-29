package image

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

// FlowToolImageClient implements ImageClient for the Flow-Tool localhost API.
// Flow-Tool supports T2I (text-to-image) and I2I (image-to-image) generation modes.
type FlowToolImageClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// FlowToolUploadRequest represents the upload request payload.
type FlowToolUploadRequest struct {
	ImageData string `json:"image_data"`
	MimeType  string `json:"mime_type"`
}

// FlowToolUploadResponse represents the upload response.
type FlowToolUploadResponse struct {
	Success bool    `json:"success"`
	MediaID string  `json:"media_id"`
	Error   *string `json:"error"`
}

// FlowToolJobRequest represents the job creation request payload.
type FlowToolJobRequest struct {
	Prompt            string   `json:"prompt"`
	Mode              string   `json:"mode"`
	Quality           string   `json:"quality,omitempty"`
	Ratio             string   `json:"ratio,omitempty"`
	ReferenceImageIDs []string `json:"reference_image_ids,omitempty"`
	StartImageID      *string  `json:"start_image_id,omitempty"`
	EndImageID        *string  `json:"end_image_id,omitempty"`
	WaitForResult     bool     `json:"wait_for_result"`
}

// FlowToolJobResponse represents the job creation / status response.
type FlowToolJobResponse struct {
	JobID      string   `json:"job_id"`
	Status     string   `json:"status"`
	Progress   int      `json:"progress"`
	Message    string   `json:"message"`
	ResultURLs []string `json:"result_urls"`
	Error      *string  `json:"error"`
}

// NewFlowToolImageClient creates a new FlowToolImageClient.
// baseURL is typically "http://localhost:8000".
func NewFlowToolImageClient(baseURL string) *FlowToolImageClient {
	return &FlowToolImageClient{
		BaseURL: strings.TrimRight(baseURL, "/"),
		HTTPClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

// uploadImage uploads a base64 image to Flow-Tool and returns the media_id.
// Supports: data URI (data:image/...), raw base64 string, and HTTP/HTTPS URLs.
func (c *FlowToolImageClient) uploadImage(base64Data string) (string, error) {
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

	uploadReq := FlowToolUploadRequest{
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

	var uploadResp FlowToolUploadResponse
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

// GenerateImage implements ImageClient.GenerateImage for Flow-Tool.
// If reference images are provided, uses I2I mode. Otherwise, uses T2I mode.
func (c *FlowToolImageClient) GenerateImage(prompt string, opts ...ImageOption) (*ImageResult, error) {
	options := &ImageOptions{}
	for _, opt := range opts {
		opt(options)
	}

	// Determine mode and prepare request
	mode := "T2I"
	var referenceImageIDs []string

	// If reference images are provided, upload them and use I2I mode
	if len(options.ReferenceImages) > 0 {
		mode = "I2I"
		for _, refImg := range options.ReferenceImages {
			mediaID, err := c.uploadImage(refImg)
			if err != nil {
				return nil, fmt.Errorf("failed to upload reference image: %w", err)
			}
			referenceImageIDs = append(referenceImageIDs, mediaID)
		}
	}

	// Map ratio from project format to Flow-Tool format
	ratio := "landscape" // default
	if options.Size != "" {
		// Map common size strings to ratio
		switch {
		case strings.Contains(options.Size, "portrait") || strings.Contains(options.Size, "9:16") || strings.Contains(options.Size, "768x1344"):
			ratio = "portrait"
		default:
			ratio = "landscape"
		}
	}
	// Check dimensions for ratio
	if options.Width > 0 && options.Height > 0 {
		if options.Height > options.Width {
			ratio = "portrait"
		}
	}

	// Map quality
	quality := "fast"
	if options.Quality == "hd" || options.Quality == "quality" || options.Quality == "high" {
		quality = "quality"
	}

	jobReq := FlowToolJobRequest{
		Prompt:            prompt,
		Mode:              mode,
		Quality:           quality,
		Ratio:             ratio,
		ReferenceImageIDs: referenceImageIDs,
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

	var jobResp FlowToolJobResponse
	if err := json.Unmarshal(body, &jobResp); err != nil {
		return nil, fmt.Errorf("parse job response: %w, body: %s", err, string(body))
	}

	if jobResp.Error != nil && *jobResp.Error != "" {
		return nil, fmt.Errorf("flowtool job error: %s", *jobResp.Error)
	}

	return &ImageResult{
		TaskID:    jobResp.JobID,
		Status:    jobResp.Status,
		Completed: false,
	}, nil
}

// GetTaskStatus implements ImageClient.GetTaskStatus for Flow-Tool.
func (c *FlowToolImageClient) GetTaskStatus(taskID string) (*ImageResult, error) {
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

	var jobResp FlowToolJobResponse
	if err := json.Unmarshal(body, &jobResp); err != nil {
		return nil, fmt.Errorf("parse status response: %w, body: %s", err, string(body))
	}

	result := &ImageResult{
		TaskID: jobResp.JobID,
		Status: jobResp.Status,
	}

	// Map Flow-Tool status to ImageResult
	switch strings.ToLower(jobResp.Status) {
	case "success", "completed", "done":
		if len(jobResp.ResultURLs) > 0 && jobResp.ResultURLs[0] != "" {
			result.Completed = true
			url := jobResp.ResultURLs[0]
			if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
				if !strings.HasPrefix(url, "/") {
					url = "/" + url
				}
				url = c.BaseURL + url
			}
			result.ImageURL = url
		} else {
			result.Completed = true
			result.Error = "image generation returned success but no result URLs were found"
		}
	case "failed", "error":
		errMsg := "image generation failed"
		if jobResp.Error != nil && *jobResp.Error != "" {
			errMsg = *jobResp.Error
		} else if jobResp.Message != "" {
			errMsg = jobResp.Message
		}
		result.Error = errMsg
		result.Completed = true
	default:
		// queued, pending, processing → still in progress
		result.Completed = false
	}

	return result, nil
}
