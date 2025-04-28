package wecom

import (
	"EvoBot/backend/utils/wecom/client"
	"sync"
)

func NewWecomKFClient(conf client.WecomConfig) (*client.WecomKF, error) {
	kFClient, err := client.NewWork(conf).GetKF()
	if err != nil {
		return nil, err
	}
	return client.NewWecomKF(client.WecomKF{
		KFClient:   kFClient,
		UserQueues: make(map[string]*client.UserMsgQueue),
		Mu:         &sync.Mutex{},
	}), nil
}
