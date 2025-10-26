package handler

import (
	"encoding/json"
	"net/http"
	"ocr-tester/domain"
	"ocr-tester/repository"

	"github.com/go-chi/chi/v5"
)

// PromptHandler はプロンプト関連のHTTPハンドラ
type PromptHandler struct {
   Repo repository.PromptRepository
}

// NewPromptHandler は PromptHandler の新しいインスタンスを生成
func NewPromptHandler(repo repository.PromptRepository) *PromptHandler {
   return &PromptHandler{Repo: repo}
}

// GetPrompts は全プロンプトの一覧を取得
// @Summary Get all prompts
// @Description Retrieve all prompts (excluding deleted)
// @Tags Prompts
// @Produce json
// @Success 200 {array} domain.PromptInfo
// @Failure 500 {object} map[string]string
// @Router /api/prompts [get]
func (h *PromptHandler) GetPrompts(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    prompts, err := h.Repo.FindAll(ctx)
    if err != nil {
        http.Error(w, `{"error":"failed to get prompts"}`, http.StatusInternalServerError)
        return
    }

    // 空の配列で初期化
    infos := make([]domain.PromptInfo, 0, len(prompts))
    for _, p := range prompts {
        infos = append(infos, p.ToInfo())
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(infos)
}

// GetPromptByID は指定IDのプロンプトを取得
// @Summary Get prompt by ID
// @Description Retrieve a specific prompt by ID
// @Tags Prompts
// @Produce json
// @Param id path string true "Prompt ID"
// @Success 200 {object} domain.Prompt
// @Failure 404 {object} map[string]string
// @Router /api/prompts/{id} [get]
func (h *PromptHandler) GetPromptByID(w http.ResponseWriter, r *http.Request) {
   ctx := r.Context()
   id := chi.URLParam(r, "id")

   prompt, err := h.Repo.FindByID(ctx, id)
   if err != nil {
       http.Error(w, `{"error":"prompt not found"}`, http.StatusNotFound)
       return
   }

   w.Header().Set("Content-Type", "application/json")
   json.NewEncoder(w).Encode(prompt)
}

// UpsertPrompt はプロンプトを作成または更新
// @Summary Create or update prompt
// @Description Create a new prompt or update existing one
// @Tags Prompts
// @Accept json
// @Produce json
// @Param prompt body domain.Prompt true "Prompt data"
// @Success 200 {object} domain.Prompt
// @Failure 400 {object} map[string]string
// @Router /api/prompts/upsert [post]
func (h *PromptHandler) UpsertPrompt(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    var req struct {
        ID            string  `json:"id,omitempty"`
        Name          string  `json:"name"`
        PromptContent string  `json:"promptContent"`
        UserID        *string `json:"userId,omitempty"`
    }

    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
        return
    }

    // 更新時はCreatedAtを保持
    var prompt *domain.Prompt
    if req.ID != "" {
        // 既存データを取得してCreatedAtを保持
        existing, err := h.Repo.FindByID(ctx, req.ID)
        if err != nil {
            // 存在しない場合は新規作成扱い
            prompt = &domain.Prompt{
                ID:            req.ID,
                Name:          req.Name,
                PromptContent: req.PromptContent,
                UserID:        req.UserID,
            }
        } else {
            prompt = existing
            prompt.Name = req.Name
            prompt.PromptContent = req.PromptContent
            prompt.UserID = req.UserID
        }
    } else {
        // 新規作成
        prompt = &domain.Prompt{
            Name:          req.Name,
            PromptContent: req.PromptContent,
            UserID:        req.UserID,
        }
    }

    if err := h.Repo.Upsert(ctx, prompt); err != nil {
        http.Error(w, `{"error":"failed to upsert prompt"}`, http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(prompt)
}

// DeletePrompt は指定IDのプロンプトを削除
// @Summary Delete prompt
// @Description Delete a prompt by ID (soft delete)
// @Tags Prompts
// @Param id path string true "Prompt ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /api/prompts/{id} [delete]
func (h *PromptHandler) DeletePrompt(w http.ResponseWriter, r *http.Request) {
   ctx := r.Context()
   id := chi.URLParam(r, "id")

   if err := h.Repo.Delete(ctx, id); err != nil {
       http.Error(w, `{"error":"failed to delete prompt"}`, http.StatusNotFound)
       return
   }

   w.WriteHeader(http.StatusNoContent)
}
