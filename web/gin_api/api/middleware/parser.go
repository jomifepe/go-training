package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var (
	errInvalidParam = func(paramName string) map[string]interface{} {
		return gin.H{"message": fmt.Sprintf("Invalid %v specified", paramName)}
	}
)

type Param struct {
	Key          string
	ExampleValue interface{}
	IsQuery      bool
}

// ExtractParam is a middleware for gin that parses parameters passed via request url path.
// It receives n Param arguments, in order to identify the name of the parameter to extract
// and also the type to parse to (using the example value).
// The parsed parameters are stored on the gin.Context for further usage.
func ExtractParam(params ...Param) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, param := range params {
			var (
				val       string
				parsedVal interface{}
				err       error
			)

			if param.IsQuery {
				val = c.Query(param.Key)
			} else {
				val = c.Param(param.Key)
			}

			if len(val) == 0 {
				c.AbortWithStatusJSON(http.StatusUnprocessableEntity, errInvalidParam(param.Key))
				return
			}

			switch param.ExampleValue.(type) {
			case int:
				parsedVal, err = strconv.Atoi(val)
			case bool:
				parsedVal, err = strconv.ParseBool(val)
			}
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnprocessableEntity, errInvalidParam(param.Key))
				return
			}
			c.Set(param.Key, parsedVal)
		}

		c.Next()
	}
}