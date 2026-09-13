package inventory_test

import (
	"fortiasaass-kio-speed/internal/inventory"
	"testing"
)

type mockDB struct{}

func (m *mockDB) Exec(query string, args ...interface{}) (interface{}, error) {
	return nil, nil
}

func (m *mockDB) Query(query string, args ...interface{}) ([]inventory.Batch, error) {
	return []inventory.Batch{
		{ID: 1, ProductID: 101, LotNumber: "L001"},
	}, nil
}

func TestGetStocks(t *testing.T) {
	repo := &inventory.Repository{} // aquí deberías inyectar mockDB si tu struct lo permite

	stocks, err := repo.GetStocks()
	if err != nil {
		t.Fatalf("Error inesperado: %v", err)
	}
	if len(stocks) == 0 {
		t.Errorf("Se esperaba al menos un lote")
	}
}
