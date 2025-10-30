package service

import (
	"context"
	"fmt"
	"io"
	"ocr-tester/domain"
	"ocr-tester/utils"
)

// OCRService
type OCRService struct {
    geminiClient *GeminiClient
}

// 新しいOCRServiceインスタンスを作成
func NewOCRService(projectID, location string) *OCRService {
    return &OCRService{
        geminiClient: NewGeminiClient(projectID, location),
    }
}

// OCR処理を実行
func (s *OCRService) ProcessOCR(ctx context.Context, req *domain.OCRRequest) (*domain.OCRResponse, error) {
    // プロンプトの検証
    if req.CustomPrompt == "" {
        return nil, fmt.Errorf("custom prompt is required")
    }

    // ファイルサイズの事前チェック
    if req.FileSize > utils.MaxFileSize {
        return nil, fmt.Errorf("file size exceeds maximum allowed size of %d MB", utils.MaxFileSize/(1024*1024))
    }

    // ファイルデータを読み込む
    fileData, err := io.ReadAll(req.File)
    if err != nil {
        return nil, fmt.Errorf("failed to read file: %w", err)
    }

    // ファイルのバリデーション
    if err := utils.ValidateFile(fileData, req.FileName); err != nil {
        return nil, err
    }

    // ファイルタイプを判定してMIMEタイプ取得
    _, mimeType := utils.DetectFileType(fileData)

    // Gemini APIでOCR処理
    ocrText, err := s.geminiClient.ProcessImage(ctx, fileData, mimeType, req.CustomPrompt)
    if err != nil {
        return nil, fmt.Errorf("OCR processing failed: %w", err)
    }

    // レスポンスを構築
    response := &domain.OCRResponse{
        OCRResult:     ocrText,
    }

    return response, nil
}
