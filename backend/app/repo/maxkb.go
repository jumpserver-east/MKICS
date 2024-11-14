package repo

import (
	"EvoBot/backend/app/model"
	"EvoBot/backend/global"
)

type IMaxkbRepo interface {
	GetConfig() (model.MaxkbConf, error)
	UpdateConfig(conf model.MaxkbConf) error
}

type MaxkbRepo struct{}

func NewIMaxkbRepo() IMaxkbRepo {
	return &MaxkbRepo{}
}

func (r *MaxkbRepo) GetConfig() (model.MaxkbConf, error) {
	var conf model.MaxkbConf
	err := global.DB.First(&conf).Error
	return conf, err
}

func (r *MaxkbRepo) UpdateConfig(conf model.MaxkbConf) error {
	return global.DB.Model(&model.MaxkbConf{}).Where("id = ?", conf.ID).Updates(conf).Error
}
