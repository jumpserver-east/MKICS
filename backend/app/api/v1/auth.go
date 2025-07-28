package v1

import (
	"MKICS/backend/app/dto/request"
	"MKICS/backend/app/dto/response"
	"MKICS/backend/constant"

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
		response.BadRequest(ctx, err)
		return
	}

	token, err := authLogic.Login(ctx, req)
	if err != nil {
		response.InternalServerError(ctx, err)
		return
	}

	response.SuccessWithData(ctx, token)
}

// @Tags auth
// @Summary User logout
// @Description
// @Accept json
// @Success 200 {object} dto.Response
// @Router /auth/logout [post]
func (b *BaseApi) Logout(ctx *gin.Context) {
	_, exists := ctx.Get(constant.UserKey)
	if !exists {
		response.Unauthorized(ctx, constant.ErrUnauthorizedOrNotLoggedIn)
		return
	}

	err := authLogic.Logout(ctx)
	if err != nil {
		response.InternalServerError(ctx, err)
		return
	}

	response.SuccessWithMsg(ctx, "success")
}
