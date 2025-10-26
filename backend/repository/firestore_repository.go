package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"ocr-tester/domain"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// firestoreRepository は OCRResultRepository の Firestore 実装。
// Firestore クライアントとコレクション名を保持し、OCR 処理結果の永続化を行う。
type firestoreRepository struct {
	client     *firestore.Client
	collection string // コレクション名（例: "ocr_results"）
}

// NewFirestoreRepository は Firestore を用いた OCRResultRepository 実装を生成する。
func NewFirestoreRepository(client *firestore.Client, collectionName string) OCRResultRepository {
	return &firestoreRepository{
		client:     client,
		collection: collectionName,
	}
}

// Save は OCRResult を Firestore に保存する。
// ID が未設定の場合は自動生成し、作成日時を現在時刻で上書きする。
func (r *firestoreRepository) Save(ctx context.Context, result *domain.OCRResult) error {
	if result.ID == "" {
		result.ID = r.client.Collection(r.collection).NewDoc().ID
	}
	result.CreatedAt = time.Now()

	_, err := r.client.Collection(r.collection).Doc(result.ID).Set(ctx, result)
	if err != nil {
		log.Printf("failed to save OCRResult to Firestore: %v", err)
		return fmt.Errorf("failed to save OCRResult to Firestore: %w", err)
	}
	return nil
}

// FindByID は指定された ID の OCRResult を Firestore から取得する。
// ドキュメントが存在しない場合はエラーを返す。
func (r *firestoreRepository) FindByID(ctx context.Context, id string) (*domain.OCRResult, error) {
	docRef := r.client.Collection(r.collection).Doc(id)
	docSnap, err := docRef.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, fmt.Errorf("document not found: %s", id)
		}
		log.Printf("failed to fetch document from Firestore (ID: %s): %v", id, err)
		return nil, fmt.Errorf("failed to fetch document from Firestore: %w", err)
	}

	var result domain.OCRResult
	if err := docSnap.DataTo(&result); err != nil {
		log.Printf("failed to decode document: %v", err)
		return nil, fmt.Errorf("failed to decode document: %w", err)
	}

	result.ID = docSnap.Ref.ID
	return &result, nil
}

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
    
    var prompts []*domain.Prompt
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
    
    return prompts, nil
}

// Delete は指定IDのプロンプトを物理削除する
func (r *promptFirestoreRepository) Delete(ctx context.Context, id string) error {
    _, err := r.client.Collection(r.collection).Doc(id).Delete(ctx)
    if err != nil {
        return fmt.Errorf("failed to delete prompt: %w", err)
    }
    return nil
}
