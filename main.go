package main

import (
	"context"
	"errors"
	"github.com/CookieNyanCloud/driveApi/internal/config"
	"github.com/CookieNyanCloud/driveApi/internal/delivery"
	ginServer "github.com/CookieNyanCloud/driveApi/internal/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

const (
	credFile = "driveapisearch.json"
)

func main() {
	conf := config.InitConf()
	ctx := context.Background()

	srv, err := drive.NewService(ctx, option.WithCredentialsFile(credFile))
	if err != nil {
		log.Fatalf("Unable to parse credantials file: %v", err)
	}
	handler := delivery.NewHandler(srv, conf)
	server := ginServer.NewServer(conf, handler.Init())

	go func() {
		if err := server.Run(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error occurred while running http serve: %v\n", err)
		}
	}()
	log.Println("started")
	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	const timeout = 5 * time.Second
	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

}
