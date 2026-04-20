package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

// Multimodal types — separate from ChatMessage to avoid breaking existing text-only code.
// Used specifically for vision requests (e.g., video review via Locally-AI / Gemini).

// ContentPart represents a single part of multimodal content (text or image).
type ContentPart struct {
	Type     string    `json:"type"`                // "text" or "image_url"
	Text     string    `json:"text,omitempty"`       // present when Type="text"
	ImageURL *ImageURLRef `json:"image_url,omitempty"` // present when Type="image_url"
}

// ImageURLRef holds the URL reference for an image content part.
type ImageURLRef struct {
	URL string `json:"url"` // "file:///path/to/image.jpg" or "data:image/...;base64,..."
}

// MultimodalMessage is a chat message with array content (text + images).
type MultimodalMessage struct {
	Role    string        `json:"role"`
	Content []ContentPart `json:"content"`
}

// MultimodalRequest is an OpenAI-compatible chat completion request with multimodal content.
type MultimodalRequest struct {
	Model    string              `json:"model"`
	Messages []MultimodalMessage `json:"messages"`
}

// SendMultimodal sends a multimodal (text + image) request to an OpenAI-compatible endpoint.
// This bypasses the existing ChatCompletion() method to avoid changing ChatMessage.Content type.
// Returns the assistant's response text content.
func SendMultimodal(baseURL, apiKey string, req *MultimodalRequest) (string, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal multimodal request: %w", err)
	}

	url := strings.TrimRight(baseURL, "/") + "/chat/completions"

	// Log request (truncated for readability)
	preview := string(jsonData)
	if len(preview) > 500 {
		preview = preview[:500] + "..."
	}
	fmt.Printf("Multimodal: Sending request to %s\n", url)
	fmt.Printf("Multimodal: Request preview: %s\n", preview)

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{
		Timeout: 30 * time.Minute, // Locally-AI has 30min timeout
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("multimodal request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		bodyPreview := string(body)
		if len(bodyPreview) > 500 {
			bodyPreview = bodyPreview[:500]
		}
		return "", fmt.Errorf("multimodal API error (status %d): %s", resp.StatusCode, bodyPreview)
	}

	// Parse standard OpenAI response
	var chatResp ChatCompletionResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", fmt.Errorf("failed to parse multimodal response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in multimodal response")
	}

	content := chatResp.Choices[0].Message.Content
	fmt.Printf("Multimodal: Response length: %d chars\n", len(content))

	return content, nil
}

// BuildVisionRequest creates a MultimodalRequest with text prompt and multiple local image files.
// imagePaths should be absolute paths on the server; they will be converted to file:// URLs.
func BuildVisionRequest(model, prompt string, imagePaths []string) *MultimodalRequest {
	content := []ContentPart{
		{Type: "text", Text: prompt},
	}

	for _, imagePath := range imagePaths {
		if imagePath != "" {
			absPath, err := filepath.Abs(imagePath)
			if err != nil {
				absPath = imagePath // fallback
			}
			// Convert Windows path to file:// URL
			fileURL := "file:///" + strings.ReplaceAll(absPath, "\\", "/")
			content = append(content, ContentPart{
				Type:     "image_url",
				ImageURL: &ImageURLRef{URL: fileURL},
			})
		}
	}

	return &MultimodalRequest{
		Model: model,
		Messages: []MultimodalMessage{
			{
				Role:    "user",
				Content: content,
			},
		},
	}
}
