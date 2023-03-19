package models

import (
	"time"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	mgm.DefaultModel   `bson:",inline"`
	UserId             primitive.ObjectID   `json:"userId" bson:"userId"`
	DeliveryEmployee   primitive.ObjectID   `json:"deliveryEmployee" bson:"deliveryEmployee"`
	ProductId          []primitive.ObjectID `json:"productId" bson:"productId"`
	OrderQuantity      []int                `json:"orderQuantity," bson:"orderQuantity,"`
	TotalCost          int                  `json:"totalCost," bson:"totalCost,"`
	GrandTotal         int                  `json:"grandTotal," bson:"grandTotal,"`
	DiscountPercentage int                  `json:"discountPercentage," bson:"discountPercentage,"`
	DiscountAmount     int                  `json:"discountAmount," bson:"discountAmount,"`
	Address            string               `json:"address" bson:"address"`
	OrderStatus        string               `json:"orderStatus" bson:"orderStatus"`
	DispatchDate       time.Time            `json:"dispatchedDate" bson:"dispatchedDate"`
}

func NewOrder(
	userId primitive.ObjectID,
	productId []primitive.ObjectID,
	orderQuantity []int,
	totalCost int,
	grandTotal int,
	discountPercentage int,
	discountAmount int,
	address string,
) *Order {
	var currTime time.Time;
	return &Order{
		UserId:             userId,
		ProductId:          productId,
		OrderQuantity:      orderQuantity,
		TotalCost:          totalCost,
		GrandTotal:         grandTotal,
		DiscountPercentage: discountPercentage,
		DiscountAmount:     discountAmount,
		Address:            address,
		OrderStatus:        "Placed",
		DispatchDate:       currTime,
	}
}

func (model *Order) CollectionDescription() string {
	return "order"
}
