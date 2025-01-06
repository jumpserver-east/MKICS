package client

import (
	mkreq "github.com/Ewall555/MaxKB-golang-sdk/api/request"   // 请求参数
	mkresp "github.com/Ewall555/MaxKB-golang-sdk/api/response" // 返回参数
	mk "github.com/Ewall555/MaxKB-golang-sdk/maxkb"            // 引入包
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

func (c *MaxKB) ChatOpen() (string, error) {
	profileresp, err := c.MaxKBClient.ApplicationChat.Profile()
	if err != nil {
		return "", err
	}
	chatid, err := c.MaxKBClient.ApplicationChat.ChatOpenByApplication_id(profileresp.ID)
	if err != nil {
		return "", err
	}
	return *chatid, err
}
