package model

type MaxkbConf struct {
	BaseModel
	BaseURL string `gorm:"column:base_url;type:varchar(255);"`
	ApiKey  string `gorm:"column:api_key;type:varchar(255);"`
}
