package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/dmi3midd/grpcsso/internal/domain"
	"github.com/redis/go-redis/v9"
)

type permissionCache struct {
	redisClient *redis.Client
}

func NewPermissionCache(redisClient *redis.Client) PermissionRepository {
	return &permissionCache{redisClient: redisClient}
}

func permKey(id string) string {
	return "perm:info:" + id
}

func rolePermissionsKey(roleId string) string {
	return "role:permissions:" + roleId
}

// IsExists checks if a permission exists by its id.
func (p *permissionCache) IsExists(ctx context.Context, permissionId string) (bool, error) {
	op := "PermissionCache.IsExists"
	exists, err := p.redisClient.Exists(ctx, permKey(permissionId)).Result()
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return exists > 0, nil
}

// GetById retrieves a permission by its id.
// Returns ErrPermissionNotFound if the permission is not found.
func (p *permissionCache) GetById(ctx context.Context, id string) (*domain.Permission, error) {
	op := "PermissionCache.GetById"
	permission, err := p.redisClient.Get(ctx, permKey(id)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, fmt.Errorf("%s: %w", op, ErrPermissionNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &domain.Permission{
		Id:   id,
		Name: permission,
	}, nil
}

// Create sets a permission permissionId->name key-value.
func (p *permissionCache) Create(ctx context.Context, permission *domain.Permission) error {
	op := "PermissionCache.Create"
	err := p.redisClient.Set(ctx, permKey(permission.Id), permission.Name, 0).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// Delete removes a permission by its id.
func (p *permissionCache) Delete(ctx context.Context, id string) error {
	op := "PermissionCache.Delete"
	err := p.redisClient.Del(ctx, permKey(id)).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// Assign adds a permission to a role.
func (p *permissionCache) Assign(ctx context.Context, roleId string, permissionId string) error {
	op := "PermissionCache.Assign"
	err := p.redisClient.SAdd(ctx, rolePermissionsKey(roleId), permissionId).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// Revoke removes a permission from a role.
func (p *permissionCache) Revoke(ctx context.Context, roleId string, permissionId string) error {
	op := "PermissionCache.Revoke"
	err := p.redisClient.SRem(ctx, rolePermissionsKey(roleId), permissionId).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// GetByRole retrieves all permissions assigned to a role.
func (p *permissionCache) GetByRole(ctx context.Context, roleId string) ([]domain.Permission, error) {
	op := "PermissionCache.GetByRole"
	permissionIds, err := p.redisClient.SMembers(ctx, rolePermissionsKey(roleId)).Result()
	if err != nil {
		return []domain.Permission{}, fmt.Errorf("%s: %w", op, err)
	}
	permissions := make([]domain.Permission, 0, len(permissionIds))
	for _, permissionId := range permissionIds {
		permission, err := p.GetById(ctx, permissionId)
		if err != nil {
			return []domain.Permission{}, fmt.Errorf("%s: %w", op, err)
		}
		permissions = append(permissions, *permission)
	}
	return permissions, nil
}
