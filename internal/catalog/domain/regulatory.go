package domain

import (
"errors"
"time"
)

// Manufacturer Representa al laboratorio fabricante o planta productiva autorizada
type Manufacturer struct {
ManufacturerID string
Name           string
Country        string
Address        string
}

func NewManufacturer(id, name, country, address string) (Manufacturer, error) {
if id == "" || name == "" {
return Manufacturer{}, errors.New("el ID y nombre del fabricante son obligatorios")
}
return Manufacturer{
ManufacturerID: id,
Name:           name,
Country:        country,
Address:        address,
}, nil
}

// MarketingAuthorizationHolder (MAH) Titular del Registro Sanitario ante la autoridad regulatoria
type MarketingAuthorizationHolder struct {
HolderID string
Name     string
Country  string
}

func NewMAH(id, name, country string) (MarketingAuthorizationHolder, error) {
if id == "" || name == "" {
return MarketingAuthorizationHolder{}, errors.New("el ID y nombre del titular MAH son obligatorios")
}
return MarketingAuthorizationHolder{
HolderID: id,
Name:     name,
Country:  country,
}, nil
}

// TherapeuticEquivalenceType Estado de la equivalencia terapéutica (CAT-010)
type TherapeuticEquivalenceType string

const (
EquivalenceProven   TherapeuticEquivalenceType = "PROVEN_EQUIVALENT" // Demostrada formalmente
EquivalenceReference TherapeuticEquivalenceType = "REFERENCE_PRODUCT"  // Producto innovador de referencia
EquivalencePending  TherapeuticEquivalenceType = "PENDING"
)

// TherapeuticEquivalence Entidad regulatoria explícita de intercambiabilidad (CAT-010)
type TherapeuticEquivalence struct {
TenantID            string
ProductID           string
ComparatorProductID string // Producto innovador de referencia con el que se compara
Status              TherapeuticEquivalenceType
EvidenceReference   string // Resolución o estudio ISP que avala la equivalencia
ValidFrom           time.Time
ValidTo             time.Time
}

func NewTherapeuticEquivalence(
tenantID, productID, comparatorID string,
status TherapeuticEquivalenceType,
evidence string,
from, to time.Time,
) (TherapeuticEquivalence, error) {
if tenantID == "" || productID == "" {
return TherapeuticEquivalence{}, errors.New("parámetros obligatorios vacíos para la equivalencia terapéutica")
}
if to.Before(from) {
return TherapeuticEquivalence{}, errors.New("la vigencia de la equivalencia terapéutica no puede terminar antes de comenzar")
}

return TherapeuticEquivalence{
TenantID:            tenantID,
ProductID:           productID,
ComparatorProductID: comparatorID,
Status:              status,
EvidenceReference:   evidence,
ValidFrom:           from,
ValidTo:             to,
}, nil
}

// PackageConfiguration Configura las jerarquías de empaque y reglas de fraccionamiento autorizado
type PackageConfiguration struct {
PackagingLevel    string // Ej: "Primary" (Blíster), "Secondary" (Caja), "Tertiary" (Master Carton)
ItemsPerPackage   int64  // Cantidad de subunidades contenidas
AllowsFranchising bool   // Si permite dispensación fraccionada (ej. venta de comprimidos sueltos)
}

func NewPackageConfiguration(level string, items int64, allowsFranchising bool) (PackageConfiguration, error) {
if level == "" || items <= 0 {
return PackageConfiguration{}, errors.New("nivel de empaque inválido o cantidad por empaque menor o igual a cero")
}
return PackageConfiguration{
PackagingLevel:    level,
ItemsPerPackage:   items,
AllowsFranchising: allowsFranchising,
}, nil
}