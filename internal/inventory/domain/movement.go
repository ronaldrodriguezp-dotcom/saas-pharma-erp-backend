package domain

import (
"errors"
"time"
)

// MovementType Define los motivos estandarizados para cualquier mutación en el ledger
type MovementType string

const (
TypePurchaseReceipt MovementType = "PURCHASE_RECEIPT"
TypeSale            MovementType = "SALE"
TypeSaleVoid        MovementType = "SALE_VOID"
TypeCustomerReturn  MovementType = "CUSTOMER_RETURN"
TypeSupplierReturn  MovementType = "SUPPLIER_RETURN"
TypeTransferOut     MovementType = "TRANSFER_OUT"
TypeTransferIn      MovementType = "TRANSFER_IN"
TypeAdjustment      MovementType = "STOCK_ADJUSTMENT"
TypeQuarantine      MovementType = "QUARANTINE"
TypeRelease         MovementType = "QUARANTINE_RELEASE"
TypeDestruction     MovementType = "DESTRUCTION"
)

// InventoryMovement Agregado de solo adición (Append-Only) para trazabilidad inmutable
type InventoryMovement struct {
TenantID        TenantID
MovementID      MovementID
LotID           LotID
LocationID      LocationID
Type            MovementType
QuantityChange  int64 // Positivo para ingresos, negativo para salidas
PreviousOnHand  int64
ResultingOnHand int64
Reason          string
ReferenceDoc    string
ActorID         string
CorrelationID   string
OccurredAt      time.Time
}

func NewInventoryMovement(
tenantID TenantID,
movementID MovementID,
lotID LotID,
locationID LocationID,
movType MovementType,
quantityChange int64,
previousOnHand int64,
reason string,
referenceDoc string,
actorID string,
correlationID string,
) (InventoryMovement, error) {
if tenantID == "" || movementID == "" || lotID == "" || locationID == "" {
return InventoryMovement{}, errors.New("parámetros obligatorios vacíos para el movimiento de inventario")
}

resultingOnHand := previousOnHand + quantityChange

// Invariante INV-001: OnHand jamás puede ser negativo bajo ninguna circunstancia
if resultingOnHand < 0 {
return InventoryMovement{}, errors.New("violación de invariante INV-001: el stock físico (OnHand) no puede ser negativo")
}

return InventoryMovement{
TenantID:        tenantID,
MovementID:      movementID,
LotID:           lotID,
LocationID:      locationID,
Type:            movType,
QuantityChange:  quantityChange,
PreviousOnHand:  previousOnHand,
ResultingOnHand: resultingOnHand,
Reason:          reason,
ReferenceDoc:    referenceDoc,
ActorID:         actorID,
CorrelationID:   correlationID,
OccurredAt:      time.Now().UTC(),
}, nil
}
