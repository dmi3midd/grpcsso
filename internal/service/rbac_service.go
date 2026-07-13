package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/dmi3midd/grpcsso/internal/domain"
	"github.com/dmi3midd/grpcsso/internal/repository"
	"github.com/jmoiron/sqlx"
)

var (
	ErrRoleNotFound       = errors.New("role not found")
	ErrPermissionNotFound = errors.New("permission not found")
)

type RBACService interface {
	// AssignRoleToUser assigns a role to a user within a transaction.
	// Returns [ErrUserNotFound] if the user does not exist.
	// Returns [ErrRoleNotFound] if the role does not exist.
	AssignRoleToUser(ctx context.Context, userId, roleId string) error
	// RemoveRoleFromUser removes a role from a user within a transaction.
	// Returns [ErrUserNotFound] if the user does not exist.
	// Returns [ErrRoleNotFound] if the role does not exist.
	RemoveRoleFromUser(ctx context.Context, userId, roleId string) error
	// GetUserRoles returns all roles assigned to a user.
	// Returns an empty slice if the user has no roles.
	GetUserRoles(ctx context.Context, userId string) ([]domain.Role, error)

	// AssignPermissionToRole assigns a permission to a role within a transaction.
	// Returns [ErrRoleNotFound] if the role does not exist.
	// Returns [ErrPermissionNotFound] if the permission does not exist.
	AssignPermissionToRole(ctx context.Context, roleId, permissionId string) error
	// RemovePermissionFromRole removes a permission from a role within a transaction.
	// Returns [ErrRoleNotFound] if the role does not exist.
	// Returns [ErrPermissionNotFound] if the permission does not exist.
	RemovePermissionFromRole(ctx context.Context, roleId, permissionId string) error
	// GetRolePermissions returns all permissions assigned to a role.
	// Returns an empty slice if the role has no permissions.
	GetRolePermissions(ctx context.Context, roleId string) ([]domain.Permission, error)
}

type rbacService struct {
	db       *sqlx.DB
	rbacRepo repository.RBACRepository
}

func NewRBACService(db *sqlx.DB, rbacRepo repository.RBACRepository) RBACService {
	return &rbacService{
		db:       db,
		rbacRepo: rbacRepo,
	}
}

// withTx executes txFn inside a database transaction.
func (s *rbacService) withTx(ctx context.Context, txFn func(tx *sqlx.Tx) error) error {
	op := "rbacService.withTx"
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()

	if err := txFn(tx); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return tx.Commit()
}

// User <=> Role

func (s *rbacService) AssignRoleToUser(ctx context.Context, userId, roleId string) error {
	op := "RBACService.AssignRoleToUser"
	return s.withTx(ctx, func(tx *sqlx.Tx) error {
		exists, err := s.rbacRepo.UserExists(ctx, tx, userId)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if !exists {
			return fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		exists, err = s.rbacRepo.RoleExists(ctx, tx, roleId)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if !exists {
			return fmt.Errorf("%s: %w", op, ErrRoleNotFound)
		}

		if err := s.rbacRepo.AssignRoleToUser(ctx, tx, userId, roleId); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		return nil
	})
}

func (s *rbacService) RemoveRoleFromUser(ctx context.Context, userId, roleId string) error {
	op := "RBACService.RemoveRoleFromUser"
	return s.withTx(ctx, func(tx *sqlx.Tx) error {
		exists, err := s.rbacRepo.UserExists(ctx, tx, userId)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if !exists {
			return fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		exists, err = s.rbacRepo.RoleExists(ctx, tx, roleId)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if !exists {
			return fmt.Errorf("%s: %w", op, ErrRoleNotFound)
		}

		if err := s.rbacRepo.RemoveRoleFromUser(ctx, tx, userId, roleId); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		return nil
	})
}

func (s *rbacService) GetUserRoles(ctx context.Context, userId string) ([]domain.Role, error) {
	op := "RBACService.GetUserRoles"
	roles, err := s.rbacRepo.GetUserRoles(ctx, s.db, userId)
	if err != nil {
		return []domain.Role{}, fmt.Errorf("%s: %w", op, err)
	}
	return roles, nil
}

// Role <=> Permission

func (s *rbacService) AssignPermissionToRole(ctx context.Context, roleId, permissionId string) error {
	op := "RBACService.AssignPermissionToRole"
	return s.withTx(ctx, func(tx *sqlx.Tx) error {
		exists, err := s.rbacRepo.RoleExists(ctx, tx, roleId)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if !exists {
			return fmt.Errorf("%s: %w", op, ErrRoleNotFound)
		}

		exists, err = s.rbacRepo.PermissionExists(ctx, tx, permissionId)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if !exists {
			return fmt.Errorf("%s: %w", op, ErrPermissionNotFound)
		}

		if err := s.rbacRepo.AssignPermissionToRole(ctx, tx, roleId, permissionId); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		return nil
	})
}

func (s *rbacService) RemovePermissionFromRole(ctx context.Context, roleId, permissionId string) error {
	op := "RBACService.RemovePermissionFromRole"
	return s.withTx(ctx, func(tx *sqlx.Tx) error {
		exists, err := s.rbacRepo.RoleExists(ctx, tx, roleId)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if !exists {
			return fmt.Errorf("%s: %w", op, ErrRoleNotFound)
		}

		exists, err = s.rbacRepo.PermissionExists(ctx, tx, permissionId)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if !exists {
			return fmt.Errorf("%s: %w", op, ErrPermissionNotFound)
		}

		if err := s.rbacRepo.RemovePermissionFromRole(ctx, tx, roleId, permissionId); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		return nil
	})
}

func (s *rbacService) GetRolePermissions(ctx context.Context, roleId string) ([]domain.Permission, error) {
	op := "RBACService.GetRolePermissions"
	permissions, err := s.rbacRepo.GetRolePermissions(ctx, s.db, roleId)
	if err != nil {
		return []domain.Permission{}, fmt.Errorf("%s: %w", op, err)
	}
	return permissions, nil
}
