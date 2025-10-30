package service

import (
	"context"
	"fmt"
	"io"
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

// ProcessingMethod
type ProcessingMethod string

const (
    ProcessingMethodPDFDirect   ProcessingMethod = "pdf_direct"
    ProcessingMethodImageDirect ProcessingMethod = "image_direct"
)

// OCRRequest
type OCRRequest struct {
    File         io.Reader
    FileName     string
    FileSize     int64
    CustomPrompt string
}

// OCRResponse
type OCRResponse struct {
    OCRResult        string           `json:"ocrResult"`
    PromptVersion    int              `json:"promptVersion"`
    ProcessingMethod ProcessingMethod `json:"processingMethod"`
}

// OCR処理を実行
func (s *OCRService) ProcessOCR(ctx context.Context, req *OCRRequest) (*OCRResponse, error) {
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

    // ファイルタイプを判定
    fileType, mimeType := utils.DetectFileType(fileData)
    
    // 処理方法を決定
    var processingMethod ProcessingMethod
    switch fileType {
    case utils.FileTypePDF:
        processingMethod = ProcessingMethodPDFDirect
    case utils.FileTypeJPEG, utils.FileTypePNG, utils.FileTypeGIF, utils.FileTypeWebP:
        processingMethod = ProcessingMethodImageDirect
    default:
        return nil, fmt.Errorf("unsupported file type")
    }

    // Gemini APIでOCR処理
    ocrText, err := s.geminiClient.ProcessImage(ctx, fileData, mimeType, req.CustomPrompt)
    if err != nil {
        return nil, fmt.Errorf("OCR processing failed: %w", err)
    }

    // レスポンスを構築
    response := &OCRResponse{
        OCRResult:        ocrText,
        PromptVersion:    1,
        ProcessingMethod: processingMethod,
    }

    return response, nil
}
