package models

type Ingredient struct {
	BaseModel
	Name string `gorm:"uniqueIndex" validate:"required" json:"name"`
}
