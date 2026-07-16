package http

import (
	"github.com/gin-gonic/gin"
)

func HandleServiceError(ctx *gin.Context, err error) {
	InternalServerErr(ctx, err)
}
