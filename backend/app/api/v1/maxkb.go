package v1

import (
	"EvoBot/backend/app/api/v1/helper"
	"EvoBot/backend/app/dto"
	"EvoBot/backend/constant"

	"github.com/gin-gonic/gin"
)

// @Tags maxkb_config
// @Summary Update maxkb configuration
// @Description Update the maxkb configuration with the provided data
// @Accept json
// @Produce json
// @Param data body dto.MaxkbConf true "maxkb configuration data"
// @Success 200 {object} dto.Response "Success response"
// @Router /maxkb/config [patch]
func (b *BaseApi) UpdateMaxkbConfig(ctx *gin.Context) {
	var req dto.MaxkbConf
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helper.ErrResponse(ctx, constant.CodeErrBadRequest)
		return
	}
	if err := maxkbLogic.UpdateConfig(ctx, req); err != nil {
		helper.ErrResponse(ctx, constant.CodeErrInternalServer)
		return
	}
	helper.SuccessWithOutData(ctx)
}

// @Tags maxkb_config
// @Summary Get maxkb configuration
// @Description Retrieve the current maxkb configuration
// @Accept json
// @Produce json
// @Success 200 {object} dto.MaxkbConf "Current maxkb configuration"
// @Router /maxkb/config [get]
func (b *BaseApi) ListMaxkbConfig(ctx *gin.Context) {
	conf, err := maxkbLogic.GetConfig(ctx)
	if err != nil {
		helper.ErrResponse(ctx, constant.CodeErrInternalServer)
		return
	}
	helper.SuccessWithData(ctx, conf)
}
