// Этот файл содержит middleware для аутентификации с использованием JWT (JSON Web Token).
// Мидлвара проверяет наличие и валидность JWT токена в заголовке Authorization каждого входящего запроса.
package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("your-secret-key")

// JWTAuth возвращает middleware, которое проверяет JWT токен в заголовке Authorization.
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем заголовок Authorization
		authHeader := c.GetHeader("Authorization")
		// Проверяем, что он начинается с Bearer
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Извлекаем сам токен, убирая "Bearer " из начала
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// Парсим токен
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			// Проверяем метод подписи токена
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return jwtKey, nil // возвращаем секретный ключ
		})

		// Если токен не валидный или ошибка при парсинге, завершаем запрос с ошибкой
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Извлекаем claims из токена (user_id)
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Устанавливаем user_id в контекст для использования в дальнейших обработчиках
			c.Set("user_id", uint(claims["user_id"].(float64))) // Преобразуем в uint
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		}
	}
}
