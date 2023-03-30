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
	Making      string       `gorm:"type:text;not null" json:"making"`
	Ingredients []Ingredient `gorm:"many2many:recipe_ingredients" json:"ingredients"`
}

type User struct {
	BaseModel
	Username string `gorm:"UniqueIndex" json:"username"`
	Password string `gorm:"not null"  json:"-"`
	IsAdmin  bool   `json:"isAdmin"`
}
