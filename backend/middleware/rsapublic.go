package middleware

import (
	"MKICS/backend/constant"
	"MKICS/backend/global"
	"MKICS/backend/utils/encrypt"
	"context"
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	keyExpire = time.Minute * 120
)

var (
	ctx = context.Background()
)

func RSAPublicKeyMiddleware() gin.HandlerFunc {
	whiteList := []string{
		"/ui/login",
		"/",
	}

	return func(c *gin.Context) {
		path := c.Request.URL.Path
		matched := false
		for _, p := range whiteList {
			if strings.HasPrefix(path, p) {
				matched = true
				break
			}
		}
		if !matched {
			c.Next()
			return
		}

		pubKey, errPub := global.RDS.Get(ctx, constant.PublicKeyCacheKey).Result()

		cookiePubKey, errCookie := c.Cookie(constant.PublicKeyCookieKey)
		if errPub == nil && errCookie == nil && cookiePubKey == pubKey {
			c.Next()
			return
		}

		privateKey, publicKey, err := encrypt.GenerateRSAKeyPair()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate RSA keys"})
			return
		}

		publicKeyEncoded := base64.StdEncoding.EncodeToString([]byte(publicKey))

		global.RDS.Set(ctx, constant.PublicKeyCacheKey, publicKeyEncoded, keyExpire)
		global.RDS.Set(ctx, constant.PrivateKeyCacheKey, privateKey, keyExpire)

		c.SetCookie(constant.PublicKeyCookieKey, publicKeyEncoded, int(keyExpire.Seconds()), "/", "", false, false)

		c.Next()
	}
}
