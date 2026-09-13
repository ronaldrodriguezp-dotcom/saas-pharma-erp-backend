package domain

import (
"testing"
"time"
)

func TestGoodsReceipt_Invariants(t *testing.T) {
gr, err := NewGoodsReceipt("GR-1", "TENANT-1", "PO-100", "SUP-1", "WH-1", DocTypeDeliveryNote, "12345", "USER-RECEIVER")
if err != nil {
t.Fatalf("Fallo al crear GoodsReceipt: %v", err)
}

futureDate := time.Now().AddDate(1, 0, 0) // 1 año en el futuro
pastDate := time.Now().AddDate(-1, 0, 0)  // 1 año en el pasado

// 1. Validar línea con matemáticas incorrectas
err = gr.AddLine("L1", "POL-1", "PRES-1", "LOTE-A", futureDate, 10, 5, 2, "Dañado") // 5+2 != 10
if err != ErrReceiptLineInvalidQty {
t.Errorf("Se esperaba error por cantidades inválidas, obtuvo: %v", err)
}

// 2. Validar aceptación de lote vencido
err = gr.AddLine("L2", "POL-1", "PRES-1", "LOTE-B", pastDate, 10, 10, 0, "")
if err != ErrReceiptExpiredLot {
t.Errorf("Se esperaba rechazo automático por lote vencido, obtuvo: %v", err)
}

// 3. Añadir línea válida
err = gr.AddLine("L3", "POL-1", "PRES-1", "LOTE-C", futureDate, 10, 8, 2, "Cajas abolladas")
if err != nil {
t.Errorf("Falló al agregar línea válida: %v", err)
}

// 4. Confirmar recepción
if err := gr.Confirm(); err != nil {
t.Errorf("La confirmación debió ser exitosa, falló: %v", err)
}

// 5. Modificar post-confirmación debe fallar
err = gr.AddLine("L4", "POL-1", "PRES-2", "LOTE-D", futureDate, 5, 5, 0, "")
if err == nil {
t.Errorf("Se esperaba error al intentar agregar líneas en estado CONFIRMED")
}
}