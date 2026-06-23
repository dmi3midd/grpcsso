package repository

import (
	"context"
	"fmt"

	"github.com/dmi3midd/grpcsso/internal/domain"
	"github.com/jmoiron/sqlx"
)

type PermissionRepository interface {
	// GetByUserIdAndClientId retrieves all permissions for a user for a specific client.
	GetByUserIdAndClientId(ctx context.Context, userId, clientId string) ([]domain.Permission, error)
	// CreateMany creates multiple permissions for a user for a specific client.
	CreateMany(ctx context.Context, permissions []domain.Permission) error
	// DeleteAllForUserAndClient deletes all permissions for a user for a specific client.
	DeleteAllForUserAndClient(ctx context.Context, userId, clientId string) error
}

type permissionRepository struct {
	db *sqlx.DB
}

func NewPermissionRepo(db *sqlx.DB) PermissionRepository {
	return &permissionRepository{
		db: db,
	}
}

func (r *permissionRepository) GetByUserIdAndClientId(ctx context.Context, userId, clientId string) ([]domain.Permission, error) {
	op := "PermissionRepository.GetByUserIdAndClientId"
	query := `
	SELECT id, user_id, client_id, permission, created_at, updated_at
	FROM permissions 
	WHERE user_id = $1 AND client_id = $2
	`
	var permissions []domain.Permission
	err := r.db.SelectContext(ctx, &permissions, query, userId, clientId)
	if err != nil {
		return []domain.Permission{}, fmt.Errorf("%s: %w", op, err)
	}

	return permissions, nil
}

func (r *permissionRepository) CreateMany(ctx context.Context, permissions []domain.Permission) error {
	op := "PermissionRepository.CreateMany"
	if len(permissions) == 0 {
		return nil
	}
	query := `
	INSERT INTO permissions (id, user_id, client_id, permission, created_at, updated_at)
	VALUES (:id, :user_id, :client_id, :permission, :created_at, :updated_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, permissions)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *permissionRepository) DeleteAllForUserAndClient(ctx context.Context, userId, clientId string) error {
	op := "PermissionRepository.DeleteAllForUserAndClient"
	query := `
	DELETE FROM permissions 
	WHERE user_id = $1 AND client_id = $2
	`
	_, err := r.db.ExecContext(ctx, query, userId, clientId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
