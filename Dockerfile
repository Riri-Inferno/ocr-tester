# ----------------------------
# ビルドステージ
# ----------------------------
FROM golang:1.25.3-alpine AS builder

# 必要なツールをインストール
RUN apk add --no-cache git

# ワークディレクトリ設定
WORKDIR /app

# Go Modules 関連ファイルをコピーして依存関係を解決
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピー
COPY . .

# 実行ファイルをビルド（静的リンクで軽量化）
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server .

# ----------------------------
# 実行ステージ
# ----------------------------
FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /app/server .

# ポート公開
EXPOSE 8080

# アプリケーション起動
ENTRYPOINT ["/app/server"]
