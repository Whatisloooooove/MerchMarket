// internal/models/user.go
// Этот файл содержит модель User и её основные поля.
// Используется для хранения информации о пользователях системы.
package models

import "time"

// User представляет пользователя системы.
type User struct {
	ID        uint      `gorm:"primaryKey"` // Уникальный идентификатор пользователя
	Email     string    `gorm:"unique"`     // Email пользователя (уникальный)
	Password  string    // Хэшированный пароль
	CreatedAt time.Time // Дата создания записи
	Wallet    Wallet    // Связь "один к одному" с кошельком пользователя
}
