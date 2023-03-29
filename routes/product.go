package routes

import (
	"Hdfc_Assignment/controllers"
	"Hdfc_Assignment/middlewares/validators"

	"github.com/gin-gonic/gin"
)

func ProductRoute(router *gin.RouterGroup, handlers ...gin.HandlerFunc) {
	emp := router.Group("/product", handlers...)
	{
		emp.POST(
			"/create",
			validators.CreateProductValidator(),
			controllers.CreateNewProduct,
		)

		emp.GET(
			"/getAllProduct",
			validators.GetProductsValidator(),
			controllers.GetAllProducts,
		)

		emp.GET(
			"/:id",
			validators.PathIdValidator(),
			controllers.GetOneProduct,
		)

		emp.PUT(
			"/update/:id",
			validators.PathIdValidator(),
			validators.UpdateProductValidator(),
			controllers.UpdateProduct,
		)

		emp.DELETE(
			"/delete/:id",
			validators.PathIdValidator(),
			controllers.DeleteProduct,
		)
	}
}
