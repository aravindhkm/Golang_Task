package middlewares

import (
	"Hdfc_Assignment/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AppRecovery() func(c *gin.Context, recovered interface{}) {
	return func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			utils.SendErrorResponse(c, http.StatusInternalServerError, err)
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"success": false}) // recovery failed
	}
}
