package logic

import (
	"MKICS/backend/app/dto/request"
	"MKICS/backend/app/dto/response"
	"MKICS/backend/constant"
	"MKICS/backend/global"
	"MKICS/backend/middleware"
	"MKICS/backend/utils/bcrypt"
	"MKICS/backend/utils/jwt"
	"MKICS/backend/utils/redis"

	"github.com/gin-gonic/gin"
)

type AuthLogic struct{}

type IAuthLogic interface {
	Login(ctx *gin.Context, req request.Login) (token *response.Token, err error)
	Logout(ctx *gin.Context) (err error)
}

func NewIAuthLogic() IAuthLogic {
	return &AuthLogic{}
}

func (u *AuthLogic) Login(ctx *gin.Context, req request.Login) (token *response.Token, err error) {
	token = &response.Token{}

	user, err := authRepo.Get(authRepo.WithByUsername(req.Username))
	if err != nil {
		return
	}

	if !bcrypt.Verify(req.Password, user.Password) {
		return nil, constant.ErrLoginFailed
	}

	jwtTokenString, err := jwt.GenerateToken(user.UUID)
	if err != nil {
		global.ZAPLOG.Error(err.Error())
		return
	}

	err = redis.SetToken(user.UUID, ctx.ClientIP(), jwtTokenString)
	if err != nil {
		global.ZAPLOG.Error(err.Error())
		return
	}

	token.AccessToken = jwtTokenString

	return
}

func (u *AuthLogic) Logout(ctx *gin.Context) (err error) {
	user := ctx.MustGet(constant.UserKey).(*middleware.CtxUser)

	err = redis.DelToken(user.UUID, user.IP)
	if err != nil {
		global.ZAPLOG.Error(err.Error())
		return
	}

	return
}
