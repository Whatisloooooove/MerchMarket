package server

import (
	"log"
	"merch_service/internal/database"
	"merch_service/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func MerchList(c *gin.Context) {
	merchlist, err := database.Connect().GetMerchList()

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error_code": http.StatusInternalServerError,
			"message":    InternalServerError,
			"data":       struct{}{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error_code": http.StatusOK,
		"message":    MerchListOK,
		"data":       merchlist,
	})
}

func WalletHistory(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"after":  123,
		"before": 321,
	})
}

// LoginRequest структура запроса на регистрацию
// и на вход в систему
// считаем что POST запросы делаются с хедером ContentType: application/json (yaml, xml можно добавить позже)
type LoginRequest struct {
	Login string `json:"login"`
	Pass  string `json:"pass"`
}

// RegHanlder возвращает обработчик регистрации с заданными параметрами
// для JWT токена
//
// Сам обработчик записывает пользователя в базу данных
// Если пользователь существует возвращает код ошибки http.StatusBadRequest
func RegHandler(config *ServerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("In registration handler!") // Remove after DEBUG !

		var json LoginRequest
		// см [validation](https://github.com/gin-gonic/gin/blob/master/docs/doc.md#model-binding-and-validation)
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Проверка существования записи в БД
		registered, err := database.Connect().CheckUser(&models.User{
			Login:    json.Login,
			Password: json.Pass,
			// Поле email добавим позже
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error_code": http.StatusInternalServerError,
				"message":    InternalServerError,
				"data":       struct{}{},
			})
			return
		}

		if !registered {
			err := database.Connect().RegisterUser(&models.User{
				Login:    json.Login,
				Password: json.Pass,
			})

			// По идее insert в бд должен быть успешен, если только он не отвалился
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error_code": http.StatusInternalServerError,
					"message":    InternalServerError,
					"data":       struct{}{},
				})
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error_code": http.StatusBadRequest,
				"message":    UserExistsError,
				"data":       struct{}{},
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"error_code": http.StatusOK,
			"message":    RegistrationOK,
			"data":       struct{}{},
		})
	}
}

// LoginHandler возвращает обработчик запроса на вход
// В случае успеха возвращает пользователю JWT токен (через response)
//
// Если пользователя не существует возврощает код ошибки http.StatusBadRequest
func LoginHandler(config *ServerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("In login handler!") // Remove after DEBUG !
		var json LoginRequest
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Проверка существования записи в БД
		registered, err := database.Connect().CheckUser(&models.User{
			Login:    json.Login,
			Password: json.Pass,
			// Поле email добавим позже
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error_code": http.StatusInternalServerError,
				"message":    InternalServerError,
				"data":       struct{}{},
			})
			return
		}

		if !registered {
			c.JSON(http.StatusBadRequest, gin.H{
				"error_code": http.StatusBadRequest,
				"message:":   UserNotFoundError,
				"data":       struct{}{},
			})
			return
		}

		SendToken(c, config, &json)
	}
}
