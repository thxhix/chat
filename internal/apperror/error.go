package apperror

type Code string

const (
	CodeBadRequest   Code = "BAD_REQUEST"
	CodeConflict     Code = "CONFLICT"
	CodeNotFound     Code = "NOT_FOUND"
	CodeUnauthorized Code = "UNAUTHORIZED"
	CodeForbidden    Code = "FORBIDDEN"
	CodeTimeout      Code = "Timeout"
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

func NewInternalError(err error) *Error {
	return &Error{
		Code:    CodeInternal,
		Message: "Internal Server Error",
		Err:     err,
	}
}

func NewBadRequestError(msg string) *Error {
	return &Error{
		Code:    CodeBadRequest,
		Message: msg,
	}
}

func NewUnauthorizedError(msg string) *Error {
	return &Error{
		Code:    CodeUnauthorized,
		Message: msg,
	}
}

func NewConflictError(msg string) *Error {
	return &Error{
		Code:    CodeConflict,
		Message: msg,
	}
}

func NewNotFoundError(msg string) *Error {
	return &Error{
		Code:    CodeNotFound,
		Message: msg,
	}
}

func NewForbiddenError(msg string) *Error {
	return &Error{
		Code:    CodeForbidden,
		Message: msg,
	}
}

func NewTimeoutError(msg string) *Error {
	return &Error{
		Code:    CodeTimeout,
		Message: msg,
	}
}
