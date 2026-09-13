package ports

import (
"context"
"time"
)

type ReceiptLineData struct {
PresentationID string
LotNumber      string
ExpirationDate time.Time
Quantity       int
LocationID     string
}

type OutboxMessage struct {
ID        string
EventType string
Payload   string
CreatedAt time.Time
}

type GoodsReceiptRepository interface {
// Definiciones base del repositorio si son requeridas por los puertos
}

type InventoryService interface {
ReceiveGoods(ctx context.Context, tenantID string, receiptID string, lines []ReceiptLineData) error
}

type OutboxRepository interface {
Save(ctx context.Context, msg OutboxMessage) error
}