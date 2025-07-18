package utils

import (
	"context"
	"encoding/json"
	"errors"
	models "marketplace/app/modles"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userIDKey contextKey = "user_id"

func WithAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			writeError(w, http.StatusUnauthorized, "Ошибка: не введен верный токен авторизации")
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			writeError(w, http.StatusUnauthorized, "Ошибка: невалидный токен авторизации")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["user_id"] == nil {
			writeError(w, http.StatusUnauthorized, "Ошибка: не соблюдены требования создания токена")
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			writeError(w, http.StatusUnauthorized, "Ошибка: поле user_id отсутствует в токене")
			return
		}
		userID := int(userIDFloat)

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next(w, r.WithContext(ctx))
	}
}

func WithTokenIfPresent(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			next(w, r)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			writeError(w, http.StatusUnauthorized, "Ошибка: невалидный токен авторизации")
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			writeError(w, http.StatusUnauthorized, "Ошибка: не введен верный токен авторизации")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			writeError(w, http.StatusUnauthorized, "Ошибка: не соблюдены требования создания токена")
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			writeError(w, http.StatusUnauthorized, "Ошибка: поле user_id отсутствует в токене")
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, int(userIDFloat))
		next(w, r.WithContext(ctx))
	}
}

// Достаем id из контекста
func ExtractUserIDFromContext(ctx context.Context) (int, error) {
	val := ctx.Value(userIDKey)
	userID, ok := val.(int)
	if !ok {
		return 0, errors.New("user_id not found in context")
	}
	return userID, nil
}

// Вставка для ответа сервера с ошибкой
func writeError(w http.ResponseWriter, code int, desc string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(models.ErrorResponse{
		Code:             code,
		ErrorDescription: desc,
	})
}
