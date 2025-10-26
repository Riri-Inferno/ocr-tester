package domain

import "time"

// OCRResult は、OCR（光学文字認識）処理の結果を表すドメインモデル
// Firestore 上では 1 件のドキュメントとして保存され、
// 入力ファイル名、抽出されたテキスト、作成日時などの情報を保持
type OCRResult struct {
	ID            string    `firestore:"id,omitempty"`       // Firestore のドキュメント ID
	FileName      string    `firestore:"fileName"`           // 処理対象のファイル名
	ExtractedText string    `firestore:"extractedText"`      // OCR により抽出されたテキスト
	CreatedAt     time.Time `firestore:"createdAt"`          // レコード作成日時
}
