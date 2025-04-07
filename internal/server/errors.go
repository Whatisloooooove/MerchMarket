package server

const (
	AuthError         = "отказано в доступе"
	TokenGenError     = "ошибка генерации токена"
	UserExistsError   = "пользователь существует"
	UserNotFoundError = "пользователя с таким логином не существует"
	WrongPassError    = "неверный пароль" // Нужен ли он вообще (см handlers.go)
)
