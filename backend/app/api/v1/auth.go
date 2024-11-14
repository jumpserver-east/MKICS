package v1

import (
	"EvoBot/backend/app/api/v1/helper"
	"EvoBot/backend/app/dto/request"
	"EvoBot/backend/constant"

	"github.com/gin-gonic/gin"
)

type BaseApi struct{}

// @Tags auth
// @Summary User login
// @Description
// @Accept json
// @Param data body request.Login true "JSON"
// @Success 200 {object} dto.Response
// @Router /auth/login [post]
func (b *BaseApi) Login(ctx *gin.Context) {
	var req request.Login
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helper.ErrResponse(ctx, constant.CodeErrBadRequest)
		return
	}
	tokens, err := authLogic.Login(ctx, req)
	if err != nil {
		helper.ErrResponse(ctx, constant.CodeErrInternalServer)
		return
	}
	helper.SuccessWithData(ctx, tokens)
}
