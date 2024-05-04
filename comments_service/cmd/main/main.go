package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"syscall"

	"github.com/agdaha/sf_final_project/comments_service/internal/config"
	"github.com/agdaha/sf_final_project/comments_service/internal/handlers/getcomments"
	"github.com/agdaha/sf_final_project/comments_service/internal/handlers/postcomment"
	"github.com/agdaha/sf_final_project/comments_service/internal/storage"
	"github.com/agdaha/sf_final_project/comments_service/internal/storage/posgres"
	"github.com/agdaha/sf_final_project/comments_service/pkg/logger"
	"github.com/agdaha/sf_final_project/comments_service/pkg/middleware"
	"github.com/agdaha/sf_final_project/comments_service/pkg/shutdown"

	"github.com/julienschmidt/httprouter"
)

const (
	postURL = "/api/comments"
	getURL  = "/api/comments/news/:newsId"
)

func main() {
	config := config.New()

	log := logger.New(config.LogLevel, config.Env)

	log.Debug(fmt.Sprintf("Настройки %v", config))

	log.Info(" подключение к БД")
	log.Debug("НАСТРОЙКИ БД", slog.Any("conStr", config.DbUrl()))
	db, err := posgres.New(config.DbUrl(), log)
	if err != nil {
		log.Error(" Проблема с подключением к БД: ", slog.Any("error", err))
		os.Exit(1)
	}

	router := httprouter.New()

	comments_handler := getcomments.New(db, log)
	comment_handler := postcomment.New(db, log)

	router.HandlerFunc(http.MethodGet, getURL, comments_handler)
	router.HandlerFunc(http.MethodPost, postURL, comment_handler)

	start(router, db, log, config)
}

func start(router *httprouter.Router, store storage.Store, log *slog.Logger, cfg *config.Config) {
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
	server = &http.Server{
		Handler:      middleware.RequestID(middlewareLogger(router)),
		WriteTimeout: cfg.Timeout,
		ReadTimeout:  cfg.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	go shutdown.Graceful([]os.Signal{syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM}, log,
		server, store)

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
