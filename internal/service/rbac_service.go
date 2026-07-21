package service

import (
	"context"
	"errors"

	"github.com/dmi3midd/grpcsso/internal/domain"
	"github.com/dmi3midd/grpcsso/internal/repository"
)

var (
	ErrRoleNotFound       = errors.New("role not found")
	ErrPermissionNotFound = errors.New("permission not found")
)

type RBACService interface {
	// GetRoleById returns a role by its id.
	// Returns [ErrRoleNotFound] if the role is not found.
	GetRoleById(ctx context.Context, roleId string) (*domain.Role, error)
	// CreateRole creates a new role.
	CreateRole(ctx context.Context, name string) (string, error)
	// DeleteRole deletes a role by its id.s
	DeleteRole(ctx context.Context, roleId string) (string, error)

	// AssignRoleToUser assigns a role to a user within a transaction.
	// Returns role and user ids.
	// Returns [ErrUserNotFound] if the user does not exist.
	// Returns [ErrRoleNotFound] if the role does not exist.
	AssignRoleToUser(ctx context.Context, roleId, userId string) (string, string, error)
	// RevokeRoleFromUser removes a role from a user within a transaction.
	// Returns role and user ids.
	// Returns [ErrUserNotFound] if the user does not exist.
	// Returns [ErrRoleNotFound] if the role does not exist.
	RevokeRoleFromUser(ctx context.Context, roleId, userId string) (string, string, error)
	// GetUserRoles returns all roles assigned to a user.
	// Returns an empty slice if the user has no roles.
	GetUserRoles(ctx context.Context, userId string) ([]domain.Role, error)

	// GetPermissionById a permission by its id.
	// Returns ErrPermissionNotFound if the permission is not found.
	GetPermissionById(ctx context.Context, permissionId string) (*domain.Permission, error)
	// CreatePermission a new permission.
	CreatePermission(ctx context.Context, name string) (string, error)
	// DeletePermission a permission by its id.
	DeletePermission(ctx context.Context, permissionId string) (string, error)

	// AssignPermissionToRole assigns a permission to a role within a transaction.
	// Returns permission and role ids.
	// Returns [ErrRoleNotFound] if the role does not exist.
	// Returns [ErrPermissionNotFound] if the permission does not exist.
	AssignPermissionToRole(ctx context.Context, permissionId, roleId string) (string, string, error)
	// RevokePermissionFromRole removes a permission from a role within a transaction.
	// Returns permission and role ids.
	// Returns [ErrRoleNotFound] if the role does not exist.
	// Returns [ErrPermissionNotFound] if the permission does not exist.
	RevokePermissionFromRole(ctx context.Context, permissionId, roleId string) (string, string, error)
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
