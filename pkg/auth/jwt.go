// internal/middlewares/jwt.go
// Этот файл содержит мидлвару для JWT аутентификации.
// Он проверяет токен, извлекает из него user_id и передаёт его дальше по цепочке.
package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("your-secret-key")

// GenerateToken генерирует JWT токен для заданного userID с истечением через 72 часа.
func GenerateToken(userID uint) (string, error) {
	// Создаем объект claims с пользовательским идентификатором
	claims := jwt.MapClaims{
		"user_id": userID,                                // Передаем user_id в качестве данных для токена
		"exp":     time.Now().Add(72 * time.Hour).Unix(), // Время истечения токена (72 часа)
	}

	// Создаем токен с использованием HMAC с SHA-256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен с использованием секрета и возвращаем строковое представление токена
	return token.SignedString(jwtKey)
}
