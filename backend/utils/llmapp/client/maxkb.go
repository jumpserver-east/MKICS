package client

import (
	"MKICS/backend/global"

	"github.com/Ewall555/MaxKB-golang-sdk/api/application"
	mkreq "github.com/Ewall555/MaxKB-golang-sdk/api/request"   // 请求参数
	mkresp "github.com/Ewall555/MaxKB-golang-sdk/api/response" // 返回参数
	mkconfig "github.com/Ewall555/MaxKB-golang-sdk/config"     // 配置参数
	mk "github.com/Ewall555/MaxKB-golang-sdk/maxkb"            // 引入包
)

type MaxKB struct {
	MaxKBClient     *mk.MaxKB
	ApplicationChat *application.ApplicationChat
}

func NewMaxkbClient(vars map[string]interface{}) (*MaxKB, error) {
	base_url := loadParamFromVars("base_url", vars)
	api_key := loadParamFromVars("api_key", vars)
	maxKBClient := mk.NewMaxKB(&mkconfig.Config{
		BaseURL: base_url,
		ApiKey:  api_key,
	})
	applicationChat := maxKBClient.GetApplicationChat()
	return &MaxKB{
		MaxKBClient:     maxKBClient,
		ApplicationChat: applicationChat,
	}, nil
}

type StreamCallback struct {
	StreamCallback func(*mkresp.Chat_messagePostStreamResponse)
}

func (c *MaxKB) ChatMessage(message string, chatid *string, asker string) (string, error) {
	formData := map[string]interface{}{
		"asker": asker,
	}
	req := mkreq.Chat_messagePostRequest{
		Message:  message,
		ReChat:   false,
		Stream:   false,
		FormData: formData,
	}
	resp, err := c.ApplicationChat.Chat_messageByChat_id(req, chatid, nil)
	if err != nil {
		global.ZAPLOG.Error(err.Error())
		return "", err
	}
	return resp.Content, nil
}

func (c *MaxKB) ChatOpen() (*string, error) {
	profileresp, err := c.ApplicationChat.Profile()
	if err != nil {
		global.ZAPLOG.Error(err.Error())
		return nil, err
	}
	chatid, err := c.ApplicationChat.ChatOpenByApplication_id(profileresp.ID)
	if err != nil {
		global.ZAPLOG.Error(err.Error())
		return nil, err
	}
	return chatid, err
}
