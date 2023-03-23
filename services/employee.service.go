package services

import (
	db "Hdfc_Assignment/models"
	"Hdfc_Assignment/utils"
	"errors"
	//"fmt"

	"github.com/kamva/mgm/v3"
	"github.com/kamva/mgm/v3/field"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// CreateEmployee create new employee record
func CreateEmployee(
	name string,
	email string,
	plainPassword string,
	mobile uint,
	address string) (*db.Employee, error) {

	password, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("cannot generate hashed password")
	}

	employee := db.NewEmployee(email, string(password), name, mobile, address)
	err = mgm.Coll(employee).Create(employee)
	if err != nil {
		return nil, errors.New("cannot create new employee")
	}

	return employee, nil
}

// GetEmployees get paginated employee list
func GetEmployees(userId primitive.ObjectID, page int, limit int) ([]db.Employee, error) {
	var employees []db.Employee

	findOptions := options.Find().
		SetSkip(int64(page * limit)).
		SetLimit(int64(limit + 1))

	err := mgm.Coll(&db.Employee{}).SimpleFind(
		&employees,
		bson.M{"_id": userId},
		findOptions,
	)

	if err != nil {
		return nil, errors.New("cannot find employees")
	}

	return employees, nil
}

// GetEmployees get paginated employee list
func GetAllEmployees(page int, limit int) ([]db.Employee, error) {
	var employees []db.Employee

	findOptions := options.Find().
		SetSkip(int64(page * limit)).
		SetLimit(int64(limit + 1))

	err := mgm.Coll(&db.Employee{}).SimpleFind(
		&employees,
		bson.M{},
		findOptions,
	)

	if err != nil {
		return nil, errors.New("cannot find employees")
	}

	return employees, nil
}

func GetEmployeeById(employeeId primitive.ObjectID) (*db.Employee, error) {
	employee := &db.Employee{}
	err := mgm.Coll(employee).First(bson.M{field.ID: employeeId}, employee)
	if err != nil {
		return nil, errors.New("cannot find employee")
	}

	return employee, nil
}

func GetEmployeeByMail(mail string) (*db.Employee, error) {
	employee := &db.Employee{}
	err := mgm.Coll(employee).First(bson.M{mail: mail}, employee)
	if err != nil {
		return nil, errors.New("cannot find employee")
	}

	return employee, nil
}

// UpdateEmployee updates a employee with id
func UpdateEmployee(employeeId primitive.ObjectID, request *utils.EmployeeRequest) error {
	employee := &db.Employee{}
	err := mgm.Coll(employee).FindByID(employeeId, employee)
	if err != nil {
		return errors.New("cannot find employee")
	}

	employee.Name = request.Name
	employee.Address = request.Address
	employee.Mobile = request.Mobile
	err = mgm.Coll(employee).Update(employee)

	if err != nil {
		return errors.New("cannot update")
	}

	return nil
}

// DeleteEmployee delete a employee with id
func DeleteEmployee(employeeId primitive.ObjectID) error {
	_, err := mgm.Coll(&db.Employee{}).DeleteOne(mgm.Ctx(), bson.M{field.ID: employeeId})

	if err != nil {
		return errors.New("cannot delete employee")
	}

	return nil
}
