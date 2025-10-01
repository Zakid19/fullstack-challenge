package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CreateOrderReq struct {
	ProductID  string  `json:"productId"`
	TotalPrice float64 `json:"totalPrice"`
}

func CreateOrderHandler(repo OrderRepo, prodClient *ProductClient, rmq *RabbitMQ) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req CreateOrderReq
		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		if req.ProductID == "" || req.TotalPrice <= 0 {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}
		prod, err := prodClient.GetProduct(req.ProductID)
		if err != nil {
			http.Error(w, "product not found", http.StatusBadRequest)
			return
		}
		if prod.Qty <= 0 {
			http.Error(w, "product out of stock", http.StatusBadRequest)
			return
		}
		order := &Order{
			ID:         uuid.New().String(),
			ProductID:  req.ProductID,
			TotalPrice: req.TotalPrice,
			Status:     "created",
			CreatedAt:  getNow(),
		}
		if err := repo.Create(r.Context(), order); err != nil {
			log.Println("failed create order:", err)
			http.Error(w, "internal", http.StatusInternalServerError)
			return
		}
		payload := map[string]interface{}{
			"orderId":    order.ID,
			"productId":  order.ProductID,
			"totalPrice": order.TotalPrice,
			"status":     order.Status,
			"createdAt":  order.CreatedAt,
		}
		if err := rmq.Publish("order.created", payload); err != nil {
			log.Println("failed publish order.created:", err)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(order)
	}
}

func GetOrdersByProductHandler(repo OrderRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productId := chi.URLParam(r, "productId")
		orders, err := repo.GetByProduct(r.Context(), productId)
		if err != nil {
			http.Error(w, "internal", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(orders)
	}
}

func getNow() (t time.Time) {
	return time.Now()
}
