package domain

import (
"errors"
"fmt"
"sort"
"time"
)

// SellabilityPolicy Agrupa las reglas de negocio para determinar si un lote puede comercializarse
type SellabilityPolicy struct{}

func (p SellabilityPolicy) IsLotSellable(lot InventoryLot, now time.Time) bool {
return lot.IsSellable(now)
}

// AllocationStrategy Define el comportamiento del motor de asignación de stock con opción FEFO
type AllocationStrategy struct {
EnableFEFO bool // True para farmacias (FEFO obligatorio), False para retail general (FIFO/Estándar)
}

// LotBalancePair Agrupa el lote y su saldo disponible para los cálculos del motor
type LotBalancePair struct {
Lot     InventoryLot
Balance InventoryBalance
}

// AllocateStock Ejecuta la distribución óptima de stock soportando FEFO configurable y fraccionamiento
func (s AllocationStrategy) AllocateStock(
pairs []LotBalancePair,
requestedQuantity int64,
now time.Time,
) ([]InventoryReservation, int64, error) {
if requestedQuantity <= 0 {
return nil, 0, errors.New("la cantidad solicitada debe ser mayor a cero")
}

policy := SellabilityPolicy{}

var eligiblePairs []LotBalancePair
for _, pair := range pairs {
if s.EnableFEFO {
if policy.IsLotSellable(pair.Lot, now) {
eligiblePairs = append(eligiblePairs, pair)
}
} else {
if pair.Lot.Status != StatusBlocked && pair.Lot.Status != StatusDamaged {
eligiblePairs = append(eligiblePairs, pair)
}
}
}

sort.Slice(eligiblePairs, func(i, j int) bool {
lotI := eligiblePairs[i].Lot
lotJ := eligiblePairs[j].Lot

if s.EnableFEFO {
if !lotI.ExpirationDate.Equal(lotJ.ExpirationDate) {
return lotI.ExpirationDate.Before(lotJ.ExpirationDate)
}
}

if !lotI.ReceivedAt.Equal(lotJ.ReceivedAt) {
return lotI.ReceivedAt.Before(lotJ.ReceivedAt)
}

return lotI.LotNumber < lotJ.LotNumber
})

var reservations []InventoryReservation
remainingQty := requestedQuantity

for _, pair := range eligiblePairs {
if remainingQty <= 0 {
break
}

available := pair.Balance.AvailableToSell(pair.Lot.IsSellable(now))
if available <= 0 {
continue
}

qtyToAllocate := remainingQty
if available < remainingQty {
qtyToAllocate = available
}

reservationID := ReservationID(fmt.Sprintf("RES-%d-%s", time.Now().UnixNano(), pair.Lot.LotID))
res, err := NewInventoryReservation(
pair.Lot.TenantID,
reservationID,
pair.Lot.LotID,
pair.Balance.LocationID,
qtyToAllocate,
30,
)
if err != nil {
return nil, 0, err
}

reservations = append(reservations, res)
remainingQty -= qtyToAllocate
}

if remainingQty > 0 {
return reservations, remainingQty, errors.New("stock insuficiente para completar la cantidad solicitada")
}

return reservations, 0, nil
}
