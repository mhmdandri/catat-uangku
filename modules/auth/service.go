package auth

import (
	"catatan-keuangan/config"
	userprofile "catatan-keuangan/modules/user_profile"
	"catatan-keuangan/modules/users"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
)

type Service interface {
	Login(loginRequest LoginRequest) (string, string, users.User, error)
	Refresh(raw string) (string, string, time.Time, users.User, error)
	Me(userID uuid.UUID) (users.User, error)
	Logout(raw string) error
	Register(registerRequest users.UserRequest) (users.User, error)
	GoogleLoginURL(state string) string
	GoogleCallback(ctx context.Context, code string, remember bool) (string, string, users.User, error)
}

type service struct {
	userRepository    users.Repository
	refreshRepository RefreshRepository
	profileService    userprofile.Service
	googleConfig      *oauth2.Config
}
type googleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

func NewService(userRepository users.Repository, refreshRepository RefreshRepository, profileService userprofile.Service) *service {
	return &service{userRepository: userRepository, refreshRepository: refreshRepository, profileService: profileService, googleConfig: &oauth2.Config{
		ClientID:     config.Cfg.GoogleClientID,
		ClientSecret: config.Cfg.GoogleClientSecret,
		RedirectURL:  config.Cfg.GoogleRedirectURL,
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     google.Endpoint,
	}}
}
func (s *service) GoogleLoginURL(state string) string {
	if state == "" {
		state = uuid.New().String()
	}
	return s.googleConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}
func (s *service) GoogleCallback(ctx context.Context, code string, remember bool) (string, string, users.User, error) {
	tok, err := s.googleConfig.Exchange(ctx, code)
	if err != nil {
		return "", "", users.User{}, errors.New("gagal tukar code ke token")
	}

	client := s.googleConfig.Client(ctx, tok)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil || resp.StatusCode != http.StatusOK {
		return "", "", users.User{}, errors.New("gagal ambil data user google")
	}
	defer resp.Body.Close()

	var gu googleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&gu); err != nil || gu.Email == "" {
		return "", "", users.User{}, errors.New("data user google tidak valid")
	}
	if !gu.VerifiedEmail {
		return "", "", users.User{}, errors.New("email google belum terverifikasi")
	}

	user, err := s.userRepository.FindByEmail(gu.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			hashed, hashErr := bcrypt.GenerateFromPassword([]byte(gu.ID), bcrypt.DefaultCost) // placeholder pwd
			if hashErr != nil {
				return "", "", users.User{}, errors.New("gagal membuat akun baru")
			}
			user, err = s.userRepository.Create(users.User{
				Name:     gu.Name,
				Email:    gu.Email,
				Password: string(hashed),
			})
			if err != nil {
				return "", "", users.User{}, err
			}
		} else {
			return "", "", users.User{}, err
		}
	}
	if _, err := s.profileService.EnsureDefaultProfile(user.ID, user.Name, user.Email); err != nil {
		return "", "", users.User{}, errors.New("gagal membuat profil user")
	}

	access, err := GenerateToken(TokenData{UserID: user.ID, Role: ""})
	if err != nil {
		return "", "", users.User{}, errors.New("gagal membuat access token")
	}

	rawRefresh, hashRefresh, err := generateRefresh()
	if err != nil {
		return "", "", users.User{}, errors.New("gagal membuat refresh token")
	}
	lifespan := 24 * time.Hour
	if remember {
		lifespan = 7 * 24 * time.Hour
	}
	expires := time.Now().Add(lifespan)
	if _, err := s.refreshRepository.Create(RefreshToken{UserID: user.ID, Token: hashRefresh, ExpiresAt: expires}); err != nil {
		return "", "", users.User{}, err
	}
	return access, rawRefresh, user, nil
}

func generateRefresh() (string, string, error) {
	base := make([]byte, 32)
	if _, err := rand.Read(base); err != nil {
		return "", "", err
	}
	raw := base64.StdEncoding.EncodeToString(base)
	sum := sha256.Sum256([]byte(raw))
	return raw, base64.StdEncoding.EncodeToString(sum[:]), nil
}

func (s *service) Login(loginRequest LoginRequest) (string, string, users.User, error) {
	user, err := s.userRepository.FindByEmail(loginRequest.Email)
	if err != nil {
		return "", "", users.User{}, errors.New("email atau password salah")
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)) != nil {
		return "", "", users.User{}, errors.New("email atau password salah")
	}
	access, err := GenerateToken(TokenData{UserID: user.ID, Role: ""})
	if err != nil {
		return "", "", users.User{}, errors.New("gagal membuat access token")
	}
	rawRefresh, hashRefresh, err := generateRefresh()
	if err != nil {
		return "", "", users.User{}, errors.New("gagal membuat refresh token")
	}
	expires := time.Now().Add(24 * time.Hour)
	if loginRequest.RememberMe {
		expires = time.Now().Add(7 * 24 * time.Hour)
	}
	if _, err := s.refreshRepository.Create(RefreshToken{UserID: user.ID, Token: hashRefresh, ExpiresAt: expires}); err != nil {
		return "", "", users.User{}, err
	}
	return access, rawRefresh, user, nil
}

func hashRaw(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return base64.StdEncoding.EncodeToString(sum[:])
}

func (s *service) Refresh(raw string) (string, string, time.Time, users.User, error) {
	hash := hashRaw(raw)
	rt, err := s.refreshRepository.FindValidByHash(hash)
	if err != nil {
		return "", "", time.Time{}, users.User{}, errors.New("refresh token invalid")
	}

	user, err := s.userRepository.FindByID(rt.UserID)
	if err != nil {
		return "", "", time.Time{}, users.User{}, err
	}

	if err := s.refreshRepository.Revoke(rt.ID); err != nil {
		return "", "", time.Time{}, users.User{}, err
	}

	access, err := GenerateToken(TokenData{UserID: user.ID, Role: ""})
	if err != nil {
		return "", "", time.Time{}, users.User{}, errors.New("gagal membuat access token")
	}
	newRaw, newHash, err := generateRefresh()
	if err != nil {
		return "", "", time.Time{}, users.User{}, errors.New("gagal membuat refresh token")
	}
	lifespan := rt.ExpiresAt.Sub(rt.CreatedAt)
	remembered := lifespan > 48*time.Hour
	if remembered {
		lifespan = 7 * 24 * time.Hour
	} else {
		lifespan = 24 * time.Hour
	}
	expires := time.Now().Add(lifespan)
	if _, err := s.refreshRepository.Create(RefreshToken{UserID: user.ID, Token: newHash, ExpiresAt: expires, RotatedFrom: &rt.ID}); err != nil {
		return "", "", time.Time{}, users.User{}, err
	}
	return access, newRaw, expires, user, nil
}

func (s *service) Me(userID uuid.UUID) (users.User, error) {
	return s.userRepository.FindByID(userID)
}
func (s *service) Logout(raw string) error {
	hash := hashRaw(raw)
	rt, err := s.refreshRepository.FindValidByHash(hash)
	if err != nil {
		return err
	}
	return s.refreshRepository.Revoke(rt.ID)
}

func (s *service) Register(registerRequest users.UserRequest) (users.User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		return users.User{}, err
	}
	checkUser, _ := s.userRepository.FindByEmail(registerRequest.Email)
	if checkUser.ID != uuid.Nil {
		return users.User{}, errors.New("email sudah terdaftar")
	}
	user := users.User{
		Name:     registerRequest.Name,
		Email:    registerRequest.Email,
		Password: string(hashed),
	}
	newUser, err := s.userRepository.Create(user)
	if err != nil {
		return users.User{}, err
	}
	if profile, err := s.profileService.EnsureDefaultProfile(newUser.ID, newUser.Name, newUser.Email); err == nil {
		newUser.Profile = profile
	} else {
		return users.User{}, errors.New("gagal membuat profil user")
	}
	return newUser, err
}
