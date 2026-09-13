package domain

import (
"errors"
"time"
)

// ReservationStatus Define la máquina de estados estricta para las reservas de inventario
type ReservationStatus string

const (
ReservationActive    ReservationStatus = "ACTIVE"
ReservationCommitted ReservationStatus = "COMMITTED"
ReservationReleased  ReservationStatus = "RELEASED"
ReservationExpired   ReservationStatus = "EXPIRED"
ReservationCancelled ReservationStatus = "CANCELLED"
)

// InventoryReservation Agregado para gestionar bloqueos temporales de stock con alta concurrencia
type InventoryReservation struct {
TenantID      TenantID
ReservationID ReservationID
LotID         LotID
LocationID    LocationID
Quantity      int64
Status        ReservationStatus
ExpiresAt     time.Time
CreatedAt     time.Time
}

func NewInventoryReservation(
tenantID TenantID,
reservationID ReservationID,
lotID LotID,
locationID LocationID,
quantity int64,
ttlMinutes time.Duration,
) (InventoryReservation, error) {
if tenantID == "" || reservationID == "" || lotID == "" || locationID == "" {
return InventoryReservation{}, errors.New("parámetros obligatorios vacíos para crear la reserva")
}
if quantity <= 0 {
return InventoryReservation{}, errors.New("la cantidad a reservar debe ser mayor a cero")
}

now := time.Now().UTC()
return InventoryReservation{
TenantID:      tenantID,
ReservationID: reservationID,
LotID:         lotID,
LocationID:    locationID,
Quantity:      quantity,
Status:        ReservationActive,
ExpiresAt:     now.Add(ttlMinutes * time.Minute),
CreatedAt:     now,
}, nil
}

// IsExpired Verifica si la reserva ha superado su tiempo límite de vida sin completarse el pago
func (r InventoryReservation) IsExpired(now time.Time) bool {
if r.Status != ReservationActive {
return false
}
return now.After(r.ExpiresAt)
}

// Commit Transiciona la reserva a estado asegurado (Invariante INV-010: una vez committed no vuelve a active)
func (r *InventoryReservation) Commit() error {
if r.Status != ReservationActive {
return errors.New("solo se pueden comprometer reservas en estado activo")
}
r.Status = ReservationCommitted
return nil
}

// Release Libera la reserva para devolver el stock a la disponibilidad general
func (r *InventoryReservation) Release() error {
if r.Status == ReservationCommitted {
return errors.New("inv-010: una reserva ya comprometida (COMMITTED) no no puede ser liberada directamente")
}
if r.Status == ReservationReleased || r.Status == ReservationCancelled || r.Status == ReservationExpired {
return nil // Ya se encuentra inactiva
}
r.Status = ReservationReleased
return nil
}
