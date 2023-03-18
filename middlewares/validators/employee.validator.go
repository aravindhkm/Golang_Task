package validators

import (
	"Hdfc_Assignment/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

func CreateEmployeeValidator() gin.HandlerFunc {
	return func(c *gin.Context) {

		var createEmployeeRequest utils.EmployeeRequest
		_ = c.ShouldBindBodyWith(&createEmployeeRequest, binding.JSON)

		if err := createEmployeeRequest.Validate(); err != nil {
			utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		c.Next()
	}
}

func GetEmployeesValidator() gin.HandlerFunc {
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

func UpdateEmployeeValidator() gin.HandlerFunc {
	return func(c *gin.Context) {

		var employeeRequest utils.EmployeeRequest
		_ = c.ShouldBindBodyWith(&employeeRequest, binding.JSON)

		if err := employeeRequest.Validate(); err != nil {
			utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		c.Next()
	}
}
