package domain

import "time"

// Prompt は、OCR処理で使用するプロンプトを管理するドメインモデル
type Prompt struct {
    ID            string     `firestore:"id,omitempty"`
    Name          string     `firestore:"name"`
    PromptContent string     `firestore:"promptContent"`
    CreatedAt     time.Time  `firestore:"createdAt"`
    UpdatedAt     time.Time  `firestore:"updatedAt"`
    UserID        *string    `firestore:"userId,omitempty"`
}

// PromptInfo は、プロンプト一覧表示用の軽量なモデル
type PromptInfo struct {
    ID        string     `json:"id"`
    Name      string     `json:"name"`
    CreatedAt time.Time  `json:"createdAt"`
    UpdatedAt time.Time  `json:"updatedAt"`
    UserID    *string    `json:"userId,omitempty"`
}

// ToInfo は、Prompt から PromptInfo への変換を行う
func (p *Prompt) ToInfo() PromptInfo {
    return PromptInfo{
        ID:        p.ID,
        Name:      p.Name,
        CreatedAt: p.CreatedAt,
        UpdatedAt: p.UpdatedAt,
        UserID:    p.UserID,
    }
}
