package main

import (
	"fmt"

	"gopkg.in/go-playground/validator.v9"
)

// TaskCollection - A collection of tasks
type TaskCollection []Task

// Task - Information about a task to be done, and if it's completed or not
type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description" validate:"required,min=1,max=124"`
	Completed   bool   `json:"completed"`
	CreatedAt   string `json:"created_at"`
}

func validateTask(t Task) (errors []string, jsonErrors string) {
	// TODO: Improve and add custom validation messages
	v := validator.New()
	err := v.Struct(t)
	if err == nil {
		return nil, ""
	}

	valErrs := err.(validator.ValidationErrors)
	if len(valErrs) == 0 {
		return nil, ""
	}
	jsonErrors += "["
	for _, e := range valErrs {
		msg := fmt.Sprintf("Invalid '%v' field", e.Field())
		errors = append(errors, msg)
		jsonErrors += fmt.Sprintf("\"%v\"", msg)
	}
	jsonErrors += "]"

	fmt.Println("Validation errors", jsonErrors)
	return
}
