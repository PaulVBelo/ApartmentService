package repository

import (
	"booking_service/internal/dtos"
	"booking_service/internal/models"
)

type Repository interface {
	AddApartment(dto *dtos.ApartmentCreateDTO) (models.Apartment, error)

	UpdateApartmentLight(id string, dto *dtos.ApartmentLightUpdateDTO) error
	UpdateApartmentHeavy(id string, dto *dtos.ApartmentHeavyUpdateDTO) error

	GetApartments(filter map[string]string) (*[]models.Apartment, error) // empty, if no filter
	GetApartment(id string) (models.Apartment, error)

	CreateBooking(dto *dtos.BookingCreateDTO) (models.Booking, error)
}
