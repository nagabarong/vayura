package pkg

import "errors"

// Custom error types for better error handling
var (
	ErrEmailExists        = errors.New("email already registered")
	ErrUsernameExists     = errors.New("username already taken")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrMissingAuth        = errors.New("missing authorization header")
	ErrValidation         = errors.New("validation error")
)

// ValidationError represents a validation error with fields
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
