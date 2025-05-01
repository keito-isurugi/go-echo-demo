FROM golang:1.24-alpine

# 必要なパッケージをインストール
RUN apk add --no-cache git

# 作業ディレクトリを設定
WORKDIR /app

# Airをインストール
RUN go install github.com/air-verse/air@latest

# go.modとgo.sumをコピー
COPY go.mod go.sum ./

# 依存関係をダウンロード
RUN go mod download

# ソースコードをコピー
COPY . .

# Airの設定ファイルをコピー
COPY .air.toml ./

# ポート8080を公開
EXPOSE 8080

# Airでアプリケーションを起動
CMD ["air"]

# 本番ビルドの場合は以下を利用
# RUN go build -o app ./cmd/main.go 
