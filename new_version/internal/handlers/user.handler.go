package handlers

import (
	"merch_service/new_version/internal/models"
	"merch_service/new_version/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type UserHandler struct {
	uServ service.UserServiceInterface
}

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

func (uh *UserHandler) LoginHandler() gin.HandlerFunc {
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

		// SendToken(c, config, &json)
	}
}
