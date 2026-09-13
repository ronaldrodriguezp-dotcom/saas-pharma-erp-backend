package domain

import (
"errors"
)

// InventoryBalance Agregado Raíz: Saldos físicos por ubicación y control de concurrencia optimista
type InventoryBalance struct {
TenantID   TenantID
BalanceID  BalanceID
LotID      LotID
LocationID LocationID
OnHand     int64 // Cantidad física estricta (jamás negativa por INV-001)
Reserved   int64 // Cantidad reservada para ventas activas (por INV-002 y INV-003)
Version    int64 // Control de versión para concurrencia optimista
}

func NewInventoryBalance(tenantID TenantID, balanceID BalanceID, lotID LotID, locationID LocationID) InventoryBalance {
return InventoryBalance{
TenantID:   tenantID,
BalanceID:  balanceID,
LotID:      lotID,
LocationID: locationID,
OnHand:     0,
Reserved:   0,
Version:    1,
}
}

// AvailableToSell Calcula de forma estricta el stock vendible real sin fuentes mutables divergentes (INV-011)
func (b InventoryBalance) AvailableToSell(lotIsSellable bool) int64 {
if !lotIsSellable {
return 0
}
available := b.OnHand - b.Reserved
if available < 0 {
return 0
}
return available
}

func (b *InventoryBalance) Reserve(qty int64) error {
if qty <= 0 {
return errors.New("la cantidad a reservar debe ser mayor a cero")
}
if (b.Reserved + qty) > b.OnHand {
return errors.New("inv-003: no es posible reservar una cantidad superior al stock físico disponible (OnHand)")
}
b.Reserved += qty
b.Version++
return nil
}

func (b *InventoryBalance) ReleaseReservation(qty int64) error {
if qty <= 0 {
return errors.New("la cantidad a liberar debe ser mayor a cero")
}
if b.Reserved < qty {
return errors.New("la cantidad a liberar supera las reservas activas actuales")
}
b.Reserved -= qty
b.Version++
return nil
}
