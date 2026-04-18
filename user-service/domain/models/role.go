package models

import "time"

type Role struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Code      string `gorm:"type:varchar(15);not null"`
	Name      string `gorm:"type:varchar(20);not null"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
