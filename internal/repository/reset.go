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
	ErrResetNotFound error = errors.New("reset not found")
)

type ResetRepository interface {
	// FindValidByTokenHash finds a Reset entity by its token hash.
	// Returns [ErrResetNotFound] if no reset is found.
	FindValidByTokenHash(ctx context.Context, tokenHash string) (*domain.Reset, error)
	// Create creates a Reset entity and returns its id.
	Create(ctx context.Context, reset *domain.Reset) (string, error)
	// Update updates the Reset entity.
	Update(ctx context.Context, reset *domain.Reset) error
}

type resetRepository struct {
	db *sqlx.DB
}

func NewResetRepo(db *sqlx.DB) ResetRepository {
	return &resetRepository{
		db: db,
	}
}

func (r *resetRepository) FindValidByTokenHash(ctx context.Context, tokenHash string) (*domain.Reset, error) {
	op := "ResetRepository.FindValidByTokenHash"
	var reset domain.Reset
	query := `
	SELECT id, user_id, token_hash, expires_at, used_at, created_at
	FROM password_resets 
	WHERE token_hash = $1
	`
	err := r.db.GetContext(ctx, &reset, query, tokenHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrResetNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &reset, nil
}

func (r *resetRepository) Create(ctx context.Context, reset *domain.Reset) (string, error) {
	op := "ResetRepository.Create"
	query := `
	INSERT INTO password_resets 
	(id, user_id, token_hash, expires_at, used_at, created_at)
	VALUES (:id, :user_id, :token_hash, :expires_at, :used_at, :created_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, reset)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return reset.Id, nil
}

func (r *resetRepository) Update(ctx context.Context, reset *domain.Reset) error {
	op := "ResetRepository.Update"
	query := `
	UPDATE password_resets 
	SET token_hash = :token_hash, used_at = :used_at
	WHERE id = :id
	`
	_, err := r.db.NamedExecContext(ctx, query, reset)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
