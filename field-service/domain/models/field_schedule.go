package models

import (
	"field-service/constants"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FieldSchedule struct {
	ID        uint                          `gorm:"primaryKey;autoIncrement"`
	UUID      uuid.UUID                     `gorm:"type:uuid;not null"`
	FieldID   uint                          `gorm:"not null"`
	TimeID    uint                          `gorm:"not null"`
	Date      time.Time                     `gorm:"type:date;not null"`
	Status    constants.FieldScheduleStatus `gorm:"type:int;not null"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *gorm.DeletedAt
	Field     Field `gorm:"foreignKey:FieldID;references:ID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	Time      Time  `gorm:"foreignKey:TimeID;references:ID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
}
