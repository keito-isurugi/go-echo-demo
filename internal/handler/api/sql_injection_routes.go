// internal/handler/api/sql_injection_routes.go
package api

import (
	"github.com/labstack/echo/v4"
	"go-echo-demo/internal/domain"
)

// RegisterSqlInjectionAPIRoutes SQLインジェクションデモAPIのルート登録
func RegisterSqlInjectionAPIRoutes(e *echo.Echo, productUsecase domain.ProductUsecase) {
	apiHandler := NewSqlInjectionAPIHandler(productUsecase)

	// APIルート
	apiGroup := e.Group("/api/sql-injection-demo")
	
	// 各検索方法別のエンドポイント
	apiGroup.POST("/search/vulnerable", apiHandler.SearchProductsVulnerable)
	apiGroup.POST("/search/secure", apiHandler.SearchProductsSecure)
	apiGroup.POST("/search/escaped", apiHandler.SearchProductsEscaped)
	
	// 全件取得
	apiGroup.GET("/products", apiHandler.GetAllProducts)
}