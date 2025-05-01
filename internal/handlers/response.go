package handlers

import "net/http"

type GeneralResponse struct {
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
	Data      any    `json:"data"`
}

func DefaultResponse() GeneralResponse {
	return GeneralResponse{
		ErrorCode: http.StatusInternalServerError,
		Message:   InternalServerError,
		Data:      struct{}{},
	}
}
