package domain

import (
"errors"
"time"
)

// BarcodeType Define el estándar del código de barras asociado a la presentación
type BarcodeType string

const (
BarcodeGTIN13 BarcodeType = "GTIN-13"
BarcodeGTIN14 BarcodeType = "GTIN-14"
BarcodeCode128 BarcodeType = "CODE-128"
BarcodeInternal BarcodeType = "INTERNAL"
)

// Barcode Entidad subordinada que modela los códigos de barras de una presentación (CAT-002, CAT-003)
type Barcode struct {
Code      string
Type      BarcodeType
IsPrimary bool
}

// ProductPresentation Entidad Central Transaccional del Catálogo
// El inventario, los lotes, las reservas, la dispensación y las ventas apuntan a esta entidad.
type ProductPresentation struct {
TenantID          string
PresentationID    string
ProductID         string // Vinculación obligatoria al Product padre (CAT-001)
Name              string // Ej. "Caja x 30 comprimidos", "Frasco x 100 mL"
BaseUnitsQuantity int64  // Cantidad total de unidades base contenidas (ej. 30 comprimidos o 100 mL)
IsFranchisable    bool   // Indica si permite dispensación/venta fraccionada
Barcodes          []Barcode
IsActive          bool
CreatedAt         time.Time
}

func NewProductPresentation(
tenantID string,
presentationID string,
productID string,
name string,
baseUnitsQuantity int64,
isFranchisable bool,
) (ProductPresentation, error) {
if tenantID == "" || presentationID == "" || productID == "" || name == "" {
return ProductPresentation{}, errors.New("parámetros obligatorios vacíos para crear la presentación del producto")
}
if baseUnitsQuantity <= 0 {
return ProductPresentation{}, errors.New("la cantidad de unidades base en la presentación debe ser mayor a cero")
}

return ProductPresentation{
TenantID:          tenantID,
PresentationID:    presentationID,
ProductID:         productID,
Name:              name,
BaseUnitsQuantity: baseUnitsQuantity,
IsFranchisable:    isFranchisable,
Barcodes:          []Barcode{},
IsActive:          true,
CreatedAt:         time.Now().UTC(),
}, nil
}

// AddBarcode Agrega un código de barras garantizando la subordinación (CAT-002)
func (p *ProductPresentation) AddBarcode(code string, bType BarcodeType, isPrimary bool) error {
if code == "" {
return errors.New("el código de barras no puede estar vacío")
}

// Si se define como primario, desactivamos el flag en los demás
if isPrimary {
for i := range p.Barcodes {
p.Barcodes[i].IsPrimary = false
}
}

p.Barcodes = append(p.Barcodes, Barcode{
Code:      code,
Type:      bType,
IsPrimary: isPrimary,
})
return nil
}