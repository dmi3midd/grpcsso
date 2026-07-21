package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/dmi3midd/grpcsso/internal/domain"
	"github.com/dmi3midd/grpcsso/internal/repository"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func (s *rbacService) GetPermissionById(ctx context.Context, permissionId string) (*domain.Permission, error) {
	op := "RBACService.GetPermissionById"
	permission, err := s.permissionRepo.GetById(ctx, nil, permissionId)
	if err != nil {
		if errors.Is(err, repository.ErrPermissionNotFound) {
			return nil, fmt.Errorf("%s: %w", op, ErrPermissionNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return permission, nil
}

func (s *rbacService) CreatePermission(ctx context.Context, name string) (string, error) {
	op := "RBACService.CreatePermission"
	v7uuid, _ := uuid.NewV7()
	id := v7uuid.String()
	permission := &domain.Permission{
		Id:   id,
		Name: name,
	}
	if err := s.permissionRepo.Create(ctx, nil, permission); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *rbacService) DeletePermission(ctx context.Context, permissionId string) (string, error) {
	op := "RBACService.DeletePermission"
	if err := s.permissionRepo.Delete(ctx, nil, permissionId); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return permissionId, nil
}

// Role <=> Permission

func (s *rbacService) AssignPermissionToRole(ctx context.Context, permissionId, roleId string) (string, string, error) {
	op := "RBACService.AssignPermissionToRole"
	err := s.txManager.WithTx(ctx, func(tx *sqlx.Tx) error {
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
	if err != nil {
		return "", "", err
	}
	return permissionId, roleId, nil
}

func (s *rbacService) RevokePermissionFromRole(ctx context.Context, permissionId, roleId string) (string, string, error) {
	op := "RBACService.RevokePermissionFromRole"
	err := s.txManager.WithTx(ctx, func(tx *sqlx.Tx) error {
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
	if err != nil {
		return "", "", err
	}
	return permissionId, roleId, nil
}

func (s *rbacService) GetRolePermissions(ctx context.Context, roleId string) ([]domain.Permission, error) {
	op := "RBACService.GetRolePermissions"
	permissions, err := s.permissionRepo.GetByRole(ctx, s.txManager.GetDB(), roleId)
	if err != nil {
		return []domain.Permission{}, fmt.Errorf("%s: %w", op, err)
	}
	return permissions, nil
}
