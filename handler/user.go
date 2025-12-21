package handler

import (
	"catatan-keuangan/modules/users"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userHandler struct {
	userService users.Service
}

func NewUserHandler(userService users.Service) *userHandler {
	return &userHandler{userService}
}

// GetAllUsers godoc
// @Summary Daftar semua user
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} UsersDataResponse
// @Failure 400 {object} ErrorResponse
// @Router /users [get]
func (h *userHandler) GetAllUsers(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		return
	}
	userData, err := h.userService.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	usersResponse := []users.UserResponse{users.FormatUserResponse(userData)}
	c.JSON(http.StatusOK, gin.H{
		"data": usersResponse,
	})
}

// GetUserByID godoc
// @Summary Detail user by ID
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} UserDataResponse
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /users/{id} [get]
func (h *userHandler) GetUserByID(c *gin.Context) {
	requesterID, ok := userIDFromContext(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id tidak valid"})
		return
	}
	if id != requesterID {
		c.JSON(http.StatusForbidden, gin.H{"error": "akses tidak diizinkan"})
		return
	}
	userData, err := h.userService.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Gagal mengambil data user",
		})
		return
	}
	userResponse := users.FormatUserResponse(userData)
	c.JSON(http.StatusOK, gin.H{
		"data": userResponse,
	})
}

// PostUserHandler godoc
// @Summary Buat user baru
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body users.UserRequest true "Data user"
// @Success 201 {object} UserDataResponse
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /users [post]
func (h *userHandler) PostUserHandler(c *gin.Context) {
	c.JSON(http.StatusForbidden, gin.H{
		"error": "akses tidak diizinkan",
	})
}

// DeleteUserHandler godoc
// @Summary Hapus user
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /users/{id} [delete]
func (h *userHandler) DeleteUserHandler(c *gin.Context) {
	requesterID, ok := userIDFromContext(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id tidak valid"})
		return
	}
	if id != requesterID {
		c.JSON(http.StatusForbidden, gin.H{"error": "akses tidak diizinkan"})
		return
	}
	_, err = h.userService.Delete(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Gagal menghapus user",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "User berhasil di hapus",
	})
}

// UpdateUserHandler godoc
// @Summary Perbarui user
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body users.UserRequest true "Data user"
// @Success 200 {object} UserDataResponse
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /users/{id} [put]
func (h *userHandler) UpdateUserHandler(c *gin.Context) {
	requesterID, ok := userIDFromContext(c)
	if !ok {
		return
	}
	var userRequest users.UserRequest
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": formatValidationError(err),
		})
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id tidak valid"})
		return
	}
	if id != requesterID {
		c.JSON(http.StatusForbidden, gin.H{"error": "akses tidak diizinkan"})
		return
	}
	updatedUser, err := h.userService.Update(id, userRequest)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Gagal memperbarui user",
		})
		return
	}
	userResponse := users.FormatUserResponse(updatedUser)
	c.JSON(http.StatusOK, gin.H{
		"data": userResponse,
	})
}

// ChangePassword godoc
// @Summary Ganti password user
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body users.ChangePasswordRequest true "Data password"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Router /users/password [put]
func (h *userHandler) ChangePassword(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		return
	}
	var req users.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "input tidak valid"})
		return
	}
	if err := h.userService.ChangePassword(userID, req); err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, gorm.ErrRecordNotFound) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "password berhasil diubah"})
}

func formatValidationError(err error) []string {
	var errors []string

	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range ve {
			var msg string
			switch fe.Tag() {
			case "required":
				msg = fmt.Sprintf("%s harus di isi", fe.Field())
			case "email":
				msg = fmt.Sprintf("%s harus valid email", fe.Field())
			case "min":
				msg = fmt.Sprintf("%s harus %s karakter", fe.Field(), fe.Param())
			default:
				msg = fmt.Sprintf("%s tidak valid", fe.Field())
			}
			errors = append(errors, msg)
		}
		return errors
	}

	return []string{err.Error()}
}
