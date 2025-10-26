FROM node:22-alpine

WORKDIR /app

# 依存関係のインストール
COPY package*.json ./
RUN npm ci

COPY . .

# Viteの開発サーバーを外部からアクセス可能に
EXPOSE 5173

# 開発サーバー起動
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]
