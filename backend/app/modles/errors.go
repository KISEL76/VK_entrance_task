package models

type ErrorResponse struct {
	Code             int    `json:"code"`
	ErrorDescription string `json:"error_description"`
}
