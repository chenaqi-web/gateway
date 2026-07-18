package repo

import (
	"fmt"
	"log"

	"backend/gateway/internal/config"
	"backend/gateway/internal/model/entity"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DBClient struct {
	DB *gorm.DB
}

func NewDBClient(cfg *config.Config) (*DBClient, error) {
	mysqlCfg := cfg.Mysql
	if mysqlCfg.Host == "" || mysqlCfg.Port == "" || mysqlCfg.DBName == "" {
		return nil, fmt.Errorf("mysql config is incomplete")
	}

	db, err := gorm.Open(mysql.Open(mysqlCfg.DSN()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open mysql: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql db: %w", err)
	}

	sqlDB.SetMaxIdleConns(mysqlCfg.MaxIdleConn)
	sqlDB.SetMaxOpenConns(mysqlCfg.MaxOpenConn)

	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}

	log.Println("mysql connected successfully")
	return &DBClient{DB: db}, nil
}

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&entity.AiChatSession{},
		&entity.AiChatMessage{},
	)
}

func (c *DBClient) Close() error {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (c *DBClient) GetDB() *gorm.DB {
	return c.DB
}
