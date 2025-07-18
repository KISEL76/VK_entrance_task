package handlers

import (
	"encoding/json"
	models "marketplace/app/modles"
	"marketplace/app/utils"
	"net/http"
)

func (h *Handler) CreateAd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Ошибка: поддерживается только метод Post для создания объявления", http.StatusMethodNotAllowed)
		return
	}

	userID, err := utils.ExtractUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Ошибка: неавторизированный пользователь", http.StatusUnauthorized)
		return
	}

	var req models.GoodsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Ошибка: невалидный JSON")
		return
	}

	var id int
	var login string

	err = h.DB.QueryRow(`
		INSERT INTO goods (title, description, image_url, price, author_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id,
		(SELECT login FROM users WHERE id = $5)
	`,
		req.Title, req.Description, req.ImageURL, req.Price, userID).Scan(&id, &login)

	if err != nil {
		writeError(w, http.StatusInternalServerError, "Ошибка при создании объявления")
		return
	}

	resp := models.GoodsResponse{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		Price:       req.Price,
		AuthorLogin: login,
		IsOwner:     true,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
