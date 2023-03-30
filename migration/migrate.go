package main

import (
	"github.com/denisyao1/welsh-academy-api/database"
	"github.com/denisyao1/welsh-academy-api/initializer"
	"github.com/denisyao1/welsh-academy-api/model"
)

func init() {
	initializer.LoadEnvVariables()
}

func main() {
	gormDB := database.NewGormDB().GetDB()
	gormDB.AutoMigrate(&model.Ingredient{})
	gormDB.AutoMigrate(&model.Recipe{})
	gormDB.AutoMigrate(&model.User{})
}
