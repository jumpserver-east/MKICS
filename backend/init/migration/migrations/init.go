package migrations

import (
	"MKICS/backend/app/model"
	"MKICS/backend/global"
	"MKICS/backend/utils/bcrypt"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var AddTableUser = &gormigrate.Migration{
	ID: "20241113-add-table-user",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&model.User{}); err != nil {
			return err
		}
		if err := tx.Create(&model.User{
			Username: global.CONF.SystemConfig.Username,
			Password: bcrypt.Encode(global.CONF.SystemConfig.Password),
		}).Error; err != nil {
			return err
		}
		return nil
	},
}
