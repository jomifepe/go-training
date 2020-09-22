package storage

import (
	"errors"
	"github.com/jomifepe/gin_api/logging"
	"github.com/jomifepe/gin_api/model"
	"github.com/sirupsen/logrus"
)

type TaskStore struct {
	DBConn
}

// NewAuthStore return an AuthStore.
func NewTaskStore(conn *DBConn) *TaskStore {
	return &TaskStore{
		DBConn{DB: conn.DB},
	}
}

func (conn *DBConn) CreateTask(t model.Task) (model.Task, error) {
	result := conn.DB.Create(&t)
	if result.Error != nil {
		logging.Logger.WithFields(logrus.Fields{
			"error": result.Error,
			"task": t,
		}).Errorln("[DB] Couldn't create task")
		return model.Task{}, result.Error
	}
	if result.RowsAffected <= 0 {
		logging.Logger.WithFields(logrus.Fields{
			"task": t,
		}).Errorln("[DB] No rows were affected when creating task")
		return model.Task{}, errors.New("now rows were affected")
	}
	logging.Logger.WithFields(logrus.Fields{
		"id": t.ID,
	}).Infoln("[DB] Created new task")
	return t, nil
}

func (conn *DBConn) GetAllTasks() ([]model.Task, error) {
	var tasks []model.Task
	result := conn.DB.Find(&tasks)
	if result.Error != nil {
		logging.Logger.WithFields(logrus.Fields{
			"error": result.Error,
		}).Errorln("[DB] Couldn't get all tasks")
		return []model.Task{}, result.Error
	}
	return tasks, nil
}

func (conn *DBConn) GetTask(id int) (model.Task, error) {
	var task model.Task
	if result := conn.DB.First(&task, "id = ?", id); result.Error != nil {
		logging.Logger.WithFields(logrus.Fields{
			"task_id": id,
			"error": result.Error,
		}).Errorln("[DB] Couldn't get task by id")
		return model.Task{}, result.Error
	}
	return task, nil
}

func (conn *DBConn) UpdateTask(t model.Task) (model.Task, error) {
	result := conn.DB.Model(&t).Select("description", "completed", "updated_at").Updates(t)
	if result.Error != nil && result.RowsAffected > 0 {
		logging.Logger.WithFields(logrus.Fields{
			"task": t,
			"error": result.Error,
		}).Errorln("[DB] Couldn't update task")
		return model.Task{}, result.Error
	}
	updatedTask, err := conn.GetTask(t.ID)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{
			"error": err,
		}).Errorln("[DB] Couldn't get updated task")
		return model.Task{}, result.Error
	}
	logging.Logger.WithFields(logrus.Fields{
		"id": t.ID,
	}).Infoln("[DB] Updated existing task")
	return updatedTask, nil
}

func (conn *DBConn) DeleteTask(id int) error {
	if result := conn.DB.Delete(&model.Task{}, id); result.Error != nil && result.RowsAffected > 0 {
		logging.Logger.WithFields(logrus.Fields{
			"task_id": id,
			"error": result.Error,
		}).Errorln("[DB] Couldn't delete task by id")
		return result.Error
	}
	logging.Logger.WithFields(logrus.Fields{
		"id": id,
	}).Infoln("[DB] Deleted existing task")
	return nil
}