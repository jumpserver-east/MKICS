package response

import "MKICS/backend/app/dto"

type WecomConfigApp struct {
	UUID string `json:"uuid"`
	dto.WecomConfig
	AgentID string `json:"agent_id"`
}
