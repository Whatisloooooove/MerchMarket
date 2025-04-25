package mock

import (
	"context"
	"merch_service/internal/models"
	"merch_service/internal/storage/entites"
	"sync"
	"time"
)

var (
	_ entites.MerchStorage       = (*MockMerchStorage)(nil)
	_ entites.UserStorage        = (*MockUserStorage)(nil)
	_ entites.TransactionStorage = (*MockTransactionStorage)(nil)
)

// MockUserStorage реализация
type MockUserStorage struct {
	mu    sync.RWMutex
	users map[string]*models.User
	coins map[string][]models.CoinsEntry
	purch map[string][]models.PurchaseEntry
}

func NewMockUserStorage() *MockUserStorage {
	return &MockUserStorage{
		users: make(map[string]*models.User),
		coins: make(map[string][]models.CoinsEntry),
		purch: make(map[string][]models.PurchaseEntry),
	}
}

func (s *MockUserStorage) Create(ctx context.Context, user *models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[user.Login]; exists {
		return models.ErrUserExists
	}

	s.users[user.Login] = user
	s.coins[user.Login] = []models.CoinsEntry{{
		Date:        time.Now(),
		CoinsBefore: 0,
		CoinsAfter:  1000,
	}}
	return nil
}

func (s *MockUserStorage) Get(ctx context.Context, login string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[login]
	if !exists {
		return nil, models.ErrUserExists
	}
	return user, nil
}

func (s *MockUserStorage) Update(ctx context.Context, user *models.User) error {
	return nil
}

func (s *MockUserStorage) Delete(ctx context.Context, user *models.User) error {
	return nil
}

func (s *MockUserStorage) GetCoinsHistory(ctx context.Context, user *models.User) ([]models.CoinsEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	history, exists := s.coins[user.Login]
	if !exists {
		return nil, models.ErrUserExists
	}
	return history, nil
}

func (s *MockUserStorage) GetPurchaseHistory(ctx context.Context, user *models.User) ([]models.PurchaseEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	history, exists := s.purch[user.Login]
	if !exists {
		return nil, models.ErrUserExists
	}
	return history, nil
}

// MockMerchStorage реализация
type MockMerchStorage struct {
	mu    sync.RWMutex
	items map[string]*models.Item
}

func NewMockMerchStorage() *MockMerchStorage {
	return &MockMerchStorage{
		items: map[string]*models.Item{
			"Футболка": {Name: "Футболка", Price: 50, Stock: 10},
			"Кружка":   {Name: "Кружка", Price: 30, Stock: 5},
		},
	}
}

func (s *MockMerchStorage) Get(ctx context.Context, name string) (*models.Item, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, exists := s.items[name]
	if !exists {
		return nil, models.ErrNoMerchInStock
	}
	return item, nil
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

func (s *MockMerchStorage) Update(ctx context.Context, user *models.User, merch *models.Item) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, exists := s.items[merch.Name]
	if !exists {
		return models.ErrNoMerchInStock
	}

	user.Coins -= item.Price * merch.Stock
	item.Stock -= merch.Stock

	return nil
}

func (s *MockMerchStorage) Create(ctx context.Context, merch *models.Item) error { return nil }
func (s *MockMerchStorage) Delete(ctx context.Context, name string) error        { return nil }

// MockTransactionStorage реализация
type MockTransactionStorage struct {
	mu           sync.RWMutex
	transactions []models.TransactionEntry
}

func NewMockTransactionStorage() *MockTransactionStorage {
	return &MockTransactionStorage{}
}

func (s *MockTransactionStorage) Create(ctx context.Context, sender, recv *models.User, amount int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	sender.Coins -= amount
	recv.Coins += amount

	s.transactions = append(s.transactions, models.TransactionEntry{
		Sender:   sender.Login,
		Reciever: recv.Login,
		Amount:   amount,
	})
	return nil
}
