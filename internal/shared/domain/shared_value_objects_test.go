package domain

import (
"testing"
)

func TestMoney_Creation(t *testing.T) {
// Prueba exitosa
m, err := NewCostMoney(13333, CLP)
if err != nil || m.Amount != 13333 {
t.Errorf("Esperaba crear Money válido, obtuvo error: %v", err)
}

// Prueba monto negativo (debe fallar)
_, err = NewCostMoney(-100, USD)
if err == nil {
t.Error("Esperaba error al crear Money con monto negativo, pero fue exitoso")
}
}

func TestContact_EmailValidation(t *testing.T) {
// Email válido
_, err := NewContact("Juan", "juan@proveedor.cl", "+56912345678")
if err != nil {
t.Errorf("Esperaba email válido, falló con: %v", err)
}

// Email inválido
_, err = NewContact("Pedro", "pedro-proveedor.cl", "")
if err != ErrInvalidEmail {
t.Errorf("Esperaba ErrInvalidEmail, obtuvo: %v", err)
}
}