package validators

import (
	"Hdfc_Assignment/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

func CreateProductValidator() gin.HandlerFunc {
	return func(c *gin.Context) {

		var createProductRequest utils.ProductRequest
		_ = c.ShouldBindBodyWith(&createProductRequest, binding.JSON)

		if err := createProductRequest.Validate(); err != nil {
			utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		c.Next()
	}
}

func GetProductsValidator() gin.HandlerFunc {
	return func(c *gin.Context) {

		page := c.DefaultQuery("page", "0")
		err := validation.Validate(page, is.Int)
		if err != nil {
			utils.SendErrorResponse(c, http.StatusBadRequest, "invalid page: "+page)
			return
		}

		c.Next()
	}
}

func UpdateProductValidator() gin.HandlerFunc {
	return func(c *gin.Context) {

		var productRequest utils.ProductRequest
		_ = c.ShouldBindBodyWith(&productRequest, binding.JSON)

		if err := productRequest.Validate(); err != nil {
			utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		c.Next()
	}
}
