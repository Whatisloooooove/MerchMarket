package entites

import "context"

// Merch описывает id, название, стоимость и количество мерча.
// Реализует сущность мерч в бд
type Merch struct {
	// ID содержит индификатор мерча
	ID int

	// Name содержит наименование мерча
	Name string

	// Price содержит стоимость одной позиции мерча
	Price int

	// Stock содержит количесвто единиц позиции мерча на складе
	Stock int
}

// MerchStorage определяет контракт для работы с товарами
type MerchStorage interface {
	// Create создает новый экземпляр Merch, возвращает ошибку
	Create(ctx context.Context, merch *Merch) error
	
	// GetByID возвращает экземпляр Merch и ошибку по ID
	GetByID(ctx context.Context, id int) (*Merch, error)

	// GetByName возвращает экземпляр Merch и ошибку по Name
	GetByName(ctx context.Context, id int) (*Merch, error)
	
	// GetList возвращает слайс Merch и ошибку
	GetList(ctx context.Context) ([]*Merch, error)
	
	// Update обновляет информацию о Merch, возвращает ошибку
	Update(ctx context.Context, merch *Merch) error
	
	// ReduceStock уменьшает количество товара на складе на qty, возвращает ошибку
	ReduceStock(ctx context.Context, id int, qty int) error
	
	// IncreaseStock увеличивает количество товара на складе на qty, возвращает ошибку
	IncreaseStock(ctx context.Context, id int, qty int) error

	// Delete удаляет товар, возвращает ошибку
	DeleteByID(ctx context.Context, id int) error
}

