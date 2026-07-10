package http

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Response struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

func newResponse(data interface{}, msg ...string) *Response {
	v, ok := data.(error)
	if ok {
		data = v.Error()
	}

	resp := &Response{
		Data: data,
	}

	if len(msg) > 0 {
		resp.Message = msg[0]
	}

	return resp
}

func OK(ctx *gin.Context, data interface{}, msg ...string) {
	resp := newResponse(data, msg...)
	ctx.JSON(http.StatusOK, resp)
}

func BadRequest(ctx *gin.Context, data interface{}) {
	resp := newResponse(data)
	ctx.JSON(http.StatusBadRequest, resp)
}

func Forbidden(ctx *gin.Context, data interface{}) {
	resp := newResponse(data)
	ctx.JSON(http.StatusForbidden, resp)
}

func InternalServerErr(ctx *gin.Context, data interface{}) {
	resp := newResponse(data)
	ctx.JSON(http.StatusInternalServerError, resp)
}

func ValidationFailed(ctx *gin.Context, err error) {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		InternalServerErr(ctx, err)

		return
	}

	errMsgs := make([]string, 0, len(validationErrors))
	for _, fieldErr := range validationErrors {
		var msg string

		switch fieldErr.Tag() {
		case "required":
			msg = fmt.Sprintf("field '%s' is required", fieldErr.Field())
		case "email":
			msg = fmt.Sprintf("field '%s' must be a valid email address", fieldErr.Field())
		case "min":
			msg = fmt.Sprintf("field '%s' must be at least %s characters long", fieldErr.Field(), fieldErr.Param())
		case "max":
			msg = fmt.Sprintf("field '%s' must be at most %s characters long", fieldErr.Field(), fieldErr.Param())
		case "oneof":
			msg = fmt.Sprintf("field '%s' must be one of: %s", fieldErr.Field(), fieldErr.Param())
		default:
			msg = fmt.Sprintf("field '%s' failed validation: %s", fieldErr.Field(), fieldErr.Tag())
		}

		errMsgs = append(errMsgs, msg)
	}

	resp := newResponse(errMsgs, "validation failed")
	ctx.JSON(http.StatusUnprocessableEntity, resp)
}
