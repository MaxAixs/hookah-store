package errs

import "errors"

type InternalError struct {
	Err error
}

func (e *InternalError) Error() string {
	return e.Err.Error()
}

type RequestError struct {
	Err error
}

func (e *RequestError) Error() string {
	return e.Err.Error()
}

var (
	ErrUserNotFound     = &InternalError{Err: errors.New("user not found")}
	ErrEmailNotFound    = &InternalError{Err: errors.New("email not found")}
	ErrInvalidSignature = &InternalError{Err: errors.New("invalid signature")}
)

var (
	ErrInvalidCredentials = &RequestError{Err: errors.New("invalid credentials")}
	ErrInternal           = &RequestError{Err: errors.New("internal error")}
)

var mapErrors = map[error]*RequestError{
	ErrUserNotFound:  ErrInvalidCredentials,
	ErrEmailNotFound: ErrInvalidCredentials,
}

func MapErr(err error) error {
	var internalErr *InternalError

	if errors.As(err, &internalErr) {
		return mapErrors[err]
	}

	return ErrInternal
}
