package domain

import (
"errors"
"time"

shared "fortiasaass-kio-speed/internal/shared/domain"
)

var (
ErrPOEmptyLines       = errors.New("la orden de compra debe tener al menos una línea")
ErrPODuplicateProduct = errors.New("la misma presentación (PresentationID) ya existe en otra línea")
ErrPOSameActor        = errors.New("segregation of duties: el creador de la orden no puede ser quien la apruebe")
)

type PurchaseOrderLine struct {
ID             string
PresentationID string 
OrderedQty     int    
UnitCost       shared.Money
}

type PurchaseOrder struct {
ID            PurchaseOrderID
TenantID      TenantID
SupplierID    SupplierID
DestinationID LocationID
Status        PurchaseOrderStatus
Lines         []PurchaseOrderLine
CreatedBy     string 
ApprovedBy    string 
CreatedAt     time.Time
UpdatedAt     time.Time
Version       int
}

func NewPurchaseOrder(id PurchaseOrderID, tenant TenantID, supplier SupplierID, dest LocationID, creator string) (*PurchaseOrder, error) {
if string(id) == "" || string(tenant) == "" || string(supplier) == "" || string(dest) == "" || creator == "" {
return nil, errors.New("identificadores incompletos para crear la orden de compra")
}

return &PurchaseOrder{
ID:            id,
TenantID:      tenant,
SupplierID:    supplier,
DestinationID: dest,
Status:        POStatusDraft,
Lines:         make([]PurchaseOrderLine, 0),
CreatedBy:     creator,
CreatedAt:     time.Now(),
UpdatedAt:     time.Now(),
Version:       1,
}, nil
}

func (po *PurchaseOrder) AddLine(lineID, presentationID string, qty int, cost shared.Money) error {
if po.Status != POStatusDraft {
return errors.New("solo se pueden agregar líneas a una orden en estado DRAFT")
}
if qty <= 0 {
return ErrPOLineInvalidQuantity
}

for _, line := range po.Lines {
if line.PresentationID == presentationID {
return ErrPODuplicateProduct
}
}

po.Lines = append(po.Lines, PurchaseOrderLine{
ID:             lineID,
PresentationID: presentationID,
OrderedQty:     qty,
UnitCost:       cost,
})
po.UpdatedAt = time.Now()
return nil
}

func (po *PurchaseOrder) Submit() error {
if po.Status != POStatusDraft {
return ErrInvalidPOStatusTransition
}
if len(po.Lines) == 0 {
return ErrPOEmptyLines
}
po.Status = POStatusSubmitted
po.UpdatedAt = time.Now()
return nil
}

// Approve aplica la regla de Segregation of Duties (SoD)
func (po *PurchaseOrder) Approve(approver string) error {
if po.Status != POStatusSubmitted {
return ErrInvalidPOStatusTransition
}
if approver == "" {
return errors.New("el aprobador es obligatorio")
}
if po.CreatedBy == approver {
return ErrPOSameActor
}
po.Status = POStatusApproved
po.ApprovedBy = approver
po.UpdatedAt = time.Now()
return nil
}