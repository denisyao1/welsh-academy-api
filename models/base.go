package models

import "time"

type BaseModel struct {
	ID        int32     `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
}
