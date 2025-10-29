package service

import (
	"context"

	"github.com/vayura/backend/internal/models"
)

// AuthService defines the interface for authentication operations
type AuthService interface {
	Register(ctx context.Context, req RegisterRequest) (*models.User, error)
	Login(ctx context.Context, req LoginRequest) (*models.User, error)
}

// UserService defines the interface for user operations
type UserService interface {
	GetProfile(ctx context.Context, userID uint) (*models.User, error)
	UpdateProfile(ctx context.Context, userID uint, req UpdateProfileRequest) (*models.User, error)
	DeleteProfile(ctx context.Context, userID uint) error
	UpdateAvatar(ctx context.Context, userID uint, avatarPath string) (*models.User, error)
}
