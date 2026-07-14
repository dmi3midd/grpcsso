package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/dmi3midd/grpcsso/internal/domain"
	"github.com/dmi3midd/grpcsso/internal/repository"
	"github.com/jmoiron/sqlx"
)

func (s *rbacService) GetPermissionById(ctx context.Context, id string) (*domain.Permission, error) {
	op := "RBACService.GetPermissionById"
	permission, err := s.permissionRepo.GetById(ctx, nil, id)
	if err != nil {
		if errors.Is(err, repository.ErrPermissionNotFound) {
			return nil, fmt.Errorf("%s: %w", op, ErrPermissionNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return permission, nil
}

func (s *rbacService) CreatePermission(ctx context.Context, permission *domain.Permission) error {
	op := "RBACService.CreatePermission"
	if err := s.permissionRepo.Create(ctx, nil, permission); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *rbacService) DeletePermission(ctx context.Context, id string) error {
	op := "RBACService.DeletePermission"
	if err := s.permissionRepo.Delete(ctx, nil, id); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
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
