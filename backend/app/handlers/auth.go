package handlers

import (
	"database/sql"
	"encoding/json"
	models "marketplace/app/modles"
	"marketplace/app/utils"
	"net/http"
	"time"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Ошибка: поддерживается только метод Post для регистрации", http.StatusMethodNotAllowed)
		return
	}

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Ошибка: невалидный JSON")
		return
	}

	if err := validateRegister(req); err != nil {
		writeError(w, http.StatusBadRequest, "Ошибка: слишком длинный/короткий пароль/логин")
		return
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Ошибка: ошибка при хэшровании пароля")
		return
	}

	var userID int
	err = h.DB.QueryRow(`
		INSERT INTO users (login, password_hash, created_at)
		VALUES ($1, $2, $3)
		RETURNING id
	`, req.Login, hash, time.Now()).Scan(&userID)

	if err != nil {
		if utils.IsUniqueViolation(err) {
			writeError(w, http.StatusConflict, "Ошибка: логин уже занят")
		} else {
			writeError(w, http.StatusInternalServerError, "Ошибка базы данных на стороне сервера")
		}
		return
	}

	resp := models.RegisterResponse{
		Id:    userID,
		Login: req.Login,
	}
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Ошибка: поддерживается только метод Post для регистрации", http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Ошибка: невалидный JSON")
		return
	}

	var id int
	var login string
	var hash string

	err := h.DB.QueryRow(`
	SELECT id, login, password_hash
	FROM users
	WHERE login = $1
	`,
		req.Login).Scan(&id, &login, &hash)

	if err == sql.ErrNoRows {
		writeError(w, http.StatusUnauthorized, "Ошибка: неверный логин или пароль")
		return
	} else if err != nil {
		writeError(w, http.StatusInternalServerError, "Ошибка базы данных на стороне сервера")
		return
	}

	if !utils.CheckPasswordHash(req.Password, hash) {
		writeError(w, http.StatusUnauthorized, "Ошибка: неверный логин или пароль")
		return
	}

	token, err := utils.GenerateJWT(id, login)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Ошибка генерации токена")
		return
	}

	resp := models.AuthResponse{Token: token}
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
