package handlers

import (
	"errors"
	"log"
	"merch_service/configs"
	"merch_service/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthRequired(config *configs.ServerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("Авторизация...")
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
			// Пересылаем новый токен, если старый больше не валиден
			if errors.Is(err, jwt.ErrTokenExpired) {
				log.Println("token expired, resending...")
				ResendToken(c, config, token)
				return
			}

			log.Println("ошибка при переотправке токена", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				error_code: http.StatusUnauthorized,
				message:    AuthError,
				data:       struct{}{},
			})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			log.Println("Claims:", claims)
			// Сораняем claims в мапу gin.Context (для дальнейшего использования обработчиком)
			c.Set("claims", claims)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				error_code: http.StatusUnauthorized,
				message:    AuthError,
				data:       struct{}{},
			})
			c.Abort()
			return
		}
		// Если авторизация успешна, вызвваем обработчик
		log.Println("Авторизация успешна, идем в хендлер...")
		c.Next()
	}
}

// SendToken возвращает JWT токены для авторизации (с таймаутом) и для обновления первого
// В случае успешной генерации пользователь получит JSON формата:
//
//		{
//	   error_code: ...,
//	   message: ...,
//		  data: {
//		  	"token": xxx.yyy.zzz,
//		  	"refresh": aaa.bbb.ccc,
//			}
//		}
func SendToken(c *gin.Context, config *configs.ServerConfig, json *models.LoginRequest) {
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
	tokenString, err := JwtToken(config, json)

	if err != nil {
		log.Println("ошибка генерации токена:", err.Error()) // FOR DEBUG ONLY
		c.JSON(http.StatusInternalServerError, gin.H{
			error_code: http.StatusInternalServerError,
			message:    TokenGenError,
			data:       struct{}{},
		})
		return
	}

	refreshTokenString, err := RefreshToken(config, json)

	if err != nil {
		log.Println("ошибка генерации токена:", err.Error()) // FOR DEBUG ONLY
		c.JSON(http.StatusInternalServerError, gin.H{
			error_code: http.StatusInternalServerError,
			message:    TokenGenError,
			data:       struct{}{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		error_code: http.StatusOK,
		message:    TokensOK,
		data: gin.H{
			refresh: refreshTokenString,
			token:   tokenString,
		},
	})
}

// JwtToken - генерирует токен авторизации с сроком истечения указанным в config
func JwtToken(config *configs.ServerConfig, json *models.LoginRequest) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"log": json.Login,
		"exp": time.Now().Add(time.Second * time.Duration(config.ExpTimeout)).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.Secret))

	return tokenString, err
}

// RefreshToken - генерирует refresh токен
func RefreshToken(config *configs.ServerConfig, json *models.LoginRequest) (string, error) {
	// TODO: Каким должен быть body рефреш токена?
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"log": json.Login,
		"typ": "refresh",
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(config.RefreshSecret))

	return refreshTokenString, err
}

// ResendToken - отправляет новый токен авторизации в случае истечения срока старого токена
func ResendToken(c *gin.Context, config *configs.ServerConfig, expiredToken *jwt.Token) {
	var login string
	if claims, ok := expiredToken.Claims.(jwt.MapClaims); ok {
		log.Println("Claims:", claims)
		login = claims["log"].(string) // Здесь точно string!
	}

	newToken, err := JwtToken(config, &models.LoginRequest{Login: login})

	if err != nil {
		log.Println("ошибка генерации токена:", err.Error()) // FOR DEBUG ONLY
		c.JSON(http.StatusInternalServerError, gin.H{
			error_code: http.StatusInternalServerError,
			message:    TokenGenError,
			data:       struct{}{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		error_code: http.StatusUnauthorized,
		message:    RefreshOK,
		data: gin.H{
			token: newToken,
		},
	})

	// Отменяем выполнение хендлеров
	c.Abort()
}
