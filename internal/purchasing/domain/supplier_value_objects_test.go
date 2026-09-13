package domain

import (
"testing"
)

func TestTaxIdentifier_Validation(t *testing.T) {
tests := []struct {
name    string
rut     string
isValid bool
}{
{"RUT Válido estándar", "11111111-1", true}, // 11.111.111-1 es válido matemáticamente
{"RUT Válido con K", "12345678-5", true},    // 12.345.678-5 es válido matemáticamente
{"RUT Inválido dígito malo", "11111111-2", false},
{"RUT Formato basura", "ABCDEFGH-K", false},
{"RUT Vacío", "", false},
{"RUT Válido con puntos", "11.111.111-1", true},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
_, err := NewTaxIdentifier(tt.rut)
if tt.isValid && err != nil {
t.Errorf("Se esperaba que %s fuera válido, pero falló: %v", tt.rut, err)
}
if !tt.isValid && err == nil {
t.Errorf("Se esperaba que %s fuera INVÁLIDO, pero pasó la validación", tt.rut)
}
})
}
}

func TestEnums_Validation(t *testing.T) {
status := SupplierStatusActive
if err := status.IsValid(); err != nil {
t.Errorf("Estado ACTIVE debería ser válido")
}

badStatus := SupplierStatus("FALSO")
if err := badStatus.IsValid(); err == nil {
t.Errorf("Estado FALSO debería ser inválido")
}
}