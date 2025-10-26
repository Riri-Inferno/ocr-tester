package domain

import "time"

// Prompt は、OCR処理で使用するプロンプト（指示文）を管理するドメインモデル
// Firestore 上では 1 件のドキュメントとして保存され、
// 用途別のプロンプトを名前付きで管理することができる
type Prompt struct {
    ID            string     `firestore:"id,omitempty"`      // Firestore のドキュメント ID
    Name          string     `firestore:"name"`              // プロンプトの名前（例：請求書OCR、領収書OCR）
    PromptContent string     `firestore:"promptContent"`     // プロンプトの内容（OCRへの指示文）
    IsDeleted     bool       `firestore:"isDeleted"`         // 削除フラグ
    CreatedAt     time.Time  `firestore:"createdAt"`         // レコード作成日時
    UpdatedAt     time.Time  `firestore:"updatedAt"`         // レコード更新日時
    UserID        *string    `firestore:"userId,omitempty"`  // ユーザーID（将来の認証機能用、オプショナル）
}

// PromptInfo は、プロンプト一覧表示用の軽量なモデル
// プロンプト本文を含まない、一覧表示に必要な情報のみを保持
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
