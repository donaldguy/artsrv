package requests

import (
	"net/http"

	"github.com/donaldguy/artsrv/art"

	"github.com/gin-gonic/gin"
)

//RegisterBody represents the expected JSON body of a registration request
type RegisterBody struct {
	Sha256    string `json:"sha256" binding:"required"`
	Size      int    `json:"size" binding:"required"`
	ChunkSize int    `json:"chunk_size"`
}

//RegisterHandler handles the registration of new art for subsequent chunked upload
func RegisterHandler(c *gin.Context) {
	var reqBody RegisterBody

	if c.ContentType() != "application/json" {
		jsonResponseFromCodeAndMessage(c, http.StatusUnsupportedMediaType, "Registration must be json")
		return
	}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		jsonResponseFromCodeAndMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := art.Register(reqBody.Sha256, reqBody.Size); err != nil {
		switch err {
		case art.ErrAlreadyRegistered:
			fallthrough
		case art.ErrAlreadyExists:
			jsonResponseFromCodeAndMessage(c, http.StatusConflict, err.Error())
		default:
			jsonResponseFromCodeAndMessage(c, http.StatusInternalServerError, err.Error())
		}

		return
	}

	jsonResponseFromCodeAndMessage(c, http.StatusCreated, "Registered")
}
