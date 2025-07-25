package logic

import (
	"MKICS/backend/app/dto/request"
	"MKICS/backend/app/dto/response"
	"MKICS/backend/constant"
	"MKICS/backend/global"
	"MKICS/backend/utils/jwt"
	"MKICS/backend/utils/passwd"
	"context"
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthLogic struct{}

type IAuthLogic interface {
	Login(c *gin.Context, req request.Login) (*response.Tokens, error)
	Logout(c *gin.Context, token string) error
}

func NewIAuthLogic() IAuthLogic {
	return &AuthLogic{}
}

func (u *AuthLogic) Login(c *gin.Context, req request.Login) (*response.Tokens, error) {
	user, err := authRepo.Get(authRepo.WithByUsername(req.Username))
	if err != nil {
		return nil, err
	}
	if !passwd.Verify(req.Password, user.Password) {
		return nil, errors.New("login failed")
	}
	accessToken, err := jwt.AccessToken(user.UUID)
	if err != nil {
		global.ZAPLOG.Error(err.Error())
		return nil, err
	}
	redisKey := constant.KeyUserTokenPrefix + user.UUID + ":" + c.RemoteIP()
	expireDuration := global.CONF.AuthConfig.JwtExpired
	if expireDuration <= 0 {
		expireDuration = 2 * time.Hour
	}
	err = global.RDS.Set(context.Background(), redisKey, accessToken, expireDuration).Err()
	if err != nil {
		global.ZAPLOG.Error(err.Error())
		return nil, err
	}
	return &response.Tokens{
		AccessToken: accessToken,
	}, nil
}

func (u *AuthLogic) Logout(c *gin.Context, token string) error {
	parts := strings.SplitN(token, " ", 2)
	claims, err := jwt.VerifyToken(parts[1])
	if err != nil {
		global.ZAPLOG.Error(err.Error())
		return err
	}
	redisKey := constant.KeyUserTokenPrefix + claims.UUID + ":" + c.RemoteIP()
	err = global.RDS.Del(context.Background(), redisKey).Err()
	if err != nil {
		global.ZAPLOG.Error(err.Error())
		return err
	}
	return nil
}
