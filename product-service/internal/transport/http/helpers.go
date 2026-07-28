package http

import (
	"errors"

	"github.com/anomalyco/hookah-store/product-service/internal/errs"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var CodeErrMap = map[error]func(ctx *gin.Context, data interface{}){
	errs.ErrInternal: InternalServerErr,

	errs.ErrProductNotFoundRequest:      BadRequest,
	errs.ErrCategoryNotFoundRequest:     BadRequest,
	errs.ErrCategoryHasProductsRequest:  BadRequest,
	errs.ErrCategoryAlreadyExistRequest: BadRequest,
	errs.ErrInvalidRequestBody:          BadRequest,
	errs.ErrInvalidProductID:            BadRequest,
	errs.ErrInvalidCategoryID:           BadRequest,
}

func HandleServiceError(ctx *gin.Context, err error) {
	if _, ok := errors.AsType[validator.ValidationErrors](err); ok {
		ValidationFailed(ctx, err)

		return
	}

	fn, ok := CodeErrMap[err]
	if !ok {
		InternalServerErr(ctx, err)

		return
	}

	fn(ctx, err)
}
