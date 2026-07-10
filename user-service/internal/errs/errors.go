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
	ErrUserAlreadyExists = &InternalError{Err: errors.New("user with that email already exists")}
	ErrUserNotFound      = &InternalError{Err: errors.New("user not found")}
)

var (
	ErrUserWithThatEmailAlreadyExists = &RequestError{Err: errors.New("user with that email already exists")}
	ErrInvalidCredentials             = &RequestError{Err: errors.New("invalid credentials")}
	ErrInternal                       = &RequestError{Err: errors.New("internal error")}
	ErrInvalidUserID                  = &RequestError{Err: errors.New("invalid user id")}
	ErrInvalidRequestBody             = &RequestError{Err: errors.New("invalid request body")}
	ErrAccessDenied                   = &RequestError{Err: errors.New("access denied")}
	ErrInvalidToken                   = &RequestError{Err: errors.New("invalid token or expired token")}
)

var ErrMap = map[*InternalError]*RequestError{
	ErrUserAlreadyExists: ErrUserWithThatEmailAlreadyExists,
	ErrUserNotFound:      ErrInvalidCredentials,
}

func MapErr(err error) *RequestError {
	for repoErr, reqErr := range ErrMap {
		if errors.Is(err, repoErr) {
			return reqErr
		}
	}

	return ErrInternal
}
