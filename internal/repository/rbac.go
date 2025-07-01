package repository

import (
	"database/sql"
	"fmt"
	"time"

	"go-echo-demo/internal/domain"
)

type RBACRepositoryImpl struct {
	db *sql.DB
}

func NewRBACRepository(db *sql.DB) domain.RBACRepository {
	return &RBACRepositoryImpl{db: db}
}

// ロール関連
func (r *RBACRepositoryImpl) GetRoles() ([]domain.Role, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM roles ORDER BY id`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles: %w", err)
	}
	defer rows.Close()

	var roles []domain.Role
	for rows.Next() {
		var role domain.Role
		err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (r *RBACRepositoryImpl) GetRoleByID(id int) (*domain.Role, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM roles WHERE id = $1`
	var role domain.Role
	err := r.db.QueryRow(query, id).Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get role by id: %w", err)
	}
	return &role, nil
}

func (r *RBACRepositoryImpl) GetRoleByName(name string) (*domain.Role, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM roles WHERE name = $1`
	var role domain.Role
	err := r.db.QueryRow(query, name).Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get role by name: %w", err)
	}
	return &role, nil
}

func (r *RBACRepositoryImpl) CreateRole(role *domain.Role) error {
	query := `INSERT INTO roles (name, description, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id`
	now := time.Now()
	err := r.db.QueryRow(query, role.Name, role.Description, now, now).Scan(&role.ID)
	if err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}
	role.CreatedAt = now
	role.UpdatedAt = now
	return nil
}

func (r *RBACRepositoryImpl) UpdateRole(role *domain.Role) error {
	query := `UPDATE roles SET name = $1, description = $2, updated_at = $3 WHERE id = $4`
	role.UpdatedAt = time.Now()
	_, err := r.db.Exec(query, role.Name, role.Description, role.UpdatedAt, role.ID)
	if err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}
	return nil
}

func (r *RBACRepositoryImpl) DeleteRole(id int) error {
	query := `DELETE FROM roles WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}
	return nil
}

// 権限関連
func (r *RBACRepositoryImpl) GetPermissions() ([]domain.Permission, error) {
	query := `SELECT id, name, description, resource, action, created_at, updated_at FROM permissions ORDER BY id`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions: %w", err)
	}
	defer rows.Close()

	var permissions []domain.Permission
	for rows.Next() {
		var permission domain.Permission
		err := rows.Scan(&permission.ID, &permission.Name, &permission.Description, &permission.Resource, &permission.Action, &permission.CreatedAt, &permission.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan permission: %w", err)
		}
		permissions = append(permissions, permission)
	}
	return permissions, nil
}

func (r *RBACRepositoryImpl) GetPermissionByID(id int) (*domain.Permission, error) {
	query := `SELECT id, name, description, resource, action, created_at, updated_at FROM permissions WHERE id = $1`
	var permission domain.Permission
	err := r.db.QueryRow(query, id).Scan(&permission.ID, &permission.Name, &permission.Description, &permission.Resource, &permission.Action, &permission.CreatedAt, &permission.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get permission by id: %w", err)
	}
	return &permission, nil
}

func (r *RBACRepositoryImpl) GetPermissionByName(name string) (*domain.Permission, error) {
	query := `SELECT id, name, description, resource, action, created_at, updated_at FROM permissions WHERE name = $1`
	var permission domain.Permission
	err := r.db.QueryRow(query, name).Scan(&permission.ID, &permission.Name, &permission.Description, &permission.Resource, &permission.Action, &permission.CreatedAt, &permission.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get permission by name: %w", err)
	}
	return &permission, nil
}

func (r *RBACRepositoryImpl) CreatePermission(permission *domain.Permission) error {
	query := `INSERT INTO permissions (name, description, resource, action, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	now := time.Now()
	err := r.db.QueryRow(query, permission.Name, permission.Description, permission.Resource, permission.Action, now, now).Scan(&permission.ID)
	if err != nil {
		return fmt.Errorf("failed to create permission: %w", err)
	}
	permission.CreatedAt = now
	permission.UpdatedAt = now
	return nil
}

func (r *RBACRepositoryImpl) UpdatePermission(permission *domain.Permission) error {
	query := `UPDATE permissions SET name = $1, description = $2, resource = $3, action = $4, updated_at = $5 WHERE id = $6`
	permission.UpdatedAt = time.Now()
	_, err := r.db.Exec(query, permission.Name, permission.Description, permission.Resource, permission.Action, permission.UpdatedAt, permission.ID)
	if err != nil {
		return fmt.Errorf("failed to update permission: %w", err)
	}
	return nil
}

func (r *RBACRepositoryImpl) DeletePermission(id int) error {
	query := `DELETE FROM permissions WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete permission: %w", err)
	}
	return nil
}

// ユーザーロール関連
func (r *RBACRepositoryImpl) GetUserRoles(userID int) ([]domain.Role, error) {
	query := `
		SELECT r.id, r.name, r.description, r.created_at, r.updated_at 
		FROM roles r 
		JOIN user_roles ur ON r.id = ur.role_id 
		WHERE ur.user_id = $1 
		ORDER BY r.id
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}
	defer rows.Close()

	var roles []domain.Role
	for rows.Next() {
		var role domain.Role
		err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user role: %w", err)
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (r *RBACRepositoryImpl) AssignRoleToUser(userID, roleID int) error {
	query := `INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2) ON CONFLICT (user_id, role_id) DO NOTHING`
	_, err := r.db.Exec(query, userID, roleID)
	if err != nil {
		return fmt.Errorf("failed to assign role to user: %w", err)
	}
	return nil
}

func (r *RBACRepositoryImpl) RemoveRoleFromUser(userID, roleID int) error {
	query := `DELETE FROM user_roles WHERE user_id = $1 AND role_id = $2`
	_, err := r.db.Exec(query, userID, roleID)
	if err != nil {
		return fmt.Errorf("failed to remove role from user: %w", err)
	}
	return nil
}

func (r *RBACRepositoryImpl) GetUsersByRole(roleID int) ([]domain.User, error) {
	query := `
		SELECT u.id, u.name, u.email, u.password, u.provider_id, u.provider_name 
		FROM users u 
		JOIN user_roles ur ON u.id = ur.user_id 
		WHERE ur.role_id = $1 
		ORDER BY u.id
	`
	rows, err := r.db.Query(query, roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get users by role: %w", err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.ProviderID, &user.ProviderName)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}

// ロール権限関連
func (r *RBACRepositoryImpl) GetRolePermissions(roleID int) ([]domain.Permission, error) {
	query := `
		SELECT p.id, p.name, p.description, p.resource, p.action, p.created_at, p.updated_at 
		FROM permissions p 
		JOIN role_permissions rp ON p.id = rp.permission_id 
		WHERE rp.role_id = $1 
		ORDER BY p.id
	`
	rows, err := r.db.Query(query, roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get role permissions: %w", err)
	}
	defer rows.Close()

	var permissions []domain.Permission
	for rows.Next() {
		var permission domain.Permission
		err := rows.Scan(&permission.ID, &permission.Name, &permission.Description, &permission.Resource, &permission.Action, &permission.CreatedAt, &permission.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan role permission: %w", err)
		}
		permissions = append(permissions, permission)
	}
	return permissions, nil
}

func (r *RBACRepositoryImpl) AssignPermissionToRole(roleID, permissionID int) error {
	query := `INSERT INTO role_permissions (role_id, permission_id) VALUES ($1, $2) ON CONFLICT (role_id, permission_id) DO NOTHING`
	_, err := r.db.Exec(query, roleID, permissionID)
	if err != nil {
		return fmt.Errorf("failed to assign permission to role: %w", err)
	}
	return nil
}

func (r *RBACRepositoryImpl) RemovePermissionFromRole(roleID, permissionID int) error {
	query := `DELETE FROM role_permissions WHERE role_id = $1 AND permission_id = $2`
	_, err := r.db.Exec(query, roleID, permissionID)
	if err != nil {
		return fmt.Errorf("failed to remove permission from role: %w", err)
	}
	return nil
}

// 権限チェック
func (r *RBACRepositoryImpl) HasPermission(userID int, resource, action string) (bool, error) {
	query := `
		SELECT COUNT(*) > 0 
		FROM user_roles ur 
		JOIN role_permissions rp ON ur.role_id = rp.role_id 
		JOIN permissions p ON rp.permission_id = p.id 
		WHERE ur.user_id = $1 AND p.resource = $2 AND p.action = $3
	`
	var hasPermission bool
	err := r.db.QueryRow(query, userID, resource, action).Scan(&hasPermission)
	if err != nil {
		return false, fmt.Errorf("failed to check permission: %w", err)
	}
	return hasPermission, nil
}

func (r *RBACRepositoryImpl) HasRole(userID int, roleName string) (bool, error) {
	query := `
		SELECT COUNT(*) > 0 
		FROM user_roles ur 
		JOIN roles r ON ur.role_id = r.id 
		WHERE ur.user_id = $1 AND r.name = $2
	`
	var hasRole bool
	err := r.db.QueryRow(query, userID, roleName).Scan(&hasRole)
	if err != nil {
		return false, fmt.Errorf("failed to check role: %w", err)
	}
	return hasRole, nil
}
