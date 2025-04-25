package client

import (
	"net/http"
)

const (
	BaseURL = "http://localhost:8080"
)

type Client struct {
	BaseURL      string
	token        string
	refreshToken string
	HTTPClient   *http.Client
}

type Credentials struct {
	Login string `json:"login"`
	Pass  string `json:"pass"`
}

func NewClient() *Client {
	return &Client{
		BaseURL:    BaseURL,
		HTTPClient: &http.Client{},
	}
}
