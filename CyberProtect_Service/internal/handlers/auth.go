// internal/handlers/auth.go
package handlers

import (
	"merch_service/internal/db"
	"merch_service/internal/models"
	"merch_service/pkg/auth"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func MeHandler(c *gin.Context) {
	userID := c.GetUint("user_id")

	var user models.User
	err := db.DB.Preload("Wallet").First(&user, userID).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":     user.ID,
		"email":  user.Email,
		"wallet": user.Wallet.Balance,
	})
}

func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var user models.User
	result := db.DB.Preload("Wallet").Where("email = ?", req.Email).First(&user)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Auto-register new user and initialize wallet with 1000 coins
			user = models.User{Email: req.Email, Password: req.Password}
			db.DB.Create(&user)
			db.DB.Create(&models.Wallet{UserID: user.ID, Balance: 1000})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
	}

	// Generate JWT token for the user
	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
