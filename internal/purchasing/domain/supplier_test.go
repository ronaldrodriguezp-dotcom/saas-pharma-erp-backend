package domain

import (
"testing"

shared "fortiasaass-kio-speed/internal/shared/domain"
)

func TestSupplier_Creation(t *testing.T) {
taxID, _ := NewTaxIdentifier("11111111-1")

supplier, err := NewSupplier(
SupplierID("SUP-001"),
TenantID("TENANT-1"),
"Laboratorios Ficticios S.A.",
"Lab Ficticio",
taxID,
SupplierTypeLaboratory,
PaymentTermsNet30,
)

if err != nil {
t.Fatalf("No se esperaba error al crear supplier, obtuvo: %v", err)
}

if supplier.Status != SupplierStatusActive {
t.Errorf("Se esperaba estado inicial ACTIVE, obtuvo: %s", supplier.Status)
}
}

func TestSupplier_ChangeStatus(t *testing.T) {
taxID, _ := NewTaxIdentifier("11111111-1")
supplier, _ := NewSupplier("SUP-001", "TENANT-1", "Empresa", "", taxID, SupplierTypeDistributor, PaymentTermsCash)

// Transición válida
err := supplier.ChangeStatus(SupplierStatusSuspended)
if err != nil || supplier.Status != SupplierStatusSuspended {
t.Errorf("Falló la transición a SUSPENDED")
}

// Transición inválida
err = supplier.ChangeStatus(SupplierStatus("INVENTADO"))
if err == nil {
t.Errorf("Se esperaba error al cambiar a estado inválido")
}
}

func TestSupplier_IsEligibleForPurchasing(t *testing.T) {
taxID, _ := NewTaxIdentifier("11111111-1")
supplier, _ := NewSupplier("SUP-001", "TENANT-1", "Empresa", "", taxID, SupplierTypeDistributor, PaymentTermsCash)

// Activo -> Elegible
if err := supplier.IsEligibleForPurchasing(); err != nil {
t.Errorf("El proveedor ACTIVE debería ser elegible")
}

// Inactivo -> No Elegible
supplier.ChangeStatus(SupplierStatusBlocked)
if err := supplier.IsEligibleForPurchasing(); err != ErrSupplierInactive {
t.Errorf("El proveedor BLOCKED NO debería ser elegible")
}
}

func TestSupplier_AddAddressAndContact(t *testing.T) {
taxID, _ := NewTaxIdentifier("11111111-1")
supplier, _ := NewSupplier("SUP-001", "TENANT-1", "Empresa", "", taxID, SupplierTypeDistributor, PaymentTermsCash)

addr, _ := shared.NewAddress("Calle 1", "Santiago", "RM", "Chile")
contact, _ := shared.NewContact("Juan", "juan@test.com", "")

supplier.AddAddress(addr)
supplier.AddContact(contact)

if len(supplier.Addresses) != 1 || len(supplier.Contacts) != 1 {
t.Errorf("No se agregaron correctamente las direcciones o contactos")
}
}