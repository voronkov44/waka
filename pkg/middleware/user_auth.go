package middleware

import (
	"context"
	"fmt"
	"net/http"
	"rest_waka/pkg/jwtx"
	"rest_waka/pkg/res"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type ctxKey string

const (
	userIDKey    ctxKey = "user_id"
	adminNameKey ctxKey = "admin_name"
)

func RequireUser(next http.Handler, jwtSecret string) http.Handler {
	return requireRole(next, jwtSecret, jwtx.RoleUser)
}

func RequireAdmin(next http.Handler, jwtSecret string) http.Handler {
	return requireRole(next, jwtSecret, jwtx.RoleAdmin)
}

func requireRole(next http.Handler, jwtSecret, wantRole string) http.Handler {
	secret := []byte(jwtSecret)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr, ok := readToken(r)
		if !ok {
			res.Json(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		claims := &jwtx.Claims{}
		t, err := jwt.ParseWithClaims(
			tokenStr,
			claims,
			func(token *jwt.Token) (interface{}, error) {
				if token.Method != jwt.SigningMethodHS256 {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return secret, nil
			},
			jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		)

		if err != nil || !t.Valid || claims.Role != wantRole {
			res.Json(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()

		switch wantRole {
		case jwtx.RoleUser:
			if claims.UserID == 0 {
				res.Json(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			ctx = context.WithValue(ctx, userIDKey, claims.UserID)

		case jwtx.RoleAdmin:
			if claims.Name == "" {
				res.Json(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			ctx = context.WithValue(ctx, adminNameKey, claims.Name)

		default:
			res.Json(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserIDFromContext(ctx context.Context) (uint64, bool) {
	v := ctx.Value(userIDKey)
	id, ok := v.(uint64)
	return id, ok
}

func AdminNameFromContext(ctx context.Context) (string, bool) {
	v := ctx.Value(adminNameKey)
	name, ok := v.(string)
	return name, ok
}

func ContextWithUserID(ctx context.Context, id uint32) context.Context {
	return context.WithValue(ctx, userIDKey, id)
}

func readToken(r *http.Request) (string, bool) {
	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
	if authHeader == "" {
		return "", false
	}

	if !strings.HasPrefix(authHeader, "Token ") {
		return "", false
	}

	tokenStr := strings.TrimSpace(strings.TrimPrefix(authHeader, "Token "))
	if tokenStr == "" {
		return "", false
	}

	return tokenStr, true
}
