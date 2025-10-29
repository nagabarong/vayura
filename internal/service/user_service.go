package service

import (
	"context"
	"time"

	"github.com/vayura/internal/models"
	"github.com/vayura/internal/repository"
	"github.com/vayura/pkg"
)

// userService implements UserService interface
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

// UpdateProfileRequest represents the update profile request
type UpdateProfileRequest struct {
	FullName string `json:"full_name"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Gender   string `json:"gender"`
	Birthday string `json:"birthday"`
}

func (s *userService) GetProfile(ctx context.Context, userID uint) (*models.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, pkg.ErrUserNotFound
	}
	return user, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userID uint, req UpdateProfileRequest) (*models.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, pkg.ErrUserNotFound
	}

	// Update fields if provided
	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Username != "" {
		// Check if username already exists
		if req.Username != user.Username {
			usernameExists, err := s.userRepo.UsernameExists(ctx, req.Username)
			if err != nil {
				return nil, err
			}
			if usernameExists {
				return nil, pkg.ErrUsernameExists
			}
		}
		user.Username = req.Username
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Gender != "" {
		user.Gender = req.Gender
	}
	if req.Birthday != "" {
		birth, err := time.Parse("2006-01-02", req.Birthday)
		if err != nil {
			return nil, &pkg.ValidationError{Field: "birthday", Message: "invalid birthday format, use YYYY-MM-DD"}
		}
		user.Birthday = birth
	}

	// Save updated user
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) DeleteProfile(ctx context.Context, userID uint) error {
	return s.userRepo.Delete(ctx, userID)
}

func (s *userService) UpdateAvatar(ctx context.Context, userID uint, avatarPath string) (*models.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, pkg.ErrUserNotFound
	}

	user.Avatar = avatarPath
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
