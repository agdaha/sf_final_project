package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"syscall"

	censorservice "github.com/agdaha/sf_final_project/api_gateway/internal/clients/censor_service"
	commentservice "github.com/agdaha/sf_final_project/api_gateway/internal/clients/comment_service"
	newsservice "github.com/agdaha/sf_final_project/api_gateway/internal/clients/news_service"
	"github.com/agdaha/sf_final_project/api_gateway/internal/config"
	getnewsdetailed "github.com/agdaha/sf_final_project/api_gateway/internal/handlers/get_news_detailed"
	getnewslist "github.com/agdaha/sf_final_project/api_gateway/internal/handlers/get_news_list"
	postcomment "github.com/agdaha/sf_final_project/api_gateway/internal/handlers/post_comment"
	"github.com/agdaha/sf_final_project/api_gateway/pkg/logger"
	"github.com/agdaha/sf_final_project/api_gateway/pkg/middleware"
	"github.com/agdaha/sf_final_project/api_gateway/pkg/shutdown"

	"github.com/julienschmidt/httprouter"
)

const (
	newsListURL = "/api/news"
	newsURL     = "/api/news/:newsid"
	commentURL  = "/api/comments"
)

func main() {
	config := config.New()

	log := logger.New(config.LogLevel, config.Env)

	log.Debug(fmt.Sprintf("Настройки %v", config))

	router := httprouter.New()
	censorService := censorservice.New(config.CensorServiceURL, log)
	commentService := commentservice.New(config.CommentsServiceURL, log)
	newsService := newsservice.New(config.NewsServiceURL, log)
	newsListHandleFunc := getnewslist.New(newsService, log)
	newsDetailedHandleFunc := getnewsdetailed.New(newsService, commentService, log)
	commentPostHandleFunc := postcomment.New(commentService, censorService, log)

	router.HandlerFunc(http.MethodGet, newsListURL, newsListHandleFunc)
	router.HandlerFunc(http.MethodGet, newsURL, newsDetailedHandleFunc)
	router.HandlerFunc(http.MethodPost, commentURL, commentPostHandleFunc)

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
	server = &http.Server{
		Handler:      middleware.RequestID(middlewareLogger(router)),
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
