package v1

import (
	"MKICS/backend/app/dto/request"
	"MKICS/backend/app/dto/response"

	"github.com/gin-gonic/gin"
)

// @Tags staff
// @Summary Add Staff
// @Description Add a new staff member
// @Accept json
// @Produce json
// @Param req body request.Staff true "Staff Information"
// @Success 200 {object} dto.Response
// @Router /staff [post]
func (u *BaseApi) StaffAdd(ctx *gin.Context) {
	var req request.Staff
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	err := staffLogic.StaffAdd(req)
	if err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "success")
}

// @Tags staff
// @Summary Update Staff
// @Description Update the staff member by UUID
// @Accept json
// @Produce json
// @Param uuid path string true "Staff UUID"
// @Param req body request.Staff true "Staff Information"
// @Success 200 {object} dto.Response
// @Router /staff/{uuid} [patch]
func (u *BaseApi) StaffUpdate(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	var req request.Staff
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	if err := staffLogic.StaffUpdate(uuid, req); err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "success")
}

// @Tags staff
// @Summary Delete Staff
// @Description Delete the staff member by UUID
// @Accept json
// @Produce json
// @Param uuid path string true "Staff UUID"
// @Success 200 {object} dto.Response
// @Router /staff/{uuid} [delete]
func (u *BaseApi) StaffDel(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	if err := staffLogic.StaffDel(uuid); err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "success")
}

// @Tags staff
// @Summary Get Staff
// @Description Get the staff information by UUID
// @Accept json
// @Produce json
// @Param uuid path string true "Staff UUID"
// @Success 200 {object} dto.Response
// @Router /staff/{uuid} [get]
func (u *BaseApi) StaffGet(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	data, err := staffLogic.StaffGet(uuid)
	if err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}

	response.SuccessWithData(ctx, data)
}

// @Tags staff
// @Summary List Staff
// @Description Retrieve a list of all staff members
// @Accept json
// @Produce json
// @Success 200 {object} dto.Response
// @Router /staff [get]
func (u *BaseApi) StaffList(ctx *gin.Context) {
	data, err := staffLogic.StaffList()
	if err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}

	response.SuccessWithData(ctx, data)
}
