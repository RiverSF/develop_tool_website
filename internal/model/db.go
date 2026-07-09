package model

import (
	"fmt"
	"time"

	"develop_tools/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var db *gorm.DB

func Init() error {
	connArgs := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		config.MysqlConfig.MysqlUser,
		config.MysqlConfig.MysqlPassword,
		config.MysqlConfig.MysqlHost,
		config.MysqlConfig.MysqlPort,
		config.MysqlConfig.MysqlDb,
	)

	logLevel := logger.Info
	if config.AppConfig.RunMode == "release" {
		logLevel = logger.Warn
	}

	var err error
	db, err = gorm.Open(mysql.Open(connArgs), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
		Logger:         logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	return sqlDB.Ping()
}

func Close() error {
	if db == nil {
		return nil
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
