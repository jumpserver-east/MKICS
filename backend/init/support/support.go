package support

import (
	"EvoBot/backend/global"
	"EvoBot/backend/utils/support"
)

func Init() {
	global.Support, _ = support.NewSupportClient(global.CONF.SupportConfig.Baseurl, global.CONF.SupportConfig.Username, global.CONF.SupportConfig.Password)
	global.ZAPLOG.Info("support init success")
}
