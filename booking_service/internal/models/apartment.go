package models

import (
	"time"
)

type Apartment struct {
	ID           string    `gorm:"column:id;type:uuid;primaryKey"`
	OwnerID      string    `gorm:"column:owner_id;type:uuid;not null"`
	Address      string    `gorm:"column:address;uniqueIndex;not null"`
	Price        float64   `gorm:"column:price;type:decimal(10,2)"`
	UpdatedAt    time.Time `gorm:"column:updated_at;type:timestamp without time zone;default:now();not null;"`

	// Relations
	Bookings     []Booking       `gorm:"foreignKey:ApartmentID;constraint:OnDelete:CASCADE"`
	Descriptions []Description   `gorm:"foreignKey:ApartmentID;constraint:OnDelete:CASCADE"`
	History      []ApartmentSCD4 `gorm:"foreignKey:ApartmentID;constraint:OnDelete:CASCADE"`
}

func (Apartment) TableName() string {
	return "apartments"
}
