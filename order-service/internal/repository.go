package internal

import (
	"database/sql"
	"time"
)

type Order struct {
	ID         int       `json:"id"`
	ProductID  string    `json:"productId"`
	TotalPrice float64   `json:"totalPrice"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"createdAt"`
}

type OrderRepo struct {
	db *sql.DB
}

func NewOrderRepo(db *sql.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

func (r *OrderRepo) Save(order *Order) error {
	query := "INSERT INTO orders (product_id, total_price, status, created_at) VALUES ($1,$2,$3,$4) RETURNING id"
	return r.db.QueryRow(query, order.ProductID, order.TotalPrice, order.Status, order.CreatedAt).Scan(&order.ID)
}

func (r *OrderRepo) FindByProductID(productID string) ([]Order, error) {
	rows, err := r.db.Query("SELECT id, product_id, total_price, status, created_at FROM orders WHERE product_id=$1", productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.ID, &o.ProductID, &o.TotalPrice, &o.Status, &o.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}
