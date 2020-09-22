package resource

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jomifepe/gin_api/api/auth"
	"github.com/jomifepe/gin_api/api/middleware"
	"github.com/jomifepe/gin_api/logging"
	"github.com/jomifepe/gin_api/model"
	"net/http"
)

var (
	errUserCreateInvalidFields = gin.H{"message": "The specified user has invalid fields"}
	errUserCreateGeneric       = gin.H{"message": "Couldn't create user"}
)

// userStore is used to define the database calls used by the route group define in this file
type userStore interface {
	GetAllUsers(omitFields ...string) ([]model.User, error)
	GetUserBy(paramName string, param interface{}, omitFields ...string) (model.User, error)
	CreateUser(u model.User) (model.User, error)
	UpdateUser(u model.User) (model.User, error)
	DeleteUser(id int) error
}

// UserResource holds a TaskStore interface, used to communicate with the database
type UserResource struct {
	Store userStore
}

// NewUserResource initializes the UserResource with an existing UserStore
func NewUserResource(store userStore) *UserResource {
	return &UserResource{
		Store: store,
	}
}

// MountUserRoutesTo defines new routes regarding Users on an existing gin.RouterGroup or gin.Engine
func (ur *UserResource) MountUserRoutesTo(r gin.IRouter) {
	idParam := middleware.Param{Key: "id", ExampleValue: -1}

	r.GET("/me", ur.handleMe)
	rg := r.Group("/users"); {
		rg.GET("", ur.handleGetUsers)
		rg.POST("", ur.handleCreateUser)
		withId := rg.Group("", middleware.ExtractParam(idParam)); {
			withId.GET("/:id", ur.handleGetUser)
			withId.PUT("/:id", ur.handleUpdateUser)
		}
	}
}

// handleMe returns the current user information by extracting its id from the access token used on
// the request and querying the database
func (ur *UserResource) handleMe(c *gin.Context) {
	metadata, err := auth.ExtractRequestTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": err,
		})
		return
	}

	user, err := ur.Store.GetUserBy("id", metadata.UserID, "password")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Invalid token",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

// handleGetUsers returns all the existing users
func (ur *UserResource) handleGetUsers(c *gin.Context) {
	users, err := ur.Store.GetAllUsers()
	if err != nil {
		logging.Logger.Errorln("[API] Failed to get all users", err)
	}

	c.JSON(http.StatusOK, users)
}

func (ur *UserResource) handleGetUser(c *gin.Context) {
	id := c.GetInt("id")

	u, err := ur.Store.GetUserBy("id", id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("No user with the id %v was found", id),
		})
		return
	}

	c.JSON(http.StatusOK, u)
}

// handleCreateUser validates the fields specified on the request body, and inserts a new user
// on the database if these are valid
func (ur *UserResource) handleCreateUser(c *gin.Context) {
	var u model.User

	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, errUserCreateInvalidFields)
		return
	}

	hash, err := auth.GeneratePassword(u.Password)
	if err != nil {
		logging.Logger.Errorln("[API] Failed to generate hash from password", err)
		c.JSON(http.StatusBadRequest, errUserCreateGeneric)
		return
	}
	u.Password = hash

	newUser, err := ur.Store.CreateUser(u)
	if err != nil {
		c.JSON(http.StatusBadRequest, errUserCreateGeneric)
		return
	}

	newUser.Password = ""
	c.JSON(http.StatusCreated, newUser)
}

func (ur *UserResource) handleUpdateUser(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"message": "TODO",
	})
}
