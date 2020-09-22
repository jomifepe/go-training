package resource

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jomifepe/gin_api/api/auth"
	"github.com/jomifepe/gin_api/api/middleware"
	"github.com/jomifepe/gin_api/logging"
	"github.com/jomifepe/gin_api/model"
	"net/http"
)

var (
	errTaskCreate        = gin.H{"message": "Couldn't create new task"}
	errTaskInvalidFields = gin.H{"message": "The specified task has invalid fields"}
	errTaskToggle        = gin.H{"message": "Couldn't toggle task completed"}
	errTaskDelete        = func(param ...interface{}) map[string]interface{} {
		return gin.H{"message": fmt.Sprintf("Couldn't delete task with id %v", param...)}
	}
	errTaskIdNotFound = func(param ...interface{}) map[string]interface{} {
		return gin.H{"message": fmt.Sprintf("No task with the id %v was found", param...)}
	}
)

// taskStore is used to define the database calls used by the route group define in this file
type taskStore interface {
	CreateTask(task model.Task) (model.Task, error)
	DeleteAccess(accessDetails auth.AccessDetails) error
	UpdateTask(task model.Task) (model.Task, error)
	GetTask(id int) (model.Task, error)
	GetAllTasks() ([]model.Task, error)
	DeleteTask(id int) error
}

// TaskResource holds a TaskStore interface, used to communicate with the database
type TaskResource struct {
	Store taskStore
}

// NewTaskResource initializes the TaskResource with an existing TaskStore
func NewTaskResource(store taskStore) *TaskResource {
	return &TaskResource{
		Store: store,
	}
}

// MountTaskRoutesTo defines new routes regarding Tasks on an existing gin.RouterGroup or gin.Engine
func (tr *TaskResource) MountTaskRoutesTo(r gin.IRouter) (rg *gin.RouterGroup) {
	idParam := middleware.Param{Key: "id", ExampleValue: -1}

	rg = r.Group("/tasks"); {
		rg.GET("", tr.handleGetTasks)
		rg.POST("", tr.handleCreateTask)
		withId := rg.Group("", middleware.ExtractParam(idParam)); {
			withId.GET("/:id", tr.handleGetTask)
			withId.PUT("/:id", tr.handleUpdateTask)
			withId.DELETE("/:id", tr.handleDeleteTask)
			withId.PUT("/:id/toggle", tr.handleTaskToggle)
		}
	}
	return
}

// handleCreateTask validates the task sent on the request body and inserts it, if it's valid, on the database
func (tr *TaskResource) handleCreateTask(c *gin.Context) {
	var t model.Task
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusUnprocessableEntity, errTaskInvalidFields)
		return
	}

	newTask, err := tr.Store.CreateTask(t)
	if err != nil {
		c.JSON(http.StatusBadRequest, errTaskCreate)
		return
	}
	c.JSON(http.StatusCreated, newTask)
}

func (tr *TaskResource) handleGetTask(c *gin.Context) {
	id := c.GetInt("id")
	t, err := tr.Store.GetTask(id)
	if err != nil {
		c.JSON(http.StatusNotFound, errTaskIdNotFound(id))
		return
	}

	c.JSON(http.StatusOK, t)
}

// handleGetTasks returns all the existing tasks
func (tr *TaskResource) handleGetTasks(c *gin.Context) {
	tc, err := tr.Store.GetAllTasks()
	if err != nil {
		logging.Logger.Errorln("[API] Failed to get all tasks", err)
	}
	c.JSON(http.StatusOK, tc)
}

// handleUpdateTask validates the task passed on the request body, and updates it (if it's valid),
// using the <id> passed on the request url path
func (tr *TaskResource) handleUpdateTask(c *gin.Context) {
	id := c.GetInt("id")
	var t model.Task
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusUnprocessableEntity, errTaskInvalidFields)
		return
	}

	t.ID = id
	updatedTask, err := tr.Store.UpdateTask(t)
	if err != nil {
		c.JSON(http.StatusBadRequest, errTaskDelete(id))
		return
	}

	c.JSON(http.StatusOK, updatedTask)
}

// handleDeleteTask deletes a task from the database using the <id> passed on the request url path
func (tr *TaskResource) handleDeleteTask(c *gin.Context) {
	id := c.GetInt("id")
	if err := tr.Store.DeleteTask(id); err != nil {
		c.JSON(http.StatusNotFound, errTaskDelete(id))
		return
	}
	c.JSON(http.StatusNoContent, "")
}

// handleTaskToggle toggles a tasks completed field, using the <id> passed on the request url path
func (tr *TaskResource) handleTaskToggle(c *gin.Context) {
	id := c.GetInt("id")
	t, err := tr.Store.GetTask(id)
	if err != nil {
		c.JSON(http.StatusNotFound, errTaskIdNotFound(id))
		return
	}

	t.Completed = !t.Completed
	updatedTask, err := tr.Store.UpdateTask(t)
	if err != nil {
		c.JSON(http.StatusNotFound, errTaskToggle)
		return
	}

	c.JSON(http.StatusOK, updatedTask)
}