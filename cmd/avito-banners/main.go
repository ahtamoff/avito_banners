package main

import (
	"avito_banners/internal/banner"
	"avito_banners/internal/config"
	"avito_banners/pkg/logging"
	"fmt"
	"net"
	"net/http"
	"time"


	"github.com/julienschmidt/httprouter"
)

// const(
// 	envLocal = "local"
// 	envDev = "dev"
// 	envProd = "prod"
// )



type test struct {
	param int
}

func main() {
	logger := logging.GetLogger()


	logger.Info("create router")
	router := httprouter.New()

	cfg := config.GetConfig()
	logger.Info("register user handler")
	handler := banner.NewHandler(logger)
	handler.Register(router)
	start(router, cfg)

}


func start(router *httprouter.Router, cfg *config.Config) {
	logger := logging.GetLogger()
	logger.Info("start application")

	var listener net.Listener
	var listenErr error

	logger.Info("listen tcp")
	listener, listenErr = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
	logger.Infof("server is listening port %s:%s", cfg.Listen.BindIP, cfg.Listen.Port)

	if listenErr != nil {
		logger.Fatal(listenErr)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	logger.Info("server is listening port 0.0.0.0:1234")

	logger.Fatal(server.Serve(listener))
}

// func setupLogger(env string) *slog.Logger {
// 	var log *slog.Logger

// 	switch env {
// 	case envLocal:
// 		log = slog.New(
// 			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
// 		)
// 	case envDev:
// 		log = slog.New(
// 			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
// 		)
// 	case envProd:
// 		log = slog.New(
// 			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
// 		)
// 	}
// 	return log
// }
