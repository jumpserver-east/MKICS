package router

import (
	v1 "EvoBot/backend/app/api/v1"
	"EvoBot/backend/middleware"

	"github.com/gin-gonic/gin"
)

type WecomRouter struct {
}

func (s *WecomRouter) InitRouter(Router *gin.RouterGroup) {
	wecomRouter := Router.Group("wecom")
	baseApi := v1.ApiGroupApp.BaseApi
	{
		wecomRouter.GET("/callback", baseApi.WecomVerifyURL)
		wecomRouter.POST("/callback", baseApi.WecomHandle)

		wecomRouter.GET("/config", middleware.AuthRequired(), baseApi.WecomConfigAppGet)
		wecomRouter.PATCH("/config", middleware.AuthRequired(), baseApi.WecomConfigAppUpdate)

		wecomRouter.GET("/receptionist/:kfid", middleware.AuthRequired(), baseApi.WecomReceptionistList)
		wecomRouter.POST("/receptionist/", middleware.AuthRequired(), baseApi.WecomReceptionistAdd)
		wecomRouter.DELETE("/receptionist/", middleware.AuthRequired(), baseApi.WecomReceptionistDel)

		wecomRouter.GET("/account/", middleware.AuthRequired(), baseApi.WecomAccountList)
		wecomRouter.GET("/account/:kfid", middleware.AuthRequired(), baseApi.WecomAddContactWay)
	}
}
