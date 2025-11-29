package server

import (
	"booking_service/internal/dtos"
	"booking_service/internal/repository"
	servererrors "booking_service/internal/server_errors"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type InnerServer struct {
	router     *gin.Engine
	repository repository.Repository
}

func NewServer(db *gorm.DB) *InnerServer {
	router := gin.Default()
	s := &InnerServer{
		router:     router,
		repository: repository.NewRepository(db),
	}
	s.routes()
	return s
}

func (s *InnerServer) routes() {
	s.router.POST("/apartments", s.postApartment)
	s.router.PATCH("/apartments/:id", s.updateApartment)
	s.router.GET("/apartments", s.getApartmentsFiltered)
	s.router.GET("/apartments/:id", s.getApartmentById)
	s.router.POST("/book", s.bookApartment)
	s.router.GET("/owners/:id/apartments", s.getApartmentsByOwner)
	s.router.GET("/users/:id/bookings", s.getBookingsByUser)
	s.router.GET("/health", health)
}

func (s *InnerServer) Run(host, port string) *http.Server {
	addr := host + ":" + port
	http_server := http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	go func() {
		logrus.WithField("Time", time.Now().String()).Infof("Server started on addr %s", addr)

		if err := http_server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.WithField("Time", time.Now().String()).Fatalf("Server failed: %v", err)
		}
	}()

	return &http_server
}

func GracefulShutdown(srv *http.Server, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logrus.Errorf("Server forced to shutdown: %v", err)
	} else {
		logrus.Info("Server exited gracefully")
	}
}

func health(c *gin.Context) {
	c.Status(http.StatusOK)
}

// === POST APARTMENT ===

func (s *InnerServer) postApartment(c *gin.Context) {
	var dto dtos.ApartmentCreateDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		logrus.WithField("Time", time.Now().String()).Infof("400: Bad Request: %v", err)
		return
	}

	ap, err := s.repository.AddApartment(&dto)
	if err != nil {
		var aee *servererrors.AlreadyExistsError
		if errors.As(err, &aee) {
			c.JSON(http.StatusConflict, gin.H{"error": "apartment with this address is already registered"})
			logrus.WithField("Time", time.Now().String()).Info("409: Conflict")
			return
		}

		var ve *servererrors.ValidationError
		if errors.As(err, &ve) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request form, beds & rooms must be positive integers"})
			logrus.WithField("Time", time.Now().String()).Info("400: Bad Request")
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		logrus.WithField("Time", time.Now().String()).Warn("500: Internal Server Error")
		return
	}

	response := dtos.MediumApartmentResponse{
		Id:      ap.ID,
		OwnerID: ap.OwnerID,
		Address: ap.Address,
		Price:   ap.Price,
		Info:    dto.Info,
	}

	c.JSON(http.StatusCreated, response)
	logrus.WithField("Time", time.Now().String()).Info("201: Created Apartment")
}

func (s *InnerServer) updateApartment(c *gin.Context) {
	id := c.Param("id")

	var dto dtos.ApartmentUpdateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	var err error
	if dto.Info != nil {
		err = s.repository.UpdateApartmentHeavy(id, &dtos.ApartmentHeavyUpdateDTO{
			OwnerID: dto.OwnerID,
			Price:   dto.Price,
			Info:    *dto.Info,
		})
	} else {
		err = s.repository.UpdateApartmentLight(id, &dtos.ApartmentLightUpdateDTO{
			OwnerID: dto.OwnerID,
			Price:   dto.Price,
		})
	}

	if err != nil {
		var ve *servererrors.ValidationError
		if errors.As(err, &ve) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request form, beds & rooms must be positive integers"})
			logrus.WithField("Time", time.Now().String()).Info("400: Bad Request")
			return
		}

		var ue *servererrors.ForbiddenAccessError
		if errors.As(err, &ue) {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden request"})
			logrus.WithField("Time", time.Now().String()).Info("403: Forbidden")
			return
		}

		var nfe *servererrors.NotFoundError
		if errors.As(err, &nfe) {
			logrus.WithField("Time", time.Now().String()).Warnf("Could not find apartment with id %s", id)
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		logrus.WithField("Time", time.Now().String()).Warn("500: Internal Server Error")
		return
	}

	c.Status(http.StatusOK)
}

func (s *InnerServer) getApartmentsFiltered(c *gin.Context) {
	// Example: GET /apartments?city=Budapest&rooms=2&beds=1

	allowed := map[string]bool{
		"city":   true,
		"rooms":  true,
		"beds":   true,
		"limit":  false,
		"offset": false,
	}

	for key := range c.Request.URL.Query() {
		if !allowed[key] {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("unknown filter parameter: '%s'", key),
			})
			return
		}
	}

	city := c.Query("city")
	rooms := c.Query("rooms")
	beds := c.Query("beds")

	filter := map[string]string{}
	if city != "" {
		filter["city"] = city
	}
	if rooms != "" {
		filter["rooms"] = rooms
	}
	if beds != "" {
		filter["beds"] = beds
	}

	aps, err := s.repository.GetApartments(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch apartments"})
		logrus.WithField("Time", time.Now().String()).Warn("500: Internal Server Error")
		return
	}

	response := []dtos.ShortApartmentResponse{}
	for _, ap := range *aps {
		response = append(response, dtos.ShortApartmentResponse{
			Id:      ap.ID,
			OwnerID: ap.OwnerID,
			Address: ap.Address,
			Price:   ap.Price,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"count":      len(response),
		"apartments": response,
	})
}

func (s *InnerServer) getApartmentById(c *gin.Context) {
	id := c.Param("id")

	ap, brs, err := s.repository.GetApartment(id)
	if err != nil {
		var nfe *servererrors.NotFoundError
		if errors.As(err, &nfe) {
			c.JSON(http.StatusNotFound, gin.H{"error": "apartment not found"})
			logrus.WithField("Time", time.Now().String()).Infof("404 - Could not find apartment with id %s", id)
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch apartments"})
		logrus.WithField("Time", time.Now().String()).Warn("500: Internal Server Error")
		return
	}

	info := map[string]string{}
	if len(ap.Descriptions) > 0 {
		var tmp map[string]string
		if err := json.Unmarshal([]byte(ap.Descriptions[0].Description), &tmp); err != nil {
			logrus.WithField("Time", time.Now().String()).
				Warnf("invalid JSON description string for ap_id=%s: %v", id, err)
		} else if tmp != nil {
			info = tmp
		}

		if ap.Descriptions[0].Rooms > 0 {
			info["rooms"] = strconv.Itoa(ap.Descriptions[0].Rooms)
		}

		if ap.Descriptions[0].Beds >= 0 {
			info["beds"] = strconv.Itoa(ap.Descriptions[0].Beds)
		}
	}

	response := dtos.MediumApartmentResponse{
		Id:      ap.ID,
		OwnerID: ap.OwnerID,
		Address: ap.Address,
		Price:   ap.Price,
		Info:    info,
	}

	c.JSON(http.StatusOK, gin.H{
		"apartment": response,
		"bookings":  brs,
	})
}

func (s *InnerServer) bookApartment(c *gin.Context) {
	var dto dtos.BookingCreateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	booking, err := s.repository.CreateBooking(&dto)
	if err != nil {
		var oe *servererrors.OverlapError
		if errors.As(err, &oe) {
			c.JSON(http.StatusConflict, gin.H{"error": "time overlap with other booking"})
			logrus.WithField("Time", time.Now().String()).Infof("Booking time overlap on apartment %s", dto.ApartmentID)
			return
		}

		var nfe *servererrors.NotFoundError
		if errors.As(err, &nfe) {
			logrus.WithField("Time", time.Now().String()).Warnf("Could not find apartment with id %s", dto.ApartmentID)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		logrus.WithField("Time", time.Now().String()).Warn("500: Internal Server Error")
		return
	}

	logrus.WithField("Time", time.Now().String()).
		WithFields(logrus.Fields{"id": booking.ID, "time_from": booking.TimeFrom, "time_to": booking.TimeTo}).
		Infof("user %s booked apartment %s", booking.UserID, booking.UserID)

	c.JSON(http.StatusOK, dtos.BookingResponse{
		Id:          booking.ID,
		ApartmentID: booking.ApartmentID,
		Address:     booking.Apartment.Address,
		UserID:      booking.UserID,
		TimeFrom:    booking.TimeFrom,
		TimeTo:      booking.TimeTo,
	})
}

func (s *InnerServer) getApartmentsByOwner(c *gin.Context) {
	id := c.Param("id")

	aps, err := s.repository.GetApartmentsByOwner(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		logrus.WithField("Time", time.Now().String()).Warn("500: Internal Server Error")
		return
	}

	response := []dtos.FullApartmentResponse{}
	for _, ap := range *aps {
		info := map[string]string{}
		if len(ap.Descriptions) > 0 {
			var tmp map[string]string
			if err := json.Unmarshal([]byte(ap.Descriptions[0].Description), &tmp); err != nil {
				logrus.WithField("Time", time.Now().String()).
					Warnf("invalid JSON description string for ap_id=%s: %v", id, err)
			} else if tmp != nil {
				info = tmp
			}

			if ap.Descriptions[0].Rooms > 0 {
				info["rooms"] = strconv.Itoa(ap.Descriptions[0].Rooms)
			}

			if ap.Descriptions[0].Beds >= 0 {
				info["beds"] = strconv.Itoa(ap.Descriptions[0].Beds)
			}
		}

		bookings := []dtos.ShortBookingResponse{}
		for _, b := range ap.Bookings {
			bookings = append(bookings, dtos.ShortBookingResponse{
				Id:       b.ID,
				UserID:   b.UserID,
				TimeFrom: b.TimeFrom,
				TimeTo:   b.TimeTo,
			})
		}

		response = append(response, dtos.FullApartmentResponse{
			Id:       ap.ID,
			Address:  ap.Address,
			OwnerID:  ap.OwnerID,
			Price:    ap.Price,
			Info:     info,
			Bookings: bookings,
		})
	}

	c.JSON(http.StatusOK, gin.H{"count": len(response), "apartments": response})
}

func (s *InnerServer) getBookingsByUser(c *gin.Context) {
	id := c.Param("id")

	bs, err := s.repository.GetBookingsByUser(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		logrus.WithField("Time", time.Now().String()).Warn("500: Internal Server Error")
		return
	}

	response := []dtos.BookingResponse{}
	for _, booking := range *bs {
		response = append(response, dtos.BookingResponse{
			Id:          booking.ID,
			ApartmentID: booking.ApartmentID,
			Address:     booking.Apartment.Address,
			UserID:      booking.UserID,
			TimeFrom:    booking.TimeFrom,
			TimeTo:      booking.TimeTo,
		})
	}

	c.JSON(http.StatusOK, gin.H{"count": len(response), "bookings": response})
}
