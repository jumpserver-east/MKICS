package logic

import (
	"EvoBot/backend/app/dto"
	"EvoBot/backend/app/dto/request"
	"EvoBot/backend/app/dto/response"
	"EvoBot/backend/app/model"
	"EvoBot/backend/global"
	"EvoBot/backend/utils/wecom"
	wecomclient "EvoBot/backend/utils/wecom/client"
	"fmt"
	"strings"

	"context"
	"time"

	"go.uber.org/zap"
)

type WecomLogic struct {
	wecomkf wecom.WecomKFClient
}

type IWecomLogic interface {
	VerifyURL(req dto.SignatureOptions) (string, error)
	Handle(body []byte) error

	ConfigAppUpdate(req request.WecomConfigApp) error
	ConfigAppGet() (response.WecomConfigApp, error)

	ReceptionistAdd(req request.ReceptionistOptions) error
	ReceptionistDel(req request.ReceptionistOptions) error
	ReceptionistList(kfid string) ([]wecomclient.ReceptionistList, error)
	CheckReceptionist(stafflist []string, receptionistlist []wecomclient.ReceptionistList) error

	AccountList() ([]wecomclient.AccountInfoSchema, error)
	AddContactWay(kfid string) (string, error)
}

func NewIWecomLogic() IWecomLogic {
	return &WecomLogic{}
}

func (u *WecomLogic) ConfigAppUpdate(req request.WecomConfigApp) error {
	conf, err := wecomRepo.Get(wecomRepo.WithType("app"))
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
	kf, err := u.loadWecomKFClientByFrom()
	if err != nil {
		return err
	}
	u.wecomkf = kf
	return nil
}

func (u *WecomLogic) ConfigAppGet() (response.WecomConfigApp, error) {
	var res response.WecomConfigApp
	conf, err := wecomRepo.Get(wecomRepo.WithType("app"))
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

func (u *WecomLogic) VerifyURL(req dto.SignatureOptions) (string, error) {
	if err := u.initializeWecomClient(); err != nil {
		return "", err
	}
	echo, err := u.wecomkf.VerifyURL(wecomclient.SignatureOptions{
		Signature: req.Signature,
		TimeStamp: req.TimeStamp,
		Nonce:     req.Nonce,
		EchoStr:   req.EchoStr,
	})
	if err != nil {
		global.ZAPLOG.Error("failed to check url", zap.Error(err))
	}
	return echo, nil
}

func (u *WecomLogic) Handle(body []byte) error {
	if err := u.initializeWecomClient(); err != nil {
		return err
	}
	msginfo, err := u.wecomkf.SyncMsg(body)
	if err != nil {
		global.ZAPLOG.Error("failed SyncMsg", zap.Error(err))
		return err
	}
	if msginfo.MessageID != "" || msginfo.KFID != "" || msginfo.KHID != "" || msginfo.StaffID != "" || msginfo.Message != "" || msginfo.MessageType != "" || msginfo.Credential != "" || msginfo.IsBot {
		global.ZAPLOG.Info("MessageInfo 内容:",
			zap.String("MessageID", msginfo.MessageID),
			zap.String("KFID", msginfo.KFID),
			zap.String("KHID", msginfo.KHID),
			zap.String("StaffID", msginfo.StaffID),
			zap.String("Message", msginfo.Message),
			zap.String("MessageType", msginfo.MessageType),
			zap.String("Credential", msginfo.Credential),
			zap.Bool("IsBot", msginfo.IsBot),
		)
	}
	return u.processMessage(msginfo)
}

func (u *WecomLogic) ReceptionistAdd(req request.ReceptionistOptions) error {
	if err := u.initializeWecomClient(); err != nil {
		return err
	}
	err := u.wecomkf.ReceptionistAdd(wecomclient.ReceptionistOptions{
		OpenKFID:   req.OpenKFID,
		UserIDList: req.UserIDList,
	})
	if err != nil {
		return err
	}
	return nil
}

func (u *WecomLogic) ReceptionistDel(req request.ReceptionistOptions) error {
	if err := u.initializeWecomClient(); err != nil {
		return err
	}
	err := u.wecomkf.ReceptionistDel(wecomclient.ReceptionistOptions{
		OpenKFID:   req.OpenKFID,
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
	receptionistIDSet := make(map[string]struct{})
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
		errorMessage := fmt.Sprintf("ErrStaffIDLIST: %s", invalidUserIDs)
		return fmt.Errorf(errorMessage)
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
func (u *WecomLogic) loadWecomKFClientByFrom() (wecom.WecomKFClient, error) {
	conf, err := wecomRepo.Get(wecomRepo.WithType("app"))
	if err != nil {
		global.ZAPLOG.Error("failed to get Wecom config", zap.Error(err))
		return nil, err
	}
	if conf.CorpID == "" || conf.AgentID == "" || conf.EncodingAESKey == "" || conf.Secret == "" || conf.Token == "" {
		global.ZAPLOG.Error("wecom config is nil", zap.Error(err))
	}
	kf, err := wecom.NewWecomKFClient(wecomclient.WecomConfig{
		CorpID:         conf.CorpID,
		Token:          conf.Token,
		EncodingAESKey: conf.EncodingAESKey,
		Secret:         conf.Secret,
		AgentID:        conf.AgentID,
	})
	if err != nil {
		global.ZAPLOG.Error("failed to New WecomKFClient", zap.Error(err))
		return nil, err
	}
	return kf, nil
}

func (u *WecomLogic) initializeWecomClient() error {
	if u.wecomkf != nil {
		return nil
	}
	var err error
	u.wecomkf, err = u.loadWecomKFClientByFrom()
	if err != nil {
		global.ZAPLOG.Error("failed to initialize WecomKFClient", zap.Error(err))
		return err
	}
	return nil
}

func (u *WecomLogic) processMessage(msginfo wecomclient.MessageInfo) error {
	if msginfo.IsBot {
		return u.handleBotMessage(msginfo)
	}
	if msginfo.KHID != "" {
		kfinfo, err := kFRepo.Get(kFRepo.WithKFID(msginfo.KFID))
		if err != nil {
			return err
		}
		err = global.RDS.Expire(context.Background(), msginfo.MessageID, time.Duration(kfinfo.ChatTimeout)*time.Second).Err()
		if err != nil {
			global.ZAPLOG.Error("重置缓存时间失败", zap.Error(err))
		}
	}
	return nil
}

func (u *WecomLogic) handleSuccessfulTransfer(msginfo wecomclient.MessageInfo, staffid string, kfinfo model.KF) error {
	global.ZAPLOG.Info("变更微信客服会话状态中")
	credential, err := u.wecomkf.ServiceStateTransToStaff(wecomclient.ServiceStateTransOptions{
		OpenKFID:       msginfo.KFID,
		ExternalUserID: msginfo.KHID,
		ServicerUserID: staffid,
	})
	if err != nil {
		global.ZAPLOG.Error("变更微信客服会话状态失败", zap.Error(err))
		currentTime := time.Now()
		formattedTime := currentTime.Format("2006-01-02 15:04:05")
		errormessage := "未知原因出现异常请联系工作人员或可尝试继续与 AI 对话，当前时间：" + formattedTime
		return u.wecomkf.SendTextMsg(wecomclient.SendTextMsgOptions{
			KFID:    msginfo.KFID,
			KHID:    msginfo.KHID,
			Message: errormessage,
		})
	}
	if err = u.wecomkf.SendTextMsgOnEvent(wecomclient.SendTextMsgOnEventOptions{
		Message:    kfinfo.StaffWelcomeMsg,
		Credential: credential,
	}); err != nil {
		return err
	}
	global.ZAPLOG.Info("变更微信客服会话状态成功")
	global.ZAPLOG.Info("更新客户模型的上一次接待人员")
	if err = kHRepo.UpdatebyKHID(model.KH{
		KHID:    msginfo.KHID,
		StaffID: staffid,
	}); err != nil {
		global.ZAPLOG.Error("更新失败，没有配置为客户上一次接待人员，下一次该客户将会被ReceiveRule匹配规则", zap.Error(err))
	}
	global.ZAPLOG.Info("降低该客服空闲权重")
	ctx := context.Background()
	weightkey := staffid + "-weight"
	if err = global.RDS.Decr(ctx, weightkey).Err(); err != nil {
		return err
	}
	global.ZAPLOG.Info("设置回复时间缓存，开始监控回复超时时间")
	if err = global.RDS.Set(ctx, msginfo.MessageID, 1, time.Duration(kfinfo.ChatTimeout)*time.Second).Err(); err != nil {
		global.ZAPLOG.Error("redis set error", zap.Error(err))
		return err
	}
	go u.monitorChat(ctx, kfinfo, msginfo.KHID, weightkey, msginfo.MessageID)
	return nil
}

func (u *WecomLogic) monitorChat(ctx context.Context, kfinfo model.KF, khid, weightkey, messageID string) {
	ticker := time.NewTicker(time.Duration(kfinfo.ChatTimeout) * time.Second / 2) // 半个超时时间
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if !u.hasKHReplied(ctx, messageID) {
				u.handleChatTimeout(ctx, kfinfo, khid, weightkey)
				return
			}
			global.ZAPLOG.Info("客户会话未过期")
		case <-ctx.Done():
			global.ZAPLOG.Info("监控已停止")
			return
		}
	}
}

func (u *WecomLogic) handleChatTimeout(ctx context.Context, kfinfo model.KF, khid, weightkey string) {
	credential, err := u.wecomkf.ServiceStateTransToEnd(wecomclient.ServiceStateTransOptions{
		OpenKFID:       kfinfo.KFID,
		ExternalUserID: khid,
	})
	if err != nil {
		global.ZAPLOG.Error("结束会话失败", zap.Error(err))
		return
	}
	global.ZAPLOG.Info("会话超时，已变更状态")
	if err := u.wecomkf.SendTextMsgOnEvent(wecomclient.SendTextMsgOnEventOptions{
		Message:    kfinfo.ChatendMsg,
		Credential: credential,
	}); err != nil {
		global.ZAPLOG.Error("SendTextMsgOnEvent", zap.Error(err))
	} else {
		global.ZAPLOG.Info("发送结束会话语", zap.String("ChatendMsg:", kfinfo.ChatendMsg))
	}
	if err := global.RDS.Incr(ctx, weightkey).Err(); err != nil {
		global.ZAPLOG.Error("Incr", zap.Error(err))
	} else {
		global.ZAPLOG.Info("结束会话变更权重", zap.String("weightkey:", weightkey))
	}
}

func (u *WecomLogic) hasKHReplied(ctx context.Context, messageID string) bool {
	exists, err := global.RDS.Exists(ctx, messageID).Result()
	if err != nil {
		global.ZAPLOG.Error("检查客户回复缓存失败", zap.Error(err))
		return false
	}
	if exists > 0 {
		return true
	}
	return false
}

func (u *WecomLogic) handleTransferToStaff(msginfo wecomclient.MessageInfo, kfinfo model.KF) error {
	if kfinfo.ReceivePriority == 1 {
		khinfo, err := kHRepo.Get(kHRepo.WithKHID(msginfo.KHID))
		if err != nil {
			global.ZAPLOG.Info("数据库没有找到该客户信息")
			if err := kHRepo.Create(model.KH{KHID: msginfo.KHID}); err != nil {
				return err
			}
			global.ZAPLOG.Info("录入客户信息:", zap.String("khid", msginfo.KHID))
		}
		if khinfo.StaffID != "" {
			isStaffWork, err := isStaffWorkByStaffID(khinfo.StaffID)
			if err != nil {
				return err
			}
			if isStaffWork {
				global.ZAPLOG.Info("选出的客服", zap.String("StaffID", khinfo.StaffID))
				return u.handleSuccessfulTransfer(msginfo, khinfo.StaffID, kfinfo)
			}
		}
		global.ZAPLOG.Info("该客户没有上一次接待人员", zap.String("KHID", msginfo.KHID))
	}
	if kfinfo.ReceiveRule == 1 {
		global.ZAPLOG.Info("轮流接待模式", zap.String("KFName", kfinfo.KFName))
		global.ZAPLOG.Info("未实现功能")
		return u.wecomkf.SendTextMsg(wecomclient.SendTextMsgOptions{
			KFID:    msginfo.KFID,
			KHID:    msginfo.KHID,
			Message: kfinfo.UnmannedMsg,
		})
	} else if kfinfo.ReceiveRule == 2 {
		global.ZAPLOG.Info("空闲接待模式", zap.String("KFName", kfinfo.KFName))
		var staffIDs []string
		for _, staffinfo := range kfinfo.Staffs {
			isStaffWork, err := isStaffWorkByStaffID(staffinfo.StaffID)
			if err != nil {
				return err
			}
			if isStaffWork {
				staffIDs = append(staffIDs, staffinfo.StaffID)
			}
		}
		selectedKF, weight := getHighestWeightKF(staffIDs)
		if selectedKF != "" {
			global.ZAPLOG.Info("选出的客服", zap.String("StaffID", selectedKF), zap.Int("权重", weight))
			err := u.handleSuccessfulTransfer(msginfo, selectedKF, kfinfo)
			if err != nil {
				return err
			}
		} else {
			global.ZAPLOG.Info("没有找到空闲的客服")
			return u.wecomkf.SendTextMsg(wecomclient.SendTextMsgOptions{
				KFID:    msginfo.KFID,
				KHID:    msginfo.KHID,
				Message: kfinfo.UnmannedMsg,
			})
		}
	} else {
		global.ZAPLOG.Info("默认接待模式未配置", zap.String("KFName", kfinfo.KFName))
		return nil
	}
	return nil
}

func getHighestWeightKF(staffIDs []string) (string, int) {
	maxWeight := -1
	selectedstaffid := ""
	for _, staffID := range staffIDs {
		key := "staffweight:" + staffID
		weight, err := global.RDS.Get(context.Background(), key).Int()
		if err != nil {
			global.ZAPLOG.Error("获取权重失败", zap.String("key", key), zap.Error(err))
			continue
		}
		if weight > maxWeight {
			maxWeight = weight
			selectedstaffid = staffID
		}
	}
	return selectedstaffid, maxWeight
}

func (u *WecomLogic) handleBotMessage(msginfo wecomclient.MessageInfo) error {
	kfinfo, err := kFRepo.Get(kFRepo.WithKFID(msginfo.KFID))
	if err != nil {
		global.ZAPLOG.Error("找不到与微信对接的客服，将默认使用机器人回复", zap.Error(err))
		return u.handleBotReply(msginfo, kfinfo)
	}
	if msginfo.MessageType == "enter_session" {
		if msginfo.Credential != "" {
			options := parseMenuTextToOptions(kfinfo.BotWelcomeMsg, msginfo.Credential)
			return u.wecomkf.SendMenuMsgOnEvent(options)
		}
		return nil
	}
	switch kfinfo.Status {
	case 1:
		keywords := strings.Split(kfinfo.TransferKeywords, ";")
		for _, keyword := range keywords {
			if msginfo.Message == keyword {
				global.ZAPLOG.Info("客户触发客服转人工关键字", zap.String("keyword", keyword))
				return u.handleTransferToStaff(msginfo, kfinfo)
			}
		}
		return u.handleBotReply(msginfo, kfinfo)
	case 2:
		global.ZAPLOG.Info("客服为仅机器人模式，无法转接人工", zap.String("KFName", kfinfo.KFName))
		return u.handleBotReply(msginfo, kfinfo)
	case 3:
		global.ZAPLOG.Info("客服为仅人工模式，机器人无法回复消息", zap.String("KFName", kfinfo.KFName))
		return nil
	default:
		global.ZAPLOG.Error("该客服模式不存在")
		return err
	}
}

func (u *WecomLogic) handleBotReply(msginfo wecomclient.MessageInfo, kfinfo model.KF) error {
	maxkbLogic := NewIMaxkbLogic() // 依赖，后续需要解耦
	chatid, err := maxkbLogic.getChatIdByKH(msginfo.KHID)
	if err != nil {
		global.ZAPLOG.Error("getChatIdByKH", zap.Error(err))
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(kfinfo.BotTimeout)*time.Second)
	defer cancel()
	resultChan := make(chan string)
	errorChan := make(chan error)
	go func() {
		fullContent, err := maxkbLogic.ChatMessage(chatid, msginfo.Message)
		if err != nil {
			errorChan <- err
		} else {
			resultChan <- fullContent
		}
	}()
	select {
	case <-ctx.Done():
		err := u.wecomkf.SendTextMsg(wecomclient.SendTextMsgOptions{
			KFID:    msginfo.KFID,
			KHID:    msginfo.KHID,
			Message: kfinfo.BotTimeoutMsg,
		})
		if err != nil {
			global.ZAPLOG.Error("failed to send wait message", zap.Error(err))
		}
		global.ZAPLOG.Info("bot超时消息:", zap.String("content", kfinfo.BotTimeoutMsg))
		select {
		case fullContent := <-resultChan:
			err = u.wecomkf.SendTextMsg(wecomclient.SendTextMsgOptions{
				KFID:    msginfo.KFID,
				KHID:    msginfo.KHID,
				Message: fullContent,
			})
			if err != nil {
				return err
			}
			global.ZAPLOG.Info("bot消息:", zap.String("fullContent", fullContent))
			return nil
		case err := <-errorChan:
			global.ZAPLOG.Error("failed ChatMessage after timeout", zap.Error(err))
			return err
		}
	case fullContent := <-resultChan:
		err = u.wecomkf.SendTextMsg(wecomclient.SendTextMsgOptions{
			KFID:    msginfo.KFID,
			KHID:    msginfo.KHID,
			Message: fullContent,
		})
		if err != nil {
			return err
		}
		global.ZAPLOG.Info("bot消息:", zap.String("fullContent", fullContent))
		return nil
	case err := <-errorChan:
		global.ZAPLOG.Error("failed ChatMessage", zap.Error(err))
		return err
	}
}

func parseMenuTextToOptions(text, credential string) wecomclient.SendMenuMsgOnEventOptions {
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
	var resp wecomclient.SendMenuMsgOnEventOptions
	if headContent != "" {
		resp.HeadContent = headContent
	}
	if menuList != nil {
		resp.MenuList = menuList
	}
	if tailContent != "" {
		resp.TailContent = tailContent
	}
	resp.Credential = credential
	return resp
}

func isStaffWorkByStaffID(staffid string) (bool, error) {
	staffinfo, err := staffRepo.Get(staffRepo.WithStaffID(staffid)) //同kfinfo
	if err != nil {
		return false, err
	}
	for _, policy := range staffinfo.Policies {
		if isWithinTime(policy.Repeat, policy.Week) {
			for _, worktime := range policy.WorkTimes {
				if isTimeInRange(worktime.StartTime, worktime.EndTime) {
					ctx := context.Background()
					key := "staffweight:" + staffid
					_, err := global.RDS.Get(ctx, key).Result()
					if err != nil {
						global.ZAPLOG.Error("redis get error", zap.Error(err))
						global.ZAPLOG.Info("初始化权重缓存")
						err = global.RDS.Set(ctx, key, policy.MaxCount, 0).Err()
						if err != nil {
							global.ZAPLOG.Error("redis set error", zap.Error(err))
							return false, err
						}
					}
					return true, nil
				}
			}
			global.ZAPLOG.Info("该接待人员的工作时间导致目前无法接待", zap.String("StaffName", staffinfo.StaffName))
			return false, nil
		} else {
			global.ZAPLOG.Info("该接待人员的策略导致目前无法接待", zap.String("StaffName", staffinfo.StaffName))
			return false, nil
		}
	}
	return false, nil
}

func isTimeInRange(startTimeStr, endTimeStr string) bool {
	layout := "15:04:05"
	startTime, err1 := time.Parse(layout, startTimeStr)
	endTime, err2 := time.Parse(layout, endTimeStr)
	if err1 != nil || err2 != nil {
		global.ZAPLOG.Error("Error parsing time:", zap.Error(err1))
		global.ZAPLOG.Error("Error parsing time:", zap.Error(err2))
		return false
	}
	now := time.Now()
	currentTime, _ := time.Parse(layout, now.Format(layout))
	if currentTime.After(startTime) && currentTime.Before(endTime) {
		return true
	}
	return false
}

func isWithinTime(repeat int, week string) bool {
	currentTime := time.Now()
	currentWeekday := int(currentTime.Weekday()) // 0: Sunday, 1: Monday, ..., 6: Saturday
	switch repeat {
	case 1:
		if len(week) != 7 {
			return false
		}
		return week[currentWeekday] == '1'
	case 2:
		return true
	case 3:
		return currentWeekday >= 1 && currentWeekday <= 5
	case 4:
		// 法定工作日有效，跳过法定节假日（此处需要法定节假日的额外数据支持，简单实现为周一到周五）
		// 简单假设法定工作日是周一到周五，实际情况可能要处理节假日
		return currentWeekday >= 1 && currentWeekday <= 5
	case 5:
		// 法定节假日有效，跳过法定工作日（此处需要法定节假日的额外数据支持）
		// 简单假设法定节假日是周六和周日，实际情况可能要处理节假日
		return currentWeekday == 0 || currentWeekday == 6
	default:
		return false
	}
}
