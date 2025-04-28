package logic

import (
	"EvoBot/backend/app/dto"
	"EvoBot/backend/app/dto/request"
	"EvoBot/backend/app/dto/response"
	"EvoBot/backend/app/model"
	"EvoBot/backend/constant"
	"EvoBot/backend/global"
	"EvoBot/backend/i18n"
	"EvoBot/backend/utils/wecom"
	wecomclient "EvoBot/backend/utils/wecom/client"
	"fmt"
	"strings"

	"context"
	"time"

	"go.uber.org/zap"
)

type WecomLogic struct {
	wecomkf *wecomclient.WecomKF
}

type IWecomLogic interface {
	VerifyURL(req dto.SignatureOptions) (string, error)
	CallbackHandler(encryptedMsg []byte) (err error)

	ConfigList() ([]response.WecomConfigApp, error)
	ConfigUpdate(uuid string, req request.WecomConfigApp) error
	ConfigGet(uuid string) (response.WecomConfigApp, error)

	ReceptionistAdd(kfid string, req request.ReceptionistOptions) error
	ReceptionistDel(kfid string, req request.ReceptionistOptions) error
	ReceptionistList(kfid string) ([]wecomclient.ReceptionistList, error)
	CheckReceptionist(stafflist []string, receptionistlist []wecomclient.ReceptionistList) error

	AccountList() ([]wecomclient.AccountInfoSchema, error)
	AddContactWay(kfid string) (string, error)
}

func NewIWecomLogic() IWecomLogic {
	return &WecomLogic{}
}

func (u *WecomLogic) ConfigUpdate(uuid string, req request.WecomConfigApp) error {
	conf, err := wecomRepo.Get(commonRepo.WithByUUID(uuid))
	if err != nil {
		global.ZAPLOG.Error("failed to get Wecom config", zap.Error(err))
		return err
	}
	conf.CorpID = req.CorpID
	conf.Secret = req.Secret
	conf.AgentID = req.AgentID
	conf.Token = req.Token
	conf.EncodingAESKey = req.EncodingAESKey
	err = wecomRepo.Update(context.Background(), &conf)
	if err != nil {
		global.ZAPLOG.Error("failed to update Wecom config", zap.Error(err))
		return err
	}
	u.wecomkf, err = wecom.NewWecomKFClient(wecomclient.WecomConfig{
		CorpID:         conf.CorpID,
		Token:          conf.Token,
		EncodingAESKey: conf.EncodingAESKey,
		Secret:         conf.Secret,
		AgentID:        conf.AgentID,
	})
	if err != nil {
		global.ZAPLOG.Error("failed to initialize WecomKFClient", zap.Error(err))
		return err
	}
	return nil
}

func (u *WecomLogic) ConfigGet(uuid string) (response.WecomConfigApp, error) {
	var res response.WecomConfigApp
	conf, err := wecomRepo.Get(commonRepo.WithByUUID(uuid))
	if err != nil {
		global.ZAPLOG.Error("failed to get Wecom config", zap.Error(err))
		return res, err
	}
	res.UUID = conf.UUID
	res.Type = conf.Type
	res.CorpID = conf.CorpID
	res.Token = conf.Token
	res.EncodingAESKey = conf.EncodingAESKey
	res.AgentID = conf.AgentID
	return res, nil
}

func (u *WecomLogic) ConfigList() ([]response.WecomConfigApp, error) {
	var resp []response.WecomConfigApp
	configs, err := wecomRepo.List()
	if err != nil {
		global.ZAPLOG.Error("failed to get Wecom config", zap.Error(err))
		return resp, err
	}
	for _, config := range configs {
		var res response.WecomConfigApp
		res.UUID = config.UUID
		res.Type = config.Type
		res.CorpID = config.CorpID
		res.Token = config.Token
		res.EncodingAESKey = config.EncodingAESKey
		res.AgentID = config.AgentID
		resp = append(resp, res)
	}
	return resp, nil
}

func (u *WecomLogic) VerifyURL(req dto.SignatureOptions) (echo string, err error) {
	if err = u.initializeWecomClient(); err != nil {
		return
	}
	return u.wecomkf.VerifyURL(wecomclient.SignatureOptions{SignatureOptions: req.SignatureOptions.SignatureOptions})
}

func (u *WecomLogic) CallbackHandler(encryptedMsg []byte) (err error) {
	if err = u.initializeWecomClient(); err != nil {
		return
	}
	callbackmessage, err := u.wecomkf.KFClient.GetCallbackMessage(encryptedMsg)
	if err != nil {
		global.ZAPLOG.Error(i18n.Tf("wecom.failed_action", "GetCallbackMessage"), zap.Error(err))
		return
	}
	switch callbackmessage.MsgType {
	case wecomclient.WecomCallbackMsgTypeEvent:
		switch callbackmessage.Event {
		case wecomclient.WecomCallbackEventKFMsgOrEvent:
			go func() {
				var options wecomclient.SyncMsgOptions
				options.Token = callbackmessage.Token
				options.OpenKfID = callbackmessage.OpenKfID
				if err := u.SyncMsg(options); err != nil {
					global.ZAPLOG.Error(i18n.Tf("wecom.failed_action", "SyncMsgs"), zap.Error(err))
					return
				}
			}()
		default:
			global.ZAPLOG.Error(i18n.T("wecom.unknow_event_type"), zap.Error(err))
		}
	default:
		global.ZAPLOG.Error(i18n.T("wecom.unknow_msgtype"), zap.Error(err))
	}
	return
}

func (u *WecomLogic) SyncMsg(options wecomclient.SyncMsgOptions) (err error) {
	wecomCursorLockKey := wecomclient.KeyWecomCursorLockPrefix + options.OpenKfID
	if !global.RDS.SetNX(context.Background(), wecomCursorLockKey, "1", 30*time.Second).Val() {
		return
	}
	defer global.RDS.Del(context.Background(), wecomCursorLockKey)
	wecomCursorKey := wecomclient.KeyWecomCursorPrefix + options.OpenKfID
	options.Cursor, _ = global.RDS.Get(context.Background(), wecomCursorKey).Result()
	for {
		resp, err := u.wecomkf.SyncMsgs(options)
		if err != nil {
			break
		}
		for _, msg := range resp.MsgList {
			if err := u.processMessage(wecomclient.Message{
				Message: msg,
			}); err != nil {
				return err
			}
		}
		if err = global.RDS.Set(context.Background(), wecomCursorKey, resp.NextCursor, 0).Err(); err != nil {
			break
		}
		if resp.HasMore == 0 {
			break
		}
		options.Cursor = resp.NextCursor
	}
	return
}

func (u *WecomLogic) ReceptionistAdd(kfid string, req request.ReceptionistOptions) error {
	if err := u.initializeWecomClient(); err != nil {
		return err
	}
	err := u.wecomkf.ReceptionistAdd(wecomclient.ReceptionistOptions{
		OpenKFID:   kfid,
		UserIDList: req.UserIDList,
	})
	if err != nil {
		return err
	}
	return nil
}

func (u *WecomLogic) ReceptionistDel(kfid string, req request.ReceptionistOptions) error {
	if err := u.initializeWecomClient(); err != nil {
		return err
	}
	err := u.wecomkf.ReceptionistDel(wecomclient.ReceptionistOptions{
		OpenKFID:   kfid,
		UserIDList: req.UserIDList,
	})
	if err != nil {
		return err
	}
	return nil
}

func (u *WecomLogic) ReceptionistList(kfid string) ([]wecomclient.ReceptionistList, error) {
	if err := u.initializeWecomClient(); err != nil {
		return nil, err
	}
	list, err := u.wecomkf.ReceptionistList(kfid)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (u *WecomLogic) CheckReceptionist(stafflist []string, receptionistlist []wecomclient.ReceptionistList) error {
	if err := u.initializeWecomClient(); err != nil {
		return err
	}
	var invalidUserIDs []string
	receptionistIDSet := make(map[string]struct{}, len(receptionistlist))
	for _, receptionist := range receptionistlist {
		receptionistIDSet[receptionist.UserID] = struct{}{}
	}
	staffs, err := staffRepo.List(commonRepo.WithUUIDsIn(stafflist))
	if err != nil {
		return err
	}
	for _, staff := range staffs {
		if _, exists := receptionistIDSet[staff.StaffID]; !exists {
			invalidUserIDs = append(invalidUserIDs, staff.StaffID)
		}
	}
	if len(invalidUserIDs) > 0 {
		return fmt.Errorf("ErrStaffIDLIST: %v", invalidUserIDs)
	}
	return nil
}

func (u *WecomLogic) AccountList() ([]wecomclient.AccountInfoSchema, error) {
	if err := u.initializeWecomClient(); err != nil {
		return nil, err
	}
	accountinfolist, err := u.wecomkf.AccountList()
	if err != nil {
		return nil, err
	}
	return accountinfolist, nil
}

func (u *WecomLogic) AddContactWay(kfid string) (string, error) {
	if err := u.initializeWecomClient(); err != nil {
		return "", err
	}
	url, err := u.wecomkf.AddContactWay(kfid)
	if err != nil {
		return "", err
	}
	return url, nil
}

// Helper
func (u *WecomLogic) initializeWecomClient() (err error) {
	if u.wecomkf != nil {
		return
	}
	conf, err := wecomRepo.Get(wecomRepo.WithType("app"))
	if err != nil {
		global.ZAPLOG.Error("failed to get Wecom config", zap.Error(err))
		return
	}
	if conf.CorpID == "" || conf.AgentID == "" || conf.EncodingAESKey == "" || conf.Secret == "" || conf.Token == "" {
		global.ZAPLOG.Error("wecom config is nil")
		return
	}
	u.wecomkf, err = wecom.NewWecomKFClient(wecomclient.WecomConfig{
		CorpID:         conf.CorpID,
		Token:          conf.Token,
		EncodingAESKey: conf.EncodingAESKey,
		Secret:         conf.Secret,
		AgentID:        conf.AgentID,
	})
	if err != nil {
		global.ZAPLOG.Error("failed to initialize WecomKFClient", zap.Error(err))
		return
	}
	return
}

func (u *WecomLogic) processMessage(msg wecomclient.Message) (err error) {
	switch msg.MsgType {
	case wecomclient.WecomMsgTypeText:
		textMessage, err := msg.GetTextMessage()
		if err != nil {
			return err
		}
		global.ZAPLOG.Debug("", zap.Any("", textMessage))
		return u.handleTextMessage(wecomclient.Text{Text: textMessage})
	case wecomclient.WecomMsgTypeEvent:
		switch msg.EventType {
		case wecomclient.WecomMsgTypeEnterSessionEvent:
			enterSessionEvent, err := msg.GetEnterSessionEvent()
			if err != nil {
				return err
			}
			global.ZAPLOG.Debug("", zap.Any("", enterSessionEvent))
			return u.handleEnterSessionEvent(wecomclient.EnterSessionEvent{EnterSessionEvent: enterSessionEvent})
		case wecomclient.WecomMsgTypeMsgSendFailEvent:
			msgSendFailEvent, err := msg.GetMsgSendFailEvent()
			if err != nil {
				return err
			}
			global.ZAPLOG.Debug("", zap.Any("", msgSendFailEvent))
			return err
		case wecomclient.WecomMsgTypeSessionStatusChangeEvent:
			sessionStatusChangeEvent, err := msg.GetSessionStatusChangeEvent()
			if err != nil {
				return err
			}
			global.ZAPLOG.Debug("", zap.Any("", sessionStatusChangeEvent))
			return u.handleSessionStatusChangeEvent(wecomclient.SessionStatusChangeEvent{SessionStatusChangeEvent: sessionStatusChangeEvent})
		}
	default:

	}
	return
}

func (u *WecomLogic) handleTextMessage(textMessage wecomclient.Text) (err error) {
	var options wecomclient.ServiceStateGetOptions
	options.OpenKFID = textMessage.OpenKFID
	options.ExternalUserID = textMessage.ExternalUserID
	serviceStateInfo, err := u.wecomkf.ServiceStateGet(options)
	if err != nil {
		global.ZAPLOG.Error(serviceStateInfo.Error(), zap.Error(err))
		return
	}
	switch serviceStateInfo.ServiceState {
	case wecomclient.SessionStatusNew:
		return u.handleNewSession(textMessage)
	case wecomclient.SessionStatusHandled:
		return u.handleBotSession(textMessage)
	case wecomclient.SessionStatusWaiting:

	case wecomclient.SessionStatusInProgress:
		return u.handleInProgress(textMessage)
	case wecomclient.SessionStatusEndedOrNotStarted:

	default:

	}
	return
}

func (u *WecomLogic) handleNewSession(textMessage wecomclient.Text) (err error) {
	var options wecomclient.ServiceStateTransOptions
	options.OpenKFID = textMessage.OpenKFID
	options.ExternalUserID = textMessage.ExternalUserID
	options.ServiceState = wecomclient.SessionStatusHandled
	serviceStateRespInfo, err := u.wecomkf.ServiceStateTrans(options)
	if err != nil {
		global.ZAPLOG.Error(serviceStateRespInfo.Error(), zap.Error(err))
		return
	}
	return u.handleBotSession(textMessage)
}

func (u *WecomLogic) handleBotSession(textMessage wecomclient.Text) (err error) {
	kFInfo, err := kFRepo.Get(kFRepo.WithKFID(textMessage.OpenKFID))
	if err != nil {
		return
	}
	switch kFInfo.Status {
	case constant.KFStatusRobotToHuman:
		for _, keyword := range strings.Split(kFInfo.TransferKeywords, ";") {
			if textMessage.Text.Text.Content == keyword {
				return u.handleTransferToStaff(textMessage, kFInfo)
			}
		}
		return u.handleBotReply(textMessage)
	case constant.KFStatusOnlyRobot:
		return u.handleBotReply(textMessage)
	case constant.KFStatusOnlyHuman:

	default:

	}
	return
}

func (u *WecomLogic) handleBotReply(textMessage wecomclient.Text) (err error) {
	kFInfo, err := kFRepo.Get(kFRepo.WithKFID(textMessage.OpenKFID))
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(kFInfo.BotTimeout)*time.Second)
	defer cancel()
	resultChan := make(chan string)
	errorChan := make(chan error)
	var sendTextMsgOptions wecomclient.SendTextMsgOptions
	sendTextMsgOptions.Touser = textMessage.ExternalUserID
	sendTextMsgOptions.OpenKfid = textMessage.OpenKFID
	go func() {
		message := textMessage.Text.Text.Content + kFInfo.BotPrompt
		llmappLogic := NewILLMAppLogic()
		fullContent, err := llmappLogic.ChatMessage(textMessage.ExternalUserID, kFInfo.BotID, message)
		if err != nil {
			errorChan <- err
		} else {
			resultChan <- fullContent
		}
	}()
	select {
	case <-ctx.Done():
		sendTextMsgOptions.Text.Content = kFInfo.BotTimeoutMsg
		err = u.wecomkf.SendTextMsg(sendTextMsgOptions)
		if err != nil {
			global.ZAPLOG.Error(i18n.Tf("wecom.failed_action", "SendTextMsg"), zap.Error(err))
		}
		select {
		case fullContent := <-resultChan:
			sendTextMsgOptions.Text.Content = fullContent
			if err = u.wecomkf.SendTextMsg(sendTextMsgOptions); err != nil {
				global.ZAPLOG.Error(i18n.Tf("wecom.failed_action", "SendTextMsg"), zap.Error(err))
				return
			}
			return
		case err = <-errorChan:
			global.ZAPLOG.Error(i18n.Tf("wecom.failed_action", "ChatMessage"), zap.Error(err))
			return
		}
	case fullContent := <-resultChan:
		sendTextMsgOptions.Text.Content = fullContent
		if err = u.wecomkf.SendTextMsg(sendTextMsgOptions); err != nil {
			return
		}
		return
	case err = <-errorChan:
		global.ZAPLOG.Error(i18n.Tf("wecom.failed_action", "ChatMessage"), zap.Error(err))
		return
	}
}

func (u *WecomLogic) handleInProgress(textMessage wecomclient.Text) (err error) {
	kFInfo, err := kFRepo.Get(kFRepo.WithKFID(textMessage.OpenKFID))
	if err != nil {
		return
	}
	chatkey := constant.KeyWecomKHStaffPrefix + textMessage.ExternalUserID + ":" + textMessage.ReceptionistUserID
	if err = global.RDS.Set(context.Background(), chatkey, 1, time.Duration(kFInfo.ChatTimeout)*time.Second).Err(); err != nil {
		global.ZAPLOG.Error("redis set error", zap.Error(err))
		return
	}
	return
}

func (u *WecomLogic) handleEnterSessionEvent(enterSessionEvent wecomclient.EnterSessionEvent) (err error) {
	var options wecomclient.ServiceStateGetOptions
	options.OpenKFID = enterSessionEvent.OpenKFID
	options.ExternalUserID = enterSessionEvent.ExternalUserID
	serviceStateInfo, err := u.wecomkf.ServiceStateGet(options)
	if err != nil {
		global.ZAPLOG.Error(serviceStateInfo.Error(), zap.Error(err))
		return
	}
	kFInfo, err := kFRepo.Get(kFRepo.WithKFID(enterSessionEvent.OpenKFID))
	if err != nil {
		return err
	}
	if serviceStateInfo.ServiceState == wecomclient.SessionStatusHandled {
		return u.wecomkf.SendMenuMsg(wecomclient.SendMenuMsgOptions{
			KFID:           enterSessionEvent.OpenKFID,
			KHID:           enterSessionEvent.ExternalUserID,
			MenuMsgOptions: parseMenuText(kFInfo.BotWelcomeMsg),
		})
	}
	return u.wecomkf.SendMenuMsgOnEvent(wecomclient.SendMenuMsgOnEventOptions{
		Credential:     enterSessionEvent.Event.WelcomeCode,
		MenuMsgOptions: parseMenuText(kFInfo.BotWelcomeMsg),
	})
}

func (u *WecomLogic) handleSessionStatusChangeEvent(sessionStatusChangeEvent wecomclient.SessionStatusChangeEvent) (err error) {
	kFInfo, err := kFRepo.Get(kFRepo.WithKFID(sessionStatusChangeEvent.OpenKFID))
	if err != nil {
		return err
	}
	switch sessionStatusChangeEvent.Event.ChangeType {
	case wecomclient.WecomEventChangeTypeJoinSession:
	case wecomclient.WecomEventChangeTypeTransferSession:
		chatkey := constant.KeyWecomKHStaffPrefix + sessionStatusChangeEvent.ExternalUserID + ":" + sessionStatusChangeEvent.Event.OldReceptionistUserID
		if err := global.RDS.Del(context.Background(), chatkey).Err(); err != nil {
			global.ZAPLOG.Error(i18n.Tf("wecom.failed_action", "RDS.Del chatkey"), zap.Error(err))
			return err
		}
		isStaffWorkByStaffID(sessionStatusChangeEvent.ReceptionistUserID)
		var serviceStateTransOptions wecomclient.ServiceStateTransOptions
		serviceStateTransOptions.OpenKFID = sessionStatusChangeEvent.OpenKFID
		serviceStateTransOptions.ExternalUserID = sessionStatusChangeEvent.ExternalUserID
		serviceStateTransOptions.ServicerUserID = sessionStatusChangeEvent.ReceptionistUserID
		return u.handleServiceStateTransInProgress(serviceStateTransOptions)
	case wecomclient.WecomEventChangeTypeEndSession:
		if err := u.wecomkf.SendMenuMsgOnEvent(wecomclient.SendMenuMsgOnEventOptions{
			Credential:     sessionStatusChangeEvent.Event.MsgCode,
			MenuMsgOptions: parseMenuText(kFInfo.ChatendMsg),
		}); err != nil {
			global.ZAPLOG.Error(i18n.Tf("wecom.failed_action", "SendMenuMsgOnEvent"), zap.Error(err))
		}
		chatkey := constant.KeyWecomKHStaffPrefix + sessionStatusChangeEvent.ExternalUserID + ":" + sessionStatusChangeEvent.Event.OldReceptionistUserID
		if err = global.RDS.Del(context.Background(), chatkey).Err(); err != nil {
			global.ZAPLOG.Error(i18n.Tf("wecom.failed_action", "RDS.Del chatkey"), zap.Error(err))
			return
		}
		return
	case wecomclient.WecomEventChangeTypeRejoinSession:
		global.ZAPLOG.Debug("WecomEventChangeTypeRejoinSession")
		isStaffWorkByStaffID(sessionStatusChangeEvent.Event.NewReceptionistUserID)
		var serviceStateTransOptions wecomclient.ServiceStateTransOptions
		serviceStateTransOptions.OpenKFID = sessionStatusChangeEvent.OpenKFID
		serviceStateTransOptions.ExternalUserID = sessionStatusChangeEvent.ExternalUserID
		serviceStateTransOptions.ServicerUserID = sessionStatusChangeEvent.Event.NewReceptionistUserID
		return u.handleServiceStateTransInProgress(serviceStateTransOptions)
	default:
		global.ZAPLOG.Debug(i18n.T("wecom.unknow_msgtype"))
	}
	return
}

func (u *WecomLogic) handleServiceStateTransInProgress(serviceStateTransOptions wecomclient.ServiceStateTransOptions) (err error) {
	kFInfo, err := kFRepo.Get(kFRepo.WithKFID(serviceStateTransOptions.OpenKFID))
	if err != nil {
		return
	}
	serviceStateTransOptions.ServiceState = wecomclient.SessionStatusInProgress
	serviceStateRespInfo, err := u.wecomkf.ServiceStateTrans(serviceStateTransOptions)
	if err != nil {
		global.ZAPLOG.Error(serviceStateRespInfo.Error(), zap.Error(err))
		return
	}
	if serviceStateRespInfo.MsgCode != "" {
		if err = u.wecomkf.SendTextMsgOnEvent(wecomclient.SendTextMsgOnEventOptions{
			Message:    kFInfo.StaffWelcomeMsg,
			Credential: serviceStateRespInfo.MsgCode,
		}); err != nil {
			return
		}
	}
	global.ZAPLOG.Info("更新客户的上一次接待人员为选出的接待人员")
	if err := kHRepo.UpdatebyKHID(model.KH{
		KHID:    serviceStateTransOptions.ExternalUserID,
		StaffID: serviceStateTransOptions.ServicerUserID,
	}); err != nil {
		global.ZAPLOG.Error("更新客户的上一次接待人员失败", zap.Error(err))
	}
	global.ZAPLOG.Debug("降低该接待人员空闲权重")
	ctx := context.Background()
	staffweightkey := constant.KeyStaffWeightPrefix + serviceStateTransOptions.ServicerUserID
	if err = global.RDS.Decr(ctx, staffweightkey).Err(); err != nil {
		return
	}
	global.ZAPLOG.Debug("缓存会话超时时间")
	chatkey := constant.KeyWecomKHStaffPrefix + serviceStateTransOptions.ExternalUserID + ":" + serviceStateTransOptions.ServicerUserID
	if err = global.RDS.Set(ctx, chatkey, 1, time.Duration(kFInfo.ChatTimeout)*time.Second).Err(); err != nil {
		global.ZAPLOG.Error("redis set error", zap.Error(err))
		return
	}
	global.ZAPLOG.Debug("handleChatTimeoutTask")
	monitorCtx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	go func(ctx context.Context) {
		defer func() {
			cancel()
			global.ZAPLOG.Debug("cancel()")
		}()
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if exists, err := global.RDS.Exists(ctx, chatkey).Result(); err == nil && exists > 0 {
					continue
				}
				u.handleChatTimeout(serviceStateTransOptions)
				if err := global.RDS.Incr(ctx, staffweightkey).Err(); err != nil {
					global.ZAPLOG.Error("Incr", zap.Error(err))
					return
				}
				global.ZAPLOG.Debug("handleChatTimeout", zap.String("staffweightkey:", staffweightkey))
				return
			case <-ctx.Done():
				global.ZAPLOG.Debug("ctx.Done()")
				return
			}
		}
	}(monitorCtx)
	return
}

func (u *WecomLogic) handleChatTimeout(serviceStateTransOptions wecomclient.ServiceStateTransOptions) {
	var getOptions wecomclient.ServiceStateGetOptions
	getOptions.OpenKFID = serviceStateTransOptions.OpenKFID
	getOptions.ExternalUserID = serviceStateTransOptions.ExternalUserID
	serviceStateInfo, err := u.wecomkf.ServiceStateGet(getOptions)
	if err != nil {
		global.ZAPLOG.Error(serviceStateInfo.Error(), zap.Error(err))
		return
	}
	if serviceStateInfo.ServiceState == wecomclient.SessionStatusEndedOrNotStarted {
		return
	}
	var cursor uint64
	chatkey := constant.KeyWecomKHStaffPrefix + serviceStateTransOptions.ExternalUserID
	result, _, err := global.RDS.Scan(context.Background(), cursor, chatkey+"*", 10).Result()
	if err != nil {
		return
	}
	if len(result) > 0 {
		return
	}
	serviceStateTransOptions.ServiceState = wecomclient.SessionStatusEndedOrNotStarted
	serviceStateRespInfo, err := u.wecomkf.ServiceStateTrans(serviceStateTransOptions)
	if err != nil {
		global.ZAPLOG.Error(serviceStateRespInfo.Error(), zap.Error(err))
		return
	}
	global.ZAPLOG.Debug("会话超时，已变更状态")
	kFInfo, err := kFRepo.Get(kFRepo.WithKFID(serviceStateTransOptions.OpenKFID))
	if err != nil {
		return
	}
	if err := u.wecomkf.SendMenuMsgOnEvent(wecomclient.SendMenuMsgOnEventOptions{
		Credential:     serviceStateRespInfo.MsgCode,
		MenuMsgOptions: parseMenuText(kFInfo.ChatendMsg),
	}); err != nil {
		global.ZAPLOG.Error("SendTextMsgOnEvent", zap.Error(err))
	}
}

func (u *WecomLogic) handleTransferToStaff(textMessage wecomclient.Text, kFInfo model.KF) (err error) {
	if kFInfo.ReceivePriority == 1 {
		kHInfo, err := kHRepo.Get(kHRepo.WithKHID(textMessage.ExternalUserID))
		if err != nil {
			if err := kHRepo.Create(model.KH{KHID: textMessage.ExternalUserID}); err != nil {
				global.ZAPLOG.Error(i18n.Tf("wecom.failed_action", "Create"), zap.Error(err))
				return err
			}
		}
		if kHInfo.StaffID != "" {
			for _, staffInfo := range kFInfo.Staffs {
				if staffInfo.StaffID == kHInfo.StaffID {
					isStaffWork, err := isStaffWorkByStaffID(kHInfo.StaffID)
					if err != nil {
						global.ZAPLOG.Error(i18n.Tf("wecom.failed_action", "isStaffWorkByStaffID"), zap.Error(err))
						return err
					}
					if isStaffWork {
						var serviceStateTransOptions wecomclient.ServiceStateTransOptions
						serviceStateTransOptions.OpenKFID = textMessage.OpenKFID
						serviceStateTransOptions.ExternalUserID = textMessage.ExternalUserID
						serviceStateTransOptions.ServicerUserID = kHInfo.StaffID
						return u.handleServiceStateTransInProgress(serviceStateTransOptions)
					}
				}
			}
		}
	}
	switch kFInfo.ReceiveRule {
	case constant.KFReceiveRuleRoundRobin:
		global.ZAPLOG.Debug("ToDo")
	case constant.KFReceiveRuleIdle:
		var staffIDs []string
		for _, staffinfo := range kFInfo.Staffs {
			isStaffWork, err := isStaffWorkByStaffID(staffinfo.StaffID)
			if err != nil {
				return err
			}
			if isStaffWork {
				staffIDs = append(staffIDs, staffinfo.StaffID)
			}
		}
		selectedstaff, _ := getHighestWeightStaff(staffIDs)
		if selectedstaff != "" {
			var serviceStateTransOptions wecomclient.ServiceStateTransOptions
			serviceStateTransOptions.OpenKFID = textMessage.OpenKFID
			serviceStateTransOptions.ExternalUserID = textMessage.ExternalUserID
			serviceStateTransOptions.ServicerUserID = selectedstaff
			return u.handleServiceStateTransInProgress(serviceStateTransOptions)
		}
		global.ZAPLOG.Debug("没有找到空闲的接待人员")
		var sendTextMsgOptions wecomclient.SendTextMsgOptions
		sendTextMsgOptions.Touser = textMessage.ExternalUserID
		sendTextMsgOptions.OpenKfid = textMessage.OpenKFID
		sendTextMsgOptions.Text.Content = kFInfo.UnmannedMsg
		return u.wecomkf.SendTextMsg(sendTextMsgOptions)
	default:

	}
	return
}

func parseMenuText(text string) wecomclient.MenuMsgOptions {
	lines := strings.Split(text, "\n")
	headContent := ""
	tailContent := ""
	menuList := []wecomclient.MenuItem{}
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "#H "):
			headContent = strings.TrimPrefix(line, "#H ")
		case strings.HasPrefix(line, "#T "):
			tailContent = strings.TrimPrefix(line, "#T ")
		case strings.HasPrefix(line, "#TXT "):
			menuList = append(menuList, wecomclient.MenuItem{
				Type: "text",
				Text: &struct {
					Content   string `json:"content"`
					NoNewLine int    `json:"no_newline,omitempty"`
				}{
					Content:   strings.TrimPrefix(line, "#TXT "),
					NoNewLine: 0,
				},
			})
		case strings.HasPrefix(line, "#CLK "):
			menuList = append(menuList, wecomclient.MenuItem{
				Type: "click",
				Click: &struct {
					ID      string `json:"id,omitempty"`
					Content string `json:"content"`
				}{
					ID:      fmt.Sprintf("%d", idx+101),
					Content: strings.TrimPrefix(line, "#CLK "),
				},
			})
		case strings.HasPrefix(line, "#VIEW "):
			viewText := strings.TrimPrefix(line, "#VIEW ")
			openParen := strings.LastIndex(viewText, "(")
			closeParen := strings.LastIndex(viewText, ")")
			if openParen > -1 && closeParen > openParen {
				content := viewText[:openParen]
				url := viewText[openParen+1 : closeParen]
				menuList = append(menuList, wecomclient.MenuItem{
					Type: "view",
					View: &struct {
						URL     string `json:"url"`
						Content string `json:"content"`
					}{
						URL:     url,
						Content: content,
					},
				})
			}
		case strings.HasPrefix(line, "#MINIPROGRAM "):
			mpText := strings.TrimPrefix(line, "#MINIPROGRAM ")
			openParen := strings.LastIndex(mpText, "(")
			closeParen := strings.LastIndex(mpText, ")")
			if openParen > -1 && closeParen > openParen {
				content := mpText[:openParen-1]
				params := mpText[openParen+1 : closeParen]
				paramParts := strings.Split(params, ",")
				if len(paramParts) == 2 {
					appid := strings.TrimSpace(paramParts[0])
					pagepath := strings.TrimSpace(paramParts[1])
					menuList = append(menuList, wecomclient.MenuItem{
						Type: "miniprogram",
						MiniProgram: &struct {
							AppID    string `json:"appid"`
							PagePath string `json:"pagepath"`
							Content  string `json:"content"`
						}{
							AppID:    appid,
							PagePath: pagepath,
							Content:  content,
						},
					})
				}
			}
		}
	}
	var resp wecomclient.MenuMsgOptions
	if headContent != "" {
		resp.HeadContent = headContent
	}
	if menuList != nil {
		resp.MenuList = menuList
	}
	if tailContent != "" {
		resp.TailContent = tailContent
	}
	return resp
}
