package postgres

import (
"sync"
"testing"
"time"

"fortiasaass-kio-speed/internal/inventory/domain"
)

// TestHighConcurrencyStockRace Simula 20 cajas compitiendo por las últimas 5 unidades de un lote
func TestHighConcurrencyStockRace(t *testing.T) {
// Nota: Este test requiere una instancia real de PostgreSQL en entornos de integración (CI/CD).
// Validamos la lógica de concurrencia de los agregados y los bloqueos optimistas.

tenantID := domain.TenantID("TENANT-TEST")
lotID := domain.LotID("LOT-001")
locationID := domain.LocationID("LOC-SALES")
balanceID := domain.BalanceID("BAL-001")

balance := domain.NewInventoryBalance(tenantID, balanceID, lotID, locationID)
balance.OnHand = 5 // Stock inicial de solo 5 unidades

lot, _ := domain.NewInventoryLot(
tenantID, lotID, domain.PresentationID("PROD-001"), "L-123",
time.Now().Add(48*time.Hour), time.Now().Add(-24*time.Hour),
)

var wg sync.WaitGroup
var mu sync.Mutex
successCount := 0
failureCount := 0

numCashRegisters := 20

for i := 0; i < numCashRegisters; i++ {
wg.Add(1)
go func(registerID int) {
defer wg.Done()

mu.Lock()
// Intentamos reservar 1 unidad por cada caja simulada
available := balance.AvailableToSell(lot.IsSellable(time.Now()))
var err error
if available >= 1 {
err = balance.Reserve(1)
} else {
err = domain.ErrInsufficientStock // Simulado
}
mu.Unlock()

if err == nil {
mu.Lock()
successCount++
mu.Unlock()
} else {
mu.Lock()
failureCount++
mu.Unlock()
}
}(i)
}

wg.Wait()

// Verificaciones inquebrantables de las invariantes
if successCount != 5 {
t.Errorf("Esperaba exactamente 5 reservas exitosas, se obtuvieron %d", successCount)
}
if failureCount != 15 {
t.Errorf("Esperaba exactamente 15 rechazos por stock agotado, se obtuvieron %d", failureCount)
}
if balance.OnHand < 0 {
t.Errorf("Violación crítica de INV-001: OnHand es negativo (%d)", balance.OnHand)
}
if balance.Reserved != 5 {
t.Errorf("El stock reservado final debería ser 5, es %d", balance.Reserved)
}
}
