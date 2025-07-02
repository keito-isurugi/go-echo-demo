package usecase

import (
	"fmt"

	"go-echo-demo/internal/domain"
)

type RBACUsecaseImpl struct {
	rbacRepo domain.RBACRepository
}

func NewRBACUsecase(rbacRepo domain.RBACRepository) domain.RBACUsecase {
	return &RBACUsecaseImpl{rbacRepo: rbacRepo}
}

// ロール管理
func (u *RBACUsecaseImpl) GetRoles() ([]domain.Role, error) {
	return u.rbacRepo.GetRoles()
}

func (u *RBACUsecaseImpl) GetRoleByID(id int) (*domain.Role, error) {
	return u.rbacRepo.GetRoleByID(id)
}

func (u *RBACUsecaseImpl) CreateRole(name, description string) (*domain.Role, error) {
	role := &domain.Role{
		Name:        name,
		Description: description,
	}
	err := u.rbacRepo.CreateRole(role)
	if err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}
	return role, nil
}

func (u *RBACUsecaseImpl) UpdateRole(id int, name, description string) (*domain.Role, error) {
	role := &domain.Role{
		ID:          id,
		Name:        name,
		Description: description,
	}
	err := u.rbacRepo.UpdateRole(role)
	if err != nil {
		return nil, fmt.Errorf("failed to update role: %w", err)
	}
	return role, nil
}

func (u *RBACUsecaseImpl) DeleteRole(id int) error {
	return u.rbacRepo.DeleteRole(id)
}

// 権限管理
func (u *RBACUsecaseImpl) GetPermissions() ([]domain.Permission, error) {
	return u.rbacRepo.GetPermissions()
}

func (u *RBACUsecaseImpl) GetPermissionByID(id int) (*domain.Permission, error) {
	return u.rbacRepo.GetPermissionByID(id)
}

func (u *RBACUsecaseImpl) CreatePermission(name, description, resource, action string) (*domain.Permission, error) {
	permission := &domain.Permission{
		Name:        name,
		Description: description,
		Resource:    resource,
		Action:      action,
	}
	err := u.rbacRepo.CreatePermission(permission)
	if err != nil {
		return nil, fmt.Errorf("failed to create permission: %w", err)
	}
	return permission, nil
}

func (u *RBACUsecaseImpl) UpdatePermission(id int, name, description, resource, action string) (*domain.Permission, error) {
	permission := &domain.Permission{
		ID:          id,
		Name:        name,
		Description: description,
		Resource:    resource,
		Action:      action,
	}
	err := u.rbacRepo.UpdatePermission(permission)
	if err != nil {
		return nil, fmt.Errorf("failed to update permission: %w", err)
	}
	return permission, nil
}

func (u *RBACUsecaseImpl) DeletePermission(id int) error {
	return u.rbacRepo.DeletePermission(id)
}

// ユーザーロール管理
func (u *RBACUsecaseImpl) GetUserRoles(userID int) ([]domain.Role, error) {
	return u.rbacRepo.GetUserRoles(userID)
}

func (u *RBACUsecaseImpl) AssignRoleToUser(userID int, roleName string) error {
	role, err := u.rbacRepo.GetRoleByName(roleName)
	if err != nil {
		return fmt.Errorf("failed to get role by name: %w", err)
	}
	if role == nil {
		return fmt.Errorf("role not found: %s", roleName)
	}
	return u.rbacRepo.AssignRoleToUser(userID, role.ID)
}

func (u *RBACUsecaseImpl) RemoveRoleFromUser(userID int, roleName string) error {
	role, err := u.rbacRepo.GetRoleByName(roleName)
	if err != nil {
		return fmt.Errorf("failed to get role by name: %w", err)
	}
	if role == nil {
		return fmt.Errorf("role not found: %s", roleName)
	}
	return u.rbacRepo.RemoveRoleFromUser(userID, role.ID)
}

// ロール権限管理
func (u *RBACUsecaseImpl) GetRolePermissions(roleID int) ([]domain.Permission, error) {
	return u.rbacRepo.GetRolePermissions(roleID)
}

func (u *RBACUsecaseImpl) AssignPermissionToRole(roleName, permissionName string) error {
	role, err := u.rbacRepo.GetRoleByName(roleName)
	if err != nil {
		return fmt.Errorf("failed to get role by name: %w", err)
	}
	if role == nil {
		return fmt.Errorf("role not found: %s", roleName)
	}

	permission, err := u.rbacRepo.GetPermissionByName(permissionName)
	if err != nil {
		return fmt.Errorf("failed to get permission by name: %w", err)
	}
	if permission == nil {
		return fmt.Errorf("permission not found: %s", permissionName)
	}

	return u.rbacRepo.AssignPermissionToRole(role.ID, permission.ID)
}

func (u *RBACUsecaseImpl) RemovePermissionFromRole(roleName, permissionName string) error {
	role, err := u.rbacRepo.GetRoleByName(roleName)
	if err != nil {
		return fmt.Errorf("failed to get role by name: %w", err)
	}
	if role == nil {
		return fmt.Errorf("role not found: %s", roleName)
	}

	permission, err := u.rbacRepo.GetPermissionByName(permissionName)
	if err != nil {
		return fmt.Errorf("failed to get permission by name: %w", err)
	}
	if permission == nil {
		return fmt.Errorf("permission not found: %s", permissionName)
	}

	return u.rbacRepo.RemovePermissionFromRole(role.ID, permission.ID)
}

// 権限チェック
func (u *RBACUsecaseImpl) HasPermission(userID int, resource, action string) (bool, error) {
	return u.rbacRepo.HasPermission(userID, resource, action)
}

func (u *RBACUsecaseImpl) HasRole(userID int, roleName string) (bool, error) {
	return u.rbacRepo.HasRole(userID, roleName)
}

func (u *RBACUsecaseImpl) CheckPermission(userID int, resource, action string) error {
	hasPermission, err := u.rbacRepo.HasPermission(userID, resource, action)
	if err != nil {
		return fmt.Errorf("failed to check permission: %w", err)
	}
	if !hasPermission {
		return fmt.Errorf("permission denied: %s:%s", resource, action)
	}
	return nil
}
