package entrypoint

import (
	"base_project/entrypoint/auth"
	"base_project/entrypoint/comment"
	"base_project/entrypoint/post"
	"base_project/repository"
	"base_project/shared/external/database"
	"base_project/shared/fhttp/middleware"
	"base_project/util/logger"
	"context"
	"fmt"

	"github.com/gorilla/mux"
)

// Controller - Interface for controllers to implement for route registration
type Controller interface {
	RegisterRoutes(router *mux.Router)
}

func StartControllers(router *mux.Router) error {
	logger.L().Info("Registering routes for controllers ...")

	ctx := context.Background()

	// initialize DB
	db := database.ConnectToDB()
	if db == nil {
		logger.L().Fatal("Can't connect to Postgres!")
	}

	gormDB, err := database.OpenGormDB(db)
	if err != nil {
		logger.L().Fatal(fmt.Sprintf("gorm driver errror: %v", err))
	}

	repo := repository.NewRepository(gormDB)

	// Initialise all controllers
	middleware.NewAuthMiddleware(ctx, repo, gormDB)
	authController, err := auth.NewController(ctx, repo, gormDB)
	if err != nil {
		return err
	}
	postController, err := post.NewController(ctx, repo, gormDB)
	if err != nil {
		return err
	}
	commentController, err := comment.NewController(ctx, repo, gormDB)
	if err != nil {
		return err
	}

	// register routes
	authController.RegisterRoutes(router.PathPrefix("/auth/v1").Subrouter())
	postController.RegisterRoutes(router.PathPrefix("/posts/v1").Subrouter())
	commentController.RegisterRoutes(router.PathPrefix("/comments/v1").Subrouter())

	return nil
}

func Initialise() *mux.Router {
	router := mux.NewRouter()

	// Register Middlewares
	router.Use(middleware.CorrelationID)
	router.Use(middleware.ResponseTime)
	router.Use(middleware.PanicRecovery())
	router.Use(middleware.LogRequest)
	router.Use(middleware.LogResponse)

	if err := StartControllers(router); err != nil {
		logger.L().Fatal(fmt.Sprintf("failed to start controllers: %v", err))
	}

	return router
}
