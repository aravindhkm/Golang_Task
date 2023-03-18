package models

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/kamva/mgm/v3"
)

type Admin struct {
	mgm.DefaultModel `bson:",inline"`
	Email            string `json:"email" bson:"email"`
	Password         string `json:"-" bson:"password"`
	Name             string `json:"name" bson:"name"`
	Role             string `json:"role" bson:"role"`
	MailVerified     bool   `json:"mail_verified" bson:"mail_verified"`
}

type AdminClaims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
	Type  string `json:"type"`
}

func NewAdmin(email string, password string, name string, role string) *Admin {
	return &Admin{
		Email:        email,
		Password:     password,
		Name:         name,
		Role:         "admin",
		MailVerified: false,
	}
}

func (model *Admin) CollectionName() string {
	return "admin"
}

// You can override Collection functions or CRUD hooks
// https://github.com/Kamva/mgm#a-models-hooks
// https://github.com/Kamva/mgm#collections
