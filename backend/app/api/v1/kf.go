package v1

import (
	"MKICS/backend/app/dto/request"
	"MKICS/backend/app/dto/response"

	"github.com/gin-gonic/gin"
)

// @Tags kf
// @Summary Add kf
// @Description Add a new kf
// @Accept json
// @Produce json
// @Param req body request.KF true "kf Information"
// @Success 200 {object} dto.Response
// @Router /kf [post]
func (u *BaseApi) KFAdd(ctx *gin.Context) {
	var req request.KF
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	if err := kFLogic.KFAdd(req); err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "success")
}

// @Tags kf
// @Summary Update kf
// @Description Update the kf by UUID
// @Accept json
// @Produce json
// @Param uuid path string true "kf UUID"
// @Param req body request.KF true "kf Information"
// @Success 200 {object} dto.Response
// @Router /kf/{uuid} [patch]
func (u *BaseApi) KFUpdate(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	var req request.KF
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	if err := kFLogic.KFUpdate(uuid, req); err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "success")
}

// @Tags kf
// @Summary Delete KF
// @Description Delete the kf by UUID
// @Accept json
// @Produce json
// @Param uuid path string true "kf UUID"
// @Success 200 {object} dto.Response
// @Router /kf/{uuid} [delete]
func (u *BaseApi) KFDel(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	if err := kFLogic.KFDel(uuid); err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "success")
}

// @Tags kf
// @Summary Get KF
// @Description Get the kf information by UUID
// @Accept json
// @Produce json
// @Param uuid path string true "kf UUID"
// @Success 200 {object} dto.Response
// @Router /kf/{uuid} [get]
func (u *BaseApi) KFGet(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	data, err := kFLogic.KFGet(uuid)
	if err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}

	response.SuccessWithData(ctx, data)
}

// @Tags kf
// @Summary List KF
// @Description Retrieve a list of all kf
// @Accept json
// @Produce json
// @Success 200 {object} dto.Response
// @Router /kf [get]
func (u *BaseApi) KFList(ctx *gin.Context) {
	data, err := kFLogic.KFList()
	if err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}

	response.SuccessWithData(ctx, data)
}
