package db

import (
	"booking_service/internal/config"
	"booking_service/internal/models"

	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectPostgres(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB_HOST, cfg.DB_PORT, cfg.DB_USER, cfg.DB_PASSWORD, cfg.DB_NAME,
	)

	var db *gorm.DB
	var err error

	maxRetries := 10
	retryInterval := 5 * time.Second

	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		logrus.Warnf("error connecting to database: %v", err)
		time.Sleep(retryInterval)
	}

	if err != nil {
		return nil, fmt.Errorf("error connecting to database after %d retries: %v", maxRetries, err)
	}

	logrus.Info("Connected to database succsessfully!")

	err = db.AutoMigrate(
		&models.Apartment{},
		&models.ApartmentSCD4{},
		&models.Description{},
		&models.Booking{},
	)

	if err != nil {
		return nil, fmt.Errorf("error during migration: %v", err)
	}

	logrus.Info("All migrations applied succsessfully!")

	return db, nil
}
