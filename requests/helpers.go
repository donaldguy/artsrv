package requests

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func jsonResponseFromCodeAndMessage(c *gin.Context, code int, msg string) {
	c.JSON(code, gin.H{
		"code":    fmt.Sprintf("%d %s", code, http.StatusText(code)),
		"message": msg,
	})
}
