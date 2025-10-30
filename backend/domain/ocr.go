package domain

import "io"

// OCRRequest は、OCR処理のリクエストを表すドメインモデル
type OCRRequest struct {
    File         io.Reader `json:"-"`
    FileName     string    `json:"fileName"`
    FileSize     int64     `json:"fileSize"`
    CustomPrompt string    `json:"customPrompt"`
}

// OCRResponse は、OCR処理の結果を表すドメインモデル
type OCRResponse struct {
    OCRResult        string           `json:"ocrResult"`
}
