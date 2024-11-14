package request

import "EvoBot/backend/app/dto"

type WecomConfigApp struct {
	dto.WecomConfig
	AgentID string `json:"agent_id"`
}

type ReceptionistOptions struct {
	OpenKFID   string   `json:"open_kfid"`
	UserIDList []string `json:"userid_list"`
}
