package service

import (
	"context"
	"fmt"
	"io"
	"ocr-tester/domain"
	"strings"
	"time"
)

// OCRService はOCR処理を提供するサービス
type OCRService struct {
	geminiClient *GeminiClient
}

// NewOCRService 新しいOCRServiceインスタンスを作成
func NewOCRService(projectID, location string) (*OCRService, error) {
	geminiClient, err := NewGeminiClient(projectID, location)
	if err != nil {
		return nil, fmt.Errorf("failed to create gemini client: %w", err)
	}

	return &OCRService{
		geminiClient: geminiClient,
	}, nil
}

// Close サービスをクローズ
func (s *OCRService) Close() error {
	if s.geminiClient != nil {
		return s.geminiClient.Close()
	}
	return nil
}

// ProcessOCR OCRリクエストを処理してレスポンスを返す
func (s *OCRService) ProcessOCR(ctx context.Context, req *domain.OCRRequest) (*domain.OCRResponse, error) {
	// ファイルデータを読み込む
	fileData, err := io.ReadAll(req.File)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// MIMEタイプを判定
	mimeType := s.detectMimeType(req.FileName, fileData)

	// サポートされているファイルタイプか確認
	if !s.isSupportedFileType(mimeType) {
		return nil, fmt.Errorf("unsupported file type: %s", mimeType)
	}

	// Geminiで処理（PDFも画像も同じメソッドで処理可能）
	extractedText, err := s.geminiClient.ProcessFile(ctx, fileData, mimeType, req.CustomPrompt)
	if err != nil {
		return nil, fmt.Errorf("OCR processing failed: %w", err)
	}

	// レスポンスを構築
	response := &domain.OCRResponse{
		ExtractedText: extractedText,
		Status:        "success",
		FileName:      req.FileName,
		FileSize:      req.FileSize,
		ProcessedAt:   s.getCurrentTimestamp(),
	}

	return response, nil
}

// detectMimeType ファイル名からMIMEタイプを判定
func (s *OCRService) detectMimeType(fileName string, _ []byte) string {
	lowerName := strings.ToLower(fileName)

	switch {
	case strings.HasSuffix(lowerName, ".pdf"):
		return "application/pdf"
	case strings.HasSuffix(lowerName, ".png"):
		return "image/png"
	case strings.HasSuffix(lowerName, ".jpg") || strings.HasSuffix(lowerName, ".jpeg"):
		return "image/jpeg"
	case strings.HasSuffix(lowerName, ".gif"):
		return "image/gif"
	case strings.HasSuffix(lowerName, ".webp"):
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}

// isSupportedFileType サポートされているファイルタイプか確認
func (s *OCRService) isSupportedFileType(mimeType string) bool {
	supportedTypes := []string{
		"application/pdf",
		"image/png",
		"image/jpeg",
		"image/gif",
		"image/webp",
	}

	for _, supported := range supportedTypes {
		if mimeType == supported {
			return true
		}
	}
	return false
}

// getCurrentTimestamp 現在のタイムスタンプを取得
func (s *OCRService) getCurrentTimestamp() string {
	return time.Now().Format(time.RFC3339)
}
