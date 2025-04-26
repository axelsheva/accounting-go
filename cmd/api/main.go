package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"db/api"
	"db/ent"

	_ "github.com/lib/pq"
)

func main() {
	// PostgreSQL connection string
	dsn := "postgresql://postgres:password@localhost:5432/postgres?sslmode=disable"

	// Create an ent client
	client, err := ent.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer client.Close()

	// Run the auto migration tool to create all schema resources
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	// Настройка Gin router
	router := api.SetupRouter(client)

	// Запускаем сервер в горутине, чтобы не блокировать основной поток
	go func() {
		if err := router.Run(":8081"); err != nil {
			log.Fatalf("Failed to run server: %v", err)
		}
	}()
	log.Println("Server is running on :8081")

	// Настраиваем обработку сигналов для корректного завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// В данной реализации мы не имеем прямого доступа к http.Server для его корректного завершения,
	// Gin уже запущен через router.Run. В реальном приложении здесь можно добавить
	// код для корректного завершения HTTP-сервера.

	log.Println("Server exiting")
}
