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

		emp.GET(
			"",
			validators.GetNotesValidator(),
			controllers.GetNotes,
		)

		emp.GET(
			"/:id",
			validators.PathIdValidator(),
			controllers.GetOneNote,
		)

		emp.PUT(
			"/:id",
			validators.PathIdValidator(),
			validators.UpdateNoteValidator(),
			controllers.UpdateNote,
		)

		emp.DELETE(
			"/:id",
			validators.PathIdValidator(),
			controllers.DeleteNote,
		)
	}
}
