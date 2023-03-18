package models

import (
	"github.com/kamva/mgm/v3"
)

type Product struct {
	mgm.DefaultModel `bson:",inline"`
	Title            string   `json:"title" bson:"title"`
	Description      string   `json:"description" bson:"description"`
	Price            int      `json:"price" bson:"price"`
	Rating           float32  `json:"rating" bson:"rating"`
	Stock            int      `json:"stock" bson:"stock"`
	Brand            string   `json:"brand" bson:"brand"`
	ProductType      string   `json:"productType" bson:"productType"`
	Category         string   `json:"category" bson:"category"`
	Thumbnail        string   `json:"thumbnail" bson:"thumbnail"`
	Images           []string `json:"images" bson:"images"`
}

func NewProduct(
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
) *Product {
	return &Product{
		Title:       title,
		Description: description,
		Price:       price,
		Rating:      rating,
		Stock:       stock,
		Brand:       brand,
		ProductType:        productType,
		Category:    category,
		Thumbnail:   thumbnail,
		Images:      image,
	}
}

func (model *Product) CollectionDescription() string {
	return "product"
}
