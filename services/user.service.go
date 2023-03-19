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

// UpdateUserOrder updates a user order record with id
func PlaceUserOrder(
	userId primitive.ObjectID,
	orderId primitive.ObjectID) error {
	user := &db.User{}
	err := mgm.Coll(user).FindByID(userId, user)
	if err != nil {
		return errors.New("cannot find note")
	}

	var getOrderId []primitive.ObjectID

	getOrderId = append(getOrderId, user.PlacedOrderIds...)
	getOrderId = append(getOrderId, orderId)

	user.PlacedOrderIds = getOrderId
	err = mgm.Coll(user).Update(user)

	if err != nil {
		return errors.New("cannot update")
	}

	return nil
}

// CompleteUserOrder updates a user order record with id
func CompleteUserOrder(
	userId primitive.ObjectID,
	orderId primitive.ObjectID) error {
	user := &db.User{}
	err := mgm.Coll(user).FindByID(userId, user)
	if err != nil {
		return errors.New("cannot find note")
	}

	var getOrderId []primitive.ObjectID
	var getCompleteId []primitive.ObjectID
	var isExist bool

	for _, arrData := range user.PlacedOrderIds {
		if arrData != orderId {
			getOrderId = append(getOrderId, arrData)
		} else if arrData == orderId {
			isExist = true
		}
	}

	if !isExist {
		return errors.New("invalid Order Id")
	}

	getCompleteId = append(getCompleteId, user.CompletedOrderIds...)
	getCompleteId = append(getCompleteId, orderId)

	user.PlacedOrderIds = getOrderId
	user.CompletedOrderIds = getCompleteId
	err = mgm.Coll(user).Update(user)

	if err != nil {
		return errors.New("cannot update")
	}

	return nil
}

// CompleteUserOrder updates a user order record with id
func CancelUserOrder(
	userId primitive.ObjectID,
	orderId primitive.ObjectID) error {
	user := &db.User{}
	err := mgm.Coll(user).FindByID(userId, user)
	if err != nil {
		return errors.New("cannot find note")
	}

	var getOrderId []primitive.ObjectID
	var getCompleteId []primitive.ObjectID
	var isExist bool

	for _, arrData := range user.PlacedOrderIds {
		if arrData != orderId {
			getOrderId = append(getOrderId, arrData)
		} else if arrData == orderId {
			isExist = true
		}
	}

	if !isExist {
		return errors.New("invalid Order Id")
	}

	getCompleteId = append(getCompleteId, user.CanceledOrderIds...)
	getCompleteId = append(getCompleteId, orderId)

	user.PlacedOrderIds = getOrderId
	user.CanceledOrderIds = getCompleteId
	err = mgm.Coll(user).Update(user)

	if err != nil {
		return errors.New("cannot update")
	}

	return nil
}

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

// FindUserById find user by id
func FindProductById(userId primitive.ObjectID) (*db.Product, error) {
	product := &db.Product{}
	err := mgm.Coll(product).FindByID(userId, product)
	if err != nil {
		return nil, errors.New("cannot find user")
	}

	return product, nil
}

// FindUserById find user by id
func FindOrderById(userId primitive.ObjectID) (*db.Order, error) {
	order := &db.Order{}
	err := mgm.Coll(order).FindByID(userId, order)
	if err != nil {
		return nil, errors.New("cannot find user")
	}

	return order, nil
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
