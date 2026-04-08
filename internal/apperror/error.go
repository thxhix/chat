package apperror

type Code string

const (
	CodeValidation   Code = "VALIDATION_ERROR"
	CodeConflict     Code = "CONFLICT"
	CodeNotFound     Code = "NOT_FOUND"
	CodeUnauthorized Code = "UNAUTHORIZED"
	CodeInternal     Code = "INTERNAL"
)

type Error struct {
	Code    Code
	Message string
	Err     error
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

func NewValidationError(msg string) error {
	return &Error{
		Code:    CodeValidation,
		Message: msg,
	}
}

func NewUnauthorizedError(msg string) error {
	return &Error{
		Code:    CodeUnauthorized,
		Message: msg,
	}
}

func NewConflictError(msg string) error {
	return &Error{
		Code:    CodeConflict,
		Message: msg,
	}
}

func NewNotFoundError(msg string) error {
	return &Error{
		Code:    CodeNotFound,
		Message: msg,
	}
}
