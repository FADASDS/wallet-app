package response

import "net/http"

type DataOkResponse[T any] struct {
	Status int `json:"status"`
	Data   *T  `json:"data,omitempty"`
}

func DataOkRsp[T any](data ...T) DataOkResponse[T] {

	return DataOkResponse[T]{Status: http.StatusOK, Data: &data[0]}

}
