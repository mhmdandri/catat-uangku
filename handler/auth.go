package handler

import (
	"catatan-keuangan/modules/auth"
	"catatan-keuangan/modules/users"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type authHandler struct {
	authService auth.Service
}

func NewAuthHandler(authService auth.Service) *authHandler {
	return &authHandler{authService}
}

func (h *authHandler) Login(c *gin.Context) {
	var loginRequest auth.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "input tidak valid",
		})
		return
	}
	access, refresh, user, err := h.authService.Login(loginRequest)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
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
		// "refresh_token": refresh,
		"user": users.FormatUserResponse(user),
	})
}

func (h *authHandler) RefreshToken(c *gin.Context) {
	rt, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token tidak ada"})
		return
	}
	access, newRefresh, user, err := h.authService.Refresh(rt)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	maxAge := 7 * 24 * 60 * 60
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("refresh_token", newRefresh, maxAge, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"access_token": access,
		// "refresh_token": refresh,
		"user": users.FormatUserResponse(user),
	})
}
