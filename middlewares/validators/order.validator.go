package validators

import (
	"Hdfc_Assignment/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

func CreateOrderValidator() gin.HandlerFunc {
	return func(c *gin.Context) {

		var createOrderRequest utils.OrderRequest
		_ = c.ShouldBindBodyWith(&createOrderRequest, binding.JSON)

		if err := createOrderRequest.Validate(); err != nil {
			utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		c.Next()
	}
}

func GetOrdersValidator() gin.HandlerFunc {
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

func UpdateOrderValidator() gin.HandlerFunc {
	return func(c *gin.Context) {

		var orderRequest utils.OrderRequest
		_ = c.ShouldBindBodyWith(&orderRequest, binding.JSON)

		if err := orderRequest.Validate(); err != nil {
			utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		c.Next()
	}
}
