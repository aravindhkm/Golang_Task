package validators

import (
	"Hdfc_Assignment/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

func CreateNoteValidator() gin.HandlerFunc {
	return func(c *gin.Context) {

		var createNoteRequest utils.NoteRequest
		_ = c.ShouldBindBodyWith(&createNoteRequest, binding.JSON)

		if err := createNoteRequest.Validate(); err != nil {
			utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		c.Next()
	}
}

func GetNotesValidator() gin.HandlerFunc {
	return func(c *gin.Context) {

		page := c.DefaultQuery("page", "0")
		err := validation.Validate(page, is.Int)
		if err != nil {
			utils.SendErrorResponse(c, http.StatusBadRequest, "invalid page: "+page)
			return
		}

		c.Next()
	}
}

func UpdateNoteValidator() gin.HandlerFunc {
	return func(c *gin.Context) {

		var noteRequest utils.NoteRequest
		_ = c.ShouldBindBodyWith(&noteRequest, binding.JSON)

		if err := noteRequest.Validate(); err != nil {
			utils.SendErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		c.Next()
	}
}
