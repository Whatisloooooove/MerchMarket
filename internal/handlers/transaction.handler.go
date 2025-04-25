package handlers

import (
	"errors"
	"merch_service/internal/models"
	"merch_service/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// TransactionHandler - структура мост, для связывания уровня хендлеров
// с транзакционным сервисом
type TransactionHandler struct {
	tServ service.TransactionServiceInterface
}

// NewTransactionHandler - конуструирует *TransactionHandler по TransactionServiceInterface
func NewTransactionHandler(tServ service.TransactionServiceInterface) *TransactionHandler {
	return &TransactionHandler{tServ}
}

// TransferHandler - функция обработчик переводов монет
func (th *TransactionHandler) TransferHandler(c *gin.Context) {
	response := DefaultResponse()
	var req models.TransactionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorCode = http.StatusBadRequest
		response.Message = InvalidAppDataError
		c.JSON(http.StatusBadRequest, response)
		return
	}

	info := c.Keys["claims"].(jwt.MapClaims)
	sender := info["id"].(int)

	err := th.tServ.Send(c, sender, req.RecieverId, req.Amount)

	switch {
	case errors.Is(err, models.ErrNotEnoughCoins):
		response.ErrorCode = http.StatusBadRequest
		response.Message = NotEnoughCoinsError
		c.JSON(http.StatusBadRequest, response)
		return
	case err != nil:
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response.ErrorCode = http.StatusOK
	response.Message = TransferOK
	c.JSON(http.StatusOK, response)
}
