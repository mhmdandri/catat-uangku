package handler

import (
	"catatan-keuangan/modules/user_profile"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userProfileHandler struct {
	profileService userprofile.Service
}

func NewUserProfileHandler(profileService userprofile.Service) *userProfileHandler {
	return &userProfileHandler{profileService}
}

// GetProfile godoc
// @Summary Ambil profil user
// @Tags Profile
// @Security BearerAuth
// @Produce json
// @Success 200 {object} ProfileDataResponse
// @Failure 401 {object} ErrorResponse
// @Router /profile [get]
func (h *userProfileHandler) GetProfile(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		return
	}
	profile, err := h.profileService.GetProfileByUserID(userID)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, gorm.ErrRecordNotFound) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": userprofile.FormatUserProfileResponse(profile),
	})
}

// UpdateProfile godoc
// @Summary Perbarui profil user
// @Tags Profile
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body userprofile.UpdateProfileRequest true "Data profil"
// @Success 200 {object} ProfileDataResponse
// @Failure 400 {object} ErrorResponse
// @Router /profile [put]
func (h *userProfileHandler) UpdateProfile(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		return
	}
	var req userprofile.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "input tidak valid"})
		return
	}
	profile, err := h.profileService.Update(userID, req)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, gorm.ErrRecordNotFound) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": userprofile.FormatUserProfileResponse(profile),
	})
}

// UploadAvatar godoc
// @Summary Upload avatar
// @Tags Profile
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param avatar formData file true "File avatar"
// @Success 200 {object} ProfileDataResponse
// @Failure 400 {object} ErrorResponse
// @Router /profile/avatar [post]
func (h *userProfileHandler) UploadAvatar(c *gin.Context) {
	userID, ok := userIDFromContext(c)
	if !ok {
		return
	}
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file avatar diperlukan"})
		return
	}
	if file.Size == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file avatar kosong"})
		return
	}
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s_%d%s", userID.String(), time.Now().UnixNano(), ext)
	saveDir := filepath.Join("uploads", "avatars")
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menyiapkan folder upload"})
		return
	}
	dstPath := filepath.Join(saveDir, filename)
	if err := c.SaveUploadedFile(file, dstPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menyimpan file"})
		return
	}
	avatarURL := fmt.Sprintf("/%s", filepath.ToSlash(dstPath))
	profile, err := h.profileService.UpdateAvatar(userID, avatarURL)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, gorm.ErrRecordNotFound) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": userprofile.FormatUserProfileResponse(profile),
	})
}

func userIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	userIDVal, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "userID tidak ditemukan di token"})
		return uuid.Nil, false
	}
	userID, err := uuid.Parse(userIDVal.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "userID tidak valid"})
		return uuid.Nil, false
	}
	return userID, true
}
