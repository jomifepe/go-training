package resource

import (
	"github.com/gin-gonic/gin"
	"github.com/jomifepe/gin_api/api/auth"
	"github.com/jomifepe/gin_api/logging"
	"github.com/jomifepe/gin_api/model"
	"net/http"
)

var (
	errAuthUnauthorizedUser    = gin.H{"message": "Unauthorized user, please sign in"}
	errAuthInvalidLoginDetails = gin.H{"message": "Please provide valid login details"}
	errAuthLoginFailed         = gin.H{"message": "Failed to sign in user"}
	errAuthLogoutSuccess       = gin.H{"message": "Successfully logged out"}
)

// authStore is used to define the database calls used by the route group define in this file
type authStore interface {
	RegisterAccess(accessDetails auth.AccessDetails) error
	GetAccess(uuid string) (auth.AccessDetails, error)
	DeleteAccess(accessDetails auth.AccessDetails) error
	GetUserBy(paramName string, param interface{}, omitFields ...string) (model.User, error)
}

// AuthResource holds a AuthStore interface, used to communicate with the database
type AuthResource struct {
	Store authStore
}

// NewAuthResource initializes the AuthResource with an existing AuthStore
func NewAuthResource(store authStore) *AuthResource {
	return &AuthResource{
		Store: store,
	}
}

// MountTaskRoutesTo defines new routes regarding Authentication on an existing gin.RouterGroup or gin.Engine
func (ar *AuthResource) MountAuthRoutesTo(r gin.IRouter, authMiddleware gin.HandlerFunc) (rg *gin.RouterGroup) {
	rg = r.Group("/auth"); {
		rg.POST("/login", ar.handleSignIn)
		rg.POST("/logout", authMiddleware, ar.handleSignOut)
	}
	return
}

// handleSignIn handles user login requests. It validates the email and password passed on the request body,
// checks if it matches and existing user on the database, generates a new access token, registers it on the database
// and returns it to the user
func (ar *AuthResource) handleSignIn(c *gin.Context) {
	var u model.AuthUser

	if err := c.ShouldBindJSON(&u); err != nil {
		logging.Logger.Errorln("[API] Failed to bind user to JSON", err)
		c.JSON(http.StatusUnprocessableEntity, errAuthInvalidLoginDetails)
		return
	}

	dbUser, err := ar.Store.GetUserBy("email", u.Email, "")
	if err != nil {
		logging.Logger.Errorln("[API] Failed to get user by email from the DB", err)
		c.JSON(http.StatusUnauthorized, errAuthInvalidLoginDetails)
		return
	}

	if err = auth.ComparePasswords(u.Password, dbUser.Password); err != nil {
		logging.Logger.Errorln("[API] Received password does not match", u.Password)
		c.JSON(http.StatusUnauthorized, errAuthInvalidLoginDetails)
		return
	}

	tokenDetails, err := auth.GenerateToken(dbUser.ID, dbUser.Email)
	if err != nil {
		logging.Logger.Errorln("[API] Failed to generate token", err)
		c.JSON(http.StatusUnprocessableEntity, errAuthLoginFailed)
		return
	}

	accessDetails := auth.AccessDetails{
		UserID:      dbUser.ID,
		AccessUUID:  tokenDetails.UUID,
		AccessToken: tokenDetails.Token,
	}
	if tErr := ar.Store.RegisterAccess(accessDetails); tErr != nil {
		logging.Logger.Errorln("[API] Failed to store token", tErr)
		c.JSON(http.StatusUnprocessableEntity, errAuthLoginFailed)
		return
	}

	c.JSON(http.StatusOK, tokenDetails)
}

// handleSignOut handles user logout requests. It reads the authorization bearer token passed on the request
// and deletes that access record from the database
func (ar *AuthResource) handleSignOut(c *gin.Context) {
	accessDetails, err := auth.ExtractRequestTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errAuthUnauthorizedUser)
		return
	}
	if err = ar.Store.DeleteAccess(accessDetails); err != nil {
		c.JSON(http.StatusUnauthorized, errAuthUnauthorizedUser)
		return
	}
	c.JSON(http.StatusOK, errAuthLogoutSuccess)
}
