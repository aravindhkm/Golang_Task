package services

import (
	model "Hdfc_Assignment/models"
	"Hdfc_Assignment/utils"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/kamva/mgm/v3"
	"github.com/kamva/mgm/v3/builder"
	"github.com/kamva/mgm/v3/field"
	"github.com/kamva/mgm/v3/operator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getCount() int {
	group := builder.Group("$_id", nil)

	res, err := mgm.Coll(&model.Product{}).Aggregate(mgm.Ctx(), bson.A{builder.S(group)}, nil)
	if err != nil {
		panic(err)
	}

	return res.RemainingBatchLength()
}

func GetTokenCount() int32 {
	var gotResult = []map[string]interface{}{}

	count := bson.M{operator.Count: "data"}
	found, err := mgm.Coll(&model.Token{}).Aggregate(mgm.Ctx(), bson.A{count}, nil)
	if err != nil {
		panic(err)
	}
	found.All(mgm.Ctx(), &gotResult)

	return gotResult[0]["data"].(int32)
}

func GetProductCount() int32 {
	var gotResult = []map[string]interface{}{}

	count := bson.M{operator.Count: "data"}
	found, err := mgm.Coll(&model.Product{}).Aggregate(mgm.Ctx(), bson.A{count}, nil)
	if err != nil {
		panic(err)
	}
	found.All(mgm.Ctx(), &gotResult)

	if len(gotResult) == 0 {
		return 0
	} else {
		return gotResult[0]["data"].(int32)
	}

}

func Initialize() {
	if GetProductCount() != 0 {
		return
	}

	byteValues, err := os.ReadFile("data.json")
	if err != nil {
		panic(err)
	}

	var docs []model.Product
	err = json.Unmarshal(byteValues, &docs)

	if err != nil {
		panic(err)
	}

	newData := []interface{}{}
	for j := range docs {
		newData = append(newData, docs[j])
	}

	res, err := mgm.Coll(&model.Product{}).InsertMany(mgm.Ctx(), newData)

	if err != nil {
		panic(err)
	}

	log.Fatalln("InsertProduct", res)
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

// GetProducts get paginated product list
func GetAllProducts(page int, limit int) ([]model.Product, error) {
	var products []model.Product

	findOptions := options.Find().
		SetSkip(int64(page * limit)).
		SetLimit(int64(limit + 1))

	err := mgm.Coll(&model.Product{}).SimpleFind(
		&products,
		bson.M{},
		findOptions,
	)

	if err != nil {
		return nil, errors.New("cannot find products")
	}

	return products, nil
}

func GetProductById(productId primitive.ObjectID) (*model.Product, error) {
	product := &model.Product{}
	err := mgm.Coll(product).First(bson.M{field.ID: productId}, product)
	if err != nil {
		return nil, errors.New("cannot find product")
	}

	return product, nil
}

// UpdateProduct updates a product with id
func UpdateProduct(productId primitive.ObjectID, request *utils.ProductRequest) error {
	product := &model.Product{}
	err := mgm.Coll(product).FindByID(productId, product)
	if err != nil {
		return errors.New("cannot find product")
	}

	product.Title = request.Title
	product.Description = request.Description
	err = mgm.Coll(product).Update(product)

	if err != nil {
		return errors.New("cannot update")
	}

	return nil
}

func UpdateProductCancelOrder(productId primitive.ObjectID, orderedStock int) error {
	product := &model.Product{}
	err := mgm.Coll(product).FindByID(productId, product)
	if err != nil {
		return errors.New("cannot find product")
	}

	product.Stock += orderedStock
	err = mgm.Coll(product).Update(product)

	if err != nil {
		return errors.New("cannot update")
	}

	return nil
}

func UpdateProductStock(productId primitive.ObjectID, orderedStock int) error {
	product := &model.Product{}
	err := mgm.Coll(product).FindByID(productId, product)
	if err != nil {
		return errors.New("cannot find product")
	}

	if product.Stock < orderedStock {
		return errors.New("invalid quantity")
	}

	fmt.Println("orderedStock", orderedStock, product.Stock)

	product.Stock -= orderedStock
	err = mgm.Coll(product).Update(product)

	if err != nil {
		return errors.New("cannot update")
	}

	return nil
}

func UpdateMultipleProductStock(productId []primitive.ObjectID, orderedStock []int) error {
	product := &model.Product{}

	for i := 0; i < len(productId); i++ {
		err := mgm.TransactionWithCtx(mgm.Ctx(), func(session mongo.Session, sc mongo.SessionContext) error {

			err := mgm.Coll(product).FindByIDWithCtx(sc, productId[i], product)
			if err != nil {
				session.AbortTransaction(sc)
				return errors.New("cannot find product")
			}

			if product.Stock < orderedStock[i] {
				session.AbortTransaction(sc)
				return errors.New("invalid quantity")
			}

			product.Stock -= orderedStock[i]
			err = mgm.Coll(product).UpdateWithCtx(sc, product)

			if err != nil {
				session.AbortTransaction(sc)
				return errors.New("cannot update")
			}

			err = session.CommitTransaction(mgm.Ctx())
			if err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			return errors.New("cannot update")
		}

	}

	return nil
}

// DeleteProduct delete a product with id
func DeleteProduct(productId primitive.ObjectID) error {
	deleteResult, err := mgm.Coll(&model.Product{}).DeleteOne(mgm.Ctx(), bson.M{field.ID: productId})

	if err != nil || deleteResult.DeletedCount <= 0 {
		return errors.New("cannot delete product")
	}

	return nil
}
