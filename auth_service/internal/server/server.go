package server

import (
	"auth_service/internal/models"
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Server struct {
	router *gin.Engine
	db     *gorm.DB
}

func NewServer(db *gorm.DB) *Server {
	router := gin.Default()
	s := &Server{
		router: router,
		db:     db,
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.router.POST("/auth/register", s.register)
	s.router.POST("/auth/login", s.login)
	s.router.GET("/health", health)
}

func (s *Server) Run(host, port string) *http.Server {
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

// =========================================================== REGISTER

func (s *Server) register(c *gin.Context) {
	var dto models.Request

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		logrus.WithField("Time", time.Now().String()).Info("400: Bad Request")
		return
	}

	if len(dto.Password) < 8 || len(dto.Password) > 16 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid password length"})
		logrus.WithField("Time", time.Now().String()).Info("400: Bad Request")
		return
	}

	enc_pass, err := models.HashPassword(dto.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "encryption failed"})
		logrus.WithField("Time", time.Now().String()).Warn("Failed to encrypt valid password")
		return
	}

	user := models.User{
		ID:        uuid.New().String(),
		Email:     dto.Email,
		EncPass:   enc_pass,
		CreatedAt: time.Now(),
	}

	err = s.db.Create(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrCheckConstraintViolated) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email format"})
			logrus.WithField("Time", time.Now().String()).Info("400: Invalid email format")
			return
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				c.JSON(http.StatusConflict, gin.H{"error": "email already in use"})
				logrus.WithField("Time", time.Now().String()).Info("409: Email already exists")
				return
			}
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		logrus.WithField("Time", time.Now().String()).Warn("Failed to create user")
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "registered"})
	logrus.WithField("Time", time.Now().String()).WithFields(logrus.Fields{
		"user_id": user.ID, "email": user.Email, "created_at": user.CreatedAt.String(),
	}).Info("201: User registered successfuly")
}

/// ===========================================================  LOGIN

func (s *Server) login(c *gin.Context) {
	var dto models.Request

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		logrus.WithField("Time", time.Now().String()).Info("400: Bad Request")
		return
	}

	if !models.IsValid(dto.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email format"})
		logrus.WithField("Time", time.Now().String()).Info("400: Invalid email format")
		return
	}

	var user models.User
	if err := s.db.Where("email = ?", dto.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid credentials"})
		logrus.WithField("Time", time.Now().String()).Info("400: Invalid credentials")
		return
	}

	if ok := models.CheckPasswordHash(dto.Password, user.EncPass); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid credentials"})
		logrus.WithField("Time", time.Now().String()).Info("400: Invalid credentials")
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "logged_in", "user_id": user.ID})
	logrus.WithField("Time", time.Now().String()).WithFields(logrus.Fields{
		"user_id": user.ID, "email": user.Email,
	}).Info("200: Logged in")
}
