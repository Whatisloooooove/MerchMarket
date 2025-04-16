package apitest

import (
	"bytes"
	"encoding/json"
	"merch_service/internal/server"
	"net/http"
	"net/http/httptest"

	"testing"
)

type TokenData struct {
	Token   string `json:"token"`
	Refresh string `json:"refresh"`
}

type TokenResponse struct {
	ErrorCode int       `json:"error_code"`
	Message   string    `json:"message"`
	Data      TokenData `json:"data"`
}

// Тест на успешное возвращение запросов при регистрации
func TestRegisterHandler(t *testing.T) {
	t.Helper()

	serv := server.NewServer()
	serv.SetupRoutes()

	ts := httptest.NewServer(serv.Router())
	defer ts.Close()

	regReq := map[string]string{
		"login":    "testReg@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(regReq)

	resp, err := http.Post(ts.URL+"/auth/register", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Запрос не выполнился: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Ожидалось 201, получили %d", resp.StatusCode)
	}
}

// Тест на успешное возвращение запросов при аутентификации
func TestAuthHandler(t *testing.T) {
	t.Helper()

	serv := server.NewServer()
	serv.SetupRoutes()

	ts := httptest.NewServer(serv.Router())
	defer ts.Close()

	// Создаём нового польщователя
	regReq := map[string]string{
		"login":    "testAuth@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(regReq)

	resp, err := http.Post(ts.URL+"/auth/register", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Запрос на регистрацию не выполнился: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Ожидалось 201, получили %d", resp.StatusCode)
	}

	// Пытаемся авторизоваться
	loginBody := map[string]string{
		"login":    "testAuth@example.com",
		"password": "password123",
	}
	loginBytes, _ := json.Marshal(loginBody)

	loginResp, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewBuffer(loginBytes))
	if err != nil {
		t.Fatalf("Запрос на аутентификацию не выполнился: %v", err)
	}
	defer loginResp.Body.Close()

	if loginResp.StatusCode != http.StatusOK {
		t.Fatalf("Ожидался статус 200, получили %d", loginResp.StatusCode)
	}

	// Проверяем, что нам вернулся токен
	var loginResult TokenResponse
	err = json.NewDecoder(loginResp.Body).Decode(&loginResult)
	if err != nil {
		t.Fatalf("Не удалось декодировать данные: %v", err)
	}

	if loginResult.Data.Token == "" {
		t.Errorf("Ожидался access token, но он пустой")
	}
	if loginResult.Data.Refresh == "" {
		t.Errorf("Ожидался refresh token, но он пустой")
	}
}
