package db

import (
	"EvoBot/backend/global"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func Init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		global.CONF.DBConfig.Username,
		global.CONF.DBConfig.Password,
		global.CONF.DBConfig.Host,
		global.CONF.DBConfig.Port,
		global.CONF.DBConfig.Database,
	)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "eb_",
			SingularTable: true,
		},
	})
	if err != nil {
		global.ZAPLOG.Error("failed to connect database ,err:" + err.Error())
		panic(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		global.ZAPLOG.Error("connect db server failed, err:" + err.Error())
		panic(err)
	}
	sqlDB.SetConnMaxIdleTime(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	global.DB = db
	global.ZAPLOG.Info("init db successfully")
}
