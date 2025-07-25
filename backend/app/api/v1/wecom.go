package v1

import (
	"MKICS/backend/app/dto"
	"MKICS/backend/app/dto/request"
	"MKICS/backend/app/dto/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Tags wecom_config
// @Summary List Wecom configuration
// @Description Retrieve the current Wecom configuration
// @Accept json
// @Produce json
// @Success 200 {object} dto.Response "Current Wecom configuration"
// @Router /wecom/config [get]
func (b *BaseApi) WecomConfigList(ctx *gin.Context) {
	data, err := wecomLogic.ConfigList()
	if err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}

	response.SuccessWithData(ctx, data)
}

// @Tags wecom_config
// @Summary Get Wecom configuration
// @Description Retrieve the current Wecom configuration
// @Accept json
// @Produce json
// @Success 200 {object} dto.Response "Current Wecom configuration"
// @Router /wecom/config/{uuid} [get]
func (b *BaseApi) WecomConfigGet(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	data, err := wecomLogic.ConfigGet(uuid)
	if err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}

	response.SuccessWithData(ctx, data)
}

// @Tags wecom_config
// @Summary Update Wecom configuration
// @Description Update the Wecom configuration with the provided data
// @Accept json
// @Produce json
// @Param data body request.WecomConfigApp true "Wecom configuration data"
// @Success 200 {object} dto.Response "Success response"
// @Router /wecom/config/{uuid} [patch]
func (b *BaseApi) WecomConfigUpdate(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	var req request.WecomConfigApp
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	if err := wecomLogic.ConfigUpdate(uuid, req); err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "success")
}

// @Tags wecom_callback
// @Summary WeCom Callback Message Handling
// @Description Handles callback messages from WeCom (Enterprise WeChat). This endpoint receives and parses the callback data sent by WeCom.
// @Accept json
// @Produce json
// @Success 200 {object} string "Parsed WeCom callback message"
// @Router /wecom/callback [get,post]
func (b *BaseApi) WecomHandle(ctx *gin.Context) {
	// https://developer.work.weixin.qq.com/document/path/90930#32-%E6%94%AF%E6%8C%81http-post%E8%AF%B7%E6%B1%82%E6%8E%A5%E6%94%B6%E4%B8%9A%E5%8A%A1%E6%95%B0%E6%8D%AE
	// 正确响应本次请求
	// 企业微信服务器在五秒内收不到响应会断掉连接，并且重新发起请求，总共重试三次。
	// 仅针对网络连接失败或者网络请求超时情况重试，建议开发者接受回调后立即应答，业务异步处理。
	if ctx.Request.Method == http.MethodGet {
		var options dto.SignatureOptions
		if err := ctx.ShouldBindQuery(&options); err != nil {
			response.BadRequest(ctx, err.Error())
			return
		}

		echo, err := wecomLogic.VerifyURL(options)
		if err != nil {
			response.InternalServerError(ctx, err.Error())
			return
		}

		ctx.String(http.StatusOK, echo)
		ctx.Abort()
		return
	}

	body, err := ctx.GetRawData()
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	wecomLogic.CallbackHandler(body)

	response.SuccessWithMsg(ctx, "success")
}

// @Tags wecom_receptionist
// @Summary WeCom receptionist add
// @Description receptionist add
// @Accept json
// @Produce json
// @Param data body request.ReceptionistOptions true "Wecom configuration data"
// @Success 200
// @Router /wecom/receptionist/{kfid} [post]
func (b *BaseApi) WecomReceptionistAdd(ctx *gin.Context) {
	kfid := ctx.Param("kfid")

	var req request.ReceptionistOptions
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	if err := wecomLogic.ReceptionistAdd(kfid, req); err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "success")
}

// @Tags wecom_receptionist
// @Summary WeCom receptionist delete
// @Description receptionist delete
// @Accept json
// @Produce json
// @Param data body request.ReceptionistOptions true "Wecom configuration data"
// @Success 200
// @Router /wecom/receptionist/{kfid} [delete]
func (b *BaseApi) WecomReceptionistDel(ctx *gin.Context) {
	kfid := ctx.Param("kfid")

	var req request.ReceptionistOptions
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	if err := wecomLogic.ReceptionistDel(kfid, req); err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "success")
}

// @Tags wecom_receptionist
// @Summary WeCom receptionist list
// @Description receptionist list
// @Accept json
// @Produce json
// @Success 200 {object} []client.ReceptionistList
// @Router /wecom/receptionist/{kfid} [get]
func (b *BaseApi) WecomReceptionistList(ctx *gin.Context) {
	kfid := ctx.Param("kfid")

	data, err := wecomLogic.ReceptionistList(kfid)
	if err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}

	response.SuccessWithData(ctx, data)
}

// @Tags wecom_account
// @Summary WeCom account list
// @Description account list
// @Accept json
// @Produce json
// @Success 200 {object} []client.AccountInfoSchema
// @Router /wecom/account [get]
func (b *BaseApi) WecomAccountList(ctx *gin.Context) {
	data, err := wecomLogic.AccountList()
	if err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}

	response.SuccessWithData(ctx, data)
}

// @Tags wecom_account
// @Summary WeCom account get
// @Description account get
// @Accept json
// @Produce json
// @Success 200 {object} string
// @Router /wecom/account/{kfid} [get]
func (b *BaseApi) WecomAddContactWay(ctx *gin.Context) {
	kfid := ctx.Param("kfid")

	data, err := wecomLogic.AddContactWay(kfid)
	if err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}

	response.SuccessWithData(ctx, data)
}
