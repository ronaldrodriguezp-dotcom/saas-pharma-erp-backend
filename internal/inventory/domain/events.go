package domain

import (
"time"
)

// DomainEvent Interfaz base para todos los eventos de dominio del inventario
type DomainEvent interface {
EventName() string
OccurredAt() time.Time
TenantID() TenantID
}

// InventoryMovementRecorded Se dispara cada vez que un movimiento es agregado al ledger inmutable
type InventoryMovementRecorded struct {
Tenant        TenantID
Movement      MovementID
Lot           LotID
Location      LocationID
Type          MovementType
QtyChange     int64
ResultingOnHand int64
Timestamp     time.Time
}

func (e InventoryMovementRecorded) EventName() string  { return "inventory.movement.recorded" }
func (e InventoryMovementRecorded) OccurredAt() time.Time { return e.Timestamp }
func (e InventoryMovementRecorded) TenantID() TenantID   { return e.Tenant }

// InventoryReserved Se dispara cuando se genera una reserva temporal desde el POS o canales de venta
type InventoryReserved struct {
Tenant      TenantID
Reservation ReservationID
Lot         LotID
Location    LocationID
Quantity    int64
ExpiresAt   time.Time
Timestamp   time.Time
}

func (e InventoryReserved) EventName() string  { return "inventory.reserved" }
func (e InventoryReserved) OccurredAt() time.Time { return e.Timestamp }
func (e InventoryReserved) TenantID() TenantID   { return e.Tenant }
