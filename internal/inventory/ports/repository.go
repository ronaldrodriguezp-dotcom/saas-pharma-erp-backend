package ports

import (
"context"
"fortiasaass-kio-speed/internal/inventory/domain"
)

// LotRepository Puerto para la persistencia del agregado InventoryLot
type LotRepository interface {
Save(ctx context.Context, lot domain.InventoryLot) error
FindByID(ctx context.Context, tenantID domain.TenantID, lotID domain.LotID) (domain.InventoryLot, error)
FindByProductID(ctx context.Context, tenantID domain.TenantID, PresentationID domain.PresentationID) ([]domain.InventoryLot, error)
}

// BalanceRepository Puerto para la persistencia y control de concurrencia de InventoryBalance
type BalanceRepository interface {
Save(ctx context.Context, balance domain.InventoryBalance) error
FindByLotAndLocation(ctx context.Context, tenantID domain.TenantID, lotID domain.LotID, locationID domain.LocationID) (domain.InventoryBalance, error)
FindBalancesByProduct(ctx context.Context, tenantID domain.TenantID, PresentationID domain.PresentationID) ([]domain.LotBalancePair, error)
}

// MovementRepository Puerto de solo adición para el Ledger inmutable
type MovementRepository interface {
Append(ctx context.Context, movement domain.InventoryMovement) error
}

// ReservationRepository Puerto para la gestión de reservas temporales en el POS
type ReservationRepository interface {
Save(ctx context.Context, reservation domain.InventoryReservation) error
FindByID(ctx context.Context, tenantID domain.TenantID, reservationID domain.ReservationID) (domain.InventoryReservation, error)
UpdateStatus(ctx context.Context, tenantID domain.TenantID, reservationID domain.ReservationID, status domain.ReservationStatus) error
}

// UnitOfWork Puerto para coordinar transacciones atómicas locales (sin llamadas HTTP externas)
type UnitOfWork interface {
Do(ctx context.Context, fn func(ctx context.Context) error) error
}
