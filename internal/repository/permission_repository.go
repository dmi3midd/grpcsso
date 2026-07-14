package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/dmi3midd/grpcsso/internal/domain"
	"github.com/jmoiron/sqlx"
)

var (
	ErrPermissionNotFound = fmt.Errorf("permission not found")
)

type PermissionRepository interface {
	// IsExists checks if a permission with the given Id exists.
	IsExists(ctx context.Context, ext sqlx.ExtContext, permissionId string) (bool, error)
	// Get a permission by its id.
	// Returns ErrPermissionNotFound if the permission is not found.
	GetById(ctx context.Context, ext sqlx.ExtContext, id string) (*domain.Permission, error)
	// Create a new permission.
	Create(ctx context.Context, ext sqlx.ExtContext, permission *domain.Permission) error
	// Delete a permission by its id.
	Delete(ctx context.Context, ext sqlx.ExtContext, id string) error

	// Assign inserts a (role_id, permission_id) record into role_permissions.
	Assign(ctx context.Context, ext sqlx.ExtContext, roleId, permissionId string) error
	// Revoke deletes a record from role_permissions.
	Revoke(ctx context.Context, ext sqlx.ExtContext, roleId, permissionId string) error
	// GetByRole returns all permissions assigned to a role.
	// Returns empty slice if no permissions are found.
	GetByRole(ctx context.Context, ext sqlx.ExtContext, roleId string) ([]domain.Permission, error)
}

type permissionRepository struct {
}

func NewPermissionRepo() PermissionRepository {
	return &permissionRepository{}
}

func (r *permissionRepository) IsExists(ctx context.Context, ext sqlx.ExtContext, permissionId string) (bool, error) {
	op := "PermissionRepository.IsExists"
	query := `SELECT EXISTS(SELECT 1 FROM permissions WHERE id = $1)`
	var exists bool
	err := sqlx.GetContext(ctx, ext, &exists, query, permissionId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return exists, nil
}

func (r *permissionRepository) GetById(ctx context.Context, ext sqlx.ExtContext, id string) (*domain.Permission, error) {
	op := "PermissionRepository.GetById"
	query := `
	SELECT id, name FROM permissions
	WHERE id = $1
	`
	permission := &domain.Permission{}
	err := sqlx.GetContext(ctx, ext, permission, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrPermissionNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return permission, nil
}

func (r *permissionRepository) Create(ctx context.Context, ext sqlx.ExtContext, permission *domain.Permission) error {
	op := "PermissionRepository.Create"
	query := `
	INSERT INTO permissions (id, name)
	VALUES (:id, :name)
	`
	_, err := sqlx.NamedExecContext(ctx, ext, query, permission)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *permissionRepository) Delete(ctx context.Context, ext sqlx.ExtContext, id string) error {
	op := "PermissionRepository.Delete"
	query := `
	DELETE FROM permissions
	WHERE id = $1
	`
	_, err := ext.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *permissionRepository) Assign(ctx context.Context, ext sqlx.ExtContext, roleId, permissionId string) error {
	op := "PermissionRepository.Assign"
	query := `
	INSERT INTO role_permissions (role_id, permission_id)
	VALUES ($1, $2)
	ON CONFLICT DO NOTHING
	`
	_, err := ext.ExecContext(ctx, query, roleId, permissionId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *permissionRepository) Revoke(ctx context.Context, ext sqlx.ExtContext, roleId, permissionId string) error {
	op := "PermissionRepository.Revoke"
	query := `
	DELETE FROM role_permissions
	WHERE role_id = $1 AND permission_id = $2
	`
	_, err := ext.ExecContext(ctx, query, roleId, permissionId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *permissionRepository) GetByRole(ctx context.Context, ext sqlx.ExtContext, roleId string) ([]domain.Permission, error) {
	op := "PermissionRepository.GetByRole"
	query := `
	SELECT p.id, p.name FROM permissions p
	INNER JOIN role_permissions rp ON rp.permission_id = p.id
	WHERE rp.role_id = $1
	`
	var permissions []domain.Permission
	err := sqlx.SelectContext(ctx, ext, &permissions, query, roleId)
	if err != nil {
		return []domain.Permission{}, fmt.Errorf("%s: %w", op, err)
	}
	return permissions, nil
}
