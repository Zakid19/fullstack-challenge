package main

import (
	"context"
	"database/sql"
	"time"
)

type Order struct {
	ID         string
	ProductID  string
	TotalPrice float64
	Status     string
	CreatedAt  time.Time
}

type OrderRepo interface {
	Create(ctx context.Context, o *Order) error
	GetByProduct(ctx context.Context, productId string) ([]Order, error)
}

type pgOrderRepo struct {
	db *sql.DB
}

func NewOrderRepo(db *sql.DB) OrderRepo {
	return &pgOrderRepo{db: db}
}

func (r *pgOrderRepo) Create(ctx context.Context, o *Order) error {
	q := `INSERT INTO orders(id, product_id, total_price, status, created_at) VALUES ($1,$2,$3,$4,$5)`
	_, err := r.db.ExecContext(ctx, q, o.ID, o.ProductID, o.TotalPrice, o.Status, o.CreatedAt)
	return err
}

func (r *pgOrderRepo) GetByProduct(ctx context.Context, productId string) ([]Order, error) {
	q := `SELECT id, product_id, total_price, status, created_at FROM orders WHERE product_id=$1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, q, productId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Order
	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.ID, &o.ProductID, &o.TotalPrice, &o.Status, &o.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, o)
	}
	return out, nil
}
