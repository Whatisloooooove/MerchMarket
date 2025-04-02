// internal/models/user.go
package models

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Email     string `gorm:"unique"`
	Password  string
	CreatedAt time.Time
	Wallet    Wallet
}

type Wallet struct {
	ID      uint `gorm:"primaryKey"`
	UserID  uint
	Balance int
}
