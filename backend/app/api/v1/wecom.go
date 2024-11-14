package v1

import (
	"EvoBot/backend/app/api/v1/helper"
	"EvoBot/backend/app/dto"
	"EvoBot/backend/app/dto/request"
	"EvoBot/backend/constant"
	"EvoBot/backend/global"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @Tags wecom_config
// @Summary Get Wecom configuration
// @Description Retrieve the current Wecom configuration
// @Accept json
// @Produce json
// @Success 200 {object} dto.Response "Current Wecom configuration"
// @Router /wecom/config [get]
func (b *BaseApi) WecomConfigAppGet(ctx *gin.Context) {
	conf, err := wecomLogic.ConfigAppGet()
	if err != nil {
		helper.ErrResponse(ctx, constant.CodeErrInternalServer)
		return
	}
	helper.SuccessWithData(ctx, conf)
}

// @Tags wecom_config
// @Summary Update Wecom configuration
// @Description Update the Wecom configuration with the provided data
// @Accept json
// @Produce json
// @Param data body request.WecomConfigApp true "Wecom configuration data"
// @Success 200 {object} dto.Response "Success response"
// @Router /wecom/config [patch]
func (b *BaseApi) WecomConfigAppUpdate(ctx *gin.Context) {
	var req request.WecomConfigApp
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helper.ErrResponse(ctx, constant.CodeErrBadRequest)
		return
	}
	if err := wecomLogic.ConfigAppUpdate(req); err != nil {
		helper.ErrResponse(ctx, constant.CodeErrInternalServer)
		return
	}
	helper.SuccessWithOutData(ctx)
}

// @Tags wecom_callback
// @Summary WeCom Callback URL Verification
// @Description Verifies the callback URL for WeCom (Enterprise WeChat). This endpoint is used to handle URL verification for WeCom's callback interface.
// @Accept json
// @Produce json
// @Param data body dto.SignatureOptions true "Wecom SignatureOptions data"
// @Success 200 {string} string "The echo string returned upon successful verification"
// @Router /wecom/callback [get]
func (b *BaseApi) WecomVerifyURL(ctx *gin.Context) {
	var options dto.SignatureOptions
	if err := ctx.ShouldBindQuery(&options); err != nil {
		helper.ErrResponse(ctx, constant.CodeErrBadRequest)
		return
	}
	echo, err := wecomLogic.VerifyURL(options)
	if err != nil {
		helper.ErrResponse(ctx, constant.CodeErrInternalServer)
		return
	}
	ctx.String(constant.CodeSuccess, echo)
}

// @Tags wecom_callback
// @Summary WeCom Callback Message Handling
// @Description Handles callback messages from WeCom (Enterprise WeChat). This endpoint receives and parses the callback data sent by WeCom.
// @Accept json
// @Produce json
// @Success 200 {object} string "Parsed WeCom callback message"
// @Router /wecom/callback [post]
func (b *BaseApi) WecomHandle(ctx *gin.Context) {
	// https://developer.work.weixin.qq.com/document/path/90930#32-%E6%94%AF%E6%8C%81http-post%E8%AF%B7%E6%B1%82%E6%8E%A5%E6%94%B6%E4%B8%9A%E5%8A%A1%E6%95%B0%E6%8D%AE
	// 正确响应本次请求
	// 企业微信服务器在五秒内收不到响应会断掉连接，并且重新发起请求，总共重试三次。
	// 仅针对网络连接失败或者网络请求超时情况重试，建议开发者接受回调后立即应答，业务异步处理。

	// 以上是企业微信对服务端的要求
	// 其中，如果 post 实际有到服务端，那么没有及时应答就会重复处理三次同一消息
	// 需要在服务端判断即使没有及时应答，服务端也应该对同一消息只处理一次，通过时间戳和随机数判断是否为同一消息
	// 一定要考虑服务端网络延迟极大的情况，在日志需要体现！！！
	body, err := ctx.GetRawData()
	if err != nil {
		helper.ErrResponse(ctx, constant.CodeErrBadRequest)
		return
	}
	timestamp := ctx.Query("timestamp")
	nonce := ctx.Query("nonce")
	messageID := fmt.Sprintf("%s-%s", timestamp, nonce)
	// 检查消息是否已经处理过
	exists, err := global.RDS.Exists(ctx, messageID).Result()
	if err != nil {
		global.ZAPLOG.Error("redis exists error", zap.Error(err))
		helper.ErrResponse(ctx, constant.CodeErrInternalServer)
		return
	}
	if exists > 0 {
		global.ZAPLOG.Info("收到重复消息事件, 不再处理")
		helper.SuccessWithOutData(ctx)
		return
	}
	// 标记消息为已处理，过期时间设置为15s，考虑负载并发
	if err = global.RDS.Set(ctx, messageID, 1, 15*time.Second).Err(); err != nil {
		global.ZAPLOG.Error("redis set error", zap.Error(err))
		helper.ErrResponse(ctx, constant.CodeErrInternalServer)
		return
	}
	if err = wecomLogic.Handle(body); err != nil {
		helper.ErrResponseWithErr(ctx, constant.CodeErrInternalServer, err)
		return
	}
	helper.SuccessWithOutData(ctx)
}

// @Tags wecom_receptionist
// @Summary WeCom receptionist add
// @Description receptionist add
// @Accept json
// @Produce json
// @Param data body request.ReceptionistOptions true "Wecom configuration data"
// @Success 200
// @Router /wecom/receptionist/ [post]
func (b *BaseApi) WecomReceptionistAdd(ctx *gin.Context) {
	var options request.ReceptionistOptions
	if err := ctx.ShouldBindJSON(&options); err != nil {
		helper.ErrResponse(ctx, constant.CodeErrBadRequest)
		return
	}
	if err := wecomLogic.ReceptionistAdd(options); err != nil {
		helper.ErrResponseWithErr(ctx, constant.CodeErrInternalServer, err)
		return
	}
	helper.SuccessWithOutData(ctx)
}

// @Tags wecom_receptionist
// @Summary WeCom receptionist delete
// @Description receptionist delete
// @Accept json
// @Produce json
// @Param data body request.ReceptionistOptions true "Wecom configuration data"
// @Success 200
// @Router /wecom/receptionist/ [delete]
func (b *BaseApi) WecomReceptionistDel(ctx *gin.Context) {
	var req request.ReceptionistOptions
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helper.ErrResponse(ctx, constant.CodeErrBadRequest)
		return
	}
	if err := wecomLogic.ReceptionistDel(req); err != nil {
		helper.ErrResponseWithErr(ctx, constant.CodeErrInternalServer, err)
		return
	}
	helper.SuccessWithOutData(ctx)
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
	list, err := wecomLogic.ReceptionistList(kfid)
	if err != nil {
		helper.ErrResponseWithErr(ctx, constant.CodeErrInternalServer, err)
		return
	}
	helper.SuccessWithData(ctx, list)
}

// @Tags wecom_account
// @Summary WeCom account list
// @Description account list
// @Accept json
// @Produce json
// @Success 200 {object} []client.AccountInfoSchema
// @Router /wecom/account/ [get]
func (b *BaseApi) WecomAccountList(ctx *gin.Context) {
	list, err := wecomLogic.AccountList()
	if err != nil {
		helper.ErrResponseWithErr(ctx, constant.CodeErrInternalServer, err)
		return
	}
	helper.SuccessWithData(ctx, list)
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
	url, err := wecomLogic.AddContactWay(kfid)
	if err != nil {
		helper.ErrResponseWithErr(ctx, constant.CodeErrInternalServer, err)
		return
	}
	helper.SuccessWithData(ctx, url)
}
