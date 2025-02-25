// Package response содержит структуры, описывающие ответы клиенту.
package response

import "net/http"

// DataOkResponse описывает структуру в состав которой входят данные для отправки клиенту в случае успеха.
type DataOkResponse[T any] struct {
	Status int `json:"status"`
	Data   *T  `json:"data,omitempty"`
}

// DataOkRsp функция, конструирующая структуру DataOkResponse.
func DataOkRsp[T any](data ...T) DataOkResponse[T] {

	return DataOkResponse[T]{Status: http.StatusOK, Data: &data[0]}

}
