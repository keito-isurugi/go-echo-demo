# SQLインジェクション デモ

## 概要

このデモは、SQLインジェクション攻撃の仕組みと対策方法を学習するための教育用コンテンツです。

**⚠️ 重要: このコードは学習目的でのみ使用してください。本番環境では絶対に使用しないでください。**

## セットアップ

### 1. データベースの準備

```bash
# PostgreSQLに接続
psql -U postgres -d go_echo_demo

# テーブルとデータを作成
\i sql_injection_demo.sql
```

### 2. アプリケーションの起動

```bash
# 依存関係のインストール
go mod tidy

# アプリケーション起動
go run cmd/main.go
```

### 3. デモページにアクセス

- フロントエンド: http://localhost:8080/sql-injection-demo
- API: http://localhost:8080/api/sql-injection-demo/

## デモの内容

### 3つの実装パターン

1. **🔴 脆弱版 (Vulnerable)**
   - 文字列連結でSQL文を構築
   - SQLインジェクション攻撃が可能
   - **絶対に本番環境で使用してはいけません**

2. **🟢 安全版 (Secure)**
   - プレースホルダー（パラメータ化クエリ）を使用
   - SQLインジェクション攻撃を完全に防げる
   - **推奨される方法**

3. **🟡 エスケープ版 (Escaped)**
   - 文字列エスケープ処理を使用
   - 基本的な攻撃は防げるが完璧ではない
   - プレースホルダーの使用を推奨

## 攻撃例

### 典型的なSQLインジェクション攻撃

```sql
' OR '1'='1
```

この文字列を検索ボックスに入力すると、脆弱版では全てのデータが表示されます。

### より危険な攻撃例

```sql
' UNION SELECT version(),current_user(),current_database(),1,1,1,1,1 --
```

データベースの情報を取得しようとする攻撃です。

## ファイル構成

```
internal/
├── domain/
│   └── product.go              # 商品エンティティとインターフェース
├── repository/
│   └── product_repository.go   # 3つの実装パターン
├── usecase/
│   └── product_usecase.go      # ビジネスロジック
└── handler/
    ├── api/
    │   ├── sql_injection_api_handler.go  # API ハンドラー
    │   └── sql_injection_routes.go       # API ルート
    └── frontend/
        ├── sql_injection_handler.go      # フロントエンド ハンドラー
        └── sql_injection_routes.go       # フロントエンド ルート

templates/
└── sql_injection_demo.html    # デモページのテンプレート

sql_injection_demo.sql          # データベース作成スクリプト
```

## API エンドポイント

### フロントエンド

- `GET /sql-injection-demo` - デモページ表示
- `POST /sql-injection-demo/search` - 商品検索

### API

- `POST /api/sql-injection-demo/search/vulnerable` - 脆弱版検索
- `POST /api/sql-injection-demo/search/secure` - 安全版検索
- `POST /api/sql-injection-demo/search/escaped` - エスケープ版検索
- `GET /api/sql-injection-demo/products` - 全商品取得

## 学習ポイント

### 1. 脆弱版の問題点

```go
// 危険: 文字列連結でSQL構築
sqlQuery := fmt.Sprintf(`
    SELECT * FROM products 
    WHERE name LIKE '%%%s%%'
`, query)
```

### 2. 安全版の実装

```go
// 安全: プレースホルダーを使用
sqlQuery := `
    SELECT * FROM products 
    WHERE name LIKE $1
`
searchPattern := fmt.Sprintf("%%%s%%", query)
err := db.Select(&products, sqlQuery, searchPattern)
```

### 3. エスケープ版の限界

```go
// 限定的: エスケープ処理
escapedQuery := strings.ReplaceAll(query, "'", "''")
sqlQuery := fmt.Sprintf(`
    SELECT * FROM products 
    WHERE name LIKE '%%%s%%'
`, escapedQuery)
```

## セキュリティのベストプラクティス

1. **プレースホルダーを使用する**
   - 最も安全で確実な方法
   - SQLインジェクション攻撃を完全に防げる

2. **入力値の検証**
   - 想定される形式かチェック
   - 長さ制限の実装

3. **最小権限の原則**
   - データベースユーザーの権限を必要最小限に
   - 不要なテーブルへのアクセスを制限

4. **定期的なセキュリティ監査**
   - コードレビューの実施
   - 自動化されたセキュリティテスト

## 注意事項

- このデモは教育目的でのみ使用してください
- 脆弱版のコードを本番環境で使用することは絶対に避けてください
- 実際のWebアプリケーションでは、適切なセキュリティ対策を実装してください