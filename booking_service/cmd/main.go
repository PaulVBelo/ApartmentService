package main

import (
	"booking_service/internal/config"
	"booking_service/internal/db"
	"booking_service/internal/server"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	// LOGGER SETUP
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
		PrettyPrint:     false,
	})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)

	// LOAD CONFIG FROM DOTENV
	cfg, err := config.LoadConfig()
	if err != nil {
		logrus.Fatal("Config not loaded")
	}

	logrus.WithFields(logrus.Fields{
		"DB_HOST": cfg.DB_HOST,
		"DB_PORT": cfg.DB_PORT,
		"DB_USER": cfg.DB_USER,
		"DB_NAME": cfg.DB_NAME,
		"SERV":    cfg.SERV_HOST + ":" + cfg.SERV_PORT,
	}).Info("Booking config loaded")

	// CONNECT TO DB
	conn, err := db.ConnectPostgres(cfg)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to connect to database")
	}

	sqlDB, err := conn.DB()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to access SQL DB handle")
	}

	// INIT SERVER
	srv := server.NewServer(conn)
	host := cfg.SERV_HOST
	port := cfg.SERV_PORT

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	http_server := srv.Run(host, port)

	// GRACEFUL SHUTDOWN
	<-quit

	logrus.Info("Shutting down server...")
	server.GracefulShutdown(http_server, 5*time.Second)

	if err := sqlDB.Close(); err != nil {
		logrus.WithError(err).Warn("Error closing DB conn")
	} else {
		logrus.Info("DB connection closed")
	}
}
