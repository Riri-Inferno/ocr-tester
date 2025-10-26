package repository

import (
	"context"
	"ocr-tester/domain"
)

// PromptRepository はプロンプトの永続化を抽象化するインターフェース
type PromptRepository interface {
    // Upsert はプロンプトを作成または更新する
    Upsert(ctx context.Context, prompt *domain.Prompt) error
    
    // FindByID は指定IDのプロンプトを取得する
    FindByID(ctx context.Context, id string) (*domain.Prompt, error)
    
    // FindAll は全プロンプトを取得する（削除済みを除く）
    FindAll(ctx context.Context) ([]*domain.Prompt, error)
    
    // Delete は指定IDのプロンプトを論理削除する
    Delete(ctx context.Context, id string) error
}
