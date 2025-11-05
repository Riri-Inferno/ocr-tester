package service

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/vertexai/genai"
)

type GeminiClient struct {
	projectID string
	location  string
	modelID   string
	client    *genai.Client
}

// NewGeminiClient 新しいGeminiClientインスタンスを作成
func NewGeminiClient(projectID, location string) (*GeminiClient, error) {
	ctx := context.Background()

	// モデルIDを環境変数から取得
	modelID := os.Getenv("GEMINI_MODEL_ID")
	if modelID == "" {
		modelID = "gemini-2.0-flash-001"
	}

	// Vertex AI クライアントを初期化（ADC使用）
	client, err := genai.NewClient(ctx, projectID, location)
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}

	return &GeminiClient{
		projectID: projectID,
		location:  location,
		modelID:   modelID,
		client:    client,
	}, nil
}

// Close クライアントをクローズ
func (c *GeminiClient) Close() error {
	return c.client.Close()
}

// ProcessFile PDFまたは画像ファイルを処理
func (c *GeminiClient) ProcessFile(ctx context.Context, fileData []byte, mimeType string, prompt string) (string, error) {
	// モデルを取得
	model := c.client.GenerativeModel(c.modelID)

	// 生成設定
	model.SetTemperature(0.2)
	model.SetTopP(0.95)
	model.SetMaxOutputTokens(8192)

	// ファイルパートを作成
	var filePart genai.Part

	switch mimeType {
	case "application/pdf":
		// PDFデータをそのまま使用
		filePart = genai.Blob{
			MIMEType: mimeType,
			Data:     fileData,
		}
	case "image/png", "image/jpeg", "image/gif", "image/webp":
		// 画像データを使用
		filePart = genai.ImageData(mimeType, fileData)
	default:
		return "", fmt.Errorf("unsupported MIME type: %s", mimeType)
	}

	// テキストパートを作成
	textPart := genai.Text(prompt)

	// コンテンツを生成
	resp, err := model.GenerateContent(ctx, filePart, textPart)
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	// レスポンスから結果を取得
	if len(resp.Candidates) == 0 {
		return "", fmt.Errorf("no candidates in response")
	}

	// テキストを抽出
	var result string
	for _, part := range resp.Candidates[0].Content.Parts {
		if text, ok := part.(genai.Text); ok {
			result += string(text)
		}
	}

	if result == "" {
		return "", fmt.Errorf("no text content found in response")
	}

	return result, nil
}
