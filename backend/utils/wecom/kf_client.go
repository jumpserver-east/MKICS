package wecom

import (
	"MKICS/backend/global"
	"MKICS/backend/utils/wecom/client"
)

func NewWecomKFClient(conf client.WecomConfig) (*client.WecomKF, error) {
	global.Wecom = client.NewWork(conf)
	kFClient, err := global.Wecom.GetKF()
	if err != nil {
		return nil, err
	}
	return client.NewWecomKF(kFClient), nil
}
