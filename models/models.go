package models

import "time"

type BaseModel struct {
	ID        int32     `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"-"`
}

type Ingredient struct {
	BaseModel
	Name string `gorm:"uniqueIndex" json:"name"`
}

type Recipe struct {
	BaseModel
	Name        string       `gorm:"uniqueIndex" json:"name"`
	Ingredients []Ingredient `gorm:"many2many:recipe_ingredients" json:"ingredients"`
}
