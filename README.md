# Go Echo Demo

Go Echoフレームワークを使用したクリーンアーキテクチャのサンプルアプリケーションです。

## 機能

- ユーザー管理（CRUD操作）
- Basic認証
- Digest認証
- JWT認証
- **マルチプロバイダーOAuth認証**（新機能）
  - Google OAuth認証
  - LINE OAuth認証
  - 拡張可能なアーキテクチャ（Facebook、X等も簡単に追加可能）

## 認証機能

### JWT認証
JWT（JSON Web Token）を使用したトークンベースの認証システムを実装しています。

### マルチプロバイダーOAuth認証
複数のOAuthプロバイダーに対応した認証システムを実装しています。

#### 対応プロバイダー

**Google OAuth**
- Googleアカウントを使用した認証
- メールアドレス、名前、プロフィール画像を取得

**LINE OAuth**
- LINEアカウントを使用した認証
- 表示名、プロフィール画像を取得

#### 使用方法

1. **OAuthプロバイダーの設定**

   **Google Cloud Console:**
   - Google Cloud Consoleにアクセス
   - プロジェクトを作成または選択
   - APIとサービス > 認証情報
   - クライアントIDを作成
   - 承認済みのリダイレクトURI: `http://localhost:8080/auth/google/callback`

   **LINE Developers:**
   - [LINE Developers Console](https://developers.line.biz/)にアクセス
   - プロバイダーを作成または選択
   - チャネルを作成（LINE Login）
   - チャネルIDとチャネルシークレットを取得
   - コールバックURL: `http://localhost:8080/auth/line/callback`

2. **環境変数の設定**
   ```bash
   cp env.example .env
   # .envファイルを編集してOAuth設定を追加
   ```

3. **ログインページにアクセス**
   ```
   http://localhost:8080/login
   ```

4. **認証方法を選択**
   - 通常のログイン（メール・パスワード）
   - Googleでログイン
   - LINEでログイン

#### テストユーザー（従来のログイン）
- メールアドレス: `user1@example.com`
- パスワード: `password123`

### API エンドポイント

#### 認証
- `POST /api/auth/login` - ログイン（JWTトークンを取得）
- `GET /api/auth/protected` - 保護されたリソース（認証が必要）
- `GET /auth/google` - Google OAuth認証開始
- `GET /auth/google/callback` - Google OAuth認証コールバック
- `GET /auth/line` - LINE OAuth認証開始
- `GET /auth/line/callback` - LINE OAuth認証コールバック

#### 使用例

```bash
# 通常ログイン
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user1@example.com","password":"password123"}'

# 保護されたリソースにアクセス
curl -X GET http://localhost:8080/api/auth/protected \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### セキュリティ設定

- JWTシークレットキー: 環境変数 `JWT_SECRET_KEY` で設定
- トークン有効期限: 24時間
- アルゴリズム: HS256
- OAuth設定: 環境変数で管理

## セットアップ

### 前提条件
- Go 1.23以上
- PostgreSQL

### 環境変数
```bash
# データベース設定
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=password
export DB_NAME=go_echo_demo

# JWT設定
export JWT_SECRET_KEY=your-secret-key-here

# Google OAuth設定
export GOOGLE_CLIENT_ID=your-google-client-id
export GOOGLE_CLIENT_SECRET=your-google-client-secret
export GOOGLE_REDIRECT_URL=http://localhost:8080/auth/google/callback

# LINE OAuth設定
export LINE_CLIENT_ID=your-line-client-id
export LINE_CLIENT_SECRET=your-line-client-secret
export LINE_REDIRECT_URL=http://localhost:8080/auth/line/callback
```

### データベース初期化
```bash
psql -U postgres -d go_echo_demo -f init_db.sql
```

### アプリケーション起動
```bash
go run cmd/main.go
```

## アーキテクチャ

### OAuthシステムの拡張性

新しいOAuthプロバイダーを追加する場合：

1. **ドメインモデルを追加**
   ```go
   // internal/domain/facebook_auth.go
   type FacebookAuthUsecase interface {
       GetProviderName() string
       GetAuthURL() string
       ExchangeCodeForToken(code string) (*oauth2.Token, error)
       GetUserInfo(token *oauth2.Token) (*OAuthUser, error)
       Authenticate(code string) (*AuthResponse, error)
   }
   ```

2. **ユースケースを実装**
   ```go
   // internal/usecase/facebook_auth.go
   func (u *FacebookAuthUsecase) GetProviderName() string {
       return "facebook"
   }
   ```

3. **インフラストラクチャに追加**
   ```go
   // internal/infrastructure/oauth.go
   if getEnv("FACEBOOK_CLIENT_ID", "") != "" {
       // Facebook認証を追加
   }
   ```

4. **環境変数を設定**
   ```env
   FACEBOOK_CLIENT_ID=your-facebook-client-id
   FACEBOOK_CLIENT_SECRET=your-facebook-client-secret
   ```

### システム構造

```
├── cmd/
│   └── main.go              # エントリーポイント
├── internal/
│   ├── domain/              # ドメインモデル
│   │   ├── user.go
│   │   ├── auth.go          # JWT認証ドメイン
│   │   ├── oauth_provider.go # OAuth共通インターフェース
│   │   ├── google_auth.go   # Google OAuthドメイン
│   │   └── line_auth.go     # LINE OAuthドメイン
│   ├── usecase/             # ビジネスロジック
│   │   ├── user.go
│   │   ├── auth.go          # JWT認証ユースケース
│   │   ├── google_auth.go   # Google OAuthユースケース
│   │   └── line_auth.go     # LINE OAuthユースケース
│   ├── repository/          # データアクセス
│   │   ├── user.go
│   │   ├── auth.go          # JWT認証リポジトリ
│   │   └── oauth.go         # OAuth共通リポジトリ
│   ├── handler/             # HTTPハンドラー
│   │   ├── api/
│   │   │   ├── user.go
│   │   │   ├── auth.go      # JWT認証API
│   │   │   └── oauth.go     # OAuth共通API
│   │   └── frontend/
│   │       ├── frontend.go
│   │       └── auth.go      # 認証フロントエンド
│   ├── middleware/          # ミドルウェア
│   │   ├── basic.go
│   │   ├── digest.go
│   │   └── jwt.go           # JWT認証ミドルウェア
│   └── infrastructure/      # 外部依存
│       ├── db.go
│       ├── user.go
│       ├── auth.go          # JWT認証インフラ
│       └── oauth.go         # OAuth共通インフラ
├── templates/               # HTMLテンプレート
│   ├── login.html           # ログインページ
│   ├── protected.html       # 保護されたページ
│   └── ...
└── static/                  # 静的ファイル
```
