package api

type ErrorResponse struct {
	Status int    `json:"status"`
	Error  string `json:"error,omitempty"`
}

func Error(status int, msg string) ErrorResponse {
	return ErrorResponse{
		Status: status,
		Error:  msg}
}
