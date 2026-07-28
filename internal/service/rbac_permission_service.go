package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/dmi3midd/grpcsso/internal/domain"
	"github.com/dmi3midd/grpcsso/internal/repository"
	"github.com/google/uuid"
)

func (s *rbacService) GetPermissionById(ctx context.Context, permissionId string) (*domain.Permission, error) {
	op := "RBACService.GetPermissionById"
	// Try to get permission from cache
	permissionFromCache, err := s.permissionCache.GetById(ctx, permissionId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if permissionFromCache != nil {
		return permissionFromCache, nil
	}
	// Try to get permission from database
	permission, err := s.permissionRepo.GetById(ctx, permissionId)
	if err != nil {
		if errors.Is(err, repository.ErrPermissionNotFound) {
			return nil, fmt.Errorf("%s: %w", op, ErrPermissionNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	// Set permission to cache
	if err := s.permissionCache.Create(ctx, permission); err != nil {
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
	// Create permission in cache
	if err := s.permissionCache.Create(ctx, permission); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	// Create permission in database
	if err := s.permissionRepo.Create(ctx, permission); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *rbacService) DeletePermission(ctx context.Context, permissionId string) (string, error) {
	op := "RBACService.DeletePermission"
	// Delete permission from cache
	if err := s.permissionCache.Delete(ctx, permissionId); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	// Delete permission from database
	if err := s.permissionRepo.Delete(ctx, permissionId); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return permissionId, nil
}

// Role <=> Permission

func (s *rbacService) AssignPermissionToRole(ctx context.Context, permissionId, roleId string) (string, string, error) {
	op := "RBACService.AssignPermissionToRole"
	// Assign permission to role in cache
	if err := s.permissionCache.Assign(ctx, roleId, permissionId); err != nil {
		return "", "", fmt.Errorf("%s: %w", op, err)
	}
	// Assign permission to role in database
	err := s.txManager.WithTx(ctx, func(txCtx context.Context) error {
		exists, err := s.roleRepo.IsExists(txCtx, roleId)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if !exists {
			return fmt.Errorf("%s: %w", op, ErrRoleNotFound)
		}

		exists, err = s.permissionRepo.IsExists(txCtx, permissionId)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if !exists {
			return fmt.Errorf("%s: %w", op, ErrPermissionNotFound)
		}

		if err := s.permissionRepo.Assign(txCtx, roleId, permissionId); err != nil {
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
	// Revoke permission from role in cache
	if err := s.permissionCache.Revoke(ctx, roleId, permissionId); err != nil {
		return "", "", fmt.Errorf("%s: %w", op, err)
	}
	// Revoke permission from role in database
	err := s.txManager.WithTx(ctx, func(txCtx context.Context) error {
		exists, err := s.roleRepo.IsExists(txCtx, roleId)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if !exists {
			return fmt.Errorf("%s: %w", op, ErrRoleNotFound)
		}

		exists, err = s.permissionRepo.IsExists(txCtx, permissionId)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if !exists {
			return fmt.Errorf("%s: %w", op, ErrPermissionNotFound)
		}

		if err := s.permissionRepo.Revoke(txCtx, roleId, permissionId); err != nil {
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
	// Get permissions from cache
	permissionsFromCache, err := s.permissionCache.GetByRole(ctx, roleId)
	if err != nil {
		return []domain.Permission{}, fmt.Errorf("%s: %w", op, err)
	}
	if len(permissionsFromCache) > 0 {
		return permissionsFromCache, nil
	}
	// Get permissions from database
	permissions, err := s.permissionRepo.GetByRole(ctx, roleId)
	if err != nil {
		return []domain.Permission{}, fmt.Errorf("%s: %w", op, err)
	}
	// Set permissions to role in cache
	for _, permission := range permissions {
		if err := s.permissionCache.Assign(ctx, roleId, permission.Id); err != nil {
			// Need to log error because it's not critical to return error
			return permissions, nil
		}
	}
	return permissions, nil
}
