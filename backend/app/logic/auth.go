package logic

import (
	"EvoBot/backend/app/dto/request"
	"EvoBot/backend/app/dto/response"
	"EvoBot/backend/constant"
	"EvoBot/backend/global"
	"EvoBot/backend/utils/jwt"
	"EvoBot/backend/utils/passwd"
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type AuthLogic struct{}

type IAuthLogic interface {
	Login(c *gin.Context, req request.Login) (*response.Tokens, error)
	//LogOut(c *gin.Context) error
}

func NewIAuthLogic() IAuthLogic {
	return &AuthLogic{}
}

func (u *AuthLogic) Login(c *gin.Context, req request.Login) (*response.Tokens, error) {
	user, err := authRepo.Get(authRepo.WithByUsername(req.Username))
	if err != nil {
		return nil, errors.WithMessage(constant.ErrRecordNotFound, err.Error())
	}
	if !passwd.Verify(req.Password, user.Password) {
		return nil, constant.ErrAuth
	}
	accessToken, err := jwt.AccessToken(user.UUID)
	if err != nil {
		global.ZAPLOG.Error("failed to generate access token", zap.Error(err))
		return nil, err
	}
	redisKey := constant.KeyUserTokenPrefix + user.UUID + ":" + c.RemoteIP()
	expireDuration := global.CONF.AuthConfig.JwtExpired
	if expireDuration <= 0 {
		expireDuration = 2 * time.Hour
	}
	err = global.RDS.Set(context.Background(), redisKey, accessToken, expireDuration).Err()
	if err != nil {
		global.ZAPLOG.Error("failed to store token in Redis", zap.Error(err))
		return nil, err
	}
	return &response.Tokens{
		AccessToken: accessToken,
	}, nil
}
