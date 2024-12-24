package client

import (
	"EvoBot/backend/global"
	"context"
	"strconv"
	"sync"

	"github.com/silenceper/wechat/v2/work/kf"
	"go.uber.org/zap"
)

type UserMsgQueue struct {
	mu sync.Mutex
}

type WecomKF struct {
	KFClient   *kf.Client
	UserQueues map[string]*UserMsgQueue
	Mu         *sync.Mutex
}

func NewWecomKF(kf WecomKF) *WecomKF {
	return &kf
}

func (k *WecomKF) VerifyURL(options SignatureOptions) (string, error) {
	echo, err := k.KFClient.VerifyURL(kf.SignatureOptions{
		Signature: options.Signature,
		TimeStamp: options.TimeStamp,
		Nonce:     options.Nonce,
		EchoStr:   options.EchoStr,
	})
	if err != nil {
		return "", err
	}
	return echo, nil
}

func (k *WecomKF) SyncMsg(body []byte) (MessageInfo, error) {
	var messageInfo MessageInfo
	callbackMessage, err := k.KFClient.GetCallbackMessage(body)
	if callbackMessage.MsgType == "event" && callbackMessage.Event == "kf_msg_or_event" && err == nil {
		var syncMsgOptions kf.SyncMsgOptions
		syncMsgOptions.OpenKfID = callbackMessage.OpenKfID
		syncMsgOptions.Token = callbackMessage.Token
		wecomcursorkey := KeyWecomCursorPrefix + callbackMessage.OpenKfID
		exists, err := global.RDS.Exists(context.Background(), wecomcursorkey).Result()
		if err != nil {
			global.ZAPLOG.Error("failed to check key existence", zap.Error(err))
		}
		if exists > 0 {
			wecomcursorvalue, err := global.RDS.Get(context.Background(), wecomcursorkey).Result()
			if err != nil {
				global.ZAPLOG.Error("failed to get wecomcursorkey", zap.Error(err))
			}
			syncMsgOptions.Cursor = wecomcursorvalue
		}
		message, err := k.KFClient.SyncMsg(syncMsgOptions)
		if err != nil {
			return messageInfo, err
		}
		if err = global.RDS.Set(context.Background(), wecomcursorkey, message.NextCursor, 0).Err(); err != nil {
			return messageInfo, err
		}
		for message.HasMore == 1 {
			syncMsgOptions.Cursor = message.NextCursor
			messageMore, err := k.KFClient.SyncMsg(syncMsgOptions)
			if err != nil {
				return messageInfo, err
			}
			if err = global.RDS.Set(context.Background(), wecomcursorkey, message.NextCursor, 0).Err(); err != nil {
				return messageInfo, err
			}
			message.MsgList = append(message.MsgList, messageMore.MsgList...)
			message.HasMore = messageMore.HasMore
		}
		if len(message.MsgList) > 0 {
			msg := message.MsgList[0]
			messageInfo.MessageID = msg.MsgID
			messageInfo.Origin = msg.Origin
			messageInfo.SendTime = msg.SendTime
			if msg.Origin == MessageTypeReceptionistMessage {
				messageInfo.StaffID = msg.ReceptionistUserID
			}
			switch msg.MsgType {
			case WecomMsgTypeText:
				messageInfo.KFID = msg.OpenKFID
				messageInfo.KHID = msg.ExternalUserID
				textmsg, _ := msg.GetTextMessage()
				messageInfo.MessageType = textmsg.MsgType
				messageInfo.Message = textmsg.Text.Content
			case WecomMsgTypeEvent:
				switch msg.EventType {
				case WecomEventTypeEnterSession:
					entersessioninfo, _ := msg.GetEnterSessionEvent()
					messageInfo.KFID = entersessioninfo.Event.OpenKFID
					messageInfo.KHID = entersessioninfo.Event.ExternalUserID
					messageInfo.MessageType = entersessioninfo.Event.EventType
					messageInfo.Credential = entersessioninfo.Event.WelcomeCode
				case WecomEventTypeSessionStatusChange:
					sessionstatuschangeinfo, _ := msg.GetSessionStatusChangeEvent()
					messageInfo.KFID = sessionstatuschangeinfo.OpenKFID
					messageInfo.KHID = sessionstatuschangeinfo.ExternalUserID
					messageInfo.Message = strconv.FormatUint(uint64(sessionstatuschangeinfo.Event.ChangeType), 10)
					messageInfo.MessageType = sessionstatuschangeinfo.Event.EventType
					if sessionstatuschangeinfo.Event.OldReceptionistUserID != "" {
						messageInfo.StaffID = sessionstatuschangeinfo.Event.OldReceptionistUserID
					}
					if sessionstatuschangeinfo.Event.NewReceptionistUserID != "" {
						messageInfo.StaffID = sessionstatuschangeinfo.Event.NewReceptionistUserID
					}
					messageInfo.Credential = sessionstatuschangeinfo.Event.MsgCode
				default:
					global.ZAPLOG.Info("此事件类型服务暂不处理", zap.String("EventType", msg.EventType))
					return messageInfo, err
				}
			default:
				global.ZAPLOG.Info("此消息类型服务暂不处理", zap.String("MsgType", msg.MsgType))
				return messageInfo, err
			}
			statusinfo, _ := k.KFClient.ServiceStateGet(kf.ServiceStateGetOptions{
				OpenKFID:       messageInfo.KFID,
				ExternalUserID: messageInfo.KHID,
			})
			if statusinfo.ServiceState == SessionStatusNew {
				_, err := k.KFClient.ServiceStateTrans(kf.ServiceStateTransOptions{
					OpenKFID:       msg.OpenKFID,
					ExternalUserID: msg.ExternalUserID,
					ServiceState:   SessionStatusHandled,
				})
				if err != nil {
					return messageInfo, err
				}
				statusinfo.ServiceState = SessionStatusHandled
			}
			messageInfo.ChatState = statusinfo.ServiceState
			return messageInfo, nil
		}
		return messageInfo, err
	}
	return messageInfo, err
}

func (k *WecomKF) SendTextMsg(info SendTextMsgOptions) error {
	k.Mu.Lock()
	if k.UserQueues == nil {
		k.UserQueues = make(map[string]*UserMsgQueue)
	}
	userQueue, exists := k.UserQueues[info.KHID]
	if !exists {
		userQueue = &UserMsgQueue{}
		k.UserQueues[info.KHID] = userQueue
	}
	k.Mu.Unlock()

	userQueue.mu.Lock()
	defer userQueue.mu.Unlock()

	sendmsgs := splitContent(info.Message)
	for _, content := range sendmsgs {
		sendmsg := struct {
			Touser   string `json:"touser"`
			OpenKfid string `json:"open_kfid"`
			MsgID    string `json:"msgid,omitempty"`
			MsgType  string `json:"msgtype"`
			Text     struct {
				Content string `json:"content"`
			} `json:"text"`
		}{
			Touser:   info.KHID,
			OpenKfid: info.KFID,
			MsgType:  WecomMsgTypeText,
		}
		sendmsg.Text.Content = content
		if _, err := k.KFClient.SendMsg(sendmsg); err != nil {
			return err
		}
	}
	return nil
}

func (k *WecomKF) SendMenuMsg(info SendMenuMsgOptions) error {
	sendmsg := struct {
		KFID    string `json:"open_kfid"`
		KHID    string `json:"touser"`
		MsgID   string `json:"msgid,omitempty"`
		MsgType string `json:"msgtype"`
		MsgMenu struct {
			HeadContent string     `json:"head_content,omitempty"`
			List        []MenuItem `json:"list"`
			TailContent string     `json:"tail_content,omitempty"`
		} `json:"msgmenu"`
	}{
		KFID:    info.KFID,
		KHID:    info.KHID,
		MsgType: "msgmenu",
	}
	sendmsg.MsgMenu.HeadContent = info.HeadContent
	sendmsg.MsgMenu.TailContent = info.TailContent
	sendmsg.MsgMenu.List = info.MenuList
	if _, err := k.KFClient.SendMsg(sendmsg); err != nil {
		return err
	}
	return nil
}

func (k *WecomKF) SendTextMsgOnEvent(info SendTextMsgOnEventOptions) error {
	sendmsg := struct {
		Code    string `json:"code"`
		MsgID   string `json:"msgid,omitempty"`
		MsgType string `json:"msgtype"`
		Text    struct {
			Content string `json:"content"`
		} `json:"text"`
	}{
		Code:    info.Credential,
		MsgType: WecomMsgTypeText,
	}
	sendmsg.Text.Content = info.Message
	if _, err := k.KFClient.SendMsgOnEvent(sendmsg); err != nil {
		return err
	}
	return nil
}

func (k *WecomKF) SendMenuMsgOnEvent(info SendMenuMsgOnEventOptions) error {
	sendmsg := struct {
		Code    string `json:"code"`
		MsgID   string `json:"msgid,omitempty"`
		MsgType string `json:"msgtype"`
		MsgMenu struct {
			HeadContent string     `json:"head_content,omitempty"`
			List        []MenuItem `json:"list"`
			TailContent string     `json:"tail_content,omitempty"`
		} `json:"msgmenu"`
	}{
		Code:    info.Credential,
		MsgType: "msgmenu",
	}
	sendmsg.MsgMenu.HeadContent = info.HeadContent
	sendmsg.MsgMenu.TailContent = info.TailContent
	sendmsg.MsgMenu.List = info.MenuList
	if _, err := k.KFClient.SendMsgOnEvent(sendmsg); err != nil {
		return err
	}
	return nil
}

func (k *WecomKF) ServiceStateGet(options ServiceStateGetOptions) (int, error) {
	res, err := k.KFClient.ServiceStateGet(kf.ServiceStateGetOptions{
		OpenKFID:       options.OpenKFID,
		ExternalUserID: options.ExternalUserID,
	})
	if err != nil {
		return -1, err
	}
	return res.ServiceState, nil
}

func (k *WecomKF) ServiceStateTrans(options ServiceStateTransOptions, servicestate int) (string, error) {
	res, err := k.KFClient.ServiceStateTrans(kf.ServiceStateTransOptions{
		OpenKFID:       options.OpenKFID,
		ExternalUserID: options.ExternalUserID,
		ServiceState:   servicestate,
		ServicerUserID: options.ServicerUserID,
	})
	if err != nil {
		return "", err
	}
	return res.MsgCode, nil
}

func (k *WecomKF) ReceptionistList(kfid string) ([]ReceptionistList, error) {
	var receptionistlsit []ReceptionistList
	receptionistinfo, err := k.KFClient.ReceptionistList(kfid)
	if err != nil {
		return nil, err
	}
	for _, receptionist := range receptionistinfo.ReceptionistList {
		receptionistlsit = append(receptionistlsit, ReceptionistList{
			UserID: receptionist.UserID,
			Status: receptionist.Status,
		})
	}
	return receptionistlsit, nil
}

func (k *WecomKF) ReceptionistAdd(options ReceptionistOptions) error {
	_, err := k.KFClient.ReceptionistAdd(kf.ReceptionistOptions{
		OpenKFID:   options.OpenKFID,
		UserIDList: options.UserIDList,
	})
	if err != nil {
		return err
	}
	return nil
}

func (k *WecomKF) ReceptionistDel(options ReceptionistOptions) error {
	_, err := k.KFClient.ReceptionistDel(kf.ReceptionistOptions{
		OpenKFID:   options.OpenKFID,
		UserIDList: options.UserIDList,
	})
	if err != nil {
		return err
	}
	return nil
}

func (k *WecomKF) AccountAdd(options AccountAddOptions) (string, error) {
	accountres, err := k.KFClient.AccountAdd(kf.AccountAddOptions{
		Name:    options.Name,
		MediaID: options.MediaID,
	})
	if err != nil {
		return "", err
	}
	return accountres.OpenKFID, nil
}

func (k *WecomKF) AccountUpdate(options AccountUpdateOptions) error {
	_, err := k.KFClient.AccountUpdate(kf.AccountUpdateOptions{
		OpenKFID: options.OpenKFID,
		Name:     options.Name,
		MediaID:  options.MediaID,
	})
	if err != nil {
		return err
	}
	return nil
}

func (k *WecomKF) AccountDel(kfid string) error {
	_, err := k.KFClient.AccountDel(kf.AccountDelOptions{
		OpenKFID: kfid,
	})
	if err != nil {
		return err
	}
	return nil
}

func (k *WecomKF) AccountList() ([]AccountInfoSchema, error) {
	accountres, err := k.KFClient.AccountList()
	if err != nil {
		return nil, err
	}
	var accountinfo []AccountInfoSchema
	for _, account := range accountres.AccountList {
		accountinfo = append(accountinfo, AccountInfoSchema{
			OpenKFID: account.OpenKFID,
			Name:     account.Name,
			Avatar:   account.Avatar,
		})
	}
	return accountinfo, nil
}

func (k *WecomKF) AddContactWay(kfid string) (string, error) {
	var opt kf.AddContactWayOptions
	opt.OpenKFID = kfid
	opt.Scene = "evobot"
	info, err := k.KFClient.AddContactWay(opt)
	if err != nil {
		return "", err
	}
	return info.URL, nil
}
