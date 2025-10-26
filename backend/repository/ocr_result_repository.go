package repository

import (
	"context"
	"ocr-tester/domain"
)

// OCRResultRepository は OCR 結果の永続化を抽象化するインターフェース。
// 具体的なデータストア（Firestore など）に依存しないリポジトリ層の契約を定義する。
type OCRResultRepository interface {
	// Save は OCR 処理結果を永続化する。
	Save(ctx context.Context, result *domain.OCRResult) error

	// FindByID は ID を指定して OCR 結果を取得する。
	FindByID(ctx context.Context, id string) (*domain.OCRResult, error)
}
