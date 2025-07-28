package llmapp

import (
	"MKICS/backend/constant"
	"MKICS/backend/utils/llmapp/client"
	"errors"
)

type LLMAppClient interface {
	ChatMessage(message string, chatid *string, asker string) (string, error)
	ChatOpen() (*string, error)
}

func NewLLMAppClient(llmappType string, vars map[string]interface{}) (LLMAppClient, error) {
	switch llmappType {
	case constant.Maxkb:
		return client.NewMaxkbClient(vars)
	default:
		return nil, errors.New("not support type")
	}
}
