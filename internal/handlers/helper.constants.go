package handlers

const (
	AuthError           = "отказано в доступе"
	TokenGenError       = "ошибка генерации токена"
	UserExistsError     = "пользователь существует"
	UserNotFoundError   = "пользователя с таким логином не существует"
	WrongPassError      = "неверный пароль" // Нужен ли он вообще (см handlers.go)
	InternalServerError = "ошибка на сервере"
	InvalidAppDataError = "неверный формат данных в запросе"
	NotEnoughMerchError = "недостаточно товара на складе"
	NotEnoughCoinsError = "недостаточно монет для покупки"
)

const (
	RegistrationOK = "регистрация успешна"
	TokensOK       = "токены успешно созданы"
	RefreshOK      = "токен авторизаций обновлен"
	MerchListOK    = "список товаров"
	HistoryCoinsOK = "история кошелька"
	HistoryPurchOK = "история покупок"
	TransferOK     = "перевод монет успешен"
	PurchaseOK     = "покупка успешна"
)

// Для централизованного контроля за API и для избежания очепяток
const (
	error_code = "error_code"
	message    = "message"
	data       = "data"
	token      = "token"
	refresh    = "refresh"
	balance    = "balance"
)
