package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Product struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Qty   int     `json:"qty"`
}

type ProductClient struct {
	baseURL string
}

func NewProductClient(baseURL string) *ProductClient {
	return &ProductClient{baseURL: baseURL}
}

func (c *ProductClient) GetProduct(id string) (*Product, error) {
	resp, err := http.Get(fmt.Sprintf("%s/products/%s", c.baseURL, id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("product not found")
	}

	var p Product
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		return nil, err
	}
	return &p, nil
}
