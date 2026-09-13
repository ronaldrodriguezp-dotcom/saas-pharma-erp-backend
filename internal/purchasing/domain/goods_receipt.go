package domain

import (
"errors"
"time"
)

type GoodsReceiptID string

type DocumentType string
const (
DocTypeInvoice        DocumentType = "INVOICE"
DocTypeDeliveryNote   DocumentType = "DELIVERY_NOTE" // Guía de Despacho
)

type GoodsReceiptStatus string
const (
ReceiptStatusDraft     GoodsReceiptStatus = "DRAFT"
ReceiptStatusConfirmed GoodsReceiptStatus = "CONFIRMED"
ReceiptStatusCancelled GoodsReceiptStatus = "CANCELLED"
)

var (
ErrReceiptEmptyLines     = errors.New("la recepción debe tener al menos una línea")
ErrReceiptLineInvalidQty = errors.New("la cantidad recibida no puede ser menor a la aceptada + rechazada")
ErrReceiptInvalidStatus  = errors.New("transición de estado de recepción inválida")
ErrReceiptExpiredLot     = errors.New("no se puede aceptar mercadería vencida")
)

// GoodsReceiptLine registra la recepción física a nivel de Lote
type GoodsReceiptLine struct {
ID               string
POLineID         string // Enlace a la línea de la orden de compra original
PresentationID   string
LotNumber        string
ExpirationDate   time.Time
ReceivedQuantity int
AcceptedQuantity int
RejectedQuantity int
RejectionReason  string
}

func (l *GoodsReceiptLine) Validate() error {
if l.AcceptedQuantity < 0 || l.RejectedQuantity < 0 {
return errors.New("las cantidades no pueden ser negativas")
}
if l.ReceivedQuantity != (l.AcceptedQuantity + l.RejectedQuantity) {
return ErrReceiptLineInvalidQty
}
if l.AcceptedQuantity > 0 && l.LotNumber == "" {
return errors.New("el número de lote es obligatorio para mercadería aceptada")
}
// Validar que no se acepte mercadería ya vencida (Truncamos a fecha sin hora para precisión)
today := time.Now().Truncate(24 * time.Hour)
expDate := l.ExpirationDate.Truncate(24 * time.Hour)
if l.AcceptedQuantity > 0 && expDate.Before(today) {
return ErrReceiptExpiredLot
}
return nil
}

// GoodsReceipt es el Aggregate Root para la entrada de mercadería
type GoodsReceipt struct {
ID              GoodsReceiptID
TenantID        TenantID
PurchaseOrderID PurchaseOrderID // Referencia a la OC aprobada
SupplierID      SupplierID
WarehouseID     LocationID
DocumentType    DocumentType
DocumentNumber  string
Status          GoodsReceiptStatus
Lines           []GoodsReceiptLine
ReceivedBy      string
ReceivedAt      time.Time
Version         int
}

func NewGoodsReceipt(id GoodsReceiptID, tenant TenantID, poID PurchaseOrderID, supp SupplierID, wh LocationID, docType DocumentType, docNum string, receiver string) (*GoodsReceipt, error) {
if string(id) == "" || string(tenant) == "" || string(poID) == "" || receiver == "" {
return nil, errors.New("identificadores incompletos para crear la recepción")
}

return &GoodsReceipt{
ID:              id,
TenantID:        tenant,
PurchaseOrderID: poID,
SupplierID:      supp,
WarehouseID:     wh,
DocumentType:    docType,
DocumentNumber:  docNum,
Status:          ReceiptStatusDraft,
Lines:           make([]GoodsReceiptLine, 0),
ReceivedBy:      receiver,
ReceivedAt:      time.Now(),
Version:         1,
}, nil
}

func (gr *GoodsReceipt) AddLine(id, poLineID, presentation, lot string, expDate time.Time, received, accepted, rejected int, rejectReason string) error {
if gr.Status != ReceiptStatusDraft {
return errors.New("solo se pueden agregar líneas a una recepción en estado DRAFT")
}

line := GoodsReceiptLine{
ID:               id,
POLineID:         poLineID,
PresentationID:   presentation,
LotNumber:        lot,
ExpirationDate:   expDate,
ReceivedQuantity: received,
AcceptedQuantity: accepted,
RejectedQuantity: rejected,
RejectionReason:  rejectReason,
}

if err := line.Validate(); err != nil {
return err
}

gr.Lines = append(gr.Lines, line)
return nil
}

func (gr *GoodsReceipt) Confirm() error {
if gr.Status != ReceiptStatusDraft {
return ErrReceiptInvalidStatus
}
if len(gr.Lines) == 0 {
return ErrReceiptEmptyLines
}
gr.Status = ReceiptStatusConfirmed
return nil
}