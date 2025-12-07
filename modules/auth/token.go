package auth

import (
	"catatan-keuangan/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenData struct {
	UserID uuid.UUID
	Role   string
}

func GenerateToken(data TokenData) (string, error) {
	claims := jwt.MapClaims{
		"iss":  config.Cfg.JWTIssuer,
		"aud":  config.Cfg.JWTAudience,
		"sub":  data.UserID.String(),
		"role": data.Role,
		"exp":  time.Now().Add(15 * time.Minute).Unix(),
		"iat":  time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Cfg.JWTSecret))
}
