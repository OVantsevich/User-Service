// Package model model User
package model

import "time"

// User model info
// @Description User account information
type User struct {
	ID       string    `json:"id"`
	Login    string    `json:"login" validate:"required,alphanum,gte=5,lte=20"`
	Email    string    `json:"email" validate:"required,email" format:"email"`
	Password string    `json:"password" validate:"required"`
	Name     string    `json:"name" validate:"required,alpha,gte=2,lte=25"`
	Age      int       `json:"age" validate:"required,gte=0,lte=100"`
	Token    string    `json:"token"`
	Role     string    `json:"role"`
	Created  time.Time `json:"created" example:"2021-05-25T00:53:16.535668Z" format:"date-time"`
	Updated  time.Time `json:"updated" example:"2021-05-25T00:53:16.535668Z" format:"date-time"`
}
