package models

import "time"

type Booking struct {
	ID          string    `gorm:"column:booking_id;primaryKey"`
	UserID      string    `gorm:"column:user_id;type:uuid;not null"`
	ApartmentID string    `gorm:"column:ap_id;type:uuid;index;not null"`
	TimeFrom    time.Time `gorm:"column:time_from;type:timestamp without time zone;default:now()"`
	TimeTo      time.Time `gorm:"column:time_to;type:timestamp without time zone;default:('9999-12-31 23:59:00'::timestamp)"`

	Apartment Apartment `gorm:"foreignKey:ApartmentID;references:ID"`
}

func (Booking) TableName() string {
	return "bookings"
}
