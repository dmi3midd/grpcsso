package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/dmi3midd/grpcsso/internal/domain"
	"github.com/redis/go-redis/v9"
)

// for crud operations need to write key-value:
// id -> name
// for assign/revoke/getByUser need to write into set:
// userId -> set{roleId1, roleId2, ...}
type roleCache struct {
	redisClient *redis.Client
}

func NewRoleCache(redisClient *redis.Client) RoleRepository {
	return &roleCache{redisClient: redisClient}
}

// IsExists checks if a role exists by its id.
func (r *roleCache) IsExists(ctx context.Context, roleId string) (bool, error) {
	op := "RoleCache.IsExists"
	exists, err := r.redisClient.Exists(ctx, roleId).Result()
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return exists > 0, nil
}

// GetById retrieves a role by its id.
// Returns ErrRoleNotFound if the role is not found.
func (r *roleCache) GetById(ctx context.Context, id string) (*domain.Role, error) {
	op := "RoleCache.GetById"
	name, err := r.redisClient.Get(ctx, id).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, fmt.Errorf("%s: %w", op, ErrRoleNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &domain.Role{
		Id:   id,
		Name: name,
	}, nil
}

// Create sets a role roleId->name key-value.
func (r *roleCache) Create(ctx context.Context, role *domain.Role) error {
	op := "RoleCache.Create"
	err := r.redisClient.Set(ctx, role.Id, role.Name, 0).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// Delete removes a role by its id.
func (r *roleCache) Delete(ctx context.Context, id string) error {
	op := "RoleCache.Delete"
	err := r.redisClient.Del(ctx, id).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// Assign adds a role to a user.
func (r *roleCache) Assign(ctx context.Context, userId string, roleId string) error {
	op := "RoleCache.Assign"
	err := r.redisClient.SAdd(ctx, userId, roleId).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// Revoke removes a role from a user.
func (r *roleCache) Revoke(ctx context.Context, userId string, roleId string) error {
	op := "RoleCache.Revoke"
	err := r.redisClient.SRem(ctx, userId, roleId).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

// GetByUser retrieves all roles assigned to a user.
func (r *roleCache) GetByUser(ctx context.Context, userId string) ([]domain.Role, error) {
	op := "RoleCache.GetByUser"
	roleIds, err := r.redisClient.SMembers(ctx, userId).Result()
	if err != nil {
		return []domain.Role{}, fmt.Errorf("%s: %w", op, err)
	}
	roles := make([]domain.Role, 0, len(roleIds))
	for _, roleId := range roleIds {
		role, err := r.GetById(ctx, roleId)
		if err != nil {
			return []domain.Role{}, fmt.Errorf("%s: %w", op, err)
		}
		roles = append(roles, *role)
	}
	return roles, nil
}
