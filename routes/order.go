package routes

import (
	"Hdfc_Assignment/controllers"
	"Hdfc_Assignment/middlewares/validators"

	"github.com/gin-gonic/gin"
)

func OrderRoute(router *gin.RouterGroup, handlers ...gin.HandlerFunc) {
	emp := router.Group("/order", handlers...)
	{
		emp.POST(
			"/placeOrder",
			validators.CreateOrderValidator(),
			controllers.CreateNewOrder,
		)

		emp.POST(
			"/cancelOrder/:id",
			validators.PathIdValidator(),
			controllers.CancelOrder,
		)

		emp.POST(
			"/completeOrder/:id",
			validators.PathIdValidator(),
			controllers.CompleteOrder,
		)

		emp.GET(
			"/getDetails",
			validators.GetOrdersValidator(),
			controllers.GetOrders,
		)

		emp.GET(
			"/getDetails/:id",
			validators.PathIdValidator(),
			controllers.GetOneOrder,
		)
	}
}
