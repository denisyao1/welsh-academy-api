package database

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormDB struct {
	db *gorm.DB
}

func NewGormDB() *GormDB {
	dns := os.Getenv("DB_DNS")
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database.")

	}
	return &GormDB{db: db}
}

func (g *GormDB) GetDB() *gorm.DB {
	return g.db
}
