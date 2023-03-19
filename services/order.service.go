package services

import (
	db "Hdfc_Assignment/models"
	"Hdfc_Assignment/utils"
	"errors"
	"time"

	"github.com/kamva/mgm/v3"
	"github.com/kamva/mgm/v3/field"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateOrder create new order record
func CreateOrder(
	userId primitive.ObjectID,
	productId []primitive.ObjectID,
	orderQuantity []int,
	totalCost int,
	grandTotal int,
	discountPercentage int,
	discountAmount int,
	address string) (*db.Order, error) {

	order := db.NewOrder(
		userId,
		productId,
		orderQuantity,
		totalCost,
		grandTotal,
		discountPercentage,
		discountAmount,
		address)
	err := mgm.Coll(order).Create(order)
	if err != nil {
		return nil, errors.New("cannot create new order")
	}

	return order, nil
}

// OrderStatus order from the record
func SetOrderStatus(
	orderId primitive.ObjectID,
	empId primitive.ObjectID,
	status string) (*db.Order, error) {
	order := &db.Order{}
	err := mgm.Coll(order).FindByID(orderId, order)
	if err != nil {
		return nil, errors.New("cannot find order")
	}

	order.OrderStatus = status

	if status == "Dispatched" {
		order.DispatchDate = time.Now()
		order.DeliveryEmployee = empId
	}
	err = mgm.Coll(order).Update(order)

	if err != nil {
		return nil, errors.New("cannot update")
	}

	return order, nil
}

// GetOrders get paginated order list
func GetOrders(userId primitive.ObjectID, page int, limit int) ([]db.Order, error) {
	var orders []db.Order

	findOptions := options.Find().
		SetSkip(int64(page * limit)).
		SetLimit(int64(limit + 1))

	err := mgm.Coll(&db.Order{}).SimpleFind(
		&orders,
		bson.M{"author": userId},
		findOptions,
	)

	if err != nil {
		return nil, errors.New("cannot find orders")
	}

	return orders, nil
}

func GetOrderById(userId primitive.ObjectID, orderId primitive.ObjectID) (*db.Order, error) {
	order := &db.Order{}
	err := mgm.Coll(order).First(bson.M{field.ID: orderId, "author": userId}, order)
	if err != nil {
		return nil, errors.New("cannot find order")
	}

	return order, nil
}

// UpdateOrder updates a order with id
func UpdateOrder(
	userId primitive.ObjectID,
	orderId primitive.ObjectID,
	request *utils.OrderRequest) error {
	order := &db.Order{}
	err := mgm.Coll(order).FindByID(orderId, order)
	if err != nil {
		return errors.New("cannot find order")
	}

	if order.ID != userId {
		return errors.New("you cannot update this order")
	}

	order.Address = request.Address
	err = mgm.Coll(order).Update(order)

	if err != nil {
		return errors.New("cannot update")
	}

	return nil
}

// DeleteOrder delete a order with id
func DeleteOrder(userId primitive.ObjectID, orderId primitive.ObjectID) error {
	deleteResult, err := mgm.Coll(&db.Order{}).DeleteOne(mgm.Ctx(), bson.M{field.ID: orderId, "author": userId})

	if err != nil || deleteResult.DeletedCount <= 0 {
		return errors.New("cannot delete order")
	}

	return nil
}
