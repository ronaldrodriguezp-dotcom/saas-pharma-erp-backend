package postgres_test

import (
"context"
"database/sql"
"os"
"testing"
"time"

_ "github.com/lib/pq"
"fortiasaass-kio-speed/internal/purchasing/adapters/postgres"
"fortiasaass-kio-speed/internal/purchasing/domain"
)

func TestGoodsReceiptIntegration(t *testing.T) {
dbURL := os.Getenv("TEST_DATABASE_URL")
if dbURL == "" {
t.Skip("Saltando prueba de integración: TEST_DATABASE_URL no está configurada")
}

db, err := sql.Open("postgres", dbURL)
if err != nil {
t.Fatalf("no se pudo conectar a la base de datos de prueba: %v", err)
}
defer db.Close()

ctx := context.Background()
repo := postgres.NewGoodsReceiptRepository(db)

tenantID := domain.TenantID("tenant-test-rengo-01")
receiptID := domain.GoodsReceiptID(fmtID("gr"))

gr := &domain.GoodsReceipt{
ID:              receiptID,
TenantID:        tenantID,
PurchaseOrderID: domain.PurchaseOrderID("po-001"),
SupplierID:      domain.SupplierID("supp-001"),
WarehouseID:     domain.LocationID("wh-01"),
DocumentType:    domain.DocumentType("INVOICE"),
DocumentNumber:  "FACT-998877",
Status:          domain.GoodsReceiptStatus("DRAFT"),
ReceivedBy:      "Farmacéutico Jefe",
ReceivedAt:      time.Now(),
Version:         1,
Lines: []domain.GoodsReceiptLine{
{
ID:               "1",
POLineID:         "pol-1",
PresentationID:   "prod-paracetamol-500mg",
LotNumber:        "LOT-2026-XYZ",
ExpirationDate:   time.Now().AddDate(1, 0, 0),
ReceivedQuantity: 100,
AcceptedQuantity: 100,
RejectedQuantity: 0,
RejectionReason:  "",
},
},
}

// 1. Probar Creación
err = repo.Create(ctx, gr)
if err != nil {
t.Fatalf("error creando GoodsReceipt en integración: %v", err)
}

// 2. Probar Recuperación asegurando aislamiento multitenant
fetched, err := repo.GetByID(ctx, tenantID, receiptID)
if err != nil {
t.Fatalf("error recuperando GoodsReceipt: %v", err)
}

if fetched.DocumentNumber != "FACT-998877" {
t.Errorf("esperaba número de documento FACT-998877, se obtuvo %s", fetched.DocumentNumber)
}

if len(fetched.Lines) != 1 {
t.Errorf("esperaba 1 línea de recepción, se obtuvieron %d", len(fetched.Lines))
}

// 3. Verificar aislamiento por tenant incorrecto
_, err = repo.GetByID(ctx, domain.TenantID("tenant-intruso"), receiptID)
if !errorsIsNotFound(err) {
t.Errorf("fallo de seguridad multitenant: el tenant intruso pudo ver la recepción")
}
}

func fmtID(prefix string) string {
return prefix + "-" + time.Now().Format("20060102150405")
}

func errorsIsNotFound(err error) bool {
return err != nil && (err == postgres.ErrGoodsReceiptNotFound || err.Error() == "recepción de mercancía no encontrada")
}