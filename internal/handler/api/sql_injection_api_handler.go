// internal/handler/api/sql_injection_api_handler.go
package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go-echo-demo/internal/domain"
)

type SqlInjectionAPIHandler struct {
	productUsecase domain.ProductUsecase
}

// NewSqlInjectionAPIHandler SQLインジェクションデモAPI ハンドラーのコンストラクタ
func NewSqlInjectionAPIHandler(productUsecase domain.ProductUsecase) *SqlInjectionAPIHandler {
	return &SqlInjectionAPIHandler{
		productUsecase: productUsecase,
	}
}

// SearchRequest 検索リクエストの構造体
type SearchRequest struct {
	Query  string `json:"query" form:"query"`
	Method string `json:"method" form:"method"`
}

// SearchResponse 検索レスポンスの構造体
type SearchResponse struct {
	Products []*domain.Product `json:"products"`
	Query    string            `json:"query"`
	Method   string            `json:"method"`
	Result   string            `json:"result"`
	Error    string            `json:"error,omitempty"`
}

// SearchProductsVulnerable 脆弱版検索API
func (h *SqlInjectionAPIHandler) SearchProductsVulnerable(c echo.Context) error {
	var req SearchRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, SearchResponse{
			Error: "リクエストが無効です",
		})
	}

	products, err := h.productUsecase.SearchVulnerable(c.Request().Context(), req.Query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, SearchResponse{
			Query:  req.Query,
			Method: "vulnerable",
			Result: "脆弱版: SQLインジェクション可能（文字列連結を使用）",
			Error:  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, SearchResponse{
		Products: products,
		Query:    req.Query,
		Method:   "vulnerable",
		Result:   "脆弱版: SQLインジェクション可能（文字列連結を使用）",
	})
}

// SearchProductsSecure 安全版検索API
func (h *SqlInjectionAPIHandler) SearchProductsSecure(c echo.Context) error {
	var req SearchRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, SearchResponse{
			Error: "リクエストが無効です",
		})
	}

	products, err := h.productUsecase.SearchSecure(c.Request().Context(), req.Query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, SearchResponse{
			Query:  req.Query,
			Method: "secure",
			Result: "安全版: プレースホルダーを使用（推奨）",
			Error:  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, SearchResponse{
		Products: products,
		Query:    req.Query,
		Method:   "secure",
		Result:   "安全版: プレースホルダーを使用（推奨）",
	})
}

// SearchProductsEscaped エスケープ版検索API
func (h *SqlInjectionAPIHandler) SearchProductsEscaped(c echo.Context) error {
	var req SearchRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, SearchResponse{
			Error: "リクエストが無効です",
		})
	}

	products, err := h.productUsecase.SearchEscaped(c.Request().Context(), req.Query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, SearchResponse{
			Query:  req.Query,
			Method: "escaped",
			Result: "エスケープ版: 文字列エスケープを使用",
			Error:  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, SearchResponse{
		Products: products,
		Query:    req.Query,
		Method:   "escaped",
		Result:   "エスケープ版: 文字列エスケープを使用",
	})
}

// GetAllProducts 全商品取得API
func (h *SqlInjectionAPIHandler) GetAllProducts(c echo.Context) error {
	products, err := h.productUsecase.GetAll(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, SearchResponse{
			Error: "商品データの取得に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, SearchResponse{
		Products: products,
		Method:   "all",
		Result:   "全件取得",
	})
}