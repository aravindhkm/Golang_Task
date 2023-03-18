package routes

import (
	"Hdfc_Assignment/docs"
	"Hdfc_Assignment/middlewares"
	"Hdfc_Assignment/services"
	"Hdfc_Assignment/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func New() *gin.Engine {
	r := gin.New()
	initRoute(r)

	r.Use(gin.LoggerWithWriter(middlewares.LogWriter()))
	r.Use(gin.CustomRecovery(middlewares.AppRecovery()))
	r.Use(middlewares.CORSMiddleware())

	v1 := r.Group("/api/v1")
	{
		HealthCheckRoute(v1)
		AuthRoute(v1)
		NoteRoute(v1, middlewares.JWTMiddleware())
		EmployeeRoute(v1, middlewares.JWTMiddleware())
	}

	docs.SwaggerInfo.BasePath = v1.BasePath() // adds /v1 to swagger base path

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return r
}

func initRoute(r *gin.Engine) {
	_ = r.SetTrustedProxies(nil)
	r.RedirectTrailingSlash = false
	r.HandleMethodNotAllowed = true

	r.NoRoute(func(c *gin.Context) {
		utils.SendErrorResponse(c, http.StatusNotFound, c.Request.RequestURI+" not found")
	})

	r.NoMethod(func(c *gin.Context) {
		utils.SendErrorResponse(c, http.StatusMethodNotAllowed, c.Request.Method+" is not allowed here")
	})
}

func InitGin() {
	gin.DisableConsoleColor()
	gin.SetMode(services.Config.Mode)
	// do some other initialization staff
}
