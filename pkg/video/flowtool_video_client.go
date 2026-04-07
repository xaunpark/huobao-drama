package video

import (
	"bytes"
	"encoding/base64"
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

	// Validate: extract raw base64 and check it's actually image data (not HTML/error page)
	base64Part := imageData
	if idx := strings.Index(imageData, ","); idx != -1 {
		base64Part = imageData[idx+1:]
	}
	// Quick size check: a real image should be at least a few hundred bytes
	rawSize := len(base64Part) * 3 / 4 // approximate decoded size
	if rawSize < 100 {
		return "", fmt.Errorf("base64 image data too small (%d bytes), likely not a valid image", rawSize)
	}

	// Decode first few bytes to check magic bytes (image signature)
	if len(base64Part) > 32 {
		snippet := base64Part[:32]
		// Pad if needed for base64 decoding
		for len(snippet)%4 != 0 {
			snippet += "="
		}
		decoded, err := base64Decode(snippet)
		if err == nil && len(decoded) >= 4 {
			isImage := false
			// PNG: 89 50 4E 47
			if decoded[0] == 0x89 && decoded[1] == 0x50 && decoded[2] == 0x4E && decoded[3] == 0x47 {
				isImage = true
			}
			// JPEG: FF D8 FF
			if decoded[0] == 0xFF && decoded[1] == 0xD8 && decoded[2] == 0xFF {
				isImage = true
			}
			// GIF: 47 49 46
			if decoded[0] == 0x47 && decoded[1] == 0x49 && decoded[2] == 0x46 {
				isImage = true
			}
			// WebP: 52 49 46 46
			if decoded[0] == 0x52 && decoded[1] == 0x49 && decoded[2] == 0x46 && decoded[3] == 0x46 {
				isImage = true
			}
			// BMP: 42 4D
			if decoded[0] == 0x42 && decoded[1] == 0x4D {
				isImage = true
			}
			if !isImage {
				// Log first bytes for debugging
				hexBytes := fmt.Sprintf("%x", decoded[:min(len(decoded), 16)])
				return "", fmt.Errorf("base64 data does not contain a valid image (magic bytes: %s, size: %d bytes). The source URL may have expired or returned an error page", hexBytes, rawSize)
			}
		}
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

	fmt.Printf("[FlowTool Upload] mime=%s, base64_size=%d, decoded_size≈%d bytes\n", mimeType, len(base64Part), rawSize)

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

	fmt.Printf("[FlowTool Upload] Success: media_id=%s\n", uploadResp.MediaID)
	return uploadResp.MediaID, nil
}

// base64Decode is a small helper to decode a base64 snippet for validation.
func base64Decode(s string) ([]byte, error) {
	// Try standard encoding first, then URL-safe encoding
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		// Try URL-safe encoding
		return base64.URLEncoding.DecodeString(s)
	}
	return decoded, nil
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
		if len(jobResp.ResultURLs) > 0 && jobResp.ResultURLs[0] != "" {
			result.Completed = true
			url := jobResp.ResultURLs[0]
			if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
				if !strings.HasPrefix(url, "/") {
					url = "/" + url
				}
				url = c.BaseURL + url
			}
			result.VideoURL = url
		} else {
			result.Completed = true
			result.Error = "video generation returned success but no result URLs were found"
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

// UpscaleVideo triggers the upscale process for a completed video in Flow-Tool.
// It uses the same job ID and does not create a new one.
func (c *FlowToolVideoClient) UpscaleVideo(taskID string) error {
	endpoint := fmt.Sprintf("%s/v1/jobs/%s/upscale", c.BaseURL, taskID)
	req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		return fmt.Errorf("create upscale request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("send upscale request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read upscale response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("upscale API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Success bool   `json:"success"`
		JobID   string `json:"job_id"`
		Message string `json:"message"`
		Error   string `json:"error"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("parse upscale response: %w, body: %s", err, string(body))
	}

	if !result.Success {
		return fmt.Errorf("upscale failed: %s", result.Error)
	}

	return nil
}
