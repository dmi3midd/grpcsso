package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dmi3midd/grpcsso/internal/domain"
	"github.com/jmoiron/sqlx"
)

var (
	ErrPermissionNotFound = fmt.Errorf("permission not found")
)

type PermissionRepository interface {
	// Get a permission by its id.
	// Returns ErrPermissionNotFound if the permission is not found.
	GetById(ctx context.Context, id string) (*domain.Permission, error)
	// Create a new permission.
	Create(ctx context.Context, permission *domain.Permission) error

	// Delete a permission by its id.
	Delete(ctx context.Context, id string) error
}

type permissionRepository struct {
	db *sqlx.DB
}

func NewPermissionRepo(db *sqlx.DB) PermissionRepository {
	return &permissionRepository{
		db: db,
	}
}

func (r *permissionRepository) GetById(ctx context.Context, id string) (*domain.Permission, error) {
	op := "PermissionRepository.GetById"
	query := `
	SELECT id, name FROM permissions
	WHERE id = $1
	`
	permission := &domain.Permission{}
	err := r.db.GetContext(ctx, permission, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s: %w", op, ErrPermissionNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return permission, nil
}

func (r *permissionRepository) Create(ctx context.Context, permission *domain.Permission) error {
	op := "PermissionRepository.Create"
	query := `
	INSERT INTO permissions (id, name)
	VALUES (:id, :name)
	`
	_, err := r.db.NamedExecContext(ctx, query, permission)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *permissionRepository) Delete(ctx context.Context, id string) error {
	op := "PermissionRepository.Delete"
	query := `
	DELETE FROM permissions
	WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
