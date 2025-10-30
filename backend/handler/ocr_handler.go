package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"ocr-tester/domain"
	"ocr-tester/service"
)

// OcrHandler
type OcrHandler struct {
    ocrService *service.OCRService
}

// NewOcrHandler
func NewOcrHandler(ocrService *service.OCRService) *OcrHandler {
    return &OcrHandler{
        ocrService: ocrService,
    }
}

// TestOCR OCRテスト処理エンドポイント
// @Summary OCR処理を実行
// @Description PDFまたは画像ファイルからテキストを抽出
// @Tags OCR
// @Accept multipart/form-data
// @Produce json
// @Param PdfFile formData file true "PDFまたは画像ファイル"
// @Param CustomPrompt formData string true "カスタムプロンプト"
// @Success 200 {object} domain.OCRResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/test-ocr [post]
func (h *OcrHandler) TestOCR(w http.ResponseWriter, r *http.Request) {
    // multipart/form-dataをパース（最大20MB）
    if err := r.ParseMultipartForm(20 << 20); err != nil {
        http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
        return
    }

    // ファイルを取得
    file, fileHeader, err := r.FormFile("PdfFile")
    if err != nil {
        http.Error(w, "File is required", http.StatusBadRequest)
        return
    }
    defer file.Close()

    // カスタムプロンプトを取得
    customPrompt := r.FormValue("CustomPrompt")
    if customPrompt == "" {
        http.Error(w, "CustomPrompt is required", http.StatusBadRequest)
        return
    }

    // OCRリクエストを構築
    ocrRequest := &domain.OCRRequest{
        File:         file,
        FileName:     fileHeader.Filename,
        FileSize:     fileHeader.Size,
        CustomPrompt: customPrompt,
    }

    // OCR処理を実行
    ctx := context.Background()
    response, err := h.ocrService.ProcessOCR(ctx, ocrRequest)
    if err != nil {
        errorResponse := map[string]string{"error": err.Error()}
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(errorResponse)
        return
    }

    // 成功レスポンスを返す
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}
