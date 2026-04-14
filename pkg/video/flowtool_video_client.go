package video

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// FlowToolVideoClient implements VideoClient for the Flow-Tool localhost API.
// Flow-Tool supports T2V, I2V_S, I2V_SE, and R2V generation modes.
type FlowToolVideoClient struct {
	BaseURL    string
	HTTPClient *http.Client
}



// flowToolJobRequest represents the job creation request payload.
type flowToolJobRequest struct {
	Prompt        string   `json:"prompt"`
	Mode          string   `json:"mode"`
	Quality       string   `json:"quality,omitempty"`
	Ratio         string   `json:"ratio,omitempty"`
	Images        []string `json:"images,omitempty"`
	WaitForResult bool     `json:"wait_for_result"`
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



// GenerateVideo implements VideoClient.GenerateVideo for Flow-Tool.
func (c *FlowToolVideoClient) GenerateVideo(imageURL, prompt string, opts ...VideoOption) (*VideoResult, error) {
	options := &VideoOptions{
		Duration:    5,
		AspectRatio: "16:9",
	}

	for _, opt := range opts {
		opt(options)
	}

	var images []string
	mode := "T2V"

	if options.GenerationMode != "" {
		if options.GenerationMode == "direct_r2v" || options.GenerationMode == "shot_r2v" || options.GenerationMode == "shot_i2v" {
			mode = "R2V"
		} else if strings.Contains(options.GenerationMode, "first_last") {
			mode = "I2V_SE"
		} else if options.GenerationMode == "t2v" {
			mode = "T2V"
		} else {
			mode = strings.ToUpper(options.GenerationMode)
		}
	} else {
		// Fallback
		if options.FirstFrameURL != "" && options.LastFrameURL != "" {
			mode = "I2V_SE"
		} else if len(options.ReferenceImageURLs) > 0 {
			mode = "R2V"
		} else if imageURL != "" {
			mode = "R2V" // User explicitly requested shot image to default to R2V
		}
	}

	// Compile images array based on resolved mode
	if mode == "I2V_SE" {
		if options.FirstFrameURL != "" {
			images = append(images, options.FirstFrameURL)
		}
		if options.LastFrameURL != "" {
			images = append(images, options.LastFrameURL)
		}
	} else if mode == "R2V" {
		if len(options.ReferenceImageURLs) > 0 {
			images = append(images, options.ReferenceImageURLs...)
		} else if imageURL != "" {
			images = append(images, imageURL)
		}
	} else if mode == "I2V_S" {
		if options.FirstFrameURL != "" {
			images = append(images, options.FirstFrameURL)
		} else if imageURL != "" {
			images = append(images, imageURL)
		}
	}

	// Map aspect ratio
	ratio := "landscape"
	if options.AspectRatio == "9:16" || options.AspectRatio == "portrait" {
		ratio = "portrait"
	}

	// Hardcode quality to fast_lp
	quality := "fast_lp"

	jobReq := flowToolJobRequest{
		Prompt:        prompt,
		Mode:          mode,
		Quality:       quality,
		Ratio:         ratio,
		Images:        images,
		WaitForResult: false,
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
	// CRITICAL: TrimSpace to handle potential whitespace in status strings
	normalizedStatus := strings.ToLower(strings.TrimSpace(jobResp.Status))
	switch normalizedStatus {
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
		// CRITICAL FIX: Propagate error field even when status is not "failed"
		// Flow-Tool may return PENDING with error for 401 auto-requeue,
		// or return an unexpected status with an error message.
		// Without this, errors are silently swallowed and jobs get stuck.
		if jobResp.Error != nil && *jobResp.Error != "" {
			result.Error = *jobResp.Error
		}
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
