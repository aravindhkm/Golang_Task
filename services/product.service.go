package services

import (
	model "Hdfc_Assignment/models"
	"Hdfc_Assignment/utils"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/kamva/mgm/v3"
	"github.com/kamva/mgm/v3/builder"
	"github.com/kamva/mgm/v3/field"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	//"golang.org/x/crypto/bcrypt"
)

func getCount() int {
	group := builder.Group("$_id", nil)

	res, err := mgm.Coll(&model.Product{}).Aggregate(mgm.Ctx(), bson.A{builder.S(group)}, nil)
	if err != nil {
		fmt.Println("error 2", err)
	}

	fmt.Println("Raw Result", res.RemainingBatchLength())

	return res.RemainingBatchLength()
}

func GetTokenCount() {
	//gotResult := []interface{}{}
	group := builder.Group("$_id", nil)

	res, err := mgm.Coll(&model.Token{}).Aggregate(mgm.Ctx(), bson.A{builder.S(group)}, nil)
	if err != nil {
		fmt.Println("error 2", err)
	}

	fmt.Println("Raw Result", res.RemainingBatchLength())
	// res.All(mgm.Ctx(),gotResult)
	// fmt.Println("res",gotResult)

}

func Initialize() {

	if getCount() != 0 {
		return
	}

	byteValues, err := os.ReadFile("data.json")
	if err != nil {
		fmt.Println("error 1", err)
	}

	var docs []model.Product
	err = json.Unmarshal(byteValues, &docs)

	if err != nil {
		fmt.Println("error 2", err)
	}

	newData := []interface{}{}

	for j := range docs {
		newData = append(newData, docs[j])
	}

	// fmt.Println("newData", newData)

	res, err := mgm.Coll(&model.Product{}).InsertMany(mgm.Ctx(), newData)

	if err != nil {
		fmt.Println("error 2", err)
	}

	fmt.Println("Res", res)
}

// CreateProduct create new product record
func CreateProduct(
	title string,
	description string,
	price int,
	rating float32,
	stock int,
	brand string,
	productType string,
	category string,
	thumbnail string,
	image []string,
) (*model.Product, error) {

	// password, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	// if err != nil {
	// 	return nil, errors.New("cannot generate hashed password")
	// }

	product := model.NewProduct(title, description, price, rating, stock, brand, productType, category, thumbnail, image)
	err := mgm.Coll(product).Create(product)
	if err != nil {
		return nil, errors.New("cannot create new product")
	}

	return product, nil
}

// GetProducts get paginated product list
func GetProducts(userId primitive.ObjectID, page int, limit int) ([]model.Product, error) {
	var products []model.Product

	findOptions := options.Find().
		SetSkip(int64(page * limit)).
		SetLimit(int64(limit + 1))

	err := mgm.Coll(&model.Product{}).SimpleFind(
		&products,
		bson.M{"author": userId},
		findOptions,
	)

	if err != nil {
		return nil, errors.New("cannot find products")
	}

	return products, nil
}

func GetProductById(userId primitive.ObjectID, productId primitive.ObjectID) (*model.Product, error) {
	product := &model.Product{}
	err := mgm.Coll(product).First(bson.M{field.ID: productId, "author": userId}, product)
	if err != nil {
		return nil, errors.New("cannot find product")
	}

	return product, nil
}

// UpdateProduct updates a product with id
func UpdateProduct(userId primitive.ObjectID, productId primitive.ObjectID, request *utils.ProductRequest) error {
	product := &model.Product{}
	err := mgm.Coll(product).FindByID(productId, product)
	if err != nil {
		return errors.New("cannot find product")
	}

	if product.ID != userId {
		return errors.New("you cannot update this product")
	}

	product.Title = request.Title
	product.Description = request.Description
	err = mgm.Coll(product).Update(product)

	if err != nil {
		return errors.New("cannot update")
	}

	return nil
}

// DeleteProduct delete a product with id
func DeleteProduct(userId primitive.ObjectID, productId primitive.ObjectID) error {
	deleteResult, err := mgm.Coll(&model.Product{}).DeleteOne(mgm.Ctx(), bson.M{field.ID: productId, "author": userId})

	if err != nil || deleteResult.DeletedCount <= 0 {
		return errors.New("cannot delete product")
	}

	return nil
}