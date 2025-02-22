package response

import "net/http"

type EmptyOkResponse struct {
	Status int `json:"status"`
}

func EmptyOkRsp() EmptyOkResponse {
	return EmptyOkResponse{Status: http.StatusOK}
}
