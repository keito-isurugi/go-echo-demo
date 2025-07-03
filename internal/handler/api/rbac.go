package api

import (
	"net/http"
	"strconv"

	"go-echo-demo/internal/domain"
	"go-echo-demo/internal/middleware"

	"github.com/labstack/echo/v4"
)

type RBACHandler struct {
	rbacUsecase domain.RBACUsecase
}

func NewRBACHandler(rbacUsecase domain.RBACUsecase) *RBACHandler {
	return &RBACHandler{rbacUsecase: rbacUsecase}
}

// RegisterRBACRoutes RBACルートを登録
func RegisterRBACRoutes(e *echo.Echo, rbacUsecase domain.RBACUsecase) {
	h := NewRBACHandler(rbacUsecase)

	// 管理者権限が必要なルートグループ（JWT認証 + adminロール）
	adminGroup := e.Group("/admin")
	// 注意: 実際の使用時は authUsecase も必要になります
	// adminGroup.Use(middleware.JWTAuthWithRole(authUsecase, rbacUsecase, "admin"))

	// ロール管理API
	adminGroup.GET("/roles", h.GetRoles)
	adminGroup.GET("/roles/:id", h.GetRole)
	adminGroup.POST("/roles", h.CreateRole)
	adminGroup.PUT("/roles/:id", h.UpdateRole)
	adminGroup.DELETE("/roles/:id", h.DeleteRole)

	// 権限管理API
	adminGroup.GET("/permissions", h.GetPermissions)
	adminGroup.GET("/permissions/:id", h.GetPermission)
	adminGroup.POST("/permissions", h.CreatePermission)
	adminGroup.PUT("/permissions/:id", h.UpdatePermission)
	adminGroup.DELETE("/permissions/:id", h.DeletePermission)

	// ユーザーロール管理API
	adminGroup.GET("/users/:user_id/roles", h.GetUserRoles)
	adminGroup.POST("/users/:user_id/roles", h.AssignRoleToUser)
	adminGroup.DELETE("/users/:user_id/roles/:role_name", h.RemoveRoleFromUser)

	// ロール権限管理API
	adminGroup.GET("/roles/:role_id/permissions", h.GetRolePermissions)
	adminGroup.POST("/roles/permissions", h.AssignPermissionToRole)
	adminGroup.DELETE("/roles/:role_name/permissions/:permission_name", h.RemovePermissionFromRole)

	// 一般ユーザー用API（JWT認証 + user/adminロール）
	userGroup := e.Group("/api")
	// 注意: 実際の使用時は authUsecase も必要になります
	// userGroup.Use(middleware.JWTAuthWithRole(authUsecase, rbacUsecase, "user", "admin"))
	userGroup.GET("/my/roles", h.GetUserRoles)
}

// ロール管理API

// GetRoles ロール一覧取得
func (h *RBACHandler) GetRoles(c echo.Context) error {
	roles, err := h.rbacUsecase.GetRoles()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "ロールの取得に失敗しました")
	}
	return c.JSON(http.StatusOK, roles)
}

// GetRole ロール詳細取得
func (h *RBACHandler) GetRole(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "無効なロールIDです")
	}

	role, err := h.rbacUsecase.GetRoleByID(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "ロールの取得に失敗しました")
	}
	if role == nil {
		return echo.NewHTTPError(http.StatusNotFound, "ロールが見つかりません")
	}

	return c.JSON(http.StatusOK, role)
}

// CreateRole ロール作成
func (h *RBACHandler) CreateRole(c echo.Context) error {
	var req struct {
		Name        string `json:"name" validate:"required"`
		Description string `json:"description"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "リクエストの解析に失敗しました")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "バリデーションエラー")
	}

	role, err := h.rbacUsecase.CreateRole(req.Name, req.Description)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "ロールの作成に失敗しました")
	}

	return c.JSON(http.StatusCreated, role)
}

// UpdateRole ロール更新
func (h *RBACHandler) UpdateRole(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "無効なロールIDです")
	}

	var req struct {
		Name        string `json:"name" validate:"required"`
		Description string `json:"description"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "リクエストの解析に失敗しました")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "バリデーションエラー")
	}

	role, err := h.rbacUsecase.UpdateRole(id, req.Name, req.Description)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "ロールの更新に失敗しました")
	}

	return c.JSON(http.StatusOK, role)
}

// DeleteRole ロール削除
func (h *RBACHandler) DeleteRole(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "無効なロールIDです")
	}

	err = h.rbacUsecase.DeleteRole(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "ロールの削除に失敗しました")
	}

	return c.NoContent(http.StatusNoContent)
}

// 権限管理API

// GetPermissions 権限一覧取得
func (h *RBACHandler) GetPermissions(c echo.Context) error {
	permissions, err := h.rbacUsecase.GetPermissions()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "権限の取得に失敗しました")
	}
	return c.JSON(http.StatusOK, permissions)
}

// GetPermission 権限詳細取得
func (h *RBACHandler) GetPermission(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "無効な権限IDです")
	}

	permission, err := h.rbacUsecase.GetPermissionByID(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "権限の取得に失敗しました")
	}
	if permission == nil {
		return echo.NewHTTPError(http.StatusNotFound, "権限が見つかりません")
	}

	return c.JSON(http.StatusOK, permission)
}

// CreatePermission 権限作成
func (h *RBACHandler) CreatePermission(c echo.Context) error {
	var req struct {
		Name        string `json:"name" validate:"required"`
		Description string `json:"description"`
		Resource    string `json:"resource" validate:"required"`
		Action      string `json:"action" validate:"required"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "リクエストの解析に失敗しました")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "バリデーションエラー")
	}

	permission, err := h.rbacUsecase.CreatePermission(req.Name, req.Description, req.Resource, req.Action)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "権限の作成に失敗しました")
	}

	return c.JSON(http.StatusCreated, permission)
}

// UpdatePermission 権限更新
func (h *RBACHandler) UpdatePermission(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "無効な権限IDです")
	}

	var req struct {
		Name        string `json:"name" validate:"required"`
		Description string `json:"description"`
		Resource    string `json:"resource" validate:"required"`
		Action      string `json:"action" validate:"required"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "リクエストの解析に失敗しました")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "バリデーションエラー")
	}

	permission, err := h.rbacUsecase.UpdatePermission(id, req.Name, req.Description, req.Resource, req.Action)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "権限の更新に失敗しました")
	}

	return c.JSON(http.StatusOK, permission)
}

// DeletePermission 権限削除
func (h *RBACHandler) DeletePermission(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "無効な権限IDです")
	}

	err = h.rbacUsecase.DeletePermission(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "権限の削除に失敗しました")
	}

	return c.NoContent(http.StatusNoContent)
}

// ユーザーロール管理API

// GetUserRoles ユーザーのロール一覧取得
func (h *RBACHandler) GetUserRoles(c echo.Context) error {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	roles, err := h.rbacUsecase.GetUserRoles(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "ユーザーロールの取得に失敗しました")
	}

	return c.JSON(http.StatusOK, roles)
}

// AssignRoleToUser ユーザーにロールを割り当て
func (h *RBACHandler) AssignRoleToUser(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "無効なユーザーIDです")
	}

	var req struct {
		RoleName string `json:"role_name" validate:"required"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "リクエストの解析に失敗しました")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "バリデーションエラー")
	}

	err = h.rbacUsecase.AssignRoleToUser(userID, req.RoleName)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "ロールの割り当てに失敗しました")
	}

	return c.NoContent(http.StatusNoContent)
}

// RemoveRoleFromUser ユーザーからロールを削除
func (h *RBACHandler) RemoveRoleFromUser(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "無効なユーザーIDです")
	}

	roleName := c.Param("role_name")
	if roleName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "ロール名が必要です")
	}

	err = h.rbacUsecase.RemoveRoleFromUser(userID, roleName)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "ロールの削除に失敗しました")
	}

	return c.NoContent(http.StatusNoContent)
}

// ロール権限管理API

// GetRolePermissions ロールの権限一覧取得
func (h *RBACHandler) GetRolePermissions(c echo.Context) error {
	roleID, err := strconv.Atoi(c.Param("role_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "無効なロールIDです")
	}

	permissions, err := h.rbacUsecase.GetRolePermissions(roleID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "ロール権限の取得に失敗しました")
	}

	return c.JSON(http.StatusOK, permissions)
}

// AssignPermissionToRole ロールに権限を割り当て
func (h *RBACHandler) AssignPermissionToRole(c echo.Context) error {
	var req struct {
		RoleName       string `json:"role_name" validate:"required"`
		PermissionName string `json:"permission_name" validate:"required"`
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "リクエストの解析に失敗しました")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "バリデーションエラー")
	}

	err := h.rbacUsecase.AssignPermissionToRole(req.RoleName, req.PermissionName)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "権限の割り当てに失敗しました")
	}

	return c.NoContent(http.StatusNoContent)
}

// RemovePermissionFromRole ロールから権限を削除
func (h *RBACHandler) RemovePermissionFromRole(c echo.Context) error {
	roleName := c.Param("role_name")
	permissionName := c.Param("permission_name")

	if roleName == "" || permissionName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "ロール名と権限名が必要です")
	}

	err := h.rbacUsecase.RemovePermissionFromRole(roleName, permissionName)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "権限の削除に失敗しました")
	}

	return c.NoContent(http.StatusNoContent)
}
