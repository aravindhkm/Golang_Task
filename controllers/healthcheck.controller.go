package controllers

import (
	"Hdfc_Assignment/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ping godoc
// @Summary      Ping
// @Description  check server
// @Tags         ping
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.Response
// @Router       /ping [get]
func HealthCheck(c *gin.Context) {
	response := &utils.Response{
		StatusCode: http.StatusOK,
		Success:    true,
		Message:    "OK!!",
	}

	response.SendResponse(c)
}
