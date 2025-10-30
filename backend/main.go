package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	"cloud.google.com/go/firestore"

	"ocr-tester/handler"
	"ocr-tester/repository"
	"ocr-tester/service"

	_ "ocr-tester/docs"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// Firestore コレクション名
const (
   promptCollectionName = "prompts"
)

func main() {
   // .env ファイルを読み込む（存在しない場合は無視）
   _ = godotenv.Load()

   // 環境変数から Firestore 設定を取得
   projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
   databaseID := os.Getenv("FIRESTORE_DATABASE_ID")

   if projectID == "" || databaseID == "" {
       log.Fatal("環境変数 GOOGLE_CLOUD_PROJECT または FIRESTORE_DATABASE_ID が設定されていません")
   }

   // Firestore クライアントを初期化
   ctx := context.Background()
   fsClient, err := firestore.NewClientWithDatabase(ctx, projectID, databaseID)
   if err != nil {
       log.Fatalf("Firestore クライアントの初期化に失敗しました: %v", err)
   }
   defer fsClient.Close()
   
   // OCRサービスを初期化（Gemini用）
   location := os.Getenv("GEMINI_LOCATION")
   if location == "" {
       location = "us-central1"
   }
   ocrService := service.NewOCRService(projectID, location)
   
   // OCRハンドラを初期化（新しい形式）
   ocrHandler := handler.NewOcrHandler(ocrService)

   promptRepo := repository.NewPromptFirestoreRepository(fsClient, promptCollectionName)
   promptHandler := handler.NewPromptHandler(promptRepo)

   // ルーターを設定
   r := chi.NewRouter()

   // 動作確認用エンドポイント
   r.Get("/api/hello", helloHandler)

   // OCRテストエンドポイント
   r.Post("/api/test-ocr", ocrHandler.TestOCR)

   // プロンプトエンドポイント
   r.Get("/api/prompts", promptHandler.GetPrompts)
   r.Get("/api/prompts/{id}", promptHandler.GetPromptByID)
   r.Post("/api/prompts/upsert", promptHandler.UpsertPrompt)
   r.Delete("/api/prompts/{id}", promptHandler.DeletePrompt)

   // Swagger UI のエンドポイント
   r.Handle("/swagger/*", httpSwagger.WrapHandler)

   // サーバー起動
   addr := ":8080"
   log.Printf("サーバーをポート %s で起動中...", addr)

   if err := http.ListenAndServe(addr, r); err != nil {
       log.Fatalf("サーバー起動エラー: %v", err)
   }
}

// helloHandler は単純な動作確認用エンドポイント。
//
// @Summary Hello World メッセージを返す
// @Description サーバーの起動確認用エンドポイント
// @Tags General
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/hello [get]
func helloHandler(w http.ResponseWriter, _ *http.Request) {
   w.Header().Set("Content-Type", "application/json")
   fmt.Fprint(w, `{"message":"Hello, World! from Go API"}`)
}
