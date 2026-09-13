package postgres

import (
"context"
"database/sql"
"errors"
"fmt"

"fortiasaass-kio-speed/internal/purchasing/domain"
"fortiasaass-kio-speed/internal/purchasing/ports"
shared "fortiasaass-kio-speed/internal/shared/domain"
)

var (
ErrSupplierNotFound = errors.New("proveedor no encontrado")
ErrConcurrencyConflict = errors.New("conflicto de concurrencia: el proveedor fue modificado por otro usuario")
)

type supplierRepository struct {
db *sql.DB
}

func NewSupplierRepository(db *sql.DB) ports.SupplierRepository {
return &supplierRepository{db: db}
}

func (r *supplierRepository) Create(ctx context.Context, s *domain.Supplier) error {
tx, err := r.db.BeginTx(ctx, nil)
if err != nil {
return err
}
defer tx.Rollback()

// 1. Insertar la raíz del agregado
qSupplier := `
INSERT INTO suppliers (id, tenant_id, legal_name, trade_name, tax_id, supplier_type, status, payment_terms, created_at, updated_at, version)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

_, err = tx.ExecContext(ctx, qSupplier,
s.ID, s.TenantID, s.LegalName, s.TradeName, s.TaxID.String(), s.SupplierType, s.Status, s.PaymentTerms, s.CreatedAt, s.UpdatedAt, s.Version)
if err != nil {
return fmt.Errorf("error insertando supplier: %w", err)
}

// 2. Insertar direcciones
qAddress := `INSERT INTO supplier_addresses (tenant_id, supplier_id, street, city, region, country, postal_code) VALUES ($1, $2, $3, $4, $5, $6, $7)`
for _, addr := range s.Addresses {
_, err = tx.ExecContext(ctx, qAddress, s.TenantID, s.ID, addr.Street, addr.City, addr.Region, addr.Country, addr.PostalCode)
if err != nil {
return fmt.Errorf("error insertando address: %w", err)
}
}

// 3. Insertar contactos
qContact := `INSERT INTO supplier_contacts (tenant_id, supplier_id, name, email, phone) VALUES ($1, $2, $3, $4, $5)`
for _, contact := range s.Contacts {
_, err = tx.ExecContext(ctx, qContact, s.TenantID, s.ID, contact.Name, contact.Email, contact.Phone)
if err != nil {
return fmt.Errorf("error insertando contact: %w", err)
}
}

return tx.Commit()
}

func (r *supplierRepository) Update(ctx context.Context, s *domain.Supplier) error {
tx, err := r.db.BeginTx(ctx, nil)
if err != nil {
return err
}
defer tx.Rollback()

// 1. Actualizar la raíz del agregado (Concurrencia Optimista)
qUpdate := `
UPDATE suppliers 
SET legal_name = $1, trade_name = $2, tax_id = $3, supplier_type = $4, status = $5, payment_terms = $6, updated_at = $7, version = version + 1
WHERE id = $8 AND tenant_id = $9 AND version = $10`

res, err := tx.ExecContext(ctx, qUpdate,
s.LegalName, s.TradeName, s.TaxID.String(), s.SupplierType, s.Status, s.PaymentTerms, s.UpdatedAt, s.ID, s.TenantID, s.Version)
if err != nil {
return err
}

rowsAffected, err := res.RowsAffected()
if err != nil {
return err
}
if rowsAffected == 0 {
return ErrConcurrencyConflict
}

// 2. Reemplazar direcciones (Estrategia sencilla para colecciones pequeñas: Delete & Insert)
_, err = tx.ExecContext(ctx, `DELETE FROM supplier_addresses WHERE tenant_id = $1 AND supplier_id = $2`, s.TenantID, s.ID)
if err != nil {
return err
}
qAddress := `INSERT INTO supplier_addresses (tenant_id, supplier_id, street, city, region, country, postal_code) VALUES ($1, $2, $3, $4, $5, $6, $7)`
for _, addr := range s.Addresses {
_, err = tx.ExecContext(ctx, qAddress, s.TenantID, s.ID, addr.Street, addr.City, addr.Region, addr.Country, addr.PostalCode)
if err != nil {
return err
}
}

// 3. Reemplazar contactos
_, err = tx.ExecContext(ctx, `DELETE FROM supplier_contacts WHERE tenant_id = $1 AND supplier_id = $2`, s.TenantID, s.ID)
if err != nil {
return err
}
qContact := `INSERT INTO supplier_contacts (tenant_id, supplier_id, name, email, phone) VALUES ($1, $2, $3, $4, $5)`
for _, contact := range s.Contacts {
_, err = tx.ExecContext(ctx, qContact, s.TenantID, s.ID, contact.Name, contact.Email, contact.Phone)
if err != nil {
return err
}
}

// Actualizamos la versión en memoria si el commit es exitoso
err = tx.Commit()
if err == nil {
s.Version++
}
return err
}

func (r *supplierRepository) GetByID(ctx context.Context, tenantID domain.TenantID, id domain.SupplierID) (*domain.Supplier, error) {
// 1. Obtener la raíz
qSupplier := `
SELECT legal_name, trade_name, tax_id, supplier_type, status, payment_terms, created_at, updated_at, version
FROM suppliers WHERE tenant_id = $1 AND id = $2`

s := &domain.Supplier{ID: id, TenantID: tenantID}
var taxIDStr, suppType, status, payTerms string

err := r.db.QueryRowContext(ctx, qSupplier, tenantID, id).Scan(
&s.LegalName, &s.TradeName, &taxIDStr, &suppType, &status, &payTerms, &s.CreatedAt, &s.UpdatedAt, &s.Version,
)
if err != nil {
if errors.Is(err, sql.ErrNoRows) {
return nil, ErrSupplierNotFound
}
return nil, err
}

// Reconstruir Value Objects
s.TaxID, _ = domain.NewTaxIdentifier(taxIDStr)
s.SupplierType = domain.SupplierType(suppType)
s.Status = domain.SupplierStatus(status)
s.PaymentTerms = domain.PaymentTerms(payTerms)

// 2. Obtener direcciones
s.Addresses = make([]shared.Address, 0)
qAddrs := `SELECT street, city, region, country, postal_code FROM supplier_addresses WHERE tenant_id = $1 AND supplier_id = $2`
rowsAddrs, err := r.db.QueryContext(ctx, qAddrs, tenantID, id)
if err == nil {
defer rowsAddrs.Close()
for rowsAddrs.Next() {
var a shared.Address
if err := rowsAddrs.Scan(&a.Street, &a.City, &a.Region, &a.Country, &a.PostalCode); err == nil {
s.Addresses = append(s.Addresses, a)
}
}
}

// 3. Obtener contactos
s.Contacts = make([]shared.Contact, 0)
qContacts := `SELECT name, email, phone FROM supplier_contacts WHERE tenant_id = $1 AND supplier_id = $2`
rowsContacts, err := r.db.QueryContext(ctx, qContacts, tenantID, id)
if err == nil {
defer rowsContacts.Close()
for rowsContacts.Next() {
var c shared.Contact
if err := rowsContacts.Scan(&c.Name, &c.Email, &c.Phone); err == nil {
s.Contacts = append(s.Contacts, c)
}
}
}

return s, nil
}

func (r *supplierRepository) ListByTenant(ctx context.Context, tenantID domain.TenantID) ([]*domain.Supplier, error) {
// Implementación base para listado (retorna agregados sin relaciones hijas por rendimiento)
// En un escenario real, se usaría un modelo de lectura (CQRS) o paginación.
q := `SELECT id, legal_name, tax_id, status FROM suppliers WHERE tenant_id = $1`
rows, err := r.db.QueryContext(ctx, q, tenantID)
if err != nil {
return nil, err
}
defer rows.Close()

var result []*domain.Supplier
for rows.Next() {
var id, legal, tax, status string
if err := rows.Scan(&id, &legal, &tax, &status); err == nil {
taxID, _ := domain.NewTaxIdentifier(tax)
result = append(result, &domain.Supplier{
ID:        domain.SupplierID(id),
TenantID:  tenantID,
LegalName: legal,
TaxID:     taxID,
Status:    domain.SupplierStatus(status),
})
}
}
return result, nil
}