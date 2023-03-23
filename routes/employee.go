package routes

import (
	"Hdfc_Assignment/controllers"
	"Hdfc_Assignment/middlewares/validators"

	"github.com/gin-gonic/gin"
)

func EmployeeRoute(router *gin.RouterGroup, handlers ...gin.HandlerFunc) {
	emp := router.Group("/employee", handlers...)
	{
		emp.POST(
			"/register",
			validators.CreateEmployeeValidator(),
			controllers.CreateNewEmployee,
		)

		emp.POST(
			"/dispatchOrder/:id",
			validators.PathIdValidator(),
			controllers.DispatchOrder,
		)

		emp.GET(
			"/getAllDetails",
			validators.GetEmployeesValidator(),
			controllers.GetAllEmployees,
		)

		emp.GET(
			"",
			validators.GetEmployeesValidator(),
			controllers.GetEmployees,
		)

		emp.GET(
			"/:id",
			validators.PathIdValidator(),
			controllers.GetOneEmployee,
		)

		emp.PUT(
			"/update/:id",
			validators.PathIdValidator(),
			validators.UpdateEmployeeValidator(),
			controllers.UpdateEmployee,
		)

		emp.DELETE(
			"/delete/:id",
			validators.PathIdValidator(),
			controllers.DeleteEmployee,
		)
	}
}
