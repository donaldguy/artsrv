package requests

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

//jsonResponseFromCodeAndMessage takes an HTTP status code and a message and sends an HTTP response
//both with that code as status header and with the status as a field of the JSON response
func jsonResponseFromCodeAndMessage(c *gin.Context, code int, msg string) {
	c.JSON(code, gin.H{
		"code":    fmt.Sprintf("%d %s", code, http.StatusText(code)),
		"message": msg,
	})
}
