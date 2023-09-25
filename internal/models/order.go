package models

import "github.com/google/uuid"

type OrderStatus int64

const (
	OrderStatusUndefined OrderStatus = iota
	OrderStatusCreated
	OrderStatusConfirmed
	OrderStatusPackaged
	OrderStatusInDelivery
	OrderStatusCompleted
	OrderStatusCancelled
)

func (s OrderStatus) ToString() string {
	switch s {
	default:
	case OrderStatusUndefined:
		return "undegined"
	case OrderStatusCreated:
		return "created"
	case OrderStatusConfirmed:
		return "confirmed"
	case OrderStatusPackaged:
		return "packaged"
	case OrderStatusInDelivery:
		return "indelivery"
	case OrderStatusCompleted:
		return "completed"
	case OrderStatusCancelled:
		return "cancelled"
	}
	return ""
}

type OrderStatusNotify struct {
	OrderId       uuid.UUID   `json:"order_id"`
	Status        OrderStatus `json:"status"`
	StatusMessage string      `json:"status_message"`
}

type Order struct {
	OrderId       uuid.UUID    `json:"order_id" validate:"omitempty"`
	Status        OrderStatus  `json:"status"`
	StatusMessage string       `json:"status_message"`
	Sum           int          `json:"sum" validate:"omitempty"`
	OrderList     []*OrderItem `json:"order_list"`
}

type OrderItem struct {
	ItemId uuid.UUID `json:"item_id" validate:"omitempty"`
	Cost   int       `json:"cost" validate:"min=1"`
	Qty    int       `json:"qty" validate:"min=1"`
	Sum    int       `json:"sum" validate:"omitempty"`
}

func (o *Order) CalculateSum() {
	sum := 0
	for _, item := range o.OrderList {
		sum += item.GetSum()
	}
	o.Sum = sum
}

func (o *Order) GetSum() int {
	o.CalculateSum()
	return o.Sum
}

func (o *OrderItem) CalculateSum() {
	o.Sum = o.Qty * o.Cost
}

func (o *OrderItem) GetSum() int {
	o.CalculateSum()
	return o.Sum
}
