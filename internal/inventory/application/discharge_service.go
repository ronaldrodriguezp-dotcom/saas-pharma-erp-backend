package application

import (
"context"
"errors"
"fmt"

"fortiasaass-kio-speed/internal/inventory/domain"
)

type DischargeItem struct {
PresentationID string
Quantity       int
}

type DischargeCommand struct {
TenantID    string
WarehouseID string
IsPerishable bool // True para activar FEFO (Farmacias), False para FIFO (Retail general)
Items       []DischargeItem
}

type InventoryRepositoryPort interface {
GetAvailableLots(ctx context.Context, tenantID, warehouseID, presentationID string) ([]domain.StockLot, error)
UpdateLotQuantity(ctx context.Context, tenantID string, lotID string, newQty int) error
}

type DischargeService struct {
repo InventoryRepositoryPort
}

func NewDischargeService(repo InventoryRepositoryPort) *DischargeService {
return &DischargeService{repo: repo}
}

func (s *DischargeService) DischargeStock(ctx context.Context, cmd DischargeCommand) error {
for _, item := range cmd.Items {
// 1. Obtener lotes disponibles para la presentación en la bodega específica bajo estricto tenant
lots, err := s.repo.GetAvailableLots(ctx, cmd.TenantID, cmd.WarehouseID, item.PresentationID)
if err != nil {
return fmt.Errorf("error obteniendo lotes para el producto %s: %w", item.PresentationID, err)
}

// 2. Seleccionar la estrategia adecuada dinámicamente (FEFO para farmacias / perecederos, FIFO general)
var strategy domain.DischargeStrategy
if cmd.IsPerishable {
strategy = &domain.FefoStrategy{}
} else {
strategy = &domain.FifoStrategy{}
}

sortedLots := strategy.SortLots(lots)

// 3. Descontar stock de forma secuencial cubriendo la cantidad solicitada
remainingToDischarge := item.Quantity

for _, lot := range sortedLots {
if remainingToDischarge <= 0 {
break
}

deductQty := remainingToDischarge
if lot.AvailableQty < remainingToDischarge {
deductQty = lot.AvailableQty
}

newAvailable := lot.AvailableQty - deductQty
err = s.repo.UpdateLotQuantity(ctx, cmd.TenantID, lot.ID, newAvailable)
if err != nil {
return fmt.Errorf("error actualizando inventario del lote %s: %w", lot.ID, err)
}

remainingToDischarge -= deductQty
}

if remainingToDischarge > 0 {
return errors.New("stock insuficiente para completar la cantidad solicitada respetando las políticas de lote")
}
}

return nil
}