package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

// GeminiClient
type GeminiClient struct {
    projectID  string
    location   string
    modelID    string
    httpClient *http.Client
}

// 新しいGeminiClientインスタンスを作成
func NewGeminiClient(projectID, location string) *GeminiClient {
    return &GeminiClient{
        projectID: projectID,
        location:  location,
        modelID:   "gemini-2.0-flash-001",
        httpClient: &http.Client{
            Timeout: 60 * time.Second,
        },
    }
}

// GeminiRequest
type GeminiRequest struct {
    Contents         []Content        `json:"contents"`
    GenerationConfig GenerationConfig `json:"generationConfig"`
}

// Content
type Content struct {
    Role  string `json:"role"`
    Parts []Part `json:"parts"`
}

// Part
type Part struct {
    Text       string      `json:"text,omitempty"`
    InlineData *InlineData `json:"inlineData,omitempty"`
}

// InlineData
type InlineData struct {
    MimeType string `json:"mimeType"`
    Data     string `json:"data"`
}

// GenerationConfig
type GenerationConfig struct {
    Temperature     float32 `json:"temperature"`
    MaxOutputTokens int     `json:"maxOutputTokens,omitempty"`
}

// GeminiResponse
type GeminiResponse struct {
    Candidates []Candidate `json:"candidates"`
    Error      *APIError   `json:"error,omitempty"`
}

// Candidate
type Candidate struct {
    Content      Content `json:"content"`
    FinishReason string  `json:"finishReason"`
}

// APIError
type APIError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Status  string `json:"status"`
}

// gcloudを使用してアクセストークンを取得
func (c *GeminiClient) getAccessToken() (string, error) {
    cmd := exec.Command("gcloud", "auth", "application-default", "print-access-token")
    output, err := cmd.Output()
    if err != nil {
        return "", fmt.Errorf("failed to get access token: %w", err)
    }
    return strings.TrimSpace(string(output)), nil
}

// 画像をOCR処理
func (c *GeminiClient) ProcessImage(ctx context.Context, imageData []byte, mimeType string, prompt string) (string, error) {
    // アクセストークンを取得
    token, err := c.getAccessToken()
    if err != nil {
        return "", fmt.Errorf("authentication failed: %w", err)
    }

    // リクエストを構築
    request := GeminiRequest{
        Contents: []Content{
            {
                Role: "user",
                Parts: []Part{
                    {
                        InlineData: &InlineData{
                            MimeType: mimeType,
                            Data:     base64.StdEncoding.EncodeToString(imageData),
                        },
                    },
                    {
                        Text: prompt,
                    },
                },
            },
        },
        GenerationConfig: GenerationConfig{
            Temperature:     0.2,
            MaxOutputTokens: 8192,
        },
    }

    // APIエンドポイントを構築
    url := fmt.Sprintf(
        "https://%s-aiplatform.googleapis.com/v1/projects/%s/locations/%s/publishers/google/models/%s:generateContent",
        c.location, c.projectID, c.location, c.modelID,
    )

    // リクエストボディをJSON化
    jsonData, err := json.Marshal(request)
    if err != nil {
        return "", fmt.Errorf("failed to marshal request: %w", err)
    }

    // HTTPリクエストを作成
    req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return "", fmt.Errorf("failed to create request: %w", err)
    }

    req.Header.Set("Authorization", "Bearer "+token)
    req.Header.Set("Content-Type", "application/json")

    // リクエストを送信
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return "", fmt.Errorf("failed to send request: %w", err)
    }
    defer resp.Body.Close()

    // レスポンスボディを読み取り
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", fmt.Errorf("failed to read response: %w", err)
    }

    // ステータスコードをチェック
    if resp.StatusCode != http.StatusOK {
        var errResp struct {
            Error APIError `json:"error"`
        }
        if err := json.Unmarshal(body, &errResp); err == nil && errResp.Error.Message != "" {
            return "", fmt.Errorf("Gemini API error (status %d): %s", resp.StatusCode, errResp.Error.Message)
        }
        return "", fmt.Errorf("Gemini API returned status %d: %s", resp.StatusCode, string(body))
    }

    // レスポンスをパース
    var geminiResp GeminiResponse
    if err := json.Unmarshal(body, &geminiResp); err != nil {
        return "", fmt.Errorf("failed to parse response: %w", err)
    }

    // エラーチェック
    if geminiResp.Error != nil {
        return "", fmt.Errorf("API error: %s", geminiResp.Error.Message)
    }

    // 結果を取得
    if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
        return "", fmt.Errorf("no content generated")
    }

    // テキストを抽出
    for _, part := range geminiResp.Candidates[0].Content.Parts {
        if part.Text != "" {
            return part.Text, nil
        }
    }

    return "", fmt.Errorf("no text content found in response")
}
