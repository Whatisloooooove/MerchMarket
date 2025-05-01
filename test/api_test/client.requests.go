package apitest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"merch_service/internal/models"
	"net/http"
)

// Register отправляет запрос на регистрацию по данным переданным в *Credentials
func (c *Client) Register(ctx context.Context, user *models.LoginRequest) (*ResponseBody, error) {
	jsonCreds, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	// В консоли делали бы так:
	// curl -X POST http://localhost:8080/auth/register \
	//      -H "Content-type: application/json" \
	//      -d '{"login":"aboba", "pass":"123"}'
	// Здесь аналогично
	req, err := http.NewRequest("POST",
		fmt.Sprintf("%s/auth/register", c.BaseURL),
		bytes.NewBuffer(jsonCreds))
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	// /auth/register не возвращает никаких данных, т.е. data: "{}"
	response, err := c.SendRequest(req, 1)

	if err != nil {
		return nil, err
	}

	return response, nil
}

// GetTokens в случае успешного запроса возвращает токены для дальнейших
// запросов и nil. Ошибку в ином случае
func (c *Client) GetTokens(ctx context.Context, user *models.LoginRequest) (*ResponseBody, error) {
	jsonCreds, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	// см. похожий блок в Register
	req, err := http.NewRequest("POST",
		fmt.Sprintf("%s/auth/login", c.BaseURL),
		bytes.NewBuffer(jsonCreds))
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	tokens := &UserTokens{}

	// Сравни с Register! Здесь мы передаем в SendRequest вторым аргументом *UserTokens
	response, err := c.SendRequest(req, tokens)

	if err != nil {
		return nil, err
	}

	return response, nil
}

// TODO authorization header
func (c *Client) Buy(ctx context.Context, purchReq *models.PurchaseRequest, tokens *UserTokens) (*ResponseBody, error) {
	purchReqBytes, err := json.Marshal(purchReq)
	if err != nil {
		return nil, err
	}

	// см. похожий блок в Register
	req, err := http.NewRequest("POST",
		fmt.Sprintf("%s/merch/buy", c.BaseURL),
		bytes.NewBuffer(purchReqBytes))
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	req.Header.Set("Authorization", tokens.Token)

	purchEntry := &PurchaseEntry{}
	response, err := c.SendRequest(req, purchEntry)

	if err != nil {
		return nil, err
	}

	return response, nil
}

// SendRequest функция отправки запроса. Так как все запросы к API имеют точку
// отправки запроса и одни и те же заголовки, для уменьшения повторения они вынесены.
// В v обычно передается указатель на структуры соотв. запросу. Например:
//
//	*UserTokens, *MerchList, *HistoryLog ...
func (c *Client) SendRequest(req *http.Request, v interface{}) (*ResponseBody, error) {
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		// При ошибке resp можно игнорировать (не закрывать)
		// см. документацию .Do
		return nil, err
	}

	defer resp.Body.Close()

	// Добавить обработку ошибок: resp.StatusCode ...

	// см. преамбулу в endpoints.go
	// декодер распарсит ответ сервера и запишет их в v (который обычно указатель)
	respBody := &ResponseBody{
		Data: v,
	}

	if err := json.NewDecoder(resp.Body).Decode(respBody); err != nil {
		return nil, err
	}

	return respBody, nil
}
