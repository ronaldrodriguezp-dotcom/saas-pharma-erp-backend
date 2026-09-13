package domain

import "errors"

type PurchaseOrderID string
type LocationID string // Referencia a la sucursal/bodega de destino

// PurchaseOrderStatus define la máquina de estados de la orden
type PurchaseOrderStatus string

const (
POStatusDraft             PurchaseOrderStatus = "DRAFT"
POStatusSubmitted         PurchaseOrderStatus = "SUBMITTED"
POStatusApproved          PurchaseOrderStatus = "APPROVED"
POStatusPartiallyReceived PurchaseOrderStatus = "PARTIALLY_RECEIVED"
POStatusReceived          PurchaseOrderStatus = "RECEIVED"
POStatusCancelled         PurchaseOrderStatus = "CANCELLED"
POStatusClosed            PurchaseOrderStatus = "CLOSED"
)

var (
ErrInvalidPOStatusTransition = errors.New("transición de estado de orden de compra inválida")
ErrPOLineInvalidQuantity     = errors.New("la cantidad ordenada debe ser mayor a cero")
)

func (s PurchaseOrderStatus) IsValid() error {
switch s {
case POStatusDraft, POStatusSubmitted, POStatusApproved, POStatusPartiallyReceived, POStatusReceived, POStatusCancelled, POStatusClosed:
return nil
}
return errors.New("estado de orden de compra desconocido")
}