package response

import "net/http"

// EmptyOkResponse описывает структуру в составе которой отсутствуют данные для отправки клиенту в случае успеха.
type EmptyOkResponse struct {
	Status int `json:"status"`
}

// EmptyOkRsp функция, конструирующая структуру EmptyOkResponse.
func EmptyOkRsp() EmptyOkResponse {
	return EmptyOkResponse{Status: http.StatusOK}
}
