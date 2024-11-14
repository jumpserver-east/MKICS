package router

import (
	v1 "EvoBot/backend/app/api/v1"
	"EvoBot/backend/middleware"

	"github.com/gin-gonic/gin"
)

type MaxkbRouter struct {
}

func (s *MaxkbRouter) InitRouter(Router *gin.RouterGroup) {
	// maxkbRouter := Router.Group("maxkb")
	maxkbRouter := Router.Group("maxkb").Use(middleware.AuthRequired())
	baseApi := v1.ApiGroupApp.BaseApi
	{
		maxkbRouter.GET("/config", baseApi.ListMaxkbConfig)
		maxkbRouter.PATCH("/config", baseApi.UpdateMaxkbConfig)
	}
}
