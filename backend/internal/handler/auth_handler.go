package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vayura/backend/internal/service"
	"github.com/vayura/backend/pkg"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
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

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.JSONBadRequest(c, err)
		return
	}

	// Convert to service request
	serviceReq := service.RegisterRequest{
		FullName: req.FullName,
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Phone:    req.Phone,
		Role:     req.Role,
		Gender:   req.Gender,
		Birthday: req.Birthday,
	}

	user, err := h.authService.Register(c.Request.Context(), serviceReq)
	if err != nil {
		if err == pkg.ErrEmailExists || err == pkg.ErrUsernameExists {
			pkg.JSONBadRequest(c, err)
		} else {
			pkg.JSONInternalServerError(c, err)
		}
		return
	}

	pkg.JSONSuccess(c, http.StatusCreated, "user registered successfully", user)
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.JSONBadRequest(c, err)
		return
	}

	serviceReq := service.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	user, err := h.authService.Login(c.Request.Context(), serviceReq)
	if err != nil {
		pkg.JSONUnauthorized(c, err)
		return
	}

	// Generate JWT
	token, err := pkg.GenerateJWT(user.ID, user.Email)
	if err != nil {
		pkg.JSONInternalServerError(c, err)
		return
	}

	// Return response with user data (excluding password)
	response := gin.H{
		"token": token,
		"user": gin.H{
			"id":        user.ID,
			"full_name": user.FullName,
			"username":  user.Username,
			"email":     user.Email,
			"phone":     user.Phone,
			"role":      user.Role,
			"gender":    user.Gender,
			"birthday":  user.Birthday,
		},
	}

	pkg.JSONSuccess(c, http.StatusOK, "login successful", response)
}
