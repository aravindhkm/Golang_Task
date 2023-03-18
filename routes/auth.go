package routes

import (
	"Hdfc_Assignment/controllers"
	"Hdfc_Assignment/middlewares/validators"

	"github.com/gin-gonic/gin"
)

func AuthRoute(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST(
			"/userRegister",
			validators.RegisterValidator(),
			controllers.UserRegister,
		)

		auth.POST(
			"/adminRegister",
			validators.RegisterValidator(),
			controllers.AdminRegister,
		)

		auth.POST(
			"/employeeRegister",
			validators.RegisterValidator(),
			controllers.EmployeeRegister,
		)

		auth.POST(
			"/login",
			validators.LoginValidator(),
			controllers.Login,
		)

		auth.POST(
			"/refresh",
			validators.RefreshValidator(),
			controllers.Refresh,
		)
	}
}
