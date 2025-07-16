package models

type AdRequest struct {
	Title       string  `json:"title" binding:"required,min=5,max=100"`
	Description string  `json:"description" binding:"omitempty,max=500"`
	ImageURL    string  `json:"image_url" binding:"omitempty,url"`
	Price       float64 `json:"price" binding:"required,gte=0"`
}

type AdResponse struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url"`
	Price       float64 `json:"price"`
	AuthorLogin string  `json:"author_login"`
	IsOwner     bool    `json:"is_owner"`
}
