package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SuccessWithDetail(ctx *gin.Context, msg string, data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	res := Response{
		Code:    http.StatusOK,
		Message: msg,
		Data:    data,
	}
	ctx.JSON(http.StatusOK, res)
	ctx.Abort()
}

func SuccessWithMsg(ctx *gin.Context, msg string) {
	SuccessWithDetail(ctx, msg, nil)
}

func SuccessWithData(ctx *gin.Context, data interface{}) {
	SuccessWithDetail(ctx, "Success", data)
}

func ErrorWithDetail(ctx *gin.Context, code int, msg string) {
	res := Response{
		Code:    code,
		Message: msg,
	}
	ctx.JSON(code, res)
	ctx.Abort()
}

func InternalServerError(ctx *gin.Context, err error) {
	ErrorWithDetail(ctx, http.StatusInternalServerError, err.Error())
}

func BadRequest(ctx *gin.Context, err error) {
	ErrorWithDetail(ctx, http.StatusBadRequest, err.Error())
}

func Unauthorized(ctx *gin.Context, err error) {
	ErrorWithDetail(ctx, http.StatusUnauthorized, err.Error())
}

func NotFound(ctx *gin.Context, err error) {
	ErrorWithDetail(ctx, http.StatusNotFound, err.Error())
}

func ResponseWithCode(ctx *gin.Context, code int) {
	ctx.JSON(code, nil)
	ctx.Abort()
}

func Success(ctx *gin.Context) {
	ResponseWithCode(ctx, http.StatusOK)
}
