package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Register отправляет запрос на регистрацию по данным переданным в *Credentials
func (c *Client) Register(ctx context.Context, user *Credentials) error {
	jsonCreds, err := json.Marshal(user)
	if err != nil {
		return err
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
		return err
	}
	req = req.WithContext(ctx)

	// /auth/register не возвращает никаких данных, т.е. data: "{}"
	// поэтому запишем ничего в никуда)
	err = c.SendRequest(req, 1) // TODO??? 1 =) Не спрашивайте почему так работает. Надо разобраться
	// UPD: Все нормально. 1 перезапишется на map[string]interface{} так как это соотв. пустым {}

	return err
}

// TODO !!! Добавить хедер Authorization
// GetTokens в случае успешного запроса возвращает токены для дальнейших
// запросов и nil. Ошибку в ином случае
func (c *Client) GetTokens(ctx context.Context, user *Credentials) (*UserTokens, error) {
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
	err = c.SendRequest(req, &tokens)

	if err != nil {
		return nil, err
	}
	return tokens, nil
}

// SendRequest функция отправки запроса. Так как все запросы к API имеют точку
// отправки запроса и одни и те же заголовки, для уменьшения повторения они вынесены.
// В v обычно передается указатель на структуры соотв. запросу. Например:
//
//	*UserTokens, *MerchList, *HistoryLog ...
func (c *Client) SendRequest(req *http.Request, v interface{}) error {
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		// При ошибке resp можно игнорировать (не закрывать)
		// см. документацию .Do
		return err
	}

	defer resp.Body.Close()

	// Добавить обработку ошибок: resp.StatusCode ...

	// см. преамбулу в endpoints.go
	// декодер распарсит ответ сервера и запишет их в v (который обычно указатель)
	respBody := &ResponseBody{
		Data: v,
	}
	if err := json.NewDecoder(resp.Body).Decode(respBody); err != nil {
		return err
	}

	log.Println("Данные ответа", *respBody) // Remove after DEBUG!!!
	return nil
}
