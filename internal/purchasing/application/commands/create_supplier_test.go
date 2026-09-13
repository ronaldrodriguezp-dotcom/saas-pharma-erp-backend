package commands

import (
"context"
"errors"
"testing"

"fortiasaass-kio-speed/internal/purchasing/domain"
)

// mockSupplierRepo simula la base de datos en memoria para probar el Handler
type mockSupplierRepo struct {
suppliers map[string]*domain.Supplier
}

func newMockSupplierRepo() *mockSupplierRepo {
return &mockSupplierRepo{
suppliers: make(map[string]*domain.Supplier),
}
}

func (m *mockSupplierRepo) Create(ctx context.Context, s *domain.Supplier) error {
// Simular la restricción de PostgreSQL: UNIQUE (tenant_id, tax_id)
for _, existing := range m.suppliers {
if existing.TenantID == s.TenantID && existing.TaxID.String() == s.TaxID.String() {
return errors.New("violación de restricción única: tenant_id + tax_id")
}
}
m.suppliers[string(s.ID)] = s
return nil
}

func (m *mockSupplierRepo) Update(ctx context.Context, s *domain.Supplier) error { return nil }
func (m *mockSupplierRepo) GetByID(ctx context.Context, tenantID domain.TenantID, id domain.SupplierID) (*domain.Supplier, error) { return nil, nil }
func (m *mockSupplierRepo) ListByTenant(ctx context.Context, tenantID domain.TenantID) ([]*domain.Supplier, error) { return nil, nil }


func TestCreateSupplierHandler_Success(t *testing.T) {
repo := newMockSupplierRepo()
handler := NewCreateSupplierHandler(repo)

cmd := CreateSupplierCommand{
ID:           "SUP-001",
TenantID:     "TENANT-A",
LegalName:    "Droguería Central SpA",
TradeName:    "Droguería Central",
TaxID:        "11111111-1", // RUT Válido
SupplierType: "DISTRIBUTOR",
PaymentTerms: "NET_30",
}

err := handler.Handle(context.Background(), cmd)
if err != nil {
t.Fatalf("No se esperaba error, se obtuvo: %v", err)
}

if len(repo.suppliers) != 1 {
t.Errorf("Se esperaba 1 proveedor en el repositorio, hay %d", len(repo.suppliers))
}
}

func TestCreateSupplierHandler_InvalidTaxID(t *testing.T) {
repo := newMockSupplierRepo()
handler := NewCreateSupplierHandler(repo)

cmd := CreateSupplierCommand{
ID:           "SUP-002",
TenantID:     "TENANT-A",
LegalName:    "Proveedor Falso",
TaxID:        "11111111-2", // RUT Inválido (dígito verificador incorrecto)
SupplierType: "DISTRIBUTOR",
PaymentTerms: "CASH",
}

err := handler.Handle(context.Background(), cmd)
if err == nil {
t.Fatal("Se esperaba error por TaxID inválido, pero pasó")
}
if len(repo.suppliers) != 0 {
t.Errorf("El repositorio no debería tener registros, se salvaguardó un proveedor inválido")
}
}

func TestCreateSupplierHandler_DuplicateTaxID_SameTenant(t *testing.T) {
repo := newMockSupplierRepo()
handler := NewCreateSupplierHandler(repo)

cmd1 := CreateSupplierCommand{
ID:           "SUP-001",
TenantID:     "TENANT-A",
LegalName:    "Droguería 1",
TaxID:        "11111111-1",
SupplierType: "DISTRIBUTOR",
PaymentTerms: "CASH",
}
_ = handler.Handle(context.Background(), cmd1)

// Intento de crear otro proveedor con el MISMO RUT en el MISMO Tenant
cmd2 := CreateSupplierCommand{
ID:           "SUP-002",
TenantID:     "TENANT-A",
LegalName:    "Droguería 2 (Clon)",
TaxID:        "11111111-1",
SupplierType: "DISTRIBUTOR",
PaymentTerms: "CASH",
}

err := handler.Handle(context.Background(), cmd2)
if err == nil {
t.Fatal("Se esperaba error por TaxID duplicado en el mismo tenant")
}
}