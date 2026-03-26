package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"notes-service/internal/database"
	"notes-service/internal/routers"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	appPort    string
	dbHost     string
	dbPort     string
	dbUser     string
	dbPassword string
	dbName     string
	debugMode  bool
)

func main() {
	flag.StringVar(&appPort, "port", "5000", "Port for the HTTP server")
	flag.StringVar(&dbHost, "db-host", "127.0.0.1", "Database host")
	flag.StringVar(&dbPort, "db-port", "3306", "Database port")
	flag.StringVar(&dbUser, "db-user", "root", "Database user")
	flag.StringVar(&dbPassword, "db-pass", "password123", "Database password")
	flag.StringVar(&dbName, "db-name", "notes_db", "Database name")
	flag.BoolVar(&debugMode, "debug", false, "Debug mode")
	flag.Parse()

	if !debugMode {
		gin.SetMode(gin.ReleaseMode)
	}
	
	log.Printf("Connecting to database %s at %s:%s...", dbName, dbHost, dbPort)
	db, err := database.New(dbHost, dbUser, dbPort, dbPassword, dbName)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	defer db.Close()

	router := routers.Router(db)
	router.SetTrustedProxies(nil)

	srv := &http.Server{
		Addr:    ":" + appPort,
		Handler: router, 
	}

	go func() {
		log.Printf("Starting server on port %s...", appPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("Server exiting")
}
