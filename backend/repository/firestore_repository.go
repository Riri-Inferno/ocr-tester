package repository

import (
	"context"
	"fmt"
	"time"

	"ocr-tester/domain"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// promptFirestoreRepository は PromptRepository の Firestore 実装
type promptFirestoreRepository struct {
    client     *firestore.Client
    collection string
}

// NewPromptFirestoreRepository は Firestore を用いた PromptRepository 実装を生成
func NewPromptFirestoreRepository(client *firestore.Client, collectionName string) PromptRepository {
    return &promptFirestoreRepository{
        client:     client,
        collection: collectionName,
    }
}

// Upsert はプロンプトを作成または更新する
func (r *promptFirestoreRepository) Upsert(ctx context.Context, prompt *domain.Prompt) error {
    if prompt.ID == "" {
        prompt.ID = r.client.Collection(r.collection).NewDoc().ID
        prompt.CreatedAt = time.Now()
    }
    prompt.UpdatedAt = time.Now()

    _, err := r.client.Collection(r.collection).Doc(prompt.ID).Set(ctx, prompt)
    if err != nil {
        return fmt.Errorf("failed to upsert prompt: %w", err)
    }
    return nil
}

// FindByID は指定IDのプロンプトを取得する
func (r *promptFirestoreRepository) FindByID(ctx context.Context, id string) (*domain.Prompt, error) {
    doc, err := r.client.Collection(r.collection).Doc(id).Get(ctx)
    if err != nil {
        if status.Code(err) == codes.NotFound {
            return nil, fmt.Errorf("prompt not found: %s", id)
        }
        return nil, fmt.Errorf("failed to get prompt: %w", err)
    }

    var prompt domain.Prompt
    if err := doc.DataTo(&prompt); err != nil {
        return nil, fmt.Errorf("failed to decode prompt: %w", err)
    }
    prompt.ID = doc.Ref.ID
    return &prompt, nil
}

// FindAll は全プロンプトを取得する（promptContentを除外）
func (r *promptFirestoreRepository) FindAll(ctx context.Context) ([]*domain.Prompt, error) {
    iter := r.client.Collection(r.collection).
        Select("name", "createdAt", "updatedAt", "userId").
        OrderBy("createdAt", firestore.Desc).
        Documents(ctx)
    
    defer iter.Stop()
    
    // 空の配列で初期化（nilではなく）
    prompts := make([]*domain.Prompt, 0)
    
    for {
        doc, err := iter.Next()
        if err == iterator.Done {
            break
        }
        if err != nil {
            return nil, fmt.Errorf("failed to iterate prompts: %w", err)
        }

        var prompt domain.Prompt
        if err := doc.DataTo(&prompt); err != nil {
            return nil, fmt.Errorf("failed to decode prompt: %w", err)
        }
        prompt.ID = doc.Ref.ID
        prompts = append(prompts, &prompt)
    }
    
    return prompts, nil  // 0件でも空の配列を返す
}

// Delete は指定IDのプロンプトを物理削除する
func (r *promptFirestoreRepository) Delete(ctx context.Context, id string) error {
    _, err := r.client.Collection(r.collection).Doc(id).Delete(ctx)
    if err != nil {
        return fmt.Errorf("failed to delete prompt: %w", err)
    }
    return nil
}
