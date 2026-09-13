package domain

import (
"errors"
"regexp"
"strconv"
"strings"
)

type TenantID string
type SupplierID string

type TaxIdentifier struct {
value string
}

type SupplierStatus string

const (
SupplierStatusActive    SupplierStatus = "ACTIVE"
SupplierStatusSuspended SupplierStatus = "SUSPENDED"
SupplierStatusBlocked   SupplierStatus = "BLOCKED"
SupplierStatusInactive  SupplierStatus = "INACTIVE"
)

type SupplierType string

const (
SupplierTypeDistributor     SupplierType = "DISTRIBUTOR"
SupplierTypeLaboratory      SupplierType = "LABORATORY"
SupplierTypeWholesaler      SupplierType = "WHOLESALER"
SupplierTypeGeneralSupplier SupplierType = "GENERAL_SUPPLIER"
SupplierTypeServiceProvider SupplierType = "SERVICE_PROVIDER"
SupplierTypeOther           SupplierType = "OTHER"
)

type PaymentTerms string

const (
PaymentTermsCash  PaymentTerms = "CASH"
PaymentTermsNet30 PaymentTerms = "NET_30"
PaymentTermsNet60 PaymentTerms = "NET_60"
PaymentTermsNet90 PaymentTerms = "NET_90"
)

var (
ErrInvalidTaxIdentifier  = errors.New("RUT/TaxIdentifier inválido")
ErrInvalidSupplierStatus = errors.New("estado de proveedor inválido")
ErrInvalidSupplierType   = errors.New("tipo de proveedor inválido")
ErrInvalidPaymentTerms   = errors.New("términos de pago inválidos")
)

func NewTaxIdentifier(rut string) (TaxIdentifier, error) {
rut = strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(rut, ".", ""), "-", ""))
if rut == "" {
return TaxIdentifier{}, ErrInvalidTaxIdentifier
}

re := regexp.MustCompile(`^[0-9]+[0-9K]$`)
if !re.MatchString(rut) {
return TaxIdentifier{}, ErrInvalidTaxIdentifier
}

body := rut[:len(rut)-1]
dv := string(rut[len(rut)-1])

if !validateMod11(body, dv) {
return TaxIdentifier{}, ErrInvalidTaxIdentifier
}

return TaxIdentifier{value: rut}, nil
}

func (t TaxIdentifier) String() string {
return t.value
}

func validateMod11(body, expectedDV string) bool {
sum := 0
multiplier := 2
for i := len(body) - 1; i >= 0; i-- {
digit, _ := strconv.Atoi(string(body[i]))
sum += digit * multiplier
multiplier++
if multiplier > 7 {
multiplier = 2
}
}
remainder := 11 - (sum % 11)
calculatedDV := strconv.Itoa(remainder)
if remainder == 11 {
calculatedDV = "0"
} else if remainder == 10 {
calculatedDV = "K"
}
return calculatedDV == expectedDV
}

func (s SupplierStatus) IsValid() error {
switch s {
case SupplierStatusActive, SupplierStatusSuspended, SupplierStatusBlocked, SupplierStatusInactive:
return nil
}
return ErrInvalidSupplierStatus
}

func (t SupplierType) IsValid() error {
switch t {
case SupplierTypeDistributor, SupplierTypeLaboratory, SupplierTypeWholesaler, SupplierTypeGeneralSupplier, SupplierTypeServiceProvider, SupplierTypeOther:
return nil
}
return ErrInvalidSupplierType
}

func (p PaymentTerms) IsValid() error {
switch p {
case PaymentTermsCash, PaymentTermsNet30, PaymentTermsNet60, PaymentTermsNet90:
return nil
}
return ErrInvalidPaymentTerms
}