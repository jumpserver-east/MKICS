package server

import (
	"MKICS/backend/global"
	"MKICS/backend/init/db"
	"MKICS/backend/init/migration"
	"MKICS/backend/init/redis"
	"MKICS/backend/init/router"
	"MKICS/backend/init/support"
	"MKICS/backend/init/viper"
	"MKICS/backend/init/zaplog"
	"fmt"
)

func Start(configPath string) {
	viper.Init(configPath)
	zaplog.Init()

	redis.Init()
	db.Init()
	migration.Init()

	support.Init()
	r := router.Init()

	r.Run(fmt.Sprintf("%s:%d", global.CONF.SystemConfig.Host, global.CONF.SystemConfig.Port))
}
