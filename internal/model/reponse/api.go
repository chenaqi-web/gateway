package reponse

type APIResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func Success(data any) APIResponse {
	return APIResponse{Code: 0, Message: "success", Data: data}
}

func Error(code int, message string) APIResponse {
	if message == "" {
		message = "error"
	}
	return APIResponse{Code: code, Message: message}
}
