package svc

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gofrs/uuid"
)

func PlaceOrder() (*Order, error) {

	foodReserve, err := http.Post("http://localhost:8082/food/reserve", "application/json", nil)
	if err != nil || foodReserve.StatusCode != http.StatusOK {
		return nil, errors.New("failed to reserve food")
	}
	agentReserve, err := http.Post("http://localhost:8081/agent/reserve", "application/json", nil)
	if err != nil || agentReserve.StatusCode != http.StatusOK {
		return nil, errors.New("failed to reserve agent")
	}
	var order Order
	uuid, err := uuid.NewV7()
	if err != nil {
		return nil, errors.New("failed to generate order id")
	}
	order.ID = uuid.String()

	foodBook, err := http.Post(fmt.Sprintf("http://localhost:8082/food/book?order_id=%s", order.ID), "application/json", nil)
	if err != nil || foodBook.StatusCode != http.StatusOK {
		return nil, errors.New("failed to book food")
	}

	agentBook, err := http.Post(fmt.Sprintf("http://localhost:8081/agent/book?order_id=%s", order.ID), "application/json", nil)
	if err != nil || agentBook.StatusCode != http.StatusOK {
		return nil, errors.New("failed to book agent")
	}
	return &order, nil
}
