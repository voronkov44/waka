package auth

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"rest_waka/pkg/jwtx"
	"strconv"
	"strings"
	"time"
)

type Service struct {
	repo          RepositoryGorm
	jwtSecret     []byte
	tokenTTL      time.Duration
	adminName     string
	adminPassword string
}

func NewService(repo RepositoryGorm, jwtSecret string, tokenTTL time.Duration, adminName, adminPassword string) (*Service, error) {
	if jwtSecret == "" {
		return nil, fmt.Errorf("empty jwt secret")
	}
	if tokenTTL <= 0 {
		return nil, fmt.Errorf("token ttl must be positive")
	}
	if strings.TrimSpace(adminName) == "" {
		return nil, fmt.Errorf("empty admin name")
	}
	if adminPassword == "" {
		return nil, fmt.Errorf("empty admin password")
	}

	return &Service{
		repo:          repo,
		jwtSecret:     []byte(jwtSecret),
		tokenTTL:      tokenTTL,
		adminName:     strings.TrimSpace(adminName),
		adminPassword: adminPassword,
	}, nil
}

// LoginTelegram - login or registry user telegram
func (s *Service) LoginTelegram(ctx context.Context, tg TelegramProfile) (string, error) {
	if tg.TgID <= 0 {
		return "", ErrInvalidArgument
	}

	u, err := s.repo.UpsertTelegram(ctx, tg)
	if err != nil {
		return "", err
	}

	token, err := s.generateUserToken(u.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) LoginAdmin(_ context.Context, req AdminLoginRequest) (string, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" || req.Password == "" {
		return "", ErrInvalidArgument
	}

	if !secureEqual(name, s.adminName) || !secureEqual(req.Password, s.adminPassword) {
		return "", ErrUnauthorized
	}

	token, err := s.generateAdminToken(name)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) Me(ctx context.Context, userID uint64) (MeResponse, error) {
	if userID == 0 {
		return MeResponse{}, ErrInvalidArgument
	}

	u, err := s.repo.Get(ctx, userID)
	if err != nil {
		return MeResponse{}, err
	}

	return MeResponse{
		ID:        u.ID,
		TgID:      u.TgID,
		Username:  u.Username,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		PhotoURL:  u.PhotoURL,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}, nil
}

func (s *Service) generateUserToken(userID uint64) (string, error) {
	now := time.Now()

	if userID <= 0 {
		return "", fmt.Errorf("bad user id: %d", userID)
	}

	claims := jwtx.Claims{
		Role:   jwtx.RoleUser,
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatUint(userID, 10),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.tokenTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *Service) generateAdminToken(name string) (string, error) {
	now := time.Now()

	if name == "" {
		return "", fmt.Errorf("empty admin name")
	}

	claims := jwtx.Claims{
		Role: jwtx.RoleAdmin,
		Name: name,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   name,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.tokenTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func secureEqual(a, b string) bool {
	ha := sha256.Sum256([]byte(a))
	hb := sha256.Sum256([]byte(b))
	return subtle.ConstantTimeCompare(ha[:], hb[:]) == 1
}
