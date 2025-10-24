package models

import "time"

type Description struct {
	ID          string    `gorm:"column:desc_id;primaryKey"`
	ApartmentID string    `gorm:"column:ap_id;type:uuid;index;not null"`
	ValidFrom   time.Time `gorm:"column:valid_from;type:timestamp without time zone;default:now()"`
	ValidTo     time.Time `gorm:"column:valid_to;type:timestamp without time zone;default:'9999-12-31 23:59:00'"`
	Rooms       int       `gorm:"column:rooms;default:-1"`
	Beds        int       `gorm:"column:beds;default:-1"`
	Description string    `gorm:"column:desc;type:text"`

	Apartment Apartment `gorm:"foreignKey:ApartmentID;references:ID"`
}

func (Description) TableName() string {
	return "description"
}
