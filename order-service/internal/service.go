package internal

import (
	"time"
)

type OrderService struct {
	repo          *OrderRepo
	productClient *ProductClient
	rabbit        *RabbitMQ
}

func NewOrderService(repo *OrderRepo, client *ProductClient, rabbit *RabbitMQ) *OrderService {
	return &OrderService{repo: repo, productClient: client, rabbit: rabbit}
}

func (s *OrderService) CreateOrder(productID string, totalPrice float64) (*Order, error) {
	product, err := s.productClient.GetProduct(productID)
	if err != nil {
		return nil, err
	}

	if product.Qty <= 0 {
		return nil, ErrOutOfStock
	}

	order := &Order{
		ProductID:  productID,
		TotalPrice: totalPrice,
		Status:     "created",
		CreatedAt:  time.Now(),
	}

	if err := s.repo.Save(order); err != nil {
		return nil, err
	}

	if s.rabbit != nil {
		s.rabbit.Publish("order.created", order)
	}

	return order, nil
}

func (s *OrderService) GetOrdersByProduct(productID string) ([]Order, error) {
	return s.repo.FindByProductID(productID)
}
