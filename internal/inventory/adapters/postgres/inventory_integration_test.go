package postgres_test

import (
"context"
"database/sql"
"os"
"testing"
"time"

_ "github.com/lib/pq"
"fortiasaass-kio-speed/internal/inventory/adapters/postgres"
"fortiasaass-kio-speed/internal/inventory/application"
)

func TestInventoryIntegration_FEFO_FIFO(t *testing.T) {
dbURL := os.Getenv("TEST_DATABASE_URL")
if dbURL == "" {
t.Skip("Saltando prueba de integración de inventario: TEST_DATABASE_URL no está configurada")
}

db, err := sql.Open("postgres", dbURL)
if err != nil {
t.Fatalf("no se pudo conectar a la base de datos de prueba: %v", err)
}
defer db.Close()

ctx := context.Background()
repo := postgres.NewInventoryRepository(db)
svc := application.NewDischargeService(repo)

tenantID := "tenant-test-rengo-pharmacy"
warehouseID := "wh-central"
presentationID := "prod-amoxicilina-500mg"

now := time.Now()
// Lote A: Vence más tarde, recibido primero
lotAID := "lot-A-" + now.Format("150405")
// Lote B: Vence antes (CRÍTICO para FEFO), recibido después
lotBID := "lot-B-" + now.Format("150405")

lines := []postgres.ReceiptLineData{
{
PresentationID: presentationID,
LotNumber:      "LOT-LATER-EXP",
ExpirationDate: now.AddDate(0, 6, 0), // 6 meses más
Quantity:       50,
LocationID:     warehouseID,
},
}

// Simulamos la recepción de mercancía para el Lote A
err = repo.ReceiveGoods(ctx, tenantID, lotAID, lines)
if err != nil {
t.Fatalf("error insertando lote A: %v", err)
}

// Insertamos Lote B con fecha de expiración más próxima (FEFO debe seleccionarlo primero)
linesB := []postgres.ReceiptLineData{
{
PresentationID: presentationID,
LotNumber:      "LOT-SOON-EXP",
ExpirationDate: now.AddDate(0, 1, 0), // 1 mes más (vence antes)
Quantity:       30,
LocationID:     warehouseID,
},
}
err = repo.ReceiveGoods(ctx, tenantID, lotBID, linesB)
if err != nil {
t.Fatalf("error insertando lote B: %v", err)
}

// Ejecutamos una descarga bajo política FEFO (IsPerishable = true) por una cantidad de 10 unidades.
// Al ser FEFO, debe descontar prioritariamente del Lote B (que vence en 1 mes).
cmd := application.DischargeCommand{
TenantID:     tenantID,
WarehouseID:  warehouseID,
IsPerishable: true,
Items: []application.DischargeItem{
{
PresentationID: presentationID,
Quantity:       10,
},
},
}

err = svc.DischargeStock(ctx, cmd)
if err != nil {
t.Fatalf("error ejecutando despacho FEFO: %v", err)
}

// Verificamos lotes disponibles para comprobar que el descuento se aplicó correctamente al lote correcto
lots, err := repo.GetAvailableLots(ctx, tenantID, warehouseID, presentationID)
if err != nil {
t.Fatalf("error obteniendo lotes post-descarga: %v", err)
}

foundLotBReduced := false
for _, lot := range lots {
if lot.ID == "lot-"+lotBID+"-0" && lot.AvailableQty == 20 { // Empezó con 30, se descontaron 10
foundLotBReduced = true
}
}

if !foundLotBReduced {
t.Errorf("la estrategia FEFO no priorizó correctamente el lote más próximo a vencer")
}
}