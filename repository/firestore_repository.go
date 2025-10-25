package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"ocr-tester/domain"

	"cloud.google.com/go/firestore"
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
