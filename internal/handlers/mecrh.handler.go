package handlers

import (
	"errors"
	"merch_service/internal/models"
	"merch_service/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// MerchHandler - структура мост, для связывания уровня хендлеров
// с сервисом мерча
type MerchHandler struct {
	mServ service.MerchServiceInterface
}

// NewMerchHandler - конуструирует *MerchHandler по MerchServiceInterface
func NewMerchHandler(mServ service.MerchServiceInterface) *MerchHandler {
	return &MerchHandler{mServ}
}

// MerchListHanlder - функция обработчик, отвечающий за возврат списка мерча
func (mh *MerchHandler) MerchListHandler(c *gin.Context) {
	response := DefaultResponse()

	merchlist, err := mh.mServ.MerchList(c)

	if err != nil {
		c.JSON(http.StatusOK, response)
		return
	}

	response.ErrorCode = http.StatusOK
	response.Message = MerchListOK
	response.Data = merchlist
	c.JSON(http.StatusOK, response)
}

// BuyMerchHandler - функция обработчик, отвечающий за покупку мерча
func (mh *MerchHandler) BuyMerchHandler(c *gin.Context) {
	response := DefaultResponse()

	var req models.PurchaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorCode = http.StatusBadRequest
		response.Message = InvalidAppDataError
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// При аутентификации (см middleware)
	// логин пользователя записывается в gin.Context
	info := c.Keys["claims"].(jwt.MapClaims)
	login := info["log"].(string)

	coins, err := mh.mServ.Buy(c, login, req.Item, req.Count)

	switch {
	case errors.Is(err, models.ErrNotEnoughMerch):
		response.ErrorCode = http.StatusBadRequest
		response.Message = NotEnoughMerchError
		c.JSON(http.StatusOK, response)
		return
	case errors.Is(err, models.ErrNotEnoughCoins):
		response.ErrorCode = http.StatusBadRequest
		response.Message = NotEnoughCoinsError
		c.JSON(http.StatusBadRequest, response)
		return
	case err != nil:
		response.Message = err.Error()
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response.ErrorCode = http.StatusOK
	response.Message = PurchaseOK
	response.Data = gin.H{
		balance: coins,
	}

	c.JSON(http.StatusOK, response)
}
