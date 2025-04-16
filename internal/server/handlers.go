package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var merch = []struct {
	Name    string
	MerchId int
	Price   int
	Stock   int
}{
	{
		MerchId: 1,
		Name:    "merch1",
		Price:   100,
		Stock:   20,
	},
	{
		MerchId: 2,
		Name:    "merch2",
		Price:   100,
		Stock:   20,
	}, {
		MerchId: 3,
		Name:    "merch3",
		Price:   100,
		Stock:   20,
	}, {
		MerchId: 4,
		Name:    "merch4",
		Price:   100,
		Stock:   20,
	}, {
		MerchId: 5,
		Name:    "merch5",
		Price:   100,
		Stock:   20,
	},
}

func MerchList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"error_code": http.StatusOK,
		"message":    MerchListOK,
		"data":       merch,
	})
}

func WalletHistory(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"after":  123,
		"before": 321,
	})
}

// Временно!
var users map[string]string = make(map[string]string)

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

		// Здесь должно быть обращение к базе данных для проверки
		// существования записи пользователя

		var json LoginRequest
		// см [validation](https://github.com/gin-gonic/gin/blob/master/docs/doc.md#model-binding-and-validation)
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userLogin := json.Login
		userPass := json.Pass

		if _, userExists := users[userLogin]; !userExists {
			users[userLogin] = userPass
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": UserExistsError})
			return
		}

		// Не нужен. Токены выдаются при /auth/login
		// SendToken(c, config, &json)
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

		if pass, userExists := users[json.Login]; userExists {
			if pass != json.Pass {
				// TODO: Наверно не очень хорошо говорить что пароль неправильный? (Могут брутфорсить, нужнен ratelimiter?)
				c.JSON(http.StatusBadRequest, gin.H{"error": WrongPassError})
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": UserNotFoundError})
			return
		}

		// TODO вернуть токен
		SendToken(c, config, &json)
	}
}
