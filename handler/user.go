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

func (h *userHandler) GetAllUsers(c *gin.Context) {
	userData, err := h.userService.FindAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	usersResponse := users.FormatUserResponses(userData)
	c.JSON(http.StatusOK, gin.H{
		"data": usersResponse,
	})
}

func (h *userHandler) GetUserByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id tidak valid"})
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

func (h *userHandler) PostUserHandler(c *gin.Context) {
	var userRequest users.UserRequest
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": formatValidationError(err),
		})
		return
	}
	newUser, err := h.userService.Create(userRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	userResponse := users.FormatUserResponse(newUser)
	c.JSON(http.StatusCreated, gin.H{
		"data": userResponse,
	})
}

func (h *userHandler) DeleteUserHandler(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id tidak valid"})
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

func (h *userHandler) UpdateUserHandler(c *gin.Context) {
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
