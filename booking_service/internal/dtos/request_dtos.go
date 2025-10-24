package dtos

import "time"

// ApartmentCreateDTO используется при создании нового апартамента
type ApartmentCreateDTO struct {
	OwnerID      string            `json:"owner_id" binding:"required,uuid4"`
	Address      string            `json:"address" binding:"required"`
	Price        float64           `json:"price" binding:"required,gt=0"`
	Info         map[string]string `json:"info" binding:"omitempty,dive,keys,required,endkeys,required"`
}

// ApartmentUpdateDTO — dto обновления для unmarshall
type ApartmentUpdateDTO struct {
	OwnerID      string             `json:"owner_id" binding:"required,uuid4"`
	Price        *float64           `json:"price" binding:"required,gt=0"`
	Info         *map[string]string `json:"info" binding:"omitempty,dive,keys,required,endkeys,required"`
}

// ApartmentLightUpdateDTO — лёгкое обновление без изменения описаний
type ApartmentLightUpdateDTO struct {
	OwnerID      string     `json:"owner_id" binding:"required,uuid4"`
	Price        *float64   `json:"price" binding:"required,gt=0"`
}

// ApartmentHeavyUpdateDTO — обновление с изменением описаний (Info)
type ApartmentHeavyUpdateDTO struct {
	OwnerID      string            `json:"owner_id" binding:"required,uuid4"`
	Price        *float64          `json:"price" binding:"required,gt=0"`
	Info         map[string]string `json:"info" binding:"required,dive,keys,required,endkeys,required"`
}

// BookingCreateDTO — создание бронирования
type BookingCreateDTO struct {
	UserID      string    `json:"user_id" binding:"required,uuid4"`
	ApartmentID string    `json:"apartment_id" binding:"required,uuid4"`
	TimeFrom    time.Time `json:"time_from" binding:"required"`
	TimeTo      time.Time `json:"time_to" binding:"required,gtfield=TimeFrom"`
}
