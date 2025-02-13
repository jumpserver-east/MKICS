package router

import (
	"EvoBot/backend/constant"
	"EvoBot/backend/global"
	rou "EvoBot/backend/router"
	"EvoBot/cmd/server/docs"
	"EvoBot/cmd/server/web"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func toIndexHtml(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.Writer.WriteHeader(http.StatusOK)
	_, _ = c.Writer.Write(web.IndexByte)
	c.Writer.Flush()
}

func isFrontendPath(c *gin.Context) bool {
	reqUri := strings.TrimSuffix(c.Request.URL.Path, "/")
	if _, ok := constant.WebUrlMap[reqUri]; ok {
		return true
	}
	for _, route := range constant.DynamicRoutes {
		if match, _ := regexp.MatchString(route, reqUri); match {
			return true
		}
	}
	return false
}

func setWebStatic(rootRouter *gin.RouterGroup) {
	rootRouter.StaticFS("/public", http.FS(web.Favicon))
	rootRouter.Static("/api/v1/images", "./uploads")
	rootRouter.Use(func(c *gin.Context) {
		c.Next()
	})
	rootRouter.GET("/assets/*filepath", func(c *gin.Context) {
		c.Writer.Header().Set("Cache-Control", fmt.Sprintf("private, max-age=%d", 3600))
		staticServer := http.FileServer(http.FS(web.Assets))
		staticServer.ServeHTTP(c.Writer, c.Request)
	})
	rootRouter.GET("/", func(c *gin.Context) {
		staticServer := http.FileServer(http.FS(web.IndexHtml))
		staticServer.ServeHTTP(c.Writer, c.Request)
	})
}

func Init() *gin.Engine {
	r := gin.Default()
	r.NoRoute(func(c *gin.Context) {
		if isFrontendPath(c) {
			toIndexHtml(c)
			return
		}
	})
	// swaggerRouter := r.Group("evobot").Use(middleware.AuthRequired())
	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	PublicGroup := r.Group("")
	{
		PublicGroup.GET("/health", func(c *gin.Context) {
			c.JSON(200, "ok")
		})
		PublicGroup.Use(gzip.Gzip(gzip.DefaultCompression))
		setWebStatic(PublicGroup)
	}
	PrivateGroup := r.Group("/api/v1")
	for _, router := range rou.RouterGroupApp {
		router.InitRouter(PrivateGroup)
	}
	global.ZAPLOG.Info("init router success")
	return r
}
