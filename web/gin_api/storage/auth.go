package storage

import (
	"errors"
	"github.com/jomifepe/gin_api/api/auth"
	"github.com/jomifepe/gin_api/logging"
	"github.com/sirupsen/logrus"
)

type AuthStore struct {
	DBConn
}

// NewAuthStore return an AuthStore.
func NewAuthStore(conn *DBConn) *AuthStore {
	return &AuthStore{
		DBConn{DB: conn.DB},
	}
}

// RegisterAccess registers a new access (auth.AccessDetails) on the database.
// Deletes existing ones with the same user id before creating.
func (conn *DBConn) RegisterAccess(t auth.AccessDetails) error {
	// Delete other existing access entries
	go func () /* no concurrency, can be async */ {
		if dr := conn.DB.Delete(&t, "user_id = ? AND access_uuid != ?", t.UserID, t.AccessUUID); dr.Error != nil {
			logging.Logger.WithFields(logrus.Fields{
				"user_id": t.UserID,
			}).Warnln("[DB] Couldn't delete existing access before creating new one")
		} else if dr.RowsAffected > 0 {
			logging.Logger.WithFields(logrus.Fields{
				"user_id": t.UserID,
			}).Warnf("[DB] Deleted %v old access entries", dr.RowsAffected)
		}
	}()

	result := conn.DB.Create(&t)
	if result.Error != nil {
		logging.Logger.WithFields(logrus.Fields{
			"error": result.Error.Error(),
		}).Errorln("[DB] Couldn't create authentication details")
		return result.Error
	}
	if result.RowsAffected <= 0 {
		logging.Logger.WithFields(logrus.Fields{
			"task": t,
		}).Errorln("[DB] No rows were affected when creating task")
		return errors.New("now rows were affected")
	}
	logging.Logger.WithFields(logrus.Fields{
		"user_id": t.UserID,
		"access_uuid": t.AccessUUID,
	}).Infoln("[DB] Created user access")
	return nil
}

func (conn *DBConn) GetAccess(uuid string) (auth.AccessDetails, error) {
	var td auth.AccessDetails
	if result := conn.DB.Where("access_uuid = ?", uuid).First(&td); result.Error != nil {
		logging.Logger.WithFields(logrus.Fields{
			"uuid": uuid,
			"error": result.Error,
		}).Errorln("[DB] Couldn't get access by uuid")
		return auth.AccessDetails{}, result.Error
	}
	return td, nil
}

func (conn *DBConn) DeleteAccess(t auth.AccessDetails) error {
	if result := conn.DB.Delete(auth.AccessDetails{}, "access_uuid = ?", t.AccessUUID); result.Error != nil {
		logging.Logger.WithFields(logrus.Fields{
			"user_id": t.UserID,
			"access_uuid": t.AccessUUID,
			"error": result.Error,
		}).Errorln("[DB] Couldn't delete user access")
		return result.Error
	}
	logging.Logger.WithFields(logrus.Fields{
		"user_id": t.UserID,
		"access_uuid": t.AccessUUID,
	}).Infoln("[DB] Deleted user access")
	return nil
}