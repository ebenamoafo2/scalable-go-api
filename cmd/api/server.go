package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}

	//A channel to listen for errors that happen during shutdown
	shutdownError := make(chan error)

	go func() {
		//intercept the signals
		quit := make(chan os.Signal, 1)

		//Notify the server of the interrupt signal
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		//Block until we receive our signal.
		s := <-quit

		app.logger.Info("shutting down server", "signal", s.String())

		//Give the server 30-second graceful shutdown period to finish
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		//Shutdown the server
		shutdownError <- srv.Shutdown(ctx)
	}()

	app.logger.Info("starting server", "addr", srv.Addr, "env", app.config.env)

	//Block until the server is shutdown
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	//Wait for the server to shut down
	err = <-shutdownError
	if err != nil {
		return err
	}

	app.logger.Info("server stopped", "addr", srv.Addr)

	return nil
}
