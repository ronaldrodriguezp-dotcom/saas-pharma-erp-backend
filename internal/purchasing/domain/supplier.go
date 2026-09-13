package domain

import (
"errors"
"time"

shared "fortiasaass-kio-speed/internal/shared/domain"
)

var (
ErrSupplierInactive = errors.New("operación denegada: el proveedor no está activo")
)

type Supplier struct {
ID            SupplierID
TenantID      TenantID
LegalName     string
TradeName     string
TaxID         TaxIdentifier
SupplierType  SupplierType
Status        SupplierStatus
Addresses     []shared.Address
Contacts      []shared.Contact
PaymentTerms  PaymentTerms
CreatedAt     time.Time
UpdatedAt     time.Time
Version       int
}

// NewSupplier crea una nueva entidad Supplier con estado inicial ACTIVE
func NewSupplier(id SupplierID, tenant TenantID, legalName, tradeName string, taxID TaxIdentifier, suppType SupplierType, terms PaymentTerms) (*Supplier, error) {
if legalName == "" {
return nil, errors.New("el nombre legal es obligatorio")
}

return &Supplier{
ID:           id,
TenantID:     tenant,
LegalName:    legalName,
TradeName:    tradeName,
TaxID:        taxID,
SupplierType: suppType,
Status:       SupplierStatusActive, // Por defecto nace activo
Addresses:    make([]shared.Address, 0),
Contacts:     make([]shared.Contact, 0),
PaymentTerms: terms,
CreatedAt:    time.Now(),
UpdatedAt:    time.Now(),
Version:      1,
}, nil
}

// ChangeStatus encapsula la transición de estado del proveedor
func (s *Supplier) ChangeStatus(newStatus SupplierStatus) error {
if err := newStatus.IsValid(); err != nil {
return err
}
s.Status = newStatus
s.UpdatedAt = time.Now()
return nil
}

// AddAddress agrega una nueva dirección validada
func (s *Supplier) AddAddress(address shared.Address) {
s.Addresses = append(s.Addresses, address)
s.UpdatedAt = time.Now()
}

// AddContact agrega un nuevo contacto validado
func (s *Supplier) AddContact(contact shared.Contact) {
s.Contacts = append(s.Contacts, contact)
s.UpdatedAt = time.Now()
}

// IsEligibleForPurchasing verifica si se le pueden emitir órdenes de compra
func (s *Supplier) IsEligibleForPurchasing() error {
if s.Status != SupplierStatusActive {
return ErrSupplierInactive
}
return nil
}