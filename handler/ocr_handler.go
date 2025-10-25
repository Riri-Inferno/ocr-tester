package handler

import (
	"encoding/json"
	"net/http"
	"ocr-tester/repository"

	"github.com/go-chi/chi/v5"
)

// OcrHandler は OCR 結果関連の HTTP ハンドラ。
// 依存として OCRResultRepository を受け取る。
type OcrHandler struct {
	Repo repository.OCRResultRepository
}

// NewOcrHandler は OcrHandler の新しいインスタンスを生成する。
func NewOcrHandler(repo repository.OCRResultRepository) *OcrHandler {
	return &OcrHandler{Repo: repo}
}

// GetOCRResultByID は指定された ID の OCR 結果を取得し、JSON 形式で返す。
// @Summary Get OCR result by ID
// @Description Retrieve a specific OCR result stored in Firestore by document ID.
// @Tags OCR
// @Accept json
// @Produce json
// @Param id path string true "OCR document ID"
// @Success 200 {object} domain.OCRResult
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/ocr-result/{id} [get]
func (h *OcrHandler) GetOCRResultByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		http.Error(w, `{"error":"missing id parameter"}`, http.StatusBadRequest)
		return
	}

	result, err := h.Repo.FindByID(ctx, id)
	if err != nil {
		// Repository 層が返すエラーの粒度によっては分岐可能だが、ここでは 404 として扱う
		http.Error(w, `{"error":"result not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, `{"error":"failed to encode response"}`, http.StatusInternalServerError)
	}
}
