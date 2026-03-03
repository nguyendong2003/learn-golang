package apperror

type ErrorType string

const (
	NotFound       ErrorType = "not_found"
	BadRequest     ErrorType = "bad_request"
	InternalServer ErrorType = "internal_server"
	Unauthorized   ErrorType = "unauthorized"
	Forbidden      ErrorType = "forbidden"
	DuplicateEntry ErrorType = "duplicate_entry"
)

type AppError struct {
	Type    ErrorType `json:"type"`
	Message string    `json:"message"`
	Err     error     `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewNotFoundError(message string) *AppError {
	return &AppError{
		Type:    NotFound,
		Message: message,
	}
}

func NewBadRequestError(message string) *AppError {
	return &AppError{
		Type:    BadRequest,
		Message: message,
	}
}

func NewInternalServerError(message string, err error) *AppError {
	return &AppError{
		Type:    InternalServer,
		Message: message,
		Err:     err,
	}
}

func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Type:    Unauthorized,
		Message: message,
	}
}
func NewForbiddenError(message string) *AppError {
	return &AppError{
		Type:    Forbidden,
		Message: message,
	}
}

func NewDuplicateEntryError(message string) *AppError {
	return &AppError{
		Type:    DuplicateEntry,
		Message: message,
	}
}
