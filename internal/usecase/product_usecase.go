// internal/usecase/product_usecase.go
package usecase

import (
	"context"
	"fmt"

	"go-echo-demo/internal/domain"
)

type productUsecase struct {
	productRepo domain.ProductRepository
}

// NewProductUsecase 商品ユースケースのコンストラクタ
func NewProductUsecase(productRepo domain.ProductRepository) domain.ProductUsecase {
	return &productUsecase{
		productRepo: productRepo,
	}
}

// SearchVulnerable 脆弱版検索
// SQLインジェクション攻撃が可能な検索機能
func (u *productUsecase) SearchVulnerable(ctx context.Context, query string) ([]*domain.Product, error) {
	if query == "" {
		// 空文字の場合は全件取得
		return u.productRepo.GetAll(ctx)
	}

	products, err := u.productRepo.SearchVulnerable(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("脆弱版検索でエラーが発生しました: %w", err)
	}

	return products, nil
}

// SearchSecure 安全版検索
// プレースホルダーを使用した安全な検索機能
func (u *productUsecase) SearchSecure(ctx context.Context, query string) ([]*domain.Product, error) {
	if query == "" {
		// 空文字の場合は全件取得
		return u.productRepo.GetAll(ctx)
	}

	products, err := u.productRepo.SearchSecure(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("安全版検索でエラーが発生しました: %w", err)
	}

	return products, nil
}

// SearchEscaped エスケープ版検索
// 文字列エスケープを使用した検索機能
func (u *productUsecase) SearchEscaped(ctx context.Context, query string) ([]*domain.Product, error) {
	if query == "" {
		// 空文字の場合は全件取得
		return u.productRepo.GetAll(ctx)
	}

	products, err := u.productRepo.SearchEscaped(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("エスケープ版検索でエラーが発生しました: %w", err)
	}

	return products, nil
}

// GetAll 全件取得
func (u *productUsecase) GetAll(ctx context.Context) ([]*domain.Product, error) {
	products, err := u.productRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("全件取得でエラーが発生しました: %w", err)
	}

	return products, nil
}