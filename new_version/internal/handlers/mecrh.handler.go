package handlers

import (
	"merch_service/new_version/internal/models"
	"merch_service/new_version/internal/service"
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
	merchlist, err := mh.mServ.MerchList(c)

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

// BuyMerchHandler - функция обработчик, отвечающий за покупку мерча
func (mh *MerchHandler) BuyMerchHandler(c *gin.Context) {
	var req models.PurchaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			error_code: http.StatusBadRequest,
			message:    InvalidAppDataError,
			data:       struct{}{},
		})
		return
	}

	// При аутентификации (см middleware)
	// логин пользователя записывается в gin.Context
	info := c.Keys["claims"].(jwt.MapClaims)
	userLogin := info["log"].(string)

	coins, err := mh.mServ.Buy(c, userLogin, req.ItemName, req.Count)

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
		message:    PurchaseOK,
		data:       coins,
	})
}
