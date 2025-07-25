package dto

type Policy struct {
	PolicyName string     `json:"policyname" validate:"required"`
	Repeat     int        `json:"repeat" validate:"required,oneof=1 2 3 4 5"`
	Week       string     `json:"week" validate:"omitempty,len=7,weekFormat"`
	MaxCount   int        `json:"max_count" validate:"required,min=1"`
	WorkTimes  []WorkTime `json:"work_times" validate:"required,min=1,dive"`
}

type WorkTime struct {
	StartTime string `json:"start_time" validate:"required,time"`
	EndTime   string `json:"end_time" validate:"required,time"`
}
