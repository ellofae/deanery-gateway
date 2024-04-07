package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/ellofae/deanery-gateway/config"
	"github.com/ellofae/deanery-gateway/core/controller/handler"
	"github.com/ellofae/deanery-gateway/core/controller/middleware"
	"github.com/ellofae/deanery-gateway/core/domain/usecase"
	"github.com/ellofae/deanery-gateway/core/session"
	"github.com/ellofae/deanery-gateway/pkg/logger"
)

func establishGatewayHandlers(mux *http.ServeMux, cfg *config.Config) {
	client_usecase := usecase.NewClientUsecase()
	client_handler := handler.NewClientHandler(client_usecase)

	client_handler.RegisterHandlers(mux)
}

func main() {
	logger := logger.GetLogger()
	cfg := config.ParseConfig(config.ConfigureViper())

	middleware.InitJWTSecretKey(cfg)
	session.InitSessionStorage(cfg)

	idleTimeout, _ := strconv.Atoi(cfg.ServerSettings.IdleTimeout)
	readTimeout, _ := strconv.Atoi(cfg.ServerSettings.ReadTimeout)
	writeTimeout, _ := strconv.Atoi(cfg.ServerSettings.WriteTimeout)

	serveMux := http.NewServeMux()
	establishGatewayHandlers(serveMux, cfg)

	srv := &http.Server{
		Addr:         cfg.ServerSettings.BindAddr,
		IdleTimeout:  time.Minute * time.Duration(idleTimeout),
		ReadTimeout:  time.Minute * time.Duration(readTimeout),
		WriteTimeout: time.Minute * time.Duration(writeTimeout),

		Handler: http.TimeoutHandler(middleware.RequestMiddleware(serveMux), 2*time.Minute, ""),
	}

	done := make(chan struct{})
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)

		for {
			if <-c == os.Interrupt {
				if err := srv.Shutdown(context.Background()); err != nil {
					logger.Printf("Error while shutting down the server occured. Error: %v.\n", err.Error())
				}
				close(done)
				return
			}
		}
	}()

	err := srv.ListenAndServe()
	if err != http.ErrServerClosed {
		logger.Printf("Error while starting the server occured. Error: %v.\n", err.Error())
	}
	<-done

	logger.Println("Server gracefully shutdown.")
}
