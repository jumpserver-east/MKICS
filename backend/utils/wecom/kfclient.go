package wecom

import (
	"EvoBot/backend/utils/wecom/client"
	"sync"
)

type WecomKFClient interface {
	VerifyURL(options client.SignatureOptions) (string, error)

	SyncMsg(body []byte) (client.MessageInfo, error)

	SendTextMsg(info client.SendTextMsgOptions) error
	SendTextMsgOnEvent(info client.SendTextMsgOnEventOptions) error
	SendMenuMsgOnEvent(info client.SendMenuMsgOnEventOptions) error

	ServiceStateTransToStaff(options client.ServiceStateTransOptions) (string, error)
	ServiceStateTransToEnd(options client.ServiceStateTransOptions) (string, error)

	ReceptionistAdd(options client.ReceptionistOptions) error
	ReceptionistList(kfid string) ([]client.ReceptionistList, error)
	ReceptionistDel(options client.ReceptionistOptions) error

	AccountAdd(options client.AccountAddOptions) (string, error)
	AccountUpdate(options client.AccountUpdateOptions) error
	AccountDel(kfid string) error
	AccountList() ([]client.AccountInfoSchema, error)
	AddContactWay(kfid string) (string, error)
}

func NewWecomKFClient(conf client.WecomConfig) (WecomKFClient, error) {
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
