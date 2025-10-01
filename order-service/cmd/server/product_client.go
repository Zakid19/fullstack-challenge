package main // ganti jadi 'server' kalau semua file lo udah pakai package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Product struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Qty   int     `json:"qty"`
}

type ProductClient struct {
	base   string
	client *http.Client
}

func NewProductClient(base string) *ProductClient {
	return &ProductClient{
		base:   base,
		client: &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *ProductClient) GetProduct(id string) (*Product, error) {
	url := fmt.Sprintf("%s/products/%s", c.base, id)

	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("product not found (status code: %d)", resp.StatusCode)
	}

	var p Product
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		return nil, err
	}

	return &p, nil
}
