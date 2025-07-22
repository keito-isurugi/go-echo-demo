// internal/domain/product.go
package domain

import (
	"context"
	"time"
)

// Product 商品エンティティ
type Product struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Price       float64   `json:"price" db:"price"`
	Category    string    `json:"category" db:"category"`
	Stock       int       `json:"stock" db:"stock"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// ProductRepository 商品リポジトリのインターフェース
type ProductRepository interface {
	// 脆弱版: SQLインジェクション可能（文字列連結）
	SearchVulnerable(ctx context.Context, query string) ([]*Product, error)
	
	// 安全版: プレースホルダーを使用
	SearchSecure(ctx context.Context, query string) ([]*Product, error)
	
	// エスケープ版: 文字列エスケープを使用
	SearchEscaped(ctx context.Context, query string) ([]*Product, error)
	
	// 全件取得
	GetAll(ctx context.Context) ([]*Product, error)
	
	// ID指定取得
	GetByID(ctx context.Context, id int) (*Product, error)
}

// ProductUsecase 商品ユースケースのインターフェース
type ProductUsecase interface {
	// 3つの検索方法を提供
	SearchVulnerable(ctx context.Context, query string) ([]*Product, error)
	SearchSecure(ctx context.Context, query string) ([]*Product, error)
	SearchEscaped(ctx context.Context, query string) ([]*Product, error)
	
	// 全件取得
	GetAll(ctx context.Context) ([]*Product, error)
}