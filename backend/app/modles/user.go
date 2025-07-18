package models

type LoginRequest struct {
	Login    string `json:"login" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type RegisterRequest struct {
	Login    string `json:"login" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}
type RegisterResponse struct {
	Id    int    `json:"id"`
	Login string `json:"login"`
}
