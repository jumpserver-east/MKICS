package v1

import (
	"MKICS/backend/app/dto"
	"MKICS/backend/app/dto/response"

	"github.com/gin-gonic/gin"
)

// @Tags policy
// @Summary Add policy
// @Description Add a new policy
// @Accept json
// @Produce json
// @Param req body dto.Policy true "policy Information"
// @Success 200 {object} dto.Response
// @Router /policy [post]
func (u *BaseApi) PolicyAdd(ctx *gin.Context) {
	var req dto.Policy
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, err)
		return
	}

	if err := policyLogic.PolicyAdd(req); err != nil {
		response.InternalServerError(ctx, err)
		return
	}

	response.SuccessWithMsg(ctx, "success")
}

// @Tags policy
// @Summary Update policy
// @Description Update the policy by UUID
// @Accept json
// @Produce json
// @Param uuid path string true "policy UUID"
// @Param req body dto.Policy true "policy Information"
// @Success 200 {object} dto.Response
// @Router /policy/{uuid} [patch]
func (u *BaseApi) PolicyUpdate(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	var req dto.Policy
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, err)
		return
	}

	if err := policyLogic.PolicyUpdate(uuid, req); err != nil {
		response.InternalServerError(ctx, err)
		return
	}

	response.SuccessWithMsg(ctx, "success")
}

// @Tags policy
// @Summary Delete policy
// @Description Delete the policy by UUID
// @Accept json
// @Produce json
// @Param uuid path string true "policy UUID"
// @Success 200 {object} dto.Response
// @Router /policy/{uuid} [delete]
func (u *BaseApi) PolicyDel(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	if err := policyLogic.PolicyDel(uuid); err != nil {
		response.InternalServerError(ctx, err)
		return
	}

	response.SuccessWithMsg(ctx, "success")
}

// @Tags policy
// @Summary Get policy
// @Description Get the policy information by UUID
// @Accept json
// @Produce json
// @Param uuid path string true "policy UUID"
// @Success 200 {object} dto.Response
// @Router /policy/{uuid} [get]
func (u *BaseApi) PolicyGet(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	data, err := policyLogic.PolicyGet(uuid)
	if err != nil {
		response.InternalServerError(ctx, err)
		return
	}

	response.SuccessWithData(ctx, data)
}

// PolicyList retrieves all policy
// @Tags policy
// @Summary List policy
// @Description Retrieve a list of all policy
// @Accept json
// @Produce json
// @Success 200 {object} dto.Response
// @Router /policy [get]
func (u *BaseApi) PolicyList(ctx *gin.Context) {
	data, err := policyLogic.PolicyList()
	if err != nil {
		response.InternalServerError(ctx, err)
		return
	}

	response.SuccessWithData(ctx, data)
}
