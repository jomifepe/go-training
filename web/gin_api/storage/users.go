package storage

import (
	"errors"
	"github.com/jomifepe/gin_api/logging"
	"github.com/jomifepe/gin_api/model"
	"github.com/sirupsen/logrus"
)

type UserStore struct {
	DBConn
}

// NewAuthStore return an AuthStore.
func NewUserStore(conn *DBConn) *UserStore {
	return &UserStore{
		DBConn{DB: conn.DB},
	}
}

func (conn *DBConn) CreateUser(u model.User) (model.User, error) {
	result := conn.DB.Create(&u)
	if result.Error != nil {
		logging.Logger.WithFields(logrus.Fields{
			"user": u,
			"error": result.Error,
		}).Errorln("[DB] Couldn't create user")
		return model.User{}, result.Error
	}
	if result.RowsAffected <= 0 {
		logging.Logger.WithFields(logrus.Fields{
			"user": u,
		}).Errorln("[DB] No rows were affected when creating user")
		return model.User{}, errors.New("now rows were affected")
	}
	logging.Logger.WithFields(logrus.Fields{
		"id": u.ID,
	}).Infoln("[DB] Created new user")
	return u, nil
}

// GetAllUsers returns all users from the database.
// By default, it omits sensitive fields, like passwords.
// In order to get all fields, pass in an empty string.
func (conn *DBConn) GetAllUsers(omitFields ...string) ([]model.User, error) {
	if len(omitFields) == 0 {
		omitFields = []string{"password"}
	}
	var users []model.User
	result := conn.DB.Omit(omitFields...).Find(&users)
	if result.Error != nil {
		logging.Logger.WithFields(logrus.Fields{
			"error": result.Error,
		}).Errorln("[DB] Couldn't get all users")
		return []model.User{}, result.Error
	}
	return users, nil
}

// GetAllUsers returns an existing user from the database, searches by <paramName> with the <param> value.
// By default, it omits sensitive fields, like passwords.
// In order to get all fields, pass in an empty string.
func (conn *DBConn) GetUserBy(paramName string, param interface{}, omitFields ...string) (model.User, error) {
	if len(omitFields) == 0 {
		omitFields = []string{"password"}
	}
	var user model.User
	if result := conn.DB.Omit(omitFields...).First(&user, paramName + " = ?", param); result.Error != nil {
		logging.Logger.WithFields(logrus.Fields{
			"error": result.Error,
			"param": param,
			"omit_fields": omitFields,
		}).Errorln("[DB] Couldn't get user with by", paramName)
		return model.User{}, result.Error
	}
	return user, nil
}

func (conn *DBConn) UpdateUser(u model.User) (model.User, error) {
	result := conn.DB.Model(&u).Select("first_name", "last_name", "email").Updates(u)
	if result.Error != nil && result.RowsAffected > 0 {
		logging.Logger.WithFields(logrus.Fields{
			"user": u,
			"error": result.Error,
		}).Errorln("[DB] Couldn't update user")
		return model.User{}, result.Error
	}
	updatedUser, err := conn.GetUserBy("id", u.ID)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{
			"error": err,
		}).Errorln("[DB] Couldn't get updated user")
		return model.User{}, result.Error
	}
	logging.Logger.WithFields(logrus.Fields{
		"id": u.ID,
	}).Infoln("[DB] Updated existing user")
	return updatedUser, nil
}

func (conn *DBConn) DeleteUser(id int) error {
	if result := conn.DB.Delete(&model.User{}, id); result.Error != nil && result.RowsAffected > 0 {
		logging.Logger.WithFields(logrus.Fields{
			"user_id": id,
			"error": result.Error,
		}).Errorln("[DB] Couldn't delete user by id")
		return result.Error
	}
	logging.Logger.WithFields(logrus.Fields{
		"id": id,
	}).Infoln("[DB] Deleted existing user")
	return nil
}
