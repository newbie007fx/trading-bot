package helper

type Response struct {
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func ErrorResponse(code int, err error, message string) *Response {
	var errors interface{}
	if err != nil {
		if v, ok := err.(*ValidationError); ok {
			errors = v.GetErrorMsg()
		}

		if message == "" {
			message = err.Error()
		}
	}
	return &Response{
		Status:  "error",
		Code:    code,
		Message: message,
		Errors:  errors,
	}
}

func SuccessResponse(code int, data interface{}, message string) *Response {
	return &Response{
		Status:  "success",
		Code:    code,
		Message: message,
	}
}
