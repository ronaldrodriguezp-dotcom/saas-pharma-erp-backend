package domain

import (
"errors"
"time"
)

// QualityStatus Define el estado sanitario y logístico del lote (Invariantes de sellabilidad)
type QualityStatus string

const (
StatusAvailable  QualityStatus = "AVAILABLE"
StatusQuarantine QualityStatus = "QUARANTINE"
StatusBlocked    QualityStatus = "BLOCKED"
StatusExpired    QualityStatus = "EXPIRED"
StatusRecalled   QualityStatus = "RECALLED"
StatusDamaged    QualityStatus = "DAMAGED"
)

// InventoryLot Agregado Raíz: Identidad y atributos sanitarios del lote farmacéutico
type InventoryLot struct {
TenantID       TenantID
LotID          LotID
PresentationID      PresentationID
LotNumber      string
ExpirationDate time.Time
Status         QualityStatus
ReceivedAt     time.Time
}

func NewInventoryLot(tenantID TenantID, lotID LotID, PresentationID PresentationID, lotNumber string, expirationDate time.Time, receivedAt time.Time) (InventoryLot, error) {
if tenantID == "" || lotID == "" || PresentationID == "" || lotNumber == "" {
return InventoryLot{}, errors.New("parámetros obligatorios vacíos para la creación del lote")
}

status := StatusAvailable
if time.Now().After(expirationDate) {
status = StatusExpired
}

return InventoryLot{
TenantID:       tenantID,
LotID:          lotID,
PresentationID:      PresentationID,
LotNumber:      lotNumber,
ExpirationDate: expirationDate,
Status:         status,
ReceivedAt:     receivedAt,
}, nil
}

// IsSellable Determina si el lote cumple con las condiciones sanitarias para la venta (FEFO / Dispensación)
func (l InventoryLot) IsSellable(now time.Time) bool {
if now.After(l.ExpirationDate) {
return false
}
switch l.Status {
case StatusAvailable:
return true
case StatusQuarantine, StatusBlocked, StatusExpired, StatusRecalled, StatusDamaged:
return false
default:
return false
}
}
