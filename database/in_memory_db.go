package database

import (
	"log"

	"github.com/denisyao1/welsh-academy-api/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// InMemorySQLite is database use for testing
type InMemorySQLite struct {
	db     *gorm.DB
	dbFile string
}

func NewInMemoryDB(shared bool) (InMemorySQLite, error) {
	dbFile := "file:testDB?mode=memory"
	if shared {
		dbFile += "&cache=shared"
	}
	db, err := gorm.Open(sqlite.Open(dbFile),
		&gorm.Config{Logger: DefaultDBLogger})
	if err != nil {
		return InMemorySQLite{}, err
	}
	return InMemorySQLite{db: db, dbFile: dbFile}, nil
}

func (m InMemorySQLite) GetDB() *gorm.DB {
	return m.db
}

func (m InMemorySQLite) MigrateAll() {
	m.db.AutoMigrate(&model.Ingredient{}, &model.Recipe{}, &model.User{})
	log.Println("Test Datase migrated successfully")
}

func (m InMemorySQLite) Migrate(models ...interface{}) {
	m.db.AutoMigrate(models...)
}

func (m InMemorySQLite) GetFile() string {
	return m.dbFile
}
