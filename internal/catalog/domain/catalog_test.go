package domain

import (
"testing"
)

func TestCatalogInvariants(t *testing.T) {
validator := CatalogInvariantValidator{}

// Prueba CAT-001 y CAT-004
prod, err := NewProduct("TENANT-1", "PROD-1", "SKU-001", "Paracetamol 500mg", "Laboratorio X", TypeMedicinal)
if err != nil {
t.Fatalf("Error inesperado creando producto: %v", err)
}

err = validator.ValidateProductCreation(prod)
if err != nil {
t.Errorf("CAT-001 falló: %v", err)
}

// Prueba presentación y subordinación GTIN (CAT-002, CAT-003)
pres, err := NewProductPresentation("TENANT-1", "PRES-1", "PROD-1", "Caja x 20", 20, true)
if err != nil {
t.Fatalf("Error creando presentación: %v", err)
}

err = validator.ValidatePresentation(pres, prod)
if err != nil {
t.Errorf("CAT-002/CAT-004 falló: %v", err)
}

// Prueba asociación de ingredientes con maestro global (CAT-005)
globalMaster := NewGlobalPlatformMaster()
globalMaster.ActiveIngredients["ING-01"] = ActiveIngredient{
IngredientID: "ING-01",
Name:         "Paracetamol",
}

strength, _ := NewStrength(500, "mg")
ingredient, _ := NewProductIngredient("TENANT-1", "PROD-1", "ING-01", strength, RoleActiveMain)

err = validator.ValidateIngredientAssociation(ingredient, globalMaster)
if err != nil {
t.Errorf("CAT-005 falló: %v", err)
}
}
