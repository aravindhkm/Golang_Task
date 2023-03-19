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
