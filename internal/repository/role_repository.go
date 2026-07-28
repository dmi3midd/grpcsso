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
	ErrRoleNotFound = errors.New("role not found")
)

type RoleRepository interface {
	// IsExists checks if a role with the given Id exists.
	IsExists(ctx context.Context, roleId string) (bool, error)
	// GetById returns a role by its id.
	// It returns [ErrRoleNotFound] if the role is not found.
	GetById(ctx context.Context, id string) (*domain.Role, error)
	// Create creates a new role.
	Create(ctx context.Context, role *domain.Role) error
	// Delete deletes a role by its id.
	Delete(ctx context.Context, id string) error

	// Assign inserts a (user_id, role_id) record into user_roles.
	Assign(ctx context.Context, userId, roleId string) error
	// Revoke deletes a record from user_roles.
	Revoke(ctx context.Context, userId, roleId string) error
	// GetByUser returns all roles assigned to a user.
	// Returns empty slice if no roles are found.
	GetByUser(ctx context.Context, userId string) ([]domain.Role, error)
}

type roleRepository struct {
	db *sqlx.DB
}

func NewRoleRepo(db *sqlx.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) IsExists(ctx context.Context, roleId string) (bool, error) {
	op := "RoleRepository.RoleExists"
	query := `SELECT EXISTS(SELECT 1 FROM roles WHERE id = $1)`
	executor := ExtractTx(ctx, r.db)
	var exists bool
	err := sqlx.GetContext(ctx, executor, &exists, query, roleId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return exists, nil
}

func (r *roleRepository) GetById(ctx context.Context, id string) (*domain.Role, error) {
	op := "RoleRepository.GetById"
	query := `
	SELECT id, name FROM roles
	WHERE id = $1
	`
	executor := ExtractTx(ctx, r.db)
	role := &domain.Role{}
	err := sqlx.GetContext(ctx, executor, role, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrRoleNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return role, nil
}

func (r *roleRepository) Create(ctx context.Context, role *domain.Role) error {
	op := "RoleRepository.Create"
	query := `
	INSERT INTO roles (id, name)
	VALUES (:id, :name)
	`
	executor := ExtractTx(ctx, r.db)
	_, err := sqlx.NamedExecContext(ctx, executor, query, role)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *roleRepository) Delete(ctx context.Context, id string) error {
	op := "RoleRepository.Delete"
	query := `
	DELETE FROM roles
	WHERE id = $1
	`
	executor := ExtractTx(ctx, r.db)
	_, err := executor.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *roleRepository) Assign(ctx context.Context, userId, roleId string) error {
	op := "RoleRepository.Assign"
	query := `
	INSERT INTO user_roles (user_id, role_id)
	VALUES ($1, $2)
	ON CONFLICT DO NOTHING
	`
	executor := ExtractTx(ctx, r.db)
	_, err := executor.ExecContext(ctx, query, userId, roleId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *roleRepository) Revoke(ctx context.Context, userId, roleId string) error {
	op := "RoleRepository.Revoke"
	query := `
	DELETE FROM user_roles
	WHERE user_id = $1 AND role_id = $2
	`
	executor := ExtractTx(ctx, r.db)
	_, err := executor.ExecContext(ctx, query, userId, roleId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *roleRepository) GetByUser(ctx context.Context, userId string) ([]domain.Role, error) {
	op := "RoleRepository.GetByUser"
	query := `
	SELECT r.id, r.name FROM roles r
	INNER JOIN user_roles ur ON ur.role_id = r.id
	WHERE ur.user_id = $1
	`
	executor := ExtractTx(ctx, r.db)
	var roles []domain.Role
	err := sqlx.SelectContext(ctx, executor, &roles, query, userId)
	if err != nil {
		return []domain.Role{}, fmt.Errorf("%s: %w", op, err)
	}
	return roles, nil
}
