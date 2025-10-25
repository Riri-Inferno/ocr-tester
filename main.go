package main

import (
	"fmt"
	"log"
	"net/http"

	// 自動生成された docs パッケージをインポート
	_ "ocr-tester/docs"
	// Swagger UIを表示するためのライブラリ
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// @title OCR Tester API
// @version 1.0
// @description これはGo言語で作成されたOCRテスト用のシンプルなAPIです。
// @host localhost:8080
// @BasePath /
func main() {
	http.HandleFunc("/api/hello", helloHandler)

    // http://localhost:8080/swagger/index.html でアクセス可能
    http.Handle("/swagger/", httpSwagger.WrapHandler)
	
	fmt.Println("サーバーをポート :8080 で起動しました...")
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatalf("サーバー起動エラー: %v", err)
	}
}

// @Summary Hello Worldメッセージを返す
// @Description シンプルなGETリクエストで起動確認用のメッセージを返す
// @Tags General
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "成功時のメッセージ"
// @Router /api/hello [get]
func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := `{"message": "Hello, World! from Go API"}`
	fmt.Fprint(w, response)
}
