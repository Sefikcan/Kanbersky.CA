package entity

import (
	"time"
)

type Currency struct {
	ID int `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Title string `gorm:"index:idx_title,unique" json:"title"`
	IsoCode string `gorm:"index:idx_iso_code,unique" json:"iso_code"`
}
