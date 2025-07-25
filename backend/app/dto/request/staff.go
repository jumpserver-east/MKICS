package request

import "MKICS/backend/app/dto"

type Staff struct {
	dto.Staff
	PolicyList []string `json:"policy_list" validate:"required"`
}
