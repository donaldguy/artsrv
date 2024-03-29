package requests

import (
	"net/http"

	"github.com/donaldguy/artsrv/art"

	"github.com/gin-gonic/gin"
)

//DownloadHandler expects a URL param id and returns a plain-text response containing the correspoding art
func DownloadHandler(c *gin.Context) {
	id := c.Param("id")
	resBody, err := art.Get(id)

	switch err {
	case nil:
		c.String(http.StatusOK, resBody)
	case art.ErrDoesNotExist:
		jsonResponseFromCodeAndMessage(c, http.StatusNotFound, err.Error())
	default:
		jsonResponseFromCodeAndMessage(c, http.StatusInternalServerError, err.Error())
	}
}
