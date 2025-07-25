package response

import "MKICS/backend/app/dto"

type KF struct {
	UUID string `json:"uuid"`
	dto.KF
	Staffs []Staff `json:"staffs"`
}
