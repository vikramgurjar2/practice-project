package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vikramgurjar2/practice-project/internal/config"
	"github.com/vikramgurjar2/practice-project/internal/http/handlers/students"
	"github.com/vikramgurjar2/practice-project/internal/storage"
	"github.com/vikramgurjar2/practice-project/internal/storage/sqlite"
)

func main() {
	//load config
	cfg := config.MustLoad()

	//database setup
	db, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal("failed to initialize database:", err)
	}

	slog.Info("database initialized successfully")
	//setup routes we will use inbuilt router mathods for routing and path
	router := http.NewServeMux()

	var store storage.Storage = db

	router.HandleFunc("POST /api/students", students.New(store))

	//setup server

	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}
	//start server
	fmt.Println("server started at", cfg.Addr)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("error starting server:", err)
		}
	}()

	<-done
	slog.Info("shutting down the server")

	//graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		slog.Error("error shutting down server", "error", err)
	}
	slog.Info("server exited properly")

}
