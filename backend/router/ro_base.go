package router

import (
	v1 "MKICS/backend/app/api/v1"
	"MKICS/backend/middleware"

	"github.com/gin-gonic/gin"
)

type BaseRouter struct{}

func (s *BaseRouter) InitRouter(Router *gin.RouterGroup) {
	baseRouter := Router.Group("auth")
	baseApi := v1.ApiGroupApp.BaseApi
	{
		baseRouter.POST("/login", baseApi.Login)
		baseRouter.Use(middleware.AuthRequired()).
			POST("/logout", baseApi.Logout)
	}
}
