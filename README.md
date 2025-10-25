# OCR Tester API / Frontend

## 🌟 プロジェクト概要

このプロジェクトは、Go 言語（Golang）で構築された Web API と、Vue.js/TypeScript で構築されたフロントエンドを組み合わせた、**OCR（光学文字認識）のテストアプリケーション**です。

バックエンドは、GCP の Vision AI および Firestore と連携し、画像アップロードから結果保存までを処理するモダンなクリーンアーキテクチャを採用しています。

## 🛠️ 技術スタック

| 分野              | 技術                                        | 特徴                                    |
| :---------------- | :------------------------------------------ | :-------------------------------------- |
| **Backend API**   | **Go (Golang)**                             | クリーンアーキテクチャ、高速実行。      |
| **Frontend**      | **Vue 3, TypeScript, Vite**                 | モダンなビルドシステムと型安全な UI。   |
| **Cloud Service** | **Google Cloud Platform (GCP)**             | クラウド環境。                          |
| **AI / DB**       | **Gemini 2.5 Flash** (予定) , **Firestore** | 高精度な画像認識と NoSQL データベース。 |
| **Infra**         | **Docker Compose**                          | コンテナ化                              |

## 🚀 ローカル開発環境のセットアップ

### 1\. 依存ツールのインストール

以下のツールがローカルにインストールされていることを確認してください。

- Docker / Docker Compose
- Go (プロジェクトで指定されたバージョン 1.25.3 以降)
- Node.js / npm (フロントエンド用)
- Google Cloud CLI (`gcloud`)

#### 2\. GCP 認証情報の設定 (最重要)

1.  **サービスアカウントキー**を GCP コンソールで作成し、プロジェクトルートに\*\*`service_account.json`**として配置します。（**`.gitignore`で管理してください\*\*）

2.  プロジェクトルートにある\*\*`.env.example`**を**`.env`**にリネームし、以下の値を**あなたの環境に合わせて\*\*設定してください。

    ```env
    # .env の内容 (例)
    GOOGLE_CLOUD_PROJECT="pj-name"
    FIRESTORE_DATABASE_ID="firestore-id"
    ```

### 3\. API（Go）の起動

Docker Compose を使用して、バックエンドサービスをコンテナで起動します。

```bash
# 依存関係を解決し、ビルド（初回のみ、またはコード変更時）
docker compose build

# バックエンドコンテナを起動
docker compose up
```

サーバーが起動すると、以下の URL でアクセス可能になります。

- **API エンドポイント**: `http://localhost:8080/api/hello`
- **Swagger UI**: `http://localhost:8080/swagger/index.html`

### 4\. フロントエンド（Vue/TS）の起動

準備中
