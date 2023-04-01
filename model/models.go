package model

import (
	"time"
)

type BaseModel struct {
	ID        int       `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"-"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"-"`
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
	Username string   `gorm:"UniqueIndex" json:"username"`
	Password string   `gorm:"not null"  json:"-"`
	IsAdmin  bool     `json:"admin"`
	Recipes  []Recipe `gorm:"many2many:user_favorites" json:"-"`
}

type UserFavorite struct {
	UserID   int
	RecipeID int
}
