package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"accounting/api"
	"accounting/ent"

	"entgo.io/ent/dialect/sql"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	jsoniter "github.com/json-iterator/go"
	_ "github.com/lib/pq"
)

func main() {
	// PostgreSQL connection string
	dsn := "postgresql://postgres:password@localhost:5432/postgres?sslmode=disable"

	// Create database driver with connection pool configuration
	drv, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}

	// Configure connection pool
	accounting := drv.DB()
	accounting.SetMaxOpenConns(100)          // Maximum number of open connections
	accounting.SetMaxIdleConns(50)           // Maximum number of idle connections
	accounting.SetConnMaxLifetime(time.Hour) // Maximum connection lifetime

	// Create an ent client
	client := ent.NewClient(ent.Driver(drv))
	defer client.Close()

	// Run the auto migration tool to create all schema resources
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	// Configure faster JSON decoder
	binding.EnableDecoderUseNumber = true

	// Use jsoniter for faster JSON processing through middleware
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	// Set Gin to release mode to disable debug logs
	gin.SetMode(gin.ReleaseMode)

	// Setup Gin router
	router := api.SetupRouter(client)

	// Add JSON middleware to use jsoniter for response marshaling
	router.Use(func(c *gin.Context) {
		c.Set("json", json)
		c.Next()
	})

	// Create a server with graceful shutdown support
	server := &http.Server{
		Addr:         ":8081",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Run server in a goroutine to not block the main thread
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to run server: %v", err)
		}
	}()
	log.Println("Server is running on :8081")

	// Configure signal handling for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server shutdown error:", err)
	}

	log.Println("Server exiting")
}
