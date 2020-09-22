package storage

import (
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jomifepe/gin_api/api/auth"
	"github.com/jomifepe/gin_api/logging"
	"github.com/jomifepe/gin_api/model"
	"gorm.io/gorm"
)

type DBConn struct {
	DB *gorm.DB
}

// MigrateDatabase auto migrates the database with existing model structs
func (conn *DBConn) MigrateDatabase() {
	logging.Logger.Infoln("[DB] Running migrations...")

	if err := conn.DB.AutoMigrate(
		&model.User{},
		&model.Task{},
		&auth.AccessDetails{},
	); err != nil {
		logging.Logger.Panicln("[DB] Failed to migrate database", err)
	}
}