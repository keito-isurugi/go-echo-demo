package usecase

import (
	"fmt"
	"strings"

	"go-echo-demo/internal/domain"
)

type CasbinRBACUsecaseImpl struct {
	casbinRepo domain.CasbinRBACRepository
	rbacRepo   domain.RBACRepository // 既存のDBベースRBACとの互換性のため
}

func NewCasbinRBACUsecase(casbinRepo domain.CasbinRBACRepository, rbacRepo domain.RBACRepository) domain.CasbinRBACUsecase {
	return &CasbinRBACUsecaseImpl{
		casbinRepo: casbinRepo,
		rbacRepo:   rbacRepo,
	}
}

// ポリシー管理
func (u *CasbinRBACUsecaseImpl) AddPolicy(role, resource, action string) error {
	return u.casbinRepo.AddPolicy(role, resource, action)
}

func (u *CasbinRBACUsecaseImpl) RemovePolicy(role, resource, action string) error {
	return u.casbinRepo.RemovePolicy(role, resource, action)
}

func (u *CasbinRBACUsecaseImpl) GetPolicies() ([][]string, error) {
	return u.casbinRepo.GetPolicies()
}

// ロール管理
func (u *CasbinRBACUsecaseImpl) AssignRoleToUser(user, role string) error {
	return u.casbinRepo.AddRoleForUser(user, role)
}

func (u *CasbinRBACUsecaseImpl) RemoveRoleFromUser(user, role string) error {
	return u.casbinRepo.RemoveRoleForUser(user, role)
}

func (u *CasbinRBACUsecaseImpl) GetUserRoles(user string) ([]string, error) {
	return u.casbinRepo.GetRolesForUser(user)
}

func (u *CasbinRBACUsecaseImpl) GetRoleUsers(role string) ([]string, error) {
	return u.casbinRepo.GetUsersForRole(role)
}

// 権限チェック
func (u *CasbinRBACUsecaseImpl) CheckPermission(user, resource, action string) error {
	allowed, err := u.casbinRepo.Enforce(user, resource, action)
	if err != nil {
		return fmt.Errorf("権限チェックに失敗しました: %w", err)
	}
	if !allowed {
		return fmt.Errorf("権限がありません: %s:%s", resource, action)
	}
	return nil
}

func (u *CasbinRBACUsecaseImpl) HasRole(user, role string) (bool, error) {
	return u.casbinRepo.HasRoleForUser(user, role)
}

// 管理機能（既存のDBベースRBACとの互換性のため）
func (u *CasbinRBACUsecaseImpl) GetRoles() ([]domain.Role, error) {
	return u.rbacRepo.GetRoles()
}

func (u *CasbinRBACUsecaseImpl) GetPermissions() ([]domain.Permission, error) {
	return u.rbacRepo.GetPermissions()
}

func (u *CasbinRBACUsecaseImpl) CreateRole(name, description string) (*domain.Role, error) {
	role := &domain.Role{
		Name:        name,
		Description: description,
	}
	err := u.rbacRepo.CreateRole(role)
	if err != nil {
		return nil, fmt.Errorf("ロールの作成に失敗しました: %w", err)
	}
	return role, nil
}

func (u *CasbinRBACUsecaseImpl) CreatePermission(name, description, resource, action string) (*domain.Permission, error) {
	permission := &domain.Permission{
		Name:        name,
		Description: description,
		Resource:    resource,
		Action:      action,
	}
	err := u.rbacRepo.CreatePermission(permission)
	if err != nil {
		return nil, fmt.Errorf("権限の作成に失敗しました: %w", err)
	}

	// Casbinにもポリシーを追加
	// 権限名からロールを推測（例: user:read -> userロール）
	parts := strings.Split(name, ":")
	if len(parts) == 2 {
		role := parts[0]
		action := parts[1]
		// 基本的なリソースマッピング
		resource := "content" // デフォルト
		switch action {
		case "read", "write", "delete":
			resource = "content"
		}
		u.casbinRepo.AddPolicy(role, resource, action)
	}

	return permission, nil
}
