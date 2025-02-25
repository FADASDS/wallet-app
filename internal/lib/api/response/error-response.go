package response

// ErrorResponse описывает структуру отправки клиенту в случае ошибки.
type ErrorResponse struct {
	Status int    `json:"status"`
	Error  string `json:"error,omitempty"`
}

// Error функция, конструирующая структуру ErrorResponse.
func Error(status int, msg string) ErrorResponse {
	return ErrorResponse{
		Status: status,
		Error:  msg}
}
