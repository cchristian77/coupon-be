package fhttp

import (
	"context"
	"coupon_be/util/config"
	"coupon_be/util/logger"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

// Server - struct that contains the config and handler for the http server
type Server struct {
	server  *http.Server
	options *serverOptions
	exit    chan os.Signal
}

// NewHTTPServer - creates a new server with the given config
// config: configuration required to initialise server, Port is mandatory and should be greater than 1
// router: gorilla mux router mandatory
// healthCheckHandler: handler registered for /healthCheckURL mandatory
// option: option to add MethodNotAllowed handler, NotFoundHandler, Timeout while stoping server
// Returns Server: wrapper of http server
// error: Error code ErrCodeConfig invalid config
// error: Error code ErrCodeHealthCheckMissing healthcheck handler is nil
// error: Error code ErrCodeRouterNil router is nil
// error: Error code ErrCodeInvalidOption Invalid options
func NewHTTPServer(router *mux.Router) (*Server, error) {
	options := &serverOptions{
		methodNotAllowedHandler:  AppHandler(DefaultMethodNotAllowedHandler),
		notFoundHandler:          AppHandler(DefaultNotFoundHandler),
		serverStopTimeoutSeconds: defaultStopTimeoutSeconds,
		readTimeoutSeconds:       defaultReadTimeoutSeconds,
		readHeaderTimeoutSeconds: defaultReadHeaderTimeoutSeconds,
	}

	if config.Env() == nil {
		return nil, errors.New("Config is nil")
	}

	if config.Env().App.Port < 1 {
		return nil, errors.New("Port is less than 1")
	}

	if router == nil {
		return nil, errors.New("Router should not be nil")
	}

	router.Handle(healthCheckURL, AppHandler(DefaultHealthCheckHandler)).Methods(http.MethodGet)

	if options.methodNotAllowedHandler != nil {
		router.MethodNotAllowedHandler = options.methodNotAllowedHandler
	}

	if options.notFoundHandler != nil {
		router.NotFoundHandler = options.notFoundHandler
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%v", config.Env().App.Port),
		Handler: router,
		// refer:
		// - https://app.deepsource.com/directory/analyzers/go/issues/GO-S2112
		// - https://adam-p.ca/blog/2022/01/golang-http-server-timeouts/
		ReadTimeout:       time.Duration(options.readTimeoutSeconds) * time.Second,
		ReadHeaderTimeout: time.Duration(options.readHeaderTimeoutSeconds) * time.Second,
	}

	return &Server{
		server:  httpServer,
		options: options,
		exit:    make(chan os.Signal),
	}, nil
}

// Start - Starts the http server
// Error: error code ErrCodeStartServer if http server throws error
func (s *Server) Start(ctx context.Context) (err error) {
	logger.L().Info(fmt.Sprintf("Starting Http Server at {%v}....", s.server.Addr))

	timeout := time.Duration(s.options.serverStopTimeoutSeconds) * time.Second

	idleConnClosed := make(chan struct{})

	// For graceful shutdown
	go func() {
		signal.Notify(s.exit, os.Interrupt, syscall.SIGTERM)

		sig := <-s.exit
		ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)

		defer cancel()

		logger.L().Info(fmt.Sprintf("Received signal {%v}, shutting down HTTP Server....", sig.String()))

		if err = s.server.Shutdown(ctx); err != nil {
			logger.L().Error(fmt.Sprintf("Error while shutting down HTTP server: %v", err))
		}

		close(idleConnClosed)
	}()

	// Start the server
	// here, we use a different var `serveErr` as ListenAndServe() always returns a non-nil error, and we want to
	// return the error only if it's not http.ErrServerClosed
	if serveErr := s.server.ListenAndServe(); !errors.Is(serveErr, http.ErrServerClosed) {
		logger.L().Error(fmt.Sprintf("Error starting up Http Server. Error: %v", serveErr.Error()))
		err = serveErr
	}

	// Wait for the shutdown process to complete
	<-idleConnClosed

	return nil
}

// Stop - stops the server
// Error: error code ErrCodeStopServer if http server throws error
func (s *Server) Stop(ctx context.Context) {
	logger.L().Info("Shutting down Http Server....")
	s.exit <- os.Interrupt
}
