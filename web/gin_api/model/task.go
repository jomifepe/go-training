package model

import (
	"fmt"
	"time"
)

// Task - Information about a task to be done, and if it's completed or not
type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description" validate:"required,min=1,max=124" binding:"required"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (u Task) String() string {
	return fmt.Sprintf("Task #%v:\nDescription: %v\nCompleted: %v\nCreated: %v\nUpdated:%v\n",
		u.ID, u.Description, u.Completed, u.CreatedAt, u.UpdatedAt)
}