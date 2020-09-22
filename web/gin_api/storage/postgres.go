package storage

import (
	"fmt"
	"github.com/jomifepe/gin_api/logging"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConfigurePostgresDB() *DBConn {
	var (
		dbHost = viper.GetString("DATABASE_HOST")
		dbPort = viper.GetString("DATABASE_PORT")
		dbUser = viper.GetString("POSTGRES_USER")
		dbPass = viper.GetString("POSTGRES_PASSWORD")
		dbName = viper.GetString("POSTGRES_DB")
	)
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)

	logging.Logger.WithFields(logrus.Fields{
		"host": dbHost,
		"port": dbPort,
		"name": dbName,
		"user": dbUser,
	}).Infof("[DB] Connecting...")

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info),
		Logger: logging.NewGORMLogger(viper.GetString("LOG_LEVEL")),
	})
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{
			"error": err,
		}).Panicln("[DB] Failed to connect")
	}

	dbConn := &DBConn{gormDB}
	dbConn.MigrateDatabase()

	logging.Logger.Infoln("[DB] Successfully connected")
	return dbConn
}