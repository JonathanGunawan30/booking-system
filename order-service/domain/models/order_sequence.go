package models

import "time"

type OrderSequence struct {
	Date       string `gorm:"type:date;primaryKey"`
	LastNumber int    `gorm:"type:int;not null;default:0"`
	UpdatedAt  time.Time
}
