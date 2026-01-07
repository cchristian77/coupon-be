package main

import (
	"context"
	api "coupon_be/entrypoint"
	"coupon_be/shared/fhttp"
	"coupon_be/util/config"
	"coupon_be/util/logger"
	"fmt"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log := logger.Initialise()
	defer log.Sync()

	if err := config.LoadConfig(); err != nil {
		logger.L().Fatal(fmt.Sprintf("failed on loading config : %v", err))
		return
	}

	server, err := fhttp.NewHTTPServer(api.Initialise())
	if err != nil {
		logger.L().Fatal(fmt.Sprintf("failed to create http server : %v", err))
		return
	}

	logger.L().Info(fmt.Sprintf("Starting HTTP Server on Port %d ...", config.Env().App.Port))
	if err = server.Start(ctx); err != nil {
		logger.L().Fatal(fmt.Sprintf("failed to start http server, err: %v", err))
	}
}
