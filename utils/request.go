package utils

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderRequest struct {
	ProductId     []primitive.ObjectID `json:"productId"`
	OrderQuantity []int                `json:"orderQuantity"`
	Address       string               `json:"address"`
}

func (a OrderRequest) Validate() error {
	return validation.ValidateStruct(&a,
		// validation.Field(&a.UserId, validation.Required, is.MongoID),
		// validation.Field(&a.ProductId, validation.Required, validation.Each(is.MongoID)),
		// validation.Field(&a.UserId, validation.Required),
		validation.Field(&a.ProductId, validation.Required),
		validation.Field(&a.OrderQuantity, validation.Required),
		validation.Field(&a.Address, validation.Required, validation.Length(5, 150)),
	)
}

var passwordRule = []validation.Rule{
	validation.Required,
	validation.Length(8, 32),
	validation.Match(regexp.MustCompile("^\\S+$")).Error("cannot contain whitespaces"),
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a RegisterRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Name, validation.Required, validation.Length(3, 64)),
		validation.Field(&a.Email, validation.Required, is.Email),
		validation.Field(&a.Password, passwordRule...),
	)
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (a LoginRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Email, validation.Required, is.Email),
		validation.Field(&a.Password, passwordRule...),
		validation.Field(&a.Role, validation.Required),
	)
}

type RefreshRequest struct {
	Token string `json:"token"`
	Role  string `json:"role"`
}

func (a RefreshRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(
			&a.Token,
			validation.Required,
			validation.Match(regexp.MustCompile("^\\S+$")).Error("cannot contain whitespaces"),
		),
		validation.Field(&a.Role, validation.Required),
	)
}

type NoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (a NoteRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Title, validation.Required),
		validation.Field(&a.Content, validation.Required),
	)
}

type EmployeeRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Mobile   uint   `json:"mobile"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Address  string `json:"address"`
}

func (a EmployeeRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Name, validation.Required, validation.Length(3, 64)),
		validation.Field(&a.Email, validation.Required, is.Email),
		validation.Field(&a.Password, passwordRule...),
		validation.Field(&a.Mobile, validation.Required),
		validation.Field(&a.Address, validation.Required),
	)
}

type UpdateEmployeeRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Mobile  uint   `json:"mobile"`
	Address string `json:"address"`
}

func (a UpdateEmployeeRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Name, validation.Required, validation.Length(3, 64)),
		validation.Field(&a.Email, validation.Required, is.Email),
		validation.Field(&a.Mobile, validation.Required),
		validation.Field(&a.Address, validation.Required),
	)
}

type ProductRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Price       int      `json:"price"`
	Rating      float32  `json:"rating"`
	Stock       int      `json:"stock"`
	Brand       string   `json:"brand"`
	Type        string   `json:"productType"`
	Category    string   `json:"category"`
	Thumbnail   string   `json:"thumbnail"`
	Images      []string `json:"images"`
}

func (a ProductRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Title, validation.Required, validation.Length(3, 64)),
		validation.Field(&a.Description, validation.Required),
		validation.Field(&a.Price, validation.Required),
		validation.Field(&a.Stock, validation.Required),
		validation.Field(&a.Brand, validation.Required),
		validation.Field(&a.Type, validation.Required),
		validation.Field(&a.Category, validation.Required),
		validation.Field(&a.Thumbnail, validation.Required),
		validation.Field(&a.Images, validation.Required),
	)
}
