package router

import (
	"EvoBot/backend/global"
	rou "EvoBot/backend/router"
	"EvoBot/cmd/server/docs"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Init() *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Requested-With"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			log.Println("origin: ", origin)
			return true
		},
	}))
	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	PublicGroup := r.Group("")
	{
		PublicGroup.GET("/health", func(c *gin.Context) {
			c.JSON(200, "ok")
		})
		PublicGroup.Use(gzip.Gzip(gzip.DefaultCompression))
	}
	PrivateGroup := r.Group("/api/v1")
	for _, router := range rou.RouterGroupApp {
		router.InitRouter(PrivateGroup)
	}
	global.ZAPLOG.Info("init router success")
	return r
}
