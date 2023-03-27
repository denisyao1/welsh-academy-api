package main

import (
	"github.com/denisyao1/welsh-academy-api/database"
	"github.com/denisyao1/welsh-academy-api/initializers"
	"github.com/denisyao1/welsh-academy-api/models"
)

func init() {
	initializers.LoadEnvVariables()
}

func main() {
	gormDB := database.NewGormDB().GetDB()
	gormDB.AutoMigrate(&models.Ingredient{})
}
