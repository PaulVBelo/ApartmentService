package models

import "time"

type Descriptions struct {
	ID           string    `gorm:"column:desc_id;primaryKey"`
	ApartmentID  string    `gorm:"column:ap_id;type:uuid;index;not null"`
	ValidFrom    time.Time `gorm:"column:valid_from;default:now();timestampt without time zone"`
	ValidTo      time.Time `gorm:"column:valid_to;timestampt without time zone;default:'9999-12-31 23:59'"`
	Rooms        int       `gorm:"column:rooms"`
	Beds         int       `gorm:"column:beds"`
	Descriptions string    `gorm:"column:desc;type:text"`

	Apartment Apartment `gorm:"foreignKey:ApartmentID;references:ID"`
}
