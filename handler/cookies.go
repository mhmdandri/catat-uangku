package handler

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func isProd() bool {
	return os.Getenv("APP_ENV") == "production"
}

func cookieDomain() string {
	if !isProd() {
		return ""
	}
	return ".mohaproject.dev"
}

func cookieSecure() bool {
	return isProd()
}

func cookieSameSite() http.SameSite {
	if isProd() {
		return http.SameSiteNoneMode
	}
	return http.SameSiteLaxMode
}

func setAuthCookie(c *gin.Context, name, value string, maxAge int) {
	c.SetSameSite(cookieSameSite())
	c.SetCookie(
		name,
		value,
		maxAge,
		"/",
		cookieDomain(),
		cookieSecure(),
		true,
	)
}
