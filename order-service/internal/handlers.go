package internal

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func CreateOrderHandler(service *OrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			ProductID  string  `json:"productId"`
			TotalPrice float64 `json:"totalPrice"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		order, err := service.CreateOrder(input.ProductID, input.TotalPrice)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(order)
	}
}

func GetOrdersByProductHandler(service *OrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productID := chi.URLParam(r, "productId")
		orders, err := service.GetOrdersByProduct(productID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(orders)
	}
}
