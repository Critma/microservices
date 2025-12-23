package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/critma/prodfiles/cmd/api/config"
	"github.com/critma/prodfiles/cmd/api/handlers"
	"github.com/critma/prodfiles/internal/store"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.SetConfig()
	if err != nil {
		log.Fatalf("Config load error: %v", err)
	}

	l := log.New(os.Stdout, cfg.LogLevel+" ", log.LstdFlags)
	local, err := store.NewLocal(cfg.BasePath, 1024*1000*5)
	if err != nil {
		l.Fatal("Local store failed")
	}
	app := config.Application{
		Config: cfg,
		Logger: l,
		Store:  local,
	}

	routerEngine := gin.Default()
	handlers.AddHandlers(routerEngine, app)

	serverAddr := fmt.Sprintf("%s:%s", app.Config.Addr, app.Config.Port)
	srv := &http.Server{
		Addr:         serverAddr,
		Handler:      routerEngine,
		ErrorLog:     l,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		app.Logger.Printf("Server start on: %v\n", serverAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Fatalf("listen: %s\n", err)
		}
	}()

	quitChan := make(chan os.Signal, 1)
	signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM)
	<-quitChan
	l.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		l.Println("Server shutdown:", err)
	}
	log.Println("Server exiting")
}
