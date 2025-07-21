package models

type GoodsRequest struct {
	Title       string  `json:"title" binding:"required,min=5,max=100"`
	Description string  `json:"description" binding:"omitempty,max=500"`
	ImageURL    string  `json:"image_url" binding:"omitempty,url"`
	Price       float64 `json:"price" binding:"required,gte=0"`
}

// возвращаем паблик структуру при отстутствии авторизации, и стандартную с авторизацией
type GoodsResponse struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url"`
	Price       float64 `json:"price"`
	AuthorLogin string  `json:"author_login"`
	IsOwner     bool    `json:"is_owner"`
}

type GoodsPublicResponse struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url"`
	Price       float64 `json:"price"`
	AuthorLogin string  `json:"author_login"`
}

type GoodsListResponse struct {
	Goods         []GoodsResponse `json:"goods"`
	Page          int             `json:"current_page"`
	GoodsQuantity int             `json:"goods_quantity"`
	PageAmount    int             `json:"page_amount"`
}

type GoodsListPublicResponse struct {
	Goods         []GoodsPublicResponse `json:"goods"`
	Page          int                   `json:"current_page"`
	GoodsQuantity int                   `json:"goods_quantity"`
	PageAmount    int                   `json:"page_amount"`
}
