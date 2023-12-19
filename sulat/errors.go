package sulat

type ResponseError struct {
	StatusCode int
	Message    string
}

func (err *ResponseError) Error() string {
	return err.Message
}

func NewResponseError(statusCode int, message string) *ResponseError {
	return &ResponseError{
		StatusCode: statusCode,
		Message:    message,
	}
}
