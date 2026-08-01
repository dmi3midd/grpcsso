package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/dmi3midd/grpcsso/internal/domain"
	"github.com/dmi3midd/grpcsso/internal/repository"
	"github.com/google/uuid"
)

func (s *rbacService) GetRoleById(ctx context.Context, roleId string) (*domain.Role, error) {
	op := "RBACService.GetRoleById"
	// Try to get role from cache
	roleFromCache, err := s.roleCache.GetById(ctx, roleId)
	if err == nil && roleFromCache != nil {
		return roleFromCache, nil
	}

	// Try to get role from database
	role, err := s.roleRepo.GetById(ctx, roleId)
	if err != nil {
		if errors.Is(err, repository.ErrRoleNotFound) {
			return nil, fmt.Errorf("%s: %w", op, ErrRoleNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Set role to cache
	_ = s.roleCache.Create(ctx, role)
	return role, nil
}

func (s *rbacService) CreateRole(ctx context.Context, name string) (string, error) {
	op := "RBACService.CreateRole"
	v7uuid, _ := uuid.NewV7()
	id := v7uuid.String()
	role := &domain.Role{
		Id:   id,
		Name: name,
	}

	// Create role in database
	if err := s.roleRepo.Create(ctx, role); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	// Create role in cache
	_ = s.roleCache.Create(ctx, role)
	return id, nil
}

func (s *rbacService) DeleteRole(ctx context.Context, roleId string) (string, error) {
	op := "RBACService.DeleteRole"
	// Delete role from database
	if err := s.roleRepo.Delete(ctx, roleId); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	// Delete role from cache
	_ = s.roleCache.Delete(ctx, roleId)
	return roleId, nil
}

// User <=> Role

func (s *rbacService) AssignRoleToUser(ctx context.Context, roleId, userId string) (string, string, error) {
	op := "RBACService.AssignRoleToUser"
	// Assign role to user in database with transaction
	err := s.txManager.WithTx(ctx, func(txCtx context.Context) error {
		exists, err := s.userRepo.IsExists(txCtx, userId)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if !exists {
			return fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		exists, err = s.roleRepo.IsExists(txCtx, roleId)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if !exists {
			return fmt.Errorf("%s: %w", op, ErrRoleNotFound)
		}

		if err := s.roleRepo.Assign(txCtx, userId, roleId); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		return nil
	})
	if err != nil {
		return "", "", err
	}

	// Assign role to user in cache
	_ = s.roleCache.Assign(ctx, userId, roleId)
	return roleId, userId, nil
}

func (s *rbacService) RevokeRoleFromUser(ctx context.Context, roleId, userId string) (string, string, error) {
	op := "RBACService.RevokeRoleFromUser"
	// Revoke role from user in database with transaction
	err := s.txManager.WithTx(ctx, func(txCtx context.Context) error {
		exists, err := s.userRepo.IsExists(txCtx, userId)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if !exists {
			return fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		exists, err = s.roleRepo.IsExists(txCtx, roleId)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if !exists {
			return fmt.Errorf("%s: %w", op, ErrRoleNotFound)
		}

		if err := s.roleRepo.Revoke(txCtx, userId, roleId); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		return nil
	})
	if err != nil {
		return "", "", err
	}

	// Revoke role from user in cache
	_ = s.roleCache.Revoke(ctx, userId, roleId)
	return roleId, userId, nil
}

func (s *rbacService) GetUserRoles(ctx context.Context, userId string) ([]domain.Role, error) {
	op := "RBACService.GetUserRoles"
	// Try to get roles from cache
	rolesFromCache, err := s.roleCache.GetByUser(ctx, userId)
	if err == nil && len(rolesFromCache) > 0 {
		return rolesFromCache, nil
	}

	// Try to get roles from database
	roles, err := s.roleRepo.GetByUser(ctx, userId)
	if err != nil {
		return []domain.Role{}, fmt.Errorf("%s: %w", op, err)
	}

	for _, role := range roles {
		_ = s.roleCache.Assign(ctx, userId, role.Id)
	}
	return roles, nil
}
