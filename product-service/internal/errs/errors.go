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
	ErrProductNotFound      = &InternalError{Err: errors.New("product not found")}
	ErrCategoryNotFound     = &InternalError{Err: errors.New("category not found")}
	ErrCategoryHasProducts  = &InternalError{Err: errors.New("category has products")}
	ErrCategoryAlreadyExist = &InternalError{Err: errors.New("category with that slug already exists")}
)

var (
	ErrProductNotFoundRequest      = &RequestError{Err: errors.New("product not found")}
	ErrCategoryNotFoundRequest     = &RequestError{Err: errors.New("category not found")}
	ErrCategoryHasProductsRequest  = &RequestError{Err: errors.New("category has products")}
	ErrCategoryAlreadyExistRequest = &RequestError{Err: errors.New("category with that slug already exists")}
	ErrInvalidRequestBody          = &RequestError{Err: errors.New("invalid request body")}
	ErrInvalidProductID            = &RequestError{Err: errors.New("invalid product id")}
	ErrInvalidCategoryID           = &RequestError{Err: errors.New("invalid category id")}
	ErrInternal                    = &RequestError{Err: errors.New("internal error")}
)

var ErrMap = map[error]*RequestError{
	ErrProductNotFound:      ErrProductNotFoundRequest,
	ErrCategoryNotFound:     ErrCategoryNotFoundRequest,
	ErrCategoryHasProducts:  ErrCategoryHasProductsRequest,
	ErrCategoryAlreadyExist: ErrCategoryAlreadyExistRequest,
}

func MapErr(err error) *RequestError {
	var internalErr *InternalError

	if errors.As(err, &internalErr) {
		return ErrMap[err]
	}

	return ErrInternal
}
