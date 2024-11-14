package middleware

import (
	"EvoBot/backend/constant"
	"EvoBot/backend/global"
	"EvoBot/backend/utils/jwt"
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AuthRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get(constant.TokenKey)
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"msg": "no login"})
			ctx.Abort()
			return
		}
		parts := strings.SplitN(token, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			ctx.JSON(http.StatusBadRequest, gin.H{"msg": "token format is incorrect"})
			ctx.Abort()
			return
		}
		claims, err := jwt.VerifyToken(parts[1])
		if err != nil {
			global.ZAPLOG.With(zap.String("token", parts[1]), zap.String("IP", ctx.RemoteIP())).Error(err.Error())
			ctx.JSON(http.StatusUnauthorized, gin.H{"msg": "invalid or expired token"})
			ctx.Abort()
			return
		}
		if !validateTokenWithRedis(claims.UUID, parts[1], ctx.RemoteIP()) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"msg": "token expired or used elsewhere"})
			ctx.Abort()
			return
		}
		ctx.Set("user", map[string]interface{}{
			"uuid": claims.UUID,
		})
		ctx.Next()
	}
}

func validateTokenWithRedis(uuid string, token, ip string) bool {
	result, err := global.RDS.Get(context.Background(), constant.KeyUserTokenPrefix+uuid+":"+ip).Result()
	if err != nil || result != token {
		global.ZAPLOG.With(zap.String("UUID", uuid), zap.String("IP", ip)).Error(err.Error())
		return false
	}
	return true
}
