package model

import "fmt"

type AuthUser struct {
	Email     string `json:"email" validate:"required,email" gorm:"unique" binding:"required"`
	Password  string `json:"password,omitempty" validate:"required,min=6,max=72"`
}

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name" validate:"required,alpha,min=1,max=1024" binding:"required"`
	LastName  string `json:"last_name" validate:"required,alpha,min=1,max=1024" binding:"required"`
	Email     string `json:"email" validate:"required,email" gorm:"unique" binding:"required"`
	Password  string `json:"password,omitempty" validate:"required,min=6,max=72"`
	Active    bool   `json:"active" gorm:"default:true"`
}

func (u User) String() string {
	return fmt.Sprintf("#%v: %v %v, %v", u.ID, u.FirstName, u.LastName, u.Email)
}
