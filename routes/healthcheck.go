package routes

import (
	"Hdfc_Assignment/controllers"

	"github.com/gin-gonic/gin"
)

func HealthCheckRoute(router *gin.RouterGroup) {
	auth := router.Group("/healthcheck")
	{
		auth.GET(
			"",
			controllers.HealthCheck,
		)
	}
}
