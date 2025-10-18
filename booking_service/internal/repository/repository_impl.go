package repository

import (
	"booking_service/internal/dtos"
	"booking_service/internal/models"

	"gorm.io/gorm"
)

type repositoryWithTM struct {
	tm *transactionManager
}

func NewRepository(db *gorm.DB) Repository {
	return &repositoryWithTM{
		tm: &transactionManager{db},
	}
}

func (r *repositoryWithTM) AddApartment(dto *dtos.ApartmentCreateDTO) (models.Apartment, error) {
	// [TODO]
}

func (r *repositoryWithTM) UpdateApartmentLight(id string, dto *dtos.ApartmentLightUpdateDTO) error {
	// [TODO]
}

func (r *repositoryWithTM) UpdateApartmentHeavy(id string, dto *dtos.ApartmentHeavyUpdateDTO) error {
	// [TODO]
}

func (r *repositoryWithTM) GetApartments(filter map[string]string) (*[]models.Apartment, error) {
	// [TODO]
}

func (r *repositoryWithTM) GetApartment(id string) (models.Apartment, error) {
	// [TODO]
}

func (r *repositoryWithTM) CreateBooking(dto *dtos.BookingCreateDTO) (models.Booking, error) {
	// [TODO]
}
