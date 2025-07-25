package global

import (
	"MKICS/backend/configs"
	"MKICS/backend/utils/support"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	DB      *gorm.DB
	RDS     redis.UniversalClient
	Support *support.Client
	ZAPLOG  *zap.Logger
	CONF    configs.ServerConfig
)
