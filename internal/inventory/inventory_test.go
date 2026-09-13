package inventory

import (
	"sync"
	"testing"
)

func TestConcurrentReservations(t *testing.T) {
	inv := NewInventory(5)
	wg := sync.WaitGroup{}
	successCount := 0
	mu := sync.Mutex{}

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if inv.TryReserve(1) {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	if successCount != 5 {
		t.Errorf("Esperaba 5 reservas exitosas, obtuve %d", successCount)
	}
	if inv.Available() != 0 {
		t.Errorf("Stock final esperado 0, obtuve %d", inv.Available())
	}
}
