package repository

import (
	"booking_service/internal/dtos"
	"booking_service/internal/models"
	servererrors "booking_service/internal/server_errors"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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
	var existing models.Apartment
	err := r.tm.db.Where("address = ?", dto.Address).First(&existing).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return models.Apartment{}, err
	}

	if err == nil {
		return models.Apartment{}, &servererrors.AlreadyExistsError{Field: "address", Value: dto.Address}
	}

	id := uuid.New().String()
	ap := models.Apartment{
		ID:        id,
		OwnerID:   dto.OwnerID,
		Address:   dto.Address,
		Price:     dto.Price,
		UpdatedAt: time.Now(),
	}

	// TRANSACTION [BEGIN]
	tx, err := r.tm.begin()
	if err != nil {
		return models.Apartment{}, err
	}

	if err := tx.Create(&ap).Error; err != nil {
		_ = r.tm.rollback(tx)
		return models.Apartment{}, err
	}

	if len(dto.Info) > 0 {
		ve := &servererrors.ValidationError{}

		var rooms, beds *int
		other := map[string]string{}

		for k, v := range dto.Info {
			switch k {
			case "rooms":
				val, err := strconv.Atoi(v)
				if err != nil {
					ve.Add("info.rooms", "must be an integer")
				} else if val <= 0 {
					ve.Add("info.rooms", "must be > 0")
				} else {
					rooms = &val
				}
			case "beds":
				val, err := strconv.Atoi(v)
				if err != nil {
					ve.Add("info.beds", "must be an integer")
				} else if val < 0 {
					ve.Add("info.beds", "must be >= 0")
				} else {
					beds = &val
				}
			default:
				other[k] = v
			}
		}

		if len(ve.Fields) > 0 {
			return models.Apartment{}, ve
		}

		descJSON, _ := json.Marshal(other)

		desc := models.Description{
			ID:          uuid.New().String(),
			ApartmentID: id,
			Description: string(descJSON),
		}
		if rooms != nil {
			desc.Rooms = *rooms
		}
		if beds != nil {
			desc.Beds = *beds
		}

		if err := tx.Create(&desc).Error; err != nil {
			_ = r.tm.rollback(tx)
			return models.Apartment{}, err
		}
	}

	// COMMIT (TRANSACTION END)
	if err := r.tm.commit(tx); err != nil {
		return models.Apartment{}, err
	}

	logrus.WithTime(time.Now()).Infof("Successfuly created apartment, id = %s", id)
	return ap, nil
}

func (r *repositoryWithTM) UpdateApartmentLight(id string, dto *dtos.ApartmentLightUpdateDTO) error {
	// TRANSACTION [BEGIN]
	tx, err := r.tm.begin()
	if err != nil {
		return err
	}

	operationTimestamp := time.Now().Add(2 * time.Second)

	var ap models.Apartment
	if err = tx.Where("id = ?", id).First(&ap).Error; err != nil {
		_ = r.tm.rollback(tx)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("apartment with id '%s' does not exist", id)
		}
		return err
	}

	if ap.OwnerID != dto.OwnerID {
		_ = r.tm.rollback(tx)
		return &servererrors.ForbiddenAccessError{
			UserId:       dto.OwnerID,
			ResourceType: "apartment",
			ResourceId:   ap.ID,
		}
	}

	scd4 := models.ApartmentSCD4{
		ID:           uuid.New().String(),
		ApartmentID:  ap.ID,
		Price:        ap.Price,
		UpdatedAt:    ap.UpdatedAt,
		InvalidSince: operationTimestamp,
	}

	if err := tx.Save(&scd4).Error; err != nil {
		_ = r.tm.rollback(tx)
		return err
	}

	if dto.Price != nil {
		ap.Price = *dto.Price
	}

	ap.UpdatedAt = operationTimestamp

	if err := tx.Save(&ap).Error; err != nil {
		_ = r.tm.rollback(tx)
		return err
	}

	// COMMIT (TRANSACTION END)
	if err := r.tm.commit(tx); err != nil {
		return err
	}

	logrus.WithTime(time.Now()).Infof("Successfuly updated apartment, id = %s", id)
	return nil
}

func (r *repositoryWithTM) UpdateApartmentHeavy(id string, dto *dtos.ApartmentHeavyUpdateDTO) error {
	// TRANSACTION [BEGIN]
	tx, err := r.tm.begin()
	if err != nil {
		return err
	}

	operationTimestamp := time.Now().Add(2 * time.Second)

	var ap models.Apartment
	if err = tx.Preload("Descriptions", "valid_to = '9999-12-31 23:59:00'").Where("id = ?", id).First(&ap).Error; err != nil {
		_ = r.tm.rollback(tx)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &servererrors.NotFoundError{Entity: "apartment", Key: id}
		}
		return err
	}

	if ap.OwnerID != dto.OwnerID {
		_ = r.tm.rollback(tx)
		return &servererrors.ForbiddenAccessError{
			UserId:       dto.OwnerID,
			ResourceType: "apartment",
			ResourceId:   ap.ID,
		}
	}

	scd4 := models.ApartmentSCD4{
		ID:           uuid.New().String(),
		ApartmentID:  ap.ID,
		Price:        ap.Price,
		UpdatedAt:    ap.UpdatedAt,
		InvalidSince: operationTimestamp,
	}

	if err := tx.Save(&scd4).Error; err != nil {
		_ = r.tm.rollback(tx)
		return err
	}

	if dto.Price != nil {
		ap.Price = *dto.Price
	}

	ap.UpdatedAt = operationTimestamp

	if err := tx.Save(&ap).Error; err != nil {
		_ = r.tm.rollback(tx)
		return err
	}

	if len(dto.Info) > 0 {
		ve := &servererrors.ValidationError{}

		var rooms, beds *int
		other := map[string]string{}

		for k, v := range dto.Info {
			switch k {
			case "rooms":
				val, err := strconv.Atoi(v)
				if err != nil {
					ve.Add("info.rooms", "must be an integer")
				} else if val <= 0 {
					ve.Add("info.rooms", "must be > 0")
				} else {
					rooms = &val
				}
			case "beds":
				val, err := strconv.Atoi(v)
				if err != nil {
					ve.Add("info.beds", "must be an integer")
				} else if val < 0 {
					ve.Add("info.beds", "must be >= 0")
				} else {
					beds = &val
				}
			default:
				other[k] = v
			}
		}

		if len(ve.Fields) > 0 {
			return ve
		}

		descJSON, _ := json.Marshal(other)

		newDesc := models.Description{
			ID:          uuid.New().String(),
			ApartmentID: id,
			Description: string(descJSON),
			ValidFrom:   operationTimestamp,
		}
		if rooms != nil {
			newDesc.Rooms = *rooms
		}
		if beds != nil {
			newDesc.Beds = *beds
		}

		oldDesc := ap.Descriptions[0]
		oldDesc.ValidTo = operationTimestamp

		if err := tx.Save(oldDesc).Error; err != nil {
			_ = r.tm.rollback(tx)
			return err
		}

		if err := tx.Create(&newDesc).Error; err != nil {
			_ = r.tm.rollback(tx)
			return err
		}
	}

	// COMMIT (TRANSACTION END)
	if err := r.tm.commit(tx); err != nil {
		return err
	}

	logrus.WithTime(time.Now()).Infof("Successfuly updated apartment and its description, id = %s", id)
	return nil
}

func (r *repositoryWithTM) GetApartments(filter map[string]string) (*[]models.Apartment, error) {
	var apartments []models.Apartment
	db := r.tm.db.Model(&models.Apartment{})

	if city, ok := filter["city"]; ok && city != "" {
		db = db.Where("split_part(address, ',', 1) ILIKE ?", "%"+city+"%")
	}

	roomsFilter, hasRooms := filter["rooms"]
	bedsFilter, hasBeds := filter["beds"]

	if hasRooms || hasBeds {
		db = db.Joins("JOIN description d ON d.ap_id = apartments.id AND d.valid_to = '9999-12-31 23:59:00'")
		if hasRooms {
			db = db.Where("d.rooms = ?", roomsFilter)
		}
		if hasBeds {
			db = db.Where("d.beds = ?", bedsFilter)
		}
	}

	if err := db.Find(&apartments).Error; err != nil {
		return nil, err
	}

	return &apartments, nil
}

func (r *repositoryWithTM) GetApartment(id string) (models.Apartment, []dtos.BookingRange, error) {
	var ap models.Apartment
	var bookings []dtos.BookingRange

	err := r.tm.db.
		Preload("Descriptions", "valid_to = '9999-12-31 23:59:00'").
		Where("id = ?", id).
		First(&ap).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Apartment{}, nil, &servererrors.NotFoundError{Key: id, Entity: "apartment"}
		}
		return models.Apartment{}, nil, err
	}

	// подгружаем активные брони (у которых окончание не в прошлом)
	if err := r.tm.db.
		Model(&models.Booking{}).
		Select("time_from AS from, time_to AS to").
		Where("ap_id = ?", id).
		Where("time_to > now()").
		Scan(&bookings).Error; err != nil {
		return ap, nil, err
	}

	return ap, bookings, nil
}

func (r *repositoryWithTM) CreateBooking(dto *dtos.BookingCreateDTO) (models.Booking, error) {
	tx, err := r.tm.begin()
	if err != nil {
		return models.Booking{}, err
	}

	var ap models.Apartment
	if err := tx.Where("id = ?", dto.ApartmentID).First(&ap).Error; err != nil {
		_ = r.tm.rollback(tx)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Booking{}, fmt.Errorf("apartment with id '%s' does not exist", dto.ApartmentID)
		}
		return models.Booking{}, err
	}

	var conflictCount int64
	if err := tx.Model(&models.Booking{}).
		Where("ap_id = ?", dto.ApartmentID).
		Where("time_from < ? AND time_to > ?", dto.TimeTo, dto.TimeFrom).
		Count(&conflictCount).Error; err != nil {
		_ = r.tm.rollback(tx)
		return models.Booking{}, err
	}

	if conflictCount > 0 {
		_ = r.tm.rollback(tx)
		return models.Booking{}, &servererrors.OverlapError{ApId: ap.ID}
	}

	booking := models.Booking{
		ID:          uuid.New().String(),
		UserID:      dto.UserID,
		ApartmentID: dto.ApartmentID,
		TimeFrom:    dto.TimeFrom,
		TimeTo:      dto.TimeTo,
	}

	if err := tx.Create(&booking).Error; err != nil {
		_ = r.tm.rollback(tx)
		return models.Booking{}, err
	}

	if err := r.tm.commit(tx); err != nil {
		return models.Booking{}, err
	}

	booking.Apartment = ap

	return booking, nil
}

func (r *repositoryWithTM) GetApartmentsByOwner(id string) (*[]models.Apartment, error) {
	var apartments []models.Apartment
	operationTimestamp := time.Now()

	err := r.tm.db.
		Preload("Descriptions", "valid_to = '9999-12-31 23:59:00'").
		Preload("Bookings", "time_to >= ?", operationTimestamp).
		Where("owner_id = ?", id).
		Find(&apartments).Error

	if err != nil {
		return nil, err
	}

	return &apartments, nil
}

func (r *repositoryWithTM) GetBookingsByUser(id string) (*[]models.Booking, error) {
	var bookings []models.Booking

	err := r.tm.db.
		Preload("Apartment").
		Where("user_id = ?", id).
		Find(&bookings).Error

	if err != nil {
		return nil, err
	}

	return &bookings, nil
}
