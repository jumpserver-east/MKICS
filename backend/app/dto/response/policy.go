package response

import "MKICS/backend/app/dto"

type Policy struct {
	UUID string `json:"uuid"`
	dto.Policy
}
