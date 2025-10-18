package models

import (
	"time"
)

type Apartment struct {
	ID           string    `gorm:"column:id;type:uuid;primaryKey"`
	OwnerID      string    `gorm:"column:owner_id;type:uuid;not null"`
	Adress       string    `gorm:"column:adress;not null"`
	CheckInTime  time.Time `gorm:"column:check_in_time;type:time without time zone"`
	CheckOutTime time.Time `gorm:"column:check_out_time;type:time without time zone"`
	Price        float64   `gorm:"column:price;type:decimal(10,2)"`
	UpdatedAt    time.Time `gorm:"column:updated_at;type:timestampt without time zone;default:now();not null;"`

	// Relations
	Bookings     []Booking       `gorm:"foreignKey:apartmentID;constraint:OnDelete:CASCADE"`
	Descriptions []Descriptions  `gorm:"foreignKey:apartmentID;constraint:OnDelete:CASCADE"`
	History      []ApartmentSCD4 `gorm:"foreignKey:apartmentID;constraint:OnDelete:CASCADE"`
}

func (Apartment) TableName() string {
	return "apartments"
}
