package handlers

import (
	"merch_service/new_version/configs"
	"merch_service/new_version/internal/models"
	"merch_service/new_version/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// UserHandler - структура мост, для связывания уровня хендлеров
// с пользовательским сервисом
type UserHandler struct {
	uServ service.UserServiceInterface
}

// NewUserHandler - конуструирует *UserHandler по UserServiceInterface
func NewUserHandler(uServ service.UserServiceInterface) *UserHandler {
	return &UserHandler{uServ}
}

// CoinsHistoryHandler - функция обработчик, отвечающий на запрос историй кошелька пользователя
// В случае успеха, в ответе возвращает историю кошелька в поле data
func (uh *UserHandler) CoinsHistoryHandler(c *gin.Context) {
	info := c.Keys["claims"].(jwt.MapClaims)
	userLogin := info["log"].(string)

	coinstHist, err := uh.uServ.CoinsHistory(c, userLogin)

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
		message:    HistoryCoinsOK,
		data:       coinstHist,
	})
}

// PurchaseHistoryHandler - функция обработчик, отвечающий на запрос покупок пользователя
// В случае успеха, в ответе возвращает список покупок в поле data
func (uh *UserHandler) PurchaseHistoryHandler(c *gin.Context) {
	info := c.Keys["claims"].(jwt.MapClaims)
	userLogin := info["log"].(string)

	pHist, err := uh.uServ.PurchaseHistory(c, userLogin)

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
		message:    HistoryPurchOK,
		data:       pHist,
	})

}

// RegHandler - обработчик, отвечающий за регистрацию пользователя
func (uh *UserHandler) RegHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.LoginRequest

		// см [validation](https://github.com/gin-gonic/gin/blob/master/docs/doc.md#model-binding-and-validation)
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				error_code: http.StatusBadRequest,
				message:    InvalidAppDataError,
				data:       struct{}{},
			})
			return
		}

		err := uh.uServ.Register(c, &req)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				error_code: http.StatusInternalServerError,
				message:    err.Error(),
				data:       struct{}{},
			})
			return
		}

		// err == nil <=> Пользователь успешно зарегистрирован
		c.JSON(http.StatusOK, gin.H{
			error_code: http.StatusOK,
			message:    RegistrationOK,
			data:       struct{}{},
		})
	}
}

// RegHandler - обработчик, отвечающий за аутентификацию пользователя
// В случае успеха, в ответе возварщает JWT и  JWTrefresh токены в поле data
func (uh *UserHandler) LoginHandler(config *configs.ServerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.LoginRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				error_code: http.StatusBadRequest,
				message:    InvalidAppDataError,
				data:       struct{}{},
			})
			return
		}

		err := uh.uServ.Login(c, &req)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				error_code: http.StatusInternalServerError,
				message:    InternalServerError,
				data:       struct{}{},
			})
			return
		}

		SendToken(c, config, &req)
	}
}
