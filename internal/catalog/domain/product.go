package domain

import (
"errors"
"time"
)

// ProductType Define la naturaleza comercial y regulatoria del producto
type ProductType string

const (
TypeMedicinal     ProductType = "MEDICINAL"     // Medicamento (Requiere perfil ISP y registro sanitario)
TypeCosmetic      ProductType = "COSMETIC"      // Cosmético e higiene personal
TypeMedicalDevice ProductType = "MEDICAL_DEVICE"// Dispositivo médico e insumo clínico
TypeSupplement    ProductType = "SUPPLEMENT"    // Suplemento alimentario / vitamínico
TypeGeneralRetail ProductType = "GENERAL_RETAIL"// Aseo, pañales y otros productos de venta general
)

// ProductStatus Estado operacional y comercial del producto en el catálogo
type ProductStatus string

const (
StatusActive    ProductStatus = "ACTIVE"
StatusSuspended ProductStatus = "SUSPENDED"
StatusDiscontinued ProductStatus = "DISCONTINUED"
)

// Product Agregado Raíz del Catálogo Maestro
type Product struct {
TenantID     string
ProductID    string // Inmutable (CAT-001)
Sku          string
CommercialName string
Brand        string
Type         ProductType
Status       ProductStatus
CreatedAt    time.Time
UpdatedAt    time.Time
}

func NewProduct(tenantID, productID, sku, commercialName, brand string, pType ProductType) (Product, error) {
if tenantID == "" || productID == "" || commercialName == "" {
return Product{}, errors.New("parámetros obligatorios vacíos para la creación del producto")
}

now := time.Now().UTC()
return Product{
TenantID:       tenantID,
ProductID:      productID,
Sku:            sku,
CommercialName: commercialName,
Brand:          brand,
Type:           pType,
Status:         StatusActive,
CreatedAt:      now,
UpdatedAt:      now,
}, nil
}

func (p *Product) Suspend() {
p.Status = StatusSuspended
p.UpdatedAt = time.Now().UTC()
}

func (p *Product) Activate() {
p.Status = StatusActive
p.UpdatedAt = time.Now().UTC()
}