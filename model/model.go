package model

import (
	"time"
)

type BaseModel struct {
	ID        int       `gorm:"primarykey" json:"id" example:"1" extensions:"x-order=1"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"-"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"-"`
}

type Ingredient struct {
	BaseModel
	Name string `gorm:"uniqueIndex" json:"name" example:"Tomato"`
}

type Recipe struct {
	BaseModel
	Name        string       `gorm:"uniqueIndex" json:"name" extensions:"x-order=2"`
	Making      string       `gorm:"type:text;not null" json:"making" extensions:"x-order=3"`
	Ingredients []Ingredient `gorm:"many2many:recipe_ingredients" json:"ingredients"`
}

type User struct {
	ID        int       `gorm:"primarykey" json:"-"`
	Username  string    `gorm:"UniqueIndex;not null" json:"username" extensions:"x-order=1"`
	Password  string    `gorm:"not null"  json:"-"`
	IsAdmin   bool      `json:"admin"`
	Recipes   []Recipe  `gorm:"many2many:user_favorites" json:"-"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"-"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"-"`
}

type UserFavorite struct {
	UserID   int
	RecipeID int
}
