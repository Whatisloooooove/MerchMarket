package server

import (
	"log"
	"merch_service/internal/database"
	"merch_service/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func MerchList(c *gin.Context) {
	merchlist, err := database.Connect().GetMerchList()

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			error_code: http.StatusInternalServerError,
			message:    InternalServerError,
			data:       struct{}{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		error_code: http.StatusOK,
		message:    MerchListOK,
		data:       merchlist,
	})
}

func WalletHistory(c *gin.Context) {
	log.Println("В хендлере пользовательской истории!")

	info := c.Keys["claims"].(jwt.MapClaims)
	userLogin := info["log"].(string)

	history, err := database.Connect().CoinsHistory(&models.User{
		Login: userLogin,
	})

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			error_code: http.StatusInternalServerError,
			message:    err.Error(),
			data:       struct{}{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		error_code: http.StatusOK,
		message:    HistoryOK,
		data:       history,
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
			c.JSON(http.StatusBadRequest, gin.H{
				error_code: http.StatusBadRequest,
				message:    InvalidAppDataError,
				data:       struct{}{},
			})
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
				error_code: http.StatusInternalServerError,
				message:    InternalServerError,
				data:       struct{}{},
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
					error_code: http.StatusInternalServerError,
					message:    InternalServerError,
					data:       struct{}{},
				})
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				error_code: http.StatusBadRequest,
				message:    UserExistsError,
				data:       struct{}{},
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			error_code: http.StatusOK,
			message:    RegistrationOK,
			data:       struct{}{},
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
			c.JSON(http.StatusBadRequest, gin.H{
				error_code: http.StatusBadRequest,
				message:    InvalidAppDataError,
				data:       struct{}{},
			})
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
				error_code: http.StatusInternalServerError,
				message:    InternalServerError,
				data:       struct{}{},
			})
			return
		}

		if !registered {
			c.JSON(http.StatusBadRequest, gin.H{
				error_code: http.StatusBadRequest,
				message:    UserNotFoundError,
				data:       struct{}{},
			})
			return
		}

		SendToken(c, config, &json)
	}
}

func TransferHandler(c *gin.Context) {
	// Логин отправляющего хранится в c.Keys["claims"] (см authorization.go:AuthRequired)
	// Достаточно узнать логин получателя и сумму
	var transactionReq TransactionRequest

	log.Println("В хендлере транзакций!")

	// Повторяющийся код, сделать умную обертку TODO
	if err := c.ShouldBindJSON(&transactionReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			error_code: http.StatusBadRequest,
			message:    InvalidAppDataError,
			data:       struct{}{},
		})
		return
	}

	info := c.Keys["claims"].(jwt.MapClaims) // (cм. замечание выше)
	log.Printf("Информация о переводе:\n\tОтправитель: %s\n\tПолучатель: %s\n\tСумма: %d\n",
		info["log"],
		transactionReq.Reciever,
		transactionReq.Amount)

	err := database.Connect().TransferCoins(&models.TransactionEntry{
		Sender:   info["log"].(string),
		Reciever: transactionReq.Reciever,
		Amount:   transactionReq.Amount,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			error_code: http.StatusInternalServerError, // or http.StatusBadRequest
			message:    InternalServerError,            // TODO, оповестить, что возможно не хватает баланса
			data:       struct{}{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		error_code: http.StatusOK,
		message:    TransferOK,
		data:       struct{}{}, // можно возвращать число оставшихся монет
	})
}
