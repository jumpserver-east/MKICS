package client

import (
	"EvoBot/backend/global"

	mkreq "github.com/Ewall555/MaxKB-golang-sdk/api/request"   // 请求参数
	mkresp "github.com/Ewall555/MaxKB-golang-sdk/api/response" // 返回参数
	mk "github.com/Ewall555/MaxKB-golang-sdk/maxkb"            // 引入包
	"go.uber.org/zap"
)

type MaxKB struct {
	MaxKBClient *mk.MaxKB
}

func NewMaxkbClient(vars map[string]interface{}) (*MaxKB, error) {
	base_url := loadParamFromVars("base_url", vars)
	api_key := loadParamFromVars("api_key", vars)
	maxKBClient := mk.New(base_url, api_key)
	return &MaxKB{
		MaxKBClient: maxKBClient,
	}, nil
}

type StreamCallback struct {
	StreamCallback func(*mkresp.Chat_messagePostStreamResponse)
}

func (c *MaxKB) ChatMessage(message string, chatid *string) (string, error) {
	req := mkreq.Chat_messagePostRequest{
		Message: message,
		ReChat:  false,
		Stream:  false,
	}
	resp, err := c.MaxKBClient.ApplicationChat.Chat_messageByChat_id(req, chatid, nil)
	if err != nil {
		panic(err)
	}
	return resp.Content, nil
}

func (c *MaxKB) ChatOpen() (*string, error) {
	profileresp, err := c.MaxKBClient.ApplicationChat.Profile()
	if err != nil {
		global.ZAPLOG.Error("get profile error", zap.Error(err))
		return nil, err
	}
	chatid, err := c.MaxKBClient.ApplicationChat.ChatOpenByApplication_id(profileresp.ID)
	if err != nil {
		global.ZAPLOG.Error("open chat error", zap.Error(err))
		return nil, err
	}
	return chatid, err
}
