package requests

import (
	"net/http"
	"strconv"

	"github.com/donaldguy/artsrv/art"

	"github.com/gin-gonic/gin"
)

type UploadChunkBody struct {
	ID   int    `json:"id"  binding:"exists"` // See: https://github.com/gin-gonic/gin/issues/491
	Size int    `json:"size"`
	Data string `json:"data" binding:"required"`
}

func UploadChunkHandler(c *gin.Context) {
	id := c.Param("id")
	var reqBody UploadChunkBody

	if c.ContentType() != "application/json" {
		jsonResponseFromCodeAndMessage(c, http.StatusUnsupportedMediaType, "Registration must be json")
		return
	}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		jsonResponseFromCodeAndMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := art.SubmitChunk(id, strconv.Itoa(reqBody.ID), reqBody.Data); err != nil {
		switch err {
		case art.ErrNotRegistered:
			jsonResponseFromCodeAndMessage(c, http.StatusNotFound, err.Error())
		case art.ErrChunkAlreadySubmitted:
			jsonResponseFromCodeAndMessage(c, http.StatusConflict, err.Error())
		default:
			jsonResponseFromCodeAndMessage(c, http.StatusInternalServerError, err.Error())
		}

		return
	}

	jsonResponseFromCodeAndMessage(c, http.StatusCreated, "Chunk Succesfully submitted")
}
