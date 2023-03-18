package validators

import (
	"Hdfc_Assignment/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func RegisterValidator() gin.HandlerFunc {
	return func(c *gin.Context) {

		var registerRequest utils.RegisterRequest
		_ = c.ShouldBindBodyWith(&registerRequest, binding.JSON)

		if err := registerRequest.Validate(); err != nil {
			utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		c.Next()
	}
}

func LoginValidator() gin.HandlerFunc {
	return func(c *gin.Context) {

		var loginRequest utils.LoginRequest
		_ = c.ShouldBindBodyWith(&loginRequest, binding.JSON)

		if err := loginRequest.Validate(); err != nil {
			utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		c.Next()
	}
}

func RefreshValidator() gin.HandlerFunc {
	return func(c *gin.Context) {

		var refreshRequest utils.RefreshRequest
		_ = c.ShouldBindBodyWith(&refreshRequest, binding.JSON)

		if err := refreshRequest.Validate(); err != nil {
			utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		c.Next()
	}
}
