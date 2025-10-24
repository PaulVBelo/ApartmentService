package dtos

import "time"

type BookingRange struct {
	From time.Time `json:"from" gorm:"column:from"`
	To   time.Time `json:"to" gorm:"column:to"`
}

type ShortApartmentResponse struct {
	Id      string  `json:"id" binding:"required,uuid4"`
	OwnerID string  `json:"owner_id" binding:"required,uuid4"`
	Address string  `json:"address" binding:"required"`
	Price   float64 `json:"price" binding:"required,gt=0"`
}

type MediumApartmentResponse struct {
	Id      string            `json:"id" binding:"required,uuid4"`
	OwnerID string            `json:"owner_id" binding:"required,uuid4"`
	Address string            `json:"address" binding:"required"`
	Price   float64           `json:"price" binding:"required,gt=0"`
	Info    map[string]string `json:"info" binding:"omitempty,dive,keys,required,endkeys,required"`
}

type FullApartmentResponse struct {
	Id       string                 `json:"id" binding:"required,uuid4"`
	OwnerID  string                 `json:"owner_id" binding:"required,uuid4"`
	Address  string                 `json:"address" binding:"required"`
	Price    float64                `json:"price" binding:"required,gt=0"`
	Info     map[string]string      `json:"info" binding:"omitempty,dive,keys,required,endkeys,required"`
	Bookings []ShortBookingResponse `json:"bookings"`
}

type ShortBookingResponse struct {
	Id       string    `json:"id" binding:"required,uuid4"`
	UserID   string    `json:"user_id" binding:"required,uuid4"`
	TimeFrom time.Time `json:"time_from" binding:"required"`
	TimeTo   time.Time `json:"time_to" binding:"required,gtfield=TimeFrom"`
}

type BookingResponse struct {
	Id          string    `json:"id" binding:"required,uuid4"`
	UserID      string    `json:"user_id" binding:"required,uuid4"`
	ApartmentID string    `json:"ap_id" binding:"required,uuid4"`
	Address     string    `json:"address" binding:"required,uuid4"`
	TimeFrom    time.Time `json:"time_from" binding:"required"`
	TimeTo      time.Time `json:"time_to" binding:"required,gtfield=TimeFrom"`
}
