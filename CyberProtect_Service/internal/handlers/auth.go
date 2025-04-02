// Содержит обработчики для аутентификации:
// LoginHandler: Обработка запроса на авторизацию. При отсутствии пользователя выполняется автогенерация (регистрация) и инициализация кошелька.
// MeHandler: Возвращает информацию о текущем пользователе (используя user_id, извлечённый из JWT).
package handlers

import (
	"merch_service/internal/db"
	"merch_service/internal/models"

	"merch_service/pkg/auth"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// LoginRequest представляет структуру входных данных для авторизации
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginHandler обрабатывает запрос авторизации.
// Если пользователь не найден, выполняется автозарегистрция и инициализация кошелька.
func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный запрос"})
		return
	}

	var user models.User
	result := db.DB.Preload("Wallet").Where("email = ?", req.Email).First(&user)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Автоматическая регистрация нового пользователя и инициализация кошелька с 1000 монет
			user = models.User{Email: req.Email, Password: req.Password}
			db.DB.Create(&user)
			db.DB.Create(&models.Wallet{UserID: user.ID, Balance: 1000})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка базы данных"})
			return
		}
	}

	// Генерация JWT токена для пользователя
	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось сгенерировать токен"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// MeHandler возвращает информацию о текущем пользователе.
// Для работы требуется, чтобы JWT мидлвара установила значение "user_id" в контексте.
func MeHandler(c *gin.Context) {
	userID := c.GetUint("user_id")

	var user models.User
	err := db.DB.Preload("Wallet").First(&user, userID).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":     user.ID,
		"email":  user.Email,
		"wallet": user.Wallet.Balance,
	})
}
