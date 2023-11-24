package domain

import "gorm.io/gorm"

type User struct {
	*gorm.Model `json:"-"`
	ID          uint   `json:"id" gorm:"unique;not null"`
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Email       string `json:"email" validate:"email"`
	Password    string `json:"password" validate:"min=8,max=20"`
	Phone       string `json:"phone"`
	Blocked     bool   `json:"blocked" gorm:"default:false"`
	Isadmin     bool   `json:"is_admin" gorm:"default:false"`
}
type Address struct {
	Id        int    `json:"id" gorm:"unique;not null"`
	UserID    uint   `json:"user_id"`
	User      User   `json:"-" gorm:"foreignkey:UserID"`
	Name      string `json:"name" validate:"required"`
	HouseName string `json:"house_name" validate:"required"`
	Street    string `json:"street" validate:"required"`
	City      string `json:"city" validate:"required"`
	State     string `json:"state" validate:"required"`
	Pin       string `json:"pin" validate:"required"`
}
