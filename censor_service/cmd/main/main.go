package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"syscall"

	"github.com/agdaha/sf_final_project/censor_service/internal/config"
	"github.com/agdaha/sf_final_project/censor_service/internal/handlers/censor"
	"github.com/agdaha/sf_final_project/censor_service/internal/storage/memdb"
	"github.com/agdaha/sf_final_project/censor_service/pkg/logger"
	"github.com/agdaha/sf_final_project/censor_service/pkg/middleware"
	"github.com/agdaha/sf_final_project/censor_service/pkg/shutdown"

	"github.com/julienschmidt/httprouter"
)

const (
	URL = "/api/censor"
)

func main() {
	config := config.New()

	log := logger.New(config.LogLevel, config.Env)

	log.Debug(fmt.Sprintf("Настройки %v", config))

	log.Info("Инициализация хранилища")
	db, err := memdb.New()
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	router := httprouter.New()
	handler := censor.New(db, log)
	router.HandlerFunc(http.MethodPost, URL, handler)
	start(router, log, config)
}

func start(router *httprouter.Router, log *slog.Logger, cfg *config.Config) {
	var server *http.Server
	var listener net.Listener

	log.Info(fmt.Sprintf("bind application to addres: %s", cfg.HTTPServer.Address))

	var err error

	listener, err = net.Listen("tcp", cfg.HTTPServer.Address)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	middlewareLogger := middleware.Logger(log)
	routWithLogger := middlewareLogger(router)
	routWithRequestId := middleware.RequestID(routWithLogger)

	server = &http.Server{
		Handler:      routWithRequestId,
		WriteTimeout: cfg.Timeout,
		ReadTimeout:  cfg.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	go shutdown.Graceful([]os.Signal{syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM}, log,
		server)

	log.Info("application initialized and started")

	if err := server.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			log.Warn("server shutdown")
		default:
			log.Error(err.Error())
			os.Exit(1)
		}
	}
}
