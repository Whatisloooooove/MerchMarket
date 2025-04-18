package handlers

import (
	"merch_service/new_version/internal/models"
	"merch_service/new_version/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type TransactionHandler struct {
	tServ service.TransactionServiceInterface
}

func (th *TransactionHandler) TransferHandler(c *gin.Context) {
	var req models.TransactionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			error_code: http.StatusBadRequest,
			message:    InvalidAppDataError,
			data:       struct{}{},
		})
		return
	}

	info := c.Keys["claims"].(jwt.MapClaims)
	sender := info["log"].(string)

	err := th.tServ.Send(c, sender, req.Reciever, req.Amount)

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
		message:    TransferOK,
		data:       struct{}{},
	})
}
