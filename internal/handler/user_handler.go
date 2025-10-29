package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vayura/internal/service"
	"github.com/vayura/pkg"
)

type UserHandler struct {
	userService    service.UserService
	storageService service.StorageService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService service.UserService, storageService service.StorageService) *UserHandler {
	return &UserHandler{
		userService:    userService,
		storageService: storageService,
	}
}

// GetProfile returns the authenticated user's profile
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := pkg.GetUserID(c)
	if !exists {
		pkg.JSONUnauthorized(c, pkg.ErrInvalidToken)
		return
	}

	user, err := h.userService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		pkg.JSONInternalServerError(c, err)
		return
	}

	pkg.JSONSuccess(c, http.StatusOK, "profile fetched successfully", user)
}

// UpdateProfile updates user profile fields
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := pkg.GetUserID(c)
	if !exists {
		pkg.JSONUnauthorized(c, pkg.ErrInvalidToken)
		return
	}

	var req service.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.JSONBadRequest(c, err)
		return
	}

	user, err := h.userService.UpdateProfile(c.Request.Context(), userID, req)
	if err != nil {
		if err == pkg.ErrUsernameExists {
			pkg.JSONBadRequest(c, err)
			return
		}
		pkg.JSONInternalServerError(c, err)
		return
	}

	pkg.JSONSuccess(c, http.StatusOK, "profile updated successfully", user)
}

// DeleteProfile deletes the authenticated user's profile
func (h *UserHandler) DeleteProfile(c *gin.Context) {
	userID, exists := pkg.GetUserID(c)
	if !exists {
		pkg.JSONUnauthorized(c, pkg.ErrInvalidToken)
		return
	}

	if err := h.userService.DeleteProfile(c.Request.Context(), userID); err != nil {
		pkg.JSONInternalServerError(c, err)
		return
	}

	pkg.JSONSuccess(c, http.StatusOK, "profile deleted successfully", nil)
}

// UploadAvatar handles avatar upload
func (h *UserHandler) UploadAvatar(c *gin.Context) {
	userID, exists := pkg.GetUserID(c)
	if !exists {
		pkg.JSONUnauthorized(c, pkg.ErrInvalidToken)
		return
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		pkg.JSONBadRequest(c, err)
		return
	}

	avatarPath, err := h.storageService.SaveAvatar(c.Request.Context(), userID, file)
	if err != nil {
		pkg.JSONBadRequest(c, err)
		return
	}

	user, err := h.userService.UpdateAvatar(c.Request.Context(), userID, avatarPath)
	if err != nil {
		pkg.JSONInternalServerError(c, err)
		return
	}

	pkg.JSONSuccess(c, http.StatusOK, "avatar updated successfully", user)
}
