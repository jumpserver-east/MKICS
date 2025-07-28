package dto

type Staff struct {
	StaffID   string `json:"staffid" validate:"required"`
	StaffName string `json:"staffname" validate:"required"`
}
