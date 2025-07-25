package response

import "MKICS/backend/app/dto"

type LLMAppConfig struct {
	UUID string `json:"uuid"`
	dto.LLMAppConfig
}
