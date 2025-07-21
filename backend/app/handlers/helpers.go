package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	models "marketplace/app/modles"
)

type Handler struct {
	DB *sql.DB
}

const maxImageSize = 5 * 1024 * 1024 // 5 MB

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

// Проверка на корректное изображение
func isImageValid(url string) (bool, string) {
	resp, err := http.Head(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return false, "Ошибка: изображение недоступно по указанной ссылке"
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" {
		return false, "Ошибка: допустимые форматы изображения — JPEG, PNG"
	}

	contentLength := resp.Header.Get("Content-Length")
	if contentLength != "" {
		size, err := strconv.Atoi(contentLength)
		if err == nil && size > maxImageSize {
			return false, "Ошибка: изображение превышает максимальный размер 5 MB"
		}
	}

	return true, ""
}
