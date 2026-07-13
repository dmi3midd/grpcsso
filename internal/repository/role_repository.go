package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/dmi3midd/grpcsso/internal/domain"
	"github.com/jmoiron/sqlx"
)

var (
	ErrRoleNotFound = errors.New("role not found")
)

type RoleRepository interface {
	// GetById returns a role by its id.
	// It returns [ErrRoleNotFound] if the role is not found.
	GetById(ctx context.Context, id string) (*domain.Role, error)
	// Create creates a new role.
	Create(ctx context.Context, role *domain.Role) error
	// Delete deletes a role by its id.
	Delete(ctx context.Context, id string) error
}

type roleRepository struct {
	db *sqlx.DB
}

func NewRoleRepository(db *sqlx.DB) RoleRepository {
	return &roleRepository{
		db: db,
	}
}

func (r *roleRepository) GetById(ctx context.Context, id string) (*domain.Role, error) {
	op := "RoleRepository.GetById"
	query := `
	SELECT id, name FROM roles
	WHERE id = $1
	`
	role := &domain.Role{}
	err := r.db.GetContext(ctx, role, query, id)
	if err != nil {
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
	_, err := r.db.NamedExecContext(ctx, query, role)
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
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
