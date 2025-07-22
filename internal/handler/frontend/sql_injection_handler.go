// internal/handler/frontend/sql_injection_handler.go
package frontend

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"go-echo-demo/internal/domain"
)

type SqlInjectionHandler struct {
	productUsecase domain.ProductUsecase
}

// NewSqlInjectionHandler SQLインジェクションデモハンドラーのコンストラクタ
func NewSqlInjectionHandler(productUsecase domain.ProductUsecase) *SqlInjectionHandler {
	return &SqlInjectionHandler{
		productUsecase: productUsecase,
	}
}

// SQLInjectionDemo SQLインジェクションデモページの表示
func (h *SqlInjectionHandler) SQLInjectionDemo(c echo.Context) error {
	// 初期表示時は全件取得
	products, err := h.productUsecase.GetAll(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "商品データの取得に失敗しました",
		})
	}

	data := map[string]interface{}{
		"title":    "SQLインジェクション デモ",
		"products": products,
		"query":    "",
		"method":   "",
		"result":   "",
	}

	return c.Render(http.StatusOK, "sql_injection_demo.html", data)
}

// SearchProducts 商品検索処理
func (h *SqlInjectionHandler) SearchProducts(c echo.Context) error {
	query := c.FormValue("query")
	method := c.FormValue("method")

	var products []*domain.Product
	var err error
	var result string

	// 検索方法に応じて処理を分岐
	switch method {
	case "vulnerable":
		products, err = h.productUsecase.SearchVulnerable(c.Request().Context(), query)
		result = "脆弱版: SQLインジェクション可能（文字列連結を使用）"
	case "secure":
		products, err = h.productUsecase.SearchSecure(c.Request().Context(), query)
		result = "安全版: プレースホルダーを使用（推奨）"
	case "escaped":
		products, err = h.productUsecase.SearchEscaped(c.Request().Context(), query)
		result = "エスケープ版: 文字列エスケープを使用"
	default:
		products, err = h.productUsecase.GetAll(c.Request().Context())
		result = "全件取得"
	}

	if err != nil {
		// エラーが発生した場合でもページを表示（SQLエラーを確認するため）
		data := map[string]interface{}{
			"title":    "SQLインジェクション デモ",
			"products": []*domain.Product{},
			"query":    query,
			"method":   method,
			"result":   result,
			"error":    err.Error(),
		}
		return c.Render(http.StatusOK, "sql_injection_demo.html", data)
	}

	data := map[string]interface{}{
		"title":    "SQLインジェクション デモ",
		"products": products,
		"query":    query,
		"method":   method,
		"result":   result,
	}

	return c.Render(http.StatusOK, "sql_injection_demo.html", data)
}

// GetProduct 単一商品取得（ID指定）
func (h *SqlInjectionHandler) GetProduct(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "無効なIDです",
		})
	}

	product, err := h.productUsecase.GetAll(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "商品の取得に失敗しました",
		})
	}

	// IDでフィルタリング（簡易実装）
	var targetProduct *domain.Product
	for _, p := range product {
		if p.ID == id {
			targetProduct = p
			break
		}
	}

	if targetProduct == nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "商品が見つかりません",
		})
	}

	return c.JSON(http.StatusOK, targetProduct)
}