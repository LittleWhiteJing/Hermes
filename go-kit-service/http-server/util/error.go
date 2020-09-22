package util

type AppError struct {
	Code 	int
	Message string
}

func NewAppError(code int, message string) error {
	return &AppError{
		Code: code,
		Message: message,
	}
}

func (e *AppError) Error() string {
	return e.Message
}

