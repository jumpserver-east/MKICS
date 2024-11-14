package model

type KH struct {
	BaseModel
	KHID    string `gorm:"uniqueIndex;column:khid;type:varchar(255);not null;" json:"khid"`
	ChatID  string `gorm:"column:chatid;type:varchar(255);" json:"chatid"`
	StaffID string `gorm:"column:staffid;" json:"staffid"`
}
