package impl

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func errorHandlingMiddleware(c *gin.Context) {
	c.Next() // Process the request

	// Check for errors added to the context
	if len(c.Errors) > 0 && c.Writer.Size() == 0 {
		errorMessages := make([]string, len(c.Errors))

		for i := range c.Errors {
			err := c.Errors[i]
			errorMessages[i] = fmt.Sprintf("Error %d: %s", i, err.Error())
		}

		errorMessage := strings.Join(errorMessages, "\n")
		_, _ = c.Writer.WriteString(errorMessage)
		return
	}
}
