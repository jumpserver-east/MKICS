package model

type KH struct {
	BaseModel
	KHID       string     `gorm:"uniqueIndex;column:khid;type:varchar(255);not null;" json:"khid"`
	StaffID    string     `gorm:"column:staffid;type:varchar(255);" json:"staffid"`
	SceneParam string     `gorm:"column:scence_param;type:varchar(32);" json:"scence_param"`
	ChatList   []ChatList `gorm:"foreignKey:KHID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type ChatList struct {
	BaseModel
	KHID   uint   `gorm:"column:kh_id"`
	BotID  string `gorm:"column:botid;type:varchar(255);not null;" json:"botid"`
	ChatID string `gorm:"column:chatid;type:varchar(255);not null;" json:"chatid"`
}
