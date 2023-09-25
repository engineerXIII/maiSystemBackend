package models

import "github.com/google/uuid"

type OrderStatus int64

const (
	Undefined OrderStatus = iota
	Created
	Confirmed
	Packaged
	InDelivery
	Completed
	Cancelled
)

func (s OrderStatus) ToString() string {
	switch s {
	default:
	case Undefined:
		return "undegined"
	case Created:
		return "created"
	case Confirmed:
		return "confirmed"
	case Packaged:
		return "packaged"
	case InDelivery:
		return "indelivery"
	case Completed:
		return "completed"
	case Cancelled:
		return "cancelled"
	}
	return ""
}

type Order struct {
	OrderId   uuid.UUID    `json:"order_id" validate:"omitempty"`
	Status    string       `json:"status"`
	sum       int          `json:"sum" validate:"min=1"`
	OrderList []*OrderItem `json:"order_list"`
}

type OrderItem struct {
	ItemId uuid.UUID `json:"order_id" validate:"omitempty"`
	Cost   int       `json:"cost" validate:"min=1"`
	Qty    int       `json:"qty" validate:"min=1"`
	sum    int       `json:"sum" validate:"min=1"`
}

func (o *Order) CalculateSum() {
	sum := 0
	for _, item := range o.OrderList {
		sum += item.GetSum()
	}
	o.sum = sum
}

func (o *Order) GetSum() int {
	o.CalculateSum()
	return o.sum
}

func (o *OrderItem) CalculateSum() {
	o.sum = o.Qty * o.Cost
}

func (o *OrderItem) GetSum() int {
	o.CalculateSum()
	return o.sum
}
