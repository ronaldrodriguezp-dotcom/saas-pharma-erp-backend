package ports

import (
"context"
"time"
"fortiasaass-kio-speed/internal/purchasing/domain"
)

type ReceiptLineData struct {
PresentationID string
LotNumber      string
ExpirationDate time.Time
Quantity       int
LocationID     string
}

type OutboxMessage struct {
EventID     string
TenantID    string
AggregateID string
EventType   string
Payload     []byte
OccurredAt  time.Time
}

type GoodsReceiptRepository interface {
Create(ctx context.Context, gr *domain.GoodsReceipt) error
GetByID(ctx context.Context, tenantID domain.TenantID, id domain.GoodsReceiptID) (*domain.GoodsReceipt, error)
Update(ctx context.Context, gr *domain.GoodsReceipt) error
}

type SupplierRepository interface {
Create(ctx context.Context, supplier *domain.Supplier) error
}

type InventoryService interface {
ReceiveGoods(ctx context.Context, tenantID domain.TenantID, receiptID domain.GoodsReceiptID, lines []ReceiptLineData) error
}

type OutboxRepository interface {
Save(ctx context.Context, msg OutboxMessage) error
}