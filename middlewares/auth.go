package middlewares

import (
	db "Hdfc_Assignment/models"
	"Hdfc_Assignment/services"
	"Hdfc_Assignment/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		tokenModel, err := services.VerifyToken(token, db.TokenTypeAccess)
		if err != nil {
			utils.SendErrorResponse(c, http.StatusUnauthorized, err.Error())
			return
		}

		c.Set("userIdHex", tokenModel.User.Hex())
		c.Set("userId", tokenModel.User)

		c.Next()
	}
}
