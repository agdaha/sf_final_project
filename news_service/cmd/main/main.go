package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"syscall"

	"github.com/agdaha/sf_final_project/news_service/internal/config"
	"github.com/agdaha/sf_final_project/news_service/internal/handlers/getonepost"
	"github.com/agdaha/sf_final_project/news_service/internal/handlers/getposts"
	"github.com/agdaha/sf_final_project/news_service/internal/reader"
	"github.com/agdaha/sf_final_project/news_service/internal/storage"
	"github.com/agdaha/sf_final_project/news_service/internal/storage/postgres"
	"github.com/agdaha/sf_final_project/news_service/pkg/logger"
	"github.com/agdaha/sf_final_project/news_service/pkg/middleware"
	"github.com/agdaha/sf_final_project/news_service/pkg/shutdown"
	"github.com/julienschmidt/httprouter"
)

const (
	getPostsURL = "/api/news"
	getPostURL  = "/api/news/:newsid"
)

func main() {
	config := config.New()

	log := logger.New(config.LogLevel, config.Env)

	log.Debug(fmt.Sprintf("Настройки %v", config))

	//Подключение к бд
	log.Info(" подключение к БД")
	store, err := postgres.New(config.DbUrl(), log)
	if err != nil {
		log.Error(" Проблема с подключением к БД", slog.Any("error", err))
		os.Exit(1)
	}

	log.Info("Запуск потоков чтения новосте с ресурсов")
	reader := reader.New(config, store, log)
	chP, chE := reader.Start()
	defer close(chP)
	defer close(chE)

	router := httprouter.New()
	postsHandler := getposts.New(store, log)
	postHandler := getonepost.New(store, log)
	router.HandlerFunc(http.MethodGet, getPostURL, postHandler)
	router.HandlerFunc(http.MethodGet, getPostsURL, postsHandler)
	start(router, store, log, config)
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
