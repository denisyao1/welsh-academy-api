package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/denisyao1/welsh-academy-api/common"
	"github.com/denisyao1/welsh-academy-api/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormDB struct {
	db *gorm.DB
}

func NewGormDB(config common.Configuration) *GormDB {
	dns := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%v", config.DB_HOST, config.DB_USER,
		config.DB_PASSWORD, config.DB_NAME, config.DB_PORT)

	dbLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{Logger: dbLogger})
	if err != nil {
		log.Fatal("Failed to connect to database.")

	}
	log.Println("Connected to database")
	return &GormDB{db: db}
}

func (g *GormDB) GetDB() *gorm.DB {
	return g.db
}

func (g *GormDB) Migrate() {
	g.db.AutoMigrate(&model.Ingredient{}, &model.Recipe{}, &model.User{})
	log.Println("Datase migrated successfully")
}
