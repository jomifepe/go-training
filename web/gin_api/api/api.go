package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jomifepe/gin_api/api/middleware"
	routes "github.com/jomifepe/gin_api/api/resource"
	"github.com/jomifepe/gin_api/logging"
	"github.com/jomifepe/gin_api/storage"
	"github.com/jomifepe/gin_api/util"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

// Start initializes the required resources, defines the API routes and starts listening for HTTP requests on <port>.
func Start(port string) {
	logging.Logger.WithFields(logrus.Fields{
		"port": port,
	}).Infoln("[API] Starting...")

	dbConn := storage.ConfigurePostgresDB()

	authStore := storage.NewAuthStore(dbConn)
	taskStore := storage.NewTaskStore(dbConn)
	userStore := storage.NewUserStore(dbConn)

	authResource := routes.NewAuthResource(authStore)
	taskResource := routes.NewTaskResource(taskStore)
	userResource := routes.NewUserResource(userStore)

	gin.SetMode(gin.ReleaseMode)
	ginEngine := gin.New()
	ginEngine.Use(middleware.Logger(logging.Logger), gin.Recovery())

	authMiddleware := middleware.NewAuthMiddleware(authStore)

	authResource.MountAuthRoutesTo(ginEngine, authMiddleware.AuthenticateToken())
	authGroup := ginEngine.Group("", authMiddleware.AuthenticateToken()); {
		taskResource.MountTaskRoutesTo(authGroup)
		userResource.MountUserRoutesTo(authGroup)
	}

	var err error
	if util.FileExists("./cert.pem") && util.FileExists("./key.pem") {
		err = ginEngine.RunTLS(":" + port, "cert.pem", "key.pem")
	} else {
		err = ginEngine.Run(":" + port)
	}
	if err != nil {
		logging.Logger.WithFields(util.OmitEmptyFields(logrus.Fields{
			"error": err,
			"port": port,
		})).Panicln("[API] Failed to start")
	}
}