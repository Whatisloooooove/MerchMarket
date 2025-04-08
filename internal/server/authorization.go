package server

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Функция 'middleware' может
// возвращать обработчик соответствующий типу gin.HandlerFunc
// Насколько я понял, это делается для соблюдения паттерна фабрики

// AuthRequired вспомогательная функция для проверки прав доступа пользователя
// к API, срабатывающая до перехода к обработке запроса
func AuthRequired(config *ServerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("HELLO! FROM AUTH MIDDLEWARE") // REMOVE AFTER DEBUG

		// Зачем нужен Bearer в ключе??? (Мы поддерживаем только JWT, или нет?)
		tokenString := c.GetHeader("Authorization")

		// JWT магия. Проверка токена на соответствие алгоритму
		// шифрования, далее на валидность.
		// Подробнее см соотв. документацию [go-jwt](https://pkg.go.dev/github.com/golang-jwt/jwt#section-readme)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return []byte(config.Secret), nil
		})

		if err != nil || !token.Valid {
			log.Println(err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": AuthError})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			log.Println("Claims:", claims)
			// Сохраняем claims в мапу gin.Context (для дальнейшего использования обработчиком)
			c.Set("claims", claims)
			log.Println("DEBUG gin context keys:", c.Keys["claims"])
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": AuthError})
			c.Abort()
			return
		}
		// Если авторизация успешна, вызвваем обработчик
		c.Next()
	}
}

// SendToken возвращает JWT токены для авторизации (с таймаутом) и для обновления первого
// В случае успешной генерации пользователь получит JSON формата:
//
//	{
//	  "token": xxx.yyy.zzz,
//	  "refresh": aaa.bbb.ccc,
//	}
func SendToken(c *gin.Context, config *ServerConfig, json *LoginRequest) {
	// TODO стоит ли передавать в CreateToken объект LoginRequest
	// полученный в RegHandler?
	// На данный момент снова парсим json
	// var json LoginRequest
	// if err := c.ShouldBindJSON(&json); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	return
	// }

	// JWT магия
	// см [jwt](https://jwt.io/introduction)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"log": json.Login,
		"exp": time.Now().Add(time.Second * time.Duration(config.ExpTimeout)).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.Secret))
	if err != nil {
		log.Println("ошибка генерации токена:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": TokenGenError})
		return
	}

	// TODO: Каким должен быть body рефреш токена?
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"log": json.Login,
		"typ": "refresh",
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(config.RefreshSecret))
	if err != nil {
		log.Println("ошибка генерации токена:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": TokenGenError})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error_code": http.StatusOK,
		"message":    TokensOK,
		"data": gin.H{
			"refresh": refreshTokenString,
			"token":   tokenString,
		},
	})
}
