package models

import (
	"github.com/kamva/mgm/v3"
)

type Employee struct {
	mgm.DefaultModel `bson:",inline"`
	Email            string `json:"email" bson:"email"`
	Password         string `json:"-" bson:"password"`
	Name             string `json:"name" bson:"name"`
	Mobile           uint   `json:"mobile" bson:"mobile"`
	Address          string `json:"address" bson:"address"`
	Role             string `json:"role" bson:"role"`
	MailVerified     bool   `json:"mail_verified" bson:"mail_verified"`
}

func NewEmployee(
	email string,
	password string,
	name string,
	mobile uint,
	address string) *Employee {
	return &Employee{
		Email:        email,
		Password:     password,
		Name:         name,
		Mobile:       mobile,
		Address:      address,
		Role:         "employee",
		MailVerified: false,
	}
}

func (model *Employee) CollectionName() string {
	return "employee"
}

// You can override Collection functions or CRUD hooks
// https://github.com/Kamva/mgm#a-models-hooks
// https://github.com/Kamva/mgm#collections
