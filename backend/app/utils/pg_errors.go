package utils

import (
	"errors"

	"github.com/lib/pq"
)

// Если ошибка является ошибкой от СУБД на уникальность записи, возвращаем true
func IsUniqueViolation(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "23505"
	}
	return false
}
