package middleware

import (
	"MKICS/backend/app/dto/response"
	"MKICS/backend/constant"
	"MKICS/backend/utils/jwt"
	"MKICS/backend/utils/redis"
	"strings"

	"github.com/gin-gonic/gin"
)

type CtxUser struct {
	UUID string
	IP   string
}

func AuthRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(constant.AuthHeaderKey)
		if authHeader == "" || !strings.HasPrefix(authHeader, constant.AuthHeaderPrefixBearer) {
			response.Unauthorized(ctx, constant.ErrMissingOrInvalidAuthHeader)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, constant.AuthHeaderPrefixBearer)

		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			response.Unauthorized(ctx, constant.ErrTokenParse)
			return
		}

		if !redis.ValidateToken(claims.UUID, ctx.ClientIP(), tokenString) {
			response.Unauthorized(ctx, constant.ErrTokenUsedElsewhere)
			return
		}

		ctx.Set(constant.UserKey, &CtxUser{
			UUID: claims.UUID,
			IP:   ctx.ClientIP(),
		})

		ctx.Next()
	}
}
