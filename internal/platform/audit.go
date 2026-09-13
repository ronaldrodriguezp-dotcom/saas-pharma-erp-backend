package platform

import (
"log"
"time"
)

// AuditEvent define la estructura inmutable para transacciones críticas de farmacia y retail
type AuditEvent struct {
Timestamp   time.Time `json:"timestamp"`
TenantID    string    `json:"tenant_id"`
BranchID    string    `json:"branch_id"`
UserRUT     string    `json:"user_rut"`
Action      string    `json:"action"`
Details     string    `json:"details"`
}

// LogCriticalAction registra de forma estructurada e inalterable los eventos críticos
func LogCriticalAction(tenantID, branchID, userRUT, action, details string) {
event := AuditEvent{
Timestamp:   time.Now().UTC(),
TenantID:    tenantID,
BranchID:    branchID,
UserRUT:     userRUT,
Action:      action,
Details:     details,
}

// En producción, esto se escribe en una tabla de base de datos de "solo escritura" (Append-Only)
log.Printf("🔒 [AUDIT LOG INMUTABLE] [%s] Tenant: %s | Sucursal: %s | Usuario: %s | Acción: %s | Detalle: %s",
event.Timestamp.Format(time.RFC3339),
event.TenantID,
event.BranchID,
event.UserRUT,
event.Action,
event.Details,
)
}