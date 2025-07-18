package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	models "marketplace/app/modles"
)

type Handler struct {
	DB *sql.DB
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

// Проверка на ограничения по длине при регистрации
func validateRegister(req models.RegisterRequest) error {
	if len(req.Login) < 3 || len(req.Login) > 32 {
		return errors.New("логин должен быть от 3 до 32 символов")
	}

	if len(req.Password) < 8 || len(req.Password) > 64 {
		return errors.New("пароль должен быть от 8 до 64 символов")
	}
	return nil
}
