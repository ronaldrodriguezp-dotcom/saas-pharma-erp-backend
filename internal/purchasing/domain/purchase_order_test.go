package domain

import (
"testing"

shared "fortiasaass-kio-speed/internal/shared/domain"
)

func TestPurchaseOrder_Invariants(t *testing.T) {
po, err := NewPurchaseOrder("PO-100", "TENANT-1", "SUP-1", "LOC-1", "USER-CREATOR")
if err != nil {
t.Fatalf("Fallo al crear PO: %v", err)
}

// 1. Submit sin líneas debe fallar
if err := po.Submit(); err != ErrPOEmptyLines {
t.Errorf("Se esperaba error por falta de líneas al hacer submit, obtuvo: %v", err)
}

// 2. Agregar línea válida
money, _ := shared.NewCostMoney(15000, shared.CLP)
err = po.AddLine("LINE-1", "PRES-001", 10, money)
if err != nil {
t.Errorf("Fallo al agregar línea: %v", err)
}

// 3. Prevenir productos duplicados en la misma OC
err = po.AddLine("LINE-2", "PRES-001", 5, money)
if err != ErrPODuplicateProduct {
t.Errorf("Se esperaba error por producto duplicado, obtuvo: %v", err)
}

// 4. Submit exitoso
if err := po.Submit(); err != nil {
t.Errorf("Submit debería haber sido exitoso, obtuvo: %v", err)
}

// 5. Agregar línea después de Submit debe fallar
err = po.AddLine("LINE-3", "PRES-002", 5, money)
if err == nil {
t.Error("No se deberían poder agregar líneas después del Submit")
}

// 6. Aprobar con el MISMO usuario creador debe fallar (Segregation of Duties)
err = po.Approve("USER-CREATOR")
if err != ErrPOSameActor {
t.Errorf("Se esperaba bloqueo por SoD, obtuvo: %v", err)
}

// 7. Aprobar con usuario distinto debe ser exitoso
err = po.Approve("USER-MANAGER")
if err != nil {
t.Errorf("Aprobación debió ser exitosa, falló con: %v", err)
}

if po.Status != POStatusApproved {
t.Errorf("El estado debió cambiar a APPROVED")
}
}