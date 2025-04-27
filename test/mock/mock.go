package mock

import (
	"context"
	"merch_service/internal/models"
	"merch_service/internal/storage/entities"
	"sync"
	"time"
)

var (
	_ entities.MerchStorage       = (*MockMerchStorage)(nil)
	_ entities.UserStorage        = (*MockUserStorage)(nil)
	_ entities.TransactionStorage = (*MockTransactionStorage)(nil)
)

// MockUserStorage реализация
type MockUserStorage struct {
	mu     sync.RWMutex
	users  map[int]*models.User
	byName map[string]*models.User
	coins  map[int][]models.CoinsEntry
	purch  map[int][]models.PurchaseEntry
}

func NewMockUserStorage() *MockUserStorage {
	return &MockUserStorage{
		users:  make(map[int]*models.User),
		byName: make(map[string]*models.User),
		coins:  make(map[int][]models.CoinsEntry),
		purch:  make(map[int][]models.PurchaseEntry),
	}
}

func (s *MockUserStorage) Create(ctx context.Context, user *models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.byName[user.Login]; exists {
		return models.ErrUserExists
	}

	user.Id = len(s.users) + 1
	s.users[user.Id] = user
	s.byName[user.Login] = user
	s.coins[user.Id] = []models.CoinsEntry{{
		Id:          1,
		Date:        time.Now(),
		CoinsBefore: 0,
		CoinsAfter:  1000,
	}}
	return nil
}

func (s *MockUserStorage) Get(ctx context.Context, id int) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[id]
	if !exists {
		return nil, models.ErrUserNotFound
	}
	return user, nil
}

func (s *MockUserStorage) GetByLogin(ctx context.Context, login string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.byName[login]
	if !exists {
		return nil, models.ErrUserNotFound
	}
	return user, nil
}

func (s *MockUserStorage) Update(ctx context.Context, user *models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[user.Id]; !exists {
		return models.ErrUserNotFound
	}

	s.users[user.Id] = user
	return nil
}

func (s *MockUserStorage) Delete(ctx context.Context, id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[id]; !exists {
		return models.ErrUserNotFound
	}

	delete(s.users, id)
	return nil
}

func (s *MockUserStorage) GetCoinsHistory(ctx context.Context, userLogin string) ([]models.CoinsEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.byName[userLogin]
	if !exists {
		return nil, models.ErrUserNotFound
	}

	history, exists := s.coins[user.Id]
	if !exists {
		return nil, models.ErrUserNotFound
	}
	return history, nil
}

func (s *MockUserStorage) GetPurchaseHistory(ctx context.Context, userLogin string) ([]models.PurchaseEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.byName[userLogin]
	if !exists {
		return nil, models.ErrUserNotFound
	}

	history, exists := s.purch[user.Id]
	if !exists {
		return nil, models.ErrUserNotFound
	}
	return history, nil
}

// MockMerchStorage реализация
type MockMerchStorage struct {
	mu    sync.RWMutex
	items map[int]*models.Item
}

func NewMockMerchStorage() *MockMerchStorage {
	return &MockMerchStorage{
		items: map[int]*models.Item{
			1: {Id: 1, Name: "Футболка", Price: 50, Stock: 10},
			2: {Id: 2, Name: "Кружка", Price: 30, Stock: 5},
		},
	}
}

func (s *MockMerchStorage) Create(ctx context.Context, merch *models.Item) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	merch.Id = len(s.items) + 1
	s.items[merch.Id] = merch
	return nil
}

func (s *MockMerchStorage) Get(ctx context.Context, id int) (*models.Item, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, exists := s.items[id]
	if !exists {
		return nil, models.ErrMerchNotFound
	}
	return item, nil
}

func (s *MockMerchStorage) GetByName(ctx context.Context, name string) (*models.Item, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, item := range s.items {
		if item.Name == name {
			return item, nil
		}
	}
	return nil, models.ErrMerchNotFound
}

func (s *MockMerchStorage) GetList(ctx context.Context) ([]*models.Item, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	list := make([]*models.Item, 0, len(s.items))
	for _, item := range s.items {
		list = append(list, item)
	}
	return list, nil
}

func (s *MockMerchStorage) Update(ctx context.Context, merch *models.Item) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.items[merch.Id]; !exists {
		return models.ErrMerchNotFound
	}

	s.items[merch.Id] = merch
	return nil
}

func (s *MockMerchStorage) Delete(ctx context.Context, id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.items[id]; !exists {
		return models.ErrMerchNotFound
	}

	delete(s.items, id)
	return nil
}

// MockTransactionStorage реализация
type MockTransactionStorage struct {
	mu           sync.RWMutex
	transactions []models.TransactionEntry
}

func NewMockTransactionStorage() *MockTransactionStorage {
	return &MockTransactionStorage{
		transactions: make([]models.TransactionEntry, 0),
	}
}

func (s *MockTransactionStorage) Create(ctx context.Context, sender, recv *models.User, amount int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if sender.Coins < amount {
		return models.ErrInsufficientCoins
	}

	sender.Coins -= amount
	recv.Coins += amount

	s.transactions = append(s.transactions, models.TransactionEntry{
		Id:         len(s.transactions) + 1,
		SenderID:   sender.Id,
		ReceiverID: recv.Id,
		Amount:     amount,
	})
	return nil
}
