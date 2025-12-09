package handler

import (
	"catatan-keuangan/config"
	"catatan-keuangan/modules/auth"
	"catatan-keuangan/modules/users"
	"fmt"
	"net/http"
	"net/url"
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

func (h *authHandler) GoogleLogin(c *gin.Context) {
	remember := c.Query("remember") == "true"
	state := uuid.NewString()
	// Simpan state di cookie agar bisa diverifikasi saat callback untuk mencegah CSRF.
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("oauth_state", state, int((10*time.Minute)/time.Second), "/", "", false, true)
	c.SetCookie("oauth_remember", fmt.Sprintf("%t", remember), int((10*time.Minute)/time.Second), "/", "", false, true)

	url := h.authService.GoogleLoginURL(state)
	c.JSON(http.StatusOK, gin.H{"url": url})
	// atau c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *authHandler) GoogleCallback(c *gin.Context) {
	redirectWithError := func(msg string) {
		if config.Cfg.FrontendRedirectURL != "" {
			if target, err := url.Parse(config.Cfg.FrontendRedirectURL); err == nil {
				q := target.Query()
				q.Set("error", msg)
				target.RawQuery = q.Encode()
				c.Redirect(http.StatusTemporaryRedirect, target.String())
				return
			}
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
	}
	rememberCookie, _ := c.Cookie("oauth_remember")
	remember := rememberCookie == "true"
	stateCookie, err := c.Cookie("oauth_state")
	stateQuery := c.Query("state")
	if err != nil || stateCookie == "" || stateCookie != stateQuery {
		c.SetCookie("oauth_state", "", -1, "/", "", false, true)
		redirectWithError("state tidak valid")
		return
	}

	code := c.Query("code")
	if code == "" {
		c.SetCookie("oauth_state", "", -1, "/", "", false, true)
		redirectWithError("code kosong")
		return
	}
	access, refresh, _, err := h.authService.GoogleCallback(c.Request.Context(), code, remember)
	if err != nil {
		c.SetCookie("oauth_state", "", -1, "/", "", false, true)
		redirectWithError(err.Error())
		return
	}
	lifespan := 24 * time.Hour
	if remember {
		lifespan = 7 * 24 * time.Hour
	}
	maxAge := int(lifespan / time.Second)
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("refresh_token", refresh, maxAge, "/", "", false, true)
	c.SetCookie("oauth_state", "", -1, "/", "", false, true)
	c.SetCookie("oauth_remember", "", -1, "/", "", false, true)

	if config.Cfg.FrontendRedirectURL != "" {
		if target, err := url.Parse(config.Cfg.FrontendRedirectURL); err == nil {
			c.Redirect(http.StatusTemporaryRedirect, target.String())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": access,
	})
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
	lifespan := 24 * time.Hour
	if loginRequest.RememberMe {
		lifespan = 7 * 24 * time.Hour
	}
	access, refresh, _, err := h.authService.Login(loginRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	maxAge := int(lifespan / time.Second)
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
	})
}

func (h *authHandler) RefreshToken(c *gin.Context) {
	rt, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token tidak ada"})
		return
	}
	access, newRefresh, expires, _, err := h.authService.Refresh(rt)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	lifespan := time.Until(expires)
	if lifespan < 0 {
		lifespan = 0
	}
	maxAge := int(lifespan / time.Second)
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
