package config

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type ServiceConfig struct {
	Port         int
	Logger       *zap.SugaredLogger
	Db           *gorm.DB
	ClientId     string
	ClientSecret string
	Mux          *chi.Mux
}

func InitServiceConfig() *ServiceConfig {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		sentry.CaptureException(err)
		log.Fatalf("error reading service configuration: %v", err)
	}

	logFilePath := path.Join(
		viper.GetString("service.log_file_path"),
		viper.GetString("service.log_file_name"),
	)

	logFile, err := os.Create(logFilePath)
	if err != nil {
		sentry.CaptureException(err)
		log.Fatalf("error creating log file: %v", err)
	}

	db, err := buildDatbaseConnection(
		viper.GetString("service.environment"),
		viper.GetString("database.host"),
		viper.GetString("database.username"),
		viper.GetString("database.password"),
		viper.GetString("database.dbname"),
		viper.GetInt("database.port"))

	if err != nil {
		sentry.CaptureException(err)
		log.Fatalf("error connecting to database server: %v", err)
	}

	return &ServiceConfig{
		Port:         viper.GetInt("service.port"),
		Logger:       buildLogger(logFile),
		Db:           db,
		ClientId:     viper.GetString("auth_service.client_id"),
		ClientSecret: viper.GetString("auth_service.client_secret"),
		Mux:          initServiceMux(),
	}
}

func buildLogger(f *os.File) *zap.SugaredLogger {
	pe := zap.NewProductionEncoderConfig()
	fileEncoder := zapcore.NewJSONEncoder(pe)
	pe.EncodeTime = zapcore.ISO8601TimeEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(pe)

	level := zap.InfoLevel

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(f), level),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
	)

	logger := zap.New(core)
	return logger.Sugar()
}

func buildDatbaseConnection(env, host, username, password, dbName string, port int) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		host,
		username,
		password,
		dbName,
		port)

	dbLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,          // Don't include params in the SQL log
			Colorful:                  false,         // Disable color
		},
	)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger:      dbLogger,
		NowFunc:     time.Now,
		PrepareStmt: true,
	})

	if err != nil {
		return nil, err
	}

	return db, nil
}
