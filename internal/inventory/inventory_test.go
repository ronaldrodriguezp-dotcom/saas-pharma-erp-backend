package inventory_test

import (
	"testing"
	"time"

	"fortiasaass-kio-speed/internal/inventory"
)

// fakeRepo simula el comportamiento de Repository para pruebas unitarias.
type fakeRepo struct{}

func (f *fakeRepo) GetStocks() ([]inventory.Batch, error) {
	return []inventory.Batch{
		{
			ID:             1,
			ProductID:      101,
			LotNumber:      "L001",
			ExpirationDate: time.Now().AddDate(0, 6, 0),
			Quantity:       50,
			CreatedAt:      time.Now(),
		},
		{
			ID:             2,
			ProductID:      102,
			LotNumber:      "L002",
			ExpirationDate: time.Now().AddDate(0, 3, 0),
			Quantity:       20,
			CreatedAt:      time.Now(),
		},
	}, nil
}

func TestGetStocksReturnsExpectedBatches(t *testing.T) {
	var repo interface {
		GetStocks() ([]inventory.Batch, error)
	} = &fakeRepo{}

	stocks, err := repo.GetStocks()
	if err != nil {
		t.Fatalf("GetStocks devolvió error inesperado: %v", err)
	}
	if len(stocks) != 2 {
		t.Fatalf("Se esperaban 2 lotes, se obtuvieron %d", len(stocks))
	}

	if stocks[0].LotNumber != "L001" {
		t.Errorf("Esperado LotNumber L001, obtenido %s", stocks[0].LotNumber)
	}
	if stocks[0].Quantity <= 0 {
		t.Errorf("Esperado quantity > 0, obtenido %d", stocks[0].Quantity)
	}
}
