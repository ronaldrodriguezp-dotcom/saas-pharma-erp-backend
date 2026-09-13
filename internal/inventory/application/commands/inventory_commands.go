package commands

import (
"context"
"errors"
"time"

)

// ReceiveStockCommand Comando para registrar el ingreso de un lote recibido por compra
type ReceiveStockCommand struct {
TenantID       string
LotID          string
PresentationID      string
LocationID     string
LotNumber      string
Quantity       int64
ExpirationDate time.Time
ReferenceDoc   string
ActorID        string
CorrelationID  string
}

// ReserveStockCommand Comando para solicitar una reserva temporal de inventario desde el POS
type ReserveStockCommand struct {
TenantID          string
ReservationID     string
LotID             string
LocationID        string
RequestedQuantity int64
}

// InventoryCommandHandler Gestiona la ejecución transaccional de los comandos de inventario
type InventoryCommandHandler struct {
// Aquí se inyectarán posteriormente los repositorios y el puerto de Unit of Work
}

func NewInventoryCommandHandler() *InventoryCommandHandler {
return &InventoryCommandHandler{}
}

func (h *InventoryCommandHandler) HandleReceiveStock(ctx context.Context, cmd ReceiveStockCommand) error {
if cmd.TenantID == "" || cmd.LotID == "" || cmd.Quantity <= 0 {
return errors.New("parámetros inválidos para la recepción de stock")
}
// Lógica de aplicación que invocará al agregado InventoryLot, InventoryBalance y Ledger Movement
return nil
}
