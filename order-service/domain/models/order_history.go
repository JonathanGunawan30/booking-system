package models

import (
	"order-service/constants"
	"time"
)

type OrderHistory struct {
	ID        uint                        `gorm:"primaryKey;autoIncrement"`
	OrderID   uint                        `gorm:"bigint;not null"`
	Status    constants.OrderStatusString `gorm:"varchar(30);not null"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
