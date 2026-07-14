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
	txManager      repository.TxManager
	userRepo       repository.UserRepository
	roleRepo       repository.RoleRepository
	permissionRepo repository.PermissionRepository
}

func NewRBACService(
	txManager repository.TxManager,
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	permissionRepo repository.PermissionRepository,
) RBACService {
	return &rbacService{
		txManager:      txManager,
		userRepo:       userRepo,
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
	}
}

// User <=> Role

func (s *rbacService) AssignRoleToUser(ctx context.Context, userId, roleId string) error {
	op := "RBACService.AssignRoleToUser"
	return s.txManager.WithTx(ctx, func(tx *sqlx.Tx) error {
		exists, err := s.userRepo.IsExists(ctx, tx, userId)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if !exists {
			return fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		exists, err = s.roleRepo.IsExists(ctx, tx, roleId)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if !exists {
			return fmt.Errorf("%s: %w", op, ErrRoleNotFound)
		}

		if err := s.roleRepo.Assign(ctx, tx, userId, roleId); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		return nil
	})
}

func (s *rbacService) RemoveRoleFromUser(ctx context.Context, userId, roleId string) error {
	op := "RBACService.RemoveRoleFromUser"
	return s.txManager.WithTx(ctx, func(tx *sqlx.Tx) error {
		exists, err := s.userRepo.IsExists(ctx, tx, userId)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if !exists {
			return fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		exists, err = s.roleRepo.IsExists(ctx, tx, roleId)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if !exists {
			return fmt.Errorf("%s: %w", op, ErrRoleNotFound)
		}

		if err := s.roleRepo.Revoke(ctx, tx, userId, roleId); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		return nil
	})
}

func (s *rbacService) GetUserRoles(ctx context.Context, userId string) ([]domain.Role, error) {
	op := "RBACService.GetUserRoles"
	roles, err := s.roleRepo.GetByUser(ctx, s.txManager.GetDB(), userId)
	if err != nil {
		return []domain.Role{}, fmt.Errorf("%s: %w", op, err)
	}
	return roles, nil
}

// Role <=> Permission

func (s *rbacService) AssignPermissionToRole(ctx context.Context, roleId, permissionId string) error {
	op := "RBACService.AssignPermissionToRole"
	return s.txManager.WithTx(ctx, func(tx *sqlx.Tx) error {
		exists, err := s.roleRepo.IsExists(ctx, tx, roleId)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if !exists {
			return fmt.Errorf("%s: %w", op, ErrRoleNotFound)
		}

		exists, err = s.permissionRepo.IsExists(ctx, tx, permissionId)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if !exists {
			return fmt.Errorf("%s: %w", op, ErrPermissionNotFound)
		}

		if err := s.permissionRepo.Assign(ctx, tx, roleId, permissionId); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		return nil
	})
}

func (s *rbacService) RemovePermissionFromRole(ctx context.Context, roleId, permissionId string) error {
	op := "RBACService.RemovePermissionFromRole"
	return s.txManager.WithTx(ctx, func(tx *sqlx.Tx) error {
		exists, err := s.roleRepo.IsExists(ctx, tx, roleId)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if !exists {
			return fmt.Errorf("%s: %w", op, ErrRoleNotFound)
		}

		exists, err = s.permissionRepo.IsExists(ctx, tx, permissionId)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if !exists {
			return fmt.Errorf("%s: %w", op, ErrPermissionNotFound)
		}

		if err := s.permissionRepo.Revoke(ctx, tx, roleId, permissionId); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		return nil
	})
}

func (s *rbacService) GetRolePermissions(ctx context.Context, roleId string) ([]domain.Permission, error) {
	op := "RBACService.GetRolePermissions"
	permissions, err := s.permissionRepo.GetByRole(ctx, s.txManager.GetDB(), roleId)
	if err != nil {
		return []domain.Permission{}, fmt.Errorf("%s: %w", op, err)
	}
	return permissions, nil
}
