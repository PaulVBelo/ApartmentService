package models

import (
	"time"
)

type ApartmentSCD4 struct {
	ID           string    `gorm:"column:id;type:uuid;primaryKey"`
	ApartmentID  string    `gorm:"column:ap_id;index;type:uuid;not null"`
	CheckInTime  time.Time `gorm:"column:check_in_time;type:time without time zone"`
	CheckOutTime time.Time `gorm:"column:check_out_time;type:time without time zone"`
	Price        float64   `gorm:"column:price;type:decimal(10,2)"`
	UpdatedAt    time.Time `gorm:"type:timestamptz;default:now();not null;"`
	InvalidSince time.Time `gorm:"type:timestamptz;default:now();not null;"`

	Apartment Apartment `gorm:"foreignKey:ApartmentID;references:ID"`
}

func (ApartmentSCD4) TableName() string {
	return "apartments_scd4"
}
