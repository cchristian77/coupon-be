package entrypoint

import (
	"context"
	"coupon_be/entrypoint/comment"
	"coupon_be/entrypoint/user"
	"coupon_be/repository"
	"coupon_be/shared/external/database"
	"coupon_be/shared/fhttp/middleware"
	"coupon_be/util/logger"
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
	userController, err := user.NewController(ctx, repo, gormDB)
	if err != nil {
		return err
	}
	commentController, err := comment.NewController(ctx, repo, gormDB)
	if err != nil {
		return err
	}

	// register routes
	userController.RegisterRoutes(router.PathPrefix("/users/v1").Subrouter())
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
