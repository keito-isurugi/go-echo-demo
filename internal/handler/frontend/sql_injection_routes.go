// internal/handler/frontend/sql_injection_routes.go
package frontend

import (
	"github.com/labstack/echo/v4"
	"go-echo-demo/internal/domain"
)

// RegisterSqlInjectionRoutes SQLインジェクションデモのルート登録
func RegisterSqlInjectionRoutes(e *echo.Echo, productUsecase domain.ProductUsecase) {
	sqlInjectionHandler := NewSqlInjectionHandler(productUsecase)

	// フロントエンドルート
	e.GET("/sql-injection-demo", sqlInjectionHandler.SQLInjectionDemo)
	e.POST("/sql-injection-demo/search", sqlInjectionHandler.SearchProducts)
	e.GET("/sql-injection-demo/product/:id", sqlInjectionHandler.GetProduct)
}