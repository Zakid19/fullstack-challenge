package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
)

func ensureOrdersTable(db *sql.DB) error {
	q := `CREATE TABLE IF NOT EXISTS orders (
    id TEXT PRIMARY KEY,
    product_id TEXT NOT NULL,
    total_price NUMERIC,
    status TEXT,
    created_at TIMESTAMP DEFAULT now()
  )`
	_, err := db.Exec(q)
	return err
}

func main() {
	dbUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "dev"),
		getEnv("DB_PASSWORD", "dev"),
		getEnv("DB_NAME", "appdb"),
	)
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(20)
	db.SetConnMaxIdleTime(5 * time.Minute)

	if err := ensureOrdersTable(db); err != nil {
		log.Fatal("failed ensure orders table:", err)
	}

	repo := NewOrderRepo(db)
	prodClient := NewProductClient(getEnv("PRODUCT_SERVICE_URL", "http://localhost:3000"))
	rmq := NewRabbitMQ(getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672"))
	if err := rmq.Connect(); err != nil {
		log.Fatal("failed to connect rabbitmq:", err)
	}

	// subscribe for logging/side-effects
	go rmq.Subscribe("order.service.logger", "order.created", func(payload map[string]interface{}) {
		log.Println("Order-service consumed order.created:", payload)
	})

	r := chi.NewRouter()
	r.Post("/orders", CreateOrderHandler(repo, prodClient, rmq))
	r.Get("/orders/product/{productId}", GetOrdersByProductHandler(repo))

	server := &http.Server{
		Addr:    ":4000",
		Handler: r,
	}
	log.Println("order-service listening on :4000")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
