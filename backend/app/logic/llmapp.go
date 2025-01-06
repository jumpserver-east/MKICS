package logic

import (
	"EvoBot/backend/app/dto"
	"EvoBot/backend/app/model"
	"EvoBot/backend/global"
	"EvoBot/backend/utils/llmapp"

	"go.uber.org/zap"
)

type LLMAppLogic struct {
	LLMAppConf dto.LLMAppConfig
}

type ILLMAppLogic interface {
	ConfigAdd(req dto.LLMAppConfig) error
	ConfigUpdate(uuid string, req dto.LLMAppConfig) error
	ConfigDel(uuid string) error
	ConfigGet(uuid string) (dto.LLMAppConfig, error)
	ConfigList() ([]dto.LLMAppConfig, error)
}

func NewILLMAppLogic() ILLMAppLogic {
	return &LLMAppLogic{}
}

func (u *LLMAppLogic) ConfigAdd(req dto.LLMAppConfig) error {
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	defer tx.Rollback()
	if err := llmappRepo.Create(model.LLMAppConfig{
		LLMAppType: req.LLMAppType,
		ApiKey:     req.ApiKey,
		BaseURL:    req.BaseURL,
		ConfigName: req.ConfigName,
	}); err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (u *LLMAppLogic) ConfigDel(uuid string) error {
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	defer tx.Rollback()
	if err := llmappRepo.Delete(commonRepo.WithByUUID(uuid)); err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (u *LLMAppLogic) ConfigUpdate(uuid string, req dto.LLMAppConfig) error {
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	defer tx.Rollback()
	conf, err := llmappRepo.Get(commonRepo.WithByUUID(uuid))
	if err != nil {
		return err
	}
	conf.UUID = uuid
	conf.LLMAppType = req.LLMAppType
	conf.ApiKey = req.ApiKey
	conf.BaseURL = req.BaseURL
	if err := llmappRepo.Update(conf); err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (u *LLMAppLogic) ConfigGet(uuid string) (dto.LLMAppConfig, error) {
	conf, err := llmappRepo.Get(commonRepo.WithByUUID(uuid))
	var resp dto.LLMAppConfig
	if err != nil {
		return resp, err
	}
	resp.LLMAppType = conf.LLMAppType
	resp.ConfigName = conf.ConfigName
	resp.BaseURL = conf.BaseURL
	resp.ApiKey = conf.ApiKey
	return resp, nil
}

func (u *LLMAppLogic) ConfigList() ([]dto.LLMAppConfig, error) {
	configs, err := llmappRepo.List()
	var resp []dto.LLMAppConfig
	if err != nil {
		global.ZAPLOG.Error("failed to get LLMApp config", zap.Error(err))
		return resp, err
	}
	for _, config := range configs {
		var res dto.LLMAppConfig
		res.LLMAppType = config.LLMAppType
		res.ConfigName = config.ConfigName
		res.BaseURL = config.BaseURL
		res.ApiKey = config.ApiKey
		resp = append(resp, res)
	}

	return resp, nil
}

func GetChatIdByKH(khid string, llmapp llmapp.LLMAppClient) (string, error) {
	khinfo, err := kHRepo.Get(kHRepo.WithKHID(khid))
	if err != nil {
		global.ZAPLOG.Error("数据库没有找到该客户信息", zap.Error(err))
		if err := kHRepo.Create(model.KH{
			KHID: khid,
		}); err != nil {
			return "", err
		}
		global.ZAPLOG.Info("录入客户信息:", zap.String("KHID", khid))
	}
	if khinfo.ChatID == "" {
		global.ZAPLOG.Info("该客户没有和机器人的聊天ID", zap.String("KHID", khid))
		newchatid, err := llmapp.ChatOpen()
		if err != nil {
			return "", err
		}
		kh := model.KH{
			KHID:   khid,
			ChatID: newchatid,
		}
		if err := kHRepo.UpdatebyKHID(kh); err != nil {
			return "", err
		}
		global.ZAPLOG.Info("生成聊天id更新客户信息:", zap.String("newchatid", newchatid), zap.String("khid", khid))
		return newchatid, nil
	}
	return khinfo.ChatID, nil
}
