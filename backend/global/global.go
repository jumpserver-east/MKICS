package global

import (
	"MKICS/backend/configs"

	"github.com/go-redis/redis/v8"
	"github.com/silenceper/wechat/v2/work"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	DB     *gorm.DB
	RDS    redis.UniversalClient
	Wecom  *work.Work
	ZAPLOG *zap.Logger
	CONF   configs.ServerConfig
)
