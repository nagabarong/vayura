package service

import (
	"context"
	"regexp"
	"time"

	"github.com/vayura/internal/models"
	"github.com/vayura/internal/repository"
	"github.com/vayura/pkg"
)

// authService implements AuthService interface
type authService struct {
	userRepo repository.UserRepository
}

// NewAuthService creates a new authentication service
func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

// RegisterRequest represents the registration request
type RegisterRequest struct {
	FullName string `json:"full_name" binding:"required,min=3"`
	Username string `json:"username" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
	Gender   string `json:"gender"`
	Birthday string `json:"birthday"` // format YYYY-MM-DD
}

// LoginRequest represents the login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (s *authService) Register(ctx context.Context, req RegisterRequest) (*models.User, error) {
	// Validation
	if len(req.FullName) < 3 {
		return nil, &pkg.ValidationError{Field: "full_name", Message: "full name must be at least 3 characters"}
	}
	if len(req.Username) < 3 {
		return nil, &pkg.ValidationError{Field: "username", Message: "username must be at least 3 characters"}
	}
	if !isValidEmail(req.Email) {
		return nil, &pkg.ValidationError{Field: "email", Message: "invalid email format"}
	}
	if len(req.Password) < 8 {
		return nil, &pkg.ValidationError{Field: "password", Message: "password must be at least 8 characters"}
	}

	// Check if email already exists
	emailExists, err := s.userRepo.EmailExists(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if emailExists {
		return nil, pkg.ErrEmailExists
	}

	// Check if username already exists
	usernameExists, err := s.userRepo.UsernameExists(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if usernameExists {
		return nil, pkg.ErrUsernameExists
	}

	// Parse birthday if provided
	var birth time.Time
	if req.Birthday != "" {
		birth, err = time.Parse("2006-01-02", req.Birthday)
		if err != nil {
			return nil, &pkg.ValidationError{Field: "birthday", Message: "invalid birthday format, use YYYY-MM-DD"}
		}
	}

	// Create user
	user := &models.User{
		FullName: req.FullName,
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
		Role:     req.Role,
		Gender:   req.Gender,
		Birthday: birth,
	}

	// Hash password
	if err := user.HashPassword(req.Password); err != nil {
		return nil, err
	}

	// Save user
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(ctx context.Context, req LoginRequest) (*models.User, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, pkg.ErrInvalidCredentials
	}

	if !user.CheckPassword(req.Password) {
		return nil, pkg.ErrInvalidCredentials
	}

	return user, nil
}

func isValidEmail(email string) bool {
	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return regex.MatchString(email)
}
