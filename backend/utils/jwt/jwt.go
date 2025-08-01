package jwt

import (
	"MKICS/backend/global"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret string
var secretOnce sync.Once

type Claims struct {
	UUID string `json:"uuid"`
	jwt.RegisteredClaims
}

func initSecret() {
	authConfig := global.CONF.AuthConfig
	if authConfig.Secret == "" {
		global.ZAPLOG.Error("secret config is empty")
		return
	}
	jwtSecret = authConfig.Secret
}

func getSecret() []byte {
	secretOnce.Do(initSecret)
	return []byte(jwtSecret)
}

func GenerateToken(uuid string) (jwtTokenString string, err error) {
	jwtExpire := global.CONF.AuthConfig.JwtExpired

	claims := Claims{
		UUID: uuid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtExpire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtSecret := getSecret()

	return jwtToken.SignedString(jwtSecret)
}

func ParseToken(jwtTokenString string) (claims *Claims, err error) {
	jwtSecret := getSecret()

	jwtToken, err := jwt.ParseWithClaims(jwtTokenString, &Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
	if err != nil {
		global.ZAPLOG.Error(err.Error())
		return
	}

	claims = jwtToken.Claims.(*Claims)

	return
}
