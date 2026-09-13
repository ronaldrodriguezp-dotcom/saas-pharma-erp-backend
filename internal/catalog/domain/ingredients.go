package domain

import (
"errors"
)

// Strength Representa la concentración de un principio activo con su respectiva unidad de medida
type Strength struct {
Value     float64
Unit      string // Ej: "mg", "g", "UI", "mL", "mcg"
Numerator float64
Denominator float64
}

func NewStrength(value float64, unit string) (Strength, error) {
if value < 0 {
return Strength{}, errors.New("el valor de la concentración no puede ser negativo")
}
if unit == "" {
return Strength{}, errors.New("la unidad de concentración es obligatoria")
}
return Strength{
Value:       value,
Unit:        unit,
Numerator:   value,
Denominator: 1.0,
}, nil
}

// ActiveIngredient Maestro Normalizado de Principios Activos (Plataforma / Global Master)
type ActiveIngredient struct {
IngredientID string
Name         string // Ej: "Paracetamol", "Ibuprofeno", "Ácido Acetilsalicílico"
CasNumber    string // Número CAS opcional para estandarización química
Description  string
}

func NewActiveIngredient(ingredientID, name, casNumber string) (ActiveIngredient, error) {
if ingredientID == "" || name == "" {
return ActiveIngredient{}, errors.New("el ID y el nombre del principio activo son obligatorios")
}
return ActiveIngredient{
IngredientID: ingredientID,
Name:         name,
CasNumber:    casNumber,
}, nil
}

// IngredientRole Define el rol del ingrediente en la formulación (Activo principal, coadyuvante, etc.)
type IngredientRole string

const (
RoleActiveMain   IngredientRole = "ACTIVE_MAIN"
RoleActiveCo     IngredientRole = "CO_ACTIVE"
RoleExcipient    IngredientRole = "EXCIPIENT"
)

// ProductIngredient Relación detallada entre el perfil del medicamento y los principios activos (CAT-005)
type ProductIngredient struct {
TenantID     string
ProductID    string
IngredientID string // Referencia al maestro ActiveIngredient
Strength     Strength
Role         IngredientRole
}

func NewProductIngredient(tenantID, productID, ingredientID string, strength Strength, role IngredientRole) (ProductIngredient, error) {
if tenantID == "" || productID == "" || ingredientID == "" {
return ProductIngredient{}, errors.New("parámetros obligatorios vacíos para asociar el ingrediente al producto")
}
return ProductIngredient{
TenantID:     tenantID,
ProductID:    productID,
IngredientID: ingredientID,
Strength:     strength,
Role:         role,
}, nil
}