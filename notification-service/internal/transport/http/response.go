package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
