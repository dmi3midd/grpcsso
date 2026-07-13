package repository

import (
	"context"
	"fmt"

	"github.com/dmi3midd/grpcsso/internal/domain"
	"github.com/jmoiron/sqlx"
)

type RBACRepository interface {
	// AssignRoleToUser inserts a (user_id, role_id) record into user_roles.
	AssignRoleToUser(ctx context.Context, ext sqlx.ExtContext, userID, roleID string) error
	// RemoveRoleFromUser deletes a record from user_roles.
	RemoveRoleFromUser(ctx context.Context, ext sqlx.ExtContext, userID, roleID string) error
	// GetUserRoles returns all roles assigned to a user.
	GetUserRoles(ctx context.Context, ext sqlx.ExtContext, userID string) ([]domain.Role, error)

	// AssignPermissionToRole inserts a (role_id, permission_id) record into role_permissions.
	AssignPermissionToRole(ctx context.Context, ext sqlx.ExtContext, roleID, permissionID string) error
	// RemovePermissionFromRole deletes a record from role_permissions.
	RemovePermissionFromRole(ctx context.Context, ext sqlx.ExtContext, roleID, permissionID string) error
	// GetRolePermissions returns all permissions assigned to a role.
	GetRolePermissions(ctx context.Context, ext sqlx.ExtContext, roleID string) ([]domain.Permission, error)

	// UserExists checks if a user with the given ID exists.
	UserExists(ctx context.Context, ext sqlx.ExtContext, userID string) (bool, error)
	// RoleExists checks if a role with the given ID exists.
	RoleExists(ctx context.Context, ext sqlx.ExtContext, roleID string) (bool, error)
	// PermissionExists checks if a permission with the given ID exists.
	PermissionExists(ctx context.Context, ext sqlx.ExtContext, permissionID string) (bool, error)
}

type rbacRepository struct{}

func NewRBACRepository() RBACRepository {
	return &rbacRepository{}
}

// User <=> Role

func (r *rbacRepository) AssignRoleToUser(ctx context.Context, ext sqlx.ExtContext, userID, roleID string) error {
	op := "RBACRepository.AssignRoleToUser"
	query := `
	INSERT INTO user_roles (user_id, role_id)
	VALUES ($1, $2)
	ON CONFLICT DO NOTHING
	`
	_, err := ext.ExecContext(ctx, query, userID, roleID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *rbacRepository) RemoveRoleFromUser(ctx context.Context, ext sqlx.ExtContext, userID, roleID string) error {
	op := "RBACRepository.RemoveRoleFromUser"
	query := `
	DELETE FROM user_roles
	WHERE user_id = $1 AND role_id = $2
	`
	_, err := ext.ExecContext(ctx, query, userID, roleID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *rbacRepository) GetUserRoles(ctx context.Context, ext sqlx.ExtContext, userID string) ([]domain.Role, error) {
	op := "RBACRepository.GetUserRoles"
	query := `
	SELECT r.id, r.name FROM roles r
	INNER JOIN user_roles ur ON ur.role_id = r.id
	WHERE ur.user_id = $1
	`
	var roles []domain.Role
	err := sqlx.SelectContext(ctx, ext, &roles, query, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return roles, nil
}

// Role <=> Permission

func (r *rbacRepository) AssignPermissionToRole(ctx context.Context, ext sqlx.ExtContext, roleID, permissionID string) error {
	op := "RBACRepository.AssignPermissionToRole"
	query := `
	INSERT INTO role_permissions (role_id, permission_id)
	VALUES ($1, $2)
	ON CONFLICT DO NOTHING
	`
	_, err := ext.ExecContext(ctx, query, roleID, permissionID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *rbacRepository) RemovePermissionFromRole(ctx context.Context, ext sqlx.ExtContext, roleID, permissionID string) error {
	op := "RBACRepository.RemovePermissionFromRole"
	query := `
	DELETE FROM role_permissions
	WHERE role_id = $1 AND permission_id = $2
	`
	_, err := ext.ExecContext(ctx, query, roleID, permissionID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *rbacRepository) GetRolePermissions(ctx context.Context, ext sqlx.ExtContext, roleID string) ([]domain.Permission, error) {
	op := "RBACRepository.GetRolePermissions"
	query := `
	SELECT p.id, p.name FROM permissions p
	INNER JOIN role_permissions rp ON rp.permission_id = p.id
	WHERE rp.role_id = $1
	`
	var permissions []domain.Permission
	err := sqlx.SelectContext(ctx, ext, &permissions, query, roleID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return permissions, nil
}

// Existence checks

func (r *rbacRepository) UserExists(ctx context.Context, ext sqlx.ExtContext, userID string) (bool, error) {
	op := "RBACRepository.UserExists"
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`
	var exists bool
	err := sqlx.GetContext(ctx, ext, &exists, query, userID)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return exists, nil
}

func (r *rbacRepository) RoleExists(ctx context.Context, ext sqlx.ExtContext, roleID string) (bool, error) {
	op := "RBACRepository.RoleExists"
	query := `SELECT EXISTS(SELECT 1 FROM roles WHERE id = $1)`
	var exists bool
	err := sqlx.GetContext(ctx, ext, &exists, query, roleID)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return exists, nil
}

func (r *rbacRepository) PermissionExists(ctx context.Context, ext sqlx.ExtContext, permissionID string) (bool, error) {
	op := "RBACRepository.PermissionExists"
	query := `SELECT EXISTS(SELECT 1 FROM permissions WHERE id = $1)`
	var exists bool
	err := sqlx.GetContext(ctx, ext, &exists, query, permissionID)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return exists, nil
}
