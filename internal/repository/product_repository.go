// internal/repository/product_repository.go
package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"go-echo-demo/internal/domain"
)

type productRepository struct {
	db *sqlx.DB
}

// NewProductRepository 商品リポジトリのコンストラクタ
func NewProductRepository(db *sqlx.DB) domain.ProductRepository {
	return &productRepository{
		db: db,
	}
}

// SearchVulnerable 脆弱版: SQLインジェクション攻撃が可能
// 危険: ユーザー入力をそのままSQL文に連結している
func (r *productRepository) SearchVulnerable(ctx context.Context, query string) ([]*domain.Product, error) {
	// 危険なSQL構築 - 文字列連結でSQLインジェクション可能
	sqlQuery := fmt.Sprintf(`
		SELECT id, name, description, price, category, stock, created_at, updated_at 
		FROM products 
		WHERE name LIKE '%%%s%%' OR description LIKE '%%%s%%'
		ORDER BY id
	`, query, query)

	fmt.Printf("脆弱版SQL: %s\n", sqlQuery)

	var products []*domain.Product
	err := r.db.SelectContext(ctx, &products, sqlQuery)
	if err != nil {
		return nil, fmt.Errorf("脆弱版検索でエラー: %w", err)
	}

	return products, nil
}

// SearchSecure 安全版: プレースホルダーを使用
// 安全: SQLインジェクション攻撃を防げる
func (r *productRepository) SearchSecure(ctx context.Context, query string) ([]*domain.Product, error) {
	// 安全なSQL構築 - プレースホルダーを使用
	sqlQuery := `
		SELECT id, name, description, price, category, stock, created_at, updated_at 
		FROM products 
		WHERE name LIKE $1 OR description LIKE $1
		ORDER BY id
	`

	searchPattern := fmt.Sprintf("%%%s%%", query)
	fmt.Printf("安全版SQL: %s, パラメータ: %s\n", sqlQuery, searchPattern)

	var products []*domain.Product
	err := r.db.SelectContext(ctx, &products, sqlQuery, searchPattern)
	if err != nil {
		return nil, fmt.Errorf("安全版検索でエラー: %w", err)
	}

	return products, nil
}

// SearchEscaped エスケープ版: 文字列エスケープを使用
// 注意: エスケープは完璧ではないが、基本的な攻撃は防げる
func (r *productRepository) SearchEscaped(ctx context.Context, query string) ([]*domain.Product, error) {
	// エスケープ処理（簡易版）
	escapedQuery := r.escapeString(query)
	
	// エスケープ済み文字列を使用してSQL構築
	sqlQuery := fmt.Sprintf(`
		SELECT id, name, description, price, category, stock, created_at, updated_at 
		FROM products 
		WHERE name LIKE '%%%s%%' OR description LIKE '%%%s%%'
		ORDER BY id
	`, escapedQuery, escapedQuery)

	fmt.Printf("エスケープ版SQL: %s\n", sqlQuery)

	var products []*domain.Product
	err := r.db.SelectContext(ctx, &products, sqlQuery)
	if err != nil {
		return nil, fmt.Errorf("エスケープ版検索でエラー: %w", err)
	}

	return products, nil
}

// escapeString 文字列エスケープ処理（簡易版）
// 注意: これは完璧なエスケープではなく、デモ用途
func (r *productRepository) escapeString(s string) string {
	// SQLインジェクションで使われる危険な文字をエスケープ
	s = strings.ReplaceAll(s, "'", "''")     // シングルクォートをエスケープ
	s = strings.ReplaceAll(s, "\"", "\"\"")  // ダブルクォートをエスケープ
	s = strings.ReplaceAll(s, "\\", "\\\\")  // バックスラッシュをエスケープ
	s = strings.ReplaceAll(s, ";", "\\;")    // セミコロンをエスケープ
	s = strings.ReplaceAll(s, "--", "\\-\\-") // SQLコメントをエスケープ
	s = strings.ReplaceAll(s, "/*", "\\/\\*") // SQLコメント開始をエスケープ
	s = strings.ReplaceAll(s, "*/", "\\*\\/") // SQLコメント終了をエスケープ
	return s
}

// GetAll 全件取得
func (r *productRepository) GetAll(ctx context.Context) ([]*domain.Product, error) {
	sqlQuery := `
		SELECT id, name, description, price, category, stock, created_at, updated_at 
		FROM products 
		ORDER BY id
	`

	var products []*domain.Product
	err := r.db.SelectContext(ctx, &products, sqlQuery)
	if err != nil {
		return nil, fmt.Errorf("全件取得でエラー: %w", err)
	}

	return products, nil
}

// GetByID ID指定取得
func (r *productRepository) GetByID(ctx context.Context, id int) (*domain.Product, error) {
	sqlQuery := `
		SELECT id, name, description, price, category, stock, created_at, updated_at 
		FROM products 
		WHERE id = $1
	`

	var product domain.Product
	err := r.db.GetContext(ctx, &product, sqlQuery, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("商品が見つかりません: ID %d", id)
		}
		return nil, fmt.Errorf("ID指定取得でエラー: %w", err)
	}

	return &product, nil
}