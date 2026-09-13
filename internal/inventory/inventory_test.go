package tests

import (
	"backend/internal/inventory"
	"testing"
)

func TestCreateProduct(t *testing.T) {
	repo := inventory.NewRepository(mockDB())
	err := repo.CreateProduct("Paracetamol 500mg", "Analgesico")
	if err != nil {
		t.Errorf("Error al crear producto: %v", err)
	}
}
