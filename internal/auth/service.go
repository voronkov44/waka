package auth

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"math"
	"strconv"
	"time"
)

// JWTClaims - набор данных, который мы кладём в токен
type JWTClaims struct {
	UserID uint32 `json:"user_id"`
	jwt.RegisteredClaims
}

type Service struct {
	repo      RepositoryGorm
	jwtSecret []byte
	tokenTTL  time.Duration
}

func NewService(repo RepositoryGorm, jwtSecret string, tokenTTL time.Duration) (*Service, error) {
	if jwtSecret == "" {
		return nil, fmt.Errorf("empty jwt secret")
	}
	if tokenTTL <= 0 {
		return nil, fmt.Errorf("token ttl must be positive")
	}
	return &Service{
		repo:      repo,
		jwtSecret: []byte(jwtSecret),
		tokenTTL:  tokenTTL,
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

	token, err := s.generateToken(u.ID)
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

func (s *Service) generateToken(userID uint64) (string, error) {
	now := time.Now()

	if userID <= 0 {
		return "", fmt.Errorf("bad user id: %d", userID)
	}
	if userID > math.MaxUint32 {
		return "", fmt.Errorf("user id too large for uint32 claim: %d", userID)
	}

	uid := uint32(userID)

	claims := JWTClaims{
		UserID: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatUint(uint64(uid), 10),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.tokenTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
