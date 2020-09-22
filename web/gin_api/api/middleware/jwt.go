package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/jomifepe/gin_api/api/auth"
	"net/http"
)

var (
	unauthorizedMessage = gin.H{"message": "Unauthorized user, please sign in"}
)

type authStore interface {
	GetAccess(uuid string) (auth.AccessDetails, error)
}

type AuthMiddleware struct {
	Store authStore
}

func NewAuthMiddleware(store authStore) *AuthMiddleware {
	return &AuthMiddleware{
		Store: store,
	}
}

// AuthenticateToken is an authentication middleware for gin that extracts an authorization token from the request
// header, parses and validates it, and checks if its UUID exists on the database. If it doesn't, aborts the request
// with a http.StatusUnauthorized status code.
func (am *AuthMiddleware) AuthenticateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		ad, err := auth.ExtractRequestTokenMetadata(c.Request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, unauthorizedMessage)
			return
		}
		err = auth.ValidateToken(ad.AccessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, unauthorizedMessage)
			return
		}
		_, err = am.Store.GetAccess(ad.AccessUUID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, unauthorizedMessage)
			return
		}

		c.Next()
	}
}