package handler

import (
	"catatan-keuangan/modules/auth"
	"catatan-keuangan/modules/users"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type authHandler struct {
	authService auth.Service
}

func NewAuthHandler(authService auth.Service) *authHandler {
	return &authHandler{authService}
}

func (h *authHandler) Register(c *gin.Context) {
	var registerRequest users.UserRequest
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "input tidak valid",
		})
		return
	}
	user, err := h.authService.Register(registerRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"data": users.FormatUserResponse(user),
	})
}

func (h *authHandler) Login(c *gin.Context) {
	var loginRequest auth.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "input tidak valid",
		})
		return
	}
	access, refresh, _, err := h.authService.Login(loginRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	maxAge := int((7 * 24 * time.Hour) / time.Second)
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		"refresh_token",
		refresh,
		maxAge,
		"/",
		"",
		false,
		true,
	)
	c.JSON(http.StatusOK, gin.H{
		"access_token": access,
		// "user":         users.FormatUserResponse(user),
	})
}

func (h *authHandler) RefreshToken(c *gin.Context) {
	rt, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token tidak ada"})
		return
	}
	access, newRefresh, _, err := h.authService.Refresh(rt)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	maxAge := 7 * 24 * 60 * 60
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("refresh_token", newRefresh, maxAge, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"access_token": access,
	})
}

func (h *authHandler) Me(c *gin.Context) {
	userIDVal, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "userID tidak ditemukan di token"})
		return
	}
	userID, err := uuid.Parse(userIDVal.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "userID tidak valid"})
		return
	}

	user, err := h.authService.Me(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users.FormatUserResponse(user)})
}
func (h *authHandler) Logout(c *gin.Context) {
	rt, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh token tidak ada"})
		return
	}
	if err := h.authService.Logout(rt); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token tidak valid"})
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "logout berhasil"})
}
