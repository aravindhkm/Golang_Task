package services

import (
	db "Hdfc_Assignment/models"
	"errors"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser create a user record
func CreateUser(name string, email string, plainPassword string) (*db.User, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("cannot generate hashed password")
	}

	user := db.NewUser(email, string(password), name, db.RoleUser)
	err = mgm.Coll(user).Create(user)
	if err != nil {
		return nil, errors.New("cannot create new user")
	}

	return user, nil
}

// CreateAdmin create a admin record
func CreateAdmin(name string, email string, plainPassword string) (*db.Admin, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("cannot generate hashed password")
	}

	admin := db.NewAdmin(email, string(password), name, db.RoleAdmin)
	err = mgm.Coll(admin).Create(admin)
	if err != nil {
		return nil, errors.New("cannot create new user")
	}

	return admin, nil
}

// CreateEmployee create a employee record
// func CreateEmployee(name string, email string, plainPassword string) (*db.Employee, error) {
// 	password, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
// 	if err != nil {
// 		return nil, errors.New("cannot generate hashed password")
// 	}

// 	employee := db.NewEmployee(email, string(password), name, db.RoleEmployee)
// 	err = mgm.Coll(employee).Create(employee)
// 	if err != nil {
// 		return nil, errors.New("cannot create new user")
// 	}

// 	return employee, nil
// }

// FindUserById find user by id
func FindUserById(userId primitive.ObjectID) (*db.User, error) {
	user := &db.User{}
	err := mgm.Coll(user).FindByID(userId, user)
	if err != nil {
		return nil, errors.New("cannot find user")
	}

	return user, nil
}

// FindUserById find user by id
func FindAdminById(userId primitive.ObjectID) (*db.Admin, error) {
	admin := &db.Admin{}
	err := mgm.Coll(admin).FindByID(userId, admin)
	if err != nil {
		return nil, errors.New("cannot find user")
	}

	return admin, nil
}

// FindUserById find user by id
func FindEmployeeById(userId primitive.ObjectID) (*db.Employee, error) {
	employee := &db.Employee{}
	err := mgm.Coll(employee).FindByID(userId, employee)
	if err != nil {
		return nil, errors.New("cannot find user")
	}

	return employee, nil
}

// FindUserByEmail find user by email
func FindUserByEmail(email string) (*db.User, error) {
	user := &db.User{}
	err := mgm.Coll(user).First(bson.M{"email": email}, user)
	if err != nil {
		return nil, errors.New("cannot find user")
	}

	return user, nil
}

// FindAdminByEmail find admin by email
func FindAdminByEmail(email string) (*db.Admin, error) {
	admin := &db.Admin{}
	err := mgm.Coll(admin).First(bson.M{"email": email}, admin)
	if err != nil {
		return nil, errors.New("cannot find user")
	}

	return admin, nil
}

// FindEmployeeByEmail find employee by email
func FindEmployeeByEmail(email string) (*db.Employee, error) {
	employee := &db.Employee{}
	err := mgm.Coll(employee).First(bson.M{"email": email}, employee)
	if err != nil {
		return nil, errors.New("cannot find user")
	}

	return employee, nil
}

// CheckUserMail search user by email, return error if someone uses
func CheckUserMail(email string, role string) error {

	if role == db.RoleAdmin {
		admin := &db.Admin{}
		adminCollection := mgm.Coll(admin)
		err := adminCollection.First(bson.M{"email": email}, admin)
		if err == nil {
			return errors.New("email is already in use")
		}
	} else if role == db.RoleUser {
		user := &db.User{}
		userCollection := mgm.Coll(user)
		err := userCollection.First(bson.M{"email": email}, user)
		if err == nil {
			return errors.New("email is already in use")
		}
	} else if role == db.RoleEmployee {
		employee := &db.User{}
		employeeCollection := mgm.Coll(employee)
		err := employeeCollection.First(bson.M{"email": email}, employee)
		if err == nil {
			return errors.New("email is already in use")
		}
	}

	return nil
}
