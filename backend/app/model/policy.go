package model

type Policy struct {
	BaseModel
	PolicyName string     `gorm:"column:policyname;type:varchar(255);unique;" json:"policyname"`
	Repeat     int        `gorm:"column:repeat;type:int;" json:"repeat"`
	Week       string     `gorm:"column:week;type:varchar(7);" json:"week"`
	MaxCount   int        `gorm:"column:max_count;type:int;" json:"max_count"`
	WorkTimes  []WorkTime `gorm:"foreignKey:PolicyID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type WorkTime struct {
	BaseModel
	PolicyID  uint   `gorm:"column:policy_id"`
	StartTime string `gorm:"column:start_time;type:time(0);" json:"start_time"`
	EndTime   string `gorm:"column:end_time;type:time(0);" json:"end_time"`
}
