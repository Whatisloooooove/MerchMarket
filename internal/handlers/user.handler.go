package handlers

import (
	"errors"
	"merch_service/configs"
	"merch_service/internal/models"
	"merch_service/internal/service"
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
	response := DefaultResponse()

	info := c.Keys["claims"].(jwt.MapClaims)
	userLogin := info["log"].(string)

	coinsHist, err := uh.uServ.CoinsHistory(c, userLogin)

	if err != nil {
		c.JSON(http.StatusOK, response)
		return
	}

	response.ErrorCode = http.StatusOK
	response.Message = HistoryCoinsOK
	response.Data = coinsHist
	c.JSON(http.StatusOK, response)
}

// PurchaseHistoryHandler - функция обработчик, отвечающий на запрос покупок пользователя
// В случае успеха, в ответе возвращает список покупок в поле data
func (uh *UserHandler) PurchaseHistoryHandler(c *gin.Context) {
	response := DefaultResponse()

	info := c.Keys["claims"].(jwt.MapClaims)
	userLogin := info["log"].(string)

	pHist, err := uh.uServ.PurchaseHistory(c, userLogin)

	if err != nil {
		c.JSON(http.StatusOK, response)
		return
	}

	response.ErrorCode = http.StatusOK
	response.Message = HistoryPurchOK
	response.Data = pHist
	c.JSON(http.StatusOK, response)
}

// RegHandler - обработчик, отвечающий за регистрацию пользователя
func (uh *UserHandler) RegHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		response := DefaultResponse()
		var req models.LoginRequest

		// см [validation](https://github.com/gin-gonic/gin/blob/master/docs/doc.md#model-binding-and-validation)
		if err := c.ShouldBindJSON(&req); err != nil {
			response.ErrorCode = http.StatusBadRequest
			response.Message = InvalidAppDataError
			c.JSON(http.StatusBadRequest, response)
			return
		}

		err := uh.uServ.Register(c, &req)

		switch {
		case errors.Is(err, models.ErrUserExists):
			response.ErrorCode = http.StatusBadRequest
			response.Message = UserExistsError
			c.JSON(http.StatusBadRequest, response)
			return
		case err != nil:
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		// err == nil <=> Пользователь успешно зарегистрирован
		response.ErrorCode = http.StatusOK
		response.Message = RegistrationOK
		c.JSON(http.StatusOK, response)
	}
}

// RegHandler - обработчик, отвечающий за аутентификацию пользователя
// В случае успеха, в ответе возварщает JWT и  JWTrefresh токены в поле data
func (uh *UserHandler) LoginHandler(config *configs.ServerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		response := DefaultResponse()

		var req models.LoginRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			response.ErrorCode = http.StatusBadRequest
			response.Message = InvalidAppDataError
			c.JSON(http.StatusBadRequest, response)
			return
		}

		err := uh.uServ.Login(c, &req)

		switch {
		case errors.Is(err, models.ErrWrongPassword):
			response.ErrorCode = http.StatusBadRequest
			response.Message = WrongPassError
			return
		case err != nil:
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		SendToken(c, config, &req)
	}
}
