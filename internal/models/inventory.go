package models

import "github.com/google/uuid"

type InventoryItem struct {
	UUID uuid.UUID
	Qty  int
}
