package handlers

import (
	"encoding/json"
	"fmt"
	"marketplace/app/utils"
	"net/http"
	"strconv"
	"strings"

	models "marketplace/app/modles"
)

// Вспомогательная функция для подсчёта общего количества товаров
func countGoods(h *Handler, whereClause string, args []interface{}, countArgsCount int) (int, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM goods g %s", whereClause)
	row := h.DB.QueryRow(query, args[:countArgsCount]...)
	var total int
	err := row.Scan(&total)
	return total, err
}

func (h *Handler) GoodsList(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	page := 1
	limit := 10

	// Получаем параметры сортировки и фильтрации из URL
	sort := queryParams.Get("sort")
	order := queryParams.Get("order")
	minPriceStr := queryParams.Get("min_price")
	maxPriceStr := queryParams.Get("max_price")

	// Парсим номер страницы
	if p, err := strconv.Atoi(queryParams.Get("page")); err == nil && p > 0 {
		page = p
	}
	offset := (page - 1) * limit

	// Разрешённые поля для сортировки
	allowedSortFields := map[string]bool{
		"price":      true,
		"created_at": true,
	}

	if sort == "" {
		sort = "created_at"
	}

	if !allowedSortFields[sort] {
		writeError(w, http.StatusBadRequest, "Неверное поле сортировки")
		return
	}

	// Проверка порядка сортировки
	order = strings.ToLower(order)
	if order != "asc" && order != "desc" && order != "" {
		writeError(w, http.StatusBadRequest, "Параметр 'order' может быть 'asc' или 'desc'")
		return
	}
	if order == "" {
		order = "desc"
	}

	// Фильтрация по цене
	var filters []string
	var args []interface{}
	argIndex := 1

	var minPrice, maxPrice float64
	var err error

	if minPriceStr != "" {
		minPrice, err = strconv.ParseFloat(minPriceStr, 64)
		if err != nil || minPrice < 0 {
			writeError(w, http.StatusBadRequest, "min_price должен быть положительным числом")
			return
		}
		filters = append(filters, fmt.Sprintf("g.price >= $%d", argIndex))
		args = append(args, minPrice)
		argIndex++
	}

	if maxPriceStr != "" {
		maxPrice, err = strconv.ParseFloat(maxPriceStr, 64)
		if err != nil || maxPrice < 0 {
			writeError(w, http.StatusBadRequest, "max_price должен быть положительным числом")
			return
		}
		filters = append(filters, fmt.Sprintf("g.price <= $%d", argIndex))
		args = append(args, maxPrice)
		argIndex++
	}

	if minPriceStr != "" && maxPriceStr != "" && minPrice > maxPrice {
		writeError(w, http.StatusBadRequest, "min_price не может быть больше max_price")
		return
	}

	// Проверяем, авторизован ли пользователь
	userID, _ := utils.ExtractUserIDFromContext(r.Context())
	isAuth := userID != 0

	// Сборка SQL WHERE
	whereClause := ""
	if len(filters) > 0 {
		whereClause = "WHERE " + strings.Join(filters, " AND ")
	}

	// Основной SQL-запрос для получения списка товаров
	query := fmt.Sprintf(`
		SELECT g.id, g.title, g.description, g.image_url, g.price, u.login, g.author_id
		FROM goods g
		JOIN users u ON g.author_id = u.id
		%s
		ORDER BY g.%s %s
		LIMIT $%d OFFSET $%d
	`, whereClause, sort, strings.ToUpper(order), argIndex, argIndex+1)

	args = append(args, limit, offset)

	// Выполняем SQL-запрос с подставленными параметрами
	rows, err := h.DB.Query(query, args...)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Ошибка при выполнении запроса к БД")
		return
	}
	defer rows.Close() // Гарантируем закрытие соединения с базой после чтения

	// Ответ для авторизованного пользователя — включает поле is_owner
	if isAuth {
		var goods []models.GoodsResponse
		for rows.Next() {
			var good models.GoodsResponse
			var authorID int
			err := rows.Scan(&good.ID, &good.Title, &good.Description, &good.ImageURL, &good.Price, &good.AuthorLogin, &authorID)
			if err != nil {
				writeError(w, http.StatusInternalServerError, "Ошибка чтения строки")
				return
			}
			good.IsOwner = (userID == authorID)
			goods = append(goods, good)
		}

		total, err := countGoods(h, whereClause, args, argIndex-1)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Ошибка при подсчёте количества товаров")
			return
		}

		// Считаем количество страниц — округление вверх
		resp := models.GoodsListResponse{
			Goods:         goods,
			Page:          page,
			GoodsQuantity: total,
			PageAmount:    (total + limit - 1) / limit,
		}
		json.NewEncoder(w).Encode(resp) // Отправляем JSON-ответ с полной информацией
	} else {
		// Ответ для неавторизованного пользователя — без поля is_owner
		var goods []models.GoodsPublicResponse
		for rows.Next() {
			var good models.GoodsPublicResponse
			err := rows.Scan(&good.ID, &good.Title, &good.Description, &good.ImageURL, &good.Price, &good.AuthorLogin, new(int))
			if err != nil {
				writeError(w, http.StatusInternalServerError, "Ошибка чтения строки")
				return
			}
			goods = append(goods, good)
		}

		total, err := countGoods(h, whereClause, args, argIndex-1)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Ошибка при подсчёте количества товаров")
			return
		}

		// Отправляем укороченный JSON-ответ без признака владения
		resp := models.GoodsListPublicResponse{
			Goods:         goods,
			Page:          page,
			GoodsQuantity: total,
			PageAmount:    (total + limit - 1) / limit,
		}
		json.NewEncoder(w).Encode(resp)
	}
}
