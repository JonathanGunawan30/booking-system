package config

import (
	"fmt"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDatabase() (*gorm.DB, error) {
	config := AppConfig
	encodedPassword := url.QueryEscape(config.Database.Password)
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		config.Database.Username,
		encodedPassword,
		config.Database.Host,
		config.Database.Port,
		config.Database.Name,
		config.Database.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Errorf("failed to connect to database: %v", err)
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		logrus.Errorf("database ping failed: %v", err)
		return nil, err
	}

	sqlDB.SetMaxIdleConns(config.Database.MaxIdleConnections)
	sqlDB.SetMaxOpenConns(config.Database.MaxOpenConnections)
	sqlDB.SetConnMaxLifetime(time.Duration(config.Database.MaxLifetimeConnection) * time.Second)
	sqlDB.SetConnMaxIdleTime(time.Duration(config.Database.MaxIdleTime) * time.Second)

	return db, nil
}
