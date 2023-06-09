package models

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	RoleUser     = "user"
	RoleAdmin    = "admin"
	RoleEmployee = "employee"
)

type User struct {
	mgm.DefaultModel  `bson:",inline"`
	Email             string               `json:"email" bson:"email"`
	Password          string               `json:"-" bson:"password"`
	Name              string               `json:"name" bson:"name"`
	Role              string               `json:"role" bson:"role"`
	MailVerified      bool                 `json:"mail_verified" bson:"mail_verified"`
	PlacedOrderIds    []primitive.ObjectID `json:"placedOrderIds" bson:"placedOrderIds"`
	CompletedOrderIds []primitive.ObjectID `json:"completedOrderIds" bson:"completedOrderIds"`
	CanceledOrderIds  []primitive.ObjectID `json:"canceledOrderIds" bson:"canceledOrderIds"`
}

type UserClaims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
	Type  string `json:"type"`
}

func NewUser(email string, password string, name string, role string) *User {
	return &User{
		Email:        email,
		Password:     password,
		Name:         name,
		Role:         role,
		MailVerified: false,
	}
}

func (model *User) CollectionName() string {
	return "users"
}

// You can override Collection functions or CRUD hooks
// https://github.com/Kamva/mgm#a-models-hooks
// https://github.com/Kamva/mgm#collections
