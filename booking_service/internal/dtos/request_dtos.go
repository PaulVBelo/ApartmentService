package dtos

import "time"

type ApartmentCreateDTO struct {
	OwnerID      string
	Adress       string
	CheckInTime  time.Time
	CheckOutTime time.Time
	Price        float64
	Info         map[string]string
}

type ApartmentLightUpdateDTO struct {
	OwnerID      string // for validation
	CheckInTime  time.Time
	CheckOutTime time.Time
	Price        float64
}

type ApartmentHeavyUpdateDTO struct {
	OwnerID      string // for validation
	CheckInTime  time.Time
	CheckOutTime time.Time
	Price        float64
	Info         map[string]string
}

type BookingCreateDTO struct {
	UserID      string
	ApartmentID string
	TimeFrom    time.Time
	TimeTo      time.Time
}
