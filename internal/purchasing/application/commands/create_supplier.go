package commands

import (
"context"
"errors"
"fmt"

"fortiasaass-kio-speed/internal/purchasing/domain"
"fortiasaass-kio-speed/internal/purchasing/ports"
)

var (
ErrSupplierAlreadyExists = errors.New("ya existe un proveedor con ese identificador fiscal")
)

// CreateSupplierCommand es el DTO de entrada (Request)
type CreateSupplierCommand struct {
ID           string
TenantID     string
LegalName    string
TradeName    string
TaxID        string
SupplierType string
PaymentTerms string
}

type CreateSupplierHandler struct {
repo ports.SupplierRepository
}

func NewCreateSupplierHandler(repo ports.SupplierRepository) *CreateSupplierHandler {
return &CreateSupplierHandler{repo: repo}
}

func (h *CreateSupplierHandler) Handle(ctx context.Context, cmd CreateSupplierCommand) error {
// 1. Instanciar Value Objects
tenantID := domain.TenantID(cmd.TenantID)
supplierID := domain.SupplierID(cmd.ID)

taxID, err := domain.NewTaxIdentifier(cmd.TaxID)
if err != nil {
return fmt.Errorf("tax_id inválido: %w", err)
}

suppType := domain.SupplierType(cmd.SupplierType)
if err := suppType.IsValid(); err != nil {
return err
}

payTerms := domain.PaymentTerms(cmd.PaymentTerms)
if err := payTerms.IsValid(); err != nil {
return err
}

// 2. Comprobar unicidad de TaxID a nivel de aplicación (Opcional, pero buena práctica antes de golpear la DB)
// En este diseño, delegaremos la responsabilidad final de unicidad a la restricción UNIQUE de PostgreSQL 
// (uq_tenant_tax_id) en el repositorio para evitar condiciones de carrera en lectura.

// 3. Crear el agregado a través de su factory
supplier, err := domain.NewSupplier(
supplierID,
tenantID,
cmd.LegalName,
cmd.TradeName,
taxID,
suppType,
payTerms,
)
if err != nil {
return fmt.Errorf("error construyendo proveedor: %w", err)
}

// 4. Persistir
err = h.repo.Create(ctx, supplier)
if err != nil {
// Aquí idealmente mapearíamos el error de clave duplicada de Postgres a ErrSupplierAlreadyExists
return fmt.Errorf("error persistiendo proveedor: %w", err)
}

return nil
}