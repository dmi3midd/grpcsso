package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/dmi3midd/grpcsso/internal/domain"
	"github.com/dmi3midd/grpcsso/internal/repository"
	"github.com/jmoiron/sqlx"
)

func (s *rbacService) GetRoleById(ctx context.Context, roleId string) (*domain.Role, error) {
	op := "RBACService.GetRoleById"
	role, err := s.roleRepo.GetById(ctx, s.txManager.GetDB(), roleId)
	if err != nil {
		if errors.Is(err, repository.ErrRoleNotFound) {
			return nil, fmt.Errorf("%s: %w", op, ErrRoleNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return role, nil
}

func (s *rbacService) CreateRole(ctx context.Context, role *domain.Role) error {
	op := "RBACService.CreateRole"
	if err := s.roleRepo.Create(ctx, s.txManager.GetDB(), role); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *rbacService) DeleteRole(ctx context.Context, roleId string) error {
	op := "RBACService.DeleteRole"
	if err := s.roleRepo.Delete(ctx, s.txManager.GetDB(), roleId); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
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
