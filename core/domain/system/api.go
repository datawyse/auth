package system

type HttpResponse struct {
	Code    int                    `mapstructure:"code" json:"code"`
	Status  bool                   `mapstructure:"status" json:"status"`
	Message string                 `mapstructure:"message" json:"message"`
	Data    map[string]interface{} `mapstructure:"data" json:"data"`
}

func NewHttpResponse(status bool, message string, data map[string]any, code int) HttpResponse {
	response := HttpResponse{
		Code:    code,
		Status:  status,
		Message: message,
		Data:    data,
	}

	return response
}

// Error makes it compatible with the `error` interface.
func (e *HttpResponse) Error() string {
	return e.Message
}
