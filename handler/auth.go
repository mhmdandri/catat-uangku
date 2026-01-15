package handler

import (
	"catatan-keuangan/config"
	"catatan-keuangan/modules/auth"
	"catatan-keuangan/modules/users"
	"net/http"
	"net/url"
	"strconv"
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

// GoogleLogin godoc
// @Summary Dapatkan URL login Google
// @Description Mengembalikan URL Google OAuth dan menyimpan state/remember di cookie.
// @Tags Auth
// @Produce json
// @Param remember query bool false "Set true untuk sesi lebih lama"
// @Success 200 {object} URLResponse
// @Failure 400 {object} ErrorResponse
// @Router /auth/google/login [get]
func (h *authHandler) GoogleLogin(c *gin.Context) {
	remember := c.Query("remember") == "true"
	state := uuid.NewString()
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		"oauth_state",
		state,
		int((10*time.Minute)/time.Second),
		"/",
		"",
		cookieSecure(),
		true,
	)
	c.SetCookie(
		"oauth_remember",
		strconv.FormatBool(remember),
		int((10*time.Minute)/time.Second),
		"/",
		"",
		cookieSecure(),
		true,
	)

	url := h.authService.GoogleLoginURL(state)
	c.JSON(http.StatusOK, gin.H{"url": url})
}

// GoogleCallback godoc
// @Summary Callback Google OAuth
// @Description Tukar code Google menjadi access token dan refresh token cookie.
// @Tags Auth
// @Produce json
// @Success 200 {object} AccessTokenResponse
// @Failure 400 {object} ErrorResponse
// @Router /auth/google/callback [get]
func (h *authHandler) GoogleCallback(c *gin.Context) {
	redirectWithError := func(msg string) {
		c.SetCookie("oauth_state", "", -1, "/", "", cookieSecure(), true)
		c.SetCookie("oauth_remember", "", -1, "/", "", cookieSecure(), true)

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
		redirectWithError("state tidak valid")
		return
	}

	code := c.Query("code")
	if code == "" {
		redirectWithError("code kosong")
		return
	}

	access, refresh, _, err := h.authService.GoogleCallback(c.Request.Context(), code, remember)
	if err != nil {
		redirectWithError(err.Error())
		return
	}

	lifespan := 4 * time.Hour
	if remember {
		lifespan = 7 * 24 * time.Hour
	}
	maxAge := int(lifespan / time.Second)
	setAuthCookie(c, "access_token", access, 15*60)
	setAuthCookie(c, "refresh_token", refresh, maxAge)
	c.SetCookie("oauth_remember", "", -1, "/", "", cookieSecure(), true)

	if config.Cfg.FrontendRedirectURL != "" {
		if target, err := url.Parse(config.Cfg.FrontendRedirectURL); err == nil {
			c.Redirect(http.StatusTemporaryRedirect, target.String())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"access_token": access})
}

// Register godoc
// @Summary Daftar user baru
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body users.UserRequest true "Data registrasi"
// @Success 201 {object} UserDataResponse
// @Failure 400 {object} ErrorResponse
// @Router /auth/register [post]
func (h *authHandler) Register(c *gin.Context) {
	var registerRequest users.UserRequest
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "input tidak valid"})
		return
	}

	user, err := h.authService.Register(registerRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": users.FormatUserResponse(user)})
}

// Login godoc
// @Summary Login user
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body auth.LoginRequest true "Data login"
// @Success 200 {object} AccessTokenResponse
// @Failure 400 {object} ErrorResponse
// @Router /auth/login [post]
func (h *authHandler) Login(c *gin.Context) {
	var loginRequest auth.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "input tidak valid"})
		return
	}

	lifespan := 4 * time.Hour
	if loginRequest.RememberMe {
		lifespan = 7 * 24 * time.Hour
	}

	access, refresh, _, err := h.authService.Login(loginRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	maxAge := int(lifespan / time.Second)
	setAuthCookie(c, "refresh_token", refresh, maxAge)

	c.JSON(http.StatusOK, gin.H{"access_token": access})
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Menggunakan cookie refresh_token untuk mendapatkan access token baru.
// @Tags Auth
// @Produce json
// @Success 200 {object} AccessTokenResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/refresh [post]
func (h *authHandler) RefreshToken(c *gin.Context) {
	rt, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token tidak ada"})
		return
	}
	decoded, decodeErr := url.QueryUnescape(rt)
	if decodeErr == nil {
		rt = decoded
	}
	access, _, _, _, err := h.authService.Refresh(rt)
	if err != nil {
		setAuthCookie(c, "access_token", "", -1)
		setAuthCookie(c, "refresh_token", "", -1)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": access})
}

// Me godoc
// @Summary Profil user dari token
// @Tags Auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} AuthMeResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/me [get]
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

	c.JSON(http.StatusOK, user)
}

// Logout godoc
// @Summary Logout dan hapus refresh token
// @Tags Auth
// @Produce json
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Router /auth/logout [post]
func (h *authHandler) Logout(c *gin.Context) {
	rt, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh token tidak ada"})
		return
	}
	decoded, decodeErr := url.QueryUnescape(rt)
	if decodeErr == nil {
		rt = decoded
	}
	if err := h.authService.Logout(rt); err != nil {
		setAuthCookie(c, "refresh_token", "", -1)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token tidak valid"})
		return
	}
	setAuthCookie(c, "refresh_token", "", -1)
	c.JSON(http.StatusOK, gin.H{"message": "logout berhasil"})
}
