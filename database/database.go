package database

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormDB struct {
	db *gorm.DB
}

func NewGormDB() *GormDB {
	dns := os.Getenv("DB_DNS")
	// @TODO : delete DB logger config
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		log.Fatal("Failed to connect to database.")

	}
	return &GormDB{db: db}
}

func (g *GormDB) GetDB() *gorm.DB {
	return g.db
}
