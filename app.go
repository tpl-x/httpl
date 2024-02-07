package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/tpl-x/httpl/internal/config"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type app struct {
	config *config.AppConfig
	logger *slog.Logger
	mux    *http.ServeMux
	svr    *http.Server
}

func newApp(config *config.AppConfig, logger *slog.Logger) *app {
	mux := http.NewServeMux()
	svr := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Server.BindPort),
		Handler: mux,
	}
	svr.RegisterOnShutdown(func() {
		logger.Info("Server shutdown", slog.Int64("timestamp", time.Now().Unix()))
	})
	return &app{
		config: config,
		logger: logger,
		mux:    mux,
		svr:    svr,
	}
}

func (a *app) registerHandler() {
	a.mux.HandleFunc("GET /ping/{user}", func(w http.ResponseWriter, r *http.Request) {
		user := r.PathValue("user")
		_, err := w.Write([]byte(fmt.Sprintf("hello,%s ,this is your `pong` ", user)))
		if err != nil {
			a.logger.Error(err.Error())
			return
		}
		a.logger.Info("pong", slog.String("user", user))
	})
}
func (a *app) start() {
	a.registerHandler()
	// start serve
	go func() {
		a.logger.Info("Server started", slog.String("addr", fmt.Sprintf(":%d", a.config.Server.BindPort)))
		if err := a.svr.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Error(err.Error())
			return
		}
	}()

	// wait signal to exit
	quit := make(chan os.Signal, 1)
	// capture SIGINT（Ctrl+C）and SIGTERM
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// set timeout to exit
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(a.config.Server.GraceExitTimeout)*time.Second)
	defer cancel()

	if err := a.svr.Shutdown(ctx); err != nil {
		a.logger.Error(err.Error())
	}

	a.logger.Info("Server gracefully stopped")

}
